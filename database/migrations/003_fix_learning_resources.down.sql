-- Rollback learning resources fixes
-- Note: This is a best-effort rollback. If data has been modified in other ways, 
-- this may not restore the exact original state.

UPDATE learning_resources SET file_path = 'java/8-classes-and-objects.md' WHERE file_path = 'java/8-oop-classes-and-objects.md';
UPDATE learning_resources SET file_path = 'java/9-inheritance.md' WHERE file_path = 'java/9-oop-inheritance.md';
UPDATE learning_resources SET file_path = 'java/10-polymorphism.md' WHERE file_path = 'java/10-oop-polymorphism.md';
UPDATE learning_resources SET file_path = 'java/15-collections.md' WHERE file_path = 'java/15-collections-framework.md';
UPDATE learning_resources SET file_path = 'python/11-modules-packages.md' WHERE file_path = 'python/11-modules-and-packages.md';
UPDATE learning_resources SET file_path = 'python/14-classes-and-objects.md' WHERE file_path = 'python/14-oop-classes-and-objects.md';
UPDATE learning_resources SET file_path = 'python/15-inheritance.md' WHERE file_path = 'python/15-oop-inheritance.md';
UPDATE learning_resources SET file_path = 'python/16-polymorphism.md' WHERE file_path = 'python/16-oop-polymorphism.md';
