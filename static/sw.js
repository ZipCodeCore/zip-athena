const CACHE_NAME = 'athena-v4';

const SHELL_FILES = [
    '/css/style.css',
    '/js/auth.js',
    '/js/dashboard.js',
    '/js/quiz.js',
    '/js/learning.js',
    '/js/admin.js',
];

self.addEventListener('install', (event) => {
    event.waitUntil(
        caches.open(CACHE_NAME).then((cache) => cache.addAll(SHELL_FILES))
    );
});

self.addEventListener('activate', (event) => {
    event.waitUntil(
        caches.keys().then((names) =>
            Promise.all(names.filter((n) => n !== CACHE_NAME).map((n) => caches.delete(n)))
        )
    );
});

self.addEventListener('fetch', (event) => {
    const url = new URL(event.request.url);

    // Skip all cross-origin requests (e.g. Bootstrap CDN) — let the browser handle them
    if (url.origin !== self.location.origin) return;

    // Skip navigation requests (page loads) - let browser handle redirects naturally
    if (event.request.mode === 'navigate') return;

    // Cache same-origin static assets — CSS, JS only (not HTML)
    if (url.pathname.endsWith('.css') || url.pathname.endsWith('.js')) {
        event.respondWith(
            caches.match(event.request).then((cached) => {
                const network = fetch(event.request).then((response) => {
                    // Only cache successful responses
                    if (response.ok) {
                        caches.open(CACHE_NAME).then((cache) => cache.put(event.request, response.clone()));
                    }
                    return response;
                });
                // Return cached version immediately, update cache in background
                return cached || network;
            })
        );
    }
    // All other requests (API, health, HTML, etc.) use default network behavior
});
