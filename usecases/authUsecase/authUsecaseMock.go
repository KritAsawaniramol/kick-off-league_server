package authUsecase

import (
	"github.com/stretchr/testify/mock"
	model "kickoff-league.com/models"
)

type AuthUsecaseMock struct {
	mock.Mock
}

func NewAuthUsecase() AuthUsecase {
	return &AuthUsecaseMock{}
}

// Login implements AuthUsecase.
func (a *AuthUsecaseMock) Login(in *model.LoginUser) (string, model.LoginResponse, error) {
	args := a.Called(in)
	return args.String(0), args.Get(1).(model.LoginResponse), args.Error(2)
}

// RegisterNormaluser implements AuthUsecase.
func (a *AuthUsecaseMock) RegisterNormaluser(in *model.RegisterNormaluser) error {
	args := a.Called(in)
	return args.Error(0)
}

// RegisterOrganizer implements AuthUsecase.
func (a *AuthUsecaseMock) RegisterOrganizer(in *model.RegisterOrganizer) error {
	args := a.Called(in)
	return args.Error(0)
}
