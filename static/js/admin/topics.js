// Topic Management Module
import API from '../utils/api.js';
import { escapeHtml } from '../utils/sanitize.js';

let currentUser = null;

export function setCurrentUser(user) {
    currentUser = user;
}

export async function loadTopics() {
    try {
        const topics = await API.get('/admin/topics');
        const tbody = document.getElementById('topicsTable');
        tbody.innerHTML = '';

        topics.forEach(topic => {
            const row = document.createElement('tr');
            let categoryBadge = '';
            if (topic.category === 'java') {
                categoryBadge = '<span class="badge bg-primary">Java</span>';
            } else if (topic.category === 'python') {
                categoryBadge = '<span class="badge bg-success">Python</span>';
            } else if (topic.category === 'all') {
                categoryBadge = '<span class="badge bg-warning">All</span>';
            } else {
                categoryBadge = `<span class="badge bg-secondary">${escapeHtml(topic.category || '')}</span>`;
            }
            row.innerHTML = `
                <td>${topic.id}</td>
                <td>${escapeHtml(topic.name)}</td>
                <td>${escapeHtml(topic.description)}</td>
                <td>${categoryBadge}</td>
                <td>${new Date(topic.created_at).toLocaleDateString()}</td>
                <td>
                    <button class="btn btn-sm btn-warning" data-action="edit-topic" data-id="${topic.id}">Edit</button>
                    <button class="btn btn-sm btn-danger" data-action="delete-topic" data-id="${topic.id}">Delete</button>
                </td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        console.error('Error loading topics:', error);
        alert('Failed to load topics');
    }
}

export function showCreateTopicModal() {
    document.getElementById('topicModalTitle').textContent = 'Create Topic';
    document.getElementById('topicId').value = '';
    document.getElementById('topicName').value = '';
    document.getElementById('topicDescription').value = '';
    document.getElementById('topicCategory').value = 'java';
    new bootstrap.Modal(document.getElementById('topicModal')).show();
}

export async function editTopic(topicId) {
    const topics = await API.get('/admin/topics');
    const topic = topics.find(t => t.id === topicId);
    if (!topic) return;

    document.getElementById('topicModalTitle').textContent = 'Edit Topic';
    document.getElementById('topicId').value = topic.id;
    document.getElementById('topicName').value = topic.name;
    document.getElementById('topicDescription').value = topic.description;
    document.getElementById('topicCategory').value = topic.category || 'java';
    new bootstrap.Modal(document.getElementById('topicModal')).show();
}

export async function saveTopic() {
    const topicId = document.getElementById('topicId').value;
    const topicData = {
        name: document.getElementById('topicName').value,
        description: document.getElementById('topicDescription').value,
        category: document.getElementById('topicCategory').value,
    };

    try {
        if (topicId) {
            // Update existing topic
            topicData.id = parseInt(topicId);
            await API.put(`/admin/topics?id=${topicId}`, topicData);
        } else {
            // Create new topic
            await API.post('/admin/topics', topicData);
        }

        bootstrap.Modal.getInstance(document.getElementById('topicModal')).hide();
        await loadTopics();
    } catch (error) {
        console.error('Error saving topic:', error);
        alert('Failed to save topic');
    }
}

export async function deleteTopic(topicId) {
    if (!confirm('Are you sure you want to delete this topic? This may affect related questions and quizzes.')) return;

    try {
        await API.delete(`/admin/topics?id=${topicId}`);
        await loadTopics();
    } catch (error) {
        console.error('Error deleting topic:', error);
        alert('Failed to delete topic');
    }
}
