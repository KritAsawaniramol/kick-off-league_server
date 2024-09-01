package organizerUsecase

import (
	model "kickoff-league.com/models"
)

type OrganizerUsecase interface {
	GetOrganizer(id uint) (*model.GetOrganizer, error)
	GetOrganizers() ([]model.OrganizersInfo, error)
	UpdateOrganizer(orgID uint, in *model.UpdateOrganizer) error
}
