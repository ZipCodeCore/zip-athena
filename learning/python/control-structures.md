# Python Control Structures

Welcome! Let's talk about control structures in Python—the mechanisms that let your code make decisions and repeat actions. Think of them as the GPS and cruise control for your program, directing where it goes and how many times it does something.

## What Are Control Structures?

Control structures determine the flow of your program. Without them, code would just run straight through from top to bottom, which gets boring quickly! Control structures give you the power to:
- Make decisions (if it's raining, bring an umbrella)
- Repeat actions (send a reminder email every day)
- Break out of repetition early (stop searching once you find what you need)

Python's control structures are particularly nice because they're clean and readable, thanks to Python's use of indentation.

## Conditional Statements

### The `if` Statement

This is your basic "if this is true, then do that" statement:

```python
temperature = 75

if temperature > 70:
    print("It's a warm day!")
```

Notice the colon `:` at the end of the if line and the indentation of the print statement. That's how Python knows what's inside the if block.

### The `if-else` Statement

Want to do one thing when a condition is true and something different when it's false? Use `if-else`:

```python
age = 16

if age >= 18:
    print("You can vote!")
else:
    print("You're not old enough to vote yet.")
```

### The `elif` Chain (else-if)

When you have multiple conditions to check, use `elif`:

```python
score = 85

if score >= 90:
    grade = 'A'
elif score >= 80:
    grade = 'B'
elif score >= 70:
    grade = 'C'
elif score >= 60:
    grade = 'D'
else:
    grade = 'F'

print(f"Your grade is: {grade}")
```

Python checks each condition from top to bottom and executes the first one that's true, then skips the rest.

### Nested Conditionals

You can put if statements inside other if statements:

```python
has_ticket = True
age = 25

if has_ticket:
    if age >= 18:
        print("Welcome to the concert!")
    else:
        print("Sorry, this is an 18+ event.")
else:
    print("You need a ticket to enter.")
```

### Conditional Expressions (Ternary Operator)

Python has a concise way to write simple if-else statements in one line:

```python
# Traditional way
if age >= 18:
    status = "adult"
else:
    status = "minor"

# Ternary way (more concise)
status = "adult" if age >= 18 else "minor"

# Another example
message = "Pass" if score >= 60 else "Fail"
```

The syntax is: `value_if_true if condition else value_if_false`

### Multiple Conditions

You can combine conditions with `and`, `or`, and `not`:

```python
age = 25
has_license = True

# Using 'and' - both must be true
if age >= 18 and has_license:
    print("You can drive!")

# Using 'or' - at least one must be true
day = "Saturday"
if day == "Saturday" or day == "Sunday":
    print("It's the weekend!")

# Using 'not' - negates the condition
is_raining = False
if not is_raining:
    print("No umbrella needed!")

# Complex combinations
if (age >= 18 and has_license) or (age >= 16 and has_permit):
    print("You're allowed to drive")
```

## Loop Statements

Loops let you repeat code multiple times without writing it over and over.

### The `while` Loop

A `while` loop keeps running as long as a condition is true:

```python
count = 1

while count <= 5:
    print(f"Count is: {count}")
    count += 1  # Don't forget to update the counter!

# Output:
# Count is: 1
# Count is: 2
# Count is: 3
# Count is: 4
# Count is: 5
```

**Warning:** Make sure your condition eventually becomes false, or you'll create an infinite loop!

```python
# Infinite loop - DON'T DO THIS!
while True:
    print("This will never stop!")
```

### The `for` Loop

The `for` loop in Python is actually a "for-each" loop—it iterates over sequences (like lists, strings, ranges):

```python
# Loop through a list
fruits = ["apple", "banana", "cherry"]
for fruit in fruits:
    print(f"I like {fruit}")

# Loop through a string
for letter in "Python":
    print(letter)

# Output: P, y, t, h, o, n (each on a new line)
```

### The `range()` Function

When you need to repeat something a specific number of times, use `range()`:

```python
# range(stop) - from 0 to stop-1
for i in range(5):
    print(i)
# Output: 0, 1, 2, 3, 4

# range(start, stop) - from start to stop-1
for i in range(2, 6):
    print(i)
# Output: 2, 3, 4, 5

# range(start, stop, step) - with a step value
for i in range(0, 10, 2):
    print(i)
# Output: 0, 2, 4, 6, 8

# Counting backwards
for i in range(10, 0, -1):
    print(i)
# Output: 10, 9, 8, 7, 6, 5, 4, 3, 2, 1
```

### Nested Loops

Loops inside loops are perfect for working with grids or tables:

```python
# Multiplication table
for i in range(1, 4):
    for j in range(1, 4):
        print(f"{i} x {j} = {i*j}")
    print()  # Blank line after each row

# Pattern printing
for row in range(5):
    for col in range(row + 1):
        print("*", end=" ")
    print()

# Output:
# * 
# * * 
# * * * 
# * * * * 
# * * * * *
```

### Looping Through Dictionaries

```python
student = {"name": "Alex", "age": 20, "grade": "A"}

# Loop through keys
for key in student:
    print(key)

# Loop through values
for value in student.values():
    print(value)

# Loop through key-value pairs (most common)
for key, value in student.items():
    print(f"{key}: {value}")
```

### Enumerate and Zip

Python provides helpful functions for more advanced looping:

```python
# enumerate() - get both index and value
fruits = ["apple", "banana", "cherry"]
for index, fruit in enumerate(fruits):
    print(f"{index}: {fruit}")
# Output:
# 0: apple
# 1: banana
# 2: cherry

# Start enumeration at 1 instead of 0
for index, fruit in enumerate(fruits, start=1):
    print(f"{index}: {fruit}")

# zip() - loop through multiple lists together
names = ["Alice", "Bob", "Charlie"]
ages = [25, 30, 35]
for name, age in zip(names, ages):
    print(f"{name} is {age} years old")
```

## Loop Control Statements

These let you control the flow within loops.

### The `break` Statement

Exit a loop immediately:

```python
# Find the first negative number
numbers = [5, 3, -2, 8, -1, 4]
for num in numbers:
    if num < 0:
        print(f"First negative: {num}")
        break  # Stop the loop

# Search for a value
names = ["Alice", "Bob", "Charlie", "Diana"]
search_name = "Charlie"
for name in names:
    if name == search_name:
        print(f"Found {name}!")
        break
else:
    print(f"{search_name} not found")  # Executes if break never happened
```

### The `continue` Statement

Skip the rest of the current iteration and move to the next one:

```python
# Print only odd numbers
for i in range(10):
    if i % 2 == 0:
        continue  # Skip even numbers
    print(i)
# Output: 1, 3, 5, 7, 9

# Skip certain values
for i in range(10):
    if i == 3 or i == 7:
        continue  # Skip 3 and 7
    print(i)
```

### The `pass` Statement

`pass` does nothing—it's a placeholder when you need a statement syntactically but don't want to execute anything:

```python
# Placeholder for future code
for i in range(10):
    if i % 2 == 0:
        pass  # TODO: handle even numbers
    else:
        print(f"{i} is odd")

# Empty function (would cause an error without pass)
def my_function():
    pass  # We'll implement this later
```

### The `else` Clause with Loops

Python lets you add an `else` clause to loops. It executes when the loop completes normally (not via `break`):

```python
# Search for a number
numbers = [1, 2, 3, 4, 5]
search = 10

for num in numbers:
    if num == search:
        print(f"Found {search}!")
        break
else:
    print(f"{search} not found in the list")

# With while loop
password_attempts = 3
correct_password = "secret"

while password_attempts > 0:
    password = input("Enter password: ")
    if password == correct_password:
        print("Access granted!")
        break
    password_attempts -= 1
else:
    print("Too many failed attempts. Account locked.")
```

## Practical Examples

### Example 1: FizzBuzz

A classic programming challenge:

```python
for i in range(1, 31):
    if i % 3 == 0 and i % 5 == 0:
        print("FizzBuzz")
    elif i % 3 == 0:
        print("Fizz")
    elif i % 5 == 0:
        print("Buzz")
    else:
        print(i)
```

### Example 2: Finding the Largest Number

```python
numbers = [45, 23, 67, 89, 12, 56, 34]

largest = numbers[0]  # Start with first number
for num in numbers:
    if num > largest:
        largest = num

print(f"Largest number: {largest}")

# Or use the built-in max() function
largest = max(numbers)
```

### Example 3: Menu System

```python
while True:
    print("\n=== Menu ===")
    print("1. Option 1")
    print("2. Option 2")
    print("3. Option 3")
    print("4. Exit")
    
    choice = input("Enter your choice (1-4): ")
    
    if choice == '1':
        print("You selected Option 1")
    elif choice == '2':
        print("You selected Option 2")
    elif choice == '3':
        print("You selected Option 3")
    elif choice == '4':
        print("Goodbye!")
        break
    else:
        print("Invalid choice. Please try again.")
```

### Example 4: Grade Calculator

```python
students = {
    "Alice": [85, 90, 88, 92],
    "Bob": [78, 82, 80, 85],
    "Charlie": [92, 95, 90, 88]
}

for student, scores in students.items():
    average = sum(scores) / len(scores)
    
    if average >= 90:
        grade = 'A'
    elif average >= 80:
        grade = 'B'
    elif average >= 70:
        grade = 'C'
    elif average >= 60:
        grade = 'D'
    else:
        grade = 'F'
    
    print(f"{student}: Average = {average:.2f}, Grade = {grade}")
```

### Example 5: Prime Number Checker

```python
def is_prime(n):
    if n < 2:
        return False
    
    for i in range(2, int(n ** 0.5) + 1):
        if n % i == 0:
            return False  # Found a divisor, not prime
    
    return True  # No divisors found, it's prime

# Test it
for num in range(1, 20):
    if is_prime(num):
        print(f"{num} is prime")
```

## Common Patterns and Best Practices

### Pattern 1: Counting and Accumulating

```python
# Count even numbers
numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
even_count = 0
for num in numbers:
    if num % 2 == 0:
        even_count += 1

# Accumulate a sum
total = 0
for num in numbers:
    total += num
# Or just use sum(): total = sum(numbers)
```

### Pattern 2: Filtering

```python
# Get all even numbers
numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
evens = []
for num in numbers:
    if num % 2 == 0:
        evens.append(num)

# Or use a list comprehension (more Pythonic)
evens = [num for num in numbers if num % 2 == 0]
```

### Pattern 3: Transforming

```python
# Square all numbers
numbers = [1, 2, 3, 4, 5]
squares = []
for num in numbers:
    squares.append(num ** 2)

# List comprehension version
squares = [num ** 2 for num in numbers]
```

### Pattern 4: Early Exit

```python
# Stop searching once found
def find_first_negative(numbers):
    for num in numbers:
        if num < 0:
            return num
    return None  # None of them were negative
```

## Common Pitfalls to Avoid

### 1. Modifying a List While Iterating

```python
# WRONG - causes unexpected behavior
numbers = [1, 2, 3, 4, 5]
for num in numbers:
    if num % 2 == 0:
        numbers.remove(num)  # Don't do this!

# RIGHT - create a new list
numbers = [1, 2, 3, 4, 5]
numbers = [num for num in numbers if num % 2 != 0]
```

### 2. Forgetting to Update Loop Variables

```python
# Infinite loop!
count = 0
while count < 5:
    print(count)
    # Forgot to add: count += 1
```

### 3. Off-by-One Errors

```python
# This prints 0-9, not 1-10
for i in range(10):
    print(i)

# To get 1-10:
for i in range(1, 11):
    print(i)
```

### 4. Indentation Errors

```python
# WRONG - indentation matters in Python!
if x > 0:
print("Positive")  # ERROR: IndentationError!

# RIGHT
if x > 0:
    print("Positive")
```

## Pythonic Control Flow

Python encourages writing clean, readable code. Here are some Pythonic patterns:

```python
# Instead of this:
found = False
for item in items:
    if item == target:
        found = True
        break
if found:
    print("Found it!")

# Do this:
if target in items:
    print("Found it!")

# Instead of this:
i = 0
while i < len(items):
    print(items[i])
    i += 1

# Do this:
for item in items:
    print(item)

# Instead of this:
if condition == True:  # Redundant - comparing to True is unnecessary
    ...

# Do this:
if condition:  # Cleaner and more Pythonic
    ...
```

## Practice Exercises

1. Write a program that asks the user to guess a number between 1 and 10. Give them 3 attempts.

2. Create a program that prints all prime numbers between 1 and 50.

3. Write a function that takes a list of numbers and returns a new list with only the positive numbers.

4. Implement a simple password validator that checks if a password is at least 8 characters, contains a number, and contains an uppercase letter.

5. Create a program that prints the Fibonacci sequence up to the 15th number.

## Wrapping Up

Control structures are the backbone of programming logic. Python's clean syntax makes them particularly easy to read and write. The key is understanding:
- When to use `if/elif/else` for decisions
- When to use `for` vs. `while` loops
- How to use `break`, `continue`, and `pass` effectively
- How to write clean, Pythonic code

Practice these patterns, experiment with different combinations, and soon you'll be writing elegant, efficient Python code. When you're ready, head over to Athena and test your knowledge with some control structure quizzes!

---

## Quiz Questions

This section contains quiz questions that test understanding of Python control structures. Each question is formatted for easy extraction and database import.

### Question 1: if-elif-else - Basic Conditionals
**Type**: code-output
**Difficulty**: 2
**Tags**: python, control-structures, if-elif-else

**Question Text**:
What will be printed by this Python code?

```python
age = 16

if age >= 18:
    print("Adult")
elif age >= 13:
    print("Teenager")
else:
    print("Child")
```

**Answer Choices**:
- A) Adult
- B) Teenager
- C) Child
- D) Nothing

**Correct Answer**: B

**Explanation**:
The value of age is 16. The first condition (age >= 18) is false, so Python moves to the next condition. The second condition (age >= 13) is true because 16 >= 13, so "Teenager" is printed. Once a condition matches in an if-elif-else chain, the rest are skipped.

---

### Question 2: for Loop with range() - Iteration
**Type**: code-output
**Difficulty**: 2
**Tags**: python, control-structures, for-loop, range

**Question Text**:
How many numbers will be printed by this code?

```python
for i in range(3, 8):
    print(i)
```

**Answer Choices**:
- A) 4
- B) 5
- C) 6
- D) 8

**Correct Answer**: B

**Explanation**:
The range(3, 8) function generates numbers starting from 3 up to (but not including) 8. So it generates: 3, 4, 5, 6, 7 - that's 5 numbers. The range() function's stop value is exclusive, meaning it stops before reaching that number.

---

### Question 3: while Loop - Loop Termination
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: python, control-structures, while-loop

**Question Text**:
What is the main difference between a `while` loop and a `for` loop in Python?

**Answer Choices**:
- A) while loops are faster than for loops
- B) for loops iterate over sequences, while loops run while a condition is true
- C) while loops can only iterate through lists
- D) for loops cannot use break statements

**Correct Answer**: B

**Explanation**:
The main difference is their purpose: `for` loops are designed to iterate over sequences (lists, strings, ranges, etc.), while `while` loops continue running as long as a condition remains true. Both have similar performance and both can use break statements. The choice between them depends on whether you're iterating over a known sequence (use for) or waiting for a condition to become false (use while).

---

### Question 4: break Statement - Early Exit
**Type**: code-output
**Difficulty**: 3
**Tags**: python, control-structures, break

**Question Text**:
What is the final value of `total`?

```python
total = 0
for num in [1, 2, 3, 4, 5]:
    if num == 4:
        break
    total += num
print(total)
```

**Answer Choices**:
- A) 15
- B) 10
- C) 6
- D) 0

**Correct Answer**: C

**Explanation**:
The loop adds each number to total until it encounters 4. The iterations are: total=0+1=1, total=1+2=3, total=3+3=6. When num becomes 4, the break statement immediately exits the loop, so 4 and 5 are never added. The final value is 6.

---

### Question 5: continue Statement - Skipping Iterations
**Type**: code-output
**Difficulty**: 4
**Tags**: python, control-structures, continue

**Question Text**:
What will this code print?

```python
for i in range(5):
    if i == 2:
        continue
    print(i, end=' ')
```

**Answer Choices**:
- A) 0 1 3 4
- B) 0 1 2 3 4
- C) 2
- D) 0 1

**Correct Answer**: A

**Explanation**:
The continue statement skips the rest of the current iteration and moves to the next one. The loop runs for i=0,1,2,3,4. When i equals 2, continue is executed, skipping the print statement for that iteration. All other values (0, 1, 3, 4) are printed with spaces between them due to `end=' '`.

---

### Question 6: else Clause with Loops - Completion Detection
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: python, control-structures, else-clause, loops

**Question Text**:
When does the `else` clause of a for loop execute?

```python
for item in items:
    if condition:
        break
else:
    print("Loop completed")
```

**Answer Choices**:
- A) After every iteration of the loop
- B) When the loop completes normally without hitting a break
- C) When the break statement is executed
- D) Only if the loop never executes

**Correct Answer**: B

**Explanation**:
The else clause of a loop executes when the loop completes normally (iterates through all items) without encountering a break statement. If break is executed, the else clause is skipped. This is unique to Python and useful for search operations where you want to know if something was found (break) or not found (else executes).

---

### Question 7: List Comprehension - Filtering
**Type**: code-output
**Difficulty**: 4
**Tags**: python, control-structures, list-comprehension

**Question Text**:
What will be the value of `evens`?

```python
numbers = [1, 2, 3, 4, 5, 6]
evens = [num for num in numbers if num % 2 == 0]
```

**Answer Choices**:
- A) [1, 3, 5]
- B) [2, 4, 6]
- C) [1, 2, 3, 4, 5, 6]
- D) [0, 2, 4, 6]

**Correct Answer**: B

**Explanation**:
This list comprehension creates a new list containing only the even numbers from the original list. The condition `if num % 2 == 0` filters for numbers where the remainder when divided by 2 is 0 (even numbers). From [1,2,3,4,5,6], the even numbers are [2,4,6].

---

### Question 8: Nested Loops - Loop Control
**Type**: code-output
**Difficulty**: 5
**Tags**: python, control-structures, nested-loops

**Question Text**:
What shape will be printed by this code?

```python
for i in range(1, 4):
    for j in range(i):
        print("*", end="")
    print()
```

**Answer Choices**:
- A) A square (3x3)
- B) A triangle: \n*\n**\n***
- C) A single line of 6 stars
- D) ***\n***\n***

**Correct Answer**: B

**Explanation**:
The outer loop runs 3 times (i=1,2,3). For each value of i, the inner loop runs i times. When i=1, it prints 1 star. When i=2, it prints 2 stars. When i=3, it prints 3 stars. The `print()` after the inner loop adds a newline, creating a triangular pattern. This is a classic nested loop pattern for printing shapes.
