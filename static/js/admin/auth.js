// Shared admin authentication helper
// Returns the current user object or redirects to login/403 if auth fails.
export async function requireAdminUser() {
    const userData = localStorage.getItem('user');
    if (!userData) {
        window.location.href = '/index.html';
        return null;
    }

    let currentUser;
    try {
        currentUser = JSON.parse(userData);
    } catch (e) {
        localStorage.removeItem('user');
        window.location.href = '/index.html';
        return null;
    }

    if (!currentUser.id) {
        alert('Session error: User ID missing. Please log in again.');
        localStorage.removeItem('user');
        window.location.href = '/index.html';
        return null;
    }

    if (!currentUser.is_admin) {
        window.location.href = '/403.html';
        return null;
    }

    return currentUser;
}

// Shared logout function for all admin pages.
export function setupLogout() {
    window.logout = function() {
        localStorage.removeItem('user');
        window.location.href = '/index.html';
    };
}
