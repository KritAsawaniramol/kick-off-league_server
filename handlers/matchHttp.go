package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	model "kickoff-league.com/models"
)

func (h *httpHandler) GetMatch(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("matchID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	match, err := h.matchUsecase.GetMatch(uint(matchID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"match": match})
}

func (h *httpHandler) UpdateMatch(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}
	updateMatch := &model.UpdateMatch{}
	err = c.BindJSON(updateMatch)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.matchUsecase.UpdateMatch(uint(matchID), c.GetUint("organizer_id"), updateMatch)
	if err != nil {
		fmt.Printf("err: %v\n", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "update match success"})
}

func (h *httpHandler) GetNextMatch(c *gin.Context) {

	nextMatchs, err := h.matchUsecase.GetNextMatch(c.GetUint("normal_user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"next_match": nextMatchs})
}

// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
