package model

type AddMemberRequest struct {
	ID               uint   `json:"id"`
	TeamID           uint   `json:"team_id" binding:"required"`
	TeamName         string `json:"team_name"`
	TeamImageProfile string `json:"team_iamge_profile"`
	ReceiverUsername string `json:"receiver_username" binding:"required"`
	Role             string `json:"role" binding:"required"`
	Status           string `json:"status"`
}