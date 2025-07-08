package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

type PomodoroServer struct {
	Pomodoro *PomodoroState
}

//go:embed frontend/dist/*
var embeddedFiles embed.FS

var embedAssets fs.FS
var distFS fs.FS

func init() {
	var err error
	// debug what embedFS gets
	entries, _ := embeddedFiles.ReadDir("frontend/dist")
	for _, e := range entries {
		fmt.Println("Entry:", e.Name(), "IsDir:", e.IsDir())
	}

	
	distFS, err = fs.Sub(embeddedFiles, "frontend/dist")
	if err != nil {
		panic(err)
	}
	embedAssets, err = fs.Sub(embeddedFiles, "frontend/dist/assets")
	if err != nil {
		panic(err)
	}
}

// Main handler ("/" request)
func MainHandler(w http.ResponseWriter, r *http.Request) {
	// Main landing page
	w.Header().Set("Content-Type", "text/html")
	config := map[string]int{"work": 25, "shortBreak": 5, "longBreak": 15}
	jsonBytes, _ := json.Marshal(config)

	// Read embedded index.html
	htmlBytes, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}

	// Inject config before </head>
	withConfig := strings.Replace(
		string(htmlBytes),
		"</head>",
		fmt.Sprintf(`<script>window.POMODORO_CONFIG = %s;</script></head>`, jsonBytes),
		1,
	)

	w.Write([]byte(withConfig))
}

// Pomodoro Handlers
func (s *PomodoroServer) StartHandler(w http.ResponseWriter, r *http.Request) {
	// This will start the handler and return OK on succesful response or err
	if !s.Pomodoro.isRunning {
		s.Pomodoro.startTimer()
		log.Print("STARTED POMODORO")
	} else {
		log.Print("POMODORO ALREADY STARTED, nothing to do")
	}
	response := map[string]bool{"running": true}
	jsonMsg, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonMsg)
}

func (s *PomodoroServer) PauseHandler(w http.ResponseWriter, r *http.Request) {
	// This will pause the Pomodoro
	if s.Pomodoro.isRunning {
		s.Pomodoro.pauseTimer()
		log.Print("Pomodoro paused")
	} else {
		log.Print("Pomodoro already in a non-running state")
	}

	response := map[string]bool{"running": false}
	jsonMsg, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Could not marshall JSON response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonMsg)
}

// func (s *PomodoroServer) StatusHandler(w http.ResponseWriter, r *http.Request) {
// 	// returns the actual status of our Pomodoro Timer
// 	response := PomodoroState{"": true}
// }
