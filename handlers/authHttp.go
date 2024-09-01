package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	model "kickoff-league.com/models"
)

// RegisterUser implements UserHandler.
func (h *httpHandler) RegisterOrganizer(c *gin.Context) {
	reqBody := new(model.RegisterOrganizer)
	if err := c.BindJSON(reqBody); err != nil {
		response(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authUsecase.RegisterOrganizer(reqBody); err != nil {
		response(c, http.StatusBadRequest, err.Error())
		return
	}
	response(c, http.StatusCreated, "register success")
}

func (h *httpHandler) RegisterNormaluser(c *gin.Context) {
	reqBody := new(model.RegisterNormaluser)
	if err := c.BindJSON(reqBody); err != nil {
		response(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.authUsecase.RegisterNormaluser(reqBody); err != nil {
		response(c, http.StatusBadRequest, err.Error())
		return
	}
	response(c, http.StatusOK, "register success")
}

// LogoutUser implements UserHandler.
func (u *httpHandler) LogoutUser(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/home")
	response(c, http.StatusOK, "logout success")
}

func (h *httpHandler) LoginUser(c *gin.Context) {
	
	reqBody := new(model.LoginUser)
	if err := c.BindJSON(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	jwt, user, err := h.authUsecase.Login(reqBody)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Set cookie
	c.SetCookie("token", jwt, 360000, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"user":    user,
	})
}
