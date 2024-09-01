package authUsecase

import (
	model "kickoff-league.com/models"
)

type AuthUsecase interface {
	Login(in *model.LoginUser) (string, model.LoginResponse, error)
	RegisterOrganizer(in *model.RegisterOrganizer) error
	RegisterNormaluser(in *model.RegisterNormaluser) error
}
