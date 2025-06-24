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
	"shahin/webserver/internal/types"
	"shahin/webserver/internal/web/components"
	"shahin/webserver/internal/web/pages"
	"strings"

	"github.com/google/uuid"
)

const DiskPath = "/files"

type Handler struct{
	KooroshClient *koorosh.Client
	db db.Crud
}

func New(client *koorosh.Client, db db.Crud)(*Handler,error){
	return &Handler{KooroshClient: client, db: db}, nil
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request)  {

	const headerKey = "query"
	const categoryHeader = "category"
	//get request check method etc
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get header query
	query := r.Header.Get(headerKey)
	category := r.Header.Get(categoryHeader)
	slog.Info("header", slog.Any("category", category))
	
	matches, err := h.KooroshClient.Search(query, category)
	if err != nil {
		slog.Error("koorosh search returned error from handler", slog.Any("error", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}


	// return results
	components.Results(matches).Render(r.Context(), w)
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
	
	// get user Files
	files, err := h.db.ReadFiles(r.Context())
	if err != nil{
		slog.Error("db failed to read files")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	
	pages.Upload(files).Render(r.Context(), w)
	
}


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

		dstPath := filepath.Join("uploads", fn)
		
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


		// Build PDF URL (assuming your server hosts /uploads/)
		pdfURL := fmt.Sprintf("http://localhost:8080/uploads/%s", fn)

		// JSON body for FastAPI
		payload := map[string]string{
			"url":      pdfURL,
			"filename": fileHeader.Filename,
			"category": category, // or any category logic you have
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			slog.Error("Failed to marshal JSON", slog.Any("error", err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Send to FastAPI
		resp, err := http.Post("http://localhost:8000/upload-pdf", "application/json", bytes.NewBuffer(jsonData))
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
	
	//dleete
	err = h.db.DeleteFile(r.Context(), payload.ID)
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















