package migrations

import (
	"fmt"

	"kickoff-league.com/database"
	"kickoff-league.com/entities"
)

func Migration(db database.Database) {
	err := db.GetDb().AutoMigrate(
		&entities.User{},
		&entities.Address{},
		&entities.Teams{},
		&entities.NormalUser{},
		&entities.Organizer{},
		// &entities.TeamMember{},
		&entities.CompatitionTeams{},
		&entities.Compatition{},
		&entities.Matches{},
		&entities.GoalRecords{},
		&entities.AddMemberRequest{},
	)
	if err != nil {
		panic("Database migration failed!")
	}
	fmt.Println("Database migration completed!")
}
