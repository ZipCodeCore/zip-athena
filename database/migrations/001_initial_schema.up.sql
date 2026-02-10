-- Initial schema migration
-- Consolidates all base tables from inline schema

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT 0,
    category TEXT DEFAULT 'java',
    total_points INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS topics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    category TEXT DEFAULT 'java',
    is_active BOOLEAN DEFAULT 1,
    display_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic_id INTEGER NOT NULL,
    question_text TEXT NOT NULL,
    difficulty INTEGER NOT NULL,
    tags TEXT,
    topic_question_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES topics(id)
);

CREATE TABLE IF NOT EXISTS answers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    question_id INTEGER NOT NULL,
    answer_text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    FOREIGN KEY (question_id) REFERENCES questions(id)
);

CREATE TABLE IF NOT EXISTS quiz_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    quiz_id INTEGER,
    parent_attempt_id INTEGER,
    attempt_number INTEGER DEFAULT 1,
    is_retake BOOLEAN DEFAULT 0,
    start_time DATETIME NOT NULL,
    completed_time DATETIME,
    total_questions INTEGER DEFAULT 0,
    correct_answers INTEGER DEFAULT 0,
    difficulty_average REAL DEFAULT 0,
    description TEXT DEFAULT '',
    question_ids TEXT,
    question_results TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id),
    FOREIGN KEY (parent_attempt_id) REFERENCES quiz_attempts(id)
);

CREATE TABLE IF NOT EXISTS question_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    topic_question_id TEXT NOT NULL,
    correct_attempts INTEGER DEFAULT 0,
    incorrect_attempts INTEGER DEFAULT 0,
    last_attempt_time DATETIME NOT NULL,
    last_attempt_correct BOOLEAN NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (question_id) REFERENCES questions(id)
);

CREATE TABLE IF NOT EXISTS quizzes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    topic_id INTEGER NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES topics(id)
);

CREATE TABLE IF NOT EXISTS user_question_states (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    due DATETIME NOT NULL,
    stability REAL NOT NULL,
    difficulty REAL NOT NULL,
    elapsed_days INTEGER NOT NULL,
    scheduled_days INTEGER NOT NULL,
    reps INTEGER NOT NULL,
    lapses INTEGER NOT NULL,
    state INTEGER NOT NULL,
    last_review DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (question_id) REFERENCES questions(id),
    UNIQUE(user_id, question_id)
);

CREATE TABLE IF NOT EXISTS tests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    topic_id INTEGER NOT NULL,
    description TEXT,
    time_limit INTEGER DEFAULT 0,
    start_date DATETIME,
    end_date DATETIME,
    is_published INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES topics(id)
);

CREATE TABLE IF NOT EXISTS test_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    test_id INTEGER NOT NULL,
    start_time DATETIME NOT NULL,
    submit_time DATETIME,
    time_remaining INTEGER DEFAULT 0,
    is_auto_submit INTEGER DEFAULT 0,
    total_questions INTEGER DEFAULT 0,
    correct_answers INTEGER DEFAULT 0,
    score REAL DEFAULT 0.0,
    answers TEXT DEFAULT '{}',
    question_results TEXT DEFAULT '{}',
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS quiz_questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    quiz_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    question_order INTEGER NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    UNIQUE(quiz_id, question_id)
);

CREATE TABLE IF NOT EXISTS test_questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    test_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    question_order INTEGER NOT NULL,
    FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    UNIQUE(test_id, question_id)
);

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

CREATE TABLE IF NOT EXISTS flagged_questions (
    question_id INTEGER PRIMARY KEY,
    flag_reason TEXT NOT NULL,
    flagged_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    accuracy REAL,
    total_attempts INTEGER,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
);
