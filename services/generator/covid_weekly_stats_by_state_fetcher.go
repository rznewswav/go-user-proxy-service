package generator

import (
	"fmt"
	"service/services/mysql"
	"service/services/redis"
	"strings"
	"time"
)

var lowercaseStateNameToDBStateName = map[string]string{
	"johor":      "Johor",
	"kedah":      "Kedah",
	"kelantan":   "Kelantan",
	"melaka":     "Melaka",
	"nsembilan":  "Negeri Sembilan",
	"pahang":     "Pahang",
	"perak":      "Perak",
	"perlis":     "Perlis",
	"pinang":     "P. Pinang",
	"sabah":      "Sabah",
	"sarawak":    "Sarawak",
	"selangor":   "Selangor",
	"terengganu": "Terengganu",
	"kl":         "Kuala Lumpur",
	"labuan":     "Labuan",
	"putrajaya":  "Putrajaya",
}

type WeeklyCovidStatsByStateMap = map[string]mysql.WeeklyCovidStatsByState

var CovidWeeklyStatsByStateFetcher = DataFetcher{
	FetcherId: "CovidWeeklyStatsByStateFetcher",
	FetcherFn: func() (interface{}, error) {
		statArray, sqlError := redis.GetOrCompute(
			"CovidWeeklyStatsByStateFetcher",
			mysql.GetWeeklyCovidStatsByState,
			redis.ExpiryConfig(5*time.Minute),
		)
		if sqlError != nil {
			return WeeklyCovidStatsByStateMap{}, sqlError
		}

		statsMap := make(WeeklyCovidStatsByStateMap)

		for _, wcsbs := range statArray {
			statsMap[wcsbs.State] = wcsbs
		}

		return statsMap, nil
	},
	GetString: func(uncasted interface{}, variable string) string {
		input, castable := uncasted.(WeeklyCovidStatsByStateMap)
		if !castable {
			return ""
		}

		if variable == "weekly.date" {
			var firstStats mysql.WeeklyCovidStatsByState
			for _, wcsbs := range input {
				firstStats = wcsbs
				break
			}
			return fmt.Sprintf(
				"%s - %s",
				strings.ToUpper(firstStats.StartWeek.Format("02 Jan")),
				strings.ToUpper(firstStats.EndWeek.Format("02 Jan 2006")),
			)
		}

		splitVariable := strings.Split(variable, ".")
		if len(splitVariable) != 4 {
			return ""
		}

		expectedWeeklyPrefix := splitVariable[0]
		dataPrefix := splitVariable[1]
		/** this is either "deaths" or "cases" */
		statsType := splitVariable[2]
		stateNameLowerCase := splitVariable[3]

		if expectedWeeklyPrefix != "weekly" && dataPrefix != "new" {
			return ""
		}

		fullStateName, hasFullStateName := lowercaseStateNameToDBStateName[stateNameLowerCase]
		if !hasFullStateName {
			return ""
		}

		stat, hasStat := input[fullStateName]

		if !hasStat {
			return ""
		}

		switch statsType {
		case "deaths":
			return englishPrinter.Sprintf("%d", stat.Death)
		case "cases":
			return englishPrinter.Sprintf("%d", stat.Confirmed)
		default:
			return ""
		}
	},
}
