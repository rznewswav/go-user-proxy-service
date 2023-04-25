package utils

import "time"

func TimeBeforeDays(numOfDays time.Duration) time.Time {
	return time.Now().
		Add(-numOfDays * time.Hour * 24)
}
