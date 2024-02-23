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
		&entities.Team{},
		&entities.NormalUser{},
		&entities.Organizer{},
		&entities.TeamMember{},
		&entities.CompatitionTeam{},
		&entities.Compatition{},
		&entities.Matches{},
		&entities.GoalRecord{},
	)
	if err != nil {
		panic("Database migration failed!")
	}
	fmt.Println("Database migration completed!")
}
