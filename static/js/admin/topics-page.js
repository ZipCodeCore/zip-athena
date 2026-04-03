// Admin Panel – Topics Page Entry Point
import * as Topics from './topics.js';
import { requireAdminUser } from './auth.js';

// Logout function
window.logout = function() {
    localStorage.removeItem('user');
    window.location.href = '/index.html';
};

// Initialize the topics admin page
window.addEventListener('DOMContentLoaded', async () => {
    const currentUser = await requireAdminUser();
    if (!currentUser) return;

    Topics.setCurrentUser(currentUser);

    document.getElementById('username').textContent = currentUser.username;

    await Topics.loadTopics();

    // Set up event delegation for button clicks
    document.addEventListener('click', (e) => {
        const action = e.target.getAttribute('data-action');
        const id = e.target.getAttribute('data-id');

        if (action && id) {
            const numId = parseInt(id);
            switch(action) {
                case 'edit-topic':
                    Topics.editTopic(numId);
                    break;
                case 'delete-topic':
                    Topics.deleteTopic(numId);
                    break;
            }
        }
    });
});

// Expose functions to global scope for modal buttons
window.showCreateTopicModal = Topics.showCreateTopicModal;
window.saveTopic = Topics.saveTopic;
