# Python Data Structures

Hey! Let's explore one of Python's superpowers: its built-in data structures. These are your tools for organizing and managing collections of data. Python makes working with data structures incredibly easy and intuitive.

## What Are Data Structures?

Data structures are ways to organize and store multiple pieces of information. Instead of having 100 separate variables for 100 students' names, you can use a single data structure to hold them all. Python gives you four main built-in data structures, each with its own personality and use cases.

## Lists: Your Flexible, Ordered Collection

Lists are probably the most versatile data structure in Python. They're ordered, changeable, and can hold any type of data.

### Creating Lists

```python
# Empty list
my_list = []

# List with values
fruits = ["apple", "banana", "cherry"]
numbers = [1, 2, 3, 4, 5]
mixed = [1, "hello", 3.14, True]  # Can mix types!
```

### Accessing Elements

Lists are zero-indexed, meaning the first element is at position 0:

```python
fruits = ["apple", "banana", "cherry"]
print(fruits[0])   # "apple"
print(fruits[1])   # "banana"
print(fruits[-1])  # "cherry" (negative indices count from the end)
print(fruits[-2])  # "banana"
```

### Slicing Lists

Want to grab a portion of a list? Use slicing:

```python
numbers = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]

print(numbers[2:5])    # [2, 3, 4] (start at 2, stop before 5)
print(numbers[:3])     # [0, 1, 2] (from beginning to position 3)
print(numbers[7:])     # [7, 8, 9] (from position 7 to end)
print(numbers[::2])    # [0, 2, 4, 6, 8] (every 2nd element)
print(numbers[::-1])   # [9, 8, 7, 6, 5, 4, 3, 2, 1, 0] (reversed!)
```

### Modifying Lists

Lists are mutable (changeable):

```python
fruits = ["apple", "banana", "cherry"]

# Change an element
fruits[1] = "blueberry"
print(fruits)  # ["apple", "blueberry", "cherry"]

# Add elements
fruits.append("date")           # Add to end
fruits.insert(1, "avocado")     # Insert at specific position
fruits.extend(["fig", "grape"]) # Add multiple items

# Remove elements
fruits.remove("apple")  # Remove by value
popped = fruits.pop()   # Remove and return last item
del fruits[0]           # Delete by index
fruits.clear()          # Remove all items
```

### List Methods

```python
numbers = [3, 1, 4, 1, 5, 9, 2, 6]

numbers.sort()           # Sort in place: [1, 1, 2, 3, 4, 5, 6, 9]
numbers.reverse()        # Reverse in place: [9, 6, 5, 4, 3, 2, 1, 1]

count = numbers.count(1) # Count occurrences: 2
index = numbers.index(5) # Find first index of 5

# Non-modifying methods (they return values)
length = len(numbers)    # Get length
```

### List Comprehensions

This is a powerful, Pythonic way to create lists:

```python
# Traditional way
squares = []
for i in range(10):
    squares.append(i ** 2)

# List comprehension way (much cleaner!)
squares = [i ** 2 for i in range(10)]
# [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]

# With a condition
evens = [i for i in range(20) if i % 2 == 0]
# [0, 2, 4, 6, 8, 10, 12, 14, 16, 18]

# Transform strings
names = ["alice", "bob", "charlie"]
upper_names = [name.upper() for name in names]
# ["ALICE", "BOB", "CHARLIE"]
```

### Nested Lists (2D Lists)

Lists can contain other lists:

```python
matrix = [
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9]
]

print(matrix[0])      # [1, 2, 3] (first row)
print(matrix[1][2])   # 6 (row 1, column 2)

# Iterate through 2D list
for row in matrix:
    for item in row:
        print(item, end=" ")
    print()  # New line after each row
```

## Tuples: Immutable Sequences

Tuples are like lists, but they can't be changed after creation. They're faster and protect data from accidental modification.

### Creating Tuples

```python
# Basic tuple
point = (3, 4)
colors = ("red", "green", "blue")

# Single element tuple (notice the comma!)
single = (42,)  # Without comma, it's just a number in parentheses

# Tuple without parentheses (works too!)
coordinates = 10, 20
```

### Working with Tuples

```python
point = (3, 4, 5)

# Accessing elements (same as lists)
x = point[0]  # 3
y = point[1]  # 4

# Tuple unpacking (super useful!)
x, y, z = point
print(x, y, z)  # 3 4 5

# Tuple methods (only 2!)
numbers = (1, 2, 2, 3, 2, 4)
count = numbers.count(2)  # 3
index = numbers.index(3)  # 3
```

### When to Use Tuples

```python
# Good use cases for tuples:

# 1. Returning multiple values from a function
def get_coordinates():
    return (10, 20)

x, y = get_coordinates()

# 2. Dictionary keys (lists can't be keys!)
locations = {
    (0, 0): "origin",
    (1, 0): "east",
    (0, 1): "north"
}

# 3. Protecting data from changes
DAYS_OF_WEEK = ("Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun")
```

## Dictionaries: Key-Value Pairs

Dictionaries store data as key-value pairs. Think of them like a real dictionary: you look up a word (key) to find its definition (value).

### Creating Dictionaries

```python
# Empty dictionary
my_dict = {}

# Dictionary with values
student = {
    "name": "Alex",
    "age": 20,
    "grade": "A",
    "is_enrolled": True
}

# Using dict() constructor
person = dict(name="Jordan", age=25)
```

### Accessing Values

```python
student = {"name": "Alex", "age": 20, "grade": "A"}

# Access with bracket notation
name = student["name"]  # "Alex"

# Access with get() (safer - doesn't crash if key missing)
age = student.get("age")           # 20
gpa = student.get("gpa", 0.0)      # 0.0 (default value)
```

### Modifying Dictionaries

```python
student = {"name": "Alex", "age": 20}

# Add or update
student["age"] = 21              # Update existing
student["major"] = "CS"          # Add new

# Update multiple at once
student.update({"grade": "A", "credits": 60})

# Remove items
grade = student.pop("grade")     # Remove and return value
del student["age"]               # Delete a key
student.clear()                  # Remove all items
```

### Dictionary Methods

```python
student = {"name": "Alex", "age": 20, "grade": "A"}

# Get all keys, values, or items
keys = student.keys()      # dict_keys(['name', 'age', 'grade'])
values = student.values()  # dict_values(['Alex', 20, 'A'])
items = student.items()    # dict_items([('name', 'Alex'), ('age', 20), ('grade', 'A')])

# Check if key exists
if "name" in student:
    print("Name found!")

# Get length
count = len(student)  # 3
```

### Iterating Over Dictionaries

```python
student = {"name": "Alex", "age": 20, "grade": "A"}

# Iterate over keys (default)
for key in student:
    print(key, student[key])

# More explicit: iterate over keys
for key in student.keys():
    print(key)

# Iterate over values
for value in student.values():
    print(value)

# Iterate over key-value pairs (most common!)
for key, value in student.items():
    print(f"{key}: {value}")
```

### Dictionary Comprehensions

```python
# Square numbers
squares = {x: x**2 for x in range(6)}
# {0: 0, 1: 1, 2: 4, 3: 9, 4: 16, 5: 25}

# Flip keys and values
original = {"a": 1, "b": 2, "c": 3}
flipped = {value: key for key, value in original.items()}
# {1: 'a', 2: 'b', 3: 'c'}

# With condition
scores = {"Alice": 85, "Bob": 92, "Charlie": 78, "Diana": 95}
high_scorers = {name: score for name, score in scores.items() if score >= 90}
# {"Bob": 92, "Diana": 95}
```

### Nested Dictionaries

```python
students = {
    "student1": {
        "name": "Alice",
        "age": 20,
        "grades": [85, 90, 88]
    },
    "student2": {
        "name": "Bob",
        "age": 21,
        "grades": [92, 88, 95]
    }
}

# Access nested data
alice_age = students["student1"]["age"]  # 20
bob_grades = students["student2"]["grades"]  # [92, 88, 95]
```

## Sets: Unique, Unordered Collections

Sets are collections of unique items with no duplicates. They're perfect for membership testing and removing duplicates.

### Creating Sets

```python
# Empty set (can't use {} - that's a dict!)
empty_set = set()

# Set with values
fruits = {"apple", "banana", "cherry"}
numbers = {1, 2, 3, 4, 5}

# From a list (removes duplicates!)
numbers = [1, 2, 2, 3, 3, 3, 4]
unique = set(numbers)  # {1, 2, 3, 4}
```

### Set Operations

```python
fruits = {"apple", "banana", "cherry"}

# Add elements
fruits.add("date")
fruits.update(["fig", "grape"])  # Add multiple

# Remove elements
fruits.remove("apple")    # Raises error if not found
fruits.discard("banana")  # Doesn't raise error if not found
fruits.pop()              # Remove and return arbitrary item
fruits.clear()            # Remove all

# Membership testing (super fast!)
if "cherry" in fruits:
    print("Found!")
```

### Set Mathematics

Sets really shine with mathematical operations:

```python
set1 = {1, 2, 3, 4, 5}
set2 = {4, 5, 6, 7, 8}

# Union (all items from both sets)
union = set1 | set2
# or: union = set1.union(set2)
# {1, 2, 3, 4, 5, 6, 7, 8}

# Intersection (items in both sets)
intersection = set1 & set2
# or: intersection = set1.intersection(set2)
# {4, 5}

# Difference (items in set1 but not set2)
difference = set1 - set2
# or: difference = set1.difference(set2)
# {1, 2, 3}

# Symmetric difference (items in one set but not both)
sym_diff = set1 ^ set2
# or: sym_diff = set1.symmetric_difference(set2)
# {1, 2, 3, 6, 7, 8}
```

### Set Comparisons

```python
set1 = {1, 2, 3}
set2 = {1, 2, 3, 4, 5}

# Subset (all items of set1 are in set2)
is_subset = set1 <= set2  # True
# or: is_subset = set1.issubset(set2)

# Superset (set2 contains all items of set1)
is_superset = set2 >= set1  # True
# or: is_superset = set2.issuperset(set1)
```

### When to Use Sets

```python
# Remove duplicates
numbers = [1, 2, 2, 3, 3, 3, 4, 4, 5]
unique_numbers = list(set(numbers))  # [1, 2, 3, 4, 5]

# Fast membership testing
valid_users = {"alice", "bob", "charlie"}
user = "alice"
if user in valid_users:  # O(1) - constant time!
    print("Valid user")

# Find common elements
python_students = {"Alice", "Bob", "Charlie", "Diana"}
java_students = {"Bob", "Charlie", "Eve", "Frank"}
both = python_students & java_students  # {"Bob", "Charlie"}
```

## Choosing the Right Data Structure

Here's a quick guide:

### Use a **List** when:
- You need an ordered collection
- Items might have duplicates
- You need to access items by index
- You'll modify the collection frequently

### Use a **Tuple** when:
- You need an immutable (unchangeable) sequence
- You're returning multiple values from a function
- You need to use the data as a dictionary key
- You want to protect data from accidental changes

### Use a **Dictionary** when:
- You need to look up values by a key
- You have associated pairs of data (like name → phone number)
- You need fast lookups
- Your data has a natural key-value structure

### Use a **Set** when:
- You only care about unique items
- Order doesn't matter
- You need fast membership testing
- You want to perform mathematical set operations

## Practical Examples

### Example 1: Student Gradebook

```python
# Dictionary of students with lists of grades
gradebook = {
    "Alice": [85, 90, 88, 92],
    "Bob": [78, 82, 80, 85],
    "Charlie": [92, 95, 90, 88]
}

# Calculate average for each student
for student, grades in gradebook.items():
    average = sum(grades) / len(grades)
    print(f"{student}: {average:.2f}")
```

### Example 2: Remove Duplicates from Survey

```python
survey_responses = ["Python", "Java", "Python", "JavaScript", 
                    "Python", "Java", "C++", "Python"]

# Get unique responses and count
unique_languages = set(survey_responses)
print(f"Languages mentioned: {unique_languages}")

# Count occurrences
for language in unique_languages:
    count = survey_responses.count(language)
    print(f"{language}: {count} votes")
```

### Example 3: Inventory System

```python
# Using different data structures together
inventory = [
    {"item": "Apple", "quantity": 50, "price": 0.50},
    {"item": "Banana", "quantity": 30, "price": 0.30},
    {"item": "Cherry", "quantity": 20, "price": 2.00}
]

# Find total value of inventory
total_value = sum(item["quantity"] * item["price"] for item in inventory)
print(f"Total inventory value: ${total_value:.2f}")

# Get all item names
item_names = [item["item"] for item in inventory]
print(f"Items in stock: {', '.join(item_names)}")
```

## Common Pitfalls

### 1. Modifying a List While Iterating

```python
# WRONG - can cause unexpected behavior
numbers = [1, 2, 3, 4, 5]
for num in numbers:
    if num % 2 == 0:
        numbers.remove(num)  # Don't do this!

# RIGHT - create a new list
numbers = [1, 2, 3, 4, 5]
numbers = [num for num in numbers if num % 2 != 0]
```

### 2. Mutable Default Arguments

```python
# WRONG
def add_item(item, my_list=[]):
    my_list.append(item)
    return my_list

# RIGHT
def add_item(item, my_list=None):
    if my_list is None:
        my_list = []
    my_list.append(item)
    return my_list
```

### 3. Shallow vs Deep Copy

```python
# Shallow copy - changes affect both!
original = [[1, 2], [3, 4]]
copy = original.copy()
copy[0][0] = 999
print(original)  # [[999, 2], [3, 4]] - Oops!

# Deep copy - independent copies
import copy
original = [[1, 2], [3, 4]]
deep = copy.deepcopy(original)
deep[0][0] = 999
print(original)  # [[1, 2], [3, 4]] - Safe!
```

## Practice Exercises

1. Create a list of your favorite movies and use a list comprehension to get only those released after 2010.

2. Write a function that takes two lists and returns a set of items that appear in both.

3. Create a dictionary representing a simple phone book, then add a function to search for a person by name.

4. Given a list of words, create a dictionary that counts how many times each word appears.

5. Use sets to find students who are taking both Math and Science classes from two separate lists.

## Wrapping Up

Python's data structures are powerful, flexible, and intuitive. Lists give you flexibility, tuples give you safety, dictionaries give you fast lookups, and sets give you uniqueness. Understanding when to use each one is a key skill that will make your code cleaner and more efficient.

Practice working with these structures, combine them in creative ways, and you'll soon find they become second nature. Ready to test your knowledge? Head over to Athena and try some data structure quizzes!
