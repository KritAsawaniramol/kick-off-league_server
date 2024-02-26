package usecases

import (
	model "kickoff-league.com/models"
)

type UserUsecase interface {
	Login(in *model.LoginUser) (string, model.User, error)
	RegisterNormaluser(in *model.RegisterNormaluser) error
	RegisterOrganizer(in *model.RegisterUser) error

	GetUsers() ([]model.User, error)
	GetUser(in uint) (model.User, error)
	CreateTeam(in *model.CreaetTeam) error
	SendAddMemberRequest(in *model.AddMemberRequest, userID uint) error
	AcceptAddMemberRequest(inReqID uint, userID uint) error
	IgnoreAddMemberRequest(inReqID uint, userID uint) error

	UpdateNormalUser(inUpdateModel *model.UpdateNormalUser, inNormalUserID uint) error

	// OrganizerRegister(in *model.)

	// GetUserByPhone(in string) (model.NormalUser, error)
	UpdateNormalUserPhone(inUserID uint, newPhone string) error
}
