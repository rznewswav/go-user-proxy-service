package main

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime"
)

type WithStackTrace interface {
	StackTrace() errors.StackTrace
}

func fileLine(f errors.Frame) string {
	fn := runtime.FuncForPC(uintptr(f) - 1)
	if fn == nil {
		return "unknown:0"
	}
	file, line := fn.FileLine(uintptr(f) - 1)
	return fmt.Sprintf("%s:%d", file, line)
}

func main() {
	err := errors.WithStack(errors.New("hello"))
	if errorWithStack, castable := err.(WithStackTrace); castable {
		stackTraces := errorWithStack.StackTrace()
		for _, trace := range stackTraces {
			println(fileLine(trace))
		}
	}
}
