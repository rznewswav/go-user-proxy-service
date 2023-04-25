package generator

import (
	"service/services/mysql"
	"service/services/redis"
	"time"
)

var CovidTotalCasesFetcher = DataFetcher{
	FetcherId: "CovidTotalCasesFetcher",
	FetcherFn: func() (interface{}, error) {
		return redis.GetOrCompute(
			"CovidTotalCasesFetcher",
			mysql.GetWeeklyCovidTotalCases,
			redis.ExpiryConfig(5*time.Minute),
		)
	},
	GetString: func(fetchedData interface{}, variable string) string {
		input, castable := fetchedData.(mysql.WeeklyCovidTotalStats)
		if !castable {
			return ""
		}

		switch variable {
		case "weekly.new.cases":
			return englishPrinter.Sprintf("%d", input.Confirmed)
		case "weekly.new.recovered":
			return englishPrinter.Sprintf("%d", input.Cured)
		case "weekly.new.deaths":
			return englishPrinter.Sprintf("%d", input.Death)
		default:
			return ""
		}
	},
}
