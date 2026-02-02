# Python Object-Oriented Programming - Polymorphism

Polymorphism—the word might sound intimidating, but the concept is actually quite natural in Python. It comes from Greek words meaning "many forms," and in programming, it means objects of different types can be used through the same interface. Python's dynamic typing makes polymorphism particularly elegant and natural.

## Concept of Polymorphism

Polymorphism allows you to use objects of different types through a common interface. You can write code that works with a general type, and it'll work with any specific type that follows that interface.

Here's a simple example:

```python
class Dog:
    def speak(self):
        return "Woof!"

class Cat:
    def speak(self):
        return "Meow!"

class Duck:
    def speak(self):
        return "Quack!"

# Polymorphic function - works with any object that has a speak() method
def make_sound(animal):
    print(animal.speak())

dog = Dog()
cat = Cat()
duck = Duck()

make_sound(dog)   # Woof!
make_sound(cat)   # Meow!
make_sound(duck)  # Quack!
```

The `make_sound()` function doesn't care what type of animal it receives—it just needs an object with a `speak()` method. That's polymorphism!

## Duck Typing

Python uses "duck typing"—if it walks like a duck and quacks like a duck, it's a duck. The type doesn't matter; what matters is whether the object has the methods and attributes you need.

```python
class Bird:
    def fly(self):
        print("Bird flying")

class Airplane:
    def fly(self):
        print("Airplane flying")

class Butterfly:
    def fly(self):
        print("Butterfly flying")

def make_it_fly(thing):
    thing.fly()  # Don't care what 'thing' is, as long as it can fly

make_it_fly(Bird())       # Bird flying
make_it_fly(Airplane())   # Airplane flying
make_it_fly(Butterfly())  # Butterfly flying
```

These classes have nothing in common (no shared parent class), but they all have a `fly()` method, so they work with `make_it_fly()`.

### The Philosophy Behind Duck Typing

```python
# Traditional approach (type checking)
def process_payment_old(payment):
    if isinstance(payment, CreditCard):
        payment.charge()
    elif isinstance(payment, PayPal):
        payment.send()
    else:
        raise TypeError("Unknown payment type")

# Duck typing approach (more Pythonic)
def process_payment(payment):
    payment.process()  # Just call the method - let it fail if it doesn't exist

# Both of these work:
class CreditCard:
    def process(self):
        print("Processing credit card payment")

class PayPal:
    def process(self):
        print("Processing PayPal payment")

class Bitcoin:
    def process(self):
        print("Processing Bitcoin payment")

# All work with the same function
process_payment(CreditCard())  # No type checking needed!
process_payment(PayPal())
process_payment(Bitcoin())
```

### Duck Typing with Built-in Protocols

Python has several "protocols" (informal interfaces) that objects can implement:

```python
# The "iterable" protocol - anything with __iter__
class Countdown:
    def __init__(self, start):
        self.start = start
    
    def __iter__(self):
        current = self.start
        while current > 0:
            yield current
            current -= 1

# Works in a for loop because it implements __iter__
for num in Countdown(5):
    print(num, end=' ')  # 5 4 3 2 1

# The "container" protocol - anything with __contains__, __len__, __getitem__
class CustomList:
    def __init__(self, items):
        self.items = items
    
    def __contains__(self, item):
        return item in self.items
    
    def __len__(self):
        return len(self.items)
    
    def __getitem__(self, index):
        return self.items[index]

mylist = CustomList([1, 2, 3, 4, 5])
print(3 in mylist)      # True (uses __contains__)
print(len(mylist))      # 5 (uses __len__)
print(mylist[2])        # 3 (uses __getitem__)
```

## Method Overriding

Method overriding is when a subclass provides its own implementation of a method defined in its parent class. This is a key form of polymorphism:

```python
class Shape:
    def area(self):
        return 0
    
    def perimeter(self):
        return 0

class Circle(Shape):
    def __init__(self, radius):
        self.radius = radius
    
    def area(self):  # Override parent's method
        return 3.14159 * self.radius ** 2
    
    def perimeter(self):  # Override parent's method
        return 2 * 3.14159 * self.radius

class Rectangle(Shape):
    def __init__(self, width, height):
        self.width = width
        self.height = height
    
    def area(self):  # Override parent's method
        return self.width * self.height
    
    def perimeter(self):  # Override parent's method
        return 2 * (self.width + self.height)

# Polymorphic usage
shapes = [Circle(5), Rectangle(4, 6), Circle(3)]

for shape in shapes:
    print(f"Area: {shape.area()}, Perimeter: {shape.perimeter()}")
```

Each shape calculates area and perimeter differently, but they all use the same method names.

### Calling Parent's Method

You can call the parent's version of a method using `super()`:

```python
class Animal:
    def make_sound(self):
        print("Some generic sound")

class Dog(Animal):
    def make_sound(self):
        super().make_sound()  # Call parent's method
        print("Woof!")

dog = Dog()
dog.make_sound()
# Output:
# Some generic sound
# Woof!
```

## Operator Overloading

Operator overloading lets you define how operators work with your objects. Python makes this easy with special methods:

### Arithmetic Operators

```python
class Vector:
    def __init__(self, x, y):
        self.x = x
        self.y = y
    
    def __add__(self, other):
        """Called for +"""
        return Vector(self.x + other.x, self.y + other.y)
    
    def __sub__(self, other):
        """Called for -"""
        return Vector(self.x - other.x, self.y - other.y)
    
    def __mul__(self, scalar):
        """Called for *"""
        return Vector(self.x * scalar, self.y * scalar)
    
    def __str__(self):
        return f"Vector({self.x}, {self.y})"

v1 = Vector(1, 2)
v2 = Vector(3, 4)

v3 = v1 + v2       # Calls __add__
print(v3)          # Vector(4, 6)

v4 = v2 - v1       # Calls __sub__
print(v4)          # Vector(2, 2)

v5 = v1 * 3        # Calls __mul__
print(v5)          # Vector(3, 6)
```

### Comparison Operators

```python
class Person:
    def __init__(self, name, age):
        self.name = name
        self.age = age
    
    def __eq__(self, other):
        """Called for =="""
        return self.age == other.age
    
    def __lt__(self, other):
        """Called for <"""
        return self.age < other.age
    
    def __le__(self, other):
        """Called for <="""
        return self.age <= other.age
    
    def __gt__(self, other):
        """Called for >"""
        return self.age > other.age
    
    def __str__(self):
        return f"{self.name} ({self.age})"

p1 = Person("Alice", 25)
p2 = Person("Bob", 30)
p3 = Person("Charlie", 25)

print(p1 == p3)  # True (same age)
print(p1 < p2)   # True (25 < 30)
print(p2 > p1)   # True (30 > 25)

# Can even sort them!
people = [p2, p3, p1]
sorted_people = sorted(people)
for person in sorted_people:
    print(person)
# Output:
# Alice (25)
# Charlie (25)
# Bob (30)
```

### Other Useful Operator Overloads

```python
class Money:
    def __init__(self, dollars):
        self.dollars = dollars
    
    def __add__(self, other):
        """Addition: money1 + money2"""
        return Money(self.dollars + other.dollars)
    
    def __iadd__(self, other):
        """In-place addition: money1 += money2"""
        self.dollars += other.dollars
        return self
    
    def __neg__(self):
        """Negation: -money"""
        return Money(-self.dollars)
    
    def __abs__(self):
        """Absolute value: abs(money)"""
        return Money(abs(self.dollars))
    
    def __str__(self):
        return f"${self.dollars:.2f}"

m1 = Money(50)
m2 = Money(30)

print(m1 + m2)    # $80.00
m1 += m2
print(m1)         # $80.00

debt = Money(-25)
print(debt)       # $-25.00
print(-debt)      # $25.00
print(abs(debt))  # $25.00
```

## Polymorphism Through Inheritance

The most common form of polymorphism uses inheritance:

```python
class PaymentMethod:
    def pay(self, amount):
        raise NotImplementedError("Subclasses must implement pay()")

class CreditCard(PaymentMethod):
    def __init__(self, card_number):
        self.card_number = card_number
    
    def pay(self, amount):
        print(f"Charging ${amount} to credit card {self.card_number[-4:]}")

class PayPal(PaymentMethod):
    def __init__(self, email):
        self.email = email
    
    def pay(self, amount):
        print(f"Sending ${amount} via PayPal to {self.email}")

class Bitcoin(PaymentMethod):
    def __init__(self, wallet_address):
        self.wallet_address = wallet_address
    
    def pay(self, amount):
        print(f"Transferring ${amount} to Bitcoin wallet {self.wallet_address[:8]}...")

# Polymorphic function - works with any PaymentMethod
def process_payment(payment_method, amount):
    payment_method.pay(amount)

# All use the same interface
payments = [
    CreditCard("1234567890123456"),
    PayPal("user@example.com"),
    Bitcoin("1A2B3C4D5E6F7G8H")
]

for payment in payments:
    process_payment(payment, 99.99)

# Output:
# Charging $99.99 to credit card 3456
# Sending $99.99 via PayPal to user@example.com
# Transferring $99.99 to Bitcoin wallet 1A2B3C4D...
```

## Polymorphism Through Abstract Classes

Using abstract base classes enforces that subclasses implement certain methods:

```python
from abc import ABC, abstractmethod

class Database(ABC):
    @abstractmethod
    def connect(self):
        pass
    
    @abstractmethod
    def execute_query(self, query):
        pass
    
    @abstractmethod
    def close(self):
        pass

class MySQLDatabase(Database):
    def connect(self):
        print("Connecting to MySQL database")
    
    def execute_query(self, query):
        print(f"Executing MySQL query: {query}")
    
    def close(self):
        print("Closing MySQL connection")

class PostgreSQLDatabase(Database):
    def connect(self):
        print("Connecting to PostgreSQL database")
    
    def execute_query(self, query):
        print(f"Executing PostgreSQL query: {query}")
    
    def close(self):
        print("Closing PostgreSQL connection")

# Polymorphic function
def run_query(database, query):
    database.connect()
    database.execute_query(query)
    database.close()

# Works with any Database subclass
mysql = MySQLDatabase()
postgres = PostgreSQLDatabase()

run_query(mysql, "SELECT * FROM users")
print()
run_query(postgres, "SELECT * FROM users")
```

## Polymorphism with Built-in Functions

Many built-in Python functions use polymorphism:

```python
class Book:
    def __init__(self, title, pages):
        self.title = title
        self.pages = pages
    
    def __len__(self):
        return self.pages
    
    def __str__(self):
        return f"{self.title} ({self.pages} pages)"

class Playlist:
    def __init__(self, name, songs):
        self.name = name
        self.songs = songs
    
    def __len__(self):
        return len(self.songs)
    
    def __str__(self):
        return f"{self.name} ({len(self.songs)} songs)"

# len() works polymorphically with any object that has __len__
book = Book("Python Guide", 350)
playlist = Playlist("Favorites", ["Song1", "Song2", "Song3"])

print(len(book))      # 350
print(len(playlist))  # 3
print(len([1,2,3]))   # 3 (built-in list)
print(len("hello"))   # 5 (built-in string)

# str() works polymorphically with any object that has __str__
print(str(book))      # Python Guide (350 pages)
print(str(playlist))  # Favorites (3 songs)
```

## Real-World Example: Plugin System

Polymorphism is perfect for plugin systems where you don't know in advance what types you'll work with:

```python
class Plugin:
    def execute(self):
        raise NotImplementedError("Plugins must implement execute()")

class PDFExportPlugin(Plugin):
    def execute(self):
        print("Exporting to PDF...")
        return "document.pdf"

class EmailNotificationPlugin(Plugin):
    def execute(self):
        print("Sending email notification...")
        return "Email sent"

class DataBackupPlugin(Plugin):
    def execute(self):
        print("Backing up data...")
        return "Backup completed"

class PluginManager:
    def __init__(self):
        self.plugins = []
    
    def register(self, plugin):
        self.plugins.append(plugin)
    
    def run_all(self):
        results = []
        for plugin in self.plugins:
            result = plugin.execute()  # Polymorphic call
            results.append(result)
        return results

# Set up and run plugins
manager = PluginManager()
manager.register(PDFExportPlugin())
manager.register(EmailNotificationPlugin())
manager.register(DataBackupPlugin())

results = manager.run_all()
print("Results:", results)
```

## Common Pitfalls

1. **Overusing operator overloading**: Make operators intuitive (e.g., `+` for addition makes sense, `+` for deleting doesn't)
2. **Not implementing all comparison operators**: If you implement `__eq__`, consider implementing `__lt__`, `__le__`, etc.
3. **Breaking the Liskov Substitution Principle**: Subclasses should be usable wherever parent classes are expected
4. **Forgetting to return `self` in `__iadd__`, `__imul__`, etc.**: In-place operators should return self
5. **Using `isinstance()` too much**: Often a sign you should redesign to use polymorphism instead

## Best Practices

1. **Favor duck typing**: Write code that works with any object that has the right methods
2. **Use abstract base classes for important interfaces**: Makes expectations clear
3. **Make operator overloading intuitive**: `+` should combine, `-` should subtract, etc.
4. **Implement special methods consistently**: If you override `__eq__`, also override `__hash__`
5. **Document expected interfaces**: Use docstrings to explain what methods objects should have
6. **Follow the principle of least surprise**: Make polymorphic behavior predictable

## Practice Exercises

1. Create a shape hierarchy with polymorphic `area()` and `draw()` methods
2. Implement a custom collection class with operator overloading for `+`, `*`, and `in`
3. Create a notification system with different notification types (email, SMS, push) using polymorphism
4. Implement a calculator that works with different numeric types using operator overloading
5. Create a file processor that can handle different file types (CSV, JSON, XML) polymorphically

---

## Quiz Questions

This section contains quiz questions that test understanding of Python polymorphism. Each question is formatted for easy extraction and database import.

### Question 1: Polymorphism Concept - Basic Understanding
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: python, oop, polymorphism

**Question Text**:
What does polymorphism mean in object-oriented programming?

**Answer Choices**:
- A) Having multiple classes in a program
- B) The ability to use objects of different types through a common interface
- C) Creating multiple objects from one class
- D) Having multiple methods with the same name

**Correct Answer**: B

**Explanation**:
Polymorphism (from Greek "poly" = many, "morph" = forms) means the ability to use objects of different types through a common interface. You can write code that works with a general interface, and it will work with any object that implements that interface. For example, a function that calls `object.speak()` works with any object that has a `speak()` method, regardless of the object's type.

---

### Question 2: Duck Typing - Python Philosophy
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: python, oop, polymorphism, duck-typing

**Question Text**:
What is "duck typing" in Python?

**Answer Choices**:
- A) A type of variable that can store duck objects
- B) If an object has the methods/attributes you need, its actual type doesn't matter
- C) A way to convert objects to duck type
- D) A strict type checking system

**Correct Answer**: B

**Explanation**:
Duck typing is Python's philosophy: "If it walks like a duck and quacks like a duck, it's a duck." You don't check the type of an object—you just try to use it. If an object has the methods and attributes you need, you can use it, regardless of its actual type. This makes Python very flexible and is central to polymorphism in Python.

---

### Question 3: Method Overriding - Replacing Parent Methods
**Type**: code-output
**Difficulty**: 3
**Tags**: python, oop, polymorphism, method-overriding

**Question Text**:
What will be the output?

```python
class Animal:
    def speak(self):
        return "Some sound"

class Dog(Animal):
    def speak(self):
        return "Woof!"

dog = Dog()
print(dog.speak())
```

**Answer Choices**:
- A) Some sound
- B) Woof!
- C) Both "Some sound" and "Woof!"
- D) Error

**Correct Answer**: B

**Explanation**:
The Dog class overrides the `speak()` method from Animal. When you call `dog.speak()`, Python uses Dog's version of the method, not Animal's. This is method overriding—a key aspect of polymorphism where subclasses provide their own implementation of methods defined in parent classes.

---

### Question 4: Operator Overloading - __add__ Method
**Type**: code-output
**Difficulty**: 4
**Tags**: python, oop, polymorphism, operator-overloading

**Question Text**:
What enables this code to work?

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
- A) Python automatically knows how to add any objects
- B) The __add__() method defines behavior for the + operator
- C) Points can always be added in Python
- D) The code will actually cause an error

**Correct Answer**: B

**Explanation**:
The `__add__()` special method (also called a magic or dunder method) defines how the `+` operator works with your objects. When you write `p1 + p2`, Python calls `p1.__add__(p2)`. This is operator overloading—defining custom behavior for operators like `+`, `-`, `*`, etc. Without `__add__()`, trying to add Point objects would raise a TypeError.

---

### Question 5: Polymorphic Functions - Interface Flexibility
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: python, oop, polymorphism, duck-typing

**Question Text**:
What is required for this function to work with an object?

```python
def process(item):
    item.save()
    item.validate()
```

**Answer Choices**:
- A) The object must inherit from a specific base class
- B) The object must have save() and validate() methods
- C) The object must be of type Item
- D) The object must be registered with the process function

**Correct Answer**: B

**Explanation**:
Due to duck typing, the function will work with any object that has `save()` and `validate()` methods. The actual type or inheritance hierarchy doesn't matter. This is polymorphism in action—the function works with any object that implements the expected interface (has the right methods), making Python code flexible and reusable.

---

### Question 6: Special Methods - Polymorphic Behavior
**Type**: code-output
**Difficulty**: 4
**Tags**: python, oop, polymorphism, special-methods

**Question Text**:
What allows `len()` to work with this custom class?

```python
class MyCollection:
    def __init__(self, items):
        self.items = items
    
    def __len__(self):
        return len(self.items)

col = MyCollection([1, 2, 3, 4, 5])
print(len(col))
```

**Answer Choices**:
- A) All classes automatically work with len()
- B) The __len__() method implements the length protocol
- C) MyCollection inherits from list
- D) len() only works with built-in types

**Correct Answer**: B

**Explanation**:
The built-in `len()` function is polymorphic—it works with any object that implements the `__len__()` method. When you call `len(col)`, Python calls `col.__len__()`. This is an example of protocol-based polymorphism: by implementing certain special methods (like `__len__`, `__iter__`, `__getitem__`), your objects can work with built-in functions and operations.

---

### Question 7: Comparison Operators - __eq__ and __lt__
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: python, oop, polymorphism, operator-overloading

**Question Text**:
If you implement `__eq__()` and `__lt__()` in a class, which operations will work?

**Answer Choices**:
- A) Only == and <
- B) ==, !=, <, >, <=, >= (all comparison operators)
- C) Only == and < (Python doesn't infer other operators)
- D) All operators except !=

**Correct Answer**: C

**Explanation**:
Python does NOT automatically infer other comparison operators from `__eq__()` and `__lt__()`. You need to implement each one separately: `__eq__()` for `==`, `__ne__()` for `!=`, `__lt__()` for `<`, `__le__()` for `<=`, `__gt__()` for `>`, and `__ge__()` for `>=`. However, you can use the `@functools.total_ordering` decorator to automatically generate missing comparison methods if you implement `__eq__()` and one other comparison method.

---

### Question 8: Polymorphism Design - Best Practices
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: python, oop, polymorphism, design

**Question Text**:
When should you use `isinstance()` checks instead of relying on duck typing?

**Answer Choices**:
- A) Always - type checking is always safer
- B) Never - duck typing is always better
- C) When you specifically need to handle different types differently or for validation
- D) Only with built-in types

**Correct Answer**: C

**Explanation**:
Generally, you should favor duck typing and polymorphism over `isinstance()` checks. However, there are legitimate cases for `isinstance()`: when you need to handle different types differently (e.g., special handling for specific types), for input validation, or when working with external APIs. Overuse of `isinstance()` often indicates poor design—you might need better polymorphism instead. The Pythonic way is to try using an object and handle exceptions if it doesn't work (EAFP: Easier to Ask Forgiveness than Permission).

---
