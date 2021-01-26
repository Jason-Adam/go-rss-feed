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

func loadConfig(fname string) (*models.RSSFeeds, error) {
	jsonFile, err := os.Open(fname)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, readErr := ioutil.ReadAll(jsonFile)
	if readErr != nil {
		log.Fatal(err)
	}

	var feeds models.RSSFeeds
	json.Unmarshal(byteValue, &feeds)

	return &feeds, nil
}

func filterFeed(feed *gofeed.Feed) {
	cutoff := time.Now().Add(-24 * time.Hour)
	filtered := []*gofeed.Item{}

	for _, f := range feed.Items {
		if time.Time(*f.PublishedParsed).After(cutoff) {
			filtered = append(filtered, f)
		}
	}
	feed.Items = filtered
	return
}

func parseFeed(url string, feedchan chan *gofeed.Feed) {
	fp := gofeed.NewParser()
	feed, feedErr := fp.ParseURL(url)
	if feedErr != nil {
		log.Printf("failed to parse %s", url)
		return
	}
	feedchan <- feed
}

func parseFeeds(urls []string) []*gofeed.Feed {
	feedchan := make(chan *gofeed.Feed, len(urls))

	for _, u := range urls {
		go parseFeed(u, feedchan)
	}

	feeds := []*gofeed.Feed{}
	for i := 0; i < len(urls); i++ {
		feeds = append(feeds, <-feedchan)
	}
	return feeds
}

func loadTemplate() *template.Template {
	templates, err := filepath.Glob("templates/*")
	if err != nil {
		log.Println(err)
	}
	t := template.Must(template.New("layout.go.html").ParseFiles(templates...))
	return t
}

func main() {
	// Load RSS Feed URLs
	rssFeeds, err := loadConfig("configs/feeds.json")
	if err != nil {
		log.Println(err)
	}
	log.Println("config file loaded successfully")

	// Fetch Feeds
	feeds := parseFeeds(rssFeeds.Urls)
	for _, f := range feeds {
		filterFeed(f)
	}

	// Load Email Templates
	t := loadTemplate()
	writer := &strings.Builder{}
	templateErr := t.Execute(writer, feeds)
	if templateErr != nil {
		log.Fatal(templateErr)
	}
	log.Println("email templates loaded successfully")

	// Email
	e := mail.NewEmailer(fromEmail, fromPassword, host, port)
	mailErr := e.Send(writer.String(), toEmail)
	if mailErr != nil {
		log.Fatal(mailErr)
	}
}
