# Java Object-Oriented Programming - Inheritance

Inheritance is one of the coolest features of object-oriented programming. It lets you create new classes based on existing ones, inheriting their properties and behaviors while adding or modifying what you need. Think of it like genetics—you inherit traits from your parents, but you're still your own unique person.

## Concept of Inheritance and Code Reuse

Inheritance is about creating a hierarchy of classes where child classes (subclasses) inherit characteristics from parent classes (superclasses). This promotes code reuse and creates natural relationships between classes.

Imagine you're building a program for a zoo. You might have classes for different animals: `Lion`, `Elephant`, `Penguin`. These all share common features (they eat, sleep, make sounds), but each has unique characteristics. Instead of writing the same code three times, you create an `Animal` superclass with the common features, then create specific animal classes that inherit from it.

```java
// Superclass (parent class)
public class Animal {
    protected String name;
    protected int age;
    
    public void eat() {
        System.out.println(name + " is eating.");
    }
    
    public void sleep() {
        System.out.println(name + " is sleeping.");
    }
}

// Subclass (child class)
public class Lion extends Animal {
    public void roar() {
        System.out.println(name + " roars loudly!");
    }
}
```

The `Lion` class automatically has `name`, `age`, `eat()`, and `sleep()` from `Animal`, plus its own unique `roar()` method.

## Extending Classes with the `extends` Keyword

The `extends` keyword creates the parent-child relationship between classes:

```java
public class Vehicle {
    protected String brand;
    protected int year;
    
    public void start() {
        System.out.println("Vehicle is starting...");
    }
    
    public void stop() {
        System.out.println("Vehicle is stopping...");
    }
}

public class Car extends Vehicle {
    private int numberOfDoors;
    
    public Car(String brand, int year, int doors) {
        this.brand = brand;  // Inherited from Vehicle
        this.year = year;    // Inherited from Vehicle
        this.numberOfDoors = doors;
    }
    
    public void honk() {
        System.out.println("Beep beep!");
    }
}
```

Now you can use a `Car` object:

```java
Car myCar = new Car("Toyota", 2023, 4);
myCar.start();  // Inherited method
myCar.honk();   // Car's own method
myCar.stop();   // Inherited method
```

## The `super` Keyword

The `super` keyword lets you access members of the parent class. It's particularly useful in constructors and when you've overridden a method.

### Using `super` in Constructors

When you create a subclass object, the superclass constructor needs to be called first. You use `super()` to do this:

```java
public class Animal {
    protected String name;
    
    public Animal(String name) {
        this.name = name;
        System.out.println("Animal constructor called");
    }
}

public class Dog extends Animal {
    private String breed;
    
    public Dog(String name, String breed) {
        super(name);  // Call parent constructor first
        this.breed = breed;
        System.out.println("Dog constructor called");
    }
}
```

Usage:

```java
Dog myDog = new Dog("Buddy", "Golden Retriever");
// Output:
// Animal constructor called
// Dog constructor called
```

**Important**: If you don't explicitly call `super()`, Java automatically calls the no-argument parent constructor. If the parent doesn't have a no-argument constructor, you'll get a compilation error.

### Using `super` to Access Parent Methods

```java
public class Shape {
    public void draw() {
        System.out.println("Drawing a shape");
    }
}

public class Circle extends Shape {
    public void draw() {
        super.draw();  // Call parent's draw() method first
        System.out.println("Drawing a circle");
    }
}
```

Output when you call `circle.draw()`:
```
Drawing a shape
Drawing a circle
```

### Using `super` to Access Parent Fields

```java
public class Parent {
    protected int value = 10;
}

public class Child extends Parent {
    private int value = 20;  // Hides parent's value
    
    public void display() {
        System.out.println("Child value: " + value);        // 20
        System.out.println("Parent value: " + super.value); // 10
    }
}
```

## Method Overriding

Method overriding is when a subclass provides its own implementation of a method that's already defined in its parent class. The method in the child class must have the same signature (name, parameters, return type) as the parent's method.

```java
public class Animal {
    public void makeSound() {
        System.out.println("Some generic animal sound");
    }
}

public class Dog extends Animal {
    @Override
    public void makeSound() {
        System.out.println("Woof! Woof!");
    }
}

public class Cat extends Animal {
    @Override
    public void makeSound() {
        System.out.println("Meow!");
    }
}
```

Usage:

```java
Dog dog = new Dog();
Cat cat = new Cat();

dog.makeSound();  // Output: Woof! Woof!
cat.makeSound();  // Output: Meow!
```

### Rules for Method Overriding

1. The method must have the same name as in the parent class
2. The method must have the same parameters as in the parent class
3. The return type must be the same or a subtype (covariant return type)
4. The access modifier must be the same or more permissive (can't be more restrictive)
5. Can't override final methods
6. Can't override static methods (that's method hiding, not overriding)

```java
public class Parent {
    protected void display() {  // protected access
        System.out.println("Parent");
    }
}

public class Child extends Parent {
    @Override
    public void display() {  // public is okay (more permissive)
        System.out.println("Child");
    }
    
    // This would be an error - more restrictive:
    // @Override
    // private void display() { }
}
```

## The `@Override` Annotation

The `@Override` annotation is optional but highly recommended. It tells the compiler "I intend to override a parent method." If you make a mistake (wrong parameter types, misspelled name), the compiler will catch it:

```java
public class Animal {
    public void makeSound() {
        System.out.println("Some sound");
    }
}

public class Dog extends Animal {
    @Override
    public void makeSound() {  // Correct override
        System.out.println("Woof!");
    }
    
    // This would cause a compilation error because of @Override:
    // @Override
    // public void makesound() {  // Wrong name - lowercase 's'
    //     System.out.println("Woof!");
    // }
}
```

Without `@Override`, the misspelled method would just create a new method instead of overriding, which could cause hard-to-find bugs.

## Constructor Chaining in Inheritance

When you create an object of a subclass, constructors are called in a chain from the top of the inheritance hierarchy down:

```java
public class Grandparent {
    public Grandparent() {
        System.out.println("Grandparent constructor");
    }
}

public class Parent extends Grandparent {
    public Parent() {
        // super() is called implicitly here
        System.out.println("Parent constructor");
    }
}

public class Child extends Parent {
    public Child() {
        // super() is called implicitly here
        System.out.println("Child constructor");
    }
}
```

Creating a `Child` object:

```java
Child child = new Child();
// Output:
// Grandparent constructor
// Parent constructor
// Child constructor
```

This ensures that the object is fully initialized from the most general (Grandparent) to the most specific (Child).

## The Object Class and Its Methods

In Java, every class automatically extends the `Object` class (unless it explicitly extends another class). This means every object in Java has certain methods inherited from `Object`.

### Common Object Methods

```java
public class Person {
    private String name;
    private int age;
    
    public Person(String name, int age) {
        this.name = name;
        this.age = age;
    }
    
    // Override toString() from Object class
    @Override
    public String toString() {
        return "Person{name='" + name + "', age=" + age + "}";
    }
    
    // Override equals() from Object class
    @Override
    public boolean equals(Object obj) {
        if (this == obj) return true;
        if (obj == null || getClass() != obj.getClass()) return false;
        
        Person person = (Person) obj;
        return age == person.age && name.equals(person.name);
    }
    
    // Override hashCode() from Object class
    @Override
    public int hashCode() {
        int result = name.hashCode();
        result = 31 * result + age;
        return result;
    }
}
```

### `toString()` Method

Returns a string representation of the object. The default implementation in `Object` returns something like `Person@15db9742`, which isn't very helpful. You should override it:

```java
Person p = new Person("Alice", 25);
System.out.println(p);  // Calls toString() automatically
// Output: Person{name='Alice', age=25}
```

### `equals()` Method

Compares two objects for equality. The default implementation just checks if they're the same object in memory (`==`). For meaningful comparison, override it:

```java
Person p1 = new Person("Alice", 25);
Person p2 = new Person("Alice", 25);

System.out.println(p1 == p2);       // false (different objects in memory)
System.out.println(p1.equals(p2));  // true (same name and age)
```

### `hashCode()` Method

Returns an integer hash code for the object. Important rule: if two objects are equal according to `equals()`, they must have the same hash code. This is crucial for hash-based collections like `HashMap`:

```java
Person p1 = new Person("Alice", 25);
Person p2 = new Person("Alice", 25);

System.out.println(p1.hashCode() == p2.hashCode());  // true
```

### `getClass()` Method

Returns the runtime class of the object:

```java
Person p = new Person("Alice", 25);
System.out.println(p.getClass().getName());  // Output: Person
```

## Abstract Classes

Sometimes you want to create a class that serves as a base for other classes but should never be instantiated on its own. That's an abstract class.

```java
public abstract class Shape {
    protected String color;
    
    public Shape(String color) {
        this.color = color;
    }
    
    // Abstract method - no implementation
    public abstract double getArea();
    
    // Concrete method - has implementation
    public void displayColor() {
        System.out.println("Color: " + color);
    }
}
```

You can't do this:
```java
Shape shape = new Shape("Red");  // Compilation error!
```

But you can create concrete subclasses:

```java
public class Circle extends Shape {
    private double radius;
    
    public Circle(String color, double radius) {
        super(color);
        this.radius = radius;
    }
    
    @Override
    public double getArea() {
        return Math.PI * radius * radius;
    }
}

public class Rectangle extends Shape {
    private double width;
    private double height;
    
    public Rectangle(String color, double width, double height) {
        super(color);
        this.width = width;
        this.height = height;
    }
    
    @Override
    public double getArea() {
        return width * height;
    }
}
```

Now you can use them:

```java
Shape circle = new Circle("Red", 5.0);
Shape rectangle = new Rectangle("Blue", 4.0, 6.0);

circle.displayColor();           // Color: Red
System.out.println(circle.getArea());     // ~78.54

rectangle.displayColor();        // Color: Blue
System.out.println(rectangle.getArea());  // 24.0
```

### When to Use Abstract Classes

- You want to share code among several closely related classes
- You expect classes that extend your abstract class to have many common methods or fields
- You want to declare non-static or non-final fields (unlike interfaces)
- You need to provide some default behavior but also require subclasses to implement certain methods

## The `final` Keyword with Inheritance

### Final Methods

A final method cannot be overridden by subclasses:

```java
public class Vehicle {
    public final void startEngine() {
        System.out.println("Engine started with key");
    }
    
    public void drive() {
        System.out.println("Vehicle is driving");
    }
}

public class Car extends Vehicle {
    // This is okay - overriding non-final method
    @Override
    public void drive() {
        System.out.println("Car is driving smoothly");
    }
    
    // This would cause a compilation error:
    // @Override
    // public void startEngine() { }  // Can't override final method
}
```

### Final Classes

A final class cannot be extended at all:

```java
public final class ImmutablePoint {
    private final int x;
    private final int y;
    
    public ImmutablePoint(int x, int y) {
        this.x = x;
        this.y = y;
    }
    
    public int getX() { return x; }
    public int getY() { return y; }
}

// This would cause a compilation error:
// public class ExtendedPoint extends ImmutablePoint { }
```

The `String` class in Java is final, which is why you can't subclass it.

## Single Inheritance in Java

Java only supports single inheritance for classes—a class can extend only one other class:

```java
public class A { }
public class B { }

// This is NOT allowed in Java:
// public class C extends A, B { }  // Compilation error!
```

This is different from multiple inheritance in languages like C++. Java uses interfaces to achieve similar functionality (we'll cover that in the Interfaces topic).

Why single inheritance? It avoids the "diamond problem" where ambiguity arises if two parent classes have methods with the same signature.

## IS-A Relationship

Inheritance represents an "IS-A" relationship: a `Dog` IS-A `Animal`, a `Car` IS-A `Vehicle`.

```java
public class Animal { }
public class Dog extends Animal { }

Dog myDog = new Dog();
Animal myAnimal = myDog;  // This works - a Dog IS-A Animal
```

This is different from a "HAS-A" relationship (composition), where one class contains another:

```java
public class Engine { }

public class Car {
    private Engine engine;  // Car HAS-A Engine
}
```

When deciding between inheritance and composition:
- Use inheritance when there's a clear IS-A relationship
- Use composition when there's a HAS-A relationship
- Favor composition over inheritance when you just need to reuse code without a true IS-A relationship

## Common Pitfalls

1. **Overusing inheritance**: Not everything needs inheritance. Sometimes composition is better.
2. **Breaking the IS-A relationship**: A `Square` extending `Rectangle` can cause issues because a square has constraints a rectangle doesn't.
3. **Forgetting to call `super()`**: If parent has no no-arg constructor, you must call `super()` with arguments.
4. **Making everything protected**: Just because fields can be accessed by subclasses doesn't mean they should be. Prefer private with getters/setters.
5. **Deep inheritance hierarchies**: Going more than 3-4 levels deep makes code hard to understand and maintain.

## Best Practices

1. **Use `@Override` annotation**: Always mark overridden methods to catch mistakes.
2. **Override `toString()`, `equals()`, and `hashCode()` together**: If you override equals, override hashCode too.
3. **Make inheritance intentional**: Use abstract classes when you want to force subclasses to implement certain methods.
4. **Document your class hierarchy**: Use comments to explain the inheritance relationships.
5. **Prefer composition over inheritance**: Unless there's a clear IS-A relationship, composition is often more flexible.
6. **Keep constructors simple**: Don't put complex logic in constructors, especially in inheritance hierarchies.

## Practice Exercises

1. Create an `Employee` superclass and `Manager` and `Developer` subclasses with appropriate methods.
2. Create a `Shape` abstract class with `Circle`, `Rectangle`, and `Triangle` subclasses.
3. Override `toString()`, `equals()`, and `hashCode()` for a `Student` class.
4. Create a three-level inheritance hierarchy (e.g., `LivingThing` → `Animal` → `Dog`) and observe constructor chaining.
5. Create a class that demonstrates the difference between overriding and method hiding (static methods).

---

## Quiz Questions

This section contains quiz questions that test understanding of Java inheritance. Each question is formatted for easy extraction and database import.

### Question 1: Inheritance Basics - extends Keyword
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: java, oop, inheritance, extends

**Question Text**:
What keyword is used to create a subclass in Java?

**Answer Choices**:
- A) implements
- B) extends
- C) inherits
- D) derives

**Correct Answer**: B

**Explanation**:
The `extends` keyword is used to create a subclass (child class) from a superclass (parent class). For example, `class Dog extends Animal` makes Dog a subclass of Animal. The `implements` keyword is used for interfaces, not inheritance. `inherits` and `derives` are not Java keywords.

---

### Question 2: super Keyword - Constructor Chaining
**Type**: code-output
**Difficulty**: 3
**Tags**: java, oop, inheritance, super, constructors

**Question Text**:
What will be the output of this code?

```java
class Parent {
    public Parent() {
        System.out.println("Parent constructor");
    }
}

class Child extends Parent {
    public Child() {
        System.out.println("Child constructor");
    }
}

public class Test {
    public static void main(String[] args) {
        Child c = new Child();
    }
}
```

**Answer Choices**:
- A) Child constructor
- B) Parent constructor
- C) Parent constructor\nChild constructor
- D) Compilation error

**Correct Answer**: C

**Explanation**:
When a Child object is created, the Parent constructor is called first (automatically via implicit `super()` call), followed by the Child constructor. This is called constructor chaining. The output shows both constructors being called in order from parent to child.

---

### Question 3: Method Overriding - Signature Matching
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, oop, inheritance, method-overriding

**Question Text**:
For a method in a subclass to properly override a method in the superclass, what must be true?

**Answer Choices**:
- A) The methods must have different names
- B) The methods must have the same name, parameters, and compatible return type
- C) The subclass method must have a more restrictive access modifier
- D) The superclass method must be abstract

**Correct Answer**: B

**Explanation**:
To override a method, the subclass method must have the same name, same parameters (type, number, and order), and the same or covariant (subtype) return type. The access modifier must be the same or more permissive (not more restrictive). The superclass method doesn't need to be abstract—you can override concrete methods too.

---

### Question 4: @Override Annotation - Error Detection
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, oop, inheritance, override-annotation

**Question Text**:
What is the benefit of using the `@Override` annotation?

**Answer Choices**:
- A) It makes the method run faster
- B) It allows the method to have a different signature than the parent method
- C) It helps the compiler detect mistakes when attempting to override a method
- D) It is required for all overridden methods

**Correct Answer**: C

**Explanation**:
The `@Override` annotation tells the compiler that you intend to override a method from the superclass. If you make a mistake (like misspelling the method name or using wrong parameters), the compiler will give an error. It doesn't affect performance, doesn't change signature requirements, and is optional (though highly recommended).

---

### Question 5: Object Class Methods - toString
**Type**: code-output
**Difficulty**: 4
**Tags**: java, oop, inheritance, object-class, tostring

**Question Text**:
What will be printed if the `toString()` method is NOT overridden in the Person class?

```java
public class Person {
    private String name;
    
    public Person(String name) {
        this.name = name;
    }
}

public class Test {
    public static void main(String[] args) {
        Person p = new Person("Alice");
        System.out.println(p);
    }
}
```

**Answer Choices**:
- A) Alice
- B) Person{name='Alice'}
- C) Something like Person@15db9742 (class name @ hash code)
- D) null

**Correct Answer**: C

**Explanation**:
When you don't override `toString()`, the default implementation from the Object class is used, which returns the class name followed by @ and the hexadecimal hash code of the object. To get a meaningful string representation, you should override `toString()` to return something like "Person{name='Alice'}".

---

### Question 6: Abstract Classes - Instantiation
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, oop, inheritance, abstract-classes

**Question Text**:
Can you create an instance of an abstract class?

**Answer Choices**:
- A) Yes, using the new keyword
- B) Yes, but only if it has no abstract methods
- C) No, abstract classes cannot be instantiated
- D) Yes, but only in the same package

**Correct Answer**: C

**Explanation**:
Abstract classes cannot be instantiated directly using the `new` keyword, even if they have no abstract methods. They are meant to be subclassed, and you can create instances of concrete (non-abstract) subclasses. The purpose of an abstract class is to serve as a template for other classes, not to be used directly.

---

### Question 7: final Keyword - Method Overriding
**Type**: code-output
**Difficulty**: 4
**Tags**: java, oop, inheritance, final

**Question Text**:
What happens with this code?

```java
class Parent {
    public final void display() {
        System.out.println("Parent");
    }
}

class Child extends Parent {
    @Override
    public void display() {
        System.out.println("Child");
    }
}
```

**Answer Choices**:
- A) Prints "Parent"
- B) Prints "Child"
- C) Compilation error - cannot override final method
- D) Runtime exception

**Correct Answer**: C

**Explanation**:
Methods marked as `final` cannot be overridden by subclasses. This code will produce a compilation error because the Child class attempts to override the final `display()` method from Parent. The `final` keyword is used when you want to prevent subclasses from changing the method's behavior.

---

### Question 8: IS-A Relationship - Inheritance vs Composition
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, oop, inheritance, is-a-relationship

**Question Text**:
Which scenario is a good candidate for using inheritance (IS-A relationship)?

**Answer Choices**:
- A) A Car class containing an Engine object
- B) A Student class extending Person class
- C) A House class containing a Door object
- D) A Library class containing Book objects

**Correct Answer**: B

**Explanation**:
A Student IS-A Person, which is a true inheritance relationship. A student has all the characteristics of a person (name, age, etc.) plus additional student-specific attributes (student ID, GPA, etc.). The other options represent HAS-A relationships (composition): a car HAS an engine, a house HAS doors, a library HAS books. These should use composition (instance variables) rather than inheritance.

---
