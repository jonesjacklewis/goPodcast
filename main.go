package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var DATABASE_NAME = "data.db"

type Enclosure struct {
	Url string `xml:"url,attr"`
}

type Item struct {
	Title     string    `xml:"title"`
	Link      string    `xml:"link"`
	Enclosure Enclosure `xml:"enclosure"`
}

type Image struct {
	Url string `xml:"url"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Image       Image  `xml:"image"`
	Items       []Item `xml:"item"`
}

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Podcast struct {
	FeedData RSS
	Url      string
}

func fetchPodcast(url string) (Podcast, error) {
	response, err := http.Get(url)

	if err != nil {
		return Podcast{}, fmt.Errorf("error when making request: %s", err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return Podcast{}, fmt.Errorf("unsuccessful response from feed: %d", response.StatusCode)
	}

	var rss RSS

	decoder := xml.NewDecoder(response.Body)

	if err := decoder.Decode(&rss); err != nil {
		return Podcast{}, fmt.Errorf("invalid XML: %s", err.Error())
	}

	return Podcast{
		FeedData: rss,
		Url:      url,
	}, nil
}

func createDatabase(db *sql.DB) error {
	podcastsTable := `
	CREATE TABLE IF NOT EXISTS podcasts (
	    Id INTEGER PRIMARY KEY,
		Title TEXT NOT NULL,
		Url TEXT NOT NULL UNIQUE,
		Description TEXT NOT NULL,
		Image TEXT NOT NULL
	)
	`

	_, err := db.Exec(podcastsTable)

	if err != nil {
		return fmt.Errorf("error executing create podcast table %s", err.Error())
	}

	episodesTable := `
	CREATE TABLE IF NOT EXISTS episodes (
	    Id INTEGER PRIMARY KEY,
		PodcastId INTEGER NOT NULL,
		EpisodeTitle TEXT NOT NULL,
		EpisodeLink TEXT NOT NULL,
		EnclosureUrl TEXT NOT NULL UNIQUE,
		FOREIGN KEY(PodcastId) REFERENCES podcasts (Id)
	)
	`

	_, err = db.Exec(episodesTable)

	if err != nil {
		return fmt.Errorf("error executing create episodes table %s", err.Error())
	}

	return nil
}

func addEpisode(podcastEpisode Item, podcastId int, db *sql.DB) error {
	addPodcastEpisodeQuery := `
	INSERT OR IGNORE
	INTO episodes
	(PodcastId, EpisodeTitle, EpisodeLink, EnclosureUrl)
	VALUES
	(?, ?, ?, ?)
	`

	preparedStatementAddEpsiode, err := db.Prepare(addPodcastEpisodeQuery)

	if err != nil {
		return fmt.Errorf("error creating prepared statement for inserting podcast episode")
	}

	defer preparedStatementAddEpsiode.Close()

	_, err = preparedStatementAddEpsiode.Exec(podcastId, podcastEpisode.Title, podcastEpisode.Link, podcastEpisode.Enclosure.Url)

	if err != nil {
		return fmt.Errorf("error executing prepared statement for podcast epsiode %s", err.Error())
	}

	return nil
}

func addPodcast(podcast Podcast, db *sql.DB) (int, error) {
	addPodcastQueryString := `
	INSERT OR IGNORE
	INTO podcasts
	(Title, Url, Description, Image)
	VALUES
	(?, ?, ?, ?)
	`

	preparedStatementAddPodcast, err := db.Prepare(addPodcastQueryString)

	if err != nil {
		return -1, fmt.Errorf("error creating prepared statement for inserting podcast %s", err.Error())
	}

	defer preparedStatementAddPodcast.Close()

	res, err := preparedStatementAddPodcast.Exec(podcast.FeedData.Channel.Title, podcast.Url, podcast.FeedData.Channel.Description, podcast.FeedData.Channel.Image.Url)

	if err != nil {
		return -1, fmt.Errorf("error inserting podcast %s", err.Error())
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return -1, fmt.Errorf("error getting rows affected by podcast insert")
	}

	podcastId := -1

	if rowsAffected == 0 {
		getPodcastIdQuery := `
		SELECT Id
		FROM podcasts
		WHERE Url = ?
		`

		preparedStatementGetPodcastId, err := db.Prepare(getPodcastIdQuery)

		if err != nil {
			return -1, fmt.Errorf("error creating prepared statement for getting podcast ID")
		}

		podcastIdRow := preparedStatementGetPodcastId.QueryRow(podcast.Url)

		err = podcastIdRow.Scan(&podcastId)

		if err != nil {
			return -1, fmt.Errorf("unable to extract podcast ID %s", err.Error())
		}
	} else {
		latestId, err := res.LastInsertId()

		if err != nil {
			return -1, fmt.Errorf("unable to use last insert ID for podcast ID %s", err.Error())
		}

		podcastId = int(latestId)
	}

	return podcastId, nil
}

func addFullPodcast(podcast Podcast, db *sql.DB) error {
	id, err := addPodcast(podcast, db)

	if err != nil {
		return err
	}

	for _, episode := range podcast.FeedData.Channel.Items {
		err = addEpisode(episode, id, db)

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {

	db, err := sql.Open("sqlite3", DATABASE_NAME)

	if err != nil {
		fmt.Printf("error opening DB %s", err.Error())
		return
	}

	defer db.Close()

	createDatabase(db)
	rssFeed := "https://feeds.megaphone.fm/lateralcast"

	podcast, err := fetchPodcast(rssFeed)

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

	err = addFullPodcast(podcast, db)

	if err != nil {
		fmt.Println(err)
		return
	}
}
