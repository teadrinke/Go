package main

import (
	"time"
	"log"
	"sync"
	"context"
	"github.com/teadrinke/Go/internal/database"
)

func startScraping(db *database.Queries, 
	concurrency int, 
	timeBetweenRequest time.Duration){
	log.Printf("Scraping on %v goroutined every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("Error fetching feeds: %v", err)
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

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed %d as fetched: %v", feed.ID, err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed %s: %v", feed.Url, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {	
		log.Println("Found post:", item.Title, "from feed:", feed.Name)
	}
	log.Printf("Fetched %d items from feed %s", len(rssFeed.Channel.Item), feed.Name)
}