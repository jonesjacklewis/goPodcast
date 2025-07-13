package main

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

func fetchRssFeed(url string) (RSS, error) {
	response, err := http.Get(url)

	if err != nil {
		return RSS{}, fmt.Errorf("error when making request: %s", err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return RSS{}, fmt.Errorf("unsuccessful response from feed: %d", response.StatusCode)
	}

	var rss RSS

	decoder := xml.NewDecoder(response.Body)

	if err := decoder.Decode(&rss); err != nil {
		return RSS{}, fmt.Errorf("invalid XML: %s", err.Error())
	}

	return rss, nil
}

func main() {
	rssFeed := "https://feeds.megaphone.fm/lateralcast"

	rss, err := fetchRssFeed(rssFeed)

	if err != nil {
		fmt.Println(err)
		return
	}

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
}
