# Contributor Guide for Learning Resources

Welcome! This guide will help you create high-quality learning resources for the Athena application. Each resource combines study material with quiz questions to provide a complete learning experience.

## Quick Start

1. Read [LEARNING_RESOURCE_FORMAT.md](LEARNING_RESOURCE_FORMAT.md) for the complete specification
2. Review existing resources for style and tone examples:
   - `java/control-structures.md` - Complete with quiz questions
   - `java/operators-and-expressions-EXAMPLE.md` - Template example
3. Choose a topic from [PHASED_PLAN.md](PHASED_PLAN.md)
4. Create your resource following the format
5. Have another contributor review your work

## Creating a New Learning Resource

### Step 1: Choose a Topic

Check [PHASED_PLAN.md](PHASED_PLAN.md) for topics that need resources. Prioritize:
- **Phase 1** topics (highest priority)
- Topics you're most knowledgeable about
- Topics that complement existing resources

### Step 2: Research and Plan

Before writing:
1. Review the corresponding section in [CURRICULUM.md](/CURRICULUM.md)
2. List all sub-topics that need coverage
3. Think about common student misconceptions
4. Gather good code examples

### Step 3: Write the Study Guide

Follow these principles:

**Tone and Style**:
- Friendly and conversational (like talking to a friend)
- Encouraging and supportive
- Clear and direct
- No condescension or overly formal language

**Structure**:
- Start with a welcoming introduction
- Break content into logical sub-sections
- Each sub-section: 2-4 paragraphs
- Include working code examples for every concept
- End sections with practical tips

**Example Opening**:
```markdown
# Java Methods

Welcome! Let's talk about methods in Java. Think of methods as reusable recipes for your code—you write the instructions once, give it a name, and then you can use it over and over.
```

**Example Section**:
```markdown
## Method Parameters

Parameters let you pass information into a method. They're like ingredients in a recipe—different inputs can produce different results, but the process stays the same.

```java
public void greet(String name) {
    System.out.println("Hello, " + name + "!");
}

greet("Alice");  // Prints: Hello, Alice!
greet("Bob");    // Prints: Hello, Bob!
```

The `String name` in the parentheses is a parameter. When you call the method, you provide an argument ("Alice" or "Bob"), which becomes the value of `name` inside the method.

**Pro tip**: Parameter names should be descriptive. Use `userName` instead of `u` or `x`.
```

### Step 4: Write Quiz Questions

Create 6-8 questions that:
- Cover all major sub-topics from your study guide
- Range in difficulty from basic (2-3) to intermediate (5-6)
- Mix question types (multiple-choice, code-output, true-false)
- Include thorough explanations

**Question Template**:
```markdown
### Question X: [Sub-topic] - [Brief Description]
**Type**: [multiple-choice|true-false|code-completion|code-output]
**Difficulty**: [1-7]
**Tags**: [language], [topic], [concept-1], [concept-2]

**Question Text**:
[Clear, specific question. Include code if relevant.]

```java
[Code example if applicable]
```

**Answer Choices**:
- A) [First option]
- B) [Second option]
- C) [Third option]
- D) [Fourth option]

**Correct Answer**: [A|B|C|D]

**Explanation**:
[Why the correct answer is right AND why the others are wrong. Should teach, not just tell.]
```

### Step 5: Quality Checks

Before submitting your resource:

**Study Guide**:
- [ ] All code examples run without errors
- [ ] Tone is friendly and encouraging
- [ ] No assumed knowledge (explain everything)
- [ ] 2-4 paragraphs per sub-section
- [ ] Practical examples throughout
- [ ] Common pitfalls section included
- [ ] Best practices section included

**Quiz Questions**:
- [ ] 6-8 questions total
- [ ] All major sub-topics covered
- [ ] Difficulty range from 2-6
- [ ] Mix of question types
- [ ] All code tested and working
- [ ] Explanations are thorough
- [ ] Tags are accurate and complete

**Overall**:
- [ ] Reviewed against similar resources for consistency
- [ ] Spell-checked
- [ ] Code formatted properly
- [ ] No broken links

## Writing Effective Questions

### Question Type Guidelines

**Multiple Choice** (most common):
- 4 answer choices
- Incorrect answers should be plausible (based on common mistakes)
- Avoid "all of the above" or "none of the above"

**Code Output**:
- Code should be short (< 10 lines)
- Focus on one concept
- Output should be deterministic

**True/False** (use sparingly):
- Make the statement specific, not general
- Avoid trick questions

**Code Completion**:
- Provide clear context
- One clear correct answer

### Difficulty Level Guide

**Level 1-2** (Basic):
- Direct recall of facts
- Simple code reading
- Basic syntax understanding

Example: "What keyword is used to declare a variable in Java?"

**Level 3-4** (Understanding):
- Apply concepts in simple scenarios
- Predict simple code behavior
- Identify common errors

Example: "What will this code print? `System.out.println(5 / 2);`"

**Level 5-6** (Analysis):
- Combine multiple concepts
- Debug more complex code
- Compare approaches

Example: "Which approach is better for checking multiple conditions and why?"

**Level 7** (Advanced):
- Complex problem-solving
- Deep understanding
- Performance considerations

Example: "Explain why this recursive solution might cause a stack overflow."

### Common Pitfalls When Writing Questions

❌ **Don't**:
- Make questions unnecessarily tricky
- Use ambiguous wording
- Test obscure edge cases
- Write overly long questions
- Use jargon without explanation

✅ **Do**:
- Test understanding, not memorization
- Write clear, specific questions
- Focus on practical concepts
- Keep questions concise
- Explain every answer thoroughly

## Code Examples Best Practices

### Guidelines

1. **Keep it short**: 3-10 lines for most examples
2. **Make it runnable**: Examples should compile and run
3. **Use realistic names**: `userName` not `x`, `studentScore` not `a`
4. **Include output**: Show what the code produces
5. **Comment sparingly**: Only when necessary for understanding

### Good Example

```java
// Good: Clear, realistic, shows input and output
String studentName = "Jordan";
int score = 87;

if (score >= 80) {
    System.out.println(studentName + " earned a B!");
}
// Output: Jordan earned a B!
```

### Bad Example

```java
// Bad: Unclear variables, no output shown, too abstract
String s = "X";
int i = 10;
if (i > 5) {
    s += "Y";
}
```

## Review Process

### Self-Review Checklist

Before submitting:
1. Read your entire resource out loud
2. Run all code examples
3. Take your own quiz (cover the answers!)
4. Check against CURRICULUM.md for completeness
5. Compare to existing resources for consistency

### Peer Review

Have another contributor review:
- Technical accuracy
- Tone and style consistency
- Question quality and clarity
- Coverage of curriculum topics

## Tips from Experienced Contributors

### For Study Guides

1. **Start with "Why?"**: Explain why the concept matters before diving into "how"
2. **Use analogies**: Compare programming concepts to real-world things
3. **Build progressively**: Simple examples first, complexity later
4. **Anticipate confusion**: Address common questions preemptively
5. **End with encouragement**: Remind students this is learnable

### For Quiz Questions

1. **Test understanding**: Not just memorization
2. **Learn from wrong answers**: Distractors should teach something
3. **Vary difficulty**: Don't make everything hard or everything easy
4. **Be consistent with study guide**: Use same terminology and examples
5. **Explain thoroughly**: Treat explanations as mini-lessons

## Common Questions

**Q: How long should a learning resource be?**
A: Typically 300-500 lines total (including quiz questions). Focus on quality over quantity.

**Q: Should I include advanced topics?**
A: Stick to the curriculum scope. Note advanced topics for future exploration but don't explain them in detail.

**Q: What if there's already a resource for my topic?**
A: Check if it has quiz questions. If not, add them! If it does, you can still improve the study guide.

**Q: How many code examples is too many?**
A: If every concept has at least one example, you're good. More examples are better than too few.

**Q: Should quiz questions use the same examples as the study guide?**
A: Sometimes yes (reinforcement), sometimes no (transfer). Mix both approaches.

**Q: What if I find an error in an existing resource?**
A: Fix it! Update both the study guide and any affected questions.

## Getting Help

- Check existing resources for examples
- Review [LEARNING_RESOURCE_FORMAT.md](LEARNING_RESOURCE_FORMAT.md)
- Look at [CURRICULUM.md](/CURRICULUM.md) for content scope
- Ask other contributors for review and feedback

## Quick Reference

### File Naming
- Use lowercase with hyphens: `control-structures.md`
- Match curriculum section names when possible
- Keep names descriptive but concise

### Markdown Formatting
```markdown
# Main Title
## Sub-topic (use ##)
### Examples (use ###)
**Bold** for emphasis
`code` for inline code
```language
code blocks
```
---
Horizontal rule before quiz section
```

### Tag Format
```
language, topic, sub-topic-1, sub-topic-2
```

Examples:
- `java, control-structures, if-statement, conditionals`
- `python, data-structures, lists, list-methods`
- `java, oop, classes, objects, constructors`

## Final Thoughts

Creating learning resources is rewarding work. You're helping students learn programming in a supportive, structured way. Take pride in your work, but don't aim for perfection on the first try. Iteration and feedback make resources better!

Happy contributing! 🚀
