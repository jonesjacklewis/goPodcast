package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jonesjacklewis/goPodcast/internal/fetching"
	"github.com/jonesjacklewis/goPodcast/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

var testXml = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:googleplay="http://www.google.com/schemas/play-podcasts/1.0" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel>
    <atom:link href="https://feeds.megaphone.fm/lateralcast" rel="self" type="application/rss+xml"/>
    <title>Lateral with Tom Scott</title>
    <link>https://www.lateralcast.com</link>
    <language>en</language>
    <copyright>Pad 26 Limited / Labyrinth Games Ltd</copyright>
    <description>Award-winning comedy panel game about weird questions with wonderful answers. Each week, Tom Scott is joined by three guests to ask each other questions with a sideways twist. There's no points or prizes - just reputation and bragging rights on the line. Enquiries and question submissions: https://www.lateralcast.com.
Ad-free shows and bonus content: https://www.lateralcast.com/club.</description>
    <image>
      <url>https://megaphone.imgix.net/podcasts/2a3f4416-77e3-11ee-8e56-f7ac65d137ba/image/b722eeab8db682b441e1a0c471c9f92c.jpg?ixlib=rails-4.3.1&amp;max-w=3000&amp;max-h=3000&amp;fit=crop&amp;auto=format,compress</url>
      <title>Lateral with Tom Scott</title>
      <link>https://www.lateralcast.com</link>
    </image>
    <itunes:explicit>no</itunes:explicit>
    <itunes:type>episodic</itunes:type>
    <itunes:subtitle></itunes:subtitle>
    <itunes:author>Tom Scott and David Bodycombe</itunes:author>
    <itunes:summary>Award-winning comedy panel game about weird questions with wonderful answers. Each week, Tom Scott is joined by three guests to ask each other questions with a sideways twist. There's no points or prizes - just reputation and bragging rights on the line. Enquiries and question submissions: https://www.lateralcast.com.
Ad-free shows and bonus content: https://www.lateralcast.com/club.</itunes:summary>
    <content:encoded>
      <![CDATA[<p>Award-winning comedy panel game about weird questions with wonderful answers. Each week, Tom Scott is joined by three guests to ask each other questions with a sideways twist. There's no points or prizes - just reputation and bragging rights on the line. Enquiries and question submissions: https://www.lateralcast.com.</p><p>Ad-free shows and bonus content: https://www.lateralcast.com/club.</p>]]>
    </content:encoded>
    <itunes:owner>
      <itunes:name>Lateral with Tom Scott</itunes:name>
      <itunes:email>lateralwithtomscott@gmail.com</itunes:email>
    </itunes:owner>
    <itunes:image href="https://megaphone.imgix.net/podcasts/2a3f4416-77e3-11ee-8e56-f7ac65d137ba/image/b722eeab8db682b441e1a0c471c9f92c.jpg?ixlib=rails-4.3.1&amp;max-w=3000&amp;max-h=3000&amp;fit=crop&amp;auto=format,compress"/>
    <itunes:category text="Comedy">
    </itunes:category>
    <itunes:category text="Education">
    </itunes:category>
    <itunes:category text="Leisure">
      <itunes:category text="Games"/>
    </itunes:category>
    <item>
      <title>144: Let's visit Greenland!</title>
      <link>https://www.lateralcast.com/episodes/144</link>
      <description>Iszi Lawrence, Dani Siller and Bill Sunderland face questions about miniature marines, technical techniques and prudent procedures.

LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.

HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: Jordan Cook-Irwin, Katherine Q, Nick Huntington-Klein, Alex Rinehart, Hendrik, James Bailey. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.
Learn more about your ad choices. Visit megaphone.fm/adchoices</description>
      <pubDate>Fri, 11 Jul 2025 04:00:00 -0000</pubDate>
      <itunes:title>Let's visit Greenland!</itunes:title>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>144</itunes:episode>
      <itunes:author>Tom Scott and David Bodycombe</itunes:author>
      <itunes:image href="https://megaphone.imgix.net/podcasts/d6b68bdc-5220-11f0-b841-b3273de395c1/image/f19aae6fb1593677a20056286febefe2.jpg?ixlib=rails-4.3.1&amp;max-w=3000&amp;max-h=3000&amp;fit=crop&amp;auto=format,compress"/>
      <itunes:subtitle>A weekly podcast about interesting questions</itunes:subtitle>
      <itunes:summary>Iszi Lawrence, Dani Siller and Bill Sunderland face questions about miniature marines, technical techniques and prudent procedures.

LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.

HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: Jordan Cook-Irwin, Katherine Q, Nick Huntington-Klein, Alex Rinehart, Hendrik, James Bailey. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.
Learn more about your ad choices. Visit megaphone.fm/adchoices</itunes:summary>
      <content:encoded>
        <![CDATA[<p>Iszi Lawrence, Dani Siller and Bill Sunderland face questions about miniature marines, technical techniques and prudent procedures.</p>
<p>LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.</p>
<p>HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: Jordan Cook-Irwin, Katherine Q, Nick Huntington-Klein, Alex Rinehart, Hendrik, James Bailey. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.</p><p> </p><p>Learn more about your ad choices. Visit <a href="https://megaphone.fm/adchoices">megaphone.fm/adchoices</a></p>]]>
      </content:encoded>
      <itunes:duration>3315</itunes:duration>
      <itunes:explicit>no</itunes:explicit>
      <guid isPermaLink="false"><![CDATA[d6b68bdc-5220-11f0-b841-b3273de395c1]]></guid>
      <enclosure url="https://traffic.megaphone.fm/PAD9460356563.mp3?updated=1752189576" length="0" type="audio/mpeg"/>
    </item>
    <item>
      <title>143: Mind-altering coasters</title>
      <link>https://www.lateralcast.com/episodes/143/</link>
      <description>Stuart Laws, Bella Hull and Olaf Falafel face questions about pally pretzels, pesky pets and prairie programmes.

LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.

HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: GC, Enigma, Thomas Irwin, OMacMacca, Otto Forsbom. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.
Learn more about your ad choices. Visit megaphone.fm/adchoices</description>
      <pubDate>Fri, 04 Jul 2025 04:00:00 -0000</pubDate>
      <itunes:title>Mind-altering coasters</itunes:title>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>143</itunes:episode>
      <itunes:author>Tom Scott and David Bodycombe</itunes:author>
      <itunes:image href="https://megaphone.imgix.net/podcasts/d6829408-5220-11f0-b841-43771805bf55/image/1ad1b579272c21308cf34a357e30a35f.jpg?ixlib=rails-4.3.1&amp;max-w=3000&amp;max-h=3000&amp;fit=crop&amp;auto=format,compress"/>
      <itunes:subtitle>A weekly podcast about interesting questions</itunes:subtitle>
      <itunes:summary>Stuart Laws, Bella Hull and Olaf Falafel face questions about pally pretzels, pesky pets and prairie programmes.

LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.

HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: GC, Enigma, Thomas Irwin, OMacMacca, Otto Forsbom. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.
Learn more about your ad choices. Visit megaphone.fm/adchoices</itunes:summary>
      <content:encoded>
        <![CDATA[<p>Stuart Laws, Bella Hull and Olaf Falafel face questions about pally pretzels, pesky pets and prairie programmes.</p>
<p>LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.</p>
<p>HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: GC, Enigma, Thomas Irwin, OMacMacca, Otto Forsbom. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.</p><p> </p><p>Learn more about your ad choices. Visit <a href="https://megaphone.fm/adchoices">megaphone.fm/adchoices</a></p>]]>
      </content:encoded>
      <itunes:duration>3185</itunes:duration>
      <itunes:explicit>no</itunes:explicit>
      <guid isPermaLink="false"><![CDATA[d6829408-5220-11f0-b841-43771805bf55]]></guid>
      <enclosure url="https://traffic.megaphone.fm/PAD1861973793.mp3?updated=1751395508" length="0" type="audio/mpeg"/>
    </item>
    <item>
      <title>142: The inverted tree</title>
      <link>https://www.lateralcast.com/episodes/142/</link>
      <description>Sam Denby, Adam Chase and Ben Doyle face questions about hilarious horrors, renamed ranks and portable plaques.

LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.

HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: Joël, Aaron Solomon, Ivan Walters, Chris Tam, Artie.. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.
Learn more about your ad choices. Visit megaphone.fm/adchoices</description>
      <pubDate>Fri, 27 Jun 2025 04:00:00 -0000</pubDate>
      <itunes:title>The inverted tree</itunes:title>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>142</itunes:episode>
      <itunes:author>Tom Scott and David Bodycombe</itunes:author>
      <itunes:image href="https://megaphone.imgix.net/podcasts/48093370-b429-11ef-887f-6bb8bde5bc0a/image/96557fca60fd7516e49ebefce540a5cb.jpg?ixlib=rails-4.3.1&amp;max-w=3000&amp;max-h=3000&amp;fit=crop&amp;auto=format,compress"/>
      <itunes:subtitle>A weekly podcast about interesting questions</itunes:subtitle>
      <itunes:summary>Sam Denby, Adam Chase and Ben Doyle face questions about hilarious horrors, renamed ranks and portable plaques.

LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.

HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: Joël, Aaron Solomon, Ivan Walters, Chris Tam, Artie.. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.
Learn more about your ad choices. Visit megaphone.fm/adchoices</itunes:summary>
      <content:encoded>
        <![CDATA[<p>Sam Denby, Adam Chase and Ben Doyle face questions about hilarious horrors, renamed ranks and portable plaques.</p>
<p>LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.</p>
<p>HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: Joël, Aaron Solomon, Ivan Walters, Chris Tam, Artie.. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.</p><p> </p><p>Learn more about your ad choices. Visit <a href="https://megaphone.fm/adchoices">megaphone.fm/adchoices</a></p>]]>
      </content:encoded>
      <itunes:duration>2542</itunes:duration>
      <itunes:explicit>no</itunes:explicit>
      <guid isPermaLink="false"><![CDATA[48093370-b429-11ef-887f-6bb8bde5bc0a]]></guid>
      <enclosure url="https://traffic.megaphone.fm/PAD4653262684.mp3?updated=1750898059" length="0" type="audio/mpeg"/>
    </item>
    <item>
      <title>141: The world's longest poem</title>
      <link>https://www.lateralcast.com/episodes/141/</link>
      <description>Matt Gray, Daniel Peake and Charlotte Yeung face questions about slovenly sales, mystery meat and risky rabbits.

LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.

HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: Leo Taanila, Daniel Peake, Stefan Teffy, Mitchell Lapham, Hugo Bouma, Ben Kitchen, Elena and Alessandro.. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.
Learn more about your ad choices. Visit megaphone.fm/adchoices</description>
      <pubDate>Fri, 20 Jun 2025 04:00:00 -0000</pubDate>
      <itunes:title>The world's longest poem</itunes:title>
      <itunes:episodeType>full</itunes:episodeType>
      <itunes:episode>141</itunes:episode>
      <itunes:author>Tom Scott and David Bodycombe</itunes:author>
      <itunes:image href="https://megaphone.imgix.net/podcasts/47dbf856-b429-11ef-887f-979ee3997071/image/1a1c06478b8e38ed0ba0f8e3322d12aa.jpg?ixlib=rails-4.3.1&amp;max-w=3000&amp;max-h=3000&amp;fit=crop&amp;auto=format,compress"/>
      <itunes:subtitle>A weekly podcast about interesting questions</itunes:subtitle>
      <itunes:summary>Matt Gray, Daniel Peake and Charlotte Yeung face questions about slovenly sales, mystery meat and risky rabbits.

LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.

HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: Leo Taanila, Daniel Peake, Stefan Teffy, Mitchell Lapham, Hugo Bouma, Ben Kitchen, Elena and Alessandro.. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.
Learn more about your ad choices. Visit megaphone.fm/adchoices</itunes:summary>
      <content:encoded>
        <![CDATA[<p>Matt Gray, Daniel Peake and Charlotte Yeung face questions about slovenly sales, mystery meat and risky rabbits.</p>
<p>LATERAL is a comedy panel game podcast about weird questions with wonderful answers, hosted by Tom Scott. For business enquiries, contestant appearances or question submissions, visit https://lateralcast.com.</p>
<p>HOST: Tom Scott. QUESTION PRODUCER: David Bodycombe. EDITED BY: Julie Hassett at The Podcast Studios, Dublin. MUSIC: Karl-Ola Kjellholm ('Private Detective'/'Agrumes', courtesy of epidemicsound.com). ADDITIONAL QUESTIONS: Leo Taanila, Daniel Peake, Stefan Teffy, Mitchell Lapham, Hugo Bouma, Ben Kitchen, Elena and Alessandro.. FORMAT: Pad 26 Limited/Labyrinth Games Ltd. EXECUTIVE PRODUCERS: David Bodycombe and Tom Scott. © Pad 26 Limited (https://www.pad26.com) / Labyrinth Games Ltd. 2025.</p><p> </p><p>Learn more about your ad choices. Visit <a href="https://megaphone.fm/adchoices">megaphone.fm/adchoices</a></p>]]>
      </content:encoded>
      <itunes:duration>3027</itunes:duration>
      <itunes:explicit>no</itunes:explicit>
      <guid isPermaLink="false"><![CDATA[47dbf856-b429-11ef-887f-979ee3997071]]></guid>
      <enclosure url="https://traffic.megaphone.fm/PAD7707747107.mp3?updated=1750354842" length="0" type="audio/mpeg"/>
    </item>
  </channel>
</rss>
`

type ErrorMessage struct {
	Error           bool                       `json:"error"`
	Message         string                     `json:"message"`
	PodcastMetaData []fetching.PodcastMetaData `json:"data"`
}

func newTestApplication(t *testing.T) (*Application, *sql.DB) {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Errorf("Expected no error when opening DB, got %s", err.Error())
	}

	err = storage.CreateDatabase(db)

	if err != nil {
		t.Errorf("Expected no error when creating DB, got %s", err.Error())
	}

	return &Application{
		Db: db,
	}, db

}

func TestHandlers(t *testing.T) {

	rssServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(testXml))
	}))

	defer rssServer.Close()

	app, db := newTestApplication(t)

	defer db.Close()

	router := app.Routes()

	apiServer := httptest.NewServer(router)

	defer apiServer.Close()

	postBody := fmt.Sprintf(`{"rssFeed": "%s"}`, rssServer.URL)
	fmt.Println(postBody)
	bodyReader := strings.NewReader(postBody)

	targetUrl := apiServer.URL + "/podcasts"

	req, err := http.NewRequest(http.MethodPost, targetUrl, bodyReader)

	if err != nil {
		t.Fatalf("failed to create POST request: %s", err.Error())
	}

	res, err := apiServer.Client().Do(req)

	if err != nil {
		t.Fatalf("failed to send POST request: %s", err.Error())
	}

	var errorMessage ErrorMessage

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&errorMessage)

	if err != nil {
		t.Fatal("Unable to interpret JSON structure")
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d got %d, %s", http.StatusOK, res.StatusCode, errorMessage.Message)
	}

	req, err = http.NewRequest(http.MethodGet, targetUrl, nil)

	if err != nil {
		t.Fatalf("failed to create GET request %s", err.Error())
	}

	res, err = apiServer.Client().Do(req)

	if err != nil {
		t.Fatalf("failed to make GET request %s", err.Error())
	}

	decoder = json.NewDecoder(res.Body)

	err = decoder.Decode(&errorMessage)

	if err != nil {
		t.Fatal("Unable to interpret JSON structure")
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d got %d, %s", http.StatusOK, res.StatusCode, errorMessage.Message)
	}

	if len(errorMessage.PodcastMetaData) == 0 {
		t.Fatalf("Expected to have podcast data")
	}

	targetUrl = apiServer.URL + "/podcasts/1"

	req, err = http.NewRequest(http.MethodGet, targetUrl, nil)

	if err != nil {
		t.Errorf("Exepceted no error, got %s", err.Error())
	}

	res, err = apiServer.Client().Do(req)

	if err != nil {
		t.Errorf("Expected no err on /podcasts/1 got %s", err.Error())
	}

	var podastResponse PodcastReponse

	decoder = json.NewDecoder(res.Body)

	err = decoder.Decode(&podastResponse)

	if err != nil {
		t.Errorf("Error decoding PodcastResponse")
	}

	if podastResponse.Error {
		t.Errorf("Expected no error on PodcastResponse")
	}

	if podastResponse.Data.FeedData.Channel.Title != "Lateral with Tom Scott" {
		t.Errorf("Unexpected title %s", podastResponse.Data.FeedData.Channel.Title)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d on /podcasts/1 got %d", http.StatusOK, res.StatusCode)
	}

	targetUrl = apiServer.URL + "/podcasts/2"

	req, err = http.NewRequest(http.MethodGet, targetUrl, nil)

	if err != nil {
		t.Errorf("Exepceted no error, got %s", err.Error())
	}

	res, err = apiServer.Client().Do(req)

	if err != nil {
		t.Errorf("Expected no err on /podcasts/2 got %s", err.Error())
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d on /podcasts/2 got %d", http.StatusNotFound, res.StatusCode)
	}

	targetUrl = apiServer.URL + "/podcasts/1/episodes"

	req, err = http.NewRequest(http.MethodGet, targetUrl, nil)

	if err != nil {
		t.Errorf("Expected no error when creating request for /podcasts/1/episodes, got %s", err.Error())
	}

	res, err = apiServer.Client().Do(req)

	if err != nil {
		t.Errorf("Expected no error when mkaing request for /podcasts/1/episodes, got %s", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected code %d for /podcasts/1/episodes got %d", http.StatusOK, res.StatusCode)
	}

	var episodesResponse EpisodesResponse

	decoder = json.NewDecoder(res.Body)

	err = decoder.Decode(&episodesResponse)

	if err != nil {
		t.Errorf("Unable to decode response for /podcasts/1/episodes")
	}

	if episodesResponse.Error {
		t.Errorf("Expected no error on GET /podcasts/1/episodes")
	}

	if len(episodesResponse.Data) == 0 {
		t.Errorf("Expected to have episodes on GET /podcasts/1/episodes")
	}

	targetUrl = apiServer.URL + "/podcasts/2/episodes"

	req, err = http.NewRequest(http.MethodGet, targetUrl, nil)

	if err != nil {
		t.Errorf("Expected no error when creating request for /podcasts/2/episodes, got %s", err.Error())
	}

	res, err = apiServer.Client().Do(req)

	if err != nil {
		t.Errorf("Expected no error when mkaing request for /podcasts/2/episodes, got %s", err.Error())
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected code %d for /podcasts/2/episodes got %d", http.StatusNotFound, res.StatusCode)
	}

	targetUrl = apiServer.URL + "/podcasts/1/episodes/1"

	req, err = http.NewRequest(http.MethodGet, targetUrl, nil)

	if err != nil {
		t.Errorf("Expected no error when creating request for /podcasts/1/episodes/1, got %s", err.Error())
	}

	res, err = apiServer.Client().Do(req)

	if err != nil {
		t.Errorf("Expected no error when mkaing request for /podcasts/1/episodes/1, got %s", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected code %d for /podcasts/1/episodes/1 got %d", http.StatusOK, res.StatusCode)
	}

	var episodeResponse EpisodeResponse

	decoder = json.NewDecoder(res.Body)

	err = decoder.Decode(&episodeResponse)

	if err != nil {
		t.Errorf("Unable to decode response for /podcasts/1/episodes/1")
	}

	if episodesResponse.Error {
		t.Errorf("Expected no error on GET /podcasts/1/episodes/1")
	}

	expectedTitle := "Let's visit Greenland!"

	if episodeResponse.Data.Title != expectedTitle {
		t.Errorf("Exepected %s got %s for /podcasts/1/episodes/1", expectedTitle, episodeResponse.Data.Title)
	}

	targetUrl = apiServer.URL + "/podcasts/2/episodes/1"

	req, err = http.NewRequest(http.MethodGet, targetUrl, nil)

	if err != nil {
		t.Errorf("Expected no error when creating request for /podcasts/2/episodes/1, got %s", err.Error())
	}

	res, err = apiServer.Client().Do(req)

	if err != nil {
		t.Errorf("Expected no error when mkaing request for /podcasts/2/episodes/1, got %s", err.Error())
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected code %d for /podcasts/2/episodes/1 got %d", http.StatusNotFound, res.StatusCode)
	}

	targetUrl = apiServer.URL + "/episodes"

	req, err = http.NewRequest(http.MethodGet, targetUrl, nil)

	if err != nil {
		t.Errorf("Expected no error on creating GET to %s got %s", targetUrl, err.Error())
	}

	res, err = apiServer.Client().Do(req)

	if err != nil {
		t.Errorf("Expected no error on GET to %s got %s", targetUrl, err.Error())
	}

	var episodesResponseEpisodesEndpoint EpisodesResponse

	decoder = json.NewDecoder(res.Body)

	err = decoder.Decode(&episodesResponseEpisodesEndpoint)

	if err != nil {
		t.Errorf("Expected no error on decode to %s got %s", targetUrl, err.Error())
	}

	if episodesResponseEpisodesEndpoint.Error {
		t.Errorf("Expected no error on GET %s", targetUrl)
	}

	if len(episodesResponseEpisodesEndpoint.Data) == 0 {
		t.Errorf("Expected episodes on GET %s but got none", targetUrl)
	}

	// here

	targetUrl = apiServer.URL + "/episodes/1"

	req, err = http.NewRequest(http.MethodGet, targetUrl, nil)

	if err != nil {
		t.Errorf("Expected no error on creating GET to %s got %s", targetUrl, err.Error())
	}

	res, err = apiServer.Client().Do(req)

	if err != nil {
		t.Errorf("Expected no error on GET to %s got %s", targetUrl, err.Error())
	}

	var episodeResponseGetById EpisodeResponse

	decoder = json.NewDecoder(res.Body)

	err = decoder.Decode(&episodeResponseGetById)

	if err != nil {
		t.Errorf("Expected no error on decode to %s got %s", targetUrl, err.Error())
	}

	if episodeResponseGetById.Error {
		t.Errorf("Expected no error on GET %s", targetUrl)
	}

}
