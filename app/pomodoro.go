package main

import (
	//"fmt"
	"sync"
	"time"
)

const (
	ThreeHourDuration    time.Duration = time.Duration(3) * time.Hour
	SixtyMinutesDuration time.Duration = time.Duration(60) * time.Minute
	TenSecondsDuration   time.Duration = time.Duration(10) * time.Second
	FiveSecondsDuration  time.Duration = time.Duration(5) * time.Second
)

type PomodoroDurations struct {
	Work       time.Duration
	ShortBreak time.Duration
	LongBreak  time.Duration
}

type PomodoroTimer struct {
	// The pomodoro should be like
	// cycle1 -           - cycle2 -            - cycle3 -			  -  cycle 4 -
	// Work -> ShortBreak -> Work -> ShortBreak -> Work -> ShortBreak -> Work -> LongBreak
	mu         sync.Mutex
	startedAt  time.Time         // Original duration of current mode	
	timer      *time.Timer		 
	CycleCount int  			 `json:"cycleCount"` // current cycle
	IsRunning  bool 			 `json:"isRunning"`
	Durations  PomodoroDurations `json:"durations"`                  // 'work,shortbreak,longbreak' durations
	Remaining  time.Duration     `json:"remaining"` // Time left in current mode
	Mode       string            `json:"mode"`
}

type PomodoroStatus struct {
	// Pomodoro status response for API calls
    Mode       string            `json:"mode"`
    Durations  PomodoroDurations `json:"durations"`
    Remaining  time.Duration     `json:"remaining"`
    CycleCount int               `json:"cycleCount"`
    IsRunning  bool              `json:"isRunning"`
}

// Constructor for new Pomodoro instances
func NewPomodoro(durations ...PomodoroDurations) *PomodoroTimer {
	var d PomodoroDurations
	if len(durations) > 0 {
		d = durations[0]
	} else {
		// set default durations
		d = PomodoroDurations{
			Work:       25 * time.Minute,
			ShortBreak: 5 * time.Minute,
			LongBreak:  15 * time.Minute,
		}
	}

	return &PomodoroTimer{
		Mode:       "work",
		Durations:  d,
		CycleCount: 0,
		Remaining: d.Work,
	}
}

func (ps *PomodoroTimer) startTimer() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	// Start counting down
	ps.Remaining = ps.Durations.Work
	ps.timer = time.NewTimer(ps.Remaining)
	ps.startedAt = time.Now()
	ps.IsRunning = true
}

func (ps *PomodoroTimer) pauseTimer() {
	// Ensure no other handler can tamper the state
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Check startTimer() has been called, and stop the timer.
	if ps.timer != nil {
		ps.timer.Stop()
	}

	// set the Remaining time
	ps.Remaining = ps.Remaining - time.Since(ps.startedAt)
	ps.IsRunning = false
}

func (ps *PomodoroTimer) resumeTimer() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// reset the timer
	if ps.IsRunning == false {
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

		ps.timer.Reset(ps.Remaining)
		ps.IsRunning = true
	}
}

func (ps *PomodoroTimer) transitionState() string {
	// This method will change the state once the timer has reached it's countdown.
	ps.mu.Lock()
	defer ps.mu.Unlock()
	var nextstate string

	switch ps.Mode {
	case "work":
		if ps.CycleCount < 3 {
			ps.CycleCount++
			ps.timer.Reset(ps.Durations.ShortBreak)
			nextstate = "shortBreak"
		} else {
			ps.CycleCount++
			ps.timer.Reset(ps.Durations.LongBreak)
			nextstate = "longBreak"
		}
	case "shortBreak":
		ps.timer.Reset(ps.Durations.Work)
		nextstate = "work"
	case "longBreak":
		ps.CycleCount = 0
		ps.timer.Reset(ps.Durations.Work)
		nextstate = "work"
	}
	ps.Mode = nextstate

	return nextstate
}

func (ps *PomodoroTimer) getCurrentStatus() PomodoroTimer {
	// returns a copy of the actual status
	ps.mu.Lock()
	defer ps.mu.Unlock()

	return PomodoroTimer{
		Mode:       ps.Mode,
		Durations:  ps.Durations,
		Remaining:  ps.Remaining,
		CycleCount: ps.CycleCount,
		IsRunning:  ps.IsRunning,
	}
}
