package usecases

import (
	model "kickoff-league.com/models"
)

type UserUsecase interface {
	Login(in *model.LoginUser) (string, model.User, error)
	Logout() error
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
	GetCompatitions(in *model.GetCompatitionsReq) ([]model.Compatition, error)
	CreateCompatition(in *model.CreateCompatition) error
	CreateTeam(in *model.CreateTeam) error
	SendAddMemberRequest(in *model.AddMemberRequest, userID uint) error
	AcceptAddMemberRequest(inReqID uint, userID uint) error
	IgnoreAddMemberRequest(inReqID uint, userID uint) error
	UpdateNormalUser(inUpdateModel *model.UpdateNormalUser, inNormalUserID uint) error // OrganizerRegister(in *model.)

	// GetUserByPhone(in string) (model.NormalUser, error)
	UpdateNormalUserPhone(inUserID uint, newPhone string) error
}
