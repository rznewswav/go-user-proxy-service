package generator

type DataFetcher struct {
	FetcherId string
	FetcherFn func() (interface{}, error)
	GetString func(fetchedData interface{}, variable string) string
}

func (df DataFetcher) Equal(otherDf DataFetcher) bool {
	return df.FetcherId == otherDf.FetcherId
}

var VariablesToFetcherMap = map[string]DataFetcher{
	"weekly.date":                  CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.johor":      CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.johor":       CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.kedah":      CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.kedah":       CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.kelantan":   CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.kelantan":    CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.melaka":     CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.melaka":      CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.nsembilan":  CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.nsembilan":   CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.pahang":     CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.pahang":      CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.perak":      CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.perak":       CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.perlis":     CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.perlis":      CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.pinang":     CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.pinang":      CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.sabah":      CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.sabah":       CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.sarawak":    CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.sarawak":     CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.selangor":   CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.selangor":    CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.terengganu": CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.terengganu":  CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.kl":         CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.kl":          CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.labuan":     CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.labuan":      CovidWeeklyStatsByStateFetcher,
	"weekly.new.deaths.putrajaya":  CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases.putrajaya":   CovidWeeklyStatsByStateFetcher,
	"weekly.new.cases":             CovidTotalCasesFetcher,
	"weekly.new.recovered":         CovidTotalCasesFetcher,
	"weekly.new.deaths":            CovidTotalCasesFetcher,
	"recovery.rate":                CovidRecoveryRateFetcher,
	"weekly.total.cases":           CovidRecoveryRateFetcher,
	"weekly.total.recovered":       CovidRecoveryRateFetcher,
}
