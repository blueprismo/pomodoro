package main

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStartHandler(t *testing.T) {
  ps := NewPomodoro()
  ps.Durations.Work = FiveSecondsDuration
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
	ps := NewPomodoro()
	ps.Durations.Work = FiveSecondsDuration
	server := &PomodoroServer{Pomodoro: ps}
	
	req := httptest.NewRequest(http.MethodGet, "/pause", nil)
	w := httptest.NewRecorder()
	server.PauseHandler(w,req)

	if w.Code != http.StatusOK {
		t.Error("Error, /pause endpoint not reachable")
	}
}

func TestStatusHandler(t *testing.T) {
	ps := NewPomodoro(PomodoroDurations{
		50 * time.Millisecond,
		10 * time.Millisecond,
		20 * time.Millisecond,
	})
	
	var pstatus PomodoroStatus
	server := &PomodoroServer{Pomodoro: ps}
	
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	server.StatusHandler(w,req)

	if w.Code != http.StatusOK {
		t.Error("Error, /status endpoint not reachable")
	}
	json.Unmarshal(w.Body.Bytes(), &pstatus)
	
	//fmt.Printf("%+v", pstatus)

	if pstatus.Remaining <= 0 {
		t.Errorf("Error, Remaining time <=0, got %v", pstatus.Remaining)
	}

	if (pstatus.Mode != "work" && pstatus.Mode != "shortBreak" && pstatus.Mode != "longBreak") {
		t.Errorf("Error, unknown mode %v", pstatus.Mode)
	}
}

