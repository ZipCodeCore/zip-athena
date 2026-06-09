CREATE TABLE users_rollback (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT 0,
    category TEXT DEFAULT 'java',
    total_points INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users_rollback (id, username, email, password, is_admin, category, total_points, created_at)
SELECT id, username, email, password, is_admin, category, total_points, created_at
FROM users;

DROP TABLE users;

ALTER TABLE users_rollback RENAME TO users;
