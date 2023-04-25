package server_routes

import (
	"bytes"
	"fmt"
	"net/http"
	"service/services/generator"
	"service/services/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mariomac/gostream/stream"
)

var GenerateImage = generator.GenerateImage
var WeeklyTemplate = generator.WeeklyTemplate

func GetResourceImage(c *gin.Context) {
	logger := logger.WithContext(c.HandlerName())

	name := c.Param("filename")
	debugLayout := c.Request.URL.Query().Has("debug")
	template, getTemplateError := getTemplate(name)
	if getTemplateError != nil {
		logger.Error("get resource error: %s")
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("cannot retrieve %s", name),
		})
		return
	}

	configStream := stream.OfSlice(template.Config)
	configStreamToString := stream.Map(configStream, func(tc generator.TextConfig) string {
		return tc.VariableName
	})
	substitutions := stream.Map(configStreamToString, func(s string) []string {
		splitString := strings.Split(s, ":")
		if len(splitString) < 2 {
			return []string{s, s}
		}
		variableName := splitString[0]
		substitutionInArray := splitString[1:]
		joinedSubstitution := strings.Join(substitutionInArray, ":")
		return []string{
			variableName,
			strings.Trim(joinedSubstitution, " "),
		}
	}).ToSlice()

	substitutionsMap := make(map[string]string)
	for index, substitution := range substitutions {
		substitutionsMap[substitution[0]] = substitution[1]
		template.Config[index].VariableName = substitution[0]
	}

	var imgByte bytes.Buffer
	if debugLayout {
		imgByte = GenerateImage(
			template,
			substitutionsMap,
			generator.DebugLayout,
		)
	} else {
		imgByte = GenerateImage(
			template,
			substitutionsMap,
		)
	}

	c.Header("Content-Disposition", "inline; filename=generatedimage.png")
	c.Data(http.StatusOK, "application/octet-stream", imgByte.Bytes())
}
