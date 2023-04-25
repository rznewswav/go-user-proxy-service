package server_routes

import (
	"fmt"
	"net/http"
	"service/services/logger"

	"github.com/gin-gonic/gin"
)

func GetTemplateByName(c *gin.Context) {
	logger := logger.WithContext(c.HandlerName())
	name := c.Param("filename")
	template, getTemplateError := getTemplate(name)
	if getTemplateError != nil {
		logger.Error("get resource error: %s")
		c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": fmt.Sprintf("cannot retrieve %s", name),
		})
		return
	}
	c.JSON(http.StatusOK, template)
}
