package main

import (
	"fmt"
	"log"
	"time"

	"kickoff-league.com/config"
	"kickoff-league.com/database"
	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
	"kickoff-league.com/usecases/addMemberUsecase"
	"kickoff-league.com/usecases/authUsecase"
	"kickoff-league.com/usecases/competitionUsecase"
	"kickoff-league.com/usecases/matchUsecase"
	"kickoff-league.com/usecases/normalUserUsecase"
	"kickoff-league.com/usecases/organizerUsecase"
	"kickoff-league.com/usecases/teamUsecase"

	"kickoff-league.com/usecases/userUsecase"
	"kickoff-league.com/util"
)

var FirstNames = [41]string{
	"John",
	"Michael",
	"David",
	"James",
	"Robert",
	"William",
	"Joseph",
	"Charles",
	"Richard",
	"Richard",
	"Thomas",
	"Daniel",
	"Matthew",
	"Christopher",
	"Anthony",
	"Brian",
	"Steven",
	"Timothy",
	"Kevin",
	"Mark",
	"Paul",
	"Andrew",
	"Edward",
	"Jason",
	"Scott",
	"Justin",
	"Ryan",
	"Eric",
	"Gregory",
	"Joshua",
	"Kenneth",
	"Jeffrey",
	"Stephen",
	"Brandon",
	"Jonathan",
	"Larry",
	"Dennis",
	"Jerry",
	"Tyler",
	"Frank",
}

var LastNames = [41]string{
	"Smith",
	"Johnson",
	"Williams",
	"Brown",
	"Jones",
	"Garcia",
	"Miller",
	"Davis",
	"Rodriguez",
	"Martinez",
	"Hernandez",
	"Lopez",
	"Gonzalez",
	"Wilson",
	"Anderson",
	"Anderson",
	"Thomas",
	"Taylor",
	"Moore",
	"Jackson",
	"Martin",
	"Lee",
	"Perez",
	"Thompson",
	"White",
	"Harris",
	"Sanchez",
	"Clark",
	"Ramirez",
	"Lewis",
	"Robinson",
	"Walker",
	"Young",
	"Hall",
	"Allen",
	"King",
	"Wright",
	"Scott",
	"Torres",
	"Nguyen",
	"Hill",
}

var ThaiFirstNames = [41]string{
	"สุรชัย",
	"วรรณพ",
	"พรพิมล",
	"ธนวรรธน์",
	"วรวัชร",
	"รัตนาภรณ์",
	"ชนกานต์",
	"ศุภวิชญ์",
	"ประพฤติ",
	"สมบัติ",
	"อริศรา",
	"อริศรา",
	"ปภัสสร",
	"วิรัช",
	"จิราภรณ์",
	"ศิริวัฒน์",
	"ณัฐพงศ์",
	"จันทร์ทอง",
	"พีระพงษ์",
	"สรรเสริญ",
	"พิชญะ",
	"วสุธา",
	"วิรัตน์",
	"ภัทรวิชญ์",
	"ปิยะพงษ์",
	"ภูมิพัฒน์",
	"สมชาย",
	"ธนกฤต",
	"ชวิน",
	"วรวุฒิ",
	"อนันต์",
	"ณรงค์",
	"วิภาวดี",
	"พัชรพล",
	"ภาณุวัฒน์",
	"ปรัชญ์",
	"สิรวิชญ์",
	"วีระพล",
	"จิรวัฒน์",
	"พงศกร",
}

var ThaiLastNames = [41]string{
	"สุวรรณ",
	"วงศ์วาริน",
	"พิมพ์พรรณ",
	"ธนาวินท์",
	"พิชัย",
	"พิชัย",
	"นาคนิรมล",
	"เทพกฤต",
	"วิชชุดา",
	"วงศ์ชัย",
	"ทองดี",
	"เกษม",
	"ชัชวาลย์",
	"ทองแดง",
	"ชาญชัย",
	"พันธุ์สุข",
	"บุญชู",
	"ประดับ",
	"ธีระวัฒน์",
	"จันทร์ดี",
	"ธีระ",
	"บุญมี",
	"สุวรรณ",
	"จันทร์รักษ์",
	"วงศ์สวัสดิ์",
	"สุวรรณภูมิ",
	"พลอยใส",
	"รัตนประดิษฐ์",
	"สุวรรณวิมล",
	"ปิยะสมบูรณ์",
	"วิชุดา",
	"กาญจน์",
	"พรหมพัฒน์",
	"สุริยะ",
	"ชินวัตร",
	"ประจักษ์",
	"รัตน์",
	"ชินวุฒิ",
	"สุริยะวงศ์",
	"ภูมิพงษ์",
}

var Position = [41]string{
	"Goalkeeper", "Defender", "Midfielder", "Forward",
	"Goalkeeper", "Defender", "Midfielder", "Forward",
	"Goalkeeper", "Defender", "Midfielder", "Forward",
	"Goalkeeper", "Defender", "Midfielder", "Forward",
	"Goalkeeper", "Defender", "Midfielder", "Forward",
	"Goalkeeper", "Defender", "Midfielder", "Forward",
	"Goalkeeper", "Defender", "Midfielder", "Forward",
	"Goalkeeper", "Defender", "Midfielder", "Forward",
	"Goalkeeper", "Defender", "Midfielder", "Forward",
	"Goalkeeper", "Defender", "Midfielder", "Forward",
	"Goalkeeper",
}

var PlayerDescriptions = [41]string{
	"Tall and strong defender known for his tough tackles.",
	"Agile goalkeeper with lightning reflexes.",
	"Dynamic midfielder with excellent passing skills.",
	"Pacey forward with a deadly finish.",
	"Versatile player who can play in multiple positions.",
	"Experienced captain known for his leadership on the pitch.",
	"Young talent with potential to become a future star.",
	"Veteran player with years of experience at the top level.",
	"Fan-favorite known for his flair and creativity.",
	"Solid defender who rarely puts a foot wrong.",
	"Hardworking midfielder who covers every blade of grass.",
	"Clinical finisher with a knack for scoring goals.",
	"Intelligent playmaker who dictates the tempo of the game.",
	"Energetic winger with blistering pace.",
	"Dependable goalkeeper who rarely makes mistakes.",
	"Physical midfielder who dominates the midfield battles.",
	"Skilful dribbler capable of beating multiple defenders.",
	"Tireless runner who never gives up on the ball.",
	"Stalwart defender who anchors the backline.",
	"Technical genius known for his exquisite ball control.",
	"Tactical mastermind who reads the game brilliantly.",
	"Speedy full-back who bombards forward with pace.",
	"Speedy full-back who bombards forward with pace.",
	"Crafty forward with an eye for goal.",
	"Commanding presence in defense with aerial dominance.",
	"Astute tactician who can adapt to different formations.",
	"Effervescent midfielder who brings energy to the team.",
	"Strong aerial presence during set pieces.",
	"Maestro in midfield orchestrating attacks.",
	"Reliable penalty taker with nerves of steel.",
	"Tenacious defender who never shies away from challenges.",
	"Playoff hero with a penchant for scoring crucial goals.",
	"Fan-favorite known for his work rate and dedication.",
	"Quick-footed winger capable of producing moments of magic.",
	"Disciplined player who follows the manager's instructions.",
	"Natural leader who leads by example on and off the pitch.",
	"Decisive in front of goal with clinical finishing.",
	"Versatile utility player capable of filling various roles.",
	"Promising young talent with a bright future ahead.",
	"Solid defensive midfielder who shields the backline.",
	"Unpredictable attacker with dazzling dribbling skills.",
}

var TeamNames = [41]string{
	"Manchester United",
	"Real Madrid",
	"FC Barcelona",
	"Bayern Munich",
	"Liverpool",
	"Juventus",
	"Paris Saint-Germain",
	"Chelsea",
	"Manchester City",
	"AC Milan",
	"Inter Milan",
	"Arsenal",
	"Borussia Dortmund",
	"Atletico Madrid",
	"Tottenham Hotspur",
	"Roma",
	"Ajax",
	"Boca Juniors",
	"River Plate",
	"River Plate",
	"Flamengo",
	"Santos",
	"Benfica",
	"Porto",
	"Sporting Lisbon",
	"Galatasaray",
	"Fenerbahce",
	"Besiktas",
	"Olympique de Marseille",
	"Olympique Lyonnais",
	"AS Monaco",
	"Valencia",
	"Sevilla",
	"Villarreal",
	"Napoli",
	"AS Roma",
	"Leicester City",
	"West Ham United",
	"Everton",
	"Newcastle United",
	"Leeds United",
}

var TeamDescriptions = [41]string{
	"A powerhouse in European football, known for its rich history and global fanbase.",
	"One of the most successful clubs in Spain, with numerous league titles and Champions League trophies.",
	"A symbol of Catalan identity, renowned for its possession-based style of play.",
	"A dominant force in German football, consistently competing at the highest level in Europe.",
	"A legendary English club with a strong tradition of success both domestically and in Europe.",
	"Italy's most successful club, recognized for its iconic black and white stripes.",
	"The top club in French football, backed by wealthy Qatari owners and star-studded squad.",
	"One of the leading clubs in England, boasting a strong track record of domestic and European success.",
	"A recent powerhouse in English football, backed by significant financial investment.",
	"A historic Italian club with a passionate fanbase and a tradition of success in Serie A.",
	"One of the giants of Italian football, known for its fierce rivalry with AC Milan.",
	"One of the most storied clubs in English football, renowned for its attractive style of play.",
	"A popular German club with a loyal fanbase and a history of success in the Bundesliga.",
	"A formidable force in Spanish football, consistently challenging for league titles and European trophies.",
	"A prominent English club with a strong following and a tradition of entertaining football.",
	"A historic Italian club with a proud tradition and a strong presence in Serie A.",
	"A Dutch powerhouse with a rich history of developing talented players and playing attractive football.",
	"One of the most successful clubs in Argentine football, with a passionate fanbase and intense rivalries.",
	"A powerhouse in Argentine football, known for its fierce rivalry with Boca Juniors.",
	"The most popular football club in Brazil, with a rich history and a tradition of producing talented players.",
	"The most popular football club in Brazil, with a rich history and a tradition of producing talented players.",
	"One of the most successful clubs in Brazilian football, boasting a rich history and passionate supporters.",
	"A Portuguese giant with a strong presence in domestic competitions and European competitions.",
	"A dominant force in Portuguese football, known for its successful youth academy and attractive style of play.",
	"A leading club in Portuguese football, with a rich history and a strong fanbase.",
	"A Turkish powerhouse with a passionate fanbase and a history of success in domestic competitions.",
	"One of Turkey's biggest clubs, with a fervent following and a tradition of success in domestic competitions.",
	"A prominent Turkish club with a passionate fanbase and a history of success in domestic competitions.",
	"A leading club in French football, with a passionate fanbase and a tradition of success in Ligue 1.",
	"One of the most successful clubs in French football, with a strong record of domestic and European success.",
	"A top club in French football, known for its attacking style of play and talented squad.",
	"A historic Spanish club with a strong tradition of success and a loyal fanbase.",
	"A prominent Spanish club with a rich history and a strong presence in European competitions.",
	"A rising force in Spanish football, with a talented squad and ambitious goals.",
	"One of Italy's top clubs, known for its passionate fanbase and a history of success in Serie A.",
	"A prominent Italian club with a tradition of success and a loyal fanbase.",
	"A rising club in Italian football, with a talented squad and aspirations of challenging the top teams.",
	"A prominent English club with a rich history and a tradition of entertaining football.",
	"A historic English club with a loyal fanbase and a tradition of success in domestic competitions.",
	"A Premier League club with a passionate fanbase and ambitions of competing at the highest level.",
	"A historic English club with a passionate fanbase and a tradition of attacking football.",
}

var CompetitionNames = [41]string{
	"FIFA World Cup",
	"UEFA Champions League",
	"English Premier League",
	"La Liga",
	"Bundesliga",
	"Serie A",
	"Ligue 1",
	"UEFA Europa League",
	"UEFA Europa League",
	"Copa Libertadores",
	"Copa America",
	"UEFA European Championship",
	"FIFA Club World Cup",
	"AFC Champions League",
	"CONCACAF Champions League",
	"Copa del Rey",
	"DFB-Pokal",
	"Coppa Italia",
	"Copa do Brasil",
	"FA Cup",
	"Copa Argentina",
	"Copa MX",
	"Scottish Premiership",
	"Eredivisie",
	"Campeonato Brasileiro Serie A",
	"Primeira Liga",
	"Belgian Pro League",
	"MLS Cup",
	"J1 League",
	"Russian Premier League",
	"Turkish Super Lig",
	"Czech First League",
	"Swiss Super League",
	"Argentine Primera Division",
	"A-League",
	"Bahraini Premier League",
	"Saudi Professional League",
	"Qatar Stars League",
	"Chinese Super League",
	"Indian Super League",
	"K-League",
}

var Clubs = [21]string{
	"arsenal",
	"aston-villa",
	"brentford",
	"brighton",
	"burnley",
	"chelsea",
	"crystal-palace",
	"everton",
	"leeds",
	"leicester-city",
	"liverpool",
	"manchester-city",
	"manchester-united",
	"newcastle",
	"norwich",
	"norwich",
	"southampton",
	"tottenham",
	"watford",
	"west-ham",
	"wolves",
}

func main() {
	cfg := config.GetConfig()

	db := database.NewPostgresDatabase(&cfg)

	repository := repositories.NewUserPostgresRepository(db.GetDb())

	userUsecase := userUsecase.NewUserUsecaseImpl(
		repository,
	)

	organizerUsecase := organizerUsecase.NewOrganizerUsecaseImpl(
		repository,
	)

	teamUsecase := teamUsecase.NewTeamUsecaseImpl(
		repository,
	)

	authUsecase := authUsecase.NewAuthUsecaseImpl(
		repository,
	)

	addMemberUsecase := addMemberUsecase.NewAddMemberUsecaseImpl(
		repository,
	)

	competitionUsecase := competitionUsecase.NewCompetitionUsecaseImpl(
		repository,
	)

	normalUserUsecase := normalUserUsecase.NewNormalUserUsecaseImpl(
		repository,
	)

	matchUsecase := matchUsecase.NewMatchUsecaseImpl(
		repository,
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

	// userPostgresRepository.ClearGoalRecordsOfMatch(1)

	mockupData(
		userUsecase,
		organizerUsecase,
		authUsecase,
		normalUserUsecase,
		teamUsecase,
		addMemberUsecase,
		competitionUsecase,
		matchUsecase,
	)

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

func TestGetTeams(normalUserID uint, teamID uint) *entities.Teams {
	team := &entities.Teams{
		OwnerID: normalUserID,
	}
	team.ID = teamID
	return team
}

func mockupData(
	userUsercase userUsecase.UserUsecase,
	organizerUsecase organizerUsecase.OrganizerUsecase,
	authUsecase authUsecase.AuthUsecase,
	normalUserUsecase normalUserUsecase.NormalUserUsecase,
	teamUsecase teamUsecase.TeamUsecase,
	addMemberUsecase addMemberUsecase.AddMemberUsecase,
	competitionUsecase competitionUsecase.CompetitionUsecase,
	matchUsecase matchUsecase.MatchUsecase,
) {
	// Create NormalUser
	nOfNormalUser := 40
	for i := 1; i <= nOfNormalUser; i++ {
		email := fmt.Sprintf("normal%d@gmail.com", i)
		password := "1234"
		username := fmt.Sprintf("normal%d", i)
		if err := authUsecase.RegisterNormaluser(
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
		// if err := u.UpdateNormalUser(&model.UpdateNormalUser{
		// 	FirstNameThai: "FirstNameThai" + fmt.Sprintf("%d", i),
		// 	LastNameThai:  "LastNameThai" + fmt.Sprintf("%d", i),
		// 	FirstNameEng:  "FirstNameEng" + fmt.Sprintf("%d", i),
		// 	LastNameEng:   "LastNameEng" + fmt.Sprintf("%d", i),
		// 	Born:          time.Now(),
		// 	Phone:         "00000000000" + fmt.Sprintf("%d", i),
		// 	Height:        175,
		// 	Weight:        70,
		// 	Sex:           "Male",
		// 	Position:      "Foward",
		// 	Nationality:   "Thailand",
		// 	Description:   "Description" + fmt.Sprintf("%d", i),
		// }, uint(i)); err != nil {
		// 	fmt.Printf("err: %v\n", err)
		// 	panic(err.Error())
		// }
		if err := normalUserUsecase.UpdateNormalUser(&model.UpdateNormalUser{
			FirstNameThai: ThaiFirstNames[i],
			LastNameThai:  ThaiLastNames[i],
			FirstNameEng:  FirstNames[i],
			LastNameEng:   LastNames[i],
			Born:          time.Now().AddDate(-19, 0, 0),
			Phone:         "00000000000" + fmt.Sprintf("%d", i),
			Height:        175,
			Weight:        70,
			Sex:           "Male",
			Position:      Position[i],
			Nationality:   "Thailand",
			Description:   PlayerDescriptions[i],
		}, uint(i)); err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err.Error())
		}

		err := userUsercase.UpdateImageProfile(uint(i), fmt.Sprintf("./images/profile/%d.jpg", i))
		if err != nil {
			panic(err)
		}

	}

	// Create team
	r := 'a'
	count := 1
	nOfTeam := 8
	for i := 1; i <= nOfNormalUser; i += nOfNormalUser / nOfTeam {
		// teamName := fmt.Sprintf("Team%d", i)
		if err := teamUsecase.CreateTeam(&model.CreateTeam{
			Name:        TeamNames[i],
			OwnerID:     uint(i),
			Description: TeamDescriptions[i],
		}); err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err)
		}
		err := teamUsecase.UpdateTeamImageProfile(uint(count), fmt.Sprintf("images/profile/%s.png", Clubs[count]), uint(i))
		if err != nil {
			panic(err)
		}
		count++
		r++
	}

	// // Add member to each team
	teamID := 1
	ownerID := 1
	for i := 1; i <= nOfNormalUser; i++ {
		fmt.Printf("i: %v\n", i)
		err := addMemberUsecase.SendAddMemberRequest(&model.AddMemberRequest{
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
		addMemberUsecase.AcceptAddMemberRequest(uint(i), uint(i))
	}

	// // Create Organizer
	nOfOrg := 2
	for i := 0; i < nOfOrg; i++ {
		err := authUsecase.RegisterOrganizer(&model.RegisterOrganizer{
			RegisterUser: model.RegisterUser{
				Email:    fmt.Sprintf("organizer%d@gmail.com", i),
				Password: "1234",
			},
			Phone:         fmt.Sprintf("%d", i) + "00000000000",
			OrganizerName: fmt.Sprintf("org%d", i),
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create Compatition
	err := competitionUsecase.CreateCompetition(&model.CreateCompetition{
		Name:                 "Football A CUP",
		Sport:                "Football",
		Type:                 "Round Robin",
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
		Sex:                  util.SexType[2],
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
		ContactType: "facebook",
		Contact:     "facbook URL",
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}

	err = competitionUsecase.CreateCompetition(&model.CreateCompetition{
		Name:                 "Football B",
		Sport:                util.Sport[0],
		Type:                 util.CompetitionType[0],
		Format:               "5vs5",
		Description:          "Description",
		Rule:                 "Rule",
		Prize:                "Prize",
		StartDate:            time.Now().AddDate(0, 0, -4),
		EndDate:              time.Now(),
		ApplicationType:      util.ApplicationType[0],
		ImageBanner:          "",
		AgeOver:              0,
		AgeUnder:             0,
		Sex:                  "Unisex",
		NumberOfTeam:         8,
		NumOfPlayerInTeamMin: 0,
		NumOfPlayerInTeamMax: 0,
		FieldSurface:         "NaturalGrass",
		OrganizerID:          2,
		Address: model.Address{
			HouseNumber: "HouseNumber",
			Village:     "Village    ",
			Subdistrict: "Subdistrict",
			District:    "District   ",
			PostalCode:  "PostalCode ",
			Country:     "Country    ",
		},
		ContactType: "facebook",
		Contact:     "https://www.google.com/webhp?hl=en&sa=X&ved=0ahUKEwjQ3r77naKFAxWYSGcHHXiJCXsQPAgJ",
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}

	err = competitionUsecase.CreateCompetition(&model.CreateCompetition{
		Name:                 "Football B CUP",
		Sport:                util.Sport[1],
		Type:                 util.CompetitionType[1],
		Format:               "5vs5",
		Description:          "Description",
		Rule:                 "Rule",
		Prize:                "Prize",
		StartDate:            time.Now().AddDate(0, 0, -4),
		EndDate:              time.Now(),
		ApplicationType:      util.ApplicationType[0],
		ImageBanner:          "",
		AgeOver:              0,
		AgeUnder:             0,
		Sex:                  util.SexType[2],
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
		ContactType: "facebook",
		Contact:     "https://www.google.com/webhp?hl=en&sa=X&ved=0ahUKEwjQ3r77naKFAxWYSGcHHXiJCXsQPAgJ",
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}

	err = competitionUsecase.CreateCompetition(&model.CreateCompetition{
		Name:                 "Football C CUP",
		Sport:                util.Sport[1],
		Type:                 util.CompetitionType[1],
		Format:               "5vs5",
		Description:          "Description",
		Rule:                 "Rule",
		Prize:                "Prize",
		StartDate:            time.Now().AddDate(0, 0, -4),
		EndDate:              time.Now(),
		ApplicationType:      util.ApplicationType[1],
		ImageBanner:          "",
		AgeOver:              0,
		AgeUnder:             0,
		Sex:                  util.SexType[1],
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
		ContactType: "facebook",
		Contact:     "https://www.google.com/webhp?hl=en&sa=X&ved=0ahUKEwjQ3r77naKFAxWYSGcHHXiJCXsQPAgJ",
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}

	// Open Application
	err = competitionUsecase.OpenApplicationCompetition(1, 1)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}

	err = competitionUsecase.OpenApplicationCompetition(2, 2)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}

	// err = u.AddJoinCode(1, 4)
	// if err != nil {
	// 	panic(err)
	// }

	// join compatition
	// nOfTeam -= 1
	for i := 1; i <= nOfTeam; i++ {
		err := competitionUsecase.JoinCompetition(&model.JoinCompetition{
			CompetitionID: 1,
			TeamID:        uint(i),
		}, uint(i+((i-1)*4)))
		if err != nil {
			fmt.Printf("nOfTeam: %v\n", nOfTeam)
			fmt.Printf("err: %v\n", err)
			panic(err)
		}
	}

	// join compatition
	// nOfTeam -= 1
	for i := 1; i <= nOfTeam; i++ {
		err := competitionUsecase.JoinCompetition(&model.JoinCompetition{
			CompetitionID: 2,
			TeamID:        uint(i),
		}, uint(i+((i-1)*4)))
		if err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err)
		}
	}

	// Start compatition
	// u.StartCompatition(1)

	// UpdateMatch
	// u.UpdateMatch(1, &model.UpdateMatch{
	// 	DateTime:   time.Now().AddDate(0, 0, -1),
	// 	Team1Goals: 2,
	// 	Team2Goals: 0,
	// 	GoalRecords: []model.GoalRecord{
	// 		model.GoalRecord{
	// 			MatchsID:   1,
	// 			TeamID:     4,
	// 			PlayerID:   5,
	// 			TimeScored: 45,
	// 		},
	// 	},
	// 	Result: "Team1Win",
	// })

	// u.UpdateMatch(1, &model.UpdateMatch{
	// 	DateTime:   time.Now().AddDate(0, 0, -1),
	// 	Team1Goals: 2,
	// 	Team2Goals: 0,
	// 	GoalRecords: []model.GoalRecord{
	// 		model.GoalRecord{
	// 			MatchsID:   1,
	// 			TeamID:     4,
	// 			PlayerID:   5,
	// 			TimeScored: 45,
	// 		},
	// 		model.GoalRecord{
	// 			MatchsID:   1,
	// 			TeamID:     4,
	// 			PlayerID:   5,
	// 			TimeScored: 50,
	// 		},
	// 	},
	// 	Result: "Team1Win",
	// })

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
