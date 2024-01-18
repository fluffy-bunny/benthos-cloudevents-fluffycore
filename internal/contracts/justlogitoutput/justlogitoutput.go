package justlogitoutput

import (
	"reflect"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
)

type (
	IJustLogItOutput interface {
		benthos_service.Output
	}
)

var TypeIJustLogItOutput = reflect.TypeOf((*IJustLogItOutput)(nil))
