// Dashboard JavaScript

// Check if user is logged in
const user = JSON.parse(localStorage.getItem('user'));
if (!user) {
    window.location.href = '/index.html';
}

// Display user info
document.getElementById('userDisplay').textContent = `Welcome, ${user.username}!`;

// Show admin button if user is admin
if (user.is_admin) {
    document.getElementById('adminButton').style.display = 'inline-block';
}

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
            const isCompleted = attempt.completed_time !== null && attempt.completed_time !== undefined;
            const status = isCompleted ? 'Completed' : 'In Progress';
            const statusBadgeClass = isCompleted ? 'bg-success' : 'bg-warning';
            
            // Get quiz description (fallback to Quiz ID if not available)
            const quizDescription = attempt.quiz_description || attempt.quiz_name || `Quiz #${attempt.id}`;
            
            // Determine if retake button should be shown (completed quiz with score < 100%)
            const showRetake = isCompleted && score < 100;
            const retakeButton = showRetake 
                ? `<button class="btn btn-sm btn-primary" onclick="retakeQuiz(${attempt.quiz_id})">Retake</button>`
                : '';

            row.innerHTML = `
                <td>${attempt.id}</td>
                <td>${quizDescription}</td>
                <td>${startTime}</td>
                <td><span class="badge ${statusBadgeClass}">${status}</span></td>
                <td>${attempt.total_questions}</td>
                <td>${attempt.correct_answers}</td>
                <td>
                    <span class="badge ${score >= 80 ? 'bg-success' : score >= 60 ? 'bg-warning' : 'bg-danger'}">
                        ${score}%
                    </span>
                </td>
                <td>${retakeButton}</td>
            `;
            table.appendChild(row);
        });
    } else {
        table.innerHTML = '<tr><td colspan="8" class="text-center text-muted">No quiz attempts yet. Start a quiz to see your progress!</td></tr>';
    }
}

// Retake a quiz
async function retakeQuiz(quizId) {
    if (!quizId) {
        alert('Unable to retake this quiz. Quiz information not available.');
        return;
    }
    
    try {
        // Fetch the quiz details to get the topic_id
        const response = await fetch(`/quizzes/${quizId}`);
        if (!response.ok) {
            throw new Error('Failed to fetch quiz details');
        }
        
        const quiz = await response.json();
        
        // Fetch the topic details to get the topic name
        const topicResponse = await fetch(`/topics/${quiz.topic_id}`);
        if (!topicResponse.ok) {
            throw new Error('Failed to fetch topic details');
        }
        
        const topic = await topicResponse.json();
        
        // Navigate to the quiz page with the topic information
        window.location.href = `/quiz.html?topic_id=${quiz.topic_id}&topic_name=${encodeURIComponent(topic.name)}`;
    } catch (error) {
        console.error('Error retaking quiz:', error);
        alert('Unable to retake this quiz. Please try again later.');
    }
}

// Start a quiz
function startQuiz(topicId, topicName) {
    window.location.href = `/quiz.html?topic_id=${topicId}&topic_name=${encodeURIComponent(topicName)}`;
}

// Load dashboard on page load
loadDashboard();
