package handlers

import (
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	model "kickoff-league.com/models"
	"kickoff-league.com/util"
)

// RemoveTeamMember implements Handler.
func (h *httpHandler) RemoveTeamMember(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	normalUserID := new(struct {
		NormalUserID uint `json:"normal_user_id"`
	})

	err = c.BindJSON(normalUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	if err := h.userUsercase.RemoveNormalUserFormTeam(uint(teamID), normalUserID.NormalUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "RemoveTeamMember success"})
}

// GetNormalUsers implements Handler.
func (h *httpHandler) GetNormalUsers(c *gin.Context) {
	normalUserList, err := h.userUsercase.GetNormalUserList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
	}
	c.JSON(http.StatusOK, normalUserList)

}

// DeleteImageProfile implements UserHandler.
func (h *httpHandler) DeleteImageProfile(c *gin.Context) {
	normalUserID := c.GetUint("normal_user_id")
	if err := h.userUsercase.RemoveImageProfile(normalUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "RemoveImageProfile successs"})
}

func createImagePath(fileExt string, dst string) string {
	// generate new uuid for image name
	uniqueID := uuid.New()

	// remove "- from imageName"
	filename := strings.Replace(uniqueID.String(), "-", "", -1)

	// generate image from filename and extension
	image := fmt.Sprintf("%s.%s", filename, fileExt)

	imagePath := fmt.Sprintf("%s/%s", dst, image)
	return imagePath
}

// UpdateImageCover implements UserHandler.
func (h *httpHandler) UpdateImageCover(c *gin.Context) {
	userID := c.GetUint("user_id")

	in, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
	}

	if err := h.userUsercase.UpdateImageCover(userID, imagePath); err != nil {
		if err := os.Remove(imagePath); err != nil {
			fmt.Printf("Error removing file: %s\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "UpdateImageProfile success"})
}

func (h *httpHandler) UpdateImageProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	in, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	// extract image extension from original file filename
	isImage, fileExt := isImage(in)
	if !isImage {
		fmt.Println("error: is not image")

		c.JSON(http.StatusBadRequest, gin.H{"message": "file is not an image(png, jpeg)"})
		return
	}

	imagePath := createImagePath(fileExt, "./images/profile")

	if err := c.SaveUploadedFile(in, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
	}

	if err := h.userUsercase.UpdateImageProfile(userID, imagePath); err != nil {
		log.Errorf(err.Error())
		if err := os.Remove(imagePath); err != nil {
			fmt.Printf("Error removing file: %s\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}

	// if err := h.userUsercase.UpdateNormalUser(reqBody, normalUserID); err != nil {
	// 	log.Errorf(err.Error())
	// 	if err := os.Remove(imagePath); err != nil {
	// 		fmt.Printf("Error removing file: %s\n", err)
	// 	}
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
	// 	return
	// }
	c.JSON(http.StatusOK, gin.H{"message": "UpdateImageProfile success"})
}

// UploadFile implements UserHandler.
func (*httpHandler) UploadImage(c *gin.Context) {
	in, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	// extract image extension from original file filename
	isImage, fileExt := isImage(in)
	if !isImage {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File is not an image(png, jpeg)"})
		return
	}

	util.PrintObjInJson(in)

	// generate new uuid for image name
	uniqueID := uuid.New()

	// remove "- from imageName"
	filename := strings.Replace(uniqueID.String(), "-", "", -1)

	// generate image from filename and extension
	image := fmt.Sprintf("%s.%s", filename, fileExt)

	if err := c.SaveUploadedFile(in, "./images/"+image); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "file upload complete!"})
}

func isImage(fileHeader *multipart.FileHeader) (bool, string) {
	file, err := fileHeader.Open()
	if err != nil {
		return false, ""
	}
	defer file.Close()
	_, format, err := image.Decode(file)
	fmt.Println("format", format)
	if err != nil {
		log.Error(err)
		return false, ""
	}
	switch format {
	case "jpeg":
		return true, "jpeg"
	case "png":
		return true, "png"
	default:
		return false, ""
	}
}

func (h *httpHandler) CreateTeam(c *gin.Context) {
	reqBody := new(model.CreateTeam)
	if err := c.BindJSON(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}
	reqBody.OwnerID = c.GetUint("user_id")
	if err := h.userUsercase.CreateTeam(reqBody); err != nil {
		switch err.(type) {
		case *util.CreateTeamError:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "CreateTeam success"})
}

// GetMatchResult implements Handler.
func (h *httpHandler) GetMatchResult(c *gin.Context) {

}

// GetNextMatch implements Handler.
func (h *httpHandler) GetNextMatch(c *gin.Context) {
	normalUserID, err := strconv.ParseUint(c.Param("normalUserID"), 10, 64)
	if err != nil {
		fmt.Println("hello world")
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	nextMatchs, err := h.userUsercase.GetNextMatch(uint(normalUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"next_match": nextMatchs})
}

// UpdateMatch implements Handler.
func (h *httpHandler) UpdateMatch(c *gin.Context) {
	matchID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
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

	err = h.userUsercase.UpdateMatch(uint(matchID), updateMatch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Create match success"})
}

// UpdateCompatition implements Handler.
func (h *httpHandler) UpdateCompatition(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}
	updateCompatition := &model.UpdateCompatition{}
	err = c.BindJSON(updateCompatition)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.userUsercase.UpdateCompatition(uint(compatitionID), updateCompatition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Create compatition success"})
}

// JoinCompatition implements Handler.
func (h *httpHandler) JoinCompatition(c *gin.Context) {
	joinModel := &model.JoinCompatition{}
	if err := c.BindJSON(joinModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
	}
	err := h.userUsercase.JoinCompatition(joinModel)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unable") {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			fmt.Printf("err: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Create compatition success"})
}

// CreateCompatition implements UserHandler.
func (h *httpHandler) CreateCompatition(c *gin.Context) {
	organizerID := c.GetUint("organizer_id")
	reqBody := new(model.CreateCompatition)
	if err := c.BindJSON(reqBody); err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	reqBody.OrganizerID = organizerID

	if err := h.userUsercase.CreateCompatition(reqBody); err != nil {
		if err.Error() == "number of Team for create competition(tounament) is not power of 2" ||
			err.Error() == "number of Team have to morn than 1" ||
			err.Error() == "undefined compatition type" {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Create compatition success"})
}

// GetCompatition implements Handler.
func (h *httpHandler) GetCompatition(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	result, err := h.userUsercase.GetCompatition(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"compatition": result})
}

// GetCompatitions implements Handler.
func (h *httpHandler) GetCompatitions(c *gin.Context) {

	organizerID, err := strconv.ParseUint(c.Query("organizerID"), 10, 64)
	if c.Query("organizerID") != "" && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	normalUserID, err := strconv.ParseUint(c.Query("normalUserID"), 10, 64)
	if c.Query("normalUserID") != "" && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	teamID, err := strconv.ParseUint(c.Query("teamID"), 10, 64)
	if c.Query("teamID") != "" && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	compatitions, err := h.userUsercase.GetCompatitions(&model.GetCompatitionsReq{
		NormalUserID: uint(normalUserID),
		TeamID:       uint(teamID),
		OrganizerID:  uint(organizerID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"compatition": compatitions})
}

// GetMyPenddingAddMemberRequest implements UserHandler.
func (h *httpHandler) GetMyPenddingAddMemberRequest(c *gin.Context) {
	addMemberRequests, err := h.userUsercase.GetMyPenddingAddMemberRequest(c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
	}
	c.JSON(http.StatusOK, gin.H{"add_member_requests": addMemberRequests})
}

// GetTeam implements UserHandler.
func (h *httpHandler) GetTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
	}
	teams, err := h.userUsercase.GetTeamWithMemberAndCompatitionByID(uint(teamID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

// GetTeamByOwnerID implements Handler.
func (h *httpHandler) GetTeamByOwnerID(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("ownerid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
	}
	teams, err := h.userUsercase.GetTeamsByOwnerID(uint(teamID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

// GetTeamList implements UserHandler.
func (h *httpHandler) GetTeams(c *gin.Context) {
	reqBody := new(model.GetTeamsReq)
	// if err := c.BindJSON(reqBody); err != nil {
	// 	log.Errorf("Error binding request body: %v", err)
	// 	response(c, http.StatusBadRequest, "Bad request")
	// 	return
	// }

	util.PrintObjInJson(reqBody)
	teams, err := h.userUsercase.GetTeams(reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

// AcceptAddMemberRequest implements UserHandler.
func (h *httpHandler) AcceptAddMemberRequest(c *gin.Context) {
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
func (h *httpHandler) IgnoreAddMemberRequest(c *gin.Context) {
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
func (h *httpHandler) SendAddMemberRequest(c *gin.Context) {

	//Get userID
	userID := c.GetUint("user_id")

	reqBody := new(model.AddMemberRequest)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.userUsercase.SendAddMemberRequest(reqBody, userID); err != nil {
		if err_str := err.Error(); err_str == "this user isn't owner's team" ||
			err_str == "this user already invited" ||
			err_str == "this user already in team" {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Create AddMemberRequest failed"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Create AddMemberRequest success"})
}

// UpdateNormalUser implements UserHandler.
func (h *httpHandler) UpdateNormalUser(c *gin.Context) {
	//Get userID
	normalUserID := c.GetUint("normal_user_id")

	reqBody := new(model.UpdateNormalUser)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.userUsercase.UpdateNormalUser(reqBody, normalUserID); err != nil {
		log.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Update normalUser failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update normalUser success"})
}

// UpdateNormalUserPhone implements UserHandler.
// func (h *httpHandler) UpdateNormalUserPhone(c *gin.Context) {
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
func (h *httpHandler) GetUsers(c *gin.Context) {
	normalUser, err := h.userUsercase.GetUsers()
	if err != nil {
		response(c, http.StatusBadRequest, "Bad request")
	}
	c.JSON(http.StatusOK, normalUser)
}

// GetUsers implements UserHandler.
// func (h *httpHandler) GetUserByPhone(c *gin.Context) {
// 	phone := c.Param("phone")

// 	user, err := h.userUsercase.GetUserByPhone(phone)
// 	if err != nil {
// 		response(c, http.StatusBadRequest, "Bad request")
// 	}

// 	c.JSON(http.StatusOK, user)
// }

// GetUser implements UserHandler.
func (h *httpHandler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
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

// GetNormalUser implements Handler.
func (h *httpHandler) GetNormalUser(c *gin.Context) {
	normalUserID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
	}
	normalUser, err := h.userUsercase.GetNormalUser(uint(normalUserID))
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	}

	c.JSON(http.StatusOK, gin.H{
		"normalUser": normalUser,
	})
}

// CancelCompatition implements Handler.
func (h *httpHandler) CancelCompatition(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.userUsercase.CancelCompatition(uint(compatitionID))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "hello world"})
}

// FinishCompatition implements Handler.
func (h *httpHandler) FinishCompatition(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.userUsercase.FinishCompatition(uint(compatitionID))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "finish com"})
}

// OpenCompatition implements Handler.
func (h *httpHandler) OpenCompatition(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.userUsercase.OpenApplicationCompatition(uint(compatitionID))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "open application for compatition success"})
}

// StartCompatition implements Handler.
func (h *httpHandler) StartCompatition(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.userUsercase.StartCompatition(uint(compatitionID))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update compatition status to \"Started\" success"})
}

// UpdateCompatitionStatus implements Handler.
func (h *httpHandler) UpdateCompatitionStatus(c *gin.Context) {
	compatitionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	var status struct {
		Status string `json:"status"`
	}

	err = c.BindJSON(&status)
	if err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	err = h.userUsercase.UpdateCompatitionStatus(uint(compatitionID), status.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update compatition status success"})
}
