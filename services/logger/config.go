package logger

import (
	"io"
	"log"
	"os"

	logger_interfaces "service/services/logger/interfaces"
	logger_transformers "service/services/logger/transformers"
)

var nilLogger = log.New(io.Discard, "nil", 0)
var globalDebugLogger *log.Logger = nilLogger
var globalInfoLogger *log.Logger = nilLogger
var globalWarnLogger *log.Logger = nilLogger
var globalErrorLogger *log.Logger = nilLogger
var debugOut = io.Discard
var infoOut = io.Discard
var warnOut = io.Discard
var errorOut = io.Discard
var TraceLogDepth int = 1

const standardFlags = log.LstdFlags
const noFlag = 0

func enableErrorLog() {
	errorOut = os.Stderr
}

func enableWarnLog() {
	warnOut = os.Stderr
	enableErrorLog()
}

func enableInfoLog() {
	infoOut = os.Stdout
	enableWarnLog()
}

func enableDebugLog() {
	debugOut = os.Stdout
	enableInfoLog()
}

func SetTraceLogDepth(depth int) {
	if depth < 1 {
		depth = 1
	}
	TraceLogDepth = depth
}

func SetLogLevel(level string) {

	switch level {
	case "error":
		enableErrorLog()
	case "warn":
		enableWarnLog()
	case "info":
		enableInfoLog()
	case "debug":
		enableDebugLog()
	}

	globalDebugLogger = log.New(debugOut, "", standardFlags)
	globalInfoLogger = log.New(infoOut, "", standardFlags)
	globalWarnLogger = log.New(warnOut, "", standardFlags)
	globalErrorLogger = log.New(errorOut, "", standardFlags)
}

func UseStandardFlag() {
	globalDebugLogger.SetFlags(standardFlags)
	globalInfoLogger.SetFlags(standardFlags)
	globalWarnLogger.SetFlags(standardFlags)
	globalErrorLogger.SetFlags(standardFlags)
}

func UseNoFlag() {
	globalDebugLogger.SetFlags(noFlag)
	globalInfoLogger.SetFlags(noFlag)
	globalWarnLogger.SetFlags(noFlag)
	globalErrorLogger.SetFlags(noFlag)
}

var jsonTransformer = logger_transformers.JsonTransformers{}
var prettyTransformer = logger_transformers.PrettyTransformers{}

var transformer logger_interfaces.ILogTransformer = &prettyTransformer

func UseJsonTransformer() {
	transformer = &jsonTransformer
	UseNoFlag()
}

func UsePrettyTransformer() {
	transformer = &prettyTransformer
	UseStandardFlag()
}
