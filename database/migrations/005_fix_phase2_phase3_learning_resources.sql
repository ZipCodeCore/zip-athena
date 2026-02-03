-- Fix Phase 2 and Phase 3 Learning Resource File Paths and Activate Them
-- Migration: 005
-- This migration corrects file paths for Phase 2 and 3 resources and activates them

-- Java Phase 2: Update file paths and activate
UPDATE learning_resources SET file_path = 'java/8-oop-classes-and-objects.md', is_active = 1 
WHERE language = 'java' AND phase = 2 AND section_number = '8';

UPDATE learning_resources SET file_path = 'java/9-oop-inheritance.md', is_active = 1 
WHERE language = 'java' AND phase = 2 AND section_number = '9';

UPDATE learning_resources SET file_path = 'java/10-oop-polymorphism.md', is_active = 1 
WHERE language = 'java' AND phase = 2 AND section_number = '10';

UPDATE learning_resources SET is_active = 1 
WHERE language = 'java' AND phase = 2 AND section_number = '11';

-- Java Phase 3: Update file paths and activate
UPDATE learning_resources SET is_active = 1 
WHERE language = 'java' AND phase = 3 AND section_number = '12';

UPDATE learning_resources SET is_active = 1 
WHERE language = 'java' AND phase = 3 AND section_number = '13';

UPDATE learning_resources SET is_active = 1 
WHERE language = 'java' AND phase = 3 AND section_number = '14';

UPDATE learning_resources SET file_path = 'java/15-collections-framework.md', is_active = 1 
WHERE language = 'java' AND phase = 3 AND section_number = '15';

-- Python Phase 2: Update file paths and activate
UPDATE learning_resources SET file_path = 'python/14-oop-classes-and-objects.md', is_active = 1 
WHERE language = 'python' AND phase = 2 AND section_number = '14';

UPDATE learning_resources SET file_path = 'python/15-oop-inheritance.md', is_active = 1 
WHERE language = 'python' AND phase = 2 AND section_number = '15';

UPDATE learning_resources SET file_path = 'python/16-oop-polymorphism.md', is_active = 1 
WHERE language = 'python' AND phase = 2 AND section_number = '16';

-- Python Phase 3: Update file paths and activate
UPDATE learning_resources SET file_path = 'python/11-modules-and-packages.md', is_active = 1 
WHERE language = 'python' AND phase = 3 AND section_number = '11';

UPDATE learning_resources SET is_active = 1 
WHERE language = 'python' AND phase = 3 AND section_number = '12';

UPDATE learning_resources SET is_active = 1 
WHERE language = 'python' AND phase = 3 AND section_number = '13';
