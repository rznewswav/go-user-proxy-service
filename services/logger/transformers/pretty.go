package logger_transformers

import (
	"encoding/json"
	"fmt"
	logger_structs "service/services/logger/structs"
	"strings"
)

type PrettyTransformers struct{}

const unableToMarshalPrettyTemplate = "(unable to marshal object: %s)"

var maxFileNameLength = 31
var maxContextLength = 8

func (transformer *PrettyTransformers) Transform(
	obj logger_structs.LogObject,
) string {
	fileName := "<nil stack>"
	verboseStack := ""
	context := obj.Context
	message := ""
	level := obj.Level

	if len(obj.Stack) > 0 {
		fileName = obj.Stack[0]
	}

	if len(obj.Stack) > 1 {
		verboseStack = strings.Join(obj.Stack, "\n")
	}

	if obj.Message != nil {
		message = *obj.Message
	}

	marshalledPayload := ""

	if obj.Payload != nil {
		y, err := json.MarshalIndent(obj.Payload, "", "  ")
		if err != nil {
			marshalledPayload = fmt.Sprintf(
				unableToMarshalPrettyTemplate,
				err.Error(),
			)
		} else {
			marshalledPayload = string(y)
		}
	}

	maxFileNameLength = maxOfInt(maxFileNameLength, len(fileName))
	maxContextLength = maxOfInt(maxContextLength, len(context))
	return fmt.Sprintf(
		"%s %*s \u001b[0;33m%*s\u001b[0m %s%s%s",
		level,
		maxFileNameLength,
		fileName,
		maxContextLength,
		context,
		message+"\n",
		marshalledPayload,
		verboseStack,
	)
}

func maxOfInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func (transformer *PrettyTransformers) Debug() string {
	return "\u001b[0;35mDEBUG\u001b[0m"
}

func (transformer *PrettyTransformers) Info() string {
	return "\u001b[1;36m INFO\u001b[0m"
}

func (transformer *PrettyTransformers) Warn() string {
	return "\u001b[1;33m WARN\u001b[0m"
}

func (transformer *PrettyTransformers) Error() string {
	return "\u001b[1;31mERROR\u001b[0m"
}
