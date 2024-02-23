package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	model "kickoff-league.com/models"
	"kickoff-league.com/usecases"
)

type userHttpHandler struct {
	userUsercase usecases.UserUsecase
}

// UpdateNormalUser implements UserHandler.
func (h *userHttpHandler) UpdateNormalUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	//verify user_id
	if userID_token, ok := c.Get("user_id"); ok {
		if userID_token != uint(userID) {

			c.JSON(http.StatusForbidden, "Forbidden")
			return
		}
	} else {
		log.Print("user_id not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	reqBody := new(model.UpdateNormalUser)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	// jsonData, err := io.ReadAll(c.Request.Body)
	if err := h.userUsercase.UpdateNormalUser(reqBody, uint(userID)); err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Update normalUser failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update normalUser success"})
}

// CreateTeam implements UserHandler.
func (h *userHttpHandler) CreateTeam(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	//verify user_id
	if userID_token, ok := c.Get("user_id"); ok {
		if userID_token != uint(userID) {
			c.JSON(http.StatusForbidden, "Forbidden")
			return
		}
	} else {
		log.Print("user_id not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	reqBody := new(model.CreaetTeam)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}
	reqBody.OwnerID = uint(userID)

	if err := h.userUsercase.CreateTeam(reqBody); err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	return
}

// UpdateNormalUserPhone implements UserHandler.
func (h *userHttpHandler) UpdateNormalUserPhone(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 32)
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	phone := c.Param("phone")

	err = h.userUsercase.UpdateNormalUserPhone(uint(userID), phone)
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update sucsses"})
	return
}

// GetUserByPhone implements UserHandler.
func (h *userHttpHandler) GetUsers(c *gin.Context) {
	normalUser, err := h.userUsercase.GetUsers()
	if err != nil {
		response(c, http.StatusBadRequest, "Bad request")
	}

	c.JSON(http.StatusOK, normalUser)
}

func NewUserHttpHandler(userUsercase usecases.UserUsecase) UserHandler {
	return &userHttpHandler{
		userUsercase: userUsercase,
	}
}

// GetUsers implements UserHandler.
// func (h *userHttpHandler) GetUserByPhone(c *gin.Context) {
// 	phone := c.Param("phone")

// 	user, err := h.userUsercase.GetUserByPhone(phone)
// 	if err != nil {
// 		response(c, http.StatusBadRequest, "Bad request")
// 	}

// 	c.JSON(http.StatusOK, user)
// }

// RegisterUser implements UserHandler.
func (h *userHttpHandler) RegisterUser(c *gin.Context) {
	reqBody := new(model.RegisterUser)
	if err := c.Bind(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}
	fmt.Printf("%v\n", reqBody)
	if err := h.userUsercase.Register(reqBody); err != nil {
		response(c, http.StatusInternalServerError, "Register failed")
		return
	}
	response(c, http.StatusOK, "Register success")
}

func (h *userHttpHandler) LoginUser(c *gin.Context) {
	reqBody := new(model.LoginUser)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	jwt, user, err := h.userUsercase.Login(reqBody)
	if err != nil {
		response(c, http.StatusUnauthorized, "Login failed")
		return
	}

	// Set cookie
	c.SetCookie("jwt", jwt, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"user":    user,
	})
}

// GetUser implements UserHandler.
func (h *userHttpHandler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
	}
	user, err := h.userUsercase.GetUser(uint(userID))
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}
