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

	Users struct {
		gorm.Model
		Email            string `gorm:"unique;not null;type:varchar(100)"`
		Role             string
		Password         string
		ImageProfilePath string
		ImageCoverPath   string
	}

	Organizers struct {
		gorm.Model
		UsersID      uint `gorm:"unique;not null"`
		Users        Users
		Name         string
		Phone        string `gorm:"unique;not null"`
		Description  string
		AddressesID  uint `gorm:"unique"`
		Addresses    Addresses
		Competitions []Competitions
	}

	NormalUsers struct {
		gorm.Model
		UsersID         uint `gorm:"unique;not null"`
		Users           Users
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
		TeamsCreated    []Teams       `gorm:"foreignKey:OwnerID;references:users_id"`
		GoalRecords     []GoalRecords `gorm:"foreignKey:PlayerID"`
		Competitions    []NormalUsersCompetitions
		RequestReceives []AddMemberRequests `gorm:"foreignKey:receiver_id;references:users_id"`
	}

	TeamsMembers struct {
		gorm.Model
		TeamsID       uint
		Teams         Teams `gorm:"foreignKey:TeamsID;references:ID"`
		NormalUsersID uint
		NormalUsers   NormalUsers `gorm:"foreignKey:NormalUsersID;references:ID"`
		Role          string
	}

	NormalUsersCompetitions struct {
		gorm.Model
		NormalUsersID  uint
		NormalUsers    NormalUsers `gorm:"foreignKey:NormalUsersID;references:ID"`
		CompetitionsID uint
		Competitions   Competitions `gorm:"foreignKey:CompetitionsID;references:ID"`
		TeamsID        uint
	}

	Teams struct {
		gorm.Model
		Name             string
		OwnerID          uint //userID
		Description      string
		TeamsMembers     []TeamsMembers
		Competitions     []CompetitionsTeams
		RequestSends     []AddMemberRequests `gorm:"foreignKey:teams_id;"`
		ImageProfilePath string
		ImageCoverPath   string
	}

	CompetitionsTeams struct {
		gorm.Model
		TeamsID        uint
		Teams          Teams `gorm:"foreignKey:TeamsID;references:ID"`
		CompetitionsID uint
		Competitions   Competitions `gorm:"foreignKey:CompetitionsID;references:ID"`
		Rank           string
		RankNumber     int
		Point          int
	}

	Address struct {
		HouseNumber string
		Village     string
		Subdistrict string
		District    string
		PostalCode  string
		Country     string
	}

	Competitions struct {
		gorm.Model
		Name                 string
		Sport                string //football or futsal
		Format               string
		Type                 string
		OrganizersID         uint
		Organizers           Organizers
		FieldSurface         string
		ApplicationType      string
		HouseNumber          string
		Village              string
		Subdistrict          string
		District             string
		PostalCode           string
		Country              string
		StartDate            time.Time
		EndDate              time.Time
		JoinCode             []JoinCode
		Description          string
		Rule                 string
		Prize                string
		ContactType          string
		Contact              string
		Status               string
		Teams                []CompetitionsTeams
		NumOfPlayerInTeamMin uint
		NumOfPlayerInTeamMax uint
		NormalUsers          []NormalUsersCompetitions

		AgeOver  uint
		AgeUnder uint
		Sex      string

		NumberOfTeam uint
		NumOfRound   int
		NumOfMatch   int
		Matchs       []Matchs

		ImageBannerPath string
	}

	JoinCode struct {
		gorm.Model
		CompetitionsID uint
		Code           string
		Status         string //used, unused
	}

	Matchs struct {
		gorm.Model
		Index          int //start with 1
		CompetitionsID uint
		Competitions   Competitions
		DateTime       time.Time
		Team1ID        uint `gorm:"foreignKey:TeamsID"`
		Team2ID        uint `gorm:"foreignKey:TeamsID"`
		Team1Goals     int
		Team2Goals     int
		Round          string        // Round 1, Round 2
		GoalRecords    []GoalRecords `gorm:"foreignKey:MatchsID"`
		VideoURL       string
		NextMatchIndex int
		NextMatchSlot  string //Team1 or Team2
		Result         string
	}

	GoalRecords struct {
		gorm.Model
		MatchsID   uint
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
