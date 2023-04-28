package resp

import (
	fstruct "github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

// smcopy shallow copy map
func smcopy(data gin.H) (populateMap gin.H) {
	populateMap = make(gin.H)
	for k, v := range data {
		populateMap[k] = v
	}
	return
}

// smcopyMap shallow copy map
func smcopyMap(data map[string]interface{}) (populateMap gin.H) {
	populateMap = make(gin.H)
	for k, v := range data {
		populateMap[k] = v
	}
	return
}

func makeGinH(data ...any) (datum gin.H) {
	datum = make(gin.H)
	if len(data) == 0 {
		return
	}

	if len(data) > 0 {
		mapMaybe := data[0]
		if asMap, isMap := mapMaybe.(map[string]interface{}); isMap {
			datum = smcopyMap(asMap)
		} else if asGinH, isGinH := mapMaybe.(gin.H); isGinH {
			datum = asGinH
		} else {
			fstruct.DefaultTagName = "json"
			asMap := fstruct.Map(mapMaybe)
			datum = smcopyMap(asMap)
		}
	}
	return
}
