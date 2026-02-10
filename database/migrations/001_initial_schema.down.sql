-- Rollback initial schema
-- Drops all tables created in 001_initial_schema.up.sql

DROP TABLE IF EXISTS question_learning_links;
DROP TABLE IF EXISTS user_learning_progress;
DROP TABLE IF EXISTS learning_resources;
DROP TABLE IF EXISTS test_questions;
DROP TABLE IF EXISTS quiz_questions;
DROP TABLE IF EXISTS test_attempts;
DROP TABLE IF EXISTS tests;
DROP TABLE IF EXISTS user_question_states;
DROP TABLE IF EXISTS quizzes;
DROP TABLE IF EXISTS question_attempts;
DROP TABLE IF EXISTS quiz_attempts;
DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS topics;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS flagged_questions;
