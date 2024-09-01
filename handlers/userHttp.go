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
)

// GetUser implements UserHandler.
func (h *httpHandler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("error: GetUser: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "error: get user failed"})
		return
	}
	user, err := h.userUsercase.GetUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}


// GetUserByPhone implements UserHandler.
func (h *httpHandler) GetUsers(c *gin.Context) {
	normalUser, err := h.userUsercase.GetUsers()
	if err != nil {
		response(c, http.StatusBadRequest, "Bad request")
		return
	}
	c.JSON(http.StatusOK, normalUser)
}


func (h *httpHandler) UpdateImageProfile(c *gin.Context) {
	
	in, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "update image profile failed"})
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "update image profile failed"})
		return
	}

	if err := h.userUsercase.UpdateImageProfile(c.GetUint("user_id"), imagePath); err != nil {
		log.Errorf(err.Error())
		if err := os.Remove(imagePath); err != nil {
			fmt.Printf("Error removing file: %s\n", err)
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update image profile success"})
}

// UpdateImageCover implements UserHandler.
func (h *httpHandler) UpdateImageCover(c *gin.Context) {
	userID := c.GetUint("user_id")

	in, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "update cover image failed"})
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "update cover image failed"})
		return
	}

	if err := h.userUsercase.UpdateImageCover(userID, imagePath); err != nil {
		if err := os.Remove(imagePath); err != nil {
			log.Printf("Error removing file: %s\n", err)
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "update cover image success"})
}

// DeleteImageProfile implements UserHandler.
func (h *httpHandler) DeleteImageProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.userUsercase.RemoveImageProfile(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "remove image profile successs"})
}


// DeleteImageCover implements Handler.
func (h *httpHandler) DeleteImageCover(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.userUsercase.RemoveImageCover(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "remove cover image successs"})
}


// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================



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
