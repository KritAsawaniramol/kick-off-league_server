package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	authentication "kickoff-league.com/Authentication"
	"kickoff-league.com/config"
	"kickoff-league.com/handlers"
	"kickoff-league.com/repositories"
	"kickoff-league.com/usecases"
)

type ginServer struct {
	app *gin.Engine
	db  *gorm.DB
	cfg *config.Config
}

func NewGinServer(cfg *config.Config, db *gorm.DB) Server {
	return &ginServer{
		app: gin.New(),
		db:  db,
		cfg: cfg,
	}
}

func (s *ginServer) Start() {
	// Initialzie routers here
	s.initialzieUserHttpHandler()

	serverUrl := fmt.Sprintf(":%d", s.cfg.App.Port)
	s.app.Run(serverUrl)
}

func (s *ginServer) initialzieUserHttpHandler() {
	// Initialize all layers

	userPostgresRepository := repositories.NewUserPostgresRepository(s.db)

	userUsercase := usecases.NewUserUsercaseImpl(
		userPostgresRepository,
	)

	auth := authentication.NewJwtAuthentication(s.cfg.JwtSecretKey)

	userHttpHandler := handlers.NewhttpHandler(userUsercase)

	s.app.Use(gin.Logger())
	s.app.Use(gin.Recovery())
	s.app.Static("/images", "./images")

	// Add CORS middleware

	s.app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	s.app.POST("/upload", userHttpHandler.UploadImage)

	// s.app.GET("/user/phone/:phone", userHttpHandler.GetUserByPhone)
	// s.app.PUT("/user/phone/:userID/:phone", userHttpHandler.UpdateNormalUserPhone)

	//Routers
	authRouter := s.app.Group("/auth")
	{
		authRouter.POST("/register/normal", userHttpHandler.RegisterNormaluser)
		authRouter.POST("/register/organizer", userHttpHandler.RegisterOrganizer)
		authRouter.POST("/login", userHttpHandler.LoginUser)
		authRouter.POST("/logout", auth.Auth(), userHttpHandler.LogoutUser)
	}

	viewRouter := s.app.Group("/view")
	{
		viewRouter.GET("/teams", userHttpHandler.GetTeams)
		viewRouter.GET("/teams/:id", userHttpHandler.GetTeam)
		viewRouter.GET("/users", userHttpHandler.GetUsers)
		viewRouter.GET("/users/:id", userHttpHandler.GetUser)

		viewRouter.GET("/normalUsers", userHttpHandler.GetNormalUsers)
		viewRouter.GET("/normalUsers/:id", userHttpHandler.GetNormalUser)

		viewRouter.GET("/compatition", userHttpHandler.GetCompatitions)
		viewRouter.GET("/compatition/:id", userHttpHandler.GetCompatition)
		// viewRouter.GET("/compatition/normalUser/:id", userHttpHandler.)

		viewRouter.GET("/match/:matchID", userHttpHandler.GetMatch)

		viewRouter.GET("/organizer", userHttpHandler.GetOrganizers)
		viewRouter.GET("/organizer/:organizerID", userHttpHandler.GetOrganizer)

	}

	imageRouter := s.app.Group("/image")
	imageRouter.Use(auth.Auth())
	{
		imageRouter.PUT("/profile", userHttpHandler.UpdateImageProfile)
		imageRouter.PUT("/cover", userHttpHandler.UpdateImageCover)
		imageRouter.PUT("/banner/:compatitionID", userHttpHandler.UpdateImageBanner)
		imageRouter.DELETE("/profile", userHttpHandler.DeleteImageProfile)
		imageRouter.DELETE("/cover", userHttpHandler.DeleteImageCover)
		imageRouter.DELETE("/banner/:compatitionID", userHttpHandler.DeleteImageBanner)

		imageRouter.PUT("/team/profile/:teamID", userHttpHandler.UpdateTeamImageProfile)
		imageRouter.PUT("/team/cover/:teamID", userHttpHandler.UpdateTeamImageCover)
		imageRouter.DELETE("/team/profile/:teamID", userHttpHandler.DeleteTeamImageProfile)
		imageRouter.DELETE("/team/cover/:teamID", userHttpHandler.DeleteTeamImageCover)
	}

	organizerRouter := s.app.Group("/organizer")
	organizerRouter.Use(auth.AuthOrganizer())
	{
		organizerRouter.POST("/compatition", userHttpHandler.CreateCompatition)
		organizerRouter.PUT("/compatition/:id", userHttpHandler.UpdateCompatition)
		organizerRouter.PUT("/compatition/start/:id", userHttpHandler.StartCompatition)
		organizerRouter.PUT("/compatition/open/:id", userHttpHandler.OpenCompatition)
		organizerRouter.PUT("/compatition/finish/:id", userHttpHandler.FinishCompatition)
		organizerRouter.PUT("/compatition/cancel/:id", userHttpHandler.CancelCompatition)
		organizerRouter.PUT("/match/:id", userHttpHandler.UpdateMatch)
		organizerRouter.PUT("/compatition/joinCode/add/:compatitionID", userHttpHandler.AddJoinCode)
		organizerRouter.PUT("/organizer/:organizerID", userHttpHandler.UpdateOrganizer)
		organizerRouter.DELETE("/compatition/:compatitionID", userHttpHandler.RemoveCompatitionTeam)
	}

	normalRouter := s.app.Group("/user")
	normalRouter.Use(auth.AuthNormalUser())
	{
		// normalRouter.GET("/addMemberRequest", userHttpHandler.SendAddMemberRequest)
		// normalRouter.GET("/team", userHttpHandler.SendAddMemberRequest)
		normalRouter.GET("/nextMatch/:normalUserID", userHttpHandler.GetNextMatch)
		normalRouter.GET("/MatchResults/:normalUserID", userHttpHandler.GetMatchResult)
		normalRouter.POST("/team", userHttpHandler.CreateTeam)
		normalRouter.GET("/requests", userHttpHandler.GetMyPenddingAddMemberRequest)
		normalRouter.POST("/sendAddMemberRequest", userHttpHandler.SendAddMemberRequest)
		normalRouter.PUT("/acceptAddMemberRequest", userHttpHandler.AcceptAddMemberRequest)
		normalRouter.PUT("/ignoreAddMemberRequest", userHttpHandler.IgnoreAddMemberRequest)
		normalRouter.PUT("/normalUser", userHttpHandler.UpdateNormalUser)
		normalRouter.GET("/teams/:ownerid", userHttpHandler.GetTeamByOwnerID)
		normalRouter.DELETE("/team/:teamID", userHttpHandler.RemoveTeamMember)
		normalRouter.PUT("/compatition/join", userHttpHandler.JoinCompatition)
		// normalRouter.POST("/team")
	}

	// publicRouter := s.app.Group("/api/view")

}
