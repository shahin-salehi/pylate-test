package handler

import (
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
	
	matches, err := h.KooroshClient.Search(query)
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

		// insert to db
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("document(s) uploaded"))
	}




















	
}

