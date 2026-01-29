# Learning Resource Format and Structure

This document defines the standard format for all learning resources in Athena. Each learning resource combines educational content with quiz questions to create a complete learning experience.

## Purpose

Learning resources in Athena serve dual purposes:
1. **Study Material**: Teach students programming concepts through clear explanations and examples
2. **Quiz Generation**: Provide questions that can be used to test understanding of the concepts

This dual-purpose structure ensures that:
- Study guides and quiz questions stay synchronized
- Questions accurately reflect the material taught
- Updates to concepts automatically update related questions
- Students can study first, then immediately test their knowledge

## Directory Structure

Resources are organized hierarchically:

```
learning/
├── README.md                           # Overview and navigation
├── PHASED_PLAN.md                      # Development roadmap
├── LEARNING_RESOURCE_FORMAT.md         # This file
├── {language}/                         # Language tag (java, python, etc.)
│   └── {topic-name}.md                 # Topic file (e.g., control-structures.md)
```

### Hierarchy Explained

- **Language Tag**: Top-level organization (Java, Python, Go, etc.)
- **Topic**: Major curriculum section (e.g., "Control Structures", "Data Types")
- **Sub-topics**: Specific concepts within a topic (e.g., "if statement", "for loop")

## File Structure

Each learning resource file follows this structure:

```markdown
# [Language] [Topic Name]

[Brief introduction paragraph explaining the topic and its importance]

## [Sub-topic 1]

[2-4 paragraphs of friendly, conversational explanation]

### Examples

```[language]
[Working code example]
```

[Explanation of the example]

## [Sub-topic 2]

[Continue with more sub-topics...]

## Common Pitfalls

[Common mistakes and how to avoid them]

## Best Practices

[Tips for using these concepts effectively]

## Practice Exercises

[Optional: Hands-on exercises for students to try]

---

## Quiz Questions

This section contains quiz questions that test understanding of the concepts taught above. Each question is formatted for easy extraction and database import.

### Question 1: [Sub-topic Name] - [Brief Description]
**Type**: [multiple-choice|true-false|code-completion|code-output]
**Difficulty**: [1-7]
**Tags**: [language], [topic], [sub-topic-1], [sub-topic-2]

**Question Text**:
[The actual question, formatted clearly]

```[language]
[Code snippet if applicable]
```

**Answer Choices** (for multiple-choice):
- A) [First option]
- B) [Second option]
- C) [Third option]
- D) [Fourth option]

**Correct Answer**: [A|B|C|D or the correct answer text]

**Explanation**:
[Clear explanation of why this is the correct answer and why others are wrong]

---

### Question 2: [Continue pattern for 6-8 questions total]
...
```

## Question Format Details

### Question Types

1. **multiple-choice**: Student selects one answer from 4 options (A, B, C, D)
2. **true-false**: Simple true/false question
3. **code-completion**: Student completes missing code
4. **code-output**: Student predicts what code will output

### Difficulty Levels (1-7)

- **Level 1-2**: Basic concepts, direct recall
- **Level 3-4**: Understanding and application
- **Level 5-6**: Analysis and synthesis
- **Level 7**: Complex problem-solving

### Tags

Tags help organize and filter questions:
- **Language tag**: java, python, go, etc.
- **Topic tag**: control-structures, data-types, oop, etc.
- **Concept tags**: if-statement, loops, variables, etc.

Example: `java, control-structures, for-loop, nested-loops`

### Best Practices for Questions

1. **Alignment**: Questions should directly test concepts explained in the study guide
2. **Progressive Difficulty**: Start with easier questions (1-3), build to harder ones (5-7)
3. **Variety**: Mix question types to test different aspects of understanding
4. **Real-world Context**: Use practical scenarios when possible
5. **Clear Explanations**: Always explain the correct answer thoroughly
6. **Common Mistakes**: Design wrong answers based on common student errors

## Example Learning Resource

Here's a minimal example of a properly formatted learning resource:

````markdown
# Python Variables and Data Types

Variables are named containers that store values in your program. Think of them as labeled boxes where you can put information, look at it later, or replace it with something new.

## Basic Data Types

Python has several built-in data types that you'll use constantly.

### Integers (int)

Integers are whole numbers without decimal points. They can be positive, negative, or zero.

```python
age = 25
temperature = -5
count = 0
```

You can perform arithmetic operations with integers:

```python
x = 10
y = 3
print(x + y)   # 13
print(x - y)   # 7
print(x * y)   # 30
print(x / y)   # 3.3333...
print(x // y)  # 3 (integer division)
print(x % y)   # 1 (remainder)
```

## Common Pitfalls

- Forgetting that `/` returns a float even with integers
- Using `==` instead of `=` for assignment

## Best Practices

- Use descriptive variable names like `student_count` instead of `sc`
- Initialize variables before using them

---

## Quiz Questions

### Question 1: Integer Operations - Division Types
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: python, data-types, integers, operators

**Question Text**:
What is the output of the following Python code?

```python
result = 17 // 5
print(result)
```

**Answer Choices**:
- A) 3.4
- B) 3
- C) 4
- D) 3.0

**Correct Answer**: B

**Explanation**:
The `//` operator performs integer (floor) division, which returns the quotient without the remainder. 17 divided by 5 is 3 with a remainder of 2, so `17 // 5` returns `3`. Option A would be the result of regular division (`17 / 5`). Options C and D are incorrect.

---

### Question 2: Variable Assignment - Multiple Assignment
**Type**: code-output
**Difficulty**: 2
**Tags**: python, variables, assignment

**Question Text**:
What will be printed by this code?

```python
x, y, z = 1, 2, 3
print(x + y + z)
```

**Answer Choices**:
- A) 123
- B) 6
- C) 1 2 3
- D) Error

**Correct Answer**: B

**Explanation**:
Python allows multiple assignment in a single line. This code assigns `x=1`, `y=2`, and `z=3`, then prints their sum: `1 + 2 + 3 = 6`. Option A would be string concatenation. Option C would require printing with commas. There's no error.

````

## Integration with Athena Application

### Current State (Phase 0)

Learning resources are standalone markdown files in the `learning/` folder. They serve as reference material for students.

### Future Integration

The structure is designed to support:

1. **Manual Question Import**: Instructors can copy quiz questions from learning resources into the database
2. **Automatic Parsing**: Future tools could parse the "Quiz Questions" section and import questions automatically
3. **Dynamic Linking**: Questions in the database could reference their source learning resource
4. **Personalized Learning**: Students could study a topic, then take a quiz with questions from that topic
5. **Spaced Repetition**: Wrong answers could trigger review of specific sub-topics from the learning resource

## Updating Existing Resources

To update existing learning resources to this format:

1. Keep all existing study guide content (it's already well-written)
2. Add a horizontal rule (`---`) after the main content
3. Add a `## Quiz Questions` section
4. Create 6-8 questions following the format above
5. Ensure questions cover all major sub-topics from the study guide

## Quality Guidelines

### Study Guide Section

- **Tone**: Friendly, conversational, encouraging
- **Length**: 2-4 paragraphs per sub-topic
- **Examples**: Include working code for every concept
- **Clarity**: Assume student is learning for the first time
- **Completeness**: Cover all aspects of the topic from CURRICULUM.md

### Quiz Questions Section

- **Quantity**: 6-8 questions per topic
- **Coverage**: Questions should cover all major sub-topics
- **Balance**: Mix of easier and harder questions
- **Accuracy**: Verify all code examples run correctly
- **Explanations**: Every answer needs a clear explanation

## Maintenance and Updates

- Learning resources are version controlled in Git
- Updates to curriculum should trigger review of related resources
- Student feedback should inform improvements
- Quiz questions should be tested for clarity and difficulty
- Code examples should be verified with each language version update

## Contributing

When creating new learning resources:

1. Follow the format exactly as specified
2. Review CURRICULUM.md for topic coverage requirements
3. Match the tone and style of existing resources
4. Test all code examples
5. Have another developer review before finalizing
6. Start with the most important topics from PHASED_PLAN.md

---

This format ensures consistency, maintainability, and future extensibility of Athena's learning resources while supporting both human learners and potential automated systems.
