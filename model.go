package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/zherenx/rssagg/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}
}

// TODO: maybe we could use a generic function?
func databaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := make([]Feed, 0, len(dbFeeds))

	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, databaseFeedToFeed(dbFeed))
	}

	return feeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

func databaseFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	feedFollows := make([]FeedFollow, 0, len(dbFeedFollows))

	for _, dbFeedFollow := range dbFeedFollows {
		feedFollows = append(feedFollows, databaseFeedFollowToFeedFollow(dbFeedFollow))
	}

	return feedFollows
}

type Post struct {
	Name        string    `json:"feed_name"`
	FeedUrl     string    `json:"feed_url"`
	ID          uuid.UUID `json:"post_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"` // TODO: probably should use a different type
	PublishedAt time.Time `json:"published_at"`
	PostUrl     string    `json:"post_url"`
}

func databasePostsForUserRowToPost(dbPostsForUserRow database.GetPostsForUserRow) Post {
	var description *string
	if dbPostsForUserRow.Description.Valid {
		description = &dbPostsForUserRow.Description.String
	}
	return Post{
		Name:        dbPostsForUserRow.Name,
		FeedUrl:     dbPostsForUserRow.FeedUrl,
		ID:          dbPostsForUserRow.ID,
		Title:       dbPostsForUserRow.Title,
		Description: description,
		PublishedAt: dbPostsForUserRow.PublishedAt,
		PostUrl:     dbPostsForUserRow.PostUrl,
	}
}

func databasePostsForUserRowsToPosts(dbPostsForUserRows []database.GetPostsForUserRow) []Post {
	posts := make([]Post, 0, len(dbPostsForUserRows))
	for _, row := range dbPostsForUserRows {
		posts = append(posts, databasePostsForUserRowToPost(row))
	}
	return posts
}
