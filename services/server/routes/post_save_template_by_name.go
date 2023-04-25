package server_routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"service/services/generator"

	"github.com/gin-gonic/gin"
)

func PostSaveTemplateByName(c *gin.Context) {
	var requestBody generator.ImageTemplate
	c.ShouldBindJSON(&requestBody)
	weeklyTemplateByte, weeklyMarshalErr := json.MarshalIndent(requestBody, "", "    ")
	name := c.Param("filename")

	if weeklyMarshalErr != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": weeklyMarshalErr.Error(),
		})
		return
	}

	templatesDir := "resources/templates"
	os.MkdirAll(templatesDir, 0755)

	templateFileName := fmt.Sprintf("%s/%s", templatesDir, name)
	writeErr := os.WriteFile(templateFileName, weeklyTemplateByte, 0755)

	if writeErr != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": writeErr.Error(),
		})
		return
	}

	// append non-indented as jsonl history
	writeHistoryDir := "resources/templates-history"
	writeHistoryFileName := fmt.Sprintf("%s/%sl", writeHistoryDir, name)
	os.MkdirAll(writeHistoryDir, 0755)
	weeklyUpdateHistoryTemplateByte, weeklyUpdateHistoryMarshalErr := json.Marshal(requestBody)
	if weeklyMarshalErr != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": weeklyUpdateHistoryMarshalErr.Error(),
		})
		return
	}

	appendFileHandler, openError := os.OpenFile(writeHistoryFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if openError != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": openError.Error(),
		})
		return
	}

	appendFileHandler.Write(weeklyUpdateHistoryTemplateByte)
	appendFileHandler.WriteString("\n")

	c.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}
