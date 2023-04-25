package tests

import (
	"service/services/common/utils"
	"testing"
)

func TestISOStringParsing(t *testing.T) {
	var expectedMilli int64 = 1673422708869

	inputDateISO := "2023-01-11T07:38:28.869Z"

	dateTime, err := utils.TimeFromJSISOString(inputDateISO)

	if err != nil {
		t.Error(err)
		return
	}

	if expectedMilli != dateTime.UnixMilli() {
		t.Errorf(
			"expected milli: %d, got milli: %d",
			expectedMilli,
			dateTime.UnixMilli(),
		)
	}

}
