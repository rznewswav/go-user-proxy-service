package logger

import (
	"fmt"
	"service/services/bugsnag"
	"service/services/common/utils"
	logger_structs "service/services/logger/structs"
	"service/services/stack"

	bnDriver "github.com/bugsnag/bugsnag-go/v2"
)

type loggerWContext logger_structs.LoggerWithContext

var logEvents = make(map[string]bool)

func LogEvents(contexts []string) {
	logger := WithContext("logger")
	logger.Info(
		"Logging debug events from: %s",
		contexts,
	)
	for _, event := range contexts {
		logEvents[event] = true
	}
}

func WithContext(context string) *loggerWContext {
	logger := loggerWContext{
		Context: context,
	}

	return &logger
}

/*
Logs DEBUG to console

Usage: Pass string formatting in the payload. eg.

	logger.Debug(
		"Initialized service: %s",
		serviceName
	)

Passing bugsnag.NotifyableError at the last parameter will
result in logging it as YAML object

	logger.Debug(
		"Initialized service: %s",
		serviceName,
		bugsnag.NotifyableError{ Error: "Test" }
	)

Output:

	Initialized service: <service name>
	- error: Test
*/
func (logger *loggerWContext) Debug(
	message string,
	formatOrPayload ...any,
) {
	level := transformer.Debug()
	useLogger := globalDebugLogger
	if logEvents[logger.Context] {
		level = transformer.Info()
		useLogger = globalInfoLogger
	}

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

	logObject := logger_structs.LogObject{
		Level:   level,
		Context: logger.Context,
		Message: &formattedMessage,
		Stack: utils.ArraySlice(
			stackTrace,
			2,
			2+TraceLogDepth,
		),
		Payload: payload,
	}

	logString := transformer.Transform(logObject)
	useLogger.Print(logString)
}

/*
Logs INFO to console

Usage: Pass string formatting in the payload. eg.

	logger.Debug(
		"Initialized service: %s",
		serviceName
	)

Passing bugsnag.NotifyableError at the last parameter will
result in logging it as YAML object

	logger.Debug(
		"Initialized service: %s",
		serviceName,
		bugsnag.NotifyableError{ Error: "Test" }
	)

Output:

	Initialized service: <service name>
	- error: Test
*/
func (logger *loggerWContext) Info(
	message string,
	formatOrPayload ...any,
) {
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

	logObject := logger_structs.LogObject{
		Level:   transformer.Info(),
		Context: logger.Context,
		Message: &formattedMessage,
		Stack: utils.ArraySlice(
			stackTrace,
			2,
			2+TraceLogDepth,
		),
		Payload: payload,
	}

	logString := transformer.Transform(logObject)
	globalInfoLogger.Print(logString)
}

/*
Logs WARN to console

Usage: Pass string formatting in the payload. eg.

	logger.Debug(
		"Initialized service: %s",
		serviceName
	)

Passing bugsnag.NotifyableError at the last parameter will
result in logging it as YAML object

	logger.Debug(
		"Initialized service: %s",
		serviceName,
		bugsnag.NotifyableError{ Error: "Test" }
	)

Output:

	Initialized service: <service name>
	- error: Test
*/
func (logger *loggerWContext) Warn(
	message string,
	formatOrPayload ...any,
) {
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

	logObject := logger_structs.LogObject{
		Level:   transformer.Warn(),
		Context: logger.Context,
		Message: &formattedMessage,
		Stack:   stackTrace,
		Payload: payload,
	}

	logString := transformer.Transform(logObject)
	globalWarnLogger.Print(logString)
}

/*
Logs ERROR to console

Usage: Pass string formatting in the payload. eg.

	logger.Debug(
		"Initialized service: %s",
		serviceName
	)

Passing bugsnag.NotifyableError at the last parameter will
result in logging it as YAML object.

	logger.Debug(
		"Initialized service: %s",
		serviceName,
		bugsnag.NotifyableError{ Error: "Test" }
	)

Output:

	Initialized service: <service name>
	- error: Test

This object will also be reported to the bugsnag on
enabled release stages.
*/
func (logger *loggerWContext) Error(
	message string,
	formatOrPayload ...any,
) {
	format, payload := SplitPayloadIntoFormatterAndNotifyableError(
		formatOrPayload,
	)
	formattedMessage := fmt.Sprintf(message, format...)

	if payload == nil {
		pl := bugsnag.New("Unknown Error")
		payload = pl
	}

	payload.Message = utils.If(
		len(payload.Message) == 0,
		formattedMessage,
		payload.Message,
	)

	var stackTrace []bnDriver.StackFrame

	if len(payload.Stacks) > 0 {
		stackTrace = payload.Stacks
	} else {
		stackTrace = stack.GetStackTrace()
	}

	logObject := logger_structs.LogObject{
		Level:   transformer.Error(),
		Context: logger.Context,
		Message: &formattedMessage,
		Stack: utils.ArraySlice(
			stackTrace,
			2,
			2+TraceLogDepth,
		),
		Payload: payload,
	}

	logString := transformer.Transform(logObject)
	globalErrorLogger.Print(logString)

	ne := logObject.Payload
	bugsnag.GetHandler().Notify(ne)
}
