package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type notificationHandler struct {
	mu            sync.RWMutex
	notifications []string
}

func main() {

	nf := new(notificationHandler)

	http.HandleFunc("/list", nf.listNotifications)
	http.HandleFunc("/notify", nf.createNotification)
	http.HandleFunc("/clear", nf.clearNotifications)
	http.HandleFunc("/verify", nf.dumpNotifications)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("Listening now at port 8080")
	http.ListenAndServe(":8080", nil)
}

func (nh *notificationHandler) listNotifications(w http.ResponseWriter, req *http.Request) {
	nh.mu.RLock()
	defer nh.mu.RUnlock()
	for _, notification := range nh.notifications {
		fmt.Fprintf(w, "<div class=\"entry\"><span>%s</span></div>", notification)
	}
}

func (nh *notificationHandler) createNotification(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("Bad Request from client. Err: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not read request data")
		return
	}

	nh.mu.Lock()
	defer nh.mu.Unlock()
	nh.notifications = append(nh.notifications, html.EscapeString(string(body)))
	fmt.Fprintf(w, "Created")
}

func (nh *notificationHandler) clearNotifications(w http.ResponseWriter, req *http.Request) {
	nh.mu.Lock()
	defer nh.mu.Unlock()
	nh.notifications = nil
	fmt.Fprintf(w, "Cleared")
}

func (nh *notificationHandler) dumpNotifications(w http.ResponseWriter, req *http.Request) {
	nh.mu.RLock()
	defer nh.mu.RUnlock()
	jsonResult, err := json.Marshal(nh.notifications)
	if err != nil {
		log.Fatalf("Could not convert to JSON Err: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not write Notifications as JSON")
		return
	}
	fmt.Fprint(w, string(jsonResult))

}
