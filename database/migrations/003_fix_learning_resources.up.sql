-- Activate phases 2 and 3 learning resources
-- Fix file_path values that don't match actual filenames on disk

UPDATE learning_resources SET is_active = 1 WHERE phase IN (2, 3) AND is_active = 0;

UPDATE learning_resources SET file_path = 'java/8-oop-classes-and-objects.md' WHERE file_path = 'java/8-classes-and-objects.md';
UPDATE learning_resources SET file_path = 'java/9-oop-inheritance.md' WHERE file_path = 'java/9-inheritance.md';
UPDATE learning_resources SET file_path = 'java/10-oop-polymorphism.md' WHERE file_path = 'java/10-polymorphism.md';
UPDATE learning_resources SET file_path = 'java/15-collections-framework.md' WHERE file_path = 'java/15-collections.md';
UPDATE learning_resources SET file_path = 'python/11-modules-and-packages.md' WHERE file_path = 'python/11-modules-packages.md';
UPDATE learning_resources SET file_path = 'python/14-oop-classes-and-objects.md' WHERE file_path = 'python/14-classes-and-objects.md';
UPDATE learning_resources SET file_path = 'python/15-oop-inheritance.md' WHERE file_path = 'python/15-inheritance.md';
UPDATE learning_resources SET file_path = 'python/16-oop-polymorphism.md' WHERE file_path = 'python/16-polymorphism.md';
