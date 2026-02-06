const CACHE_NAME = 'athena-v5';

const SHELL_FILES = [
    '/index.html',
    '/dashboard.html',
    '/quiz.html',
    '/learning.html',
    '/instructor.html',
    '/admin.html',
    '/css/style.css',
    '/js/auth.js',
    '/js/dashboard.js',
    '/js/quiz.js',
    '/js/learning.js',
    '/js/admin.js',
    '/manifest.json'
];

self.addEventListener('install', (event) => {
    event.waitUntil(
        caches.open(CACHE_NAME).then((cache) => cache.addAll(SHELL_FILES))
    );
    // Force the waiting service worker to become the active service worker
    self.skipWaiting();
});

self.addEventListener('activate', (event) => {
    event.waitUntil(
        caches.keys().then((names) =>
            Promise.all(names.filter((n) => n !== CACHE_NAME).map((n) => caches.delete(n)))
        ).then(() => {
            // Claim all clients immediately
            return self.clients.claim();
        })
    );
});

self.addEventListener('fetch', (event) => {
    const url = new URL(event.request.url);

    // Skip all cross-origin requests (e.g. Bootstrap CDN) — let the browser handle them
    if (url.origin !== self.location.origin) return;

    // For navigation requests, try network first, then fall back to cache
    if (event.request.mode === 'navigate') {
        event.respondWith(
            fetch(event.request)
                .then((response) => {
                    // Cache the page if successful
                    if (response.ok) {
                        caches.open(CACHE_NAME).then((cache) => cache.put(event.request, response.clone()));
                    }
                    return response;
                })
                .catch(() => {
                    // If offline, try to serve from cache
                    return caches.match(event.request).then((cached) => {
                        return cached || caches.match('/index.html');
                    });
                })
        );
        return;
    }

    // Cache same-origin static assets — CSS, JS, HTML
    if (url.pathname.endsWith('.css') || url.pathname.endsWith('.js') || url.pathname.endsWith('.html')) {
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
    // All other requests (API, health, etc.) use default network behavior
});
