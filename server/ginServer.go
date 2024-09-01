package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kickoff-league.com/authentication"
	"kickoff-league.com/config"
	"kickoff-league.com/handlers"
	"kickoff-league.com/repositories"
	"kickoff-league.com/usecases/addMemberUsecase"
	"kickoff-league.com/usecases/authUsecase"
	"kickoff-league.com/usecases/competitionUsecase"
	"kickoff-league.com/usecases/matchUsecase"
	"kickoff-league.com/usecases/normalUserUsecase"
	"kickoff-league.com/usecases/organizerUsecase"
	"kickoff-league.com/usecases/teamUsecase"
	"kickoff-league.com/usecases/userUsecase"
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
	s.InitialzieHttpHandler()

	// serverUrl := fmt.Sprintf(":%d", s.cfg.App.Port)
	// s.app.Run(serverUrl)

	serverUrl := fmt.Sprintf(":%d", s.cfg.App.Port)
	srv := &http.Server{
		Addr:    serverUrl,
		Handler: s.app,
	}

	go func() {
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal, 1)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		// catching ctx.Done(). timeout of 5 seconds.
		<-ctx.Done()
		log.Println("timeout of 5 seconds.")
	}()

	// service connections
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	log.Println("Server exiting")

}

func (s *ginServer) InitialzieHttpHandler() {

	// Initialize all layers
	repository := repositories.NewUserPostgresRepository(s.db)

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

	auth := authentication.NewJwtAuthentication(s.cfg.JwtSecretKey)

	httpHandler := handlers.NewhttpHandler(
		userUsecase,
		organizerUsecase,
		authUsecase,
		normalUserUsecase,
		teamUsecase,
		addMemberUsecase,
		competitionUsecase,
		matchUsecase,
	)

	s.app.Use(gin.Logger())
	s.app.Use(gin.Recovery())
	s.app.Static("/images", "./images")

	// Define your allowed origins
	allowedOrigins := []string{
		"http://localhost:5173",
	}

	// Configure CORS middleware with multiple allowed origins
	config := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	// Apply the CORS middleware to the Gin router
	s.app.Use(cors.New(config))

	//Routers
	authRouter := s.app.Group("/auth")
	{
		authRouter.POST("/register/organizer", httpHandler.RegisterOrganizer)
		authRouter.POST("/login", httpHandler.LoginUser)
		authRouter.POST("/register/normal", httpHandler.RegisterNormaluser)
		authRouter.POST("/logout", auth.Auth(), httpHandler.LogoutUser)
		// ========================================================================
		// ========================================================================
		// ========================================================================

	}

	viewRouter := s.app.Group("/view")
	{
		viewRouter.GET("/competition/:id", httpHandler.GetCompetition)
		viewRouter.GET("/competition", httpHandler.GetCompetitions)
		viewRouter.GET("/match/:matchID", httpHandler.GetMatch)
		viewRouter.GET("/normalUsers", httpHandler.GetNormalUsers)
		viewRouter.GET("/normalUsers/:id", httpHandler.GetNormalUser)
		viewRouter.GET("/organizer", httpHandler.GetOrganizers)
		viewRouter.GET("/organizer/:organizerID", httpHandler.GetOrganizer)
		viewRouter.GET("/teams", httpHandler.GetTeams)
		viewRouter.GET("/teams/:id", httpHandler.GetTeam)
		viewRouter.GET("/users/:id", httpHandler.GetUser)
		viewRouter.GET("/users", httpHandler.GetUsers)

		// ========================================================================
		// ========================================================================
		// ========================================================================

		// viewRouter.GET("/compatition/normalUser/:id", userHttpHandler.)

	}

	imageRouter := s.app.Group("/image")
	imageRouter.Use(auth.Auth())
	{
		imageRouter.DELETE("/banner/:compatitionID", auth.AuthOrganizer(), httpHandler.DeleteImageBanner)
		imageRouter.PATCH("/banner/:compatitionID", auth.AuthOrganizer(), httpHandler.UpdateImageBanner)
		imageRouter.DELETE("/team/profile/:teamID", httpHandler.DeleteTeamImageProfile)
		imageRouter.DELETE("/team/cover/:teamID", httpHandler.DeleteTeamImageCover)
		imageRouter.PATCH("/team/cover/:teamID", httpHandler.UpdateTeamImageCover)
		imageRouter.PATCH("/team/profile/:teamID", httpHandler.UpdateTeamImageProfile)
		imageRouter.PATCH("/profile", httpHandler.UpdateImageProfile)
		imageRouter.PATCH("/cover", httpHandler.UpdateImageCover)
		imageRouter.DELETE("/profile", httpHandler.DeleteImageProfile)
		imageRouter.DELETE("/cover", httpHandler.DeleteImageCover)

		// ========================================================================
		// ========================================================================
		// ========================================================================

	}

	organizerRouter := s.app.Group("/organizer")
	organizerRouter.Use(auth.AuthOrganizer())
	{
		organizerRouter.POST("/competition", httpHandler.CreateCompetition)
		organizerRouter.PATCH("/competition/finish/:id", httpHandler.FinishCompetition)
		organizerRouter.PATCH("/competition/cancel/:id", httpHandler.CancelCompetition)
		organizerRouter.PATCH("/competition/open/:id", httpHandler.OpenApplicationCompetition)
		organizerRouter.PATCH("/competition/start/:id", httpHandler.StartCompetition)
		organizerRouter.PATCH("/competition/:id", httpHandler.UpdateCompetition)
		organizerRouter.PUT("/competition/joinCode/add/:compatitionID", httpHandler.AddJoinCode)
		organizerRouter.PUT("/match/:id", httpHandler.UpdateMatch)
		organizerRouter.PUT("", httpHandler.UpdateOrganizer)
		organizerRouter.DELETE("/competition/:competitionID", httpHandler.RemoveCompatitionTeam)

		// ========================================================================
		// ========================================================================
		// ========================================================================
	}

	normalRouter := s.app.Group("/user")
	normalRouter.Use(auth.AuthNormalUser())
	{

		normalRouter.POST("/competition/join", httpHandler.JoinCompetition)
		normalRouter.GET("/nextMatch", httpHandler.GetNextMatch)
		normalRouter.PATCH("/normalUser", httpHandler.UpdateNormalUser)
		normalRouter.GET("/teams", httpHandler.GetTeamByOwnerID)
		normalRouter.POST("/team", httpHandler.CreateTeam)
		normalRouter.DELETE("/team/:teamID", httpHandler.RemoveTeamMember)
		normalRouter.POST("/sendAddMemberRequest", httpHandler.SendAddMemberRequest)
		normalRouter.PATCH("/acceptAddMemberRequest", httpHandler.AcceptAddMemberRequest)
		normalRouter.PATCH("/ignoreAddMemberRequest", httpHandler.IgnoreAddMemberRequest)
		normalRouter.GET("/requests", httpHandler.GetMyPenddingAddMemberRequest)

		// ========================================================================
		// ========================================================================
		// ========================================================================

		// normalRouter.GET("/addMemberRequest", userHttpHandler.SendAddMemberRequest)
		// normalRouter.GET("/team", userHttpHandler.SendAddMemberRequest)
		// normalRouter.POST("/team")
	}

	// publicRouter := s.app.Group("/api/view")

}
