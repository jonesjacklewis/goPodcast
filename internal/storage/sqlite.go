package storage

import (
	"database/sql"
	"fmt"

	"github.com/jonesjacklewis/goPodcast/internal/fetching"
)

func CreateDatabase(db *sql.DB) error {
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

func AddEpisode(podcastEpisode fetching.Item, podcastId int, db *sql.DB) error {
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

func AddPodcast(podcast fetching.Podcast, db *sql.DB) (int, error) {
	addPodcastQueryString := `
	INSERT OR IGNORE
	INTO podcasts
	(Title, Url, Description, Image)
	VALUES
	(?, ?, ?, ?)
	`

	preparedStatementaddPodcast, err := db.Prepare(addPodcastQueryString)

	if err != nil {
		return -1, fmt.Errorf("error creating prepared statement for inserting podcast %s", err.Error())
	}

	defer preparedStatementaddPodcast.Close()

	res, err := preparedStatementaddPodcast.Exec(podcast.FeedData.Channel.Title, podcast.Url, podcast.FeedData.Channel.Description, podcast.FeedData.Channel.Image.Url)

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

func AddFullPodcast(podcast fetching.Podcast, db *sql.DB) error {
	id, err := AddPodcast(podcast, db)

	if err != nil {
		return err
	}

	for _, episode := range podcast.FeedData.Channel.Items {
		err = AddEpisode(episode, id, db)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetAllPodcasts(db *sql.DB) ([]fetching.PodcastMetaData, error) {
	getAllPodcastsQuery := `
	SELECT p.Id, p.Title, p.Url, p.Description, p.Image, COUNT(e.Id)
	FROM podcasts p
	LEFT JOIN episodes e
	ON p.Id = e.PodcastId
	GROUP BY p.Id
	`

	rows, err := db.Query(getAllPodcastsQuery)

	if err != nil {
		return []fetching.PodcastMetaData{}, fmt.Errorf("error getting all podcasts %s", err.Error())
	}

	defer rows.Close()

	var podcastMetadata []fetching.PodcastMetaData

	for rows.Next() {
		var podcastMetadataItem fetching.PodcastMetaData

		if err = rows.Scan(&podcastMetadataItem.Id, &podcastMetadataItem.Title, &podcastMetadataItem.Url, &podcastMetadataItem.Description, &podcastMetadataItem.Image, &podcastMetadataItem.NumberOfEpisodes); err == nil {
			podcastMetadata = append(podcastMetadata, podcastMetadataItem)
		}

	}

	return podcastMetadata, nil
}
