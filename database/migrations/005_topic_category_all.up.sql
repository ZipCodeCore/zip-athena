-- Migrate topic categories to the three-value scheme: java, python, all
-- sql topics become "all" (visible to every student category)
UPDATE topics SET category = 'all' WHERE category = 'sql';

-- java-tracing topics become "java"
UPDATE topics SET category = 'java' WHERE category = 'java-tracing';

-- python-tracing topics become "python"
UPDATE topics SET category = 'python' WHERE category = 'python-tracing';
