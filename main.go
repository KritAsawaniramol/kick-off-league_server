package main

import (
	"kickoff-league.com/config"
	"kickoff-league.com/database"
	"kickoff-league.com/migrations"
	"kickoff-league.com/server"
	"kickoff-league.com/util"
)

var CFG config.Config

func main() {
	CFG = config.GetConfig()

	db := database.NewPostgresDatabase(&CFG)

	//migration
	migrations.Migration(db)

	server.NewGinServer(&CFG, db.GetDb()).Start()

}
