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
    document.getElementById('instructorButton').style.display = 'inline-block';
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
        document.getElementById('quizzesLast7Days').textContent = progressData.quizzes_last_7_days || 0;
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
        const response = await fetch(`/topics?user_id=${user.id}`);
        const topics = await response.json();

        const topicsList = document.getElementById('topicsList');
        topicsList.innerHTML = '';

        if (topics && topics.length > 0) {
            topics.forEach(topic => {
                const topicCard = document.createElement('div');
                topicCard.className = 'col-md-4';

                // Determine badge color and code tracing flag based on category and name.
                // Code tracing topics are identified by their name since all topics now use
                // the three-value category scheme (java/python/all) for access control.
                const isCodeTracing = topic.name.toLowerCase().includes('code tracing');

                let badgeColor = 'secondary';
                let badgeText = topic.category ? topic.category.toUpperCase() : 'GENERAL';
                if (topic.category === 'java') {
                    badgeColor = 'primary';
                    badgeText = 'JAVA';
                } else if (topic.category === 'python') {
                    badgeColor = 'success';
                    badgeText = 'PYTHON';
                } else if (topic.category === 'all') {
                    badgeColor = 'warning';
                    badgeText = 'ALL';
                }

                const cardClass = isCodeTracing ? 'card topic-card code-tracing-card h-100' : 'card topic-card h-100';

                topicCard.innerHTML = `
                    <div class="${cardClass}" onclick="startQuiz(${topic.id}, '${topic.name}')">
                        <div class="card-body">
                            <div class="d-flex justify-content-between align-items-start mb-2">
                                <h5 class="card-title">${topic.name}</h5>
                                <span class="badge bg-${badgeColor}">${badgeText}</span>
                            </div>
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
            
            // Format correct answers as "6/10"
            const correctDisplay = `${attempt.correct_answers}/${attempt.total_questions}`;
            
            // Format difficulty average (round to 1 decimal place)
            const difficultyDisplay = attempt.difficulty_average 
                ? attempt.difficulty_average.toFixed(1) 
                : 'N/A';
            
            // Determine action button based on status and user role
            let actionButton = '';
            // Only show action buttons if user is NOT an admin
            if (!user.is_admin) {
                if (!isCompleted) {
                    // In progress - show Resume button
                    actionButton = `<button class="btn btn-sm btn-success" onclick="resumeQuiz(${attempt.id})">Resume</button>`;
                } else if (score < 100) {
                    // Completed but not perfect - show Retake button
                    actionButton = `<button class="btn btn-sm btn-primary" onclick="retakeQuiz(${attempt.id})">Retake</button>`;
                }
            }

            row.innerHTML = `
                <td data-label="Quiz ID">${attempt.id}</td>
                <td data-label="Quiz Description">${quizDescription}</td>
                <td data-label="Username">${attempt.username || 'N/A'}</td>
                <td data-label="Started">${startTime}</td>
                <td data-label="Status"><span class="badge ${statusBadgeClass}">${status}</span></td>
                <td data-label="Correct">${correctDisplay}</td>
                <td data-label="Difficulty">${difficultyDisplay}</td>
                <td data-label="Score">
                    <span class="badge ${score >= 80 ? 'bg-success' : score >= 60 ? 'bg-warning' : 'bg-danger'}">
                        ${score}%
                    </span>
                </td>
                <td data-label="Action">${actionButton}</td>
            `;
            table.appendChild(row);
        });
    } else {
        table.innerHTML = '<tr><td colspan="9" class="text-center text-muted">No quiz attempts yet. Start a quiz to see your progress!</td></tr>';
    }
}

// Resume an in-progress quiz
async function resumeQuiz(attemptId) {
    if (!attemptId) {
        alert('Unable to resume this quiz. Quiz attempt information not available.');
        return;
    }
    
    try {
        // Get the quiz attempt details with questions
        const response = await fetch(`/students/quiz-history/${attemptId}`);
        if (!response.ok) {
            throw new Error('Failed to load quiz attempt details');
        }
        
        const attemptDetails = await response.json();
        
        if (attemptDetails.completed_time) {
            alert('This quiz is already completed. Use Retake to try again.');
            return;
        }
        
        // Store resume data in localStorage
        localStorage.setItem('resumeQuiz', JSON.stringify({
            attemptId: attemptDetails.id,
            questionResults: attemptDetails.question_results || {},
            questions: attemptDetails.questions || [],
            totalQuestions: attemptDetails.total_questions,
            correctAnswers: attemptDetails.correct_answers
        }));
        
        // Navigate to the quiz page
        window.location.href = '/quiz.html';
    } catch (error) {
        console.error('Error resuming quiz:', error);
        alert('Unable to resume this quiz. Please try again later.');
    }
}

// Retake a quiz
async function retakeQuiz(attemptId) {
    if (!attemptId) {
        alert('Unable to retake this quiz. Quiz attempt information not available.');
        return;
    }
    
    try {
        // Call the retake API endpoint
        const response = await fetch('/quiz/retake', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                parent_attempt_id: attemptId,
                user_id: user.id
            })
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to create retake');
        }
        
        const retakeData = await response.json();
        
        // Store the retake quiz data
        localStorage.setItem('currentQuiz', JSON.stringify(retakeData));
        
        // Navigate to the quiz page
        window.location.href = '/quiz.html';
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

// Load assigned quizzes on startup to check for pending quizzes
loadAssignedQuizzesBadge();

// Load assigned exams on startup to check for pending exams
loadAssignedExamsBadge();

// Show the assigned quizzes panel
function showAssignedQuizzes() {
    document.getElementById('assignedQuizzesSection').style.display = 'block';
    loadAssignedQuizzesList();
}

// Hide the assigned quizzes panel
function hideAssignedQuizzes() {
    document.getElementById('assignedQuizzesSection').style.display = 'none';
}

// Load assigned quizzes badge (count of pending/incomplete quizzes)
async function loadAssignedQuizzesBadge() {
    if (user.is_admin) return; // admins don't take assigned quizzes

    try {
        const response = await fetch(`/assigned-quizzes?user_id=${user.id}`);
        if (!response.ok) return;

        const quizzes = await response.json();
        if (!quizzes || quizzes.length === 0) return;

        document.getElementById('assignedQuizzesBtn').style.display = 'inline-block';

        const pending = quizzes.filter(q => !q.completed_at).length;
        const badge = document.getElementById('assignedQuizzesBadge');
        if (pending > 0) {
            badge.textContent = pending;
            badge.style.display = 'inline-block';
        }
    } catch (error) {
        console.error('Error loading assigned quizzes badge:', error);
    }
}

// Load and display the list of assigned quizzes
async function loadAssignedQuizzesList() {
    const container = document.getElementById('assignedQuizzesList');
    container.innerHTML = '<div class="text-center text-muted">Loading...</div>';

    try {
        const response = await fetch(`/assigned-quizzes?user_id=${user.id}`);
        if (!response.ok) {
            container.innerHTML = '<div class="alert alert-danger">Failed to load assigned quizzes.</div>';
            return;
        }

        const quizzes = await response.json();
        if (!quizzes || quizzes.length === 0) {
            container.innerHTML = '<div class="text-center text-muted">No quizzes have been assigned to you yet.</div>';
            return;
        }

        let html = '<div class="list-group">';
        quizzes.forEach(q => {
            const categoryBadge = q.quiz_category
                ? `<span class="badge bg-secondary ms-1">${q.quiz_category.toUpperCase()}</span>`
                : '';

            let statusHtml = '';
            let actionBtn = '';

            if (q.completed_at) {
                const scoreVal = q.score !== undefined && q.score !== null ? q.score.toFixed(0) : '?';
                const scoreClass = q.score >= 80 ? 'bg-success' : q.score >= 60 ? 'bg-warning' : 'bg-danger';
                statusHtml = `<span class="badge bg-success">Completed</span>
                    <span class="badge ${scoreClass} ms-1">${scoreVal}%</span>
                    <small class="text-muted ms-2">(${q.correct_answers}/${q.total_questions} correct)</small>`;
                actionBtn = `<button class="btn btn-sm btn-outline-primary" onclick="startAssignedQuiz(${q.quiz_id})">Retake</button>`;
            } else {
                statusHtml = '<span class="badge bg-warning text-dark">Not Started</span>';
                actionBtn = `<button class="btn btn-sm btn-warning" onclick="startAssignedQuiz(${q.quiz_id})">Start Quiz</button>`;
            }

            html += `
                <div class="list-group-item">
                    <div class="d-flex justify-content-between align-items-center">
                        <div>
                            <strong>${q.quiz_name}</strong>${categoryBadge}
                            <br>
                            <small class="text-muted">${q.quiz_description || ''}</small>
                            <br>
                            ${statusHtml}
                        </div>
                        <div>
                            ${actionBtn}
                        </div>
                    </div>
                </div>
            `;
        });
        html += '</div>';
        container.innerHTML = html;
    } catch (error) {
        console.error('Error loading assigned quizzes list:', error);
        container.innerHTML = '<div class="alert alert-danger">Error loading assigned quizzes.</div>';
    }
}

// Start an assigned quiz
async function startAssignedQuiz(quizId) {
    try {
        const response = await fetch(`/assigned-quizzes/start?quiz_id=${quizId}&user_id=${user.id}`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'}
        });

        if (!response.ok) {
            const errorText = await response.text();
            alert('Failed to start quiz: ' + errorText);
            return;
        }

        const quizData = await response.json();
        localStorage.setItem('currentQuiz', JSON.stringify(quizData));
        window.location.href = '/quiz.html';
    } catch (error) {
        console.error('Error starting assigned quiz:', error);
        alert('Error starting quiz. Please try again.');
    }
}

// ============ Assigned Exams Functions ============

// Load assigned exams badge (count of pending/incomplete exams)
async function loadAssignedExamsBadge() {
    if (user.is_admin) return; // admins don't take assigned exams

    try {
        const response = await fetch(`/assigned-exams?user_id=${user.id}`);
        if (!response.ok) return;

        const exams = await response.json();
        if (!exams || exams.length === 0) return;

        // Show Exam Center link whenever there are any exams
        document.getElementById('examCenterBtn').style.display = 'inline-block';
        // Hide old inline button (superseded by Exam Center page)
        document.getElementById('assignedExamsBtn').style.display = 'none';

        const pending = exams.filter(e => !e.completed_at).length;
        const badge = document.getElementById('examCenterBadge');
        if (pending > 0) {
            badge.textContent = pending;
            badge.style.display = 'inline-block';
        }

        // Keep old badge in sync for any leftover code references
        const oldBadge = document.getElementById('assignedExamsBadge');
        if (pending > 0 && oldBadge) {
            oldBadge.textContent = pending;
            oldBadge.style.display = 'inline-block';
        }
    } catch (error) {
        console.error('Error loading assigned exams badge:', error);
    }
}

// Show the assigned exams panel
function showAssignedExams() {
    document.getElementById('assignedExamsSection').style.display = 'block';
    loadAssignedExamsList();
}

// Hide the assigned exams panel
function hideAssignedExams() {
    document.getElementById('assignedExamsSection').style.display = 'none';
}

// Load and display the list of assigned exams
async function loadAssignedExamsList() {
    const container = document.getElementById('assignedExamsList');
    container.innerHTML = '<div class="text-center text-muted">Loading...</div>';

    try {
        const response = await fetch(`/assigned-exams?user_id=${user.id}`);
        if (!response.ok) {
            container.innerHTML = '<div class="alert alert-danger">Failed to load assigned exams.</div>';
            return;
        }

        const exams = await response.json();
        if (!exams || exams.length === 0) {
            container.innerHTML = '<div class="text-center text-muted">No exams have been assigned to you yet.</div>';
            return;
        }

        let html = '<div class="list-group">';
        exams.forEach(e => {
            let statusHtml = '';
            let actionBtn = '';

            if (e.completed_at) {
                const scoreVal = e.score !== undefined && e.score !== null ? e.score.toFixed(0) : '?';
                const scoreNum = e.score !== undefined && e.score !== null ? e.score : 0;
                const scoreClass = scoreNum >= 80 ? 'bg-success' : scoreNum >= 60 ? 'bg-warning' : 'bg-danger';
                statusHtml = `<span class="badge bg-success">Completed</span>
                    <span class="badge ${scoreClass} ms-1">${scoreVal}%</span>
                    <small class="text-muted ms-2">(${e.correct_answers}/${e.total_questions} correct)</small>`;
                actionBtn = `<button class="btn btn-sm btn-outline-primary" onclick="startAssignedExam(${e.exam_id}, '${e.access_token}')">Retake</button>`;
            } else {
                statusHtml = '<span class="badge bg-warning text-dark">Not Started</span>';
                actionBtn = `<button class="btn btn-sm btn-danger" onclick="startAssignedExam(${e.exam_id}, '${e.access_token}')">Start Exam</button>`;
            }

            html += `
                <div class="list-group-item">
                    <div class="d-flex justify-content-between align-items-center">
                        <div>
                            <strong>${e.exam_name}</strong>
                            <span class="badge bg-danger ms-1">${e.question_count} questions</span>
                            <br>
                            ${statusHtml}
                        </div>
                        <div>
                            ${actionBtn}
                        </div>
                    </div>
                </div>
            `;
        });
        html += '</div>';
        container.innerHTML = html;
    } catch (error) {
        console.error('Error loading assigned exams list:', error);
        container.innerHTML = '<div class="alert alert-danger">Error loading assigned exams.</div>';
    }
}

// Start an assigned exam
async function startAssignedExam(examId, token) {
    try {
        const response = await fetch(`/exam/start?token=${encodeURIComponent(token)}&user_id=${user.id}`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'}
        });

        if (!response.ok) {
            const errorText = await response.text();
            alert('Failed to start exam: ' + errorText);
            return;
        }

        const examData = await response.json();
        localStorage.setItem('currentExam', JSON.stringify(examData));
        window.location.href = '/exam.html?token=' + encodeURIComponent(token);
    } catch (error) {
        console.error('Error starting assigned exam:', error);
        alert('Error starting exam. Please try again.');
    }
}
