# Java Object-Oriented Programming - Classes and Objects

Welcome to the world of Object-Oriented Programming (OOP)! If you've been writing code so far with just variables and methods, you're about to level up. OOP is like organizing your code into neat little packages that mirror real-world things. Instead of having scattered variables and functions everywhere, you'll create blueprints (classes) and build actual things (objects) from them.

Think of a class as a cookie cutter and objects as the cookies you make with it. The cookie cutter defines the shape, but each cookie is its own thing.

## Understanding Object-Oriented Programming Concepts

Object-Oriented Programming is a way of organizing code around "objects" that combine data and behavior. It's based on a few key principles:

- **Encapsulation**: Bundling data and methods that work on that data together
- **Inheritance**: Creating new classes based on existing ones
- **Polymorphism**: Using the same interface for different types
- **Abstraction**: Hiding complex implementation details

These might sound abstract now, but they'll make sense as we go through examples. Let's start with the basics: classes and objects.

## Classes and Objects

### What's a Class?

A class is a blueprint or template for creating objects. It defines what properties (data) and behaviors (methods) objects of that type will have.

```java
public class Dog {
    // Properties (instance variables)
    String name;
    int age;
    String breed;
    
    // Behavior (instance methods)
    void bark() {
        System.out.println(name + " says: Woof!");
    }
    
    void celebrate() {
        System.out.println(name + " is " + age + " years old!");
    }
}
```

This `Dog` class is a blueprint. It says "every dog has a name, age, and breed, and every dog can bark and celebrate."

### Creating Objects (Instantiation)

Once you have a class, you can create objects from it:

```java
public class Main {
    public static void main(String[] args) {
        // Create a Dog object
        Dog myDog = new Dog();
        myDog.name = "Buddy";
        myDog.age = 3;
        myDog.breed = "Golden Retriever";
        
        myDog.bark();        // Output: Buddy says: Woof!
        myDog.celebrate();   // Output: Buddy is 3 years old!
        
        // Create another Dog object
        Dog yourDog = new Dog();
        yourDog.name = "Max";
        yourDog.age = 5;
        yourDog.breed = "Labrador";
        
        yourDog.bark();      // Output: Max says: Woof!
    }
}
```

Notice that `myDog` and `yourDog` are completely separate objects. Changing one doesn't affect the other, even though they're both created from the same class.

### Constructors

Setting up each object's properties one by one gets tedious. That's where constructors come in—special methods that initialize objects when they're created:

```java
public class Dog {
    String name;
    int age;
    String breed;
    
    // Constructor
    public Dog(String name, int age, String breed) {
        this.name = name;
        this.age = age;
        this.breed = breed;
    }
    
    void bark() {
        System.out.println(name + " says: Woof!");
    }
}
```

Now you can create dogs more easily:

```java
Dog myDog = new Dog("Buddy", 3, "Golden Retriever");
Dog yourDog = new Dog("Max", 5, "Labrador");

myDog.bark();  // Output: Buddy says: Woof!
```

### Constructor Overloading

You can have multiple constructors with different parameters:

```java
public class Dog {
    String name;
    int age;
    String breed;
    
    // Constructor with all parameters
    public Dog(String name, int age, String breed) {
        this.name = name;
        this.age = age;
        this.breed = breed;
    }
    
    // Constructor with just name
    public Dog(String name) {
        this.name = name;
        this.age = 0;  // Default values
        this.breed = "Unknown";
    }
    
    // No-argument constructor
    public Dog() {
        this.name = "Unnamed";
        this.age = 0;
        this.breed = "Unknown";
    }
}
```

This gives you flexibility in how you create objects:

```java
Dog dog1 = new Dog("Buddy", 3, "Golden Retriever");
Dog dog2 = new Dog("Max");
Dog dog3 = new Dog();
```

### The `this` Keyword

You've seen `this` in the examples above. It refers to the current object—the one whose method or constructor is being called. It's especially useful when parameter names match field names:

```java
public Dog(String name, int age) {
    this.name = name;  // this.name is the instance variable
                        // name is the parameter
    this.age = age;
}
```

You can also use `this` to call other constructors from within a constructor:

```java
public class Dog {
    String name;
    int age;
    String breed;
    
    public Dog(String name, int age, String breed) {
        this.name = name;
        this.age = age;
        this.breed = breed;
    }
    
    public Dog(String name) {
        this(name, 0, "Unknown");  // Calls the three-parameter constructor
    }
}
```

## Access Modifiers

Access modifiers control who can see and use your class members. Think of them as privacy settings.

### The Four Access Levels

1. **`public`**: Anyone can access it
2. **`private`**: Only code within the same class can access it
3. **`protected`**: Code in the same package or subclasses can access it
4. **Package-private (no modifier)**: Only code in the same package can access it

```java
public class BankAccount {
    public String accountNumber;      // Anyone can access
    private double balance;           // Only this class can access
    protected String ownerName;       // Same package or subclasses
    String bankName;                  // Package-private (same package only)
    
    public void deposit(double amount) {
        if (amount > 0) {
            balance += amount;  // We can access private balance here
        }
    }
    
    public double getBalance() {
        return balance;  // Provide controlled access to private data
    }
}
```

Why make things private? To protect your data and ensure it's only modified in safe ways.

## Encapsulation

Encapsulation is the practice of keeping fields private and providing public methods to access and modify them. It's one of the core principles of OOP.

### Why Encapsulation?

Imagine a `BankAccount` class where anyone could directly change the balance:

```java
// BAD: Direct access to fields
public class BankAccount {
    public double balance;
}

// Someone could do this:
account.balance = -1000;  // Negative balance? That shouldn't happen!
```

With encapsulation, you control how data is accessed and modified:

```java
// GOOD: Encapsulated data
public class BankAccount {
    private double balance;  // Private field
    
    // Getter method
    public double getBalance() {
        return balance;
    }
    
    // Setter method with validation
    public void setBalance(double balance) {
        if (balance >= 0) {
            this.balance = balance;
        } else {
            System.out.println("Balance cannot be negative!");
        }
    }
    
    public void deposit(double amount) {
        if (amount > 0) {
            balance += amount;
        }
    }
    
    public boolean withdraw(double amount) {
        if (amount > 0 && amount <= balance) {
            balance -= amount;
            return true;
        }
        return false;
    }
}
```

Now the balance can only be changed through your controlled methods, which can validate the data.

### Getters and Setters

Methods that get and set private fields are called getters and setters (or accessor and mutator methods):

```java
public class Person {
    private String name;
    private int age;
    
    // Getter for name
    public String getName() {
        return name;
    }
    
    // Setter for name
    public void setName(String name) {
        this.name = name;
    }
    
    // Getter for age
    public int getAge() {
        return age;
    }
    
    // Setter for age with validation
    public void setAge(int age) {
        if (age >= 0 && age <= 150) {
            this.age = age;
        }
    }
}
```

## Static Members

### Static Variables

Sometimes you want data that belongs to the class itself, not to individual objects. That's what `static` is for:

```java
public class Student {
    private String name;
    private static int studentCount = 0;  // Shared by all Student objects
    
    public Student(String name) {
        this.name = name;
        studentCount++;  // Increment for each new student
    }
    
    public static int getStudentCount() {
        return studentCount;
    }
}
```

Usage:

```java
Student s1 = new Student("Alice");
Student s2 = new Student("Bob");
Student s3 = new Student("Charlie");

System.out.println(Student.getStudentCount());  // Output: 3
```

Notice we call `Student.getStudentCount()` on the class itself, not on an object.

### Static Methods

Static methods belong to the class, not to instances. They can't access instance variables because there's no specific object to work with:

```java
public class MathHelper {
    // Static method - can be called without creating an object
    public static int add(int a, int b) {
        return a + b;
    }
    
    public static int multiply(int a, int b) {
        return a * b;
    }
}

// Usage
int sum = MathHelper.add(5, 3);         // Output: 8
int product = MathHelper.multiply(4, 7); // Output: 28
```

The `Math` class in Java's standard library is a great example—all its methods are static:

```java
double result = Math.sqrt(16);    // 4.0
double power = Math.pow(2, 3);    // 8.0
```

### Static Blocks

Static blocks run once when the class is first loaded, before any objects are created:

```java
public class Database {
    private static Connection connection;
    
    // Static block - runs once when class is loaded
    static {
        System.out.println("Initializing database connection...");
        connection = createConnection();
    }
    
    private static Connection createConnection() {
        // Setup code here
        return new Connection();
    }
}
```

### Static Imports

You can import static members to use them without the class name:

```java
import static java.lang.Math.*;

public class Calculator {
    public double calculateCircleArea(double radius) {
        return PI * pow(radius, 2);  // No need to write Math.PI or Math.pow
    }
}
```

## The `final` Keyword

### Final Variables (Constants)

The `final` keyword means "can't be changed":

```java
public class Circle {
    private final double PI = 3.14159;  // Can't be changed
    private final double radius;
    
    public Circle(double radius) {
        this.radius = radius;  // Can only be set once, in constructor
    }
    
    public double getArea() {
        return PI * radius * radius;
    }
}
```

For constants shared across all instances, combine `static` and `final`:

```java
public class Constants {
    public static final double PI = 3.14159;
    public static final int MAX_STUDENTS = 100;
    public static final String SCHOOL_NAME = "Java Academy";
}

// Usage
double area = Constants.PI * radius * radius;
```

### Final Methods

A final method can't be overridden by subclasses:

```java
public class Vehicle {
    public final void startEngine() {
        System.out.println("Engine started");
    }
}

public class Car extends Vehicle {
    // This would cause a compiler error:
    // public void startEngine() { ... }  // Can't override final method
}
```

### Final Classes

A final class can't be extended:

```java
public final class SecurityUtility {
    // This class cannot be subclassed
    public static String encrypt(String data) {
        // Encryption code
        return encryptedData;
    }
}

// This would cause a compiler error:
// public class MySecurityUtility extends SecurityUtility { }
```

The `String` class in Java is final—you can't create a subclass of String.

## Nested and Inner Classes

### Inner Classes

You can define a class inside another class:

```java
public class OuterClass {
    private String outerField = "Outer";
    
    public class InnerClass {
        public void display() {
            System.out.println("Accessing: " + outerField);  // Can access outer's private members
        }
    }
}

// Usage
OuterClass outer = new OuterClass();
OuterClass.InnerClass inner = outer.new InnerClass();
inner.display();  // Output: Accessing: Outer
```

### Static Nested Classes

A static nested class doesn't need an instance of the outer class:

```java
public class OuterClass {
    private static String staticField = "Static Outer";
    
    public static class StaticNestedClass {
        public void display() {
            System.out.println("Accessing: " + staticField);
        }
    }
}

// Usage
OuterClass.StaticNestedClass nested = new OuterClass.StaticNestedClass();
nested.display();  // Output: Accessing: Static Outer
```

### Local Classes

You can even define a class inside a method:

```java
public class OuterClass {
    public void someMethod() {
        final String localVar = "Local";
        
        class LocalClass {
            public void display() {
                System.out.println("Accessing: " + localVar);
            }
        }
        
        LocalClass local = new LocalClass();
        local.display();
    }
}
```

### Anonymous Classes

Anonymous classes are one-time-use classes defined inline:

```java
interface Greeting {
    void greet();
}

public class Main {
    public static void main(String[] args) {
        Greeting greeting = new Greeting() {
            public void greet() {
                System.out.println("Hello!");
            }
        };
        
        greeting.greet();  // Output: Hello!
    }
}
```

## Common Pitfalls

1. **Forgetting to initialize objects**: Using `Dog myDog;` without `new Dog()` leads to NullPointerException
2. **Making everything public**: Breaks encapsulation and makes code harder to maintain
3. **Confusing static and instance members**: You can't access instance variables from static methods
4. **Not using constructors effectively**: Leads to repetitive initialization code
5. **Circular dependencies**: Class A depends on Class B which depends on Class A

## Best Practices

1. **Keep fields private**: Use getters and setters for controlled access
2. **Use meaningful names**: `Customer` is better than `C` or `Cust`
3. **One class per file**: Makes code easier to find and maintain
4. **Initialize in constructors**: Set up objects completely when they're created
5. **Use `this` for clarity**: Makes it clear when you're referring to instance members
6. **Document your classes**: Use Javadoc comments to explain what the class does
7. **Keep classes focused**: Each class should have one clear responsibility

## Practice Exercises

1. Create a `Book` class with properties (title, author, ISBN, price) and methods (display info, apply discount)
2. Create a `Rectangle` class with private width and height, with methods to calculate area and perimeter
3. Create a `BankAccount` class with encapsulation, ensuring balance can't go negative
4. Create a `Student` class with a static counter to track total number of students created
5. Create a `Car` class with nested `Engine` class

---

## Quiz Questions

This section contains quiz questions that test understanding of Java classes and objects. Each question is formatted for easy extraction and database import.

### Question 1: Class Definition - Basic Concepts
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: java, oop, classes, objects

**Question Text**:
What is a class in Java?

**Answer Choices**:
- A) A variable that stores data
- B) A blueprint or template for creating objects
- C) A method that performs an action
- D) A type of loop structure

**Correct Answer**: B

**Explanation**:
A class is a blueprint or template that defines the properties (fields) and behaviors (methods) that objects created from that class will have. It's like a cookie cutter that defines the shape of cookies (objects). Option A describes a variable, Option C describes a method, and Option D describes a control structure.

---

### Question 2: Object Instantiation - Creating Objects
**Type**: code-output
**Difficulty**: 2
**Tags**: java, oop, objects, instantiation

**Question Text**:
What will be the output of this code?

```java
public class Dog {
    String name = "Unknown";
}

public class Main {
    public static void main(String[] args) {
        Dog dog1 = new Dog();
        dog1.name = "Buddy";
        
        Dog dog2 = new Dog();
        dog2.name = "Max";
        
        System.out.println(dog1.name);
    }
}
```

**Answer Choices**:
- A) Unknown
- B) Max
- C) Buddy
- D) Compilation error

**Correct Answer**: C

**Explanation**:
Each object created with `new Dog()` is independent. Setting `dog2.name = "Max"` doesn't affect `dog1.name`. The code creates two separate Dog objects, sets dog1's name to "Buddy" and dog2's name to "Max", then prints dog1's name, which is "Buddy".

---

### Question 3: Constructors - Initialization
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, oop, constructors

**Question Text**:
Which statement about constructors is TRUE?

**Answer Choices**:
- A) Constructors must have a return type
- B) A class can have only one constructor
- C) Constructors are called automatically when an object is created
- D) Constructor names can be different from the class name

**Correct Answer**: C

**Explanation**:
Constructors are automatically called when you create an object using the `new` keyword. They don't have return types (not even void), which makes option A false. A class can have multiple constructors through overloading, making option B false. Constructor names must match the class name exactly, making option D false.

---

### Question 4: this Keyword - Self Reference
**Type**: code-output
**Difficulty**: 3
**Tags**: java, oop, this-keyword

**Question Text**:
What is the purpose of the `this` keyword in this constructor?

```java
public class Person {
    private String name;
    
    public Person(String name) {
        this.name = name;
    }
}
```

**Answer Choices**:
- A) It refers to the parent class
- B) It refers to the current object and distinguishes the instance variable from the parameter
- C) It creates a new object
- D) It's optional and serves no purpose

**Correct Answer**: B

**Explanation**:
The `this` keyword refers to the current object. In this case, `this.name` refers to the instance variable, while `name` alone refers to the parameter. This is necessary because both have the same name. Without `this`, you couldn't distinguish between them, and the assignment wouldn't work correctly.

---

### Question 5: Access Modifiers - Encapsulation
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, oop, access-modifiers, encapsulation

**Question Text**:
Which access modifier should be used for a field that should only be accessible within its own class?

**Answer Choices**:
- A) public
- B) private
- C) protected
- D) No modifier (package-private)

**Correct Answer**: B

**Explanation**:
The `private` access modifier restricts access to within the same class only. This is the most restrictive access level and is the foundation of encapsulation. Public allows access from anywhere, protected allows access from the same package and subclasses, and package-private allows access from the same package.

---

### Question 6: Static Members - Class vs Instance
**Type**: code-output
**Difficulty**: 4
**Tags**: java, oop, static, class-members

**Question Text**:
What will be printed?

```java
public class Counter {
    private static int count = 0;
    
    public Counter() {
        count++;
    }
    
    public static int getCount() {
        return count;
    }
}

public class Main {
    public static void main(String[] args) {
        Counter c1 = new Counter();
        Counter c2 = new Counter();
        Counter c3 = new Counter();
        System.out.println(Counter.getCount());
    }
}
```

**Answer Choices**:
- A) 0
- B) 1
- C) 3
- D) Compilation error

**Correct Answer**: C

**Explanation**:
Static variables are shared across all instances of a class. Each time a Counter object is created, the constructor increments the static `count` variable. After creating three Counter objects, count is 3. The static getCount() method returns this shared count value.

---

### Question 7: Final Keyword - Constants
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, oop, final, constants

**Question Text**:
What happens if you try to modify a `final` variable after it's been initialized?

**Answer Choices**:
- A) The change is allowed but a warning is shown
- B) A compilation error occurs
- C) A runtime exception is thrown
- D) The change is silently ignored

**Correct Answer**: B

**Explanation**:
The `final` keyword makes a variable constant—once initialized, it cannot be changed. Attempting to reassign a final variable results in a compilation error. This happens at compile time, not runtime, so the code won't even run. This is Java's way of enforcing immutability for variables marked as final.

---

### Question 8: Encapsulation - Getters and Setters
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, oop, encapsulation, getters-setters

**Question Text**:
What is the main benefit of using getters and setters instead of making fields public?

**Answer Choices**:
- A) They make the code run faster
- B) They allow controlled access and validation of data
- C) They reduce the amount of code needed
- D) They are required by the Java compiler

**Correct Answer**: B

**Explanation**:
Getters and setters provide controlled access to private fields. They allow you to add validation logic (e.g., ensuring age is positive), change the internal implementation without affecting external code, and add logging or other side effects when fields are accessed or modified. They don't improve performance, actually add more code, and aren't required by the compiler—they're a design choice for better encapsulation.

---
