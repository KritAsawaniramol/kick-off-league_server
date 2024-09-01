package main

import (
	"fmt"

	"kickoff-league.com/config"
	"kickoff-league.com/database"
	"kickoff-league.com/migrations"
	"kickoff-league.com/server"
)

var CFG config.Config

func main() {
	CFG = config.GetConfig()

	
	db := database.NewPostgresDatabase(&CFG)
	
	//migration
	migrations.Migration(db)
	
	server.NewGinServer(&CFG, db.GetDb()).Start()
	
	fmt.Printf("CFG: %v\n", CFG)
}
