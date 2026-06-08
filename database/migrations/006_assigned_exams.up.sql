-- Add assigned exams feature
-- Exams are distinct from quizzes: they are generated from random questions
-- filtered by category and difficulty level.

CREATE TABLE IF NOT EXISTS exams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    categories TEXT NOT NULL DEFAULT '[]',       -- JSON array of category strings
    difficulty_levels TEXT NOT NULL DEFAULT '[]', -- JSON array of ints (1=easy,2=medium,3=hard)
    question_ids TEXT NOT NULL DEFAULT '[]',      -- JSON array of question IDs selected at generation time
    question_count INTEGER NOT NULL DEFAULT 30,
    access_token TEXT NOT NULL UNIQUE,            -- Unique token for exam URL
    created_by INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS exam_assignments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    exam_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    assigned_by INTEGER NOT NULL,
    FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by) REFERENCES users(id),
    UNIQUE(exam_id, user_id)
);

CREATE TABLE IF NOT EXISTS exam_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    exam_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    total_questions INTEGER DEFAULT 0,
    correct_answers INTEGER DEFAULT 0,
    score REAL DEFAULT 0.0,
    answers TEXT DEFAULT '{}',         -- JSON: question_id -> answer_id
    question_results TEXT DEFAULT '{}', -- JSON: question_id -> bool (correct/incorrect)
    FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_exams_token ON exams(access_token);
CREATE INDEX IF NOT EXISTS idx_exam_assignments_exam ON exam_assignments(exam_id);
CREATE INDEX IF NOT EXISTS idx_exam_assignments_user ON exam_assignments(user_id);
CREATE INDEX IF NOT EXISTS idx_exam_attempts_exam ON exam_attempts(exam_id);
CREATE INDEX IF NOT EXISTS idx_exam_attempts_user ON exam_attempts(user_id);
