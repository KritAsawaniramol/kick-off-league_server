package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	model "kickoff-league.com/models"
)

// RegisterUser implements UserHandler.
func (h *httpHandler) RegisterOrganizer(c *gin.Context) {
	reqBody := new(model.RegisterOrganizer)
	if err := c.BindJSON(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}
	if err := h.userUsercase.RegisterOrganizer(reqBody); err != nil {
		if err.Error() == "invalid email format" ||
			err.Error() == "this email is already in use" ||
			err.Error() == "this phone is already in use" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "InternalServerError",
			})
		}
		return
	}
	response(c, http.StatusOK, "register success")
}

func (h *httpHandler) RegisterNormaluser(c *gin.Context) {
	reqBody := new(model.RegisterNormaluser)
	if err := c.BindJSON(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}
	if err := h.userUsercase.RegisterNormaluser(reqBody); err != nil {
		if err.Error() == "invalid email format" ||
			err.Error() == "this email is already in use" ||
			err.Error() == "this username is already in use" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "InternalServerError",
			})
		}
		return
	}
	response(c, http.StatusOK, "register success")
}

// LogoutUser implements UserHandler.
func (u *httpHandler) LogoutUser(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/home")
	c.JSON(http.StatusOK, gin.H{
		"message": "logout success",
	})
}

func (h *httpHandler) LoginUser(c *gin.Context) {
	reqBody := new(model.LoginUser)
	if err := c.BindJSON(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	jwt, user, err := h.userUsercase.Login(reqBody)
	if err != nil {
		if err.Error() == "invalid email format" ||
			err.Error() == "incorrect email or password" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "InternalServerError",
			})
		}
		return
	}

	// Set cookie
	fmt.Println("jwt", jwt)
	c.SetCookie("token", jwt, 360000, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"user":    user,
	})
}
