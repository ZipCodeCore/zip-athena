// Exam Center JavaScript

const user = JSON.parse(localStorage.getItem('user'));
if (!user) {
    window.location.href = '/index.html';
}

document.getElementById('userDisplay').textContent = `Welcome, ${user.username}!`;

function logout() {
    localStorage.removeItem('user');
    window.location.href = '/index.html';
}

async function loadExams() {
    try {
        const response = await fetch(`/assigned-exams?user_id=${user.id}`);
        if (!response.ok) {
            throw new Error('Failed to load exams');
        }
        const exams = await response.json();

        const pending = exams.filter(e => !e.completed_at);
        const submitted = exams.filter(e => !!e.completed_at);

        renderPendingExams(pending);
        renderSubmittedExams(submitted);
    } catch (error) {
        console.error('Error loading exams:', error);
        document.getElementById('pendingExamsList').innerHTML =
            '<div class="alert alert-danger">Error loading exams. Please try again.</div>';
        document.getElementById('submittedExamsList').innerHTML =
            '<div class="alert alert-danger">Error loading exams. Please try again.</div>';
    }
}

function renderPendingExams(exams) {
    const container = document.getElementById('pendingExamsList');
    if (!exams || exams.length === 0) {
        container.innerHTML = '<p class="text-muted mb-0">No pending exams. You\'re all caught up!</p>';
        return;
    }

    let html = '<div class="list-group">';
    exams.forEach(e => {
        html += `
            <div class="list-group-item">
                <div class="d-flex justify-content-between align-items-center">
                    <div>
                        <strong>${escapeHtml(e.exam_name)}</strong>
                        <span class="badge bg-danger ms-2">${e.question_count} questions</span>
                        <br>
                        <small class="text-muted">Assigned: ${new Date(e.assigned_at).toLocaleDateString()}</small>
                        <br>
                        <span class="badge bg-warning text-dark">Not Started</span>
                    </div>
                    <div>
                        <button class="btn btn-danger btn-sm"
                            data-exam-id="${e.exam_id}"
                            data-token="${escapeHtml(e.access_token)}"
                            onclick="startExam(this.dataset.examId, this.dataset.token)">
                            <i class="bi bi-play-fill"></i> Start Exam
                        </button>
                    </div>
                </div>
            </div>
        `;
    });
    html += '</div>';
    container.innerHTML = html;
}

function renderSubmittedExams(exams) {
    const container = document.getElementById('submittedExamsList');
    if (!exams || exams.length === 0) {
        container.innerHTML = '<p class="text-muted mb-0">No submitted exams yet.</p>';
        return;
    }

    let html = '<div class="list-group">';
    exams.forEach(e => {
        const scoreVal = e.score !== undefined && e.score !== null ? e.score.toFixed(0) : '?';
        const scoreNum = e.score !== undefined && e.score !== null ? e.score : 0;
        const scoreClass = scoreNum >= 80 ? 'bg-success' : scoreNum >= 60 ? 'bg-warning' : 'bg-danger';
        const completedDate = e.completed_at ? new Date(e.completed_at).toLocaleString() : 'Unknown';

        // Link directly with attempt_id when available for reliability
        const reviewLink = e.attempt_id
            ? `/exam-review.html?attempt_id=${e.attempt_id}&user_id=${user.id}`
            : `/exam-review.html?exam_id=${e.exam_id}&user_id=${user.id}`;

        html += `
            <div class="list-group-item">
                <div class="d-flex justify-content-between align-items-center">
                    <div>
                        <strong>${escapeHtml(e.exam_name)}</strong>
                        <span class="badge bg-secondary ms-2">${e.question_count} questions</span>
                        <br>
                        <small class="text-muted">Submitted: ${completedDate}</small>
                        <br>
                        <span class="badge bg-success">Submitted</span>
                        <span class="badge ${scoreClass} ms-1">${scoreVal}%</span>
                        <small class="text-muted ms-2">(${e.correct_answers}/${e.total_questions} correct)</small>
                    </div>
                    <div class="d-flex flex-column gap-2 align-items-end">
                        <a href="${reviewLink}" class="btn btn-outline-primary btn-sm">
                            <i class="bi bi-eye"></i> View Submission
                        </a>
                        <button class="btn btn-outline-danger btn-sm"
                            data-exam-id="${e.exam_id}"
                            data-token="${escapeHtml(e.access_token)}"
                            onclick="startExam(this.dataset.examId, this.dataset.token)">
                            <i class="bi bi-arrow-repeat"></i> Retake
                        </button>
                    </div>
                </div>
            </div>
        `;
    });
    html += '</div>';
    container.innerHTML = html;
}

async function startExam(examId, token) {
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
        console.error('Error starting exam:', error);
        alert('Error starting exam. Please try again.');
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text == null ? '' : String(text);
    return div.innerHTML;
}

// Load exams on page load
loadExams();
