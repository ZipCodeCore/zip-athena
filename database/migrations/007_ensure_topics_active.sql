-- Migration 007: Ensure all existing topics are marked as active
-- This migration ensures that all topics in the database are marked as is_active=1
-- to fix the issue where students were not seeing topics on their dashboard

-- Update all topics to be active (this is safe because we want all existing topics to be visible)
UPDATE topics SET is_active = 1 WHERE is_active = 0 OR is_active IS NULL;

-- Ensure the is_active column has a default value for new topics
-- Note: This should already be in the schema, but we're being explicit
-- SQLite doesn't support ALTER COLUMN DEFAULT directly, so we skip this if it fails
