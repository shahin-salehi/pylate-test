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
	pool, err := db.NewDatabase(connectionString)
	
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


	defer pool.Conn.Close()

	//init handler 
	handler, err := handler.New(client, pool)
	if err != nil{
		slog.Error("init handler returned error", slog.Any("error", err))
		os.Exit(1)
	}

	mux := http.NewServeMux()
	// serve static files
	fs := http.FileServer(http.Dir("../../static"))
	fsPublic := http.FileServer(http.Dir("../../public"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.Handle("/public/", http.StripPrefix("/public/", fsPublic))
	//pages
	mux.Handle("/", templ.Handler(pages.Index()))
	mux.Handle("/login", templ.Handler(pages.Login()))

	// handlers
	mux.HandleFunc("/query", handler.Search)
	mux.HandleFunc("/view", handler.View)
	mux.HandleFunc("/data", handler.Files)
	mux.HandleFunc("/upload-pdf", handler.UploadPDF)
	

	http.ListenAndServe(":8080", mux)

}
