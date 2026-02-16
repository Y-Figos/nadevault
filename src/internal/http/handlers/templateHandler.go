package handlers

import (
	"net/http"

	"github.com/Y-Figos/nadevault/internal/domain"
	"github.com/Y-Figos/nadevault/internal/http/response"
	"github.com/Y-Figos/nadevault/internal/http/web"
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
