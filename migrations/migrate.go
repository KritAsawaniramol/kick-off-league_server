package migrations

import (
	"fmt"

	"kickoff-league.com/database"
	"kickoff-league.com/entities"
)

func Migration(db database.Database) {

	err := db.GetDb().AutoMigrate(
		&entities.NormalUsersCompetitions{},
		&entities.Users{},
		&entities.TeamsMembers{},
		&entities.Teams{},
		&entities.NormalUsers{},
		&entities.Organizers{},
		&entities.CompetitionsTeams{},
		&entities.Competitions{},
		&entities.Addresses{},
		&entities.Matchs{},
		&entities.GoalRecords{},
		&entities.AddMemberRequests{},
		&entities.JoinCode{},
	)
	if err != nil {
		panic(err)
	}

	db.GetDb().Migrator().AddColumn(&entities.TeamsMembers{}, "Status")
	if err != nil {
		panic("Database migration failed!")
	}
	fmt.Println("Database migration completed!")
}
