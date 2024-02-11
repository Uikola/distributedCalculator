package main

import (
	"context"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/config"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/db"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/db/repository/postgres"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/handler"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/handler/computing_resource"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/handler/operation"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/handler/task"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/usecase/cresource_usecase"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/usecase/operation_usecase"
	"github.com/Uikola/yandexDAEC/orchestrator/internal/usecase/task_usecase"
	"github.com/Uikola/yandexDAEC/orchestrator/pkg/zlog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	run()
}

func run() {
	cfg := config.MustConfig()

	switch cfg.Env {
	case "dev":
		log.Logger = zlog.Default(true, "dev", zerolog.InfoLevel)
	case "debug":
		log.Logger = zlog.Default(true, "debug", zerolog.DebugLevel)
	case "prod":
		log.Logger = zlog.Default(true, "prod", zerolog.InfoLevel)
	}

	log.Info().Msg("starting application")

	dataBase := db.InitDB(cfg)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	operationRepository := postgres.NewOperationRepository(dataBase)
	taskRepository := postgres.NewTaskRepository(dataBase)
	cResourceRepository := postgres.NewCResourceRepository(dataBase)
	operationUseCase := operation_usecase.New(operationRepository)
	taskUseCase := task_usecase.New(taskRepository, cResourceRepository)
	cResourceUseCase := cresource_usecase.New(cResourceRepository, taskRepository)
	h := handler.New(operation.New(operationUseCase), task.New(taskUseCase), computing_resource.New(cResourceUseCase))

	handler.Router(h, router)
	log.Info().Msg("starting server")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("failed to start server")
		}
	}()

	log.Info().Msg("server started")

	<-done

	log.Info().Msg("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to stop server")
		return
	}

	log.Info().Msg("cleaning up computing resources")
	err := cResourceRepository.CleanUpCResources(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to clean up computing resources")
	}
	log.Info().Msg("done")

	defer dataBase.Close()
	log.Info().Msg("server stopped")
}
