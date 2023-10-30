package cloudeventoutput

import (
	"reflect"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
)

type (
	ICloudEventOutput interface {
		benthos_service.Output
	}
)

var TypeICloudEventOutput = reflect.TypeOf((*ICloudEventOutput)(nil))
