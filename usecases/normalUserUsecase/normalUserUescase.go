package normalUserUsecase

import (
	model "kickoff-league.com/models"
)

type NormalUserUsecase interface {
	GetNormalUserList() ([]model.NormalUserList, error)
	GetNormalUser(id uint) (*model.NormalUserProfile, error)
	UpdateNormalUser(inUpdateModel *model.UpdateNormalUser, inNormalUserID uint) error // OrganizerRegister(in *model.)
}
