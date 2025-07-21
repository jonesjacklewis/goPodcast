package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

type PodcastReponse struct {
	Error bool             `json:"error"`
	Data  fetching.Podcast `json:"data"`
}

type EpisodesResponse struct {
	Error bool               `json:"error"`
	Data  []fetching.Episode `json:"data"`
}

type EpisodeResponse struct {
	Error bool             `json:"error"`
	Data  fetching.Episode `json:"data"`
}

func (app *Application) podcastsHandlerGetEpisode(w http.ResponseWriter, r *http.Request) {

	podcastIdStr := chi.URLParam(r, "podcastId")

	podcastId, err := strconv.Atoi(podcastIdStr)

	if err != nil {
		app.writeError(w, http.StatusBadRequest, "Podcast ID must be an integer")
		return
	}

	episodeIdStr := chi.URLParam(r, "episodeId")

	episodeId, err := strconv.Atoi(episodeIdStr)

	if err != nil {
		app.writeError(w, http.StatusBadRequest, "Episode ID must be an integer")
		return
	}

	episode, err := storage.GetEpisodeForPodcast(app.Db, podcastId, episodeId)

	if err == sql.ErrNoRows {
		app.writeError(w, http.StatusNotFound, fmt.Sprintf("No episode found for podcastId %d and episodeId %d", podcastId, episodeId))
		return
	}

	if err != nil {
		app.writeError(w, http.StatusInternalServerError, "Error extracting episode details")
		return
	}

	var episodeResponse EpisodeResponse = EpisodeResponse{
		Error: false,
		Data:  episode,
	}

	app.writeJson(w, 200, episodeResponse)
}

func (app *Application) podcastsHandlerGetEpisodes(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := strconv.Atoi(id); err != nil {
		app.writeError(w, http.StatusBadRequest, "ID must be an integer")
		return
	}

	episodes, err := storage.GetPodcastEpisodes(app.Db, id)

	if err != nil {
		app.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(episodes) == 0 {
		app.writeError(w, http.StatusNotFound, fmt.Sprintf("No episodes for podcast with ID = %s", id))
		return
	}

	var episodesResponse EpisodesResponse = EpisodesResponse{
		Data:  episodes,
		Error: false,
	}

	app.writeJson(w, 200, episodesResponse)
}

func (app *Application) podcastsHandlerGetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := strconv.Atoi(id); err != nil {
		app.writeError(w, http.StatusBadRequest, "ID must be an integer")
		return
	}

	podcastData, err := storage.GetPodcastById(app.Db, id)

	if err != nil {
		app.writeError(w, http.StatusInternalServerError, "Error retrieving podcast by ID")
		return
	}

	podcastResponse := PodcastReponse{
		Error: false,
		Data:  podcastData,
	}

	if podcastData.FeedData.Channel.Title == "" {

		err = app.writeError(w, http.StatusNotFound, "Podcast not found")

		if err != nil {
			app.writeError(w, http.StatusInternalServerError, "Error creating JSON for podcast by ID")
			return
		}

		return
	}

	app.writeJson(w, http.StatusOK, podcastResponse)

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

func (app *Application) episodesHandlerGet(w http.ResponseWriter, _ *http.Request) {
	episodeData, err := storage.GetAllEpisodes(app.Db)

	if err != nil {
		app.writeError(w, http.StatusInternalServerError, "Unable to get all episodes")
		return
	}

	podcastResponse := EpisodesResponse{
		Error: false,
		Data:  episodeData,
	}

	err = app.writeJson(w, http.StatusOK, podcastResponse)

	if err != nil {
		app.writeError(w, http.StatusInternalServerError, "Unable to convert episode data to JSON")
		return
	}
}
