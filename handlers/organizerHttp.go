package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"

	model "kickoff-league.com/models"
)

func (h *httpHandler) GetOrganizers(c *gin.Context) {
	org, err := h.organizerUsecase.GetOrganizers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"organizer": org})
}

func (h *httpHandler) GetOrganizer(c *gin.Context) {
	organizerID, err := strconv.ParseUint(c.Param("organizerID"), 10, 32)
	if err != nil {
		log.Printf("err: GetOrganizer: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	org, err := h.organizerUsecase.GetOrganizer(uint(organizerID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organizer": org})
}

func (h *httpHandler) UpdateOrganizer(c *gin.Context) {

	updateOrganizer := &model.UpdateOrganizer{}

	if err := c.BindJSON(updateOrganizer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	if err := h.organizerUsecase.UpdateOrganizer(c.GetUint("organizer_id"), updateOrganizer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update organizer success"})
}

// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
