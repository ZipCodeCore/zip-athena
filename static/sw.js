const CACHE_NAME = 'athena-v2';

const SHELL_FILES = [
    '/',
    '/index.html',
    '/dashboard.html',
    '/quiz.html',
    '/learning.html',
    '/admin.html',
    '/instructor.html',
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

    // Cache same-origin static assets — HTML pages, CSS, JS
    if (url.pathname.endsWith('.css') || url.pathname.endsWith('.js') ||
        url.pathname === '/' || url.pathname.endsWith('.html')) {
        event.respondWith(
            caches.match(event.request).then((cached) => {
                const network = fetch(event.request).then((response) => {
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
    // All other requests (API, health, etc.) use default network behavior
});
