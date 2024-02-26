package repositories

import "kickoff-league.com/entities"

type Userrepository interface {
	InsertUserData(in *entities.User) error
	GetUserByEmail(in string) (*entities.User, error)
	GetUsers() ([]entities.User, error)
	GetUserByID(in uint) (*entities.User, error)

	InsertTeam(in *entities.Teams) error
	InsertNormalUser(in *entities.NormalUser) error
	InsertOrganizer(in *entities.Organizer) error
	InsertUserWihtNormalUserAndAddress(in_normalUser *entities.NormalUser, in_user *entities.User) error
	InsertUserWihtOrganizerAndAddress(in_organizer *entities.Organizer, in_user *entities.User) error
	InsertAddMemberRequest(in *entities.AddMemberRequest) error
	UpdateAddMemberRequestStatusByID(inID uint, inStatus string) error
	UpdateNormalUser(inNormalUser *entities.NormalUser) error
	GetAddMemberRequestByID(in uint) (*entities.AddMemberRequest, error)
	InsertTeamsMember(in *entities.TeamsMember) error
	UpdateNormalUserPhone(in_userID uint, newPhone string) error
	// GetNormalUserByUserID(in uint) (*entities.NormalUser, error)
	GetNormalUserWithUserByUserID(in uint) (*entities.NormalUser, error)
	UpdateAddMemberRequestStatusAndSoftDelete(inReq *entities.AddMemberRequest, inStatus string) error

	GetTeamWithMemberAndRequestSendByID(in uint) (*entities.Teams, error)
	GetTeamByID(in uint) (*entities.Teams, error)
	GetNormalUserByUsername(in string) (*entities.NormalUser, error)
	GetNormalUserByPhone(in string) (*entities.NormalUser, error)
	GetNormalUserByUserID(in uint) (*entities.NormalUser, error)
	GetOrganizerByUserID(in uint) (*entities.Organizer, error)
	SoftDeleteAddMemberRequestByID(in uint) error
	// GetNormalUser
}
