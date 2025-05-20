<!-- web/ui/src/components/Register.svelte -->
<script>
    import { isAuthenticated, currentUser } from '../stores/auth.js';
    import { navigate } from 'svelte-routing'; // Assuming you'll use svelte-routing

    let username = '';
    let email = '';
    let password = '';
    let message = '';
    let isError = false;

    async function handleRegister() {
        isError = false;
        message = 'Registering...';
        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, email, password }),
            });
            const data = await response.json();
            if (response.ok) {
                message = data.message + " Redirecting to login...";
                // Optionally auto-login or redirect to login page
                setTimeout(() => navigate("/login"), 2000);
            } else {
                isError = true;
                message = data.error || 'Registration failed';
            }
        } catch (error) {
            isError = true;
            message = 'An error occurred: ' + error.message;
            console.error("Registration error:", error);
        }
    }
</script>

<div class="auth-form">
    <h2>Register</h2>
    <form on:submit|preventDefault={handleRegister}>
        <div>
            <label for="username">Username:</label>
            <input type="text" id="username" bind:value={username} required />
        </div>
        <div>
            <label for="email">Email:</label>
            <input type="email" id="email" bind:value={email} required />
        </div>
        <div>
            <label for="password">Password:</label>
            <input type="password" id="password" bind:value={password} required />
        </div>
        <button type="submit">Register</button>
    </form>
    {#if message}
        <p class:error={isError} class:success={!isError}>{message}</p>
    {/if}
</div>

<style>
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
