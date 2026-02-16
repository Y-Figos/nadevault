package http

import (
	"time"

	"github.com/Y-Figos/nadevault/internal/http/handlers"
	"github.com/Y-Figos/nadevault/internal/http/web"
	"github.com/Y-Figos/nadevault/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(repo repository.NadeDataRepository, renderer *web.Renderer) *chi.Mux {
	nadeHandler := handlers.NewNadeHandler(repo)
	templateRenderer := handlers.NewTemplateRenderer(repo, renderer)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(15 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handlers.HealthHandler)
		r.Get("/nades/{nadeID}", nadeHandler.GetNadeByID)
		r.Route("/maps", func(r chi.Router) {
			r.Get("/{mapCode}/nades", nadeHandler.ListNadesByMapID)
		})
		r.Post("/nades", nadeHandler.AddNade)

	})
	r.Get("/maps/{mapCode}", templateRenderer.RenderMapPage)
	r.Get("/maps/{mapCode}/nades", templateRenderer.RenderNadeListPartial)

	return r
}
