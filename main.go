package main

import (
	"log"
	"net/http"
	"time"

	ps "pubsub.com/pubsub/pubsub"
	"pubsub.com/pubsub/session"
	"pubsub.com/pubsub/web"
)

func main() {
	repo := session.LoadRepo()
	loggedInUsers := make(map[string]ps.User)
	handler := web.NewAppHandler(repo, loggedInUsers)
	mux := web.MakeServer(handler)

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	log.Fatal(s.ListenAndServe())
}
