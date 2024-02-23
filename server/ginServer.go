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

	userHttpHandler := handlers.NewUserHttpHandler(userUsercase)

	s.app.Use(gin.Logger())

	// s.app.GET("/user/phone/:phone", userHttpHandler.GetUserByPhone)
	s.app.PUT("/user/phone/:userID/:phone", userHttpHandler.UpdateNormalUserPhone)

	//Routers
	authRouter := s.app.Group("/auth")
	{
		authRouter.POST("/register", userHttpHandler.RegisterUser)
		authRouter.POST("/login", userHttpHandler.LoginUser)
	}

	adminRouter := s.app.Group("/admin")
	adminRouter.Use(auth.AuthAdmin())
	{

	}

	organizerRouter := s.app.Group("/organizer")
	organizerRouter.Use(auth.AuthOrganizer())
	{
		organizerRouter.GET("/users", userHttpHandler.GetUsers)

	}

	normalRouter := s.app.Group("/user")
	normalRouter.Use(auth.AuthNormalUser())
	{
		normalRouter.PUT("/normalUser/:id", userHttpHandler.UpdateNormalUser)
		normalRouter.GET("/:id", userHttpHandler.GetUser)
		normalRouter.POST("/team/:id", userHttpHandler.CreateTeam)
		// normalRouter.POST("/team")

	}

	// publicRouter := s.app.Group("/api/view")

	protectedRoutes := s.app.Group("/api")
	protectedRoutes.Use(auth.Auth())
	{

	}
}
