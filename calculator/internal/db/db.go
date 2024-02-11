package db

import (
	"database/sql"

	"github.com/Uikola/yandexDAEC/calculator/internal/config"
	"github.com/rs/zerolog/log"
)

// InitDB инициализирует инстанцию базы данных.
func InitDB(cfg *config.Config) *sql.DB {
	db, err := sql.Open(cfg.DriverName, cfg.ConnString)
	if err != nil {
		log.Info().Err(err).Msg("failed to connect to the database")
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Info().Err(err).Msg("failed to ping the database")
		return nil
	}
	return db
}
