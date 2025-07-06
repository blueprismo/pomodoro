package main

import (
	"log"
	"encoding/json"
	"net/http"
)

type PomodoroServer struct {
	Pomodoro *PomodoroState
}

func (s *PomodoroServer) StartHandler(w http.ResponseWriter, r *http.Request) {
	// This will start the handler and return OK on succesful response or err
	if !s.Pomodoro.isRunning {
		s.Pomodoro.startTimer()
		log.Print("STARTED POMODORO")
	} else {
		log.Print("POMODORO ALREADY STARTED!")
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
