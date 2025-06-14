package main

import (
	"log/slog"
	"net/http"
	"os"
	"shahin/webserver/internal/db"
	"shahin/webserver/internal/web/pages"
	"shahin/webserver/internal/grpc"
	"shahin/webserver/internal/handler"

	"github.com/a-h/templ"
)

func main(){
	slog.Info("starting webserver...")
	
	// set this in deployment 
	connectionString := "postgres://admin:password@localhost:9876/documents"   //:= os.Getenv("DATABASE_URL")
	_, err := db.NewDatabase(connectionString)
	
	if err != nil {
		slog.Error("db init returned error in main, shutting down.", slog.Any("error", err))
		os.Exit(1)
	}

	// init grpc client obj
	client, err := koorosh.New("localhost:50051")
	if err!= nil {
		slog.Error("failed to create grpc client object", slog.Any("error",err))
		os.Exit(1)
	}
	defer client.Close()


	//init handler 
	handler, err := handler.New(client)
	if err != nil{
		slog.Error("init handler returned error", slog.Any("error", err))
		os.Exit(1)
	}

	mux := http.NewServeMux()
	// serve static files
	fs := http.FileServer(http.Dir("../../static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	//pages
	mux.Handle("/", templ.Handler(pages.Index()))
	mux.Handle("/data", templ.Handler(pages.Upload()))
	mux.Handle("/login", templ.Handler(pages.Login()))

	// handlers
	mux.HandleFunc("/query", handler.Search)

	http.ListenAndServe(":8080", mux)

}
