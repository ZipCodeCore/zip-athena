# Python Object-Oriented Programming - Classes and Objects

Welcome to Object-Oriented Programming (OOP) in Python! If you've been writing Python code with just variables and functions, you're about to discover a whole new way of organizing your code. OOP lets you bundle data and functionality together into objects, making your code more organized, reusable, and easier to understand.

Think of classes as blueprints for creating objects. Just like an architect's blueprint can be used to build many houses, a class can be used to create many objects with the same structure but different data.

## Understanding Object-Oriented Programming

Object-Oriented Programming is a programming paradigm that organizes code around "objects"—combinations of data (attributes) and functions (methods) that work with that data. Python has been an object-oriented language from the beginning, and you've actually been using objects all along (strings, lists, dictionaries are all objects!).

The main principles of OOP are:
- **Encapsulation**: Bundling data and methods together
- **Inheritance**: Creating new classes based on existing ones
- **Polymorphism**: Using the same interface for different types
- **Abstraction**: Hiding complex implementation details

Let's start with the basics and build up from there!

## Classes and Objects

### Class Definition

A class is a blueprint for creating objects. You define it using the `class` keyword:

```python
class Dog:
    pass  # Empty class for now
```

Now we can create objects (instances) of this class:

```python
dog1 = Dog()
dog2 = Dog()

print(type(dog1))  # <class '__main__.Dog'>
```

### Adding Attributes

Let's make our class more useful by adding attributes (data):

```python
class Dog:
    pass

dog1 = Dog()
dog1.name = "Buddy"
dog1.age = 3
dog1.breed = "Golden Retriever"

dog2 = Dog()
dog2.name = "Max"
dog2.age = 5
dog2.breed = "Labrador"

print(f"{dog1.name} is a {dog1.age}-year-old {dog1.breed}")
# Output: Buddy is a 3-year-old Golden Retriever
```

But there's a better way to initialize objects...

### The `__init__()` Method

The `__init__()` method is a special method (called a constructor) that's automatically called when you create a new object:

```python
class Dog:
    def __init__(self, name, age, breed):
        self.name = name
        self.age = age
        self.breed = breed

# Now creating dogs is much cleaner
dog1 = Dog("Buddy", 3, "Golden Retriever")
dog2 = Dog("Max", 5, "Labrador")

print(f"{dog1.name} is {dog1.age} years old")
# Output: Buddy is 3 years old
```

### The `self` Parameter

You'll see `self` everywhere in Python classes. It represents the instance (object) itself and must be the first parameter of all instance methods:

```python
class Person:
    def __init__(self, name, age):
        self.name = name  # self.name is an instance attribute
        self.age = age
    
    def greet(self):
        print(f"Hello, my name is {self.name}")
    
    def have_birthday(self):
        self.age += 1
        print(f"Happy birthday! I'm now {self.age}")

person = Person("Alice", 25)
person.greet()          # Output: Hello, my name is Alice
person.have_birthday()  # Output: Happy birthday! I'm now 26
```

When you call `person.greet()`, Python automatically passes `person` as the first argument (`self`).

**Important**: You can technically name it anything, but everyone uses `self` by convention. Don't be that person who uses something else!

### Instance Methods

Instance methods are functions defined inside a class that operate on instances:

```python
class BankAccount:
    def __init__(self, owner, balance=0):
        self.owner = owner
        self.balance = balance
    
    def deposit(self, amount):
        if amount > 0:
            self.balance += amount
            print(f"Deposited ${amount}. New balance: ${self.balance}")
        else:
            print("Deposit amount must be positive")
    
    def withdraw(self, amount):
        if amount > self.balance:
            print("Insufficient funds")
        elif amount > 0:
            self.balance -= amount
            print(f"Withdrew ${amount}. New balance: ${self.balance}")
        else:
            print("Withdrawal amount must be positive")
    
    def get_balance(self):
        return self.balance

account = BankAccount("Alice", 100)
account.deposit(50)   # Deposited $50. New balance: $150
account.withdraw(30)  # Withdrew $30. New balance: $120
print(account.get_balance())  # 120
```

## Class Attributes vs Instance Attributes

### Instance Attributes

Instance attributes are unique to each object, defined in `__init__()` or elsewhere in instance methods:

```python
class Dog:
    def __init__(self, name):
        self.name = name  # Instance attribute - unique to each dog

dog1 = Dog("Buddy")
dog2 = Dog("Max")

print(dog1.name)  # Buddy
print(dog2.name)  # Max
```

### Class Attributes

Class attributes are shared by all instances of the class:

```python
class Dog:
    species = "Canis familiaris"  # Class attribute - shared by all dogs
    
    def __init__(self, name):
        self.name = name  # Instance attribute

dog1 = Dog("Buddy")
dog2 = Dog("Max")

print(dog1.species)  # Canis familiaris
print(dog2.species)  # Canis familiaris
print(Dog.species)   # Canis familiaris
```

You can access class attributes through the class itself or through instances. But be careful when modifying:

```python
class Counter:
    count = 0  # Class attribute
    
    def __init__(self):
        Counter.count += 1  # Modify class attribute correctly
        self.id = Counter.count

c1 = Counter()
c2 = Counter()
c3 = Counter()

print(Counter.count)  # 3
print(c1.id, c2.id, c3.id)  # 1 2 3
```

## Class Methods with `@classmethod` Decorator

Class methods work with the class itself rather than instances. They receive the class as the first argument (conventionally named `cls`):

```python
class Person:
    population = 0
    
    def __init__(self, name):
        self.name = name
        Person.population += 1
    
    @classmethod
    def get_population(cls):
        return cls.population
    
    @classmethod
    def create_anonymous(cls):
        return cls("Anonymous")

p1 = Person("Alice")
p2 = Person("Bob")

print(Person.get_population())  # 2

p3 = Person.create_anonymous()
print(p3.name)  # Anonymous
print(Person.get_population())  # 3
```

Class methods are often used as alternative constructors:

```python
class Date:
    def __init__(self, year, month, day):
        self.year = year
        self.month = month
        self.day = day
    
    @classmethod
    def from_string(cls, date_string):
        year, month, day = map(int, date_string.split('-'))
        return cls(year, month, day)
    
    def __str__(self):
        return f"{self.year}-{self.month:02d}-{self.day:02d}"

# Regular constructor
date1 = Date(2024, 3, 15)

# Alternative constructor using class method
date2 = Date.from_string("2024-03-15")

print(date1)  # 2024-03-15
print(date2)  # 2024-03-15
```

## Static Methods with `@staticmethod` Decorator

Static methods don't receive the instance (`self`) or class (`cls`) as the first argument. They're like regular functions that belong to the class namespace:

```python
class Math:
    @staticmethod
    def add(a, b):
        return a + b
    
    @staticmethod
    def multiply(a, b):
        return a * b
    
    @staticmethod
    def is_even(n):
        return n % 2 == 0

print(Math.add(5, 3))      # 8
print(Math.multiply(4, 7))  # 28
print(Math.is_even(10))     # True
```

Use static methods for utility functions that don't need access to instance or class data but logically belong with the class.

## Encapsulation

Encapsulation is about controlling access to data. While Python doesn't have true private attributes like some languages, it has conventions.

### Name Mangling (Strong Privacy)

Using double underscores `__` before an attribute name triggers name mangling:

```python
class BankAccount:
    def __init__(self, owner, balance):
        self.owner = owner
        self.__balance = balance  # Name mangling
    
    def deposit(self, amount):
        if amount > 0:
            self.__balance += amount
    
    def get_balance(self):
        return self.__balance

account = BankAccount("Alice", 100)
print(account.get_balance())  # 100

# This won't work:
# print(account.__balance)  # AttributeError

# But this will (name mangling makes it _ClassName__attribute):
# print(account._BankAccount__balance)  # 100 (don't do this!)
```

### Single Underscore Convention (Weak Privacy)

A single underscore `_` before a name is a convention meaning "this is internal, please don't use it directly":

```python
class Person:
    def __init__(self, name, age):
        self._age = age  # Convention: this is "protected"
        self.name = name
    
    def get_age(self):
        return self._age

person = Person("Alice", 25)
print(person.get_age())  # 25
# You CAN access person._age, but you shouldn't - it's a convention
```

### Property Decorators (Getters and Setters)

The `@property` decorator provides a Pythonic way to implement getters and setters:

```python
class Temperature:
    def __init__(self, celsius):
        self._celsius = celsius
    
    @property
    def celsius(self):
        """Getter for celsius"""
        return self._celsius
    
    @celsius.setter
    def celsius(self, value):
        """Setter for celsius with validation"""
        if value < -273.15:
            raise ValueError("Temperature below absolute zero is not possible")
        self._celsius = value
    
    @property
    def fahrenheit(self):
        """Calculated property"""
        return self._celsius * 9/5 + 32

temp = Temperature(25)
print(temp.celsius)     # 25
print(temp.fahrenheit)  # 77.0

temp.celsius = 30
print(temp.celsius)     # 30
print(temp.fahrenheit)  # 86.0

# This will raise an error:
# temp.celsius = -300  # ValueError: Temperature below absolute zero...
```

Properties make your code look cleaner—you access them like attributes but they're actually method calls!

## Special Methods (Magic Methods / Dunder Methods)

Special methods (also called magic methods or dunder methods, from "double underscore") customize how objects behave:

### `__str__()` and `__repr__()`

```python
class Book:
    def __init__(self, title, author):
        self.title = title
        self.author = author
    
    def __str__(self):
        """User-friendly string representation"""
        return f"{self.title} by {self.author}"
    
    def __repr__(self):
        """Developer-friendly string representation"""
        return f"Book('{self.title}', '{self.author}')"

book = Book("1984", "George Orwell")
print(str(book))   # 1984 by George Orwell
print(repr(book))  # Book('1984', 'George Orwell')
print(book)        # Uses __str__: 1984 by George Orwell
```

### `__len__()`, `__getitem__()`, and More

```python
class Playlist:
    def __init__(self, name):
        self.name = name
        self.songs = []
    
    def add_song(self, song):
        self.songs.append(song)
    
    def __len__(self):
        """Called by len()"""
        return len(self.songs)
    
    def __getitem__(self, index):
        """Called by playlist[index]"""
        return self.songs[index]
    
    def __contains__(self, song):
        """Called by 'song in playlist'"""
        return song in self.songs

playlist = Playlist("My Favorites")
playlist.add_song("Song A")
playlist.add_song("Song B")
playlist.add_song("Song C")

print(len(playlist))        # 3
print(playlist[0])          # Song A
print("Song B" in playlist) # True
```

### Comparison Methods

```python
class Point:
    def __init__(self, x, y):
        self.x = x
        self.y = y
    
    def __eq__(self, other):
        """Called by =="""
        return self.x == other.x and self.y == other.y
    
    def __lt__(self, other):
        """Called by <"""
        return (self.x ** 2 + self.y ** 2) < (other.x ** 2 + other.y ** 2)

p1 = Point(1, 2)
p2 = Point(1, 2)
p3 = Point(3, 4)

print(p1 == p2)  # True
print(p1 < p3)   # True (distance from origin)
```

## Operator Overloading

You can define how operators work with your objects:

```python
class Vector:
    def __init__(self, x, y):
        self.x = x
        self.y = y
    
    def __add__(self, other):
        """Called by +"""
        return Vector(self.x + other.x, self.y + other.y)
    
    def __mul__(self, scalar):
        """Called by *"""
        return Vector(self.x * scalar, self.y * scalar)
    
    def __str__(self):
        return f"Vector({self.x}, {self.y})"

v1 = Vector(1, 2)
v2 = Vector(3, 4)

v3 = v1 + v2        # Uses __add__
print(v3)           # Vector(4, 6)

v4 = v1 * 3         # Uses __mul__
print(v4)           # Vector(3, 6)
```

## Common Pitfalls

1. **Forgetting `self`**: Every instance method needs `self` as the first parameter
2. **Not using `__init__()`**: Initialize attributes in `__init__()` for clean, predictable objects
3. **Modifying class attributes incorrectly**: Use `ClassName.attribute` to modify class attributes, not `self.attribute`
4. **Confusing class and instance methods**: Use `@classmethod` when you need to work with the class, not instances
5. **Trying to make everything private**: Python's philosophy is "we're all consenting adults"—use conventions, not strict privacy

## Best Practices

1. **Use clear, descriptive names**: `class BankAccount` is better than `class BA`
2. **Initialize all instance attributes in `__init__()`**: Makes object state predictable
3. **Follow naming conventions**: 
   - `ClassName` (PascalCase for classes)
   - `method_name` (snake_case for methods)
   - `_internal_method` (leading underscore for internal use)
4. **Implement `__str__()` and `__repr__()`**: Makes debugging easier
5. **Use `@property` for computed attributes**: Cleaner than getter/setter methods
6. **Keep classes focused**: Each class should have one clear responsibility

## Practice Exercises

1. Create a `Student` class with name, ID, and grades, with methods to add grades and calculate GPA
2. Create a `Rectangle` class with width and height, using `@property` for area
3. Create a `ShoppingCart` class that uses `__len__()` and `__getitem__()` to act like a list
4. Create a `Time` class with hours, minutes, seconds and implement `__add__()` to add times together
5. Create a counter class using class attributes to track total instances created

---

## Quiz Questions

This section contains quiz questions that test understanding of Python classes and objects. Each question is formatted for easy extraction and database import.

### Question 1: Class Definition - Basic Concepts
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: python, oop, classes, objects

**Question Text**:
What is a class in Python?

**Answer Choices**:
- A) A variable that stores data
- B) A blueprint for creating objects that bundles data and methods
- C) A type of loop
- D) A built-in function

**Correct Answer**: B

**Explanation**:
A class is a blueprint or template for creating objects. It defines the attributes (data) and methods (functions) that objects created from the class will have. It's a way to bundle data and functionality together in object-oriented programming.

---

### Question 2: __init__ Method - Constructor
**Type**: code-output
**Difficulty**: 2
**Tags**: python, oop, init, constructor

**Question Text**:
What will be printed?

```python
class Dog:
    def __init__(self, name):
        self.name = name

dog = Dog("Buddy")
print(dog.name)
```

**Answer Choices**:
- A) Dog
- B) Buddy
- C) name
- D) Error

**Correct Answer**: B

**Explanation**:
The `__init__()` method is automatically called when creating a new object. It initializes the object with the provided arguments. In this case, when `Dog("Buddy")` is called, `__init__()` sets `self.name = "Buddy"`, so `dog.name` is "Buddy".

---

### Question 3: self Parameter - Instance Reference
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: python, oop, self

**Question Text**:
What does the `self` parameter in instance methods represent?

**Answer Choices**:
- A) The class itself
- B) The current instance of the class
- C) The parent class
- D) A keyword that must always be used

**Correct Answer**: B

**Explanation**:
The `self` parameter represents the current instance (object) that the method is being called on. When you call `dog.bark()`, Python automatically passes the `dog` object as `self` to the `bark()` method. While "self" is just a naming convention, it's universally used in Python code.

---

### Question 4: Class vs Instance Attributes - Shared vs Unique
**Type**: code-output
**Difficulty**: 4
**Tags**: python, oop, class-attributes, instance-attributes

**Question Text**:
What will be the output?

```python
class Counter:
    count = 0
    
    def __init__(self):
        Counter.count += 1

c1 = Counter()
c2 = Counter()
c3 = Counter()
print(Counter.count)
```

**Answer Choices**:
- A) 0
- B) 1
- C) 2
- D) 3

**Correct Answer**: D

**Explanation**:
`count` is a class attribute shared by all instances. Each time a Counter object is created, `__init__()` increments `Counter.count`. After creating three objects, `Counter.count` is 3. Class attributes are accessed and modified using the class name (`Counter.count`), not `self.count` when you want to modify the shared value.

---

### Question 5: @classmethod Decorator - Class Methods
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: python, oop, classmethod

**Question Text**:
What is the first parameter of a method decorated with `@classmethod`?

**Answer Choices**:
- A) self (the instance)
- B) cls (the class)
- C) class (the class name)
- D) No parameter is needed

**Correct Answer**: B

**Explanation**:
Methods decorated with `@classmethod` receive the class itself as the first parameter, conventionally named `cls`. This is different from instance methods which receive `self` (the instance). Class methods are useful for alternative constructors and operations that work with the class rather than individual instances.

---

### Question 6: @property Decorator - Getters and Setters
**Type**: code-output
**Difficulty**: 4
**Tags**: python, oop, property, encapsulation

**Question Text**:
What does the `@property` decorator do?

```python
class Circle:
    def __init__(self, radius):
        self._radius = radius
    
    @property
    def area(self):
        return 3.14 * self._radius ** 2

c = Circle(5)
print(c.area)
```

**Answer Choices**:
- A) Makes area a private attribute
- B) Allows area to be accessed like an attribute even though it's a method
- C) Prevents area from being modified
- D) Creates a class variable

**Correct Answer**: B

**Explanation**:
The `@property` decorator turns a method into a "getter" for a computed attribute. You can access it like `c.area` (without parentheses) even though it's defined as a method. This is Pythonic way to implement getters and computed properties. The value is calculated each time it's accessed.

---

### Question 7: Special Methods - __str__ vs __repr__
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: python, oop, special-methods, str, repr

**Question Text**:
What is the difference between `__str__()` and `__repr__()`?

**Answer Choices**:
- A) There is no difference, they do the same thing
- B) __str__() is for user-friendly output, __repr__() is for developer-friendly output
- C) __str__() is faster than __repr__()
- D) __repr__() is deprecated

**Correct Answer**: B

**Explanation**:
`__str__()` should return a user-friendly string representation (what you'd show to end users), while `__repr__()` should return a developer-friendly representation (ideally, code that could recreate the object). `print()` uses `__str__()` if available, otherwise `__repr__()`. For example, `__str__()` might return "Book: 1984" while `__repr__()` might return "Book('1984', 'George Orwell')".

---

### Question 8: Operator Overloading - Custom Operators
**Type**: code-output
**Difficulty**: 5
**Tags**: python, oop, operator-overloading

**Question Text**:
What allows this code to work?

```python
class Point:
    def __init__(self, x, y):
        self.x = x
        self.y = y
    
    def __add__(self, other):
        return Point(self.x + other.x, self.y + other.y)

p1 = Point(1, 2)
p2 = Point(3, 4)
p3 = p1 + p2
```

**Answer Choices**:
- A) The + operator automatically works with all classes
- B) The __add__() special method defines behavior for the + operator
- C) Points can always be added in Python
- D) This code will cause an error

**Correct Answer**: B

**Explanation**:
The `__add__()` special method (dunder method) defines how the `+` operator works with your objects. When you write `p1 + p2`, Python calls `p1.__add__(p2)`. This is called operator overloading—you're defining custom behavior for operators. Without `__add__()`, trying to add Point objects would raise a TypeError.

---
