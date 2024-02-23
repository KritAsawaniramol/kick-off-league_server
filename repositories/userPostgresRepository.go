package repositories

import (
	"github.com/labstack/gommon/log"

	"gorm.io/gorm"
	"kickoff-league.com/entities"
)

type userPostgresRepository struct {
	db *gorm.DB
}

// UpdateNormalUser implements Userrepository.
func (u *userPostgresRepository) UpdateNormalUser(inNormalUser *entities.NormalUser, inUserID uint) error {
	log.Print(*inNormalUser)
	result := u.db.Model(&entities.NormalUser{}).Where("user_id = ?", inUserID).Updates(inNormalUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetOrganizerByUserID implements Userrepository.
func (u *userPostgresRepository) GetOrganizerByUserID(in uint) (*entities.Organizer, error) {
	organizer := &entities.Organizer{}
	err := u.db.Model(&entities.Organizer{}).Where("user_id = ?", in).First(&organizer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &entities.Organizer{}, err
	}
	return organizer, nil
}

// GetNormalUserByUserID implements Userrepository.
func (u *userPostgresRepository) GetNormalUserByUserID(in uint) (*entities.NormalUser, error) {
	normalUser := &entities.NormalUser{}
	err := u.db.Model(&entities.NormalUser{}).Where("user_id = ?", in).First(&normalUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &entities.NormalUser{}, err
	}
	return normalUser, nil
}

// InsertTeam implements Userrepository.
func (u *userPostgresRepository) InsertTeam(in *entities.Team) error {
	result := u.db.Create(in)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

// InsertUserWihtOrganizerAndAddress implements Userrepository.
func (u *userPostgresRepository) InsertUserWihtOrganizerAndAddress(in_organizer *entities.Organizer, in_user *entities.User) error {
	err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&in_user).Error; err != nil {
			return err
		}

		address := &entities.Address{}

		if err := tx.Create(address).Error; err != nil {
			return err
		}

		in_organizer.UserID = in_user.ID
		in_organizer.AddressID = address.ID

		if err := tx.Create(&in_organizer).Error; err != nil {
			return err
		}
		return nil
	},
	)
	return err
}

// InsertUserWihtNormalUserAndAddress implements Userrepository.
func (u *userPostgresRepository) InsertUserWihtNormalUserAndAddress(in_normalUser *entities.NormalUser, in_user *entities.User) error {
	err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&in_user).Error; err != nil {
			return err
		}

		address := &entities.Address{}

		if err := tx.Create(address).Error; err != nil {
			return err
		}

		in_normalUser.UserID = in_user.ID
		in_normalUser.AddressID = address.ID

		if err := tx.Create(&in_normalUser).Error; err != nil {
			return err
		}
		return nil
	},
	)
	return err
}

// GetNormalUserByPhone implements Userrepository.
func (u *userPostgresRepository) GetNormalUserByPhone(in string) (*entities.NormalUser, error) {

	normalUser := &entities.NormalUser{}
	err := u.db.Model(&entities.NormalUser{}).Where("phone = ?", in).First(&normalUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &entities.NormalUser{}, err
	}
	return normalUser, nil
}

// UpdateNormalUserPhone implements Userrepository.
func (u *userPostgresRepository) UpdateNormalUserPhone(in_userID uint, newPhone string) error {
	// Check if new Phone already exists
	existingUser := &entities.NormalUser{}
	err := u.db.Model(&entities.NormalUser{}).Where("phone = ?", newPhone).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf(err.Error())
		return err
	}

	// If email doesn't exist, update

	if existingUser.ID == 0 {
		if err := u.db.Model(&entities.NormalUser{}).Where("user_id = ?", in_userID).Update("phone", newPhone).Error; err != nil {
			return err
		}
	}

	return nil
	// // Check if new email already exists
	// existingUser := User{}
	// err := db.Model(&User{}).Where("email = ?", newEmail).First(&existingUser).Error
	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	return err
	// }

	// // If email doesn't exist, update
	// if existingUser.ID == 0 {
	// 	result := db.Model(&User{ID: userID}).Update("email", newEmail)
	// 	return result.Error
	// }

	// // Email already exists, handle accordingly (e.g., error or ignore)
	// return nil
}

// InsertOrganizer implements Userrepository.
func (u *userPostgresRepository) InsertOrganizer(in *entities.Organizer) error {
	result := u.db.Create(in)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// InsertNormalUser implements Userrepository.
func (u *userPostgresRepository) InsertNormalUser(in *entities.NormalUser) error {
	result := u.db.Create(in)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetNormalUserWithUserByUserID implements Userrepository.
func (u *userPostgresRepository) GetNormalUserWithUserByUserID(in uint) (*entities.NormalUser, error) {
	var normalUser entities.NormalUser
	result := u.db.Preload("User").First(&normalUser, in)
	if result.Error != nil {
		return nil, result.Error
	}
	return &normalUser, nil
}

func NewUserPostgresRepository(db *gorm.DB) Userrepository {
	return &userPostgresRepository{db: db}
}

// GetUsers implements Userrepository.
func (u *userPostgresRepository) GetUsers() ([]entities.User, error) {
	users := []entities.User{}
	result := u.db.Find(&users).Order("id DESC")
	if result.Error != nil {
		return []entities.User{}, result.Error
	}
	return users, nil
}

// GetUserByEmail implements Userrepository.
func (u *userPostgresRepository) GetUserByEmail(in string) (*entities.User, error) {
	selectedUser := &entities.User{}
	result := u.db.Where("email = ?", in).First(selectedUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return selectedUser, nil
}

// GetUserByID implements Userrepository.
func (u *userPostgresRepository) GetUserByID(in uint) (*entities.User, error) {
	selectedUser := &entities.User{}
	result := u.db.Where("id = ?", in).First(selectedUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return selectedUser, nil
}

func (r *userPostgresRepository) InsertUserData(in *entities.User) error {

	result := r.db.Create(in)

	if result.Error != nil {
		log.Errorf("InsertUserData: %v", result.Error)
		return result.Error
	}

	log.Debugf("InsertUserData: %v", result.RowsAffected)
	return nil
}
