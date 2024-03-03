package model

import (
	"time"
)

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

type OrganizersInfo struct {
	ID          uint
	Name        string
	Phone       string
	Description string
	Address     Address
	// imageProfile string
}

type User struct {
	ID        uint
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt time.Time
	Datail    interface{}
}

type AddMemberRequest struct {
	ID               uint   `json:"id"`
	TeamID           uint   `json:"team_id" binding:"required"`
	TeamName         string `json:"team_name"`
	ReceiverUsername string `json:"receiver_username" binding:"required"`
	Role             string `json:"role" binding:"required"`
	Status           string `json:"status"`
}

type NormalUser struct {
	NormalUserInfo
	Team        []Team
	TeamCreated []Team
	GoalRecord  GoalRecord
	// imageProfile string
}

type NormalUserInfo struct {
	ID               uint      `json:"id"`
	FirstNameThai    string    `json:"first_name_thai"`
	LastNameThai     string    `json:"last_name_thai"`
	FirstNameEng     string    `json:"first_name_eng"`
	LastNameEng      string    `json:"last_name_eng"`
	Born             time.Time `json:"born"`
	Phone            string    `json:"phone"`
	Height           uint      `json:"height"`
	Weight           uint      `json:"weight"`
	Sex              string    `json:"sex"`
	Position         string    `json:"position"`
	Nationality      string    `json:"nationality"`
	Description      string    `json:"description"`
	ImageProfilePath string    `json:"image_profile_path"`
	ImageCoverPath   string    `json:"image_cover_path"`
	Address          `json:"address"`
}

type UpdateNormalUser struct {
	FirstNameThai    string    `json:"first_name_thai"`
	LastNameThai     string    `json:"last_name_thai"`
	FirstNameEng     string    `json:"first_name_eng"`
	LastNameEng      string    `json:"last_name_eng"`
	Born             time.Time `json:"born"`
	Phone            string    `json:"phone"`
	Height           uint      `json:"height"`
	Weight           uint      `json:"weight"`
	Sex              string    `json:"sex"`
	Position         string    `json:"position"`
	Nationality      string    `json:"nationality"`
	Description      string    `json:"description"`
	ImageProfilePath string    `json:"image_profile_path"`
	ImageCoverPath   string    `json:"image_cover_path"`
}

type GoalRecord struct {
	MatchesID  uint
	TeamID     uint
	PlayerID   uint
	TimeScored uint
}

type CreaetTeam struct {
	Name        string        `json:"name" binding:"required"`
	OwnerID     uint          `json:"owner_id"`
	Member      []NormalUser  `json:"member"`
	Compatition []Compatition `json:"compatition"`
	Description string        `json:"description"`
	Role        string        `json:"role"`
}

type CreateTeam struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name" binding:"required"`
	OwnerID     uint          `json:"owner_id" binding:"required"`
	Member      []Member      `json:"member"`
	Compatition []Compatition `json:"compatition"`
	Description string        `json:"description"`
}

type Team struct {
	ID           uint                   `json:"id"`
	Name         string                 `json:"name" binding:"required"`
	OwnerID      uint                   `json:"owner_id" binding:"required"`
	Members      []Member               `json:"member"`
	Compatitions []CompatitionBasicInfo `json:"compatitions"`
	Description  string                 `json:"description"`
}

type CompatitionBasicInfo struct {
	ID                uint              `json:"id"`
	Name              string            `json:"name"`
	Format            CompetitionFormat `json:"format"`
	Description       string            `json:"description"`
	StartDate         time.Time         `json:"start_date"`
	EndDate           time.Time         `json:"end_date"`
	RegisterStartDate time.Time         `json:"register_start_date"`
	RegisterEndDate   time.Time         `json:"register_end_date"`
	ApplicationFee    float64           `json:"application_fee"`
	AgeOver           uint              `json:"age_over"`
	AgeUnder          uint              `json:"age_under"`
	Sex               SexType           `json:"sex"`
	FieldSurface      FieldSurfaces     `json:"field_surface"`
	Status            CompetitionStatus `json:"status"`
	NumberOfTeam      uint              `json:"number_of_team"`
	OrganizerID       uint              `json:"organizer_id"`
}

type Compatition struct {
	CompatitionBasicInfo
	AddressID      uint   `json:"address_id"`
	Teams          []Team `json:"teams"`
	NumOfPlayerMin uint   `json:"num_of_player_min"`
	NumOfPlayerMax uint   `json:"num_of_player_max"`
}

type CreateCompatition struct {
	Name        string            `json:"name" binding:"required"`
	Format      CompetitionFormat `json:"format"`
	Description string            `json:"description"`
	StartDate   time.Time         `json:"start_date"`
	EndDate     time.Time         `json:"end_date"`
	AgeOver     uint              `json:"age_over"`
	AgeUnder    uint              `json:"age_under"`
	Sex         SexType           `json:"sex"`
	OrganizerID uint              `json:"organizer_id"`
	Address     Address           `json:"address"`
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
	FirstNameThai string `json:"first_name_thai"`
	LastNameThai  string `json:"last_name_thai"`
	FirstNameEng  string `json:"first_name_eng"`
	LastNameEng   string `json:"last_name_eng"`
	Position      string `json:"position"`
	Sex           string `json:"sex"`
	Role          string `json:"role"`
}

type GetCompatitionsReq struct {
	TeamID       uint   `json:"team_id"`
	NormalUserID uint   `json:"owner_id"`
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
	Description    string `json:"description"`
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

type CompetitionFormat string

const (
	Tournament               CompetitionFormat = "Tournament"
	GroupStage               CompetitionFormat = "GroupStage"
	TournamentsAndGroupStage CompetitionFormat = "TournamentsAndGroupStage"
)

type FieldSurfaces string

const (
	NaturalGrass   FieldSurfaces = "naturalGrass"
	ArtificialTurf FieldSurfaces = "artificialTurf"
	FlatSurface    FieldSurfaces = "flatSurface"
	Other          FieldSurfaces = "other"
)

type CompetitionStatus string

const (
	ComingSoon         CompetitionStatus = "ComingSoon"
	ApplicationsOpened CompetitionStatus = "ApplicationsOpened"
	ApplicationsEnded  CompetitionStatus = "ApplicationsEnded"
	Started            CompetitionStatus = "Started"
	Finished           CompetitionStatus = "Finished"
	Cancelled          CompetitionStatus = "Cancelled"
)

type SexType string

const (
	Male   SexType = "Male"
	Female SexType = "Female"
	Unisex SexType = "Unisex"
)

type MatchesResult string

const (
	Team1Win MatchesResult = "Team1Win"
	Team2Win MatchesResult = "Team2Win"
	Draw     MatchesResult = "Draw"
)
