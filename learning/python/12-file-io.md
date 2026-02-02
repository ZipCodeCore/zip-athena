# Python File I/O

Hey there! Let's dive into one of the most practical topics in programming: working with files. Whether you're reading configuration files, saving user data, processing logs, or working with CSV spreadsheets, file I/O (Input/Output) is a skill you'll use constantly. Think of it as teaching your program how to read and write—just like you learned in school, but for computers!

## 12.1 Opening Files with the open() Function

Before you can do anything with a file, you need to open it. Python makes this super simple with the built-in `open()` function. Think of it like opening a book—you need to open it before you can read or write in it.

The basic syntax is straightforward:

```python
file = open('example.txt', 'r')
# Do something with the file
file.close()
```

The `open()` function takes two main arguments: the filename (including the path if needed) and the mode. The mode tells Python what you want to do with the file—read it, write to it, or append to it. We'll explore all the different modes in the next section.

Here's a simple example that opens a file, reads its contents, and properly closes it:

```python
file = open('greeting.txt', 'r')
content = file.read()
print(content)
file.close()
```

**Important:** Always close your files when you're done! Leaving files open can cause problems like data loss or preventing other programs from accessing the file. It's like leaving books scattered everywhere instead of putting them back on the shelf.

## 12.2 File Modes (r, w, a, r+, w+, a+, rb, wb)

File modes are like different keys that unlock different capabilities. Let's break them down:

**Text Modes:**
- `'r'` - Read mode (default). Opens the file for reading. The file must exist, or you'll get an error.
- `'w'` - Write mode. Creates a new file or **overwrites** an existing file. Be careful—this erases the old content!
- `'a'` - Append mode. Adds new content to the end of the file without erasing what's already there.
- `'r+'` - Read and write mode. File must exist. You can both read and write.
- `'w+'` - Write and read mode. Creates a new file or overwrites existing. You can read and write.
- `'a+'` - Append and read mode. Creates the file if it doesn't exist. You can read and append.

**Binary Modes:**
- `'rb'` - Read binary. For reading non-text files like images or videos.
- `'wb'` - Write binary. For writing binary data.
- `'ab'` - Append binary. For adding binary data to the end of a file.

```python
# Reading a file (must exist)
file = open('data.txt', 'r')
print(file.read())
file.close()

# Writing to a file (creates new or overwrites)
file = open('output.txt', 'w')
file.write('Hello, World!')
file.close()

# Appending to a file (adds to the end)
file = open('log.txt', 'a')
file.write('New log entry\n')
file.close()

# Working with binary files (like images)
file = open('image.jpg', 'rb')
data = file.read()
file.close()
```

A common mistake is using `'w'` mode when you mean to use `'a'` mode—you'll accidentally erase your entire file! Double-check your mode before opening important files.

## 12.3 Reading Files (read(), readline(), readlines(), iterating over file object)

Python gives you several ways to read files, depending on what you need. It's like having different tools for different jobs!

**The `read()` method** reads the entire file as a single string:

```python
file = open('story.txt', 'r')
entire_content = file.read()
print(entire_content)
file.close()
```

You can also specify how many characters to read:

```python
file = open('story.txt', 'r')
first_100_chars = file.read(100)
print(first_100_chars)
file.close()
```

**The `readline()` method** reads one line at a time:

```python
file = open('poem.txt', 'r')
first_line = file.readline()
second_line = file.readline()
print(first_line)
print(second_line)
file.close()
```

**The `readlines()` method** reads all lines into a list:

```python
file = open('shopping_list.txt', 'r')
all_lines = file.readlines()
for line in all_lines:
    print(line.strip())  # strip() removes the newline character
file.close()
```

**Iterating over the file object** is the most Pythonic and memory-efficient way:

```python
file = open('large_file.txt', 'r')
for line in file:
    print(line.strip())
file.close()
```

This last method is great for huge files because it reads one line at a time instead of loading everything into memory. It's like reading a book page by page instead of trying to memorize the whole thing at once!

## 12.4 Writing Files (write(), writelines())

Writing to files is just as easy as reading from them. The main methods are `write()` and `writelines()`.

**The `write()` method** writes a string to the file:

```python
file = open('message.txt', 'w')
file.write('Hello, Python!')
file.close()
```

Note that `write()` doesn't automatically add newlines, so you need to add them yourself:

```python
file = open('notes.txt', 'w')
file.write('First line\n')
file.write('Second line\n')
file.write('Third line\n')
file.close()
```

**The `writelines()` method** writes a list of strings to the file:

```python
lines = ['Apple\n', 'Banana\n', 'Cherry\n']
file = open('fruits.txt', 'w')
file.writelines(lines)
file.close()
```

Here's a practical example that creates a simple report:

```python
file = open('report.txt', 'w')
file.write('Sales Report\n')
file.write('=' * 30 + '\n')
file.write('Total Sales: $1,500\n')
file.write('Total Customers: 25\n')
file.close()
```

Remember: using `'w'` mode will erase any existing content! If you want to add to a file instead, use `'a'` (append) mode.

## 12.5 Closing Files and the with Statement (Context Managers)

Here's a golden rule: **always close your files**. But here's an even better rule: **use the `with` statement** and Python will close them automatically!

The `with` statement is a context manager that ensures your file gets closed properly, even if an error occurs. It's like having a responsible assistant who always puts things away for you:

```python
# The old way (manual closing)
file = open('data.txt', 'r')
content = file.read()
file.close()

# The better way (automatic closing)
with open('data.txt', 'r') as file:
    content = file.read()
# File is automatically closed here!
```

The magic happens when you leave the `with` block—Python automatically closes the file, even if an exception occurs:

```python
try:
    with open('important.txt', 'r') as file:
        content = file.read()
        # Even if an error occurs here...
        result = 1 / 0  # This causes an error!
except ZeroDivisionError:
    print("Oops, division by zero!")
# ... the file still gets closed properly!
```

You can even work with multiple files at once:

```python
with open('source.txt', 'r') as source, open('destination.txt', 'w') as dest:
    content = source.read()
    dest.write(content.upper())
```

Using `with` is considered best practice in Python. It makes your code cleaner, safer, and less prone to resource leaks. Once you get used to it, you'll never want to go back to manual file closing!

## 12.6 File Positioning (seek(), tell())

Sometimes you need to jump around in a file, like using a bookmark in a book. Python provides `seek()` and `tell()` methods for this.

**The `tell()` method** returns your current position in the file (in bytes from the beginning):

```python
with open('sample.txt', 'r') as file:
    print(file.tell())  # 0 (at the beginning)
    file.read(10)
    print(file.tell())  # 10 (moved 10 bytes forward)
```

**The `seek()` method** lets you jump to a specific position:

```python
with open('sample.txt', 'r') as file:
    file.seek(0)  # Go to the beginning
    first_char = file.read(1)
    
    file.seek(10)  # Jump to position 10
    next_char = file.read(1)
    
    file.seek(0, 2)  # Jump to the end (0 bytes from end)
```

The `seek()` method takes two arguments: offset and whence. The whence parameter tells Python where to measure from:
- `0` - Beginning of file (default)
- `1` - Current position
- `2` - End of file

Here's a practical example that reads the last 100 bytes of a log file:

```python
with open('server.log', 'r') as file:
    file.seek(-100, 2)  # Go to 100 bytes before the end
    recent_logs = file.read()
    print(recent_logs)
```

Note: When working with text files (not binary), seeking can be tricky due to character encoding. For most text file operations, you won't need `seek()` and `tell()`, but they're super useful for binary files!

## 12.7 Working with File Paths (os.path module, pathlib module)

Hardcoding file paths like `'C:\\Users\\John\\file.txt'` is a recipe for trouble. Different operating systems use different path separators (Windows uses `\`, while Mac and Linux use `/`). Python has tools to handle this elegantly!

**Using the `os.path` module** (the traditional approach):

```python
import os

# Join path components (works on any OS)
file_path = os.path.join('data', 'files', 'input.txt')
print(file_path)  # data/files/input.txt (or data\files\input.txt on Windows)

# Get the current working directory
current_dir = os.getcwd()
print(current_dir)

# Build a full path
full_path = os.path.join(current_dir, 'output.txt')

# Check if a file exists
if os.path.exists('data.txt'):
    print('File exists!')

# Get file information
print(os.path.basename('/home/user/file.txt'))  # file.txt
print(os.path.dirname('/home/user/file.txt'))   # /home/user
print(os.path.splitext('file.txt'))             # ('file', '.txt')
```

**Using the `pathlib` module** (the modern, object-oriented approach):

```python
from pathlib import Path

# Create a Path object
file_path = Path('data') / 'files' / 'input.txt'
print(file_path)

# Get the current working directory
current_dir = Path.cwd()

# Check if a file exists
if file_path.exists():
    print('File exists!')

# Get file information
print(file_path.name)        # input.txt
print(file_path.stem)        # input
print(file_path.suffix)      # .txt
print(file_path.parent)      # data/files

# Read and write files directly
content = file_path.read_text()
file_path.write_text('New content')
```

`pathlib` is generally preferred in modern Python because it's more intuitive and object-oriented. Here's a practical example that finds all Python files in a directory:

```python
from pathlib import Path

# Find all .py files in the current directory and subdirectories
for py_file in Path('.').rglob('*.py'):
    print(f'Found: {py_file}')
    print(f'Size: {py_file.stat().st_size} bytes')
```

## 12.8 Reading and Writing Different Formats (Text, CSV, JSON)

Real-world programs often work with structured data formats. Let's explore the most common ones!

**Working with Plain Text:**

```python
# Writing
with open('notes.txt', 'w') as file:
    file.write('This is a simple text file.\n')
    file.write('It can contain multiple lines.\n')

# Reading
with open('notes.txt', 'r') as file:
    for line in file:
        print(line.strip())
```

**Working with CSV (Comma-Separated Values):**

CSV files are like spreadsheets in text format. Python's `csv` module makes them easy to work with:

```python
import csv

# Writing CSV
with open('students.csv', 'w', newline='') as file:
    writer = csv.writer(file)
    writer.writerow(['Name', 'Age', 'Grade'])
    writer.writerow(['Alice', 20, 'A'])
    writer.writerow(['Bob', 22, 'B'])
    writer.writerow(['Charlie', 21, 'A'])

# Reading CSV
with open('students.csv', 'r') as file:
    reader = csv.reader(file)
    for row in reader:
        print(row)

# Using DictReader for named columns
with open('students.csv', 'r') as file:
    reader = csv.DictReader(file)
    for row in reader:
        print(f"{row['Name']} is {row['Age']} years old")
```

**Working with JSON (JavaScript Object Notation):**

JSON is perfect for storing structured data like dictionaries and lists:

```python
import json

# Writing JSON
data = {
    'name': 'Alice',
    'age': 25,
    'hobbies': ['reading', 'coding', 'gaming'],
    'address': {
        'city': 'Boston',
        'zip': '02101'
    }
}

with open('person.json', 'w') as file:
    json.dump(data, file, indent=4)  # indent makes it readable

# Reading JSON
with open('person.json', 'r') as file:
    loaded_data = json.load(file)
    print(loaded_data['name'])
    print(loaded_data['hobbies'])

# Working with JSON strings (not files)
json_string = json.dumps(data)
parsed_data = json.loads(json_string)
```

Here's a practical example that reads a CSV file and exports it to JSON:

```python
import csv
import json

# Read CSV
students = []
with open('students.csv', 'r') as file:
    reader = csv.DictReader(file)
    for row in reader:
        students.append(row)

# Write to JSON
with open('students.json', 'w') as file:
    json.dump(students, file, indent=2)
```

## 12.9 Exception Handling with Files

File operations can fail in many ways: the file might not exist, you might not have permission, the disk might be full, etc. Always use exception handling to make your code robust!

The most common exceptions you'll encounter:
- `FileNotFoundError` - The file doesn't exist
- `PermissionError` - You don't have permission to access the file
- `IsADirectoryError` - You tried to open a directory as a file
- `IOError` - General I/O operation failure

Here's how to handle them gracefully:

```python
try:
    with open('data.txt', 'r') as file:
        content = file.read()
        print(content)
except FileNotFoundError:
    print("Sorry, the file doesn't exist!")
except PermissionError:
    print("You don't have permission to read this file.")
except Exception as e:
    print(f"An unexpected error occurred: {e}")
```

A more comprehensive example with multiple operations:

```python
def read_config(filename):
    """Safely read a configuration file."""
    try:
        with open(filename, 'r') as file:
            data = file.read()
            return data
    except FileNotFoundError:
        print(f"Config file '{filename}' not found. Using defaults.")
        return None
    except PermissionError:
        print(f"Permission denied to read '{filename}'.")
        return None
    except Exception as e:
        print(f"Error reading config: {e}")
        return None

def write_log(filename, message):
    """Safely write to a log file."""
    try:
        with open(filename, 'a') as file:
            from datetime import datetime
            timestamp = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
            file.write(f"[{timestamp}] {message}\n")
        return True
    except PermissionError:
        print(f"Cannot write to '{filename}' - permission denied.")
        return False
    except Exception as e:
        print(f"Error writing log: {e}")
        return False

# Using the functions
config = read_config('settings.txt')
write_log('app.log', 'Application started')
```

You can also use the `finally` clause to ensure cleanup happens:

```python
file = None
try:
    file = open('important.txt', 'r')
    # Do something with the file
    content = file.read()
except FileNotFoundError:
    print("File not found!")
except Exception as e:
    print(f"Error: {e}")
finally:
    # This always runs, even if there was an error
    if file:
        file.close()
        print("File closed safely")
```

However, when using the `with` statement, you don't need the `finally` clause because the file is automatically closed:

```python
try:
    with open('important.txt', 'r') as file:
        content = file.read()
except FileNotFoundError:
    print("File not found!")
```

## Common Pitfalls

Here are some mistakes that trip up beginners (and sometimes even experienced programmers!):

1. **Forgetting to close files** - Use `with` statements to avoid this!
2. **Using 'w' mode when you meant 'a' mode** - You'll overwrite your entire file by accident.
3. **Not handling exceptions** - Your program will crash when files don't exist or permissions are denied.
4. **Reading huge files with `read()`** - This loads everything into memory. Use iteration for large files.
5. **Hardcoding file paths** - Use `os.path` or `pathlib` for cross-platform compatibility.
6. **Forgetting newline characters** - `write()` doesn't add them automatically, but `print()` does!
7. **Opening files in binary mode for text data** - You'll get bytes instead of strings.

## Best Practices

Follow these tips for clean, professional file handling:

1. **Always use `with` statements** - They ensure files are closed properly.
2. **Handle exceptions appropriately** - Don't let missing files crash your program.
3. **Use `pathlib` for path operations** - It's cleaner and more modern than `os.path`.
4. **Be explicit about encoding** - Use `open('file.txt', 'r', encoding='utf-8')` to avoid encoding issues.
5. **Use the right tool for the job** - CSV module for CSV files, JSON module for JSON files, etc.
6. **Don't keep files open longer than necessary** - Open, do your work, close (or use `with`).
7. **Validate file paths before using them** - Check if files exist before trying to read them.
8. **Use meaningful mode strings** - Comment your code if the mode isn't obvious.

## Practice Exercises

Try these exercises to solidify your file I/O skills:

1. Write a program that counts the number of lines, words, and characters in a text file.

2. Create a simple contact manager that stores names and phone numbers in a CSV file, with functions to add, search, and display contacts.

3. Build a program that reads a JSON configuration file with settings like theme colors and window size, then applies those settings.

4. Write a log analyzer that reads a log file and counts how many ERROR, WARNING, and INFO messages it contains.

5. Create a file backup utility that copies the contents of one file to another with a timestamp in the filename.

## Wrapping Up

File I/O is one of those fundamental skills that you'll use in almost every real-world Python project. Whether you're building a web scraper that saves data to CSV, a configuration system that reads JSON files, or a simple script that processes text files, you now have all the tools you need!

Remember: use `with` statements, handle exceptions, and choose the right format for your data. With these principles in mind, you'll be reading and writing files like a pro in no time!

---

## Quiz Questions

This section contains quiz questions that test understanding of Python File I/O concepts. Each question is formatted for easy extraction and database import.

### Question 1: File Opening Modes - Write vs Append
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: python, file-io, file-modes

**Question Text**:
What is the main difference between opening a file with 'w' mode versus 'a' mode?

**Answer Choices**:
- A) 'w' mode is faster than 'a' mode
- B) 'w' mode overwrites the file, while 'a' mode appends to the end
- C) 'w' mode is for reading, 'a' mode is for writing
- D) There is no difference between them

**Correct Answer**: B

**Explanation**:
The 'w' (write) mode creates a new file or completely overwrites an existing file, erasing all previous content. The 'a' (append) mode adds new content to the end of the file without erasing existing content. This is a critical difference—using 'w' when you meant 'a' can result in data loss! Option A is incorrect as performance is essentially the same. Option C is wrong because both modes are for writing. Option D is obviously false.

---

### Question 2: Context Managers - with Statement
**Type**: true-false
**Difficulty**: 3
**Tags**: python, file-io, context-managers, with-statement

**Question Text**:
True or False: When using the `with` statement to open a file, you must explicitly call `file.close()` to close the file.

**Answer Choices**:
- A) True
- B) False

**Correct Answer**: B

**Explanation**:
False! The beauty of the `with` statement is that it automatically closes the file when you exit the `with` block, even if an exception occurs. This is one of the main reasons to use `with` statements—they handle cleanup for you. Calling `file.close()` explicitly when using `with` is redundant and goes against the purpose of using a context manager.

---

### Question 3: Reading Files - Output Prediction
**Type**: code-output
**Difficulty**: 3
**Tags**: python, file-io, reading-files, readline

**Question Text**:
Given a file 'numbers.txt' containing:
```
1
2
3
4
5
```

What will this code print?

```python
with open('numbers.txt', 'r') as file:
    print(file.readline())
    print(file.readline())
```

**Answer Choices**:
- A) 1 2
- B) 1\n2\n
- C) 1 (followed by newline) 2 (followed by newline)
- D) [1, 2]

**Correct Answer**: C

**Explanation**:
The `readline()` method reads one line at a time, including the newline character at the end. The first call reads "1\n" and the second reads "2\n". When printed, each print() statement outputs the line (which already has a newline) and adds its own newline, resulting in extra blank lines. The output would be "1" on one line, a blank line, then "2" on another line, and another blank line. Option A would require stripping newlines or using read() with a limit. Option B shows the representation but not what's actually printed. Option D suggests a list, which isn't what readline() returns.

---

### Question 4: File Modes - Binary Mode
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: python, file-io, binary-mode, file-modes

**Question Text**:
When should you use binary mode ('rb' or 'wb') instead of text mode?

**Answer Choices**:
- A) When working with large files to improve performance
- B) When working with non-text files like images, videos, or audio
- C) When you want to read the file faster
- D) When the file contains numbers

**Correct Answer**: B

**Explanation**:
Binary mode ('rb', 'wb', 'ab') should be used for non-text files such as images, videos, audio files, executables, and other binary data. These files don't contain human-readable text and shouldn't be decoded as text. Text mode is the default and handles encoding/decoding automatically. Option A is incorrect—binary mode doesn't inherently improve performance for large files. Option C is wrong as speed isn't the determining factor. Option D is misleading—numbers can be stored as text and read in text mode just fine.

---

### Question 5: Exception Handling - File Operations
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: python, file-io, exception-handling, errors

**Question Text**:
Which exception is raised when you try to open a file that doesn't exist in read mode?

**Answer Choices**:
- A) IOError
- B) FileNotFoundError
- C) PermissionError
- D) ValueError

**Correct Answer**: B

**Explanation**:
`FileNotFoundError` is the specific exception raised when you attempt to open a non-existent file in read mode. While `IOError` is a general I/O exception, Python 3 uses more specific exceptions like `FileNotFoundError` (which is actually a subclass of `OSError`). `PermissionError` is raised when you lack permissions to access a file that exists. `ValueError` is raised for invalid arguments, not for file access issues.

---

### Question 6: Writing Files - Code Output
**Type**: code-output
**Difficulty**: 5
**Tags**: python, file-io, writing-files, write-method

**Question Text**:
After running this code, what will be the contents of 'output.txt'?

```python
with open('output.txt', 'w') as file:
    file.write('Hello')
    file.write('World')
```

**Answer Choices**:
- A) Hello World (with a space)
- B) HelloWorld (no space)
- C) Hello\nWorld (on separate lines)
- D) ['Hello', 'World']

**Correct Answer**: B

**Explanation**:
The `write()` method writes the exact string you give it—it doesn't add spaces, newlines, or any other separators. Since we called `write('Hello')` followed by `write('World')`, the result is "HelloWorld" with no space between them. If you want a space, you'd need to write `file.write('Hello ')` or use separate writes with explicit spaces/newlines. Option C would require writing `'Hello\nWorld'` or calling write with '\n'.

---

### Question 7: File Positioning - seek() and tell()
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: python, file-io, seek, tell, file-positioning

**Question Text**:
What does `file.seek(0)` do?

**Answer Choices**:
- A) Moves the file pointer to the end of the file
- B) Moves the file pointer to the beginning of the file
- C) Closes the file
- D) Deletes the file contents

**Correct Answer**: B

**Explanation**:
`file.seek(0)` moves the file pointer to the beginning of the file (position 0, which is 0 bytes from the start). This is useful when you want to read a file again after already reading it once. To move to the end, you'd use `file.seek(0, 2)`. `seek()` doesn't close files (that's `close()`) or delete contents (that's opening in 'w' mode or using truncate).

---

### Question 8: JSON vs CSV - Format Selection
**Type**: multiple-choice
**Difficulty**: 6
**Tags**: python, file-io, json, csv, data-formats

**Question Text**:
Which file format (JSON or CSV) would be most appropriate for storing hierarchical/nested data structures like a dictionary containing lists of dictionaries?

**Answer Choices**:
- A) CSV, because it's more compact
- B) JSON, because it naturally supports nested structures
- C) Either one works equally well
- D) CSV, because it's easier to read

**Correct Answer**: B

**Explanation**:
JSON is ideal for hierarchical or nested data structures because it's designed to represent complex objects with nesting, arrays, and different data types. CSV is fundamentally flat/tabular—it represents data in rows and columns like a spreadsheet. While you could serialize nested data to CSV, you'd lose the structure and need complex parsing logic to reconstruct it. JSON preserves the structure naturally. For simple tabular data (like a list of students with name, age, grade), CSV is perfect. For complex nested structures, JSON is the clear winner.

---

### Question 9: pathlib vs os.path - Modern Practices
**Type**: multiple-choice
**Difficulty**: 6
**Tags**: python, file-io, pathlib, os-path, file-paths

**Question Text**:
What is the main advantage of using `pathlib.Path` over `os.path` for file path operations?

**Answer Choices**:
- A) pathlib is faster than os.path
- B) pathlib provides an object-oriented interface with more intuitive syntax
- C) os.path is deprecated and will be removed in future Python versions
- D) pathlib works on Windows while os.path doesn't

**Correct Answer**: B

**Explanation**:
The main advantage of `pathlib.Path` is its object-oriented, intuitive interface. Instead of calling separate functions like `os.path.join()`, `os.path.exists()`, etc., you work with Path objects that have methods and support operators like `/` for joining paths. For example, `Path('data') / 'file.txt'` is more readable than `os.path.join('data', 'file.txt')`. Option A is incorrect—performance is similar. Option C is false—os.path is not deprecated and both modules coexist. Option D is wrong—both work on all platforms.

---

### Question 10: Reading Large Files - Memory Efficiency
**Type**: multiple-choice
**Difficulty**: 7
**Tags**: python, file-io, reading-files, memory-management, best-practices

**Question Text**:
You need to process a 5GB log file line by line. Which approach is most memory-efficient?

```python
# Option A
with open('huge.log', 'r') as file:
    for line in file:
        process(line)

# Option B
with open('huge.log', 'r') as file:
    lines = file.readlines()
    for line in lines:
        process(line)

# Option C
with open('huge.log', 'r') as file:
    content = file.read()
    for line in content.split('\n'):
        process(line)
```

**Answer Choices**:
- A) Option A - iterating over the file object
- B) Option B - using readlines()
- C) Option C - using read() and split()
- D) All three options use the same amount of memory

**Correct Answer**: A

**Explanation**:
Option A is the most memory-efficient because iterating directly over a file object reads one line at a time into memory. It's essentially streaming the file. Option B loads ALL lines into memory at once as a list (5GB!), which could cause memory issues or even crash the program. Option C is even worse—it loads the entire file into memory as a single string AND creates a list from splitting it. For large files, always iterate over the file object directly. This is a critical distinction for real-world applications dealing with large datasets.
