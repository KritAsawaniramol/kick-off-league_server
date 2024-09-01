package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	model "kickoff-league.com/models"
	"kickoff-league.com/util"
)

func (h *httpHandler) GetTeams(c *gin.Context) {
	reqBody := new(model.GetTeamsReq)
	teams, err := h.teamUsecase.GetTeams(reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

func (h *httpHandler) GetTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	teams, err := h.teamUsecase.GetTeamWithMemberAndCompatitionByID(uint(teamID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

func (h *httpHandler) GetTeamByOwnerID(c *gin.Context) {
	teams, err := h.teamUsecase.GetTeamsByOwnerID(c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

func (h *httpHandler) CreateTeam(c *gin.Context) {
	reqBody := new(model.CreateTeam)
	if err := c.BindJSON(reqBody); err != nil {
		log.Printf("error: CreateTeam: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}
	reqBody.OwnerID = c.GetUint("user_id")
	if err := h.teamUsecase.CreateTeam(reqBody); err != nil {
		switch err.(type) {
		case *util.CreateTeamError:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "create team success"})
}

func (h *httpHandler) RemoveTeamMember(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 64)
	if err != nil {
		log.Printf("error: RemoveTeamMember: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	normalUserID := new(struct {
		NormalUserID uint `json:"normal_user_id"`
	})

	err = c.BindJSON(normalUserID)
	if err != nil {
		log.Printf("error: RemoveTeamMember: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	if err := h.teamUsecase.RemoveNormalUserFormTeam(uint(teamID), normalUserID.NormalUserID, c.GetUint("user_id")); err != nil {
		log.Printf("error: RemoveTeamMember: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "remove member success"})
}

func (h *httpHandler) RemoveCompatitionTeam(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("competitionID"), 10, 64)
	if err != nil {
		log.Printf("Error: RemoveCompatitionTeam failed: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	teamID := new(struct {
		TeamID uint `json:"team_id"`
	})

	err = c.BindJSON(teamID)
	if err != nil {
		log.Printf("Error: RemoveCompatitionTeam failed: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	if err := h.teamUsecase.RemoveTeamFormCompatition(teamID.TeamID, uint(compatitionID), c.GetUint("organizer_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "remove team success"})
}

func (h *httpHandler) DeleteTeamImageProfile(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 64)
	if err != nil {
		log.Printf("Error: DeleteTeamImageProfile failed: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}
	if err := h.teamUsecase.RemoveTeamImageProfile(uint(teamID), c.GetUint("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "remove team's image profile successs"})
}

func (h *httpHandler) DeleteTeamImageCover(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 64)
	if err != nil {
		log.Printf("Error: DeleteTeamImageCover failed: %s\n", err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}
	if err := h.teamUsecase.RemoveTeamImageCover(uint(teamID), c.GetUint("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "remove team's image cover successs"})
}

func (h *httpHandler) UpdateTeamImageCover(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 64)
	if err != nil {
		log.Printf("Error: UpdateTeamImageCover failed: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "update team's cover image failed"})
		return
	}

	in, err := c.FormFile("image")
	if err != nil {
		log.Printf("Error: UpdateTeamImageCover failed: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "update team's cover image failed"})
		return
	}

	// extract image extension from original file filename
	isImage, fileExt := isImage(in)
	if !isImage {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file is not an image(png, jpeg)"})
		return
	}

	imagePath := createImagePath(fileExt, "./images/cover")

	if err := c.SaveUploadedFile(in, imagePath); err != nil {
		log.Printf("Error: UpdateTeamImageCover failed: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "update team's cover image failed"})
		return
	}

	if err := h.teamUsecase.UpdateTeamImageCover(uint(teamID), imagePath, c.GetUint("user_id")); err != nil {
		if err := os.Remove(imagePath); err != nil {
			log.Printf("Error removing file: %s\n", err)
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "update team's cover image success"})
}

func (h *httpHandler) UpdateTeamImageProfile(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 64)
	if err != nil {
		log.Printf("Error: UpdateTeamImageProfile failed: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "update team's profile image failed"})
		return
	}

	in, err := c.FormFile("image")
	if err != nil {
		log.Printf("Error: UpdateTeamImageProfile failed: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "update team's profile image failed"})
		return
	}

	// extract image extension from original file filename
	isImage, fileExt := isImage(in)
	if !isImage {
		log.Printf("Error: UpdateTeamImageProfile failed: is not image\n")
		c.JSON(http.StatusBadRequest, gin.H{"message": "file is not an image(png, jpeg)"})
		return
	}

	imagePath := createImagePath(fileExt, "./images/profile")
	if err := c.SaveUploadedFile(in, imagePath); err != nil {
		log.Printf("Error: UpdateTeamImageProfile failed: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "update team's profile image failed"})
		return
	}

	if err := h.teamUsecase.UpdateTeamImageProfile(uint(teamID), imagePath, c.GetUint("user_id")); err != nil {
		log.Errorf(err.Error())
		if err := os.Remove(imagePath); err != nil {
			fmt.Printf("Error removing file: %s\n", err)
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "update team's profile image success"})
}

// ========================================================================
// ========================================================================
// ========================================================================
// ========================================================================
// ========================================================================
// ========================================================================
