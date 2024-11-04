package teamUsecase

import (
	"errors"
	"log"
	"os"
	"strings"

	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
	"kickoff-league.com/util"
)

type teamUsecaseImpl struct {
	repository repositories.Repository
}

func NewTeamUsecaseImpl(
	repository repositories.Repository,
) TeamUsecase {
	return &teamUsecaseImpl{
		repository: repository,
	}
}

// GetTeamWithMemberAndCompatitionByID implements TeamUsecase.
func (t *teamUsecaseImpl) GetTeamWithMemberAndCompatitionByID(id uint) (*model.Team, error) {
	selectedTeams, err := t.repository.GetTeamWithMemberAndCompetitionByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			return &model.Team{}, nil
		}
		return &model.Team{}, err
	}
	memberList := []model.Member{}
	for _, member := range selectedTeams.TeamsMembers {
		memberList = append(memberList, model.Member{
			ID:               member.NormalUsers.ID,
			UsersID:          member.NormalUsers.UsersID,
			FirstNameThai:    member.NormalUsers.FirstNameThai,
			LastNameThai:     member.NormalUsers.LastNameThai,
			FirstNameEng:     member.NormalUsers.FirstNameEng,
			LastNameEng:      member.NormalUsers.LastNameEng,
			Position:         member.NormalUsers.Position,
			Sex:              member.NormalUsers.Sex,
			Role:             member.Role,
			ImageProfilePath: member.NormalUsers.Users.ImageProfilePath,
			ImageCoverPath:   member.NormalUsers.Users.ImageCoverPath,
		})
	}

	compatition_model := []model.CompatitionBasicInfo{}

	for _, v := range selectedTeams.Competitions {
		compatition_model = append(compatition_model, model.CompatitionBasicInfo{
			ID:           v.CompetitionsID,
			Name:         v.Competitions.Name,
			Format:       v.Competitions.Format,
			OrganizerID:  v.Competitions.OrganizersID,
			StartDate:    v.Competitions.StartDate,
			EndDate:      v.Competitions.EndDate,
			AgeOver:      v.Competitions.AgeOver,
			AgeUnder:     v.Competitions.AgeUnder,
			Sex:          v.Competitions.Sex,
			FieldSurface: v.Competitions.FieldSurface,
			Description:  v.Competitions.Description,
			Status:       v.Competitions.Status,
			NumberOfTeam: v.Competitions.NumberOfTeam,
			ImageBanner:  v.Competitions.ImageBannerPath,
			Rank:         v.Rank,
			RankNumber:   v.RankNumber,
			Sport:        v.Competitions.Sport,
		})
	}

	return &model.Team{
		ID:               selectedTeams.ID,
		Name:             selectedTeams.Name,
		OwnerID:          selectedTeams.OwnerID,
		Members:          memberList,
		Compatitions:     compatition_model,
		Description:      selectedTeams.Description,
		ImageCoverPath:   selectedTeams.ImageCoverPath,
		ImageProfilePath: selectedTeams.ImageProfilePath,
	}, nil
}

// GetTeams implements TeamUsecase.
func (t *teamUsecaseImpl) GetTeams(in *model.GetTeamsReq) ([]model.TeamList, error) {
	team := entities.Teams{
		// 0 is select all id
		OwnerID: in.NormalUserID,
	}
	team.ID = in.TeamID

	limit := int(in.PageSize)
	if limit <= 0 {
		limit = -1
	}

	offset := int(in.PageSize * in.Page)
	if offset <= 0 {
		offset = -1
	}

	if in.Ordering == "" {
		in.Ordering = "id"
	}

	teams, err := t.repository.GetTeams(&team, strings.Trim(in.Ordering, " "), in.Decs, limit, offset)
	if err != nil {
		log.Printf("Error: GetTeams failed: %v\n", err.Error())
		return []model.TeamList{}, errors.New("error: get teams failed")
	}

	teamList := []model.TeamList{}
	for _, team := range teams {
		teamList = append(teamList, model.TeamList{
			ID:               team.ID,
			Name:             team.Name,
			Description:      team.Description,
			NumberOfMember:   uint(t.repository.GetNumberOfTeamsMember(team.ID)),
			ImageProfilePath: team.ImageProfilePath,
			ImageCoverPath:   team.ImageCoverPath,
		})
	}
	return teamList, nil
}

// CreateTeam implements TeamUsecase.
func (t *teamUsecaseImpl) CreateTeam(in *model.CreateTeam) error {
	normalUser, err := t.repository.GetNormalUserDetails(&entities.NormalUsers{
		UsersID: in.OwnerID,
	})
	if err != nil {
		log.Printf("Error: CreateTeam failed: %v\n", err.Error())
		return errors.New("error: create team failed")
	}

	//check required data
	requiredData := []string{}
	if normalUser.FirstNameThai == "" {
		requiredData = append(requiredData, "first_name_thai")
	}
	if normalUser.LastNameThai == "" {
		requiredData = append(requiredData, "last_name_thai")
	}
	if normalUser.FirstNameEng == "" {
		requiredData = append(requiredData, "first_name_eng")
	}
	if normalUser.LastNameEng == "" {
		requiredData = append(requiredData, "last_name_eng")
	}
	if normalUser.Born.IsZero() {
		requiredData = append(requiredData, "born")
	}
	if normalUser.Sex == "" {
		requiredData = append(requiredData, "sex")
	}
	if normalUser.Nationality == "" {
		requiredData = append(requiredData, "nationality")
	}
	if normalUser.Phone == "" {
		requiredData = append(requiredData, "phone")
	}
	if len(requiredData) != 0 {
		return &util.CreateTeamError{
			RequiredData: requiredData,
		}
	}

	team := entities.Teams{
		Name:             in.Name,
		OwnerID:          in.OwnerID,
		Description:      in.Description,
		Competitions:     []entities.CompetitionsTeams{},
		ImageProfilePath: "./images/default/defaultProfile.jpg",
		ImageCoverPath:   "./images/default/defaultCover.jpg",
	}

	if err := t.repository.InsertTeam(&team); err != nil {
		return err
	}

	return nil
}

// RemoveNormalUserFormTeam implements TeamUsecase.
func (t *teamUsecaseImpl) RemoveNormalUserFormTeam(teamID uint, nomalUserID uint, ownerID uint) error {

	team, err := t.repository.GetTeam(&entities.Teams{
		OwnerID: ownerID,
	})
	if err != nil {
		return err
	}
	if team.OwnerID != ownerID {
		return errors.New("error: you can't remove member from this team")
	}
	err = t.repository.DeleteTeamMember(nomalUserID, teamID)
	if err != nil {
		return err
	}
	return nil
}

// RemoveTeamFormCompatition implements TeamUsecase.
func (t *teamUsecaseImpl) RemoveTeamFormCompatition(teamID uint, compatitionID uint, orgID uint) error {
	queryCompetition := &entities.Competitions{}
	queryCompetition.ID = compatitionID
	competition, err := t.repository.GetCompetition(queryCompetition)
	if err != nil {
		return err
	}

	if competition.OrganizersID != orgID {
		return errors.New("error: you can't remove team from this competition")
	}

	if err := t.repository.DeleteTeamFromCompetition(compatitionID, teamID); err != nil {
		return err
	}
	return nil
}

// RemoveTeamImageProfile implements TeamUsecase.
func (t *teamUsecaseImpl) RemoveTeamImageProfile(teamID uint, userID uint) error {
	team := &entities.Teams{}
	team.ID = teamID
	team.OwnerID = userID

	team, err := t.repository.GetTeam(team)
	if err != nil {
		return err
	}

	if !strings.Contains(team.ImageProfilePath, "/default/") {
		if err := os.Remove(team.ImageProfilePath); err != nil {
			log.Printf("error: RemoveTeamImageProfile: %s\n", err.Error())
		}
	}

	if err := t.repository.UpdateSelectedFields(
		team,
		"ImageProfilePath",
		&entities.Teams{ImageProfilePath: "./images/default/defaultProfile.jpg"},
	); err != nil {
		return err
	}
	return nil
}

// RemoveTeamImageCover implements TeamUsecase.
func (t *teamUsecaseImpl) RemoveTeamImageCover(teamID uint, userID uint) error {
	team := &entities.Teams{}
	team.ID = teamID
	team.OwnerID = userID

	team, err := t.repository.GetTeam(team)
	if err != nil {
		return err
	}

	if !strings.Contains(team.ImageCoverPath, "/default/") {
		if err := os.Remove(team.ImageCoverPath); err != nil {
			log.Printf("error: RemoveTeamImageCover: %s\n", err.Error())
		}
	}

	if err := t.repository.UpdateSelectedFields(
		team,
		"ImageCoverPath",
		&entities.Teams{ImageCoverPath: "./images/default/defaultCover.jpg"},
	); err != nil {
		return err
	}
	return nil
}

// UpdateTeamImageCover implements TeamUsecase.
func (t *teamUsecaseImpl) UpdateTeamImageCover(teamID uint, newImagePath string, userID uint) error {
	team := &entities.Teams{}
	team.ID = teamID
	team.OwnerID = userID

	team, err := t.repository.GetTeam(team)
	if err != nil {
		return err
	}

	if !strings.Contains(team.ImageCoverPath, "/default/") {
		if err := os.Remove(team.ImageCoverPath); err != nil {
			log.Printf("error: UpdateTeamImageCover: %s\n", err.Error())
		}
	}

	if err := t.repository.UpdateSelectedFields(team,
		"ImageCoverPath",
		&entities.Teams{ImageCoverPath: newImagePath}); err != nil {
		return err
	}
	return nil
}

// UpdateTeamImageProfile implements TeamUsecase.
func (t *teamUsecaseImpl) UpdateTeamImageProfile(teamID uint, newImagePath string, userID uint) error {

	team := &entities.Teams{}
	team.ID = teamID
	team.OwnerID = userID

	team, err := t.repository.GetTeam(team)
	if err != nil {
		return err
	}

	if !strings.Contains(team.ImageProfilePath, "/default/") {
		if err := os.Remove(team.ImageProfilePath); err != nil {
			log.Printf("error: UpdateTeamImageProfile: %s\n", err.Error())
		}
	}

	if err := t.repository.UpdateSelectedFields(team, "ImageProfilePath", &entities.Teams{ImageProfilePath: newImagePath}); err != nil {
		return err
	}
	return nil
}

// GetTeamsByOwnerID implements TeamUsecase.
func (t *teamUsecaseImpl) GetTeamsByOwnerID(in uint) ([]model.TeamList, error) {
	team := entities.Teams{
		// 0 is select all id
		OwnerID: in,
	}

	teams, err := t.repository.GetTeams(&team, "id", false, -1, -1)
	if err != nil {
		log.Printf("Error: GetTeamsByOwnerID failed: %v\n", err.Error())
		return []model.TeamList{}, errors.New("error: get teams failed")
	}

	teamList := []model.TeamList{}
	for _, team := range teams {
		teamList = append(teamList, model.TeamList{
			ID:               team.ID,
			Name:             team.Name,
			Description:      team.Description,
			NumberOfMember:   uint(t.repository.GetNumberOfTeamsMember(team.ID)),
			OwnerID:          team.OwnerID,
			ImageProfilePath: team.ImageProfilePath,
			ImageCoverPath:   team.ImageCoverPath,
		})
	}
	return teamList, nil
}

// ========================================================================
// ========================================================================
// ========================================================================
// ========================================================================
// ========================================================================
// ========================================================================
