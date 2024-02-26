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
	s.app.Use(gin.Recovery())
	// Add CORS middleware
	s.app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	// s.app.GET("/user/phone/:phone", userHttpHandler.GetUserByPhone)
	// s.app.PUT("/user/phone/:userID/:phone", userHttpHandler.UpdateNormalUserPhone)

	//Routers
	authRouter := s.app.Group("/auth")
	{
		authRouter.POST("/register", userHttpHandler.RegisterNormaluser)
		authRouter.POST("/register/organizer", userHttpHandler.RegisterOrganizer)
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
		normalRouter.PUT("/acceptAddMemberRequest", userHttpHandler.AcceptAddMemberRequest)
		normalRouter.PUT("/ignoreAddMemberRequest", userHttpHandler.IgnoreAddMemberRequest)
		normalRouter.PUT("/normalUser", userHttpHandler.UpdateNormalUser)
		normalRouter.POST("/team", userHttpHandler.CreateTeam)
		normalRouter.POST("/sendAddMemberRequest", userHttpHandler.SendAddMemberRequest)

		normalRouter.GET("/:id", userHttpHandler.GetUser)
		// normalRouter.POST("/team")
	}

	// publicRouter := s.app.Group("/api/view")

	protectedRoutes := s.app.Group("/api")
	protectedRoutes.Use(auth.Auth())
	{

	}
}
