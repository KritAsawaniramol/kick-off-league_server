package competitionUsecase

import (
	model "kickoff-league.com/models"
)

type CompetitionUsecase interface {
	GetCompetition(in uint) (*model.GetCompatition, error)
	GetCompetitions(in *model.GetCompatitionsReq) ([]model.GetCompatitions, error)
	CreateCompetition(in *model.CreateCompetition) error
	UpdateCompetition(id uint, orgID uint, in *model.UpdateCompatition) error
	FinishCompetition(id uint, orgID uint) error
	JoinCompetition(in *model.JoinCompetition, userID uint) error
	AddJoinCode(compatitionID uint, orgID uint, n int) error
	OpenApplicationCompetition(id uint, orgID uint) error
	StartCompetition(id uint, orgID uint) error
	CancelCompatition(id uint, orgID uint) error
	UpdateImageBanner(compatitionID uint, orgID uint, newImagePath string) error
	RemoveImageBanner(compatitionID uint, orgID uint) error
}
