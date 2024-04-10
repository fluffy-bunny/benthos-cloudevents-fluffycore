package batcher

import (
	"context"
	"sync"

	centrifuge "github.com/centrifugal/centrifuge-go"
	contracts_centrifuge "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/contracts/centrifuge"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	status "github.com/gogo/status"
	async "github.com/reugn/async"
	zerolog "github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"

	queue "github.com/adrianbrad/queue"
	codes "google.golang.org/grpc/codes"
)

type (
	publicationBatcher struct {
		mutex        sync.Mutex
		Publications []centrifuge.Publication
		BatchSize    int32
	}
	batchManagement struct {
		publicationBatcher *publicationBatcher
		mutexBatch         sync.Mutex
		Batchs             queue.Queue[*contracts_centrifuge.Batch]
	}
	service struct {
		contracts_centrifuge.UnimplementedSubscriptionHandlers

		config              *contracts_centrifuge.CentrifugeConfig
		mutexEngine         sync.Mutex
		running             bool
		cancelCtx           context.CancelFunc
		ctxEngine           context.Context
		ctxCatchup          context.Context
		cancelCtxCatchup    context.CancelFunc
		future              async.Future[string]
		ctn                 di.Container
		client              contracts_centrifuge.ICentrifugeClient
		sub                 *centrifuge.Subscription
		err                 error
		log                 zerolog.Logger
		wgPublsh            sync.WaitGroup
		mutexHistoryCatchup sync.Mutex

		errHistoricalCatchUp error

		BatchManagement *batchManagement
	}
)

func NewBatchManagement(NumberOfBatches int32, batchSize int32) *batchManagement {
	return &batchManagement{
		publicationBatcher: NewPublicationBatcher(batchSize),
		Batchs:             queue.NewBlocking[*contracts_centrifuge.Batch](nil, queue.WithCapacity(int(NumberOfBatches))),
	}
}

var stemService = &service{}

func (s *service) Ctor(ctn di.Container) contracts_centrifuge.ICentrifugeStreamBatcher {

	return &service{
		ctn: ctn,
	}

}
func init() {
	var _ contracts_centrifuge.ICentrifugeStreamBatcher = (*service)(nil)
}
func AddTransientCentrifugeStreamBatcher(cb di.ContainerBuilder) {
	di.AddTransient[contracts_centrifuge.ICentrifugeStreamBatcher](cb, stemService.Ctor)
}
func (s *service) Configure(ctx context.Context, config *contracts_centrifuge.CentrifugeConfig) error {
	log := zerolog.Ctx(ctx).With().Logger()
	if s.running {
		err := status.Error(codes.AlreadyExists, "config cannot be changed while running")
		log.Error().Err(err).Msg("Configure")
		return err
	}
	s.config = config
	s.BatchManagement = NewBatchManagement(config.NumberOfBatches, config.BatchSize)
	return nil
}
func (s *service) GetBatch(ctx context.Context) (*contracts_centrifuge.Batch, error) {
	return nil, nil
}
func (s *service) Start(ctx context.Context) error {
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	s.mutexEngine.Lock()
	defer s.mutexEngine.Unlock()
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	if s.running {
		return status.Error(codes.AlreadyExists, "already running")
	}
	log := zerolog.Ctx(ctx).With().Logger()
	ctxEngine := context.Background()
	ctxEngine = log.WithContext(ctxEngine)
	s.ctxEngine, s.cancelCtx = context.WithCancel(ctxEngine)

	s.client = di.Get[contracts_centrifuge.ICentrifugeClient](s.ctn)
	s.client.Configure(s.ctxEngine, s.config.CentrifugeClientConfig)
	s.running = true
	s.future = s.asyncAction()
	return nil
}
func (s *service) Stop(ctx context.Context) error {
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	s.mutexEngine.Lock()
	defer s.mutexEngine.Unlock()
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	if !s.running {
		return nil
	}
	s.cancelCtx()
	s.running = false
	_, err := s.future.Join()
	s.client.Dispose(ctx)
	s.client = nil
	return err
}
func StringPtr(s string) *string {
	return &s
}
func (s *service) asyncAction() async.Future[string] {
	log := zerolog.Ctx(s.ctxEngine).With().Logger()
	s.log = log
	promise := async.NewPromise[string]()
	go func() {
		err := s.Subscribe(s.ctxEngine)
		if err != nil {
			log.Error().Err(err).Msg("failed to Subscribe")
			s.err = err
			promise.Failure(err)
			return
		}

		// wait for the s.ctxEngine has been canceled
		<-s.ctxEngine.Done()
		if s.sub != nil {
			s.sub.Unsubscribe()
		}
		s.sub = nil
		promise.Success(StringPtr("OK"))
	}()

	return promise.Future()
}

func (s *service) Subscribe(ctx context.Context) error {
	log := zerolog.Ctx(ctx).With().Logger()
	client, err := s.client.GetClient()
	if err != nil {
		log.Error().Err(err).Msg("failed to GetClient")
		return err
	}
	sub, err := client.NewSubscription(s.config.Channel,
		centrifuge.SubscriptionConfig{
			Recoverable: true,
			JoinLeave:   true,
			Positioned:  true,
		})
	if err != nil {
		log.Error().Err(err).Msg("failed to NewSubscription")
		return err
	}
	// register for all the events
	sub.OnError(s.OnSubscriptionErrorHandler)
	sub.OnJoin(s.OnJoinHandler)
	sub.OnLeave(s.OnLeaveHandler)
	sub.OnPublication(s.OnPublicationHandler)
	sub.OnSubscribed(s.OnSubscribedHandler)
	sub.OnSubscribing(s.OnSubscribingHandler)
	sub.OnUnsubscribed(s.OnUnsubscribedHandler)
	s.sub = sub
	err = sub.Subscribe()
	if err != nil {
		log.Error().Err(err).Msg("failed to Subscribe")
		return err
	}
	return nil
}
func (s *service) OnSubscribedHandler(e centrifuge.SubscribedEvent) {
	s.log.Debug().Msg("OnSubscribedHandler")
	// here we establish where centrifuge is in the stream
	// fire up a go routine to catch up with history
	s.goCatchupHistory(e)
}
func (s *service) goCatchupHistory(e centrifuge.SubscribedEvent) {
	if s.config.CatchupStreamPosition == nil {
		return
	}
	if s.config.CatchupStreamPosition.Offset == e.StreamPosition.Offset {
		return

	}
	// that a lock no matter what, we will release it when we have caught up with history.
	s.wgPublsh.Add(1)
	if s.cancelCtxCatchup != nil {
		s.cancelCtxCatchup()
	}
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	// we only allow one go routine to catch up with history at a time
	s.mutexHistoryCatchup.Lock()
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	s.ctxCatchup, s.cancelCtxCatchup = context.WithCancel(context.Background())
	subscribedStreamPosition := e.StreamPosition
	type StreamPositionTracker struct {
		CurrentHistoricalStreamPosition *centrifuge.StreamPosition
		CaughtUpStreamPosition          *centrifuge.StreamPosition
	}
	s.config.CatchupStreamPosition.Epoch = subscribedStreamPosition.Epoch
	tracker := &StreamPositionTracker{
		CurrentHistoricalStreamPosition: s.config.CatchupStreamPosition,
		CaughtUpStreamPosition:          subscribedStreamPosition,
	}
	go func() {
		defer func() {
			// release our wait group reference
			s.wgPublsh.Done()
			// unlock the history catch up mutex
			s.mutexHistoryCatchup.Unlock()
		}()
		log := log.With().Caller().Interface("subscribedStreamPosition", subscribedStreamPosition).Logger()
		currentHistoricalStreamPosition := tracker.CurrentHistoricalStreamPosition
		log.Debug().Interface("currentHistoricalStreamPosition", currentHistoricalStreamPosition).Msg("currentHistoricalStreamPosition - initial")

		for {
			if IsCanceled(s.ctxCatchup) {
				log.Debug().Msg("IsCanceled")
				return
			}
			if tracker.CurrentHistoricalStreamPosition.Offset == tracker.CaughtUpStreamPosition.Offset {
				log.Debug().Msg("caught up")
				break
			}
			historyResult, err := s.sub.History(context.Background(),
				centrifuge.WithHistorySince(currentHistoricalStreamPosition),
				centrifuge.WithHistoryLimit(int32(s.config.BatchSize)),
			)
			if err != nil {
				s.errHistoricalCatchUp = err
				log.Error().Err(err).Msg("failed to History")
				break
			}
			if len(historyResult.Publications) == 0 {
				break
			}
			for _, publication := range historyResult.Publications {
				s.BatchManagement.Add(publication)
			}

		}
	}()

}
func IsCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		// The context has been canceled
		return true
	default:
		// The context has not been canceled
		return false
	}
}

func (b *batchManagement) Offer(batch *contracts_centrifuge.Batch) {
	b.mutexBatch.Lock()
	defer b.mutexBatch.Unlock()
	b.Batchs.Offer(batch)
}

func NewPublicationBatcher(batchSize int32) *publicationBatcher {
	return &publicationBatcher{
		BatchSize: batchSize,
	}
}

// Add adds a publication to the batcher. Returns true if the batch is full and ready to be sent.
func (p *publicationBatcher) Add(publication centrifuge.Publication) bool {
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	p.mutex.Lock()
	defer p.mutex.Unlock()
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	p.Publications = append(p.Publications, publication)
	return int32(len(p.Publications)) >= p.BatchSize
}

func (p *publicationBatcher) PullBatch() *contracts_centrifuge.Batch {
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	p.mutex.Lock()
	defer p.mutex.Unlock()
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	batch := &contracts_centrifuge.Batch{
		Publications: p.Publications,
	}
	p.Publications = nil
	return batch
}
func (b *batchManagement) Add(publication centrifuge.Publication) {
	ready := b.publicationBatcher.Add(publication)
	if ready {
		// pull it from the batcher and add it to our batch queue
		batch := b.publicationBatcher.PullBatch()
		b.Batchs.Offer(batch)
	}
}
func (b *batchManagement) PullBatch(publication centrifuge.Publication) *contracts_centrifuge.Batch {
	batch, err := b.Batchs.Get()
	if err != nil {
		if err == queue.ErrNoElementsAvailable {
			return nil
		}
	}
	return batch
}
