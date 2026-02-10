-- Rollback performance indexes

DROP INDEX IF EXISTS idx_question_learning_links_resource;
DROP INDEX IF EXISTS idx_question_learning_links_question;
DROP INDEX IF EXISTS idx_user_learning_progress_resource;
DROP INDEX IF EXISTS idx_user_learning_progress_user;
DROP INDEX IF EXISTS idx_learning_resources_active;
DROP INDEX IF EXISTS idx_learning_resources_phase;
DROP INDEX IF EXISTS idx_learning_resources_language;
