package main

import (
	"fmt"
	"time"

	"kickoff-league.com/config"
	"kickoff-league.com/database"
	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
	"kickoff-league.com/usecases"
	"kickoff-league.com/util"
)

func main() {
	cfg := config.GetConfig()

	db := database.NewPostgresDatabase(&cfg)

	userPostgresRepository := repositories.NewUserPostgresRepository(db.GetDb())

	userUsercase := usecases.NewUserUsercaseImpl(
		userPostgresRepository,
	)
	// c := &entities.Compatitions{}
	// c.ID = 3
	// t := &entities.Teams{}
	// t.ID = 1
	// userPostgresRepository.AppendTeamtoCompatition(c, t)
	// userUsercase.StartCompatition(3)
	// err := userUsercase.JoinCompatition(&model.JoinCompatition{
	// 	CompatitionID: 3,
	// 	TeamID:        1,
	// 	Code:          "",
	// })

	// util.PrintObjInJson(normalUser)

	// err := userPostgresRepository.DeleteTeamMember(&entities.TeamsMembers{
	// 	TeamsID:       1,
	// 	NormalUsersID: 1,
	// })
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	panic(err)
	// }

	// userUsercase.CreateJoinCode(1, 3)

	mockupData(userUsercase)

	// team := entities.Teams{}
	// db.GetDb().Preload("TeamsMembers").First(&team, 1)

	// util.PrintObjInJson(team)

	// team, err := userPostgresRepository.GetTeam(1)
	// if err != nil {
	// 	log.Error(err)
	// }

	// normalUser, err := userPostgresRepository.GetNormalUserByUserID(2)
	// if err != nil {
	// 	log.Error(err)
	// }
	// util.PrintObjInJson(team)

	// normalUser.Teams = append(normalUser.Teams, *team)

	// log.Print("hello world")
	// if err := userPostgresRepository.UpdateNormalUser(normalUser); err != nil {
	// 	log.Error(err)
	// }

	// err := userPostgresRepository.InsertTeamsMembers(&entities.TeamsMembers{
	// 	TeamsID:       3,
	// 	NormalUsersID: 1,
	// 	Role:          "player",
	// })
	// if err != nil {
	// 	log.Error(err)
	// }

	// teams, err := userUsercase.GetTeamList(&model.GetTeamList{})
	// if err != nil {
	// 	log.Error(err)
	// } else {
	// 	util.PrintObjInJson(teams)
	// }

	// teams, err := userUsercase.GetTeamWithMemberAndCompatitionByID(1)
	// if err != nil {
	// 	log.Error(err)
	// }

	// count := userPostgresRepository.GetNumberOfTeamsMember(TestGetTeams(userUsercase, 0, 1))
	// fmt.Printf("count: %d\n", count)

	// teams, err := userPostgresRepository.GetTeams(TestGetTeams(userUsercase, 0), "name", false, -1, -1)
	// if err != nil {
	// 	log.Error(err)
	// }
	// util.PrintObjInJson(teams)

	// team, err := userPostgresRepository.GetTeamWithMemberByID(uint(1))
	// if err != nil {
	// 	log.Error(err)
	// }

	// team_byte, err := json.MarshalIndent(team, "", "    ")
	// if err != nil {
	// 	log.Error(err)
	// }

	// log.Print(team_byte)

	// if team, err := userPostgresRepository.GetTeamWithMemberAndRequestSendByID(uint(1)); err != nil {
	// 	log.Errorf(err.Error())
	// 	panic(err)
	// } else {
	// 	util.PrintObjInJson(team)
	// }

	// userMigrate(db)

	// CFG = config.GetConfig()

	// db := database.NewPostgresDatabase(&CFG)

	// //migration
	// migrations.Migration(db)

	// server.NewGinServer(&CFG, db.GetDb()).Start()
}

func TestGetTeams(u usecases.UserUsecase, normalUserID uint, teamID uint) *entities.Teams {
	team := &entities.Teams{
		OwnerID: normalUserID,
	}
	team.ID = teamID
	return team
}

func mockupData(u usecases.UserUsecase) {
	// Create NormalUser
	nOfNormalUser := 20
	for i := 1; i <= nOfNormalUser; i++ {
		email := fmt.Sprintf("normal%d@gmail.com", i)
		password := "1234"
		username := fmt.Sprintf("normal%d", i)
		if err := u.RegisterNormaluser(
			&model.RegisterNormaluser{
				RegisterUser: model.RegisterUser{
					Email:    email,
					Password: password,
				},
				Username: username,
			},
		); err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err.Error())
		}

		// Create Assign info in each NormalUser Account
		if err := u.UpdateNormalUser(&model.UpdateNormalUser{
			FirstNameThai: "FirstNameThai" + fmt.Sprintf("%d", i),
			LastNameThai:  "LastNameThai" + fmt.Sprintf("%d", i),
			FirstNameEng:  "FirstNameEng" + fmt.Sprintf("%d", i),
			LastNameEng:   "LastNameEng" + fmt.Sprintf("%d", i),
			Born:          time.Now(),
			Phone:         "00000000000" + fmt.Sprintf("%d", i),
			Height:        175,
			Weight:        70,
			Sex:           "Male",
			Position:      "Foward",
			Nationality:   "Thailand",
			Description:   "Description" + fmt.Sprintf("%d", i),
		}, uint(i)); err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err.Error())
		}
	}

	//Create team
	r := 'a'
	nOfTeam := 4
	for i := 1; i <= nOfNormalUser; i += nOfNormalUser / nOfTeam {
		teamName := fmt.Sprintf("Team%d", i)
		if err := u.CreateTeam(&model.CreateTeam{
			Name:        fmt.Sprintf("%c", r) + teamName,
			OwnerID:     uint(i),
			Description: "Description" + fmt.Sprintf("%d", i),
		}); err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err)
		}
		r++
	}

	// Add member to each team
	teamID := 1
	ownerID := 1
	for i := 1; i <= nOfNormalUser; i++ {
		fmt.Printf("i: %v\n", i)
		err := u.SendAddMemberRequest(&model.AddMemberRequest{
			TeamID:           uint(teamID),
			ReceiverUsername: fmt.Sprintf("normal%d", i),
			Role:             "player",
		}, uint(ownerID))
		if err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err)
		}
		if i%5 == 0 {
			ownerID += 5
			teamID += 1
		}
	}

	// Accept request
	for i := 1; i <= nOfNormalUser; i++ {
		u.AcceptAddMemberRequest(uint(i), uint(i))
	}

	// Create Organizer
	nOfOrg := 2
	for i := 0; i < nOfOrg; i++ {
		u.RegisterOrganizer(&model.RegisterOrganizer{
			RegisterUser: model.RegisterUser{
				Email:    fmt.Sprintf("organizer%d@gmail.com", i),
				Password: "1234",
			},
			Phone:         "00000000000" + fmt.Sprintf("%d", i),
			OrganizerName: fmt.Sprintf("org%d", i),
		})
	}

	// Create Compatition
	// err := u.CreateCompatition(&model.CreateCompatition{
	// 	Name:                 "Football A",
	// 	Sport:                "Football",
	// 	Type:                 "Round Robin",
	// 	Format:               "5vs5",
	// 	Description:          "Description",
	// 	Rule:                 "Rule",
	// 	Prize:                "Prize",
	// 	StartDate:            time.Now().AddDate(0, 0, -4),
	// 	EndDate:              time.Now(),
	// 	ApplicationType:      "free",
	// 	ImageBanner:          "",
	// 	AgeOver:              0,
	// 	AgeUnder:             0,
	// 	Sex:                  "Unisex",
	// 	NumberOfTeam:         4,
	// 	NumOfPlayerInTeamMin: 0,
	// 	NumOfPlayerInTeamMax: 0,
	// 	FieldSurface:         "NaturalGrass",
	// 	OrganizerID:          1,
	// 	Address: model.Address{
	// 		HouseNumber: "HouseNumber",
	// 		Village:     "Village    ",
	// 		Subdistrict: "Subdistrict",
	// 		District:    "District   ",
	// 		PostalCode:  "PostalCode ",
	// 		Country:     "Country    ",
	// 	},
	// 	ContractType: "facebook",
	// 	Contract:     "facbook URL",
	// })
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	panic(err)
	// }

	err := u.CreateCompatition(&model.CreateCompatition{
		Name:                 "Football B",
		Sport:                "Futsal",
		Type:                 util.CompetitionFormat[0],
		Format:               "5vs5",
		Description:          "Description",
		Rule:                 "Rule",
		Prize:                "Prize",
		StartDate:            time.Now().AddDate(0, 0, -4),
		EndDate:              time.Now(),
		ApplicationType:      "free",
		ImageBanner:          "",
		AgeOver:              0,
		AgeUnder:             0,
		Sex:                  "Unisex",
		NumberOfTeam:         8,
		NumOfPlayerInTeamMin: 0,
		NumOfPlayerInTeamMax: 0,
		FieldSurface:         "NaturalGrass",
		OrganizerID:          1,
		Address: model.Address{
			HouseNumber: "HouseNumber",
			Village:     "Village    ",
			Subdistrict: "Subdistrict",
			District:    "District   ",
			PostalCode:  "PostalCode ",
			Country:     "Country    ",
		},
		ContractType: "facebook",
		Contract:     "facbook URL",
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}

	// Open Application
	err = u.OpenApplicationCompatition(1)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}

	//join compatition
	nOfTeam -= 1
	for i := 1; i <= nOfTeam; i++ {
		err := u.JoinCompatition(&model.JoinCompatition{
			CompatitionID: 1,
			TeamID:        uint(i),
		})
		if err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err)
		}
	}

	// Start compatition
	u.StartCompatition(1)

	// UpdateMatch
	u.UpdateMatch(1, &model.UpdateMatch{
		DateTime:   time.Now().AddDate(0, 0, -1),
		Team1Goals: 2,
		Team2Goals: 0,
		GoalRecords: []model.GoalRecord{
			model.GoalRecord{
				MatchsID:   1,
				TeamID:     4,
				PlayerID:   5,
				TimeScored: 45,
			},
		},
		Result: "Team1Win",
	})
	// Next Matcht
	// nextMatchs, err := u.GetNextMatch(1)
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	panic(err)
	// }

	// //GetNormalUser
	// normalUser, err := u.GetNormalUser(6)
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	panic(err)
	// }

	// util.PrintObjInJson(normalUser)
}

// func userMigrate(db database.Database) {

// 	db.GetDb().Migrator().CreateTable(
// 		&entities.User{},
// 		&entities.Address{},
// 		&entities.Teams{},
// 		&entities.NormalUser{},
// 		&entities.Organizer{},
// 		&entities.TeamsMember{},
// 		&entities.CompatitionTeams{},
// 		&entities.Compatition{},
// 		&entities.Matchs{},
// 		&entities.GoalRecords{},
// 		&entities.AddMemberRequest{},
// 	)
// 	db.GetDb().CreateInBatches([]entities.User{
// 		{
// 			Email:    "normal01@gmail.com",
// 			Role:     "normal",
// 			Password: "1234",
// 		},
// 		{
// 			Email:    "normal02@gmail.com",
// 			Role:     "normal",
// 			Password: "1234",
// 		},
// 		{
// 			Email:    "normal03@gmail.com",
// 			Role:     "normal",
// 			Password: "1234",
// 		},
// 		{
// 			Email:    "normal03@gmail.com",
// 			Role:     "normal",
// 			Password: "1234",
// 		},
// 	}, 4)

// 	db.GetDb().CreateInBatches([]entities.User{
// 		{
// 			Email:    "normal01@gmail.com",
// 			Role:     "normal",
// 			Password: "1234",
// 		},
// 		{
// 			Email:    "normal02@gmail.com",
// 			Role:     "normal",
// 			Password: "1234",
// 		},
// 		{
// 			Email:    "normal03@gmail.com",
// 			Role:     "normal",
// 			Password: "1234",
// 		},
// 		{
// 			Email:    "normal03@gmail.com",
// 			Role:     "normal",
// 			Password: "1234",
// 		},
// 	}, 4)

// }
