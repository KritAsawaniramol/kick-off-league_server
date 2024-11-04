package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	model "kickoff-league.com/models"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

func (h *httpHandler) DeleteImageBanner(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("compatitionID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "RemoveImageBanner fail"})
		return
	}

	orgID := c.GetUint("organizer_id")

	if err := h.competitionUsecase.RemoveImageBanner(uint(compatitionID), orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "RemoveImageBanner fail"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "RemoveImageBanner successs"})
}

func (h *httpHandler) UpdateImageBanner(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("compatitionID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	in, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest (image)"})
		return
	}

	// extract image extension from original file filename
	isImage, fileExt := isImage(in)
	if !isImage {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file is not an image(png, jpeg)"})
		return
	}

	imagePath := createImagePath(fileExt, "./images/banner")

	if err := c.SaveUploadedFile(in, imagePath); err != nil {
		log.Printf("error: UpdateImageBanner: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "UpdateImageBanner fail"})
		return
	}
	if err := h.competitionUsecase.UpdateImageBanner(uint(compatitionID), c.GetUint("organizer_id"), imagePath); err != nil {
		if err := os.Remove(imagePath); err != nil {
			fmt.Printf("Error removing file: %s\n", err)
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "UpdateImageBanner fail"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "UpdateImageBanner success"})
}

func (h *httpHandler) CreateCompetition(c *gin.Context) {
	reqBody := new(model.CreateCompetition)
	if err := c.BindJSON(reqBody); err != nil {
		log.Printf("error: CreateCompetition: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	reqBody.OrganizerID = c.GetUint("organizer_id")

	if err := h.competitionUsecase.CreateCompetition(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Create compatition success"})
}

func (h *httpHandler) GetCompetition(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("error: GetCompetition: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	result, err := h.competitionUsecase.GetCompetition(uint(teamID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"compatition": result})
}

func (h *httpHandler) FinishCompetition(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("err: FinishCompetition: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.competitionUsecase.FinishCompetition(uint(compatitionID), c.GetUint("organizer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "update status success"})
}

func (h *httpHandler) CancelCompetition(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("error: CancelCompetition: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.competitionUsecase.CancelCompatition(uint(compatitionID), c.GetUint("organizer_id"))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "update status success"})
}

func (h *httpHandler) OpenApplicationCompetition(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("error: OpenApplicationCompetition: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.competitionUsecase.OpenApplicationCompetition(uint(compatitionID), c.GetUint("organizer_id"))
	if err != nil {
		log.Printf("err: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "open application success"})
}

func (h *httpHandler) StartCompetition(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("error: StartCompetition: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.competitionUsecase.StartCompetition(uint(compatitionID), c.GetUint("organizer_id"))
	if err != nil {
		log.Printf("error: StartCompetition: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update compatition status to \"Started\" success"})
}

func (h *httpHandler) GetCompetitions(c *gin.Context) {
	organizerID, err := strconv.ParseUint(c.Query("organizerID"), 10, 64)
	if c.Query("organizerID") != "" && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	// normalUserID, err := strconv.ParseUint(c.Query("normalUserID"), 10, 64)
	// if c.Query("normalUserID") != "" && err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
	// 	return
	// }
	// teamID, err := strconv.ParseUint(c.Query("teamID"), 10, 64)
	// if c.Query("teamID") != "" && err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
	// 	return
	// }

	compatitions, err := h.competitionUsecase.GetCompetitions(&model.GetCompatitionsReq{
		// NormalUserID: uint(normalUserID),
		// TeamID:       uint(teamID),
		OrganizerID: uint(organizerID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"compatition": compatitions})
}

func (h *httpHandler) UpdateCompetition(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("error: UpdateCompetition: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}
	updateCompatition := &model.UpdateCompatition{}
	err = c.BindJSON(updateCompatition)
	if err != nil {
		log.Printf("error: UpdateCompetition: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.competitionUsecase.UpdateCompetition(uint(compatitionID), c.GetUint("organizer_id"), updateCompatition)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "update compatition success"})
}

func (h *httpHandler) AddJoinCode(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("compatitionID"), 10, 64)
	if err != nil {
		log.Printf("error: AddJoinCode: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	number := new(struct {
		Number int `json:"number"`
	})

	err = c.BindJSON(number)
	if err != nil {
		log.Printf("error: AddJoinCode: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.competitionUsecase.AddJoinCode(uint(compatitionID), c.GetUint("organizer_id"), number.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "create join code success"})
}

func (h *httpHandler) JoinCompetition(c *gin.Context) {
	joinModel := &model.JoinCompetition{}
	if err := c.BindJSON(joinModel); err != nil {
		log.Printf("error: JoinCompetition: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}
	err := h.competitionUsecase.JoinCompetition(joinModel, c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "join compatition success"})
}

// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
