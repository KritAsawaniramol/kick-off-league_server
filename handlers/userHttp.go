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

// AcceptAddMemberRequest implements UserHandler.
func (h *userHttpHandler) AcceptAddMemberRequest(c *gin.Context) {
	reqBody := new(struct {
		ReqID uint `json:"requestID"`
	})
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.userUsercase.AcceptAddMemberRequest(reqBody.ReqID, c.GetUint("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "AcceptAddMemberRequest success"})
}

// IgnoreAddMemberRequest implements UserHandler.
func (h *userHttpHandler) IgnoreAddMemberRequest(c *gin.Context) {
	reqBody := new(struct {
		ReqID uint `json:"requestID"`
	})
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.userUsercase.IgnoreAddMemberRequest(reqBody.ReqID, c.GetUint("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "IgnoreAddMemberRequest success"})
}

// SendAddMemberRequest implements UserHandler.
func (h *userHttpHandler) SendAddMemberRequest(c *gin.Context) {

	//Get userID
	userID := c.GetUint("user_id")

	reqBody := new(model.AddMemberRequest)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.userUsercase.SendAddMemberRequest(reqBody, userID); err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Create AddMemberRequest failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Create AddMemberRequest success"})

}

// UpdateNormalUser implements UserHandler.
func (h *userHttpHandler) UpdateNormalUser(c *gin.Context) {

	//Get userID
	normalUserID, ok := c.Get("normalUser_id")
	if !ok {
		log.Print("normalUser_id not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	reqBody := new(model.UpdateNormalUser)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.userUsercase.UpdateNormalUser(reqBody, uint(normalUserID.(float64))); err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Update normalUser failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update normalUser success"})
}

// CreateTeam implements UserHandler.
func (h *userHttpHandler) CreateTeam(c *gin.Context) {

	//Get Current userID
	userID, ok := c.Get("user_id")
	if !ok {
		log.Print("user_id not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userID_uint := uint(userID.(float64))

	reqBody := new(model.CreaetTeam)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}
	reqBody.OwnerID = userID_uint

	if err := h.userUsercase.CreateTeam(reqBody); err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	return
}

// UpdateNormalUserPhone implements UserHandler.
// func (h *userHttpHandler) UpdateNormalUserPhone(c *gin.Context) {
// 	userID, err := strconv.ParseUint(c.Param("userID"), 10, 32)
// 	if err != nil {
// 		log.Errorf(err.Error())
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
// 		return
// 	}
// 	phone := c.Param("phone")

// 	err = h.userUsercase.UpdateNormalUserPhone(uint(userID), phone)
// 	if err != nil {
// 		log.Errorf(err.Error())
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Update sucsses"})
// }

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
func (h *userHttpHandler) RegisterOrganizer(c *gin.Context) {
	reqBody := new(model.RegisterUser)
	if err := c.Bind(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}
	if err := h.userUsercase.RegisterOrganizer(reqBody); err != nil {
		response(c, http.StatusInternalServerError, "Register failed")
		return
	}
	response(c, http.StatusOK, "Register success")
}

func (h *userHttpHandler) RegisterNormaluser(c *gin.Context) {
	reqBody := new(model.RegisterNormaluser)
	if err := c.Bind(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}
	fmt.Printf("%v\n", reqBody)
	if err := h.userUsercase.RegisterNormaluser(reqBody); err != nil {
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
