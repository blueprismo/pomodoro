package main

import (
	"testing"
	"time"
)

func TestInitTimer(t *testing.T) {
	var desiredInitState string = "work"

	pt := NewPomodoro()
	pt.durations.Work = SixtyMinutesDuration

	if pt.mode != desiredInitState {
		t.Errorf("Initmode NOT ok, got %v want %v", pt.mode, desiredInitState)
	}

	if pt.remaining >= ThreeHourDuration {
		t.Errorf("Error, init working time too long: got %v want less than %v", pt.remaining, ThreeHourDuration)
	}

	if pt.cycleCount != 0 {
		t.Errorf("Error, Cycle count should be 0, got %v", pt.cycleCount)
	}
}

func TestWorkCountdown(t *testing.T) {
	pt := NewPomodoro(PomodoroDurations{15 * time.Millisecond,5 * time.Millisecond,10 * time.Millisecond})
	pt.startTimer()

	// check the timer status
	select {
	case <-pt.timer.C:
		t.Log("Timer expired as expected")
	case <-time.After(1 * time.Second):
		t.Error("Timer did not expire in expected time")
	}
}

func TestPauseTimer(t *testing.T) {
	// create a new pomodoro
	pt := NewPomodoro()
	pt.durations.Work = 100 * time.Millisecond
	pt.durations.ShortBreak = 20 * time.Millisecond
	pt.durations.LongBreak = 40 * time.Millisecond

	// Test that the timer can be paused
	pt.startTimer()
	time.Sleep(10 * time.Millisecond)
	pt.pauseTimer()

	// Check resulting state
	if pt.isRunning {
		t.Errorf("ERROR: Timer is running")
	}
	// We slept 10 out of 100 Milliseconds, so we should be around 90 Milliseconds, give 5ms buffer
	if pt.remaining > 100 *time.Millisecond || pt.remaining < 85 *time.Millisecond {
		t.Errorf("Expected remaining time between 100ms and 85ms, got %v", pt.remaining)
	}
}

func TestResumeTimer(t *testing.T) {
	pt := NewPomodoro()
	pt.durations.Work = 100 * time.Millisecond
	pt.durations.ShortBreak = 20 * time.Millisecond
	pt.durations.LongBreak = 40 * time.Millisecond

	pt.startTimer()
	time.Sleep(10 * time.Millisecond)
	pt.pauseTimer()

	remaining := pt.remaining

	if remaining > 100 * time.Millisecond || remaining < 85*time.Millisecond {
		t.Errorf("Unexpected remaining time after pause: got %v", remaining)
	}

	pt.resumeTimer()

	select {
	case <-pt.timer.C:
		t.Log("Timer expired after resume as expected")
	case <-time.After(remaining + 1*time.Second): // give some buffer
		t.Error("Timer did not expire after resume in expected time")
	}

	if !pt.isRunning {
		t.Error("Timer should be running after resume")
	}
}

func assertTransitionState(t *testing.T, currentCicle, desiredCicle int, currentmode, desiredmode string) {
	// helper function to check transition state results
	if currentCicle != desiredCicle {
		t.Errorf("Cycle count unexpected: want %v, got %v", desiredCicle, currentCicle)
	}
	if currentmode != desiredmode {
		t.Errorf("Unexpected mode: want %v, got %v", desiredCicle, currentCicle)
	}
}

func expectTransition(t *testing.T, pt *PomodoroState, expectedCycle int, expectedmode string) {
	// function to check transitioning state
	select {
	case <-pt.timer.C:
		t.Logf("Timer expired, transitioning to %s", expectedmode)
		pt.transitionState()
	case <-time.After(1 * time.Second):
		t.Fatalf("Timer did not expire in expected time for mode %s", expectedmode)
	}
	assertTransitionState(t, pt.cycleCount, expectedCycle, pt.mode, expectedmode)
}

func TestTransitionState(t *testing.T) {
	// In this test we are going to test the whole state machine and all the transitions until the first one.
	// Work -> ShortBreak -> Work -> ShortBreak -> Work -> ShortBreak -> Work -> LongBreak -> reset
	// BEGIN -> Work -[1 tr]-> ShortBreak -[2 tr]> Work -[3tr]-> ShortBreak -[4tr]-> Work -[5tr]-> ShortBreak -[6tr]-> Work -[7tr]-> LongBreak -[8tr]-> BEGIN
	pt := NewPomodoro()
	pt.durations.Work = 20 * time.Millisecond
	pt.durations.ShortBreak = 5 * time.Millisecond
	pt.durations.LongBreak = 15 * time.Millisecond

	// start pomodoro and evaluate first transition (work -> shortBreak)
	pt.startTimer()

	expectTransition(t, pt, 1, "shortBreak")
	expectTransition(t, pt, 1, "work")
	expectTransition(t, pt, 2, "shortBreak")
	expectTransition(t, pt, 2, "work")
	expectTransition(t, pt, 3, "shortBreak")
	expectTransition(t, pt, 3, "work")
	expectTransition(t, pt, 4, "longBreak")
	expectTransition(t, pt, 0, "work") // reset to BEGIN -> work
}


