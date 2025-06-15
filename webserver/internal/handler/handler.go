package handler

import (
	"log/slog"
	"net/http"
	"shahin/webserver/internal/grpc"
	"shahin/webserver/internal/types"
	"shahin/webserver/internal/web/components"
	"shahin/webserver/internal/web/pages"
)


type Handler struct{
	KooroshClient *koorosh.Client
}

func New(client *koorosh.Client)(*Handler,error){
	return &Handler{KooroshClient: client}, nil
}

// we need to return html here 
func (h *Handler) Search(w http.ResponseWriter, r *http.Request)  {

	const headerKey = "query"
	//get request check method etc
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get header query
	query := r.Header.Get(headerKey)
	
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

