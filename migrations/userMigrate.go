package migrations

import (
	"kickoff-league.com/config"
	"kickoff-league.com/database"
	"kickoff-league.com/entities"
)

func main() {
	cfg := config.GetConfig()

	db := database.NewPostgresDatabase(&cfg)

	userMigrate(db)
}

func userMigrate(db database.Database) {
	db.GetDb().Migrator().CreateTable(&entities.User{})
	db.GetDb().CreateInBatches([]entities.User{
		{
			Email:    "heloda8243@fahih.com",
			Role:     "Normal",
			Password: "123",
		},
		{
			Email:    "flaviog@mac.com",
			Role:     "Organizer",
			Password: "324",
		},
		{
			Email:    "cgreuter@comcast.net",
			Role:     "Normal",
			Password: "908",
		},
		{
			Email:    "bartak@yahoo.com",
			Role:     "Normal",
			Password: "323",
		},
	}, 4)

}
