package main

import (
	"context"
	"github.com/Uikola/yandexDAEC/calculator/internal/config"
	"github.com/Uikola/yandexDAEC/calculator/internal/db"
	"github.com/Uikola/yandexDAEC/calculator/internal/db/repository"
	"github.com/Uikola/yandexDAEC/calculator/pkg/discovery"
	"github.com/Uikola/yandexDAEC/calculator/pkg/kafka/consumer"
	"github.com/Uikola/yandexDAEC/calculator/pkg/zlog"
	"github.com/forscht/namegen"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustConfig()
	log := zlog.Default(true, "dev", zerolog.InfoLevel)

	database := db.InitDB(cfg, log)

	// генерация имени для сервиса
	name := namegen.New().WithNumberOfWords(1).WithStyle(namegen.Lowercase).Generate()
	err := discovery.RegistryService(name)
	if err != nil {
		panic(err)
	}

	repo := repository.NewRepository(database)
	operations, err := repo.ListOperation(context.Background())

	go consumer.StartConsumer(name, operations)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}
