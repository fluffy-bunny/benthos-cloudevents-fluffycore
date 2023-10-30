package benthos

import "reflect"

type (
	IBenthosRegistration interface {
		Register() error
	}
)

var TypeIBenthosRegistration = reflect.TypeOf((*IBenthosRegistration)(nil))
