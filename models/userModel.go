package model

import (
	"time"
)

type RegisterUser struct {
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterOrganizer struct {
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password" binding:"required"`
	Role          string `json:"role" binding:"required"`
	OrganizerName string `json:"name" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID        uint
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt time.Time
	Datail    map[string]interface{}
}

type NormalUser_2 struct {
	UserID uint
	Email  string `json:"email"`

	FirstNameThai string
	LastNameThai  string
	FirstNameEng  string
	LastNameEng   string
	Born          time.Time
	Height        uint
	Weight        uint
	Sex           string
	Position      string
	Nationality   string
	Phone         string
	AddressID     uint
	Address       Address
	Team          []Team
	TeamCreated   []Team
	GoalRecord    GoalRecord
	// imageProfile string
}

type NormalUser struct {
	UserID uint
	Email  string `json:"email"`

	FirstNameThai string
	LastNameThai  string
	FirstNameEng  string
	LastNameEng   string
	Born          time.Time
	Height        uint
	Weight        uint
	Sex           string
	Position      string
	Nationality   string
	Phone         string
	AddressID     uint
	Address       Address
	Team          []Team
	TeamCreated   []Team
	GoalRecord    GoalRecord
	// imageProfile string
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

type Address struct {
	HouseNumber string
	Village     string
	Subdistrict string
	District    string
	PostalCode  string
	Country     string
}

type GoalRecord struct {
	MatchesID  uint
	TeamID     uint
	PlayerID   uint
	TimeScored uint
}

type CreaetTeam struct {
	Name        string        `json:"name" binding:"required"`
	OwnerID     uint          `json:"ownerID"`
	Member      []NormalUser  `json:"member"`
	Compatition []Compatition `json:"compatition"`
}

type Team struct {
	Name        string        `json:"name" binding:"required"`
	OwnerID     uint          `json:"ownerID" binding:"required"`
	Member      []NormalUser  `json:"member"`
	Compatition []Compatition `json:"compatition"`
}

type Compatition struct {
	Name              string
	Format            CompetitionFormat
	OrganizerID       uint
	StartDate         time.Time
	EndDate           time.Time
	RegisterStartDate time.Time
	RegisterEndDate   time.Time
	ApplicationFee    float64
	AgeOver           uint
	AgeUnder          uint
	Sex               SexType
	FieldSurface      FieldSurfaces
	Description       string
	Status            CompetitionStatus
	AddressID         uint
	NumberOfTeam      uint
	Teams             []Team `gorm:"many2many:compatition_team;"`
	NumOfPlayerMin    uint
	NumOfPlayerMax    uint
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
