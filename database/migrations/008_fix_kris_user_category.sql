-- Migration: Fix Kris user to be a data category student (not admin)
-- This migration fixes the issue where Kris was incorrectly set as admin with java category
-- After this fix, Kris will only see Python and SQL topics (data category)

UPDATE users 
SET is_admin = 0, 
    category = 'data'
WHERE username = 'kris';
