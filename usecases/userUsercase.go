package usecases

import (
	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
)

type UserUsecase interface {
	Login(in *model.LoginUser) (string, model.User, error)
	Register(in *model.RegisterUser) error
	GetUsers() ([]entities.User, error)
	GetUser(in uint) (model.User, error)
	CreateTeam(in *model.CreaetTeam) error
	UpdateNormalUser(inUpdateModel *model.UpdateNormalUser, inUserID uint) error
	// OrganizerRegister(in *model.)

	// GetUserByPhone(in string) (model.NormalUser, error)
	UpdateNormalUserPhone(inUserID uint, newPhone string) error
}
