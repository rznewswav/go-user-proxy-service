package mysql

import (
	"context"
	"service/services/bugsnag"
	"time"
)

type WeeklyCovidTotalStats struct {
	// hazard: with ?parseTime=true enabled in url string,
	// db column of type DATE will be parsed as UTC
	//
	// eg. db record 2023-04-02
	// sql package will parse as 2023-04-02 00:00:00 +0000 UTC
	StartWeek time.Time
	EndWeek   time.Time
	Confirmed int
	Death     int
	Cured     int
}

func GetWeeklyCovidTotalCases() (WeeklyCovidTotalStats, error) {
	queryCtx, queryCtxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer queryCtxCancel()
	row, queryError := Client.QueryContext(queryCtx, string(SQLGetWeeklyCovidTotalCases))
	if queryError != nil {
		return WeeklyCovidTotalStats{}, bugsnag.FromError("SQL Latest Covid Recovery Rate Error", queryError)
	}
	if row == nil {
		return WeeklyCovidTotalStats{}, bugsnag.FromError("SQL Latest Covid Recovery Rate Row Nil Error", ErrRowNil)
	}

	defer row.Close()

	if !row.Next() {
		return WeeklyCovidTotalStats{}, nil
	}

	var totalStats WeeklyCovidTotalStats
	row.Scan(
		&totalStats.StartWeek,
		&totalStats.EndWeek,
		&totalStats.Confirmed,
		&totalStats.Death,
		&totalStats.Cured,
	)
	return totalStats, nil
}
