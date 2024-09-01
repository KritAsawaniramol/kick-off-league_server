package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	model "kickoff-league.com/models"
)

func (h *httpHandler) GetNormalUsers(c *gin.Context) {
	normalUserList, err := h.normalUserUsecase.GetNormalUserList()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, normalUserList)
}

func (h *httpHandler) GetNormalUser(c *gin.Context) {
	normalUserID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}
	normalUser, err := h.normalUserUsecase.GetNormalUser(uint(normalUserID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"normalUser": normalUser,
	})
}

func (h *httpHandler) UpdateNormalUser(c *gin.Context) {
	//Get userID
	normalUserID := c.GetUint("normal_user_id")
	reqBody := new(model.UpdateNormalUser)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.normalUserUsecase.UpdateNormalUser(reqBody, normalUserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Update normalUser failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update normalUser success"})
}

// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================


