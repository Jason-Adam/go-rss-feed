package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jason-adam/go-rss-feed/mail"
	"github.com/jason-adam/go-rss-feed/models"
	"github.com/mmcdole/gofeed"
)

var (
	fromEmail    string = os.Getenv("MAIL_FROM")
	fromPassword string = os.Getenv("MAIL_PASSWORD")
	toEmail      string = os.Getenv("MAIL_TO")
	port         string = os.Getenv("MAIL_PORT")
	host         string = os.Getenv("MAIL_HOST")
)

type CompletedFeed struct {
	Success bool
	Feed    *gofeed.Feed
}

func loadConfig(fname string) (*models.RSSFeeds, error) {
	jsonFile, err := os.Open(fname)
	if err != nil {
		log.Printf("unable to open config file due to error: %s", err)
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, readErr := ioutil.ReadAll(jsonFile)
	if readErr != nil {
		log.Printf("unable to read config file due to error: %s", readErr)
		return nil, readErr
	}

	var feeds models.RSSFeeds
	uErr := json.Unmarshal(byteValue, &feeds)
	if uErr != nil {
		log.Printf("unable to unmarshal config into struct due to error: %s", uErr)
		return nil, uErr
	}

	return &feeds, nil
}

func filterFeed(feed *gofeed.Feed) {
	cutoff := time.Now().Add(-24 * time.Hour)
	filtered := []*gofeed.Item{}

	if feed.Len() > 0 {
		for _, f := range feed.Items {
			if f.PublishedParsed != nil && time.Time(*f.PublishedParsed).After(cutoff) {
				filtered = append(filtered, f)
			}
		}
	}

	feed.Items = filtered
	return
}

func parseFeed(url string, feedchan chan *CompletedFeed) {
	fp := gofeed.NewParser()
	feed, feedErr := fp.ParseURL(url)
	if feedErr != nil {
		log.Printf("failed to parse %s", url)
		feedchan <- &CompletedFeed{
			Success: false,
			Feed:    nil,
		}
	}
	log.Printf("successfully parsed feed %s", url)
	feedchan <- &CompletedFeed{
		Success: true,
		Feed:    feed,
	}
}

func parseFeeds(urls []string) []*gofeed.Feed {
	feedchan := make(chan *CompletedFeed, len(urls))

	for _, u := range urls {
		go parseFeed(u, feedchan)
	}

	feeds := []*gofeed.Feed{}
	for i := 0; i < len(urls); i++ {
		f := <-feedchan
		if f.Success == true {
			feeds = append(feeds, f.Feed)
		}
	}
	return feeds
}

func loadTemplate() (*template.Template, error) {
	templates, err := filepath.Glob("templates/*")
	if err != nil {
		log.Printf("unable to load templates due to error: %s", err)
		return nil, err
	}

	tt, pErr := template.New("layout.go.html").ParseFiles(templates...)
	if pErr != nil {
		log.Printf("unable to generate template due to parse error: %s", pErr)
		return nil, pErr
	}

	t := template.Must(tt, nil)
	return t, nil
}

func main() {
	// Load RSS Feed URLs
	rssFeeds, err := loadConfig("configs/feeds.json")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("config file loaded successfully")

	// Fetch Feeds
	feeds := parseFeeds(rssFeeds.Urls)
	for _, f := range feeds {
		filterFeed(f)
	}

	// Remove feeds with no new posts
	f := feeds[:0]
	for _, i := range feeds {
		if i.Len() > 0 {
			f = append(f, i)
		}
	}

	// Load Email Templates
	t, tErr := loadTemplate()
	if tErr != nil {
		log.Fatalln(tErr)
	}

	writer := &strings.Builder{}
	templateErr := t.Execute(writer, f)
	if templateErr != nil {
		log.Fatalln(templateErr)
	}

	log.Println("email templates loaded successfully")

	// Email
	e := mail.NewEmailer(fromEmail, fromPassword, host, port)
	mailErr := e.Send(writer.String(), toEmail)
	if mailErr != nil {
		log.Fatal(mailErr)
	}
}
