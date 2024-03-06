package main

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"kickoff-league.com/config"
	"kickoff-league.com/database"
	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
	"kickoff-league.com/usecases"
)

func main() {
	cfg := config.GetConfig()

	db := database.NewPostgresDatabase(&cfg)

	userPostgresRepository := repositories.NewUserPostgresRepository(db.GetDb())

	userUsercase := usecases.NewUserUsercaseImpl(
		userPostgresRepository,
	)

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
	n := 10
	for i := 1; i <= n; i++ {
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
			defer log.Errorf(err.Error())
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
			Sex:           "unisex",
			Position:      "Foward",
			Nationality:   "Thailand",
			Description:   "Description" + fmt.Sprintf("%d", i),
		}, uint(i)); err != nil {
			defer log.Errorf(err.Error())
			panic(err.Error())
		}
	}

	//Create team
	r := 'a'
	for i := 1; i <= n; i++ {
		teamName := fmt.Sprintf("Team%d", i)
		if err := u.CreateTeam(&model.CreateTeam{
			Name:        fmt.Sprintf("%c", r) + teamName,
			OwnerID:     uint(i),
			Description: "Description" + fmt.Sprintf("%d", i),
		}); err != nil {
			defer log.Errorf(err.Error())
			panic(err)
		}
		r++
	}

	//Add member to each team
	// for i := 0; i < n; i++ {
	// 	u.SendAddMemberRequest(&model.AddMemberRequest{
	// 		TeamID           :,
	// 		ReceiverUsername :,
	// 	})
	// }

	// Create Organizer
	n = 10
	for i := 0; i < n; i++ {
		u.RegisterOrganizer(&model.RegisterOrganizer{
			// Email:    fmt.Sprintf("organizer%d@gmail.com", i),
			// Password: "1234",
			// Phone:    "00000000000" + fmt.Sprintf("%d", i),
			RegisterUser: model.RegisterUser{
				Email:    fmt.Sprintf("organizer%d@gmail.com", i),
				Password: "1234",
			},
			OrganizerName: fmt.Sprintf("org%d", i),
			Phone:         "00000000000" + fmt.Sprintf("%d", i),
		})
	}

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
// 		&entities.Matches{},
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
