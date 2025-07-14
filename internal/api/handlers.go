package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jonesjacklewis/goPodcast/internal/fetching"
	"github.com/jonesjacklewis/goPodcast/internal/storage"
)

type Application struct {
	Db *sql.DB
}

type AddPodcastInput struct {
	RssFeed string `json:"rssFeed"`
}

func (app *Application) podcastsHandlerGet(w http.ResponseWriter, _ *http.Request) {

	podcastData, err := storage.GetAllPodcasts(app.Db)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, "{\"error\":true, \"message\":\"Unable to get all podcasts\"}")
		return
	}

	jsonPodcastData, err := json.Marshal(podcastData)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, "{\"error\":true, \"message\":\"Unable to convert to JSON\"}")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonPodcastData))
}

func (app *Application) podcastsHandlerPost(w http.ResponseWriter, r *http.Request) {
	var addPodcastInput AddPodcastInput

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&addPodcastInput)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, "{\"error\":true, \"message\":\"Unable to interpret JSON object\"}")
		return
	}

	podcast, err := fetching.FetchPodcast(addPodcastInput.RssFeed)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"error\":true, \"message\":\"Unable to fetch podcast for feed %s\"}", addPodcastInput.RssFeed)
		return
	}

	err = storage.AddFullPodcast(podcast, app.Db)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"error\":true, \"message\":\"Unable to add podcast for feed %s\"}", addPodcastInput.RssFeed)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"error\":false, \"message\":\"Successfuly added %s\"}", podcast.FeedData.Channel.Title)
}
