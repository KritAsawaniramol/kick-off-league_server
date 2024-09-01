package model

// type MyNormalUser struct {
// 	NormalUser
// }

type User struct {
	ID               uint            `json:"id"`
	Email            string          `json:"email"`
	Role             string          `json:"role"`
	NormalUserInfo   *NormalUserInfo `json:"normal_user,omitempty"`
	OrganizersInfo   *OrganizersInfo `json:"organizer,omitempty"`
	ImageProfilePath string          `json:"image_profile_path"`
	ImageCoverPath   string          `json:"image_cover_path"`
}

// type NormalUser struct {
// 	NormalUserInfo `json:"normal_user_info"`
// 	Team           []Team     `json:"team"`
// 	TeamCreated    []Team     `json:"team_created"`
// 	GoalRecord     GoalRecord `json:"goal_record"`
// }

type Address struct {
	HouseNumber string `json:"house_number"`
	Village     string `json:"village"`
	Subdistrict string `json:"subdistrict"`
	District    string `json:"district"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
}

// type CompetitionType string

// const (
// 	Tournament CompetitionType = "Tournament"
// 	RoundRobin CompetitionType = "Round Robin"
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

// type MatchsResult string

// const (
// 	Team1Win MatchsResult = "Team1Win"
// 	Team2Win MatchsResult = "Team2Win"
// 	Draw     MatchsResult = "Draw"
// )

// var SexType = []string{"Male", "Female", "Unisex"}
// var FieldSurfaces = []string{"NaturalGrass", "ArtificialTurf", "FlatSurface"}
// var CompetitionType = []string{"Tournament", "Round Robin"}
// var MatchsResult = []string{"Team1Win", "Team2Win", "Draw"}
