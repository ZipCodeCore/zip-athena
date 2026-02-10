// Centralized API calls with error handling and auth
const API = {
    async request(endpoint, options = {}) {
        const token = localStorage.getItem('token');
        const userData = localStorage.getItem('user');
        let userId = null;
        
        if (userData) {
            try {
                const user = JSON.parse(userData);
                userId = user.id;
            } catch (e) {
                console.error('Failed to parse user data:', e);
            }
        }
        
        const config = {
            headers: {
                'Content-Type': 'application/json',
                ...(token && { 'Authorization': `Bearer ${token}` }),
                ...(userId && { 'X-User-ID': userId.toString() }),
                ...options.headers
            },
            ...options
        };

        const response = await fetch(endpoint, config);
        if (!response.ok) {
            const error = await response.json().catch(() => ({ message: response.statusText }));
            throw new Error(error.message || 'Request failed');
        }
        return response.json();
    },

    get: (url, options = {}) => API.request(url, options),
    post: (url, data, options = {}) => API.request(url, { method: 'POST', body: JSON.stringify(data), ...options }),
    put: (url, data, options = {}) => API.request(url, { method: 'PUT', body: JSON.stringify(data), ...options }),
    delete: (url, options = {}) => API.request(url, { method: 'DELETE', ...options })
};

export default API;
