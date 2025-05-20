// web/ui/src/stores/auth.js
import { writable } from 'svelte/store';

export const isAuthenticated = writable(false);
export const currentUser = writable(null); // To store user details like username, email

// Function to check initial auth status (e.g., by calling /api/me)
export async function checkAuthStatus() {
    try {
        const response = await fetch('/api/me', {
            credentials: 'include', // Important to send cookies
        });
        if (response.ok) {
            const user = await response.json();
            if (user && user.id) {
                isAuthenticated.set(true);
                currentUser.set(user);
            } else {
                isAuthenticated.set(false);
                currentUser.set(null);
            }
        } else {
            isAuthenticated.set(false);
            currentUser.set(null);
        }
    } catch (error) {
        console.error("Error checking auth status:", error);
        isAuthenticated.set(false);
        currentUser.set(null);
    }
}

// Call it once when the app loads
checkAuthStatus();
