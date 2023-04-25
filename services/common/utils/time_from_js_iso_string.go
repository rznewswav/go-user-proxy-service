package utils

import "time"

func TimeFromJSISOString(str string) (time.Time, error) {
	javascriptISOString := "2006-01-02T15:04:05.999Z"
	dataTime, err := time.Parse(javascriptISOString, str)
	return dataTime, err
}
