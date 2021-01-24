package main

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

var rssFeeds = []string{
	"https://c-for-dummies.com/blog/?feed=rss2",
	"https://threedots.tech/index.xml",
	"https://blog.golang.org/feed.atom?format=xml",
	"https://dave.cheney.net/feed/atom",
	"https://changelog.com/gotime/feed",
}

type slimFeedItem struct {
	title   string
	link    string
	pubDate time.Time
}

func newSlimFeedItem(title string, link string, pubDate time.Time) *slimFeedItem {
	return &slimFeedItem{
		title:   title,
		link:    link,
		pubDate: pubDate,
	}
}

type finalArticle struct {
	site          string
	slimFeedItems []*slimFeedItem
}

func newFinalArticle(site string, slimFeedItems []*slimFeedItem) *finalArticle {
	return &finalArticle{
		site:          site,
		slimFeedItems: slimFeedItems,
	}
}

func mail(content string) error {
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	addr := host + ":" + port
	to := os.Getenv("MAIL_TO")
	from := os.Getenv("MAIL_FROM")
	pw := os.Getenv("MAIL_PASSWORD")
	auth := smtp.PlainAuth("", from, pw, host)

	subject := "RSS Feeds for " + time.Now().Format("Jan 02, 2006")

	msg := strings.Builder{}
	msg.WriteString("From: \"Feed Update\" <" + from + ">\n")
	msg.WriteString("To: " + to + "\n")
	msg.WriteString("Subject: " + subject + "\n")
	msg.WriteString("MIME-version: 1.0;\n")
	msg.WriteString("Content-Type: text/html;charset=\"UTF-8\";\n")
	msg.WriteString("\n")
	msg.WriteString(content)

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg.String()))
	if err != nil {
		return err
	}
	return nil
}

func filterFeed(feed *gofeed.Feed) []*slimFeedItem {
	cutoff := time.Now().Add(-24 * time.Hour)
	filtered := make([]*slimFeedItem, 0)

	for _, f := range feed.Items {
		if time.Time(*f.PublishedParsed).After(cutoff) {
			slimF := newSlimFeedItem(f.Title, f.Link, *f.PublishedParsed)
			filtered = append(filtered, slimF)
		}
	}

	return filtered
}

func parseFeed(link string) (feed *gofeed.Feed, err error) {
	fp := gofeed.NewParser()
	feed, feedErr := fp.ParseURL(link)
	if feedErr != nil {
		return nil, feedErr
	}
	return feed, nil
}

func main() {
	final := []*finalArticle{}

	for _, l := range rssFeeds {
		feed, err := parseFeed(l)
		if err != nil {
			fmt.Println(err)
		}
		ff := filterFeed(feed)
		fa := newFinalArticle(l, ff)
		final = append(final, fa)
	}

	for _, fa := range final {
		for _, f := range fa.slimFeedItems {
			fmt.Println(f.link)
		}
	}
}
