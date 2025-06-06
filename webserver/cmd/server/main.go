package main

import (
	"log/slog"
	"net/http"
	"os"
	"shahin/webserver/internal/db"
)

func main(){
	slog.Info("starting webserver...")
	
	// set this in deployment 
	connectionString := os.Getenv("DATABASE_URL")
	_, err := db.NewDatabase(connectionString)
	
	if err != nil {
		slog.Error("db init returned error in main, shutting down.", slog.Any("error", err))
		os.Exit(1)
	}

	http.ListenAndServe(":8080", nil)

}
