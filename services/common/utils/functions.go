package utils

import (
	"fmt"
	"runtime"
	"strings"
)

func GetCurrentFuncInfo() (uintptr, string, int, bool) {
	pc, file, line, ok := runtime.Caller(1)
	return pc, file, line, ok
}

func GetCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	// split string using delimiter
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	// access last item in array
	return fmt.Sprintf("%s", parts[len(parts)-1])
}
