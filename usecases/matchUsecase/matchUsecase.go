package matchUsecase

import (
	model "kickoff-league.com/models"
)

type MatchUsecase interface {
	UpdateMatch(id uint, orgID uint,updateMatch *model.UpdateMatch) error
	GetMatch(id uint) (*model.Match, error)
	GetNextMatch(id uint) ([]model.NextMatch, error)
}
