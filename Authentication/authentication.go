package authentication

import "github.com/gin-gonic/gin"

type Authentication interface {
	Auth() gin.HandlerFunc
	AuthAdmin() gin.HandlerFunc
	AuthNormalUser() gin.HandlerFunc
	AuthOrganizer() gin.HandlerFunc
}
