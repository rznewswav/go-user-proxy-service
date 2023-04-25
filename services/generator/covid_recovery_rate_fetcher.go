package generator

import (
	"service/services/mysql"
	"service/services/redis"
	"strconv"
	"time"
)

var CovidRecoveryRateFetcher = DataFetcher{
	FetcherId: "CovidRecoveryRateFetcher",
	FetcherFn: func() (interface{}, error) {
		return redis.GetOrCompute(
			"CovidRecoveryRateFetcher",
			mysql.GetLatestCovidRecoveryRate,
			redis.ExpiryConfig(5*time.Minute),
		)
	},
	GetString: func(fetchedData interface{}, variable string) string {
		cumulative, castable := fetchedData.(mysql.CovidTotalCumulative)
		if !castable {
			return ""
		}

		switch variable {
		case "recovery.rate":
			return strconv.FormatFloat(cumulative.RecoveryRate, 'f', 2, 64) + "%"
		case "weekly.total.cases":
			return englishPrinter.Sprintf("%d", cumulative.Confirmed)
		case "weekly.total.recovered":
			return englishPrinter.Sprintf("%d", cumulative.Cured)
		default:
			return ""
		}
	},
}
