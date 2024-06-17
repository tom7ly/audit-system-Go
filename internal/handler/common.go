package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleRequestParsingError(c *gin.Context, err error) {
	if err == io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body cannot be empty"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
