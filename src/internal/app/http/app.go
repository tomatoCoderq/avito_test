package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tomatoCoderq/avito_task/src/internal/modules/prs"
	"github.com/tomatoCoderq/avito_task/src/internal/modules/stats"
	"github.com/tomatoCoderq/avito_task/src/internal/modules/teams"
	"github.com/tomatoCoderq/avito_task/src/internal/modules/users"

	"github.com/tomatoCoderq/avito_task/src/internal/storage/sql"
)

type App struct {
	port       int
	httpServer *http.Server
}

func New(
	port int,
	address string,

) *App {
	_ = godotenv.Load(".env")

	DB_NAME := os.Getenv("DB_NAME")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")

	connectionString := "postgres://" + DB_USER + ":" + DB_PASSWORD + "@" + DB_HOST + ":" + DB_PORT + "/" + DB_NAME

	router := gin.Default()

	repo, err := sql.New(connectionString)
	if err != nil {
		panic(err)
	}

	teamsRepo := teams.NewRepo(repo)
	teamsService := teams.RegisterService(teamsRepo)
	teamsController := teams.RegisterController(teamsService)

	router.Handle(http.MethodPost, "/team/add", teamsController.TeamCreate)
	router.Handle(http.MethodGet, "/team/get", teamsController.TeamGetByName)
	router.Handle(http.MethodPost, "/team/addUsers", teamsController.AddUsers)
	router.Handle(http.MethodPost, "/team/deactivateUsers", teamsController.DeactivateUsers)

	usersRepo := users.NewRepo(repo)
	usersService := users.RegisterService(usersRepo)
	usersController := users.RegisterController(usersService)

	router.Handle(http.MethodPost, "/users/setIsActive", usersController.SetIsActive)
	router.Handle(http.MethodGet, "/users/getReview", usersController.GetReview)

	prsRepo := prs.NewRepo(repo)
	prsService := prs.RegisterService(prsRepo)
	prsController := prs.RegisterController(prsService)

	router.Handle(http.MethodPost, "/pullRequest/create", prsController.Create)
	router.Handle(http.MethodGet, "/pullRequest/get", prsController.GetByID)
	router.Handle(http.MethodPost, "/pullRequest/merge", prsController.Merge)
	router.Handle(http.MethodPost, "/pullRequest/reassign", prsController.Reassign)

	statsRepo := stats.NewRepo(repo)
	statsService := stats.RegisterService(statsRepo)
	statsController := stats.RegisterController(statsService)

	router.Handle(http.MethodGet, "/stats/users", statsController.GetUserStats)
	router.Handle(http.MethodGet, "/stats/overview", statsController.GetOverview)
	router.Handle(http.MethodGet, "/stats/teams", statsController.GetTeamStats)

	httpServer := &http.Server{
		Addr:    address,
		ReadHeaderTimeout: 10 * time.Second,
		Handler: router,
	}

	return &App{
		port:       port,
		httpServer: httpServer,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server error: %w", err)
	}

	return nil
}

func (a *App) Stop() {
	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		fmt.Printf("Stopped")
	}
}
