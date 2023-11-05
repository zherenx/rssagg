package main

import (
	"context"
	"log"
	"sync"
	"time"

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

	// TODO: save to db instead of loging
	log.Println(rssFeed)
}
