// Admin Panel – Users Page Entry Point
import * as Users from './users.js';
import { requireAdminUser } from './auth.js';

// Logout function
window.logout = function() {
    localStorage.removeItem('user');
    window.location.href = '/index.html';
};

// Initialize the admin panel
window.addEventListener('DOMContentLoaded', async () => {
    const currentUser = await requireAdminUser();
    if (!currentUser) return;

    Users.setCurrentUser(currentUser);

    document.getElementById('username').textContent = currentUser.username;

    await Users.loadUsers();

    // Set up event delegation for button clicks
    document.addEventListener('click', (e) => {
        const action = e.target.getAttribute('data-action');
        const id = e.target.getAttribute('data-id');
        
        if (action && id) {
            const numId = parseInt(id);
            switch(action) {
                case 'edit-user':
                    Users.editUser(numId);
                    break;
                case 'delete-user':
                    Users.deleteUser(numId);
                    break;
            }
        }
    });
});

// Expose functions to global scope for modal buttons
window.showCreateUserModal = Users.showCreateUserModal;
window.saveUser = Users.saveUser;
