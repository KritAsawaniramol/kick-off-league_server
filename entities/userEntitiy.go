package entities

import (
	"time"

	"gorm.io/gorm"
)

type (
	Address struct {
		gorm.Model
		HouseNumber string
		Village     string
		Subdistrict string
		District    string
		PostalCode  string
		Country     string
	}

	User struct {
		gorm.Model
		Email    string `gorm:"unique;not null;type:varchar(100)"`
		Role     string
		Password string
	}

	Organizer struct {
		gorm.Model
		UserID      uint `gorm:"unique;not null"`
		User        User
		Name        string
		Phone       string
		Description string
		AddressID   uint `gorm:"unique"`
		Address     Address
		Compatition []Compatition
		// imageProfile string
	}

	NormalUserTeam struct {
		UserID uint
		User   User
		TeamID uint
		Team   Team
	}

	TeamMember struct {
		TeamID       uint
		Team         Team
		NormalUserID uint
		NormalUser   NormalUser
	}

	NormalUser struct {
		gorm.Model
		UserID        uint `gorm:"unique;not null"`
		User          User
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
		Description   string
		Phone         string
		AddressID     uint `gorm:"unique"`
		Address       Address
		Team          []Team     `gorm:"many2many:team_member;"`
		TeamCreated   []Team     `gorm:"foreignKey:OwnerID";references:user_id`
		GoalRecord    GoalRecord `gorm:"foreignKey:PlayerID"`
		// imageProfile  string
	}
	Team struct {
		gorm.Model
		Name        string
		Member      []NormalUser `gorm:"many2many:team_member;"`
		OwnerID     uint
		Compatition []Compatition `gorm:"many2many:compatition_team;"`
	}

	CompatitionTeam struct {
		gorm.Model
		TeamID        uint
		Team          Team
		CompatitionID uint
		Compatition   Compatition
	}

	Compatition struct {
		gorm.Model
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

	Event struct {
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
		Team1ID       uint `gorm:"foreignKey:TeamID"`
		Team2ID       uint `gorm:"foreignKey:TeamID"`
		Team1Goals    int
		Team2Goals    int
		Events        []Event    `gorm:"foreignKey:MatchesID"`
		GoalRecord    GoalRecord `gorm:"foreignKey:MatchesID"`
		Result        MatchesResult
	}

	GoalRecord struct {
		gorm.Model
		MatchesID  uint
		TeamID     uint
		PlayerID   uint
		TimeScored uint
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
