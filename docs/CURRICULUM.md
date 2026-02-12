# Programming Curriculum Plan

This document outlines a comprehensive curriculum for teaching programming fundamentals from beginner to intermediate level. The curriculum is designed to support the Athena quiz application used at ZIPCodeWilmington.

## Table of Contents
- [Java Curriculum](#java-curriculum)
- [Python Curriculum](#python-curriculum)

---

## Java Curriculum

### 1. Introduction to Java and Development Environment
- 1.1 What is Java and its key characteristics (platform independence, compiled language)
- 1.2 Installing JDK (Java Development Kit)
- 1.3 Setting up an IDE (IntelliJ IDEA, Eclipse, or VS Code)
- 1.4 Understanding the Java compilation process (`.java` to `.class` files)
- 1.5 Writing and running your first Java program (`Hello, World!`)
- 1.6 Understanding the `main` method
- 1.7 Basic program structure and syntax rules

### 2. Variables and Data Types
- 2.1 Primitive data types
  - Integer types: `byte`, `short`, `int`, `long`
  - Floating-point types: `float`, `double`
  - Character type: `char`
  - Boolean type: `boolean`
- 2.2 Variable declaration and initialization
- 2.3 Variable naming conventions (camelCase)
- 2.4 Type casting and type conversion
  - Implicit (widening) conversions
  - Explicit (narrowing) conversions
- 2.5 Constants using `final` keyword
- 2.6 Literals (integer, floating-point, character, string, boolean)

### 3. Operators and Expressions
- 3.1 Arithmetic operators (`+`, `-`, `*`, `/`, `%`)
- 3.2 Assignment operators (`=`, `+=`, `-=`, `*=`, `/=`, `%=`)
- 3.3 Increment and decrement operators (`++`, `--`)
- 3.4 Comparison/Relational operators (`==`, `!=`, `>`, `<`, `>=`, `<=`)
- 3.5 Logical operators (`&&`, `||`, `!`)
- 3.6 Bitwise operators (`&`, `|`, `^`, `~`, `<<`, `>>`, `>>>`)
- 3.7 Ternary operator (`? :`)
- 3.8 Operator precedence and associativity
- 3.9 Expression evaluation

### 4. Control Structures
- 4.1 Conditional statements
  - `if` statement
  - `if-else` statement
  - `else-if` ladder
  - Nested `if` statements
  - `switch` statement (including switch expressions in modern Java)
- 4.2 Loop statements
  - `while` loop
  - `do-while` loop
  - `for` loop
  - Enhanced `for` loop (for-each)
  - Nested loops
- 4.3 Jump statements
  - `break` statement
  - `continue` statement
  - `return` statement
- 4.4 Understanding control flow and program logic

### 5. Methods (Functions)
- 5.1 Method declaration and definition
- 5.2 Method parameters and arguments
- 5.3 Return types and `return` statement
- 5.4 Method overloading
- 5.5 Variable scope (local vs. class level)
- 5.6 Method invocation and call stack
- 5.7 Passing arguments by value
- 5.8 Recursive methods
- 5.9 Variable arguments (varargs)
- 5.10 Method naming conventions

### 6. Arrays
- 6.1 One-dimensional arrays
  - Array declaration and initialization
  - Accessing array elements
  - Array length property
  - Iterating through arrays
- 6.2 Multi-dimensional arrays
  - Two-dimensional arrays (matrices)
  - Jagged arrays
- 6.3 Arrays class and utility methods
  - Sorting arrays
  - Searching arrays
  - Copying arrays
  - Comparing arrays
- 6.4 Common array operations and algorithms

### 7. Strings
- 7.1 String class and string literals
- 7.2 String immutability
- 7.3 String operations
  - Concatenation
  - Comparison (`equals()`, `equalsIgnoreCase()`, `compareTo()`)
  - Searching (`indexOf()`, `contains()`)
  - Extraction (`substring()`, `charAt()`)
  - Modification (`toUpperCase()`, `toLowerCase()`, `trim()`, `replace()`)
  - Splitting and joining strings
- 7.4 StringBuilder and StringBuffer classes
- 7.5 String formatting
- 7.6 Regular expressions basics

### 8. Object-Oriented Programming - Classes and Objects
- 8.1 Understanding object-oriented programming concepts
- 8.2 Classes and objects
  - Class definition
  - Object instantiation
  - Instance variables (fields)
  - Instance methods
  - Constructors (default and parameterized)
  - Constructor overloading
  - The `this` keyword
- 8.3 Access modifiers
  - `public`
  - `private`
  - `protected`
  - Package-private (default)
- 8.4 Encapsulation
  - Getters and setters
  - Data hiding
- 8.5 Static members
  - Static variables
  - Static methods
  - Static blocks
  - Static imports
- 8.6 The `final` keyword with classes and methods
- 8.7 Nested and inner classes
  - Static nested classes
  - Inner classes
  - Local classes
  - Anonymous classes

### 9. Object-Oriented Programming - Inheritance
- 9.1 Concept of inheritance and code reuse
- 9.2 Extending classes with `extends` keyword
- 9.3 The `super` keyword
- 9.4 Method overriding
- 9.5 The `@Override` annotation
- 9.6 Constructor chaining in inheritance
- 9.7 The Object class and its methods
  - `toString()`
  - `equals()`
  - `hashCode()`
  - `getClass()`
- 9.8 Abstract classes
  - Abstract methods
  - When to use abstract classes
- 9.9 The `final` keyword with inheritance
- 9.10 Single inheritance in Java
- 9.11 IS-A relationship

### 10. Object-Oriented Programming - Polymorphism
- 10.1 Concept of polymorphism
- 10.2 Compile-time polymorphism (method overloading)
- 10.3 Runtime polymorphism (method overriding)
- 10.4 Dynamic method dispatch
- 10.5 Upcasting and downcasting
- 10.6 The `instanceof` operator
- 10.7 Reference type vs. object type

### 11. Interfaces
- 11.1 Defining interfaces
- 11.2 Implementing interfaces with `implements` keyword
- 11.3 Multiple interface implementation
- 11.4 Interface inheritance (extending interfaces)
- 11.5 Default methods in interfaces
- 11.6 Static methods in interfaces
- 11.7 Functional interfaces
- 11.8 Marker interfaces
- 11.9 HAS-A relationship vs. IS-A relationship

### 12. Packages
- 12.1 Understanding packages
- 12.2 Creating packages
- 12.3 Importing packages and classes
  - Single type import
  - On-demand import
  - Static import
- 12.4 Package naming conventions
- 12.5 Access modifiers and packages
- 12.6 Built-in Java packages
  - `java.lang`
  - `java.util`
  - `java.io`
  - `java.math`

### 13. Exception Handling
- 13.1 Understanding exceptions and errors
- 13.2 Exception hierarchy
  - `Throwable` class
  - `Error` class
  - `Exception` class
  - Checked vs. unchecked exceptions
- 13.3 Try-catch blocks
- 13.4 Multiple catch blocks
- 13.5 The `finally` block
- 13.6 Try-with-resources statement
- 13.7 Throwing exceptions with `throw` keyword
- 13.8 Declaring exceptions with `throws` keyword
- 13.9 Creating custom exceptions
- 13.10 Best practices in exception handling
- 13.11 Common exceptions
  - `NullPointerException`
  - `ArrayIndexOutOfBoundsException`
  - `NumberFormatException`
  - `IllegalArgumentException`
  - `IOException`

### 14. File I/O
- 14.1 File and directory operations
  - `File` class
  - Creating, deleting, and checking files/directories
- 14.2 Character streams
  - `FileReader` and `FileWriter`
  - `BufferedReader` and `BufferedWriter`
- 14.3 Byte streams
  - `FileInputStream` and `FileOutputStream`
  - `BufferedInputStream` and `BufferedOutputStream`
- 14.4 Object serialization
  - `Serializable` interface
  - `ObjectInputStream` and `ObjectOutputStream`
- 14.5 NIO (New I/O) basics
  - `Path` and `Paths`
  - `Files` class
- 14.6 Reading and writing different file formats
  - Text files
  - CSV files
  - Binary files

### 15. Collections Framework
- 15.1 Introduction to Collections Framework
- 15.2 Collection hierarchy
- 15.3 List interface
  - `ArrayList` - using ArrayList
  - `LinkedList` - using LinkedList
  - When to use List implementations
- 15.4 Set interface
  - `HashSet` - using HashSet
  - `TreeSet` - using TreeSet
  - `LinkedHashSet` - using LinkedHashSet
- 15.5 Queue interface
  - `LinkedList` as Queue
  - `PriorityQueue` - using PriorityQueue
- 15.6 Map interface
  - `HashMap` - using HashMap
  - `TreeMap` - using TreeMap
  - `LinkedHashMap` - using LinkedHashMap
- 15.7 Collections utility class
  - Sorting collections
  - Searching collections
  - Shuffling and reversing
- 15.8 Iterators and enhanced for-loop
- 15.9 Comparable and Comparator interfaces

### 16. Generics - Using Generics
- 16.1 What are generics and why use them
- 16.2 Type safety benefits
- 16.3 Using generic classes
  - `ArrayList<String>`, `HashMap<Integer, String>`, etc.
- 16.4 Using generic methods
- 16.5 Type parameters and naming conventions (`T`, `E`, `K`, `V`)
- 16.6 Bounded type parameters in usage
- 16.7 Wildcards in usage
  - Unbounded wildcards (`?`)
  - Upper bounded wildcards (`? extends Type`)
  - Lower bounded wildcards (`? super Type`)
- 16.8 Generic type erasure concept
- 16.9 Using generics with collections
- **Note**: Writing custom generic classes/methods is beyond this curriculum scope

### 17. Basic Multithreading Concepts
- 17.1 Introduction to concurrency
- 17.2 Processes vs. threads
- 17.3 Creating threads
  - Extending `Thread` class
  - Implementing `Runnable` interface
- 17.4 Thread lifecycle
- 17.5 Basic thread methods
  - `start()`, `run()`, `sleep()`, `join()`
- 17.6 Thread synchronization basics
  - The `synchronized` keyword
  - Race conditions
- 17.7 Thread safety awareness

### 18. Java Standard Library Essentials
- 18.1 Math class
- 18.2 Random class
- 18.3 Date and time basics
  - `LocalDate`, `LocalTime`, `LocalDateTime` (using)
  - `DateTimeFormatter` (using)
- 18.4 Scanner class for input
- 18.5 System class
  - Standard input/output
  - System properties
- 18.6 Wrapper classes
  - Integer, Double, Boolean, etc.
  - Autoboxing and unboxing
- 18.7 Enum types

### 19. Lambda Expressions and Streams - Basic Usage
- 19.1 Introduction to lambda expressions
- 19.2 Lambda syntax
- 19.3 Functional interfaces and lambdas
- 19.4 Using lambda expressions with collections
- 19.5 Introduction to Stream API
  - Creating streams
  - Intermediate operations: `filter()`, `map()`, `sorted()`
  - Terminal operations: `forEach()`, `collect()`, `reduce()`
- 19.6 Using streams for basic data processing
- **Note**: Advanced stream operations and writing custom functional interfaces is beyond this curriculum scope

### 20. Best Practices and Code Quality
- 20.1 Code formatting and style
- 20.2 Naming conventions
- 20.3 Writing readable code
- 20.4 Comments and documentation (Javadoc basics)
- 20.5 Debugging techniques
- 20.6 Common pitfalls and how to avoid them
- 20.7 Understanding compile-time vs. runtime errors
- 20.8 Code organization principles

---

## Python Curriculum

### 1. Introduction to Python and Development Environment
- 1.1 What is Python and its key characteristics (interpreted, dynamic typing, versatile)
- 1.2 Python 2 vs. Python 3
- 1.3 Installing Python
- 1.4 Setting up an IDE/editor (PyCharm, VS Code, or IDLE)
- 1.5 Using the Python REPL (interactive shell)
- 1.6 Writing and running your first Python program (`Hello, World!`)
- 1.7 Python script files (`.py`)
- 1.8 Understanding indentation and its importance
- 1.9 Basic syntax rules

### 2. Variables and Data Types
- 2.1 Variables and dynamic typing
- 2.2 Variable naming conventions (snake_case)
- 2.3 Basic data types
  - Integers (`int`)
  - Floating-point numbers (`float`)
  - Strings (`str`)
  - Booleans (`bool`)
  - NoneType (`None`)
- 2.4 Type checking with `type()` function
- 2.5 Type conversion/casting
  - `int()`, `float()`, `str()`, `bool()`
- 2.6 Variable reassignment and dynamic typing
- 2.7 Constants convention (UPPER_CASE)
- 2.8 Multiple assignment
- 2.9 Unpacking

### 3. Operators and Expressions
- 3.1 Arithmetic operators (`+`, `-`, `*`, `/`, `//`, `%`, `**`)
- 3.2 Assignment operators (`=`, `+=`, `-=`, `*=`, `/=`, etc.)
- 3.3 Comparison operators (`==`, `!=`, `>`, `<`, `>=`, `<=`)
- 3.4 Logical operators (`and`, `or`, `not`)
- 3.5 Identity operators (`is`, `is not`)
- 3.6 Membership operators (`in`, `not in`)
- 3.7 Bitwise operators (`&`, `|`, `^`, `~`, `<<`, `>>`)
- 3.8 Operator precedence
- 3.9 Expression evaluation
- 3.10 Truthy and falsy values

### 4. Strings
- 4.1 String creation and literals
  - Single quotes, double quotes, triple quotes
- 4.2 String indexing and slicing
- 4.3 String immutability
- 4.4 String methods
  - Case conversion: `upper()`, `lower()`, `capitalize()`, `title()`
  - Searching: `find()`, `index()`, `count()`
  - Testing: `startswith()`, `endswith()`, `isalpha()`, `isdigit()`
  - Modification: `strip()`, `replace()`, `split()`, `join()`
- 4.5 String concatenation and repetition
- 4.6 String formatting
  - Old style: `%` formatting
  - `str.format()` method
  - f-strings (formatted string literals)
- 4.7 Raw strings
- 4.8 String encoding basics (Unicode, UTF-8)

### 5. Control Structures
- 5.1 Conditional statements
  - `if` statement
  - `if-else` statement
  - `elif` (else-if) chains
  - Nested conditionals
  - Conditional expressions (ternary operator)
- 5.2 Loop statements
  - `while` loop
  - `for` loop (iteration)
  - `range()` function
  - Nested loops
- 5.3 Loop control statements
  - `break` statement
  - `continue` statement
  - `pass` statement
- 5.4 `else` clause with loops
- 5.5 Understanding iteration and control flow

### 6. Data Structures - Lists
- 6.1 Creating lists
- 6.2 List indexing and slicing
- 6.3 List mutability
- 6.4 List methods
  - Adding elements: `append()`, `insert()`, `extend()`
  - Removing elements: `remove()`, `pop()`, `clear()`
  - Sorting: `sort()`, `reverse()`
  - Searching: `index()`, `count()`
- 6.5 List operations
  - Concatenation and repetition
  - Membership testing
  - Length with `len()`
- 6.6 List comprehensions
- 6.7 Nested lists (2D lists)
- 6.8 Copying lists (shallow vs. deep copy)
- 6.9 Common list algorithms

### 7. Data Structures - Tuples
- 7.1 Creating tuples
- 7.2 Tuple immutability
- 7.3 Tuple indexing and slicing
- 7.4 Tuple unpacking
- 7.5 Single element tuples
- 7.6 Tuple methods: `count()`, `index()`
- 7.7 When to use tuples vs. lists
- 7.8 Named tuples (using)

### 8. Data Structures - Dictionaries
- 8.1 Creating dictionaries
- 8.2 Key-value pairs
- 8.3 Dictionary mutability
- 8.4 Accessing values with keys
- 8.5 Dictionary methods
  - Adding/updating: setting values, `update()`
  - Removing: `pop()`, `popitem()`, `clear()`
  - Accessing: `get()`, `keys()`, `values()`, `items()`
- 8.6 Dictionary operations
  - Membership testing (keys)
  - Length
- 8.7 Iterating over dictionaries
- 8.8 Dictionary comprehensions
- 8.9 Nested dictionaries
- 8.10 Default dictionaries (using `defaultdict`)

### 9. Data Structures - Sets
- 9.1 Creating sets
- 9.2 Set mutability
- 9.3 Set uniqueness property
- 9.4 Set methods
  - Adding: `add()`, `update()`
  - Removing: `remove()`, `discard()`, `pop()`, `clear()`
- 9.5 Set operations
  - Union (`|`, `union()`)
  - Intersection (`&`, `intersection()`)
  - Difference (`-`, `difference()`)
  - Symmetric difference (`^`, `symmetric_difference()`)
  - Subset and superset testing
- 9.6 Frozen sets
- 9.7 Set comprehensions
- 9.8 When to use sets

### 10. Functions
- 10.1 Defining functions with `def` keyword
- 10.2 Function parameters and arguments
- 10.3 Return values and `return` statement
- 10.4 Default parameter values
- 10.5 Keyword arguments
- 10.6 Variable-length arguments
  - `*args` for positional arguments
  - `**kwargs` for keyword arguments
- 10.7 Function documentation strings (docstrings)
- 10.8 Scope and namespaces
  - Local scope
  - Enclosing scope
  - Global scope
  - Built-in scope
- 10.9 The `global` and `nonlocal` keywords
- 10.10 Lambda functions (anonymous functions)
- 10.11 Functions as first-class objects
- 10.12 Recursion

### 11. Modules and Packages
- 11.1 What are modules
- 11.2 Importing modules
  - `import module`
  - `from module import name`
  - `import module as alias`
  - `from module import *` (and why to avoid it)
- 11.3 Creating your own modules
- 11.4 Module search path
- 11.5 The `__name__` variable and `if __name__ == "__main__":`
- 11.6 Standard library modules
  - `math`
  - `random`
  - `datetime`
  - `os`
  - `sys`
  - `json`
  - `re` (regular expressions)
- 11.7 Packages and `__init__.py`
- 11.8 Installing third-party packages with `pip`

### 12. File I/O
- 12.1 Opening files with `open()` function
- 12.2 File modes (`r`, `w`, `a`, `r+`, `w+`, `a+`, `rb`, `wb`)
- 12.3 Reading files
  - `read()`
  - `readline()`
  - `readlines()`
  - Iterating over file object
- 12.4 Writing files
  - `write()`
  - `writelines()`
- 12.5 Closing files and the `with` statement (context managers)
- 12.6 File positioning: `seek()`, `tell()`
- 12.7 Working with file paths
  - `os.path` module
  - `pathlib` module (using `Path`)
- 12.8 Reading and writing different formats
  - Text files
  - CSV files (using `csv` module)
  - JSON files (using `json` module)
- 12.9 Exception handling with files

### 13. Exception Handling
- 13.1 Understanding exceptions
- 13.2 Common built-in exceptions
  - `SyntaxError`
  - `TypeError`
  - `ValueError`
  - `KeyError`
  - `IndexError`
  - `FileNotFoundError`
  - `ZeroDivisionError`
  - `AttributeError`
- 13.3 Try-except blocks
- 13.4 Multiple except clauses
- 13.5 The `else` clause with try-except
- 13.6 The `finally` clause
- 13.7 Raising exceptions with `raise`
- 13.8 Creating custom exceptions
- 13.9 Exception chaining
- 13.10 Best practices in exception handling
- 13.11 Assertions with `assert`

### 14. Object-Oriented Programming - Classes and Objects
- 14.1 Understanding object-oriented programming
- 14.2 Classes and objects
  - Class definition with `class` keyword
  - Object instantiation
  - Instance attributes
  - Instance methods
  - The `self` parameter
  - The `__init__()` method (constructor)
- 14.3 Class attributes vs. instance attributes
- 14.4 Class methods with `@classmethod` decorator
- 14.5 Static methods with `@staticmethod` decorator
- 14.6 Encapsulation
  - Name mangling with double underscore
  - Single underscore convention for "private"
  - Property decorators (`@property`, setter, deleter)
- 14.7 Special methods (magic methods/dunder methods)
  - `__str__()` and `__repr__()`
  - `__len__()`
  - `__getitem__()`, `__setitem__()`
  - `__eq__()`, `__lt__()`, etc.
  - `__call__()`
- 14.8 Operator overloading

### 15. Object-Oriented Programming - Inheritance
- 15.1 Concept of inheritance
- 15.2 Single inheritance
  - Creating subclasses
  - The `super()` function
  - Method overriding
- 15.3 The `isinstance()` and `issubclass()` functions
- 15.4 Multiple inheritance
  - Defining classes with multiple parents
  - Method Resolution Order (MRO)
  - The `__mro__` attribute
  - Diamond problem
- 15.5 Composition vs. inheritance
- 15.6 Abstract base classes
  - `abc` module
  - `@abstractmethod` decorator

### 16. Object-Oriented Programming - Polymorphism
- 16.1 Concept of polymorphism
- 16.2 Duck typing ("If it walks like a duck...")
- 16.3 Method overriding
- 16.4 Operator overloading (revisited)
- 16.5 Polymorphism through inheritance
- 16.6 Polymorphism through interfaces (abstract classes)

### 17. Working with Mixins - Using Mixins
- 17.1 What are mixins
- 17.2 Understanding mixin classes
- 17.3 Using mixins for code reuse
- 17.4 Combining mixins with inheritance
- 17.5 Common mixin patterns in Python
- 17.6 Examples of mixins in standard library
- 17.7 Best practices when using mixins
  - Mixin naming conventions
  - Keeping mixins focused
  - MRO considerations
- **Note**: Writing custom mixins is beyond this curriculum scope

### 18. Iterators and Generators
- 18.1 Understanding iteration protocol
- 18.2 Iterators
  - `__iter__()` and `__next__()` methods
  - `iter()` and `next()` functions
  - `StopIteration` exception
- 18.3 Generators
  - Generator functions with `yield`
  - Generator expressions
  - Benefits of generators (memory efficiency, lazy evaluation)
- 18.4 The `itertools` module (basic usage)

### 19. Comprehensions and Functional Programming
- 19.1 List comprehensions (revisited in depth)
- 19.2 Dictionary comprehensions
- 19.3 Set comprehensions
- 19.4 Generator expressions
- 19.5 Functional programming concepts
  - `map()` function
  - `filter()` function
  - `reduce()` function (from `functools`)
  - `zip()` function
  - `enumerate()` function
  - `sorted()` function with custom keys
- 19.6 Lambda functions in functional programming

### 20. Decorators - Using Decorators
- 20.1 What are decorators
- 20.2 Function decorators
- 20.3 Using common decorators
  - `@property`, `@staticmethod`, `@classmethod`
  - `@functools.wraps`
  - `@functools.lru_cache`
- 20.4 Decorator syntax and application
- 20.5 Chaining decorators
- 20.6 Class decorators (using)
- **Note**: Writing custom decorators is beyond this curriculum scope

### 21. Context Managers - Using Context Managers
- 21.1 Understanding context managers
- 21.2 The `with` statement
- 21.3 Common context managers
  - File operations
  - Thread locks
  - Database connections
- 21.4 The `contextlib` module (using)
- **Note**: Writing custom context managers with `__enter__` and `__exit__` is beyond this curriculum scope

### 22. Regular Expressions
- 22.1 Introduction to regular expressions
- 22.2 The `re` module
- 22.3 Basic patterns
  - Literal characters
  - Metacharacters (`.`, `^`, `$`, `*`, `+`, `?`, `[]`, `()`, `|`, `\`)
  - Character classes and ranges
  - Escape sequences
- 22.4 `re` module functions
  - `search()`, `match()`, `fullmatch()`
  - `findall()`, `finditer()`
  - `sub()`, `subn()`, `split()`
- 22.5 Match objects
- 22.6 Groups and capturing
- 22.7 Flags (case-insensitive, multiline, etc.)
- 22.8 Common regex patterns

### 23. Working with Dates and Times
- 23.1 The `datetime` module
- 23.2 Date objects: `date` class
- 23.3 Time objects: `time` class
- 23.4 Datetime objects: `datetime` class
- 23.5 Timedelta objects for time arithmetic
- 23.6 Formatting dates with `strftime()`
- 23.7 Parsing dates with `strptime()`
- 23.8 Time zones basics
  - Using `timezone` class
  - Creating timezone-aware datetime objects
  - UTC and timezone conversion
- 23.9 The `time` module basics

### 24. Common Standard Library Modules
- 24.1 `collections` module
  - `Counter`, `defaultdict`, `OrderedDict`, `deque`, `namedtuple`
- 24.2 `math` module
  - Mathematical functions and constants
- 24.3 `random` module
  - Random number generation
  - Random selection from sequences
- 24.4 `os` and `os.path` modules
  - Operating system interface
  - File and directory operations
- 24.5 `sys` module
  - System-specific parameters and functions
  - Command-line arguments
- 24.6 `copy` module
  - Shallow and deep copying
- 24.7 `pickle` module
  - Object serialization

### 25. Testing Basics
- 25.1 Importance of testing
- 25.2 Using `assert` statements
- 25.3 Introduction to `unittest` module
  - Test cases and test suites
  - `setUp()` and `tearDown()`
  - Assertion methods
- 25.4 Writing basic unit tests
- 25.5 Running tests
- 25.6 Test-driven development concept

### 26. Virtual Environments and Package Management
- 26.1 Why use virtual environments
- 26.2 Creating virtual environments with `venv`
- 26.3 Activating and deactivating environments
- 26.4 Installing packages with `pip`
- 26.5 Requirements files (`requirements.txt`)
- 26.6 Listing installed packages
- 26.7 Package versioning

### 27. Best Practices and Code Quality
- 27.1 PEP 8 style guide
  - Naming conventions
  - Indentation and whitespace
  - Line length
  - Imports organization
- 27.2 Code readability
- 27.3 Comments and docstrings
- 27.4 The Zen of Python (`import this`)
- 27.5 Common pitfalls and how to avoid them
- 27.6 Debugging techniques
  - Print debugging
  - Using `pdb` (Python debugger) basics
- 27.7 Code organization and project structure
- 27.8 Type hints (annotations) - basic usage

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
