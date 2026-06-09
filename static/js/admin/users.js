// User Management Module
import API from '../utils/api.js';
import { normalizeCategory, DEFAULT_CATEGORY } from '../utils/constants.js';
import { escapeHtml } from '../utils/sanitize.js';

let currentUser = null;

export function setCurrentUser(user) {
    currentUser = user;
}

export async function loadUsers() {
    const adminsTable = document.getElementById('adminsTable');
    const activeStudentsTable = document.getElementById('activeStudentsTable');
    const inactiveStudentsTable = document.getElementById('inactiveStudentsTable');
    if (!adminsTable || !activeStudentsTable || !inactiveStudentsTable) return;
    try {
        console.log('Loading users with X-User-ID:', currentUser.id.toString());
        const users = await API.get('/admin/users');
        console.log('Loaded users:', users.length, 'users');
        adminsTable.innerHTML = '';
        activeStudentsTable.innerHTML = '';
        inactiveStudentsTable.innerHTML = '';

        users.forEach(user => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${user.id}</td>
                <td>${escapeHtml(user.username)}</td>
                <td>${escapeHtml(user.email)}</td>
                <td>${escapeHtml(normalizeCategory(user.category))}</td>
                <td>${user.total_points}</td>
                <td>${new Date(user.created_at).toLocaleDateString()}</td>
                <td>
                    <button class="btn btn-sm btn-warning" data-action="edit-user" data-id="${user.id}">Edit</button>
                    <button class="btn btn-sm btn-danger" data-action="delete-user" data-id="${user.id}" 
                            ${user.id === currentUser.id ? 'disabled' : ''}>Delete</button>
                </td>
            `;
            if (user.is_admin) {
                adminsTable.appendChild(row);
            } else if (user.is_active) {
                activeStudentsTable.appendChild(row);
            } else {
                inactiveStudentsTable.appendChild(row);
            }
        });
    } catch (error) {
        console.error('Error loading users:', error);
        const msg = '<tr><td colspan="7" class="text-center text-danger py-3">Failed to load users. Please refresh the page to try again.</td></tr>';
        if (adminsTable) adminsTable.innerHTML = msg;
        if (activeStudentsTable) activeStudentsTable.innerHTML = msg;
        if (inactiveStudentsTable) inactiveStudentsTable.innerHTML = msg;
    }
}

export async function deleteUser(userId) {
    if (!confirm('Are you sure you want to delete this user?')) return;

    try {
        await API.delete(`/admin/users?id=${userId}`);
        await loadUsers();
    } catch (error) {
        console.error('Error deleting user:', error);
        alert('Failed to delete user');
    }
}

export async function editUser(userId) {
    const users = await API.get('/admin/users');
    const user = users.find(u => u.id === userId);
    if (!user) return;

    document.getElementById('userModalTitle').textContent = 'Edit User';
    document.getElementById('userId').value = user.id;
    document.getElementById('userUsername').value = user.username;
    document.getElementById('userEmail').value = user.email;
    document.getElementById('userPassword').value = '';
    document.getElementById('userPassword').required = false;
    document.getElementById('userTotalPoints').value = user.total_points;
    document.getElementById('userCategory').value = normalizeCategory(user.category);
    document.getElementById('userIsAdmin').checked = user.is_admin;
    document.getElementById('userIsActive').checked = user.is_active;
    new bootstrap.Modal(document.getElementById('userModal')).show();
}

export function showCreateUserModal() {
    document.getElementById('userModalTitle').textContent = 'Create User';
    document.getElementById('userId').value = '';
    document.getElementById('userUsername').value = '';
    document.getElementById('userEmail').value = '';
    document.getElementById('userPassword').value = '';
    document.getElementById('userPassword').required = true;
    document.getElementById('userTotalPoints').value = '0';
    document.getElementById('userCategory').value = DEFAULT_CATEGORY;
    document.getElementById('userIsAdmin').checked = false;
    document.getElementById('userIsActive').checked = true;
    new bootstrap.Modal(document.getElementById('userModal')).show();
}

export async function saveUser() {
    const userId = document.getElementById('userId').value;
    const totalPoints = parseInt(document.getElementById('userTotalPoints').value);
    
    // Validate total points
    if (isNaN(totalPoints) || totalPoints < 0) {
        alert('Total Points must be a valid non-negative number');
        return;
    }
    
    const userData = {
        username: document.getElementById('userUsername').value,
        email: document.getElementById('userEmail').value,
        total_points: totalPoints,
        category: document.getElementById('userCategory').value,
        is_admin: document.getElementById('userIsAdmin').checked,
        is_active: document.getElementById('userIsActive').checked
    };

    // Add password if creating or if it's been changed
    const password = document.getElementById('userPassword').value;
    if (password) {
        userData.password = password;
    }

    try {
        if (userId) {
            // Update existing user
            await API.put(`/admin/users?id=${userId}`, userData);
        } else {
            // Create new user
            if (!password) {
                alert('Password is required for new users');
                return;
            }
            await API.post('/admin/users', userData);
        }

        bootstrap.Modal.getInstance(document.getElementById('userModal')).hide();
        await loadUsers();
    } catch (error) {
        console.error('Error saving user:', error);
        alert('Failed to save user: ' + error.message);
    }
}
