<script>
	let config = window.POMODORO_CONFIG || {
		work: 25,
		shortBreak: 5,
		longBreak: 15,
	};
	let { work, shortBreak, longBreak } = config;
	console.log(
		"work, shortbreak and longbreak times (minutes): " +
			work +
			", " +
			shortBreak +
			", " +
			longBreak,
	);

	// Load initial state.
	let isRunning = false;
	let isFetching = false;
	let nextstate, timer;
	let cycleCount = 0;
	let mode = "work"; // "work", "shortBreak", "longBreak"
	let secondsLeft = work; //%2;

	function format(t) {
		const m = Math.floor(t / 60)
			.toString()
			.padStart(2, "0");
		const s = (t % 60).toString().padStart(2, "0");
		return `${m}:${s}`;
	}

	function countdown() {
		timer = setInterval(() => {
			// console.log("Tick:", secondsLeft);  // Debug line
			if (secondsLeft > 0) {
				secondsLeft--;
			} else {
				console.log("Pomodoro transitioning");
				nextMode();
			}
		}, 1000); // every 1 second
	}

	async function startstop() {
		if (isFetching) return; // No-op if we're mid-request (throttle)

		isFetching = true;
		try {
			if (!isRunning) {
				console.log("/start hit");
				const response = await fetch("/start");
				if (!response.ok)
					throw new Error(`Response status: ${response.status}`);

				const json = await response.json();
				countdown();
				console.log(json);
				isRunning = true;
			} else {
				const response = await fetch("/pause");
				if (!response.ok)
					throw new Error(`Response status: ${response.status}`);

				const json = await response.json();
				console.log(json);
				clearInterval(timer);
				isRunning = false;
			}
		} catch (err) {
			console.error(err.message);
		} finally {
			isFetching = false;
		}
	}

	function nextMode(desiredMode) {
		if (desiredMode != undefined) {
			switch (desiredMode) {
				case "work":
					clearInterval(timer);
					secondsLeft = work * 60;
					isRunning = false;
					break;
				case "shortBreak":
					clearInterval(timer);
					secondsLeft = shortBreak * 60;
					isRunning = false;
					break;
				case "longBreak":
					clearInterval(timer);
					secondsLeft = longBreak * 60;
					isRunning = false;
					break;
			}
			mode = desiredMode;
			return mode;
		}

		console.log(`Current mode before change: ${mode}`);
		switch (mode) {
			case "work":
				if (cycleCount < 3) {
					cycleCount++;
					clearInterval(timer);
					nextstate = "shortBreak";
					secondsLeft = shortBreak * 60;
				} else {
					cycleCount++;
					clearInterval(timer);
					nextstate = "longBreak";
					secondsLeft = longBreak * 60;
				}
				isRunning = false;
				break;

			case "shortBreak":
				clearInterval(timer);
				nextstate = "work";
				secondsLeft = work * 60;
				isRunning = false;
				break;

			case "longBreak":
				clearInterval(timer);
				cycleCount = 0;
				nextstate = "work";
				secondsLeft = work * 60;
				isRunning = false;
				break;

			default:
				console.log(`Error, undefined state ${mode}`);
				break;
		}

		console.log(`Current mode after change: ${nextstate}`);
		return (mode = nextstate);
	}
</script>

<main class={mode}>
	<div id="timerContainer">
		<div class="buttonsHeader buttons">
			<button class:selected={mode === "work"} on:click={() => nextMode("work")}>Work</button>
			<button class:selected={mode === "shortBreak"} on:click={() => nextMode("shortBreak")}>Short Break</button>
			<button class:selected={mode === "longBreak"}  on:click={() => nextMode("longBreak")}>Long Break</button>
		</div>
		<h1>{format(secondsLeft)}</h1>
		<p>{mode}</p>

		<div class="buttons">
			<button
				class="startStop"
				class:is-start={!isRunning}
				class:is-pause={isRunning}
				on:click={startstop}
			>
				{isRunning ? "Pause" : "Start"}
			</button>
		</div>
	</div>
</main>

<style>
	/* fonts */
	@import url("https://fonts.googleapis.com/css2?family=Roboto+Mono:ital,wght@0,100..700;1,100..700&family=Rubik:ital,wght@0,300..900;1,300..900&display=swap");

	/* basics */
	main {
		text-align: center;
		font-family: "Rubik", "Roboto Mono", "Sans Serif";
		padding: 2rem;
		height: 100vh;
		color: white;
		transition: background-color 0.5s ease;
	}

	h1 {
		font-size: 5rem;
		margin: 0.5rem 0;
	}
	p {
		font-size: 1.5rem;
		margin-bottom: 2rem;
	}
	/* Timer style*/
	#timerContainer {
		display: flex;
		flex-direction: column;
		flex-wrap: nowrap;
		width: 50%;
		margin: auto;
		background-color: rgba(255, 255, 255, 0.4);
		border-radius: 12px;
		padding: 1em;
	}

	.buttons button {
		margin: 0.5rem;
		padding: 0.75rem 1.5rem;
		font-size: 1rem;
		cursor: pointer;
	}

	.buttonsHeader button {
		border: none;
		border-radius: 12px;
		color: white;
		background-color: transparent;
	}

	.buttonsHeader > button.selected {
		background-color: rgba(0, 0, 0, 0.1);
	}

	.startStop {
		background-color: white;
	}

	.is-start,
	.is-pause {
		font-size: 1.5rem;
		font-weight: bold;
		text-transform: uppercase;
		border: none;
		border-radius: 12px;
		padding: 1rem 2.5rem;
	}
	.is-start {
		box-shadow: 0px 4px 0px rgb(224, 224, 224);
		transition: transform 0.1s ease, box-shadow 0.1s ease;
	}
	.is-start:active {
		transform: translateY(2px);
		box-shadow: 0px 2px 0px rgba(0, 0, 0, 0.2);
	}

	.work {
		background-color: tomato;
	}
	.work .startStop {
		color: tomato;
	}

	.shortBreak {
		background-color: rgb(56, 133, 138);
	}
	.shortBreak .startStop {
		color: rgb(56, 133, 138);
	}
	.longBreak {
		background-color: rgb(57, 112, 151);
	}
	.longBreak .startStop {
		color: rgb(57, 112, 151);
	}
</style>
