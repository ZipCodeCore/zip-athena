// Dashboard JavaScript

// Check if user is logged in
const user = JSON.parse(localStorage.getItem('user'));
if (!user) {
    window.location.href = '/index.html';
}

// Display user info
document.getElementById('userDisplay').textContent = `Welcome, ${user.username}!`;

// Logout function
function logout() {
    localStorage.removeItem('user');
    window.location.href = '/index.html';
}

// Load dashboard data
async function loadDashboard() {
    try {
        // Load progress
        const progressResponse = await fetch(`/progress?user_id=${user.id}`);
        const progressData = await progressResponse.json();

        // Update stats
        document.getElementById('totalPoints').textContent = progressData.total_points || 0;
        document.getElementById('totalQuizzes').textContent = progressData.total_quizzes || 0;
        document.getElementById('totalQuestions').textContent = progressData.total_questions || 0;
        document.getElementById('totalCorrect').textContent = progressData.total_correct || 0;

        // Load quiz attempts
        loadQuizAttempts(progressData.quiz_attempts || []);

        // Load topics
        await loadTopics();
    } catch (error) {
        console.error('Error loading dashboard:', error);
    }
}

// Load topics
async function loadTopics() {
    try {
        const response = await fetch('/topics');
        const topics = await response.json();

        const topicsList = document.getElementById('topicsList');
        topicsList.innerHTML = '';

        if (topics && topics.length > 0) {
            topics.forEach(topic => {
                const topicCard = document.createElement('div');
                topicCard.className = 'col-md-4';
                topicCard.innerHTML = `
                    <div class="card topic-card h-100" onclick="startQuiz(${topic.id}, '${topic.name}')">
                        <div class="card-body">
                            <h5 class="card-title">${topic.name}</h5>
                            <p class="card-text">${topic.description || 'Click to start a quiz'}</p>
                            <button class="btn btn-primary btn-sm">Start Quiz</button>
                        </div>
                    </div>
                `;
                topicsList.appendChild(topicCard);
            });
        } else {
            topicsList.innerHTML = '<div class="col-12"><p class="text-muted">No topics available yet.</p></div>';
        }
    } catch (error) {
        console.error('Error loading topics:', error);
    }
}

// Load quiz attempts
function loadQuizAttempts(quizAttempts) {
    const table = document.getElementById('quizAttemptsTable');
    table.innerHTML = '';

    if (quizAttempts && quizAttempts.length > 0) {
        // Sort by most recent first
        quizAttempts.sort((a, b) => new Date(b.start_time) - new Date(a.start_time));

        quizAttempts.forEach(attempt => {
            const row = document.createElement('tr');
            const score = attempt.total_questions > 0 
                ? Math.round((attempt.correct_answers / attempt.total_questions) * 100) 
                : 0;
            
            const startTime = new Date(attempt.start_time).toLocaleString();
            const completedTime = attempt.completed_time 
                ? new Date(attempt.completed_time).toLocaleString() 
                : 'In Progress';

            row.innerHTML = `
                <td>${attempt.id}</td>
                <td>${startTime}</td>
                <td>${completedTime}</td>
                <td>${attempt.total_questions}</td>
                <td>${attempt.correct_answers}</td>
                <td>
                    <span class="badge ${score >= 80 ? 'bg-success' : score >= 60 ? 'bg-warning' : 'bg-danger'}">
                        ${score}%
                    </span>
                </td>
            `;
            table.appendChild(row);
        });
    } else {
        table.innerHTML = '<tr><td colspan="6" class="text-center text-muted">No quiz attempts yet. Start a quiz to see your progress!</td></tr>';
    }
}

// Start a quiz
function startQuiz(topicId, topicName) {
    window.location.href = `/quiz.html?topic_id=${topicId}&topic_name=${encodeURIComponent(topicName)}`;
}

// Load dashboard on page load
loadDashboard();
