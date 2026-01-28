# Java Variables and Data Types

Welcome! Let's dive into one of the most fundamental concepts in Java: variables and data types. Think of variables as labeled boxes where you store information, and data types as the rules about what kind of information each box can hold.

## What Are Variables?

A variable is simply a named location in memory that holds a value. It's like a container with a label on it. You can put something in the container, look at what's inside, or replace it with something else.

```java
int age = 25;
String name = "Alex";
double price = 19.99;
```

In these examples:
- `age` is a variable holding the number 25
- `name` is a variable holding the text "Alex"
- `price` is a variable holding the decimal number 19.99

## Primitive Data Types

Java has 8 primitive data types. These are the basic building blocks for storing simple values. Let's go through each one.

### Integer Types

These store whole numbers (no decimals).

#### `byte`
- Size: 8 bits (1 byte)
- Range: -128 to 127
- Use when: You need to save memory and know the number will be small

```java
byte age = 42;
byte temperature = -10;
```

#### `short`
- Size: 16 bits (2 bytes)
- Range: -32,768 to 32,767
- Use when: You need more range than byte but want to save memory

```java
short year = 2024;
short altitude = -1500;
```

#### `int`
- Size: 32 bits (4 bytes)
- Range: -2,147,483,648 to 2,147,483,647
- Use when: This is your go-to for most integer needs!

```java
int population = 1000000;
int accountBalance = -500;
```

**Pro tip:** When in doubt, use `int` for whole numbers. It's the most commonly used integer type.

#### `long`
- Size: 64 bits (8 bytes)
- Range: Really, really big (±9 quintillion)
- Use when: You need to store numbers bigger than `int` can handle

```java
long worldPopulation = 7800000000L;  // Notice the 'L' at the end
long milliseconds = 1234567890123L;
```

**Note:** Always add an `L` or `l` at the end of long literals (use uppercase `L` to avoid confusion with the number 1).

### Floating-Point Types

These store decimal numbers.

#### `float`
- Size: 32 bits (4 bytes)
- Precision: About 6-7 decimal digits
- Use when: You need decimals but don't need high precision

```java
float pi = 3.14f;  // Notice the 'f' at the end
float price = 19.99f;
```

**Note:** Always add an `f` or `F` at the end of float literals.

#### `double`
- Size: 64 bits (8 bytes)
- Precision: About 15 decimal digits
- Use when: This is your default for decimal numbers!

```java
double precise = 3.14159265359;
double distance = 384400.0;  // Distance to moon in km
```

**Pro tip:** Use `double` for most decimal calculations unless you specifically need to save memory.

### Character Type

#### `char`
- Size: 16 bits (2 bytes)
- Stores: A single character
- Use when: You need to store one character

```java
char grade = 'A';
char initial = 'J';
char symbol = '@';
```

**Important:** Use single quotes `' '` for char, not double quotes!

You can also use Unicode values:
```java
char heart = '\u2764';  // ❤ symbol
char copyright = '\u00A9';  // © symbol
```

### Boolean Type

#### `boolean`
- Size: Not precisely defined (typically 1 byte)
- Values: Only `true` or `false`
- Use when: You need a yes/no, on/off, true/false value

```java
boolean isStudent = true;
boolean hasLicense = false;
boolean isRaining = true;
```

## Variable Declaration and Initialization

### Declaration
This creates the variable but doesn't give it a value yet:

```java
int count;
double price;
String name;
```

### Initialization
This gives the variable its first value:

```java
count = 0;
price = 9.99;
name = "Alex";
```

### Declaration + Initialization
You can do both at once (most common approach):

```java
int count = 0;
double price = 9.99;
String name = "Alex";
```

### Multiple Variables
You can declare multiple variables of the same type in one line:

```java
int x = 1, y = 2, z = 3;
double pi = 3.14, e = 2.71;
```

## Variable Naming Conventions

Java has some rules and conventions for naming variables:

### Rules (Must Follow)
- Must start with a letter, `$`, or `_`
- Can contain letters, numbers, `$`, or `_`
- Cannot use Java keywords (like `int`, `class`, `public`)
- Case-sensitive (`age` and `Age` are different)

### Conventions (Should Follow)
- Use **camelCase** for variable names: start with lowercase, capitalize each new word

```java
// Good examples
int studentAge;
double accountBalance;
String firstName;
boolean isComplete;

// Bad examples (but technically valid)
int StudentAge;  // Should start with lowercase
double account_balance;  // Use camelCase, not snake_case
String FirstName;  // Should start with lowercase
```

**Pro tip:** Make names descriptive! `customerAge` is way better than `ca` or `x`.

## Type Casting and Conversion

Sometimes you need to convert one type to another. There are two ways to do this:

### Implicit Casting (Widening)
Java automatically converts smaller types to larger types:

```java
int myInt = 100;
long myLong = myInt;  // int -> long (automatic)

float myFloat = myLong;  // long -> float (automatic)
double myDouble = myFloat;  // float -> double (automatic)
```

This works because you're not losing information—a larger container can hold everything from a smaller one.

### Explicit Casting (Narrowing)
You must manually convert larger types to smaller types:

```java
double myDouble = 9.78;
int myInt = (int) myDouble;  // myInt becomes 9 (decimal part lost!)

long myLong = 100000;
int anotherInt = (int) myLong;  // Works, but be careful with large values
```

**Warning:** You can lose data when narrowing! The decimal part disappears, and large values might overflow.

### Practical Examples

```java
// Be careful with division!
int a = 5;
int b = 2;
int result = a / b;  // result is 2, not 2.5!

// Fix it by casting
double correctResult = (double) a / b;  // Now it's 2.5

// Another way
double alsoCorrect = a / (double) b;  // Also 2.5
```

## Constants Using `final`

Want a value that never changes? Use the `final` keyword:

```java
final double PI = 3.14159;
final int MAX_STUDENTS = 30;
final String SCHOOL_NAME = "ZipCode";
```

**Conventions for constants:**
- Use `final` keyword
- Use ALL_CAPS with underscores
- Initialize when declaring (you can't change it later anyway!)

```java
final int MAX_SPEED = 100;
MAX_SPEED = 120;  // ERROR! Can't change a final variable
```

## Literals

Literals are the actual values you write in your code.

### Integer Literals

```java
int decimal = 100;      // Regular decimal
int hex = 0x64;         // Hexadecimal (same as 100)
int binary = 0b1100100; // Binary (same as 100)
int withUnderscores = 1_000_000;  // Use underscores for readability!
```

### Floating-Point Literals

```java
double d1 = 123.4;
double d2 = 1.234e2;  // Scientific notation (123.4)
float f = 123.4f;     // Don't forget the 'f'!
```

### Character Literals

```java
char letter = 'A';
char newline = '\n';  // Escape sequences
char tab = '\t';
char backslash = '\\';
char singleQuote = '\'';
```

### String Literals

```java
String message = "Hello, World!";
String quote = "He said, \"Hi!\"";  // Escape the quotes
String multiLine = "Line 1\nLine 2\nLine 3";
```

### Boolean Literals

```java
boolean yes = true;
boolean no = false;
```

## Practical Examples

Let's put it all together with some real-world examples:

### Example 1: Calculate Area

```java
final double PI = 3.14159;
double radius = 5.0;
double area = PI * radius * radius;
System.out.println("Area: " + area);  // Area: 78.53975
```

### Example 2: Temperature Conversion

```java
double fahrenheit = 98.6;
double celsius = (fahrenheit - 32) * 5.0 / 9.0;
System.out.println(fahrenheit + "°F = " + celsius + "°C");
```

### Example 3: Student Grade

```java
String studentName = "Jordan";
int score = 87;
char grade;

if (score >= 90) {
    grade = 'A';
} else if (score >= 80) {
    grade = 'B';
} else if (score >= 70) {
    grade = 'C';
} else {
    grade = 'F';
}

System.out.println(studentName + " earned a " + grade);
```

### Example 4: Using Different Integer Types

```java
byte age = 25;                    // Small number
short yearOfBirth = 1998;         // Medium number
int population = 328_000_000;     // Large number
long distanceToSun = 93_000_000L; // Really large number

System.out.println("Age: " + age);
System.out.println("Born: " + yearOfBirth);
System.out.println("Population: " + population);
System.out.println("Distance to sun (miles): " + distanceToSun);
```

## Common Pitfalls to Avoid

### 1. Uninitialized Variables

```java
int x;
System.out.println(x);  // ERROR! Variable might not have been initialized
```

Always initialize your variables before using them!

### 2. Integer Division

```java
int result = 5 / 2;  // result is 2, not 2.5!
```

Remember: when both operands are integers, you get integer division.

### 3. Overflow

```java
int big = 2_000_000_000;
int tooBig = big + big;  // Overflow! Wrong result
```

Use `long` for numbers that might exceed int's range.

### 4. Float Precision

```java
float f1 = 0.1f;
float f2 = 0.2f;
float sum = f1 + f2;  // Might not be exactly 0.3!
```

Floating-point math isn't always perfectly precise. Use `double` for better precision, or `BigDecimal` for money calculations.

## Best Practices

1. **Choose the right type** - Use `int` for integers, `double` for decimals as defaults
2. **Name variables clearly** - `studentAge` beats `sa` every time
3. **Use camelCase** - Follow Java naming conventions
4. **Initialize variables** - Don't leave them undefined
5. **Use `final` for constants** - Make your intent clear
6. **Add underscores to large numbers** - `1_000_000` is more readable than `1000000`

## Practice Exercises

1. Declare variables to store a person's name, age, height (in meters), and whether they're a student. Choose appropriate types!

2. Write a program that converts inches to centimeters (1 inch = 2.54 cm).

3. Create constants for the number of days in a week, hours in a day, and minutes in an hour. Calculate the total minutes in a week.

4. Experiment with type casting: create a double variable with value 15.8, cast it to int, and observe what happens.

## Wrapping Up

Variables and data types are the foundation of Java programming. Every program you write will use them extensively. Understanding when to use each type, how to name variables properly, and how type conversion works will serve you well throughout your Java journey.

Take your time to practice these concepts. Try the exercises, experiment with the code examples, and when you feel ready, head over to Athena for some quiz questions to test your knowledge!
