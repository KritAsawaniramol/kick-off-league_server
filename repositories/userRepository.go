package repositories

import "kickoff-league.com/entities"

type Userrepository interface {
	GetUserByEmail(email string) (*entities.Users, error)
	GetUsers() ([]entities.Users, error)
	GetUserByID(in uint) (*entities.Users, error)

	GetNormalUser(*entities.NormalUsers) (*entities.NormalUsers, error)
	GetNormalUserWithAddressByUserID(in uint) (*entities.NormalUsers, error)
	GetNormalUserByUsername(username string) (*entities.NormalUsers, error)
	GetNumberOfTeamsMember(in uint) int64
	GetNormalUsers(in *entities.NormalUsers) ([]entities.NormalUsers, error)

	GetTeamWithMemberAndRequestSendByID(in uint) (*entities.Teams, error)
	GetTeamWithAllAssociationsByID(in *entities.Teams) (*entities.Teams, error)
	GetTeams(in *entities.Teams, orderString string, decs bool, limit int, offset int) ([]entities.Teams, error)
	GetTeamWithMemberAndCompatitionByID(in uint) (*entities.Teams, error)
	GetTeam(in uint) (*entities.Teams, error)
	GetTeamMembersByTeamID(in uint, orderString string, decs bool, limit int, offset int) ([]entities.TeamsMembers, error)
	GetOrganizerWithAddressByUserID(in uint) (*entities.Organizers, error)
	GetOrganizer(*entities.Organizers) (*entities.Organizers, error)
	GetAddMemberRequestByID(in *entities.AddMemberRequests) ([]entities.AddMemberRequests, error)
	GetCompatitions(in *entities.Compatitions, orderString string, decs bool, limit int, offset int) ([]entities.Compatitions, error)
	GetCompatition(in *entities.Compatitions) (*entities.Compatitions, error)

	InsertCompatition(in *entities.Compatitions) error
	InsertMatchs(in []entities.Matches) error
	InsertTeamsMembers(in *entities.TeamsMembers) error
	InsertTeam(in *entities.Teams) error
	InsertUserWihtNormalUserAndAddress(in_normalUser *entities.NormalUsers, in_user *entities.Users) error
	InsertUserWihtOrganizerAndAddress(in_organizer *entities.Organizers, in_user *entities.Users) error
	InsertAddMemberRequest(in *entities.AddMemberRequests) error

	UpdateSelectedFields(model interface{}, fieldname string, value interface{}) error
	UpdateAddMemberRequestStatusByID(inID uint, inStatus string) error
	UpdateNormalUser(inNormalUser *entities.NormalUsers) error
	UpdateUser(inUser *entities.Users) error
	UpdateNormalUserPhone(in_userID uint, newPhone string) error
	UpdateAddMemberRequestStatusAndSoftDelete(inReq *entities.AddMemberRequests, inStatus string) error
}
