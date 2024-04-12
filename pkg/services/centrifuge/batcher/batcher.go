package batcher

import (
	"context"
	"sync"

	queue "github.com/adrianbrad/queue"
	centrifuge "github.com/centrifugal/centrifuge-go"
	contracts_centrifuge "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/contracts/centrifuge"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	status "github.com/gogo/status"
	async "github.com/reugn/async"
	zerolog "github.com/rs/zerolog"
	blockingQueues "github.com/theodesp/blockingQueues"
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

		Batchs *blockingQueues.BlockingQueue
	}
	service struct {
		contracts_centrifuge.UnimplementedSubscriptionHandlers

		config *contracts_centrifuge.CentrifugeConfig
		// mutexEngine guards the configure, stop, and start methods
		mutexEngine      sync.Mutex
		running          bool
		cancelCtx        context.CancelFunc
		ctxEngine        context.Context
		ctxCatchup       context.Context
		cancelCtxCatchup context.CancelFunc
		future           async.Future[string]
		// container to our di so we can dynamically get services
		ctn di.Container
		// client to connect to centrifuge
		client contracts_centrifuge.ICentrifugeClient
		// sub is the subscription to the channel
		sub *centrifuge.Subscription
		err error
		// log is the logger
		log zerolog.Logger
		// wgPublsh is a wait group to ensure we have caught up with history before we start processing current publications
		wgPublsh sync.WaitGroup
		// mutexHistoryCatchup is a mutex to ensure only one go routine is catching up with history at a time
		mutexHistoryCatchup sync.Mutex

		errHistoricalCatchUp error

		// BatchManagement is a struct that accepts publications and batches them up
		BatchManagement *batchManagement
	}
)

func NewBatchManagement(NumberOfBatches int, batchSize int32) *batchManagement {
	batch, err := blockingQueues.NewArrayBlockingQueue(uint64(NumberOfBatches))
	if err != nil {
		panic(err)
	}

	return &batchManagement{
		publicationBatcher: NewPublicationBatcher(batchSize),
		Batchs:             batch,
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
func (s *service) IsRunning() bool {
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	s.mutexEngine.Lock()
	defer s.mutexEngine.Unlock()
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	return s.running
}
func (s *service) BatchState() contracts_centrifuge.BatchState {
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	s.mutexEngine.Lock()
	defer s.mutexEngine.Unlock()
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	return s.BatchManagement.BatchState()
}

func (s *service) Configure(ctx context.Context, config *contracts_centrifuge.CentrifugeConfig) error {
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	s.mutexEngine.Lock()
	defer s.mutexEngine.Unlock()
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--

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
func (s *service) GetBatch(ctx context.Context, flush bool) (*contracts_centrifuge.Batch, error) {
	batch := s.BatchManagement.PullBatch(flush)
	return batch, nil
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
	// stop any publishing from comming in.
	s.wgPublsh.Add(1)
	defer s.wgPublsh.Done()
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
	s.goCatchupHistory(e.StreamPosition)
}
func (s *service) OnPublicationHandler(e centrifuge.PublicationEvent) {
	s.log.Debug().Msg("OnPublicationHandler")
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	// we are blocking publication to use from centrifuge until we have caught up with history
	s.wgPublsh.Wait()
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	// we may still be behind here.  If the stream position of this message is more than one greater than ours then we must catchup to that.

	if e.Publication.Offset > s.config.HistoricalStreamPosition.Offset+1 {
		s.goCatchupHistory(&centrifuge.StreamPosition{
			Epoch:  s.config.HistoricalStreamPosition.Epoch,
			Offset: e.Publication.Offset,
		})
		return
	}
	s.internalOnPublicationHandler(&e)
}
func (s *service) internalOnPublicationHandler(e *centrifuge.PublicationEvent) {
	s.BatchManagement.Add(e.Publication)
	s.config.HistoricalStreamPosition.Offset = e.Publication.Offset
}

func (s *service) goCatchupHistory(streamPosition *centrifuge.StreamPosition) {
	log := s.log.With().Interface("catchup_stream_position", streamPosition).Logger()
	log.Info().Msg("goCatchupHistory")
	if s.config.HistoricalStreamPosition == nil {
		// our first time through
		s.config.HistoricalStreamPosition = streamPosition
		log.Info().Msg("no historical data requested")
		return
	}
	log = log.With().Interface("historical_stream_position", s.config.HistoricalStreamPosition).Logger()
	if s.config.HistoricalStreamPosition.Offset == streamPosition.Offset {
		s.log.Info().Msg("already caught up")
		return
	}
	if len(s.config.HistoricalStreamPosition.Epoch) == 0 {
		// bad config
		s.log.Error().Msg("bad config - HistoricalStreamPosition.Epoch is empty")
		return
	}
	if s.config.HistoricalStreamPosition.Epoch != streamPosition.Epoch {
		// this will fail with the following error
		// {"error":"112: unrecoverable position"}
		s.config.HistoricalStreamPosition.Epoch = streamPosition.Epoch
		s.config.HistoricalStreamPosition.Offset = 0
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
	subscribedStreamPosition := streamPosition
	type StreamPositionTracker struct {
		SubscribedStreamPosition *centrifuge.StreamPosition
	}
	tracker := &StreamPositionTracker{
		SubscribedStreamPosition: subscribedStreamPosition,
	}
	go func() {
		defer func() {
			// release our wait group reference
			s.wgPublsh.Done()
			// unlock the history catch up mutex
			s.mutexHistoryCatchup.Unlock()
		}()
		log := log.With().Caller().Interface("subscribedStreamPosition", subscribedStreamPosition).Logger()
		currentHistoricalStreamPosition := s.config.HistoricalStreamPosition
		log.Debug().Interface("currentHistoricalStreamPosition", currentHistoricalStreamPosition).Msg("currentHistoricalStreamPosition - initial")

		for {
			if IsCanceled(s.ctxCatchup) {
				log.Debug().Msg("IsCanceled")
				return
			}
			if s.config.HistoricalStreamPosition.Offset == tracker.SubscribedStreamPosition.Offset {
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
				// normalize our historical pulls as if we got the event directly from centrifuge
				s.internalOnPublicationHandler(&centrifuge.PublicationEvent{
					Publication: publication,
				})
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

func NewPublicationBatcher(batchSize int32) *publicationBatcher {
	return &publicationBatcher{
		BatchSize: batchSize,
	}
}
func (p *publicationBatcher) IsDirty() bool {
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	p.mutex.Lock()
	defer p.mutex.Unlock()
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	return len(p.Publications) > 0
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

func (b *batchManagement) BatchState() contracts_centrifuge.BatchState {
	if b.Batchs.IsEmpty() {
		if b.publicationBatcher.IsDirty() {
			return contracts_centrifuge.BatchState_PartiallyAvailable
		}
		return contracts_centrifuge.BatchState_Empty
	}
	return contracts_centrifuge.BatchState_AtLeastOneAvailable
}

func (b *batchManagement) Add(publication centrifuge.Publication) {
	ready := b.publicationBatcher.Add(publication)
	if ready {
		// pull it from the batcher and add it to our batch queue
		batch := b.publicationBatcher.PullBatch()
		b.Batchs.Put(batch) // PUT is the blocking call
	}
}
func (b *batchManagement) PullBatch(flush bool) *contracts_centrifuge.Batch {
	isEmpty := b.Batchs.IsEmpty()
	if isEmpty {
		if flush {
			// force the pull from the internal batcher, doesn't have to be full.
			batch := b.publicationBatcher.PullBatch()
			if len(batch.Publications) > 0 {
				return batch
			}
		}
		return nil
	}

	batch, err := b.Batchs.Get() // Get is the blocking call
	if err != nil {
		if err == queue.ErrNoElementsAvailable {
			return nil
		}
	}
	return batch.(*contracts_centrifuge.Batch)
}
