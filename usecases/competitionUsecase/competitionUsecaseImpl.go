package competitionUsecase

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"sort"
	"time"

	"github.com/google/uuid"
	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
	"kickoff-league.com/util"
)

type competitionUsecaseImpl struct {
	repository repositories.Repository
}

func NewCompetitionUsecaseImpl(
	repository repositories.Repository,
) CompetitionUsecase {
	return &competitionUsecaseImpl{
		repository: repository,
	}
}

// RemoveImageBanner implements CompetitionUsecase.
func (c *competitionUsecaseImpl) RemoveImageBanner(compatitionID uint, orgID uint) error {
	compatition := &entities.Competitions{}
	compatition.ID = compatitionID
	compatition.OrganizersID = orgID
	if err := c.repository.UpdateSelectedFields(
		compatition, "ImageBannerPath",
		&entities.Competitions{ImageBannerPath: "./images/default/defaultBanner.png"},
	); err != nil {
		return err
	}
	return nil
}

// UpdateImageBanner implements CompetitionUsecase.
func (c *competitionUsecaseImpl) UpdateImageBanner(compatitionID uint, orgID uint, newImagePath string) error {
	compatition := &entities.Competitions{}
	compatition.ID = compatitionID
	compatition.OrganizersID = orgID
	if err := c.repository.UpdateSelectedFields(
		compatition, "ImageBannerPath",
		&entities.Competitions{ImageBannerPath: newImagePath},
	); err != nil {
		return err
	}
	return nil
}

// CreateCompetition implements CompetitionUsecase.
func (c *competitionUsecaseImpl) CreateCompetition(in *model.CreateCompetition) error {
	compatition := &entities.Competitions{
		Name:                 in.Name,
		Sport:                in.Sport,
		Type:                 in.Type,
		Format:               in.Format,
		Description:          in.Description,
		Rule:                 in.Rule,
		Prize:                in.Prize,
		StartDate:            in.StartDate,
		EndDate:              in.EndDate,
		ApplicationType:      in.ApplicationType,
		ImageBannerPath:      "./images/default/defaultBanner.png",
		AgeOver:              in.AgeOver,
		AgeUnder:             in.AgeUnder,
		Sex:                  in.Sex,
		NumberOfTeam:         in.NumberOfTeam,
		NumOfPlayerInTeamMin: in.NumOfPlayerInTeamMin,
		NumOfPlayerInTeamMax: in.NumOfPlayerInTeamMax,
		FieldSurface:         in.FieldSurface,
		OrganizersID:         in.OrganizerID,
		HouseNumber:          in.Address.HouseNumber,
		Village:              in.Address.Village,
		Subdistrict:          in.Address.Subdistrict,
		District:             in.Address.District,
		PostalCode:           in.Address.PostalCode,
		Country:              in.Address.Country,
		ContactType:          in.ContactType,
		Contact:              in.Contact,
		Status:               "Coming soon",
	}
	if in.Type == "Tournament" {
		if util.CheckNumberPowerOfTwo(int(in.NumberOfTeam)) != 0 {
			return errors.New("error: number of Team for create competition(tounament) is not power of 2")
		}
		if in.NumberOfTeam < 2 {
			return errors.New("error: number of Team have to morn than 1")
		}

	}

	if in.Type != "Tournament" && in.Type != "Round Robin" {
		return errors.New("error: undefined compatition type")
	}
	if err := c.repository.InsertCompetition(compatition); err != nil {
		return err
	}
	return nil
}

// GetCompatition implements CompetitionUsecase.
func (c *competitionUsecaseImpl) GetCompatition(in uint) (*model.GetCompatition, error) {
	compatitionEntity := &entities.Competitions{}
	compatitionEntity.ID = in
	result, err := c.repository.GetCompetitionDetails(compatitionEntity)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	joinCode := []model.JoinCode{}

	for _, v := range result.JoinCode {
		joinCode = append(joinCode, model.JoinCode{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			Code:      v.Code,
			Status:    v.Status,
		})
	}

	Matchs := []model.Match{}
	goalsScored := map[uint]int{}
	goalsConceded := map[uint]int{}

	temes := []model.Team{}

	for _, v := range result.Teams {
		// members := []model.Member{}
		// for _, member := range v.Teams.TeamsMembers {
		// 	members = append(members, model.Member{
		// 		ID:            member.ID,
		// 		UsersID:       member.NormalUsers.UsersID,
		// 		FirstNameThai: member.NormalUsers.FirstNameThai,
		// 		LastNameThai:  member.NormalUsers.LastNameThai,
		// 		FirstNameEng:  member.NormalUsers.FirstNameEng,
		// 		LastNameEng:   member.NormalUsers.LastNameEng,
		// 		Position:      member.NormalUsers.Position,
		// 		Sex:           member.NormalUsers.Sex,
		// 	})
		// }

		temes = append(temes, model.Team{
			ID:               v.TeamsID,
			Name:             v.Teams.Name,
			OwnerID:          v.Teams.OwnerID,
			Description:      v.Teams.Description,
			Members:          nil,
			Rank:             v.Rank,
			RankNumber:       v.RankNumber,
			Point:            v.Point,
			ImageProfilePath: v.Teams.ImageProfilePath,
			ImageCoverPath:   v.Teams.ImageCoverPath,
		})
		goalsScored[v.TeamsID] = 0
		goalsConceded[v.TeamsID] = 0
	}

	for _, v := range result.Matchs {
		goalRecords := []model.GoalRecord{}

		for _, goalRecord := range v.GoalRecords {
			goalRecords = append(goalRecords, model.GoalRecord{
				MatchsID: goalRecord.MatchsID,
				TeamID:   goalRecord.TeamsID,
				PlayerID: goalRecord.PlayerID,
			})
		}

		m := model.Match{
			ID:             v.ID,
			Index:          v.Index,
			DateTime:       v.DateTime,
			Team1ID:        v.Team1ID,
			Team2ID:        v.Team2ID,
			Team1Goals:     v.Team1Goals,
			Team2Goals:     v.Team2Goals,
			Round:          v.Round,
			NextMatchIndex: v.NextMatchIndex,
			NextMatchSlot:  v.NextMatchSlot,
			GoalRecords:    goalRecords,
			Result:         v.Result,
			VideoURL:       v.VideoURL,
		}

		for i := 0; i < len(result.Teams); i++ {
			if m.Team1ID == result.Teams[i].TeamsID {
				m.Team1Name = result.Teams[i].Teams.Name
			}
			if m.Team2ID == result.Teams[i].TeamsID {
				m.Team2Name = result.Teams[i].Teams.Name
			}
		}

		Matchs = append(Matchs, m)

		goalsScored[v.Team1ID] += v.Team1Goals
		goalsConceded[v.Team1ID] += v.Team2Goals

		goalsScored[v.Team2ID] += v.Team2Goals
		goalsConceded[v.Team2ID] += v.Team1Goals
	}

	sort.SliceStable(Matchs, func(i, j int) bool {
		return Matchs[i].ID < Matchs[j].ID
	})

	for i := 0; i < len(temes); i++ {
		temes[i].GoalsScored = goalsScored[temes[i].ID]
		temes[i].GoalsConceded = goalsConceded[temes[i].ID]
	}

	temesForSort := []model.Team{}
	temesForSort = append(temesForSort, temes...)

	if result.Type == util.CompetitionType[1] {
		sort.Slice(temesForSort, func(i, j int) bool {
			iv, jv := temesForSort[i], temesForSort[j]
			switch {
			case iv.Point != jv.Point:
				return iv.Point > jv.Point

			default:

				return (iv.GoalsScored - iv.GoalsConceded) > (jv.GoalsScored - jv.GoalsConceded)
			}
		})

		for i := 0; i < len(temes); i++ {
			for j := 0; j < len(temesForSort); j++ {
				if temes[i].ID == temesForSort[j].ID {
					temes[i].Rank = fmt.Sprint(j + 1)
					temes[i].RankNumber = j + 1
				}
			}
		}
	}

	getCompatitionModel := model.GetCompatition{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		Name:      result.Name,
		Sport:     result.Sport,
		Format:    result.Format,
		Type:      result.Type,
		OrganizerInfo: model.OrganizersInfo{
			ID:          result.OrganizersID,
			Name:        result.Organizers.Name,
			Phone:       result.Organizers.Phone,
			Description: result.Organizers.Description,
			Address: model.Address{
				HouseNumber: result.Organizers.Addresses.HouseNumber,
				Village:     result.Organizers.Addresses.Village,
				Subdistrict: result.Organizers.Addresses.Subdistrict,
				District:    result.Organizers.Addresses.District,
				PostalCode:  result.Organizers.Addresses.PostalCode,
				Country:     result.Organizers.Addresses.Country,
			},
			ImageProfilePath: result.Organizers.Users.ImageProfilePath,
			ImageCoverPath:   result.Organizers.Users.ImageCoverPath,
		},
		FieldSurface:    string(result.FieldSurface),
		ApplicationType: result.ApplicationType,
		Address: model.Address{
			HouseNumber: result.HouseNumber,
			Village:     result.Village,
			Subdistrict: result.Subdistrict,
			District:    result.District,
			PostalCode:  result.PostalCode,
			Country:     result.Country,
		},
		ImageBanner:          result.ImageBannerPath,
		StartDate:            result.StartDate,
		EndDate:              result.EndDate,
		JoinCode:             joinCode,
		Description:          result.Description,
		Rule:                 result.Rule,
		Prize:                result.Prize,
		ContactType:          result.ContactType,
		Contact:              result.Contact,
		AgeOver:              result.AgeOver,
		AgeUnder:             result.AgeUnder,
		Sex:                  result.Sex,
		Status:               string(result.Status),
		NumberOfTeam:         result.NumberOfTeam,
		NumOfPlayerInTeamMin: result.NumOfPlayerInTeamMin,
		NumOfPlayerInTeamMax: result.NumOfPlayerInTeamMax,
		Teams:                temes,
		NumOfRound:           result.NumOfRound,
		NumOfMatch:           result.NumOfMatch,
		Matchs:               Matchs,
	}

	return &getCompatitionModel, nil
}

// FinishCompetition implements CompetitionUsecase.
func (c *competitionUsecaseImpl) FinishCompetition(id uint, orgID uint) error {
	competition := &entities.Competitions{}
	competition.ID = id
	competition.OrganizersID = orgID
	competition, err := c.repository.GetCompetition(competition)
	if err != nil {
		return err
	}

	if competition.Status != util.CompetitionStatus[2] {
		return fmt.Errorf(`can't update compatition status to "Finished" (current status: %v)`, competition.Status)
	}

	if err := c.repository.UpdateSelectedFields(
		competition, "status",
		&entities.Competitions{
			OrganizersID: orgID,
			Status:       util.CompetitionStatus[3],
		},
	); err != nil {
		return err
	}
	return nil
}

// CancelCompatition implements CompetitionUsecase.
func (c *competitionUsecaseImpl) CancelCompatition(id uint, orgID uint) error {
	competition := &entities.Competitions{}
	competition.ID = id
	competition.OrganizersID = orgID
	competition, err := c.repository.GetCompetition(competition)
	if err != nil {
		return err
	}

	if competition.Status == util.CompetitionStatus[3] {
		return fmt.Errorf(`can't update compatition status to "Cancel" (current status: %v)`, competition.Status)
	}

	if competition.Status == util.CompetitionStatus[4] {
		return fmt.Errorf(`the current competition status is “Cancelled”`)
	}
	if err := c.repository.UpdateSelectedFields(
		competition,
		"status",
		&entities.Competitions{
			OrganizersID: orgID,
			Status:       util.CompetitionStatus[4],
		},
	); err != nil {
		return err
	}
	return nil
}

// OpenApplicationCompetition implements CompetitionUsecase.
func (c *competitionUsecaseImpl) OpenApplicationCompetition(id uint, orgID uint) error {
	competition := &entities.Competitions{}
	competition.ID = id
	competition.OrganizersID = orgID
	competition, err := c.repository.GetCompetition(competition)
	if err != nil {
		return err
	}

	if competition.Status != "Coming soon" {
		return fmt.Errorf("can't update compatition status to \"Applications opening\" (current status: %v)", competition.Status)
	}

	if err := c.repository.UpdateSelectedFields(
		competition,
		"status",
		&entities.Competitions{
			OrganizersID: orgID,
			Status:       util.CompetitionStatus[1],
		},
	); err != nil {
		return err
	}

	return nil
}

// StartCompetition implements CompetitionUsecase.
func (c *competitionUsecaseImpl) StartCompetition(id uint, orgID uint) error {
	competition := &entities.Competitions{}
	competition.ID = id
	competition.OrganizersID = orgID
	competition, err := c.repository.GetCompetitionDetails(competition)
	if err != nil {
		return err
	}

	if competition.Status != "Applications opening" {
		return fmt.Errorf("can't update compatition status to \"Stared\" (current status: %v)", competition.Status)
	}

	competition.Teams = shuffleTeam(competition.Teams)

	matchs := []entities.Matchs{}
	numOfRound := 0
	if competition.Type == "Tournament" {
		if util.CheckNumberPowerOfTwo(int(competition.NumberOfTeam)) != 0 {
			return errors.New("number of Team for create competition(tounament) is not power of 2")
		}
		if competition.NumberOfTeam < 2 {
			return errors.New("number of Team have to morn than 1")
		}
		numOfRound = int(math.Log2(float64(competition.NumberOfTeam)))
		count := 0

		for i := 0; i < numOfRound; i++ {
			round := numOfRound - i
			numOfMatchInRound := int(math.Pow(2, float64(round)) / 2)
			for j := 0; j < int(numOfMatchInRound); j++ {
				match := entities.Matchs{
					Round: fmt.Sprintf("Round %d", i+1),
				}
				if i != numOfRound-1 {
					if j%2 == 0 {
						match.NextMatchSlot = "Team1"
					} else {
						match.NextMatchSlot = "Team2"
					}
				}
				if i != 0 {
					matchs[count].NextMatchIndex = len(matchs) + 1
					matchs[count+1].NextMatchIndex = len(matchs) + 1
					count += 2
				} else {
					if j*2 < len(competition.Teams) {
						match.Team1ID = competition.Teams[j*2].TeamsID
					}
					if (j*2)+1 < len(competition.Teams) {
						match.Team2ID = competition.Teams[(j*2)+1].TeamsID
					}
				}
				match.Index = len(matchs) + 1
				match.DateTime = time.Date(
					0001, 01, 01, 0, 0, 0, 0, time.Local)
				matchs = append(matchs, match)
			}
		}
	} else if competition.Type == "Round Robin" {
		numOfRound = int(competition.NumberOfTeam - 1)
		matchs = roundRobin(int(competition.NumberOfTeam))
		for i := 0; i < len(matchs); i++ {
			if int(matchs[i].Team1ID) != 0 && int(matchs[i].Team1ID) <= len(competition.Teams) {
				matchs[i].Team1ID = competition.Teams[matchs[i].Team1ID-1].ID
			} else {
				matchs[i].Team1ID = 0
			}

			if int(matchs[i].Team2ID) != 0 && int(matchs[i].Team2ID) <= len(competition.Teams) {
				matchs[i].Team2ID = competition.Teams[matchs[i].Team2ID-1].ID
			} else {
				matchs[i].Team2ID = 0
			}
		}
	} else {
		return errors.New("undefined compatition type")
	}

	fmt.Printf("number of match %d\n", len(matchs))

	for _, v := range matchs {
		fmt.Printf("match %d. %s. Next match slot: %s, Next match: %v\n", v.Index, v.Round, v.NextMatchSlot, v.NextMatchIndex)
	}

	query := &entities.Competitions{}
	query.ID = id
	query.OrganizersID = orgID

	if err := c.repository.StartCompetitionAndAppendMatchToCompetition(
		query,
		&entities.Competitions{
			Status:     "Started",
			NumOfRound: numOfRound,
			NumOfMatch: len(matchs),
		},
		matchs,
	); err != nil {
		return err
	}

	return nil
}

// UpdateCompatition implements CompetitionUsecase.
func (c *competitionUsecaseImpl) UpdateCompatition(id uint, orgID uint, in *model.UpdateCompatition) error {

	query := &entities.Competitions{}
	query.ID = id
	query.OrganizersID = orgID
	competition, err := c.repository.GetCompetitionDetails(query)
	if err != nil {
		return err
	}

	if competition.Status != util.CompetitionStatus[0] {
		return fmt.Errorf(`competition can only be updated in "%s"(current status: %s)`, util.CompetitionStatus[0], competition.Status)
	}

	if err := c.repository.UpdateCompetition(query, &entities.Competitions{
		Name:                 in.Name,
		Sport:                in.Sport,
		Format:               in.Format,
		Type:                 in.Type,
		FieldSurface:         in.FieldSurface,
		ApplicationType:      in.ApplicationType,
		HouseNumber:          in.Address.HouseNumber,
		Village:              in.Address.Village,
		Subdistrict:          in.Address.Subdistrict,
		District:             in.Address.District,
		PostalCode:           in.Address.PostalCode,
		Country:              in.Address.Country,
		StartDate:            in.StartDate,
		EndDate:              in.EndDate,
		Description:          in.Description,
		Rule:                 in.Rule,
		Prize:                in.Prize,
		ContactType:          in.ContactType,
		Contact:              in.Contact,
		AgeOver:              in.AgeOver,
		AgeUnder:             in.AgeUnder,
		Sex:                  in.Sex,
		NumberOfTeam:         in.NumberOfTeam,
		NumOfPlayerInTeamMin: in.NumOfPlayerInTeamMin,
		NumOfPlayerInTeamMax: in.NumOfPlayerInTeamMax,
	}); err != nil {
		return err
	}
	return nil
}

// GetCompatitions implements CompetitionUsecase.
func (c *competitionUsecaseImpl) GetCompatitions(in *model.GetCompatitionsReq) ([]model.GetCompatitions, error) {

	competition := &entities.Competitions{
		OrganizersID: in.OrganizerID,
	}

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

	competitions, err := c.repository.GetCompetitions(competition, "id", true, limit, offset)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	competitionsModel := []model.GetCompatitions{}
	for _, v := range competitions {
		competitionsModel = append(competitionsModel, model.GetCompatitions{
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
			Status:          string(v.Status),
			Sex:             v.Sex,
			StartDate:       v.StartDate,
			EndDate:         v.EndDate,
			OrganizerID:     v.OrganizersID,
			OrganizerName:   v.Organizers.Name,
			ApplicationType: v.ApplicationType,
			AgeOver:         v.AgeOver,
			AgeUnder:        v.AgeUnder,
			ImageBanner:     v.ImageBannerPath,
		})
	}
	return competitionsModel, nil
}

// AddJoinCode implements CompetitionUsecase.
func (c *competitionUsecaseImpl) AddJoinCode(compatitionID uint, orgID uint, n int) error {

	competition := &entities.Competitions{}
	competition.ID = compatitionID
	competition.OrganizersID = orgID
	competition, err := c.repository.GetCompetition(competition)
	if err != nil {
		return err
	}

	if competition.Status == util.CompetitionStatus[2] ||
		competition.Status == util.CompetitionStatus[3] ||
		competition.Status == util.CompetitionStatus[4] {
		return fmt.Errorf(`competition can only join join code in "%s" or "%s" (current status: %s)`,
			util.CompetitionStatus[0], util.CompetitionStatus[1], competition.Status)
	}

	codes := []entities.JoinCode{}
	for i := 0; i < n; i++ {
		code := uuid.New().String()
		codes = append(codes, entities.JoinCode{
			CompetitionsID: compatitionID,
			Code:           code,
			Status:         util.JoinCodeStatus[0],
		})
	}

	if err := c.repository.AppendJoinCodeToCompetition(compatitionID, codes); err != nil {
		return err
	}
	return nil
}

// JoinCompetition implements CompetitionUsecase.
func (c *competitionUsecaseImpl) JoinCompetition(in *model.JoinCompetition, userID uint) error {
	util.PrintObjInJson(in)

	compatitionsEntity := &entities.Competitions{}
	teamEntity := &entities.Teams{}
	compatitionsEntity.ID = in.CompetitionID
	teamEntity.ID = in.TeamID

	compatition, err := c.repository.GetCompetitionDetails(compatitionsEntity)
	if err != nil {
		return err
	}

	team, err := c.repository.GetTeamWithAllAssociationsByID(teamEntity)
	if err != nil {
		return err
	}

	for _, team := range compatition.Teams {
		if team.ID == in.TeamID {
			return errors.New("this team has joined")
		}
	}

	if team.OwnerID != userID {
		return errors.New("only the owner of this team can use this team to join competiion")
	}

	// case 1: participatiing team are full
	if int(compatition.NumberOfTeam) <= len(compatition.Teams) {
		return errors.New("unable to join. the participating teams are full")
	}

	// case 2: applications don't opening
	if compatition.Status != "Applications opening" {
		return errors.New("unable to join. applications isn't opening")
	}

	validCode := false
	var joinCodeID uint

	if compatition.ApplicationType != util.ApplicationType[0] && compatition.ApplicationType != util.ApplicationType[1] {
		return errors.New("applicationType is unexpected")
	}

	// case 4: team member not enough
	if len(team.TeamsMembers) < int(compatition.NumOfPlayerInTeamMin) && compatition.NumOfPlayerInTeamMin != 0 {
		return errors.New("unable to join. your team does not have enough members")
	}

	// case 5: team member not enough
	if len(team.TeamsMembers) > int(compatition.NumOfPlayerInTeamMax) && compatition.NumOfPlayerInTeamMax != 0 {
		return errors.New("unable to join. your team has exceeded the maximum number of members")
	}

	normalUserIDs := []uint{}
	for _, member := range team.TeamsMembers {
		age := calculateAge(member.NormalUsers.Born)

		// case 6: age condition
		if (age < int(compatition.AgeOver) || age > int(compatition.AgeUnder)) && int(compatition.AgeOver) != 0 {
			return errors.New("unable to join. your team has older members. or lower than specified")
		}

		// case 7: gender condition
		if member.NormalUsers.Sex != string(compatition.Sex) && string(compatition.Sex) != "Unisex" {
			return errors.New("unable to join. your team has members whose sex does not match the gender assigned to the competition")
		}

		// case 8: has members who have entered this competition
		for _, teamJoined := range compatition.Teams {
			for _, teamJoinedMember := range teamJoined.Teams.TeamsMembers {
				if teamJoinedMember.NormalUsersID == member.NormalUsers.ID {
					return errors.New("unable to join. your team already has members who have entered this competition")
				}
			}
		}
		normalUserIDs = append(normalUserIDs, member.NormalUsersID)
	}

	if compatition.ApplicationType == util.ApplicationType[1] {

		// case 3: join code is invalid
		if in.Code == "" {
			return errors.New("unable to join. required code to join")
		}
		for i := 0; i < len(compatition.JoinCode); i++ {
			if compatition.JoinCode[i].Code == in.Code {
				if compatition.JoinCode[i].Status == util.JoinCodeStatus[1] {
					return errors.New("unable to join. join code is used")
				}
				validCode = true
				joinCodeID = compatition.JoinCode[i].ID
				break
			}
		}

		if !validCode {
			return errors.New("unable to join. code isn't valid")
		} else {
			err = c.repository.UpdateJoinCode(joinCodeID, &entities.JoinCode{
				Status: util.JoinCodeStatus[1],
			})
			if err != nil {
				return err
			}
		}
	}

	if err := c.repository.AddCompetitionToTeamAndNormalUsers(
		&entities.CompetitionsTeams{
			TeamsID:        team.ID,
			CompetitionsID: compatition.ID,
			Rank:           "0",
			RankNumber:     0,
			Point:          0,
		},
		normalUserIDs,
	); err != nil {
		// Rollback join code
		err = c.repository.UpdateJoinCode(joinCodeID, &entities.JoinCode{
			Status: util.JoinCodeStatus[0],
		})
		if err != nil {
			return err
		}
		return nil
	}

	// for _, member := range team.TeamsMembers {
	// 	err := c.repository.InsertNormalUserCompetition(
	// 		&entities.NormalUsersCompetitions{
	// 			NormalUsersID:  member.NormalUsersID,
	// 			CompetitionsID: in.CompetitionID,
	// 			TeamsID:        in.TeamID,
	// 		},
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// err = c.repository.InsertCompetitionsTeams(&entities.CompetitionsTeams{
	// 	TeamsID:        team.ID,
	// 	CompetitionsID: compatition.ID,
	// 	Rank:           "0",
	// 	RankNumber:     0,
	// 	Point:          0,
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================

func calculateAge(birthDate time.Time) int {
	now := time.Now()
	years := now.Year() - birthDate.Year()

	// Check if the birthday has occurred this year or not
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		years--
	}

	return years
}

func shuffleTeam(src []entities.CompetitionsTeams) []entities.CompetitionsTeams {
	dest := make([]entities.CompetitionsTeams, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}

func roundRobin(n int) []entities.Matchs {
	matchs := []entities.Matchs{}
	lst := make([]int, n-1)
	for i := 0; i < len(lst); i++ {
		lst[i] = i + 2
	}
	if n%2 == 1 {
		lst = append(lst, 0) // 0 denotes a bye
		n++
	}
	// count := 0
	for r := 1; r < n; r++ {
		fmt.Printf("Round %2d", r)
		lst2 := append([]int{1}, lst...)
		for i := 0; i < n/2; i++ {
			matchs = append(matchs, entities.Matchs{
				Index:   len(matchs) + 1,
				Round:   fmt.Sprintf("Round %d", r),
				Team1ID: uint(lst2[i]),
				Team2ID: uint(lst2[n-1-i]),
				DateTime: time.Date(
					0001, 01, 01, 0, 0, 0, 0, time.Local),
			})
			fmt.Printf(" (%2d vs %-2d)", lst2[i], lst2[n-1-i])
		}
		fmt.Println()
		rotate(lst)
	}
	return matchs
}

func rotate(lst []int) {
	len := len(lst)
	last := lst[len-1]
	for i := len - 1; i >= 1; i-- {
		lst[i] = lst[i-1]
	}
	lst[0] = last
}
