<script>
	let config = window.POMODORO_CONFIG || { work: 25, shortBreak: 5, longBreak: 15 };
	let { work, shortBreak, longBreak } = config;
	console.log("work, shortbreak and longbreak times (minutes): " + work + ", " + shortBreak + ", " + longBreak)
	let isRunning = false;
	let timer = null;
	let mode = "Work"; // "Work", "Short Break", "Long Break"
	let secondsLeft = work * 60;

	function format(t) {
		const m = Math.floor(t / 60).toString().padStart(2, '0');
		const s = (t % 60).toString().padStart(2, '0');
		return `${m}:${s}`;
	}

	async function start() {
		try {
  		  const response = await fetch('/start');
  		  if (!response.ok) {
  		    throw new Error(`Response status: ${response.status}`);
  		  }
		  console.log("/start hit")
  		  const json = await response.json();
  		  console.log(json);
  		} catch (error) {
  		  console.error(error.message);
		}
	}

	function stop() {
		isRunning = false;
		clearInterval(timer);
	}

	function switchMode(newMode) {
		stop();
		mode = newMode;
		if (mode === "Work") secondsLeft = work * 60;
		if (mode === "Short Break") secondsLeft = shortBreak * 60;
		if (mode === "Long Break") secondsLeft = longBreak * 60;
	}
</script>

<main>
	<h1>{format(secondsLeft)}</h1>
	<p>{mode}</p>

	<div class="buttons">
		<button on:click={start}>Start</button>
		<button on:click={stop}>Stop</button>
		<button on:click={() => switchMode("Work")}>Work</button>
		<button on:click={() => switchMode("Short Break")}>Short Break</button>
		<button on:click={() => switchMode("Long Break")}>Long Break</button>
	</div>
</main>

<style>
	main {
		text-align: center;
		font-family: sans-serif;
		padding: 2rem;
	}
	h1 {
		font-size: 5rem;
		margin: 0.5rem 0;
	}
	p {
		font-size: 1.5rem;
		margin-bottom: 2rem;
	}
	.buttons button {
		margin: 0.5rem;
		padding: 0.75rem 1.5rem;
		font-size: 1rem;
		cursor: pointer;
	}
</style>
