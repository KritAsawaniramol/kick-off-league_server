package handlers

import (
	"github.com/gin-gonic/gin"
	"kickoff-league.com/usecases"
)

type Handler interface {

	//auth
	RegisterOrganizer(c *gin.Context)
	RegisterNormaluser(c *gin.Context)
	LoginUser(c *gin.Context)
	LogoutUser(c *gin.Context)

	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	GetTeams(c *gin.Context)
	GetTeam(c *gin.Context)
	GetTeamByOwnerID(c *gin.Context)
	GetMyPenddingAddMemberRequest(c *gin.Context)
	GetNormalUsers(c *gin.Context)
	GetNormalUser(c *gin.Context)
	// uploadImage(c *gin.Context)
	// UpdateNormalUserPhone(c *gin.Context)
	DeleteImageProfile(c *gin.Context)
	UpdateImageCover(c *gin.Context)
	UpdateImageProfile(c *gin.Context)
	UploadImage(c *gin.Context)
	CreateTeam(c *gin.Context)

	CreateCompatition(c *gin.Context)
	GetCompatition(c *gin.Context)
	GetCompatitions(c *gin.Context)
	JoinCompatition(c *gin.Context)
	UpdateCompatition(c *gin.Context)
	UpdateMatch(c *gin.Context)
	GetMatch(c *gin.Context)

	UpdateNormalUser(c *gin.Context)
	SendAddMemberRequest(c *gin.Context)
	AcceptAddMemberRequest(c *gin.Context)
	IgnoreAddMemberRequest(c *gin.Context)
	GetNextMatch(c *gin.Context)
	GetMatchResult(c *gin.Context)

	AddJoinCode(c *gin.Context)
	StartCompatition(c *gin.Context)
	OpenCompatition(c *gin.Context)
	FinishCompatition(c *gin.Context)
	CancelCompatition(c *gin.Context)
	RemoveTeamMember(c *gin.Context)

	// GetUserByPhone(c *gin.Context)
}

func NewhttpHandler(userUsercase usecases.UserUsecase) Handler {
	return &httpHandler{
		userUsercase: userUsercase,
	}
}

type httpHandler struct {
	userUsercase usecases.UserUsecase
}
