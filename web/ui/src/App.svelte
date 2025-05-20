<script>
	import { Router, Route, Link } from "svelte-routing";
	import { isAuthenticated, currentUser, checkAuthStatus } from './stores/auth.js';

	import Login from './components/Login.svelte';
	import Register from './components/Register.svelte';

	export let name; // This prop might not be used directly anymore if main content is route-based
	let backendMessage = ''; // Keep for the /api/hello example or remove if not needed

	async function fetchBackendMessage() {
		try {
			const response = await fetch('http://localhost:3000/api/hello', {credentials: 'include'});
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			backendMessage = data.message;
		} catch (error) {
			backendMessage = `Failed to fetch message: ${error.message}`;
			console.error('Fetch error:', error);
		}
	}

	async function handleLogout() {
		try {
			const response = await fetch('/api/logout', {
				method: 'POST',
				credentials: 'include',
			});
			if (response.ok) {
				await checkAuthStatus(); // Re-check auth status, which should update stores
			} else {
				console.error("Logout failed");
				// Optionally show an error message to the user
			}
		} catch (error) {
			console.error("Logout error:", error);
		}
	}

	// Check auth status when App mounts, though auth.js already does this
	// checkAuthStatus(); // Can be redundant if auth.js handles it on import

</script>

<Router>
	<nav>
		<div class="nav-links">
			<Link to="/">Home</Link>
			{#if $isAuthenticated}
				<Link to="/account">Account ({$currentUser?.username || 'User'})</Link>
				<button on:click={handleLogout}>Logout</button>
			{:else}
				<Link to="/register">Register</Link>
				<Link to="/login">Login</Link>
			{/if}
		</div>
	</nav>

	<main>
		<Route path="/login" component={Login} />
		<Route path="/register" component={Register} />
		<Route path="/">
			<h1>Hello {$isAuthenticated ? ($currentUser?.username || 'User') : (name || 'Guest')}!</h1>
			<p>Visit the <a href="https://svelte.dev/tutorial">Svelte tutorial</a> to learn how to build Svelte apps.</p>
			<button on:click={fetchBackendMessage}>Get Message from Backend</button>
			{#if backendMessage}
				<p>Backend says: <strong>{backendMessage}</strong></p>
			{/if}
		</Route>
		<Route path="/account">
			{#if $isAuthenticated}
				<h2>Account Management</h2>
				<p>Welcome, {$currentUser?.username}!</p>
				<p>Email: {$currentUser?.email}</p>
				<p>(Account update functionality to be implemented)</p>
			{:else}
				<p>Please <Link to="/login">login</Link> to view your account.</p>
			{/if}
		</Route>
	</main>
</Router>

<style>
  nav {
    display: flex;
    justify-content: space-between; /* Adjusted for home link on left */
    align-items: center;
    padding: 1em;
    background-color: #f0f0f0;
  }

  .nav-links a,
  .nav-links button {
    margin-left: 1em;
    text-decoration: none;
    color: #333;
    background: none;
    border: none;
    cursor: pointer;
    font-size: inherit;
    font-family: inherit;
  }

  .nav-links a:hover,
  .nav-links button:hover {
    color: #ff3e00;
  }

	main {
		text-align: center;
		padding: 1em;
		/* max-width: 240px; */ /* Removed to allow wider content for forms */
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
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