package model

import "time"

type NormalUserInfo struct {
	ID            uint      `json:"id"`
	FirstNameThai string    `json:"first_name_thai"`
	LastNameThai  string    `json:"last_name_thai"`
	FirstNameEng  string    `json:"first_name_eng"`
	LastNameEng   string    `json:"last_name_eng"`
	Username      string    `json:"username"`
	Born          time.Time `json:"born"`
	Phone         string    `json:"phone"`
	Height        uint      `json:"height"`
	Weight        uint      `json:"weight"`
	Sex           string    `json:"sex"`
	Position      string    `json:"position"`
	Nationality   string    `json:"nationality"`
	Description   string    `json:"description"`
	Address       `json:"address"`
}

type NormalUserProfile struct {
	NormalUserInfo   `json:"normal_user_info"`
	UserID           uint   `json:"user_id"`
	ImageProfilePath string `json:"image_profile_path"`
	ImageCoverPath   string `json:"image_cover_path"`
	TeamJoined       []Team
	NormalUserStat   `json:"normal_user_stat"`
}

type NormalUserStat struct {
	WinRate       float64       `json:"win_rate"`
	TotalMatch    int           `json:"total_match"`
	Win           int           `json:"win"`
	Lose          int           `json:"lose"`
	Goals         int           `json:"goals"`
	GoalsPerMatch float64       `json:"goals_per_match"`
	RecentMatch   []RecentMatch `json:"recent_match"` // Last 20 matches
}

type NormalUserList struct {
	ID               uint      `json:"id"`
	FirstNameThai    string    `json:"first_name_thai"`
	LastNameThai     string    `json:"last_name_thai"`
	FirstNameEng     string    `json:"first_name_eng"`
	LastNameEng      string    `json:"last_name_eng"`
	Born             time.Time `json:"born"`
	Height           uint      `json:"height"`
	Weight           uint      `json:"weight"`
	Sex              string    `json:"sex"`
	Position         string    `json:"position"`
	Nationality      string    `json:"nationality"`
	Description      string    `json:"description"`
	ImageProfilePath string    `json:"image_profile_path"`
	ImageCoverPath   string    `json:"image_cover_path"`
}

type UpdateNormalUser struct {
	FirstNameThai string    `json:"first_name_thai"`
	LastNameThai  string    `json:"last_name_thai"`
	FirstNameEng  string    `json:"first_name_eng"`
	LastNameEng   string    `json:"last_name_eng"`
	Born          time.Time `json:"born"`
	Phone         string    `json:"phone"`
	Height        uint      `json:"height"`
	Weight        uint      `json:"weight"`
	Sex           string    `json:"sex"`
	Position      string    `json:"position"`
	Nationality   string    `json:"nationality"`
	Description   string    `json:"description"`
	Username      string    `json:"username"`
	Address       Address   `jsoon:"address"`
}
