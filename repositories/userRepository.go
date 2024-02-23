package repositories

import "kickoff-league.com/entities"

type Userrepository interface {
	InsertUserData(in *entities.User) error
	GetUserByEmail(in string) (*entities.User, error)
	GetUsers() ([]entities.User, error)
	GetUserByID(in uint) (*entities.User, error)

	InsertTeam(in *entities.Team) error
	InsertNormalUser(in *entities.NormalUser) error
	InsertOrganizer(in *entities.Organizer) error
	InsertUserWihtNormalUserAndAddress(in_normalUser *entities.NormalUser, in_user *entities.User) error
	InsertUserWihtOrganizerAndAddress(in_organizer *entities.Organizer, in_user *entities.User) error

	UpdateNormalUser(inNormalUser *entities.NormalUser, inUserID uint) error

	UpdateNormalUserPhone(in_userID uint, newPhone string) error
	// GetNormalUserByUserID(in uint) (*entities.NormalUser, error)
	GetNormalUserWithUserByUserID(in uint) (*entities.NormalUser, error)

	GetNormalUserByPhone(in string) (*entities.NormalUser, error)
	GetNormalUserByUserID(in uint) (*entities.NormalUser, error)
	GetOrganizerByUserID(in uint) (*entities.Organizer, error)

	// GetNormalUser
}
