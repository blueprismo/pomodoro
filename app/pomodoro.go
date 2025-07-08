package main

import (
	//"fmt"
	"time"
	"sync"
)

const (
	ThreeHourDuration         time.Duration = time.Duration(3) * time.Hour
	SixtyMinutesDuration      time.Duration = time.Duration(60) * time.Minute
	TwentyFiveMinutesDuration time.Duration = time.Duration(25) * time.Minute
	FifteenMinutesDuration    time.Duration = time.Duration(15) * time.Minute
	FiveMinutesDuration       time.Duration = time.Duration(5) * time.Minute
	TenSecondsDuration        time.Duration = time.Duration(10) * time.Second
	FiveSecondsDuration       time.Duration = time.Duration(5) * time.Second
)

type phases = map[string]time.Duration

var threePhases = phases{
	"work":       TwentyFiveMinutesDuration,
	"longBreak":  FifteenMinutesDuration,
	"shortBreak": FiveMinutesDuration,
}

type PomodoroState struct {
	// The pomodoro should be like
	// cycle1 -           - cycle2 -            - cycle3 -			  -  cycle 4 -
	// Work -> ShortBreak -> Work -> ShortBreak -> Work -> ShortBreak -> Work -> LongBreak
	mu           sync.Mutex	  // Allow multiple http requests without race conditions
	phase        string
	startedAt    time.Time     // Original duration of current phase
	workDuration time.Duration // work Duration
	shortBreak   time.Duration // short break duration
	longBreak    time.Duration // long break duration
	remaining    time.Duration // Time left in current phase
	timer        *time.Timer
	cycleCount   int
	isRunning    bool
}

// Constructor for new Pomodoro instances
func NewPomodoro(duration time.Duration) *PomodoroState {
	return &PomodoroState{
		phase:        "work",
		remaining:    duration,
		workDuration: duration,
		shortBreak:   FiveMinutesDuration,
		longBreak:    FifteenMinutesDuration,
		cycleCount:   0,
	}
}

func (ps *PomodoroState) startTimer() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	// Start counting down
	ps.timer = time.NewTimer(ps.remaining)
	ps.startedAt = time.Now()
	ps.isRunning = true
}

func (ps *PomodoroState) pauseTimer() {
	// Ensure no other handler can tamper the state
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Check startTimer() has been called, and stop the timer.
	if ps.timer != nil {
		ps.timer.Stop()
	}
	
	// set the remaining time
	ps.remaining = ps.remaining - time.Since(ps.startedAt)
	ps.isRunning = false
}

func (ps *PomodoroState) resumeTimer() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// reset the timer
	if ps.isRunning == false {
		// Try to stop the timer first (safe even if it's already stopped)
		stopped := ps.timer.Stop()
		if !stopped {
			// Drain the channel if timer already fired
			// For a Timer created with NewTimer, Reset should be invoked only on stopped or expired timers with drained channels!
			select {
			case <-ps.timer.C:
			default:
			}
		}

		ps.timer.Reset(ps.remaining)
		ps.isRunning = true
	}
}

func (ps *PomodoroState) transitionState() string {
	// This method will change the state once the timer has reached it's countdown.
	ps.mu.Lock()
	defer ps.mu.Unlock()
	var returnstate string

	switch ps.phase {
	case "work":
		
		if ps.cycleCount < 3 {
			ps.cycleCount++
			ps.timer.Reset(ps.shortBreak)
			returnstate = "shortBreak"
		} else {
			ps.cycleCount++
			ps.timer.Reset(ps.longBreak)
			returnstate = "longBreak"
		}
	case "shortBreak":
		ps.timer.Reset(ps.workDuration)
		returnstate = "work"
	case "longBreak":
		ps.cycleCount = 0
		returnstate = "work"
	}
	ps.phase = returnstate

	return returnstate
}
