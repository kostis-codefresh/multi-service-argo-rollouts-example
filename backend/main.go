package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "up")
		// http.Error(w, "down!", http.StatusServiceUnavailable)
	})

	http.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "yes")
	})

	http.HandleFunc("/api/v1/interest", func(w http.ResponseWriter, r *http.Request) {
		randomSource := rand.NewSource(time.Now().UnixNano())
		calculatedInterest := rand.New(randomSource)
		fmt.Fprint(w, (calculatedInterest.Intn(26) + 10))
	})

	http.HandleFunc("/api/", http.NotFound)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html>I calculate interests. Call <a href='api/v1/interest'>api/v1/interest</a> to get your quote.</html>")
	})

	fmt.Println("Listening now at port 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
