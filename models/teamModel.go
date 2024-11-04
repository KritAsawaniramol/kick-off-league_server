package model

type CreateTeam struct {
	Name        string        `json:"name" binding:"required"`
	OwnerID     uint          `json:"owner_id"`
	Member      []Member      `json:"member"`
	Compatition []Compatition `json:"compatition"`
	Description string        `json:"description"`
}

type Team struct {
	ID               uint                   `json:"id"`
	Name             string                 `json:"name" binding:"required"`
	OwnerID          uint                   `json:"owner_id" binding:"required"`
	Description      string                 `json:"description"`
	Members          []Member               `json:"member,omitempty"`
	Compatitions     []CompatitionBasicInfo `json:"compatitions,omitempty"`
	NumberOfMember   int                   `json:"number_of_member"`
	Rank             string                 `json:"rank"`
	RankNumber       int                    `json:"rank_number"`
	Point            int                    `json:"point"`
	GoalsScored      int                    `json:"goals_scored"`
	GoalsConceded    int                    `json:"goals_conceded"`
	ImageProfilePath string                 `json:"image_profile_path"`
	ImageCoverPath   string                 `json:"image_cover_path"`
}

type Member struct {
	ID               uint   `json:"id"`
	UsersID          uint   `json:"user_id"`
	FirstNameThai    string `json:"first_name_thai"`
	LastNameThai     string `json:"last_name_thai"`
	FirstNameEng     string `json:"first_name_eng"`
	LastNameEng      string `json:"last_name_eng"`
	Position         string `json:"position"`
	Sex              string `json:"sex"`
	Role             string `json:"role"`
	ImageProfilePath string `json:"image_profile_path"`
	ImageCoverPath   string `json:"image_cover_path"`
}

type GetTeamsReq struct {
	TeamID       uint   `json:"team_id"`
	NormalUserID uint   `json:"owner_id"`
	Ordering     string `json:"ordering"`
	Decs         bool   `json:"decs"`
	Page         uint   `json:"page"`
	PageSize     uint   `json:"page_size"`
}

type TeamList struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	NumberOfMember   uint   `json:"number_of_member"`
	OwnerID          uint   `json:"owner_id"`
	ImageProfilePath string `json:"image_profile_path"`
	ImageCoverPath   string `json:"image_cover_path"`
}
