package utils

import "time"

func TimeRandomRangeHour(
	start time.Time,
	end time.Time,
) time.Time {
	timeDifference := start.Sub(end)
	hours := timeDifference.Abs().Hours()
	return start.Add(
		time.Duration(hours*Random.Float64()) * time.Hour,
	)
}
