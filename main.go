package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/zherenx/rssagg/internal/database"

	_ "github.com/lib/pq"
)

// hold a connection to a database
type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// connect to the database
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	v1Router.Post("/create_user", apiCfg.handlerCreateUser)
	// v1Router.Get("/get_user", apiCfg.handlerGetUser)
	v1Router.Get("/get_user", apiCfg.AuthMiddleware(handlerGetUser))

	v1Router.Post("/feeds", apiCfg.AuthMiddleware(apiCfg.handlerCreateFeedWithUser))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiCfg.AuthMiddleware(apiCfg.HandlerCreateFeedFollowForUser))
	v1Router.Get("/feed_follows", apiCfg.AuthMiddleware(apiCfg.HandlerGetFeedFollowsOfUser))
	// Note: http delete requests do not typically have a body in the payload,
	// it's more conventional to pass the id in the http path
	v1Router.Delete("/feed_follows/{feedFollowId}", apiCfg.AuthMiddleware(apiCfg.HandlerDeleteFeedFollowForUser))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
