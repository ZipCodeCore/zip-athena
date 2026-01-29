# Java Operators and Expressions - Example Template

Welcome! This file demonstrates the complete learning resource format for Athena. It combines study material with integrated quiz questions to create a comprehensive learning experience.

## What Are Operators?

Operators are special symbols that perform operations on variables and values. Think of them as the verbs of programming—they tell Java what action to perform. Just like you use + for addition in math, Java uses the same symbol to add numbers together.

Java has several types of operators, each serving a different purpose. Let's explore the most important ones!

## Arithmetic Operators

These are the operators you use for math operations.

### Basic Math Operations

```java
int a = 10;
int b = 3;

System.out.println(a + b);  // Addition: 13
System.out.println(a - b);  // Subtraction: 7
System.out.println(a * b);  // Multiplication: 30
System.out.println(a / b);  // Division: 3 (integer division!)
System.out.println(a % b);  // Modulo (remainder): 1
```

**Important note about integer division**: When you divide two integers, Java gives you an integer result. So `10 / 3` gives you `3`, not `3.333...`. The decimal part is discarded!

### Getting Decimal Results

If you want decimal results, at least one operand should be a double or float:

```java
double c = 10.0;
int d = 3;
System.out.println(c / d);  // 3.3333...

// Or cast one to double
System.out.println((double) a / b);  // 3.3333...
```

## Assignment Operators

Assignment operators let you store values in variables.

### Basic Assignment

The equals sign `=` assigns a value to a variable:

```java
int score = 95;
String name = "Alex";
```

### Compound Assignment Operators

These are shortcuts that combine an operation with assignment:

```java
int x = 10;

x += 5;  // Same as: x = x + 5    (now x is 15)
x -= 3;  // Same as: x = x - 3    (now x is 12)
x *= 2;  // Same as: x = x * 2    (now x is 24)
x /= 4;  // Same as: x = x / 4    (now x is 6)
x %= 4;  // Same as: x = x % 4    (now x is 2)
```

These shortcuts make your code cleaner and often express intent more clearly.

## Increment and Decrement Operators

Want to add or subtract 1 from a variable? Java has special operators for that:

```java
int count = 5;

count++;  // Increment: count is now 6 (same as count = count + 1)
count--;  // Decrement: count is now 5 (same as count = count - 1)
```

### Prefix vs. Postfix

There's a subtle difference between `++count` (prefix) and `count++` (postfix):

```java
int a = 5;
int b = ++a;  // Increment FIRST, then assign: a=6, b=6

int x = 5;
int y = x++;  // Assign FIRST, then increment: x=6, y=5
```

**Pro tip**: When in doubt, put the increment/decrement on its own line to avoid confusion!

## Comparison Operators

These operators compare values and return `true` or `false`:

```java
int age = 25;

age == 25  // Equal to: true
age != 30  // Not equal to: true
age > 20   // Greater than: true
age < 20   // Less than: false
age >= 25  // Greater than or equal to: true
age <= 24  // Less than or equal to: false
```

**Common mistake**: Don't confuse `==` (comparison) with `=` (assignment)!

```java
if (age = 25) { ... }  // ERROR! This assigns, not compares
if (age == 25) { ... } // CORRECT! This compares
```

## Logical Operators

Logical operators combine multiple boolean conditions:

### AND Operator (`&&`)

Both conditions must be true:

```java
int age = 25;
boolean hasLicense = true;

if (age >= 18 && hasLicense) {
    System.out.println("Can drive");  // Both conditions are true!
}
```

### OR Operator (`||`)

At least one condition must be true:

```java
String day = "Saturday";

if (day.equals("Saturday") || day.equals("Sunday")) {
    System.out.println("It's the weekend!");
}
```

### NOT Operator (`!`)

Reverses a boolean value:

```java
boolean isRaining = false;

if (!isRaining) {
    System.out.println("No umbrella needed!");  // Prints because NOT false = true
}
```

### Short-Circuit Evaluation

Java uses short-circuit evaluation with `&&` and `||`:

```java
// If age < 18, Java never checks hasLicense (no need, first condition is false)
if (age >= 18 && hasLicense) { ... }

// If isWeekend is true, Java never checks isHoliday (no need, already true)
if (isWeekend || isHoliday) { ... }
```

This is both efficient and useful for avoiding errors!

## Operator Precedence

When you have multiple operators, Java follows precedence rules (like PEMDAS in math):

```java
int result = 5 + 3 * 2;  // result is 11, not 16!
// Multiplication happens before addition

int result2 = (5 + 3) * 2;  // result2 is 16
// Parentheses force addition first
```

**Precedence order** (high to low):
1. Parentheses `()`
2. Increment/Decrement `++`, `--`
3. Multiplication/Division/Modulo `*`, `/`, `%`
4. Addition/Subtraction `+`, `-`
5. Comparison `<`, `>`, `<=`, `>=`
6. Equality `==`, `!=`
7. Logical AND `&&`
8. Logical OR `||`
9. Assignment `=`, `+=`, `-=`, etc.

**Pro tip**: When in doubt, use parentheses to make your intent clear!

## Common Pitfalls

### 1. Integer Division

```java
int result = 5 / 2;  // result is 2, not 2.5!
```

**Fix**: Use at least one double: `5.0 / 2` or `(double) 5 / 2`

### 2. Comparing Strings

```java
String s1 = "hello";
String s2 = "hello";

if (s1 == s2) { ... }  // DON'T do this! Compares references, not content

if (s1.equals(s2)) { ... }  // CORRECT! Compares actual content
```

### 3. Assignment vs. Comparison

```java
if (x = 5) { ... }  // ERROR! Assignment, not comparison
if (x == 5) { ... } // CORRECT! Comparison
```

### 4. Confusing Increment Order

```java
int x = 5;
int y = x++ + x;  // What's y? It's confusing!
```

**Fix**: Keep it simple, put increments on their own line.

## Best Practices

1. **Use parentheses** for clarity when combining operators
2. **Watch out for integer division** when you need decimals
3. **Use `.equals()` for strings**, not `==`
4. **Use compound assignment** (`+=`, `-=`) for cleaner code
5. **Don't nest complex expressions** - break them into multiple lines

## Practice Exercises

1. What's the value of `x` after this code? `int x = 10; x = x++ + ++x;`
2. Write an expression to check if a number is between 1 and 100 (inclusive).
3. Why does `5 / 2` equal 2 in Java? How would you get 2.5?
4. What's the difference between `&&` and `||`?

---

## Quiz Questions

This section contains quiz questions that test understanding of Java operators and expressions. Each question is formatted for easy extraction and database import.

### Question 1: Integer Division - Basic Arithmetic
**Type**: code-output
**Difficulty**: 2
**Tags**: java, operators, arithmetic, division

**Question Text**:
What is the output of this code?

```java
System.out.println(17 / 5);
```

**Answer Choices**:
- A) 3.4
- B) 3
- C) 4
- D) 2

**Correct Answer**: B

**Explanation**:
When dividing two integers in Java, the result is always an integer (the decimal part is truncated). 17 divided by 5 is 3.4, but Java drops the .4, giving us 3. To get 3.4, you'd need to use at least one double: `17.0 / 5` or `(double)17 / 5`.

---

### Question 2: Modulo Operator - Remainder
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: java, operators, modulo

**Question Text**:
What does the modulo operator `%` return?

**Answer Choices**:
- A) The quotient of division
- B) The remainder after division
- C) The percentage of two numbers
- D) The average of two numbers

**Correct Answer**: B

**Explanation**:
The modulo operator `%` returns the remainder after division. For example, `17 % 5` returns 2 because 17 divided by 5 is 3 with a remainder of 2. It's commonly used to check if a number is even (`n % 2 == 0`) or odd (`n % 2 == 1`).

---

### Question 3: Compound Assignment - Shorthand
**Type**: code-output
**Difficulty**: 3
**Tags**: java, operators, assignment

**Question Text**:
What is the value of `x` after executing this code?

```java
int x = 10;
x += 5;
x *= 2;
```

**Answer Choices**:
- A) 20
- B) 25
- C) 30
- D) 15

**Correct Answer**: C

**Explanation**:
Starting with x=10, the statement `x += 5` is equivalent to `x = x + 5`, making x=15. Then `x *= 2` is equivalent to `x = x * 2`, making x=30. The final value is 30.

---

### Question 4: Increment Operators - Prefix vs Postfix
**Type**: code-output
**Difficulty**: 4
**Tags**: java, operators, increment

**Question Text**:
What will be printed?

```java
int a = 5;
int b = ++a;
System.out.println(a + " " + b);
```

**Answer Choices**:
- A) 5 5
- B) 6 5
- C) 6 6
- D) 5 6

**Correct Answer**: C

**Explanation**:
The prefix increment `++a` increments `a` BEFORE using its value. So `a` becomes 6, then that value (6) is assigned to `b`. Both `a` and `b` are 6. If it were postfix (`a++`), then the answer would be "6 5" because the assignment would happen before the increment.

---

### Question 5: Logical Operators - Short-Circuit
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, operators, logical, short-circuit

**Question Text**:
In the expression `(x > 0 && y / x > 2)`, what happens if `x` is 0?

**Answer Choices**:
- A) Division by zero error occurs
- B) The entire expression returns false without evaluating y / x
- C) The expression returns true
- D) Both conditions are always evaluated

**Correct Answer**: B

**Explanation**:
Java uses short-circuit evaluation for `&&`. If the first condition (`x > 0`) is false, Java knows the entire expression must be false (both must be true for AND), so it doesn't evaluate the second part. This prevents the division by zero error. This is different from the bitwise AND operator `&`, which always evaluates both sides.

---

### Question 6: Comparison Operators - String Equality
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, operators, comparison, strings

**Question Text**:
What is the correct way to compare two String variables for equality in Java?

**Answer Choices**:
- A) str1 == str2
- B) str1.equals(str2)
- C) str1.compareTo(str2)
- D) compare(str1, str2)

**Correct Answer**: B

**Explanation**:
For strings, use the `.equals()` method to compare content: `str1.equals(str2)`. The `==` operator compares references (memory addresses), not the actual string content, which can lead to incorrect results. The `compareTo()` method returns an integer for ordering, not a boolean for equality.

---

### Question 7: Operator Precedence - Order of Operations
**Type**: code-output
**Difficulty**: 4
**Tags**: java, operators, precedence

**Question Text**:
What is the value of `result`?

```java
int result = 2 + 3 * 4;
```

**Answer Choices**:
- A) 20
- B) 14
- C) 24
- D) 11

**Correct Answer**: B

**Explanation**:
Java follows standard mathematical operator precedence. Multiplication has higher precedence than addition, so `3 * 4` is calculated first (12), then `2 + 12` = 14. If you wanted 20, you'd need parentheses: `(2 + 3) * 4`. This is the same as PEMDAS/order of operations in mathematics.

---

### Question 8: Logical NOT Operator - Boolean Negation
**Type**: code-output
**Difficulty**: 3
**Tags**: java, operators, logical, not

**Question Text**:
What does this code print?

```java
boolean isRaining = false;
if (!isRaining) {
    System.out.println("Go outside");
} else {
    System.out.println("Stay inside");
}
```

**Answer Choices**:
- A) Go outside
- B) Stay inside
- C) Nothing
- D) Error

**Correct Answer**: A

**Explanation**:
The `!` operator negates (reverses) a boolean value. Since `isRaining` is false, `!isRaining` becomes true. Therefore, the if condition is true, and "Go outside" is printed. The NOT operator is useful for checking negative conditions.

