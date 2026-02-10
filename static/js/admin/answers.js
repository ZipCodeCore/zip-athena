// Answer Management Module
import API from '../utils/api.js';
import { escapeHtml } from '../utils/sanitize.js';

let currentUser = null;

export function setCurrentUser(user) {
    currentUser = user;
}

export async function loadAnswers(questionId) {
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
            strong.textContent = (answer.is_correct ? '✓ ' : '') + escapeHtml(answer.answer_text);
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

export function showAddAnswerForm() {
    document.getElementById('addAnswerForm').style.display = 'block';
    document.getElementById('newAnswerText').value = '';
    document.getElementById('newAnswerIsCorrect').checked = false;
}

export function cancelAddAnswer() {
    document.getElementById('addAnswerForm').style.display = 'none';
}

export async function addAnswer() {
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
        await API.post('/admin/answers', answerData);
        cancelAddAnswer();
        await loadAnswers(questionId);
    } catch (error) {
        console.error('Error adding answer:', error);
        alert('Failed to add answer');
    }
}

export async function deleteAnswer(answerId) {
    if (!confirm('Are you sure you want to delete this answer?')) return;

    const questionId = document.getElementById('questionId').value;

    try {
        await API.delete(`/admin/answers?id=${answerId}`);
        await loadAnswers(questionId);
    } catch (error) {
        console.error('Error deleting answer:', error);
        alert('Failed to delete answer');
    }
}

export function showEditAnswerForm(answer) {
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
        <button type="button" class="btn btn-sm btn-primary" id="saveEditedAnswerBtn">Save Changes</button>
        <button type="button" class="btn btn-sm btn-secondary" id="cancelEditAnswerBtn">Cancel</button>
    `;
    
    // Insert the edit form after the add answer form
    const addAnswerForm = document.getElementById('addAnswerForm');
    addAnswerForm.parentNode.insertBefore(editForm, addAnswerForm.nextSibling);
    
    // Add event listeners
    editForm.querySelector('#saveEditedAnswerBtn').addEventListener('click', saveEditedAnswer);
    editForm.querySelector('#cancelEditAnswerBtn').addEventListener('click', cancelEditAnswer);
    
    return editForm;
}

export function cancelEditAnswer() {
    const editForm = document.getElementById('editAnswerForm');
    if (editForm) {
        editForm.style.display = 'none';
    }
}

export async function saveEditedAnswer() {
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
        await API.put(`/admin/answers?id=${encodeURIComponent(answerId)}`, answerData);
        cancelEditAnswer();
        await loadAnswers(questionId);
    } catch (error) {
        console.error('Error updating answer:', error);
        alert('Failed to update answer');
    }
}
