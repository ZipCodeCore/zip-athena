// Learning Wizard JavaScript

let currentLanguage = 'java'; // 'java' or 'python'
let currentResourceId = null;
let allPhases = [];

// Check authentication
function checkAuth() {
    const user = JSON.parse(localStorage.getItem('user'));
    if (!user) {
        window.location.href = '/index.html';
        return false;
    }
    
    // Display user info
    const userDisplay = document.getElementById('userDisplay');
    if (userDisplay) {
        userDisplay.textContent = `Welcome, ${user.username}!`;
    }
    
    return true;
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    if (!checkAuth()) return;
    
    // Check if there's a resource ID in URL parameter
    const urlParams = new URLSearchParams(window.location.search);
    const resourceId = urlParams.get('resource');
    
    if (resourceId) {
        // Load specific resource
        loadResource(parseInt(resourceId));
    }
    
    loadOutline();
    setupEventListeners();
});

function setupEventListeners() {
    // Language toggle
    const toggle = document.getElementById('languageToggle');
    const label = document.getElementById('languageLabel');
    
    toggle.addEventListener('change', function() {
        if (toggle.checked) {
            currentLanguage = 'python';
            label.textContent = 'Python';
        } else {
            currentLanguage = 'java';
            label.textContent = 'Java';
        }
        loadOutline();
    });

    // Search functionality
    const searchInput = document.getElementById('searchInput');
    let searchTimeout;
    searchInput.addEventListener('input', function() {
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            if (searchInput.value.trim()) {
                searchResources(searchInput.value.trim());
            } else {
                loadOutline();
            }
        }, 300);
    });
}

async function loadOutline() {
    const container = document.getElementById('outlineContainer');
    container.innerHTML = `
        <div class="loading-spinner">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
    `;

    try {
        const response = await fetch(`/learning/outline?language=${currentLanguage}`);
        
        if (!response.ok) {
            throw new Error('Failed to load curriculum outline');
        }

        const data = await response.json();
        allPhases = data.phases;
        renderOutline(data.phases);
    } catch (error) {
        console.error('Error loading outline:', error);
        container.innerHTML = `
            <div class="alert alert-danger">
                <strong>Error!</strong> Failed to load curriculum outline.
            </div>
        `;
    }
}

function renderOutline(phases) {
    const container = document.getElementById('outlineContainer');
    
    if (!phases || phases.length === 0) {
        container.innerHTML = `
            <div class="alert alert-info">
                No learning resources available for ${currentLanguage}.
            </div>
        `;
        return;
    }

    let html = '';
    
    phases.forEach((phase, index) => {
        const isExpanded = index === 0; // Expand first phase by default
        const collapseId = `phase${phase.phase}`;
        
        html += `
            <div class="phase-section">
                <div class="phase-header ${isExpanded ? '' : 'collapsed'}" 
                     data-bs-toggle="collapse" 
                     data-bs-target="#${collapseId}">
                    Phase ${phase.phase}: ${phase.title}
                    <span class="float-end">
                        <span class="badge bg-light text-dark">${phase.resources.length}</span>
                    </span>
                </div>
                <div class="phase-description">${phase.description}</div>
                <div class="collapse ${isExpanded ? 'show' : ''}" id="${collapseId}">
                    <ul class="resource-list">
        `;
        
        phase.resources.forEach(resource => {
            const badgeClass = resource.language === 'java' ? 'java-badge' : 'python-badge';
            const viewedClass = resource.view_count > 0 ? 'viewed' : '';
            const activeClass = currentResourceId === resource.id ? 'active' : '';
            
            html += `
                <li class="resource-item ${viewedClass} ${activeClass}" 
                    onclick="loadResource(${resource.id})">
                    <span class="language-badge ${badgeClass}">${resource.language.toUpperCase()}</span>
                    <span>${resource.section_number}. ${resource.section_title}</span>
                    ${resource.view_count > 0 ? '<span class="text-success">✓</span>' : ''}
                </li>
            `;
        });
        
        html += `
                    </ul>
                </div>
            </div>
        `;
    });
    
    container.innerHTML = html;
}

async function loadResource(resourceId) {
    currentResourceId = resourceId;
    
    const contentContainer = document.getElementById('contentContainer');
    contentContainer.innerHTML = `
        <div class="loading-spinner">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
    `;

    try {
        // Load resource metadata
        console.log('Fetching metadata for resource:', resourceId);
        const metaResponse = await fetch(`/learning/resource/${resourceId}`);
        console.log('Metadata response status:', metaResponse.status, metaResponse.statusText);
        if (!metaResponse.ok) {
            const errorText = await metaResponse.text();
            console.error('Metadata fetch failed:', errorText);
            throw new Error('Failed to load resource metadata: ' + metaResponse.status);
        }
        const metadata = await metaResponse.json();
        console.log('Metadata loaded:', metadata);

        // Load resource content
        console.log('Fetching content for resource:', resourceId);
        const contentResponse = await fetch(`/learning/resource/${resourceId}/content`);
        console.log('Content response status:', contentResponse.status, contentResponse.statusText);
        if (!contentResponse.ok) {
            const errorText = await contentResponse.text();
            console.error('Content fetch failed:', errorText);
            throw new Error('Failed to load resource content: ' + contentResponse.status);
        }
        const markdownContent = await contentResponse.text();
        console.log('Content loaded, length:', markdownContent.length);

        // Render the content
        renderResourceContent(metadata, markdownContent);

        // Track the view
        trackView(resourceId);

        // Update outline to highlight current resource
        loadOutline();

    } catch (error) {
        console.error('Error loading resource:', error);
        contentContainer.innerHTML = `
            <div class="alert alert-danger">
                <strong>Error!</strong> Failed to load learning resource.
                <br><small>${error.message}</small>
            </div>
        `;
    }
}

function renderResourceContent(metadata, markdownContent) {
    const contentContainer = document.getElementById('contentContainer');
    
    // Parse markdown to HTML
    const htmlContent = marked.parse(markdownContent);
    
    const languageClass = metadata.language === 'java' ? 'java-badge' : 'python-badge';
    
    contentContainer.innerHTML = `
        <div class="action-buttons d-flex justify-content-between align-items-center">
            <div>
                <button class="btn btn-sm btn-outline-secondary" onclick="showOutline()">
                    ← Back to Outline
                </button>
                <span class="language-badge ${languageClass} ms-2">${metadata.language.toUpperCase()}</span>
                <span class="text-muted ms-2">Phase ${metadata.phase}</span>
            </div>
            <div>
                <button class="btn btn-sm btn-primary" onclick="takeQuiz('${metadata.language}', ${metadata.phase})">
                    Take Quiz
                </button>
            </div>
        </div>
        <div class="markdown-content">
            ${htmlContent}
        </div>
    `;

    // Apply syntax highlighting to code blocks
    if (typeof hljs !== 'undefined') {
        contentContainer.querySelectorAll('pre code').forEach((block) => {
            hljs.highlightElement(block);
        });
    }

    // Scroll to top
    contentContainer.scrollTop = 0;
}

async function trackView(resourceId) {
    try {
        await fetch('/learning/track-view', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                resource_id: resourceId,
                completion_percentage: 50 // For now, just mark as 50% on view
            })
        });
    } catch (error) {
        console.error('Error tracking view:', error);
    }
}

async function searchResources(query) {
    const container = document.getElementById('outlineContainer');
    container.innerHTML = `
        <div class="loading-spinner">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Searching...</span>
            </div>
        </div>
    `;

    try {
        const langParam = currentLanguage === 'both' ? '' : `&language=${currentLanguage}`;
        const response = await fetch(`/learning/search?q=${encodeURIComponent(query)}${langParam}`);
        
        if (!response.ok) {
            throw new Error('Search failed');
        }

        const data = await response.json();
        renderSearchResults(data.results, query);
    } catch (error) {
        console.error('Error searching:', error);
        container.innerHTML = `
            <div class="alert alert-danger">
                <strong>Error!</strong> Search failed.
            </div>
        `;
    }
}

function renderSearchResults(results, query) {
    const container = document.getElementById('outlineContainer');
    
    if (!results || results.length === 0) {
        container.innerHTML = `
            <div class="alert alert-info">
                No results found for "${query}".
                <button class="btn btn-sm btn-link" onclick="loadOutline()">Clear search</button>
            </div>
        `;
        return;
    }

    let html = `
        <div class="mb-3">
            <strong>Search Results:</strong> ${results.length} found
            <button class="btn btn-sm btn-link float-end" onclick="loadOutline()">Clear</button>
        </div>
        <ul class="resource-list">
    `;
    
    results.forEach(resource => {
        const badgeClass = resource.language === 'java' ? 'java-badge' : 'python-badge';
        const activeClass = currentResourceId === resource.id ? 'active' : '';
        
        html += `
            <li class="resource-item ${activeClass}" onclick="loadResource(${resource.id})">
                <span class="language-badge ${badgeClass}">${resource.language.toUpperCase()}</span>
                <span>${resource.section_number}. ${resource.section_title}</span>
                <small class="text-muted">(Phase ${resource.phase})</small>
            </li>
        `;
    });
    
    html += '</ul>';
    container.innerHTML = html;
}

function showOutline() {
    currentResourceId = null;
    const contentContainer = document.getElementById('contentContainer');
    contentContainer.innerHTML = `
        <div class="welcome-screen">
            <div style="font-size: 80px; margin-bottom: 20px;">📚</div>
            <h2>Welcome to the Learning Center</h2>
            <p class="lead">Select a topic from the outline to begin studying.</p>
            <p class="text-muted">Your progress is automatically tracked as you explore the materials.</p>
        </div>
    `;
    loadOutline();
}

function takeQuiz(language, phase) {
    // Navigate to quiz page with topic filter
    // This is a placeholder - you'll need to implement based on your quiz system
    window.location.href = `/quiz.html?language=${language}&phase=${phase}`;
}

function logout() {
    // Clear session and redirect to login
    localStorage.removeItem('user');
    window.location.href = '/index.html';
}

// Export functions for use in other contexts
window.loadResourceFromQuiz = loadResource;
