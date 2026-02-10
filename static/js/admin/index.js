// Admin Panel Main Entry Point
import * as Users from './users.js';
import * as Questions from './questions.js';
import * as Topics from './topics.js';
import * as Quizzes from './quizzes.js';
import * as Answers from './answers.js';

let currentUser = null;

// Initialize the admin panel
window.addEventListener('DOMContentLoaded', async () => {
    const userData = localStorage.getItem('user');
    if (!userData) {
        console.error('No user data in localStorage');
        window.location.href = '/index.html';
        return;
    }

    try {
        currentUser = JSON.parse(userData);
        console.log('Current user:', {id: currentUser.id, username: currentUser.username, is_admin: currentUser.is_admin});
    } catch (e) {
        console.error('Failed to parse user data:', e);
        localStorage.removeItem('user');
        window.location.href = '/index.html';
        return;
    }
    
    if (!currentUser.id) {
        console.error('User ID is missing. Please log out and log back in.');
        alert('Session error: User ID missing. Please log in again.');
        localStorage.removeItem('user');
        window.location.href = '/index.html';
        return;
    }
    
    if (!currentUser.is_admin) {
        console.warn('User is not an admin');
        window.location.href = '/403.html';
        return;
    }

    // Set current user in modules
    Users.setCurrentUser(currentUser);
    Questions.setCurrentUser(currentUser);
    Topics.setCurrentUser(currentUser);
    Quizzes.setCurrentUser(currentUser);
    Answers.setCurrentUser(currentUser);

    document.getElementById('username').textContent = currentUser.username;

    // Load initial data
    console.log('Loading admin panel data with user ID:', currentUser.id);
    await Users.loadUsers();
    await Questions.loadQuestions();
    await Quizzes.loadQuizzes();
    await Topics.loadTopics();

    // Set up event delegation for button clicks
    setupEventDelegation();
});

// Logout function
window.logout = function() {
    localStorage.removeItem('user');
    window.location.href = '/index.html';
};

// Setup event delegation for all buttons
function setupEventDelegation() {
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
                case 'edit-question':
                    Questions.editQuestion(numId);
                    break;
                case 'delete-question':
                    Questions.deleteQuestion(numId);
                    break;
                case 'edit-topic':
                    Topics.editTopic(numId);
                    break;
                case 'delete-topic':
                    Topics.deleteTopic(numId);
                    break;
                case 'edit-quiz':
                    Quizzes.editQuiz(numId);
                    break;
                case 'delete-quiz':
                    Quizzes.deleteQuiz(numId);
                    break;
            }
        }
    });
}

// Expose functions to global scope for modal buttons
window.showCreateUserModal = Users.showCreateUserModal;
window.saveUser = Users.saveUser;

window.showCreateQuestionModal = Questions.showCreateQuestionModal;
window.saveQuestion = Questions.saveQuestion;

window.showCreateTopicModal = Topics.showCreateTopicModal;
window.saveTopic = Topics.saveTopic;

window.showCreateQuizModal = Quizzes.showCreateQuizModal;
window.saveQuiz = Quizzes.saveQuiz;
window.showAddQuestionToQuizForm = Quizzes.showAddQuestionToQuizForm;
window.cancelAddQuestionToQuiz = Quizzes.cancelAddQuestionToQuiz;
window.addQuestionToQuiz = Quizzes.addQuestionToQuiz;

window.showAddAnswerForm = Answers.showAddAnswerForm;
window.cancelAddAnswer = Answers.cancelAddAnswer;
window.addAnswer = Answers.addAnswer;
