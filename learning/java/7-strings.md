# Java Strings

Welcome to the world of Strings! In Java, strings are everywhere—usernames, messages, file names, web addresses, you name it. A String is a sequence of characters, and Java provides powerful tools for working with text.

Unlike primitive types, String is actually a class in Java, which means it comes with a rich set of methods to manipulate text. Let's explore how to work with strings effectively!

## String Class and String Literals

A String represents text—a sequence of characters like "Hello", "Java", or "abc123". In Java, strings are objects, but they're so commonly used that Java gives them special treatment.

### Creating Strings

There are two main ways to create strings:

```java
// Method 1: String literal (most common)
String greeting = "Hello, World!";
String name = "Alex";

// Method 2: Using the new keyword (less common)
String message = new String("Hello");
```

**Pro tip**: Use string literals (the first method) in most cases. They're simpler and more efficient.

### String Literals and the String Pool

When you create a string literal, Java stores it in a special memory area called the "string pool." If you create another string with the same value, Java reuses the existing one:

```java
String s1 = "Hello";
String s2 = "Hello";
String s3 = new String("Hello");

System.out.println(s1 == s2);  // true (same object from pool)
System.out.println(s1 == s3);  // false (s3 is a new object)
System.out.println(s1.equals(s3));  // true (same content)
```

This is why you should use `.equals()` to compare string content, not `==`.

### Escape Sequences

Special characters need escape sequences (backslash followed by a character):

```java
String quote = "He said, \"Hello!\"";  // Quotes inside string
String path = "C:\\Users\\Alex\\file.txt";  // Backslashes
String multiLine = "Line 1\nLine 2\nLine 3";  // Newlines
String tabbed = "Name:\tAlex";  // Tab character

System.out.println(quote);     // He said, "Hello!"
System.out.println(multiLine); // Three separate lines
```

Common escape sequences:
- `\"` - Double quote
- `\'` - Single quote
- `\\` - Backslash
- `\n` - Newline
- `\t` - Tab
- `\r` - Carriage return

### Empty Strings

An empty string is a string with no characters:

```java
String empty1 = "";
String empty2 = new String();

System.out.println(empty1.length());  // 0
System.out.println(empty1.isEmpty()); // true
```

**Important**: Empty string (`""`) is different from `null`:

```java
String empty = "";
String nothing = null;

System.out.println(empty.length());   // 0
System.out.println(nothing.length()); // NullPointerException!
```

## String Immutability

Here's something crucial: **Strings in Java are immutable**—once created, they cannot be changed. Any operation that appears to modify a string actually creates a new string.

### What Immutability Means

```java
String s = "Hello";
s.toUpperCase();
System.out.println(s);  // Still "Hello"! String unchanged

// To actually use the uppercase version, assign it
String upper = s.toUpperCase();
System.out.println(upper);  // "HELLO"

// Or reassign to the same variable
s = s.toUpperCase();
System.out.println(s);  // "HELLO"
```

When you call methods like `toUpperCase()`, they return a *new* string. The original string remains unchanged.

### Why Immutability?

Immutability brings several benefits:

1. **Thread safety**: Multiple threads can safely share strings
2. **Security**: String values can't be changed unexpectedly
3. **String pooling**: Identical strings can share memory
4. **Caching**: Hash codes can be cached safely

### Implications of Immutability

```java
String result = "";
for (int i = 0; i < 1000; i++) {
    result += i;  // Creates 1000 new string objects! Inefficient!
}
```

This is very inefficient because each `+=` creates a new string. For building strings in loops, use `StringBuilder` instead (we'll cover that later).

### String Concatenation

Even concatenation creates new strings:

```java
String first = "Hello";
String second = "World";
String combined = first + " " + second;  // Creates new string "Hello World"

System.out.println(first);    // Still "Hello"
System.out.println(combined); // "Hello World"
```

## String Operations

Java's String class provides dozens of useful methods. Let's explore the most important ones.

### Concatenation

Combine strings using `+` or the `concat()` method:

```java
String firstName = "Alex";
String lastName = "Smith";

// Method 1: Using + operator
String fullName = firstName + " " + lastName;  // "Alex Smith"

// Method 2: Using concat() method
String fullName2 = firstName.concat(" ").concat(lastName);  // "Alex Smith"

// Concatenating with numbers
String message = "Age: " + 25;  // "Age: 25" (number converted to string)
```

The `+` operator is simpler and more commonly used.

### Comparison

Compare strings using `equals()`, `equalsIgnoreCase()`, or `compareTo()`:

```java
String s1 = "Hello";
String s2 = "Hello";
String s3 = "hello";
String s4 = "World";

// Check equality (case-sensitive)
System.out.println(s1.equals(s2));      // true
System.out.println(s1.equals(s3));      // false (different case)

// Check equality (case-insensitive)
System.out.println(s1.equalsIgnoreCase(s3));  // true

// Lexicographic comparison (dictionary order)
System.out.println(s1.compareTo(s2));   // 0 (equal)
System.out.println(s1.compareTo(s4));   // negative (Hello < World)
System.out.println(s4.compareTo(s1));   // positive (World > Hello)
```

**Remember**: Use `.equals()` for equality, never `==`!

### Searching

Find substrings or characters within a string:

```java
String text = "Java programming is fun!";

// Check if string contains a substring
System.out.println(text.contains("programming"));  // true
System.out.println(text.contains("Python"));       // false

// Check if string starts/ends with something
System.out.println(text.startsWith("Java"));   // true
System.out.println(text.endsWith("fun!"));     // true

// Find index of substring or character
System.out.println(text.indexOf("programming"));     // 5
System.out.println(text.indexOf("xyz"));             // -1 (not found)
System.out.println(text.indexOf('a'));               // 1 (first 'a')
System.out.println(text.lastIndexOf('a'));           // 8 (last 'a')

// Find from a specific starting position
System.out.println(text.indexOf('a', 2));  // 3 (first 'a' after index 2)
```

### Extraction

Extract parts of a string:

```java
String text = "Java Programming";

// Get a single character
char ch = text.charAt(0);      // 'J'
char ch2 = text.charAt(5);     // 'P'

// Get a substring
String sub1 = text.substring(5);      // "Programming" (from index 5 to end)
String sub2 = text.substring(0, 4);   // "Java" (from index 0 to 3)
String sub3 = text.substring(5, 9);   // "Prog" (from index 5 to 8)

// Note: substring(start, end) includes start but excludes end!
```

**Important**: `substring(start, end)` includes the character at `start` but excludes the character at `end`.

### Modification (Creates New Strings)

Transform strings in various ways:

```java
String text = "  Hello World  ";

// Change case
System.out.println(text.toUpperCase());  // "  HELLO WORLD  "
System.out.println(text.toLowerCase());  // "  hello world  "

// Remove whitespace from beginning and end
System.out.println(text.trim());         // "Hello World"

// Replace characters or substrings
System.out.println(text.replace('l', 'L'));          // "  HeLLo WorLd  "
System.out.println(text.replace("World", "Java"));   // "  Hello Java  "
System.out.println(text.replaceAll("\\s+", "-"));    // "--Hello-World--"
System.out.println(text.replaceFirst("l", "L"));     // "  HeLlo World  "

// Note: Original string 'text' remains unchanged!
```

### Length and Empty Check

```java
String s1 = "Hello";
String s2 = "";

System.out.println(s1.length());   // 5
System.out.println(s2.length());   // 0
System.out.println(s2.isEmpty());  // true
System.out.println(s1.isEmpty());  // false

// Java 11+ also has isBlank() for whitespace-only strings
String s3 = "   ";
System.out.println(s3.isEmpty());  // false (has spaces)
System.out.println(s3.isBlank());  // true (only whitespace)
```

### Splitting and Joining

Split strings into arrays or join arrays into strings:

```java
// Splitting
String csv = "apple,banana,cherry,date";
String[] fruits = csv.split(",");
for (String fruit : fruits) {
    System.out.println(fruit);
}
// Output: apple, banana, cherry, date (each on separate line)

String sentence = "The quick brown fox";
String[] words = sentence.split(" ");
System.out.println(words.length);  // 4

// Split with limit
String data = "a:b:c:d:e";
String[] parts = data.split(":", 3);
// parts = ["a", "b", "c:d:e"] (splits into max 3 parts)

// Joining (Java 8+)
String[] items = {"apple", "banana", "cherry"};
String joined = String.join(", ", items);
System.out.println(joined);  // "apple, banana, cherry"

String path = String.join("/", "home", "user", "documents");
System.out.println(path);  // "home/user/documents"
```

## StringBuilder and StringBuffer Classes

When you need to build strings efficiently (especially in loops), use `StringBuilder` or `StringBuffer`. These classes create mutable strings.

### Why Use StringBuilder?

Remember, strings are immutable. Concatenating in a loop is inefficient:

```java
// INEFFICIENT - creates many temporary string objects
String result = "";
for (int i = 0; i < 1000; i++) {
    result += i + ",";
}

// EFFICIENT - modifies one StringBuilder object
StringBuilder sb = new StringBuilder();
for (int i = 0; i < 1000; i++) {
    sb.append(i).append(",");
}
String result2 = sb.toString();
```

The `StringBuilder` version is much faster and uses less memory.

### StringBuilder Methods

```java
StringBuilder sb = new StringBuilder();

// Append (add to end)
sb.append("Hello");
sb.append(" ");
sb.append("World");
sb.append(123);
System.out.println(sb.toString());  // "Hello World123"

// Insert at specific position
sb.insert(5, " Beautiful");
System.out.println(sb.toString());  // "Hello Beautiful World123"

// Delete characters
sb.delete(5, 15);  // Delete from index 5 to 14
System.out.println(sb.toString());  // "Hello World123"

// Replace a section
sb.replace(6, 11, "Java");
System.out.println(sb.toString());  // "Hello Java123"

// Reverse
sb.reverse();
System.out.println(sb.toString());  // "321avaJ olleH"

// Get length and capacity
System.out.println(sb.length());    // Current length
System.out.println(sb.capacity());  // Current capacity
```

### StringBuilder vs StringBuffer

- **StringBuilder**: Faster, not thread-safe (use in single-threaded programs)
- **StringBuffer**: Slower, thread-safe (use when multiple threads access it)

In most cases, use `StringBuilder`:

```java
// StringBuilder - preferred for most cases
StringBuilder sb = new StringBuilder("Hello");
sb.append(" World");

// StringBuffer - use when thread safety is needed
StringBuffer buffer = new StringBuffer("Hello");
buffer.append(" World");
```

### Practical Example: Building HTML

```java
StringBuilder html = new StringBuilder();
html.append("<html>\n");
html.append("  <head>\n");
html.append("    <title>").append("My Page").append("</title>\n");
html.append("  </head>\n");
html.append("  <body>\n");
html.append("    <h1>").append("Welcome!").append("</h1>\n");
html.append("  </body>\n");
html.append("</html>");

System.out.println(html.toString());
```

## String Formatting

String formatting creates formatted strings from templates and values.

### Using String.format()

```java
// Basic formatting
String name = "Alex";
int age = 25;
String formatted = String.format("Name: %s, Age: %d", name, age);
System.out.println(formatted);  // "Name: Alex, Age: 25"

// Formatting numbers
double price = 19.99;
String priceStr = String.format("Price: $%.2f", price);
System.out.println(priceStr);  // "Price: $19.99"

// Multiple values
String message = String.format("%s scored %d out of %d (%.1f%%)",
    "Alice", 95, 100, 95.0);
System.out.println(message);  // "Alice scored 95 out of 100 (95.0%)"
```

### Format Specifiers

Common format specifiers:
- `%s` - String
- `%d` - Integer (decimal)
- `%f` - Floating-point
- `%b` - Boolean
- `%c` - Character
- `%n` - Newline (platform-independent)

### Formatting Options

```java
// Width and alignment
System.out.printf("|%10s|%n", "Hi");     // "|        Hi|" (right-aligned)
System.out.printf("|%-10s|%n", "Hi");    // "|Hi        |" (left-aligned)

// Padding numbers
System.out.printf("%05d%n", 42);         // "00042" (pad with zeros)
System.out.printf("%,d%n", 1000000);     // "1,000,000" (thousands separator)

// Decimal places
System.out.printf("%.2f%n", 3.14159);    // "3.14" (2 decimal places)
System.out.printf("%.0f%n", 3.14159);    // "3" (0 decimal places)

// Scientific notation
System.out.printf("%e%n", 1000000.0);    // "1.000000e+06"
```

### printf() vs format()

`printf()` prints directly, while `format()` returns the formatted string:

```java
// printf - prints to console
System.out.printf("Hello %s%n", "World");  // Prints: Hello World

// format - returns string
String s = String.format("Hello %s", "World");
System.out.println(s);  // Prints: Hello World
```

### Practical Example: Receipt

```java
String item1 = "Apple";
String item2 = "Banana";
String item3 = "Orange";
double price1 = 1.50;
double price2 = 0.75;
double price3 = 2.00;

System.out.println("===== RECEIPT =====");
System.out.printf("%-10s $%5.2f%n", item1, price1);
System.out.printf("%-10s $%5.2f%n", item2, price2);
System.out.printf("%-10s $%5.2f%n", item3, price3);
System.out.println("-------------------");
double total = price1 + price2 + price3;
System.out.printf("%-10s $%5.2f%n", "TOTAL", total);
```

Output:
```
===== RECEIPT =====
Apple      $ 1.50
Banana     $ 0.75
Orange     $ 2.00
-------------------
TOTAL      $ 4.25
```

## Regular Expressions Basics

Regular expressions (regex) are patterns for matching text. Java's String class has built-in regex support.

### Basic Pattern Matching

```java
String text = "Hello123";

// Check if string matches a pattern
System.out.println(text.matches("Hello\\d+"));     // true (Hello followed by digits)
System.out.println(text.matches("Hello"));         // false (digits don't match)

String email = "user@example.com";
System.out.println(email.matches(".*@.*\\..*"));   // true (basic email pattern)
```

### Common Regex Patterns

```java
// Digits
"\\d"       // Any single digit (0-9)
"\\d+"      // One or more digits
"\\d{3}"    // Exactly 3 digits

// Word characters (letters, digits, underscore)
"\\w"       // Any word character
"\\w+"      // One or more word characters

// Whitespace
"\\s"       // Any whitespace character
"\\s+"      // One or more whitespace characters

// Any character
"."         // Any character (except newline)
".*"        // Zero or more of any character
".+"        // One or more of any character
```

### Replacing with Regex

```java
String text = "Hello   World  !";

// Replace one or more spaces with single space
String cleaned = text.replaceAll("\\s+", " ");
System.out.println(cleaned);  // "Hello World !"

// Remove all non-digits
String phone = "(555) 123-4567";
String digits = phone.replaceAll("\\D", "");
System.out.println(digits);  // "5551234567"

// Validate and extract
String date = "2024-03-15";
if (date.matches("\\d{4}-\\d{2}-\\d{2}")) {
    System.out.println("Valid date format");
}
```

### Practical Examples

```java
// Validate email (simplified)
String email = "user@example.com";
boolean isValidEmail = email.matches("\\w+@\\w+\\.\\w+");

// Validate phone number
String phone = "555-1234";
boolean isValidPhone = phone.matches("\\d{3}-\\d{4}");

// Extract words
String sentence = "The quick brown fox";
String[] words = sentence.split("\\s+");

// Remove special characters
String input = "Hello, World!@#";
String clean = input.replaceAll("[^a-zA-Z\\s]", "");
System.out.println(clean);  // "Hello World"
```

### Character Classes

```java
"[abc]"      // Matches a, b, or c
"[^abc]"     // Matches any character EXCEPT a, b, c
"[a-z]"      // Matches any lowercase letter
"[A-Z]"      // Matches any uppercase letter
"[0-9]"      // Matches any digit (same as \\d)
"[a-zA-Z]"   // Matches any letter

// Example: Validate username (letters and numbers only, 3-15 chars)
String username = "User123";
boolean valid = username.matches("[a-zA-Z0-9]{3,15}");
```

## Common Pitfalls and Best Practices

### Pitfall 1: Using == Instead of equals()

```java
String s1 = "Hello";
String s2 = new String("Hello");
System.out.println(s1 == s2);        // false! Compares references
System.out.println(s1.equals(s2));   // true! Compares content
```

**Fix**: Always use `.equals()` for string comparison.

### Pitfall 2: Null Pointer Exceptions

```java
String s = null;
System.out.println(s.length());  // NullPointerException!
```

**Fix**: Check for null before calling methods:
```java
if (s != null && s.length() > 0) {
    // Safe to use s
}
```

### Pitfall 3: Inefficient Concatenation in Loops

```java
String result = "";
for (int i = 0; i < 1000; i++) {
    result += i;  // Inefficient!
}
```

**Fix**: Use StringBuilder:
```java
StringBuilder sb = new StringBuilder();
for (int i = 0; i < 1000; i++) {
    sb.append(i);
}
String result = sb.toString();
```

### Pitfall 4: Forgetting Immutability

```java
String s = "hello";
s.toUpperCase();
System.out.println(s);  // Still "hello"!
```

**Fix**: Assign the result:
```java
s = s.toUpperCase();
System.out.println(s);  // "HELLO"
```

### Best Practices

1. **Use `.equals()` for comparison**, never `==`
2. **Check for null** before calling string methods
3. **Use StringBuilder** for building strings in loops
4. **Use string literals** instead of `new String()`
5. **Remember immutability** - methods return new strings
6. **Use appropriate methods** - `isEmpty()`, `isBlank()`, `trim()`, etc.
7. **Validate input** with regex when appropriate

## Real-World Example: Text Processor

```java
public class TextProcessor {
    public static void main(String[] args) {
        String text = "  The Quick BROWN fox jumps over the Lazy dog!  ";
        
        // Clean and normalize
        String cleaned = text.trim().toLowerCase();
        System.out.println("Cleaned: " + cleaned);
        
        // Count words
        String[] words = cleaned.split("\\s+");
        System.out.println("Word count: " + words.length);
        
        // Remove punctuation
        String noPunct = cleaned.replaceAll("[^a-z\\s]", "");
        System.out.println("No punctuation: " + noPunct);
        
        // Count specific word
        int count = 0;
        for (String word : words) {
            if (word.replaceAll("[^a-z]", "").equals("the")) {
                count++;
            }
        }
        System.out.println("'the' appears: " + count + " times");
        
        // Build formatted output
        StringBuilder result = new StringBuilder();
        result.append("Original length: ").append(text.length()).append("\n");
        result.append("Cleaned length: ").append(cleaned.length()).append("\n");
        result.append("Word count: ").append(words.length);
        
        System.out.println("\n" + result.toString());
    }
}
```

---

## Quiz Questions

### Question 1: String Comparison
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: java, strings, comparison, equals

**Question Text**:
What is the correct way to compare two String variables for equality in Java?

**Answer Choices**:
- A) `str1 == str2`
- B) `str1.equals(str2)`
- C) `str1.compareTo(str2) == 0`
- D) Both B and C are correct

**Correct Answer**: D

**Explanation**:
Both `str1.equals(str2)` and `str1.compareTo(str2) == 0` correctly compare string content. The `equals()` method is the standard way, while `compareTo()` returns 0 for equal strings. However, `==` compares references (memory addresses), not content, and should not be used. Using `.equals()` is the preferred and most readable approach.

---

### Question 2: String Immutability
**Type**: code-output
**Difficulty**: 3
**Tags**: java, strings, immutability

**Question Text**:
What does this code print?

```java
String s = "hello";
s.toUpperCase();
System.out.println(s);
```

**Answer Choices**:
- A) HELLO
- B) hello
- C) Hello
- D) Compilation error

**Correct Answer**: B

**Explanation**:
Strings in Java are immutable, meaning they cannot be changed after creation. The `toUpperCase()` method returns a new string with uppercase letters but doesn't modify the original string. Since we don't assign the result to anything, it's discarded, and `s` remains "hello". To actually change s, you would need: `s = s.toUpperCase();`.

---

### Question 3: String Methods
**Type**: code-output
**Difficulty**: 3
**Tags**: java, strings, methods, substring

**Question Text**:
What is the output of this code?

```java
String text = "Java Programming";
System.out.println(text.substring(5, 9));
```

**Answer Choices**:
- A) Prog
- B) Progr
- C) Programming
- D) ava P

**Correct Answer**: A

**Explanation**:
The `substring(start, end)` method returns characters from index `start` up to but NOT including index `end`. So `substring(5, 9)` returns characters at indices 5, 6, 7, and 8, which are "Prog". The character at index 9 ('r') is not included. This is a common source of off-by-one errors.

---

### Question 4: StringBuilder vs String
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, strings, stringbuilder, performance

**Question Text**:
Why should you use StringBuilder instead of String concatenation in a loop?

**Answer Choices**:
- A) StringBuilder is easier to use
- B) String concatenation doesn't work in loops
- C) StringBuilder is mutable and more efficient for repeated modifications
- D) StringBuilder is thread-safe

**Correct Answer**: C

**Explanation**:
StringBuilder is mutable, meaning it can be modified without creating new objects. String concatenation in a loop creates a new String object for each iteration, which is very inefficient. StringBuilder modifies one object, making it much faster and more memory-efficient. Note that StringBuilder is NOT thread-safe (StringBuffer is), and string concatenation does work in loops, just inefficiently.

---

### Question 5: String Pool
**Type**: code-output
**Difficulty**: 5
**Tags**: java, strings, string-pool, references

**Question Text**:
What does this code print?

```java
String s1 = "Hello";
String s2 = "Hello";
String s3 = new String("Hello");
System.out.println(s1 == s2);
System.out.println(s1 == s3);
```

**Answer Choices**:
- A) true, true
- B) false, false
- C) true, false
- D) false, true

**Correct Answer**: C

**Explanation**:
String literals are stored in a string pool. Both s1 and s2 refer to the same "Hello" in the pool, so `s1 == s2` is true. However, `new String("Hello")` creates a new object outside the pool, so s3 is a different object and `s1 == s3` is false. This demonstrates why you should use `.equals()` instead of `==` for string comparison—`.equals()` would return true for all comparisons here.

---

### Question 6: String Formatting
**Type**: code-output
**Difficulty**: 4
**Tags**: java, strings, formatting, printf

**Question Text**:
What does this code print?

```java
double price = 12.5;
System.out.printf("Price: $%.2f", price);
```

**Answer Choices**:
- A) Price: $12.5
- B) Price: $12.50
- C) Price: 12.50
- D) $12.50

**Correct Answer**: B

**Explanation**:
The format specifier `%.2f` formats a floating-point number with exactly 2 decimal places. So 12.5 becomes 12.50. The entire format string includes "Price: $", so the complete output is "Price: $12.50". The `printf` method formats and prints the string directly to the console.

---

### Question 7: String Splitting
**Type**: code-output
**Difficulty**: 3
**Tags**: java, strings, split, arrays

**Question Text**:
What is the length of the array?

```java
String text = "apple,banana,cherry";
String[] parts = text.split(",");
System.out.println(parts.length);
```

**Answer Choices**:
- A) 1
- B) 2
- C) 3
- D) 19

**Correct Answer**: C

**Explanation**:
The `split(",")` method splits the string at each comma, creating an array of the parts. "apple,banana,cherry" splits into three parts: ["apple", "banana", "cherry"], so the array length is 3. The split method returns an array containing all the substrings between the delimiter.

---

### Question 8: Regular Expressions
**Type**: code-output
**Difficulty**: 6
**Tags**: java, strings, regex, validation

**Question Text**:
What does this code print?

```java
String phone = "123-4567";
System.out.println(phone.matches("\\d{3}-\\d{4}"));
```

**Answer Choices**:
- A) true
- B) false
- C) 123-4567
- D) Compilation error

**Correct Answer**: A

**Explanation**:
The regex pattern `\\d{3}-\\d{4}` matches exactly 3 digits, followed by a hyphen, followed by exactly 4 digits. The string "123-4567" matches this pattern perfectly, so `matches()` returns true. `\\d` represents any digit, and `{n}` means exactly n occurrences. This is a common pattern for validating simple phone number formats.
