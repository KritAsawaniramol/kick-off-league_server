package handlers

import "github.com/gin-gonic/gin"

type UserHandler interface {
	RegisterOrganizer(c *gin.Context)
	RegisterNormaluser(c *gin.Context)
	LoginUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	GetTeams(c *gin.Context)
	GetTeam(c *gin.Context)
	GetMyPenddingAddMemberRequest(c *gin.Context)
	// uploadImage(c *gin.Context)
	// UpdateNormalUserPhone(c *gin.Context)
	UpdateImageCover(c *gin.Context)
	UpdateImageProfile(c *gin.Context)
	UploadImage(c *gin.Context)
	CreateTeam(c *gin.Context)
	CreateCompatition(c *gin.Context)
	UpdateNormalUser(c *gin.Context)
	SendAddMemberRequest(c *gin.Context)
	AcceptAddMemberRequest(c *gin.Context)
	IgnoreAddMemberRequest(c *gin.Context)
	// GetUserByPhone(c *gin.Context)
}
