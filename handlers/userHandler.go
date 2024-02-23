package handlers

import "github.com/gin-gonic/gin"

type UserHandler interface {
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateNormalUserPhone(c *gin.Context)
	CreateTeam(c *gin.Context)
	UpdateNormalUser(c *gin.Context)
	// GetUserByPhone(c *gin.Context)
}
