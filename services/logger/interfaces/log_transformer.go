package logger_interfaces

import logger_structs "service/services/logger/structs"

type ILogTransformer interface {
	Transform(obj logger_structs.LogObject) string
	Debug() string
	Info() string
	Warn() string
	Error() string
}
