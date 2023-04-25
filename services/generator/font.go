package generator

import (
	"fmt"
	"strings"

	"github.com/tdewolff/canvas"
)

const DefaultFontPath = "resources/fonts/nunito-v23-latin-regular.ttf"
const DefaultFontName = "Arial"

var loadedFonts = make(map[string]*canvas.Font)

func LoadFont(pathOrName ...string) (font *canvas.Font) {
	var fontPathOrName = DefaultFontPath
	// var fontPathOrName = defaultFontName
	if len(pathOrName) > 0 {
		fontPathOrName = pathOrName[0]
	}
	if len(pathOrName) > 1 {
		fmt.Printf("cannot load multiple fonts at once, loading only %s, ignoring: %s", fontPathOrName, strings.Join(pathOrName[1:], ", "))
	}

	if alreadyLoaded, hasLoadedBefore := loadedFonts[fontPathOrName]; hasLoadedBefore {
		return alreadyLoaded
	}

	var fontError error
	if strings.HasPrefix(fontPathOrName, "resources") {
		font, fontError = canvas.LoadFontFile(fontPathOrName, canvas.FontRegular)
		if fontError != nil {
			panic(fontError)
		}
	} else {
		font, fontError = canvas.LoadLocalFont(fontPathOrName, canvas.FontRegular)
		if fontError != nil {
			panic(fontError)
		}
	}

	loadedFonts[fontPathOrName] = font

	return font
}
