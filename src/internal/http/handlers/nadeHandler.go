package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Y-Figos/nadevault/internal/domain"
	"github.com/Y-Figos/nadevault/internal/http/response"
	"github.com/Y-Figos/nadevault/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
		response.WriteError(w, http.StatusBadRequest, "invalid_param", "invalid request parameter", nil)
		return
	}
	nade, err := h.repo.GetNadeByID(r.Context(), nadeID)
	if err != nil {
		response.WriteAppError(w, err)
		return
	}
	response.WriteJSON(w, http.StatusOK, nade)
}

func (h *NadeHandler) GetNadeByPublicID(w http.ResponseWriter, r *http.Request) {
	nadeID := chi.URLParam(r, "nadeID")
	nadeUUID, err := uuid.Parse(nadeID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid_uuid", "invalid UUID format", nil)
		return
	}
	nade, err := h.repo.GetNadeByPublicID(r.Context(), nadeUUID)
	if err != nil {
		response.WriteAppError(w, err)
		return
	}
	response.WriteJSON(w, http.StatusOK, nade)
}

func (h *NadeHandler) ListNadesByMapID(w http.ResponseWriter, r *http.Request) {
	mapCode := chi.URLParam(r, "mapCode")
	mapData, err := h.repo.GetMapByCode(r.Context(), mapCode)
	if err != nil {
		response.WriteAppError(w, err)
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
		response.WriteAppError(w, err)
		return
	}
	response.WriteJSON(w, http.StatusOK, nades)
}

func clamplimit(limit int64, def int64, max int64) int64 {
	if limit <= 0 {
		limit = def
	} else if limit > max {
		limit = max
	}
	return limit
}

func (h *NadeHandler) AddNade(w http.ResponseWriter, r *http.Request) {
	var req CreateNadeRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid_json", "invalid JSON body", nil)
		return
	}

	if err := validateCreateNade(&req); err != nil {
		log.Printf("validation error: %v", err)
		response.WriteAppError(w, err)
		return
	}

	// map DTO -> domain
	nade := domain.Nade{
		Info: domain.Info{
			Name:       req.Name,
			Description: req.Desc,
			MapID:      req.MapID,
			Type:       domain.NadeType(req.Type),
			CommonSide: domain.Side(req.CommonSide),
			From:       req.From,
			To:         req.To,
			InputModifiers: domain.InputModifiers{
				MouseClick: domain.MouseClick(req.MouseClick),
				IsJumping:  req.IsJumping,
				IsRunning:  req.IsRunning,
				IsWalking:  req.IsWalking,
			},
		},
		Metadata: domain.Metadata{
			IsPublic:  true,
			CreatedBy: req.CreatedBy,
		},
	}

	id, err := h.repo.AddNade(r.Context(), nade)
	if err != nil {
		response.WriteAppError(w, err)
		return
	}
	log.Printf("nade added with public ID: %s", id)
	response.WriteJSON(w, http.StatusCreated, map[string]string{"public_id": id})
}

func (h *NadeHandler )GetMapList(w http.ResponseWriter, r *http.Request) {
	mapList, err := h.repo.GetMapList(r.Context())
	if err != nil {
		response.WriteAppError(w, err)
		return
	}
	response.WriteJSON(w, http.StatusOK, mapList)
}
