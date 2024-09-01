package repositories

import (
	"errors"

	"github.com/labstack/gommon/log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"kickoff-league.com/entities"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewUserPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

// GetNormalUserByUsername implements Userrepository.
func (u *postgresRepository) GetNormalUserByUsername(username string) (*entities.NormalUsers, error) {
	selectedUser := &entities.NormalUsers{}
	if err := u.db.Where("username = ?", username).First(selectedUser).Error; err != nil {
		log.Printf("error: GetNormalUserByUsername: %s\n", err.Error())
		return nil, errors.New("error: normalUser not found")
	}
	return selectedUser, nil
}

// GetUserByEmail implements Userrepository.
func (u *postgresRepository) GetUserByEmail(email string) (*entities.Users, error) {
	selectedUser := &entities.Users{}
	if err := u.db.Where("email = ?", email).First(selectedUser).Error; err != nil {
		log.Printf("error: GetUserByEmail: %s\n", err.Error())
		return nil, errors.New("error: user not found")
	}
	return selectedUser, nil
}

// GetNormalUserDetails implements Userrepository.
func (u *postgresRepository) GetNormalUserDetails(in *entities.NormalUsers) (*entities.NormalUsers, error) {
	normalUser := new(entities.NormalUsers)
	if err := u.db.Where(&in).
		Preload(clause.Associations).
		Preload("Teams.Teams").
		Preload("Teams.Teams.TeamsMembers").
		Preload("Competitions.Competitions").
		Preload("Competitions.Competitions.Matchs").
		First(normalUser).Error; err != nil {
		log.Printf("error: GetNormalUser: %s\n", err.Error())
		return nil, errors.New("error: user not found")
	}
	return normalUser, nil
}

// InsertUserWihtNormalUserAndAddress implements Userrepository.
func (u *postgresRepository) InsertUserWihtNormalUserAndAddress(in_normalUser *entities.NormalUsers, in_user *entities.Users) error {
	err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&in_user).Error; err != nil {
			log.Printf("error: InsertUserWihtNormalUserAndAddress: %s\n", err.Error())
			return errors.New("error: insert user with normalUser and address failed")
		}

		address := &entities.Addresses{}

		if err := tx.Create(address).Error; err != nil {
			log.Printf("error: InsertUserWihtNormalUserAndAddress: %s\n", err.Error())
			return errors.New("error: insert user with normalUser and address failed")
		}

		in_normalUser.UsersID = in_user.ID
		in_normalUser.AddressesID = address.ID

		if err := tx.Create(&in_normalUser).Error; err != nil {
			log.Printf("error: InsertUserWihtNormalUserAndAddress: %s\n", err.Error())
			return errors.New("error: insert user with normalUser and address failed")
		}
		return nil
	},
	)
	if err != nil {
		log.Printf("error: InsertUserWihtNormalUserAndAddress: %s\n", err.Error())
		return errors.New("error: insert user with normalUser and address failed")
	}
	return nil
}

// GetNormalUserWithAddressByUserID implements Userrepository.
func (u *postgresRepository) GetNormalUserWithAddressByUserID(in uint) (*entities.NormalUsers, error) {
	normalUser := &entities.NormalUsers{}
	if err := u.db.Model(&entities.NormalUsers{}).Preload("Addresses").Where("users_id = ?", in).First(&normalUser).Error; err != nil {
		log.Printf("error: GetNormalUserWithAddressByUserID: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("error: record not found")
		}
		return nil, errors.New("error: get normaluser failed")
	}
	return normalUser, nil
}

// GetCompetition implements Repository.
func (u *postgresRepository) GetCompetition(in *entities.Competitions) (*entities.Competitions, error) {
	competition := &entities.Competitions{}
	if err := u.db.Where(&in).First(competition).Error; err != nil {
		log.Printf("error: GetCompetition: %s\n", err.Error())
		return nil, errors.New("error: competition not found")
	}
	return competition, nil
}

// UpdateSelectedFields implements Userrepository.
func (u *postgresRepository) UpdateSelectedFields(model interface{}, fieldname string, value interface{}) error {
	result := u.db.Where(model).Select(fieldname).Updates(value)
	if result.Error != nil {
		log.Printf("error: UpdateSelectedFields: %s\n", result.Error)
		return errors.New("error: update failed")
	}

	if result.RowsAffected == 0 {
		log.Printf("error: UpdateSelectedFields: no records found to update")
		return errors.New("error: no records found to update")
	}

	return nil
}

// InserrtCompetition implements Userrepository.
func (u *postgresRepository) InsertCompetition(in *entities.Competitions) error {
	if err := u.db.Create(in).Error; err != nil {
		log.Printf("error: InsertCompetition: %s\n", err.Error())
		return errors.New("error: insert competition failed")
	}
	return nil
}

// GetCompetitionByID implements Userrepository.
func (u *postgresRepository) GetCompetitionDetails(in *entities.Competitions) (*entities.Competitions, error) {
	competition := &entities.Competitions{}
	if err := u.db.Where(&in).
		Preload(clause.Associations).
		Preload("JoinCode").
		Preload("Organizers.Addresses").
		Preload("Organizers.Users").
		Preload("Matchs.GoalRecords").
		Preload("Teams.Teams.TeamsMembers").
		First(&competition).Error; err != nil {
		log.Printf("error: GetCompetitionDetails: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("error: record not found")
		}
		return nil, errors.New("error: get competition failed")
	}
	return competition, nil
}

// InsertUserWihtOrganizerAndAddress implements Userrepository.
func (u *postgresRepository) InsertUserWihtOrganizerAndAddress(in_organizer *entities.Organizers, in_user *entities.Users) error {
	if err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&in_user).Error; err != nil {
			log.Printf("error: InsertUserWihtOrganizerAndAddress: %s\n", err.Error())
			return errors.New("error: insert user with organizer and address failed")
		}
		address := &entities.Addresses{}
		if err := tx.Create(address).Error; err != nil {
			log.Printf("error: InsertUserWihtOrganizerAndAddress: %s\n", err.Error())
			return errors.New("error: insert user with organizer and address failed")
		}
		in_organizer.UsersID = in_user.ID
		in_organizer.AddressesID = address.ID
		if err := tx.Create(&in_organizer).Error; err != nil {
			log.Printf("error: InsertUserWihtOrganizerAndAddress: %s\n", err.Error())
			return errors.New("error: insert user with organizer and address failed")
		}
		return nil
	},
	); err != nil {
		log.Printf("error: InsertUserWihtOrganizerAndAddress: %s\n", err.Error())
		return errors.New("error: insert user with organizer and address failed")
	}
	return nil
}

// UpdateCompetition implements Userrepository.
func (u *postgresRepository) UpdateCompetition(query *entities.Competitions, in *entities.Competitions) error {
	// c := &entities.Competitions{}
	// c.ID = id
	if err := u.db.Where(query).Updates(in).Error; err != nil {
		log.Printf("error: UpdateCompetition: %s\n", err.Error())
		return errors.New("error: update competition failed")
	}
	return nil
}

func (u *postgresRepository) StartCompetitionAndAppendMatchToCompetition(query *entities.Competitions, in *entities.Competitions, matchs []entities.Matchs) error {
	if err := u.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where(query).Updates(in)
		if result.Error != nil {
			log.Printf("error: StartCompetitionAndAppendMatchToCompetition: %s\n", result.Error.Error())
			return result.Error
		}
		if result.RowsAffected == 0 {
			log.Printf("error: StartCompetitionAndAppendMatchToCompetition: no records found to update")
			return errors.New("error: no records found to update")
		}

		if err := tx.Model(query).Association("Matchs").Append(matchs); err != nil {
			log.Printf("error: StartCompetitionAndAppendMatchToCompetition: %s\n", err.Error())
			return err
		}

		return nil
	}); err != nil {
		return errors.New("error: start competition and append match to competition failed")
	}
	return nil
}

// AppendMatchToCompetition implements Userrepository.
// func (u *postgresRepository) AppendMatchToCompetition(competition *entities.Competitions, matchs []entities.Matchs) error {
// 	if err := u.db.Model(competition).Association("Matchs").Append(matchs); err != nil {
// 		return err
// 	}
// 	return nil
// }

// GetCompetitions implements Userrepository.
func (u *postgresRepository) GetCompetitions(in *entities.Competitions, orderString string, decs bool, limit int, offset int) ([]entities.Competitions, error) {
	competitions := []entities.Competitions{}
	// if err := u.db.Where(&in).Preload("Organizers").Order(clause.OrderByColumn{Column: clause.Column{Name: orderString}, Desc: decs}).Offset(offset).Limit(limit).Find(&competitions).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil, errors.New("record not found")
	// 	}
	// 	return nil, err
	// }
	if err := u.db.Where(&in).
		Preload("Teams.Teams").
		Preload("Organizers").
		Order(clause.OrderByColumn{Column: clause.Column{Name: orderString}, Desc: decs}).
		Offset(offset).
		Limit(limit).
		Find(&competitions).
		Error; err != nil {
		log.Printf("error: GetCompetitions: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get competitions failed")
	}
	return competitions, nil
}

// AppendJoinCodeToCompetition implements Userrepository.
func (u *postgresRepository) AppendJoinCodeToCompetition(id uint, joinCodes []entities.JoinCode) error {
	competition := &entities.Competitions{}
	competition.ID = id
	if err := u.db.Model(competition).Association("JoinCode").Append(joinCodes); err != nil {
		log.Printf("error: AppendJoinCodeToCompetition: %s\n", err.Error())
		return errors.New("error: add join code to competition failed")
	}
	return nil
}

func (u *postgresRepository) AddCompetitionToTeamAndNormalUsers(inComTeam *entities.CompetitionsTeams, normalUserID []uint) error {
	err := u.db.Transaction(func(tx *gorm.DB) error {

		competition := &entities.Competitions{}
		competition.ID = inComTeam.CompetitionsID
		if err := tx.First(competition).Error; err != nil {
			return err
		}

		team := &entities.Teams{}
		team.ID = inComTeam.TeamsID
		if err := tx.First(team).Error; err != nil {
			return err
		}

		for _, id := range normalUserID {
			nomalUser := &entities.NormalUsers{}
			nomalUser.ID = id
			if err := tx.First(nomalUser).Error; err != nil {
				return err
			}

			if err := tx.Create(&entities.NormalUsersCompetitions{
				NormalUsersID:  id,
				CompetitionsID: inComTeam.CompetitionsID,
				TeamsID:        inComTeam.TeamsID,
			}).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(inComTeam).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("error: AddCompetitionToTeamAndNormalUser: %s\n", err.Error())
		return errors.New("error: add competition to team and normaluser")
	}
	return nil
}

// UpdateJoinCode implements Userrepository.
func (u *postgresRepository) UpdateJoinCode(id uint, in *entities.JoinCode) error {
	joinCode := &entities.JoinCode{}
	joinCode.ID = id
	if err := u.db.Where(joinCode).Updates(in).Error; err != nil {
		log.Printf("error: UpdateJoinCode: %s\n", err.Error())
		return errors.New("error: update join code's status failed")
	}
	return nil
}

// GetMatchDetail implements Userrepository.
func (u *postgresRepository) GetMatchDetail(in *entities.Matchs) (*entities.Matchs, error) {
	match := &entities.Matchs{}
	if err := u.db.Where(in).Preload("GoalRecords").Preload("Competitions").First(match).Error; err != nil {
		log.Printf("error: GetMatch: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return match, nil
}

// GetTeam implements Userrepository.
func (u *postgresRepository) GetTeam(in *entities.Teams) (*entities.Teams, error) {
	team := &entities.Teams{}
	if err := u.db.Where(in).First(team).Error; err != nil {
		log.Printf("error: GetTeam: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return team, nil
}

// GetNormalUserCompetitions implements Userrepository.
func (u *postgresRepository) GetNormalUserCompetitions(in *entities.NormalUsersCompetitions) ([]entities.NormalUsersCompetitions, error) {

	result := []entities.NormalUsersCompetitions{}
	if err := u.db.Where(in).Preload("NormalUsers").Preload("Competitions").Find(&result).Error; err != nil {
		log.Printf("error: GetNormalUserCompetitions: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entities.NormalUsersCompetitions{}, errors.New("record not found")
		}
		return []entities.NormalUsersCompetitions{}, err
	}
	return result, nil
}

// UpdateMatch implements Userrepository.
func (u *postgresRepository) UpdateMatch(id uint, in *entities.Matchs) error {
	match := &entities.Matchs{}
	match.ID = id
	if err := u.db.Where(match).Select("*").Updates(in).Error; err != nil {
		log.Printf("error: UpdateMatch: %s\n", err.Error())
		return errors.New("error: update match failed")
	}
	return nil
}

// GetNormalUsers implements Userrepository.
func (u *postgresRepository) GetNormalUsers(in *entities.NormalUsers) ([]entities.NormalUsers, error) {
	normalUsers := []entities.NormalUsers{}
	if err := u.db.Where(&in).Preload("Users").Find(&normalUsers).Error; err != nil {
		log.Printf("error: GetNormalUsers: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get normalUser failed")
	}
	return normalUsers, nil
}

// UpdateNormalUser implements Userrepository.
func (u *postgresRepository) UpdateNormalUser(inNormalUser *entities.NormalUsers) error {
	// result := u.db.Model(&inNormalUser).Updates(inNormalUser)
	result := u.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&inNormalUser)
	if result.Error != nil {
		log.Printf("error: UpdateNormalUser: %s\n", result.Error.Error())
		return errors.New("error: update normalUser failed")
	}
	return nil
}

// GetOrganizers implements Userrepository.
func (u *postgresRepository) GetOrganizers() ([]entities.Organizers, error) {
	org := []entities.Organizers{}
	if err := u.db.Where(&entities.Organizers{}).Preload("Addresses").Preload("Users").Find(&org).Error; err != nil {
		log.Printf("error: GetOrganizers: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get organizers failed")
	}
	return org, nil
}

// GetOrganizer implements Userrepository.
func (u *postgresRepository) GetOrganizer(in *entities.Organizers) (*entities.Organizers, error) {
	org := new(entities.Organizers)
	if err := u.db.Where(in).Preload("Addresses").Preload("Competitions").Preload("Users").First(org).Error; err != nil {
		log.Printf("error: GetOrganizer: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get organizer failed")
	}
	return org, nil
}

// UpdateOrganizer implements Userrepository.
func (u *postgresRepository) UpdateOrganizer(id uint, in *entities.Organizers) error {
	org := &entities.Organizers{}
	org.ID = id
	if err := u.db.Where(org).Select("*").Session(&gorm.Session{FullSaveAssociations: true}).Updates(in).Error; err != nil {
		log.Printf("error: UpdateOrganizer: %s\n", err.Error())
		return errors.New("error: update organizer failed")
	}
	return nil
}

// GetTeams implements Userrepository.
func (u *postgresRepository) GetTeams(in *entities.Teams, orderString string, decs bool, limit int, offset int) ([]entities.Teams, error) {
	teams := []entities.Teams{}
	if err := u.db.Where(&in).Order(clause.OrderByColumn{Column: clause.Column{Name: orderString}, Desc: decs}).Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
		log.Printf("error: GetTeams: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get teams failed")
	}
	return teams, nil
}

func (u *postgresRepository) GetTeamWithMemberAndCompetitionByID(in uint) (*entities.Teams, error) {
	team := entities.Teams{}
	if err := u.db.Preload("TeamsMembers.NormalUsers").Preload("TeamsMembers.NormalUsers.Users").Preload("Competitions.Competitions").First(&team, in).Error; err != nil {
		log.Printf("error: GetTeamWithMemberAndCompetitionByID: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get team failed")
	}
	return &team, nil
}

// InsertTeam implements Userrepository.
func (u *postgresRepository) InsertTeam(in *entities.Teams) error {
	result := u.db.Create(in)
	if result.Error != nil {
		log.Printf("Error: InsertTeam: %s\n", result.Error.Error())
		return errors.New("error: create team failed")
	}
	return nil
}

// DeleteTeamMember implements Userrepository.
func (u *postgresRepository) DeleteTeamMember(nomalUserID uint, teamID uint) error {
	if err := u.db.
		Where("normal_users_id = ?", nomalUserID).
		Where("teams_id = ?", teamID).
		Delete(&entities.TeamsMembers{}).Error; err != nil {
		log.Printf("Error: DeleteTeamMember: %s\n", err.Error())
		return errors.New("error: remove normaluser from team failed")
	}
	return nil
}

func (u *postgresRepository) DeleteTeamFromCompetition(competitionID uint, teamID uint) error {
	err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("teams_id = ?", teamID).Where("competitions_id = ?", competitionID).Delete(&entities.NormalUsersCompetitions{}).Error; err != nil {
			return err
		}

		if err := tx.Where("competitions_id = ?", competitionID).Where("teams_id = ?", teamID).Delete(&entities.CompetitionsTeams{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("Error: DeleteTeamFromCompetition: %s\n", err.Error())
		return errors.New("error: delete team from competition failed")
	}
	return nil
}

// // DeleteNormalUserCompetition implements Userrepository.
// func (u *postgresRepository) DeleteNormalUserCompetitionByTeamIDAndCompetitionID(competitionID uint, teamID uint) error {
// 	if err := u.db.Where("teams_id = ?", teamID).Where("competitions_id = ?", competitionID).Delete(&entities.NormalUsersCompetitions{}).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// // DeleteCompetitionsTeam implements Userrepository.
// func (u *postgresRepository) DeleteCompetitionsTeam(competitionID uint, teamID uint) error {
// 	if err := u.db.Where("competitions_id = ?", competitionID).Where("teams_id = ?", teamID).Delete(&entities.CompetitionsTeams{}).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// GetUserByID implements Userrepository.
func (u *postgresRepository) GetUserByID(in uint) (*entities.Users, error) {
	selectedUser := &entities.Users{}
	if err := u.db.Where("id = ?", in).First(selectedUser).Error; err != nil {
		log.Printf("Error: GetUserByID: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get user failed")
	}
	return selectedUser, nil
}

// GetOrganizerWithAddressByUserID implements Userrepository.
func (u *postgresRepository) GetOrganizerWithAddressByUserID(in uint) (*entities.Organizers, error) {
	organizer := &entities.Organizers{}
	err := u.db.Model(&entities.Organizers{}).Preload("Addresses").Where("users_id = ?", in).First(&organizer).Error
	if err != nil {
		log.Printf("Error: GetOrganizerWithAddressByUserID: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get organizer failed")
	}
	return organizer, nil
}

// GetUsers implements Userrepository.
func (u *postgresRepository) GetUsers() ([]entities.Users, error) {
	users := []entities.Users{}
	if err := u.db.Find(&users).Order("id DESC").Error; err != nil {
		log.Printf("Error: GetUsers: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get users failed")
	}
	return users, nil
}

// GetTeamWithMemberAndRequestSentByID implements Userrepository.
func (u *postgresRepository) GetTeamWithMemberAndRequestSentByID(in uint) (*entities.Teams, error) {
	team := new(entities.Teams)
	if err := u.db.Model(&entities.Teams{}).
		Preload("TeamsMembers").
		Preload("Competitions").
		Preload("RequestSends").
		Where("id = ?", in).
		First(team).Error; err != nil {
		log.Printf("Error: GetTeamWithMemberAndRequestSentByID: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get team failed")
	}
	return team, nil
}

// InsertAddMemberRequest implements Userrepository.
func (u *postgresRepository) InsertAddMemberRequest(in *entities.AddMemberRequests) error {
	result := u.db.Create(in)
	if result.Error != nil {
		log.Printf("Error: InsertAddMemberRequest: %s\n", result.Error.Error())
		return errors.New("error: create add member request failed")
	}
	return nil
}

// UpdateAddMemberRequestStatusByID implements Userrepository.
func (u *postgresRepository) UpdateAddMemberRequestStatusByID(inID uint, inStatus string) error {
	result := u.db.Model(&entities.AddMemberRequests{}).Where("id = ?", inID).Update("status", inStatus)
	if result.Error != nil {
		log.Printf("Error: UpdateAddMemberRequestStatusByID: %s\n", result.Error.Error())
		return errors.New("error: update add member request failed")
	}
	return nil
}

// UpdateAddmemberRequestAndInsertTeamsMembers implements .
func (u *postgresRepository) UpdateAddmemberRequestAndInsertTeamsMembers(
	in *entities.TeamsMembers, memberReqID uint, memberReqStatus string,
) error {
	err := u.db.Transaction(func(tx *gorm.DB) error {

		result := tx.Model(&entities.AddMemberRequests{}).
		Where("id = ?", memberReqID).
		Update("status", memberReqStatus)
		if result.Error != nil {
			return result.Error
		}

		team := &entities.Teams{}
		team.ID = in.TeamsID
		if err := tx.First(team).Error; err != nil {
			return err
		}
		normalUser := &entities.NormalUsers{}
		normalUser.ID = in.NormalUsersID
		if err := tx.First(normalUser).Error; err != nil {
			return err
		}
		if err := tx.Create(in).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("Error: UpdateAddmemberRequestAndInsertTeamsMembers: %s\n", err.Error())
		return errors.New("error: fail to update add member request and add member to team")
	}
	return nil
}


// GetAddMemberRequestByID implements Userrepository.
func (u *postgresRepository) GetAddMemberRequestByID(in *entities.AddMemberRequests) ([]entities.AddMemberRequests, error) {
	addMemberRequests := []entities.AddMemberRequests{}
	log.Print(in)
	if err := u.db.Where(in).Preload("Teams").Find(&addMemberRequests).Error; err != nil {
		log.Printf("Error: GetAddMemberRequestByID: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get add member request failed")
	}
	return addMemberRequests, nil
}

// update status and soft delete by ID within a transaction
func (u *postgresRepository) UpdateAddMemberRequestStatusAndSoftDelete(inReq *entities.AddMemberRequests, inStatus string) error {
	if err := u.db.Transaction(func(tx *gorm.DB) error {
		inReq.Status = inStatus

		//Update Req status
		if err := tx.Save(inReq).Error; err != nil {
			return err
		}

		// Soft delete by ID
		if err := tx.Delete(inReq).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Printf("Error: UpdateAddMemberRequestStatusAndSoftDelete: %s\n", err.Error())
		return errors.New("error: update add member request failed")
	}
	return nil
}

func (u *postgresRepository) GetNumberOfTeamsMember(in uint) int64 {
	var count int64
	u.db.Model(&entities.TeamsMembers{}).Where("teams_id = ?", in).Count(&count)
	return count
}

// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================

// UpdateCompetitionsTeams implements Userrepository.
func (u *postgresRepository) UpdateCompetitionsTeams(in *entities.CompetitionsTeams) error {
	competitionTeam := &entities.CompetitionsTeams{
		TeamsID:        in.TeamsID,
		CompetitionsID: in.CompetitionsID,
	}
	if err := u.db.Where(competitionTeam).Updates(in).Error; err != nil {
		log.Printf("Error: UpdateCompetitionsTeams: %s\n", err.Error())
		return errors.New("error: update competitionTeams failed")
	}
	return nil
}

// ClearGoalRecordsOfMatch implements Userrepository.
func (u *postgresRepository) ClearGoalRecordsOfMatch(matchID uint) error {
	match := &entities.Matchs{}
	match.ID = matchID

	// Soft Delete
	if err := u.db.Model(match).Association("GoalRecords").Unscoped().Clear(); err != nil {
		log.Printf("Error: ClearGoalRecordsOfMatch: %s\n", err.Error())
		return errors.New("error: clear goal recoreds of this match failed")
	}
	return nil
}

// ReplaceGoalRecordsOfMatch implements Userrepository.
func (u *postgresRepository) ReplaceGoalRecordsOfMatch(matchID uint, goalRecords []entities.GoalRecords) error {
	match := &entities.Matchs{}
	match.ID = matchID
	if err := u.db.Model(match).Association("GoalRecords").Replace(goalRecords); err != nil {
		log.Printf("Error: ReplaceGoalRecordsOfMatch: %s\n", err.Error())
		return errors.New("error: replace goal recoreds of this match failed")	
	}
	return nil
}

// AppendGoalRecordsToMatch implements Userrepository.
func (u *postgresRepository) AppendGoalRecordsToMatch(id uint, goalRecords []entities.GoalRecords) error {
	match := &entities.Matchs{}
	match.ID = id
	if err := u.db.Model(match).Association("GoalRecords").Append(goalRecords); err != nil {
		log.Printf("Error: AppendGoalRecordsToMatch: %s\n", err.Error())
		return errors.New("error: append goal recoreds of this match failed")	
	}
	return nil
}


// GetTeamsWithCompetitionAndMatch implements Userrepository.
func (u *postgresRepository) GetTeamsWithCompetitionAndMatch(in *entities.Teams) (*entities.Teams, error) {
	team := &entities.Teams{}
	// if err := u.db.Where(in).Preload("Competitions").Preload("Compeitions.Matchs").First(team).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil, errors.New("record not found")
	// 	}
	// 	return nil, err
	// }
	if err := u.db.Where(in).Preload("Competitions.Competitions").Preload("Competitions.Competitions.Matchs").First(team).Error; err != nil {
		log.Printf("Error: GetTeamsWithCompetitionAndMatch: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("error: get teams failed")
	}
	return team, nil
}

// GetMatchs implements Userrepository.
func (u *postgresRepository) GetMatchs(in *entities.Matchs) ([]entities.Matchs, error) {
	match := []entities.Matchs{}
	if err := u.db.Where(in).Find(match).Error; err != nil {
		log.Printf("Error: GetMatchs: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("errors: get matches failed")
	}
	return match, nil
}



// GetTeam implements Userrepository.
func (u *postgresRepository) GetTeamWithAllAssociationsByID(in *entities.Teams) (*entities.Teams, error) {
	team := new(entities.Teams)
	if err := u.db.Where(in).Preload(clause.Associations).Preload("TeamsMembers.NormalUsers").First(team).Error; err != nil {
		log.Printf("Error: GetTeamWithAllAssociationsByID: %s\n", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, errors.New("errors: get matches failed")
	}
	return team, nil
}


