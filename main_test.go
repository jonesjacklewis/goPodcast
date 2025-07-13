package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

func TestFetchRssFeed(t *testing.T) {

	t.Run("When Invalid URL provided, then should return an error", func(t *testing.T) {
		_, err := fetchRssFeed("http://this-is-an-invalid-url.invalidurl")

		if err == nil {
			t.Error("Expected an error, but none received")
		}

		if !strings.HasPrefix(err.Error(), "error when making request") {
			t.Errorf("Expected error to start with 'error when making request', got %s", err.Error())
		}
	})

	t.Run("Whens status code not 200, should return an error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		defer server.Close()

		_, err := fetchRssFeed(server.URL)

		expectedErrorMessage := "unsuccessful response from feed: 404"

		if err == nil {
			t.Errorf("Expected to receive an error")
		}

		if err.Error() != expectedErrorMessage {
			t.Errorf("Expected '%s' got '%s'", expectedErrorMessage, err.Error())
		}
	})

	t.Run("Whens invalid RSS, should return an error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application.xml")
			w.Write([]byte("Hello World"))
		}))

		defer server.Close()

		_, err := fetchRssFeed(server.URL)

		if err == nil {
			t.Errorf("Expected to receive an error")
		}

		if !strings.HasPrefix(err.Error(), "invalid XML") {
			t.Errorf("Expected error to start with 'invalid XML', got %s", err.Error())
		}
	})

	t.Run("When Valid RSS then should return valid RSS Struct with a title", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application.xml")
			w.Write([]byte(testXml))
		}))

		defer server.Close()

		xml, err := fetchRssFeed(server.URL)

		if err != nil {
			t.Errorf("Expected no err got %s", err)
		}

		if xml.Channel.Title != "Lateral with Tom Scott" {
			t.Errorf("Expected 'Lateral with Tom Scott' got %s", xml.Channel.Title)
		}
	})
}
