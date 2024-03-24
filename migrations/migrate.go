package migrations

import (
	"fmt"

	"kickoff-league.com/database"
	"kickoff-league.com/entities"
)

func Migration(db database.Database) {
	err := db.GetDb().AutoMigrate(
		&entities.Users{},
		&entities.TeamsMembers{},
		&entities.Teams{},
		&entities.NormalUsers{},
		&entities.Organizers{},
		&entities.Compatitions{},
		&entities.CompatitionsTeams{},
		&entities.Addresses{},
		&entities.Matches{},
		&entities.GoalRecords{},
		&entities.AddMemberRequests{},
		&entities.JoinCode{},
	)

	db.GetDb().Migrator().AddColumn(&entities.TeamsMembers{}, "Status")
	if err != nil {
		panic("Database migration failed!")
	}
	fmt.Println("Database migration completed!")
}
