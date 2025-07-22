package session

import (
	"context"
	"log/slog"
	"net/http"
	"shahin/webserver/internal/types"
)


func (s *Session) ValidateUserID(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.store.Get(r, sessionName)
		if err != nil{
			// log and reroute
			slog.Error("Decoding error when getting sessionID", slog.Any("class", "session"), slog.Any("function", "ValidateUserID"), slog.Any("error", err))

			// clear
			session.Options.MaxAge = -1
			session.Save(r, w)

			// redirect
			http.Redirect(w,r,"/login", http.StatusFound)
			return
		}

		if session.IsNew || session.Values["user_id"] == nil{
			// get out
			slog.Warn("Invalid or missing session")
			http.Redirect(w,r,"/login", http.StatusFound)
			return
		}

		// ok pass to next
		// add user to context since we're nice guys
		ctx := context.WithValue(r.Context(), types.ContextGroup, session.Values["group_id"])
		ctx = context.WithValue(ctx, types.ContextUser, session.Values["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// Protect wrapper for endpoints HandleFunc
func (s *Session) Protect(endpoint http.HandlerFunc) http.Handler{
	return s.ValidateUserID(http.HandlerFunc(endpoint))
}
