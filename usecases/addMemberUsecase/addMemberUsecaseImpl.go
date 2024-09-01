package addMemberUsecase

import (
	"errors"
	"log"

	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
)

type addMemberUsecaseImpl struct {
	repository repositories.Repository
}

func NewAddMemberUsecaseImpl(
	repository repositories.Repository,
) AddMemberUsecase {
	return &addMemberUsecaseImpl{
		repository: repository,
	}
}

// SendAddMemberRequest implements AddMemberUsecase.
func (a *addMemberUsecaseImpl) SendAddMemberRequest(inAddMemberRequest *model.AddMemberRequest, inUserID uint) error {

	team, err := a.repository.GetTeamWithMemberAndRequestSentByID(inAddMemberRequest.TeamID)
	if err != nil {
		return err
	}

	if team.OwnerID != inUserID {
		return errors.New("this user isn't owner of this team")
	}

	// Get Receiver normaluser
	receiver, err := a.repository.GetNormalUserDetails(&entities.NormalUsers{
		Username: inAddMemberRequest.ReceiverUsername,
	})
	if err != nil {
		return err
	}

	for _, member := range team.TeamsMembers {
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
	if err := a.repository.InsertAddMemberRequest(addMemberRequest); err != nil {
		return err
	}

	return nil
}

// AcceptAddMemberRequest implements AddMemberUsecase.
func (a *addMemberUsecaseImpl) AcceptAddMemberRequest(inReqID uint, userID uint) error {
	addMemberRequestSearch := &entities.AddMemberRequests{}
	addMemberRequestSearch.ID = inReqID
	addMemberRequestList, err := a.repository.GetAddMemberRequestByID(addMemberRequestSearch)
	if err != nil {
		return err
	}
	addMemberRequest := &addMemberRequestList[0]

	if addMemberRequest.ReceiverID != userID {
		log.Printf("ReceiverID: %d, userID: %d\n", addMemberRequest.ReceiverID, userID)
		return errors.New("UserID does not match ReceiverID")
	}

	normalUser, err := a.repository.GetNormalUserDetails(&entities.NormalUsers{UsersID: userID})
	if err != nil {
		return err
	}


	if err := a.repository.UpdateAddmemberRequestAndInsertTeamsMembers(&entities.TeamsMembers{
		TeamsID:       addMemberRequest.TeamsID,
		NormalUsersID: normalUser.ID,
		Role:          addMemberRequest.Role,
	}, inReqID, "accepted",
	); err != nil {
		return err
	}

	return nil
}


// IgnoreAddMemberRequest implements AddMemberUsecase.
func (a *addMemberUsecaseImpl) IgnoreAddMemberRequest(inReqID uint, userID uint) error {
	addMemberRequestSearch := &entities.AddMemberRequests{}
	addMemberRequestSearch.ID = inReqID
	addMemberRequestList, err := a.repository.GetAddMemberRequestByID(addMemberRequestSearch)
	if err != nil {
		return err
	}
	addMemberRequest := &addMemberRequestList[0]
	if addMemberRequest.ReceiverID != userID {
		return errors.New("userID does not match ReceiverID")
	}
	if err := a.repository.UpdateAddMemberRequestStatusAndSoftDelete(addMemberRequest, "ignore"); err != nil {
		return err
	}
	return nil
}


// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================

// GetMyPenddingAddMemberRequest implements AddMemberUsecase.
func (a *addMemberUsecaseImpl) GetMyPenddingAddMemberRequest(userID uint) ([]model.AddMemberRequest, error) {
	addMemberRequestSearch := &entities.AddMemberRequests{
		ReceiverID: userID,
		Status:     "pending",
	}
	addMemberRequestList, err := a.repository.GetAddMemberRequestByID(addMemberRequestSearch)
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
