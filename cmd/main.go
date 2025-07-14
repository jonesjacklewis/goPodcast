package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jonesjacklewis/goPodcast/internal/api"
	"github.com/jonesjacklewis/goPodcast/internal/fetching"
	"github.com/jonesjacklewis/goPodcast/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

var DATABASE_NAME = "data.db"
var USE_OLD = false

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

	app := &api.Application{
		Db: db,
	}

	http.ListenAndServe(":8080", app.Routes())

}
