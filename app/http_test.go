package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	//"fmt"
)

func TestStartHandler(t *testing.T) {
  ps := NewPomodoro(FiveSecondsDuration)
  server := &PomodoroServer{Pomodoro: ps}
	expected := `{"running":true}`

  req := httptest.NewRequest(http.MethodGet, "/start", nil)
  w := httptest.NewRecorder()
  server.StartHandler(w, req)

  if w.Code != http.StatusOK {
		t.Errorf("Error, pomodoro could not start")
	}
	if w.Body.String() != expected {
		t.Errorf("Error! could not get the started status")
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}
}

func TestPauseHandler(t *testing.T) {
	ps := NewPomodoro(FiveSecondsDuration)
	server := &PomodoroServer{Pomodoro: ps}
	
	req := httptest.NewRequest(http.MethodGet, "/pause", nil)
	w := httptest.NewRecorder()
	server.PauseHandler(w,req)

	if w.Code != http.StatusOK {
		t.Error("Error, /pause endpoint not reachable")
	}

}

