package utils

import (
	"fmt"
	"runtime"
)

func Caller() string {
	// skip this frame
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v:%v", file, line)
}
