-- Seed Learning Resources
-- This file populates the learning_resources table with the current learning materials

-- Java Learning Resources

-- Phase 1: Core Fundamentals
INSERT INTO learning_resources (language, phase, section_number, section_title, file_path, display_order, is_active) VALUES
('java', 1, '2', 'Variables and Data Types', 'java/2-variables-and-data-types.md', 1, 1),
('java', 1, '3', 'Operators and Expressions', 'java/3-operators-and-expressions.md', 2, 1),
('java', 1, '4', 'Control Structures', 'java/4-control-structures.md', 3, 1),
('java', 1, '5', 'Methods (Functions)', 'java/5-methods.md', 4, 1),
('java', 1, '6', 'Arrays', 'java/6-arrays.md', 5, 1),
('java', 1, '7', 'Strings', 'java/7-strings.md', 6, 1);

-- Phase 2: Object-Oriented Programming
INSERT INTO learning_resources (language, phase, section_number, section_title, file_path, display_order, is_active) VALUES
('java', 2, '8', 'OOP - Classes and Objects', 'java/8-oop-classes-and-objects.md', 1, 1),
('java', 2, '9', 'OOP - Inheritance', 'java/9-oop-inheritance.md', 2, 1),
('java', 2, '10', 'OOP - Polymorphism', 'java/10-oop-polymorphism.md', 3, 1),
('java', 2, '11', 'Interfaces', 'java/11-interfaces.md', 4, 1);

-- Phase 3: Advanced Data Handling
INSERT INTO learning_resources (language, phase, section_number, section_title, file_path, display_order, is_active) VALUES
('java', 3, '12', 'Packages', 'java/12-packages.md', 1, 1),
('java', 3, '13', 'Exception Handling', 'java/13-exception-handling.md', 2, 1),
('java', 3, '14', 'File I/O', 'java/14-file-io.md', 3, 1),
('java', 3, '15', 'Collections Framework', 'java/15-collections-framework.md', 4, 1);

-- Phase 4: Advanced Topics
INSERT INTO learning_resources (language, phase, section_number, section_title, file_path, display_order, is_active) VALUES
('java', 4, '16', 'Generics - Using Generics', 'java/16-generics.md', 1, 0),
('java', 4, '17', 'Basic Multithreading Concepts', 'java/17-multithreading.md', 2, 0),
('java', 4, '18', 'Java Standard Library Essentials', 'java/18-standard-library.md', 3, 0),
('java', 4, '19', 'Lambda Expressions and Streams', 'java/19-lambda-streams.md', 4, 0);

-- Phase 5: Specialized Topics
INSERT INTO learning_resources (language, phase, section_number, section_title, file_path, display_order, is_active) VALUES
('java', 5, '20', 'Best Practices and Code Quality', 'java/20-best-practices.md', 1, 0);


-- Python Learning Resources

-- Phase 1: Core Fundamentals
INSERT INTO learning_resources (language, phase, section_number, section_title, file_path, display_order, is_active) VALUES
('python', 1, '2', 'Variables and Data Types', 'python/2-variables-and-data-types.md', 1, 1),
('python', 1, '3', 'Operators and Expressions', 'python/3-operators-and-expressions.md', 2, 1),
('python', 1, '4', 'Strings', 'python/4-strings.md', 3, 1),
('python', 1, '5', 'Control Structures', 'python/5-control-structures.md', 4, 1),
('python', 1, '6', 'Data Structures', 'python/6-data-structures.md', 5, 1),
('python', 1, '10', 'Functions', 'python/10-functions.md', 6, 1);

-- Phase 2: Object-Oriented Programming
INSERT INTO learning_resources (language, phase, section_number, section_title, file_path, display_order, is_active) VALUES
('python', 2, '14', 'OOP - Classes and Objects', 'python/14-oop-classes-and-objects.md', 1, 1),
('python', 2, '15', 'OOP - Inheritance', 'python/15-oop-inheritance.md', 2, 1),
('python', 2, '16', 'OOP - Polymorphism', 'python/16-oop-polymorphism.md', 3, 1);

-- Phase 3: Advanced Data Handling
INSERT INTO learning_resources (language, phase, section_number, section_title, file_path, display_order, is_active) VALUES
('python', 3, '11', 'Modules and Packages', 'python/11-modules-and-packages.md', 1, 1),
('python', 3, '12', 'File I/O', 'python/12-file-io.md', 2, 1),
('python', 3, '13', 'Exception Handling', 'python/13-exception-handling.md', 3, 1);

-- Phase 4: Advanced Topics
INSERT INTO learning_resources (language, phase, section_number, section_title, file_path, display_order, is_active) VALUES
('python', 4, '17', 'Working with Mixins', 'python/17-mixins.md', 1, 0),
('python', 4, '18', 'Iterators and Generators', 'python/18-iterators-generators.md', 2, 0),
('python', 4, '19', 'Comprehensions and Functional Programming', 'python/19-comprehensions-functional.md', 3, 0),
('python', 4, '20', 'Decorators - Using Decorators', 'python/20-decorators.md', 4, 0);

-- Phase 5: Specialized Topics
INSERT INTO learning_resources (language, phase, section_number, section_title, file_path, display_order, is_active) VALUES
('python', 5, '21', 'Context Managers', 'python/21-context-managers.md', 1, 0),
('python', 5, '22', 'Regular Expressions', 'python/22-regex.md', 2, 0),
('python', 5, '23', 'Working with Dates and Times', 'python/23-datetime.md', 3, 0),
('python', 5, '24', 'Common Standard Library Modules', 'python/24-stdlib-modules.md', 4, 0),
('python', 5, '25', 'Testing Basics', 'python/25-testing.md', 5, 0),
('python', 5, '26', 'Virtual Environments and Package Management', 'python/26-venv-packages.md', 6, 0),
('python', 5, '27', 'Best Practices and Code Quality', 'python/27-best-practices.md', 7, 0);
