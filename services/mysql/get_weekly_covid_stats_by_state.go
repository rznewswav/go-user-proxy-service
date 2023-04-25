package mysql

import (
	"context"
	"service/services/bugsnag"
	"time"
)

type WeeklyCovidStatsByState struct {
	State string
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

func GetWeeklyCovidStatsByState() ([]WeeklyCovidStatsByState, error) {
	queryCtx, queryCtxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer queryCtxCancel()
	row, queryError := Client.QueryContext(queryCtx, string(SQLGetCovidLatestWeekStatsByState))
	if queryError != nil {
		return []WeeklyCovidStatsByState{}, bugsnag.FromError("SQL Weekly Covid Stats By State Error", queryError)
	}
	if row == nil {
		return []WeeklyCovidStatsByState{}, bugsnag.FromError("SQL Weekly Covid Stats By State Row Nil Error", ErrRowNil)
	}

	defer row.Close()

	scannedRows := make([]WeeklyCovidStatsByState, 0)

	for row.Next() {
		var scannedRow WeeklyCovidStatsByState
		scanError := row.Scan(
			&scannedRow.State,
			&scannedRow.StartWeek,
			&scannedRow.EndWeek,
			&scannedRow.Confirmed,
			&scannedRow.Death,
			&scannedRow.Cured,
		)

		if scanError != nil {
			return scannedRows, bugsnag.FromError("SQL Weekly Covid Stats By State Scan Error", scanError)
		}
		scannedRows = append(scannedRows, scannedRow)
	}
	return scannedRows, nil
}
