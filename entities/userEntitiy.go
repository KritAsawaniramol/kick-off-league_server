package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	Addresses struct {
		gorm.Model
		HouseNumber string
		Village     string
		Subdistrict string
		District    string
		PostalCode  string
		Country     string
	}

	CompatitionAddresses struct {
		gorm.Model
		CompatitionsID uint
		HouseNumber    string
		Village        string
		Subdistrict    string
		District       string
		PostalCode     string
		Country        string
	}

	Users struct {
		gorm.Model
		Email    string `gorm:"unique;not null;type:varchar(100)"`
		Role     string
		Password string
	}

	Organizers struct {
		gorm.Model
		UsersID      uint `gorm:"unique;not null"`
		Name         string
		Phone        string `gorm:"unique;not null"`
		Description  string
		AddressesID  uint `gorm:"unique"`
		Addresses    Addresses
		Compatitions []Compatitions
		ImageURL     string
	}

	NormalUsers struct {
		gorm.Model
		UsersID         uint   `gorm:"unique;not null"`
		Username        string `gorm:"unique;not null"`
		FirstNameThai   string
		LastNameThai    string
		FirstNameEng    string
		LastNameEng     string
		Born            time.Time
		Height          uint
		Weight          uint
		Sex             string
		Position        string
		Nationality     string
		Description     string
		Phone           string
		AddressesID     uint `gorm:"unique"`
		Addresses       Addresses
		Teams           []TeamsMembers
		TeamsCreated    []Teams             `gorm:"foreignKey:OwnerID;references:users_id"`
		GoalRecords     []GoalRecords       `gorm:"foreignKey:PlayerID"`
		RequestReceives []AddMemberRequests `gorm:"foreignKey:receiver_id;references:users_id"`
		ImagePath       string
	}

	Teams struct {
		gorm.Model
		Name         string
		OwnerID      uint
		Description  string
		TeamsMembers []TeamsMembers
		Compatitions []Compatitions      `gorm:"many2many:compatition_teams;"`
		RequestSends []AddMemberRequests `gorm:"foreignKey:teams_id;"`
	}

	TeamsMembers struct {
		gorm.Model
		TeamsID       uint
		Teams         Teams `gorm:"foreignKey:TeamsID;references:ID"`
		NormalUsersID uint
		NormalUsers   NormalUsers `gorm:"foreignKey:NormalUsersID;references:ID"`
		Role          string
	}

	CompatitionsTeams struct {
		gorm.Model
		TeamsID        uint
		Teams          Teams
		CompatitionsID uint
		Compatitions   Compatitions
	}

	Compatitions struct {
		gorm.Model
		Name              string
		Format            CompetitionFormat
		OrganizersID      uint
		Organizers        Organizers
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
		AddressesID       uint `gorm:"unique"`
		Addresses         Addresses
		NumberOfTeam      uint
		Teams             []Teams `gorm:"many2many:compatition_teams;"`
		NumOfPlayerMin    uint
		NumOfPlayerMax    uint
	}

	Events struct {
		gorm.Model
		MatchesID   uint
		Time        string
		Description string
		RedCard     bool
		YellowCard  bool
	}

	Matches struct {
		gorm.Model
		CompetitionID uint
		DateTime      time.Time
		Team1ID       uint `gorm:"foreignKey:TeamsID"`
		Team2ID       uint `gorm:"foreignKey:TeamsID"`
		Team1Goals    int
		Team2Goals    int
		Events        []Events    `gorm:"foreignKey:MatchesID"`
		GoalRecords   GoalRecords `gorm:"foreignKey:MatchesID"`
		Result        MatchesResult
	}

	GoalRecords struct {
		gorm.Model
		MatchesID  uint
		TeamsID    uint
		PlayerID   uint
		TimeScored uint
	}

	AddMemberRequests struct {
		gorm.Model
		TeamsID    uint
		Teams      Teams
		ReceiverID uint
		Role       string
		Status     string
	}
)

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
