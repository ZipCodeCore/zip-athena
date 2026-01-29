# Java Control Structures

Hey there! Let's talk about control structures in Java. Think of them as the traffic signals and decision points in your code—they determine which direction your program goes and how many times certain actions happen.

## What Are Control Structures?

In simple terms, control structures let your program make decisions and repeat actions. Without them, your code would just run line by line from top to bottom, which gets boring fast! Control structures give you the power to:
- Make choices based on conditions (should I do this or that?)
- Repeat actions multiple times (do this 10 times)
- Skip certain parts of code when needed

## Conditional Statements

### The `if` Statement

This is your basic "if this is true, then do that" statement. It's like checking if it's raining before deciding whether to grab an umbrella.

```java
int temperature = 75;

if (temperature > 70) {
    System.out.println("It's a warm day!");
}
```

Simple, right? If the temperature is greater than 70, we print a message. If not, nothing happens.

### The `if-else` Statement

Sometimes you want to do one thing if a condition is true, and something different if it's false. That's where `if-else` comes in handy.

```java
int age = 16;

if (age >= 18) {
    System.out.println("You can vote!");
} else {
    System.out.println("You're not old enough to vote yet.");
}
```

### The `else-if` Ladder

What if you have more than two possibilities? Use `else-if` to chain multiple conditions together:

```java
int score = 85;

if (score >= 90) {
    System.out.println("Grade: A");
} else if (score >= 80) {
    System.out.println("Grade: B");
} else if (score >= 70) {
    System.out.println("Grade: C");
} else if (score >= 60) {
    System.out.println("Grade: D");
} else {
    System.out.println("Grade: F");
}
```

Java checks each condition from top to bottom and executes the first one that's true. Once it finds a match, it skips the rest.

### Nested `if` Statements

You can put `if` statements inside other `if` statements. It's like making a decision within a decision:

```java
boolean hasTicket = true;
int age = 25;

if (hasTicket) {
    if (age >= 18) {
        System.out.println("Welcome to the concert!");
    } else {
        System.out.println("Sorry, this is an 18+ event.");
    }
} else {
    System.out.println("You need a ticket to enter.");
}
```

### The `switch` Statement

When you have many possible values for a single variable, `switch` can be cleaner than a long `else-if` ladder:

```java
String day = "Monday";

switch (day) {
    case "Monday":
        System.out.println("Start of the work week!");
        break;
    case "Friday":
        System.out.println("TGIF!");
        break;
    case "Saturday":
    case "Sunday":
        System.out.println("Weekend time!");
        break;
    default:
        System.out.println("It's a regular day.");
}
```

**Important:** Don't forget the `break` statements! Without them, execution "falls through" to the next case.

In modern Java (14+), you can also use switch expressions, which are even nicer:

```java
String day = "Friday";
String mood = switch (day) {
    case "Monday" -> "Tired";
    case "Friday" -> "Excited";
    case "Saturday", "Sunday" -> "Relaxed";
    default -> "Normal";
};
System.out.println("Mood: " + mood);
```

**Note:** Switch expressions require Java 14 or later. If you're using an older version, stick with the traditional switch statement above.

## Loop Statements

Loops let you repeat actions. Instead of writing the same code 100 times, you write it once and loop it 100 times!

### The `while` Loop

This loop keeps going as long as a condition is true. It's like saying "while the cookie jar isn't empty, keep eating cookies."

```java
int count = 1;

while (count <= 5) {
    System.out.println("Count is: " + count);
    count++;  // Don't forget to update the condition!
}
```

**Warning:** If you forget to update the variable that affects your condition, you'll create an infinite loop. Not fun!

### The `do-while` Loop

This is like a `while` loop, but it always runs at least once because it checks the condition at the end:

```java
int number = 10;

do {
    System.out.println("Number is: " + number);
    number++;
} while (number <= 5);
```

Even though `number` starts at 10 (which doesn't meet the condition), the loop runs once and prints "Number is: 10".

### The `for` Loop

This is the most common loop when you know exactly how many times you want to repeat something:

```java
for (int i = 0; i < 5; i++) {
    System.out.println("Iteration: " + i);
}
```

Let's break down that `for` statement:
- `int i = 0` - Start with i at 0
- `i < 5` - Keep going while i is less than 5
- `i++` - After each iteration, increment i by 1

So this prints "Iteration: 0" through "Iteration: 4".

### The Enhanced `for` Loop (for-each)

When you want to go through each item in an array or collection, use this simpler syntax:

```java
String[] fruits = {"Apple", "Banana", "Cherry"};

for (String fruit : fruits) {
    System.out.println("I like " + fruit);
}
```

Read it as: "for each fruit in fruits, do this." Much cleaner than dealing with array indices!

### Nested Loops

You can put loops inside loops. This is great for working with grids or tables:

```java
for (int row = 1; row <= 3; row++) {
    for (int col = 1; col <= 3; col++) {
        System.out.print(row + "," + col + " ");
    }
    System.out.println();
}
```

This creates a simple grid:
```
1,1 1,2 1,3 
2,1 2,2 2,3 
3,1 3,2 3,3 
```

## Jump Statements

These let you jump out of loops or skip iterations.

### The `break` Statement

Use `break` to exit a loop early:

```java
for (int i = 0; i < 10; i++) {
    if (i == 5) {
        break;  // Stop the loop when i reaches 5
    }
    System.out.println(i);
}
// Prints: 0, 1, 2, 3, 4
```

### The `continue` Statement

Use `continue` to skip the rest of the current iteration and jump to the next one:

```java
for (int i = 0; i < 10; i++) {
    if (i % 2 == 0) {
        continue;  // Skip even numbers
    }
    System.out.println(i);
}
// Prints: 1, 3, 5, 7, 9 (only odd numbers)
```

### The `return` Statement

`return` exits from a method entirely. We'll cover this more when we talk about methods, but here's a quick example:

```java
public static int findFirstNegative(int[] numbers) {
    for (int num : numbers) {
        if (num < 0) {
            return num;  // Exit method and return this value
        }
    }
    return 0;  // No negative found
}
```

## Understanding Control Flow

Control flow is the order in which your program executes statements. Think of it as the path your program takes through the code:

```java
int x = 10;

if (x > 5) {
    System.out.println("x is greater than 5");
    
    for (int i = 0; i < x; i++) {
        if (i == 3) {
            continue;
        }
        System.out.println("i = " + i);
    }
}

System.out.println("Done!");
```

Here's what happens:
1. Check if x > 5 (it is, so enter the if block)
2. Print "x is greater than 5"
3. Start the for loop
4. Each iteration, check if i == 3 (if so, skip that iteration)
5. Print the value of i
6. After the loop finishes, print "Done!"

## Common Patterns and Best Practices

### Avoid Deep Nesting

Too many nested if statements or loops make code hard to read:

```java
// Hard to read
if (condition1) {
    if (condition2) {
        if (condition3) {
            // Do something
        }
    }
}

// Better - combine conditions
if (condition1 && condition2 && condition3) {
    // Do something
}
```

### Use Meaningful Conditions

Make your conditions clear:

```java
// Unclear
if (x > 0 && x < 100 && y == 5) {
    // ...
}

// Better
boolean isValidAge = x > 0 && x < 100;
boolean isCorrectStatus = y == 5;

if (isValidAge && isCorrectStatus) {
    // ...
}
```

### Choose the Right Loop

- Use `for` when you know the number of iterations
- Use `while` when you don't know how many iterations you'll need
- Use `for-each` when iterating through collections or arrays
- Use `do-while` when you need at least one iteration

## Practice Exercises

Try these on your own to reinforce what you've learned:

1. Write a program that prints all numbers from 1 to 100, but for multiples of 3, print "Fizz" instead of the number, and for multiples of 5, print "Buzz". For numbers that are multiples of both 3 and 5, print "FizzBuzz".

2. Create a simple menu system using a switch statement that responds to user choices (1-4, with 4 being exit).

3. Write a program that finds the largest number in an array using a for loop.

4. Create a nested loop that prints a multiplication table (1-10).

## Wrapping Up

Control structures are fundamental to programming. They turn your code from a simple list of instructions into a dynamic, responsive program that can make decisions and handle repetition efficiently.

Master these concepts, and you'll be able to solve a huge range of programming problems. Keep practicing, experiment with different combinations, and don't be afraid to make mistakes—that's how you learn!

Ready to test your knowledge? Head over to the Athena quiz system and try some control structure questions!

---

## Quiz Questions

This section contains quiz questions that test understanding of Java control structures. Each question is formatted for easy extraction and database import.

### Question 1: if-else Statement - Basic Condition
**Type**: code-output
**Difficulty**: 2
**Tags**: java, control-structures, if-else

**Question Text**:
What will be printed when this code is executed?

```java
int score = 75;

if (score >= 80) {
    System.out.println("Great!");
} else if (score >= 60) {
    System.out.println("Good!");
} else {
    System.out.println("Keep trying!");
}
```

**Answer Choices**:
- A) Great!
- B) Good!
- C) Keep trying!
- D) Nothing will be printed

**Correct Answer**: B

**Explanation**:
The score is 75, which is not >= 80, so the first condition is false. The second condition checks if score >= 60, which is true (75 >= 60), so "Good!" is printed. Once a condition in an if-else chain matches, the rest are skipped.

---

### Question 2: switch Statement - Fall-through Behavior
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, control-structures, switch

**Question Text**:
What happens if you forget to include a `break` statement in a switch case?

**Answer Choices**:
- A) The program will crash with a syntax error
- B) Only that case will execute, then the switch will exit
- C) Execution will "fall through" to the next case(s)
- D) The default case will execute instead

**Correct Answer**: C

**Explanation**:
Without a `break` statement, execution continues ("falls through") to the next case, executing its code as well. This is called fall-through behavior. It's not a syntax error, but it's usually unintended. The `break` statement is what tells Java to exit the switch block.

---

### Question 3: for Loop - Iteration Count
**Type**: code-output
**Difficulty**: 3
**Tags**: java, control-structures, for-loop

**Question Text**:
How many times will "Hello" be printed?

```java
for (int i = 0; i < 5; i++) {
    System.out.println("Hello");
}
```

**Answer Choices**:
- A) 4 times
- B) 5 times
- C) 6 times
- D) Infinite times

**Correct Answer**: B

**Explanation**:
The loop starts with i=0 and continues while i < 5. It will execute for i=0, 1, 2, 3, and 4 (five values), printing "Hello" 5 times. When i becomes 5, the condition (i < 5) is false, so the loop stops.

---

### Question 4: while Loop - Loop Control
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, control-structures, while-loop

**Question Text**:
What is the main difference between a `while` loop and a `do-while` loop?

**Answer Choices**:
- A) while loops are faster than do-while loops
- B) A do-while loop always executes at least once, while a while loop may not execute at all
- C) while loops can only iterate through arrays
- D) do-while loops cannot use break statements

**Correct Answer**: B

**Explanation**:
The key difference is that a `do-while` loop checks its condition AFTER executing the loop body, guaranteeing at least one execution. A `while` loop checks its condition BEFORE the first execution, so if the condition is initially false, the loop body never runs. Both have the same performance characteristics and can use break statements.

---

### Question 5: break Statement - Loop Termination
**Type**: code-output
**Difficulty**: 4
**Tags**: java, control-structures, break, for-loop

**Question Text**:
What will be the final value of `sum` after this code executes?

```java
int sum = 0;
for (int i = 1; i <= 10; i++) {
    if (i == 5) {
        break;
    }
    sum += i;
}
System.out.println(sum);
```

**Answer Choices**:
- A) 10
- B) 15
- C) 55
- D) 0

**Correct Answer**: A

**Explanation**:
The loop adds i to sum for i=1, 2, 3, and 4, giving us sum = 1+2+3+4 = 10. When i becomes 5, the break statement is executed, immediately exiting the loop before adding 5 to sum. Therefore, sum remains 10.

---

### Question 6: Nested Loops - Pattern Output
**Type**: code-output
**Difficulty**: 5
**Tags**: java, control-structures, nested-loops, for-loop

**Question Text**:
What pattern will this code print?

```java
for (int i = 1; i <= 3; i++) {
    for (int j = 1; j <= i; j++) {
        System.out.print("*");
    }
    System.out.println();
}
```

**Answer Choices**:
- A) *\n**\n***
- B) ***\n**\n*
- C) *\n*\n*
- D) ***\n***\n***

**Correct Answer**: A

**Explanation**:
The outer loop runs 3 times (i=1, 2, 3). For each iteration of the outer loop, the inner loop runs i times. When i=1, it prints one star. When i=2, it prints two stars. When i=3, it prints three stars. Each inner loop is followed by a println(), creating a new line. This creates a triangular pattern.

---

### Question 7: continue Statement - Loop Control
**Type**: code-output
**Difficulty**: 4
**Tags**: java, control-structures, continue, for-loop

**Question Text**:
What numbers will be printed by this code?

```java
for (int i = 1; i <= 5; i++) {
    if (i == 3) {
        continue;
    }
    System.out.println(i);
}
```

**Answer Choices**:
- A) 1 2 4 5
- B) 1 2 3 4 5
- C) 1 2
- D) 3

**Correct Answer**: A

**Explanation**:
The `continue` statement skips the rest of the current iteration and moves to the next one. When i equals 3, the continue statement is executed, skipping the println() for that iteration. All other values (1, 2, 4, 5) are printed normally.

---

### Question 8: Nested if Statements - Logical Flow
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, control-structures, nested-if, logical-operators

**Question Text**:
Which of the following is better practice for checking multiple conditions?

```java
// Option 1
if (age >= 18) {
    if (hasLicense) {
        if (hasInsurance) {
            System.out.println("Can drive");
        }
    }
}

// Option 2
if (age >= 18 && hasLicense && hasInsurance) {
    System.out.println("Can drive");
}
```

**Answer Choices**:
- A) Option 1, because it's more readable
- B) Option 2, because it's more concise and easier to read
- C) They are exactly the same in all aspects
- D) Option 1, because it performs better

**Correct Answer**: B

**Explanation**:
Option 2 is generally preferred because it's more concise, easier to read, and expresses the intent more clearly: all three conditions must be true. Both options have the same logical behavior and performance (short-circuit evaluation), but Option 2 reduces nesting and is considered more maintainable. Excessive nesting (Option 1) can make code harder to follow.
