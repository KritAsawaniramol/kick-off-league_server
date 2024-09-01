package repositories

import "kickoff-league.com/entities"

type Repository interface {
	GetUserByEmail(email string) (*entities.Users, error)
	GetUsers() ([]entities.Users, error)
	GetUserByID(in uint) (*entities.Users, error)

	GetNormalUserDetails(*entities.NormalUsers) (*entities.NormalUsers, error)
	GetNormalUserWithAddressByUserID(in uint) (*entities.NormalUsers, error)
	GetNormalUserByUsername(username string) (*entities.NormalUsers, error)
	GetNumberOfTeamsMember(in uint) int64
	GetNormalUsers(in *entities.NormalUsers) ([]entities.NormalUsers, error)

	GetNormalUserCompetitions(in *entities.NormalUsersCompetitions) ([]entities.NormalUsersCompetitions, error)

	GetTeamsWithCompetitionAndMatch(in *entities.Teams) (*entities.Teams, error)
	GetTeamWithMemberAndRequestSentByID(in uint) (*entities.Teams, error)
	GetTeamWithAllAssociationsByID(in *entities.Teams) (*entities.Teams, error)
	GetTeams(in *entities.Teams, orderString string, decs bool, limit int, offset int) ([]entities.Teams, error)
	GetTeamWithMemberAndCompetitionByID(in uint) (*entities.Teams, error)
	GetTeam(in *entities.Teams) (*entities.Teams, error)
	// GetTeamMembersByTeamID(in uint, orderString string, decs bool, limit int, offset int) ([]entities.TeamsMembers, error)

	GetOrganizerWithAddressByUserID(in uint) (*entities.Organizers, error)
	GetOrganizer(*entities.Organizers) (*entities.Organizers, error)
	GetOrganizers() ([]entities.Organizers, error)
	GetAddMemberRequestByID(in *entities.AddMemberRequests) ([]entities.AddMemberRequests, error)
	GetCompetitions(in *entities.Competitions, orderString string, decs bool, limit int, offset int) ([]entities.Competitions, error)
	GetCompetition(in *entities.Competitions) (*entities.Competitions, error)
	GetCompetitionDetails(in *entities.Competitions) (*entities.Competitions, error)
	AppendJoinCodeToCompetition(id uint, joinCodes []entities.JoinCode) error
	UpdateJoinCode(id uint, in *entities.JoinCode) error

	// ============================
	ClearGoalRecordsOfMatch(matchID uint) error
	ReplaceGoalRecordsOfMatch(matchID uint, goalRecords []entities.GoalRecords) error
	AppendGoalRecordsToMatch(id uint, goalRecords []entities.GoalRecords) error
	UpdateMatch(uint, *entities.Matchs) error

	DeleteTeamMember(nomalUserID uint, teamID uint) error
	// DeleteNormalUserCompetitionByTeamIDAndCompetitionID(competitionID uint, teamID uint) error
	// DeleteCompetitionsTeam(competitionID uint, teamID uint) error
	DeleteTeamFromCompetition(competitionID uint, teamID uint) error

	GetMatchDetail(in *entities.Matchs) (*entities.Matchs, error)
	GetMatchs(in *entities.Matchs) ([]entities.Matchs, error)
	// AppendMatchToCompetition(competition *entities.Competitions, matchs []entities.Matchs) error

	UpdateOrganizer(id uint, in *entities.Organizers) error
	UpdateCompetitionsTeams(in *entities.CompetitionsTeams) error

	AddCompetitionToTeamAndNormalUsers(inComTeam *entities.CompetitionsTeams, normalUserID []uint) error
	// InsertNormalUserCompetition(in *entities.NormalUsersCompetitions) error
	// InsertCompetitionsTeams(in *entities.CompetitionsTeams) error


	// UpdateCompetition(id uint, orgID uint, in *entities.Competitions) error
	UpdateCompetition(query *entities.Competitions, in *entities.Competitions) error

	StartCompetitionAndAppendMatchToCompetition(query *entities.Competitions, in *entities.Competitions, matchs []entities.Matchs) error

	InsertCompetition(in *entities.Competitions) error
	UpdateAddmemberRequestAndInsertTeamsMembers(in *entities.TeamsMembers,  memberReqID uint, memberReqStatus string) error
	InsertTeam(in *entities.Teams) error
	InsertUserWihtNormalUserAndAddress(in_normalUser *entities.NormalUsers, in_user *entities.Users) error
	InsertUserWihtOrganizerAndAddress(in_organizer *entities.Organizers, in_user *entities.Users) error
	InsertAddMemberRequest(in *entities.AddMemberRequests) error

	UpdateSelectedFields(model interface{}, fieldname string, value interface{}) error
	UpdateAddMemberRequestStatusByID(inID uint, inStatus string) error
	UpdateNormalUser(inNormalUser *entities.NormalUsers) error
	UpdateAddMemberRequestStatusAndSoftDelete(inReq *entities.AddMemberRequests, inStatus string) error
}
