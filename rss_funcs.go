package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ablanchetmd/gator/internal/config"
	"github.com/ablanchetmd/gator/internal/database"
	"github.com/lib/pq" // PostgreSQL driver
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	rssbyte, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(rssbyte, &feed)
	unescapeFeed(&feed)

	if err != nil {
		return nil, err
	}
	return &feed, nil
}

func scrapeFeeds(s *config.State) error {
	ctx := context.Background()
	current_time := sql.NullTime{Time: time.Now(), Valid: true}
	feed, err := s.Db.GetNextFeedToFetch(ctx, current_time)
	if err != nil {
		return err
	}

	s.Db.MarkFeedAsFetched(ctx, database.MarkFeedAsFetchedParams{
		ID:            feed.ID,
		LastFetchedAt: current_time,
	})

	feedData, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}

	for _, item := range feedData.Channel.Item {

		parsedPublishedTime, err := time.Parse(time.RFC1123Z, item.PubDate)

		if err != nil {
			fmt.Printf("Error parsing time: %v for %s in %s \n", err, item.PubDate, feed.Url)
			continue
		}

		_, err = s.Db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: (item.Title != "")},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: (item.Description != "")},
			PublishedAt: parsedPublishedTime,
			FeedID:      feed.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				// Check for unique constraint violation				
				if pqErr.Code != "23505" {
					fmt.Printf("PostgreSQL error: %s\n", pqErr.Message)
				}
				continue
			}else{
				fmt.Printf("Error creating post: %v for %s in %s \n", err, item.Title, feed.Url)
				continue
			}					
		}
		fmt.Printf("Created post %s in %s \n", item.Title, feed.Url)
	}

	return nil
}

func unescapeRSSItem(item *RSSItem) {
	item.Title = html.UnescapeString(item.Title)
	item.Description = html.UnescapeString(item.Description)
}

func unescapeFeed(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		unescapeRSSItem(&feed.Channel.Item[i])
	}
}
