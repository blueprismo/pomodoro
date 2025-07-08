package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// HTML landing
	mux.HandleFunc("/", MainHandler)

	// API
	ps := NewPomodoro()
	server := &PomodoroServer{Pomodoro: ps}
	mux.HandleFunc("/start", server.StartHandler)
	mux.HandleFunc("/pause", server.PauseHandler)

	// Static assets (JS, CSS, etc.)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(embedAssets))))
	
	log.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
