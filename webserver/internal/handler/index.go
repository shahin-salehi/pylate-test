package handler

import (
	"log/slog"
	"net/http"
	"shahin/webserver/internal/types"
	"shahin/webserver/internal/web/pages"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request){

	groupID := r.Context().Value(types.ContextGroup)

	gid, ok := groupID.(int64)
	if !ok {
		slog.Error("big error groupid type assert", slog.Any("function", "upload-pdf"))
		slog.Info("here it is", slog.Any("groupID", groupID))
	}
	
	// get index categories
	categories, err := h.db.GetCategories(r.Context(), gid)
	if err != nil {
		slog.Error("categories init returned error, shutting down", slog.Any("error", err))
	}
	

	pages.Index(categories).Render(r.Context(), w)

}
