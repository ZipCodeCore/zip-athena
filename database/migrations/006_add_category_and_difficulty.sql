-- Add category columns and difficulty average
-- Migration: 006
-- This migration adds category support for users and topics, and difficulty tracking for quiz attempts

-- Add category column to users table (default to 'java' for existing users)
ALTER TABLE users ADD COLUMN category TEXT DEFAULT 'java';

-- Add category column to topics table (default to 'java' for existing topics)
ALTER TABLE topics ADD COLUMN category TEXT DEFAULT 'java';

-- Add difficulty_average column to quiz_attempts table
ALTER TABLE quiz_attempts ADD COLUMN difficulty_average REAL DEFAULT 0;
