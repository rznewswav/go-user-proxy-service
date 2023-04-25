package generator

import (
	"github.com/tdewolff/canvas"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var englishPrinter = message.NewPrinter(language.English)

var WeeklyTemplate = ImageTemplate{
	"resources/images/Confirmed (Weekly & Total).png",
	[]TextConfig{
		{
			"{weekly new cases}",
			0.5,
			0.44,
			55,
			128,
			18,
			canvas.Center,
			canvas.Black,
			"resources/fonts/nunito-v23-latin-regular.ttf",
		},
		{
			"{weekly total cases}",
			0.5,
			0.20,
			55,
			128,
			18,
			canvas.Center,
			canvas.Black,
			"resources/fonts/nunito-v23-latin-regular.ttf",
		},
	},
}
