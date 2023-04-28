package resp

import "github.com/gin-gonic/gin"

// smcopy shallow copy map
func smcopy(data gin.H) (populateMap gin.H) {
	populateMap = make(gin.H)
	for k, v := range data {
		populateMap[k] = v
	}
	return
}
