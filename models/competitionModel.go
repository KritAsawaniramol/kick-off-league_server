package model

import "time"

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
	ImageBanner  string `json:"image_banner"`
	Rank         string `json:"rank"`
	RankNumber   int    `json:"rank_number"`
	Sport        string `json:"sport"`
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
	ImageBanner     string    `json:"image_banner"`
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
	ContactType          string         `json:"Contact_type"  binding:"required"`
	Contact              string         `json:"Contact"  binding:"required"`
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
	Matchs               []Match        `json:"match"`
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
	ContactType          string    `json:"Contact_type" `
	Contact              string    `json:"Contact" `
	AgeOver              uint      `json:"age_over"`
	AgeUnder             uint      `json:"age_under"`
	Sex                  string    `json:"sex"`
	NumberOfTeam         uint      `json:"number_of_team"`
	NumOfPlayerInTeamMin uint      `json:"num_of_player_min"`
	NumOfPlayerInTeamMax uint      `json:"num_of_player_max"`
}

type JoinCode struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Code      string    `json:"code"`
	Status    string    `json:"status"`
}

type CreateCompetition struct {
	Name   string `json:"name" binding:"required"`
	Sport  string `json:"sport" binding:"required"`
	Type   string `json:"type" binding:"required"`   // Tournament, Round Robbin,....
	Format string `json:"format" binding:"required"` // 1vs1, 2vs2,...

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

	ContactType string `json:"Contact_type"  binding:"required"`
	Contact     string `json:"Contact"  binding:"required"`
}



type JoinCompetition struct {
	CompetitionID uint   `json:"compatition_id" binding:"required"`
	TeamID        uint   `json:"team_id" binding:"required"`
	Code          string `json:"code"`
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

type CompatitionList struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Description    string `jsFcon:"description"`
	NumberOfMember uint   `json:"number_of_member"`
	// img         string
}
