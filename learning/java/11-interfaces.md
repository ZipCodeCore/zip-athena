# Java Interfaces

Interfaces are one of the most important concepts in Java. Think of them as contracts or promises—when a class implements an interface, it's promising to provide certain behaviors. Unlike classes, interfaces define *what* needs to be done, not *how* to do it. They're like job descriptions: they specify what methods a class must have, but each class can implement those methods in its own way.

## Defining Interfaces

An interface is defined using the `interface` keyword. It contains method signatures (no implementations) that classes must provide:

```java
public interface Drawable {
    void draw();  // No method body - just the signature
    void erase();
    void resize(int width, int height);
}
```

Any class that implements `Drawable` must provide implementations for all three methods.

Here's a more complete example:

```java
public interface Animal {
    void eat();
    void sleep();
    void makeSound();
}
```

### Interface Characteristics

- All methods in an interface are implicitly `public` and `abstract` (before Java 8)
- Interfaces cannot have instance variables (only constants)
- Interfaces cannot be instantiated directly
- A class can implement multiple interfaces
- An interface can extend other interfaces

```java
public interface Vehicle {
    // These are implicitly public static final
    int MAX_SPEED = 200;
    String DEFAULT_COLOR = "White";
    
    // These are implicitly public abstract
    void start();
    void stop();
    void accelerate(int speed);
}
```

## Implementing Interfaces with `implements` Keyword

A class implements an interface using the `implements` keyword and must provide implementations for all methods:

```java
public class Car implements Vehicle {
    private int currentSpeed;
    
    @Override
    public void start() {
        System.out.println("Car engine starting...");
        currentSpeed = 0;
    }
    
    @Override
    public void stop() {
        System.out.println("Car stopping...");
        currentSpeed = 0;
    }
    
    @Override
    public void accelerate(int speed) {
        currentSpeed += speed;
        if (currentSpeed > MAX_SPEED) {
            currentSpeed = MAX_SPEED;
        }
        System.out.println("Car accelerating to " + currentSpeed + " km/h");
    }
}
```

Now you can use it:

```java
Vehicle myCar = new Car();  // Polymorphism!
myCar.start();              // Output: Car engine starting...
myCar.accelerate(50);       // Output: Car accelerating to 50 km/h
myCar.stop();               // Output: Car stopping...
```

### Another Example: Multiple Classes Implementing Same Interface

```java
public interface Playable {
    void play();
    void pause();
    void stop();
}

public class MusicPlayer implements Playable {
    @Override
    public void play() {
        System.out.println("Playing music...");
    }
    
    @Override
    public void pause() {
        System.out.println("Music paused");
    }
    
    @Override
    public void stop() {
        System.out.println("Music stopped");
    }
}

public class VideoPlayer implements Playable {
    @Override
    public void play() {
        System.out.println("Playing video...");
    }
    
    @Override
    public void pause() {
        System.out.println("Video paused");
    }
    
    @Override
    public void stop() {
        System.out.println("Video stopped");
    }
}
```

Both classes implement the same interface but in different ways:

```java
Playable music = new MusicPlayer();
Playable video = new VideoPlayer();

music.play();  // Output: Playing music...
video.play();  // Output: Playing video...
```

## Multiple Interface Implementation

Unlike classes (where Java only allows single inheritance), a class can implement multiple interfaces:

```java
public interface Flyable {
    void fly();
}

public interface Swimmable {
    void swim();
}

public class Duck implements Flyable, Swimmable {
    @Override
    public void fly() {
        System.out.println("Duck flying in the sky");
    }
    
    @Override
    public void swim() {
        System.out.println("Duck swimming in the pond");
    }
}
```

Usage:

```java
Duck duck = new Duck();
duck.fly();   // Output: Duck flying in the sky
duck.swim();  // Output: Duck swimming in the pond

// Can use polymorphism with either interface
Flyable flyingBird = new Duck();
flyingBird.fly();

Swimmable swimmer = new Duck();
swimmer.swim();
```

This is incredibly powerful—a Duck IS-A Flyable AND IS-A Swimmable!

### Real-World Example: Smart Device

```java
public interface Connectable {
    void connect();
    void disconnect();
}

public interface Controllable {
    void turnOn();
    void turnOff();
    void setSettings(String settings);
}

public class SmartTV implements Connectable, Controllable {
    private boolean isOn;
    private boolean isConnected;
    
    @Override
    public void connect() {
        isConnected = true;
        System.out.println("SmartTV connected to WiFi");
    }
    
    @Override
    public void disconnect() {
        isConnected = false;
        System.out.println("SmartTV disconnected from WiFi");
    }
    
    @Override
    public void turnOn() {
        isOn = true;
        System.out.println("SmartTV turned on");
    }
    
    @Override
    public void turnOff() {
        isOn = false;
        System.out.println("SmartTV turned off");
    }
    
    @Override
    public void setSettings(String settings) {
        System.out.println("SmartTV settings updated: " + settings);
    }
}
```

## Interface Inheritance (Extending Interfaces)

Interfaces can extend other interfaces using the `extends` keyword. Unlike classes, an interface can extend multiple interfaces:

```java
public interface Animal {
    void eat();
    void sleep();
}

public interface WildAnimal extends Animal {
    void hunt();
    void roam();
}

public interface Pet extends Animal {
    void play();
    void obey();
}

// An interface extending multiple interfaces
public interface ServiceAnimal extends Pet {
    void performDuty();
}
```

A class implementing a child interface must implement all methods from the parent interfaces too:

```java
public class Dog implements Pet {
    @Override
    public void eat() {
        System.out.println("Dog eating food");
    }
    
    @Override
    public void sleep() {
        System.out.println("Dog sleeping");
    }
    
    @Override
    public void play() {
        System.out.println("Dog playing fetch");
    }
    
    @Override
    public void obey() {
        System.out.println("Dog obeying command");
    }
}
```

## Default Methods in Interfaces (Java 8+)

Before Java 8, all interface methods were abstract. Java 8 introduced default methods—methods with implementations in interfaces:

```java
public interface Vehicle {
    // Abstract methods
    void start();
    void stop();
    
    // Default method with implementation
    default void honk() {
        System.out.println("Beep beep!");
    }
    
    default void displayInfo() {
        System.out.println("This is a vehicle");
    }
}
```

Classes implementing this interface get the default methods for free but can override them if needed:

```java
public class Car implements Vehicle {
    @Override
    public void start() {
        System.out.println("Car starting");
    }
    
    @Override
    public void stop() {
        System.out.println("Car stopping");
    }
    
    // Using the default honk() method - no need to implement it
    
    // But we can override it if we want
    @Override
    public void honk() {
        System.out.println("Car horn: HONK HONK!");
    }
}
```

### Why Default Methods?

Default methods allow you to add new methods to interfaces without breaking existing implementations. Before Java 8, adding a method to an interface would break all classes implementing it. With default methods, you can add new functionality while maintaining backward compatibility.

```java
public interface List {
    // Existing methods
    void add(Object item);
    void remove(Object item);
    
    // New method added in Java 8 with default implementation
    default void sort() {
        // Default sorting implementation
        System.out.println("Sorting with default algorithm");
    }
}
```

Old classes implementing `List` still work because `sort()` has a default implementation.

## Static Methods in Interfaces (Java 8+)

Interfaces can also have static methods:

```java
public interface MathOperations {
    // Abstract method
    int calculate(int a, int b);
    
    // Static method
    static int add(int a, int b) {
        return a + b;
    }
    
    static int multiply(int a, int b) {
        return a * b;
    }
    
    static double divide(int a, int b) {
        if (b == 0) {
            throw new ArithmeticException("Cannot divide by zero");
        }
        return (double) a / b;
    }
}
```

Static methods are called on the interface itself, not on objects:

```java
int sum = MathOperations.add(5, 3);         // Output: 8
int product = MathOperations.multiply(4, 6); // Output: 24
double quotient = MathOperations.divide(10, 2); // Output: 5.0
```

Static methods in interfaces are useful for utility functions related to the interface.

## Functional Interfaces (Java 8+)

A functional interface is an interface with exactly one abstract method. They're used extensively with lambda expressions:

```java
@FunctionalInterface
public interface Calculator {
    int calculate(int a, int b);  // Single abstract method
    
    // Can have default and static methods
    default void printResult(int result) {
        System.out.println("Result: " + result);
    }
}
```

The `@FunctionalInterface` annotation is optional but recommended—it makes the compiler ensure the interface has exactly one abstract method.

Using with lambda expressions:

```java
Calculator addition = (a, b) -> a + b;
Calculator subtraction = (a, b) -> a - b;
Calculator multiplication = (a, b) -> a * b;

System.out.println(addition.calculate(5, 3));      // Output: 8
System.out.println(subtraction.calculate(10, 4));  // Output: 6
System.out.println(multiplication.calculate(3, 7)); // Output: 21
```

### Common Functional Interfaces

Java provides many built-in functional interfaces in the `java.util.function` package:

```java
import java.util.function.*;

// Predicate<T> - takes one argument, returns boolean
Predicate<Integer> isEven = n -> n % 2 == 0;
System.out.println(isEven.test(4));  // true

// Function<T, R> - takes one argument, returns a result
Function<String, Integer> getLength = s -> s.length();
System.out.println(getLength.apply("Hello"));  // 5

// Consumer<T> - takes one argument, returns nothing
Consumer<String> printer = s -> System.out.println(s);
printer.accept("Hello World");  // Output: Hello World

// Supplier<T> - takes no arguments, returns a result
Supplier<Double> randomValue = () -> Math.random();
System.out.println(randomValue.get());  // Random number
```

## Marker Interfaces

A marker interface is an empty interface with no methods. It's used to "mark" or "tag" classes:

```java
public interface Serializable {
    // Empty - just marks the class
}
```

Classes implement marker interfaces to indicate they have certain properties:

```java
public class Student implements Serializable {
    private String name;
    private int id;
    
    // No methods to implement - it's just a marker
}
```

The JVM or framework can check if a class implements the marker interface and treat it specially. For example, `Serializable` tells Java that objects of this class can be serialized (converted to bytes for storage or transmission).

Other marker interfaces in Java:
- `Cloneable` - indicates objects can be cloned
- `Remote` - indicates objects can be used in remote method invocation

Modern Java tends to use annotations instead of marker interfaces, but you'll still see them in older code.

## HAS-A vs. IS-A Relationship

Understanding when to use interfaces versus composition is crucial.

### IS-A Relationship (Inheritance/Interfaces)

Use inheritance or interfaces when there's a true "IS-A" relationship:

```java
public interface Drivable {
    void drive();
}

public class Car implements Drivable {
    // A Car IS-A Drivable thing
    @Override
    public void drive() {
        System.out.println("Driving a car");
    }
}
```

### HAS-A Relationship (Composition)

Use composition when there's a "HAS-A" relationship:

```java
public class Engine {
    public void start() {
        System.out.println("Engine starting");
    }
}

public class Car {
    private Engine engine;  // Car HAS-A Engine
    
    public Car() {
        engine = new Engine();
    }
    
    public void start() {
        engine.start();
    }
}
```

### Combining Both

Often you'll use both:

```java
public interface Vehicle {
    void start();
    void stop();
}

public class Engine {
    public void ignite() {
        System.out.println("Engine igniting");
    }
}

public class Car implements Vehicle {  // IS-A Vehicle
    private Engine engine;  // HAS-A Engine
    
    public Car() {
        engine = new Engine();
    }
    
    @Override
    public void start() {
        engine.ignite();
        System.out.println("Car started");
    }
    
    @Override
    public void stop() {
        System.out.println("Car stopped");
    }
}
```

## Interfaces vs. Abstract Classes

When should you use an interface vs. an abstract class?

### Use an Interface When:
- You want to define a contract for unrelated classes
- You need multiple inheritance (a class can implement multiple interfaces)
- You want to specify behavior without any implementation (except default methods)
- You're designing for future flexibility

### Use an Abstract Class When:
- You have shared code you want subclasses to inherit
- You need non-static, non-final fields
- You want to provide some default behavior with implementation
- Your classes are closely related

### Comparison Example

```java
// Interface - defines "what" a shape can do
public interface Shape {
    double getArea();
    double getPerimeter();
}

// Abstract class - provides shared code for all shapes
public abstract class AbstractShape implements Shape {
    protected String color;
    
    public AbstractShape(String color) {
        this.color = color;
    }
    
    public String getColor() {
        return color;
    }
    
    // Concrete method available to all shapes
    public void display() {
        System.out.println(color + " shape with area " + getArea());
    }
    
    // Abstract methods from Shape interface
    public abstract double getArea();
    public abstract double getPerimeter();
}

// Concrete class
public class Circle extends AbstractShape {
    private double radius;
    
    public Circle(String color, double radius) {
        super(color);
        this.radius = radius;
    }
    
    @Override
    public double getArea() {
        return Math.PI * radius * radius;
    }
    
    @Override
    public double getPerimeter() {
        return 2 * Math.PI * radius;
    }
}
```

## Common Pitfalls

1. **Forgetting to implement all methods**: If you implement an interface, you must implement ALL its abstract methods
2. **Confusing interfaces with classes**: You can't instantiate interfaces directly
3. **Overusing marker interfaces**: Modern Java prefers annotations
4. **Not using @Override**: Makes it harder to catch mistakes
5. **Creating too many interfaces**: Only create interfaces when you have multiple implementations

## Best Practices

1. **Name interfaces by capability**: `Readable`, `Writable`, `Drawable`, `Comparable`
2. **Keep interfaces focused**: Each interface should represent one capability (Single Responsibility Principle)
3. **Use interfaces for polymorphism**: Write methods that accept interface types
4. **Prefer interfaces to abstract classes**: Interfaces are more flexible
5. **Use @FunctionalInterface**: For single-method interfaces that work with lambdas
6. **Document interface contracts**: Explain what implementing classes must do

## Practice Exercises

1. Create a `Sortable` interface and implement it in `Student` and `Product` classes
2. Create a `Payable` interface with `Employee` and `Contractor` implementations
3. Create a hierarchy: `Animal` interface, `Mammal` and `Bird` interfaces extending it, concrete classes implementing them
4. Create a functional interface for a custom operation and use it with lambda expressions
5. Design an interface for a remote control that can control TV, AC, and lights

---

## Quiz Questions

This section contains quiz questions that test understanding of Java interfaces. Each question is formatted for easy extraction and database import.

### Question 1: Interface Basics - Definition
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: java, oop, interfaces

**Question Text**:
What is an interface in Java?

**Answer Choices**:
- A) A class that cannot be instantiated
- B) A contract that specifies what methods a class must implement
- C) A type of variable
- D) A method that returns multiple values

**Correct Answer**: B

**Explanation**:
An interface is a contract or specification that defines what methods a class must implement, but not how to implement them. It specifies the "what" (method signatures) but not the "how" (implementation). Classes that implement an interface promise to provide implementations for all its methods. Option A describes an abstract class, not specifically an interface.

---

### Question 2: Implementing Interfaces - implements Keyword
**Type**: code-output
**Difficulty**: 2
**Tags**: java, oop, interfaces, implements

**Question Text**:
What keyword is used to implement an interface in Java?

**Answer Choices**:
- A) extends
- B) implements
- C) inherits
- D) uses

**Correct Answer**: B

**Explanation**:
The `implements` keyword is used when a class implements an interface. For example: `class Dog implements Animal`. The `extends` keyword is used for class inheritance or interface inheritance, not for implementing interfaces. `inherits` and `uses` are not Java keywords.

---

### Question 3: Multiple Interface Implementation - Flexibility
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, oop, interfaces, multiple-implementation

**Question Text**:
How many interfaces can a Java class implement?

**Answer Choices**:
- A) Only one
- B) Two at most
- C) Multiple interfaces
- D) None - classes cannot implement interfaces

**Correct Answer**: C

**Explanation**:
Unlike class inheritance (where a class can only extend one class), a class can implement multiple interfaces separated by commas: `class Duck implements Flyable, Swimmable`. This is one of the key advantages of interfaces—they allow a form of multiple inheritance in Java. Each interface represents a different capability or contract.

---

### Question 4: Interface Methods - Implementation Requirement
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, oop, interfaces, abstract-methods

**Question Text**:
What happens if a class implements an interface but doesn't provide implementations for all its methods?

**Answer Choices**:
- A) The program runs but those methods do nothing
- B) The missing methods use default implementations from Object class
- C) Compilation error unless the class is declared abstract
- D) A warning is shown but the code compiles

**Correct Answer**: C

**Explanation**:
If a class implements an interface, it MUST provide implementations for all abstract methods in that interface. If it doesn't, the class must be declared abstract. The compiler enforces this rule strictly—it's not just a warning. The only exception is if the interface methods have default implementations (Java 8+), in which case the class can choose to use those defaults.

---

### Question 5: Default Methods - Java 8 Feature
**Type**: code-output
**Difficulty**: 4
**Tags**: java, oop, interfaces, default-methods

**Question Text**:
What is true about default methods in interfaces?

```java
public interface Vehicle {
    void start();
    
    default void honk() {
        System.out.println("Beep!");
    }
}
```

**Answer Choices**:
- A) Classes implementing Vehicle must override the honk() method
- B) Classes implementing Vehicle get the honk() method but can optionally override it
- C) Default methods are not allowed in interfaces
- D) The honk() method will cause a compilation error

**Correct Answer**: B

**Explanation**:
Default methods (introduced in Java 8) provide a default implementation in the interface. Classes implementing the interface automatically get this implementation and can use it as-is, or they can optionally override it with their own implementation. This feature allows adding new methods to interfaces without breaking existing implementations.

---

### Question 6: Static Methods in Interfaces - Utility Functions
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, oop, interfaces, static-methods

**Question Text**:
How do you call a static method defined in an interface?

```java
public interface MathUtils {
    static int add(int a, int b) {
        return a + b;
    }
}
```

**Answer Choices**:
- A) Through an instance of a class implementing the interface
- B) Directly on the interface name: MathUtils.add(5, 3)
- C) Static methods cannot be defined in interfaces
- D) By creating an instance of the interface

**Correct Answer**: B

**Explanation**:
Static methods in interfaces (added in Java 8) are called directly on the interface name, just like static methods in classes: `MathUtils.add(5, 3)`. You don't need an instance of a class implementing the interface. These methods are typically used for utility functions related to the interface. You cannot create instances of interfaces, and these methods are not inherited by implementing classes.

---

### Question 7: Functional Interfaces - Lambda Expressions
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, oop, interfaces, functional-interfaces

**Question Text**:
What makes an interface a functional interface?

**Answer Choices**:
- A) It must be annotated with @FunctionalInterface
- B) It has exactly one abstract method
- C) All its methods must be static
- D) It cannot have any default methods

**Correct Answer**: B

**Explanation**:
A functional interface is an interface with exactly one abstract method. The `@FunctionalInterface` annotation is optional but recommended—it tells the compiler to verify this constraint. Functional interfaces can have multiple default and static methods; only the abstract method count matters. Functional interfaces are crucial for lambda expressions, which provide the implementation for that single abstract method.

---

### Question 8: Interface vs Abstract Class - Design Choice
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, oop, interfaces, abstract-classes, design

**Question Text**:
When should you prefer an interface over an abstract class?

**Answer Choices**:
- A) When you need to provide shared implementation code
- B) When you need to define non-static, non-final fields
- C) When you want multiple inheritance or defining a contract for unrelated classes
- D) When you want to use constructors

**Correct Answer**: C

**Explanation**:
Interfaces are preferred when you want to define a contract (what methods must exist) without specifying implementation, especially for unrelated classes. A class can implement multiple interfaces but can only extend one abstract class, so interfaces provide more flexibility. Abstract classes are better when you need shared implementation code, instance fields, or constructors. Options A, B, and D all describe situations where abstract classes are more appropriate.

---
