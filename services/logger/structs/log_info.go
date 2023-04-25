package logger_structs

import (
	bugsnag_structs "service/services/bugsnag/structs"

	"github.com/bugsnag/bugsnag-go/v2"
)

type LogObject struct {
	Context string `json:"context"`
	/* `Message` field should not be populated with `Error`. Only either one of them should be populated. */
	Message   *string                          `json:"message"`
	Level     string                           `json:"level"`
	Stack     []bugsnag.StackFrame             `json:"stack"`
	Payload   *bugsnag_structs.NotifyableError `json:"payload"`
	Timestamp int64                            `json:"timestamp"`
}
