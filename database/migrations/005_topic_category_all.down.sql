-- Revert topic category migration (best-effort; lossy for custom topics)
-- This rollback can only reliably restore topics that were seeded by the application.
-- Any "all" category topics created by users after migration 005 will remain "all".
-- The original topic names are used to identify previously-migrated rows.
UPDATE topics SET category = 'sql'           WHERE category = 'all'    AND name = 'SQL Database';
UPDATE topics SET category = 'java-tracing'  WHERE category = 'java'   AND name = 'Java Code Tracing';
UPDATE topics SET category = 'python-tracing' WHERE category = 'python' AND name = 'Python Code Tracing';
