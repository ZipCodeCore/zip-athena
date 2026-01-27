// Admin panel JavaScript

let currentUser = null;

// Check if user is logged in and is admin
window.addEventListener('DOMContentLoaded', async () => {
    const userData = localStorage.getItem('user');
    if (!userData) {
        window.location.href = '/index.html';
        return;
    }

    currentUser = JSON.parse(userData);
    if (!currentUser.is_admin) {
        alert('Access denied. Admin privileges required.');
        window.location.href = '/dashboard.html';
        return;
    }

    document.getElementById('username').textContent = currentUser.username;

    // Load initial data
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
        const response = await fetch('/admin/users', {
            headers: {
                'X-User-ID': currentUser.id.toString()
            }
        });

        if (!response.ok) throw new Error('Failed to load users');

        const users = await response.json();
        const tbody = document.getElementById('usersTable');
        tbody.innerHTML = '';

        users.forEach(user => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${user.id}</td>
                <td>${user.username}</td>
                <td>${user.email}</td>
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
    document.getElementById('userTotalPoints').value = user.total_points;
    document.getElementById('userIsAdmin').checked = user.is_admin;
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
        is_admin: document.getElementById('userIsAdmin').checked
    };

    try {
        const response = await fetch(`/admin/users?id=${userId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'X-User-ID': currentUser.id.toString()
            },
            body: JSON.stringify(userData)
        });

        if (response.ok) {
            bootstrap.Modal.getInstance(document.getElementById('userModal')).hide();
            await loadUsers();
        } else {
            throw new Error('Failed to save user');
        }
    } catch (error) {
        console.error('Error saving user:', error);
        alert('Failed to save user');
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
        const response = await fetch(`/admin/questions?id=${questionId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'X-User-ID': currentUser.id.toString()
            },
            body: JSON.stringify(questionData)
        });

        if (response.ok) {
            bootstrap.Modal.getInstance(document.getElementById('questionModal')).hide();
            await loadQuestions();
        } else {
            throw new Error('Failed to save question');
        }
    } catch (error) {
        console.error('Error saving question:', error);
        alert('Failed to save question');
    }
}

function showCreateQuestionModal() {
    alert('Question creation UI coming soon. Use the API to create questions for now.');
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
        const response = await fetch('/topics');

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
    const topics = await (await fetch('/topics')).json();
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
