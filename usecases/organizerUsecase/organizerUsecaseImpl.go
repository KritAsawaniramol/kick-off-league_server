package organizerUsecase

import (
	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
)

type organizerUsecaseImpl struct {
	repository repositories.Repository
}

func NewOrganizerUsecaseImpl(
	repository repositories.Repository,
) OrganizerUsecase {
	return &organizerUsecaseImpl{
		repository: repository,
	}
}

// GetOrganizers implements OrganizerUsecase.
func (o *organizerUsecaseImpl) GetOrganizers() ([]model.OrganizersInfo, error) {
	org, err := o.repository.GetOrganizers()
	if err != nil {
		if err.Error() != "record not found" {
			return nil, err
		} else {
			return nil, nil
		}
	}

	orgModel := []model.OrganizersInfo{}
	for _, v := range org {
		orgModel = append(orgModel, model.OrganizersInfo{
			ID:          v.ID,
			Name:        v.Name,
			Phone:       v.Phone,
			Description: v.Description,
			Address: model.Address{
				HouseNumber: v.Addresses.HouseNumber,
				Village:     v.Addresses.Village,
				Subdistrict: v.Addresses.Subdistrict,
				District:    v.Addresses.District,
				PostalCode:  v.Addresses.PostalCode,
				Country:     v.Addresses.Country,
			},
			ImageProfilePath: v.Users.ImageProfilePath,
			ImageCoverPath:   v.Users.ImageCoverPath,
		})
	}

	return orgModel, err
}

// GetOrganizer implements OrganizerUsecase.
func (o *organizerUsecaseImpl) GetOrganizer(id uint) (*model.GetOrganizer, error) {
	getOrganizer := &entities.Organizers{}
	getOrganizer.ID = id
	org, err := o.repository.GetOrganizer(getOrganizer)
	if err != nil {
		if err.Error() != "record not found" {
			return nil, err
		} else {
			return nil, nil
		}
	}
	getCompatitions := []model.GetCompatitions{}
	for _, v := range org.Competitions {
		getCompatitions = append(getCompatitions, model.GetCompatitions{
			ID:     v.ID,
			Name:   v.Name,
			Sport:  v.Sport,
			Format: v.Format,
			Address: model.Address{
				HouseNumber: v.HouseNumber,
				Village:     v.Village,
				Subdistrict: v.Subdistrict,
				District:    v.District,
				PostalCode:  v.PostalCode,
				Country:     v.Country,
			},
			Status:          v.Status,
			ApplicationType: v.ApplicationType,
			Sex:             v.Sex,
			StartDate:       v.StartDate,
			EndDate:         v.EndDate,
			OrganizerID:     org.ID,
			OrganizerName:   org.Name,
			AgeOver:         v.AgeOver,
			AgeUnder:        v.AgeUnder,
			ImageBanner:     v.ImageBannerPath,
		})
	}
	return &model.GetOrganizer{
		ID:          org.ID,
		Name:        org.Name,
		Phone:       org.Phone,
		Description: org.Description,
		Address: model.Address{
			HouseNumber: org.Addresses.HouseNumber,
			Village:     org.Addresses.Village,
			Subdistrict: org.Addresses.Subdistrict,
			District:    org.Addresses.District,
			PostalCode:  org.Addresses.PostalCode,
			Country:     org.Addresses.Country,
		},
		ImageProfilePath: org.Users.ImageProfilePath,
		ImageCoverPath:   org.Users.ImageCoverPath,
		Compatition:      getCompatitions,
	}, nil
}

// UpdateOrganizer implements OrganizerUsecase.
func (o *organizerUsecaseImpl) UpdateOrganizer(orgID uint, in *model.UpdateOrganizer) error {
	getOrg := &entities.Organizers{}
	getOrg.ID = orgID
	org, err := o.repository.GetOrganizer(getOrg)
	if err != nil {
		return err
	}
	org.Name = in.Name
	org.Phone = in.Phone
	org.Description = in.Description
	org.Addresses.HouseNumber = in.Address.HouseNumber
	org.Addresses.Village = in.Address.Village
	org.Addresses.Subdistrict = in.Address.Subdistrict
	org.Addresses.District = in.Address.District
	org.Addresses.PostalCode = in.Address.PostalCode
	org.Addresses.Country = in.Address.Country

	err = o.repository.UpdateOrganizer(orgID, org)
	if err != nil {
		return err
	}
	return nil
}

// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
