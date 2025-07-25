package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *Application) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // For development, allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browser
	}))

	r.Route("/podcasts", func(c chi.Router) {
		c.Get("/", http.HandlerFunc(app.podcastsHandlerGet))
		c.Get("/{id}", http.HandlerFunc(app.podcastsHandlerGetById))
		c.Get("/{id}/episodes", http.HandlerFunc(app.podcastsHandlerGetEpisodes))
		c.Get("/{podcastId}/episodes/{episodeId}", http.HandlerFunc(app.podcastsHandlerGetEpisode))
		c.Post("/", http.HandlerFunc(app.podcastsHandlerPost))
	})

	r.Route("/episodes", func(c chi.Router) {
		c.Get("/", http.HandlerFunc(app.episodesHandlerGet))
		c.Get("/{id}", http.HandlerFunc(app.episodesHandlerGetById))
	})
	return r
}
