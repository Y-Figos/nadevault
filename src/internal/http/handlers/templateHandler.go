package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Y-Figos/nadevault/internal/domain"
	"github.com/Y-Figos/nadevault/internal/http/response"
	"github.com/Y-Figos/nadevault/internal/http/web"
	"github.com/Y-Figos/nadevault/internal/http/web/templates"
	"github.com/Y-Figos/nadevault/internal/repository"
	"github.com/go-chi/chi/v5"
)

type TemplateRenderer struct {
	repo     repository.NadeDataRepository
	renderer *web.Renderer
}

func NewTemplateRenderer(repo repository.NadeDataRepository, renderer *web.Renderer) *TemplateRenderer {
	return &TemplateRenderer{
		repo:     repo,
		renderer: renderer}
}

func (tr *TemplateRenderer) RenderMapPage(w http.ResponseWriter, r *http.Request) {
	mapCode := chi.URLParam(r, "mapCode")
	csMap, err := tr.repo.GetMapByCode(r.Context(), mapCode)
	if err != nil {
		response.WriteAppError(w, err)
		return
	}
	tr.renderer.RenderTemplate(w, http.StatusOK, "base", csMap)
}

func (tr *TemplateRenderer) RenderNadeListPartial(w http.ResponseWriter, r *http.Request) {
	mapCode := chi.URLParam(r, "mapCode")
	csMap, err := tr.repo.GetMapByCode(r.Context(), mapCode)
	if err != nil {
		response.WriteAppError(w, err)
		return
	}
	nadeList, err := tr.repo.ListNadesByMapID(r.Context(), csMap.ID, 10, 0)
	if err != nil {
		response.WriteAppError(w, err)
		return
	}
	nades := struct {
		Nades []domain.Nade
	}{Nades: nadeList}
	tr.renderer.RenderTemplate(w, http.StatusOK, "_nade_list", nades)
}

func (tr *TemplateRenderer) AddNadeForm(w http.ResponseWriter, r *http.Request) {
	csMaps, err := tr.repo.GetMapList(r.Context())
	if err != nil {
		response.WriteAppError(w, err)
		return
	}
	// tr.renderer.RenderTemplate(w, http.StatusOK, "base", data)
	templates.AddNadeForm(csMaps).Render(r.Context(), w)
}

func (tr *TemplateRenderer) AddNade(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		response.WriteError(w, http.StatusBadRequest, "invalid_form", "invalid form data", nil)
		return
	}
	mapid, err := strconv.ParseInt(r.FormValue("map_id"), 10, 16)
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		response.WriteError(w, http.StatusBadRequest, "invalid_map_id", "invalid map ID", nil)
		return
	}
	mapData, err := tr.repo.GetMapByID(r.Context(), int16(mapid))
	if err != nil {
		log.Printf("Error fetching map: %v", err)
		response.WriteAppError(w, err)
		return
	}
	// map DTO -> domain
	nade := domain.Nade{
		Name:       r.FormValue("name"),
		Desc:       r.FormValue("desc"),
		MapID:      int16(mapid),
		MapName:    mapData.DisplayName,
		Type:       domain.NadeType(r.FormValue("type")),
		CommonSide: domain.Side(r.FormValue("side")),
		From:       r.FormValue("from"),
		To:         r.FormValue("to"),
		MouseClick: domain.MouseClick(r.FormValue("mouse_click")),
		IsJumping:  r.FormValue("is_jumping") == "on",
		IsRunning:  r.FormValue("is_running") == "on",
		IsWalking:  r.FormValue("is_walking") == "on",
		IsPublic:   true,
		CreatedBy:  r.FormValue("created_by"),
	}

	pubID, err := tr.repo.AddNade(r.Context(), nade)
	if err != nil {
		log.Printf("Error adding nade: %v", err)
		response.WriteAppError(w, err)
		return
	}

	templates.UploadForm(nade, pubID).Render(r.Context(), w)
}

func (tr *TemplateRenderer) RenderHomePage(w http.ResponseWriter, r *http.Request) {
	templates.HomePage().Render(r.Context(), w)
}
