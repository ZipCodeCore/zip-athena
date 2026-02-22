-- Rollback assigned_quizzes migration

DROP TABLE IF EXISTS assigned_quizzes;

-- Note: SQLite 3.35+ supports ALTER TABLE DROP COLUMN
-- If running an older SQLite, this step can be skipped as the extra column is harmless
ALTER TABLE quizzes DROP COLUMN IF EXISTS category;
