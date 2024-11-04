package model

type AddMemberRequest struct {
	ID               uint   `json:"id"`
	TeamName         string `json:"team_name"`
	TeamImageProfile string `json:"team_iamge_profile"`
	TeamID           uint   `json:"team_id" binding:"required"`
	ReceiverUsername string `json:"receiver_username" binding:"required"`
	Role             string `json:"role" binding:"required"`
	Status           string `json:"status"`
}