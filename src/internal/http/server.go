package http

import (
	"net/http"
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
	fs := http.FileServer(http.Dir("./internal/http/web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handlers.HealthHandler)
		r.Get("/nades/{nadeID}", nadeHandler.GetNadeByID)
		r.Route("/maps", func(r chi.Router) {
			r.Get("/{mapCode}/nades", nadeHandler.ListNadesByMapID)
		})
		r.Get("/maps", nadeHandler.GetMapList)
		r.Post("/nades", nadeHandler.AddNade)

	})

	r.Get("/", templateRenderer.RenderHomePage)
	r.Get("/maps/{mapCode}", templateRenderer.RenderMapPage)
	r.Get("/maps/{mapCode}/nades", templateRenderer.RenderNadeListPartial)
	r.Get("/admin/nades", templateRenderer.AddNadeForm)
	r.Post("/admin/nades/new", templateRenderer.AddNade)
	
	return r
}
