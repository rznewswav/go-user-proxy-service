package logger_test

import (
	bugsnag_structs "service/services/bugsnag"
	"service/services/logger"
	"testing"
)

func TestSeparateNotifyableErrorArrayLength0(t *testing.T) {
	get := func(payload ...any) ([]any, *bugsnag_structs.NotifiableError) {
		return logger.SplitPayloadIntoFormatterAndNotifyableError(
			payload,
		)
	}

	inputNotifyableError := bugsnag_structs.NotifiableError{
		Message: "test",
	}

	outPayload, outNotifyableError := get(
		inputNotifyableError,
	)

	if len(outPayload) != 0 {
		t.Fatalf(
			"expected outpayload to be length 0",
		)
		return
	}

	if outNotifyableError.Message != inputNotifyableError.Message {
		t.Fatalf(
			"expected outNotifyableError to be the same as inputNotifyableError",
		)
		return
	}
}

func TestSeparateNotifyableErrorArrayLength0NoError(
	t *testing.T,
) {
	get := func(payload ...any) ([]any, *bugsnag_structs.NotifiableError) {
		return logger.SplitPayloadIntoFormatterAndNotifyableError(
			payload,
		)
	}

	outPayload, outNotifyableError := get()

	if len(outPayload) != 0 {
		t.Fatalf(
			"expected outpayload to be length 0",
		)
		return
	}

	if outNotifyableError != nil {
		t.Fatalf(
			"expected outNotifyableError to be nil",
		)
		return
	}
}
func TestSeparateNotifyableErrorArrayLength1(t *testing.T) {
	get := func(payload ...any) ([]any, *bugsnag_structs.NotifiableError) {
		return logger.SplitPayloadIntoFormatterAndNotifyableError(
			payload,
		)
	}

	outPayload, outNotifyableError := get("a")

	if len(outPayload) != 1 || outPayload[0] != "a" {
		t.Fatalf(
			"expected outpayload to be length 1 and with item only \"a\"",
		)
		return
	}

	if outNotifyableError != nil {
		t.Fatalf(
			"expected outNotifyableError to be nil",
		)
		return
	}
}

func TestSeparateNotifyableErrorArrayLength2(t *testing.T) {
	get := func(payload ...any) ([]any, *bugsnag_structs.NotifiableError) {
		return logger.SplitPayloadIntoFormatterAndNotifyableError(
			payload,
		)
	}

	inputNotifyableError := bugsnag_structs.NotifiableError{
		Message: "test",
	}

	outPayload, outNotifyableError := get(
		"a",
		"b",
		inputNotifyableError,
	)

	if len(outPayload) != 2 || outPayload[0] != "a" ||
		outPayload[1] != "b" {
		t.Fatalf(
			"expected outpayload to be length 2 and with item only [\"a\", \"b\"]",
		)
		return
	}

	if outNotifyableError.Message != inputNotifyableError.Message {
		t.Fatalf(
			"expected outNotifyableError to be the same as inputNotifyableError",
		)
		return
	}
}
