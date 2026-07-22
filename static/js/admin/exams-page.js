import { requireAdminUser, setupLogout } from './auth.js';

setupLogout();

let userId = null;

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text == null ? '' : String(text);
    return div.innerHTML;
}

async function loadExamTopics() {
    const container = document.getElementById('examCategoriesCheckboxes');
    try {
        const response = await fetch('/instructor/exams/topics', {headers: {'X-User-ID': userId}});
        const topics = await response.json();
        if (!topics || topics.length === 0) {
            container.innerHTML = '<small class="text-muted">No topics available.</small>';
            return;
        }
        container.innerHTML = '';

        // Group by category for display
        const byCategory = {};
        topics.forEach((topic) => {
            const cat = topic.category || 'Other';
            if (!byCategory[cat]) byCategory[cat] = [];
            byCategory[cat].push(topic);
        });

        Object.keys(byCategory).sort().forEach((cat) => {
            const groupDiv = document.createElement('div');
            groupDiv.className = 'mb-1';

            const catLabel = document.createElement('small');
            catLabel.className = 'text-muted fw-bold d-block';
            catLabel.textContent = cat.toUpperCase();
            groupDiv.appendChild(catLabel);

            byCategory[cat].forEach((topic) => {
                const topicId = `examTopic_${topic.id}`;
                const div = document.createElement('div');
                div.className = 'form-check form-check-inline';

                const input = document.createElement('input');
                input.className = 'form-check-input exam-topic-check';
                input.type = 'checkbox';
                input.value = String(topic.id);
                input.id = topicId;

                const label = document.createElement('label');
                label.className = 'form-check-label';
                label.htmlFor = topicId;
                label.textContent = topic.name;

                div.appendChild(input);
                div.appendChild(label);
                groupDiv.appendChild(div);
            });

            container.appendChild(groupDiv);
        });
    } catch (e) {
        container.innerHTML = '<small class="text-danger">Error loading topics.</small>';
    }
}

async function generateExam() {
    const name = document.getElementById('newExamName').value.trim();
    const resultDiv = document.getElementById('generateExamResult');

    if (!name) {
        alert('Exam name is required.');
        return;
    }

    const checkedTopics = Array.from(document.querySelectorAll('.exam-topic-check:checked')).map((cb) => parseInt(cb.value, 10));
    if (checkedTopics.length === 0) {
        alert('Please select at least one topic.');
        return;
    }

    const checkedDifficultyLevels = Array.from(document.querySelectorAll('.exam-difficulty-check:checked')).map((cb) => parseInt(cb.value, 10));
    if (checkedDifficultyLevels.length === 0) {
        alert('Please select at least one difficulty level.');
        return;
    }

    const questionCount = parseInt(document.getElementById('examQuestionCount').value, 10) || 30;

    resultDiv.style.display = 'none';
    try {
        const response = await fetch('/instructor/exams/generate', {
            method: 'POST',
            headers: {'Content-Type': 'application/json', 'X-User-ID': userId},
            body: JSON.stringify({
                name,
                topic_ids: checkedTopics,
                difficulty_levels: checkedDifficultyLevels,
                question_count: questionCount,
            }),
        });
        const result = await response.json();

        if (response.ok && result.success) {
            const examUrl = window.location.origin + result.exam_url;
            resultDiv.className = 'alert alert-success';
            resultDiv.innerHTML = `✅ Exam "<strong>${escapeHtml(result.exam_name)}</strong>" created with ${result.question_count} questions.<br>
                <strong>Exam URL:</strong> <a href="${escapeHtml(result.exam_url)}" target="_blank">${escapeHtml(examUrl)}</a>`;
            resultDiv.style.display = 'block';
            document.getElementById('newExamName').value = '';
            loadAssignedExamDropdown();
        } else {
            resultDiv.className = 'alert alert-danger';
            resultDiv.innerHTML = `❌ Failed to generate exam: ${escapeHtml(result.error || JSON.stringify(result))}`;
            resultDiv.style.display = 'block';
        }
    } catch (e) {
        resultDiv.className = 'alert alert-danger';
        resultDiv.innerHTML = '❌ Network error: ' + e.message;
        resultDiv.style.display = 'block';
    }
}

async function loadAssignedExamDropdown() {
    const select = document.getElementById('assignExamSelect');
    select.innerHTML = '<option value="">-- Select an exam --</option>';
    try {
        const response = await fetch('/instructor/exams', {headers: {'X-User-ID': userId}});
        const exams = await response.json();
        exams.forEach((exam) => {
            const opt = document.createElement('option');
            opt.value = exam.id;
            opt.dataset.token = exam.access_token;
            opt.textContent = `${exam.name} (${exam.question_count} questions)`;
            select.appendChild(opt);
        });
    } catch (e) {
        console.error('Error loading exams for dropdown:', e);
    }
}

function updateExamUrlDisplay() {
    const select = document.getElementById('assignExamSelect');
    const opt = select.options[select.selectedIndex];
    const token = opt ? opt.dataset.token : '';
    const urlInput = document.getElementById('examUrlDisplay');
    if (token) {
        urlInput.value = window.location.origin + '/exam.html?token=' + token;
    } else {
        urlInput.value = '';
    }
}

function copyExamUrl() {
    const urlInput = document.getElementById('examUrlDisplay');
    if (urlInput.value) {
        navigator.clipboard.writeText(urlInput.value).then(() => {
            alert('URL copied to clipboard!');
        }).catch(() => {
            urlInput.select();
            document.execCommand('copy');
        });
    }
}

async function loadExamAssignStudentsList() {
    const container = document.getElementById('examAssignStudentsList');
    container.innerHTML = '<small class="text-muted">Loading...</small>';
    try {
        const response = await fetch('/instructor/students', {headers: {'X-User-ID': userId}});
        const students = await response.json();
        if (!students || students.length === 0) {
            container.innerHTML = '<small class="text-muted">No active students found.</small>';
            return;
        }
        container.innerHTML = '';
        students.forEach((student) => {
            if (student.is_admin) return;

            const studentId = Number(student.user_id ?? student.id);
            if (!Number.isFinite(studentId)) return;

            const div = document.createElement('div');
            div.className = 'form-check';

            const input = document.createElement('input');
            input.className = 'form-check-input exam-student-check';
            input.type = 'checkbox';
            input.value = String(studentId);
            input.id = `examStudent_${studentId}`;

            const label = document.createElement('label');
            label.className = 'form-check-label';
            label.htmlFor = input.id;
            label.textContent = student.username;

            div.appendChild(input);
            div.appendChild(label);
            container.appendChild(div);
        });
    } catch (e) {
        container.innerHTML = '<small class="text-danger">Error loading students.</small>';
    }
}

function selectAllExamStudents() {
    document.querySelectorAll('.exam-student-check').forEach((cb) => {
        cb.checked = true;
    });
}

function deselectAllExamStudents() {
    document.querySelectorAll('.exam-student-check').forEach((cb) => {
        cb.checked = false;
    });
}

async function assignExamToStudents() {
    const examId = parseInt(document.getElementById('assignExamSelect').value, 10);
    const resultDiv = document.getElementById('assignExamResult');

    if (!examId) {
        alert('Please select an exam.');
        return;
    }

    const checked = document.querySelectorAll('.exam-student-check:checked');
    if (checked.length === 0) {
        alert('Please select at least one student.');
        return;
    }

    const userIds = Array.from(checked).map((cb) => parseInt(cb.value, 10));
    resultDiv.style.display = 'none';
    try {
        const response = await fetch(`/instructor/exams/assign?exam_id=${examId}`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json', 'X-User-ID': userId},
            body: JSON.stringify({user_ids: userIds}),
        });
        const result = await response.json();

        if (response.ok) {
            resultDiv.className = 'alert alert-success';
            resultDiv.innerHTML = `✅ Exam assigned to ${result.assigned} student(s).`;
            resultDiv.style.display = 'block';
            loadExamStatus();
        } else {
            resultDiv.className = 'alert alert-danger';
            resultDiv.innerHTML = '❌ Failed to assign exam.';
            resultDiv.style.display = 'block';
        }
    } catch (e) {
        resultDiv.className = 'alert alert-danger';
        resultDiv.innerHTML = '❌ Network error: ' + e.message;
        resultDiv.style.display = 'block';
    }
}

async function deleteExam(examId, examName) {
    if (!confirm(`Delete exam "${examName}"? This will also remove all assignments and attempts.`)) {
        return;
    }
    try {
        const response = await fetch(`/instructor/exams?exam_id=${examId}`, {
            method: 'DELETE',
            headers: {'X-User-ID': userId},
        });
        if (response.ok) {
            loadExamStatus();
            loadAssignedExamDropdown();
        } else {
            const text = await response.text();
            alert('Failed to delete exam: ' + text);
        }
    } catch (e) {
        alert('Network error: ' + e.message);
    }
}

async function loadExamStatus() {
    const container = document.getElementById('examSummaryTable');
    container.innerHTML = '<p class="text-muted">Loading...</p>';
    try {
        const response = await fetch('/instructor/exams', {headers: {'X-User-ID': userId}});
        const exams = await response.json();

        if (!exams || exams.length === 0) {
            container.innerHTML = '<p class="text-muted">No exams have been created yet.</p>';
            return;
        }

        let html = `<div class="table-responsive"><table class="table table-striped">
            <thead><tr>
                <th>Exam Name</th><th>Categories</th><th>Difficulty Levels</th>
                <th>Questions</th><th>Exam URL</th><th>Actions</th>
            </tr></thead><tbody>`;
        exams.forEach((exam) => {
            const examId = Number(exam.id);
            if (!Number.isFinite(examId)) {
                return;
            }

            const categoriesDisplay = Array.isArray(exam.categories) ? exam.categories.join(', ') : '';
            const difficultyLevelsDisplay = Array.isArray(exam.difficulty_levels) ? exam.difficulty_levels.map((d) => 'Level ' + d).join(', ') : '';
            const examUrl = window.location.origin + '/exam.html?token=' + exam.access_token;
            html += `<tr>
                <td><strong>${escapeHtml(exam.name)}</strong></td>
                <td><span class="badge bg-secondary">${escapeHtml(categoriesDisplay)}</span></td>
                <td>${escapeHtml(difficultyLevelsDisplay)}</td>
                <td>${exam.question_count}</td>
                <td><a href="/exam.html?token=${encodeURIComponent(exam.access_token)}" target="_blank" class="small">${escapeHtml(examUrl)}</a></td>
                <td>
                    <button class="btn btn-sm btn-primary me-1" data-exam-id="${examId}" data-exam-name="${encodeURIComponent(exam.name)}" onclick="showExamDetail(parseInt(this.dataset.examId, 10), decodeURIComponent(this.dataset.examName))">View Scores</button>
                    <button class="btn btn-sm btn-danger" data-exam-id="${examId}" data-exam-name="${encodeURIComponent(exam.name)}" onclick="deleteExam(parseInt(this.dataset.examId, 10), decodeURIComponent(this.dataset.examName))">Delete</button>
                </td>
            </tr>`;
        });
        html += '</tbody></table></div>';
        container.innerHTML = html;
    } catch (e) {
        container.innerHTML = '<p class="text-danger">Error loading exam status: ' + e.message + '</p>';
    }
}

async function unassignExamFromUser(examId, studentUserId, username) {
    if (!confirm(`Unassign exam from "${username}"?`)) {
        return;
    }
    try {
        const response = await fetch(`/instructor/exams/unassign?exam_id=${examId}&user_id=${studentUserId}`, {
            method: 'DELETE',
            headers: {'X-User-ID': userId},
        });
        if (response.ok) {
            showExamDetail(examId, _currentDetailExamName);
        } else {
            const text = await response.text();
            alert('Failed to unassign: ' + text);
        }
    } catch (e) {
        alert('Network error: ' + e.message);
    }
}

let _currentDetailExamId = null;
let _currentDetailExamName = null;

async function showExamDetail(examId, examName) {
    _currentDetailExamId = examId;
    _currentDetailExamName = examName;
    const panel = document.getElementById('examDetailPanel');
    const title = document.getElementById('examDetailTitle');
    const content = document.getElementById('examDetailContent');

    title.textContent = `Scores: ${examName}`;
    content.innerHTML = '<p class="text-muted">Loading...</p>';
    panel.style.display = 'block';
    panel.scrollIntoView({behavior: 'smooth'});

    try {
        const response = await fetch(`/instructor/exams/status?exam_id=${examId}`, {headers: {'X-User-ID': userId}});
        const status = await response.json();

        if (!status.student_results || status.student_results.length === 0) {
            content.innerHTML = '<p class="text-muted">No students assigned to this exam yet.</p>';
            return;
        }

        let html = `<div class="table-responsive"><table class="table table-striped">
            <thead><tr>
                <th>Student</th><th>Assigned</th><th>Completed</th><th>Score</th><th>Correct / Total</th><th>Actions</th>
            </tr></thead><tbody>`;
        status.student_results.forEach((result) => {
            const completed = result.completed_at
                ? new Date(result.completed_at).toLocaleString()
                : '<span class="text-muted">Not completed</span>';
            const scoreHtml = result.score !== undefined && result.score !== null
                ? `<span class="badge ${result.score >= 80 ? 'bg-success' : result.score >= 60 ? 'bg-warning' : 'bg-danger'}">${result.score.toFixed(0)}%</span>`
                : '<span class="text-muted">—</span>';
            const correct = result.correct_answers !== undefined && result.correct_answers !== null
                ? `${result.correct_answers} / ${result.total_questions}`
                : '—';
            const viewSubmissionBtn = result.attempt_id
                ? `<a href="/exam-review.html?attempt_id=${result.attempt_id}&user_id=${userId}" class="btn btn-sm btn-outline-primary me-1" target="_blank">View Submission</a>`
                : '';
            html += `<tr>
                <td>${escapeHtml(result.username)}</td>
                <td>${new Date(result.assigned_at).toLocaleDateString()}</td>
                <td>${completed}</td>
                <td>${scoreHtml}</td>
                <td>${correct}</td>
                <td>${viewSubmissionBtn}<button class="btn btn-sm btn-outline-danger" data-exam-id="${examId}" data-user-id="${result.user_id}" data-username="${encodeURIComponent(result.username)}" onclick="unassignExamFromUser(parseInt(this.dataset.examId,10), parseInt(this.dataset.userId,10), decodeURIComponent(this.dataset.username))">Unassign</button></td>
            </tr>`;
        });
        html += '</tbody></table></div>';
        content.innerHTML = html;
    } catch (e) {
        content.innerHTML = '<p class="text-danger">Error loading results: ' + e.message + '</p>';
    }
}

function hideExamDetail() {
    document.getElementById('examDetailPanel').style.display = 'none';
}

window.addEventListener('DOMContentLoaded', async () => {
    const currentUser = await requireAdminUser();
    if (!currentUser) return;

    userId = String(currentUser.id);
    document.getElementById('username').textContent = currentUser.username;
    document.getElementById('assignExamSelect').addEventListener('change', updateExamUrlDisplay);

    await loadExamTopics();
    await loadAssignedExamDropdown();
    await loadExamAssignStudentsList();
});

window.generateExam = generateExam;
window.loadExamStatus = loadExamStatus;
window.copyExamUrl = copyExamUrl;
window.selectAllExamStudents = selectAllExamStudents;
window.deselectAllExamStudents = deselectAllExamStudents;
window.assignExamToStudents = assignExamToStudents;
window.deleteExam = deleteExam;
window.showExamDetail = showExamDetail;
window.hideExamDetail = hideExamDetail;
window.unassignExamFromUser = unassignExamFromUser;
