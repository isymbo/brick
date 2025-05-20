<!-- web/ui/src/components/Login.svelte -->
<script>
    import { isAuthenticated, currentUser, checkAuthStatus } from '../stores/auth.js';
    import { navigate } from 'svelte-routing'; // Assuming svelte-routing

    let username = ''; // Can be username or email
    let password = '';
    let message = '';
    let isError = false;

    async function handleLogin() {
        isError = false;
        message = 'Logging in...';
        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password }),
                credentials: 'include', // Important to send/receive cookies
            });
            const data = await response.json();
            if (response.ok) {
                message = data.message;
                await checkAuthStatus(); // Update auth state
                if ($isAuthenticated) {
                    navigate("/"); // Navigate to home or dashboard
                }
            } else {
                isError = true;
                message = data.error || 'Login failed';
                isAuthenticated.set(false);
                currentUser.set(null);
            }
        } catch (error) {
            isError = true;
            message = 'An error occurred: ' + error.message;
            isAuthenticated.set(false);
            currentUser.set(null);
            console.error("Login error:", error);
        }
    }
</script>

<div class="auth-form">
    <h2>Login</h2>
    <form on:submit|preventDefault={handleLogin}>
        <div>
            <label for="username">Username or Email:</label>
            <input type="text" id="username" bind:value={username} required />
        </div>
        <div>
            <label for="password">Password:</label>
            <input type="password" id="password" bind:value={password} required />
        </div>
        <button type="submit">Login</button>
    </form>
    {#if message}
        <p class:error={isError} class:success={!isError}>{message}</p>
    {/if}
</div>

<style>
    /* Styles are similar to Register.svelte, can be shared in a global CSS */
    .auth-form {
        max-width: 400px;
        margin: 2em auto;
        padding: 1em;
        border: 1px solid #ccc;
        border-radius: 5px;
    }
    .auth-form div {
        margin-bottom: 1em;
    }
    .auth-form label {
        display: block;
        margin-bottom: 0.25em;
    }
    .auth-form input {
        width: calc(100% - 12px);
        padding: 0.5em;
        border: 1px solid #ddd;
        border-radius: 3px;
    }
    .error {
        color: red;
    }
    .success {
        color: green;
    }
</style>
