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
