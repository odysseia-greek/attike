<script>
	import TraceQuery from './components/TraceQuery.svelte';
	import Sidebar from './components/Sidebar.svelte';
</script>

<main>
	<Sidebar /> <!-- Include the header for smaller screens -->
	<!-- Use the TraceQuery component -->
	<h1>Praxiteles</h1>
	<TraceQuery />
</main>

<style>
	body {
		margin: 0; /* Reset margin */
		padding: 0; /* Reset padding */
	}

	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
		background: #727981;
		min-height: 100vh; /* Make sure main covers the entire viewport height */
	}


	h1 {
		color: #ff7100;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>
