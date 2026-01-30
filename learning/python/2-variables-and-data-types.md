# Python Variables and Data Types

Welcome! Let's explore one of the most fundamental concepts in Python: variables and data types. Think of variables as labeled containers where you store information, and data types as the kind of information each container can hold. The cool thing about Python? It's smart enough to figure out what type of data you're storing automatically!

## What Are Variables?

A variable is a named location in memory that holds a value. Unlike some other languages, Python doesn't make you declare what type of data you're storing—it figures it out for you. This is called **dynamic typing**.

```python
age = 25
name = "Alex"
price = 19.99
is_student = True
```

In these examples:
- `age` is a variable holding the number 25
- `name` is a variable holding the text "Alex"
- `price` is a variable holding the decimal number 19.99
- `is_student` is a variable holding the boolean value True

Notice we didn't have to say "this is an integer" or "this is a string"—Python figured it out!

## Dynamic Typing

One of Python's superpowers is dynamic typing. You can assign a value of any type to a variable, and later assign a completely different type to the same variable.

```python
x = 42          # x is an integer
print(x)        # 42

x = "hello"     # Now x is a string!
print(x)        # hello

x = 3.14        # Now x is a float!
print(x)        # 3.14
```

This flexibility is powerful, but it also means you need to keep track of what type your variables are. 
**Pro tip**: Use descriptive variable names to help remember what they contain!
**Even More Pro tip**: Don't use this! It's generally leads to difficult to understand code!

## Basic Data Types

Python has several built-in data types. Let's explore the most important ones.

### Integers (`int`)

Integers are whole numbers—no decimal points.

```python
count = 100
temperature = -10
year = 2024
big_number = 1_000_000  # Underscores make it readable!
```

Python integers can be as large as you want—there's no fixed size limit like in some other languages. Need to represent the number of atoms in the universe? Python's got you covered!

```python
really_big = 99999999999999999999999999999999
print(really_big)  # Works just fine!
```

### Floats (`float`)

Floats are numbers with decimal points.

```python
pi = 3.14159
price = 19.99
temperature = -40.5
scientific = 1.23e-4  # Scientific notation: 0.000123
```

**Important note**: Floating-point arithmetic isn't always perfectly precise due to how computers store decimals. This is true in all programming languages, not just Python.

```python
result = 0.1 + 0.2
print(result)  # Prints: 0.30000000000000004 (not exactly 0.3!)
```

### Strings (`str`)

Strings are sequences of characters—text data. You can use single quotes, double quotes, or even triple quotes.

```python
name = 'Alex'
message = "Hello, World!"
quote = "He said, 'Hi!'"  # Single quotes inside double quotes
another_quote = 'She said, "Hello!"'  # Double quotes inside single quotes

# Triple quotes for multi-line strings
paragraph = """This is a longer string
that spans multiple lines.
Very useful for documentation!"""
```

**Pro tip**: Use single quotes or double quotes consistently in your code. Python doesn't care which, but pick one style and stick with it!

### Booleans (`bool`)

Booleans represent truth values: `True` or `False`. Notice the capital letters—that's important in Python!

```python
is_student = True
has_license = False
is_raining = True
```

Booleans are often the result of comparisons:

```python
age = 25
is_adult = age >= 18  # True
is_teenager = 13 <= age <= 19  # False
```

### NoneType (`None`)

`None` is a special value that represents the absence of a value. It's Python's equivalent of "null" in other languages.

```python
result = None  # No value yet
optional_value = None

if result is None:
    print("No result yet!")
```

**Important**: Use `is None` or `is not None` to check for None, not `== None`.

## Type Checking with `type()`

Want to know what type a variable is? Use the `type()` function!

```python
age = 25
print(type(age))  # <class 'int'>

name = "Alex"
print(type(name))  # <class 'str'>

price = 19.99
print(type(price))  # <class 'float'>

is_valid = True
print(type(is_valid))  # <class 'bool'>

nothing = None
print(type(nothing))  # <class 'NoneType'>
```

This is super useful when debugging or when you're not sure what type of data you're working with!

## Type Conversion (Casting)

Sometimes you need to convert one type to another. Python provides built-in functions for this.

### Converting to Integer: `int()`

There is a built-in function `int()` that converts values to integers.

```python
# From float (decimal part is truncated, not rounded!)
x = int(3.9)  # x is 3, not 4!
y = int(3.1)  # y is 3

# From string
age_str = "25"
age_int = int(age_str)  # age_int is 25

# Boolean to int
num1 = int(True)   # 1
num2 = int(False)  # 0
```

**Warning**: Converting a string that doesn't represent a valid number will cause an error!

```python
number = int("hello")  # ERROR! ValueError
```

### Converting to Float: `float()`

Sometimes you need a decimal number. Use `float()` to convert values to floats.

```python
# From integer
x = float(5)  # x is 5.0

# From string
price_str = "19.99"
price_float = float(price_str)  # price_float is 19.99

# Boolean to float
num1 = float(True)   # 1.0
num2 = float(False)  # 0.0
```

### Converting to String: `str()`

If you have numbers, booleans, or other types and want to convert them to strings, use `str()`.

```python
# From integer
age = 25
age_str = str(age)  # age_str is "25"

# From float
price = 19.99
price_str = str(price)  # price_str is "19.99"

# From boolean
flag_str = str(True)  # flag_str is "True"

# Everything can become a string!
none_str = str(None)  # none_str is "None"
```

### Converting to Boolean: `bool()`

You can convert values to booleans using `bool()`. In Python, certain values are considered "falsy" (convert to `False`), while others are "truthy" (convert to `True`).

```python
# Non-zero numbers are True
bool(1)    # True
bool(-5)   # True
bool(0)    # False!

# Non-empty strings are True
bool("hello")  # True
bool("")       # False!

# None is False
bool(None)  # False
```

We'll talk more about "truthy" and "falsy" values later!

## Variable Naming Conventions

Python has rules and conventions for naming variables.

### Rules (Must Follow)

- Must start with a letter (a-z, A-Z) or underscore `_`
- Can contain letters, numbers, and underscores
- Cannot use Python keywords (like `if`, `for`, `class`, `True`, etc.)
- Case-sensitive (`age` and `Age` are different variables)

```python
# Valid names
student_name = "Alex"
age = 25
_private = "secret"
name2 = "Bob"

# Invalid names
2name = "Bob"      # ERROR! Can't start with a number
my-name = "Alex"   # ERROR! Can't use hyphens
class = "Math"     # ERROR! 'class' is a Python keyword
```

### Conventions (Should Follow)

Python uses **snake_case** for variable names: all lowercase with underscores between words.

```python
# Good examples (following Python conventions)
student_name = "Alex"
account_balance = 1000.50
is_valid = True
max_retries = 3

# Bad examples (technically valid, but not Pythonic)
studentName = "Alex"    # camelCase (use in JavaScript, not Python!)
StudentName = "Alex"    # PascalCase (reserved for class names)
STUDENT_NAME = "Alex"   # ALL_CAPS (reserved for constants)
```

**Pro tip**: Make names descriptive! `customer_age` is way better than `ca` or `x`.

## Constants Convention

Python doesn't have true constants (values that can't be changed), but there's a convention: use ALL_CAPS with underscores for values that should not be changed.

```python
PI = 3.14159
MAX_RETRIES = 3
DEFAULT_TIMEOUT = 30
API_KEY = "abc123xyz"
```

This tells other programmers "don't change these values!" But Python won't stop you if you do:

```python
MAX_RETRIES = 3
MAX_RETRIES = 10  # Python allows this, but you shouldn't do it!
```

## Multiple Assignment

Python lets you assign multiple variables in one line.

### Assign the Same Value

```python
x = y = z = 0
print(x, y, z)  # 0 0 0
```

### Assign Different Values

```python
name, age, city = "Alex", 25, "New York"
print(name)  # Alex
print(age)   # 25
print(city)  # New York

# Also works with different types
x, y, z = 1, 2.5, "three"
```

This is super convenient and very Pythonic!

## Unpacking

Unpacking is a powerful Python feature that lets you extract values from sequences.

### Basic Unpacking

```python
# Unpacking a list
coordinates = [10, 20]
x, y = coordinates
print(x)  # 10
print(y)  # 20

# Unpacking a tuple
person = ("Alex", 25, "New York")
name, age, city = person
print(name)  # Alex
print(age)   # 25
```

### Swapping Values

Here's a Python trick that looks like magic:

```python
a = 5
b = 10

# Swap values (no temporary variable needed!)
a, b = b, a

print(a)  # 10
print(b)  # 5
```

### Extended Unpacking

Use `*` to capture multiple values:

```python
# Get first and rest
numbers = [1, 2, 3, 4, 5]
first, *rest = numbers
print(first)  # 1
print(rest)   # [2, 3, 4, 5]

# Get first, middle, and last
first, *middle, last = numbers
print(first)   # 1
print(middle)  # [2, 3, 4]
print(last)    # 5
```

## Variable Reassignment

Because Python is dynamically typed, you can reassign variables to new values of any type at any time.

```python
# Start with an integer
count = 10
print(count, type(count))  # 10 <class 'int'>

# Reassign to a float
count = 10.5
print(count, type(count))  # 10.5 <class 'float'>

# Reassign to a string
count = "ten"
print(count, type(count))  # ten <class 'str'>
```

This flexibility is powerful, but use it wisely! Frequently changing a variable's type can make code confusing.

## Practical Examples

Let's put it all together with some real-world examples.

### Example 1: Temperature Conversion

```python
# Convert Fahrenheit to Celsius
fahrenheit = 98.6
celsius = (fahrenheit - 32) * 5 / 9
print(f"{fahrenheit}°F = {celsius:.1f}°C")  # 98.6°F = 37.0°C
```

### Example 2: User Information

```python
# Store user data
username = "alex_dev"
age = 28
height_meters = 1.75
is_premium_member = True

# Calculate age in days (approximately)
DAYS_PER_YEAR = 365
age_in_days = age * DAYS_PER_YEAR

print(f"User: {username}")
print(f"Age: {age} years ({age_in_days} days)")
print(f"Height: {height_meters}m")
print(f"Premium: {is_premium_member}")
```

### Example 3: Type Conversion

```python
# User input comes as a string
user_input = "25"
age = int(user_input)  # Convert to integer

birth_year = 2024 - age
print(f"You were born in approximately {birth_year}")
```

### Example 4: Multiple Assignment

```python
# Initialize game state
player_x, player_y = 0, 0
score, lives = 0, 3
is_game_over = False

print(f"Position: ({player_x}, {player_y})")
print(f"Score: {score}, Lives: {lives}")
```

## Common Pitfalls to Avoid

### 1. Uninitialized Variables

```python
print(x)  # ERROR! NameError: name 'x' is not defined
```

Always assign a value before using a variable!

### 2. Type Confusion

```python
age = "25"  # This is a string, not a number!
next_year_age = age + 1  # ERROR! Can't add string and int

# Fix it by converting
age = int("25")
next_year_age = age + 1  # Works! next_year_age is 26
```

### 3. Integer Division

In Python 3, `/` always gives you a float, even with integers:

```python
result = 5 / 2
print(result)  # 2.5 (float)

# Use // for integer division
result = 5 // 2
print(result)  # 2 (int)
```

### 4. String Concatenation

You can only concatenate strings with strings:

```python
age = 25
message = "I am " + age  # ERROR! Can't concatenate str and int

# Fix it by converting
message = "I am " + str(age)  # Works!

# Or use f-strings (better!)
message = f"I am {age}"  # Works and easier to read!
```

## Best Practices

1. **Use descriptive variable names** - `customer_age` beats `ca`
2. **Follow snake_case convention** - Python style guide (PEP 8)
3. **Use ALL_CAPS for constants** - Makes intent clear
4. **Don't change variable types unnecessarily** - Confusing!
5. **Initialize variables before using them** - Avoid errors
6. **Use type conversion explicitly** - Don't rely on implicit conversion
7. **Use f-strings for formatting** - Modern and readable

## Practice Exercises

1. Create variables to store your name, age, height in meters, and whether you're a student. Use appropriate types!

2. Write code that converts your height from meters to feet (1 meter = 3.28084 feet).

3. Create constants for the speed of light (299,792,458 m/s) and calculate how far light travels in one minute.

4. Swap the values of two variables using Python's unpacking feature.

5. Take a string "123" and convert it to an integer, multiply by 2, then convert back to a string.

## Wrapping Up

Variables and data types are the foundation of Python programming. Understanding dynamic typing, type conversion, and Python's naming conventions will serve you well throughout your Python journey. The flexibility of Python's type system makes it easy to write code quickly, but remember: with great power comes great responsibility—use type conversions wisely!

Take your time to practice these concepts. Try the exercises, experiment with the code examples, and when you feel ready, test your knowledge with the quiz questions below!

---

## Quiz Questions

### Question 1: Dynamic Typing - Basic Concept
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: python, variables, dynamic-typing

**Question Text**:
What is dynamic typing in Python?

**Answer Choices**:
- A) Python requires you to declare the type of each variable before using it
- B) Python automatically determines the type of a variable based on the value assigned to it
- C) Python converts all numbers to strings automatically
- D) Python only allows one data type per program

**Correct Answer**: B

**Explanation**:
Dynamic typing means Python automatically determines the type of a variable based on the value you assign to it. You don't need to explicitly declare types like in Java or C++. For example, `x = 5` automatically makes `x` an integer, while `x = "hello"` makes it a string. This is one of Python's key features that makes it easier to write code quickly.

---

### Question 2: Variable Naming Conventions
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: python, variables, naming, conventions

**Question Text**:
Which variable name follows Python naming conventions (PEP 8)?

**Answer Choices**:
- A) studentName
- B) StudentName
- C) student_name
- D) STUDENT_NAME

**Correct Answer**: C

**Explanation**:
Python uses snake_case for variable names: lowercase letters with underscores between words. So `student_name` is correct. `studentName` is camelCase (used in JavaScript/Java), `StudentName` is PascalCase (used for class names in Python), and `STUDENT_NAME` is ALL_CAPS (reserved for constants in Python).

---

### Question 3: Type Conversion - String to Integer
**Type**: code-output
**Difficulty**: 3
**Tags**: python, type-conversion, casting, int

**Question Text**:
What is the output of this code?

```python
value = "42"
result = int(value) + 8
print(result)
```

**Answer Choices**:
- A) 428
- B) 50
- C) "428"
- D) Error

**Correct Answer**: B

**Explanation**:
The `int()` function converts the string "42" to the integer 42. Then 42 + 8 equals 50. If we hadn't converted the string to an integer first, we would get an error because Python doesn't allow adding strings and integers directly. The result is an integer, not a string, so it prints as 50, not "50".

---

### Question 4: Multiple Assignment
**Type**: code-output
**Difficulty**: 3
**Tags**: python, variables, multiple-assignment, unpacking

**Question Text**:
What will this code print?

```python
a, b, c = 10, 20, 30
a, c = c, a
print(a, b, c)
```

**Answer Choices**:
- A) 10 20 30
- B) 30 20 10
- C) 20 10 30
- D) 10 30 20

**Correct Answer**: B

**Explanation**:
First line assigns a=10, b=20, c=30. The second line swaps `a` and `c` using Python's unpacking feature. So `a` becomes 30 and `c` becomes 10, while `b` stays 20. The final values are a=30, b=20, c=10, which prints as "30 20 10". This is a common Python idiom for swapping values without needing a temporary variable.

---

### Question 5: Boolean Conversion
**Type**: code-output
**Difficulty**: 4
**Tags**: python, bool, type-conversion, truthy-falsy

**Question Text**:
What is the output?

```python
print(bool(0), bool(""), bool(None), bool("False"))
```

**Answer Choices**:
- A) False False False False
- B) False False False True
- C) True True True True
- D) False False True True

**Correct Answer**: B

**Explanation**:
In Python, `0`, empty string `""`, and `None` are all "falsy" values—they convert to `False` when passed to `bool()`. However, the string `"False"` is a non-empty string, so it's "truthy" and converts to `True`! This is important: the string "False" is not the same as the boolean False. Any non-empty string (even "0" or "False") evaluates to True.

---

### Question 6: Type Checking
**Type**: code-output
**Difficulty**: 3
**Tags**: python, type, type-checking

**Question Text**:
What does this print?

```python
x = 5
y = 5.0
print(type(x) == type(y))
```

**Answer Choices**:
- A) True
- B) False
- C) Error
- D) None

**Correct Answer**: B

**Explanation**:
Even though `x` and `y` have the same numeric value (5), they have different types. `x = 5` is an integer (`int`), while `y = 5.0` is a float (`float`). When you compare their types with `type(x) == type(y)`, you're comparing `<class 'int'>` with `<class 'float'>`, which are different, so the result is False.

---

### Question 7: Constants Convention
**Type**: true-false
**Difficulty**: 2
**Tags**: python, constants, naming, conventions

**Question Text**:
True or False: In Python, variables named in ALL_CAPS cannot be modified after they are first assigned.

**Answer Choices**:
- A) True
- B) False

**Correct Answer**: B

**Explanation**:
False! Python doesn't have true constants. ALL_CAPS naming is just a convention that tells other programmers "this value shouldn't be changed." But Python won't stop you from modifying it. For example, you can write `PI = 3.14` and later `PI = 42`—Python allows it. The ALL_CAPS naming is a signal to humans, not a rule enforced by Python. Other languages like Java have the `final` keyword for true constants, but Python relies on conventions.

---

### Question 8: Extended Unpacking
**Type**: code-output
**Difficulty**: 5
**Tags**: python, unpacking, extended-unpacking, variables

**Question Text**:
What will be printed?

```python
numbers = [1, 2, 3, 4, 5]
first, *middle, last = numbers
print(len(middle))
```

**Answer Choices**:
- A) 1
- B) 2
- C) 3
- D) 5

**Correct Answer**: C

**Explanation**:
Extended unpacking with `*middle` captures all values between the first and last. The variable `first` gets 1, `last` gets 5, and `middle` gets everything in between: [2, 3, 4]. The length of the list [2, 3, 4] is 3. Extended unpacking is very useful for handling sequences where you want to extract specific elements and group the rest together.
