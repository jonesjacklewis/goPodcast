package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) Routes() http.Handler {
	r := chi.NewRouter()

	r.Route("/podcasts", func(c chi.Router) {
		c.Get("/", http.HandlerFunc(app.podcastsHandlerGet))
		c.Post("/", http.HandlerFunc(app.podcastsHandlerPost))
	})

	return r
}
