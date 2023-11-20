package bloblang

import (
	"reflect"

	"github.com/benthosdev/benthos/v4/public/bloblang"
)

type (
	IBlobLangFuncs interface {
		DeadLetterIt(args *bloblang.ParsedParams) (bloblang.Function, error)
	}
)

var TypeIBlobLangFuncs = reflect.TypeOf((*IBlobLangFuncs)(nil))
