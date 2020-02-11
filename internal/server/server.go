package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/sschiz/apartament/internal/apartment"
	apthttp "github.com/sschiz/apartament/internal/apartment/delivery/http"
	"github.com/sschiz/apartament/internal/apartment/repository/postgres"
	"github.com/sschiz/apartament/internal/apartment/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer *http.Server
	aptUC      apartment.UseCase
}

func NewApp() *App {
	db := initDB()
	ar := postgres.NewApartmentRepository(db)

	return &App{
		aptUC: usecase.NewApartmentUseCase(ar),
	}
}

func (a *App) Run(port string) error {
	router := gin.Default()
	router.Use(
		gin.Logger(),
		//gin.Recovery(),
	)

	api := router.Group("/api")

	apthttp.RegisterHTTPEndpoints(api, a.aptUC)

	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB() *sqlx.DB {
	return sqlx.MustConnect("postgres", viper.GetString("database.DSN"))
}
