package userUsecase

import (
	"log"
	"os"
	"strings"

	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
)

type userUsecaseImpl struct {
	repository repositories.Repository
}

func NewUserUsecaseImpl(
	repository repositories.Repository,
) UserUsecase {
	return &userUsecaseImpl{
		repository: repository,
	}
}

// GetUser implements UserUsecase.
func (u *userUsecaseImpl) GetUser(in uint) (model.User, error) {
	// get user from email
	user, err := u.repository.GetUserByID(in)
	if err != nil {
		return model.User{}, err
	}

	userModel := model.User{
		ID:               user.ID,
		Email:            user.Email,
		Role:             user.Role,
		ImageProfilePath: user.ImageProfilePath,
		ImageCoverPath:   user.ImageCoverPath,
	}

	if user.Role == "normal" {
		normalUser, err := u.repository.GetNormalUserWithAddressByUserID(user.ID)
		if err != nil {
			return model.User{}, err
		}

		userModel.NormalUserInfo = &model.NormalUserInfo{
			ID:            normalUser.ID,
			FirstNameThai: normalUser.FirstNameThai,
			LastNameThai:  normalUser.LastNameThai,
			FirstNameEng:  normalUser.FirstNameEng,
			LastNameEng:   normalUser.LastNameEng,
			Username:      normalUser.Username,
			Born:          normalUser.Born,
			Phone:         normalUser.Phone,
			Height:        normalUser.Height,
			Weight:        normalUser.Weight,
			Sex:           normalUser.Sex,
			Position:      normalUser.Position,
			Nationality:   normalUser.Nationality,
			Description:   normalUser.Description,
			Address: model.Address{
				HouseNumber: normalUser.Addresses.HouseNumber,
				Village:     normalUser.Addresses.Village,
				Subdistrict: normalUser.Addresses.Subdistrict,
				District:    normalUser.Addresses.District,
				PostalCode:  normalUser.Addresses.PostalCode,
				Country:     normalUser.Addresses.Country,
			},
		}
		// userModel.OrganizersInfo = model.OrganizersInfo{}

	} else if user.Role == "organizer" {
		organizer, err := u.repository.GetOrganizerWithAddressByUserID(user.ID)
		if err != nil {
			return model.User{}, err
		}
		userModel.OrganizersInfo = &model.OrganizersInfo{
			ID:          organizer.ID,
			Name:        organizer.Name,
			Phone:       organizer.Phone,
			Description: organizer.Description,
			Address: model.Address{
				HouseNumber: organizer.Addresses.HouseNumber,
				Village:     organizer.Addresses.Village,
				Subdistrict: organizer.Addresses.Subdistrict,
				District:    organizer.Addresses.District,
				PostalCode:  organizer.Addresses.PostalCode,
				Country:     organizer.Addresses.Country,
			},
		}
	}
	return userModel, nil
}

// GetUsers implements UserUsecase.
func (u *userUsecaseImpl) GetUsers() ([]model.User, error) {
	users_entity, err := u.repository.GetUsers()
	if err != nil {
		return []model.User{}, err
	}
	users_model := []model.User{}
	for _, e := range users_entity {
		m := model.User{
			ID:               e.ID,
			Email:            e.Email,
			Role:             e.Role,
			ImageProfilePath: e.ImageProfilePath,
			ImageCoverPath:   e.ImageCoverPath,
		}
		users_model = append(users_model, m)
	}
	return users_model, nil
}

// UpdateImageProfile implements UserUsecase.
func (u *userUsecaseImpl) UpdateImageProfile(userID uint, newImagePath string) error {


	user, err := u.repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	if !strings.Contains(user.ImageProfilePath, "/default/") {
		if err := os.Remove(user.ImageProfilePath); err != nil {
			log.Printf("error: UpdateImageProfile: %s\n", err.Error())
		}
	}

	if err := u.repository.UpdateSelectedFields(user, "ImageProfilePath", &entities.Users{ImageProfilePath: newImagePath}); err != nil {
		return err
	}
	return nil
}

// UpdateImageCover implements UserUsecase.
func (u *userUsecaseImpl) UpdateImageCover(userID uint, newImagePath string) error {
	user, err := u.repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	if !strings.Contains(user.ImageCoverPath, "/default/") {
		if err := os.Remove(user.ImageCoverPath); err != nil {
			log.Printf("error: UpdateImageCover: %s\n", err.Error())
		}
	}

	if err := u.repository.UpdateSelectedFields(user, "ImageCoverPath", &entities.Users{ImageCoverPath: newImagePath}); err != nil {
		return err
	}
	return nil
}

// RemoveImageProfile implements UserUsecase.
func (u *userUsecaseImpl) RemoveImageProfile(userID uint) error {
	user, err := u.repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	if !strings.Contains(user.ImageProfilePath, "/default/") {
		if err := os.Remove(user.ImageProfilePath); err != nil {
			log.Printf("error: RemoveImageProfile: %s\n", err.Error())
		}
	}

	if err := u.repository.UpdateSelectedFields(user, "ImageProfilePath", &entities.Users{ImageProfilePath: "./images/default/defaultProfile.jpg"}); err != nil {
		return err
	}
	return nil
}

// RemoveImageCover implements UserUsecase.
func (u *userUsecaseImpl) RemoveImageCover(userID uint) error {
	user, err := u.repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	if !strings.Contains(user.ImageCoverPath, "/default/") {
		if err := os.Remove(user.ImageCoverPath); err != nil {
			log.Printf("error: RemoveImageProfile: %s\n", err.Error())
		}
	}

	if err := u.repository.UpdateSelectedFields(user, "ImageCoverPath", &entities.Users{ImageCoverPath: "./images/default/defaultCover.jpg"}); err != nil {
		return err
	}
	return nil
}



// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================




