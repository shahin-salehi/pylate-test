package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"shahin/webserver/internal/types"

	"golang.org/x/crypto/bcrypt"
)


func (h *Handler) Login(w http.ResponseWriter, r *http.Request){
	// check
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	payload := new(types.LoginRequest)
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil{
		slog.Error("failed to unmarshal login payload", slog.Any("error", err))
		return
	}

	if payload.Email == "" || payload.Password == "" {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	
	// call db 
	user, err := h.db.GetUserByEmail(r.Context(), payload.Email)
	if err != nil{
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password))
	if err != nil{
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// set session
	err = h.session.SetUserID(w, r, user.ID)
	if err != nil{
		slog.Error("failed to store session", slog.Any("error", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// redirect user
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	} 

	err := h.session.Clear(w, r)
	if err != nil {
		slog.Error("failed to clear session", slog.Any("error", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	} 

	// only allow aiops group

	payload := new(types.User)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil{
		slog.Error("failed to unmarshal payload", slog.Any("function", "Signup"), slog.Any("error", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	if payload.Username == "" || payload.Email == "" || payload.PasswordHash == "" {
		http.Error(w, "bad input", http.StatusBadRequest)
		return
	}

	// the name is a little misleading its not always a password hash
	// sometimes its just a password
	hash, err := bcrypt.GenerateFromPassword([]byte(payload.PasswordHash), bcrypt.DefaultCost)
	if err != nil{
		slog.Error("failed to hash password", slog.Any("error", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	user := types.User{
		Username: payload.Username,
		Email: payload.Email,
		PasswordHash: string(hash),
	}
	// pull from context and pass

	_, err = h.db.RegisterUser(r.Context(), user)
	if err != nil{
		slog.Error("db returned error when registering user", slog.Any("error", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// maybe json this and make it fancy
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered"))
}

