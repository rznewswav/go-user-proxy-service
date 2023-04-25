package main

import (
	"os"
	generator_service "service/services/generator"
)

var GenerateImage = generator_service.GenerateImage
var WeeklyTemplate = generator_service.WeeklyTemplate

func main() {
	imgByte := GenerateImage(WeeklyTemplate, generator_service.M{
		"{weekly new cases}": "501,784",
	})

	os.WriteFile("/tmp/out.png", imgByte.Bytes(), 0755)
}
