package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zherenx/rssagg/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {

	log.Printf("Start scraping on %v goroutines every %s duration...\n", concurrency, timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	// Note: the following way will exec the loop immediately
	// while "for range ticker.C" will wait timeBetweenRequests first
	for ; ; <-ticker.C {

		log.Printf("Fetching next %v feeds...\n", concurrency)

		feeds, err := db.GetNextNFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	log.Printf("Attempting to fetch rss feed from %s\n", feed.Url)

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	_, err = db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	for _, item := range rssFeed.Channel.Items {

		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		} else {
			description.Valid = false
		}

		// TODO: should support all/more time formats
		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println("Error parsing date string:", err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: t,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			// TODO: this feels like a hack, should we do something like
			// comparing the time in some way to prevent posting duplicate
			// request altogether?
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Error writing post to database:", err)
			continue
		}
	}

	log.Printf("Feed %s collected, %v posts found\n", feed.Name, len(rssFeed.Channel.Items))
}
