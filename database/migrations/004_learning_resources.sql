-- Migration: Learning Resources and Integration
-- This migration adds support for learning materials within Athena

-- Learning Resources table
CREATE TABLE IF NOT EXISTS learning_resources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    language TEXT NOT NULL CHECK(language IN ('java', 'python')),
    phase INTEGER NOT NULL CHECK(phase BETWEEN 1 AND 5),
    section_number TEXT NOT NULL,
    section_title TEXT NOT NULL,
    file_path TEXT NOT NULL,
    display_order INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(language, section_number)
);

-- Link questions to learning resources
CREATE TABLE IF NOT EXISTS question_learning_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    question_id INTEGER NOT NULL,
    learning_resource_id INTEGER NOT NULL,
    relevance_score REAL DEFAULT 1.0 CHECK(relevance_score >= 0 AND relevance_score <= 1),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    FOREIGN KEY (learning_resource_id) REFERENCES learning_resources(id) ON DELETE CASCADE,
    UNIQUE(question_id, learning_resource_id)
);

-- Track which resources students have viewed
CREATE TABLE IF NOT EXISTS user_learning_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    learning_resource_id INTEGER NOT NULL,
    first_viewed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_viewed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    view_count INTEGER DEFAULT 1,
    completion_percentage REAL DEFAULT 0 CHECK(completion_percentage >= 0 AND completion_percentage <= 100),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (learning_resource_id) REFERENCES learning_resources(id) ON DELETE CASCADE,
    UNIQUE(user_id, learning_resource_id)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_learning_resources_language ON learning_resources(language);
CREATE INDEX IF NOT EXISTS idx_learning_resources_phase ON learning_resources(phase);
CREATE INDEX IF NOT EXISTS idx_learning_resources_active ON learning_resources(is_active);
CREATE INDEX IF NOT EXISTS idx_question_learning_links_question ON question_learning_links(question_id);
CREATE INDEX IF NOT EXISTS idx_question_learning_links_resource ON question_learning_links(learning_resource_id);
CREATE INDEX IF NOT EXISTS idx_user_learning_progress_user ON user_learning_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_user_learning_progress_resource ON user_learning_progress(learning_resource_id);
