package benthos

import (
	"context"
	"reflect"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

type (
	IBenthosRegistration interface {
		Register() error
	}
	IBenthosStream interface {
		Run(ctx context.Context) (err error)
		Stop(ctx context.Context) (err error)
	}
	UnimplementedIBenthosRegistration struct {
	}
	UnimplementedBenthosOutput struct {
	}
)

var TypeIBenthosRegistration = reflect.TypeOf((*IBenthosRegistration)(nil))

func (UnimplementedIBenthosRegistration) Register() error {
	return status.Error(codes.Unimplemented, "method Register not implemented")
}

func (UnimplementedBenthosOutput) Connect(context.Context) error {
	return nil
}

func (UnimplementedBenthosOutput) Write(context.Context, *benthos_service.Message) error {
	return status.Error(codes.Unimplemented, "method Connect not implemented")
}

func (UnimplementedBenthosOutput) Close(ctx context.Context) error {
	return nil
}
