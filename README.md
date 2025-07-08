# Intro

Pomodoro focus timer, with golang in the backend and Svelte in the frontend.

## TODO

1. Have a background goroutine that manages the countdown. When the timer reaches zero, it transitions to the next phase.
2. The frontend can poll the backend (using HTTP requests) to get the current state and remaining time.
3. The backend needs to handle starting, pausing, and resetting the timer.

## API endpoints

- /start: starts the timer
- /status: returns the current status
- /pause: pauses the timer
- /reset: resets the timer
