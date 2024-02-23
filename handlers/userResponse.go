package handlers

import (
	"github.com/gin-gonic/gin"
)

type baseResponse struct {
	Message string `json:"message"`
}

func response(c *gin.Context, responseCode int, message string) {
	c.JSON(responseCode, gin.H{"message": message})
}

// c.JSON in echo has return error
// func response(c echo.Context, responseCode int, message string) error {
// 	return c.JSON(responseCode, &baseResponse{
// 		Message: message,
// 	})
// }
