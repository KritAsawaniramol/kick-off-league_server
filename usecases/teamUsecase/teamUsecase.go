package teamUsecase

import (
	model "kickoff-league.com/models"
)

type TeamUsecase interface {
	RemoveTeamImageProfile(teamID uint, userID uint) error	
	RemoveTeamImageCover(teamID uint, userID uint) error
	GetTeamWithMemberAndCompatitionByID(id uint) (*model.Team, error)
	RemoveTeamFormCompatition(teamID uint, compatitionID uint, orgID uint) error
	GetTeams(in *model.GetTeamsReq) ([]model.TeamList, error)
	GetTeamsByOwnerID(in uint) ([]model.TeamList, error)
	CreateTeam(in *model.CreateTeam) error
	UpdateTeamImageCover(teamID uint, newImagePath string, userID uint) error
	UpdateTeamImageProfile(teamID uint, newImagePath string, userID uint) error	
	RemoveNormalUserFormTeam(teamID uint, nomalUserID uint, ownerID uint) error
	
	// GetMyPenddingAddMemberRequest(userID uint) ([]model.AddMemberRequest, error)
	// SendAddMemberRequest(in *model.AddMemberRequest, userID uint) error
	// AcceptAddMemberRequest(inReqID uint, userID uint) error
	// IgnoreAddMemberRequest(inReqID uint, userID uint) error
}
