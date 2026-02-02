# Python Exception Handling

Hey there! Let's dive into one of the most important topics in Python—exception handling. You know how life throws curveballs at you? Well, your programs will face unexpected situations too. A file might not exist, a user might enter text when you expected a number, or you might try to divide by zero. Exception handling is how you catch these curveballs gracefully instead of letting your program crash and burn.

## 13.1 Understanding Exceptions

An exception is Python's way of saying "Houston, we have a problem!" When something goes wrong during program execution, Python creates an exception object that contains information about what happened. If you don't handle it, your program crashes with an error message.

Think of exceptions like fire alarms in a building. When smoke is detected, the alarm goes off—that's Python raising an exception. You can either ignore it (and let everything burn down), or you can respond appropriately (handle it gracefully).

Here's what happens when an exception occurs without handling:

```python
# This will crash!
number = int("hello")  # Can't convert "hello" to an integer
print("This line never runs")
```

When you run this, you'll see something like:
```
ValueError: invalid literal for int() with base 10: 'hello'
```

The program stops immediately, and nothing after the error runs. Not ideal if you're writing software that needs to stay running!

Exceptions aren't always bad, though. They're actually Python's way of helping you identify and fix problems. They provide valuable information about what went wrong and where, making debugging much easier.

## 13.2 Common Built-in Exceptions

Python has dozens of built-in exceptions, each for different types of errors. Let's look at the most common ones you'll encounter:

### SyntaxError

This happens when Python can't understand your code because you've broken the language rules:

```python
# Missing colon
if True
    print("Oops")

# This gives: SyntaxError: invalid syntax
```

**Important:** SyntaxError is special—it's caught at parse time, before your code even runs, so you can't catch it with try-except in the same code block.

### TypeError

This occurs when you use an operation on an inappropriate type:

```python
# Can't add a string and an integer
result = "5" + 5  # TypeError: can only concatenate str (not "int") to str

# Can't call a non-callable object
my_var = 42
my_var()  # TypeError: 'int' object is not callable
```

### ValueError

Raised when a function receives an argument of the right type but inappropriate value:

```python
# Right type (string), wrong value (not a number)
number = int("xyz")  # ValueError: invalid literal for int()

# Unpacking with wrong number of values
a, b = [1, 2, 3]  # ValueError: too many values to unpack
```

### KeyError

Happens when you try to access a dictionary key that doesn't exist:

```python
person = {"name": "Alice", "age": 30}
city = person["city"]  # KeyError: 'city'
```

### IndexError

Raised when you try to access a list index that's out of range:

```python
fruits = ["apple", "banana"]
print(fruits[5])  # IndexError: list index out of range
```

### FileNotFoundError

Occurs when trying to open a file that doesn't exist:

```python
with open("nonexistent_file.txt") as f:
    content = f.read()  # FileNotFoundError: [Errno 2] No such file or directory
```

### ZeroDivisionError

Raised when dividing by zero:

```python
result = 10 / 0  # ZeroDivisionError: division by zero
```

### AttributeError

Happens when you try to access an attribute or method that doesn't exist:

```python
text = "hello"
text.nonexistent_method()  # AttributeError: 'str' object has no attribute 'nonexistent_method'
```

## 13.3 Try-Except Blocks

Now for the good stuff—how to actually handle these exceptions! The `try-except` block is your safety net. You try to do something risky, and if it fails, you catch the exception and handle it gracefully.

```python
try:
    age = int(input("Enter your age: "))
    print(f"Next year you'll be {age + 1}")
except ValueError:
    print("That's not a valid number! Please enter digits only.")
```

The structure is simple:
- **try block**: Contains the code that might raise an exception
- **except block**: Contains the code that runs if an exception occurs

Here's a practical example with file handling:

```python
try:
    with open("data.txt", "r") as file:
        data = file.read()
        print(data)
except FileNotFoundError:
    print("Sorry, the file doesn't exist. Creating a new one...")
    with open("data.txt", "w") as file:
        file.write("Default content")
```

You can also catch the exception object itself to get more information:

```python
try:
    result = 10 / 0
except ZeroDivisionError as e:
    print(f"Error occurred: {e}")
    print(f"Error type: {type(e).__name__}")
```

## 13.4 Multiple Except Clauses

Often, different exceptions require different responses. You can have multiple `except` blocks to handle various exceptions differently:

```python
def divide_numbers(a, b):
    try:
        result = int(a) / int(b)
        return result
    except ValueError:
        print("Please provide valid numbers!")
        return None
    except ZeroDivisionError:
        print("Cannot divide by zero!")
        return None
    except TypeError:
        print("Wrong data type provided!")
        return None

print(divide_numbers("10", "2"))   # Works: 5.0
print(divide_numbers("abc", "2"))  # ValueError
print(divide_numbers("10", "0"))   # ZeroDivisionError
```

You can also catch multiple exceptions in one block if you want to handle them the same way:

```python
try:
    value = int(input("Enter a number: "))
    result = 100 / value
    items = [1, 2, 3]
    print(items[value])
except (ValueError, ZeroDivisionError, IndexError) as e:
    print(f"Something went wrong: {type(e).__name__}")
    print(f"Details: {e}")
```

A bare `except` catches all exceptions, but use it sparingly—it can hide bugs:

```python
try:
    # Some risky operation
    risky_operation()
except Exception as e:  # Better than bare except
    print(f"An error occurred: {e}")

# Avoid this unless you have a really good reason:
try:
    risky_operation()
except:  # Catches EVERYTHING, even KeyboardInterrupt!
    print("Something went wrong")
```

**Best practice:** Always catch specific exceptions when possible. It makes your code more maintainable and helps you understand what can go wrong.

## 13.5 The Else Clause with Try-Except

The `else` clause in a try-except block is a hidden gem! It runs only if the try block completes successfully without raising any exceptions. This keeps your try block focused on just the risky code:

```python
try:
    age = int(input("Enter your age: "))
except ValueError:
    print("Invalid input!")
else:
    print(f"Great! You are {age} years old.")
    if age >= 18:
        print("You're an adult!")
```

The else block is useful for code that should only run if the try block succeeded but doesn't need exception handling itself:

```python
def read_config_file(filename):
    try:
        file = open(filename, 'r')
    except FileNotFoundError:
        print(f"Config file {filename} not found. Using defaults.")
        return {}
    else:
        # Only runs if file opened successfully
        with file:
            config = file.read()
            print(f"Config loaded from {filename}")
            return config

config = read_config_file("settings.conf")
```

Think of the `else` clause as saying "if everything went fine, do this extra stuff." It makes your code clearer by separating the risky operation from the success-case logic.

## 13.6 The Finally Clause

The `finally` block is like that friend who always shows up, no matter what. It executes whether an exception occurred or not, making it perfect for cleanup operations:

```python
try:
    file = open("data.txt", "r")
    data = file.read()
    print(data)
except FileNotFoundError:
    print("File not found!")
finally:
    # This ALWAYS runs
    print("Cleanup complete!")
```

The `finally` block is commonly used to release resources like file handles, database connections, or network sockets:

```python
file = None
try:
    file = open("important_data.txt", "r")
    # Process the file
    data = file.read()
    result = int(data)  # Might raise ValueError
    print(f"Result: {result}")
except FileNotFoundError:
    print("File not found!")
except ValueError:
    print("File doesn't contain a valid number!")
finally:
    # Ensure the file is closed, even if an exception occurred
    if file:
        file.close()
        print("File closed.")
```

You can combine try, except, else, and finally all together:

```python
def divide_and_log(a, b):
    try:
        result = a / b
    except ZeroDivisionError:
        print("Cannot divide by zero!")
        return None
    else:
        print(f"Division successful: {a} / {b} = {result}")
        return result
    finally:
        print("Operation complete.")

divide_and_log(10, 2)  # Success case
divide_and_log(10, 0)  # Exception case
```

**Pro tip:** With context managers (the `with` statement), you often don't need finally for cleanup—the context manager handles it automatically!

## 13.7 Raising Exceptions with Raise

Sometimes you need to raise exceptions yourself. Maybe a user provided invalid input, or a function received inappropriate arguments. The `raise` keyword lets you throw exceptions deliberately:

```python
def set_age(age):
    if age < 0:
        raise ValueError("Age cannot be negative!")
    if age > 150:
        raise ValueError("Age seems unrealistic!")
    print(f"Age set to {age}")

set_age(25)    # Works fine
set_age(-5)    # Raises ValueError
```

You can raise any built-in exception or create your own. When raising, you can provide a custom error message:

```python
def withdraw_money(balance, amount):
    if amount > balance:
        raise ValueError(f"Insufficient funds! You have ${balance} but tried to withdraw ${amount}")
    return balance - amount

try:
    new_balance = withdraw_money(100, 150)
except ValueError as e:
    print(f"Transaction failed: {e}")
```

Sometimes you want to catch an exception, do something with it, and then re-raise it:

```python
def process_data(filename):
    try:
        with open(filename, 'r') as f:
            data = f.read()
        # Process data
        return data
    except FileNotFoundError:
        print(f"Logging: File {filename} not found")
        raise  # Re-raise the same exception

try:
    process_data("missing.txt")
except FileNotFoundError:
    print("Handled at higher level")
```

You can also raise a different exception type:

```python
def get_user_age(user_id):
    try:
        # Simulate database lookup
        ages = {"user1": 25, "user2": 30}
        return ages[user_id]
    except KeyError:
        # Convert KeyError to a more meaningful exception
        raise ValueError(f"User {user_id} not found in database")
```

## 13.8 Creating Custom Exceptions

For complex applications, you'll want to create your own exception types. This makes your code more expressive and easier to debug. Custom exceptions are just classes that inherit from `Exception`:

```python
class InsufficientFundsError(Exception):
    """Raised when a withdrawal exceeds account balance"""
    pass

class InvalidAccountError(Exception):
    """Raised when account number is invalid"""
    pass

def withdraw(balance, amount):
    if amount > balance:
        raise InsufficientFundsError(f"Cannot withdraw ${amount}. Balance: ${balance}")
    return balance - amount

try:
    new_balance = withdraw(100, 150)
except InsufficientFundsError as e:
    print(f"Transaction denied: {e}")
```

You can add extra functionality to custom exceptions:

```python
class ValidationError(Exception):
    """Base class for validation errors"""
    
    def __init__(self, field, message):
        self.field = field
        self.message = message
        super().__init__(f"{field}: {message}")

class EmailValidationError(ValidationError):
    """Raised for invalid email addresses"""
    
    def __init__(self, email):
        super().__init__("email", f"'{email}' is not a valid email address")

def register_user(email):
    if "@" not in email:
        raise EmailValidationError(email)
    print(f"User registered with email: {email}")

try:
    register_user("invalid-email")
except EmailValidationError as e:
    print(f"Registration failed - {e.field}: {e.message}")
```

Custom exceptions can have their own attributes and methods:

```python
class GameOverError(Exception):
    """Raised when the game ends"""
    
    def __init__(self, score, level):
        self.score = score
        self.level = level
        message = f"Game Over! Score: {score}, Level: {level}"
        super().__init__(message)
    
    def save_high_score(self):
        print(f"Saving high score: {self.score}")

try:
    # Game logic
    raise GameOverError(score=1500, level=5)
except GameOverError as e:
    print(e)
    e.save_high_score()
```

**Best practice:** Create custom exceptions for domain-specific errors in your application. It makes error handling clearer and more maintainable.

## 13.9 Exception Chaining

Exception chaining helps you preserve the original error context when raising a new exception. Python 3 provides two ways to chain exceptions: implicit chaining and explicit chaining.

**Implicit chaining** happens automatically when an exception is raised inside an except block:

```python
def divide(a, b):
    try:
        return a / b
    except ZeroDivisionError:
        # This creates an implicit chain
        raise ValueError("Invalid operation attempted")

try:
    result = divide(10, 0)
except ValueError as e:
    print(f"Error: {e}")
    print(f"Original cause: {e.__cause__}")
    print(f"Context: {e.__context__}")
```

**Explicit chaining** uses `raise ... from` to show that one exception was directly caused by another:

```python
class DatabaseError(Exception):
    pass

def get_user_data(user_id):
    try:
        # Simulate database operation
        users = {}
        return users[user_id]
    except KeyError as e:
        # Explicitly chain the exceptions
        raise DatabaseError(f"User {user_id} not found") from e

try:
    get_user_data("user123")
except DatabaseError as e:
    print(f"Database error: {e}")
    print(f"Caused by: {e.__cause__}")
```

You can suppress the original exception using `raise ... from None`:

```python
def process_config(filename):
    try:
        with open(filename, 'r') as f:
            return f.read()
    except FileNotFoundError:
        # Suppress the original exception - only show our message
        raise ValueError(f"Configuration file '{filename}' is required") from None

try:
    config = process_config("missing.conf")
except ValueError as e:
    print(e)  # Only shows our message, no FileNotFoundError details
```

Here's a practical example showing why chaining matters:

```python
class DataProcessingError(Exception):
    """Custom exception for data processing errors"""
    pass

def process_json_file(filename):
    import json
    
    try:
        with open(filename, 'r') as f:
            data = json.load(f)
        return data
    except FileNotFoundError as e:
        raise DataProcessingError(f"Cannot process {filename}") from e
    except json.JSONDecodeError as e:
        raise DataProcessingError(f"Invalid JSON in {filename}") from e

try:
    data = process_json_file("data.json")
except DataProcessingError as e:
    print(f"Error: {e}")
    if e.__cause__:
        print(f"Root cause: {type(e.__cause__).__name__}: {e.__cause__}")
```

Exception chaining preserves the full error context, making debugging much easier!

## 13.10 Best Practices in Exception Handling

Let's talk about how to use exception handling effectively. Following these practices will make your code more robust, maintainable, and easier to debug.

**1. Catch specific exceptions:** Always catch the most specific exception possible. Catching generic exceptions can hide bugs:

```python
# Bad - too broad
try:
    value = int(user_input)
    result = 100 / value
except:  # Catches everything!
    print("Error")

# Good - specific
try:
    value = int(user_input)
    result = 100 / value
except ValueError:
    print("Please enter a valid number")
except ZeroDivisionError:
    print("Cannot divide by zero")
```

**2. Don't use exceptions for flow control:** Exceptions should handle exceptional cases, not normal program flow:

```python
# Bad - using exceptions for control flow
try:
    user = users[user_id]
except KeyError:
    user = create_new_user(user_id)

# Good - use conditions
if user_id in users:
    user = users[user_id]
else:
    user = create_new_user(user_id)
```

**3. Keep try blocks small:** Only put code that might raise an exception in the try block:

```python
# Bad - try block too large
try:
    data = read_file("data.txt")
    processed = clean_data(data)
    result = analyze_data(processed)
    save_results(result)
except Exception as e:
    print(f"Error: {e}")  # Which operation failed?

# Good - targeted exception handling
data = read_file("data.txt")
processed = clean_data(data)

try:
    result = analyze_data(processed)
except ValueError as e:
    print(f"Analysis failed: {e}")
else:
    save_results(result)
```

**4. Log exceptions properly:** Include enough information for debugging:

```python
import logging

try:
    process_data(filename)
except Exception as e:
    logging.error(f"Failed to process {filename}: {e}", exc_info=True)
    # exc_info=True includes the full traceback
```

**5. Clean up resources:** Use context managers or finally blocks to ensure cleanup:

```python
# Good - context manager handles cleanup automatically
with open("data.txt", "r") as f:
    data = f.read()

# Also good - finally ensures cleanup
file = None
try:
    file = open("data.txt", "r")
    data = file.read()
finally:
    if file:
        file.close()
```

**6. Document exceptions:** Tell users what exceptions your functions might raise:

```python
def divide(a, b):
    """
    Divide two numbers.
    
    Args:
        a: The dividend (number)
        b: The divisor (number)
    
    Returns:
        The quotient (number)
    
    Raises:
        ZeroDivisionError: If b is zero
        TypeError: If a or b are not numbers
    """
    return a / b
```

## 13.11 Assertions with Assert

Assertions are a debugging aid that test conditions in your code. Think of them as sanity checks—they verify that something that should always be true actually is true. If the condition is false, an `AssertionError` is raised:

```python
def calculate_average(numbers):
    assert len(numbers) > 0, "Cannot calculate average of empty list"
    return sum(numbers) / len(numbers)

print(calculate_average([1, 2, 3]))  # Works: 2.0
print(calculate_average([]))         # AssertionError: Cannot calculate average of empty list
```

Assertions are typically used during development to catch programming errors:

```python
def set_discount(price, discount_percent):
    assert 0 <= discount_percent <= 100, "Discount must be between 0 and 100"
    assert price > 0, "Price must be positive"
    
    discount_amount = price * (discount_percent / 100)
    return price - discount_amount

final_price = set_discount(100, 20)  # $80
# set_discount(100, 150)  # AssertionError: Discount must be between 0 and 100
```

**Important difference:** Assertions vs exceptions:
- **Assertions:** For catching programmer errors and bugs during development
- **Exceptions:** For handling expected runtime errors (user input, file operations, etc.)

```python
# Use exceptions for runtime errors
def withdraw_money(balance, amount):
    if amount > balance:
        raise ValueError("Insufficient funds")  # User error - handle with exception
    return balance - amount

# Use assertions for programmer errors
def calculate_percentage(part, total):
    assert total > 0, "Total must be positive - programmer error!"
    return (part / total) * 100
```

You can use assertions to verify function contracts:

```python
def binary_search(sorted_list, target):
    # Precondition: list must be sorted
    assert sorted_list == sorted(sorted_list), "List must be sorted!"
    
    left, right = 0, len(sorted_list) - 1
    
    while left <= right:
        mid = (left + right) // 2
        if sorted_list[mid] == target:
            return mid
        elif sorted_list[mid] < target:
            left = mid + 1
        else:
            right = mid - 1
    
    return -1

# Works fine
print(binary_search([1, 3, 5, 7, 9], 5))  # 2

# Assertion fails
print(binary_search([5, 1, 9, 3, 7], 5))  # AssertionError: List must be sorted!
```

**Critical warning:** Assertions can be disabled globally with the `-O` (optimize) flag when running Python:

```bash
python -O my_program.py  # All assertions are ignored!
```

This means:
- Never use assertions for data validation or security checks
- Never use assertions for anything that must happen in production
- Only use assertions for debugging and catching developer mistakes

```python
# Bad - don't use assertions for validation
def login(username, password):
    assert username and password, "Credentials required"  # Could be disabled!

# Good - use exceptions for validation
def login(username, password):
    if not username or not password:
        raise ValueError("Credentials required")  # Always checked
```

Assertions are a powerful tool when used correctly—they make your code self-documenting and help catch bugs early!

## Common Pitfalls

**1. Catching too broad:** Don't catch `Exception` or use bare `except` unless absolutely necessary. You might hide bugs!

**2. Silently ignoring exceptions:** Always log or handle exceptions meaningfully. Empty except blocks are evil:
```python
# Bad!
try:
    risky_operation()
except:
    pass  # What went wrong? We'll never know!
```

**3. Overusing exceptions:** Don't use exceptions for normal control flow. They're slower and harder to read than regular if statements.

**4. Forgetting to clean up resources:** Always close files, connections, and other resources, even if an exception occurs.

**5. Using assertions for runtime validation:** Remember, assertions can be disabled! Use exceptions for anything that must be checked in production.

## Best Practices Summary

- Catch specific exceptions, not generic ones
- Keep try blocks focused and small
- Use finally or context managers for cleanup
- Document what exceptions your functions raise
- Use custom exceptions for application-specific errors
- Chain exceptions to preserve error context
- Log exceptions with full tracebacks
- Use assertions only for debugging, not validation
- Don't use exceptions for normal program flow

## Practice Exercises

Ready to test your skills? Try these exercises:

1. Write a function that safely converts user input to an integer, returning a default value if conversion fails.

2. Create a custom exception hierarchy for a simple banking system (e.g., `BankingError`, `InsufficientFundsError`, `InvalidAccountError`).

3. Write a function that reads a JSON file with proper exception handling for file not found, invalid JSON, and permission errors.

4. Implement a retry mechanism that attempts an operation up to 3 times before giving up, with different handling for different exception types.

## Wrapping Up

Exception handling is what separates robust, production-ready code from fragile programs that crash at the first sign of trouble. Master these concepts and you'll write Python code that gracefully handles errors, provides helpful feedback to users, and makes debugging a breeze.

Remember: exceptions aren't your enemy—they're your friends! They help you identify and handle problems systematically. Embrace them, handle them properly, and your code will be rock-solid.

Now go forth and handle those exceptions like a pro!

---

## Quiz Questions

This section contains quiz questions that test understanding of Python exception handling. Each question is formatted for easy extraction and database import.

### Question 1: Basic Exception Handling - Try-Except
**Type**: code-output
**Difficulty**: 2
**Tags**: python, exception-handling, try-except

**Question Text**:
What will be printed when this code is executed?

```python
try:
    result = 10 / 2
    print("Success!")
except ZeroDivisionError:
    print("Cannot divide by zero!")
```

**Answer Choices**:
- A) Success!
- B) Cannot divide by zero!
- C) Nothing will be printed
- D) An error will occur

**Correct Answer**: A

**Explanation**:
Since we're dividing 10 by 2 (not by zero), no exception is raised. The try block executes successfully, printing "Success!". The except block only runs if a ZeroDivisionError occurs, which doesn't happen here.

---

### Question 2: Exception Types - ValueError
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: python, exception-handling, valueerror, built-in-exceptions

**Question Text**:
Which of the following will raise a `ValueError`?

**Answer Choices**:
- A) `int(42)`
- B) `int("42")`
- C) `int("hello")`
- D) `int(3.14)`

**Correct Answer**: C

**Explanation**:
`int("hello")` raises a ValueError because "hello" cannot be converted to an integer—it's the right type (string) but has an invalid value for conversion. Option A converts an integer to an integer (works fine), option B converts a valid numeric string (works fine), and option D converts a float to an integer (also works, returning 3).

---

### Question 3: Multiple Except Clauses - Order Matters
**Type**: code-output
**Difficulty**: 4
**Tags**: python, exception-handling, multiple-except

**Question Text**:
What will be printed by this code?

```python
try:
    numbers = [1, 2, 3]
    print(numbers[10])
except IndexError:
    print("Index error!")
except ValueError:
    print("Value error!")
except Exception:
    print("Some other error!")
```

**Answer Choices**:
- A) Index error!
- B) Value error!
- C) Some other error!
- D) Nothing will be printed

**Correct Answer**: A

**Explanation**:
Accessing `numbers[10]` on a list with only 3 elements raises an IndexError. Python checks the except clauses in order and executes the first matching one. Since IndexError is caught first and matches, it prints "Index error!" and doesn't check the remaining except clauses.

---

### Question 4: Finally Clause - Execution Guarantee
**Type**: true-false
**Difficulty**: 3
**Tags**: python, exception-handling, finally

**Question Text**:
The code in a `finally` block will always execute, regardless of whether an exception occurred or not.

**Answer Choices**:
- A) True
- B) False

**Correct Answer**: A

**Explanation**:
True! The finally block always executes, whether an exception was raised, caught, or not. It even executes if there's a return statement in the try or except blocks. This makes finally perfect for cleanup operations like closing files or releasing resources. The only ways to prevent finally from executing are killing the program externally or encountering a system-level failure.

---

### Question 5: Raising Exceptions - Custom Messages
**Type**: code-output
**Difficulty**: 4
**Tags**: python, exception-handling, raise

**Question Text**:
What will be printed when this function is called with `validate_age(-5)`?

```python
def validate_age(age):
    try:
        if age < 0:
            raise ValueError("Age cannot be negative")
        print(f"Valid age: {age}")
    except ValueError as e:
        print(f"Error: {e}")

validate_age(-5)
```

**Answer Choices**:
- A) Valid age: -5
- B) Error: Age cannot be negative
- C) ValueError: Age cannot be negative
- D) Nothing will be printed

**Correct Answer**: B

**Explanation**:
When age is -5, the condition `age < 0` is true, so a ValueError is raised with the message "Age cannot be negative". This exception is caught by the except block, which prints "Error: " followed by the exception message (e), resulting in "Error: Age cannot be negative".

---

### Question 6: Custom Exceptions - Class Inheritance
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: python, exception-handling, custom-exceptions

**Question Text**:
When creating a custom exception class, which built-in class should it inherit from?

**Answer Choices**:
- A) BaseException
- B) Error
- C) Exception
- D) RuntimeError

**Correct Answer**: C

**Explanation**:
Custom exceptions should inherit from the `Exception` class (or one of its subclasses). While technically you could inherit from BaseException, this is not recommended—BaseException is the base for system-exiting exceptions like KeyboardInterrupt and SystemExit that you typically don't want to catch. Exception is the proper base for all user-defined and most built-in exceptions. Error and RuntimeError are not the standard base classes for custom exceptions.

---

### Question 7: Else Clause - Execution Condition
**Type**: code-output
**Difficulty**: 5
**Tags**: python, exception-handling, else-clause

**Question Text**:
When will "Success!" be printed in this code?

```python
try:
    value = int(input_string)
except ValueError:
    print("Invalid input!")
else:
    print("Success!")
```

**Answer Choices**:
- A) Always
- B) Only if a ValueError is raised
- C) Only if no exception is raised in the try block
- D) Never

**Correct Answer**: C

**Explanation**:
The else block in a try-except statement only executes if the try block completes successfully without raising any exceptions. If `int(input_string)` succeeds (input_string is a valid number), no exception is raised, and "Success!" is printed. If it raises a ValueError, the except block runs instead, and the else block is skipped.

---

### Question 8: Assertions - Use Cases
**Type**: multiple-choice
**Difficulty**: 6
**Tags**: python, exception-handling, assertions, assert

**Question Text**:
Which of the following is the BEST use case for assertions?

**Answer Choices**:
- A) Validating user input in a web application
- B) Checking security credentials before granting access
- C) Verifying that a function receives correctly sorted data during development
- D) Handling file not found errors

**Correct Answer**: C

**Explanation**:
Assertions are best used for debugging and catching programmer errors during development, such as verifying preconditions like "this list should be sorted." They should NOT be used for user input validation (A), security checks (B), or runtime error handling (D) because assertions can be globally disabled with Python's -O flag. For production code, use proper exception handling with if statements and raise for these cases.

---

### Question 9: Exception Chaining - Explicit Chaining
**Type**: code-output
**Difficulty**: 6
**Tags**: python, exception-handling, exception-chaining

**Question Text**:
What is the purpose of using `raise ... from` syntax?

```python
try:
    result = int("abc")
except ValueError as e:
    raise TypeError("Conversion failed") from e
```

**Answer Choices**:
- A) It prevents the original exception from being raised
- B) It explicitly chains exceptions, showing one exception was caused by another
- C) It converts one exception type to another without preserving the original
- D) It's the same as using `raise` alone

**Correct Answer**: B

**Explanation**:
The `raise ... from` syntax explicitly chains exceptions, creating a direct causal relationship where the new exception was caused by the original. This preserves the full error context via the `__cause__` attribute, making debugging easier. Option A is incorrect (both exceptions are preserved), option C is incorrect (the original IS preserved), and option D is incorrect (raise alone creates implicit chaining via `__context__`, which is different).

---

### Question 10: Best Practices - Try Block Size
**Type**: true-false
**Difficulty**: 4
**Tags**: python, exception-handling, best-practices

**Question Text**:
It's best practice to keep try blocks as small as possible, containing only the code that might raise an exception.

**Answer Choices**:
- A) True
- B) False

**Correct Answer**: A

**Explanation**:
True! Keeping try blocks small and focused makes it clear exactly which operation might fail and why. Large try blocks make debugging harder because you can't easily tell which line caused the exception. By keeping try blocks minimal, you improve code readability and make exception handling more precise and meaningful.
