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
	"kickoff-league.com/usecases"
	"kickoff-league.com/util"
)

type userHttpHandler struct {
	userUsercase usecases.UserUsecase
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
func (h *userHttpHandler) UpdateImageCover(c *gin.Context) {
	normalUserID := c.GetUint("normal_user_id")

	in, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// extract image extension from original file filename
	isImage, fileExt := isImage(in)
	if !isImage {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File is not an image(png, jpeg)"})
		return
	}

	imagePath := createImagePath(fileExt, "./images/cover")

	if err := c.SaveUploadedFile(in, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	reqBody := new(model.UpdateNormalUser)
	reqBody.ImageCoverPath = imagePath

	if err := h.userUsercase.UpdateNormalUser(reqBody, normalUserID); err != nil {
		log.Errorf(err.Error())
		if err := os.Remove(imagePath); err != nil {
			fmt.Printf("Error removing file: %s\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "UpdateImageProfile success"})
}

func (h *userHttpHandler) UpdateImageProfile(c *gin.Context) {
	normalUserID := c.GetUint("normal_user_id")

	in, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// extract image extension from original file filename
	isImage, fileExt := isImage(in)
	if !isImage {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File is not an image(png, jpeg)"})
		return
	}

	imagePath := createImagePath(fileExt, "./images/profile")

	if err := c.SaveUploadedFile(in, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	reqBody := new(model.UpdateNormalUser)
	reqBody.ImageProfilePath = imagePath

	if err := h.userUsercase.UpdateNormalUser(reqBody, normalUserID); err != nil {
		log.Errorf(err.Error())
		if err := os.Remove(imagePath); err != nil {
			fmt.Printf("Error removing file: %s\n", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "UpdateImageProfile success"})
}

// UploadFile implements UserHandler.
func (*userHttpHandler) UploadImage(c *gin.Context) {
	in, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "File upload complete!"})
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

// // uploadImage implements UserHandler.
// func (*userHttpHandler) uploadImage(c *gin.Context) {
// 	file, err := c.FormFile("image") {

// 	}
// }

// CreateTeam implements UserHandler.
func (h *userHttpHandler) CreateTeam(c *gin.Context) {
	reqBody := new(model.CreaetTeam)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}
	reqBody.OwnerID = c.GetUint("user_id")
	if err := h.userUsercase.CreateTeam(reqBody); err != nil {
		log.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "CreateTeam success"})
}

// CreateCompatition implements UserHandler.
func (h *userHttpHandler) CreateCompatition(c *gin.Context) {

	organizerID := c.GetUint("organizer_id")

	reqBody := new(model.CreateCompatition)
	if err := c.BindJSON(reqBody); err != nil {
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	reqBody.OrganizerID = organizerID

	if err := h.userUsercase.CreateCompatition(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Create compatition success"})
}

// GetMyPenddingAddMemberRequest implements UserHandler.
func (h *userHttpHandler) GetMyPenddingAddMemberRequest(c *gin.Context) {
	addMemberRequests, err := h.userUsercase.GetMyPenddingAddMemberRequest(c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
	}
	c.JSON(http.StatusOK, gin.H{"add_member_requests": addMemberRequests})
}

// GetTeam implements UserHandler.
func (h *userHttpHandler) GetTeam(c *gin.Context) {
	TeamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
	}
	teams, err := h.userUsercase.GetTeamWithMemberAndCompatitionByID(uint(TeamID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

// GetTeamList implements UserHandler.
func (h *userHttpHandler) GetTeams(c *gin.Context) {
	reqBody := new(model.GetTeamsReq)
	if err := c.BindJSON(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		response(c, http.StatusBadRequest, "Bad request")
		return
	}

	util.PrintObjInJson(reqBody)
	teams, err := h.userUsercase.GetTeams(reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"teams": teams})
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
	reqBody := new(model.RegisterOrganizer)
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
