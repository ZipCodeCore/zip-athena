// HTML escaping utility to prevent XSS attacks
export function escapeHtml(text) {
    if (text === null || text === undefined) {
        return '';
    }
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Escape HTML in an object's properties
export function escapeObject(obj, keys) {
    const escaped = { ...obj };
    keys.forEach(key => {
        if (escaped[key] !== undefined) {
            escaped[key] = escapeHtml(escaped[key]);
        }
    });
    return escaped;
}
