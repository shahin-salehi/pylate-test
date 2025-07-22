package session

import (
	"log/slog"
	"net/http"
	"github.com/gorilla/sessions"
	"shahin/webserver/internal/db"
)

type Session struct{
	store *sessions.CookieStore // yum yum yum
	db db.Crud
}

const sessionName = "shahin-session"

func InitStore(db db.Crud) (*Session, error){
	secret := "123"  // os.Getenv("SESSION_KEY)
	// key, err := base64.StdEncoding.DecodeString(secret)
	// handle err..

	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path: "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure: false, // change to false if testing without HTTPS
	}

	return &Session{store: store, db: db}, nil

}

func (s *Session) SetUserID(w http.ResponseWriter, r *http.Request, ID int64) error {
	// get the session from the request 
	session, err := s.store.Get(r, sessionName)
	if err != nil{
		slog.Error("session exists but could not be decoded", slog.Any("function", "SetUserID"), slog.Any("error", err))
	}

	groupID, err := s.db.GetUserGroup(r.Context(), ID)
	if err != nil{
		slog.Error("couldnt get user group in session handler", slog.Any("function", "SetUserID"), slog.Any("error", err))
	}

	// set val
	session.Values["group_id"] = groupID 
	session.Values["user_id"] = ID

	return sessions.Save(r, w)
}

func (s *Session) GetUserID(r *http.Request) (int64, bool){
	session, err := s.store.Get(r, sessionName)
	if err != nil{
		slog.Error("session exists but could not be decoded", slog.Any("function", "SetUserID"), slog.Any("error", err))
	}
	
	id, ok := session.Values["user_id"].(int64)
	return id,ok
}

func (s *Session) Clear(w http.ResponseWriter, r *http.Request) error {
	session, err := s.store.Get(r, sessionName)
	if err != nil{
		slog.Error("session exists but could not be decoded", slog.Any("function", "SetUserID"), slog.Any("error", err))
	}

	// read on stackoverflow a while back
	session.Options.MaxAge = -1

	return sessions.Save(r, w)
}
	
