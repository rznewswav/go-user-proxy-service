package stack

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/mariomac/gostream/stream"
	"gopkg.in/yaml.v2"

	bnErrors "github.com/bugsnag/bugsnag-go/v2/errors"
)

var cwd = ""

func GetStackTrace() []bugsnag.StackFrame {
	err := bnErrors.New("", 1)
	sf := err.StackFrames()

	stackFrameStream := stream.OfSlice(sf)

	return stream.Map(
		stackFrameStream,
		func(sf bnErrors.StackFrame) bugsnag.StackFrame {
			if strings.Contains(sf.File, cwd) {
				if relativePath, err := filepath.Rel(cwd, sf.File); err == nil {
					sf.File = relativePath
				}
			}
			return bugsnag.StackFrame{
				Method:     sf.Func().Name(),
				File:       sf.File,
				LineNumber: sf.LineNumber,
				InProject: strings.Contains(
					sf.Package,
					"service",
				),
			}
		},
	).ToSlice()
}

func Init() {
	logger := log.New(os.Stderr, "STACK", log.LstdFlags)

	if wd, err := os.Getwd(); err != nil {
		logger.Panicf(
			"cannot init stack trace: %s", err,
		)
	} else {
		cwd = wd
	}

	for strings.Contains(cwd, "/cmd") {
		cwd = filepath.Dir(cwd)
	}
}

func SimpleString(stack []bugsnag.StackFrame) string {
	stackString := ""
	mappedStack := stream.Map(
		stream.OfSlice(stack),
		func(st bugsnag.StackFrame) string {
			return fmt.Sprintf(
				"%s:%d",
				st.File,
				st.LineNumber,
			)
		},
	).ToSlice()

	y, err := yaml.Marshal(mappedStack)
	if err != nil {
		stackString = fmt.Sprintf(
			"cannot marshal yaml: %s",
			err.Error(),
		)
	} else {
		stackString = string(y)
	}

	return stackString
}
