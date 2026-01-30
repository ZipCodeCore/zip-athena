# Java Operators and Expressions

Welcome! Let's explore operators and expressions—the building blocks that make your Java programs actually *do* things. If variables are the nouns of programming, operators are the verbs. They let you perform calculations, make comparisons, and combine values in powerful ways.

Think of operators as the tools in your programming toolbox. Just like a carpenter has different tools for different jobs, you'll learn to choose the right operator for each task.

## Arithmetic Operators

Arithmetic operators perform mathematical calculations. You've used most of these in math class!

### Basic Math Operations (+, -, *, /, %)

The five basic arithmetic operators work on numbers:

```java
int a = 10;
int b = 3;

System.out.println(a + b);  // Addition: 13
System.out.println(a - b);  // Subtraction: 7
System.out.println(a * b);  // Multiplication: 30
System.out.println(a / b);  // Division: 3 (integer division!)
System.out.println(a % b);  // Modulo (remainder): 1
```

**Wait, why is 10 / 3 equal to 3?** This is called *integer division*. When both operands are integers, Java gives you an integer result. The decimal part is simply discarded, not rounded. So 10 divided by 3 is 3.333..., but Java drops the .333 and gives you 3.

The modulo operator `%` gives you the remainder. When you divide 10 by 3, you get 3 with a remainder of 1. That's what `10 % 3` returns. It's super useful for checking if numbers are even (`n % 2 == 0`) or for cycling through ranges.

### Getting Decimal Results

Want decimal answers? Make at least one operand a decimal type:

```java
double x = 10.0;
int y = 3;
System.out.println(x / y);  // 3.3333333...

// Or cast one to double
int a = 10;
int b = 3;
System.out.println((double) a / b);  // 3.3333333...
```

### Real-World Example: Splitting a Bill

```java
double totalBill = 87.50;
int people = 4;
double perPerson = totalBill / people;
System.out.println("Each person pays: $" + perPerson);  // $21.875
```

## Assignment Operators

Assignment operators store values in variables. They're how you put data into your containers.

### Basic Assignment (=)

The equals sign assigns a value to a variable:

```java
int score = 95;
String name = "Alex";
boolean isPassing = true;
```

**Important**: The equals sign in programming is *not* like equals in math! In Java, `x = 5` doesn't mean "x equals 5"—it means "assign 5 to x" or "make x become 5."

### Compound Assignment Operators (+=, -=, *=, /=, %=)

These operators combine an arithmetic operation with assignment, making your code cleaner and more expressive:

```java
int x = 10;

x += 5;  // Same as: x = x + 5    (x is now 15)
x -= 3;  // Same as: x = x - 3    (x is now 12)
x *= 2;  // Same as: x = x * 2    (x is now 24)
x /= 4;  // Same as: x = x / 4    (x is now 6)
x %= 4;  // Same as: x = x % 4    (x is now 2)
```

Why use these? They're not just shorter—they're clearer. When you write `score += bonus`, you're explicitly saying "add the bonus to the score," which is more readable than `score = score + bonus`.

### Real-World Example: Bank Account

```java
double balance = 1000.00;
double deposit = 250.00;
double withdrawal = 100.00;

balance += deposit;     // Deposit money
balance -= withdrawal;  // Withdraw money
System.out.println("New balance: $" + balance);  // $1150.00
```

## Increment and Decrement Operators

These special operators add or subtract 1 from a variable. They're incredibly common in loops and counters.

### Basic Usage (++, --)

```java
int count = 5;

count++;  // Increment: count is now 6 (same as count = count + 1)
count--;  // Decrement: count is now 5 (same as count = count - 1)
```

Simple, right? But there's a twist...

### Prefix vs. Postfix

You can put `++` or `--` before (prefix) or after (postfix) the variable, and there's a difference:

```java
int a = 5;
int b = ++a;  // Prefix: increment FIRST, then use. a=6, b=6

int x = 5;
int y = x++;  // Postfix: use FIRST, then increment. x=6, y=5
```

Think of prefix as "increment then use" and postfix as "use then increment."

Here's a clearer example:

```java
int counter = 10;
System.out.println(counter++);  // Prints 10, then counter becomes 11
System.out.println(++counter);  // counter becomes 12, then prints 12
```

**Pro tip**: When you're learning, keep increment/decrement on their own line to avoid confusion:

```java
// Clear and simple
counter++;
System.out.println(counter);
```

## Comparison (Relational) Operators

Comparison operators compare two values and return a boolean result (`true` or `false`). You'll use these constantly in if-statements and loops.

### The Six Comparison Operators (==, !=, >, <, >=, <=)

```java
int age = 25;
int votingAge = 18;

System.out.println(age == 25);      // Equal to: true
System.out.println(age != 30);      // Not equal to: true
System.out.println(age > votingAge);  // Greater than: true
System.out.println(age < votingAge);  // Less than: false
System.out.println(age >= 25);      // Greater than or equal to: true
System.out.println(age <= 24);      // Less than or equal to: false
```

**Critical mistake to avoid**: Don't confuse `==` (comparison) with `=` (assignment)!

```java
if (age = 25) { ... }   // ERROR! This tries to assign, not compare
if (age == 25) { ... }  // CORRECT! This compares
```

### Comparing Strings

**Warning**: Don't use `==` for strings! It compares references (memory locations), not the actual text:

```java
String s1 = "hello";
String s2 = "hello";

if (s1 == s2) { ... }       // DON'T! Compares references
if (s1.equals(s2)) { ... }  // CORRECT! Compares content
```

### Real-World Example: Age Check

```java
int userAge = 20;
int minAge = 18;
int maxAge = 65;

boolean canVote = userAge >= minAge;
boolean isAdult = userAge >= 18;
boolean isSenior = userAge >= 65;

System.out.println("Can vote: " + canVote);    // true
System.out.println("Is senior: " + isSenior);  // false
```

## Logical Operators

Logical operators combine multiple boolean conditions. They're essential for complex decision-making.

### AND Operator (&&)

Both conditions must be true for the result to be true:

```java
int age = 25;
boolean hasLicense = true;

if (age >= 18 && hasLicense) {
    System.out.println("Can drive legally");  // Prints because BOTH are true
}

int score = 85;
if (score >= 90 && score <= 100) {
    System.out.println("Grade: A");  // Doesn't print - first condition is false
}
```

Think of `&&` as requiring *all* conditions to be satisfied.

### OR Operator (||)

At least one condition must be true:

```java
String day = "Saturday";

if (day.equals("Saturday") || day.equals("Sunday")) {
    System.out.println("It's the weekend!");  // Prints!
}

boolean isRaining = false;
boolean isCold = true;

if (isRaining || isCold) {
    System.out.println("Bring a jacket");  // Prints because isCold is true
}
```

Think of `||` as requiring *any* condition to be satisfied.

### NOT Operator (!)

The NOT operator reverses a boolean value:

```java
boolean isRaining = false;

if (!isRaining) {
    System.out.println("No umbrella needed!");  // Prints because NOT false = true
}

boolean hasTicket = true;
if (!hasTicket) {
    System.out.println("Buy a ticket");  // Doesn't print
}
```

You can read `!` as "not" when you read the code: `!isRaining` reads as "not raining."

### Short-Circuit Evaluation

Java is smart about evaluating logical operators. With `&&`, if the first condition is false, Java doesn't even check the second one—why bother when both must be true?

```java
int x = 0;
int y = 10;

// Java checks x > 0, sees it's false, and STOPS (doesn't do y/x)
if (x > 0 && y / x > 2) {
    System.out.println("This prevents division by zero!");
}
```

This short-circuit behavior prevents errors and makes your code more efficient!

### Real-World Example: Login Validation

```java
String username = "alex";
String password = "secret123";
boolean accountActive = true;

if (username.equals("alex") && password.equals("secret123") && accountActive) {
    System.out.println("Login successful!");
} else {
    System.out.println("Login failed!");
}
```

## Bitwise Operators

Bitwise operators work directly with the binary (0s and 1s) representation of numbers. They're more advanced but useful for performance-critical code and low-level operations.

### Bitwise AND (&), OR (|), XOR (^), NOT (~)

```java
int a = 5;   // Binary: 0101
int b = 3;   // Binary: 0011

System.out.println(a & b);  // AND: 0001 = 1
System.out.println(a | b);  // OR:  0111 = 7
System.out.println(a ^ b);  // XOR: 0110 = 6
System.out.println(~a);     // NOT: ...11111010 = -6 (two's complement)
```

**Bitwise AND (`&`)**: Each bit is 1 only if both corresponding bits are 1.
**Bitwise OR (`|`)**: Each bit is 1 if either corresponding bit is 1.
**Bitwise XOR (`^`)**: Each bit is 1 if the corresponding bits are different.
**Bitwise NOT (`~`)**: Flips all bits.

### Shift Operators (<<, >>, >>>)

Shift operators move bits left or right:

```java
int x = 8;  // Binary: 1000

System.out.println(x << 1);  // Left shift: 10000 = 16 (multiply by 2)
System.out.println(x >> 1);  // Right shift: 0100 = 4 (divide by 2)
System.out.println(x << 2);  // Left shift 2: 100000 = 32 (multiply by 4)
```

**Left shift (`<<`)**: Shifts bits left, filling with 0s. Each shift left multiplies by 2.
**Right shift (`>>`)**: Shifts bits right, preserving sign. Each shift right divides by 2.
**Unsigned right shift (`>>>`)**: Shifts bits right, filling with 0s regardless of sign.

### When to Use Bitwise Operators

Most beginners won't use these often, but they're useful for:
- Setting flags or permissions (e.g., file permissions)
- Optimizing multiplication/division by powers of 2
- Working with color values (RGB)
- Low-level hardware programming

```java
// Example: Checking if a number is even (last bit is 0)
int num = 42;
if ((num & 1) == 0) {
    System.out.println("Even");  // Faster than num % 2 == 0
}
```

## Ternary Operator

The ternary operator is a shorthand for simple if-else statements. It's called "ternary" because it has three parts.

### Basic Syntax (? :)

```java
condition ? valueIfTrue : valueIfFalse
```

Here's a simple example:

```java
int age = 20;
String status = (age >= 18) ? "Adult" : "Minor";
System.out.println(status);  // "Adult"
```

This is equivalent to:

```java
String status;
if (age >= 18) {
    status = "Adult";
} else {
    status = "Minor";
}
```

### More Examples

```java
// Find the maximum of two numbers
int a = 15, b = 20;
int max = (a > b) ? a : b;
System.out.println("Max: " + max);  // 20

// Determine pass/fail
int score = 72;
String result = (score >= 60) ? "Pass" : "Fail";
System.out.println(result);  // "Pass"

// Even or odd
int num = 7;
String parity = (num % 2 == 0) ? "Even" : "Odd";
System.out.println(parity);  // "Odd"
```

### Nested Ternary (Use Sparingly!)

You can nest ternary operators, but it gets hard to read quickly:

```java
int score = 85;
String grade = (score >= 90) ? "A" :
               (score >= 80) ? "B" :
               (score >= 70) ? "C" :
               (score >= 60) ? "D" : "F";
System.out.println(grade);  // "B"
```

**Pro tip**: If you're nesting more than once, use a regular if-else instead for readability!

## Operator Precedence and Associativity

When you combine multiple operators, Java needs to know which to evaluate first. This is called *precedence*, just like PEMDAS in math.

### Precedence Order

Here's the precedence from highest to lowest:

1. **Postfix**: `x++`, `x--`
2. **Prefix/Unary**: `++x`, `--x`, `!`, `~`, `+`, `-` (unary)
3. **Multiplicative**: `*`, `/`, `%`
4. **Additive**: `+`, `-`
5. **Shift**: `<<`, `>>`, `>>>`
6. **Relational**: `<`, `>`, `<=`, `>=`
7. **Equality**: `==`, `!=`
8. **Bitwise AND**: `&`
9. **Bitwise XOR**: `^`
10. **Bitwise OR**: `|`
11. **Logical AND**: `&&`
12. **Logical OR**: `||`
13. **Ternary**: `? :`
14. **Assignment**: `=`, `+=`, `-=`, `*=`, `/=`, `%=`, etc.

### Examples

```java
int result = 5 + 3 * 2;  // result is 11, not 16!
// Multiplication happens before addition: 5 + (3 * 2)

int x = 10;
int y = 2 + x * 3;  // y is 32, not 36
// Multiplication first: 2 + (10 * 3)

boolean b = 5 > 3 && 2 < 4;  // true
// Comparisons happen before logical AND: (5 > 3) && (2 < 4)
```

### Use Parentheses for Clarity

When in doubt, use parentheses! They have the highest precedence and make your intent crystal clear:

```java
int result = (5 + 3) * 2;  // result is 16
// Parentheses force addition first

int avg = (a + b + c) / 3;  // Clear: sum first, then divide
```

### Associativity

When operators have the same precedence, associativity determines the order:

```java
int x = 10 - 5 - 2;  // Left-to-right: (10 - 5) - 2 = 3, not 10 - (5 - 2) = 7
int y = 16 / 4 / 2;  // Left-to-right: (16 / 4) / 2 = 2, not 16 / (4 / 2) = 8

// Assignment is right-to-left
int a, b, c;
a = b = c = 5;  // Right-to-left: a = (b = (c = 5))
```

Most operators are left-to-right, but assignment operators are right-to-left.

## Expression Evaluation

An expression is any valid combination of variables, literals, and operators that produces a value. Let's understand how Java evaluates them.

### Simple Expressions

```java
int x = 5;
int y = 10;
int result = x + y;  // Expression: x + y evaluates to 15
```

### Complex Expressions

Java evaluates expressions following precedence and associativity rules:

```java
int a = 5, b = 10, c = 15;
int result = a + b * c - 10;
// Step 1: b * c = 150 (multiplication first)
// Step 2: a + 150 = 155 (addition)
// Step 3: 155 - 10 = 145 (subtraction)
```

### Side Effects

Some expressions have side effects—they change the state of variables:

```java
int x = 5;
int y = x++ + ++x;  // Side effects! x changes during evaluation
// Step 1: x++ uses 5, then increments x to 6
// Step 2: ++x increments x to 7, then uses 7
// Step 3: y = 5 + 7 = 12
// Final: x=7, y=12
```

**Pro tip**: Avoid complex expressions with multiple side effects. Break them into multiple lines:

```java
int x = 5;
int temp1 = x++;  // temp1 = 5, x = 6
++x;              // x = 7
int y = temp1 + x;  // y = 12
```

### Type Promotion in Expressions

When mixing types, Java automatically promotes to the "wider" type:

```java
int i = 10;
double d = 5.5;
double result = i + d;  // i promoted to 10.0, result is 15.5

byte b = 5;
short s = 10;
int sum = b + s;  // Both promoted to int, result is int
```

The promotion hierarchy: `byte` → `short` → `int` → `long` → `float` → `double`

### Boolean Expressions

Expressions that evaluate to boolean are essential for control flow:

```java
int age = 20;
int income = 50000;

boolean isAdult = age >= 18;  // true
boolean canGetLoan = isAdult && income > 30000;  // true
boolean isYoungProfessional = age < 30 && income > 40000;  // true
```

## Common Pitfalls and Best Practices

### Pitfall 1: Integer Division

```java
int result = 5 / 2;  // result is 2, not 2.5!
```

**Fix**: Use at least one double: `5.0 / 2` or `(double) 5 / 2`

### Pitfall 2: Comparing Strings with ==

```java
String s1 = "hello";
String s2 = new String("hello");
if (s1 == s2) { ... }  // DON'T! Often gives wrong result
```

**Fix**: Use `.equals()`: `if (s1.equals(s2)) { ... }`

### Pitfall 3: Assignment in Conditions

```java
if (x = 5) { ... }  // ERROR! Assignment, not comparison
```

**Fix**: Use `==` for comparison: `if (x == 5) { ... }`

### Pitfall 4: Confusing Increment Order

```java
int x = 5;
int y = x++ + x;  // Confusing! What's the value?
```

**Fix**: Use separate statements for clarity

### Best Practices

1. **Use parentheses** when precedence isn't obvious
2. **Watch for integer division** when you need decimal results
3. **Use `.equals()` for strings**, never `==`
4. **Prefer compound assignment** (`+=`, `-=`) for clarity
5. **Keep expressions simple**—break complex ones into multiple lines
6. **Put increments/decrements on their own line** when learning
7. **Use meaningful variable names** even in expressions

## Real-World Example: Grade Calculator

```java
public class GradeCalculator {
    public static void main(String[] args) {
        // Student scores
        int homework = 85;
        int midterm = 78;
        int finalExam = 92;
        
        // Weights
        double hwWeight = 0.30;
        double midtermWeight = 0.30;
        double finalWeight = 0.40;
        
        // Calculate weighted average
        double average = homework * hwWeight + 
                        midterm * midtermWeight + 
                        finalExam * finalWeight;
        
        // Determine letter grade using ternary
        String grade = (average >= 90) ? "A" :
                      (average >= 80) ? "B" :
                      (average >= 70) ? "C" :
                      (average >= 60) ? "D" : "F";
        
        // Check if passing (>= 60)
        boolean isPassing = average >= 60;
        
        System.out.println("Average: " + average);
        System.out.println("Grade: " + grade);
        System.out.println("Passing: " + isPassing);
    }
}
```

---

## Quiz Questions

### Question 1: Integer Division Basics
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
When dividing two integers in Java, the result is always an integer with the decimal part truncated (not rounded). 17 divided by 5 is 3.4, but Java drops the .4, giving us 3. To get 3.4, you would need to use at least one double: `17.0 / 5` or `(double)17 / 5`.

---

### Question 2: Modulo Operator
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: java, operators, modulo, arithmetic

**Question Text**:
What does the expression `23 % 7` evaluate to?

**Answer Choices**:
- A) 3
- B) 2
- C) 3.285
- D) 7

**Correct Answer**: B

**Explanation**:
The modulo operator `%` returns the remainder after division. 23 divided by 7 is 3 with a remainder of 2 (because 7 × 3 = 21, and 23 - 21 = 2). So `23 % 7` equals 2. This operator is commonly used to check if numbers are even (`n % 2 == 0`) or to cycle through ranges.

---

### Question 3: Compound Assignment
**Type**: code-output
**Difficulty**: 3
**Tags**: java, operators, assignment, compound

**Question Text**:
What is the value of `x` after this code executes?

```java
int x = 10;
x += 5;
x *= 2;
x -= 8;
```

**Answer Choices**:
- A) 22
- B) 42
- C) 30
- D) 12

**Correct Answer**: A

**Explanation**:
Following step by step: Start with x=10. Then `x += 5` makes x=15. Then `x *= 2` makes x=30. Finally `x -= 8` makes x=22. Compound assignment operators perform the operation and assign the result back to the variable.

---

### Question 4: Prefix vs Postfix Increment
**Type**: code-output
**Difficulty**: 4
**Tags**: java, operators, increment, prefix, postfix

**Question Text**:
What will be printed?

```java
int a = 5;
int b = a++;
System.out.println(a + " " + b);
```

**Answer Choices**:
- A) 5 5
- B) 6 5
- C) 6 6
- D) 5 6

**Correct Answer**: B

**Explanation**:
The postfix increment `a++` assigns the current value to `b` FIRST (b=5), then increments `a` (a=6). So the output is "6 5". If it were prefix (`++a`), the increment would happen first, and both would be 6. Remember: postfix = "use then increment", prefix = "increment then use".

---

### Question 5: Logical AND Short-Circuit
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, operators, logical, short-circuit, and

**Question Text**:
In the expression `(x > 0 && y / x > 2)`, what happens if `x` is 0?

**Answer Choices**:
- A) Division by zero error occurs
- B) The entire expression returns false without evaluating y / x
- C) The expression returns true
- D) Both conditions are always evaluated

**Correct Answer**: B

**Explanation**:
Java uses short-circuit evaluation for `&&`. When the first condition (`x > 0`) is false, Java immediately knows the entire AND expression must be false (both conditions must be true for AND), so it doesn't evaluate the second part (`y / x`). This prevents the division by zero error. This is different from the non-short-circuit bitwise AND operator `&`.

---

### Question 6: Operator Precedence
**Type**: code-output
**Difficulty**: 4
**Tags**: java, operators, precedence, expressions

**Question Text**:
What is the value of `result`?

```java
int result = 2 + 3 * 4 - 5;
```

**Answer Choices**:
- A) 15
- B) 9
- C) 20
- D) 14

**Correct Answer**: B

**Explanation**:
Following operator precedence (multiplication before addition/subtraction): First `3 * 4 = 12`, then `2 + 12 = 14`, finally `14 - 5 = 9`. Java follows the same precedence rules as mathematics (like PEMDAS). If you wanted a different result, you'd need parentheses to change the order: `(2 + 3) * (4 - 5)` would give -5.

---

### Question 7: Ternary Operator
**Type**: code-output
**Difficulty**: 3
**Tags**: java, operators, ternary, conditional

**Question Text**:
What does this code print?

```java
int score = 75;
String result = (score >= 60) ? "Pass" : "Fail";
System.out.println(result);
```

**Answer Choices**:
- A) Pass
- B) Fail
- C) 75
- D) true

**Correct Answer**: A

**Explanation**:
The ternary operator `? :` is a shorthand if-else. It evaluates the condition `(score >= 60)`, which is true since 75 >= 60. When the condition is true, it returns the value after the `?` (which is "Pass"). If the condition were false, it would return the value after the `:` (which is "Fail"). So "Pass" is printed.

---

### Question 8: Bitwise Shift Operations
**Type**: code-output
**Difficulty**: 6
**Tags**: java, operators, bitwise, shift

**Question Text**:
What is the output of this code?

```java
int x = 16;
System.out.println(x >> 2);
```

**Answer Choices**:
- A) 8
- B) 4
- C) 32
- D) 64

**Correct Answer**: B

**Explanation**:
The right shift operator `>>` shifts bits to the right. In binary, 16 is 10000. Shifting right by 2 positions gives 00100, which is 4. Each right shift divides by 2, so shifting right by 2 is equivalent to dividing by 4: 16 / 4 = 4. Left shift `<<` multiplies by powers of 2, while right shift divides by powers of 2.
