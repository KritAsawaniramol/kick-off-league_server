package normalUserUsecase

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
	"kickoff-league.com/util"
)

func NewNormalUserUsecaseImpl(
	repository repositories.Repository,
) NormalUserUsecase {
	return &normalUserUsecaseImpl{
		repository: repository,
	}
}

type normalUserUsecaseImpl struct {
	repository repositories.Repository
}

// GetNormalUserList implements NormalUserUsecase.
func (n *normalUserUsecaseImpl) GetNormalUserList() ([]model.NormalUserList, error) {
	normalUsers_entity, err := n.repository.GetNormalUsers(&entities.NormalUsers{})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	normalUserList := []model.NormalUserList{}
	for _, v := range normalUsers_entity {
		normalUserList = append(normalUserList, model.NormalUserList{
			ID:               v.ID,
			FirstNameThai:    v.FirstNameThai,
			LastNameThai:     v.LastNameThai,
			FirstNameEng:     v.FirstNameEng,
			LastNameEng:      v.LastNameEng,
			Born:             v.Born,
			Height:           v.Height,
			Weight:           v.Weight,
			Sex:              v.Sex,
			Position:         v.Position,
			Nationality:      v.Nationality,
			Description:      v.Description,
			ImageProfilePath: v.Users.ImageProfilePath,
			ImageCoverPath:   v.Users.ImageCoverPath,
		})
	}
	return normalUserList, nil
}

// GetNormalUser implements NormalUserUsecase.
func (n *normalUserUsecaseImpl) GetNormalUser(id uint) (*model.NormalUserProfile, error) {
	normalUserEntity := &entities.NormalUsers{}
	normalUserEntity.ID = id
	resultNormalUser, err := n.repository.GetNormalUserDetails(normalUserEntity)
	if err != nil {
		return nil, err
	}

	resultUser, err := n.repository.GetUserByID(resultNormalUser.UsersID)
	if err != nil {
		return nil, err
	}

	totalMatch := 0
	win := 0
	lose := 0
	recentMatch := []model.RecentMatch{}
	for _, compatition := range resultNormalUser.Competitions {

		teamID := compatition.TeamsID
		for _, match := range compatition.Competitions.Matchs {
			result := ""

			if match.Team1ID == teamID && teamID != 0 && match.Result != "" {
				totalMatch += 1
				if match.Result == util.MatchsResult[0] {
					win += 1
					result = "Win"
				} else if match.Result == util.MatchsResult[1] {
					lose += 1
					result = "Loss"
				}
				queryTeam2 := &entities.Teams{}
				queryTeam2.ID = match.Team2ID
				vsTeam, err := n.repository.GetTeam(queryTeam2)
				if err != nil {
					if err.Error() != "record not found" {
						return nil, err
					} else {
						vsTeam = &entities.Teams{
							Name: "",
						}
					}
				}
				recentMatch = append(recentMatch, model.RecentMatch{
					ID:             match.ID,
					DateTime:       match.DateTime,
					VsTeamName:     vsTeam.Name,
					Result:         result,
					Score:          fmt.Sprintf("%d - %d", match.Team1Goals, match.Team2Goals),
					CompatitionsID: match.CompetitionsID,
					TournamentName: compatition.Competitions.Name,
				})
			} else if match.Team2ID == teamID && teamID != 0 && match.Result != "" {
				totalMatch += 1
				if match.Result == util.MatchsResult[1] {
					win += 1
					result = "Win"
				} else if match.Result == util.MatchsResult[0] {
					lose += 1
					result = "Loss"
				}
				queryTeam1 := &entities.Teams{}
				queryTeam1.ID = match.Team1ID
				vsTeam, err := n.repository.GetTeam(queryTeam1)
				if err != nil {
					if err.Error() != "record not found" {
						return nil, err
					} else {
						vsTeam = &entities.Teams{
							Name: "",
						}
					}
				}
				recentMatch = append(recentMatch, model.RecentMatch{
					ID:             match.ID,
					DateTime:       match.DateTime,
					VsTeamName:     vsTeam.Name,
					Result:         result,
					Score:          fmt.Sprintf("%d - %d", match.Team2Goals, match.Team1Goals),
					TournamentName: compatition.Competitions.Name,
					CompatitionsID: match.CompetitionsID,
				})
			}
		}
	}
	winRate := (float64(win) / float64(totalMatch)) * 100
	goalPerCompatition := float64(len(resultNormalUser.GoalRecords)) / float64(totalMatch)

	// Handling NaN value
	if math.IsNaN(float64(goalPerCompatition)) {
		goalPerCompatition = 0
	}
	if math.IsNaN(float64(winRate)) {
		winRate = 0
	}

	// Custom sort function
	sortByMatchDateTime := func(i, j int) bool {
		return recentMatch[i].DateTime.After(recentMatch[j].DateTime)
	}
	// Sorting the array using custom sort function
	sort.Slice(recentMatch, sortByMatchDateTime)

	if len(recentMatch) > 20 {
		recentMatch = recentMatch[:20]
	}
	teamJoined := []model.Team{}
	for _, team := range resultNormalUser.Teams {
		teamJoined = append(teamJoined, model.Team{
			ID:             team.Teams.ID,
			Name:           team.Teams.Name,
			OwnerID:        team.Teams.OwnerID,
			Description:    team.Teams.Description,
			NumberOfMember: len(team.Teams.TeamsMembers),
			// Members:          members,
			ImageProfilePath: team.Teams.ImageProfilePath,
			ImageCoverPath:   team.Teams.ImageCoverPath,
		})
	}

	normalUserProfile := &model.NormalUserProfile{
		NormalUserInfo: model.NormalUserInfo{
			ID:            resultNormalUser.ID,
			FirstNameThai: resultNormalUser.FirstNameThai,
			LastNameThai:  resultNormalUser.LastNameThai,
			FirstNameEng:  resultNormalUser.FirstNameEng,
			LastNameEng:   resultNormalUser.LastNameEng,
			Username:      resultNormalUser.Username,
			Born:          resultNormalUser.Born,
			Phone:         resultNormalUser.Phone,
			Height:        resultNormalUser.Height,
			Weight:        resultNormalUser.Weight,
			Sex:           resultNormalUser.Sex,
			Position:      resultNormalUser.Position,
			Nationality:   resultNormalUser.Nationality,
			Description:   resultNormalUser.Description,

			Address: model.Address{
				HouseNumber: resultNormalUser.Addresses.HouseNumber,
				Village:     resultNormalUser.Addresses.Village,
				Subdistrict: resultNormalUser.Addresses.Subdistrict,
				District:    resultNormalUser.Addresses.District,
				PostalCode:  resultNormalUser.Addresses.PostalCode,
				Country:     resultNormalUser.Addresses.Country,
			},
		},

		UserID:           resultUser.ID,
		ImageProfilePath: resultUser.ImageProfilePath,
		ImageCoverPath:   resultUser.ImageCoverPath,
		NormalUserStat: model.NormalUserStat{
			WinRate:       winRate,
			TotalMatch:    totalMatch,
			Win:           win,
			Lose:          lose,
			Goals:         len(resultNormalUser.GoalRecords),
			GoalsPerMatch: goalPerCompatition,
			RecentMatch:   recentMatch,
		},
		TeamJoined: teamJoined,
	}
	return normalUserProfile, nil
}

// UpdateNormalUser implements NormalUserUsecase.
func (n *normalUserUsecaseImpl) UpdateNormalUser(inUpdateModel *model.UpdateNormalUser, inNormalUserID uint) error {
	normalUser := &entities.NormalUsers{
		FirstNameThai: inUpdateModel.FirstNameThai,
		LastNameThai:  inUpdateModel.LastNameThai,
		FirstNameEng:  inUpdateModel.FirstNameEng,
		LastNameEng:   inUpdateModel.LastNameEng,
		Born:          inUpdateModel.Born,
		Height:        inUpdateModel.Height,
		Weight:        inUpdateModel.Weight,
		Sex:           inUpdateModel.Sex,
		Position:      inUpdateModel.Position,
		Nationality:   inUpdateModel.Nationality,
		Description:   inUpdateModel.Description,
		Phone:         inUpdateModel.Phone,
		Username:      inUpdateModel.Username,
		Addresses: entities.Addresses{
			HouseNumber: inUpdateModel.Address.HouseNumber,
			Village:     inUpdateModel.Address.Village,
			Subdistrict: inUpdateModel.Address.Subdistrict,
			District:    inUpdateModel.Address.District,
			PostalCode:  inUpdateModel.Address.PostalCode,
			Country:     inUpdateModel.Address.Country,
		},
	}

	if isUsernameAlreadyInUser(normalUser.Username, n.repository) {
		return errors.New("error: this username is already in use")
	}

	normalUser.ID = inNormalUserID

	if err := n.repository.UpdateNormalUser(normalUser); err != nil {
		return err
	}
	return nil
}

// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================

func isUsernameAlreadyInUser(username string, u repositories.Repository) bool {
	if _, err := u.GetNormalUserByUsername(username); err != nil {
		return false
	}
	return true
}
