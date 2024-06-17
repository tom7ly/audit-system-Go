package handler

import (
	"audit-system/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var auditLogService *service.AuditLogService

func InitAuditLogHandler(als *service.AuditLogService) {
	auditLogService = als
}

func GetAllAuditLogs(c *gin.Context) {
	logs, err := auditLogService.GetAllAuditLogs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func GetAuditLogsByEmail(c *gin.Context) {
	email := c.Param("email")
	logs, err := auditLogService.GetAuditLogsByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
