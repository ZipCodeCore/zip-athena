# Python Functions

Welcome! Functions are one of the most powerful tools in programming. They let you organize code into reusable blocks, making your programs cleaner, more maintainable, and easier to understand. Think of functions as mini-programs within your program—each one does a specific job. Let's explore everything Python functions have to offer!

## What Are Functions?

A function is a reusable block of code that performs a specific task. Instead of writing the same code over and over, you write it once in a function and call it whenever you need it.

```python
# Without a function (repetitive!)
print("Hello, Alice!")
print("Welcome to our program")
print()
print("Hello, Bob!")
print("Welcome to our program")
print()

# With a function (reusable!)
def greet(name):
    print(f"Hello, {name}!")
    print("Welcome to our program")
    print()

greet("Alice")
greet("Bob")
```

Functions make your code DRY: **D**on't **R**epeat **Y**ourself!

## Defining Functions with `def`

Create a function using the `def` keyword:

```python
# Basic function definition
def say_hello():
    print("Hello!")

# Call the function
say_hello()  # Output: Hello!
```

The syntax is:
```python
def function_name(parameters):
    # function body
    # code to execute
    return value  # optional
```

**Naming conventions**: Use snake_case for function names (lowercase with underscores), just like variables.

```python
# Good function names
def calculate_total():
    pass

def send_email():
    pass

def get_user_info():
    pass

# Bad (but valid)
def CalculateTotal():  # Should be snake_case
    pass

def sendEmail():  # Should be snake_case
    pass
```

## Function Parameters and Arguments

Parameters let you pass data into functions. They make functions flexible and reusable!

### Positional Parameters

The order matters:

```python
def greet(name, greeting):
    print(f"{greeting}, {name}!")

greet("Alex", "Hello")     # Hello, Alex!
greet("Hello", "Alex")     # Alex, Hello! (wrong order!)
```

### Multiple Parameters

```python
def calculate_rectangle_area(width, height):
    area = width * height
    return area

result = calculate_rectangle_area(5, 3)
print(result)  # 15
```

### The Difference: Parameters vs Arguments

- **Parameters**: Variables in the function definition
- **Arguments**: Values you pass when calling the function

```python
def greet(name):     # 'name' is a parameter
    print(f"Hello, {name}!")

greet("Alex")        # "Alex" is an argument
```

## Return Values and the `return` Statement

Functions can send data back using `return`:

```python
def add(a, b):
    return a + b

result = add(5, 3)
print(result)  # 8
```

### Functions Without `return`

If a function doesn't explicitly return a value, it returns `None`:

```python
def say_hello():
    print("Hello!")

result = say_hello()  # Prints: Hello!
print(result)         # Prints: None
```

### Returning Multiple Values

Python lets you return multiple values as a tuple:

```python
def get_user_info():
    name = "Alex"
    age = 25
    city = "New York"
    return name, age, city

# Unpack the returned values
user_name, user_age, user_city = get_user_info()
print(user_name)  # Alex
print(user_age)   # 25

# Or keep as a tuple
info = get_user_info()
print(info)  # ('Alex', 25, 'New York')
```

### Early Return

You can return from a function at any point:

```python
def check_positive(number):
    if number <= 0:
        return False  # Exit early!
    return True

print(check_positive(5))   # True
print(check_positive(-3))  # False
```

## Default Parameter Values

Give parameters default values for optional arguments:

```python
def greet(name, greeting="Hello"):
    print(f"{greeting}, {name}!")

greet("Alex")              # Hello, Alex! (uses default)
greet("Alex", "Hi")        # Hi, Alex! (overrides default)
greet("Alex", "Welcome")   # Welcome, Alex!
```

### Multiple Defaults

```python
def create_user(name, age=0, city="Unknown"):
    return {
        "name": name,
        "age": age,
        "city": city
    }

print(create_user("Alex"))
# {'name': 'Alex', 'age': 0, 'city': 'Unknown'}

print(create_user("Alex", 25))
# {'name': 'Alex', 'age': 25, 'city': 'Unknown'}

print(create_user("Alex", 25, "New York"))
# {'name': 'Alex', 'age': 25, 'city': 'New York'}
```

**Important**: Default parameters must come after non-default parameters!

```python
# Wrong!
def greet(greeting="Hello", name):  # ERROR!
    print(f"{greeting}, {name}!")

# Correct!
def greet(name, greeting="Hello"):
    print(f"{greeting}, {name}!")
```

## Keyword Arguments

Specify arguments by parameter name, regardless of order:

```python
def describe_pet(animal_type, pet_name):
    print(f"I have a {animal_type} named {pet_name}.")

# Positional arguments
describe_pet("dog", "Buddy")  # I have a dog named Buddy.

# Keyword arguments (order doesn't matter!)
describe_pet(pet_name="Buddy", animal_type="dog")  # Same result

# Mix positional and keyword (positional first!)
describe_pet("dog", pet_name="Buddy")
```

Keyword arguments make function calls clearer and more readable:

```python
# Hard to understand
create_account("Alex", 25, True, False, "premium")

# Much clearer!
create_account(
    name="Alex",
    age=25,
    is_active=True,
    is_verified=False,
    tier="premium"
)
```

## Variable-Length Arguments

Sometimes you don't know how many arguments you'll receive. Python has special syntax for this!

### `*args` - Variable Positional Arguments

Collect any number of positional arguments into a tuple:

```python
def sum_all(*numbers):
    total = 0
    for num in numbers:
        total += num
    return total

print(sum_all(1, 2, 3))           # 6
print(sum_all(1, 2, 3, 4, 5))     # 15
print(sum_all(10, 20))            # 30
```

The name `args` is just a convention—you can use any name (but always include the `*`):

```python
def print_names(*names):
    for name in names:
        print(f"Hello, {name}!")

print_names("Alex", "Jordan", "Casey")
# Hello, Alex!
# Hello, Jordan!
# Hello, Casey!
```

### `**kwargs` - Variable Keyword Arguments

Collect any number of keyword arguments into a dictionary:

```python
def print_info(**info):
    for key, value in info.items():
        print(f"{key}: {value}")

print_info(name="Alex", age=25, city="New York")
# name: Alex
# age: 25
# city: New York
```

### Combining All Parameter Types

You can use all types together, but they must be in this order:
1. Regular positional parameters
2. `*args`
3. Keyword-only parameters (after `*args`)
4. `**kwargs`

```python
def complex_function(a, b, *args, option=None, **kwargs):
    print(f"a: {a}, b: {b}")
    print(f"args: {args}")
    print(f"option: {option}")
    print(f"kwargs: {kwargs}")

complex_function(1, 2, 3, 4, 5, option="test", x=10, y=20)
# a: 1, b: 2
# args: (3, 4, 5)
# option: test
# kwargs: {'x': 10, 'y': 20}
```

## Function Documentation Strings (Docstrings)

Docstrings document what your function does. They're the first string in a function, usually triple-quoted:

```python
def calculate_area(width, height):
    """
    Calculate the area of a rectangle.
    
    Parameters:
        width (float): The width of the rectangle
        height (float): The height of the rectangle
    
    Returns:
        float: The area of the rectangle
    """
    return width * height

# Access the docstring
print(calculate_area.__doc__)
```

### Good Docstring Example

```python
def greet(name, greeting="Hello"):
    """
    Greet a person with a custom or default greeting.
    
    Args:
        name (str): The name of the person to greet
        greeting (str, optional): The greeting to use. Defaults to "Hello".
    
    Returns:
        None: Prints the greeting to console
    
    Example:
        >>> greet("Alex")
        Hello, Alex!
        >>> greet("Alex", "Hi")
        Hi, Alex!
    """
    print(f"{greeting}, {name}!")
```

**Best practice**: Always write docstrings for functions that others will use!

## Scope and Namespaces

Scope determines where variables are accessible. Python has the **LEGB** rule:
- **L**ocal: Inside the current function
- **E**nclosing: Inside enclosing functions
- **G**lobal: Module level
- **B**uilt-in: Python built-ins

### Local Scope

Variables created inside a function are local:

```python
def my_function():
    x = 10  # Local variable
    print(x)

my_function()  # 10
print(x)       # ERROR! NameError: 'x' is not defined
```

### Global Scope

Variables created outside functions are global:

```python
x = 10  # Global variable

def my_function():
    print(x)  # Can read global variables

my_function()  # 10
print(x)       # 10
```

### Enclosing Scope

Variables in outer functions:

```python
def outer():
    x = 10  # Enclosing variable
    
    def inner():
        print(x)  # Can access enclosing scope
    
    inner()

outer()  # 10
```

## The `global` and `nonlocal` Keywords

### `global` Keyword

Modify global variables from inside a function:

```python
count = 0  # Global variable

def increment():
    global count  # Declare we're modifying the global variable
    count += 1

increment()
print(count)  # 1

increment()
print(count)  # 2
```

**Without `global`**:

```python
count = 0

def increment():
    count = count + 1  # ERROR! UnboundLocalError

increment()
```

### `nonlocal` Keyword

Modify variables in enclosing (but not global) scope:

```python
def outer():
    count = 0
    
    def inner():
        nonlocal count  # Refer to enclosing scope's count
        count += 1
    
    inner()
    print(count)  # 1
    inner()
    print(count)  # 2

outer()
```

**Best practice**: Minimize use of `global` and `nonlocal`. They make code harder to understand. Consider using parameters and return values instead!

## Lambda Functions (Anonymous Functions)

Lambda functions are small, anonymous functions defined in one line:

```python
# Regular function
def square(x):
    return x ** 2

# Lambda equivalent
square = lambda x: x ** 2

print(square(5))  # 25
```

### Syntax

```python
lambda parameters: expression
```

### Common Use Cases

Lambda functions are often used with `map()`, `filter()`, and `sorted()`:

```python
# With map()
numbers = [1, 2, 3, 4, 5]
squared = list(map(lambda x: x ** 2, numbers))
print(squared)  # [1, 4, 9, 16, 25]

# With filter()
numbers = [1, 2, 3, 4, 5, 6, 7, 8]
evens = list(filter(lambda x: x % 2 == 0, numbers))
print(evens)  # [2, 4, 6, 8]

# With sorted()
pairs = [(1, 'one'), (3, 'three'), (2, 'two')]
sorted_pairs = sorted(pairs, key=lambda pair: pair[1])
print(sorted_pairs)  # [(1, 'one'), (3, 'three'), (2, 'two')]
```

### Multiple Parameters

```python
# Lambda with multiple parameters
add = lambda x, y: x + y
print(add(3, 5))  # 8

# With map and multiple iterables
nums1 = [1, 2, 3]
nums2 = [4, 5, 6]
result = list(map(lambda x, y: x + y, nums1, nums2))
print(result)  # [5, 7, 9]
```

**Best practice**: Use lambda for simple functions. For complex logic, use `def`!

## Functions as First-Class Objects

In Python, functions are first-class objects. You can:
- Assign them to variables
- Pass them as arguments
- Return them from other functions
- Store them in data structures

### Assigning Functions to Variables

```python
def greet(name):
    return f"Hello, {name}!"

# Assign function to a variable
say_hello = greet
print(say_hello("Alex"))  # Hello, Alex!
```

### Passing Functions as Arguments

```python
def apply_operation(func, x, y):
    return func(x, y)

def add(a, b):
    return a + b

def multiply(a, b):
    return a * b

print(apply_operation(add, 5, 3))       # 8
print(apply_operation(multiply, 5, 3))  # 15
```

### Returning Functions

```python
def make_multiplier(n):
    def multiplier(x):
        return x * n
    return multiplier

times_three = make_multiplier(3)
times_five = make_multiplier(5)

print(times_three(10))  # 30
print(times_five(10))   # 50
```

### Storing in Data Structures

```python
def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

operations = {
    'add': add,
    'subtract': subtract
}

result = operations['add'](10, 5)
print(result)  # 15
```

## Recursion

A recursive function is one that calls itself. Recursion is powerful for solving problems that can be broken down into smaller, similar sub-problems.

### Basic Recursion

```python
def countdown(n):
    if n <= 0:  # Base case
        print("Blastoff!")
    else:
        print(n)
        countdown(n - 1)  # Recursive call

countdown(5)
# 5
# 4
# 3
# 2
# 1
# Blastoff!
```

Every recursive function needs:
1. **Base case**: When to stop recursing
2. **Recursive case**: When to call itself

### Classic Example: Factorial

```python
def factorial(n):
    # Base case
    if n == 0 or n == 1:
        return 1
    # Recursive case
    return n * factorial(n - 1)

print(factorial(5))  # 120 (5 * 4 * 3 * 2 * 1)
```

### Fibonacci Sequence

```python
def fibonacci(n):
    # Base cases
    if n <= 0:
        return 0
    if n == 1:
        return 1
    # Recursive case
    return fibonacci(n - 1) + fibonacci(n - 2)

# First 10 Fibonacci numbers
for i in range(10):
    print(fibonacci(i), end=" ")
# Output: 0 1 1 2 3 5 8 13 21 34
```

### Recursion vs Iteration

Many recursive problems can be solved with loops:

```python
# Recursive factorial
def factorial_recursive(n):
    if n <= 1:
        return 1
    return n * factorial_recursive(n - 1)

# Iterative factorial
def factorial_iterative(n):
    result = 1
    for i in range(1, n + 1):
        result *= i
    return result
```

**When to use recursion**: Problems naturally suited to it (tree traversal, divide-and-conquer). Otherwise, iteration is often clearer and more efficient.

**Warning**: Deep recursion can cause stack overflow! Python has a recursion limit (usually ~1000 calls).

## Practical Examples

### Example 1: Temperature Converter

```python
def celsius_to_fahrenheit(celsius):
    """Convert Celsius to Fahrenheit."""
    return (celsius * 9/5) + 32

def fahrenheit_to_celsius(fahrenheit):
    """Convert Fahrenheit to Celsius."""
    return (fahrenheit - 32) * 5/9

print(f"100°C = {celsius_to_fahrenheit(100)}°F")  # 212.0°F
print(f"98.6°F = {fahrenheit_to_celsius(98.6):.1f}°C")  # 37.0°C
```

### Example 2: Input Validation

```python
def get_valid_age():
    """Keep asking until user enters a valid age."""
    while True:
        try:
            age = int(input("Enter your age: "))
            if 0 < age < 120:
                return age
            else:
                print("Age must be between 1 and 119")
        except ValueError:
            print("Please enter a valid number")

user_age = get_valid_age()
print(f"Your age: {user_age}")
```

### Example 3: Data Processing

```python
def filter_and_process(numbers, threshold=10):
    """
    Filter numbers above threshold and return their squares.
    
    Args:
        numbers (list): List of numbers to process
        threshold (int): Minimum value to include
    
    Returns:
        list: Squared values of numbers above threshold
    """
    filtered = [n for n in numbers if n > threshold]
    squared = [n ** 2 for n in filtered]
    return squared

data = [5, 12, 8, 20, 3, 15]
result = filter_and_process(data, threshold=10)
print(result)  # [144, 400, 225]
```

### Example 4: Decorator Pattern (Advanced)

```python
def timer(func):
    """Decorator to measure function execution time."""
    import time
    
    def wrapper(*args, **kwargs):
        start = time.time()
        result = func(*args, **kwargs)
        end = time.time()
        print(f"{func.__name__} took {end - start:.4f} seconds")
        return result
    
    return wrapper

@timer
def slow_function():
    import time
    time.sleep(1)
    return "Done!"

result = slow_function()
# Output: slow_function took 1.0001 seconds
```

## Common Pitfalls

### 1. Mutable Default Arguments

```python
# DON'T DO THIS!
def add_item(item, items=[]):
    items.append(item)
    return items

print(add_item(1))  # [1]
print(add_item(2))  # [1, 2] (unexpected!)

# DO THIS INSTEAD
def add_item(item, items=None):
    if items is None:
        items = []
    items.append(item)
    return items
```

### 2. Modifying Global Variables

```python
# Confusing!
count = 0

def increment():
    global count
    count += 1

# Better!
def increment(count):
    return count + 1

count = increment(count)
```

### 3. Forgetting to Return

```python
def add(a, b):
    result = a + b  # Forgot return!

answer = add(5, 3)
print(answer)  # None (oops!)
```

## Best Practices

1. **Function should do one thing** - Single responsibility principle
2. **Use descriptive names** - `calculate_total()` not `calc()`
3. **Keep functions short** - If it's too long, break it up
4. **Avoid side effects** - Functions should be predictable
5. **Write docstrings** - Document parameters and return values
6. **Use type hints** (Python 3.5+) - Makes code clearer
7. **Return early** - Reduce nesting with early returns
8. **Avoid mutable defaults** - Use `None` instead

## Practice Exercises

1. Write a function that takes a list of numbers and returns the sum of even numbers only.

2. Create a function that checks if a string is a palindrome (reads same forwards and backwards).

3. Write a recursive function to calculate the sum of digits of a number (e.g., 123 → 6).

4. Create a function that accepts `*args` and returns a dictionary counting how many times each value appears.

5. Write a function that takes a function as an argument and applies it to a list of values.

---

## Quiz Questions

### Question 1: Basic Function Definition
**Type**: code-output
**Difficulty**: 2
**Tags**: python, functions, return

**Question Text**:
What does this code print?

```python
def multiply(a, b):
    result = a * b

answer = multiply(3, 4)
print(answer)
```

**Answer Choices**:
- A) 12
- B) None
- C) 7
- D) Error

**Correct Answer**: B

**Explanation**:
The function `multiply` calculates `result` but doesn't return it. Functions without an explicit `return` statement return `None`. To get 12, you'd need to add `return result` or `return a * b`. This is a common mistake—always remember to return values you want to use outside the function!

---

### Question 2: Default Parameters
**Type**: code-output
**Difficulty**: 3
**Tags**: python, functions, parameters, defaults

**Question Text**:
What is the output?

```python
def greet(name, greeting="Hello"):
    return f"{greeting}, {name}!"

print(greet("Alex"))
print(greet("Alex", "Hi"))
```

**Answer Choices**:
- A) Hello, Alex! and Hello, Alex!
- B) Hello, Alex! and Hi, Alex!
- C) Hi, Alex! and Hi, Alex!
- D) Error

**Correct Answer**: B

**Explanation**:
Default parameters provide fallback values. The first call `greet("Alex")` uses the default greeting "Hello", printing "Hello, Alex!". The second call `greet("Alex", "Hi")` overrides the default with "Hi", printing "Hi, Alex!". Default parameters must come after required parameters in the function definition.

---

### Question 3: Variable-Length Arguments (*args)
**Type**: code-output
**Difficulty**: 4
**Tags**: python, functions, args, parameters

**Question Text**:
What does this print?

```python
def sum_all(*numbers):
    total = 0
    for num in numbers:
        total += num
    return total

print(sum_all(1, 2, 3, 4))
```

**Answer Choices**:
- A) 1234
- B) 10
- C) [1, 2, 3, 4]
- D) Error

**Correct Answer**: B

**Explanation**:
The `*numbers` parameter collects all positional arguments into a tuple: (1, 2, 3, 4). The function then sums them: 1 + 2 + 3 + 4 = 10. The `*args` syntax is useful when you don't know how many arguments will be passed. You can call `sum_all(1, 2)` or `sum_all(1, 2, 3, 4, 5)` and it works with any number of arguments.

---

### Question 4: Scope and Global
**Type**: code-output
**Difficulty**: 5
**Tags**: python, functions, scope, global

**Question Text**:
What is printed?

```python
x = 10

def modify():
    x = 5
    print(x)

modify()
print(x)
```

**Answer Choices**:
- A) 5 and 5
- B) 5 and 10
- C) 10 and 5
- D) 10 and 10

**Correct Answer**: B

**Explanation**:
Inside `modify()`, `x = 5` creates a new local variable `x` that shadows the global `x`. So the function prints 5. The global `x` remains unchanged at 10, so the second print shows 10. To modify the global `x`, you'd need `global x` before the assignment. This demonstrates Python's local vs global scope—local variables don't affect global ones unless you explicitly declare them as global.

---

### Question 5: Lambda Functions
**Type**: code-output
**Difficulty**: 4
**Tags**: python, functions, lambda

**Question Text**:
What is the output?

```python
square = lambda x: x ** 2
print(square(5))
```

**Answer Choices**:
- A) 10
- B) 25
- C) x ** 2
- D) Error

**Correct Answer**: B

**Explanation**:
Lambda functions are anonymous functions defined in one line. `lambda x: x ** 2` creates a function that takes one parameter `x` and returns `x ** 2`. Calling `square(5)` evaluates to 5 ** 2 = 25. Lambda functions are useful for simple operations, especially with `map()`, `filter()`, and `sorted()`. For complex logic, use a regular `def` function instead.

---

### Question 6: Keyword Arguments
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: python, functions, arguments, keyword

**Question Text**:
Which function call is valid?

```python
def describe_pet(animal_type, pet_name):
    return f"{pet_name} is a {animal_type}"
```

**Answer Choices**:
- A) describe_pet(pet_name="Buddy", "dog")
- B) describe_pet("dog", pet_name="Buddy")
- C) describe_pet(animal_type="dog", animal_type="cat")
- D) describe_pet(pet_name="Buddy")

**Correct Answer**: B

**Explanation**:
When mixing positional and keyword arguments, positional arguments must come first. Option B is correct: "dog" is positional, `pet_name="Buddy"` is keyword. Option A fails because you can't have positional arguments after keyword ones. Option C duplicates `animal_type`. Option D is missing the required `animal_type` parameter. Keyword arguments make function calls more readable but must follow the positional arguments.

---

### Question 7: Recursion
**Type**: code-output
**Difficulty**: 5
**Tags**: python, functions, recursion

**Question Text**:
What does this function return when called with `factorial(4)`?

```python
def factorial(n):
    if n <= 1:
        return 1
    return n * factorial(n - 1)
```

**Answer Choices**:
- A) 10
- B) 24
- C) 4
- D) Error

**Correct Answer**: B

**Explanation**:
This recursive function calculates factorial. `factorial(4)` returns `4 * factorial(3)`, which returns `4 * 3 * factorial(2)`, which returns `4 * 3 * 2 * factorial(1)`, which returns `4 * 3 * 2 * 1 = 24`. The base case (`if n <= 1: return 1`) stops the recursion. Recursion works by breaking a problem into smaller instances of itself until reaching a base case.

---

### Question 8: Functions as First-Class Objects
**Type**: code-output
**Difficulty**: 6
**Tags**: python, functions, first-class, higher-order

**Question Text**:
What is printed?

```python
def apply(func, x):
    return func(x)

result = apply(lambda n: n * 2, 5)
print(result)
```

**Answer Choices**:
- A) 5
- B) 10
- C) lambda n: n * 2
- D) Error

**Correct Answer**: B

**Explanation**:
Functions in Python are first-class objects, meaning they can be passed as arguments. The `apply` function takes a function `func` and a value `x`, then calls `func(x)`. We pass the lambda function `lambda n: n * 2` and value 5. The lambda is called with 5, returning 5 * 2 = 10. This demonstrates higher-order functions—functions that take or return other functions, a powerful functional programming concept.
