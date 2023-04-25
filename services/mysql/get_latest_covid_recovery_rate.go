package mysql

import (
	"context"
	"service/services/bugsnag"
	"time"
)

type CovidTotalCumulative struct {
	Cured        int
	Confirmed    int
	RecoveryRate float64
}

func GetLatestCovidRecoveryRate() (CovidTotalCumulative, error) {
	queryCtx, queryCtxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer queryCtxCancel()
	row, queryError := Client.QueryContext(queryCtx, string(SQLGetCovidLatestCumulative))
	if queryError != nil {
		return CovidTotalCumulative{}, bugsnag.FromError("SQL Latest Covid Recovery Rate Error", queryError)
	}
	if row == nil {
		return CovidTotalCumulative{}, bugsnag.FromError("SQL Latest Covid Recovery Rate Row Nil Error", ErrRowNil)
	}

	defer row.Close()

	if !row.Next() {
		return CovidTotalCumulative{}, nil
	}

	var cumulative CovidTotalCumulative
	row.Scan(
		&cumulative.Cured,
		&cumulative.Confirmed,
		&cumulative.RecoveryRate,
	)
	return cumulative, nil
}
