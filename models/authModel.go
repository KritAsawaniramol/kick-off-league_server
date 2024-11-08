package model

type RegisterUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterNormaluser struct {
	RegisterUser
	Username string `json:"username" binding:"required"`
}

type RegisterOrganizer struct {
	RegisterUser
	OrganizerName string `json:"name" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}


type LoginResponse struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	NormalUserID uint   `json:"normal_user_id,omitempty"`
	OrganizerID  uint   `json:"organizer_id,omitempty"`
}