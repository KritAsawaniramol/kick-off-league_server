package repositories

import (
	"errors"
	"fmt"

	"github.com/labstack/gommon/log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"kickoff-league.com/entities"
	"kickoff-league.com/util"
)

type userPostgresRepository struct {
	db *gorm.DB
}

// AppendGoalRecordsToMatch implements Userrepository.
func (u *userPostgresRepository) AppendGoalRecordsToMatch(id uint, goalRecords []entities.GoalRecords) error {
	match := &entities.Matchs{}
	match.ID = id
	if err := u.db.Model(match).Association("GoalRecords").Append(goalRecords); err != nil {
		return err
	}
	return nil
}

// UpdateMatch implements Userrepository.
func (h *userPostgresRepository) UpdateMatch(id uint, in *entities.Matchs) error {
	match := &entities.Matchs{}
	match.ID = id
	if err := h.db.Model(match).Updates(in).Error; err != nil {
		return err
	}
	return nil
}

// AppendMatchToCompatition implements Userrepository.
func (u *userPostgresRepository) AppendMatchToCompatition(compatition *entities.Compatitions, matchs []entities.Matchs) error {
	if err := u.db.Model(compatition).Association("Matchs").Append(matchs); err != nil {
		return err
	}
	return nil
}

// UpdateCompatition implements Userrepository.
func (u *userPostgresRepository) UpdateCompatition(id uint, in *entities.Compatitions) error {
	c := &entities.Compatitions{}
	c.ID = id
	if err := u.db.Where(c).Updates(in).Error; err != nil {
		return err
	}
	return nil
}

// AppendTeamtoCompatition implements Userrepository.
func (u *userPostgresRepository) AppendTeamtoCompatition(compatition *entities.Compatitions, newTeam *entities.Teams) error {
	if err := u.db.Model(compatition).Association("Teams").Append(newTeam); err != nil {
		return err
	}
	return nil
}

// InsertMatchs implements Userrepository.
func (*userPostgresRepository) InsertMatchs(in []entities.Matchs) error {
	panic("unimplemented")
}

// UpdateUser implements Userrepository.
func (*userPostgresRepository) UpdateUser(inUser *entities.Users) error {
	panic("unimplemented")
}

// UpdateSelectedFields implements Userrepository.
func (u *userPostgresRepository) UpdateSelectedFields(model interface{}, fieldname string, value interface{}) error {
	//model.ID is required
	if err := u.db.Model(model).Select(fieldname).Updates(value).Error; err != nil {
		return err
	}
	return nil
}

// GetUsers implements Userrepository.
func (u *userPostgresRepository) GetUsers() ([]entities.Users, error) {
	users := []entities.Users{}
	if err := u.db.Find(&users).Order("id DESC").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return users, nil
}

// GetCompatitionByID implements Userrepository.
func (u *userPostgresRepository) GetCompatition(in *entities.Compatitions) (*entities.Compatitions, error) {
	compatition := &entities.Compatitions{}
	if err := u.db.Where(&in).Preload(clause.Associations).Preload("Organizers.Addresses").Preload("Matchs.GoalRecords").Preload("Teams.TeamsMembers").First(&compatition).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return compatition, nil
}

// GetCompatitions implements Userrepository.
func (u *userPostgresRepository) GetCompatitions(in *entities.Compatitions, orderString string, decs bool, limit int, offset int) ([]entities.Compatitions, error) {
	compatitions := []entities.Compatitions{}
	util.PrintObjInJson(in)
	if err := u.db.Where(&in).Preload("Organizers").Order(clause.OrderByColumn{Column: clause.Column{Name: orderString}, Desc: decs}).Offset(offset).Limit(limit).Find(&compatitions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return compatitions, nil
}

// GetNormalUsers implements Userrepository.
func (u *userPostgresRepository) GetNormalUsers(in *entities.NormalUsers) ([]entities.NormalUsers, error) {
	normalUsers := []entities.NormalUsers{}
	if err := u.db.Where(&in).Find(&normalUsers).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return normalUsers, nil
}

// GetTeams implements Userrepository.
func (u *userPostgresRepository) GetTeams(in *entities.Teams, orderString string, decs bool, limit int, offset int) ([]entities.Teams, error) {
	teams := []entities.Teams{}
	if err := u.db.Where(&in).Order(clause.OrderByColumn{Column: clause.Column{Name: orderString}, Desc: decs}).Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return teams, nil
}

// GetOrganizer implements Userrepository.
func (u *userPostgresRepository) GetOrganizer(in *entities.Organizers) (*entities.Organizers, error) {
	org := new(entities.Organizers)
	if err := u.db.Where(in).First(org).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return org, nil
}

// GetOrganizerWithAddressByUserID implements Userrepository.
func (u *userPostgresRepository) GetOrganizerWithAddressByUserID(in uint) (*entities.Organizers, error) {
	organizer := &entities.Organizers{}
	err := u.db.Model(&entities.Organizers{}).Preload("Addresses").Where("users_id = ?", in).First(&organizer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return organizer, nil
}

// GetNormalUser implements Userrepository.
func (u *userPostgresRepository) GetNormalUser(in *entities.NormalUsers) (*entities.NormalUsers, error) {
	normalUser := new(entities.NormalUsers)
	if err := u.db.Where(&in).Preload(clause.Associations).Preload("Teams.Teams").First(normalUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return normalUser, nil
}

// GetNormalUserWithAddressByUserID implements Userrepository.
func (u *userPostgresRepository) GetNormalUserWithAddressByUserID(in uint) (*entities.NormalUsers, error) {
	normalUser := &entities.NormalUsers{}
	if err := u.db.Model(&entities.NormalUsers{}).Preload("Addresses").Where("users_id = ?", in).First(&normalUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return normalUser, nil
}

// GetTeamMembersByTeamID implements Userrepository.
func (u *userPostgresRepository) GetTeamMembersByTeamID(in uint, orderString string, decs bool, limit int, offset int) ([]entities.TeamsMembers, error) {
	teamMember := &entities.TeamsMembers{}
	teamMember.TeamsID = in
	teamMembers := []entities.TeamsMembers{}
	if err := u.db.Where(&teamMember).Order(clause.OrderByColumn{Column: clause.Column{Name: orderString}, Desc: decs}).Offset(offset).Limit(limit).Find(&teamMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return teamMembers, nil
}

// InsertTeamsMembers implements .
func (u *userPostgresRepository) InsertTeamsMembers(in *entities.TeamsMembers) error {
	err := u.db.Transaction(func(tx *gorm.DB) error {
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
		return err
	}
	return nil
}

// GetTeam implements Userrepository.
func (u *userPostgresRepository) GetTeam(in uint) (*entities.Teams, error) {
	team := &entities.Teams{}
	team.ID = in
	if err := u.db.Where("id = ?", in).First(team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return team, nil
}

func (u *userPostgresRepository) GetNumberOfTeamsMember(in uint) int64 {
	var count int64
	u.db.Model(&entities.TeamsMembers{}).Where("teams_id = ?", in).Count(&count)
	return count
}

// GetAddMemberRequestByID implements Userrepository.
func (u *userPostgresRepository) GetAddMemberRequestByID(in *entities.AddMemberRequests) ([]entities.AddMemberRequests, error) {
	addMemberRequests := []entities.AddMemberRequests{}
	log.Print(in)
	if err := u.db.Where(in).Preload("Teams").Find(&addMemberRequests).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return addMemberRequests, nil
}

// update status and soft delete by ID within a transaction
func (u *userPostgresRepository) UpdateAddMemberRequestStatusAndSoftDelete(inReq *entities.AddMemberRequests, inStatus string) error {
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
		return err
	}
	return nil
}

func (u *userPostgresRepository) GetTeamWithMemberAndCompatitionByID(in uint) (*entities.Teams, error) {
	team := entities.Teams{}
	if err := u.db.Preload("TeamsMembers.NormalUsers").Preload("Compatitions").First(&team, in).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	util.PrintObjInJson(team)
	return &team, nil
}

// GetTeam implements Userrepository.
func (u *userPostgresRepository) GetTeamWithAllAssociationsByID(in *entities.Teams) (*entities.Teams, error) {
	team := new(entities.Teams)
	if err := u.db.Where(in).Preload(clause.Associations).Preload("TeamsMembers.NormalUsers").First(team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return team, nil
}

// GetTeamWithMemberAndRequestSendByID implements Userrepository.
func (u *userPostgresRepository) GetTeamWithMemberAndRequestSendByID(in uint) (*entities.Teams, error) {
	team := new(entities.Teams)
	if err := u.db.Model(&entities.Teams{}).Preload("TeamsMembers").Preload("Compatitions").Preload("RequestSends").Where("id = ?", in).First(team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return team, nil
}

// UpdateAddMemberRequestStatusByID implements Userrepository.
func (u *userPostgresRepository) UpdateAddMemberRequestStatusByID(inID uint, inStatus string) error {
	result := u.db.Model(&entities.AddMemberRequests{}).Where("id = ?", inID).Update("status", inStatus)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// InsertAddMemberRequest implements Userrepository.
func (u *userPostgresRepository) InsertAddMemberRequest(in *entities.AddMemberRequests) error {
	util.PrintObjInJson(in)
	result := u.db.Create(in)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateNormalUser implements Userrepository.
func (u *userPostgresRepository) UpdateNormalUser(inNormalUser *entities.NormalUsers) error {
	// result := u.db.Model(&inNormalUser).Updates(inNormalUser)
	result := u.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&inNormalUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// InsertTeam implements Userrepository.
func (u *userPostgresRepository) InsertTeam(in *entities.Teams) error {
	result := u.db.Create(in)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// InserrtCompatition implements Userrepository.
func (u *userPostgresRepository) InsertCompatition(in *entities.Compatitions) error {
	util.PrintObjInJson(in)
	if err := u.db.Create(in).Error; err != nil {
		return err
	}
	return nil
}

// InsertUserWihtOrganizerAndAddress implements Userrepository.
func (u *userPostgresRepository) InsertUserWihtOrganizerAndAddress(in_organizer *entities.Organizers, in_user *entities.Users) error {
	if err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&in_user).Error; err != nil {
			return err
		}

		address := &entities.Addresses{}
		fmt.Println("address before: ")
		util.PrintObjInJson(address)

		if err := tx.Create(address).Error; err != nil {
			return err
		}
		fmt.Println("address after: ")
		util.PrintObjInJson(address)

		in_organizer.UsersID = in_user.ID
		in_organizer.AddressesID = address.ID

		if err := tx.Create(&in_organizer).Error; err != nil {
			return err
		}
		return nil
	},
	); err != nil {
		return err
	}
	return nil
}

// InsertUserWihtNormalUserAndAddress implements Userrepository.
func (u *userPostgresRepository) InsertUserWihtNormalUserAndAddress(in_normalUser *entities.NormalUsers, in_user *entities.Users) error {
	err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&in_user).Error; err != nil {
			return err
		}

		address := &entities.Addresses{}

		if err := tx.Create(address).Error; err != nil {
			return err
		}

		in_normalUser.UsersID = in_user.ID
		in_normalUser.AddressesID = address.ID

		if err := tx.Create(&in_normalUser).Error; err != nil {
			return err
		}
		return nil
	},
	)
	return err
}

// UpdateNormalUserPhone implements Userrepository.
func (u *userPostgresRepository) UpdateNormalUserPhone(in_userID uint, newPhone string) error {
	// Check if new Phone already exists
	existingUser := &entities.NormalUsers{}
	err := u.db.Model(&entities.NormalUsers{}).Where("phone = ?", newPhone).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {

		log.Errorf(err.Error())
		return err
	}

	if existingUser.ID == 0 {
		if err := u.db.Model(&entities.NormalUsers{}).Where("user_id = ?", in_userID).Update("phone", newPhone).Error; err != nil {
			return err
		}
	}

	return nil

}

func NewUserPostgresRepository(db *gorm.DB) Userrepository {
	return &userPostgresRepository{db: db}
}

// GetUserByEmail implements Userrepository.
func (u *userPostgresRepository) GetUserByEmail(email string) (*entities.Users, error) {
	selectedUser := &entities.Users{}
	if err := u.db.Where("email = ?", email).First(selectedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return selectedUser, nil
}

// GetNormalUserByUsername implements Userrepository.
func (u *userPostgresRepository) GetNormalUserByUsername(username string) (*entities.NormalUsers, error) {
	selectedUser := &entities.NormalUsers{}
	if err := u.db.Where("username = ?", username).First(selectedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return selectedUser, nil
}

// GetUserByID implements Userrepository.
func (u *userPostgresRepository) GetUserByID(in uint) (*entities.Users, error) {
	selectedUser := &entities.Users{}
	if err := u.db.Where("id = ?", in).First(selectedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return selectedUser, nil
}
