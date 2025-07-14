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

type PodcastsReponse struct {
	Error bool                       `json:"error"`
	Data  []fetching.PodcastMetaData `json:"data"`
}

func (app *Application) podcastsHandlerGet(w http.ResponseWriter, _ *http.Request) {

	podcastData, err := storage.GetAllPodcasts(app.Db)

	if err != nil {
		app.writeError(w, http.StatusInternalServerError, "Unable to get all podcasts")
		return
	}

	podcastResponse := PodcastsReponse{
		Error: false,
		Data:  podcastData,
	}

	err = app.writeJson(w, http.StatusOK, podcastResponse)

	if err != nil {
		app.writeError(w, http.StatusInternalServerError, "Unable to convert podcast data to JSON")
		return
	}
}

func (app *Application) podcastsHandlerPost(w http.ResponseWriter, r *http.Request) {
	var addPodcastInput AddPodcastInput

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&addPodcastInput)

	if err != nil {
		app.writeError(w, http.StatusBadRequest, "Unable to interpret JSON structure")
		return
	}

	podcast, err := fetching.FetchPodcast(addPodcastInput.RssFeed)

	if err != nil {
		app.writeError(w, http.StatusBadRequest, fmt.Sprintf("Unable to fetch podcast for feed %s", addPodcastInput.RssFeed))
		return
	}

	err = storage.AddFullPodcast(podcast, app.Db)

	if err != nil {
		app.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to add podcast for feed %s", addPodcastInput.RssFeed))
		return
	}

	app.writeJson(w, http.StatusOK, app.createSuccessMessage(fmt.Sprintf("Successfully added podcast %s", podcast.FeedData.Channel.Title)))
}
