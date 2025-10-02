package main

import (
	"PointsInterviewTestServer/internal"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := internal.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
