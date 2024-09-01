package model

import "time"

type NextMatch struct {
	RivalTeamID           uint      `json:"rival_team_id"`
	RivalTeamName         string    `json:"rival_team_name"`
	RivalTeamImageProfile string    `json:"rival_team_image_profile"`
	RivalTeamImageCover   string    `json:"rival_team_image_cover"`
	CompatitionsID        uint      `json:"compatition_id"`
	CompatitionsName      string    `json:"compatition_name"`
	CompatitionsAddress   Address   `json:"compatition_address"`
	MatchID               uint      `json:"match_id"`
	MatchDateTime         time.Time `json:"match_date_time"`
}

type GoalRecord struct {
	MatchsID   uint `json:"matches_id" binding:"required"`
	TeamID     uint `json:"team_id" binding:"required"`
	PlayerID   uint `json:"player_id" binding:"required"`
	TimeScored uint `json:"time_scored" binding:"required"`
}

type RecentMatch struct {
	ID             uint      `json:"id"`
	DateTime       time.Time `json:"date_time"`
	VsTeamName     string    `json:"vs_team_name"`
	Result         string    `json:"result"`
	Score          string    `json:"score"` // 1 - 1, 3 - 1
	TournamentName string    `json:"tournament_name"`
}

type Match struct {
	ID             uint         `json:"id"`
	Index          int          `json:"index"`
	CompatitionsID uint         `json:"compatition_id"`
	DateTime       time.Time    `json:"date_time"`
	Team1ID        uint         `json:"team1_id"`
	Team2ID        uint         `json:"team2_id"`
	Team1Name      string       `json:"team1_name"`
	Team2Name      string       `json:"team2_name"`
	Team1Goals     int          `json:"team1_goals"`
	Team2Goals     int          `json:"team2_goals"`
	Round          string       `json:"round"`
	NextMatchIndex int          `json:"next_match_index"`
	NextMatchSlot  string       `json:"next_match_slot"` //Team1 or Team2
	GoalRecords    []GoalRecord `json:"goal_records"`
	Result         string       `json:"result"`
	VideoURL       string       `json:"video_url"`
	Team1Player    []Member     `json:"team1_player"`
	Team2Player    []Member     `json:"team2_player"`
}

type UpdateMatch struct {
	DateTime    time.Time    `json:"date_time"`
	Team1Goals  int          `json:"team1_goals"`
	Team2Goals  int          `json:"team2_goals"`
	GoalRecords []GoalRecord `json:"goal_records"`
	VideoURL    string       `json:"video_url"`
	Result      string       `json:"result"`
}