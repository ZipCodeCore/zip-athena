# Programming Curriculum Plan

This document outlines a comprehensive curriculum for teaching programming fundamentals from beginner to intermediate level. The curriculum is designed to support the Athena quiz application used at ZIPCodeWilmington.

## Table of Contents
- [Java Curriculum](#java-curriculum)
- [Python Curriculum](#python-curriculum)

---

## Java Curriculum

### 1. Introduction to Java and Development Environment
- What is Java and its key characteristics (platform independence, compiled language)
- Installing JDK (Java Development Kit)
- Setting up an IDE (IntelliJ IDEA, Eclipse, or VS Code)
- Understanding the Java compilation process (`.java` to `.class` files)
- Writing and running your first Java program (`Hello, World!`)
- Understanding the `main` method
- Basic program structure and syntax rules

### 2. Variables and Data Types
- Primitive data types
  - Integer types: `byte`, `short`, `int`, `long`
  - Floating-point types: `float`, `double`
  - Character type: `char`
  - Boolean type: `boolean`
- Variable declaration and initialization
- Variable naming conventions (camelCase)
- Type casting and type conversion
  - Implicit (widening) conversions
  - Explicit (narrowing) conversions
- Constants using `final` keyword
- Literals (integer, floating-point, character, string, boolean)

### 3. Operators and Expressions
- Arithmetic operators (`+`, `-`, `*`, `/`, `%`)
- Assignment operators (`=`, `+=`, `-=`, `*=`, `/=`, `%=`)
- Increment and decrement operators (`++`, `--`)
- Comparison/Relational operators (`==`, `!=`, `>`, `<`, `>=`, `<=`)
- Logical operators (`&&`, `||`, `!`)
- Bitwise operators (`&`, `|`, `^`, `~`, `<<`, `>>`, `>>>`)
- Ternary operator (`? :`)
- Operator precedence and associativity
- Expression evaluation

### 4. Control Structures
- Conditional statements
  - `if` statement
  - `if-else` statement
  - `else-if` ladder
  - Nested `if` statements
  - `switch` statement (including switch expressions in modern Java)
- Loop statements
  - `while` loop
  - `do-while` loop
  - `for` loop
  - Enhanced `for` loop (for-each)
  - Nested loops
- Jump statements
  - `break` statement
  - `continue` statement
  - `return` statement
- Understanding control flow and program logic

### 5. Methods (Functions)
- Method declaration and definition
- Method parameters and arguments
- Return types and `return` statement
- Method overloading
- Variable scope (local vs. class level)
- Method invocation and call stack
- Passing arguments by value
- Recursive methods
- Variable arguments (varargs)
- Method naming conventions

### 6. Arrays
- One-dimensional arrays
  - Array declaration and initialization
  - Accessing array elements
  - Array length property
  - Iterating through arrays
- Multi-dimensional arrays
  - Two-dimensional arrays (matrices)
  - Jagged arrays
- Arrays class and utility methods
  - Sorting arrays
  - Searching arrays
  - Copying arrays
  - Comparing arrays
- Common array operations and algorithms

### 7. Strings
- String class and string literals
- String immutability
- String operations
  - Concatenation
  - Comparison (`equals()`, `equalsIgnoreCase()`, `compareTo()`)
  - Searching (`indexOf()`, `contains()`)
  - Extraction (`substring()`, `charAt()`)
  - Modification (`toUpperCase()`, `toLowerCase()`, `trim()`, `replace()`)
  - Splitting and joining strings
- StringBuilder and StringBuffer classes
- String formatting
- Regular expressions basics

### 8. Object-Oriented Programming - Classes and Objects
- Understanding object-oriented programming concepts
- Classes and objects
  - Class definition
  - Object instantiation
  - Instance variables (fields)
  - Instance methods
  - Constructors (default and parameterized)
  - Constructor overloading
  - The `this` keyword
- Access modifiers
  - `public`
  - `private`
  - `protected`
  - Package-private (default)
- Encapsulation
  - Getters and setters
  - Data hiding
- Static members
  - Static variables
  - Static methods
  - Static blocks
  - Static imports
- The `final` keyword with classes and methods
- Nested and inner classes
  - Static nested classes
  - Inner classes
  - Local classes
  - Anonymous classes

### 9. Object-Oriented Programming - Inheritance
- Concept of inheritance and code reuse
- Extending classes with `extends` keyword
- The `super` keyword
- Method overriding
- The `@Override` annotation
- Constructor chaining in inheritance
- The Object class and its methods
  - `toString()`
  - `equals()`
  - `hashCode()`
  - `getClass()`
- Abstract classes
  - Abstract methods
  - When to use abstract classes
- The `final` keyword with inheritance
- Single inheritance in Java
- IS-A relationship

### 10. Object-Oriented Programming - Polymorphism
- Concept of polymorphism
- Compile-time polymorphism (method overloading)
- Runtime polymorphism (method overriding)
- Dynamic method dispatch
- Upcasting and downcasting
- The `instanceof` operator
- Reference type vs. object type

### 11. Interfaces
- Defining interfaces
- Implementing interfaces with `implements` keyword
- Multiple interface implementation
- Interface inheritance (extending interfaces)
- Default methods in interfaces
- Static methods in interfaces
- Functional interfaces
- Marker interfaces
- HAS-A relationship vs. IS-A relationship

### 12. Packages
- Understanding packages
- Creating packages
- Importing packages and classes
  - Single type import
  - On-demand import
  - Static import
- Package naming conventions
- Access modifiers and packages
- Built-in Java packages
  - `java.lang`
  - `java.util`
  - `java.io`
  - `java.math`

### 13. Exception Handling
- Understanding exceptions and errors
- Exception hierarchy
  - `Throwable` class
  - `Error` class
  - `Exception` class
  - Checked vs. unchecked exceptions
- Try-catch blocks
- Multiple catch blocks
- The `finally` block
- Try-with-resources statement
- Throwing exceptions with `throw` keyword
- Declaring exceptions with `throws` keyword
- Creating custom exceptions
- Best practices in exception handling
- Common exceptions
  - `NullPointerException`
  - `ArrayIndexOutOfBoundsException`
  - `NumberFormatException`
  - `IllegalArgumentException`
  - `IOException`

### 14. File I/O
- File and directory operations
  - `File` class
  - Creating, deleting, and checking files/directories
- Character streams
  - `FileReader` and `FileWriter`
  - `BufferedReader` and `BufferedWriter`
- Byte streams
  - `FileInputStream` and `FileOutputStream`
  - `BufferedInputStream` and `BufferedOutputStream`
- Object serialization
  - `Serializable` interface
  - `ObjectInputStream` and `ObjectOutputStream`
- NIO (New I/O) basics
  - `Path` and `Paths`
  - `Files` class
- Reading and writing different file formats
  - Text files
  - CSV files
  - Binary files

### 15. Collections Framework
- Introduction to Collections Framework
- Collection hierarchy
- List interface
  - `ArrayList` - using ArrayList
  - `LinkedList` - using LinkedList
  - When to use List implementations
- Set interface
  - `HashSet` - using HashSet
  - `TreeSet` - using TreeSet
  - `LinkedHashSet` - using LinkedHashSet
- Queue interface
  - `LinkedList` as Queue
  - `PriorityQueue` - using PriorityQueue
- Map interface
  - `HashMap` - using HashMap
  - `TreeMap` - using TreeMap
  - `LinkedHashMap` - using LinkedHashMap
- Collections utility class
  - Sorting collections
  - Searching collections
  - Shuffling and reversing
- Iterators and enhanced for-loop
- Comparable and Comparator interfaces

### 16. Generics - Using Generics
- What are generics and why use them
- Type safety benefits
- Using generic classes
  - `ArrayList<String>`, `HashMap<Integer, String>`, etc.
- Using generic methods
- Type parameters and naming conventions (`T`, `E`, `K`, `V`)
- Bounded type parameters in usage
- Wildcards in usage
  - Unbounded wildcards (`?`)
  - Upper bounded wildcards (`? extends Type`)
  - Lower bounded wildcards (`? super Type`)
- Generic type erasure concept
- Using generics with collections
- **Note**: Writing custom generic classes/methods is beyond this curriculum scope

### 17. Basic Multithreading Concepts
- Introduction to concurrency
- Processes vs. threads
- Creating threads
  - Extending `Thread` class
  - Implementing `Runnable` interface
- Thread lifecycle
- Basic thread methods
  - `start()`, `run()`, `sleep()`, `join()`
- Thread synchronization basics
  - The `synchronized` keyword
  - Race conditions
- Thread safety awareness

### 18. Java Standard Library Essentials
- Math class
- Random class
- Date and time basics
  - `LocalDate`, `LocalTime`, `LocalDateTime` (using)
  - `DateTimeFormatter` (using)
- Scanner class for input
- System class
  - Standard input/output
  - System properties
- Wrapper classes
  - Integer, Double, Boolean, etc.
  - Autoboxing and unboxing
- Enum types

### 19. Lambda Expressions and Streams - Basic Usage
- Introduction to lambda expressions
- Lambda syntax
- Functional interfaces and lambdas
- Using lambda expressions with collections
- Introduction to Stream API
  - Creating streams
  - Intermediate operations: `filter()`, `map()`, `sorted()`
  - Terminal operations: `forEach()`, `collect()`, `reduce()`
- Using streams for basic data processing
- **Note**: Advanced stream operations and writing custom functional interfaces is beyond this curriculum scope

### 20. Best Practices and Code Quality
- Code formatting and style
- Naming conventions
- Writing readable code
- Comments and documentation (Javadoc basics)
- Debugging techniques
- Common pitfalls and how to avoid them
- Understanding compile-time vs. runtime errors
- Code organization principles

---

## Python Curriculum

### 1. Introduction to Python and Development Environment
- What is Python and its key characteristics (interpreted, dynamic typing, versatile)
- Python 2 vs. Python 3
- Installing Python
- Setting up an IDE/editor (PyCharm, VS Code, or IDLE)
- Using the Python REPL (interactive shell)
- Writing and running your first Python program (`Hello, World!`)
- Python script files (`.py`)
- Understanding indentation and its importance
- Basic syntax rules

### 2. Variables and Data Types
- Variables and dynamic typing
- Variable naming conventions (snake_case)
- Basic data types
  - Integers (`int`)
  - Floating-point numbers (`float`)
  - Strings (`str`)
  - Booleans (`bool`)
  - NoneType (`None`)
- Type checking with `type()` function
- Type conversion/casting
  - `int()`, `float()`, `str()`, `bool()`
- Variable reassignment and dynamic typing
- Constants convention (UPPER_CASE)
- Multiple assignment
- Unpacking

### 3. Operators and Expressions
- Arithmetic operators (`+`, `-`, `*`, `/`, `//`, `%`, `**`)
- Assignment operators (`=`, `+=`, `-=`, `*=`, `/=`, etc.)
- Comparison operators (`==`, `!=`, `>`, `<`, `>=`, `<=`)
- Logical operators (`and`, `or`, `not`)
- Identity operators (`is`, `is not`)
- Membership operators (`in`, `not in`)
- Bitwise operators (`&`, `|`, `^`, `~`, `<<`, `>>`)
- Operator precedence
- Expression evaluation
- Truthy and falsy values

### 4. Strings
- String creation and literals
  - Single quotes, double quotes, triple quotes
- String indexing and slicing
- String immutability
- String methods
  - Case conversion: `upper()`, `lower()`, `capitalize()`, `title()`
  - Searching: `find()`, `index()`, `count()`
  - Testing: `startswith()`, `endswith()`, `isalpha()`, `isdigit()`
  - Modification: `strip()`, `replace()`, `split()`, `join()`
- String concatenation and repetition
- String formatting
  - Old style: `%` formatting
  - `str.format()` method
  - f-strings (formatted string literals)
- Raw strings
- String encoding basics (Unicode, UTF-8)

### 5. Control Structures
- Conditional statements
  - `if` statement
  - `if-else` statement
  - `elif` (else-if) chains
  - Nested conditionals
  - Conditional expressions (ternary operator)
- Loop statements
  - `while` loop
  - `for` loop (iteration)
  - `range()` function
  - Nested loops
- Loop control statements
  - `break` statement
  - `continue` statement
  - `pass` statement
- `else` clause with loops
- Understanding iteration and control flow

### 6. Data Structures - Lists
- Creating lists
- List indexing and slicing
- List mutability
- List methods
  - Adding elements: `append()`, `insert()`, `extend()`
  - Removing elements: `remove()`, `pop()`, `clear()`
  - Sorting: `sort()`, `reverse()`
  - Searching: `index()`, `count()`
- List operations
  - Concatenation and repetition
  - Membership testing
  - Length with `len()`
- List comprehensions
- Nested lists (2D lists)
- Copying lists (shallow vs. deep copy)
- Common list algorithms

### 7. Data Structures - Tuples
- Creating tuples
- Tuple immutability
- Tuple indexing and slicing
- Tuple unpacking
- Single element tuples
- Tuple methods: `count()`, `index()`
- When to use tuples vs. lists
- Named tuples (using)

### 8. Data Structures - Dictionaries
- Creating dictionaries
- Key-value pairs
- Dictionary mutability
- Accessing values with keys
- Dictionary methods
  - Adding/updating: setting values, `update()`
  - Removing: `pop()`, `popitem()`, `clear()`
  - Accessing: `get()`, `keys()`, `values()`, `items()`
- Dictionary operations
  - Membership testing (keys)
  - Length
- Iterating over dictionaries
- Dictionary comprehensions
- Nested dictionaries
- Default dictionaries (using `defaultdict`)

### 9. Data Structures - Sets
- Creating sets
- Set mutability
- Set uniqueness property
- Set methods
  - Adding: `add()`, `update()`
  - Removing: `remove()`, `discard()`, `pop()`, `clear()`
- Set operations
  - Union (`|`, `union()`)
  - Intersection (`&`, `intersection()`)
  - Difference (`-`, `difference()`)
  - Symmetric difference (`^`, `symmetric_difference()`)
  - Subset and superset testing
- Frozen sets
- Set comprehensions
- When to use sets

### 10. Functions
- Defining functions with `def` keyword
- Function parameters and arguments
- Return values and `return` statement
- Default parameter values
- Keyword arguments
- Variable-length arguments
  - `*args` for positional arguments
  - `**kwargs` for keyword arguments
- Function documentation strings (docstrings)
- Scope and namespaces
  - Local scope
  - Enclosing scope
  - Global scope
  - Built-in scope
- The `global` and `nonlocal` keywords
- Lambda functions (anonymous functions)
- Functions as first-class objects
- Recursion

### 11. Modules and Packages
- What are modules
- Importing modules
  - `import module`
  - `from module import name`
  - `import module as alias`
  - `from module import *` (and why to avoid it)
- Creating your own modules
- Module search path
- The `__name__` variable and `if __name__ == "__main__":`
- Standard library modules
  - `math`
  - `random`
  - `datetime`
  - `os`
  - `sys`
  - `json`
  - `re` (regular expressions)
- Packages and `__init__.py`
- Installing third-party packages with `pip`

### 12. File I/O
- Opening files with `open()` function
- File modes (`r`, `w`, `a`, `r+`, `w+`, `a+`, `rb`, `wb`)
- Reading files
  - `read()`
  - `readline()`
  - `readlines()`
  - Iterating over file object
- Writing files
  - `write()`
  - `writelines()`
- Closing files and the `with` statement (context managers)
- File positioning: `seek()`, `tell()`
- Working with file paths
  - `os.path` module
  - `pathlib` module (using `Path`)
- Reading and writing different formats
  - Text files
  - CSV files (using `csv` module)
  - JSON files (using `json` module)
- Exception handling with files

### 13. Exception Handling
- Understanding exceptions
- Common built-in exceptions
  - `SyntaxError`
  - `TypeError`
  - `ValueError`
  - `KeyError`
  - `IndexError`
  - `FileNotFoundError`
  - `ZeroDivisionError`
  - `AttributeError`
- Try-except blocks
- Multiple except clauses
- The `else` clause with try-except
- The `finally` clause
- Raising exceptions with `raise`
- Creating custom exceptions
- Exception chaining
- Best practices in exception handling
- Assertions with `assert`

### 14. Object-Oriented Programming - Classes and Objects
- Understanding object-oriented programming
- Classes and objects
  - Class definition with `class` keyword
  - Object instantiation
  - Instance attributes
  - Instance methods
  - The `self` parameter
  - The `__init__()` method (constructor)
- Class attributes vs. instance attributes
- Class methods with `@classmethod` decorator
- Static methods with `@staticmethod` decorator
- Encapsulation
  - Name mangling with double underscore
  - Single underscore convention for "private"
  - Property decorators (`@property`, setter, deleter)
- Special methods (magic methods/dunder methods)
  - `__str__()` and `__repr__()`
  - `__len__()`
  - `__getitem__()`, `__setitem__()`
  - `__eq__()`, `__lt__()`, etc.
  - `__call__()`
- Operator overloading

### 15. Object-Oriented Programming - Inheritance
- Concept of inheritance
- Single inheritance
  - Creating subclasses
  - The `super()` function
  - Method overriding
- The `isinstance()` and `issubclass()` functions
- Multiple inheritance
  - Defining classes with multiple parents
  - Method Resolution Order (MRO)
  - The `__mro__` attribute
  - Diamond problem
- Composition vs. inheritance
- Abstract base classes
  - `abc` module
  - `@abstractmethod` decorator

### 16. Object-Oriented Programming - Polymorphism
- Concept of polymorphism
- Duck typing ("If it walks like a duck...")
- Method overriding
- Operator overloading (revisited)
- Polymorphism through inheritance
- Polymorphism through interfaces (abstract classes)

### 17. Working with Mixins - Using Mixins
- What are mixins
- Understanding mixin classes
- Using mixins for code reuse
- Combining mixins with inheritance
- Common mixin patterns in Python
- Examples of mixins in standard library
- Best practices when using mixins
  - Mixin naming conventions
  - Keeping mixins focused
  - MRO considerations
- **Note**: Writing custom mixins is beyond this curriculum scope

### 18. Iterators and Generators
- Understanding iteration protocol
- Iterators
  - `__iter__()` and `__next__()` methods
  - `iter()` and `next()` functions
  - `StopIteration` exception
- Generators
  - Generator functions with `yield`
  - Generator expressions
  - Benefits of generators (memory efficiency, lazy evaluation)
- The `itertools` module (basic usage)

### 19. Comprehensions and Functional Programming
- List comprehensions (revisited in depth)
- Dictionary comprehensions
- Set comprehensions
- Generator expressions
- Functional programming concepts
  - `map()` function
  - `filter()` function
  - `reduce()` function (from `functools`)
  - `zip()` function
  - `enumerate()` function
  - `sorted()` function with custom keys
- Lambda functions in functional programming

### 20. Decorators - Using Decorators
- What are decorators
- Function decorators
- Using common decorators
  - `@property`, `@staticmethod`, `@classmethod`
  - `@functools.wraps`
  - `@functools.lru_cache`
- Decorator syntax and application
- Chaining decorators
- Class decorators (using)
- **Note**: Writing custom decorators is beyond this curriculum scope

### 21. Context Managers - Using Context Managers
- Understanding context managers
- The `with` statement
- Common context managers
  - File operations
  - Thread locks
  - Database connections
- The `contextlib` module (using)
- **Note**: Writing custom context managers with `__enter__` and `__exit__` is beyond this curriculum scope

### 22. Regular Expressions
- Introduction to regular expressions
- The `re` module
- Basic patterns
  - Literal characters
  - Metacharacters (`.`, `^`, `$`, `*`, `+`, `?`, `[]`, `()`, `|`, `\`)
  - Character classes and ranges
  - Escape sequences
- `re` module functions
  - `search()`, `match()`, `fullmatch()`
  - `findall()`, `finditer()`
  - `sub()`, `subn()`, `split()`
- Match objects
- Groups and capturing
- Flags (case-insensitive, multiline, etc.)
- Common regex patterns

### 23. Working with Dates and Times
- The `datetime` module
- Date objects: `date` class
- Time objects: `time` class
- Datetime objects: `datetime` class
- Timedelta objects for time arithmetic
- Formatting dates with `strftime()`
- Parsing dates with `strptime()`
- Time zones basics
  - Using `timezone` class
  - Creating timezone-aware datetime objects
  - UTC and timezone conversion
- The `time` module basics

### 24. Common Standard Library Modules
- `collections` module
  - `Counter`, `defaultdict`, `OrderedDict`, `deque`, `namedtuple`
- `math` module
  - Mathematical functions and constants
- `random` module
  - Random number generation
  - Random selection from sequences
- `os` and `os.path` modules
  - Operating system interface
  - File and directory operations
- `sys` module
  - System-specific parameters and functions
  - Command-line arguments
- `copy` module
  - Shallow and deep copying
- `pickle` module
  - Object serialization

### 25. Testing Basics
- Importance of testing
- Using `assert` statements
- Introduction to `unittest` module
  - Test cases and test suites
  - `setUp()` and `tearDown()`
  - Assertion methods
- Writing basic unit tests
- Running tests
- Test-driven development concept

### 26. Virtual Environments and Package Management
- Why use virtual environments
- Creating virtual environments with `venv`
- Activating and deactivating environments
- Installing packages with `pip`
- Requirements files (`requirements.txt`)
- Listing installed packages
- Package versioning

### 27. Best Practices and Code Quality
- PEP 8 style guide
  - Naming conventions
  - Indentation and whitespace
  - Line length
  - Imports organization
- Code readability
- Comments and docstrings
- The Zen of Python (`import this`)
- Common pitfalls and how to avoid them
- Debugging techniques
  - Print debugging
  - Using `pdb` (Python debugger) basics
- Code organization and project structure
- Type hints (annotations) - basic usage

---

## Notes on Curriculum Scope

### Java Curriculum
- **Includes**: Core Java fundamentals through using generics, collections, and basic lambda/stream operations
- **Excludes**: Writing custom generic classes, advanced concurrency, JavaFX/Swing, networking, JDBC, frameworks (Spring, Hibernate), writing annotations, reflection API

### Python Curriculum  
- **Includes**: Python fundamentals through using mixins, decorators, and context managers
- **Excludes**: Writing custom mixins, writing complex decorators, metaclasses, advanced concurrency (threading/multiprocessing), async/await, web frameworks (Django, Flask), NumPy/Pandas, GUI programming, networking beyond basics

## Implementation in Athena

This curriculum provides a comprehensive framework for creating quiz questions in the Athena application. Each topic and subtopic can be mapped to:

1. **Topics**: Major sections (e.g., "Control Structures", "Object-Oriented Programming")
2. **Question Categories**: Individual subtopics for granular quiz generation
3. **Tags**: Cross-cutting concepts that appear in multiple topics (e.g., "loops", "inheritance", "exception handling")
4. **Difficulty Levels**: Questions can range from level 1 (basic concepts) to level 7 (complex applications of intermediate concepts)

The structured progression from fundamentals to intermediate concepts ensures students build knowledge incrementally, which aligns well with spaced repetition learning algorithms used in Athena.
