package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"service/services/config"
	"testing"
	"time"
)

type MysqlTestConfig struct {
	DatabaseURL string `env:"TEST_MYSQL_URL" default:"tcp(127.0.0.1:3306)/newswav_news"`
}

func TestQueryWeeklyCovidStatsByState(t *testing.T) {
	dbConfig := config.QuietBuild(MysqlTestConfig{})
	defer func() {
		if Client != nil {
			Client.Close()
		}
	}()

	if dbClient,
		dbConnectionError := sql.Open("mysql", dbConfig.DatabaseURL); dbConnectionError != nil {
		t.Fatal(dbConnectionError)
		return
	} else {
		Client = dbClient
	}

	pingCtx, pingCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCtxCancel()

	if pingError := Client.PingContext(pingCtx); pingError != nil {
		t.Fatal(pingError)
		return
	}

	weeklyStat, weeklyStatError := GetWeeklyCovidStatsByState()
	if weeklyStatError != nil {
		t.Fatal(weeklyStatError)
		return
	}

	for _, wcsbs := range weeklyStat {
		fmt.Printf(
			"State: %s (%s - %s) %d Confirmed, %d Deaths, %d Cured\n",
			wcsbs.State,
			wcsbs.StartWeek.Format("02 Jan"),
			wcsbs.EndWeek.Format("02 Jan 2006"),
			wcsbs.Confirmed,
			wcsbs.Death,
			wcsbs.Cured,
		)
	}
}
