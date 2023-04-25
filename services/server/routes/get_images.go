package server_routes

import (
	"fmt"
	"net/http"
	"regexp"
	"service/services/bugsnag"
	"service/services/generator"
	"service/services/logger"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mariomac/gostream/stream"
)

var whitelistedImages = map[string]string{
	"covid-weekly-by-state.png":      "weekly-by-state.json",
	"covid-weekly-new-recovered.png": "weekly-with-recovered.json",
	"covid-weekly-new-total.png":     "weekly.json",
}

type M = map[string]interface{}

var etagRegex = regexp.MustCompile(`(?m)W/\"expirein:(\d+)\"`)

func GetCovidImages(c *gin.Context) {
	logger := logger.WithContext(c.HandlerName())
	currentEtag := c.Request.Header.Get("if-none-match")
	if len(currentEtag) > 0 {
		matches := etagRegex.FindStringSubmatch(currentEtag)
		if len(matches) >= 2 {
			firstGroupMatch := matches[1]
			nowMilli := time.Now().UnixMilli()
			parsedInt, parseError := strconv.ParseInt(string(firstGroupMatch), 10, 64)
			if parseError == nil && nowMilli < parsedInt {
				c.Status(http.StatusNotModified)
				return
			}
		}
	}
	name := c.Param("filename")
	var configJsonName string

	if config, filenameExists := whitelistedImages[name]; !filenameExists {
		c.JSON(http.StatusNotFound, M{
			"success": false,
			"message": M{
				"code":    "NOT_FOUND",
				"message": fmt.Sprintf("file %s is not found", name),
			},
		})
		return
	} else {
		configJsonName = config
	}

	template, getTemplateError := getTemplate(configJsonName)
	if getTemplateError != nil {
		logger.Error("get resource error: %s")
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("cannot retrieve %s", name),
		})
		return
	}

	configStream := stream.OfSlice(template.Config)
	variables := stream.Map(configStream, func(tc generator.TextConfig) string {
		return strings.Split(tc.VariableName, ":")[0]
	}).ToSlice()

	variablesToFetcherId := make(map[string]string)
	uniqueFetchers := make(map[string]generator.DataFetcher)
	fetchedValues := make(map[string]interface{})
	for _, v := range variables {
		if fetcher, hasFetcher := generator.VariablesToFetcherMap[v]; hasFetcher {
			variablesToFetcherId[v] = fetcher.FetcherId
			uniqueFetchers[fetcher.FetcherId] = fetcher
		}
	}

	for fetcherId, fetcher := range uniqueFetchers {
		if fetchedValue, fetchingError := fetcher.FetcherFn(); fetchingError != nil {
			logger.Error("Error fetching for %s", fetcher.FetcherId, bugsnag.FromError("Fetcher Error", fetchingError))
			continue
		} else {
			fetchedValues[fetcherId] = fetchedValue
		}
	}

	substitutionsMap := make(map[string]string)
	for index, variable := range variables {
		template.Config[index].VariableName = variable
		fetcherId, hasFetcherId := variablesToFetcherId[variable]
		if !hasFetcherId {
			continue
		}
		fetcher, hasFetcher := uniqueFetchers[fetcherId]
		if !hasFetcher {
			continue
		}
		fetchedValue, hasFetchedValue := fetchedValues[fetcherId]
		if !hasFetchedValue {
			continue
		}
		substitutionsMap[variable] = fetcher.GetString(fetchedValue, variable)

	}

	imgByte := GenerateImage(
		template,
		substitutionsMap,
	)

	contentDispositionValue := fmt.Sprintf("inline; filename=\"%s\"", name)
	c.Header("Content-Disposition", contentDispositionValue)
	c.Header("Cache-Control", "public, max-age=300")

	shouldExpiresIn := time.Now().Add(5 * time.Minute)
	etagValue := fmt.Sprintf("W/\"expirein:%d\"", shouldExpiresIn.UnixMilli())
	c.Header("Etag", etagValue)
	c.Data(http.StatusOK, "image/png", imgByte.Bytes())
}
