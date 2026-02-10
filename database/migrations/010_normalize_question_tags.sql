-- Migration 010: Create question_tags table to normalize JSON field
-- This replaces the tags TEXT (JSON array) column in questions

CREATE TABLE IF NOT EXISTS question_tags (
    question_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (question_id, tag),
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_question_tags_tag ON question_tags(tag);
CREATE INDEX IF NOT EXISTS idx_question_tags_question ON question_tags(question_id);
