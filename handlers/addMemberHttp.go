package handlers

import (
	"net/http"

	model "kickoff-league.com/models"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

func (h *httpHandler) GetMyPenddingAddMemberRequest(c *gin.Context) {
	addMemberRequests, err := h.addMemberUsecase.GetMyPenddingAddMemberRequest(c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"add_member_requests": addMemberRequests})
}

func (h *httpHandler) SendAddMemberRequest(c *gin.Context) {
	reqBody := new(model.AddMemberRequest)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "error: send add member request failed")
		return
	}

	if err := h.addMemberUsecase.SendAddMemberRequest(reqBody, c.GetUint("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "send add member request success"})
}


func (h *httpHandler) AcceptAddMemberRequest(c *gin.Context) {
	reqBody := new(struct {
		ReqID uint `json:"requestID"`
	})
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.addMemberUsecase.AcceptAddMemberRequest(reqBody.ReqID, c.GetUint("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "accept AddMemberRequest success"})
}


func (h *httpHandler) IgnoreAddMemberRequest(c *gin.Context) {
	reqBody := new(struct {
		ReqID uint `json:"requestID"`
	})
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.addMemberUsecase.IgnoreAddMemberRequest(reqBody.ReqID, c.GetUint("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ignore AddMemberRequest success"})
}


// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================


