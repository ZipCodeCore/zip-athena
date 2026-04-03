// Admin Panel – Quizzes Page Entry Point
import * as Quizzes from './quizzes.js';
import { requireAdminUser } from './auth.js';

// Logout function
window.logout = function() {
    localStorage.removeItem('user');
    window.location.href = '/index.html';
};

// Initialize the quizzes admin page
window.addEventListener('DOMContentLoaded', async () => {
    const currentUser = await requireAdminUser();
    if (!currentUser) return;

    Quizzes.setCurrentUser(currentUser);

    document.getElementById('username').textContent = currentUser.username;

    await Quizzes.loadQuizzes();

    // Set up event delegation for button clicks
    document.addEventListener('click', (e) => {
        const action = e.target.getAttribute('data-action');
        const id = e.target.getAttribute('data-id');

        if (action && id) {
            const numId = parseInt(id);
            switch(action) {
                case 'edit-quiz':
                    Quizzes.editQuiz(numId);
                    break;
                case 'delete-quiz':
                    Quizzes.deleteQuiz(numId);
                    break;
            }
        }
    });
});

// Expose functions to global scope for modal buttons
window.showCreateQuizModal = Quizzes.showCreateQuizModal;
window.saveQuiz = Quizzes.saveQuiz;
window.showAddQuestionToQuizForm = Quizzes.showAddQuestionToQuizForm;
window.cancelAddQuestionToQuiz = Quizzes.cancelAddQuestionToQuiz;
window.addQuestionToQuiz = Quizzes.addQuestionToQuiz;
