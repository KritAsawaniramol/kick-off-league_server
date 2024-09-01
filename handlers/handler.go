package handlers

import (
	"github.com/gin-gonic/gin"
	"kickoff-league.com/usecases/addMemberUsecase"
	"kickoff-league.com/usecases/authUsecase"
	"kickoff-league.com/usecases/competitionUsecase"
	"kickoff-league.com/usecases/matchUsecase"
	"kickoff-league.com/usecases/normalUserUsecase"
	"kickoff-league.com/usecases/organizerUsecase"
	"kickoff-league.com/usecases/teamUsecase"
	"kickoff-league.com/usecases/userUsecase"
)

type Handler interface {

	//User
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	DeleteImageProfile(c *gin.Context)
	DeleteImageCover(c *gin.Context)
	UpdateImageCover(c *gin.Context)
	UpdateImageProfile(c *gin.Context)

	//Auth
	RegisterOrganizer(c *gin.Context)
	RegisterNormaluser(c *gin.Context)
	LoginUser(c *gin.Context)
	LogoutUser(c *gin.Context)

	//Team
	GetTeams(c *gin.Context)
	GetTeam(c *gin.Context)
	GetTeamByOwnerID(c *gin.Context)
	CreateTeam(c *gin.Context)
	RemoveTeamMember(c *gin.Context)
	RemoveCompatitionTeam(c *gin.Context)
	DeleteTeamImageCover(c *gin.Context)
	DeleteTeamImageProfile(c *gin.Context)
	UpdateTeamImageCover(c *gin.Context)
	UpdateTeamImageProfile(c *gin.Context)

	//NormalUser
	GetNormalUsers(c *gin.Context)
	GetNormalUser(c *gin.Context)
	UpdateNormalUser(c *gin.Context)

	//Competition
	DeleteImageBanner(c *gin.Context)
	UpdateImageBanner(c *gin.Context)
	CreateCompetition(c *gin.Context)
	GetCompetition(c *gin.Context)
	GetCompetitions(c *gin.Context)
	JoinCompetition(c *gin.Context)
	UpdateCompetition(c *gin.Context)
	AddJoinCode(c *gin.Context)
	StartCompetition(c *gin.Context)
	OpenApplicationCompetition(c *gin.Context)
	FinishCompetition(c *gin.Context)
	CancelCompetition(c *gin.Context)

	//Match
	UpdateMatch(c *gin.Context)
	GetMatch(c *gin.Context)
	GetNextMatch(c *gin.Context)

	// Organizer
	GetOrganizers(c *gin.Context)
	GetOrganizer(c *gin.Context)
	UpdateOrganizer(c *gin.Context)

	// Add member
	GetMyPenddingAddMemberRequest(c *gin.Context)
	SendAddMemberRequest(c *gin.Context)
	AcceptAddMemberRequest(c *gin.Context)
	IgnoreAddMemberRequest(c *gin.Context)

}

func NewhttpHandler(
	userUsercase userUsecase.UserUsecase,
	organizerUsecase organizerUsecase.OrganizerUsecase,
	authUsecase authUsecase.AuthUsecase,
	normalUserUsecase normalUserUsecase.NormalUserUsecase,
	teamUsecase teamUsecase.TeamUsecase,
	addMemberUsecase addMemberUsecase.AddMemberUsecase,
	competitionUsecase competitionUsecase.CompetitionUsecase,
	matchUsecase matchUsecase.MatchUsecase,
) Handler {
	return &httpHandler{
		userUsercase:       userUsercase,
		organizerUsecase:   organizerUsecase,
		authUsecase:        authUsecase,
		normalUserUsecase:  normalUserUsecase,
		teamUsecase:        teamUsecase,
		addMemberUsecase:   addMemberUsecase,
		competitionUsecase: competitionUsecase,
		matchUsecase:       matchUsecase,
	}
}

type httpHandler struct {
	userUsercase       userUsecase.UserUsecase
	normalUserUsecase  normalUserUsecase.NormalUserUsecase
	organizerUsecase   organizerUsecase.OrganizerUsecase
	authUsecase        authUsecase.AuthUsecase
	teamUsecase        teamUsecase.TeamUsecase
	addMemberUsecase   addMemberUsecase.AddMemberUsecase
	competitionUsecase competitionUsecase.CompetitionUsecase
	matchUsecase       matchUsecase.MatchUsecase
}
