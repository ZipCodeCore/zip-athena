// Quiz Management Module
import API from '../utils/api.js';
import { escapeHtml } from '../utils/sanitize.js';

let currentUser = null;

export function setCurrentUser(user) {
    currentUser = user;
}

export async function loadQuizzes() {
    const tbody = document.getElementById('quizzesTable');
    if (!tbody) return;
    try {
        const quizzes = await API.get('/admin/quizzes');
        tbody.innerHTML = '';

        quizzes.forEach(quiz => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${quiz.id}</td>
                <td>${escapeHtml(quiz.name)}</td>
                <td>${quiz.topic_id}</td>
                <td>${escapeHtml(quiz.description)}</td>
                <td>${new Date(quiz.created_at).toLocaleDateString()}</td>
                <td>
                    <button class="btn btn-sm btn-warning" data-action="edit-quiz" data-id="${quiz.id}">Edit</button>
                    <button class="btn btn-sm btn-danger" data-action="delete-quiz" data-id="${quiz.id}">Delete</button>
                </td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        console.error('Error loading quizzes:', error);
        tbody.innerHTML = '<tr><td colspan="6" class="text-center text-danger py-3">Failed to load quizzes. Please refresh the page to try again.</td></tr>';
    }
}

export function showCreateQuizModal() {
    document.getElementById('quizModalTitle').textContent = 'Create Quiz';
    document.getElementById('quizId').value = '';
    document.getElementById('quizName').value = '';
    document.getElementById('quizTopicId').value = '';
    document.getElementById('quizDescription').value = '';
    
    // Hide quiz questions section for new quizzes
    document.getElementById('quizQuestionsSection').style.display = 'none';
    
    new bootstrap.Modal(document.getElementById('quizModal')).show();
}

export async function editQuiz(quizId) {
    // Fetch quiz details
    const quizzes = await API.get('/admin/quizzes');
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

export async function saveQuiz() {
    const quizId = document.getElementById('quizId').value;
    const quizData = {
        name: document.getElementById('quizName').value,
        topic_id: parseInt(document.getElementById('quizTopicId').value),
        description: document.getElementById('quizDescription').value
    };

    try {
        if (quizId) {
            // Update existing quiz
            quizData.id = parseInt(quizId);
            await API.put(`/admin/quizzes?id=${quizId}`, quizData);
        } else {
            // Create new quiz
            await API.post('/admin/quizzes', quizData);
        }

        bootstrap.Modal.getInstance(document.getElementById('quizModal')).hide();
        await loadQuizzes();
    } catch (error) {
        console.error('Error saving quiz:', error);
        alert('Failed to save quiz');
    }
}

export async function deleteQuiz(quizId) {
    if (!confirm('Are you sure you want to delete this quiz?')) return;

    try {
        await API.delete(`/admin/quizzes?id=${quizId}`);
        await loadQuizzes();
    } catch (error) {
        console.error('Error deleting quiz:', error);
        alert('Failed to delete quiz');
    }
}

export async function loadQuizQuestions(quizId) {
    try {
        const questions = await API.get(`/admin/quiz-questions?quiz_id=${quizId}`);
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

export function showAddQuestionToQuizForm() {
    document.getElementById('addQuestionToQuizForm').style.display = 'block';
    document.getElementById('quizQuestionId').value = '';
    
    // Auto-increment order based on current questions
    const questionsList = document.getElementById('quizQuestionsList');
    const currentCount = questionsList.querySelectorAll('.border').length;
    document.getElementById('quizQuestionOrder').value = currentCount;
}

export function cancelAddQuestionToQuiz() {
    document.getElementById('addQuestionToQuizForm').style.display = 'none';
}

export async function addQuestionToQuiz() {
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
        await API.post('/admin/quiz-questions', {
            quiz_id: parseInt(quizId),
            question_id: questionId,
            order: order
        });
        
        cancelAddQuestionToQuiz();
        await loadQuizQuestions(quizId);
    } catch (error) {
        console.error('Error adding question to quiz:', error);
        alert('Failed to add question to quiz: ' + error.message);
    }
}

export async function removeQuestionFromQuiz(quizId, questionId) {
    if (!confirm('Are you sure you want to remove this question from the quiz?')) return;

    try {
        await API.delete('/admin/quiz-questions', {
            body: JSON.stringify({
                quiz_id: quizId,
                question_id: questionId
            })
        });
        
        await loadQuizQuestions(quizId);
    } catch (error) {
        console.error('Error removing question from quiz:', error);
        alert('Failed to remove question from quiz');
    }
}
