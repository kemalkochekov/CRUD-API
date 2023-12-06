package main

import (
	"CRUD_Go_Backend/internal/config"
	"CRUD_Go_Backend/internal/handlers"
	"CRUD_Go_Backend/internal/repository"
	"CRUD_Go_Backend/internal/repository/connectionDatabase"
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Could not set up environment variable: %s", err)
	}

	queryParamKey := os.Getenv("QUERY_PARAM_KEY")
	port := os.Getenv("PORT")

	dbConfig, err := config.FromEnv()
	if err != nil {
		log.Fatalf("Could not get environment variable: %v", err)
	}

	database, err := connectionDatabase.NewDB(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect Database %s", err)
	}

	defer database.GetPool(ctx).Close()

	err = connectionDatabase.MigrationUp(dbConfig)
	if err != nil {
		log.Fatalf("Failed to %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	go func() {
		<-quit
		log.Printf("Graceful Shut down")
		// Perform graceful shut down
		if err := connectionDatabase.MigrationDownAndCloseSql(dbConfig); err != nil {
			log.Printf("Error during migration down and closing SQL: %v", err)
		}

		// Close the database pool
		database.GetPool(ctx).Close()

		os.Exit(0)
	}()

	studentStorage := repository.NewStudentStorage(database)
	classInfoStorage := repository.NewClassInfoStorage(database)

	router := handlers.NewRouter(&studentStorage, &classInfoStorage, queryParamKey)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
