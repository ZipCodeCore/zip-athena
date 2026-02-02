# Python Object-Oriented Programming - Inheritance

Inheritance is a powerful feature of object-oriented programming that lets you create new classes based on existing ones. The new class (child or subclass) inherits attributes and methods from the existing class (parent or superclass), and can add its own or modify inherited ones. It's like how you might inherit traits from your parents but also have your own unique characteristics.

## Concept of Inheritance

Inheritance promotes code reuse and creates natural hierarchies. Instead of copying code between similar classes, you create a parent class with shared functionality and then create child classes that inherit and extend it.

Think of a zoo management system. You might have `Lion`, `Elephant`, and `Penguin` classes. They all share common animal behaviors (eating, sleeping, making sounds), but each has unique characteristics. Instead of duplicating the common code, create an `Animal` parent class:

```python
class Animal:
    def __init__(self, name, age):
        self.name = name
        self.age = age
    
    def eat(self):
        print(f"{self.name} is eating")
    
    def sleep(self):
        print(f"{self.name} is sleeping")

class Lion(Animal):  # Lion inherits from Animal
    def roar(self):
        print(f"{self.name} roars loudly!")

class Elephant(Animal):
    def trumpet(self):
        print(f"{self.name} trumpets!")

# Using inheritance
lion = Lion("Simba", 5)
lion.eat()    # Inherited from Animal
lion.roar()   # Lion's own method

elephant = Elephant("Dumbo", 3)
elephant.eat()     # Inherited from Animal
elephant.trumpet() # Elephant's own method
```

## Single Inheritance

Single inheritance is when a class inherits from one parent class. It's the most common form of inheritance.

### Basic Single Inheritance

```python
class Vehicle:
    def __init__(self, brand, year):
        self.brand = brand
        self.year = year
    
    def start(self):
        print(f"{self.brand} vehicle is starting")
    
    def stop(self):
        print(f"{self.brand} vehicle is stopping")

class Car(Vehicle):  # Car inherits from Vehicle
    def __init__(self, brand, year, num_doors):
        super().__init__(brand, year)  # Call parent constructor
        self.num_doors = num_doors
    
    def honk(self):
        print("Beep beep!")

# Using the Car class
my_car = Car("Toyota", 2023, 4)
my_car.start()  # Inherited method
my_car.honk()   # Car's own method
print(f"My car is a {my_car.year} {my_car.brand}")
```

### The `super()` Function

The `super()` function gives you access to methods and attributes of the parent class. It's especially useful in `__init__()` to call the parent's constructor:

```python
class Person:
    def __init__(self, name, age):
        self.name = name
        self.age = age
    
    def introduce(self):
        print(f"Hi, I'm {self.name} and I'm {self.age} years old")

class Student(Person):
    def __init__(self, name, age, student_id):
        super().__init__(name, age)  # Initialize parent attributes
        self.student_id = student_id  # Add student-specific attribute
    
    def study(self):
        print(f"{self.name} is studying")

student = Student("Alice", 20, "S12345")
student.introduce()  # Inherited from Person
student.study()      # Student's own method
```

### Method Overriding

Child classes can override parent methods to provide their own implementation:

```python
class Animal:
    def __init__(self, name):
        self.name = name
    
    def make_sound(self):
        print(f"{self.name} makes a sound")

class Dog(Animal):
    def make_sound(self):  # Override parent method
        print(f"{self.name} says Woof!")

class Cat(Animal):
    def make_sound(self):  # Override parent method
        print(f"{self.name} says Meow!")

dog = Dog("Buddy")
cat = Cat("Whiskers")

dog.make_sound()  # Buddy says Woof!
cat.make_sound()  # Whiskers says Meow!
```

You can also call the parent's method from within the override using `super()`:

```python
class Animal:
    def move(self):
        print("Animal is moving")

class Bird(Animal):
    def move(self):
        super().move()  # Call parent's method first
        print("Bird is flying")

bird = Bird()
bird.move()
# Output:
# Animal is moving
# Bird is flying
```

## Multiple Inheritance

Python supports multiple inheritance—a class can inherit from multiple parent classes:

```python
class Flyer:
    def fly(self):
        print("Flying in the sky")

class Swimmer:
    def swim(self):
        print("Swimming in water")

class Duck(Flyer, Swimmer):  # Inherits from both
    def quack(self):
        print("Quack quack!")

duck = Duck()
duck.fly()    # From Flyer
duck.swim()   # From Swimmer
duck.quack()  # Duck's own method
```

### Method Resolution Order (MRO)

When a class inherits from multiple parents, Python uses the Method Resolution Order (MRO) to determine which method to call:

```python
class A:
    def greet(self):
        print("Hello from A")

class B:
    def greet(self):
        print("Hello from B")

class C(A, B):  # Inherits from both A and B
    pass

c = C()
c.greet()  # Which greet() is called?
# Output: Hello from A (A comes first in the inheritance list)

# Check the MRO
print(C.__mro__)
# Output: (<class 'C'>, <class 'A'>, <class 'B'>, <class 'object'>)
```

The MRO follows a depth-first, left-to-right search: C → A → B → object

### The Diamond Problem

The diamond problem occurs when a class inherits from two classes that both inherit from the same base class:

```python
class Animal:
    def __init__(self, name):
        self.name = name
        print(f"Animal.__init__ called for {name}")

class Mammal(Animal):
    def __init__(self, name):
        super().__init__(name)
        print("Mammal.__init__ called")

class Bird(Animal):
    def __init__(self, name):
        super().__init__(name)
        print("Bird.__init__ called")

class Bat(Mammal, Bird):
    def __init__(self, name):
        super().__init__(name)
        print("Bat.__init__ called")

bat = Bat("Dracula")
# Output:
# Animal.__init__ called for Dracula
# Bird.__init__ called
# Mammal.__init__ called
# Bat.__init__ called

print(Bat.__mro__)
# (<class 'Bat'>, <class 'Mammal'>, <class 'Bird'>, <class 'Animal'>, <class 'object'>)
```

Python's MRO (using the C3 linearization algorithm) ensures each class is only called once and in a predictable order.

## The `isinstance()` and `issubclass()` Functions

### `isinstance()`

Checks if an object is an instance of a class or its subclasses:

```python
class Animal:
    pass

class Dog(Animal):
    pass

dog = Dog()

print(isinstance(dog, Dog))     # True
print(isinstance(dog, Animal))  # True (Dog is a subclass of Animal)
print(isinstance(dog, object))  # True (everything inherits from object)
print(isinstance(dog, str))     # False
```

### `issubclass()`

Checks if a class is a subclass of another:

```python
class Animal:
    pass

class Dog(Animal):
    pass

class Cat(Animal):
    pass

print(issubclass(Dog, Animal))   # True
print(issubclass(Cat, Animal))   # True
print(issubclass(Dog, Cat))      # False
print(issubclass(Animal, object)) # True
```

## Composition vs Inheritance

Sometimes you don't need inheritance—composition (having an object as an attribute) might be better.

### When to Use Inheritance (IS-A Relationship)

Use inheritance when there's a true "IS-A" relationship:

```python
class Vehicle:
    pass

class Car(Vehicle):  # A Car IS-A Vehicle
    pass
```

### When to Use Composition (HAS-A Relationship)

Use composition when there's a "HAS-A" relationship:

```python
class Engine:
    def start(self):
        print("Engine starting")

class Car:
    def __init__(self):
        self.engine = Engine()  # Car HAS-A Engine
    
    def start(self):
        self.engine.start()
        print("Car is ready")

car = Car()
car.start()
```

**General Rule**: Favor composition over inheritance when you just need functionality, not a type relationship.

### Example: Mixing Both

```python
class Engine:
    def __init__(self, horsepower):
        self.horsepower = horsepower
    
    def start(self):
        print(f"{self.horsepower}HP engine starting")

class Vehicle:
    def __init__(self, brand):
        self.brand = brand
    
    def info(self):
        print(f"Vehicle: {self.brand}")

class Car(Vehicle):  # Inheritance (IS-A relationship)
    def __init__(self, brand, horsepower):
        super().__init__(brand)
        self.engine = Engine(horsepower)  # Composition (HAS-A relationship)
    
    def start(self):
        self.engine.start()
        print(f"{self.brand} car is ready to drive")

car = Car("Toyota", 200)
car.info()
car.start()
```

## Abstract Base Classes

Abstract base classes (ABCs) define interfaces that subclasses must implement. Use the `abc` module:

```python
from abc import ABC, abstractmethod

class Shape(ABC):  # Abstract base class
    @abstractmethod
    def area(self):
        pass  # Subclasses MUST implement this
    
    @abstractmethod
    def perimeter(self):
        pass  # Subclasses MUST implement this
    
    # Concrete method (optional for subclasses to override)
    def describe(self):
        print(f"I am a shape with area {self.area()}")

class Circle(Shape):
    def __init__(self, radius):
        self.radius = radius
    
    def area(self):
        return 3.14159 * self.radius ** 2
    
    def perimeter(self):
        return 2 * 3.14159 * self.radius

class Rectangle(Shape):
    def __init__(self, width, height):
        self.width = width
        self.height = height
    
    def area(self):
        return self.width * self.height
    
    def perimeter(self):
        return 2 * (self.width + self.height)

# Can't instantiate abstract class:
# shape = Shape()  # TypeError: Can't instantiate abstract class

# But can instantiate concrete subclasses:
circle = Circle(5)
print(f"Circle area: {circle.area()}")
circle.describe()

rectangle = Rectangle(4, 6)
print(f"Rectangle area: {rectangle.area()}")
rectangle.describe()
```

### Why Use Abstract Base Classes?

1. **Enforce interface**: Ensure subclasses implement required methods
2. **Document intent**: Make it clear what subclasses should do
3. **Catch errors early**: Get an error if a subclass forgets to implement required methods

```python
from abc import ABC, abstractmethod

class PaymentProcessor(ABC):
    @abstractmethod
    def process_payment(self, amount):
        pass
    
    @abstractmethod
    def refund(self, transaction_id):
        pass

class CreditCardProcessor(PaymentProcessor):
    def process_payment(self, amount):
        print(f"Processing ${amount} credit card payment")
    
    def refund(self, transaction_id):
        print(f"Refunding transaction {transaction_id}")

class PayPalProcessor(PaymentProcessor):
    def process_payment(self, amount):
        print(f"Processing ${amount} PayPal payment")
    
    def refund(self, transaction_id):
        print(f"Refunding PayPal transaction {transaction_id}")

# Both classes implement the required interface
cc = CreditCardProcessor()
cc.process_payment(99.99)

paypal = PayPalProcessor()
paypal.process_payment(49.99)
```

## Common Pitfalls

1. **Forgetting to call `super().__init__()`**: Parent's initialization code won't run
2. **Deep inheritance hierarchies**: More than 3-4 levels gets confusing
3. **Overusing inheritance**: Sometimes composition is better
4. **Multiple inheritance confusion**: Can lead to complex MRO issues
5. **Modifying parent class attributes incorrectly**: Be careful with mutable attributes

## Best Practices

1. **Call `super().__init__()` in child classes**: Ensures proper initialization
2. **Keep inheritance hierarchies shallow**: Easier to understand and maintain
3. **Use abstract base classes for interfaces**: Makes expectations clear
4. **Favor composition over inheritance**: Unless there's a true IS-A relationship
5. **Document your class hierarchy**: Use docstrings to explain relationships
6. **Check the MRO when using multiple inheritance**: Use `ClassName.__mro__`

## Practice Exercises

1. Create an `Employee` parent class with `Manager` and `Developer` child classes
2. Create a `Shape` abstract base class with `Circle`, `Rectangle`, and `Triangle` concrete classes
3. Create a class using multiple inheritance (e.g., `FlyingFish` from `Fish` and `Flyer`)
4. Implement a three-level inheritance hierarchy and trace how `super()` works
5. Create a class that uses both inheritance and composition (e.g., `Car` inherits from `Vehicle` and has an `Engine`)

---

## Quiz Questions

This section contains quiz questions that test understanding of Python inheritance. Each question is formatted for easy extraction and database import.

### Question 1: Inheritance Basics - Creating Subclasses
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: python, oop, inheritance

**Question Text**:
How do you create a class that inherits from another class in Python?

**Answer Choices**:
- A) class ChildClass inherits ParentClass:
- B) class ChildClass(ParentClass):
- C) class ChildClass extends ParentClass:
- D) class ChildClass <- ParentClass:

**Correct Answer**: B

**Explanation**:
In Python, you create a subclass by putting the parent class name in parentheses after the child class name: `class ChildClass(ParentClass):`. This is different from languages like Java which use the `extends` keyword. Options A, C, and D use incorrect syntax.

---

### Question 2: super() Function - Parent Access
**Type**: code-output
**Difficulty**: 3
**Tags**: python, oop, inheritance, super

**Question Text**:
What will be the output?

```python
class Parent:
    def __init__(self, name):
        self.name = name
        print(f"Parent: {name}")

class Child(Parent):
    def __init__(self, name, age):
        super().__init__(name)
        self.age = age
        print(f"Child: {age}")

c = Child("Alice", 10)
```

**Answer Choices**:
- A) Child: 10
- B) Parent: Alice
- C) Parent: Alice\nChild: 10
- D) Error

**Correct Answer**: C

**Explanation**:
`super().__init__(name)` calls the parent class's `__init__()` method first, which prints "Parent: Alice". Then the Child's `__init__()` continues and prints "Child: 10". The super() function is used to call methods from the parent class, and it's especially important in constructors to properly initialize inherited attributes.

---

### Question 3: Method Overriding - Replacing Parent Methods
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: python, oop, inheritance, method-overriding

**Question Text**:
What happens when a child class defines a method with the same name as a method in the parent class?

**Answer Choices**:
- A) Error - method names must be unique
- B) The child's method overrides the parent's method
- C) Both methods are called automatically
- D) The parent's method is always used

**Correct Answer**: B

**Explanation**:
When a child class defines a method with the same name as a parent class method, it overrides (replaces) the parent's method. When you call that method on a child class instance, the child's version is executed, not the parent's. You can still call the parent's version using `super().method_name()` if needed.

---

### Question 4: Multiple Inheritance - Multiple Parents
**Type**: code-output
**Difficulty**: 4
**Tags**: python, oop, inheritance, multiple-inheritance

**Question Text**:
What does this code demonstrate?

```python
class Flyer:
    def move(self):
        print("Flying")

class Swimmer:
    def move(self):
        print("Swimming")

class Duck(Flyer, Swimmer):
    pass

duck = Duck()
duck.move()
```

**Answer Choices**:
- A) Error - can't inherit from multiple classes
- B) Prints "Swimming"
- C) Prints "Flying"
- D) Prints both "Flying" and "Swimming"

**Correct Answer**: C

**Explanation**:
This demonstrates multiple inheritance and the Method Resolution Order (MRO). When Duck inherits from both Flyer and Swimmer, and both have a `move()` method, Python uses the MRO to decide which to call. Flyer comes first in the inheritance list `(Flyer, Swimmer)`, so its `move()` method is called, printing "Flying". You can check the MRO with `Duck.__mro__`.

---

### Question 5: isinstance() Function - Type Checking
**Type**: code-output
**Difficulty**: 3
**Tags**: python, oop, inheritance, isinstance

**Question Text**:
What will be printed?

```python
class Animal:
    pass

class Dog(Animal):
    pass

dog = Dog()
print(isinstance(dog, Animal))
```

**Answer Choices**:
- A) False
- B) True
- C) Error
- D) Dog

**Correct Answer**: B

**Explanation**:
`isinstance()` checks if an object is an instance of a class or any of its parent classes. Since Dog inherits from Animal, a Dog instance IS-A Animal, so `isinstance(dog, Animal)` returns True. This works up the entire inheritance chain—even `isinstance(dog, object)` would return True since all classes inherit from object.

---

### Question 6: Abstract Base Classes - Interface Definition
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: python, oop, inheritance, abstract-classes

**Question Text**:
What is the purpose of abstract base classes in Python?

```python
from abc import ABC, abstractmethod

class Shape(ABC):
    @abstractmethod
    def area(self):
        pass
```

**Answer Choices**:
- A) To make classes run faster
- B) To define interfaces that subclasses must implement
- C) To prevent any class from being created
- D) To make all methods private

**Correct Answer**: B

**Explanation**:
Abstract base classes (ABCs) define interfaces that subclasses must implement. The `@abstractmethod` decorator marks methods that subclasses MUST implement. You cannot instantiate an abstract class—you can only create instances of concrete subclasses that implement all abstract methods. This ensures all subclasses follow the same interface, catching errors early if required methods are missing.

---

### Question 7: Composition vs Inheritance - Design Choice
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: python, oop, inheritance, composition

**Question Text**:
When should you use composition instead of inheritance?

**Answer Choices**:
- A) When there's an IS-A relationship
- B) When there's a HAS-A relationship or you just need functionality
- C) Never - inheritance is always better
- D) Only when using multiple inheritance

**Correct Answer**: B

**Explanation**:
Use composition (having an object as an attribute) when there's a HAS-A relationship or when you just need functionality without a true type relationship. For example, a Car HAS-A Engine (composition), but a Car IS-A Vehicle (inheritance). Composition is often more flexible than inheritance and is generally preferred unless there's a clear IS-A relationship. The principle "favor composition over inheritance" is a common design guideline.

---

### Question 8: Method Resolution Order (MRO) - Multiple Inheritance
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: python, oop, inheritance, mro

**Question Text**:
How can you view the Method Resolution Order (MRO) of a class?

**Answer Choices**:
- A) ClassName.mro
- B) ClassName.__mro__
- C) super().__mro__
- D) MRO is not accessible in Python

**Correct Answer**: B

**Explanation**:
You can view a class's Method Resolution Order using `ClassName.__mro__` or `ClassName.mro()`. The MRO determines the order in which Python searches for methods in the inheritance hierarchy, especially important with multiple inheritance. It uses the C3 linearization algorithm to ensure a consistent and predictable order. Understanding MRO is crucial when working with complex inheritance hierarchies.

---
