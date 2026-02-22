-- Add category to quizzes table and create assigned_quizzes table

ALTER TABLE quizzes ADD COLUMN category TEXT DEFAULT '';

CREATE TABLE IF NOT EXISTS assigned_quizzes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    quiz_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    assigned_by INTEGER NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by) REFERENCES users(id),
    UNIQUE(quiz_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_assigned_quizzes_user ON assigned_quizzes(user_id);
CREATE INDEX IF NOT EXISTS idx_assigned_quizzes_quiz ON assigned_quizzes(quiz_id);
