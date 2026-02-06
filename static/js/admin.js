// Admin panel JavaScript

// Constants
const DEFAULT_CATEGORY = 'java';

let currentUser = null;

// Check if user is logged in and is admin
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

    document.getElementById('username').textContent = currentUser.username;

    // Load initial data
    console.log('Loading admin panel data with user ID:', currentUser.id);
    await loadUsers();
    await loadQuestions();
    await loadQuizzes();
    await loadTopics();
});

function logout() {
    localStorage.removeItem('user');
    window.location.href = '/index.html';
}

// User Management
async function loadUsers() {
    try {
        console.log('Loading users with X-User-ID:', currentUser.id.toString());
        const response = await fetch('/admin/users', {
            headers: {
                'X-User-ID': currentUser.id.toString()
            }
        });

        console.log('Users response status:', response.status);
        if (!response.ok) throw new Error('Failed to load users');

        const users = await response.json();
        console.log('Loaded users:', users.length, 'users');
        const tbody = document.getElementById('usersTable');
        tbody.innerHTML = '';

        users.forEach(user => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${user.id}</td>
                <td>${user.username}</td>
                <td>${user.email}</td>
                <td>${user.category || DEFAULT_CATEGORY}</td>
                <td>${user.is_admin ? 'Yes' : 'No'}</td>
                <td>${user.total_points}</td>
                <td>${new Date(user.created_at).toLocaleDateString()}</td>
                <td>
                    <button class="btn btn-sm btn-warning" onclick="editUser(${user.id})">Edit</button>
                    <button class="btn btn-sm btn-danger" onclick="deleteUser(${user.id})" 
                            ${user.id === currentUser.id ? 'disabled' : ''}>Delete</button>
                </td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        console.error('Error loading users:', error);
        alert('Failed to load users');
    }
}

async function deleteUser(userId) {
    if (!confirm('Are you sure you want to delete this user?')) return;

    try {
        const response = await fetch(`/admin/users?id=${userId}`, {
            method: 'DELETE',
            headers: {
                'X-User-ID': currentUser.id.toString()
            }
        });

        if (response.ok) {
            await loadUsers();
        } else {
            throw new Error('Failed to delete user');
        }
    } catch (error) {
        console.error('Error deleting user:', error);
        alert('Failed to delete user');
    }
}

async function editUser(userId) {
    const users = await (await fetch('/admin/users', {
        headers: { 'X-User-ID': currentUser.id.toString() }
    })).json();
    
    const user = users.find(u => u.id === userId);
    if (!user) return;

    document.getElementById('userModalTitle').textContent = 'Edit User';
    document.getElementById('userId').value = user.id;
    document.getElementById('userUsername').value = user.username;
    document.getElementById('userEmail').value = user.email;
    document.getElementById('userPassword').value = '';
    document.getElementById('userPassword').required = false;
    document.getElementById('userTotalPoints').value = user.total_points;
    document.getElementById('userCategory').value = user.category || DEFAULT_CATEGORY;
    document.getElementById('userIsAdmin').checked = user.is_admin;
    new bootstrap.Modal(document.getElementById('userModal')).show();
}

function showCreateUserModal() {
    document.getElementById('userModalTitle').textContent = 'Create User';
    document.getElementById('userId').value = '';
    document.getElementById('userUsername').value = '';
    document.getElementById('userEmail').value = '';
    document.getElementById('userPassword').value = '';
    document.getElementById('userPassword').required = true;
    document.getElementById('userTotalPoints').value = '0';
    document.getElementById('userCategory').value = DEFAULT_CATEGORY;
    document.getElementById('userIsAdmin').checked = false;
    new bootstrap.Modal(document.getElementById('userModal')).show();
}

async function saveUser() {
    const userId = document.getElementById('userId').value;
    const totalPoints = parseInt(document.getElementById('userTotalPoints').value);
    
    // Validate total points
    if (isNaN(totalPoints) || totalPoints < 0) {
        alert('Total Points must be a valid non-negative number');
        return;
    }
    
    const userData = {
        username: document.getElementById('userUsername').value,
        email: document.getElementById('userEmail').value,
        total_points: totalPoints,
        category: document.getElementById('userCategory').value,
        is_admin: document.getElementById('userIsAdmin').checked
    };

    // Add password if creating or if it's been changed
    const password = document.getElementById('userPassword').value;
    if (password) {
        userData.password = password;
    }

    try {
        let response;
        if (userId) {
            // Update existing user
            response = await fetch(`/admin/users?id=${userId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'X-User-ID': currentUser.id.toString()
                },
                body: JSON.stringify(userData)
            });
        } else {
            // Create new user
            if (!password) {
                alert('Password is required for new users');
                return;
            }
            response = await fetch('/admin/users', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-User-ID': currentUser.id.toString()
                },
                body: JSON.stringify(userData)
            });
        }

        if (response.ok) {
            bootstrap.Modal.getInstance(document.getElementById('userModal')).hide();
            await loadUsers();
        } else {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to save user');
        }
    } catch (error) {
        console.error('Error saving user:', error);
        alert('Failed to save user: ' + error.message);
    }
}

// Question Management
async function loadQuestions() {
    try {
        const response = await fetch('/admin/questions', {
            headers: {
                'X-User-ID': currentUser.id.toString()
            }
        });

        if (!response.ok) throw new Error('Failed to load questions');

        const questions = await response.json();
        const tbody = document.getElementById('questionsTable');
        tbody.innerHTML = '';

        questions.forEach(q => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${q.id}</td>
                <td>${q.question_text.substring(0, 50)}...</td>
                <td>${q.topic_id}</td>
                <td>${q.difficulty}</td>
                <td>${q.tags ? q.tags.join(', ') : ''}</td>
                <td>
                    <button class="btn btn-sm btn-warning" onclick="editQuestion(${q.id})">Edit</button>
                    <button class="btn btn-sm btn-danger" onclick="deleteQuestion(${q.id})">Delete</button>
                </td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        console.error('Error loading questions:', error);
        alert('Failed to load questions');
    }
}

async function deleteQuestion(questionId) {
    if (!confirm('Are you sure you want to delete this question?')) return;

    try {
        const response = await fetch(`/admin/questions?id=${questionId}`, {
            method: 'DELETE',
            headers: {
                'X-User-ID': currentUser.id.toString()
            }
        });

        if (response.ok) {
            await loadQuestions();
        } else {
            throw new Error('Failed to delete question');
        }
    } catch (error) {
        console.error('Error deleting question:', error);
        alert('Failed to delete question');
    }
}

async function editQuestion(questionId) {
    const questions = await (await fetch('/admin/questions', {
        headers: { 'X-User-ID': currentUser.id.toString() }
    })).json();
    
    const question = questions.find(q => q.id === questionId);
    if (!question) return;

    document.getElementById('questionModalTitle').textContent = 'Edit Question';
    document.getElementById('questionId').value = question.id;
    document.getElementById('questionText').value = question.question_text;
    document.getElementById('questionTopicId').value = question.topic_id;
    document.getElementById('questionDifficulty').value = question.difficulty;
    document.getElementById('questionTags').value = question.tags ? question.tags.join(', ') : '';
    document.getElementById('questionTopicQuestionId').value = question.topic_question_id;
    
    // Show answers section and load answers
    document.getElementById('answersSection').style.display = 'block';
    await loadAnswers(question.id);
    
    new bootstrap.Modal(document.getElementById('questionModal')).show();
}

function showCreateQuestionModal() {
    document.getElementById('questionModalTitle').textContent = 'Create Question';
    document.getElementById('questionId').value = '';
    document.getElementById('questionText').value = '';
    document.getElementById('questionTopicId').value = '';
    document.getElementById('questionDifficulty').value = '1';
    document.getElementById('questionTags').value = '';
    document.getElementById('questionTopicQuestionId').value = '';
    
    // Hide answers section for new questions
    document.getElementById('answersSection').style.display = 'none';
    
    new bootstrap.Modal(document.getElementById('questionModal')).show();
}

async function saveQuestion() {
    const questionId = document.getElementById('questionId').value;
    const tagsInput = document.getElementById('questionTags').value;
    const tags = tagsInput ? tagsInput.split(',').map(t => t.trim()).filter(t => t) : [];
    
    const topicId = parseInt(document.getElementById('questionTopicId').value);
    const difficulty = parseInt(document.getElementById('questionDifficulty').value);
    
    // Validate inputs
    if (isNaN(topicId) || topicId < 1) {
        alert('Topic ID must be a valid positive number');
        return;
    }
    if (isNaN(difficulty) || difficulty < 1 || difficulty > 7) {
        alert('Difficulty must be a number between 1 and 7');
        return;
    }
    
    const questionData = {
        question_text: document.getElementById('questionText').value,
        topic_id: topicId,
        difficulty: difficulty,
        tags: tags,
        topic_question_id: document.getElementById('questionTopicQuestionId').value
    };

    try {
        let response;
        if (questionId) {
            // Update existing question
            response = await fetch(`/admin/questions?id=${questionId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'X-User-ID': currentUser.id.toString()
                },
                body: JSON.stringify(questionData)
            });
        } else {
            // Create new question
            response = await fetch('/admin/questions', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-User-ID': currentUser.id.toString()
                },
                body: JSON.stringify(questionData)
            });
        }

        if (response.ok) {
            bootstrap.Modal.getInstance(document.getElementById('questionModal')).hide();
            await loadQuestions();
        } else {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to save question');
        }
    } catch (error) {
        console.error('Error saving question:', error);
        alert('Failed to save question: ' + error.message);
    }
}

// Quiz Management
async function loadQuizzes() {
    try {
        const response = await fetch('/admin/quizzes', {
            headers: {
                'X-User-ID': currentUser.id.toString()
            }
        });

        if (!response.ok) throw new Error('Failed to load quizzes');

        const quizzes = await response.json();
        const tbody = document.getElementById('quizzesTable');
        tbody.innerHTML = '';

        quizzes.forEach(quiz => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${quiz.id}</td>
                <td>${quiz.name}</td>
                <td>${quiz.topic_id}</td>
                <td>${quiz.description}</td>
                <td>${new Date(quiz.created_at).toLocaleDateString()}</td>
                <td>
                    <button class="btn btn-sm btn-warning" onclick="editQuiz(${quiz.id})">Edit</button>
                    <button class="btn btn-sm btn-danger" onclick="deleteQuiz(${quiz.id})">Delete</button>
                </td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        console.error('Error loading quizzes:', error);
        alert('Failed to load quizzes');
    }
}

function showCreateQuizModal() {
    document.getElementById('quizModalTitle').textContent = 'Create Quiz';
    document.getElementById('quizId').value = '';
    document.getElementById('quizName').value = '';
    document.getElementById('quizTopicId').value = '';
    document.getElementById('quizDescription').value = '';
    
    // Hide quiz questions section for new quizzes
    document.getElementById('quizQuestionsSection').style.display = 'none';
    
    new bootstrap.Modal(document.getElementById('quizModal')).show();
}

async function editQuiz(quizId) {
    // Fetch quiz details
    const quizzes = await (await fetch('/admin/quizzes', {
        headers: { 'X-User-ID': currentUser.id.toString() }
    })).json();
    
    const quiz = quizzes.find(q => q.id === quizId);
    if (!quiz) return;

    document.getElementById('quizModalTitle').textContent = 'Edit Quiz';
    document.getElementById('quizId').value = quiz.id;
    document.getElementById('quizName').value = quiz.name;
    document.getElementById('quizTopicId').value = quiz.topic_id;
    document.getElementById('quizDescription').value = quiz.description;
    
    // Show quiz questions section and load questions
    document.getElementById('quizQuestionsSection').style.display = 'block';
    await loadQuizQuestions(quiz.id);
    
    new bootstrap.Modal(document.getElementById('quizModal')).show();
}

async function saveQuiz() {
    const quizId = document.getElementById('quizId').value;
    const quizData = {
        name: document.getElementById('quizName').value,
        topic_id: parseInt(document.getElementById('quizTopicId').value),
        description: document.getElementById('quizDescription').value
    };

    try {
        let response;
        if (quizId) {
            // Update existing quiz
            quizData.id = parseInt(quizId);
            response = await fetch(`/admin/quizzes?id=${quizId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'X-User-ID': currentUser.id.toString()
                },
                body: JSON.stringify(quizData)
            });
        } else {
            // Create new quiz
            response = await fetch('/admin/quizzes', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-User-ID': currentUser.id.toString()
                },
                body: JSON.stringify(quizData)
            });
        }

        if (response.ok) {
            bootstrap.Modal.getInstance(document.getElementById('quizModal')).hide();
            await loadQuizzes();
        } else {
            throw new Error('Failed to save quiz');
        }
    } catch (error) {
        console.error('Error saving quiz:', error);
        alert('Failed to save quiz');
    }
}

async function deleteQuiz(quizId) {
    if (!confirm('Are you sure you want to delete this quiz?')) return;

    try {
        const response = await fetch(`/admin/quizzes?id=${quizId}`, {
            method: 'DELETE',
            headers: {
                'X-User-ID': currentUser.id.toString()
            }
        });

        if (response.ok) {
            await loadQuizzes();
        } else {
            throw new Error('Failed to delete quiz');
        }
    } catch (error) {
        console.error('Error deleting quiz:', error);
        alert('Failed to delete quiz');
    }
}

// Topic Management
async function loadTopics() {
    try {
        if (!currentUser || !currentUser.id) {
            console.error('User not loaded yet');
            return;
        }
        const response = await fetch(`/topics?user_id=${currentUser.id}`);

        if (!response.ok) throw new Error('Failed to load topics');

        const topics = await response.json();
        const tbody = document.getElementById('topicsTable');
        tbody.innerHTML = '';

        topics.forEach(topic => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${topic.id}</td>
                <td>${topic.name}</td>
                <td>${topic.description}</td>
                <td>${new Date(topic.created_at).toLocaleDateString()}</td>
                <td>
                    <button class="btn btn-sm btn-warning" onclick="editTopic(${topic.id})">Edit</button>
                    <button class="btn btn-sm btn-danger" onclick="deleteTopic(${topic.id})">Delete</button>
                </td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        console.error('Error loading topics:', error);
        alert('Failed to load topics');
    }
}

function showCreateTopicModal() {
    document.getElementById('topicModalTitle').textContent = 'Create Topic';
    document.getElementById('topicId').value = '';
    document.getElementById('topicName').value = '';
    document.getElementById('topicDescription').value = '';
    new bootstrap.Modal(document.getElementById('topicModal')).show();
}

async function editTopic(topicId) {
    if (!currentUser || !currentUser.id) {
        console.error('User not loaded yet');
        return;
    }
    const topics = await (await fetch(`/topics?user_id=${currentUser.id}`)).json();
    const topic = topics.find(t => t.id === topicId);
    if (!topic) return;

    document.getElementById('topicModalTitle').textContent = 'Edit Topic';
    document.getElementById('topicId').value = topic.id;
    document.getElementById('topicName').value = topic.name;
    document.getElementById('topicDescription').value = topic.description;
    new bootstrap.Modal(document.getElementById('topicModal')).show();
}

async function saveTopic() {
    const topicId = document.getElementById('topicId').value;
    const topicData = {
        name: document.getElementById('topicName').value,
        description: document.getElementById('topicDescription').value
    };

    try {
        let response;
        if (topicId) {
            // Update existing topic
            topicData.id = parseInt(topicId);
            response = await fetch(`/admin/topics?id=${topicId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'X-User-ID': currentUser.id.toString()
                },
                body: JSON.stringify(topicData)
            });
        } else {
            // Create new topic (use standard endpoint since there's no admin create endpoint)
            response = await fetch('/topics', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(topicData)
            });
        }

        if (response.ok) {
            bootstrap.Modal.getInstance(document.getElementById('topicModal')).hide();
            await loadTopics();
        } else {
            throw new Error('Failed to save topic');
        }
    } catch (error) {
        console.error('Error saving topic:', error);
        alert('Failed to save topic');
    }
}

async function deleteTopic(topicId) {
    if (!confirm('Are you sure you want to delete this topic? This may affect related questions and quizzes.')) return;

    try {
        const response = await fetch(`/admin/topics?id=${topicId}`, {
            method: 'DELETE',
            headers: {
                'X-User-ID': currentUser.id.toString()
            }
        });

        if (response.ok) {
            await loadTopics();
        } else {
            throw new Error('Failed to delete topic');
        }
    } catch (error) {
        console.error('Error deleting topic:', error);
        alert('Failed to delete topic');
    }
}

// Answer Management Functions
async function loadAnswers(questionId) {
    try {
        const response = await fetch(`/admin/answers?question_id=${questionId}`, {
            headers: {
                'X-User-ID': currentUser.id.toString()
            }
        });

        if (!response.ok) throw new Error('Failed to load answers');

        const answers = await response.json();
        const answersList = document.getElementById('answersList');
        answersList.innerHTML = '';

        if (!answers || answers.length === 0) {
            answersList.innerHTML = '<p class="text-muted">No answers yet. Add some answer alternatives.</p>';
            return;
        }

        answers.forEach(answer => {
            const answerDiv = document.createElement('div');
            answerDiv.className = 'border rounded p-2 mb-2';
            answerDiv.id = `answer-${answer.id}`;
            
            const container = document.createElement('div');
            container.className = 'd-flex justify-content-between align-items-start';
            
            const textDiv = document.createElement('div');
            const strong = document.createElement('strong');
            strong.textContent = (answer.is_correct ? '✓ ' : '') + answer.answer_text;
            textDiv.appendChild(strong);
            
            if (answer.is_correct) {
                const badge = document.createElement('span');
                badge.className = 'badge bg-success ms-2';
                badge.textContent = 'Correct';
                textDiv.appendChild(badge);
            }
            
            const buttonGroup = document.createElement('div');
            
            const editBtn = document.createElement('button');
            editBtn.className = 'btn btn-sm btn-warning me-2';
            editBtn.textContent = 'Edit';
            editBtn.onclick = () => showEditAnswerForm(answer);
            
            const deleteBtn = document.createElement('button');
            deleteBtn.className = 'btn btn-sm btn-danger';
            deleteBtn.textContent = 'Delete';
            deleteBtn.onclick = () => deleteAnswer(answer.id);
            
            buttonGroup.appendChild(editBtn);
            buttonGroup.appendChild(deleteBtn);
            
            container.appendChild(textDiv);
            container.appendChild(buttonGroup);
            answerDiv.appendChild(container);
            answersList.appendChild(answerDiv);
        });
    } catch (error) {
        console.error('Error loading answers:', error);
        alert('Failed to load answers');
    }
}

function showAddAnswerForm() {
    document.getElementById('addAnswerForm').style.display = 'block';
    document.getElementById('newAnswerText').value = '';
    document.getElementById('newAnswerIsCorrect').checked = false;
}

function cancelAddAnswer() {
    document.getElementById('addAnswerForm').style.display = 'none';
}

async function addAnswer() {
    const questionId = document.getElementById('questionId').value;
    if (!questionId) {
        alert('Please save the question first before adding answers');
        return;
    }

    const answerData = {
        question_id: parseInt(questionId),
        answer_text: document.getElementById('newAnswerText').value,
        is_correct: document.getElementById('newAnswerIsCorrect').checked
    };

    if (!answerData.answer_text) {
        alert('Answer text is required');
        return;
    }

    try {
        const response = await fetch('/admin/answers', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-User-ID': currentUser.id.toString()
            },
            body: JSON.stringify(answerData)
        });

        if (response.ok) {
            cancelAddAnswer();
            await loadAnswers(questionId);
        } else {
            throw new Error('Failed to add answer');
        }
    } catch (error) {
        console.error('Error adding answer:', error);
        alert('Failed to add answer');
    }
}

async function deleteAnswer(answerId) {
    if (!confirm('Are you sure you want to delete this answer?')) return;

    const questionId = document.getElementById('questionId').value;

    try {
        const response = await fetch(`/admin/answers?id=${answerId}`, {
            method: 'DELETE',
            headers: {
                'X-User-ID': currentUser.id.toString()
            }
        });

        if (response.ok) {
            await loadAnswers(questionId);
        } else {
            throw new Error('Failed to delete answer');
        }
    } catch (error) {
        console.error('Error deleting answer:', error);
        alert('Failed to delete answer');
    }
}

function showEditAnswerForm(answer) {
    // Hide the add answer form if it's open
    document.getElementById('addAnswerForm').style.display = 'none';
    
    // Create or get the edit form
    let editForm = document.getElementById('editAnswerForm');
    if (!editForm) {
        editForm = createEditAnswerForm();
    }
    
    // Populate the form with the answer data
    document.getElementById('editAnswerId').value = answer.id;
    document.getElementById('editAnswerText').value = answer.answer_text;
    document.getElementById('editAnswerIsCorrect').checked = answer.is_correct;
    
    // Show the edit form
    editForm.style.display = 'block';
    
    // Scroll to the form
    editForm.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
}

function createEditAnswerForm() {
    const editForm = document.createElement('div');
    editForm.id = 'editAnswerForm';
    editForm.className = 'mb-3 p-3 border rounded bg-light';
    editForm.style.display = 'none';
    editForm.innerHTML = `
        <h6>Edit Answer</h6>
        <input type="hidden" id="editAnswerId">
        <div class="mb-2">
            <label for="editAnswerText" class="form-label">Answer Text</label>
            <input type="text" class="form-control" id="editAnswerText" required>
        </div>
        <div class="mb-2 form-check">
            <input type="checkbox" class="form-check-input" id="editAnswerIsCorrect">
            <label class="form-check-label" for="editAnswerIsCorrect">Is Correct Answer</label>
        </div>
        <button type="button" class="btn btn-sm btn-primary" onclick="saveEditedAnswer()">Save Changes</button>
        <button type="button" class="btn btn-sm btn-secondary" onclick="cancelEditAnswer()">Cancel</button>
    `;
    
    // Insert the edit form after the add answer form
    const addAnswerForm = document.getElementById('addAnswerForm');
    addAnswerForm.parentNode.insertBefore(editForm, addAnswerForm.nextSibling);
    
    return editForm;
}

function cancelEditAnswer() {
    const editForm = document.getElementById('editAnswerForm');
    if (editForm) {
        editForm.style.display = 'none';
    }
}

async function saveEditedAnswer() {
    const answerId = document.getElementById('editAnswerId').value;
    const questionId = document.getElementById('questionId').value;
    
    // Validate IDs
    if (!answerId || isNaN(parseInt(answerId))) {
        alert('Invalid answer ID');
        return;
    }
    
    if (!questionId || isNaN(parseInt(questionId))) {
        alert('Invalid question ID');
        return;
    }
    
    const answerData = {
        question_id: parseInt(questionId),
        answer_text: document.getElementById('editAnswerText').value,
        is_correct: document.getElementById('editAnswerIsCorrect').checked
    };

    if (!answerData.answer_text) {
        alert('Answer text is required');
        return;
    }

    try {
        const response = await fetch(`/admin/answers?id=${encodeURIComponent(answerId)}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'X-User-ID': currentUser.id.toString()
            },
            body: JSON.stringify(answerData)
        });

        if (response.ok) {
            cancelEditAnswer();
            await loadAnswers(questionId);
        } else {
            throw new Error('Failed to update answer');
        }
    } catch (error) {
        console.error('Error updating answer:', error);
        alert('Failed to update answer');
    }
}

// Quiz Question Management Functions
async function loadQuizQuestions(quizId) {
    try {
        const response = await fetch(`/admin/quiz-questions?quiz_id=${quizId}`, {
            headers: {
                'X-User-ID': currentUser.id.toString()
            }
        });

        if (!response.ok) throw new Error('Failed to load quiz questions');

        const questions = await response.json();
        const questionsList = document.getElementById('quizQuestionsList');
        questionsList.innerHTML = '';

        if (!questions || questions.length === 0) {
            questionsList.innerHTML = '<p class="text-muted">No questions in this quiz yet. Add some questions.</p>';
            return;
        }

        questions.forEach((question, index) => {
            const questionDiv = document.createElement('div');
            questionDiv.className = 'border rounded p-2 mb-2';
            
            const container = document.createElement('div');
            container.className = 'd-flex justify-content-between align-items-start';
            
            const textDiv = document.createElement('div');
            const strong = document.createElement('strong');
            strong.textContent = `#${index + 1}: `;
            textDiv.appendChild(strong);
            
            const questionText = document.createTextNode(question.question_text.substring(0, 60) + '...');
            textDiv.appendChild(questionText);
            
            const br = document.createElement('br');
            textDiv.appendChild(br);
            
            const small = document.createElement('small');
            small.className = 'text-muted';
            small.textContent = `ID: ${question.id} | Topic ID: ${question.topic_id}`;
            textDiv.appendChild(small);
            
            const removeBtn = document.createElement('button');
            removeBtn.className = 'btn btn-sm btn-danger';
            removeBtn.textContent = 'Remove';
            removeBtn.onclick = () => removeQuestionFromQuiz(quizId, question.id);
            
            container.appendChild(textDiv);
            container.appendChild(removeBtn);
            questionDiv.appendChild(container);
            questionsList.appendChild(questionDiv);
        });
    } catch (error) {
        console.error('Error loading quiz questions:', error);
        alert('Failed to load quiz questions');
    }
}

function showAddQuestionToQuizForm() {
    document.getElementById('addQuestionToQuizForm').style.display = 'block';
    document.getElementById('quizQuestionId').value = '';
    
    // Auto-increment order based on current questions
    const questionsList = document.getElementById('quizQuestionsList');
    const currentCount = questionsList.querySelectorAll('.border').length;
    document.getElementById('quizQuestionOrder').value = currentCount;
}

function cancelAddQuestionToQuiz() {
    document.getElementById('addQuestionToQuizForm').style.display = 'none';
}

async function addQuestionToQuiz() {
    const quizId = document.getElementById('quizId').value;
    if (!quizId) {
        alert('Please save the quiz first before adding questions');
        return;
    }

    const questionId = parseInt(document.getElementById('quizQuestionId').value);
    const order = parseInt(document.getElementById('quizQuestionOrder').value);

    if (!questionId || questionId < 1) {
        alert('Please enter a valid question ID');
        return;
    }
    
    if (isNaN(order) || order < 0) {
        alert('Please enter a valid order (non-negative number)');
        return;
    }

    try {
        const response = await fetch('/admin/quiz-questions', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-User-ID': currentUser.id.toString()
            },
            body: JSON.stringify({
                quiz_id: parseInt(quizId),
                question_id: questionId,
                order: order
            })
        });

        if (response.ok) {
            cancelAddQuestionToQuiz();
            await loadQuizQuestions(quizId);
        } else {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to add question to quiz');
        }
    } catch (error) {
        console.error('Error adding question to quiz:', error);
        alert('Failed to add question to quiz: ' + error.message);
    }
}

async function removeQuestionFromQuiz(quizId, questionId) {
    if (!confirm('Are you sure you want to remove this question from the quiz?')) return;

    try {
        const response = await fetch('/admin/quiz-questions', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'X-User-ID': currentUser.id.toString()
            },
            body: JSON.stringify({
                quiz_id: quizId,
                question_id: questionId
            })
        });

        if (response.ok) {
            await loadQuizQuestions(quizId);
        } else {
            throw new Error('Failed to remove question from quiz');
        }
    } catch (error) {
        console.error('Error removing question from quiz:', error);
        alert('Failed to remove question from quiz');
    }
}
