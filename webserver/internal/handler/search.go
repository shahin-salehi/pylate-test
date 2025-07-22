package handler

import (
	"log/slog"
	"net/http"
	"shahin/webserver/internal/web/components"
	"shahin/webserver/internal/types"
)


func (h *Handler) Search(w http.ResponseWriter, r *http.Request)  {

	const headerKey = "query"
	const categoryHeader = "category"
	//get request check method etc
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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


	// get header query
	query := r.Header.Get(headerKey)
	category := r.Header.Get(categoryHeader)
	slog.Info("header", slog.Any("category", category))
	
	matches, err := h.KooroshClient.Search(gid, query, category)
	if err != nil {
		slog.Error("koorosh search returned error from handler", slog.Any("error", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}


	// return results
	components.Results(matches).Render(r.Context(), w)
}
