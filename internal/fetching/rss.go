package fetching

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

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

type PodcastMetaData struct {
	Id               int
	Title            string
	Url              string
	Description      string
	NumberOfEpisodes int
	Image            string
}

type Episode struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Link         string `json:"link"`
	EnclosureUrl string `json:"enclosureUrl"`
	PodcastName  string `json:"podcastName"`
}

func FetchPodcast(url string) (Podcast, error) {
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
