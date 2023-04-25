package logger_transformers

import (
	"encoding/json"
	"fmt"
	logger_structs "service/services/logger/structs"
	"time"
)

type JsonTransformers struct{}

const unableToMarshalJsonTemplate = "{\"error\": \"%s\"}"

func (transformer *JsonTransformers) Transform(
	obj logger_structs.LogObject,
) string {
	obj.Timestamp = time.Now().UnixMilli()

	marshalled, err := json.Marshal(obj)

	if err != nil {
		return fmt.Sprintf(
			unableToMarshalJsonTemplate,
			err.Error(),
		)
	}

	return string(marshalled)
}

func (transformer *JsonTransformers) Debug() string {
	return "debug"
}

func (transformer *JsonTransformers) Info() string {
	return "info"
}

func (transformer *JsonTransformers) Warn() string {
	return "warn"
}

func (transformer *JsonTransformers) Error() string {
	return "error"
}
