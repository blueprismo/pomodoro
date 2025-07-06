package main

import (
	"time"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(`
            <html><body>
                <button onclick="fetch('/start').then(r => r.json()).then(console.log)">Start Pomodoro</button>
            </body></html>
        `))
    })
	// Register your /start API endpoint
    ps := NewPomodoro(25 * time.Minute)
    server := &PomodoroServer{Pomodoro: ps}
    mux.HandleFunc("/start", server.StartHandler)

	err := http.ListenAndServe(":8080",mux)


	if err != nil {
		log.Fatal(err)
	}
}