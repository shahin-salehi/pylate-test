package main

import (
	"log/slog"
	"net/http"
	"os"
	"shahin/webserver/internal/db"
	"shahin/webserver/internal/grpc"
	"shahin/webserver/internal/handler"
	"shahin/webserver/internal/session"
	"shahin/webserver/internal/web/pages"

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

	crud := db.NewCrud(pool)

	// init session store
	sess, err := session.InitStore(crud)
	if err != nil{
		slog.Error("failed to init session store, shutting down.")
		os.Exit(1)
	}

	// init grpc client obj
	client, err := koorosh.New("localhost:50051")
	if err!= nil {
		slog.Error("failed to create grpc client object", slog.Any("error",err))
		os.Exit(1)
	}
	defer client.Close()


	defer pool.Close()

	//init handler 
	handler, err := handler.New(client, crud, sess)
	if err != nil{
		slog.Error("init handler returned error", slog.Any("error", err))
		os.Exit(1)
	}
	mux := http.NewServeMux()

	// serve static files
	fs := http.FileServer(http.Dir("./static"))
	fsPublic := http.FileServer(http.Dir("./public"))
	fsPrivate := http.FileServer(http.Dir("./files"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.Handle("/public/", http.StripPrefix("/public/", fsPublic))
	mux.Handle("/files/", http.StripPrefix("/files/", fsPrivate))

	//pages testers
	mux.Handle("/", sess.Protect(handler.Index))
	mux.Handle("/login", templ.Handler(pages.Login()))

	// admin pages
	mux.Handle("/register", templ.Handler(pages.Register()))

	// handlers
	// TODO CHANGE TO /api/ 
	mux.Handle("/query", sess.Protect(handler.Search))
	mux.Handle("/view", sess.Protect(handler.View))
	mux.Handle("/data", sess.Protect(handler.Files))
	mux.Handle("/upload-pdf", sess.Protect(handler.UploadPDF))
	mux.Handle("/delete", sess.Protect(handler.DeletePDF))
	mux.HandleFunc("/api/login", handler.Login)
	mux.HandleFunc("/api/logout", handler.Logout)
	mux.Handle("/register-user", sess.Protect(handler.Signup))

	slog.Info("Listening...")
	

	err = http.ListenAndServe(":8080", mux)
	if err != nil{
		slog.Error("why you coming fast man, sorry i crashed", slog.Any("error", err))

	}
}
