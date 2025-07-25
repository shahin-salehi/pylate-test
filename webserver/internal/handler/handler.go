package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"shahin/webserver/internal/db"
	"shahin/webserver/internal/grpc"
	"shahin/webserver/internal/session"
	"shahin/webserver/internal/types"
	"shahin/webserver/internal/web/pages"
	"strings"

	"github.com/google/uuid"
)
const VolumePath = "/files"

type Handler struct{
	KooroshClient *koorosh.Client
	db db.Crud
	session *session.Session
}

func New(client *koorosh.Client, db db.Crud, sess *session.Session)(*Handler,error){
	return &Handler{KooroshClient: client, db: db, session: sess}, nil
}


	
func (h *Handler) View(w http.ResponseWriter, r *http.Request){
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	filename := r.FormValue("filename") 
	search := r.FormValue("search")


	instruct := types.ViewerInstructions{
		Path: filename,
		Highlight: search,
	}

	pages.Reader(instruct).Render(r.Context(), w)
	
}


func (h *Handler) Files(w http.ResponseWriter, r *http.Request){
	// get group
	// put this in a shared variable at this point
	groupID := r.Context().Value(types.ContextGroup)
	slog.Info("groupid", slog.Any("value", groupID))

	gid, ok := groupID.(int64)
	if !ok {
		slog.Error("big error groupid type assert")
		slog.Info("here it is", slog.Any("groupID", groupID))
	}
	
	// get user Files
	files, err := h.db.ReadFiles(r.Context(), gid)
	if err != nil{
		slog.Error("db failed to read files")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	
	pages.Upload(files).Render(r.Context(), w)
	
}


/*
  UploadPDF This function is central.
  it dictates where files go in prod/helm/cloud
  this will need to change first to persistant storage
*/
func (h *Handler) UploadPDF(w http.ResponseWriter, r *http.Request){
	err := r.ParseMultipartForm(32 << 20) // 32MB 2^20
	if err != nil{
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	category := r.MultipartForm.Value["tag"][0]
	
	files := r.MultipartForm.File["files"]
	for _, fileHeader := range files {
		
		// check on filename if already exists
		// h.db.fileExists(filename)
		// return bad request if != nil

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "error reading file", http.StatusInternalServerError)
			return
		}
		
		defer file.Close()

		ext := filepath.Ext(fileHeader.Filename)
		if ext == "" {
			ext = ".pdf" // fallback
		}
		
		id  := uuid.New()
		fn := id.String() + strings.ToLower(ext)


		dstPath := filepath.Join("./"+VolumePath+"/", fn)
		
		// write
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Error while saving file", http.StatusInternalServerError)
			return
		}
		
		defer dst.Close()
		

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "error writing file to disk", http.StatusInternalServerError)
			return
		}

		/*

		  WARNING THIS WILL CRASH OUR ENDPOINT IF OVERWHELMED

		*/
		groupID := r.Context().Value(types.ContextGroup)

		gid, ok := groupID.(int64)
		if !ok {
			slog.Error("big error groupid type assert", slog.Any("function", "upload-pdf"))
			slog.Info("here it is", slog.Any("groupID", groupID))
		}

		// Build PDF URL (assuming your server hosts )
		pdfURL := fmt.Sprintf("http://website:8080%s/%s", VolumePath, fn)

		// JSON body for FastAPI
		pdf := types.NewPDF{
			Url: pdfURL,
			Filename: fileHeader.Filename,
			Category: category,
			Owner: gid,
		}

		jsonData, err := json.Marshal(pdf)
		if err != nil {
			slog.Error("Failed to marshal JSON", slog.Any("error", err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Send to FastAPI
		resp, err := http.Post("http://upload:8000/upload-pdf", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			slog.Error("Failed to contact FastAPI", slog.Any("error",err))
			http.Error(w, "error contacting indexing service", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusAccepted {
			body, _ := io.ReadAll(resp.Body)
			msg := fmt.Sprintf("FastAPI returned status %d: %s\n", resp.StatusCode, string(body))
			slog.Error("bad status", slog.Any("error", msg))
			http.Error(w, "indexing service error", http.StatusBadGateway)
			return
		}
	}
	
		
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("document(s) uploaded"))
}




// doesnt actually delete from disk though
func (h *Handler) DeletePDF(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	} 
	
	payload := new(types.DeleteRequest)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil{
		slog.Error("failed to decode payload body", slog.Any("error", err))
		http.Error(w, "failed to decode payload", http.StatusUnprocessableEntity)
		return
	}

	// get group
	groupID := r.Context().Value(types.ContextGroup)
	slog.Info("groupid", slog.Any("value", groupID))

	gid, ok := groupID.(int64)
	if !ok {
		slog.Error("big error groupid type assert")
		slog.Info("here it is", slog.Any("groupID", groupID))
	}
	
	//dleete
	err = h.db.DeleteFile(r.Context(), payload.ID, gid)
	if err != nil {
		slog.Error("db returned error for delete file")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	
	//ok
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"server": "file deleted"`))
}




