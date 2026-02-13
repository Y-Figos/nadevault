package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/Y-Figos/nadevault/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type NadeHandler struct {
	repo repository.NadeDataRepository
}

func NewNadeHandler(repo repository.NadeDataRepository) *NadeHandler {
	return &NadeHandler{repo: repo}
}

func (h *NadeHandler) GetNadeByID(w http.ResponseWriter, r *http.Request) {
	nadeID, err := strconv.ParseInt(chi.URLParam(r, "nadeID"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid nade ID", http.StatusBadRequest)
		return
	}
	nade, err := h.repo.GetNadeByID(r.Context(), nadeID)
	if err != nil {
		if pgx.ErrNoRows == err {
			http.Error(w, "Nade not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nade)
}

func (h *NadeHandler) ListNadesByMapID(w http.ResponseWriter, r *http.Request) {
	mapCode := chi.URLParam(r, "mapCode")
	mapData, err := h.repo.GetMapByCode(r.Context(), mapCode)
	if err != nil {
		if pgx.ErrNoRows == err {
			http.Error(w, "Map not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	if err != nil {
		limit = 20
	}
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 32)
	if err != nil || offset < 0 {
		offset = 0
	}

	limit = clamplimit(limit, 20, 100)

	nades, err := h.repo.ListNadesByMapID(r.Context(), mapData.ID, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nades)
}

func clamplimit(limit int64, def int64, max int64) int64 {
	if limit <= 0 {
		limit = def
	} else if limit > max {
		limit = max
	}
	return limit
}
