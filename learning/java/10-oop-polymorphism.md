# Java Object-Oriented Programming - Polymorphism

Polymorphism is one of the most powerful features of object-oriented programming, and it might sound intimidating at first, but it's actually quite intuitive. The word "polymorphism" comes from Greek: "poly" means many, and "morph" means forms. So polymorphism means "many forms"—the ability of objects to take on different forms and the ability of a single interface to represent different types.

Think of it like a smartphone. You can use the same charging cable (interface) with different phones (objects of different types), even though internally they might charge differently. That's polymorphism in action!

## Concept of Polymorphism

Polymorphism allows you to use a single interface to represent different underlying forms (data types). In Java, this means you can write code that works with a superclass type, and that code will work with any subclass of that superclass.

Here's a simple example:

```java
Animal animal1 = new Dog();
Animal animal2 = new Cat();
Animal animal3 = new Bird();

animal1.makeSound();  // Woof!
animal2.makeSound();  // Meow!
animal3.makeSound();  // Tweet!
```

Even though all three variables are of type `Animal`, each calls a different implementation of `makeSound()`. That's polymorphism—same method call, different behavior.

## Compile-Time Polymorphism (Method Overloading)

Compile-time polymorphism, also called static polymorphism, happens when the compiler determines which method to call based on the method signature. The most common form is method overloading.

### Method Overloading

Method overloading means having multiple methods with the same name but different parameters:

```java
public class Calculator {
    // Overloaded add methods
    public int add(int a, int b) {
        return a + b;
    }
    
    public double add(double a, double b) {
        return a + b;
    }
    
    public int add(int a, int b, int c) {
        return a + b + c;
    }
    
    public String add(String a, String b) {
        return a + b;
    }
}
```

Usage:

```java
Calculator calc = new Calculator();

System.out.println(calc.add(5, 3));           // Output: 8 (int version)
System.out.println(calc.add(5.5, 3.2));       // Output: 8.7 (double version)
System.out.println(calc.add(1, 2, 3));        // Output: 6 (three int version)
System.out.println(calc.add("Hello", "World")); // Output: HelloWorld (String version)
```

The compiler knows which method to call based on the arguments you pass. This decision happens at compile time, hence "compile-time polymorphism."

### Rules for Method Overloading

Methods must differ in:
1. **Number of parameters**: `add(int a)` vs. `add(int a, int b)`
2. **Type of parameters**: `add(int a)` vs. `add(double a)`
3. **Order of parameters**: `add(int a, String b)` vs. `add(String a, int b)`

These do NOT count as overloading:
- Different return types alone: `int add(int a)` vs. `double add(int a)` - INVALID
- Different parameter names alone: `add(int a)` vs. `add(int x)` - INVALID

### Constructor Overloading

Constructors can also be overloaded:

```java
public class Person {
    private String name;
    private int age;
    private String email;
    
    public Person(String name) {
        this.name = name;
        this.age = 0;
        this.email = "";
    }
    
    public Person(String name, int age) {
        this.name = name;
        this.age = age;
        this.email = "";
    }
    
    public Person(String name, int age, String email) {
        this.name = name;
        this.age = age;
        this.email = email;
    }
}
```

## Runtime Polymorphism (Method Overriding)

Runtime polymorphism, also called dynamic polymorphism, happens when the actual method to be called is determined at runtime based on the actual object type, not the reference type.

### Method Overriding

We covered this in the Inheritance topic, but it's central to runtime polymorphism:

```java
public class Animal {
    public void makeSound() {
        System.out.println("Some generic sound");
    }
    
    public void eat() {
        System.out.println("Eating...");
    }
}

public class Dog extends Animal {
    @Override
    public void makeSound() {
        System.out.println("Woof! Woof!");
    }
    
    @Override
    public void eat() {
        System.out.println("Dog eating kibble");
    }
}

public class Cat extends Animal {
    @Override
    public void makeSound() {
        System.out.println("Meow!");
    }
    
    @Override
    public void eat() {
        System.out.println("Cat eating fish");
    }
}
```

Now the magic:

```java
Animal myPet;  // Reference of type Animal

myPet = new Dog();
myPet.makeSound();  // Output: Woof! Woof!
myPet.eat();        // Output: Dog eating kibble

myPet = new Cat();
myPet.makeSound();  // Output: Meow!
myPet.eat();        // Output: Cat eating fish
```

The same variable `myPet` can refer to different types of objects, and the correct method is called based on the actual object type at runtime. This is the essence of runtime polymorphism.

## Dynamic Method Dispatch

Dynamic method dispatch is the mechanism by which a call to an overridden method is resolved at runtime rather than compile time. It's how Java implements runtime polymorphism.

```java
public class Shape {
    public void draw() {
        System.out.println("Drawing a shape");
    }
    
    public double getArea() {
        return 0.0;
    }
}

public class Circle extends Shape {
    private double radius;
    
    public Circle(double radius) {
        this.radius = radius;
    }
    
    @Override
    public void draw() {
        System.out.println("Drawing a circle");
    }
    
    @Override
    public double getArea() {
        return Math.PI * radius * radius;
    }
}

public class Rectangle extends Shape {
    private double width;
    private double height;
    
    public Rectangle(double width, double height) {
        this.width = width;
        this.height = height;
    }
    
    @Override
    public void draw() {
        System.out.println("Drawing a rectangle");
    }
    
    @Override
    public double getArea() {
        return width * height;
    }
}
```

Using dynamic method dispatch:

```java
Shape[] shapes = new Shape[3];
shapes[0] = new Circle(5.0);
shapes[1] = new Rectangle(4.0, 6.0);
shapes[2] = new Circle(3.0);

for (Shape shape : shapes) {
    shape.draw();                          // Calls the correct draw() method
    System.out.println("Area: " + shape.getArea());  // Calls the correct getArea() method
}
```

Output:
```
Drawing a circle
Area: 78.53981633974483
Drawing a rectangle
Area: 24.0
Drawing a circle
Area: 28.274333882308138
```

The JVM determines at runtime which version of `draw()` and `getArea()` to call based on the actual object type.

## Upcasting and Downcasting

### Upcasting

Upcasting is converting a subclass reference to a superclass reference. It's automatic and always safe:

```java
Dog myDog = new Dog();
Animal myAnimal = myDog;  // Upcasting - automatic
```

This is safe because a Dog IS-A Animal—it has all the methods and fields that Animal has.

### Downcasting

Downcasting is converting a superclass reference to a subclass reference. It requires an explicit cast and can fail at runtime:

```java
Animal myAnimal = new Dog();           // Upcasting
Dog myDog = (Dog) myAnimal;            // Downcasting - requires explicit cast

myDog.bark();  // Now we can access Dog-specific methods
```

But be careful:

```java
Animal myAnimal = new Cat();
Dog myDog = (Dog) myAnimal;  // Compiles, but throws ClassCastException at runtime!
```

The compiler allows the cast, but at runtime, Java sees you're trying to treat a Cat as a Dog and throws an exception.

### Safe Downcasting with `instanceof`

Always check the type before downcasting:

```java
Animal myAnimal = getAnimal();  // Returns some animal

if (myAnimal instanceof Dog) {
    Dog myDog = (Dog) myAnimal;
    myDog.bark();
} else if (myAnimal instanceof Cat) {
    Cat myCat = (Cat) myAnimal;
    myCat.meow();
}
```

Modern Java (Java 16+) has pattern matching for instanceof:

```java
Animal myAnimal = getAnimal();

if (myAnimal instanceof Dog myDog) {
    myDog.bark();  // myDog is automatically cast
} else if (myAnimal instanceof Cat myCat) {
    myCat.meow();  // myCat is automatically cast
}
```

## The `instanceof` Operator

The `instanceof` operator checks if an object is an instance of a specific class or implements an interface:

```java
Dog myDog = new Dog();

System.out.println(myDog instanceof Dog);     // true
System.out.println(myDog instanceof Animal);  // true (Dog extends Animal)
System.out.println(myDog instanceof Object);  // true (all classes extend Object)
System.out.println(myDog instanceof Cat);     // false
```

Using with null:

```java
Animal myAnimal = null;
System.out.println(myAnimal instanceof Animal);  // false (null is not an instance of anything)
```

Practical use case:

```java
public void processAnimal(Animal animal) {
    animal.eat();  // All animals can eat
    
    // Special handling for specific types
    if (animal instanceof Dog) {
        Dog dog = (Dog) animal;
        dog.fetch();
    } else if (animal instanceof Cat) {
        Cat cat = (Cat) animal;
        cat.climb();
    } else if (animal instanceof Bird) {
        Bird bird = (Bird) animal;
        bird.fly();
    }
}
```

## Reference Type vs. Object Type

This is crucial to understanding polymorphism:

- **Reference type**: The type of the variable (determined at compile time)
- **Object type**: The actual type of the object in memory (determined at runtime)

```java
Animal myAnimal = new Dog();
// Reference type: Animal
// Object type: Dog
```

**Important rules:**
1. The reference type determines which methods you can call (compile-time check)
2. The object type determines which method implementation is executed (runtime behavior)

Example:

```java
public class Animal {
    public void eat() {
        System.out.println("Animal eating");
    }
}

public class Dog extends Animal {
    @Override
    public void eat() {
        System.out.println("Dog eating");
    }
    
    public void bark() {
        System.out.println("Woof!");
    }
}

Animal myAnimal = new Dog();
myAnimal.eat();   // Calls Dog's eat() - object type determines implementation
// myAnimal.bark();  // Compilation error! Reference type is Animal, which has no bark() method

Dog myDog = (Dog) myAnimal;  // Downcast to access Dog-specific methods
myDog.bark();  // Now this works
```

## Polymorphism with Arrays and Collections

One of the most powerful uses of polymorphism is working with collections of different object types through a common superclass:

```java
Animal[] animals = new Animal[5];
animals[0] = new Dog();
animals[1] = new Cat();
animals[2] = new Bird();
animals[3] = new Dog();
animals[4] = new Cat();

for (Animal animal : animals) {
    animal.makeSound();  // Each animal makes its own sound
    animal.eat();        // Each animal eats in its own way
}
```

With ArrayList:

```java
ArrayList<Animal> zoo = new ArrayList<>();
zoo.add(new Dog());
zoo.add(new Cat());
zoo.add(new Bird());

for (Animal animal : zoo) {
    animal.makeSound();
}
```

This is incredibly powerful—you can write code once that works with any subclass, present or future!

## Polymorphism in Practice

### Example: Payment System

```java
public abstract class Payment {
    protected double amount;
    
    public Payment(double amount) {
        this.amount = amount;
    }
    
    public abstract void processPayment();
    public abstract String getPaymentDetails();
}

public class CreditCardPayment extends Payment {
    private String cardNumber;
    
    public CreditCardPayment(double amount, String cardNumber) {
        super(amount);
        this.cardNumber = cardNumber;
    }
    
    @Override
    public void processPayment() {
        System.out.println("Processing credit card payment of $" + amount);
        // Credit card processing logic
    }
    
    @Override
    public String getPaymentDetails() {
        return "Credit Card ending in " + cardNumber.substring(cardNumber.length() - 4);
    }
}

public class PayPalPayment extends Payment {
    private String email;
    
    public PayPalPayment(double amount, String email) {
        super(amount);
        this.email = email;
    }
    
    @Override
    public void processPayment() {
        System.out.println("Processing PayPal payment of $" + amount);
        // PayPal processing logic
    }
    
    @Override
    public String getPaymentDetails() {
        return "PayPal account: " + email;
    }
}

public class CryptoPayment extends Payment {
    private String walletAddress;
    
    public CryptoPayment(double amount, String walletAddress) {
        super(amount);
        this.walletAddress = walletAddress;
    }
    
    @Override
    public void processPayment() {
        System.out.println("Processing cryptocurrency payment of $" + amount);
        // Crypto processing logic
    }
    
    @Override
    public String getPaymentDetails() {
        return "Crypto wallet: " + walletAddress;
    }
}
```

Using polymorphism:

```java
public class PaymentProcessor {
    public void processPayments(List<Payment> payments) {
        for (Payment payment : payments) {
            System.out.println(payment.getPaymentDetails());
            payment.processPayment();
            System.out.println("Payment processed successfully!\n");
        }
    }
}

// Usage
List<Payment> payments = new ArrayList<>();
payments.add(new CreditCardPayment(99.99, "1234567890123456"));
payments.add(new PayPalPayment(49.99, "user@example.com"));
payments.add(new CryptoPayment(199.99, "1A2B3C4D5E6F..."));

PaymentProcessor processor = new PaymentProcessor();
processor.processPayments(payments);  // Handles all payment types polymorphically
```

## Common Pitfalls

1. **Confusing overloading with overriding**: Overloading is compile-time (same name, different parameters), overriding is runtime (same signature, different implementation)
2. **Unsafe downcasting**: Always use `instanceof` before downcasting to avoid ClassCastException
3. **Forgetting @Override**: Without it, you might accidentally overload instead of override
4. **Private methods can't be overridden**: They're not visible to subclasses
5. **Static methods can't be overridden**: They can be hidden, but that's different from overriding

## Best Practices

1. **Program to an interface**: Use the most general type possible for variables and parameters
2. **Use `instanceof` sparingly**: If you're using it a lot, reconsider your design—you might need polymorphism
3. **Prefer composition over type checking**: Instead of `if (x instanceof Type)`, design better class hierarchies
4. **Always use @Override**: Catches mistakes and documents intent
5. **Keep hierarchies shallow**: Deep inheritance makes code hard to understand
6. **Document polymorphic behavior**: Make it clear how subclasses should override methods

## Practice Exercises

1. Create a `Shape` hierarchy with polymorphic `draw()` and `getArea()` methods
2. Create a notification system with `Email`, `SMS`, and `Push` notification types, all extending `Notification`
3. Implement a simple game with different types of `Character` objects that attack differently
4. Create a calculator that uses method overloading for different numeric types
5. Build a transportation system with `Vehicle` superclass and various vehicle types with polymorphic `move()` methods

---

## Quiz Questions

This section contains quiz questions that test understanding of Java polymorphism. Each question is formatted for easy extraction and database import.

### Question 1: Polymorphism Concept - Basic Understanding
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: java, oop, polymorphism

**Question Text**:
What does polymorphism mean in object-oriented programming?

**Answer Choices**:
- A) A class can have multiple constructors
- B) An object can take many forms and use the same interface to represent different types
- C) A method can return different types
- D) Multiple classes can have the same name

**Correct Answer**: B

**Explanation**:
Polymorphism (from Greek "poly" = many, "morph" = forms) means "many forms." In OOP, it allows objects of different types to be accessed through the same interface. For example, you can have an `Animal` reference that points to a `Dog`, `Cat`, or `Bird` object, and calling `makeSound()` will execute the appropriate method for each type. Option A describes constructor overloading, which is a form of polymorphism but not the complete definition.

---

### Question 2: Method Overloading - Compile-Time Polymorphism
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, oop, polymorphism, method-overloading

**Question Text**:
Which of the following represents valid method overloading?

**Answer Choices**:
- A) `void print(int a)` and `int print(int a)`
- B) `void print(int a)` and `void print(double a)`
- C) `void print(int a)` and `void print(int b)`
- D) `void print(int a)` and `void Print(int a)`

**Correct Answer**: B

**Explanation**:
Method overloading requires the same method name with different parameter types, number of parameters, or order of parameters. Option B has the same name but different parameter types (int vs. double), which is valid overloading. Option A differs only in return type, which isn't allowed. Option C differs only in parameter name, which doesn't count. Option D differs in method name case, so these are different methods entirely.

---

### Question 3: Method Overriding - Runtime Polymorphism
**Type**: code-output
**Difficulty**: 3
**Tags**: java, oop, polymorphism, method-overriding

**Question Text**:
What will be the output of this code?

```java
class Animal {
    public void makeSound() {
        System.out.println("Some sound");
    }
}

class Dog extends Animal {
    @Override
    public void makeSound() {
        System.out.println("Woof!");
    }
}

public class Test {
    public static void main(String[] args) {
        Animal a = new Dog();
        a.makeSound();
    }
}
```

**Answer Choices**:
- A) Some sound
- B) Woof!
- C) Compilation error
- D) Both "Some sound" and "Woof!"

**Correct Answer**: B

**Explanation**:
This demonstrates runtime polymorphism. Even though the reference type is `Animal`, the actual object type is `Dog`. At runtime, Java calls the overridden `makeSound()` method from the Dog class, not the Animal class. This is called dynamic method dispatch—the method to execute is determined by the object type, not the reference type.

---

### Question 4: Dynamic Method Dispatch - Method Resolution
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, oop, polymorphism, dynamic-dispatch

**Question Text**:
When is the decision made about which overridden method to call?

**Answer Choices**:
- A) At compile time based on the reference type
- B) At runtime based on the actual object type
- C) At class loading time
- D) It's randomly chosen

**Correct Answer**: B

**Explanation**:
Dynamic method dispatch means the JVM determines which overridden method to call at runtime based on the actual object type, not the reference type. This is the mechanism that enables runtime polymorphism. The compiler only checks that the method exists for the reference type, but the actual method executed is determined when the program runs based on what type of object the reference points to.

---

### Question 5: Upcasting and Downcasting - Type Conversion
**Type**: code-output
**Difficulty**: 4
**Tags**: java, oop, polymorphism, upcasting, downcasting

**Question Text**:
What happens when this code runs?

```java
Animal animal = new Cat();
Dog dog = (Dog) animal;
```

**Answer Choices**:
- A) Works fine - dog is now a Dog object
- B) Compilation error
- C) ClassCastException at runtime
- D) NullPointerException

**Correct Answer**: C

**Explanation**:
While the code compiles (the compiler sees you're casting an Animal reference to Dog), at runtime the actual object is a Cat. You cannot cast a Cat to a Dog—they're incompatible types. This results in a ClassCastException at runtime. To safely downcast, always use `instanceof` to check the type first: `if (animal instanceof Dog) { Dog dog = (Dog) animal; }`.

---

### Question 6: instanceof Operator - Type Checking
**Type**: code-output
**Difficulty**: 3
**Tags**: java, oop, polymorphism, instanceof

**Question Text**:
What will be printed?

```java
Dog myDog = new Dog();
System.out.println(myDog instanceof Animal);
```

Assuming Dog extends Animal.

**Answer Choices**:
- A) false
- B) true
- C) Compilation error
- D) ClassCastException

**Correct Answer**: B

**Explanation**:
The `instanceof` operator checks if an object is an instance of a specific class or any of its superclasses. Since Dog extends Animal, a Dog object IS-A Animal, so `myDog instanceof Animal` returns true. The instanceof operator checks the inheritance hierarchy, not just the exact type.

---

### Question 7: Reference vs Object Type - Method Access
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, oop, polymorphism, reference-type

**Question Text**:
Given `Animal myAnimal = new Dog();`, and Dog has a `bark()` method that Animal doesn't have, what happens if you try to call `myAnimal.bark()`?

**Answer Choices**:
- A) It works fine and calls Dog's bark() method
- B) Compilation error - bark() is not defined in Animal
- C) Runtime error - bark() is not available
- D) It calls a default bark() method from Object class

**Correct Answer**: B

**Explanation**:
The reference type (Animal) determines which methods you can call at compile time. Even though the actual object is a Dog, the compiler only knows it's an Animal reference and doesn't allow calling Dog-specific methods. To call `bark()`, you need to downcast: `((Dog) myAnimal).bark()` or `Dog myDog = (Dog) myAnimal; myDog.bark();`. This is a compile-time restriction based on the reference type.

---

### Question 8: Polymorphism with Collections - Practical Application
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, oop, polymorphism, collections

**Question Text**:
What is the main advantage of using polymorphism with collections?

```java
List<Animal> animals = new ArrayList<>();
animals.add(new Dog());
animals.add(new Cat());
animals.add(new Bird());
```

**Answer Choices**:
- A) It makes the code run faster
- B) It allows you to write code that works with many different types through a common interface
- C) It reduces memory usage
- D) It's required by the Java compiler

**Correct Answer**: B

**Explanation**:
The main advantage of polymorphism is code reusability and flexibility. You can write one piece of code (like a loop that calls `makeSound()` on each animal) that works with any subclass of Animal, both existing and future ones. This makes code more maintainable and extensible. It doesn't improve performance or reduce memory usage, and it's not required—it's a design choice that makes code more flexible.

---
