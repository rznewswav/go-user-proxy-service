package mysql

type SQLStatement string

const SQLGetCovidLatestWeekStatsByState SQLStatement = "" +
	"SELECT " +
	"    DISTINCT state, " +
	"        start_week, " +
	"        end_week, " +
	"        confirmed, " +
	"        death, " +
	"        cured " +
	"FROM corona_malaysia_weekly_state_stats " +
	"WHERE start_week = ( " +
	"    SELECT " +
	"        MAX(cmwss.start_week) AS start_week " +
	"    FROM corona_malaysia_weekly_state_stats cmwss " +
	");"

const SQLGetWeeklyCovidTotalCases SQLStatement = "" +
	"SELECT " +
	"    DISTINCT start_week, " +
	"        end_week, " +
	"        confirmed, " +
	"        death, " +
	"        cured " +
	"FROM corona_malaysia_weekly_stats " +
	"WHERE start_week = ( " +
	"    SELECT " +
	"        MAX(cmws.start_week) AS start_week " +
	"    FROM corona_malaysia_weekly_stats cmws " +
	") LIMIT 1;"

const SQLGetCovidLatestRecoveryRate SQLStatement = "" +
	"SELECT " +
	"    cured_sum / GREATEST(confirmed_sum, 1) * 100 AS recovery_rate " +
	"FROM " +
	"	corona_malaysia_stats " +
	"ORDER BY " +
	"	stats_date DESC " +
	"LIMIT 1;"

const SQLGetCovidLatestCumulative SQLStatement = "" +
	"SELECT " +
	"    cured_sum, " +
	"    confirmed_sum, " +
	"    cured_sum / GREATEST(confirmed_sum, 1) * 100 AS recovery_rate " +
	"FROM " +
	"	corona_malaysia_stats " +
	"ORDER BY " +
	"	stats_date DESC " +
	"LIMIT 1;"
