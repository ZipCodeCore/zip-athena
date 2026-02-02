# Python Modules and Packages

Hey there! Let's dive into one of Python's most powerful features: modules and packages. Think of them as your code's organizational system—like having labeled filing cabinets instead of throwing all your papers into one giant pile. They help you write cleaner, more maintainable code and let you reuse code across different projects. Plus, they give you access to thousands of pre-written tools that other developers have created!

## 11.1 What Are Modules?

A module is simply a Python file that contains code—functions, classes, variables—that you can use in other Python programs. It's like having a toolbox where each tool (module) serves a specific purpose.

Think about it this way: imagine you've written some really useful functions for working with dates. Instead of copying and pasting those functions into every project that needs them, you can put them in a module and just import that module whenever you need it. Much cleaner!

Python itself comes with tons of built-in modules (the "standard library"), and you can also create your own or install modules that other people have written. Every `.py` file is technically a module that can be imported by other Python files.

Here's a simple example of using a built-in module:

```python
import math

radius = 5
area = math.pi * radius ** 2
print(f"The area of a circle with radius {radius} is {area:.2f}")
# Output: The area of a circle with radius 5 is 78.54
```

In this example, `math` is a module that comes with Python. It contains mathematical functions and constants like `pi`. By importing it, we get access to all those goodies!

## 11.2 Importing Modules

There are several ways to import modules in Python, each with its own use case. Let's explore them all!

### Basic Import

The most straightforward way is to import the entire module:

```python
import random

# Now you can use functions from the random module
number = random.randint(1, 10)
print(f"Random number: {number}")

# Pick a random item from a list
colors = ['red', 'blue', 'green', 'yellow']
chosen_color = random.choice(colors)
print(f"Chosen color: {chosen_color}")
```

When you import a module this way, you need to use the module name as a prefix (like `random.randint()`). This makes it clear where each function comes from, which is great for readability.

### Import Specific Items

Sometimes you only need one or two things from a module. You can import them directly:

```python
from math import sqrt, pi

# Now you can use sqrt and pi directly without the math. prefix
print(sqrt(16))  # Output: 4.0
print(f"Pi is approximately {pi:.2f}")  # Output: Pi is approximately 3.14
```

This is convenient when you're using just a few items frequently. However, be careful not to import something that might conflict with names you're already using in your code!

### Import with Alias

If a module name is long or you want to avoid name conflicts, you can give it a shorter alias:

```python
import datetime as dt

# Now use 'dt' instead of 'datetime'
today = dt.date.today()
print(f"Today is {today}")  # Output: Today is 2024-01-15 (or current date)

# This is especially popular with data science libraries
import numpy as np
import pandas as pd
```

The `as` keyword lets you choose any name you like, but try to use common conventions (like `np` for numpy) so other developers can easily understand your code.

### Import Everything (Use Sparingly!)

You can import everything from a module using the wildcard `*`:

```python
from math import *

# Now all functions from math are available directly
print(sqrt(25))  # Output: 5.0
print(cos(0))    # Output: 1.0
```

**Warning:** This is generally considered bad practice! It can make your code confusing because you won't know where functions came from, and you might accidentally override something. Use this only in interactive sessions or when you really know what you're doing.

## 11.3 Creating Your Own Modules

Creating your own modules is super easy—just write some Python code in a file, and voilà, you have a module! Let's create a simple module for temperature conversions.

First, create a file called `temperature.py`:

```python
# temperature.py

def celsius_to_fahrenheit(celsius):
    """Convert Celsius to Fahrenheit."""
    return (celsius * 9/5) + 32

def fahrenheit_to_celsius(fahrenheit):
    """Convert Fahrenheit to Celsius."""
    return (fahrenheit - 32) * 5/9

def celsius_to_kelvin(celsius):
    """Convert Celsius to Kelvin."""
    return celsius + 273.15

# You can also include constants
ABSOLUTE_ZERO_CELSIUS = -273.15
WATER_BOILING_POINT_CELSIUS = 100
```

Now, in another file in the same directory, you can import and use it:

```python
# main.py
import temperature

temp_c = 25
temp_f = temperature.celsius_to_fahrenheit(temp_c)
print(f"{temp_c}°C is {temp_f}°F")

# Output: 25°C is 77.0°F

# Or import specific functions
from temperature import fahrenheit_to_celsius, WATER_BOILING_POINT_CELSIUS

temp = fahrenheit_to_celsius(98.6)
print(f"98.6°F is {temp:.1f}°C")
# Output: 98.6°F is 37.0°C
```

That's it! You've created your first reusable module. You can use this `temperature.py` file in any Python project now.

## 11.4 Module Search Path

When you import a module, Python needs to find it. But where does it look? Python searches several locations in a specific order:

1. **Current directory**: The folder where your script is running
2. **PYTHONPATH**: Environment variable (if set)
3. **Standard library directories**: Where Python's built-in modules live
4. **Site-packages**: Where third-party packages are installed

You can see the search path like this:

```python
import sys

print("Python will search for modules in these locations:")
for path in sys.path:
    print(path)
```

This will print something like:

```
/home/user/my_project
/usr/lib/python3.9
/usr/lib/python3.9/site-packages
...
```

Understanding this is helpful when you're wondering why Python can't find a module you created. Make sure your module is in one of these directories, or add your module's directory to `sys.path`:

```python
import sys
sys.path.append('/path/to/my/modules')

import my_custom_module  # Now Python can find it!
```

However, modifying `sys.path` is usually a last resort. It's better to organize your code properly using packages (which we'll cover soon) or install your module properly.

## 11.5 The __name__ Variable and if __name__ == "__main__":

Here's something that confuses beginners but is actually brilliant once you understand it. Every module has a built-in variable called `__name__`. When you run a file directly, `__name__` is set to `"__main__"`. When you import it as a module, `__name__` is set to the module's name.

Why does this matter? It lets you write code that runs when the file is executed directly but not when it's imported. Here's an example:

```python
# calculator.py

def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

def multiply(a, b):
    return a * b

def divide(a, b):
    if b == 0:
        return "Cannot divide by zero!"
    return a / b

# This code only runs when calculator.py is run directly
if __name__ == "__main__":
    print("Calculator Module Test")
    print(f"10 + 5 = {add(10, 5)}")
    print(f"10 - 5 = {subtract(10, 5)}")
    print(f"10 * 5 = {multiply(10, 5)}")
    print(f"10 / 5 = {divide(10, 5)}")
```

Now if you run `python calculator.py` directly, you'll see the test output. But if you import it in another file:

```python
# another_file.py
import calculator

result = calculator.add(100, 50)
print(result)  # Output: 150

# The test code from calculator.py does NOT run!
```

This is super useful for:
- Including test code in your modules
- Creating modules that can work as both libraries and standalone scripts
- Demonstrating how to use your module's functions

## 11.6 Standard Library Modules

Python comes with a fantastic collection of built-in modules called the standard library. You don't need to install anything—they're ready to use! Let's explore some of the most commonly used ones.

### The math Module

For mathematical operations beyond basic arithmetic:

```python
import math

# Constants
print(math.pi)      # 3.141592653589793
print(math.e)       # 2.718281828459045

# Functions
print(math.sqrt(16))        # 4.0
print(math.ceil(4.3))       # 5 (round up)
print(math.floor(4.7))      # 4 (round down)
print(math.factorial(5))    # 120
print(math.sin(math.pi/2))  # 1.0
print(math.log(100, 10))    # 2.0 (log base 10 of 100)
```

### The random Module

For generating random numbers and making random choices:

```python
import random

# Random float between 0 and 1
print(random.random())

# Random integer in a range
dice_roll = random.randint(1, 6)
print(f"You rolled a {dice_roll}")

# Random choice from a list
fruits = ['apple', 'banana', 'cherry', 'date']
print(random.choice(fruits))

# Shuffle a list (modifies in-place)
deck = ['A', 'K', 'Q', 'J', '10']
random.shuffle(deck)
print(deck)

# Random sample (without replacement)
winners = random.sample(['Alice', 'Bob', 'Charlie', 'Diana', 'Eve'], 3)
print(f"Winners: {winners}")
```

### The datetime Module

For working with dates and times:

```python
import datetime

# Current date and time
now = datetime.datetime.now()
print(f"Current date and time: {now}")

# Current date only
today = datetime.date.today()
print(f"Today's date: {today}")

# Create specific date
birthday = datetime.date(1990, 5, 15)
print(f"Birthday: {birthday}")

# Calculate age
age = today.year - birthday.year
print(f"Age: {age}")

# Format dates
formatted = now.strftime("%B %d, %Y at %I:%M %p")
print(formatted)  # e.g., "January 15, 2024 at 02:30 PM"

# Date arithmetic
tomorrow = today + datetime.timedelta(days=1)
next_week = today + datetime.timedelta(weeks=1)
print(f"Tomorrow: {tomorrow}")
print(f"Next week: {next_week}")
```

### The os Module

For interacting with the operating system:

```python
import os

# Get current working directory
print(f"Current directory: {os.getcwd()}")

# List files in a directory
files = os.listdir('.')
print(f"Files in current directory: {files}")

# Check if a file or directory exists
if os.path.exists('myfile.txt'):
    print("File exists!")

# Get file size
if os.path.exists('myfile.txt'):
    size = os.path.getsize('myfile.txt')
    print(f"File size: {size} bytes")

# Join path components (works on any OS)
path = os.path.join('folder', 'subfolder', 'file.txt')
print(path)  # folder/subfolder/file.txt (or folder\subfolder\file.txt on Windows)

# Environment variables
home = os.getenv('HOME')
print(f"Home directory: {home}")
```

### The sys Module

For system-specific parameters and functions:

```python
import sys

# Python version
print(f"Python version: {sys.version}")

# Command-line arguments
print(f"Script name: {sys.argv[0]}")
if len(sys.argv) > 1:
    print(f"Arguments: {sys.argv[1:]}")

# Exit the program with a status code
# sys.exit(0)  # Uncomment to exit with success code

# Module search path
print("Module search paths:")
for path in sys.path[:3]:  # Show first 3
    print(f"  {path}")

# Platform information
print(f"Operating system: {sys.platform}")
```

### The json Module

For working with JSON data (super common for APIs and configuration files):

```python
import json

# Python dict to JSON string
person = {
    'name': 'Alice',
    'age': 30,
    'city': 'New York',
    'hobbies': ['reading', 'cycling', 'cooking']
}

json_string = json.dumps(person, indent=2)
print(json_string)

# JSON string to Python dict
json_data = '{"name": "Bob", "age": 25, "city": "Boston"}'
parsed = json.loads(json_data)
print(f"Name: {parsed['name']}, Age: {parsed['age']}")

# Write to file
with open('person.json', 'w') as f:
    json.dump(person, f, indent=2)

# Read from file
with open('person.json', 'r') as f:
    loaded_person = json.load(f)
    print(f"Loaded: {loaded_person}")
```

### The re Module

For regular expressions (pattern matching in strings):

```python
import re

# Search for a pattern
text = "My email is alice@example.com"
match = re.search(r'\w+@\w+\.\w+', text)
if match:
    print(f"Found email: {match.group()}")  # Output: alice@example.com

# Find all matches
text = "Call me at 555-1234 or 555-5678"
phone_numbers = re.findall(r'\d{3}-\d{4}', text)
print(f"Phone numbers: {phone_numbers}")  # ['555-1234', '555-5678']

# Replace patterns
text = "Hello, World!"
new_text = re.sub(r'World', 'Python', text)
print(new_text)  # Output: Hello, Python!

# Split by pattern
text = "apple,banana;cherry:date"
fruits = re.split(r'[,;:]', text)
print(fruits)  # ['apple', 'banana', 'cherry', 'date']
```

## 11.7 Packages and __init__.py

A package is a way to organize related modules into a directory hierarchy. Think of it as a folder containing modules, and possibly sub-folders with more modules.

To create a package, you need:
1. A directory with a meaningful name
2. An `__init__.py` file inside that directory (can be empty)
3. One or more module files

Here's an example structure:

```
my_game/
    __init__.py
    player.py
    enemy.py
    weapons.py
```

Let's create this:

```python
# my_game/__init__.py
# This file can be empty, or contain initialization code

print("Game package initialized!")
__version__ = "1.0.0"

# You can also control what gets imported with *
__all__ = ['player', 'enemy', 'weapons']
```

```python
# my_game/player.py

class Player:
    def __init__(self, name, health=100):
        self.name = name
        self.health = health
    
    def take_damage(self, amount):
        self.health -= amount
        print(f"{self.name} took {amount} damage. Health: {self.health}")
```

```python
# my_game/enemy.py

class Enemy:
    def __init__(self, name, health=50):
        self.name = name
        self.health = health
    
    def attack(self, damage):
        return damage
```

Now you can import from your package:

```python
# main.py
from my_game.player import Player
from my_game.enemy import Enemy

# Or import the whole package
import my_game.player
import my_game.enemy

# Create game objects
hero = Player("Hero", 100)
goblin = Enemy("Goblin", 30)

hero.take_damage(10)
# Output: Hero took 10 damage. Health: 90
```

You can also create nested packages:

```
my_game/
    __init__.py
    player.py
    enemies/
        __init__.py
        goblin.py
        dragon.py
    items/
        __init__.py
        weapons.py
        potions.py
```

Then import like this:

```python
from my_game.enemies.goblin import Goblin
from my_game.items.weapons import Sword
```

The `__init__.py` file is crucial—it tells Python that this directory should be treated as a package. In Python 3.3+, it can be empty, but it's often used to:
- Run initialization code
- Set the `__all__` variable to control `from package import *`
- Import specific items to make them available at the package level

## 11.8 Installing Third-Party Packages with pip

Python's standard library is great, but there's a whole universe of third-party packages created by other developers. These are available on PyPI (Python Package Index), and you install them using `pip` (Python's package manager).

### Installing a Package

Open your terminal or command prompt and type:

```bash
pip install requests
```

This installs the `requests` library, which makes working with HTTP requests super easy:

```python
import requests

# Make a GET request
response = requests.get('https://api.github.com')

if response.status_code == 200:
    print("Success!")
    data = response.json()
    print(data)
```

### Common pip Commands

```bash
# Install a package
pip install package_name

# Install a specific version
pip install package_name==1.2.3

# Upgrade a package
pip install --upgrade package_name

# Uninstall a package
pip uninstall package_name

# List installed packages
pip list

# Show package information
pip show package_name

# Save installed packages to a file
pip freeze > requirements.txt

# Install packages from a file
pip install -r requirements.txt
```

### Popular Third-Party Packages

Here are some packages you'll likely encounter:

- **requests**: HTTP library for making web requests
- **numpy**: Numerical computing with arrays and matrices
- **pandas**: Data analysis and manipulation
- **matplotlib**: Creating plots and visualizations
- **flask** or **django**: Web frameworks
- **pytest**: Testing framework
- **pillow**: Image processing
- **beautifulsoup4**: Web scraping

Installing and using a simple package:

```bash
pip install colorama
```

```python
from colorama import Fore, Style, init

# Initialize colorama
init()

# Print colored text
print(Fore.RED + "This is red text!")
print(Fore.GREEN + "This is green text!")
print(Fore.BLUE + "This is blue text!")
print(Style.RESET_ALL + "Back to normal")
```

### Virtual Environments

Before installing packages, it's a good practice to use virtual environments. They create isolated Python environments for different projects, so packages from one project don't interfere with another:

```bash
# Create a virtual environment
python -m venv myenv

# Activate it (Windows)
myenv\Scripts\activate

# Activate it (Mac/Linux)
source myenv/bin/activate

# Now install packages
pip install requests

# Deactivate when done
deactivate
```

When you activate a virtual environment, any packages you install go into that environment only, not your global Python installation. This keeps things clean and organized!

## Common Pitfalls

**Circular Imports**: Be careful not to create modules that import each other. If `module_a.py` imports `module_b.py` and vice versa, you'll get an error.

**Name Conflicts**: Don't name your own files the same as standard library modules. If you create `random.py`, you won't be able to import Python's built-in `random` module!

**Relative vs Absolute Imports**: In packages, use absolute imports (`from my_package.module import something`) rather than relative imports when possible—they're clearer and less error-prone.

**Forgetting __init__.py**: In older Python versions (before 3.3), packages absolutely must have an `__init__.py` file. Keep this habit even in newer versions for compatibility.

## Best Practices

**One Module, One Purpose**: Each module should have a clear, focused purpose. Don't create a giant module that does everything.

**Use Descriptive Names**: Module names should be lowercase with underscores if needed (`my_module.py`, not `MyModule.py` or `mymodule.py`).

**Document Your Modules**: Add docstrings at the top of each module explaining what it does:

```python
"""
This module provides temperature conversion utilities.

Functions:
    celsius_to_fahrenheit: Convert Celsius to Fahrenheit
    fahrenheit_to_celsius: Convert Fahrenheit to Celsius
"""
```

**Keep Imports at the Top**: Put all your import statements at the beginning of your file (except in rare special cases).

**Avoid `import *`**: Be explicit about what you're importing. It makes code more readable and prevents name conflicts.

## Practice Exercises

1. Create a module called `string_utils.py` with functions to reverse a string, count vowels, and check if a string is a palindrome. Then import and test it.

2. Write a program that uses the `datetime` module to calculate how many days until your next birthday.

3. Create a simple package called `shapes` with modules for `circle.py`, `rectangle.py`, and `triangle.py`. Each should have functions to calculate area and perimeter.

4. Use the `json` module to save a dictionary of your favorite movies to a file, then load it back and print it.

5. Explore the `random` module by creating a simple dice rolling simulator that rolls two dice and displays the results.

6. Install a third-party package like `colorama` or `emoji` and create a fun program that uses it.

## Wrapping Up

Modules and packages are essential to writing organized, reusable Python code. They let you break your programs into logical pieces, use code written by others, and share your own code easily. The standard library gives you tons of functionality right out of the box, and pip opens up a world of third-party packages.

Start by creating your own simple modules, practice importing in different ways, and explore the standard library. Before you know it, you'll be creating complex, well-organized Python projects like a pro!

---

## Quiz Questions

This section contains quiz questions that test understanding of Python modules and packages. Each question is formatted for easy extraction and database import.

### Question 1: Module Basics - Import Statement
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: python, modules, import, basics

**Question Text**:
What is the correct way to import the `math` module and use its `sqrt` function?

**Answer Choices**:
- A) `include math` then `sqrt(16)`
- B) `import math` then `math.sqrt(16)`
- C) `using math` then `sqrt(16)`
- D) `require math` then `math.sqrt(16)`

**Correct Answer**: B

**Explanation**:
In Python, you use the `import` keyword to import modules. After importing `math`, you access its functions using the dot notation: `math.sqrt(16)`. Options A, C, and D use keywords from other programming languages (C's include, C#'s using, Ruby's require) that don't exist in Python.

---

### Question 2: Import Variations - from...import
**Type**: code-output
**Difficulty**: 3
**Tags**: python, modules, import, from-import

**Question Text**:
What will this code print?

```python
from math import pi, sqrt

print(sqrt(9))
print(pi)
```

**Answer Choices**:
- A) Error: sqrt and pi are not defined
- B) 3.0 and 3.141592653589793
- C) math.sqrt(9) and math.pi
- D) 9 and 3.14

**Correct Answer**: B

**Explanation**:
The `from math import pi, sqrt` statement imports `pi` and `sqrt` directly into the current namespace, so you can use them without the `math.` prefix. `sqrt(9)` returns `3.0` (as a float), and `pi` is approximately 3.141592653589793. No error occurs because the names are properly imported.

---

### Question 3: Module Search Path - Import Location
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: python, modules, import, sys-path

**Question Text**:
When you run `import mymodule`, in what order does Python search for `mymodule.py`?

**Answer Choices**:
- A) Standard library → current directory → site-packages
- B) Current directory → PYTHONPATH → standard library → site-packages
- C) Site-packages → standard library → current directory
- D) Current directory only

**Correct Answer**: B

**Explanation**:
Python searches for modules in this order: 1) The current directory (where the script is running), 2) Directories in the PYTHONPATH environment variable (if set), 3) Standard library directories, and 4) Site-packages (where third-party packages are installed). This search path is stored in `sys.path` and can be viewed by importing sys and printing `sys.path`.

---

### Question 4: __name__ Variable - Script vs Module
**Type**: code-output
**Difficulty**: 5
**Tags**: python, modules, __name__, __main__

**Question Text**:
You have a file `calculator.py` with this code:

```python
def add(a, b):
    return a + b

if __name__ == "__main__":
    print("Testing:", add(5, 3))
```

What happens when you run `import calculator` in another file?

**Answer Choices**:
- A) It prints "Testing: 8"
- B) It prints nothing (the test code doesn't run)
- C) It causes an error
- D) It prints "Testing:" without the number

**Correct Answer**: B

**Explanation**:
When you import a module, Python sets `__name__` to the module's name (in this case, `"calculator"`), not `"__main__"`. Therefore, the condition `if __name__ == "__main__":` is False, and the code inside that block doesn't execute. The test code only runs when you execute `calculator.py` directly (e.g., `python calculator.py`), not when you import it.

---

### Question 5: Standard Library - datetime Module
**Type**: code-output
**Difficulty**: 4
**Tags**: python, modules, datetime, standard-library

**Question Text**:
What will be the output type and approximate value when running this code?

```python
import datetime

today = datetime.date.today()
birthday = datetime.date(2020, 1, 15)
difference = today - birthday
print(type(difference))
```

**Answer Choices**:
- A) `<class 'int'>` representing days
- B) `<class 'datetime.timedelta'>` representing the time difference
- C) `<class 'str'>` with a formatted date
- D) `<class 'float'>` representing years

**Correct Answer**: B

**Explanation**:
When you subtract two `datetime.date` objects, Python returns a `datetime.timedelta` object, which represents the difference between two dates or times. The timedelta object contains information about days, seconds, and microseconds. To get the number of days as an integer, you would access `difference.days`. It's not automatically converted to int, str, or float.

---

### Question 6: Packages - __init__.py Purpose
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: python, packages, __init__, package-structure

**Question Text**:
What is the primary purpose of the `__init__.py` file in a Python package?

**Answer Choices**:
- A) It contains the main executable code for the package
- B) It marks the directory as a Python package and can contain initialization code
- C) It is required to store all function definitions for the package
- D) It prevents other files from accessing the package

**Correct Answer**: B

**Explanation**:
The `__init__.py` file serves to tell Python that a directory should be treated as a package (collection of modules). While it can be empty, it often contains initialization code that runs when the package is imported, defines the `__all__` variable to control what's exported, or imports items to make them available at the package level. It's not required to contain all function definitions—those go in separate module files. It doesn't contain the main executable code or prevent access.

---

### Question 7: Third-Party Packages - pip Commands
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: python, packages, pip, package-management

**Question Text**:
Which pip command would you use to save a list of all currently installed packages to a file called `requirements.txt`?

**Answer Choices**:
- A) `pip list > requirements.txt`
- B) `pip freeze > requirements.txt`
- C) `pip save requirements.txt`
- D) `pip export requirements.txt`

**Correct Answer**: B

**Explanation**:
The `pip freeze` command outputs installed packages in a format that can be used with `pip install -r requirements.txt` to recreate the environment. By redirecting output with `>`, you save it to a file. While `pip list > requirements.txt` would create a file, the format is less suitable for reinstallation. Options C and D use commands that don't exist in pip. The freeze format shows packages as `package==version`, making it perfect for version-controlled dependency lists.

---

### Question 8: Import Aliases - Common Conventions
**Type**: true-false
**Difficulty**: 2
**Tags**: python, modules, import, aliases

**Question Text**:
True or False: The following import statement is valid Python syntax:

```python
import datetime as dt
```

**Correct Answer**: True

**Explanation**:
This is absolutely valid syntax. The `as` keyword allows you to create an alias (nickname) for a module. After this import, you can use `dt` instead of typing out `datetime` every time. For example, `dt.date.today()` instead of `datetime.date.today()`. This is a common practice, especially with libraries that have long names. Common conventions include `import numpy as np`, `import pandas as pd`, and `import matplotlib.pyplot as plt`.

---

### Question 9: Import Variations - Wildcard Import
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: python, modules, import, wildcard, best-practices

**Question Text**:
Why is using `from module import *` generally discouraged?

**Answer Choices**:
- A) It makes the code run slower
- B) It imports too much data and uses too much memory
- C) It can cause name conflicts and makes code harder to understand
- D) It only works with standard library modules

**Correct Answer**: C

**Explanation**:
Using `from module import *` imports all public names from a module into your current namespace. This creates two problems: 1) Name conflicts—you might accidentally override existing variables or functions without realizing it, and 2) Code readability—when you see a function call, you don't know which module it came from. Performance and memory are not significantly affected, and it works with any module, not just standard library ones. Explicit imports (`from module import specific_name`) are preferred.

---

### Question 10: Creating Modules - Module Execution
**Type**: code-output
**Difficulty**: 6
**Tags**: python, modules, module-creation, __name__

**Question Text**:
Given two files with this content:

```python
# utils.py
print("Utils loaded")
value = 42

if __name__ == "__main__":
    print("Running utils directly")
    print(f"Value is {value}")
```

```python
# main.py
print("Main starting")
import utils
print(f"Imported value: {utils.value}")
```

What output will you see when you run `python main.py`?

**Answer Choices**:
- A) Main starting / Utils loaded / Running utils directly / Value is 42 / Imported value: 42
- B) Main starting / Utils loaded / Imported value: 42
- C) Utils loaded / Main starting / Imported value: 42
- D) Main starting / Imported value: 42

**Correct Answer**: B

**Explanation**:
When you run `main.py`, Python executes it from top to bottom. First, it prints "Main starting". Then it encounters `import utils`, which loads the utils module. When a module is imported, Python executes all code at the module level, so "Utils loaded" is printed and `value = 42` is set. However, the `if __name__ == "__main__":` block in utils.py does NOT execute because when imported, `__name__` is set to "utils", not "__main__". Finally, main.py prints "Imported value: 42". The test code in utils only runs if you execute utils.py directly.

---

### Question 11: Standard Library - json Module
**Type**: code-output
**Difficulty**: 5
**Tags**: python, modules, json, standard-library, data-serialization

**Question Text**:
What will be the type of `result` after this code executes?

```python
import json

data = {"name": "Alice", "age": 30}
json_string = json.dumps(data)
result = json.loads(json_string)
print(type(result))
```

**Answer Choices**:
- A) `<class 'str'>`
- B) `<class 'dict'>`
- C) `<class 'json'>`
- D) `<class 'list'>`

**Correct Answer**: B

**Explanation**:
`json.dumps()` converts a Python dictionary to a JSON string (serialization). Then `json.loads()` converts a JSON string back to a Python dictionary (deserialization). So `result` is a Python dictionary (`dict`), not a string, json object (no such type), or list. The 's' in `dumps` and `loads` stands for "string"—these functions work with JSON strings. For files, you'd use `json.dump()` and `json.load()` (without the 's').

---

### Question 12: Module Organization - Best Practices
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: python, modules, best-practices, code-organization

**Question Text**:
Which of the following is considered a best practice when creating Python modules?

**Answer Choices**:
- A) Name modules using CamelCase like MyModule.py
- B) Put all import statements at the bottom of the file
- C) Use lowercase names with underscores like my_module.py
- D) Include all your functions in one large module for convenience

**Correct Answer**: C

**Explanation**:
Python module naming convention is to use lowercase letters with underscores for readability (snake_case), like `my_module.py` or `string_utils.py`. CamelCase is used for class names, not module names. Import statements should go at the top of the file (after the module docstring) for visibility and clarity. Modules should be focused and serve a single purpose—putting everything in one large module defeats the purpose of modularization and makes code harder to maintain. Following PEP 8 style guidelines helps make your code more readable and maintainable.

---
