package usecases

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/mail"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"kickoff-league.com/config"
	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
	"kickoff-league.com/util"
)

type userUsecaseImpl struct {
	userrepository repositories.Userrepository
}

// RemoveTeamFormCompatition implements UserUsecase.
func (u *userUsecaseImpl) RemoveTeamFormCompatition(teamID uint, compatitionID uint) error {
	err := u.userrepository.DeleteNormalUserCompatitionByTeamIDAndCompatitionID(compatitionID, teamID)
	if err != nil {
		return err
	}
	err = u.userrepository.DeleteCompatitionsTeam(compatitionID, teamID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrganizer implements UserUsecase.
func (u *userUsecaseImpl) UpdateOrganizer(orgID uint, in *model.UpdateOrganizer) error {
	getOrg := &entities.Organizers{}
	getOrg.ID = orgID
	org, err := u.userrepository.GetOrganizer(getOrg)
	if err != nil {
		return err
	}
	org.Name = in.Name
	org.Phone = in.Phone
	org.Description = in.Description
	org.Addresses.HouseNumber = in.Address.HouseNumber
	org.Addresses.Village = in.Address.Village
	org.Addresses.Subdistrict = in.Address.Subdistrict
	org.Addresses.District = in.Address.District
	org.Addresses.PostalCode = in.Address.PostalCode
	org.Addresses.Country = in.Address.Country

	err = u.userrepository.UpdateOrganizer(orgID, org)
	if err != nil {
		return err
	}
	return nil
}

// GetOrganizers implements UserUsecase.
func (u *userUsecaseImpl) GetOrganizers() ([]model.OrganizersInfo, error) {
	org, err := u.userrepository.GetOrganizers()
	if err != nil {
		if err.Error() != "record not found" {
			return nil, err
		} else {
			return nil, nil
		}
	}

	orgModel := []model.OrganizersInfo{}
	for _, v := range org {
		orgModel = append(orgModel, model.OrganizersInfo{
			ID:          v.ID,
			Name:        v.Name,
			Phone:       v.Phone,
			Description: v.Description,
			Address: model.Address{
				HouseNumber: v.Addresses.HouseNumber,
				Village:     v.Addresses.Village,
				Subdistrict: v.Addresses.Subdistrict,
				District:    v.Addresses.District,
				PostalCode:  v.Addresses.PostalCode,
				Country:     v.Addresses.Country,
			},
			ImageProfilePath: v.Users.ImageProfilePath,
			ImageCoverPath:   v.Users.ImageCoverPath,
		})
	}

	return orgModel, err
}

// GetOrganizer implements UserUsecase.
func (u *userUsecaseImpl) GetOrganizer(id uint) (*model.GetOrganizer, error) {
	getOrganizer := &entities.Organizers{}
	getOrganizer.ID = id
	org, err := u.userrepository.GetOrganizer(getOrganizer)
	if err != nil {
		if err.Error() != "record not found" {
			return nil, err
		} else {
			return nil, nil
		}
	}
	getCompatitions := []model.GetCompatitions{}
	for _, v := range org.Compatitions {
		getCompatitions = append(getCompatitions, model.GetCompatitions{
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
			Status:          v.Status,
			ApplicationType: v.ApplicationType,
			Sex:             v.Sex,
			StartDate:       v.StartDate,
			EndDate:         v.EndDate,
			OrganizerID:     org.ID,
			OrganizerName:   org.Name,
			AgeOver:         v.AgeOver,
			AgeUnder:        v.AgeUnder,
			ImageBanner:     v.ImageBannerPath,
		})
	}
	return &model.GetOrganizer{
		ID:          org.ID,
		Name:        org.Name,
		Phone:       org.Phone,
		Description: org.Description,
		Address: model.Address{
			HouseNumber: org.Addresses.HouseNumber,
			Village:     org.Addresses.Village,
			Subdistrict: org.Addresses.Subdistrict,
			District:    org.Addresses.District,
			PostalCode:  org.Addresses.PostalCode,
			Country:     org.Addresses.Country,
		},
		ImageProfilePath: org.Users.ImageProfilePath,
		ImageCoverPath:   org.Users.ImageCoverPath,
		Compatition:      getCompatitions,
	}, nil
}

// GetMatch implements UserUsecase.
func (u *userUsecaseImpl) GetMatch(id uint) (*model.Match, error) {
	getMatch := &entities.Matchs{}
	getMatch.ID = id
	match, err := u.userrepository.GetMatch(getMatch)
	if err != nil {
		return nil, err
	}
	team1, err := u.userrepository.GetTeam(match.Team1ID)
	if err != nil {
		return nil, err
	}
	team2, err := u.userrepository.GetTeam(match.Team2ID)
	if err != nil {
		return nil, err
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
		CompatitionsID: match.CompatitionsID,
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
	}, nil
}

// CreateJoinCode implements UserUsecase.
func (u *userUsecaseImpl) AddJoinCode(compatitionID uint, n int) error {
	codes := []entities.JoinCode{}
	for i := 0; i < n; i++ {
		code := uuid.New().String()
		codes = append(codes, entities.JoinCode{
			CompatitionsID: compatitionID,
			Code:           code,
			Status:         util.JoinCodeStatus[0],
		})
	}
	err := u.userrepository.AppendJoinCodeToCompatition(compatitionID, codes)
	if err != nil {
		return err
	}
	return nil
}

// RemoveNormalUserFormTeam implements UserUsecase.
func (u *userUsecaseImpl) RemoveNormalUserFormTeam(teamID uint, nomalUserID uint) error {
	err := u.userrepository.DeleteTeamMember(teamID, nomalUserID)
	if err != nil {
		return err
	}
	return nil
}

// GetNextMatch implements UserUsecase.
func (u *userUsecaseImpl) GetNextMatch(id uint) ([]model.NextMatch, error) {
	nextMatchs := []model.NextMatch{}
	normalUser := &entities.NormalUsers{}
	normalUser.ID = id
	resultNormalUser, err := u.userrepository.GetNormalUser(normalUser)
	if err != nil {
		return nil, err
	}
	for _, t := range resultNormalUser.Teams {
		team := &entities.Teams{}
		team.ID = t.TeamsID
		resultTeam, err := u.userrepository.GetTeamsWithCompatitionAndMatch(team)
		if err != nil {
			return nil, err
		}
		for _, compatition := range resultTeam.Compatitions {
			if compatition.Compatitions.Status == util.CompetitionStatus[2] {
				for _, match := range compatition.Compatitions.Matchs {
					if match.Team1ID == t.ID && match.Team2ID != 0 && match.Result == "" {
						rivalTeam, err := u.userrepository.GetTeam(match.Team2ID)
						if err != nil {
							return nil, err
						}
						nextMatchs = append(nextMatchs, model.NextMatch{
							RivalTeamID:           match.Team2ID,
							RivalTeamName:         rivalTeam.Name,
							RivalTeamImageProfile: rivalTeam.ImageProfilePath,
							RivalTeamImageCover:   rivalTeam.ImageCoverPath,
							CompatitionsID:        compatition.ID,
							CompatitionsName:      compatition.Compatitions.Name,
							CompatitionsAddress: model.Address{
								HouseNumber: compatition.Compatitions.HouseNumber,
								Village:     compatition.Compatitions.Village,
								Subdistrict: compatition.Compatitions.Subdistrict,
								District:    compatition.Compatitions.District,
								PostalCode:  compatition.Compatitions.PostalCode,
								Country:     compatition.Compatitions.Country,
							},
							MatchID:       match.ID,
							MatchDateTime: match.DateTime,
						})
					} else if match.Team2ID == t.ID && match.Team1ID != 0 && match.Result == "" {
						rivalTeam, err := u.userrepository.GetTeam(match.Team1ID)
						if err != nil {
							return nil, err
						}
						nextMatchs = append(nextMatchs, model.NextMatch{
							RivalTeamID:           match.Team1ID,
							RivalTeamName:         rivalTeam.Name,
							RivalTeamImageProfile: rivalTeam.ImageProfilePath,
							RivalTeamImageCover:   rivalTeam.ImageCoverPath,
							CompatitionsID:        compatition.ID,
							CompatitionsName:      compatition.Compatitions.Name,
							CompatitionsAddress: model.Address{
								HouseNumber: compatition.Compatitions.HouseNumber,
								Village:     compatition.Compatitions.Village,
								Subdistrict: compatition.Compatitions.Subdistrict,
								District:    compatition.Compatitions.District,
								PostalCode:  compatition.Compatitions.PostalCode,
								Country:     compatition.Compatitions.Country,
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
		return nextMatchs[i].MatchDateTime.After(nextMatchs[j].MatchDateTime)
	}
	// Sorting the array using custom sort function
	sort.Slice(nextMatchs, sortByMatchDateTime)
	return nextMatchs, nil
}

// UpdateMatch implements UserUsecase.
func (u *userUsecaseImpl) UpdateMatch(id uint, updateMatch *model.UpdateMatch) error {
	goalRecords := []entities.GoalRecords{}

	getMatch := &entities.Matchs{}
	getMatch.ID = id
	match, err := u.userrepository.GetMatch(getMatch)
	if err != nil || match == nil {
		fmt.Printf("err: %v\n", err)
		return err
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

	// err = u.userrepository.UpdateMatch(id, &entities.Matchs{
	// 	DateTime:   updateMatch.DateTime,
	// 	Team1Goals: updateMatch.Team1Goals,
	// 	Team2Goals: updateMatch.Team2Goals,
	// 	Result:     updateMatch.Result,
	// })
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	return err
	// }
	err = u.userrepository.UpdateMatch(id, match)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	err = u.userrepository.ClearGoalRecordsOfMatch(id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	err = u.userrepository.AppendGoalRecordsToMatch(id, goalRecords)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	if match.Compatitions.Type == util.CompetitionType[0] {
		round, err := strconv.Atoi(strings.Split(match.Round, " ")[1])
		if err != nil {
			return err
		}
		if round == 1 {
			team1Rank := 0
			team2Rank := 0
			if updateMatch.Result == util.MatchsResult[0] {
				team1Rank = 1
				team2Rank = 2
			} else if updateMatch.Result == util.MatchsResult[1] {
				team1Rank = 2
				team2Rank = 1
			}
			err := u.userrepository.UpdateCompatitionsTeams(&entities.CompatitionsTeams{
				TeamsID:        match.Team1ID,
				CompatitionsID: match.CompatitionsID,
				Rank:           fmt.Sprint(team1Rank),
				RankNumber:     team1Rank,
			})
			if err != nil {
				return err
			}
			err = u.userrepository.UpdateCompatitionsTeams(&entities.CompatitionsTeams{
				TeamsID:        match.Team2ID,
				CompatitionsID: match.CompatitionsID,
				Rank:           fmt.Sprint(team2Rank),
				RankNumber:     team2Rank,
			})
			if err != nil {
				return err
			}
		} else {
			numberOfTeamInRound := int(math.Pow(float64(round), 2))
			loserRank := fmt.Sprintf("%d-%d", (numberOfTeamInRound/2)+1, numberOfTeamInRound)
			if updateMatch.Result == util.MatchsResult[0] {
				err := u.userrepository.UpdateCompatitionsTeams(&entities.CompatitionsTeams{
					TeamsID:        match.Team2ID,
					CompatitionsID: match.CompatitionsID,
					Rank:           loserRank,
					RankNumber:     numberOfTeamInRound,
				})
				if err != nil {
					return err
				}

				err = u.userrepository.UpdateCompatitionsTeams(&entities.CompatitionsTeams{
					TeamsID:        match.Team1ID,
					CompatitionsID: match.CompatitionsID,
					Rank:           "",
					RankNumber:     0,
				})
				if err != nil {
					return err
				}

			} else if updateMatch.Result == util.MatchsResult[1] {
				if updateMatch.Result == util.MatchsResult[0] {
					err := u.userrepository.UpdateCompatitionsTeams(&entities.CompatitionsTeams{
						TeamsID:        match.Team1ID,
						CompatitionsID: match.CompatitionsID,
						Rank:           loserRank,
						RankNumber:     numberOfTeamInRound,
					})
					if err != nil {
						return err
					}

					err = u.userrepository.UpdateCompatitionsTeams(&entities.CompatitionsTeams{
						TeamsID:        match.Team2ID,
						CompatitionsID: match.CompatitionsID,
						Rank:           "",
						RankNumber:     0,
					})
					if err != nil {
						return err
					}
				}
			}
		}
	} else if match.Compatitions.Type == util.CompetitionType[1] {

		getCompatition := &entities.Compatitions{}
		getCompatition.ID = match.CompatitionsID
		team1Point := 0
		team2Point := 0
		compatition, err := u.userrepository.GetCompatition(getCompatition)
		if err != nil {
			return err
		}
		for _, m := range compatition.Matchs {
			if m.Team1ID == match.Team1ID {
				if m.Result == util.MatchsResult[0] {
					team1Point += 3
				} else if m.Result == util.MatchsResult[2] {
					team1Point += 1
				}
			} else if m.Team1ID == match.Team2ID {
				if m.Result == util.MatchsResult[1] {
					team1Point += 3
				} else if m.Result == util.MatchsResult[2] {
					team1Point += 1
				}
			}
			if m.Team2ID == match.Team1ID {
				if m.Result == util.MatchsResult[0] {
					team2Point += 3
				} else if m.Result == util.MatchsResult[2] {
					team2Point += 1
				}
			} else if m.Team2ID == match.Team2ID {
				if m.Result == util.MatchsResult[1] {
					team2Point += 3
				} else if m.Result == util.MatchsResult[2] {
					team2Point += 1
				}
			}
		}

		fmt.Printf("team2Point: %v\n", team2Point)

		err = u.userrepository.UpdateCompatitionsTeams(&entities.CompatitionsTeams{
			TeamsID:        match.Team1ID,
			CompatitionsID: match.CompatitionsID,
			Point:          team1Point,
		})
		if err != nil {
			return err
		}
		err = u.userrepository.UpdateCompatitionsTeams(&entities.CompatitionsTeams{
			TeamsID:        match.Team2ID,
			CompatitionsID: match.CompatitionsID,
			Point:          team2Point,
		})
		if err != nil {
			return err
		}
	}

	// assign team to next match (if there are)
	if match.NextMatchIndex != 0 {
		nextMatch, err := u.userrepository.GetMatch(&entities.Matchs{Index: match.NextMatchIndex})
		if err != nil || match == nil {
			fmt.Printf("err: %v\n", err)
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

		err = u.userrepository.UpdateMatch(nextMatch.ID, nextMatch)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
	}

	return nil
}

// UpdateCompatition implements UserUsecase.
func (u *userUsecaseImpl) UpdateCompatition(id uint, in *model.UpdateCompatition) error {
	err := u.userrepository.UpdateCompatition(id, &entities.Compatitions{
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
	})
	if err != nil {
		return err
	}
	return nil
}

// JoinCompatition implements UserUsecase.
func (u *userUsecaseImpl) JoinCompatition(in *model.JoinCompatition) error {
	compatitionsEntity := &entities.Compatitions{}
	teamEntity := &entities.Teams{}
	compatitionsEntity.ID = in.CompatitionID
	teamEntity.ID = in.TeamID
	compatition, err := u.userrepository.GetCompatition(compatitionsEntity)
	if err != nil {
		return err
	}

	if int(compatition.NumberOfTeam) <= len(compatition.Teams) {
		return errors.New("unable to join. the participating team is full")
	}

	if compatition.Status != "Applications opening" {
		return errors.New("unable to join. applications isn't opening")
	}

	validCode := false
	var joinCodeID uint

	if compatition.ApplicationType == util.ApplicationType[1] {
		if in.Code == "" {
			return errors.New("unable to join. required code to join")
		}
		for i := 0; i < len(compatition.JoinCode); i++ {
			if compatition.JoinCode[i].Code == in.Code {
				validCode = true
				if compatition.JoinCode[i].Status == util.JoinCodeStatus[1] {
					return errors.New("unable to join. join code is used")
				}
				joinCodeID = compatition.JoinCode[i].ID
				break
			}
		}

		if !validCode {
			return errors.New("unable to join. code isn't valid")
		}
	}

	team, err := u.userrepository.GetTeamWithAllAssociationsByID(teamEntity)
	if err != nil {
		return err
	}

	if len(team.TeamsMembers) < int(compatition.NumOfPlayerInTeamMin) && compatition.NumOfPlayerInTeamMin != 0 {
		return errors.New("unable to join. your team does not have enough members")
	}

	if len(team.TeamsMembers) > int(compatition.NumOfPlayerInTeamMax) && compatition.NumOfPlayerInTeamMax != 0 {
		return errors.New("uunable to join. your team has exceeded the maximum number of members")
	}

	for _, member := range team.TeamsMembers {
		age := calculateAge(member.NormalUsers.Born)
		if age < int(compatition.AgeOver) && age != 0 {
			return errors.New("unable to join. your team has older members. or lower than specified")
		}

		if age > int(compatition.AgeUnder) && age != 0 {
			return errors.New("unable to join. your team has older members. or lower than specified")
		}

		if member.NormalUsers.Sex != string(compatition.Sex) && string(compatition.Sex) != "Unisex" {
			return errors.New("unable to join. your team has members whose sex does not match the gender assigned to the competition")
		}

		for _, teamJoined := range compatition.Teams {
			for _, teamJoinedMember := range teamJoined.Teams.TeamsMembers {
				if teamJoinedMember.NormalUsersID == member.NormalUsers.ID {
					return errors.New("unable to join. your team already has members who have entered this competition")
				}
			}
		}
	}

	for _, member := range team.TeamsMembers {
		err := u.userrepository.InsertNormalUserCompatition(
			&entities.NormalUsersCompatitions{
				NormalUsersID:  member.NormalUsersID,
				CompatitionsID: in.CompatitionID,
				TeamsID:        in.TeamID,
			},
		)
		if err != nil {
			return err
		}
	}

	if compatition.ApplicationType == util.ApplicationType[1] && joinCodeID != 0 && validCode {
		err = u.userrepository.UpdateJoinCode(joinCodeID, &entities.JoinCode{
			Status: util.JoinCodeStatus[1],
		})
		if err != nil {
			return err
		}
	}

	// err = u.userrepository.AppendTeamtoCompatition(compatition, team)
	// if err != nil {
	// 	return err
	// }

	err = u.userrepository.InsertCompatitionsTeams(&entities.CompatitionsTeams{
		TeamsID:        team.ID,
		CompatitionsID: compatition.ID,
		Rank:           "0",
		RankNumber:     0,
		Point:          0,
	})
	if err != nil {
		return err
	}

	return nil
}

func calculateAge(birthDate time.Time) int {
	now := time.Now()
	years := now.Year() - birthDate.Year()

	// Check if the birthday has occurred this year or not
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		years--
	}

	return years
}

// GetNormalUserList implements UserUsecase.
func (u *userUsecaseImpl) GetNormalUserList() ([]model.NormalUserList, error) {
	normalUsers_entity, err := u.userrepository.GetNormalUsers(&entities.NormalUsers{})
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

// UpdateTeamImageCover implements UserUsecase.
func (u *userUsecaseImpl) UpdateTeamImageCover(teamID uint, newImagePath string) error {
	team := &entities.Compatitions{}
	team.ID = teamID
	if err := u.userrepository.UpdateSelectedFields(team, "ImageCoverPath", &entities.Teams{ImageCoverPath: newImagePath}); err != nil {
		return err
	}
	return nil
}

// UpdateTeamImageProfile implements UserUsecase.
func (u *userUsecaseImpl) UpdateTeamImageProfile(teamID uint, newImagePath string) error {
	team := &entities.Compatitions{}
	team.ID = teamID
	if err := u.userrepository.UpdateSelectedFields(team, "ImageProfilePath", &entities.Teams{ImageProfilePath: newImagePath}); err != nil {
		return err
	}
	return nil
}

// UpdateImageBanner implements UserUsecase.
func (u *userUsecaseImpl) UpdateImageBanner(compatitionID uint, newImagePath string) error {
	compatition := &entities.Compatitions{}
	compatition.ID = compatitionID
	if err := u.userrepository.UpdateSelectedFields(compatition, "ImageBannerPath", &entities.Compatitions{ImageBannerPath: newImagePath}); err != nil {
		return err
	}
	return nil
}

// RemoveImageProfile implements UserUsecase.
func (u *userUsecaseImpl) UpdateImageProfile(userID uint, newImagePath string) error {
	user := &entities.Users{}
	user.ID = userID
	util.PrintObjInJson(user)
	if err := u.userrepository.UpdateSelectedFields(user, "ImageProfilePath", &entities.Users{ImageProfilePath: newImagePath}); err != nil {
		return err
	}
	return nil
}

// RemoveImageProfile implements UserUsecase.
func (u *userUsecaseImpl) UpdateImageCover(userID uint, newImagePath string) error {
	user := &entities.Users{}
	user.ID = userID
	if err := u.userrepository.UpdateSelectedFields(user, "ImageCoverPath", &entities.Users{ImageCoverPath: newImagePath}); err != nil {
		return err
	}
	return nil
}

// RemoveImageProfile implements UserUsecase.
func (u *userUsecaseImpl) RemoveImageProfile(userID uint) error {
	user := &entities.Users{}
	user.ID = userID
	if err := u.userrepository.UpdateSelectedFields(user, "ImageProfilePath", &entities.Users{ImageProfilePath: "./images/default/defaultProfile.jpg"}); err != nil {
		return err
	}
	return nil
}

// RemoveImageProfile implements UserUsecase.
func (u *userUsecaseImpl) RemoveImageCover(userID uint) error {
	user := &entities.Users{}
	user.ID = userID
	if err := u.userrepository.UpdateSelectedFields(user, "ImageCoverPath", &entities.Users{ImageCoverPath: "./images/default/defaultCover.jpg"}); err != nil {
		return err
	}
	return nil
}

// RemoveTeamImageCover implements UserUsecase.
func (u *userUsecaseImpl) RemoveTeamImageCover(teamID uint) error {
	team := &entities.Teams{}
	team.ID = teamID
	if err := u.userrepository.UpdateSelectedFields(team, "ImageCoverPath", &entities.Teams{ImageCoverPath: "./images/default/defaultCover.jpg"}); err != nil {
		return err
	}
	return nil
}

// RemoveTeamImageProfile implements UserUsecase.
func (u *userUsecaseImpl) RemoveTeamImageProfile(teamID uint) error {
	team := &entities.Teams{}
	team.ID = teamID
	if err := u.userrepository.UpdateSelectedFields(team, "ImageProfilePath", &entities.Users{ImageProfilePath: "./images/default/defaultProfile.jpg"}); err != nil {
		return err
	}
	return nil
}

func (u *userUsecaseImpl) RemoveImageBanner(compatitionID uint) error {
	compatition := &entities.Compatitions{}
	compatition.ID = compatitionID
	if err := u.userrepository.UpdateSelectedFields(compatition, "ImageBannerPath", &entities.Compatitions{ImageBannerPath: "./images/default/defaultBanner.png"}); err != nil {
		return err
	}
	return nil
}

// GetCompatitions implements UserUsecase.
func (u *userUsecaseImpl) GetCompatitions(in *model.GetCompatitionsReq) ([]model.GetCompatitions, error) {

	compatition := &entities.Compatitions{
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

	compatitions, err := u.userrepository.GetCompatitions(compatition, "id", true, limit, offset)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	compatitionsModel := []model.GetCompatitions{}
	for _, v := range compatitions {
		compatitionsModel = append(compatitionsModel, model.GetCompatitions{
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
	return compatitionsModel, nil
}

type ByPoint []model.Team

func (o ByPoint) Len() int      { return len(o) }
func (o ByPoint) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o ByPoint) Less(i, j int) bool {
	if o[i].Point == o[j].Point {
		return (o[i].GoalsScored - o[i].GoalsConceded) < (o[j].GoalsScored - o[j].GoalsConceded)
	} else {
		return o[i].Point < o[j].Point
	}
}

// GetCompatition implements UserUsecase.
func (u *userUsecaseImpl) GetCompatition(in uint) (*model.GetCompatition, error) {
	compatitionEntity := &entities.Compatitions{}
	compatitionEntity.ID = in
	result, err := u.userrepository.GetCompatition(compatitionEntity)
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
		members := []model.Member{}
		for _, member := range v.Teams.TeamsMembers {
			members = append(members, model.Member{
				ID:            member.ID,
				UsersID:       member.NormalUsers.UsersID,
				FirstNameThai: member.NormalUsers.FirstNameThai,
				LastNameThai:  member.NormalUsers.LastNameThai,
				FirstNameEng:  member.NormalUsers.FirstNameEng,
				LastNameEng:   member.NormalUsers.LastNameEng,
				Position:      member.NormalUsers.Position,
				Sex:           member.NormalUsers.Sex,
			})
		}

		temes = append(temes, model.Team{
			ID:               v.ID,
			Name:             v.Teams.Name,
			OwnerID:          v.Teams.OwnerID,
			Description:      v.Teams.Description,
			Members:          members,
			Rank:             v.Rank,
			RankNumber:       v.RankNumber,
			Point:            v.Point,
			ImageProfilePath: v.Teams.ImageProfilePath,
			ImageCoverPath:   v.Teams.ImageCoverPath,
		})
		goalsScored[v.TeamsID] = 0
		goalsConceded[v.TeamsID] = 0
	}

	temesForSort := temes
	if result.Type == util.CompetitionType[1] {
		sort.Sort(ByPoint(temesForSort))
		for i := 0; i < len(temesForSort); i++ {
			temes[i].Rank = fmt.Sprint(i + 1)
			temes[i].RankNumber = i + 1
		}
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

		Matchs = append(Matchs, model.Match{
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
		})

		goalsScored[v.Team1ID] += v.Team1Goals
		goalsConceded[v.Team1ID] += v.Team2Goals

		goalsScored[v.Team2ID] += v.Team2Goals
		goalsConceded[v.Team2ID] += v.Team1Goals
	}

	for _, t := range temes {
		t.GoalsScored = goalsScored[t.ID]
		t.GoalsConceded = goalsScored[t.ID]
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

// CreateCompatition implements UserUsecase.
func (u *userUsecaseImpl) CreateCompatition(in *model.CreateCompatition) error {

	compatition := &entities.Compatitions{
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

	// matchs := []entities.Matchs{}
	// numOfRound := 0
	if in.Type == "Tournament" {
		if checkNumberPowerOfTwo(int(in.NumberOfTeam)) != 0 {
			return errors.New("number of Team for create competition(tounament) is not power of 2")
		}
		if in.NumberOfTeam < 2 {
			return errors.New("number of Team have to morn than 1")
		}

	}

	if in.Type != "Tournament" && in.Type != "Round Robin" {
		return errors.New("undefined compatition type")
	}
	if err := u.userrepository.InsertCompatition(compatition); err != nil {
		return err
	}

	return nil
}

// CancelCompatition implements UserUsecase.
func (u *userUsecaseImpl) CancelCompatition(id uint) error {
	err := u.userrepository.UpdateCompatition(id, &entities.Compatitions{
		Status: util.CompetitionStatus[4],
	})
	if err != nil {
		return err
	}
	return nil
}

// FinishCompatition implements UserUsecase.
func (u *userUsecaseImpl) FinishCompatition(id uint) error {
	compatition := &entities.Compatitions{}
	compatition.ID = id
	compatition, err := u.userrepository.GetCompatition(compatition)
	if err != nil {
		return err
	}

	if compatition.Status != util.CompetitionStatus[2] {
		return fmt.Errorf("can't update compatition status to \"Finished\" (status now: %v)", compatition.Status)
	}

	err = u.userrepository.UpdateCompatition(id, &entities.Compatitions{
		Status: util.CompetitionStatus[3],
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecaseImpl) OpenApplicationCompatition(id uint) error {
	compatition := &entities.Compatitions{}
	compatition.ID = id
	compatition, err := u.userrepository.GetCompatition(compatition)
	if err != nil {
		return err
	}

	if compatition.Status != "Coming soon" {
		return fmt.Errorf("can't update compatition status to \"Applications opening\" (status now: %v)", compatition.Status)
	}

	err = u.userrepository.UpdateCompatition(id, &entities.Compatitions{
		Status: "Applications opening",
	})
	if err != nil {
		return err
	}
	return nil
}

// StartCompatition implements UserUsecase.
func (u *userUsecaseImpl) StartCompatition(id uint) error {
	compatition := &entities.Compatitions{}
	compatition.ID = id
	compatition, err := u.userrepository.GetCompatition(compatition)
	if err != nil {
		return err
	}

	if compatition.Status != "Applications opening" {
		return fmt.Errorf("can't update compatition status to \"Stared\" (status now: %v)", compatition.Status)
	}

	compatition.Teams = shuffleTeam(compatition.Teams)

	matchs := []entities.Matchs{}
	numOfRound := 0
	if compatition.Type == "Tournament" {
		if checkNumberPowerOfTwo(int(compatition.NumberOfTeam)) != 0 {
			return errors.New("number of Team for create competition(tounament) is not power of 2")
		}
		if compatition.NumberOfTeam < 2 {
			return errors.New("number of Team have to morn than 1")
		}
		numOfRound = int(math.Log2(float64(compatition.NumberOfTeam)))
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
					if j*2 < len(compatition.Teams) {
						match.Team1ID = compatition.Teams[j*2].ID
					}
					if (j*2)+1 < len(compatition.Teams) {
						match.Team2ID = compatition.Teams[(j*2)+1].ID
					}
				}
				match.Index = len(matchs) + 1
				match.DateTime = time.Date(
					0001, 01, 01, 0, 0, 0, 0, time.Local)
				matchs = append(matchs, match)
			}
		}
	} else if compatition.Type == "Round Robin" {
		numOfRound = int(compatition.NumberOfTeam - 1)
		// numOfMatch := (int(compatition.NumberOfTeam) * numOfRound) / 2
		// numOfMatchInRound := numOfMatch / numOfRound
		matchs = roundRobin(int(compatition.NumberOfTeam))
		for i := 0; i < len(matchs); i++ {
			if int(matchs[i].Team1ID) != 0 && int(matchs[i].Team1ID) <= len(compatition.Teams) {
				matchs[i].Team1ID = compatition.Teams[matchs[i].Team1ID-1].ID
			} else {
				matchs[i].Team1ID = 0
			}

			if int(matchs[i].Team2ID) != 0 && int(matchs[i].Team2ID) <= len(compatition.Teams) {
				matchs[i].Team2ID = compatition.Teams[matchs[i].Team2ID-1].ID
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

	// compatition.Matchs = matchs
	// compatition.NumOfRound = numOfRound
	// compatition.NumOfMatch = len(matchs)

	u.userrepository.AppendMatchToCompatition(compatition, matchs)
	err = u.userrepository.UpdateCompatition(id, &entities.Compatitions{
		Status:     "Started",
		NumOfRound: numOfRound,
		NumOfMatch: len(matchs),
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	return nil
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

func checkNumberPowerOfTwo(n int) int {
	return n & (n - 1)
}

//	func shuffleTeam(src []entities.Teams) []entities.Teams {
//		dest := make([]entities.Teams, len(src))
//		perm := rand.Perm(len(src))
//		for i, v := range perm {
//			dest[v] = src[i]
//		}
//		return dest
//	}
func shuffleTeam(src []entities.CompatitionsTeams) []entities.CompatitionsTeams {
	dest := make([]entities.CompatitionsTeams, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}

func rotate(lst []int) {
	len := len(lst)
	last := lst[len-1]
	for i := len - 1; i >= 1; i-- {
		lst[i] = lst[i-1]
	}
	lst[0] = last
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// GetMyPenddingAddMemberRequest implements UserUsecase.
func (u *userUsecaseImpl) GetMyPenddingAddMemberRequest(userID uint) ([]model.AddMemberRequest, error) {
	addMemberRequestSearch := &entities.AddMemberRequests{
		ReceiverID: userID,
		Status:     "pending",
	}
	addMemberRequestList, err := u.userrepository.GetAddMemberRequestByID(addMemberRequestSearch)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	addMemberRequestModelList := []model.AddMemberRequest{}
	for _, v := range addMemberRequestList {
		addMemberRequestModelList = append(addMemberRequestModelList, model.AddMemberRequest{
			ID:               v.ID,
			TeamID:           v.TeamsID,
			TeamName:         v.Teams.Name,
			TeamImageProfile: v.Teams.ImageProfilePath,
			Role:             v.Role,
			Status:           v.Status,
		})
	}

	return addMemberRequestModelList, nil
}

// GetTeamMember implements UserUsecase.
func (u *userUsecaseImpl) GetTeamMembers(id uint) (*model.Team, error) {
	// teamsMembers, err := u.userrepository.GetTeamMembersByTeamID(id, "id", false, -1, -1)
	// if err != nil {
	// 	return nil, err
	// }
	// util.PrintObjInJson(teamsMembers)
	return nil, nil
}

// GetTeamWithMemberAndCompatitionByID implements UserUsecase.
func (u *userUsecaseImpl) GetTeamWithMemberAndCompatitionByID(id uint) (*model.Team, error) {
	t := &entities.Teams{}
	t.ID = id
	selectedTeams, err := u.userrepository.GetTeamWithMemberAndCompatitionByID(id)
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

	for _, v := range selectedTeams.Compatitions {
		compatition_model = append(compatition_model, model.CompatitionBasicInfo{
			ID:           v.ID,
			Name:         v.Compatitions.Name,
			Format:       v.Compatitions.Format,
			OrganizerID:  v.Compatitions.OrganizersID,
			StartDate:    v.Compatitions.StartDate,
			EndDate:      v.Compatitions.EndDate,
			AgeOver:      v.Compatitions.AgeOver,
			AgeUnder:     v.Compatitions.AgeUnder,
			Sex:          v.Compatitions.Sex,
			FieldSurface: v.Compatitions.FieldSurface,
			Description:  v.Compatitions.Description,
			Status:       v.Compatitions.Status,
			NumberOfTeam: v.Compatitions.NumberOfTeam,
			ImageBanner:  v.Compatitions.ImageBannerPath,
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

// GetTeamList implements UserUsecase.
func (u *userUsecaseImpl) GetTeams(in *model.GetTeamsReq) ([]model.TeamList, error) {

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

	teams, err := u.userrepository.GetTeams(&team, strings.Trim(in.Ordering, " "), in.Decs, limit, offset)
	if err != nil {
		return []model.TeamList{}, err
	}

	teamList := []model.TeamList{}
	for _, team := range teams {
		teamList = append(teamList, model.TeamList{
			ID:               team.ID,
			Name:             team.Name,
			Description:      team.Description,
			NumberOfMember:   uint(u.userrepository.GetNumberOfTeamsMember(team.ID)),
			ImageProfilePath: team.ImageProfilePath,
			ImageCoverPath:   team.ImageCoverPath,
		})
	}
	return teamList, nil
}

// GetTeamList implements UserUsecase.
func (u *userUsecaseImpl) GetTeamsByOwnerID(in uint) ([]model.TeamList, error) {

	team := entities.Teams{
		// 0 is select all id
		OwnerID: in,
	}

	teams, err := u.userrepository.GetTeams(&team, "id", false, -1, -1)
	if err != nil {
		return []model.TeamList{}, err
	}

	teamList := []model.TeamList{}
	for _, team := range teams {
		teamList = append(teamList, model.TeamList{
			ID:               team.ID,
			Name:             team.Name,
			Description:      team.Description,
			NumberOfMember:   uint(u.userrepository.GetNumberOfTeamsMember(team.ID)),
			OwnerID:          team.OwnerID,
			ImageProfilePath: team.ImageProfilePath,
			ImageCoverPath:   team.ImageCoverPath,
		})
	}
	return teamList, nil
}

// GetMyTeam implements UserUsecase.
func (*userUsecaseImpl) GetMyTeam() ([]model.Team, error) {
	panic("unimplemented")
}

// IgnoreAddMemberRequest implements UserUsecase.
func (u *userUsecaseImpl) IgnoreAddMemberRequest(inReqID uint, userID uint) error {
	addMemberRequestSearch := &entities.AddMemberRequests{}
	addMemberRequestSearch.ID = inReqID
	addMemberRequestList, err := u.userrepository.GetAddMemberRequestByID(addMemberRequestSearch)
	if err != nil {
		return err
	}
	addMemberRequest := &addMemberRequestList[0]
	if addMemberRequest.ReceiverID != userID {
		return errors.New("UserID does not match ReceiverID")
	}
	if err := u.userrepository.UpdateAddMemberRequestStatusAndSoftDelete(addMemberRequest, "ignore"); err != nil {
		return err
	}
	return nil
}

// AcceptAddMemberRequest implements UserUsecase.
func (u *userUsecaseImpl) AcceptAddMemberRequest(inReqID uint, userID uint) error {
	addMemberRequestSearch := &entities.AddMemberRequests{}
	addMemberRequestSearch.ID = inReqID
	addMemberRequestList, err := u.userrepository.GetAddMemberRequestByID(addMemberRequestSearch)
	if err != nil {
		return err
	}
	addMemberRequest := &addMemberRequestList[0]

	if addMemberRequest.ReceiverID != userID {
		log.Printf("ReceiverID: %d, userID: %d\n", addMemberRequest.ReceiverID, userID)
		return errors.New("UserID does not match ReceiverID")
	}

	if err := u.userrepository.UpdateAddMemberRequestStatusByID(inReqID, "accepted"); err != nil {
		return err
	}

	normalUser, err := u.userrepository.GetNormalUser(&entities.NormalUsers{UsersID: userID})
	if err != nil {
		return err
	}

	if err := u.userrepository.InsertTeamsMembers(&entities.TeamsMembers{
		TeamsID:       addMemberRequest.TeamsID,
		NormalUsersID: normalUser.ID,
		Role:          addMemberRequest.Role,
	}); err != nil {
		return err
	}

	return nil
}

// CreateAddMemberRequest implements UserUsecase.
func (u *userUsecaseImpl) SendAddMemberRequest(inAddMemberRequest *model.AddMemberRequest, inUserID uint) error {
	// Get Receiver normaluser
	receiver, err := u.userrepository.GetNormalUser(&entities.NormalUsers{
		Username: inAddMemberRequest.ReceiverUsername,
	})
	if err != nil {
		return err
	}

	team, err := u.userrepository.GetTeamWithMemberAndRequestSendByID(inAddMemberRequest.TeamID)
	if err != nil {
		return err
	}

	if team.OwnerID != inUserID {
		return errors.New("this user isn't owner's team")
	}

	for _, member := range team.TeamsMembers {
		fmt.Println(member.NormalUsers.UsersID)
		fmt.Println(receiver.UsersID)
		if member.NormalUsersID == receiver.ID {
			return errors.New("this user already in team")
		}
	}

	for _, requestSend := range team.RequestSends {
		if requestSend.ReceiverID == receiver.UsersID && requestSend.Status == "pending" {
			return errors.New("this user already invited")
		}
	}

	addMemberRequest := &entities.AddMemberRequests{
		TeamsID:    inAddMemberRequest.TeamID,
		ReceiverID: receiver.UsersID,
		Status:     "pending",
		Role:       inAddMemberRequest.Role,
	}

	// Create Request
	if err := u.userrepository.InsertAddMemberRequest(addMemberRequest); err != nil {
		return err
	}

	return nil

}

// UpdateNormalUser implements UserUsecase.
func (u *userUsecaseImpl) UpdateNormalUser(inUpdateModel *model.UpdateNormalUser, inNormalUserID uint) error {

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
	}

	if isUsernameAlreadyInUser(normalUser.Username, u.userrepository) {
		return errors.New("this username is already in use")
	}

	normalUser.ID = inNormalUserID

	if err := u.userrepository.UpdateNormalUser(normalUser); err != nil {
		return err
	}
	return nil
}

// CreateTeam implements UserUsecase.
func (u *userUsecaseImpl) CreateTeam(in *model.CreateTeam) error {

	normalUser, err := u.userrepository.GetNormalUser(&entities.NormalUsers{
		UsersID: in.OwnerID,
	})
	if err != nil {
		return err
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
		Compatitions:     []entities.CompatitionsTeams{},
		ImageProfilePath: "./images/default/defaultProfile.jpg",
		ImageCoverPath:   "./images/default/defaultCover.jpg",
	}

	if err := u.userrepository.InsertTeam(&team); err != nil {
		return err
	}

	return nil
}

// UpdateNormalUserPhone implements UserUsecase.
func (u *userUsecaseImpl) UpdateNormalUserPhone(inUserID uint, newPhone string) error {
	err := u.userrepository.UpdateNormalUserPhone(inUserID, newPhone)
	if err != nil {
		return err
	}
	return nil
}

func NewUserUsercaseImpl(
	userrepository repositories.Userrepository,
) UserUsecase {
	return &userUsecaseImpl{
		userrepository: userrepository,
	}
}

func (u *userUsecaseImpl) Login(in *model.LoginUser) (string, model.LoginResponse, error) {

	if !isEmail(in.Email) {
		return "", model.LoginResponse{}, errors.New("invalid email format")
	}

	//get user from email
	user, err := u.userrepository.GetUserByEmail(in.Email)
	if err != nil {
		return "", model.LoginResponse{}, errors.New("incorrect email or password")
	}
	//compare password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(in.Password)); err != nil {
		return "", model.LoginResponse{}, errors.New("incorrect email or password")
	}

	loginResponse := model.LoginResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	// pass = return jwt
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	if user.Role == "normal" {
		normalUser, err := u.userrepository.GetNormalUserWithAddressByUserID(user.ID)
		if err != nil {
			return "", model.LoginResponse{}, err
		}
		if user.ImageProfilePath != "" {
			user.ImageProfilePath = user.ImageProfilePath[1:]
		}
		if user.ImageCoverPath != "" {
			user.ImageCoverPath = user.ImageCoverPath[1:]
		}
		loginResponse.NormalUserID = normalUser.ID
		claims["normal_user_id"] = normalUser.ID
	} else if user.Role == "organizer" {
		organizer, err := u.userrepository.GetOrganizerWithAddressByUserID(user.ID)
		if err != nil {
			return "", model.LoginResponse{}, err
		}
		loginResponse.OrganizerID = organizer.ID
		claims["organizer_id"] = organizer.ID
	}
	//bypass
	jwtSecretKey := config.GetConfig().JwtSecretKey
	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", model.LoginResponse{}, err
	}
	return t, loginResponse, nil
}

// GetNormalUser implements UserUsecase.
func (u *userUsecaseImpl) GetNormalUser(id uint) (*model.NormalUserProfile, error) {
	normalUserEntity := &entities.NormalUsers{}
	normalUserEntity.ID = id
	resultNormalUser, err := u.userrepository.GetNormalUser(normalUserEntity)
	if err != nil {
		return nil, err
	}

	resultUser, err := u.userrepository.GetUserByID(resultNormalUser.UsersID)
	if err != nil {
		return nil, err
	}

	totalMatch := 0
	win := 0
	lose := 0
	recentMatch := []model.RecentMatch{}
	for _, compatition := range resultNormalUser.Compatitions {
		teamID := compatition.TeamsID
		for _, match := range compatition.Compatitions.Matchs {
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
				vsTeam, err := u.userrepository.GetTeam(match.Team2ID)
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
					TournamentName: compatition.Compatitions.Name,
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
				vsTeam, err := u.userrepository.GetTeam(match.Team1ID)
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
					Score:          fmt.Sprintf("%d - %d", match.Team2ID, match.Team1ID),
					TournamentName: compatition.Compatitions.Name,
				})
			}
		}
	}
	winRate := (float64(win) / float64(totalMatch)) * 100
	goalPerCompatition := float64(len(resultNormalUser.GoalRecords)) / float64(len(resultNormalUser.Compatitions))

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
		members := []model.Member{}
		for _, member := range team.Teams.TeamsMembers {
			members = append(members, model.Member{
				ID:            member.ID,
				UsersID:       member.NormalUsers.UsersID,
				FirstNameThai: member.NormalUsers.FirstNameThai,
				LastNameThai:  member.NormalUsers.LastNameThai,
				FirstNameEng:  member.NormalUsers.FirstNameEng,
				LastNameEng:   member.NormalUsers.LastNameEng,
				Position:      member.NormalUsers.Position,
				Sex:           member.NormalUsers.Sex,
			})
		}

		teamJoined = append(teamJoined, model.Team{
			ID:               team.Teams.ID,
			Name:             team.Teams.Name,
			OwnerID:          team.Teams.OwnerID,
			Description:      team.Teams.Description,
			Members:          members,
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
			WinRate:             winRate,
			TotalMatch:          totalMatch,
			Win:                 win,
			Lose:                lose,
			Goals:               len(resultNormalUser.GoalRecords),
			GoalsPerCompatition: goalPerCompatition,
			RecentMatch:         recentMatch,
		},
		TeamJoined: teamJoined,
	}
	return normalUserProfile, nil
}

// GetUser implements UserUsecase.
func (u *userUsecaseImpl) GetUser(in uint) (model.User, error) {
	// get user from email
	user, err := u.userrepository.GetUserByID(in)
	if err != nil {
		return model.User{}, err
	}

	userModel := model.User{
		ID:               user.ID,
		Email:            user.Email,
		Role:             user.Role,
		ImageProfilePath: user.ImageProfilePath,
		ImageCoverPath:   user.ImageCoverPath,
	}

	if user.Role == "normal" {
		normalUser, err := u.userrepository.GetNormalUserWithAddressByUserID(user.ID)
		if err != nil {
			return model.User{}, err
		}

		userModel.NormalUserInfo = &model.NormalUserInfo{
			ID:            normalUser.ID,
			FirstNameThai: normalUser.FirstNameThai,
			LastNameThai:  normalUser.LastNameThai,
			FirstNameEng:  normalUser.FirstNameEng,
			LastNameEng:   normalUser.LastNameEng,
			Username:      normalUser.Username,
			Born:          normalUser.Born,
			Phone:         normalUser.Phone,
			Height:        normalUser.Height,
			Weight:        normalUser.Weight,
			Sex:           normalUser.Sex,
			Position:      normalUser.Position,
			Nationality:   normalUser.Nationality,
			Description:   normalUser.Description,
			Address: model.Address{
				HouseNumber: normalUser.Addresses.HouseNumber,
				Village:     normalUser.Addresses.Village,
				Subdistrict: normalUser.Addresses.Subdistrict,
				District:    normalUser.Addresses.District,
				PostalCode:  normalUser.Addresses.PostalCode,
				Country:     normalUser.Addresses.Country,
			},
		}
		// userModel.OrganizersInfo = model.OrganizersInfo{}

	} else if user.Role == "organizer" {
		organizer, err := u.userrepository.GetOrganizerWithAddressByUserID(user.ID)
		if err != nil {
			return model.User{}, err
		}
		userModel.OrganizersInfo = &model.OrganizersInfo{
			ID:          organizer.ID,
			Name:        organizer.Name,
			Phone:       organizer.Phone,
			Description: organizer.Description,
			Address: model.Address{
				HouseNumber: organizer.Addresses.HouseNumber,
				Village:     organizer.Addresses.Village,
				Subdistrict: organizer.Addresses.Subdistrict,
				District:    organizer.Addresses.District,
				PostalCode:  organizer.Addresses.PostalCode,
				Country:     organizer.Addresses.Country,
			},
		}

	}
	return userModel, nil
}
func (u *userUsecaseImpl) RegisterNormaluser(in *model.RegisterNormaluser) error {

	if !isEmail(in.Email) {
		return errors.New("invalid email format")
	}

	if isEmailAlreadyInUse(in.Email, u.userrepository) {
		return errors.New("this email is already in use")
	}
	if isUsernameAlreadyInUser(in.Username, u.userrepository) {
		return errors.New("this username is already in use")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entities.Users{
		Email:            in.Email,
		Role:             "normal",
		Password:         string(hashedPassword),
		ImageProfilePath: "./images/default/defaultProfile.jpg",
		ImageCoverPath:   "./images/default/defaultCover.jpg",
	}

	normalUser := &entities.NormalUsers{
		UsersID:  user.ID,
		Username: in.Username,
	}

	if err := u.userrepository.InsertUserWihtNormalUserAndAddress(normalUser, user); err != nil {
		return err
	}

	return nil
}

func isEmailAlreadyInUse(email string, u repositories.Userrepository) bool {
	if _, err := u.GetUserByEmail(email); err != nil {
		return false
	}
	return true
}

func isUsernameAlreadyInUser(username string, u repositories.Userrepository) bool {
	if _, err := u.GetNormalUserByUsername(username); err != nil {
		return false
	}
	return true
}

func isPhoneAlreadyInUse(phone string, u repositories.Userrepository) bool {
	if _, err := u.GetNormalUser(&entities.NormalUsers{
		Phone: phone,
	}); err != nil {
		if _, err := u.GetOrganizer(&entities.Organizers{
			Phone: phone,
		}); err != nil {
			return false
		}
	}
	return true
}

func (u *userUsecaseImpl) RegisterOrganizer(in *model.RegisterOrganizer) error {

	if !isEmail(in.Email) {
		return errors.New("invalid email format")
	}

	if isEmailAlreadyInUse(in.Email, u.userrepository) {
		return errors.New("this email is already in use")
	}

	if isPhoneAlreadyInUse(in.Phone, u.userrepository) {
		return errors.New("this phone is already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entities.Users{
		Email:            in.Email,
		Role:             "organizer",
		Password:         string(hashedPassword),
		ImageProfilePath: "./images/default/defaultProfile.jpg",
		ImageCoverPath:   "./images/default/defaultCover.jpg",
	}

	organizer := &entities.Organizers{
		UsersID: user.ID,
		Name:    in.OrganizerName,
		Phone:   in.Phone,
	}

	if err := u.userrepository.InsertUserWihtOrganizerAndAddress(organizer, user); err != nil {
		return err
	}

	return nil
}

// GetUsers implements UserUsecase.
func (u *userUsecaseImpl) GetUsers() ([]model.User, error) {
	users_entity, err := u.userrepository.GetUsers()
	if err != nil {
		return []model.User{}, err
	}
	users_model := []model.User{}
	for _, e := range users_entity {
		m := model.User{
			ID:               e.ID,
			Email:            e.Email,
			Role:             e.Role,
			ImageProfilePath: e.ImageProfilePath,
			ImageCoverPath:   e.ImageCoverPath,
		}
		users_model = append(users_model, m)
	}

	return users_model, nil
}
