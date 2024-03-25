package model

import (
	"time"
)

type OrganizersInfo struct {
	ID               uint
	Name             string
	Phone            string
	Description      string
	Address          Address
	ImageProfilePath string
	ImageCoverPath   string
}

// type MyNormalUser struct {
// 	NormalUser
// }

type LoginResponse struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	NormalUserID uint   `json:"normal_user_id,omitempty"`
	OrganizerID  uint   `json:"organizer_id,omitempty"`
}

type User struct {
	ID               uint            `json:"id"`
	Email            string          `json:"email"`
	Role             string          `json:"role"`
	NormalUserInfo   *NormalUserInfo `json:"normal_user,omitempty"`
	OrganizersInfo   *OrganizersInfo `json:"organizer,omitempty"`
	ImageProfilePath string          `json:"image_profile_path"`
	ImageCoverPath   string          `json:"image_cover_path"`
}

type NormalUser struct {
	NormalUserInfo
	Team        []Team     `json:"team"`
	TeamCreated []Team     `json:"team_created"`
	GoalRecord  GoalRecord `json:"goal_record"`
}

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
	NormalUserInfo
	UserID           uint   `json:"user_id"`
	ImageProfilePath string `json:"image_profile_path"`
	ImageCoverPath   string `json:"image_cover_path"`
	NormalUserStat   `json:"normal_user_stat"`
}

type NormalUserStat struct {
	WinRate     float32       `json:"win_rate"`
	TotalMatch  int           `json:"total_match"`
	Win         int           `json:"win"`
	Lose        int           `json:"lose"`
	Draw        int           `json:"draw"`
	Goals       int           `json:"goals"`
	RecentMatch []RecentMatch `json:"recent_match"`
}

type NormalUserList struct {
	ID            uint      `json:"id"`
	FirstNameThai string    `json:"first_name_thai"`
	LastNameThai  string    `json:"last_name_thai"`
	FirstNameEng  string    `json:"first_name_eng"`
	LastNameEng   string    `json:"last_name_eng"`
	Born          time.Time `json:"born"`
	Height        uint      `json:"height"`
	Weight        uint      `json:"weight"`
	Sex           string    `json:"sex"`
	Position      string    `json:"position"`
	Nationality   string    `json:"nationality"`
	Description   string    `json:"description"`
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
}

type AddMemberRequest struct {
	ID               uint   `json:"id"`
	TeamID           uint   `json:"team_id" binding:"required"`
	TeamName         string `json:"team_name"`
	ReceiverUsername string `json:"receiver_username" binding:"required"`
	Role             string `json:"role" binding:"required"`
	Status           string `json:"status"`
}

type GoalRecord struct {
	MatchesID  uint
	TeamID     uint
	PlayerID   uint
	TimeScored uint
}

type CreateTeam struct {
	Name        string        `json:"name" binding:"required"`
	OwnerID     uint          `json:"owner_id"`
	Member      []Member      `json:"member"`
	Compatition []Compatition `json:"compatition"`
	Description string        `json:"description"`
}

type Team struct {
	ID           uint                   `json:"id"`
	Name         string                 `json:"name" binding:"required"`
	OwnerID      uint                   `json:"owner_id" binding:"required"`
	Description  string                 `json:"description"`
	Members      []Member               `json:"member"`
	Compatitions []CompatitionBasicInfo `json:"compatitions,omitempty"`
}

type CompatitionBasicInfo struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Format      string    `json:"format"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`

	AgeOver      uint   `json:"age_over"`
	AgeUnder     uint   `json:"age_under"`
	Sex          string `json:"sex"`
	FieldSurface string `json:"field_surface"`
	Status       string `json:"status"`
	NumberOfTeam uint   `json:"number_of_team"`
	OrganizerID  uint   `json:"organizer_id"`
}

type Compatition struct {
	CompatitionBasicInfo
	AddressID            uint   `json:"address_id"`
	Teams                []Team `json:"teams"`
	NumOfPlayerInTeamMin uint   `json:"num_of_player_min"`
	NumOfPlayerInTeamMax uint   `json:"num_of_player_max"`
}

type GetCompatitions struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Sport           string    `json:"sport"`
	Format          string    `json:"format"` // 1 vs 1, 2 vs 2,...
	Address         Address   `json:"address"`
	Status          string    `json:"status"`
	ApplicationType string    `json:"application_type"`
	Sex             string    `json:"sex"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	OrganizerID     uint      `json:"organizer_id"`
	OrganizerName   string    `json:"organizer_name"`
	AgeOver         uint      `json:"age_over"`
	AgeUnder        uint      `json:"age_under"`
}

type GetCompatition struct {
	ID                   uint           `json:"id"`
	CreatedAt            time.Time      `json:"created_at"`
	Name                 string         `json:"name" binding:"required"`
	Sport                string         `json:"sport" binding:"required"`
	Format               string         `json:"format" binding:"required"` // 1 vs 1, 2 vs 2,...
	Type                 string         `json:"type" binding:"required"`   // Tournament, Round Robbin,....
	OrganizerInfo        OrganizersInfo `json:"organizer_info"`
	FieldSurface         string         `json:"field_surface"`
	ApplicationType      string         `json:"application_type" binding:"required"` // free, with code
	Address              Address        `json:"address" binding:"required"`
	ImageBanner          string         `json:"image_banner"`
	StartDate            time.Time      `json:"start_date" binding:"required"`
	EndDate              time.Time      `json:"end_date" binding:"required"`
	JoinCode             []JoinCode     `json:"join_code"`
	Description          string         `json:"description"`
	Rule                 string         `json:"rule"`
	Prize                string         `json:"prize"`
	ContractType         string         `json:"contract_type"  binding:"required"`
	Contract             string         `json:"contract"  binding:"required"`
	AgeOver              uint           `json:"age_over"`
	AgeUnder             uint           `json:"age_under"`
	Sex                  string         `json:"sex" binding:"required"`
	Status               string         `json:"status"`
	NumberOfTeam         uint           `json:"number_of_team" binding:"required"`
	NumOfPlayerInTeamMin uint           `json:"num_of_player_min"`
	NumOfPlayerInTeamMax uint           `json:"num_of_player_max"`
	Teams                []Team         `json:"teams"`
	NumOfRound           int            `json:"number_of_round"`
	NumOfMatch           int            `json:"number_of_match"`
	Matches              []Match        `json:"matches"`
}

type UpdateCompatition struct {
	Name                 string    `json:"name"`
	Sport                string    `json:"sport"`
	Format               string    `json:"format"` // 1 vs 1, 2 vs 2,...
	Type                 string    `json:"type"`   // Tournament, Round Robbin,....
	FieldSurface         string    `json:"field_surface"`
	ApplicationType      string    `json:"application_type"` // free, with code
	Address              Address   `json:"address"`
	ImageBanner          string    `json:"image_banner"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Description          string    `json:"description"`
	Rule                 string    `json:"rule"`
	Prize                string    `json:"prize"`
	ContractType         string    `json:"contract_type" `
	Contract             string    `json:"contract" `
	AgeOver              uint      `json:"age_over"`
	AgeUnder             uint      `json:"age_under"`
	Sex                  string    `json:"sex"`
	NumberOfTeam         uint      `json:"number_of_team"`
	NumOfPlayerInTeamMin uint      `json:"num_of_player_min"`
	NumOfPlayerInTeamMax uint      `json:"num_of_player_max"`
}

type RecentMatch struct {
	ID         uint      `json:"id"`
	DateTime   time.Time `json:"date_time"`
	VsTeamName string    `json:"vs_team_name"`
	Result     string    `json:"result"`
	Score      string    `json:"score"`
}

type Match struct {
	ID             uint         `json:"id"`
	Index          int          `json:"index"`
	DateTime       time.Time    `json:"date_time"`
	Team1ID        uint         `json:"team1_id"`
	Team2ID        uint         `json:"team2_id"`
	Team1Goals     int          `json:"team1_goals"`
	Team2Goals     int          `json:"team2_goals"`
	Round          string       `json:"round"`
	NextMatchIndex int          `json:"next_match_index"`
	NextMatchSlot  string       `json:"next_match_slot"` //Team1 or Team2
	GoalRecords    []GoalRecord `json:"goal_records"`
	Result         string       `json:"result"`
}

type JoinCode struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Code      string    `json:"Code"`
	Status    string    `json:"status"`
}

type CreateCompatition struct {
	Name   string `json:"name" binding:"required"`
	Sport  string `json:"sport" binding:"required"`
	Type   string `json:"type" binding:"required"`   // Tournament, Round Robbin,....
	Format string `json:"format" binding:"required"` // 1 vs 1, 2 vs 2,...

	// Desscription
	Description string `json:"description"`
	Rule        string `json:"rule"`
	Prize       string `json:"prize"`

	StartDate       time.Time `json:"start_date" binding:"required"`
	EndDate         time.Time `json:"end_date" binding:"required"`
	ApplicationType string    `json:"application_type" binding:"required"` // free, with code
	ImageBanner     string    `json:"image_banner"`

	// Condition
	AgeOver      uint   `json:"age_over"`
	AgeUnder     uint   `json:"age_under"`
	Sex          string `json:"sex" binding:"required"`
	NumberOfTeam uint   `json:"number_of_team" binding:"required"`

	NumOfPlayerInTeamMin uint `json:"num_of_player_min"`
	NumOfPlayerInTeamMax uint `json:"num_of_player_max"`

	FieldSurface string  `json:"field_surface"`
	OrganizerID  uint    `json:"organizer_id"`
	Address      Address `json:"address" binding:"required"`

	ContractType string `json:"contract_type"  binding:"required"`
	Contract     string `json:"contract"  binding:"required"`
}

type Address struct {
	HouseNumber string `json:"house_number"`
	Village     string `json:"village"`
	Subdistrict string `json:"subdistrict"`
	District    string `json:"district"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
}

type Member struct {
	ID            uint   `json:"id"`
	UsersID       uint   `json:"user_id"`
	FirstNameThai string `json:"first_name_thai"`
	LastNameThai  string `json:"last_name_thai"`
	FirstNameEng  string `json:"first_name_eng"`
	LastNameEng   string `json:"last_name_eng"`
	Position      string `json:"position"`
	Sex           string `json:"sex"`
	Role          string `json:"role"`
}

type JoinCompatition struct {
	CompatitionID uint   `json:"compatition_id" binding:"required"`
	TeamID        uint   `json:"team_id" binding:"required"`
	Code          string `json:"code" binding:"required"`
}

type GetCompatitionsReq struct {
	TeamID       uint   `json:"team_id"`
	NormalUserID uint   `json:"owner_id"`
	OrganizerID  uint   `json:"organizer_id"`
	Ordering     string `json:"ordering"`
	Decs         bool   `json:"decs"`
	Page         uint   `json:"page"`
	PageSize     uint   `json:"page_size"`
}

type GetTeamsReq struct {
	TeamID       uint   `json:"team_id"`
	NormalUserID uint   `json:"owner_id"`
	Ordering     string `json:"ordering"`
	Decs         bool   `json:"decs"`
	Page         uint   `json:"page"`
	PageSize     uint   `json:"page_size"`
}

type CompatitionList struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Description    string `jsFcon:"description"`
	NumberOfMember uint   `json:"number_of_member"`
	// img         string
}

type TeamList struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	NumberOfMember uint   `json:"number_of_member"`
	// img         string
}

// type CompetitionFormat string

// const (
// 	Tournament CompetitionFormat = "Tournament"
// 	RoundRobin CompetitionFormat = "Round Robin"
// )

// type FieldSurfaces string

// const (
// 	NaturalGrass   FieldSurfaces = "NaturalGrass"
// 	ArtificialTurf FieldSurfaces = "ArtificialTurf"
// 	FlatSurface    FieldSurfaces = "FlatSurface"
// 	Other          FieldSurfaces = "Other"
// )

// type SexType string

// const (
// 	Male   SexType = "Male"
// 	Female SexType = "Female"
// 	Unisex SexType = "Unisex"
// )

// type MatchesResult string

// const (
// 	Team1Win MatchesResult = "Team1Win"
// 	Team2Win MatchesResult = "Team2Win"
// 	Draw     MatchesResult = "Draw"
// )

var SexType = []string{"Male", "Female", "Unisex"}
var FieldSurfaces = []string{"NaturalGrass", "ArtificialTurf", "FlatSurface"}
var CompetitionFormat = []string{"Tournament", "Round Robin"}
var MatchesResult = []string{"Team1Win", "Team2Win", "Draw"}
