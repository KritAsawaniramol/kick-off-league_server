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
	RemoveTeamFormCompatition(teamID uint, compatitionID uint) error
	CancelCompatition(id uint) error
	AddJoinCode(compatitionID uint, n int) error
	GetMatch(id uint) (*model.Match, error)

	GetOrganizer(id uint) (*model.GetOrganizer, error)
	GetOrganizers() ([]model.OrganizersInfo, error)

	CreateTeam(in *model.CreateTeam) error
	SendAddMemberRequest(in *model.AddMemberRequest, userID uint) error
	AcceptAddMemberRequest(inReqID uint, userID uint) error
	IgnoreAddMemberRequest(inReqID uint, userID uint) error
	UpdateOrganizer(orgID uint, in *model.UpdateOrganizer) error
	UpdateNormalUser(inUpdateModel *model.UpdateNormalUser, inNormalUserID uint) error // OrganizerRegister(in *model.)
	// GetUserByPhone(in string) (model.NormalUser, error)
	UpdateNormalUserPhone(inUserID uint, newPhone string) error
	UpdateImageCover(userID uint, newImagePath string) error
	UpdateImageProfile(userID uint, newImagePath string) error
	UpdateImageBanner(compatitionID uint, newImagePath string) error
	UpdateTeamImageCover(teamID uint, newImagePath string) error
	UpdateTeamImageProfile(teamID uint, newImagePath string) error

	RemoveImageBanner(compatitionID uint) error
	RemoveImageProfile(userID uint) error
	RemoveImageCover(userID uint) error
	RemoveTeamImageProfile(teamID uint) error
	RemoveTeamImageCover(teamID uint) error
}
