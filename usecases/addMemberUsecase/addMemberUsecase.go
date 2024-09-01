package addMemberUsecase

import (
	model "kickoff-league.com/models"
)

type AddMemberUsecase interface {
	GetMyPenddingAddMemberRequest(userID uint) ([]model.AddMemberRequest, error)
	SendAddMemberRequest(inAddMemberRequest *model.AddMemberRequest, inUserID uint) error
	AcceptAddMemberRequest(inReqID uint, userID uint) error
	IgnoreAddMemberRequest(inReqID uint, userID uint) error
}
