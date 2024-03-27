package usecases

import (
	model "kickoff-league.com/models"
)

type UserUsecase interface {
	Login(in *model.LoginUser) (string, model.LoginResponse, error)
	RegisterNormaluser(in *model.RegisterNormaluser) error
	RegisterOrganizer(in *model.RegisterOrganizer) error

	GetUsers() ([]model.User, error)
	GetUser(in uint) (model.User, error)
	// GetNormalUser(in uint) (model.NormalUser, error)
	// GetTeam(in uint) (model.NormalUser, error)

	RemoveImageProfile(normalUserID uint) error
	GetMyPenddingAddMemberRequest(userID uint) ([]model.AddMemberRequest, error)
	GetTeamMembers(id uint) (*model.Team, error)
	GetTeamWithMemberAndCompatitionByID(id uint) (*model.Team, error)
	GetTeams(in *model.GetTeamsReq) ([]model.TeamList, error)
	GetTeamsByOwnerID(in uint) ([]model.TeamList, error)
	GetNormalUserList() ([]model.NormalUserList, error)
	GetNormalUser(id uint) (*model.NormalUserProfile, error)
	GetNextMatch(id uint) ([]model.NextMatch, error)
	UpdateMatch(id uint, updateMatch *model.UpdateMatch) error
	JoinCompatition(in *model.JoinCompatition) error
	CreateCompatition(in *model.CreateCompatition) error
	GetCompatition(in uint) (*model.GetCompatition, error)
	GetCompatitions(in *model.GetCompatitionsReq) ([]model.GetCompatitions, error)
	UpdateCompatition(id uint, in *model.UpdateCompatition) error
	OpenApplicationCompatition(id uint) error
	StartCompatition(id uint) error
	FinishCompatition(id uint) error
	RemoveNormalUserFormTeam(teamID uint, nomalUserID uint) error
	CancelCompatition(id uint) error
	UpdateCompatitionStatus(id uint, status string) error
	CreateJoinCode(compatitionID uint, n int) error

	CreateTeam(in *model.CreateTeam) error
	SendAddMemberRequest(in *model.AddMemberRequest, userID uint) error
	AcceptAddMemberRequest(inReqID uint, userID uint) error
	IgnoreAddMemberRequest(inReqID uint, userID uint) error
	UpdateNormalUser(inUpdateModel *model.UpdateNormalUser, inNormalUserID uint) error // OrganizerRegister(in *model.)
	UpdateUser(in *model.User) error
	// GetUserByPhone(in string) (model.NormalUser, error)
	UpdateNormalUserPhone(inUserID uint, newPhone string) error
	UpdateImageCover(userID uint, newImagePath string) error
	UpdateImageProfile(userID uint, newImagePath string) error
}
