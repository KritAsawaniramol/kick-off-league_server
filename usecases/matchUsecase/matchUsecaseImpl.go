package matchUsecase

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
	"kickoff-league.com/util"
)

type matchUsecaseImpl struct {
	repository repositories.Repository
}

func NewMatchUsecaseImpl(
	repository repositories.Repository,
) MatchUsecase {
	return &matchUsecaseImpl{
		repository: repository,
	}
}

// GetMatch implements MatchUsecase.
func (m *matchUsecaseImpl) GetMatch(id uint) (*model.Match, error) {
	getMatch := &entities.Matchs{}
	getMatch.ID = id
	match, err := m.repository.GetMatchDetail(getMatch)
	if err != nil {
		return nil, err
	}
	team1 := &entities.Teams{}
	team2 := &entities.Teams{}
	team1.ID = match.Team1ID
	team2.ID = match.Team2ID
	team1, err = m.repository.GetTeam(team1)
	if err != nil {
		return nil, err
	}
	team2, err = m.repository.GetTeam(team2)
	if err != nil {
		return nil, err
	}

	team1PlayerInCompatition := []entities.NormalUsersCompetitions{}
	team2PlayerInCompatition := []entities.NormalUsersCompetitions{}

	if match.Team1ID != 0 {
		team1PlayerInCompatition, err = m.repository.GetNormalUserCompetitions(&entities.NormalUsersCompetitions{
			CompetitionsID: match.CompetitionsID,
			TeamsID:        match.Team1ID,
		})
		if err != nil {
			if err.Error() != "record not found" {
				return nil, err
			} else {
				team1PlayerInCompatition = []entities.NormalUsersCompetitions{}
			}
		}
	}

	if match.Team2ID != 0 {
		team2PlayerInCompatition, err = m.repository.GetNormalUserCompetitions(&entities.NormalUsersCompetitions{
			CompetitionsID: match.CompetitionsID,
			TeamsID:        match.Team2ID,
		})
		if err != nil {
			if err.Error() != "record not found" {
				return nil, err
			} else {
				team2PlayerInCompatition = []entities.NormalUsersCompetitions{}
			}
		}
	}

	team1PlayerModel := []model.Member{}
	team2PlayerModel := []model.Member{}

	for _, v := range team1PlayerInCompatition {
		team1PlayerModel = append(team1PlayerModel, model.Member{
			ID:            v.NormalUsers.ID,
			UsersID:       v.NormalUsers.UsersID,
			FirstNameThai: v.NormalUsers.FirstNameThai,
			LastNameThai:  v.NormalUsers.LastNameThai,
			FirstNameEng:  v.NormalUsers.FirstNameEng,
			LastNameEng:   v.NormalUsers.LastNameEng,
			Position:      v.NormalUsers.Position,
			Sex:           v.NormalUsers.Sex,
		})
	}
	for _, v := range team2PlayerInCompatition {
		team2PlayerModel = append(team2PlayerModel, model.Member{
			ID:            v.NormalUsers.ID,
			UsersID:       v.NormalUsers.UsersID,
			FirstNameThai: v.NormalUsers.FirstNameThai,
			LastNameThai:  v.NormalUsers.LastNameThai,
			FirstNameEng:  v.NormalUsers.FirstNameEng,
			LastNameEng:   v.NormalUsers.LastNameEng,
			Position:      v.NormalUsers.Position,
			Sex:           v.NormalUsers.Sex,
		})
	}

	goalRecord := []model.GoalRecord{}
	for _, v := range match.GoalRecords {
		goalRecord = append(goalRecord, model.GoalRecord{
			MatchsID:   v.MatchsID,
			TeamID:     v.TeamsID,
			PlayerID:   v.PlayerID,
			TimeScored: v.TimeScored,
		})
	}

	return &model.Match{
		ID:             match.ID,
		Index:          match.Index,
		CompatitionsID: match.CompetitionsID,
		GoalRecords:    goalRecord,
		DateTime:       match.DateTime,
		Team1ID:        match.Team1ID,
		Team2ID:        match.Team2ID,
		Team1Name:      team1.Name,
		Team2Name:      team2.Name,
		Team1Goals:     match.Team1Goals,
		Team2Goals:     match.Team2Goals,
		Round:          match.Round,
		NextMatchIndex: match.NextMatchIndex,
		NextMatchSlot:  match.NextMatchSlot,
		Result:         match.Result,
		VideoURL:       match.VideoURL,
		Team1Player:    team1PlayerModel,
		Team2Player:    team2PlayerModel,
	}, nil
}

// GetNextMatch implements MatchUsecase.
func (m *matchUsecaseImpl) GetNextMatch(id uint) ([]model.NextMatch, error) {
	nextMatchs := []model.NextMatch{}
	normalUser := &entities.NormalUsers{}
	normalUser.ID = id
	resultNormalUser, err := m.repository.GetNormalUserDetails(normalUser)
	if err != nil {
		return nil, err
	}
	for _, t := range resultNormalUser.Teams {
		team := &entities.Teams{}
		team.ID = t.TeamsID
		resultTeam, err := m.repository.GetTeamsWithCompetitionAndMatch(team)
		if err != nil {
			return nil, err
		}
		for _, compatition := range resultTeam.Competitions {
			if compatition.Competitions.Status == util.CompetitionStatus[2] {
				for _, match := range compatition.Competitions.Matchs {
					if match.Team1ID == t.ID && match.Team2ID != 0 && match.Result == "" {
						queryTeam2 := &entities.Teams{}
						queryTeam2.ID = match.Team2ID
						rivalTeam, err := m.repository.GetTeam(queryTeam2)
						if err != nil {
							return nil, err
						}
						nextMatchs = append(nextMatchs, model.NextMatch{
							RivalTeamID:           match.Team2ID,
							RivalTeamName:         rivalTeam.Name,
							RivalTeamImageProfile: rivalTeam.ImageProfilePath,
							RivalTeamImageCover:   rivalTeam.ImageCoverPath,
							CompatitionsID:        compatition.ID,
							CompatitionsName:      compatition.Competitions.Name,
							CompatitionsAddress: model.Address{
								HouseNumber: compatition.Competitions.HouseNumber,
								Village:     compatition.Competitions.Village,
								Subdistrict: compatition.Competitions.Subdistrict,
								District:    compatition.Competitions.District,
								PostalCode:  compatition.Competitions.PostalCode,
								Country:     compatition.Competitions.Country,
							},
							MatchID:       match.ID,
							MatchDateTime: match.DateTime,
						})
					} else if match.Team2ID == t.ID && match.Team1ID != 0 && match.Result == "" {
						queryTeam1 := &entities.Teams{}
						queryTeam1.ID = match.Team1ID
						rivalTeam, err := m.repository.GetTeam(queryTeam1)
						if err != nil {
							return nil, err
						}
						nextMatchs = append(nextMatchs, model.NextMatch{
							RivalTeamID:           match.Team1ID,
							RivalTeamName:         rivalTeam.Name,
							RivalTeamImageProfile: rivalTeam.ImageProfilePath,
							RivalTeamImageCover:   rivalTeam.ImageCoverPath,
							CompatitionsID:        compatition.ID,
							CompatitionsName:      compatition.Competitions.Name,
							CompatitionsAddress: model.Address{
								HouseNumber: compatition.Competitions.HouseNumber,
								Village:     compatition.Competitions.Village,
								Subdistrict: compatition.Competitions.Subdistrict,
								District:    compatition.Competitions.District,
								PostalCode:  compatition.Competitions.PostalCode,
								Country:     compatition.Competitions.Country,
							},
							MatchID:       match.ID,
							MatchDateTime: match.DateTime,
						})
					}
				}
			}
		}
	}

	// Custom sort function
	sortByMatchDateTime := func(i, j int) bool {
		if nextMatchs[i].MatchDateTime.Year() <= 1 {
			return false
		}
		return nextMatchs[i].MatchDateTime.Before(nextMatchs[j].MatchDateTime)
	}
	// Sorting the array using custom sort function
	sort.Slice(nextMatchs, sortByMatchDateTime)
	return nextMatchs, nil
}

// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================

// UpdateMatch implements MatchUsecase.
func (m *matchUsecaseImpl) UpdateMatch(id uint, orgID uint, updateMatch *model.UpdateMatch) error {
	util.PrintObjInJson(updateMatch)
	goalRecords := []entities.GoalRecords{}
	getMatch := &entities.Matchs{}
	getMatch.ID = id
	match, err := m.repository.GetMatchDetail(getMatch)
	if err != nil || match == nil {
		return err
	}

	getCompatition := &entities.Competitions{}
	getCompatition.ID = match.CompetitionsID
	compatition, err := m.repository.GetCompetitionDetails(getCompatition)
	if err != nil {
		return err
	}

	if compatition.OrganizersID != orgID {
		return errors.New("error: you can't update this match")
	}

	for _, goalRecord := range updateMatch.GoalRecords {
		goalRecords = append(goalRecords, entities.GoalRecords{
			MatchsID:   goalRecord.MatchsID,
			TeamsID:    goalRecord.TeamID,
			PlayerID:   goalRecord.PlayerID,
			TimeScored: goalRecord.TimeScored,
		})
	}

	match.DateTime = updateMatch.DateTime
	match.Team1Goals = updateMatch.Team1Goals
	match.Team2Goals = updateMatch.Team2Goals
	match.Result = updateMatch.Result
	match.VideoURL = updateMatch.VideoURL

	err = m.repository.UpdateMatch(id, match)
	if err != nil {
		return err
	}

	err = m.repository.ClearGoalRecordsOfMatch(id)
	if err != nil {
		return err
	}

	err = m.repository.AppendGoalRecordsToMatch(id, goalRecords)
	if err != nil {
		return err
	}
	if match.Competitions.Type == util.CompetitionType[0] {
		round, err := strconv.Atoi(strings.Split(match.Round, " ")[1])
		if err != nil {
			return err
		}

		if round == match.Competitions.NumOfRound {
			team1Rank := 0
			team2Rank := 0

			if updateMatch.Result == util.MatchsResult[0] {
				team1Rank = 1
				team2Rank = 2
			} else if updateMatch.Result == util.MatchsResult[1] {
				team1Rank = 2
				team2Rank = 1
			}

			err := m.repository.UpdateCompetitionsTeams(&entities.CompetitionsTeams{
				TeamsID:        match.Team1ID,
				CompetitionsID: match.CompetitionsID,
				Rank:           fmt.Sprint(team1Rank),
				RankNumber:     team1Rank,
			})
			if err != nil {
				return err
			}
			err = m.repository.UpdateCompetitionsTeams(&entities.CompetitionsTeams{
				TeamsID:        match.Team2ID,
				CompetitionsID: match.CompetitionsID,
				Rank:           fmt.Sprint(team2Rank),
				RankNumber:     team2Rank,
			})
			if err != nil {
				return err
			}
		} else {

			numberOfTeamInRound := int(match.Competitions.NumberOfTeam) / int(math.Pow(2, float64(round)-1))
			loserRank := fmt.Sprintf("%d-%d", (numberOfTeamInRound/2)+1, numberOfTeamInRound)
			if updateMatch.Result == util.MatchsResult[0] {
				err := m.repository.UpdateCompetitionsTeams(&entities.CompetitionsTeams{
					TeamsID:        match.Team2ID,
					CompetitionsID: match.CompetitionsID,
					Rank:           loserRank,
					RankNumber:     numberOfTeamInRound,
				})
				if err != nil {
					return err
				}

				err = m.repository.UpdateCompetitionsTeams(&entities.CompetitionsTeams{
					TeamsID:        match.Team1ID,
					CompetitionsID: match.CompetitionsID,
					Rank:           "",
					RankNumber:     0,
				})
				if err != nil {
					return err
				}

			} else if updateMatch.Result == util.MatchsResult[1] {
				err := m.repository.UpdateCompetitionsTeams(&entities.CompetitionsTeams{
					TeamsID:        match.Team1ID,
					CompetitionsID: match.CompetitionsID,
					Rank:           loserRank,
					RankNumber:     numberOfTeamInRound,
				})
				if err != nil {
					return err
				}

				err = m.repository.UpdateCompetitionsTeams(&entities.CompetitionsTeams{
					TeamsID:        match.Team2ID,
					CompetitionsID: match.CompetitionsID,
					Rank:           "",
					RankNumber:     0,
				})
				if err != nil {
					return err
				}
				// }
			}
		}
	} else if match.Competitions.Type == util.CompetitionType[1] {
		team1Point := 0
		team2Point := 0

		for _, m := range compatition.Matchs {

			if match.Team1ID == m.Team1ID && m.Result == util.MatchsResult[0] {
				team1Point += 3
			} else if match.Team1ID == m.Team1ID && m.Result == util.MatchsResult[2] {
				team1Point += 1
			} else if match.Team1ID == m.Team2ID && m.Result == util.MatchsResult[1] {
				team1Point += 3
			}

			if match.Team2ID == m.Team1ID && m.Result == util.MatchsResult[0] {
				team2Point += 3
			} else if match.Team2ID == m.Team1ID && m.Result == util.MatchsResult[2] {
				team2Point += 1
			} else if match.Team2ID == m.Team2ID && m.Result == util.MatchsResult[1] {
				team2Point += 3
			}
		}

		err = m.repository.UpdateCompetitionsTeams(&entities.CompetitionsTeams{
			TeamsID:        match.Team1ID,
			CompetitionsID: match.CompetitionsID,
			Point:          team1Point,
		})
		if err != nil {
			return err
		}
		err = m.repository.UpdateCompetitionsTeams(&entities.CompetitionsTeams{
			TeamsID:        match.Team2ID,
			CompetitionsID: match.CompetitionsID,
			Point:          team2Point,
		})
		if err != nil {
			return err
		}
	}

	// assign team to next match (if there are)
	fmt.Printf("match.NextMatchIndex: %v\n", match.NextMatchIndex)
	if match.NextMatchIndex != 0 {
		nextMatch, err := m.repository.GetMatchDetail(&entities.Matchs{
			CompetitionsID: match.CompetitionsID,
			Index:          match.NextMatchIndex,
		},
		)
		if err != nil {
			return err
		}

		if updateMatch.Result == util.MatchsResult[0] {
			if match.NextMatchSlot == "Team1" {
				nextMatch.Team1ID = match.Team1ID
			} else if match.NextMatchSlot == "Team2" {
				nextMatch.Team2ID = match.Team1ID
			}
		} else if updateMatch.Result == util.MatchsResult[1] {
			if match.NextMatchSlot == "Team1" {
				nextMatch.Team1ID = match.Team2ID
			} else if match.NextMatchSlot == "Team2" {
				nextMatch.Team2ID = match.Team2ID
			}
		}

		err = m.repository.UpdateMatch(nextMatch.ID, nextMatch)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
	}

	return nil
}
