// Answer Management Module
import API from '../utils/api.js';
import { escapeHtml } from '../utils/sanitize.js';

let currentUser = null;

export function setCurrentUser(user) {
    currentUser = user;
}

export async function loadAnswers(questionId) {
    try {
        const answers = await API.get(`/admin/answers?question_id=${questionId}`);
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
    
    // Create heading
    const heading = document.createElement('h6');
    heading.textContent = 'Edit Answer';
    editForm.appendChild(heading);
    
    // Create hidden input
    const hiddenInput = document.createElement('input');
    hiddenInput.type = 'hidden';
    hiddenInput.id = 'editAnswerId';
    editForm.appendChild(hiddenInput);
    
    // Create answer text field
    const textDiv = document.createElement('div');
    textDiv.className = 'mb-2';
    const textLabel = document.createElement('label');
    textLabel.htmlFor = 'editAnswerText';
    textLabel.className = 'form-label';
    textLabel.textContent = 'Answer Text';
    const textInput = document.createElement('input');
    textInput.type = 'text';
    textInput.className = 'form-control';
    textInput.id = 'editAnswerText';
    textInput.required = true;
    textDiv.appendChild(textLabel);
    textDiv.appendChild(textInput);
    editForm.appendChild(textDiv);
    
    // Create checkbox field
    const checkDiv = document.createElement('div');
    checkDiv.className = 'mb-2 form-check';
    const checkInput = document.createElement('input');
    checkInput.type = 'checkbox';
    checkInput.className = 'form-check-input';
    checkInput.id = 'editAnswerIsCorrect';
    const checkLabel = document.createElement('label');
    checkLabel.className = 'form-check-label';
    checkLabel.htmlFor = 'editAnswerIsCorrect';
    checkLabel.textContent = 'Is Correct Answer';
    checkDiv.appendChild(checkInput);
    checkDiv.appendChild(checkLabel);
    editForm.appendChild(checkDiv);
    
    // Create save button
    const saveBtn = document.createElement('button');
    saveBtn.type = 'button';
    saveBtn.className = 'btn btn-sm btn-primary me-2';
    saveBtn.textContent = 'Save Changes';
    saveBtn.addEventListener('click', saveEditedAnswer);
    editForm.appendChild(saveBtn);
    
    // Create cancel button
    const cancelBtn = document.createElement('button');
    cancelBtn.type = 'button';
    cancelBtn.className = 'btn btn-sm btn-secondary';
    cancelBtn.textContent = 'Cancel';
    cancelBtn.addEventListener('click', cancelEditAnswer);
    editForm.appendChild(cancelBtn);
    
    // Insert the edit form after the add answer form
    const addAnswerForm = document.getElementById('addAnswerForm');
    addAnswerForm.parentNode.insertBefore(editForm, addAnswerForm.nextSibling);
    
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
