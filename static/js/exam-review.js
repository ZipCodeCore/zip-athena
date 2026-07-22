// Exam Review JavaScript

const user = JSON.parse(localStorage.getItem('user'));
if (!user) {
    window.location.href = '/index.html';
}

document.getElementById('userDisplay').textContent = `Welcome, ${user.username}!`;

function logout() {
    localStorage.removeItem('user');
    window.location.href = '/index.html';
}

function showPanel(id) {
    ['loadingPanel', 'errorPanel', 'reviewPanel'].forEach(p => {
        const el = document.getElementById(p);
        if (el) el.style.display = p === id ? 'block' : 'none';
    });
}

function showError(msg) {
    document.getElementById('errorMsg').textContent = msg;
    showPanel('errorPanel');
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text == null ? '' : String(text);
    return div.innerHTML;
}

function renderQuestionText(container, text) {
    container.innerHTML = '';
    const codeBlockRegex = /```(\w+)?\r?\n?([\s\S]*?)```/g;
    let lastIndex = 0;
    let match;

    while ((match = codeBlockRegex.exec(text)) !== null) {
        if (match.index > lastIndex) {
            const preamble = text.substring(lastIndex, match.index).trim();
            if (preamble) {
                const p = document.createElement('p');
                p.textContent = preamble;
                container.appendChild(p);
            }
        }
        const pre = document.createElement('pre');
        const code = document.createElement('code');
        code.textContent = match[2];
        pre.appendChild(code);
        pre.style.cssText = 'background:#f4f4f4;padding:10px;border-radius:4px;overflow-x:auto;font-size:0.85em;';
        container.appendChild(pre);
        lastIndex = match.index + match[0].length;
    }

    const remaining = text.substring(lastIndex).trim();
    if (remaining) {
        const p = document.createElement('p');
        p.textContent = remaining;
        container.appendChild(p);
    }

    if (!container.hasChildNodes()) {
        container.textContent = text;
    }
}

(async function init() {
    const params = new URLSearchParams(window.location.search);

    const attemptId = params.get('attempt_id') ? parseInt(params.get('attempt_id'), 10) : null;

    if (!attemptId || attemptId <= 0) {
        showError('Missing or invalid attempt ID. Please navigate here from the Exam Center or the Admin Exams page.');
        return;
    }

    // Adjust back link for instructors
    if (user.is_admin) {
        const backLink = document.getElementById('examCenterLink');
        if (backLink) {
            backLink.href = '/admin-exams.html';
            backLink.textContent = '← Admin Exams';
        }
    }

    try {
        const response = await fetch(`/exam/submission?attempt_id=${attemptId}&user_id=${user.id}`);
        if (!response.ok) {
            const errText = await response.text();
            showError('Could not load submission: ' + errText);
            return;
        }
        const data = await response.json();
        renderReview(data);
    } catch (e) {
        showError('Network error: ' + e.message);
    }
})();

function renderReview(data) {
    // Fill summary
    document.getElementById('examNameTitle').textContent = data.exam_name || 'Exam Review';

    if (data.completed_at) {
        document.getElementById('completedAtBadge').textContent =
            'Submitted: ' + new Date(data.completed_at).toLocaleString();
    }

    const score = data.score !== undefined && data.score !== null ? data.score : 0;
    const scoreClass = score >= 80 ? 'text-success' : score >= 60 ? 'text-warning' : 'text-danger';
    document.getElementById('scoreDisplay').innerHTML =
        `<span class="${scoreClass}">${score.toFixed(0)}%</span>`;

    const correct = data.correct_answers || 0;
    const total = data.total_questions || 0;
    const incorrect = total - correct;

    document.getElementById('correctDisplay').textContent = correct;
    document.getElementById('incorrectDisplay').textContent = incorrect;
    document.getElementById('totalDisplay').textContent = total;

    // Render questions
    const container = document.getElementById('questionsReview');
    container.innerHTML = '';

    if (!data.questions || data.questions.length === 0) {
        container.innerHTML = '<p class="text-muted">No question details available.</p>';
        showPanel('reviewPanel');
        return;
    }

    data.questions.forEach((q, idx) => {
        const cardDiv = document.createElement('div');
        cardDiv.className = `card mb-3 question-card ${q.is_correct ? 'correct-card' : 'wrong-card'}`;

        const cardHeader = document.createElement('div');
        cardHeader.className = 'card-header d-flex justify-content-between align-items-center';

        const questionNum = document.createElement('span');
        questionNum.className = 'question-number';
        questionNum.textContent = `Question ${idx + 1} of ${data.questions.length}`;

        const statusBadge = document.createElement('span');
        if (q.is_correct) {
            statusBadge.className = 'badge bg-success';
            statusBadge.innerHTML = '<i class="bi bi-check-circle-fill"></i> Correct';
        } else {
            statusBadge.className = 'badge bg-danger';
            statusBadge.innerHTML = '<i class="bi bi-x-circle-fill"></i> Incorrect';
        }

        const categoryBadge = document.createElement('span');
        categoryBadge.className = 'badge bg-secondary ms-2';
        categoryBadge.textContent = q.category || '';

        cardHeader.appendChild(questionNum);
        const badgeGroup = document.createElement('div');
        badgeGroup.appendChild(categoryBadge);
        badgeGroup.appendChild(document.createTextNode(' '));
        badgeGroup.appendChild(statusBadge);
        cardHeader.appendChild(badgeGroup);

        const cardBody = document.createElement('div');
        cardBody.className = 'card-body';

        const questionTextEl = document.createElement('div');
        questionTextEl.className = 'mb-3 fw-semibold';
        renderQuestionText(questionTextEl, q.question_text);
        cardBody.appendChild(questionTextEl);

        // Render answers
        const answersDiv = document.createElement('div');
        if (q.answers && q.answers.length > 0) {
            q.answers.forEach(a => {
                const answerDiv = document.createElement('div');
                answerDiv.className = 'answer-option';

                const isSelected = a.id === q.selected_answer_id;
                const isCorrectAnswer = a.is_correct;

                if (isSelected && isCorrectAnswer) {
                    answerDiv.classList.add('answer-selected-correct');
                } else if (isSelected && !isCorrectAnswer) {
                    answerDiv.classList.add('answer-selected-wrong');
                } else if (!isSelected && isCorrectAnswer) {
                    answerDiv.classList.add('answer-correct-not-selected');
                } else {
                    answerDiv.classList.add('answer-neutral');
                }

                let prefix = '';
                if (isSelected && isCorrectAnswer) {
                    prefix = '✅ ';
                } else if (isSelected && !isCorrectAnswer) {
                    prefix = '❌ ';
                } else if (!isSelected && isCorrectAnswer) {
                    prefix = '✓ ';
                }

                answerDiv.textContent = prefix + a.answer_text;
                answersDiv.appendChild(answerDiv);
            });
        } else {
            answersDiv.innerHTML = '<small class="text-muted">No answer details available.</small>';
        }
        cardBody.appendChild(answersDiv);

        // If no answer was selected
        if (q.selected_answer_id === 0 && !q.is_correct) {
            const noAnswer = document.createElement('small');
            noAnswer.className = 'text-muted d-block mt-2';
            noAnswer.textContent = 'No answer was selected for this question.';
            cardBody.appendChild(noAnswer);
        }

        cardDiv.appendChild(cardHeader);
        cardDiv.appendChild(cardBody);
        container.appendChild(cardDiv);
    });

    showPanel('reviewPanel');
}
