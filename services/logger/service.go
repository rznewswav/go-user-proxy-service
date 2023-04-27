package logger

import (
	"fmt"
	"service/services/bugsnag"
	"service/services/common/utils"
	logger_structs "service/services/logger/structs"
	"service/services/stack"

	bnDriver "github.com/bugsnag/bugsnag-go/v2"
)

type WithContext struct {
	Context string
}

var logEvents = make(map[string]bool)

func LogEvents(contexts []string) {
	logger := For("logger")
	logger.Info(
		"Logging debug events from: %s",
		contexts,
	)
	for _, event := range contexts {
		logEvents[event] = true
	}
}

func For(context string) *WithContext {
	logger := WithContext{
		Context: context,
	}

	return &logger
}

func getLogObject(message string, formatOrPayload []any, level string, context string) logger_structs.LogObject {
	format, payload := SplitPayloadIntoFormatterAndNotifyableError(
		formatOrPayload,
	)
	formattedMessage := fmt.Sprintf(message, format...)
	if payload != nil && len(payload.Message) == 0 {
		payload.Message = formattedMessage
	}

	var stackTrace []bnDriver.StackFrame
	if payload != nil && len(payload.Stacks) > 0 {
		stackTrace = payload.Stacks
	} else {
		stackTrace = stack.GetStackTrace()
	}

	offset := 3
	stackTrace = utils.ArraySlice(
		stackTrace,
		offset,
		offset+TraceLogDepth,
	)

	logObject := logger_structs.LogObject{
		Level:   level,
		Context: context,
		Message: &formattedMessage,
		Stack:   stack.ToStackString(stackTrace),
		Payload: payload,
	}
	return logObject
}

// Debug Logs DEBUG to console
//
// Usage: Pass string formatting in the payload.
// E.g.
//
//	logger.Debug(
//		"Initialized service: %s",
//		serviceName
//	)
//
// Passing bugsnag.NotifyableError at the last parameter will
// result in logging it as YAML object
//
//	logger.Debug(
//		"Initialized service: %s",
//		serviceName,
//		bugsnag.NotifyableError{ Error: "Test" }
//	)
//
// Output:
//
//	Initialized service: <service name>
//	- error: Test
func (logger *WithContext) Debug(
	message string,
	formatOrPayload ...any,
) {
	level := transformer.Debug()
	useLogger := globalDebugLogger
	if logEvents[logger.Context] {
		level = transformer.Info()
		useLogger = globalInfoLogger
	}

	logObject := getLogObject(message, formatOrPayload, level, logger.Context)
	logString := transformer.Transform(logObject)
	useLogger.Print(logString)
}

// Info Logs INFO to console
//
// Usage: Pass string formatting in the payload.
// E.g.
//
//	logger.Debug(
//		"Initialized service: %s",
//		serviceName
//	)
//
// Passing bugsnag.NotifyableError at the last parameter will
// result in logging it as YAML object
//
//	logger.Debug(
//		"Initialized service: %s",
//		serviceName,
//		bugsnag.NotifyableError{ Error: "Test" }
//	)
//
// Output:
//
//	Initialized service: <service name>
//	- error: Test
func (logger *WithContext) Info(
	message string,
	formatOrPayload ...any,
) {
	logObject := getLogObject(message, formatOrPayload, transformer.Info(), logger.Context)

	logString := transformer.Transform(logObject)
	globalInfoLogger.Print(logString)
}

// Warn Logs WARN to console
//
// Usage: Pass string formatting in the payload.
// E.g.
//
//	logger.Debug(
//		"Initialized service: %s",
//		serviceName
//	)
//
// Passing bugsnag.NotifyableError at the last parameter will
// result in logging it as YAML object
//
//	logger.Debug(
//		"Initialized service: %s",
//		serviceName,
//		bugsnag.NotifyableError{ Error: "Test" }
//	)
//
// Output:
//
//	Initialized service: <service name>
//	- error: Test
func (logger *WithContext) Warn(
	message string,
	formatOrPayload ...any,
) {
	logObject := getLogObject(message, formatOrPayload, transformer.Warn(), logger.Context)
	logString := transformer.Transform(logObject)
	globalWarnLogger.Print(logString)
}

// Error Logs ERROR to console
//
// Usage: Pass string formatting in the payload.
// E.g.
//
//	logger.Debug(
//		"Initialized service: %s",
//		serviceName
//	)
//
// Passing bugsnag.NotifyableError at the last parameter will
// result in logging it as YAML object.
//
//	Logger.Debug(
//		"Initialized service: %s",
//		serviceName,
//		bugsnag.NotifyableError{ Error: "Test" }
//	)
//
// Output:
//
//	Initialized service: <service name>
//	- error: Test
//
// This object will also be reported to the bugsnag on
// enabled release stages.
func (logger *WithContext) Error(
	message string,
	formatOrPayload ...any,
) {
	logObject := getLogObject(message, formatOrPayload, transformer.Info(), logger.Context)
	logString := transformer.Transform(logObject)
	globalErrorLogger.Print(logString)

	ne := logObject.Payload
	bugsnag.GetHandler().Notify(ne)
}
