package logger_transformers

import (
	"fmt"
	logger_structs "service/services/logger/structs"
	stack "service/services/stack"

	"gopkg.in/yaml.v2"
)

type PrettyTransformers struct{}

const unableToMarshalPrettyTemplate = "(unable to marshal object: %s)"

func (transformer *PrettyTransformers) Transform(
	obj logger_structs.LogObject,
) string {
	fileName := "<nil stack>"
	verboseStack := ""
	context := obj.Context
	message := ""
	level := obj.Level

	if len(obj.Stack) > 0 {
		stackObject := obj.Stack[0]
		fileName = fmt.Sprintf(
			"%s:%d",
			stackObject.File,
			stackObject.LineNumber,
		)
	}

	if len(obj.Stack) > 1 {
		verboseStack = stack.SimpleString(obj.Stack)
	}

	if obj.Message != nil {
		message = *obj.Message
	}

	marshalledPayload := ""

	if obj.Payload != nil {
		y, err := yaml.Marshal(obj.Payload)
		if err != nil {
			marshalledPayload = fmt.Sprintf(
				unableToMarshalPrettyTemplate,
				err.Error(),
			)
		} else {
			marshalledPayload = string(y)
		}
	}
	return fmt.Sprintf(
		"%s %s \u001b[0;33m(%s)\u001b[0m %s%s%s",
		level,
		fileName,
		context,
		message+"\n",
		marshalledPayload,
		verboseStack,
	)
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
