# Java Methods

Welcome! Methods are one of the most powerful tools in Java. Think of a method as a recipe—it's a named set of instructions that you can use over and over. Instead of writing the same code repeatedly, you write it once in a method and call it whenever you need it.

Methods help you organize code, avoid repetition, and make programs easier to understand and maintain. Let's dive in!

## Method Declaration and Definition

A method is a block of code that performs a specific task. When you create a method, you're defining what it does. When you use it, you're calling (or invoking) it.

### Basic Method Syntax

```java
returnType methodName(parameters) {
    // Method body - code that runs when method is called
    return value;  // If returnType is not void
}
```

Let's break down each part:
- **Return type**: What kind of value the method gives back (or `void` if it doesn't return anything)
- **Method name**: What you call the method (use camelCase)
- **Parameters**: Input values the method needs (in parentheses)
- **Method body**: The code between the curly braces
- **Return statement**: Sends a value back to whoever called the method

### Simple Example

```java
public class Calculator {
    // Method that adds two numbers
    public static int add(int a, int b) {
        int sum = a + b;
        return sum;
    }
    
    public static void main(String[] args) {
        int result = add(5, 3);  // Call the method
        System.out.println(result);  // Prints 8
    }
}
```

In this example, `add` is a method that takes two integers, adds them together, and returns the result. The method definition explains *what* it does, and calling it actually *does* it.

### Methods with No Return Value (void)

Some methods perform actions but don't return a value. Use `void` for the return type:

```java
public static void greet(String name) {
    System.out.println("Hello, " + name + "!");
}

public static void main(String[] args) {
    greet("Alex");  // Prints: Hello, Alex!
    greet("Jordan");  // Prints: Hello, Jordan!
}
```

No `return` statement needed for `void` methods! They just do their job and finish.

## Method Parameters and Arguments

Parameters are the variables listed in the method definition. Arguments are the actual values you pass when calling the method.

### Parameters: The Method's Inputs

```java
// num1 and num2 are parameters
public static int multiply(int num1, int num2) {
    return num1 * num2;
}
```

Parameters act like placeholders. When someone calls the method, they provide the actual values.

### Arguments: The Actual Values

```java
int result = multiply(4, 5);  // 4 and 5 are arguments
System.out.println(result);  // 20
```

Here's a helpful way to think about it: parameters are the recipe's ingredients list, while arguments are the actual ingredients you use.

### Multiple Parameters

Methods can have multiple parameters of different types:

```java
public static void displayInfo(String name, int age, double height) {
    System.out.println("Name: " + name);
    System.out.println("Age: " + age);
    System.out.println("Height: " + height + " meters");
}

// Call with three arguments
displayInfo("Alex", 25, 1.75);
```

The order matters! The first argument goes to the first parameter, the second to the second, and so on.

### No Parameters

Some methods don't need any input:

```java
public static void printWelcome() {
    System.out.println("Welcome to Java programming!");
}

printWelcome();  // No arguments needed
```

## Return Types and Return Statement

The return type specifies what kind of value a method sends back to its caller.

### Methods That Return Values

```java
public static double calculateArea(double radius) {
    double area = Math.PI * radius * radius;
    return area;  // Return a double value
}

double circleArea = calculateArea(5.0);
System.out.println("Area: " + circleArea);  // Area: 78.53981633974483
```

The return type in the method signature (`double`) must match the type of value you return.

### Multiple Return Statements

A method can have multiple return statements, but only one will execute:

```java
public static String getGrade(int score) {
    if (score >= 90) {
        return "A";
    } else if (score >= 80) {
        return "B";
    } else if (score >= 70) {
        return "C";
    } else if (score >= 60) {
        return "D";
    } else {
        return "F";
    }
}
```

Once a return statement executes, the method immediately exits. Code after a return is unreachable.

### Void Methods

`void` means the method doesn't return anything:

```java
public static void printStars(int count) {
    for (int i = 0; i < count; i++) {
        System.out.print("*");
    }
    System.out.println();  // New line
}

printStars(5);  // Prints: *****
```

Void methods can have a `return;` statement (no value) to exit early:

```java
public static void printPositive(int num) {
    if (num <= 0) {
        return;  // Exit early if not positive
    }
    System.out.println(num);
}
```

## Method Overloading

Method overloading means having multiple methods with the same name but different parameters. Java figures out which one to call based on the arguments you provide.

### Why Overload?

It makes your code more intuitive. For example, `System.out.println()` can print strings, integers, doubles, etc. That's method overloading!

### Overloading Examples

```java
public class MathHelper {
    // Add two integers
    public static int add(int a, int b) {
        return a + b;
    }
    
    // Add three integers
    public static int add(int a, int b, int c) {
        return a + b + c;
    }
    
    // Add two doubles
    public static double add(double a, double b) {
        return a + b;
    }
    
    public static void main(String[] args) {
        System.out.println(add(5, 3));        // Calls first version: 8
        System.out.println(add(5, 3, 2));     // Calls second version: 10
        System.out.println(add(5.5, 3.2));    // Calls third version: 8.7
    }
}
```

Java determines which method to call based on:
1. Number of arguments
2. Types of arguments
3. Order of argument types

### Overloading Rules

Methods must differ in their parameter list:

```java
// VALID - different number of parameters
public static int max(int a, int b) { ... }
public static int max(int a, int b, int c) { ... }

// VALID - different parameter types
public static int max(int a, int b) { ... }
public static double max(double a, double b) { ... }

// INVALID - return type alone doesn't count
public static int max(int a, int b) { ... }
public static double max(int a, int b) { ... }  // ERROR!
```

The return type alone doesn't matter for overloading—the parameters must be different.

## Variable Scope

Scope determines where a variable can be accessed in your code. Understanding scope prevents bugs and makes your code clearer.

### Local Variables (Method Scope)

Variables declared inside a method are local to that method:

```java
public static void myMethod() {
    int x = 10;  // Local variable
    System.out.println(x);  // Works fine
}

public static void anotherMethod() {
    System.out.println(x);  // ERROR! x doesn't exist here
}
```

Local variables:
- Exist only within their method
- Are created when the method is called
- Are destroyed when the method returns
- Must be initialized before use

### Class-Level Variables (Instance and Static)

Variables declared at the class level (outside methods) are accessible to all methods:

```java
public class Student {
    // Class-level variables
    static int studentCount = 0;  // Static variable (shared by all)
    String name;  // Instance variable (each object has its own)
    
    public static void incrementCount() {
        studentCount++;  // Can access static variable
    }
    
    public void displayName() {
        System.out.println(name);  // Can access instance variable
    }
}
```

### Block Scope

Variables declared in blocks (like if, for, while) are only visible in that block:

```java
public static void checkAge(int age) {
    if (age >= 18) {
        String message = "Adult";  // Only exists in this if block
        System.out.println(message);
    }
    // System.out.println(message);  // ERROR! message doesn't exist here
}
```

### Shadowing

A local variable can "shadow" (hide) a class-level variable with the same name:

```java
public class Example {
    static int value = 10;  // Class-level
    
    public static void method() {
        int value = 20;  // Local variable shadows class-level one
        System.out.println(value);  // Prints 20 (local)
        System.out.println(Example.value);  // Prints 10 (class-level)
    }
}
```

**Best practice**: Avoid shadowing—use different names to prevent confusion!

## Method Invocation and Call Stack

When you call a method, Java uses a call stack to keep track of method calls.

### How Method Calls Work

```java
public static void main(String[] args) {
    System.out.println("Start");
    method1();
    System.out.println("End");
}

public static void method1() {
    System.out.println("In method1");
    method2();
    System.out.println("Back in method1");
}

public static void method2() {
    System.out.println("In method2");
}
```

Output:
```
Start
In method1
In method2
Back in method1
End
```

### The Call Stack

The call stack is like a stack of plates. When you call a method, Java puts it on top of the stack. When a method returns, Java removes it from the stack:

1. `main` starts → Stack: [main]
2. `main` calls `method1` → Stack: [main, method1]
3. `method1` calls `method2` → Stack: [main, method1, method2]
4. `method2` returns → Stack: [main, method1]
5. `method1` returns → Stack: [main]
6. `main` ends → Stack: []

### Stack Overflow

If methods call each other too many times (often through recursion without a base case), you get a stack overflow error:

```java
public static void infiniteMethod() {
    infiniteMethod();  // Calls itself forever - BAD!
}
```

Don't do this! Make sure recursive methods have a stopping condition.

## Passing Arguments by Value

Java is a "pass-by-value" language. This means when you pass an argument to a method, Java passes a copy of the value, not the original variable.

### Primitive Types

With primitive types, the method gets a copy of the value:

```java
public static void changeValue(int num) {
    num = 100;  // Changes the copy, not the original
    System.out.println("Inside method: " + num);  // 100
}

public static void main(String[] args) {
    int x = 5;
    changeValue(x);
    System.out.println("In main: " + x);  // Still 5! Original unchanged
}
```

The method can't change the original variable—only its own copy.

### Reference Types (Objects)

With objects, Java still passes by value, but the value is a reference (memory address). So the method can modify the object itself:

```java
public static void modifyArray(int[] arr) {
    arr[0] = 999;  // Modifies the actual array
}

public static void main(String[] args) {
    int[] numbers = {1, 2, 3};
    modifyArray(numbers);
    System.out.println(numbers[0]);  // 999! Array was modified
}
```

Important distinction: You can't change *which* object the variable points to, but you can modify the object's contents:

```java
public static void reassignArray(int[] arr) {
    arr = new int[]{7, 8, 9};  // Creates new array, but...
}

public static void main(String[] args) {
    int[] numbers = {1, 2, 3};
    reassignArray(numbers);
    System.out.println(numbers[0]);  // Still 1! Original array unchanged
}
```

## Recursive Methods

A recursive method is one that calls itself. It's a powerful technique for solving problems that can be broken down into smaller, similar problems.

### Anatomy of Recursion

Every recursive method needs:
1. **Base case**: A condition that stops the recursion
2. **Recursive case**: The method calling itself with a simpler problem

### Classic Example: Factorial

```java
public static int factorial(int n) {
    // Base case: factorial of 0 or 1 is 1
    if (n <= 1) {
        return 1;
    }
    // Recursive case: n! = n * (n-1)!
    return n * factorial(n - 1);
}

System.out.println(factorial(5));  // 120 (5 * 4 * 3 * 2 * 1)
```

How it works for `factorial(5)`:
```
factorial(5) = 5 * factorial(4)
             = 5 * (4 * factorial(3))
             = 5 * (4 * (3 * factorial(2)))
             = 5 * (4 * (3 * (2 * factorial(1))))
             = 5 * (4 * (3 * (2 * 1)))
             = 120
```

### Another Example: Fibonacci

```java
public static int fibonacci(int n) {
    // Base cases
    if (n <= 1) {
        return n;
    }
    // Recursive case
    return fibonacci(n - 1) + fibonacci(n - 2);
}

System.out.println(fibonacci(6));  // 8 (sequence: 0,1,1,2,3,5,8)
```

### When to Use Recursion

Recursion is great for:
- Tree/graph traversal
- Divide and conquer algorithms
- Problems with self-similar structure

But be careful! Recursion can be slower and use more memory than iteration. Always ensure you have a base case to prevent infinite recursion.

## Variable Arguments (Varargs)

Sometimes you want a method to accept any number of arguments. That's what varargs are for!

### Varargs Syntax

Use three dots `...` after the type:

```java
public static int sum(int... numbers) {
    int total = 0;
    for (int num : numbers) {
        total += num;
    }
    return total;
}

// Call with any number of arguments
System.out.println(sum(1, 2, 3));        // 6
System.out.println(sum(1, 2, 3, 4, 5));  // 15
System.out.println(sum(10));             // 10
System.out.println(sum());               // 0
```

Inside the method, varargs are treated as an array.

### Varargs Rules

1. You can have only one varargs parameter per method
2. The varargs parameter must be the last parameter

```java
// VALID
public static void print(String prefix, int... numbers) { ... }

// INVALID - varargs must be last
public static void print(int... numbers, String suffix) { ... }

// INVALID - only one varargs allowed
public static void print(int... nums1, int... nums2) { ... }
```

### Practical Example

```java
public static String buildMessage(String separator, String... words) {
    if (words.length == 0) {
        return "";
    }
    StringBuilder result = new StringBuilder(words[0]);
    for (int i = 1; i < words.length; i++) {
        result.append(separator).append(words[i]);
    }
    return result.toString();
}

System.out.println(buildMessage(" ", "Hello", "world"));  // "Hello world"
System.out.println(buildMessage("-", "one", "two", "three"));  // "one-two-three"
```

## Method Naming Conventions

Good method names make your code self-documenting and easy to understand.

### Naming Rules

Follow these conventions:
1. **Use camelCase**: Start with lowercase, capitalize each new word
2. **Use verbs**: Methods perform actions, so use action words
3. **Be descriptive**: Name should explain what the method does
4. **Keep it concise**: But not at the cost of clarity

### Good Examples

```java
// Good names - clear and descriptive
public static int calculateTotal(int[] prices) { ... }
public static boolean isEven(int number) { ... }
public static void printReceipt(Order order) { ... }
public static String getUserName() { ... }
public static void setAge(int age) { ... }
```

### Common Patterns

- **Getters**: `getName()`, `getAge()`, `isActive()`
- **Setters**: `setName()`, `setAge()`, `setActive()`
- **Boolean methods**: `isEmpty()`, `hasPermission()`, `canVote()`
- **Action methods**: `calculate()`, `process()`, `save()`, `load()`

### Poor Examples (Avoid These)

```java
// BAD - not descriptive
public static void doStuff() { ... }
public static int x(int y) { ... }

// BAD - unclear abbreviations
public static void prtRcpt() { ... }

// BAD - doesn't follow camelCase
public static void Calculate_Total() { ... }
public static void PRINT_MESSAGE() { ... }
```

## Real-World Example: Temperature Converter

Let's put it all together with a practical example:

```java
public class TemperatureConverter {
    // Convert Celsius to Fahrenheit
    public static double celsiusToFahrenheit(double celsius) {
        return (celsius * 9.0 / 5.0) + 32;
    }
    
    // Convert Fahrenheit to Celsius
    public static double fahrenheitToCelsius(double fahrenheit) {
        return (fahrenheit - 32) * 5.0 / 9.0;
    }
    
    // Overloaded method - convert and display
    public static void convertAndDisplay(double temp, String fromUnit) {
        if (fromUnit.equals("C")) {
            double fahrenheit = celsiusToFahrenheit(temp);
            System.out.println(temp + "°C = " + fahrenheit + "°F");
        } else if (fromUnit.equals("F")) {
            double celsius = fahrenheitToCelsius(temp);
            System.out.println(temp + "°F = " + celsius + "°C");
        }
    }
    
    // Varargs method - convert multiple temperatures
    public static void convertMultiple(String unit, double... temperatures) {
        for (double temp : temperatures) {
            convertAndDisplay(temp, unit);
        }
    }
    
    public static void main(String[] args) {
        // Direct conversion
        double f = celsiusToFahrenheit(100);
        System.out.println("100°C = " + f + "°F");  // 212.0°F
        
        // Using convenience method
        convertAndDisplay(32, "F");  // 32°F = 0°C
        
        // Convert multiple
        convertMultiple("C", 0, 25, 100);
    }
}
```

---

## Quiz Questions

### Question 1: Basic Method Syntax
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: java, methods, syntax, declaration

**Question Text**:
What is the correct way to declare a method that takes two integers and returns their sum?

**Answer Choices**:
- A) `int sum(int a, int b) { return a + b; }`
- B) `sum(int a, int b) { return a + b; }`
- C) `int sum(a, b) { return a + b; }`
- D) `return int sum(int a, int b) { a + b; }`

**Correct Answer**: A

**Explanation**:
A method declaration needs a return type (int), method name (sum), parameter list with types (int a, int b), and a return statement. Option A has all required components in the correct order. Option B is missing the return type, option C doesn't specify parameter types, and option D has incorrect syntax with "return" before the method declaration.

---

### Question 2: Void Methods
**Type**: true-false
**Difficulty**: 2
**Tags**: java, methods, void, return

**Question Text**:
A method with return type `void` cannot contain any `return` statement.

**Answer Choices**:
- A) True
- B) False

**Correct Answer**: B

**Explanation**:
False. A void method can contain a `return;` statement (without a value) to exit the method early. For example: `public static void checkAge(int age) { if (age < 0) return; System.out.println(age); }` uses return to exit early for invalid input. Void methods just cannot return a value.

---

### Question 3: Method Overloading
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, methods, overloading

**Question Text**:
Which of the following is a valid method overload for `public static int add(int a, int b)`?

**Answer Choices**:
- A) `public static double add(int a, int b)`
- B) `public static int add(int x, int y)`
- C) `public static int add(double a, double b)`
- D) `public static int add(int a, int b, int c)`

**Correct Answer**: C

**Explanation**:
Method overloading requires different parameter lists (number, types, or order of parameters). Option C changes parameter types from int to double, making it a valid overload. Option A only changes return type (invalid), option B only changes parameter names (invalid), and option D changes the number of parameters (valid, but C is also correct). Both C and D would be valid, but C is the best single answer.

---

### Question 4: Variable Scope
**Type**: code-output
**Difficulty**: 3
**Tags**: java, methods, scope, variables

**Question Text**:
What happens when this code runs?

```java
public static void test() {
    int x = 10;
}

public static void main(String[] args) {
    test();
    System.out.println(x);
}
```

**Answer Choices**:
- A) Prints 10
- B) Prints 0
- C) Compilation error
- D) Runtime error

**Correct Answer**: C

**Explanation**:
This code will not compile. The variable `x` is declared inside the `test()` method, making it local to that method. Local variables cannot be accessed outside their method. The `main` method cannot see `x`, resulting in a compilation error: "cannot find symbol x". To fix this, x would need to be declared at the class level or returned from the test method.

---

### Question 5: Pass-by-Value with Primitives
**Type**: code-output
**Difficulty**: 4
**Tags**: java, methods, pass-by-value, parameters

**Question Text**:
What does this code print?

```java
public static void modify(int num) {
    num = 100;
}

public static void main(String[] args) {
    int value = 50;
    modify(value);
    System.out.println(value);
}
```

**Answer Choices**:
- A) 100
- B) 50
- C) 150
- D) 0

**Correct Answer**: B

**Explanation**:
Java passes primitive types by value, meaning the method receives a copy of the value. When `modify` changes `num` to 100, it only changes its local copy, not the original `value` variable in main. So `value` remains 50 and that's what gets printed. To actually change the value, the method would need to return the new value and you'd assign it: `value = modify(value);`.

---

### Question 6: Recursion Base Case
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, methods, recursion, base-case

**Question Text**:
What is the purpose of a base case in a recursive method?

**Answer Choices**:
- A) To make the method run faster
- B) To prevent infinite recursion and provide a stopping condition
- C) To handle the first call to the method
- D) To initialize variables used in the method

**Correct Answer**: B

**Explanation**:
A base case is essential in recursive methods to prevent infinite recursion. It provides a condition where the method stops calling itself and returns a value directly. Without a base case, the method would call itself forever, eventually causing a StackOverflowError. For example, in calculating factorial, the base case is `if (n <= 1) return 1;` which stops the recursion.

---

### Question 7: Varargs Parameter
**Type**: code-output
**Difficulty**: 3
**Tags**: java, methods, varargs, parameters

**Question Text**:
What does this code output?

```java
public static int sum(int... numbers) {
    int total = 0;
    for (int n : numbers) {
        total += n;
    }
    return total;
}

public static void main(String[] args) {
    System.out.println(sum(1, 2, 3, 4));
}
```

**Answer Choices**:
- A) 4
- B) 10
- C) 24
- D) Compilation error

**Correct Answer**: B

**Explanation**:
The varargs parameter `int... numbers` allows the method to accept any number of integer arguments. These are treated as an array inside the method. The sum method adds all numbers: 1 + 2 + 3 + 4 = 10. Varargs are flexible - you could call this with `sum(1, 2)` or `sum(1, 2, 3, 4, 5)` or even `sum()` (which would return 0).

---

### Question 8: Method Return Flow
**Type**: code-output
**Difficulty**: 5
**Tags**: java, methods, return, control-flow

**Question Text**:
What is printed by this code?

```java
public static String checkNumber(int num) {
    if (num > 0) {
        return "Positive";
    } else if (num < 0) {
        return "Negative";
    }
    System.out.println("This line");
    return "Zero";
}

public static void main(String[] args) {
    System.out.println(checkNumber(0));
}
```

**Answer Choices**:
- A) Zero
- B) This line
- C) This line, then Zero on a new line
- D) Positive

**Correct Answer**: C

**Explanation**:
When `checkNumber(0)` is called, the value is neither greater than nor less than 0, so neither if condition is true. The code then prints "This line" and returns "Zero". In main, that returned value "Zero" is printed. So the output is two lines: "This line" followed by "Zero". If the number were positive or negative, the method would return early and never reach the println statement.
