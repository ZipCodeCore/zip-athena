# Java File I/O

Hey there! Welcome to the world of File I/O (Input/Output) in Java. If you've ever wondered how programs save data, read configuration files, or work with documents, you're about to find out! File I/O is all about reading from and writing to files on your computer's disk. It's a crucial skill that transforms your programs from temporary calculations into systems that can persist and retrieve data.

Java gives you tons of options for file operations—some old-school, some modern, some simple, some powerful. Don't worry if it seems overwhelming at first; we'll break it down step by step, starting with the basics and working our way up to the cool new stuff.

## 14.1 File and Directory Operations

Before you start reading and writing files, you need to know how to work with them at a basic level—checking if they exist, creating directories, listing files, and so on. The `File` class has been Java's traditional way of handling these operations.

The `File` class represents a file or directory path on your filesystem. Think of it as a pointer to a location on disk—it doesn't actually contain the file's contents, just information about the file itself.

```java
import java.io.File;
import java.io.IOException;

public class FileBasics {
    public static void main(String[] args) {
        // Create a File object (doesn't create the actual file yet!)
        File myFile = new File("example.txt");
        
        // Check if the file exists
        if (myFile.exists()) {
            System.out.println("File exists!");
            System.out.println("Size: " + myFile.length() + " bytes");
            System.out.println("Can read: " + myFile.canRead());
            System.out.println("Can write: " + myFile.canWrite());
        } else {
            System.out.println("File doesn't exist yet.");
        }
        
        // Create the file
        try {
            if (myFile.createNewFile()) {
                System.out.println("File created: " + myFile.getName());
            }
        } catch (IOException e) {
            System.out.println("Error creating file: " + e.getMessage());
        }
    }
}
```

Notice the try-catch block? File operations can fail for many reasons (no permission, disk full, etc.), so Java requires you to handle potential `IOException`s.

Working with directories is similar. You can create directories, list their contents, and navigate the file system:

```java
import java.io.File;

public class DirectoryOperations {
    public static void main(String[] args) {
        // Create a directory
        File myDir = new File("myFolder");
        if (myDir.mkdir()) {
            System.out.println("Directory created!");
        }
        
        // Create nested directories
        File nestedDir = new File("parent/child/grandchild");
        if (nestedDir.mkdirs()) {  // mkdirs creates parent directories too
            System.out.println("Nested directories created!");
        }
        
        // List files in current directory
        File currentDir = new File(".");
        String[] files = currentDir.list();
        System.out.println("Files in current directory:");
        for (String fileName : files) {
            System.out.println("  " + fileName);
        }
        
        // Get more information about each file
        File[] fileObjects = currentDir.listFiles();
        for (File file : fileObjects) {
            String type = file.isDirectory() ? "DIR" : "FILE";
            System.out.println(type + ": " + file.getName());
        }
    }
}
```

The `File` class also lets you delete files, rename them, and get path information. Just remember: `File` is about file metadata and filesystem operations, not reading/writing the actual content.

## 14.2 Character Streams

Character streams are perfect for reading and writing text files. They handle character encoding automatically, so you don't have to worry about converting bytes to characters yourself. This makes them ideal for configuration files, logs, and any human-readable text.

The basic character streams are `FileReader` and `FileWriter`. They work, but they're not very efficient because they read/write one character at a time. That's where the buffered versions come in—they read/write in chunks, which is much faster.

```java
import java.io.*;

public class CharacterStreamExample {
    public static void main(String[] args) {
        // Writing to a file using FileWriter
        try (FileWriter writer = new FileWriter("message.txt")) {
            writer.write("Hello, File I/O!\n");
            writer.write("This is pretty cool.\n");
            System.out.println("File written successfully!");
        } catch (IOException e) {
            System.out.println("Error writing file: " + e.getMessage());
        }
        
        // Reading from a file using FileReader
        try (FileReader reader = new FileReader("message.txt")) {
            int character;
            System.out.println("Reading file:");
            while ((character = reader.read()) != -1) {
                System.out.print((char) character);
            }
        } catch (IOException e) {
            System.out.println("Error reading file: " + e.getMessage());
        }
    }
}
```

Notice the try-with-resources syntax? That `try (...)` automatically closes the file when you're done, even if an error occurs. Super handy!

Now let's look at the buffered versions, which are what you'll use most of the time in real code:

```java
import java.io.*;

public class BufferedCharacterStreams {
    public static void main(String[] args) {
        // Writing with BufferedWriter (much more efficient!)
        try (BufferedWriter writer = new BufferedWriter(new FileWriter("log.txt"))) {
            writer.write("Log Entry 1");
            writer.newLine();  // Platform-independent newline
            writer.write("Log Entry 2");
            writer.newLine();
            writer.write("Log Entry 3");
            writer.newLine();
        } catch (IOException e) {
            System.out.println("Error writing: " + e.getMessage());
        }
        
        // Reading with BufferedReader (can read whole lines at once!)
        try (BufferedReader reader = new BufferedReader(new FileReader("log.txt"))) {
            String line;
            System.out.println("Reading log file:");
            while ((line = reader.readLine()) != null) {
                System.out.println(line);
            }
        } catch (IOException e) {
            System.out.println("Error reading: " + e.getMessage());
        }
    }
}
```

The `BufferedReader.readLine()` method is incredibly useful—it reads an entire line at once and returns `null` when there's no more to read. The `BufferedWriter.newLine()` method is great too because it uses the correct newline character for your operating system (\n on Unix/Mac, \r\n on Windows).

Here's a practical example that copies a text file line by line:

```java
import java.io.*;

public class FileCopier {
    public static void main(String[] args) {
        try (BufferedReader reader = new BufferedReader(new FileReader("input.txt"));
             BufferedWriter writer = new BufferedWriter(new FileWriter("output.txt"))) {
            
            String line;
            int lineCount = 0;
            while ((line = reader.readLine()) != null) {
                writer.write(line);
                writer.newLine();
                lineCount++;
            }
            System.out.println("Copied " + lineCount + " lines successfully!");
            
        } catch (IOException e) {
            System.out.println("Error: " + e.getMessage());
        }
    }
}
```

## 14.3 Byte Streams

While character streams are great for text, byte streams are what you need for binary data—images, audio files, compiled programs, or any non-text data. Byte streams work with raw bytes instead of characters, giving you complete control over the data.

The basic byte streams are `FileInputStream` and `FileOutputStream`. Like their character stream cousins, they work better when wrapped in buffered versions.

```java
import java.io.*;

public class ByteStreamExample {
    public static void main(String[] args) {
        // Writing bytes to a file
        try (FileOutputStream fos = new FileOutputStream("data.bin")) {
            byte[] data = {72, 101, 108, 108, 111};  // "Hello" in ASCII
            fos.write(data);
            System.out.println("Binary data written!");
        } catch (IOException e) {
            System.out.println("Error: " + e.getMessage());
        }
        
        // Reading bytes from a file
        try (FileInputStream fis = new FileInputStream("data.bin")) {
            int byteData;
            System.out.print("Reading bytes: ");
            while ((byteData = fis.read()) != -1) {
                System.out.print((char) byteData);
            }
            System.out.println();
        } catch (IOException e) {
            System.out.println("Error: " + e.getMessage());
        }
    }
}
```

Now let's use the buffered versions, which are much faster for larger files:

```java
import java.io.*;

public class BufferedByteStreams {
    public static void main(String[] args) {
        // Copy an image file using buffered byte streams
        String sourceFile = "photo.jpg";
        String destFile = "photo_copy.jpg";
        
        try (BufferedInputStream bis = new BufferedInputStream(new FileInputStream(sourceFile));
             BufferedOutputStream bos = new BufferedOutputStream(new FileOutputStream(destFile))) {
            
            byte[] buffer = new byte[8192];  // 8KB buffer
            int bytesRead;
            
            while ((bytesRead = bis.read(buffer)) != -1) {
                bos.write(buffer, 0, bytesRead);
            }
            
            System.out.println("File copied successfully!");
            
        } catch (IOException e) {
            System.out.println("Error: " + e.getMessage());
        }
    }
}
```

The key difference here is that we're reading into a byte array buffer rather than one byte at a time. This is way more efficient—imagine reading an 8KB chunk versus reading 8,192 individual bytes!

Here's a practical example that demonstrates both reading and writing with data interpretation:

```java
import java.io.*;

public class BinaryDataHandler {
    public static void main(String[] args) {
        // Write some binary data
        try (DataOutputStream dos = new DataOutputStream(
                new BufferedOutputStream(new FileOutputStream("numbers.dat")))) {
            
            dos.writeInt(42);
            dos.writeDouble(3.14159);
            dos.writeBoolean(true);
            dos.writeUTF("Binary data rocks!");
            System.out.println("Data written!");
            
        } catch (IOException e) {
            System.out.println("Error writing: " + e.getMessage());
        }
        
        // Read the binary data back
        try (DataInputStream dis = new DataInputStream(
                new BufferedInputStream(new FileInputStream("numbers.dat")))) {
            
            int number = dis.readInt();
            double pi = dis.readDouble();
            boolean flag = dis.readBoolean();
            String message = dis.readUTF();
            
            System.out.println("Read data:");
            System.out.println("  Integer: " + number);
            System.out.println("  Double: " + pi);
            System.out.println("  Boolean: " + flag);
            System.out.println("  String: " + message);
            
        } catch (IOException e) {
            System.out.println("Error reading: " + e.getMessage());
        }
    }
}
```

`DataInputStream` and `DataOutputStream` are super useful because they let you read and write primitive data types directly in binary format. Just remember: you must read the data in the same order you wrote it!

## 14.4 Object Serialization

Object serialization is like magic—it lets you save entire Java objects to a file and recreate them later. This is perfect for saving game state, user preferences, or any complex data structures.

For a class to be serializable, it just needs to implement the `Serializable` interface. This is a marker interface (no methods to implement), but it tells Java, "Hey, it's okay to save instances of this class to disk."

```java
import java.io.Serializable;

public class Player implements Serializable {
    private static final long serialVersionUID = 1L;  // Version control
    
    private String name;
    private int level;
    private int health;
    private transient int sessionScore;  // transient = don't serialize this
    
    public Player(String name, int level, int health) {
        this.name = name;
        this.level = level;
        this.health = health;
        this.sessionScore = 0;
    }
    
    @Override
    public String toString() {
        return "Player{name='" + name + "', level=" + level + 
               ", health=" + health + ", sessionScore=" + sessionScore + "}";
    }
    
    // Getters and setters...
    public String getName() { return name; }
    public int getLevel() { return level; }
    public int getHealth() { return health; }
    public void setHealth(int health) { this.health = health; }
}
```

The `serialVersionUID` is important—it helps Java verify that the saved object matches the current class definition. The `transient` keyword marks fields that shouldn't be saved (like temporary session data).

Now let's save and load objects:

```java
import java.io.*;

public class SerializationExample {
    public static void main(String[] args) {
        // Create and save a player
        Player player = new Player("Hero", 10, 100);
        
        // Serialize (save) the object
        try (ObjectOutputStream oos = new ObjectOutputStream(
                new FileOutputStream("player.dat"))) {
            
            oos.writeObject(player);
            System.out.println("Player saved: " + player);
            
        } catch (IOException e) {
            System.out.println("Error saving: " + e.getMessage());
        }
        
        // Deserialize (load) the object
        try (ObjectInputStream ois = new ObjectInputStream(
                new FileInputStream("player.dat"))) {
            
            Player loadedPlayer = (Player) ois.readObject();
            System.out.println("Player loaded: " + loadedPlayer);
            
        } catch (IOException | ClassNotFoundException e) {
            System.out.println("Error loading: " + e.getMessage());
        }
    }
}
```

Notice that reading an object requires catching `ClassNotFoundException` too—that happens if the class definition isn't available when deserializing.

You can also serialize collections of objects:

```java
import java.io.*;
import java.util.*;

public class CollectionSerialization {
    public static void main(String[] args) {
        // Create a list of players
        List<Player> team = new ArrayList<>();
        team.add(new Player("Warrior", 15, 120));
        team.add(new Player("Mage", 12, 80));
        team.add(new Player("Archer", 10, 90));
        
        // Save the entire list
        try (ObjectOutputStream oos = new ObjectOutputStream(
                new FileOutputStream("team.dat"))) {
            
            oos.writeObject(team);
            System.out.println("Team saved with " + team.size() + " players");
            
        } catch (IOException e) {
            System.out.println("Error: " + e.getMessage());
        }
        
        // Load the list
        try (ObjectInputStream ois = new ObjectInputStream(
                new FileInputStream("team.dat"))) {
            
            @SuppressWarnings("unchecked")
            List<Player> loadedTeam = (List<Player>) ois.readObject();
            System.out.println("Team loaded with " + loadedTeam.size() + " players:");
            for (Player p : loadedTeam) {
                System.out.println("  " + p);
            }
            
        } catch (IOException | ClassNotFoundException e) {
            System.out.println("Error: " + e.getMessage());
        }
    }
}
```

Serialization is incredibly convenient, but remember: it's Java-specific. If you need to share data with other languages or systems, use a standard format like JSON or XML instead.

## 14.5 NIO (New I/O) Basics

Java NIO (New I/O), introduced in Java 7, is the modern way to work with files. It's more powerful, more flexible, and often more convenient than the old File class. The main classes you'll use are `Path`, `Paths`, and `Files`.

A `Path` represents a file or directory path (similar to the old `File` class). You create paths using the `Paths.get()` factory method:

```java
import java.nio.file.*;
import java.io.IOException;

public class PathBasics {
    public static void main(String[] args) {
        // Create a Path
        Path filePath = Paths.get("example.txt");
        Path dirPath = Paths.get("myFolder", "subFolder", "file.txt");
        
        System.out.println("File name: " + filePath.getFileName());
        System.out.println("Parent: " + filePath.getParent());
        System.out.println("Absolute path: " + filePath.toAbsolutePath());
        
        // Path operations
        Path path1 = Paths.get("folder/file.txt");
        Path path2 = Paths.get("folder");
        Path path3 = path2.resolve("file.txt");  // Combine paths
        
        System.out.println("Combined: " + path3);
        System.out.println("Are they equal? " + path1.equals(path3));
    }
}
```

The `Files` class is where the real power is—it has static methods for almost every file operation you can imagine:

```java
import java.nio.file.*;
import java.io.IOException;
import java.util.List;

public class FilesClassDemo {
    public static void main(String[] args) {
        Path file = Paths.get("test.txt");
        
        try {
            // Create a file
            if (!Files.exists(file)) {
                Files.createFile(file);
                System.out.println("File created!");
            }
            
            // Write to a file (super easy!)
            Files.writeString(file, "Hello, NIO!\nThis is much simpler!");
            System.out.println("Content written!");
            
            // Read from a file (entire file at once!)
            String content = Files.readString(file);
            System.out.println("Content: " + content);
            
            // Read all lines into a List
            List<String> lines = Files.readAllLines(file);
            System.out.println("Number of lines: " + lines.size());
            for (String line : lines) {
                System.out.println("  > " + line);
            }
            
            // Get file attributes
            System.out.println("Size: " + Files.size(file) + " bytes");
            System.out.println("Is regular file? " + Files.isRegularFile(file));
            System.out.println("Is readable? " + Files.isReadable(file));
            
        } catch (IOException e) {
            System.out.println("Error: " + e.getMessage());
        }
    }
}
```

See how simple that is? `Files.readString()` and `Files.writeString()` are incredibly convenient for working with text files. No need for readers, writers, or buffers!

Here's a more comprehensive example showing directory operations and file copying:

```java
import java.nio.file.*;
import java.io.IOException;
import java.util.stream.Stream;

public class AdvancedNIOOperations {
    public static void main(String[] args) {
        try {
            // Create directories
            Path dir = Paths.get("testDir");
            if (!Files.exists(dir)) {
                Files.createDirectory(dir);
                System.out.println("Directory created!");
            }
            
            // Create a file and write content
            Path sourceFile = dir.resolve("source.txt");
            Files.writeString(sourceFile, "This file will be copied!");
            
            // Copy the file
            Path destFile = dir.resolve("destination.txt");
            Files.copy(sourceFile, destFile, StandardCopyOption.REPLACE_EXISTING);
            System.out.println("File copied!");
            
            // List files in directory
            System.out.println("\nFiles in directory:");
            try (Stream<Path> paths = Files.list(dir)) {
                paths.forEach(path -> System.out.println("  " + path.getFileName()));
            }
            
            // Delete files
            Files.deleteIfExists(destFile);
            System.out.println("\nDestination file deleted!");
            
            // Walk directory tree (careful with this!)
            System.out.println("\nDirectory tree:");
            try (Stream<Path> paths = Files.walk(dir)) {
                paths.forEach(path -> System.out.println("  " + path));
            }
            
        } catch (IOException e) {
            System.out.println("Error: " + e.getMessage());
        }
    }
}
```

NIO also supports watching directories for changes, symbolic links, and advanced file attributes. It's powerful stuff! For most day-to-day file operations, NIO is now the preferred approach.

## 14.6 Reading and Writing Different File Formats

Real-world applications work with many file formats—CSV, JSON, XML, properties files, and more. Let's look at how to handle some common formats.

**Properties Files** are simple key-value configuration files. Java has built-in support:

```java
import java.io.*;
import java.util.Properties;

public class PropertiesExample {
    public static void main(String[] args) {
        Properties props = new Properties();
        
        // Write properties
        try (FileOutputStream out = new FileOutputStream("config.properties")) {
            props.setProperty("database.url", "jdbc:mysql://localhost:3306/mydb");
            props.setProperty("database.user", "admin");
            props.setProperty("max.connections", "10");
            props.store(out, "Database Configuration");
            System.out.println("Properties saved!");
        } catch (IOException e) {
            System.out.println("Error writing: " + e.getMessage());
        }
        
        // Read properties
        Properties loadedProps = new Properties();
        try (FileInputStream in = new FileInputStream("config.properties")) {
            loadedProps.load(in);
            
            System.out.println("\nLoaded properties:");
            System.out.println("URL: " + loadedProps.getProperty("database.url"));
            System.out.println("User: " + loadedProps.getProperty("database.user"));
            System.out.println("Max connections: " + loadedProps.getProperty("max.connections"));
            
            // Get with default value
            String timeout = loadedProps.getProperty("timeout", "30");
            System.out.println("Timeout (with default): " + timeout);
            
        } catch (IOException e) {
            System.out.println("Error reading: " + e.getMessage());
        }
    }
}
```

**CSV Files** (Comma-Separated Values) are commonly used for data exchange. While there are libraries for complex CSV parsing, simple CSV files can be handled manually:

```java
import java.io.*;
import java.nio.file.*;
import java.util.*;

public class CSVExample {
    public static void main(String[] args) {
        // Write CSV
        List<String[]> data = new ArrayList<>();
        data.add(new String[]{"Name", "Age", "City"});
        data.add(new String[]{"Alice", "30", "New York"});
        data.add(new String[]{"Bob", "25", "Los Angeles"});
        data.add(new String[]{"Charlie", "35", "Chicago"});
        
        try (BufferedWriter writer = Files.newBufferedWriter(Paths.get("people.csv"))) {
            for (String[] row : data) {
                writer.write(String.join(",", row));
                writer.newLine();
            }
            System.out.println("CSV file written!");
        } catch (IOException e) {
            System.out.println("Error writing: " + e.getMessage());
        }
        
        // Read CSV
        try (BufferedReader reader = Files.newBufferedReader(Paths.get("people.csv"))) {
            String line;
            System.out.println("\nReading CSV:");
            while ((line = reader.readLine()) != null) {
                String[] columns = line.split(",");
                System.out.println("  " + Arrays.toString(columns));
            }
        } catch (IOException e) {
            System.out.println("Error reading: " + e.getMessage());
        }
    }
}
```

**JSON Files** require an external library like Jackson or Gson, but here's what it looks like conceptually with manual parsing for simple cases:

```java
import java.io.*;
import java.nio.file.*;

public class SimpleJSONExample {
    public static void main(String[] args) {
        // Write simple JSON manually
        String json = "{\n" +
                      "  \"name\": \"John Doe\",\n" +
                      "  \"age\": 30,\n" +
                      "  \"email\": \"john@example.com\"\n" +
                      "}";
        
        try {
            Files.writeString(Paths.get("user.json"), json);
            System.out.println("JSON file written!");
            
            // Read JSON
            String content = Files.readString(Paths.get("user.json"));
            System.out.println("\nJSON content:");
            System.out.println(content);
            
        } catch (IOException e) {
            System.out.println("Error: " + e.getMessage());
        }
    }
}
```

For production code with JSON, definitely use a library like Jackson:
```java
// Note: This requires Jackson library (jackson-databind)
// Add to Maven: com.fasterxml.jackson.core:jackson-databind:2.15.0

import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.*;

class User {
    public String name;
    public int age;
    public String email;
}

public class JacksonExample {
    public static void main(String[] args) {
        ObjectMapper mapper = new ObjectMapper();
        
        // Write JSON
        try {
            User user = new User();
            user.name = "John Doe";
            user.age = 30;
            user.email = "john@example.com";
            
            mapper.writeValue(new File("user.json"), user);
            System.out.println("JSON written with Jackson!");
            
            // Read JSON
            User loadedUser = mapper.readValue(new File("user.json"), User.class);
            System.out.println("Loaded: " + loadedUser.name + ", " + loadedUser.age);
            
        } catch (IOException e) {
            System.out.println("Error: " + e.getMessage());
        }
    }
}
```

**Text Files with Different Encodings** sometimes need special handling:

```java
import java.io.*;
import java.nio.file.*;
import java.nio.charset.*;

public class EncodingExample {
    public static void main(String[] args) {
        try {
            // Write with UTF-8 encoding
            String text = "Hello, 世界! 🌍";
            Files.writeString(Paths.get("utf8.txt"), text, StandardCharsets.UTF_8);
            
            // Read with UTF-8 encoding
            String content = Files.readString(Paths.get("utf8.txt"), StandardCharsets.UTF_8);
            System.out.println("UTF-8 content: " + content);
            
            // Write with different encoding
            Files.writeString(Paths.get("iso.txt"), "Hello World", StandardCharsets.ISO_8859_1);
            
            System.out.println("Files written with different encodings!");
            
        } catch (IOException e) {
            System.out.println("Error: " + e.getMessage());
        }
    }
}
```

## Common Pitfalls

**Not Closing Resources**: Always close streams! Use try-with-resources to do it automatically:
```java
// BAD - might leave file open if error occurs
FileWriter writer = new FileWriter("file.txt");
writer.write("data");
writer.close();

// GOOD - automatically closes
try (FileWriter writer = new FileWriter("file.txt")) {
    writer.write("data");
}
```

**Reading and Writing Order**: When using `DataInputStream`/`DataOutputStream` or serialization, you must read in the same order you wrote:
```java
// Write: int, then String
dos.writeInt(42);
dos.writeUTF("Hello");

// Must read: int, then String (same order!)
int num = dis.readInt();
String msg = dis.readUTF();
```

**Assuming Files Exist**: Always check or handle the exception when a file might not exist:
```java
if (Files.exists(path)) {
    // Safe to read
} else {
    // Handle missing file
}
```

**Platform-Specific Paths**: Use `File.separator` or `Path` operations instead of hardcoding slashes:
```java
// BAD - won't work on Windows
String path = "folder/subfolder/file.txt";

// GOOD - works everywhere
Path path = Paths.get("folder", "subfolder", "file.txt");
```

## Best Practices

1. **Use try-with-resources** for all I/O operations to ensure resources are closed
2. **Prefer NIO for new code** - it's more modern and flexible
3. **Buffer your streams** - always wrap with Buffered classes for better performance
4. **Handle exceptions properly** - file operations can fail in many ways
5. **Use appropriate streams** - character streams for text, byte streams for binary
6. **Close streams in finally blocks** if not using try-with-resources
7. **Validate paths** before operations to give better error messages
8. **Use `Files.exists()` carefully** - it's not always reliable for concurrent access

## Practice Exercises

Try these to reinforce your learning:

1. Write a program that copies a text file while adding line numbers to each line.

2. Create a simple address book that serializes a list of Contact objects (name, phone, email) to a file.

3. Write a program that reads a CSV file of student grades and calculates the average for each student.

4. Build a log file analyzer that reads a log file and counts how many ERROR vs WARNING vs INFO messages appear.

5. Create a configuration manager using Properties that can save and load application settings.

6. Write a program using NIO that recursively copies an entire directory tree.

## Wrapping Up

File I/O in Java is a rich topic with many options. Start with the basics—reading and writing text files with buffered streams or NIO—and gradually explore more advanced features as you need them. The key is understanding when to use character vs. byte streams, when to serialize objects, and how to leverage the modern NIO API for cleaner code.

Don't be intimidated by all the classes and options. In practice, you'll use a few patterns repeatedly: BufferedReader/Writer for text, BufferedInputStream/OutputStream for binary, and the Files class for most operations in modern code.

Keep practicing, and soon file operations will feel natural. Now head over to the quiz section and test your knowledge!

---

## Quiz Questions

This section contains quiz questions that test understanding of Java File I/O concepts. Each question is formatted for easy extraction and database import.

### Question 1: File Class - File Existence Check
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: java, file-io, file-operations

**Question Text**:
What happens when you create a `File` object in Java?

```java
File myFile = new File("example.txt");
```

**Answer Choices**:
- A) A new file named "example.txt" is created on disk
- B) The program checks if the file exists and throws an exception if it doesn't
- C) A File object is created that represents the path, but no actual file is created
- D) The file is created and opened for writing

**Correct Answer**: C

**Explanation**:
Creating a `File` object only creates an object in memory that represents a path to a file or directory. It does NOT create an actual file on disk. To create the actual file, you need to call methods like `createNewFile()` or use streams. The File object is just a wrapper around a path string with methods to query and manipulate file metadata.

---

### Question 2: Character Streams - BufferedReader Efficiency
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, file-io, character-streams, buffering

**Question Text**:
Why is `BufferedReader` more efficient than `FileReader` for reading large text files?

**Answer Choices**:
- A) BufferedReader uses multi-threading to read faster
- B) BufferedReader reads chunks of data into memory at once, reducing disk I/O operations
- C) BufferedReader compresses the file before reading it
- D) BufferedReader skips whitespace automatically

**Correct Answer**: B

**Explanation**:
`BufferedReader` improves performance by reading large chunks (buffers) of data from disk into memory at once, rather than making a separate disk I/O operation for each character. Disk I/O is slow, so reducing the number of operations significantly improves performance. It doesn't use multi-threading, compression, or whitespace skipping—it's all about batching I/O operations.

---

### Question 3: Character Streams - Reading Lines
**Type**: code-output
**Difficulty**: 3
**Tags**: java, file-io, character-streams, buffered-reader

**Question Text**:
Given a file "data.txt" with the following content:
```
Line 1
Line 2
Line 3
```

What will this code print?

```java
try (BufferedReader reader = new BufferedReader(new FileReader("data.txt"))) {
    String line;
    int count = 0;
    while ((line = reader.readLine()) != null) {
        count++;
    }
    System.out.println(count);
} catch (IOException e) {
    System.out.println("Error");
}
```

**Answer Choices**:
- A) 0
- B) 3
- C) Line 1Line 2Line 3
- D) Error

**Correct Answer**: B

**Explanation**:
The `readLine()` method reads one line at a time and returns `null` when there are no more lines. The while loop counts each line read, so with three lines in the file, `count` will be 3. The code doesn't print the lines themselves, just the count. If the file existed and was readable, no error would occur.

---

### Question 4: Byte Streams - Use Case
**Type**: multiple-choice
**Difficulty**: 3
**Tags**: java, file-io, byte-streams

**Question Text**:
Which of the following file types would be BEST handled using byte streams rather than character streams?

**Answer Choices**:
- A) A text log file with application errors
- B) A CSV file with comma-separated values
- C) A JPEG image file
- D) An XML configuration file

**Correct Answer**: C

**Explanation**:
Byte streams are designed for binary data, while character streams are for text. A JPEG image is binary data, so it must be handled with byte streams (like `FileInputStream` or `BufferedInputStream`). The other options are all text formats that would be better handled with character streams, which properly handle character encoding and make reading text more convenient.

---

### Question 5: Object Serialization - Transient Keyword
**Type**: true-false
**Difficulty**: 4
**Tags**: java, file-io, serialization, transient

**Question Text**:
When an object is serialized, fields marked with the `transient` keyword are included in the serialized data.

**Answer Choices**:
- A) True
- B) False

**Correct Answer**: B

**Explanation**:
False. Fields marked as `transient` are explicitly excluded from serialization. This is useful for fields that don't need to be saved (like cached values, temporary session data, or fields that can be recalculated). When the object is deserialized, transient fields will have their default values (null for objects, 0 for numbers, false for booleans, etc.).

---

### Question 6: Object Serialization - Interface Requirement
**Type**: code-output
**Difficulty**: 4
**Tags**: java, file-io, serialization, serializable

**Question Text**:
What will happen when you try to serialize an object of this class?

```java
public class Person {
    private String name;
    private int age;
    
    public Person(String name, int age) {
        this.name = name;
        this.age = age;
    }
}

// Later...
Person person = new Person("Alice", 30);
ObjectOutputStream oos = new ObjectOutputStream(new FileOutputStream("person.dat"));
oos.writeObject(person);
```

**Answer Choices**:
- A) The object will be successfully serialized
- B) A NotSerializableException will be thrown
- C) The file will be created but will be empty
- D) A ClassNotFoundException will be thrown

**Correct Answer**: B

**Explanation**:
A `NotSerializableException` will be thrown because the `Person` class doesn't implement the `Serializable` interface. For a class to be serializable, it must explicitly implement `Serializable` (even though it's a marker interface with no methods). This is a safety feature—Java won't serialize objects unless you explicitly mark them as safe to serialize.

---

### Question 7: NIO Files Class - Reading Entire File
**Type**: code-output
**Difficulty**: 3
**Tags**: java, file-io, nio, files-class

**Question Text**:
What does `Files.readString()` return when reading a file?

```java
Path file = Paths.get("message.txt");
// Assume message.txt contains: "Hello\nWorld"
String content = Files.readString(file);
```

**Answer Choices**:
- A) An array of strings, one per line
- B) A single string containing the entire file content including newlines
- C) A List<String> with each line as an element
- D) Only the first line of the file

**Correct Answer**: B

**Explanation**:
`Files.readString()` reads the entire file content into a single String, including all newline characters. For the example file, it would return "Hello\nWorld" (one string with a newline in the middle). If you want lines as separate elements, use `Files.readAllLines()` instead, which returns a `List<String>`. This method is very convenient for small to medium files.

---

### Question 8: Try-with-Resources - Resource Management
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, file-io, try-with-resources, exception-handling

**Question Text**:
What is the primary advantage of using try-with-resources for file I/O?

```java
try (BufferedReader reader = new BufferedReader(new FileReader("file.txt"))) {
    // Read file
}
```

**Answer Choices**:
- A) It makes the code run faster
- B) It automatically closes the resource even if an exception occurs
- C) It prevents IOException from being thrown
- D) It allows reading files that don't exist

**Correct Answer**: B

**Explanation**:
The primary advantage of try-with-resources is automatic resource management. The resource (BufferedReader in this case) is automatically closed when the try block exits, whether normally or due to an exception. This prevents resource leaks where files remain open. It doesn't improve performance, doesn't prevent exceptions, and doesn't allow reading non-existent files—it just ensures proper cleanup.

---

### Question 9: File Formats - Properties Files
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, file-io, properties, configuration

**Question Text**:
What method would you use to read a property value from a `Properties` object with a default value if the property doesn't exist?

**Answer Choices**:
- A) `properties.get("key", "default")`
- B) `properties.getProperty("key", "default")`
- C) `properties.read("key") ?? "default"`
- D) `properties.fetch("key").orElse("default")`

**Correct Answer**: B

**Explanation**:
The `getProperty(String key, String defaultValue)` method returns the property value if it exists, or the default value if it doesn't. This is useful for optional configuration values. Option A is incorrect because `get()` returns an Object without a default parameter. Options C and D use syntax that doesn't exist in Java's Properties class.

---

### Question 10: Byte Streams - Buffer Size
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, file-io, byte-streams, performance

**Question Text**:
When copying a large file using buffered byte streams, what is the purpose of the byte array buffer?

```java
byte[] buffer = new byte[8192];
int bytesRead;
while ((bytesRead = bis.read(buffer)) != -1) {
    bos.write(buffer, 0, bytesRead);
}
```

**Answer Choices**:
- A) To store the entire file in memory before writing
- B) To read and write data in chunks, reducing the number of I/O operations
- C) To compress the data before writing to disk
- D) To encrypt the data during transfer

**Correct Answer**: B

**Explanation**:
The buffer allows reading and writing data in chunks (8KB in this case) rather than one byte at a time. This dramatically reduces the number of I/O operations, which are expensive. The buffer doesn't store the entire file (just 8KB at a time), doesn't compress or encrypt data—it simply batches I/O operations for efficiency. The `bytesRead` variable ensures we only write the actual bytes read, which is important for the last chunk which might be smaller than 8KB.

---

### Question 11: NIO Path Operations - Path Resolution
**Type**: code-output
**Difficulty**: 4
**Tags**: java, file-io, nio, path

**Question Text**:
What will be the result of this code?

```java
Path base = Paths.get("documents");
Path file = base.resolve("reports").resolve("annual.txt");
System.out.println(file);
```

**Answer Choices**:
- A) reports/annual.txt
- B) documents/reports/annual.txt
- C) documents/annual.txt
- D) annual.txt

**Correct Answer**: B

**Explanation**:
The `resolve()` method combines paths. Starting with "documents", resolving "reports" gives "documents/reports", and resolving "annual.txt" gives "documents/reports/annual.txt". The resolve method appends the argument to the current path, building a complete path from components. This is the proper way to construct paths in a platform-independent manner.

---

### Question 12: DataInputStream/DataOutputStream - Read/Write Order
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, file-io, data-streams, binary

**Question Text**:
What will happen if you write data in one order but read it in a different order?

```java
// Writing
dos.writeInt(100);
dos.writeDouble(3.14);

// Reading
double d = dis.readDouble();  // Wrong order!
int i = dis.readInt();
```

**Answer Choices**:
- A) The values will be automatically reordered
- B) An exception will be thrown immediately
- C) The data will be corrupted and you'll get incorrect values
- D) Java will return null for mismatched types

**Correct Answer**: C

**Explanation**:
When reading binary data, you must read in the exact same order and types as you wrote. If you don't, you'll get corrupted data. The first 8 bytes (the int) will be interpreted as a double, giving you a nonsensical value. No exception is thrown because Java doesn't know what order you originally wrote the data—it just reads bytes and interprets them as requested. This is a common and tricky bug in binary I/O.
