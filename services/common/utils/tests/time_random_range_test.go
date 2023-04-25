package tests

import (
	"service/services/common/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeRandomRange(t *testing.T) {
	endTime := time.Now()
	startTime := time.Now().Add(-48 * time.Hour)

	firstCandidate := utils.TimeRandomRangeHour(
		startTime,
		endTime,
	)
	secondCandidate := utils.TimeRandomRangeHour(
		startTime,
		endTime,
	)

	assert.NotEqual(t, firstCandidate, secondCandidate)
}
