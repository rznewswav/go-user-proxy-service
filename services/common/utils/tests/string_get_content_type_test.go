package tests

import (
	"service/services/common/utils"
	"testing"
)

func TestGetContentType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"A123_456", "article"},
		{"V123_456", "video"},
		{"P123_456", "podcast"},
		{"F_1_1", "feed"},
		{"ARandomString", "article"},
	}

	for _, test := range tests {
		actual := utils.StringGetContentType(test.input)
		if actual != test.expected {
			t.Errorf(
				"Expected %s for input %s, got %s",
				test.expected,
				test.input,
				actual,
			)
		}
	}
}
