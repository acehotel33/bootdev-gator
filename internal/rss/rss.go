package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Add("User-Agent", "gator")

	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	rssFeed := &RSSFeed{}
	if err := xml.Unmarshal(body, rssFeed); err != nil {
		return &RSSFeed{}, err
	}
	rssFeed.unescapeString()

	return rssFeed, nil
}

func (rssFeed *RSSFeed) unescapeString() {
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}
}
