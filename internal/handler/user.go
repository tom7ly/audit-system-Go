package handler

import (
	"audit-system/internal/model"
	"audit-system/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var userService *service.UserService

func InitUserHandler(us *service.UserService) {
	userService = us
}

func CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		handleRequestParsingError(c, err)
		return
	}
	createdUser, err := userService.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, createdUser)
}

func GetUsers(c *gin.Context) {
	users, err := userService.GetUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := userService.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	email := c.Param("email")
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		handleRequestParsingError(c, err)
		return
	}
	updatedUser, err := userService.UpdateUser(c.Request.Context(), email, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
func DeleteUser(c *gin.Context) {
	email := c.Param("email")
	if err := userService.DeleteUser(c.Request.Context(), email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
