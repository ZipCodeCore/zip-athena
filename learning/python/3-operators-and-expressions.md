# Python Operators and Expressions

Welcome! Operators are the tools that let you perform operations on data—they're the action words of programming. Think of them as the verbs in your code. Just like you use + for addition in math, Python uses operators to manipulate values and variables. Let's explore how Python's operators work!

## What Are Operators?

Operators are special symbols that perform operations on variables and values. An expression is a combination of values, variables, and operators that Python evaluates to produce a result.

```python
result = 10 + 5  # + is the operator, 10 + 5 is the expression
```

Python has several categories of operators. Let's dive into each one!

## Arithmetic Operators

These are the operators you use for mathematical operations.

### Basic Arithmetic

```python
a = 10
b = 3

print(a + b)   # Addition: 13
print(a - b)   # Subtraction: 7
print(a * b)   # Multiplication: 30
print(a / b)   # Division: 3.3333... (always returns a float!)
print(a // b)  # Floor division: 3 (integer division, rounds down)
print(a % b)   # Modulo (remainder): 1
print(a ** b)  # Exponentiation: 1000 (10 to the power of 3)
```

### Division in Python

Python has two division operators, and understanding the difference is crucial:

```python
# Regular division (/) always returns a float
print(10 / 2)   # 5.0 (not 5!)
print(10 / 3)   # 3.3333...
print(7 / 2)    # 3.5

# Floor division (//) returns an integer (rounds down)
print(10 // 2)  # 5
print(10 // 3)  # 3
print(7 // 2)   # 3
print(-7 // 2)  # -4 (rounds down to more negative)
```

**Pro tip**: Use `/` when you want decimal precision, use `//` when you want integer results.

### Modulo Operator

The modulo operator `%` gives you the remainder after division. It's super useful!

```python
# Check if a number is even or odd
number = 17
if number % 2 == 0:
    print("Even")
else:
    print("Odd")  # This prints!

# Get the last digit of a number
last_digit = 12345 % 10  # 5

# Cycle through values
for i in range(10):
    print(i % 3)  # Prints: 0, 1, 2, 0, 1, 2, 0, 1, 2, 0
```

### Exponentiation

The `**` operator raises a number to a power:

```python
print(2 ** 3)    # 8 (2 cubed)
print(5 ** 2)    # 25 (5 squared)
print(10 ** 6)   # 1000000 (10 to the 6th power)

# Also works with floats
print(4 ** 0.5)  # 2.0 (square root of 4)
print(27 ** (1/3))  # 3.0 (cube root of 27)
```

## Assignment Operators

Assignment operators let you store values in variables.

### Basic Assignment

The equals sign `=` assigns a value to a variable:

```python
score = 95
name = "Alex"
```

**Important**: `=` is assignment, not equality! To test equality, use `==`.

### Compound Assignment Operators

These are shortcuts that combine an operation with assignment:

```python
x = 10

x += 5   # Same as: x = x + 5    (now x is 15)
x -= 3   # Same as: x = x - 3    (now x is 12)
x *= 2   # Same as: x = x * 2    (now x is 24)
x /= 4   # Same as: x = x / 4    (now x is 6.0)
x //= 2  # Same as: x = x // 2   (now x is 3.0)
x %= 2   # Same as: x = x % 2    (now x is 1.0)
x **= 3  # Same as: x = x ** 3   (now x is 1.0)
```

All arithmetic operators have compound versions! They make your code cleaner and often express your intent more clearly.

```python
# Updating a score
score = 100
score += 10  # Much cleaner than: score = score + 10

# Doubling a value
count = 5
count *= 2  # Much cleaner than: count = count * 2
```

## Comparison Operators

Comparison operators compare values and return `True` or `False`.

```python
a = 10
b = 5

print(a == b)   # Equal to: False
print(a != b)   # Not equal to: True
print(a > b)    # Greater than: True
print(a < b)    # Less than: False
print(a >= b)   # Greater than or equal to: True
print(a <= b)   # Less than or equal to: False
```

### Comparing Different Types

```python
# Comparing numbers of different types works
print(5 == 5.0)    # True (values are equal)
print(5 is 5.0)    # False (types are different)

# Comparing strings (lexicographic order)
print("apple" < "banana")  # True
print("abc" == "ABC")      # False (case-sensitive!)

# Comparing with None
x = None
print(x is None)     # True (correct way)
print(x == None)     # True (works, but "is" is preferred)
```

### Chaining Comparisons

Python lets you chain comparisons in a natural, mathematical way:

```python
age = 25

# This works in Python!
if 18 <= age <= 65:
    print("Working age")  # This prints!

# Same as:
if age >= 18 and age <= 65:
    print("Working age")

# More examples
x = 5
print(1 < x < 10)      # True
print(0 <= x <= 100)   # True
```

This is a unique Python feature that makes code very readable!

## Logical Operators

Logical operators combine boolean conditions. Python uses words, not symbols: `and`, `or`, `not`.

### The `and` Operator

Both conditions must be true:

```python
age = 25
has_license = True

if age >= 18 and has_license:
    print("Can drive")  # Both conditions are true!

# False if either is false
print(True and True)    # True
print(True and False)   # False
print(False and True)   # False
print(False and False)  # False
```

### The `or` Operator

At least one condition must be true:

```python
day = "Saturday"

if day == "Saturday" or day == "Sunday":
    print("It's the weekend!")

# True if either is true
print(True or True)     # True
print(True or False)    # True
print(False or True)    # True
print(False or False)   # False
```

### The `not` Operator

Reverses a boolean value:

```python
is_raining = False

if not is_raining:
    print("No umbrella needed!")  # Prints!

print(not True)   # False
print(not False)  # True
```

### Short-Circuit Evaluation

Python uses short-circuit evaluation with `and` and `or`:

```python
# With 'and': if first is False, second isn't evaluated
age = 15
if age >= 18 and check_license():  # check_license() never called!
    print("Can drive")

# With 'or': if first is True, second isn't evaluated
is_weekend = True
if is_weekend or is_holiday():  # is_holiday() never called!
    print("Day off!")
```

This is efficient and helps avoid errors like division by zero:

```python
x = 0
# This is safe! Second part not evaluated when x == 0
if x != 0 and 10 / x > 2:
    print("Large quotient")
```

## Identity Operators

Identity operators check if two variables refer to the same object in memory, not just equal values.

### The `is` Operator

```python
a = [1, 2, 3]
b = [1, 2, 3]
c = a

print(a == b)   # True (same values)
print(a is b)   # False (different objects in memory)
print(a is c)   # True (c points to the same object as a)
```

### Use `is` with `None`

The most common use of `is` is checking for `None`:

```python
result = None

if result is None:  # Correct!
    print("No result")

# Using == works but is not preferred
if result == None:  # Works, but use 'is' instead
    print("No result")
```

### The `is not` Operator

```python
x = 5
if x is not None:
    print("x has a value")  # This prints
```

**Important**: Use `==` for value equality, use `is` for identity (same object).

## Membership Operators

Membership operators test if a value is in a sequence (like a string, list, or tuple).

### The `in` Operator

```python
# In strings
print("a" in "cat")        # True
print("x" in "cat")        # False

# In lists
numbers = [1, 2, 3, 4, 5]
print(3 in numbers)        # True
print(10 in numbers)       # False

# In tuples
colors = ("red", "green", "blue")
print("red" in colors)     # True

# In dictionaries (checks keys)
person = {"name": "Alex", "age": 25}
print("name" in person)    # True
print("Alex" in person)    # False (checks keys, not values)
```

### The `not in` Operator

```python
# Check if something is NOT in a sequence
forbidden = ["admin", "root", "system"]
username = "alex"

if username not in forbidden:
    print("Username accepted")  # This prints
```

Membership operators make checking for values incredibly easy and readable!

## Bitwise Operators

Bitwise operators work on individual bits of integers. These are more advanced but useful for certain tasks.

```python
a = 60   # Binary: 0011 1100
b = 13   # Binary: 0000 1101

print(a & b)   # AND: 12      (0000 1100)
print(a | b)   # OR: 61       (0011 1101)
print(a ^ b)   # XOR: 49      (0011 0001)
print(~a)      # NOT: -61     (inverts all bits)
print(a << 2)  # Left shift: 240   (multiply by 2^2)
print(a >> 2)  # Right shift: 15   (divide by 2^2)
```

**When to use**: Bitwise operators are useful for flags, permissions, and low-level programming. For most everyday Python programming, you won't need them.

## Operator Precedence

When you have multiple operators, Python follows precedence rules (like PEMDAS in math):

```python
result = 5 + 3 * 2  # result is 11, not 16!
# Multiplication happens before addition

result2 = (5 + 3) * 2  # result2 is 16
# Parentheses force addition first
```

### Precedence Order (Highest to Lowest)

1. **Parentheses**: `()`
2. **Exponentiation**: `**`
3. **Unary operators**: `+x`, `-x`, `~x`, `not x`
4. **Multiplication, Division, Modulo**: `*`, `/`, `//`, `%`
5. **Addition, Subtraction**: `+`, `-`
6. **Bitwise shifts**: `<<`, `>>`
7. **Bitwise AND**: `&`
8. **Bitwise XOR**: `^`
9. **Bitwise OR**: `|`
10. **Comparisons**: `<`, `<=`, `>`, `>=`, `==`, `!=`
11. **Identity**: `is`, `is not`
12. **Membership**: `in`, `not in`
13. **Logical NOT**: `not`
14. **Logical AND**: `and`
15. **Logical OR**: `or`

```python
# Complex expression
result = 2 + 3 * 4 ** 2 / 2 - 1
# Step 1: 4 ** 2 = 16
# Step 2: 3 * 16 = 48
# Step 3: 48 / 2 = 24.0
# Step 4: 2 + 24.0 = 26.0
# Step 5: 26.0 - 1 = 25.0
print(result)  # 25.0
```

**Pro tip**: When in doubt, use parentheses to make your intent clear! Readable code is better than clever code.

## Expression Evaluation

Python evaluates expressions from left to right, following precedence rules.

```python
# Left to right with same precedence
result = 10 - 5 - 2  # (10 - 5) - 2 = 3, not 10 - (5 - 2) = 7

# Mixed operations
x = 5
y = 10
result = x + y * 2  # 5 + 20 = 25 (multiplication first)

# With parentheses
result = (x + y) * 2  # 15 * 2 = 30 (addition first)
```

## Truthy and Falsy Values

In Python, every value can be evaluated in a boolean context. Some values are "truthy" (evaluate to True) and some are "falsy" (evaluate to False).

### Falsy Values

These values are considered False:
- `False` (the boolean)
- `None`
- `0` (zero of any numeric type)
- `0.0`
- `""` (empty string)
- `[]` (empty list)
- `()` (empty tuple)
- `{}` (empty dictionary)
- `set()` (empty set)

```python
# All of these print "Falsy!"
if not 0:
    print("Falsy!")

if not "":
    print("Falsy!")

if not []:
    print("Falsy!")
```

### Truthy Values

Everything else is truthy! Including:
- Non-zero numbers
- Non-empty strings
- Non-empty collections
- Any object (by default)

```python
# All of these print "Truthy!"
if 1:
    print("Truthy!")

if "hello":
    print("Truthy!")

if [1, 2, 3]:
    print("Truthy!")

if "0":  # String "0" is truthy!
    print("Truthy!")
```

This lets you write concise code:

```python
# Instead of this:
if len(my_list) > 0:
    process(my_list)

# You can write this:
if my_list:  # Truthy if list is not empty
    process(my_list)
```

## Practical Examples

Let's put it all together!

### Example 1: Calculate Discount

```python
price = 100
discount_percent = 20

discount = price * discount_percent / 100
final_price = price - discount

print(f"Original: ${price}")
print(f"Discount: ${discount}")
print(f"Final: ${final_price}")  # Final: $80.0
```

### Example 2: Check Eligibility

```python
age = 25
is_citizen = True
has_criminal_record = False

# Can vote if: adult, citizen, and no criminal record
can_vote = age >= 18 and is_citizen and not has_criminal_record
print(f"Can vote: {can_vote}")  # True
```

### Example 3: Temperature Warning

```python
temperature = 35

if temperature > 30:
    print("Hot!")
elif temperature < 10:
    print("Cold!")
else:
    print("Pleasant")

# Check if in comfortable range
is_comfortable = 15 <= temperature <= 25
print(f"Comfortable: {is_comfortable}")  # False
```

### Example 4: String Validation

```python
username = "alex_dev"

# Check username validity
is_long_enough = len(username) >= 3
has_no_spaces = " " not in username
is_alphanumeric = username.replace("_", "").isalnum()

is_valid = is_long_enough and has_no_spaces and is_alphanumeric
print(f"Valid username: {is_valid}")  # True
```

## Common Pitfalls to Avoid

### 1. Assignment vs. Comparison

```python
x = 5
if x = 5:  # ERROR! SyntaxError
    print("Five")

if x == 5:  # CORRECT!
    print("Five")
```

### 2. Division Types

```python
# Be aware of which division you're using
result1 = 7 / 2    # 3.5 (float)
result2 = 7 // 2   # 3 (integer)
```

### 3. Chaining Comparisons Incorrectly

```python
# This doesn't work as expected
x = 5
if x == 3 or 4:  # Always True! Evaluates as: (x == 3) or (4)
    print("Oops!")

# Correct way
if x == 3 or x == 4:
    print("Correct!")
```

### 4. String "0" is Truthy

```python
value = "0"
if value:  # True! String "0" is not empty
    print("Truthy!")

# Convert to int first if you want numeric behavior
if int(value):  # False! int 0 is falsy
    print("Won't print")
```

## Best Practices

1. **Use parentheses for clarity** - Make complex expressions easier to understand
2. **Prefer `is` for None checks** - `if x is None:` not `if x == None:`
3. **Use `in` for membership** - More Pythonic than loops
4. **Leverage truthy/falsy** - `if my_list:` instead of `if len(my_list) > 0:`
5. **Use compound operators** - `x += 1` instead of `x = x + 1`
6. **Chain comparisons** - `if 0 <= x <= 100:` instead of `if x >= 0 and x <= 100:`
7. **Watch division types** - Know when to use `/` vs `//`

## Practice Exercises

1. Write an expression to check if a number is divisible by both 3 and 5.

2. What's the difference between `5 / 2` and `5 // 2`?

3. Why does `"0"` evaluate to `True` in a boolean context?

4. Use the modulo operator to check if a year is a leap year (divisible by 4 but not 100, unless also divisible by 400).

5. Rewrite this using compound assignment: `total = total + price * quantity`

---

## Quiz Questions

### Question 1: Division Operators
**Type**: code-output
**Difficulty**: 2
**Tags**: python, operators, arithmetic, division

**Question Text**:
What is the output of this code?

```python
print(17 / 5, 17 // 5)
```

**Answer Choices**:
- A) 3.4 3
- B) 3 3
- C) 3.4 3.0
- D) 3.5 3

**Correct Answer**: A

**Explanation**:
In Python, `/` always performs regular division and returns a float (3.4), while `//` performs floor division and returns an integer by rounding down (3). So `17 / 5` gives 3.4 and `17 // 5` gives 3. This is an important distinction in Python 3—the `/` operator always returns a float, even if both operands are integers.

---

### Question 2: Modulo Operator
**Type**: code-output
**Difficulty**: 2
**Tags**: python, operators, modulo, arithmetic

**Question Text**:
What does this code print?

```python
print(23 % 5)
```

**Answer Choices**:
- A) 4
- B) 3
- C) 5
- D) 4.6

**Correct Answer**: B

**Explanation**:
The modulo operator `%` returns the remainder after division. 23 divided by 5 is 4 with a remainder of 3 (because 5 × 4 = 20, and 23 - 20 = 3). The modulo operator is commonly used to check if a number is even/odd (`n % 2`), to cycle through values, or to get the last digit of a number (`n % 10`).

---

### Question 3: Compound Assignment
**Type**: code-output
**Difficulty**: 3
**Tags**: python, operators, assignment, compound

**Question Text**:
What is the value of `x` after this code executes?

```python
x = 10
x += 5
x *= 2
x //= 3
```

**Answer Choices**:
- A) 10
- B) 15
- C) 20
- D) 10.0

**Correct Answer**: A

**Explanation**:
Let's trace through: Start with x=10. Then `x += 5` makes x=15. Then `x *= 2` makes x=30. Finally `x //= 3` performs floor division (30 // 3 = 10), making x=10. Note that `//=` uses floor division, which returns an integer when both operands are integers. The final value is 10.

---

### Question 4: Comparison Chaining
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: python, operators, comparison, chaining

**Question Text**:
Which of these is a valid way to check if `x` is between 10 and 20 (inclusive)?

**Answer Choices**:
- A) if 10 <= x <= 20:
- B) if x >= 10 and x <= 20:
- C) if 10 <= x and x <= 20:
- D) All of the above

**Correct Answer**: D

**Explanation**:
All three options are valid in Python! Option A uses Python's comparison chaining feature, which is unique and very readable. Options B and C use explicit `and` operators, which work but are more verbose. Comparison chaining `10 <= x <= 20` is the most Pythonic way. This feature makes Python code read more like natural mathematics.

---

### Question 5: Logical Operators
**Type**: code-output
**Difficulty**: 3
**Tags**: python, operators, logical, and, or

**Question Text**:
What does this print?

```python
print(True or False and False)
```

**Answer Choices**:
- A) True
- B) False
- C) Error
- D) None

**Correct Answer**: A

**Explanation**:
This evaluates to True because `and` has higher precedence than `or`. The expression is evaluated as `True or (False and False)`. First, `False and False` is evaluated to False. Then `True or False` evaluates to True. If you wanted a different result, you'd need parentheses: `(True or False) and False` would be False.

---

### Question 6: Identity vs Equality
**Type**: code-output
**Difficulty**: 5
**Tags**: python, operators, identity, is, equality

**Question Text**:
What is the output?

```python
a = [1, 2, 3]
b = [1, 2, 3]
print(a == b, a is b)
```

**Answer Choices**:
- A) True True
- B) True False
- C) False True
- D) False False

**Correct Answer**: B

**Explanation**:
The `==` operator checks value equality, and both lists have the same values [1, 2, 3], so `a == b` is True. However, `is` checks identity (whether they're the same object in memory). Even though `a` and `b` have the same contents, they are two different list objects in memory, so `a is b` is False. Use `==` for comparing values, use `is` mainly for checking `None`.

---

### Question 7: Truthy and Falsy
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: python, operators, truthy, falsy, bool

**Question Text**:
Which of the following values is truthy in Python?

**Answer Choices**:
- A) 0
- B) ""
- C) "0"
- D) []

**Correct Answer**: C

**Explanation**:
The string "0" is truthy because it's a non-empty string! In Python, falsy values include: False, None, 0 (any zero), empty string "", empty list [], empty tuple (), empty dict {}, and empty set(). Everything else is truthy, including the string "0", the string "False", and any non-empty collection. This is a common gotcha—the string "0" looks like zero but behaves as truthy!

---

### Question 8: Membership Operators
**Type**: code-output
**Difficulty**: 3
**Tags**: python, operators, membership, in

**Question Text**:
What does this code print?

```python
text = "Python"
print("py" in text)
```

**Answer Choices**:
- A) True
- B) False
- C) Error
- D) None

**Correct Answer**: B

**Explanation**:
This prints False because string membership in Python is case-sensitive. The string "Python" contains "Py" (capital P), not "py" (lowercase p). If you want case-insensitive checking, you'd need to convert both to the same case: `"py" in text.lower()` would return True. The `in` operator is very useful for checking if a substring exists in a string or if an element exists in a list.
