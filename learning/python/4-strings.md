# Python Strings

Welcome! Strings are one of the most commonly used data types in programming. Whether you're building a website, analyzing data, or creating a game, you'll work with text constantly. Python makes working with strings powerful and fun. Let's explore everything strings have to offer!

## What Are Strings?

A string is a sequence of characters—basically, text. In Python, strings are immutable, ordered sequences that can contain letters, numbers, symbols, and even emoji! 🎉

## String Creation and Literals

Python gives you multiple ways to create strings. This flexibility makes it easy to handle different situations.

### Single and Double Quotes

You can use either single `'` or double `"` quotes—Python treats them the same:

```python
name = 'Alex'
message = "Hello, World!"
```

**Pro tip**: Pick one style and be consistent! Most Python code uses double quotes, but single quotes are fine too.

### Why Two Quote Types?

Having both quote types makes it easy to include quotes within strings:

```python
# Use double quotes when string contains single quotes
sentence1 = "It's a beautiful day!"

# Use single quotes when string contains double quotes
sentence2 = 'She said, "Hello!"'

# Or use escape characters
sentence3 = "She said, \"Hello!\""
sentence4 = 'It\'s a beautiful day!'
```

### Triple Quotes

Triple quotes (`"""` or `'''`) create multi-line strings:

```python
# Great for paragraphs or documentation
paragraph = """This is a long string
that spans multiple lines.
Very useful for text blocks!"""

# Also works with single quotes
another = '''Line 1
Line 2
Line 3'''

# Commonly used for docstrings
def greet(name):
    """
    This function greets a person by name.
    It returns a friendly greeting string.
    """
    return f"Hello, {name}!"
```

Triple-quoted strings preserve formatting, including line breaks and indentation.

### Escape Sequences

Special characters start with a backslash `\`:

```python
# Common escape sequences
newline = "First line\nSecond line"
tab = "Name:\tAlex"
backslash = "Path: C:\\Users\\Alex"
quote = "He said, \"Hi!\""

# Other useful escapes
carriage_return = "Text\r"
unicode = "Heart: \u2764"  # ❤

print(newline)
# Output:
# First line
# Second line
```

## String Indexing and Slicing

Strings are sequences, so you can access individual characters or ranges of characters.

### Indexing

Access individual characters using square brackets `[]`. Python uses **zero-based indexing**:

```python
text = "Python"

print(text[0])   # 'P' (first character)
print(text[1])   # 'y'
print(text[5])   # 'n' (last character)
```

### Negative Indexing

Negative indices count from the end:

```python
text = "Python"

print(text[-1])  # 'n' (last character)
print(text[-2])  # 'o' (second from end)
print(text[-6])  # 'P' (first character)
```

**Visualization**:
```
String:  P  y  t  h  o  n
Index:   0  1  2  3  4  5
Neg:    -6 -5 -4 -3 -2 -1
```

### Slicing

Extract substrings using slice notation: `string[start:end:step]`

```python
text = "Python Programming"

# Basic slicing [start:end] (end is exclusive!)
print(text[0:6])    # 'Python'
print(text[7:18])   # 'Programming'

# Omitting start or end
print(text[:6])     # 'Python' (from beginning)
print(text[7:])     # 'Programming' (to end)
print(text[:])      # 'Python Programming' (entire string)

# Negative indices in slicing
print(text[-11:])   # 'Programming'
print(text[:-12])   # 'Python'

# Step parameter
print(text[::2])    # 'Pto rgamn' (every 2nd character)
print(text[::3])    # 'Ph gmi' (every 3rd character)

# Reverse a string!
print(text[::-1])   # 'gnimmargorP nohtyP'
```

The slicing syntax is powerful! Remember: `[start:end:step]`
- `start`: index to begin (inclusive)
- `end`: index to stop (exclusive!)
- `step`: increment (default is 1)

## String Immutability

Strings in Python are **immutable**—once created, they cannot be changed. You can't modify individual characters.

```python
text = "Hello"
text[0] = "J"  # ERROR! TypeError: 'str' object does not support item assignment
```

Instead, create a new string:

```python
text = "Hello"
text = "J" + text[1:]  # Create new string: "Jello"

# Or use string methods
text = "Hello"
text = text.replace("H", "J")  # "Jello"
```

Why immutable? It makes strings safer to use (no accidental changes), faster to work with, and allows them to be used as dictionary keys!

## String Methods

Python strings have tons of built-in methods. Let's explore the most useful ones!

### Case Conversion

```python
text = "Python Programming"

print(text.upper())       # 'PYTHON PROGRAMMING'
print(text.lower())       # 'python programming'
print(text.capitalize())  # 'Python programming'
print(text.title())       # 'Python Programming'
print(text.swapcase())    # 'pYTHON pROGRAMMING'
```

### Searching and Testing

```python
text = "Python Programming"

# Check if starts/ends with something
print(text.startswith("Python"))    # True
print(text.endswith("ing"))         # True

# Find position of substring (-1 if not found)
print(text.find("Pro"))             # 7
print(text.find("Java"))            # -1

# Check if substring exists (better to use 'in')
print("Pro" in text)                # True

# Count occurrences
print(text.count("m"))              # 2
```

### Testing String Content

```python
# Test what the string contains
print("hello".isalpha())      # True (only letters)
print("hello123".isalnum())   # True (letters and numbers)
print("12345".isdigit())      # True (only digits)
print("   ".isspace())        # True (only whitespace)
print("Hello".islower())      # False
print("HELLO".isupper())      # True
print("Hello World".istitle()) # True
```

These are super useful for input validation!

### Modification (Creates New Strings)

```python
text = "  Hello World  "

# Remove whitespace
print(text.strip())       # 'Hello World' (both ends)
print(text.lstrip())      # 'Hello World  ' (left)
print(text.rstrip())      # '  Hello World' (right)

# Replace substring
print(text.replace("World", "Python"))  # '  Hello Python  '

# Split into list
print("a,b,c".split(","))              # ['a', 'b', 'c']
print("one two three".split())         # ['one', 'two', 'three']

# Join list into string
words = ["Python", "is", "awesome"]
print(" ".join(words))                  # 'Python is awesome'
print("-".join(words))                  # 'Python-is-awesome'
```

### Alignment and Padding

```python
text = "Python"

# Center, left, right align
print(text.center(10))       # '  Python  '
print(text.ljust(10))        # 'Python    '
print(text.rjust(10))        # '    Python'

# Add specific character
print(text.center(10, "*"))  # '**Python**'
print(text.zfill(10))        # '0000Python' (zero-padding)
```

## String Concatenation and Repetition

### Concatenation with `+`

```python
first = "Hello"
last = "World"
greeting = first + " " + last  # 'Hello World'

# Building longer strings
name = "Alex"
message = "Welcome, " + name + "!"  # 'Welcome, Alex!'
```

**Note**: Concatenating many strings in a loop is inefficient. Use `join()` instead:

```python
# Inefficient (don't do this!)
result = ""
for word in words:
    result += word + " "

# Efficient (do this!)
result = " ".join(words)
```

### Repetition with `*`

```python
print("Ha" * 3)          # 'HaHaHa'
print("-" * 20)          # '--------------------'
print("=" * 40 + "\n")   # Separator line

# Useful for formatting
separator = "=" * 50
print(separator)
print("TITLE")
print(separator)
```

## String Formatting

Python has three main ways to format strings. Each has its uses!

### Old Style: `%` Formatting

```python
name = "Alex"
age = 25

# %s for string, %d for integer, %f for float
message = "My name is %s and I am %d years old." % (name, age)
print(message)

# Float precision
pi = 3.14159
print("Pi is approximately %.2f" % pi)  # 'Pi is approximately 3.14'
```

This style is old but you'll still see it in legacy code.

### Modern: `.format()` Method

```python
name = "Alex"
age = 25

# Positional arguments
message = "My name is {} and I am {} years old.".format(name, age)

# Numbered placeholders
message = "I am {1} and my name is {0}.".format(name, age)

# Named placeholders
message = "My name is {n} and I am {a} years old.".format(n=name, a=age)

# Formatting numbers
pi = 3.14159
print("Pi: {:.2f}".format(pi))  # 'Pi: 3.14'
```

### Best: f-Strings (Python 3.6+)

F-strings are the modern, preferred way:

```python
name = "Alex"
age = 25

# Just prefix with 'f' and use {variables}
message = f"My name is {name} and I am {age} years old."

# Expressions inside {}
print(f"{age} in 5 years will be {age + 5}")

# Formatting
pi = 3.14159
print(f"Pi: {pi:.2f}")  # 'Pi: 3.14'

# Call functions
print(f"Uppercase: {name.upper()}")

# Multiple lines
message = f"""
Name: {name}
Age: {age}
Status: {'Adult' if age >= 18 else 'Minor'}
"""
```

F-strings are faster, more readable, and more powerful. **Use them!**

### Advanced f-String Formatting

```python
# Number formatting
number = 1234567.89
print(f"{number:,.2f}")      # '1,234,567.89' (thousands separator)
print(f"{number:.2e}")       # '1.23e+06' (scientific notation)

# Padding and alignment
name = "Alex"
print(f"{name:>10}")         # '      Alex' (right align)
print(f"{name:<10}")         # 'Alex      ' (left align)
print(f"{name:^10}")         # '   Alex   ' (center)
print(f"{name:*^10}")        # '***Alex***' (center with *)

# Binary, octal, hex
num = 42
print(f"{num:b}")            # '101010' (binary)
print(f"{num:o}")            # '52' (octal)
print(f"{num:x}")            # '2a' (hex)
```

## Raw Strings

Raw strings treat backslashes as literal characters, not escape sequences. Prefix with `r`:

```python
# Regular string (backslash escapes)
path = "C:\new\test"  # Actually has newline and tab!
print(path)
# Output: C:
# ew    est

# Raw string (backslash is literal)
path = r"C:\new\test"
print(path)  # Output: C:\new\test

# Useful for regular expressions
pattern = r"\d+\.\d+"  # Match numbers like 3.14

# Windows file paths
path = r"C:\Users\Alex\Documents"
```

## String Encoding Basics

Strings in Python 3 are Unicode by default, meaning they can contain any character from any language!

### Unicode

```python
# Any character works!
english = "Hello"
spanish = "Hola"
chinese = "你好"
emoji = "Hello 👋 🌍"

print(english, spanish, chinese, emoji)
```

### Encoding and Decoding

Strings are text (Unicode). Bytes are binary data. Convert between them:

```python
# String to bytes (encode)
text = "Hello"
byte_data = text.encode("utf-8")
print(byte_data)  # b'Hello'
print(type(byte_data))  # <class 'bytes'>

# Bytes to string (decode)
text_again = byte_data.decode("utf-8")
print(text_again)  # 'Hello'

# Other encodings
text = "Café"
utf8_bytes = text.encode("utf-8")
latin1_bytes = text.encode("latin-1")
```

**When you need this**: Reading files, network communication, working with APIs. For most everyday programming, you work with strings and don't need to worry about encoding.

## Practical Examples

### Example 1: Input Validation

```python
username = input("Enter username: ")

# Validate username
is_long_enough = len(username) >= 3
has_no_spaces = " " not in username
is_alphanumeric = username.isalnum()

if is_long_enough and has_no_spaces and is_alphanumeric:
    print("Valid username!")
else:
    print("Invalid username!")
```

### Example 2: Text Processing

```python
text = "  Python is AMAZING!  "

# Clean and normalize
cleaned = text.strip().lower()
print(cleaned)  # 'python is amazing!'

# Extract words
words = cleaned.split()
print(words)  # ['python', 'is', 'amazing!']

# Count specific words
word_count = cleaned.count("python")
print(f"'python' appears {word_count} time(s)")
```

### Example 3: Building Reports

```python
# Generate a formatted report
name = "Alex Johnson"
sales = 15000.75
target = 12000
percentage = (sales / target) * 100

report = f"""
{'=' * 40}
SALES REPORT
{'=' * 40}
Name:     {name}
Sales:    ${sales:,.2f}
Target:   ${target:,.2f}
Progress: {percentage:.1f}%
Status:   {'ACHIEVED ✓' if sales >= target else 'BELOW TARGET'}
{'=' * 40}
"""

print(report)
```

### Example 4: Password Strength

```python
password = "MyP@ssw0rd"

# Check password requirements
has_upper = any(c.isupper() for c in password)
has_lower = any(c.islower() for c in password)
has_digit = any(c.isdigit() for c in password)
has_special = any(not c.isalnum() for c in password)
is_long = len(password) >= 8

if all([has_upper, has_lower, has_digit, has_special, is_long]):
    print("Strong password!")
else:
    print("Weak password!")
```

## Common Pitfalls to Avoid

### 1. Strings Are Immutable

```python
text = "hello"
text[0] = "H"  # ERROR! Can't modify

# Create a new string instead
text = "H" + text[1:]  # 'Hello'
```

### 2. String Concatenation in Loops

```python
# Slow! Creates many intermediate strings
result = ""
for i in range(1000):
    result += str(i)

# Fast! Much more efficient
result = "".join(str(i) for i in range(1000))
```

### 3. Forgetting Strings Are Sequences

```python
text = "Python"
print(len(text))   # 6 (not an error!)
print(text[0])     # 'P'
print('P' in text) # True
```

### 4. Slice Indices

```python
text = "Python"
# End index is exclusive!
print(text[0:3])   # 'Pyt' (not 'Pyth')
```

## Best Practices

1. **Use f-strings** for formatting (Python 3.6+)
2. **Use `in` for substring checking** - clearer than `.find()`
3. **Use `.join()` for building strings** from lists
4. **Use raw strings** for paths and regex: `r"C:\path"`
5. **Use triple quotes** for multi-line strings
6. **Be consistent** with quote style (pick ' or ")
7. **Use string methods** instead of reinventing the wheel
8. **Remember immutability** - methods return new strings

## Practice Exercises

1. Write code to reverse a string using slicing.

2. Create a function that checks if a string is a palindrome (reads the same forwards and backwards).

3. Format a number like 1234567.89 to display as "$1,234,567.89".

4. Split a sentence into words and count how many words start with a capital letter.

5. Take a string and create a new string with all vowels removed.

---

## Quiz Questions

### Question 1: String Indexing
**Type**: code-output
**Difficulty**: 2
**Tags**: python, strings, indexing

**Question Text**:
What is the output of this code?

```python
text = "Python"
print(text[-1])
```

**Answer Choices**:
- A) P
- B) n
- C) o
- D) Error

**Correct Answer**: B

**Explanation**:
Negative indices count from the end of the string. `text[-1]` accesses the last character, which is 'n'. `text[-2]` would be 'o', `text[-3]` would be 'h', and so on. This is a very useful Python feature for accessing elements from the end of a sequence without knowing its length.

---

### Question 2: String Slicing
**Type**: code-output
**Difficulty**: 3
**Tags**: python, strings, slicing

**Question Text**:
What does this print?

```python
text = "Python Programming"
print(text[7:11])
```

**Answer Choices**:
- A) Prog
- B) Progr
- C) Programming
- D) rogramming

**Correct Answer**: A

**Explanation**:
String slicing uses `[start:end]` where start is inclusive and end is exclusive. `text[7:11]` extracts characters at indices 7, 8, 9, and 10 (but not 11). At index 7 is 'P', at 8 is 'r', at 9 is 'o', at 10 is 'g', giving us 'Prog'. Remember: the end index is always exclusive!

---

### Question 3: String Immutability
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: python, strings, immutability

**Question Text**:
What happens when you try to execute this code?

```python
text = "Hello"
text[0] = "J"
```

**Answer Choices**:
- A) text becomes "Jello"
- B) Nothing happens, text stays "Hello"
- C) A TypeError is raised
- D) Only the first letter changes to lowercase

**Correct Answer**: C

**Explanation**:
Strings in Python are immutable, meaning they cannot be changed after creation. Attempting to modify a character using index assignment (`text[0] = "J"`) raises a TypeError. To get "Jello", you would need to create a new string: `text = "J" + text[1:]` or `text = text.replace("H", "J")`. Immutability makes strings safer and more efficient.

---

### Question 4: String Methods - Case Conversion
**Type**: code-output
**Difficulty**: 2
**Tags**: python, strings, methods, case

**Question Text**:
What is the output?

```python
text = "python programming"
print(text.title())
```

**Answer Choices**:
- A) python programming
- B) PYTHON PROGRAMMING
- C) Python Programming
- D) Python programming

**Correct Answer**: C

**Explanation**:
The `.title()` method converts a string to title case, where the first letter of each word is capitalized. So "python programming" becomes "Python Programming". This is different from `.capitalize()` which only capitalizes the first letter of the entire string ("Python programming"), and `.upper()` which capitalizes all letters ("PYTHON PROGRAMMING").

---

### Question 5: String Formatting - f-strings
**Type**: code-output
**Difficulty**: 4
**Tags**: python, strings, formatting, f-strings

**Question Text**:
What does this code print?

```python
name = "Alex"
age = 25
print(f"{age + 5}")
```

**Answer Choices**:
- A) {age + 5}
- B) age + 5
- C) 30
- D) 25 + 5

**Correct Answer**: C

**Explanation**:
F-strings (formatted string literals) evaluate expressions inside curly braces. Since `age = 25`, the expression `age + 5` evaluates to 30, which is then converted to a string and printed. F-strings are very powerful—you can put any valid Python expression inside the braces, including function calls, method calls, and arithmetic operations.

---

### Question 6: String Methods - Split and Join
**Type**: code-output
**Difficulty**: 4
**Tags**: python, strings, methods, split, join

**Question Text**:
What is the output?

```python
words = "Python-is-awesome".split("-")
print(" ".join(words))
```

**Answer Choices**:
- A) Python-is-awesome
- B) Python is awesome
- C) Pythonisawesome
- D) ['Python', 'is', 'awesome']

**Correct Answer**: B

**Explanation**:
The `.split("-")` method splits the string at each hyphen, creating a list: ['Python', 'is', 'awesome']. Then `" ".join(words)` joins the list elements with spaces between them, producing "Python is awesome". The `split()` and `join()` methods are inverses: split breaks strings into lists, join combines lists into strings.

---

### Question 7: Raw Strings
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: python, strings, raw-strings, escaping

**Question Text**:
What is the purpose of a raw string (prefixed with `r`)?

**Answer Choices**:
- A) To make the string read-only
- B) To treat backslashes as literal characters, not escape sequences
- C) To remove all whitespace from the string
- D) To encode the string in raw binary format

**Correct Answer**: B

**Explanation**:
Raw strings (e.g., `r"C:\new\folder"`) treat backslashes as literal characters instead of escape sequences. In a regular string, `\n` is a newline and `\t` is a tab. In a raw string, `\n` is literally the two characters backslash-n. Raw strings are especially useful for Windows file paths and regular expressions, where backslashes are common and you don't want them interpreted as escape sequences.

---

### Question 8: String Repetition
**Type**: code-output
**Difficulty**: 2
**Tags**: python, strings, repetition, operators

**Question Text**:
What does this code print?

```python
print("Ha" * 3 + "!")
```

**Answer Choices**:
- A) HaHaHa!
- B) Ha * 3!
- C) Ha!Ha!Ha!
- D) HaHa!

**Correct Answer**: A

**Explanation**:
The `*` operator repeats a string. `"Ha" * 3` creates "HaHaHa" (repeating "Ha" three times), and then `+ "!"` concatenates an exclamation mark, giving "HaHaHa!". String repetition is useful for creating separators, padding, or patterns. It's one of Python's convenient string operations.
