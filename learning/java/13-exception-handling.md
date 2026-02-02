# Java Exception Handling

Hey there! Let's dive into one of the most important topics in Java—exception handling. You know how things can go wrong in real life? Maybe you try to open a file that doesn't exist, or divide a number by zero. In programming, these situations are called exceptions, and Java gives us powerful tools to handle them gracefully instead of letting our program crash and burn.

## 13.1 Understanding Exceptions and Errors

Think of an exception as Java's way of saying "Houston, we have a problem!" When something unexpected happens during program execution, Java creates an exception object and throws it. If we don't catch it, our program terminates abruptly—not a great user experience!

Exceptions can happen for all sorts of reasons: trying to read a file that doesn't exist, attempting to access an array index that's out of bounds, or dividing by zero. The good news? Java lets us anticipate these problems and handle them elegantly.

Here's what happens when an exception occurs without handling:

```java
public class ExceptionDemo {
    public static void main(String[] args) {
        int[] numbers = {1, 2, 3};
        System.out.println(numbers[5]);  // Oops! Index 5 doesn't exist
        System.out.println("This line never executes");
    }
}
```

This code crashes with an `ArrayIndexOutOfBoundsException`, and the last print statement never runs. Exception handling lets us prevent this disaster!

Errors, on the other hand, are serious problems that applications typically shouldn't try to catch—things like running out of memory or stack overflow. We'll focus mainly on exceptions, which are problems we can and should handle.

## 13.2 Exception Hierarchy

Java organizes exceptions in a family tree, with `Throwable` at the top. Understanding this hierarchy helps you know which exceptions to catch and when.

The hierarchy looks like this:
- **Throwable** (the grandparent of all exceptions)
  - **Error** (serious problems you usually don't handle)
    - OutOfMemoryError, StackOverflowError, etc.
  - **Exception** (problems you can handle)
    - **IOException**, **SQLException**, etc. (checked exceptions)
    - **RuntimeException** (unchecked exceptions)
      - NullPointerException, ArrayIndexOutOfBoundsException, etc.

The key distinction is between **checked** and **unchecked** exceptions. Checked exceptions (like IOException) must be either caught or declared in your method signature—Java forces you to handle them. Unchecked exceptions (RuntimeException and its subclasses) don't have this requirement. They typically represent programming errors that you could have prevented with better code.

```java
// Checked exception - must be handled
import java.io.*;

public class CheckedExample {
    // This won't compile without handling IOException
    public static void readFile() throws IOException {
        FileReader file = new FileReader("data.txt");
    }
}

// Unchecked exception - no handling required
public class UncheckedExample {
    public static void divide() {
        int result = 10 / 0;  // ArithmeticException (unchecked)
    }
}
```

Think of checked exceptions as Java being your protective parent: "You can't do that risky thing unless you promise to handle what might go wrong!" Unchecked exceptions are more like "Hey, that's a bug in your code—fix it!"

## 13.3 Try-Catch Blocks

The try-catch block is your primary weapon against exceptions. You "try" some risky code, and if an exception occurs, you "catch" it and handle it gracefully.

```java
public class TryCatchDemo {
    public static void main(String[] args) {
        try {
            int[] numbers = {1, 2, 3};
            System.out.println("Accessing element...");
            System.out.println(numbers[5]);  // This will throw an exception
            System.out.println("This won't print");
        } catch (ArrayIndexOutOfBoundsException e) {
            System.out.println("Oops! Invalid index: " + e.getMessage());
        }
        
        System.out.println("Program continues running!");
    }
}
```

When the exception occurs inside the try block, Java immediately jumps to the catch block. The code after the problematic line in the try block gets skipped, but the program continues after the catch block—no crash!

The catch block receives the exception object (we called it `e`), which contains useful information like the error message and stack trace. You can print it, log it, or use it to decide what to do next.

```java
try {
    String text = null;
    System.out.println(text.length());  // NullPointerException!
} catch (NullPointerException e) {
    System.out.println("Error: " + e.getMessage());
    e.printStackTrace();  // Prints detailed error information
}
```

The beauty of try-catch is that it separates your normal code flow from your error handling code, making everything cleaner and more readable.

## 13.4 Multiple Catch Blocks

Sometimes different exceptions need different handling strategies. You can have multiple catch blocks, each handling a specific type of exception.

```java
public class MultipleCatchDemo {
    public static void main(String[] args) {
        try {
            String str = args[0];  // Might throw ArrayIndexOutOfBoundsException
            int num = Integer.parseInt(str);  // Might throw NumberFormatException
            int result = 100 / num;  // Might throw ArithmeticException
            System.out.println("Result: " + result);
        } catch (ArrayIndexOutOfBoundsException e) {
            System.out.println("Please provide a number as argument!");
        } catch (NumberFormatException e) {
            System.out.println("That's not a valid number!");
        } catch (ArithmeticException e) {
            System.out.println("Cannot divide by zero!");
        }
    }
}
```

Java checks each catch block from top to bottom and executes the first one that matches. Important: you must order catch blocks from most specific to most general exception types.

Since Java 7, you can also catch multiple exception types in a single catch block using the pipe (|) operator:

```java
try {
    // Some risky code
    performOperation();
} catch (IOException | SQLException e) {
    System.out.println("Database or file error: " + e.getMessage());
    logError(e);
}
```

This is handy when you want to handle different exceptions the same way, keeping your code DRY (Don't Repeat Yourself).

## 13.5 The Finally Block

The `finally` block is like the friend who always shows up, no matter what. It executes whether an exception occurs or not, making it perfect for cleanup tasks like closing files or database connections.

```java
import java.io.*;

public class FinallyDemo {
    public static void main(String[] args) {
        FileReader reader = null;
        try {
            reader = new FileReader("data.txt");
            // Read file content
            int data = reader.read();
            System.out.println("First character: " + (char) data);
        } catch (IOException e) {
            System.out.println("Error reading file: " + e.getMessage());
        } finally {
            System.out.println("Finally block always executes!");
            // Clean up resources
            try {
                if (reader != null) {
                    reader.close();
                }
            } catch (IOException e) {
                System.out.println("Error closing file");
            }
        }
    }
}
```

The finally block runs in all these scenarios:
- No exception occurs (normal execution)
- An exception occurs and is caught
- An exception occurs and is NOT caught (still runs before the program terminates)
- Even if there's a return statement in the try or catch block!

This makes finally perfect for releasing resources like file handles, database connections, or network sockets. Even if something goes wrong, you know your cleanup code will run.

## 13.6 Try-With-Resources Statement

Closing resources in a finally block works, but it's verbose and error-prone. Java 7 introduced try-with-resources, which automatically closes resources for you—much cleaner!

```java
import java.io.*;

public class TryWithResourcesDemo {
    public static void main(String[] args) {
        // Resources declared here are auto-closed
        try (FileReader reader = new FileReader("data.txt");
             BufferedReader br = new BufferedReader(reader)) {
            
            String line = br.readLine();
            System.out.println("First line: " + line);
            
        } catch (IOException e) {
            System.out.println("Error: " + e.getMessage());
        }
        // No need for finally block - reader and br are automatically closed!
    }
}
```

Any resource used in a try-with-resources statement must implement the `AutoCloseable` interface. Most I/O classes (FileReader, BufferedReader, Scanner, etc.) already implement this, so you can use them directly.

You can declare multiple resources separated by semicolons. They're closed in the reverse order of their creation (last declared, first closed).

```java
try (Scanner scanner = new Scanner(new File("input.txt"));
     PrintWriter writer = new PrintWriter("output.txt")) {
    
    while (scanner.hasNextLine()) {
        String line = scanner.nextLine();
        writer.println(line.toUpperCase());
    }
} catch (IOException e) {
    e.printStackTrace();
}
// scanner and writer are both automatically closed
```

This is not just convenient—it's also safer! If an exception occurs while closing a resource in a finally block, it can mask the original exception. Try-with-resources handles this correctly by preserving both exceptions.

## 13.7 Throwing Exceptions with the throw Keyword

Sometimes you need to throw your own exceptions. Maybe you're validating input and find something invalid, or you want to signal that a method can't complete its task. That's where the `throw` keyword comes in.

```java
public class ThrowDemo {
    public static void validateAge(int age) {
        if (age < 0) {
            throw new IllegalArgumentException("Age cannot be negative!");
        }
        if (age < 18) {
            throw new IllegalArgumentException("Must be 18 or older!");
        }
        System.out.println("Age " + age + " is valid");
    }
    
    public static void main(String[] args) {
        try {
            validateAge(15);
        } catch (IllegalArgumentException e) {
            System.out.println("Validation failed: " + e.getMessage());
        }
    }
}
```

You can throw any exception object. Common ones to throw include:
- `IllegalArgumentException` - for invalid method parameters
- `IllegalStateException` - when an object is in the wrong state
- `NullPointerException` - for unexpected null values (though preventing nulls is better)
- `UnsupportedOperationException` - for operations that aren't supported

Here's a practical example with a bank account:

```java
public class BankAccount {
    private double balance;
    
    public void withdraw(double amount) {
        if (amount < 0) {
            throw new IllegalArgumentException("Amount must be positive");
        }
        if (amount > balance) {
            throw new IllegalStateException("Insufficient funds! Balance: " + balance);
        }
        balance -= amount;
        System.out.println("Withdrew $" + amount);
    }
}
```

When you throw an exception, the current method stops executing immediately, and the exception propagates up the call stack until someone catches it.

## 13.8 Declaring Exceptions with the throws Keyword

When your method might throw a checked exception but you don't want to handle it yourself, you can declare it using the `throws` keyword in your method signature. This passes the responsibility to the caller.

```java
import java.io.*;

public class ThrowsDemo {
    // This method declares that it might throw IOException
    public static String readFirstLine(String filename) throws IOException {
        FileReader file = new FileReader(filename);
        BufferedReader reader = new BufferedReader(file);
        String line = reader.readLine();
        reader.close();
        return line;
    }
    
    public static void main(String[] args) {
        // Now the caller must handle the exception
        try {
            String line = readFirstLine("data.txt");
            System.out.println(line);
        } catch (IOException e) {
            System.out.println("Could not read file: " + e.getMessage());
        }
    }
}
```

You can declare multiple exceptions separated by commas:

```java
public void complexOperation() throws IOException, SQLException, InterruptedException {
    // Code that might throw any of these exceptions
}
```

Important distinction:
- `throw` is used to actually throw an exception object
- `throws` is used to declare that a method might throw an exception

```java
// Using throw (throws an exception)
if (value < 0) {
    throw new IllegalArgumentException("Negative value!");
}

// Using throws (declares possibility)
public void readFile(String path) throws IOException {
    // Implementation
}
```

For unchecked exceptions (RuntimeException subclasses), you don't have to declare them with throws, but you can if you want to document them. For checked exceptions, you must either catch them or declare them—there's no escape!

## 13.9 Creating Custom Exceptions

Sometimes the built-in exceptions don't quite fit your needs. That's when you create your own custom exceptions. This makes your code more expressive and helps users of your API understand what went wrong.

```java
// Custom checked exception
public class InsufficientFundsException extends Exception {
    private double amount;
    
    public InsufficientFundsException(double amount) {
        super("Insufficient funds: requested $" + amount);
        this.amount = amount;
    }
    
    public double getAmount() {
        return amount;
    }
}

// Custom unchecked exception
public class InvalidAccountException extends RuntimeException {
    public InvalidAccountException(String message) {
        super(message);
    }
}
```

To create a custom exception:
1. Extend `Exception` for a checked exception, or `RuntimeException` for an unchecked one
2. Provide constructors that call the parent constructor using `super()`
3. Optionally add fields and methods for additional context

Here's how you'd use these custom exceptions:

```java
public class BankAccount {
    private String accountNumber;
    private double balance;
    
    public void withdraw(double amount) throws InsufficientFundsException {
        if (accountNumber == null) {
            throw new InvalidAccountException("Account not initialized");
        }
        if (amount > balance) {
            throw new InsufficientFundsException(amount);
        }
        balance -= amount;
    }
    
    public static void main(String[] args) {
        BankAccount account = new BankAccount();
        account.balance = 1000;
        account.accountNumber = "12345";
        
        try {
            account.withdraw(1500);
        } catch (InsufficientFundsException e) {
            System.out.println("Transaction failed: " + e.getMessage());
            System.out.println("Shortfall: $" + (e.getAmount() - 1000));
        }
    }
}
```

Custom exceptions make your code self-documenting. When someone sees `InsufficientFundsException`, they immediately understand what went wrong, without having to decipher a generic exception message.

## 13.10 Best Practices in Exception Handling

Let's talk about how to use exceptions effectively. Following these best practices will make your code more robust and maintainable.

**1. Catch Specific Exceptions**: Don't catch `Exception` or `Throwable` unless you have a very good reason. Catch the most specific exception you can handle.

```java
// Bad - too broad
try {
    processFile();
} catch (Exception e) {
    // What went wrong? File not found? Parse error? Who knows!
}

// Good - specific
try {
    processFile();
} catch (FileNotFoundException e) {
    System.out.println("File not found: " + e.getMessage());
} catch (IOException e) {
    System.out.println("Error reading file: " + e.getMessage());
}
```

**2. Don't Swallow Exceptions**: Always do something meaningful in your catch block. An empty catch block hides problems!

```java
// Bad - silently swallowing the exception
try {
    riskyOperation();
} catch (Exception e) {
    // Nothing here - error disappears into the void!
}

// Good - at least log it
try {
    riskyOperation();
} catch (Exception e) {
    System.err.println("Operation failed: " + e.getMessage());
    e.printStackTrace();
}
```

**3. Use Exceptions for Exceptional Situations**: Don't use exceptions for normal control flow. They're slower than regular if-else statements and make code harder to understand.

```java
// Bad - using exceptions for flow control
try {
    int index = 0;
    while (true) {
        System.out.println(array[index++]);
    }
} catch (ArrayIndexOutOfBoundsException e) {
    // End of array reached
}

// Good - use normal loop control
for (int i = 0; i < array.length; i++) {
    System.out.println(array[i]);
}
```

**4. Clean Up Resources**: Always close resources (files, connections, streams) in a finally block or use try-with-resources.

**5. Provide Context in Custom Exceptions**: When creating custom exceptions, include relevant information that helps diagnose the problem.

```java
public class DataValidationException extends Exception {
    private String fieldName;
    private Object invalidValue;
    
    public DataValidationException(String fieldName, Object invalidValue, String message) {
        super(message);
        this.fieldName = fieldName;
        this.invalidValue = invalidValue;
    }
    
    // Getters for additional context
}
```

**6. Document Exceptions**: Use Javadoc @throws tags to document which exceptions your method might throw and under what conditions.

```java
/**
 * Withdraws money from the account.
 * 
 * @param amount the amount to withdraw
 * @throws InsufficientFundsException if balance is less than amount
 * @throws IllegalArgumentException if amount is negative
 */
public void withdraw(double amount) throws InsufficientFundsException {
    // Implementation
}
```

## 13.11 Common Exceptions

Let's tour the exceptions you'll encounter most frequently in Java. Knowing these helps you understand what went wrong and how to prevent it.

**NullPointerException**: The infamous NPE! Thrown when you try to use an object reference that's null.

```java
String text = null;
System.out.println(text.length());  // NullPointerException!

// Prevention:
if (text != null) {
    System.out.println(text.length());
}
// Or use Java 14+ null-safe methods
```

**ArrayIndexOutOfBoundsException**: You tried to access an array index that doesn't exist.

```java
int[] numbers = {1, 2, 3};
System.out.println(numbers[5]);  // Only indices 0-2 exist!

// Prevention: check bounds first
if (index >= 0 && index < numbers.length) {
    System.out.println(numbers[index]);
}
```

**NumberFormatException**: Thrown when trying to convert a string to a number, but the string isn't a valid number.

```java
int num = Integer.parseInt("abc");  // Not a number!

// Prevention: validate or catch
try {
    int num = Integer.parseInt(userInput);
} catch (NumberFormatException e) {
    System.out.println("Please enter a valid number");
}
```

**IllegalArgumentException**: A method received an inappropriate argument.

```java
public void setAge(int age) {
    if (age < 0 || age > 150) {
        throw new IllegalArgumentException("Invalid age: " + age);
    }
    this.age = age;
}
```

**IOException**: Signals a failed or interrupted I/O operation. Common subclasses include FileNotFoundException.

```java
try (FileReader file = new FileReader("data.txt")) {
    // Read file
} catch (FileNotFoundException e) {
    System.out.println("File not found!");
} catch (IOException e) {
    System.out.println("Error reading file!");
}
```

**ArithmeticException**: Math went wrong, typically division by zero.

```java
int result = 10 / 0;  // ArithmeticException

// Prevention:
if (divisor != 0) {
    int result = dividend / divisor;
}
```

**ClassCastException**: You tried to cast an object to a type it's not compatible with.

```java
Object obj = "Hello";
Integer num = (Integer) obj;  // Can't cast String to Integer!

// Prevention: check type first
if (obj instanceof Integer) {
    Integer num = (Integer) obj;
}
```

Understanding these common exceptions helps you write defensive code that anticipates problems and handles them gracefully.

## Wrapping Up

Exception handling is your safety net when coding. It lets you write robust programs that handle the unexpected without crashing. Remember:

- Use try-catch to handle exceptions gracefully
- Let checked exceptions guide you to handle serious problems
- Create custom exceptions to make your code more expressive
- Follow best practices: be specific, don't swallow exceptions, and always clean up resources
- Try-with-resources is your friend for automatic resource management

Exceptions might seem annoying at first ("Why can't my code just work?"), but they're actually doing you a favor by forcing you to think about what can go wrong and how to handle it. Master exception handling, and you'll write more reliable, production-ready code.

Ready to test your exception-handling skills? Time for the quiz!

---

## Quiz Questions

This section contains quiz questions that test understanding of Java exception handling. Each question is formatted for easy extraction and database import.

### Question 1: Exception Hierarchy - Checked vs Unchecked
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, exception-handling, exception-hierarchy, checked-exceptions

**Question Text**:
Which of the following is true about checked exceptions in Java?

**Answer Choices**:
- A) They extend RuntimeException and don't need to be declared
- B) They must be either caught or declared using the throws keyword
- C) They represent programming errors that should be fixed in code
- D) They cannot be caught using try-catch blocks

**Correct Answer**: B

**Explanation**:
Checked exceptions must be either caught in a try-catch block or declared in the method signature using the throws keyword—Java enforces this at compile time. Option A describes unchecked exceptions (which extend RuntimeException). Option C describes the typical use case for unchecked exceptions. Option D is false—checked exceptions can and should be caught when appropriate.

---

### Question 2: Try-Catch Blocks - Basic Usage
**Type**: code-output
**Difficulty**: 2
**Tags**: java, exception-handling, try-catch

**Question Text**:
What will be printed when this code is executed?

```java
try {
    int[] numbers = {1, 2, 3};
    System.out.println("A");
    System.out.println(numbers[5]);
    System.out.println("B");
} catch (ArrayIndexOutOfBoundsException e) {
    System.out.println("C");
}
System.out.println("D");
```

**Answer Choices**:
- A) A B C D
- B) A C D
- C) A D
- D) A B D

**Correct Answer**: B

**Explanation**:
The code first prints "A", then attempts to access index 5 which doesn't exist, throwing an ArrayIndexOutOfBoundsException. The exception is caught, so "C" is printed. The line printing "B" is never reached because the exception occurred before it. After the catch block, execution continues normally and "D" is printed. Output: A, C, D.

---

### Question 3: Finally Block - Execution Guarantee
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, exception-handling, finally

**Question Text**:
When does the finally block NOT execute?

**Answer Choices**:
- A) When an exception is thrown but not caught
- B) When there's a return statement in the try block
- C) When System.exit() is called
- D) When no exception occurs

**Correct Answer**: C

**Explanation**:
The finally block executes in almost all scenarios: when no exception occurs, when an exception is caught, when an exception is not caught, and even when there's a return statement. The only situations where finally doesn't execute are catastrophic events like System.exit() being called (which terminates the JVM), power failure, or infinite loops. Option C is the only scenario listed where finally won't execute.

---

### Question 4: Try-With-Resources - Automatic Closing
**Type**: true-false
**Difficulty**: 3
**Tags**: java, exception-handling, try-with-resources

**Question Text**:
True or False: In a try-with-resources statement, resources are automatically closed even if an exception occurs within the try block.

**Correct Answer**: True

**Explanation**:
This is true and is one of the main benefits of try-with-resources. Any resource declared in the try-with-resources statement (resources that implement AutoCloseable) will be automatically closed at the end of the try block, whether an exception occurs or not. The resources are closed in the reverse order of their declaration. This eliminates the need for explicit finally blocks to close resources and prevents resource leaks.

---

### Question 5: Throw vs Throws - Understanding the Difference
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, exception-handling, throw, throws

**Question Text**:
What is the difference between the `throw` and `throws` keywords?

**Answer Choices**:
- A) throw is used in method signatures, throws is used to throw exceptions
- B) throw is used to actually throw an exception, throws is used to declare that a method might throw an exception
- C) They are interchangeable and mean the same thing
- D) throw is for checked exceptions, throws is for unchecked exceptions

**Correct Answer**: B

**Explanation**:
The `throw` keyword is used to actually throw an exception object (e.g., `throw new IOException()`), while `throws` is used in a method signature to declare that the method might throw one or more exceptions (e.g., `public void readFile() throws IOException`). Option A has them backwards. Option C is incorrect—they have different purposes. Option D is false—both can be used with either type of exception.

---

### Question 6: Multiple Catch Blocks - Order Matters
**Type**: code-output
**Difficulty**: 5
**Tags**: java, exception-handling, multiple-catch

**Question Text**:
What happens when you compile and run this code?

```java
try {
    String text = null;
    System.out.println(text.length());
} catch (Exception e) {
    System.out.println("Exception caught");
} catch (NullPointerException e) {
    System.out.println("Null pointer caught");
}
```

**Answer Choices**:
- A) Prints "Exception caught"
- B) Prints "Null pointer caught"
- C) Compilation error: unreachable catch block
- D) Both messages are printed

**Correct Answer**: C

**Explanation**:
This code won't compile because the second catch block is unreachable. Since NullPointerException is a subclass of Exception, it would already be caught by the first catch block. Java requires you to order catch blocks from most specific to most general. To fix this, the NullPointerException catch should come before the Exception catch. This is a compile-time error, not a runtime error.

---

### Question 7: Custom Exceptions - Creating Your Own
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, exception-handling, custom-exceptions

**Question Text**:
To create a custom checked exception, your exception class should extend which class?

**Answer Choices**:
- A) RuntimeException
- B) Exception
- C) Throwable
- D) Error

**Correct Answer**: B

**Explanation**:
To create a custom checked exception, you should extend the Exception class (or one of its subclasses, but not RuntimeException). Extending RuntimeException (option A) would create an unchecked exception. While you could technically extend Throwable (option C), this is not recommended—extend Exception instead. Extending Error (option D) is for serious system errors and should not be used for custom application exceptions.

---

### Question 8: Best Practices - Exception Handling Patterns
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, exception-handling, best-practices

**Question Text**:
Which of the following is considered a BAD practice in exception handling?

**Answer Choices**:
- A) Using try-with-resources to automatically close resources
- B) Catching specific exceptions rather than catching Exception
- C) Using exceptions to control normal program flow (like ending a loop)
- D) Providing detailed error messages in custom exceptions

**Correct Answer**: C

**Explanation**:
Using exceptions for normal control flow (option C) is considered bad practice. Exceptions should be used for exceptional situations, not routine program logic. They are slower than normal control structures and make code harder to understand and maintain. Options A, B, and D are all good practices: try-with-resources prevents resource leaks, catching specific exceptions makes error handling more precise, and detailed error messages help with debugging.

---

### Question 9: Common Exceptions - NullPointerException
**Type**: code-output
**Difficulty**: 2
**Tags**: java, exception-handling, common-exceptions, nullpointerexception

**Question Text**:
What will be the output of this code?

```java
String text = null;
try {
    int length = text.length();
    System.out.println("Length: " + length);
} catch (NullPointerException e) {
    System.out.println("Null detected!");
} finally {
    System.out.println("Done");
}
```

**Answer Choices**:
- A) Length: 0<br>Done
- B) Null detected!<br>Done
- C) Null detected!
- D) Program crashes

**Correct Answer**: B

**Explanation**:
When text.length() is called on a null reference, it throws a NullPointerException. This exception is caught by the catch block, which prints "Null detected!". Then, the finally block executes and prints "Done". Both messages are printed. The finally block always executes whether an exception occurs or not, so "Done" will be printed in either case.

---

### Question 10: Try-With-Resources - Resource Declaration
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, exception-handling, try-with-resources, autocloseable

**Question Text**:
What requirement must a class meet to be used in a try-with-resources statement?

**Answer Choices**:
- A) It must extend the Resource class
- B) It must implement the AutoCloseable or Closeable interface
- C) It must be declared as final
- D) It must have a public close() method with no parameters

**Correct Answer**: B

**Explanation**:
To be used in a try-with-resources statement, a class must implement either the AutoCloseable interface or the Closeable interface (Closeable extends AutoCloseable). This ensures the resource has a close() method that can be called automatically. There is no Resource class to extend (option A). The class doesn't need to be final (option C). While option D is partially correct (the close() method is required), simply having the method isn't enough—the class must formally implement one of these interfaces.

---

### Question 11: Exception Propagation - Call Stack
**Type**: multiple-choice
**Difficulty**: 6
**Tags**: java, exception-handling, exception-propagation, call-stack

**Question Text**:
Consider this code:

```java
public static void main(String[] args) {
    try {
        methodA();
        System.out.println("Back in main");
    } catch (Exception e) {
        System.out.println("Caught in main");
    }
}

public static void methodA() throws Exception {
    methodB();
    System.out.println("Back in methodA");
}

public static void methodB() throws Exception {
    throw new Exception("Error!");
}
```

What will be printed?

**Answer Choices**:
- A) Caught in main<br>Back in methodA<br>Back in main
- B) Back in methodA<br>Back in main
- C) Caught in main
- D) Error! (program crashes)

**Correct Answer**: C

**Explanation**:
When methodB() throws the exception, it immediately propagates up to methodA(). Since methodA() doesn't catch it, the exception continues propagating to main(). The line "Back in methodA" never executes because the exception occurred in methodB() and was never caught there. In main(), the exception is caught and "Caught in main" is printed. The line "Back in main" after methodA() never executes because the exception was thrown before methodA() could return normally. Only "Caught in main" is printed.

---

### Question 12: Finally vs Try-With-Resources
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, exception-handling, finally, try-with-resources

**Question Text**:
What is an advantage of using try-with-resources over manually closing resources in a finally block?

**Answer Choices**:
- A) Try-with-resources is faster at runtime
- B) Try-with-resources properly handles exceptions thrown during resource closing
- C) Try-with-resources can be used with any type of object
- D) Try-with-resources doesn't require a catch block

**Correct Answer**: B

**Explanation**:
The key advantage of try-with-resources is that it properly handles "suppressed exceptions"—exceptions that occur while closing resources. If both the try block and the close() method throw exceptions, try-with-resources preserves both (the close exception becomes a suppressed exception), while a manual finally block might lose the original exception. Option A is false—performance is similar. Option C is false—only AutoCloseable resources work. Option D is misleading—you can omit catch with either approach (try-finally and try-with-resources both allow omitting catch, but throws clause is still needed for checked exceptions).
