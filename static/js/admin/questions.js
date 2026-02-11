// Question Management Module
import API from '../utils/api.js';
import { escapeHtml } from '../utils/sanitize.js';

let currentUser = null;
let importTopics = []; // Store topics for the import dropdown

export function setCurrentUser(user) {
    currentUser = user;
}

export async function loadQuestions() {
    try {
        const questions = await API.get('/admin/questions');
        const tbody = document.getElementById('questionsTable');
        tbody.innerHTML = '';

        questions.forEach(q => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${q.id}</td>
                <td>${escapeHtml(q.question_text.substring(0, 50))}...</td>
                <td>${q.topic_id}</td>
                <td>${q.difficulty}</td>
                <td>${escapeHtml(q.tags ? q.tags.join(', ') : '')}</td>
                <td>
                    <button class="btn btn-sm btn-warning" data-action="edit-question" data-id="${q.id}">Edit</button>
                    <button class="btn btn-sm btn-danger" data-action="delete-question" data-id="${q.id}">Delete</button>
                </td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        console.error('Error loading questions:', error);
        alert('Failed to load questions');
    }
}

// Load topics for import dropdown
export async function loadTopicsForImport() {
    try {
        importTopics = await API.get('/topics');
        const select = document.getElementById('importTopicSelect');
        select.innerHTML = '<option value="">-- Select a topic --</option>';
        
        importTopics.forEach(topic => {
            const option = document.createElement('option');
            option.value = topic.id;
            option.textContent = `${topic.name} (ID: ${topic.id})`;
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Error loading topics:', error);
        alert('Failed to load topics');
    }
}

export async function deleteQuestion(questionId) {
    if (!confirm('Are you sure you want to delete this question?')) return;

    try {
        await API.delete(`/admin/questions?id=${questionId}`);
        await loadQuestions();
    } catch (error) {
        console.error('Error deleting question:', error);
        alert('Failed to delete question');
    }
}

export async function editQuestion(questionId) {
    const questions = await API.get('/admin/questions');
    const question = questions.find(q => q.id === questionId);
    if (!question) return;

    document.getElementById('questionModalTitle').textContent = 'Edit Question';
    document.getElementById('questionId').value = question.id;
    document.getElementById('questionText').value = question.question_text;
    document.getElementById('questionTopicId').value = question.topic_id;
    document.getElementById('questionDifficulty').value = question.difficulty;
    document.getElementById('questionTags').value = question.tags ? question.tags.join(', ') : '';
    document.getElementById('questionTopicQuestionId').value = question.topic_question_id;
    
    // Show answers section and load answers
    document.getElementById('answersSection').style.display = 'block';
    
    // Import loadAnswers dynamically to avoid circular dependency
    const { loadAnswers } = await import('./answers.js');
    await loadAnswers(question.id);
    
    new bootstrap.Modal(document.getElementById('questionModal')).show();
}

export function showCreateQuestionModal() {
    document.getElementById('questionModalTitle').textContent = 'Create Question';
    document.getElementById('questionId').value = '';
    document.getElementById('questionText').value = '';
    document.getElementById('questionTopicId').value = '';
    document.getElementById('questionDifficulty').value = '1';
    document.getElementById('questionTags').value = '';
    document.getElementById('questionTopicQuestionId').value = '';
    
    // Hide answers section for new questions
    document.getElementById('answersSection').style.display = 'none';
    
    new bootstrap.Modal(document.getElementById('questionModal')).show();
}

export async function saveQuestion() {
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
        if (questionId) {
            // Update existing question
            await API.put(`/admin/questions?id=${questionId}`, questionData);
        } else {
            // Create new question
            await API.post('/admin/questions', questionData);
        }

        bootstrap.Modal.getInstance(document.getElementById('questionModal')).hide();
        await loadQuestions();
    } catch (error) {
        console.error('Error saving question:', error);
        alert('Failed to save question: ' + error.message);
    }
}

// Show import questions panel
window.showImportQuestionsPanel = async function() {
    await loadTopicsForImport();
    document.getElementById('importQuestionsPanel').style.display = 'block';
    document.getElementById('importResults').style.display = 'none';
    document.getElementById('importQuestionsForm').reset();
};

// Hide import questions panel
window.hideImportQuestionsPanel = function() {
    document.getElementById('importQuestionsPanel').style.display = 'none';
    document.getElementById('importResults').style.display = 'none';
    document.getElementById('importQuestionsForm').reset();
};

// Handle import form submission
document.addEventListener('DOMContentLoaded', function() {
    const importForm = document.getElementById('importQuestionsForm');
    if (importForm) {
        importForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const topicId = document.getElementById('importTopicSelect').value;
            const fileInput = document.getElementById('importFileInput');
            const file = fileInput.files[0];
            
            if (!topicId) {
                alert('Please select a topic');
                return;
            }
            
            if (!file) {
                alert('Please select a JSON file');
                return;
            }
            
            // Show loading state
            const submitBtn = e.target.querySelector('button[type="submit"]');
            const originalText = submitBtn.textContent;
            submitBtn.disabled = true;
            submitBtn.textContent = 'Importing...';
            
            try {
                const formData = new FormData();
                formData.append('file', file);
                formData.append('topic_id', topicId);
                
                const response = await fetch('/admin/questions/import-json', {
                    method: 'POST',
                    body: formData,
                    headers: {
                        'X-User-ID': localStorage.getItem('userId') || '1'
                    }
                });
                
                const result = await response.json();
                
                // Display results
                const resultsDiv = document.getElementById('importResults');
                const alertDiv = document.getElementById('importResultsAlert');
                resultsDiv.style.display = 'block';
                
                if (response.ok && result.success) {
                    alertDiv.className = 'alert alert-success';
                    let message = `<strong>Import Successful!</strong><br>`;
                    message += `Topic: ${escapeHtml(result.topic)}<br>`;
                    message += `Successfully imported: ${result.imported} questions<br>`;
                    if (result.skipped > 0) {
                        message += `Skipped (duplicates): ${result.skipped}<br>`;
                    }
                    if (result.errors > 0) {
                        message += `<span class="text-warning">Errors: ${result.errors}</span><br>`;
                        if (result.messages && result.messages.length > 0) {
                            message += '<br><strong>Error details:</strong><ul>';
                            result.messages.forEach(msg => {
                                message += `<li>${escapeHtml(msg)}</li>`;
                            });
                            message += '</ul>';
                        }
                    }
                    alertDiv.innerHTML = message;
                    
                    // Reload questions table
                    await loadQuestions();
                    
                    // Reset form
                    importForm.reset();
                } else {
                    alertDiv.className = 'alert alert-danger';
                    alertDiv.innerHTML = `<strong>Import Failed!</strong><br>${escapeHtml(result.message || 'Unknown error')}`;
                }
            } catch (error) {
                console.error('Error importing questions:', error);
                const resultsDiv = document.getElementById('importResults');
                const alertDiv = document.getElementById('importResultsAlert');
                resultsDiv.style.display = 'block';
                alertDiv.className = 'alert alert-danger';
                alertDiv.innerHTML = `<strong>Import Failed!</strong><br>Error: ${escapeHtml(error.message)}`;
            } finally {
                submitBtn.disabled = false;
                submitBtn.textContent = originalText;
            }
        });
    }
});
