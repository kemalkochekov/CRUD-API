package config

import (
	"CRUD_Go_Backend/internal/pkg/pkgErrors"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func FromEnv() (DatabaseConfig, error) {
	dbConfig := DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
	var err error

	dbConfig.Port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return DatabaseConfig{}, pkgErrors.ErrParse
	}

	if dbConfig.Host == "" || dbConfig.User == "" || dbConfig.Password == "" || dbConfig.DBName == "" {
		return DatabaseConfig{}, pkgErrors.ErrDbConfigNotFound
	}

	return dbConfig, nil
}
