// Admin Panel – Questions Page Entry Point
import * as Questions from './questions.js';
import * as Answers from './answers.js';
import { requireAdminUser } from './auth.js';

// Logout function
window.logout = function() {
    localStorage.removeItem('user');
    window.location.href = '/index.html';
};

// Initialize the questions admin page
window.addEventListener('DOMContentLoaded', async () => {
    const currentUser = await requireAdminUser();
    if (!currentUser) return;

    Questions.setCurrentUser(currentUser);
    Answers.setCurrentUser(currentUser);

    document.getElementById('username').textContent = currentUser.username;

    await Questions.loadQuestions();

    // Set up event delegation for button clicks
    document.addEventListener('click', (e) => {
        const action = e.target.getAttribute('data-action');
        const id = e.target.getAttribute('data-id');

        if (action && id) {
            const numId = parseInt(id);
            switch(action) {
                case 'edit-question':
                    Questions.editQuestion(numId);
                    break;
                case 'delete-question':
                    Questions.deleteQuestion(numId);
                    break;
            }
        }
    });
});

// Expose functions to global scope for modal and panel buttons
window.showCreateQuestionModal = Questions.showCreateQuestionModal;
window.saveQuestion = Questions.saveQuestion;

window.showAddAnswerForm = Answers.showAddAnswerForm;
window.cancelAddAnswer = Answers.cancelAddAnswer;
window.addAnswer = Answers.addAnswer;
