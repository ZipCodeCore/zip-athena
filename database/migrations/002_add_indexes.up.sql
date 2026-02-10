-- Performance indexes for learning resources and related tables

CREATE INDEX IF NOT EXISTS idx_learning_resources_language ON learning_resources(language);
CREATE INDEX IF NOT EXISTS idx_learning_resources_phase ON learning_resources(phase);
CREATE INDEX IF NOT EXISTS idx_learning_resources_active ON learning_resources(is_active);
CREATE INDEX IF NOT EXISTS idx_user_learning_progress_user ON user_learning_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_user_learning_progress_resource ON user_learning_progress(learning_resource_id);
CREATE INDEX IF NOT EXISTS idx_question_learning_links_question ON question_learning_links(question_id);
CREATE INDEX IF NOT EXISTS idx_question_learning_links_resource ON question_learning_links(learning_resource_id);
