package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jonesjacklewis/goPodcast/internal/fetching"
	"github.com/jonesjacklewis/goPodcast/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

var DATABASE_NAME = "data.db"
var USE_OLD = false

type application struct {
	db *sql.DB
}

type addPodcastInput struct {
	RssFeed string `json:"rssFeed"`
}

func oldMain() {

	db, err := sql.Open("sqlite3", DATABASE_NAME)

	if err != nil {
		fmt.Printf("error opening DB %s", err.Error())
		return
	}

	defer db.Close()

	storage.CreateDatabase(db)
	rssFeed := "https://feeds.megaphone.fm/lateralcast"

	podcast, err := fetching.FetchPodcast(rssFeed)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Podcast URL: %s\n", podcast.Url)

	rss := podcast.FeedData

	fmt.Printf("Podcast Title: %s\n", rss.Channel.Title)
	fmt.Printf("Podcast Description: %s\n", rss.Channel.Description)
	fmt.Printf("Number of Episodes: %d\n", len(rss.Channel.Items))

	size := len(rss.Channel.Items)

	if size > 10 {
		size = 10
	}

	for i := 0; i < size; i++ {
		fmt.Println("===========")
		item := rss.Channel.Items[i]

		fmt.Printf("Episode Title: %s\n", item.Title)
		fmt.Printf("Episode Link: %s\n", item.Link)
		fmt.Printf("Enclosure URL: %s\n", item.Enclosure.Url)
		fmt.Println("===========")
	}

	err = storage.AddFullPodcast(podcast, db)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprint(w, "hello world")
}

func (app *application) podcastsHandlerGet(w http.ResponseWriter, _ *http.Request) {

	podcastData, err := storage.GetAllPodcasts(app.db)

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

func (app *application) podcastsHandlerPost(w http.ResponseWriter, r *http.Request) {
	var addPodcastInput addPodcastInput

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

	err = storage.AddFullPodcast(podcast, app.db)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"error\":true, \"message\":\"Unable to add podcast for feed %s\"}", addPodcastInput.RssFeed)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"error\":false, \"message\":\"Successfuly added %s\"}", podcast.FeedData.Channel.Title)
}

func (app *application) podcastsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.podcastsHandlerGet(w, r)
		return

	}

	if r.Method == http.MethodPost {
		app.podcastsHandlerPost(w, r)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"error\":true, \"message\":\"Method %s not supported\"}", r.Method)
}

func main() {
	if USE_OLD {
		oldMain()
		return
	}

	db, err := sql.Open("sqlite3", DATABASE_NAME)

	if err != nil {
		fmt.Printf("error opening DB %s", err.Error())
		return
	}

	defer db.Close()

	storage.CreateDatabase(db)

	app := &application{
		db: db,
	}

	r := chi.NewRouter()

	r.Route("/podcasts", func(c chi.Router) {
		c.Get("/", http.HandlerFunc(app.podcastsHandlerGet))
		c.Post("/", http.HandlerFunc(app.podcastsHandlerPost))
	})

	http.ListenAndServe(":8080", r)

}
