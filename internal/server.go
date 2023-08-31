package server

import (
	"avito-trainee/internal/config"
	"avito-trainee/internal/migrations"
	"avito-trainee/internal/repository"
	"avito-trainee/internal/repository/postgresDB"
	"avito-trainee/internal/service"
	"avito-trainee/internal/transport/rest"
	"avito-trainee/internal/transport/rest/helpers"
	"avito-trainee/pkg/logger"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		log.Panicln(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger := logger.NewLogManager()
	helpers := helpers.NewManager(logger)

	ormDB, err := postgresDB.NewPostgresDB(postgresDB.Config{
		Host:        cfg.DB.Host,
		Port:        cfg.DB.Port,
		Username:    cfg.DB.Username,
		Password:    cfg.DB.Password,
		DBName:      cfg.DB.DBName,
		SSLMode:     cfg.DB.SSLMode,
		Environment: cfg.Env,
	})
	if err != nil {
		logger.Error(err)
		return
	}
	defer postgresDB.CloseDB(ormDB)

	migrations.RunMigrations(ormDB)

	repositories := repository.NewRepository(ormDB)

	timeZone, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
		return
	}

	scheduler := cron.New(cron.WithLocation(timeZone))
	defer scheduler.Stop()

	services := service.NewService(&service.Deps{
		Repository:   repositories,
		Logger:       logger,
		CronSheduler: scheduler,
	})

	handlers := rest.NewHandler(services, cfg, helpers)
	httpServer := rest.NewServer(cfg, handlers.InitRoutes(cfg))

	services.CronService.Init()
	services.CronService.Start()

	go func() {
		if err := httpServer.RunHttp(); err != nil {
			log.Fatalf("Error while run HTTP server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-quit

	httpServer.Shutdown(ctx)
	postgresDB.CloseDB(ormDB)

}
