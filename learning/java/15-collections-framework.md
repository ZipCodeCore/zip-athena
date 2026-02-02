# Java Collections Framework

Welcome to one of the most powerful and practical parts of Java! The Collections Framework is like having a Swiss Army knife for managing groups of objects. Instead of wrestling with arrays and writing your own data structure code from scratch, Java gives you a whole toolbox of pre-built, tested, and optimized collection classes. Trust me, once you get comfortable with collections, you'll wonder how you ever lived without them!

## 15.1 Introduction to Collections Framework

The Collections Framework is a unified architecture for representing and manipulating collections of objects. Think of it as a set of standardized tools that help you store, retrieve, manipulate, and communicate groups of data efficiently.

Before the Collections Framework (way back in Java 1.1), developers had to work with simple arrays or basic classes like Vector and Hashtable. While arrays are still useful, they have limitations—they're fixed in size, and you can't easily add or remove elements. The Collections Framework solves these problems and gives you much more flexibility.

Here's what makes the framework so powerful: it provides interfaces (like List, Set, Map) that define standard behaviors, and multiple implementations of each interface (like ArrayList, HashSet, HashMap) that are optimized for different use cases. This means you can write code that works with a List interface, and you can swap out ArrayList for LinkedList later without changing your logic.

The framework also provides algorithms for common operations like sorting, searching, and shuffling. Instead of writing these algorithms yourself (and potentially making mistakes), you can use the tested, optimized versions from the Collections utility class.

### Example: Why Collections Matter

```java
import java.util.ArrayList;
import java.util.List;

public class CollectionsIntro {
    public static void main(String[] args) {
        // Old way with arrays - fixed size!
        String[] arrayNames = new String[3];
        arrayNames[0] = "Alice";
        arrayNames[1] = "Bob";
        arrayNames[2] = "Charlie";
        // What if we want to add a 4th name? We'd need a new array!
        
        // New way with ArrayList - dynamic size!
        List<String> listNames = new ArrayList<>();
        listNames.add("Alice");
        listNames.add("Bob");
        listNames.add("Charlie");
        listNames.add("Diana");  // No problem!
        listNames.remove("Bob");  // Easy removal too!
        
        System.out.println("Names in list: " + listNames);
        System.out.println("Number of names: " + listNames.size());
    }
}
```

See how much easier that is? No worrying about array sizes, and you get convenient methods like `add()`, `remove()`, and `size()` built right in.

## 15.2 Collection Hierarchy

Understanding the hierarchy helps you choose the right collection for your needs. At the top level, we have the `Collection` interface (note the singular form), which is the root of the collection hierarchy. Then we have several main branches:

**The Main Interfaces:**
- **List**: An ordered collection (also called a sequence). Lists can contain duplicate elements, and you can access elements by their position (index).
- **Set**: A collection that cannot contain duplicate elements. Sets model the mathematical set abstraction.
- **Queue**: A collection designed for holding elements prior to processing. Queues typically (but not always) order elements in FIFO (first-in-first-out) manner.
- **Map**: Not technically part of the Collection hierarchy (it doesn't extend Collection), but it's a crucial part of the framework. Maps store key-value pairs, where each key maps to exactly one value.

Here's a simple mental model: if you need ordered elements with duplicates allowed, think List. If you need unique elements and don't care about order, think Set. If you need to process items in a specific order (like a waiting line), think Queue. And if you need to look things up by a key (like a dictionary), think Map.

### Example: Choosing the Right Type

```java
import java.util.*;

public class CollectionTypes {
    public static void main(String[] args) {
        // List - allows duplicates, maintains order
        List<String> playlist = new ArrayList<>();
        playlist.add("Song A");
        playlist.add("Song B");
        playlist.add("Song A");  // Duplicate is fine!
        System.out.println("Playlist: " + playlist);
        // Output: [Song A, Song B, Song A]
        
        // Set - no duplicates allowed
        Set<String> uniqueSongs = new HashSet<>();
        uniqueSongs.add("Song A");
        uniqueSongs.add("Song B");
        uniqueSongs.add("Song A");  // This won't be added again
        System.out.println("Unique songs: " + uniqueSongs);
        // Output: [Song A, Song B] (order may vary)
        
        // Map - key-value pairs
        Map<String, Integer> songPlays = new HashMap<>();
        songPlays.put("Song A", 10);
        songPlays.put("Song B", 5);
        System.out.println("Song A plays: " + songPlays.get("Song A"));
        // Output: 10
    }
}
```

## 15.3 List Interface (ArrayList, LinkedList)

The List interface represents an ordered collection where you have precise control over where each element is inserted. You can access elements by their integer index (position in the list) and search for elements. The two most common implementations are ArrayList and LinkedList.

**ArrayList** is backed by a dynamic array. It's great for storing and accessing data randomly because getting an element by index is super fast (O(1) time complexity). However, adding or removing elements from the middle of the list is slower because elements need to be shifted. Think of ArrayList as your go-to choice for most situations—it's simple, fast for reads, and handles most use cases well.

**LinkedList** is backed by a doubly-linked list structure. It's better when you frequently add or remove elements from the beginning or middle of the list because it just needs to update a few references (no shifting required). The trade-off is that accessing elements by index is slower because it needs to traverse the list from the beginning or end. Use LinkedList when you're doing lots of insertions and deletions, especially at the beginning of the list.

Both implement the List interface, so you can write code against the List interface and swap implementations easily. This is good practice—program to an interface, not an implementation!

### Example: ArrayList vs LinkedList

```java
import java.util.*;

public class ListComparison {
    public static void main(String[] args) {
        // ArrayList - great for random access
        List<String> arrayList = new ArrayList<>();
        arrayList.add("First");
        arrayList.add("Second");
        arrayList.add("Third");
        
        // Fast access by index
        System.out.println("Element at index 1: " + arrayList.get(1));
        
        // Add at beginning - slower for ArrayList
        arrayList.add(0, "New First");
        System.out.println("ArrayList: " + arrayList);
        
        // LinkedList - great for frequent insertions/deletions
        List<String> linkedList = new LinkedList<>();
        linkedList.add("A");
        linkedList.add("B");
        linkedList.add("C");
        
        // Add at beginning - faster for LinkedList
        linkedList.add(0, "New A");
        System.out.println("LinkedList: " + linkedList);
        
        // Common List operations work on both
        System.out.println("ArrayList size: " + arrayList.size());
        System.out.println("LinkedList contains 'B': " + linkedList.contains("B"));
        
        // Removing elements
        arrayList.remove("Second");
        linkedList.remove(1);  // Remove by index
        
        System.out.println("After removal - ArrayList: " + arrayList);
        System.out.println("After removal - LinkedList: " + linkedList);
    }
}
```

## 15.4 Set Interface (HashSet, TreeSet, LinkedHashSet)

The Set interface models the mathematical set abstraction—a collection that contains no duplicate elements. If you try to add a duplicate, it simply won't be added (the add method will return false). Sets are perfect for when you need to ensure uniqueness, like tracking unique visitors to a website or storing a collection of unique tags.

**HashSet** is the most commonly used Set implementation. It uses a hash table for storage, which makes adding, removing, and checking for containment very fast (average O(1) time). However, it doesn't maintain any order—when you iterate through a HashSet, the elements can come out in any order. Use HashSet when you need fast performance and don't care about order.

**TreeSet** stores elements in sorted order (natural ordering or via a Comparator). It uses a red-black tree structure, which means operations are slower than HashSet (O(log n) time), but you get automatic sorting for free. If you need your unique elements sorted, TreeSet is your friend.

**LinkedHashSet** combines the best of both worlds—it has the fast performance of HashSet but also maintains insertion order. It uses a bit more memory (it maintains a linked list alongside the hash table), but if you need predictable iteration order with fast performance, this is the one.

### Example: Different Set Implementations

```java
import java.util.*;

public class SetExamples {
    public static void main(String[] args) {
        // HashSet - no guaranteed order
        Set<String> hashSet = new HashSet<>();
        hashSet.add("Banana");
        hashSet.add("Apple");
        hashSet.add("Cherry");
        hashSet.add("Apple");  // Duplicate - won't be added
        System.out.println("HashSet: " + hashSet);
        // Output order may vary: [Apple, Cherry, Banana] or different
        
        // TreeSet - sorted order
        Set<String> treeSet = new TreeSet<>();
        treeSet.add("Banana");
        treeSet.add("Apple");
        treeSet.add("Cherry");
        System.out.println("TreeSet: " + treeSet);
        // Output: [Apple, Banana, Cherry] - alphabetically sorted
        
        // LinkedHashSet - insertion order maintained
        Set<String> linkedHashSet = new LinkedHashSet<>();
        linkedHashSet.add("Banana");
        linkedHashSet.add("Apple");
        linkedHashSet.add("Cherry");
        System.out.println("LinkedHashSet: " + linkedHashSet);
        // Output: [Banana, Apple, Cherry] - order preserved
        
        // Common Set operations
        System.out.println("HashSet size: " + hashSet.size());
        System.out.println("Contains Apple? " + hashSet.contains("Apple"));
        
        // Set operations
        Set<Integer> set1 = new HashSet<>(Arrays.asList(1, 2, 3, 4, 5));
        Set<Integer> set2 = new HashSet<>(Arrays.asList(4, 5, 6, 7, 8));
        
        // Union
        Set<Integer> union = new HashSet<>(set1);
        union.addAll(set2);
        System.out.println("Union: " + union);  // [1, 2, 3, 4, 5, 6, 7, 8]
        
        // Intersection
        Set<Integer> intersection = new HashSet<>(set1);
        intersection.retainAll(set2);
        System.out.println("Intersection: " + intersection);  // [4, 5]
    }
}
```

## 15.5 Queue Interface (LinkedList as Queue, PriorityQueue)

The Queue interface is designed for holding elements prior to processing. Queues typically order elements in FIFO (first-in-first-out) manner, like a line at a coffee shop—the first person in line is the first to be served. However, this isn't always the case (as we'll see with PriorityQueue).

**LinkedList as Queue**: Remember LinkedList from the List section? It also implements the Queue interface! This makes it perfect for standard FIFO queues because it's efficient at adding to the end and removing from the beginning. You'll use methods like `offer()` to add elements, `poll()` to remove and return the head, and `peek()` to look at the head without removing it.

**PriorityQueue** is a special kind of queue where elements are ordered by their priority (either natural ordering or via a Comparator). The element with the highest priority (smallest value in natural ordering) is always at the head of the queue. This is super useful for scenarios like task scheduling where some tasks are more urgent than others, or implementing algorithms like Dijkstra's shortest path.

Think of a regular Queue (LinkedList) as a fair line where everyone is served in order. Think of a PriorityQueue as an emergency room where patients with more serious conditions get treated first, regardless of when they arrived.

### Example: Queue Operations

```java
import java.util.*;

public class QueueExamples {
    public static void main(String[] args) {
        // LinkedList as a FIFO Queue
        Queue<String> queue = new LinkedList<>();
        queue.offer("First Customer");
        queue.offer("Second Customer");
        queue.offer("Third Customer");
        
        System.out.println("Queue: " + queue);
        
        // Peek at the head without removing
        System.out.println("Next to serve: " + queue.peek());
        
        // Serve customers in FIFO order
        while (!queue.isEmpty()) {
            System.out.println("Serving: " + queue.poll());
        }
        
        System.out.println("\n--- Priority Queue Example ---");
        
        // PriorityQueue - orders by priority
        PriorityQueue<Integer> priorityQueue = new PriorityQueue<>();
        priorityQueue.offer(5);
        priorityQueue.offer(2);
        priorityQueue.offer(8);
        priorityQueue.offer(1);
        
        System.out.println("Priority Queue: " + priorityQueue);
        
        // Elements come out in priority order (smallest first)
        while (!priorityQueue.isEmpty()) {
            System.out.println("Processing: " + priorityQueue.poll());
        }
        // Output: 1, 2, 5, 8
        
        // Custom priority with Comparator
        PriorityQueue<String> taskQueue = new PriorityQueue<>(
            (a, b) -> b.length() - a.length()  // Longer strings have higher priority
        );
        taskQueue.offer("Short");
        taskQueue.offer("Medium task");
        taskQueue.offer("Very long task here");
        
        System.out.println("\nProcessing tasks by length:");
        while (!taskQueue.isEmpty()) {
            System.out.println(taskQueue.poll());
        }
    }
}
```

## 15.6 Map Interface (HashMap, TreeMap, LinkedHashMap)

Maps are all about key-value pairs. Instead of accessing elements by index (like in a List), you access them by a key. Think of it like a real dictionary—you look up a word (the key) to find its definition (the value). Each key maps to exactly one value, and keys must be unique (but values can be duplicated).

**HashMap** is your go-to Map implementation. It uses hashing for lightning-fast lookups, insertions, and deletions (average O(1) time). The keys are unordered—when you iterate through a HashMap, they can come out in any sequence. Use HashMap when you need fast key-value lookups and don't care about order.

**TreeMap** keeps its keys sorted (either natural ordering or via a Comparator). This is great when you need your map entries in a specific order, or when you want to perform range operations (like getting all keys between two values). The trade-off is slower operations (O(log n) time) due to the tree structure.

**LinkedHashMap** maintains insertion order (or access order if you configure it that way). It's almost as fast as HashMap but gives you predictable iteration order. Perfect for implementing LRU (Least Recently Used) caches or when you just want your map to remember the order you put things in.

One important note: Maps don't extend the Collection interface (they're separate), but they're still considered part of the Collections Framework because they work hand-in-hand with collections.

### Example: Map Implementations

```java
import java.util.*;

public class MapExamples {
    public static void main(String[] args) {
        // HashMap - fast, unordered
        Map<String, Integer> scores = new HashMap<>();
        scores.put("Alice", 95);
        scores.put("Bob", 87);
        scores.put("Charlie", 92);
        
        System.out.println("Alice's score: " + scores.get("Alice"));
        System.out.println("HashMap: " + scores);
        
        // Update a value
        scores.put("Alice", 98);  // Updates existing key
        System.out.println("Updated Alice's score: " + scores.get("Alice"));
        
        // Check if key exists
        if (scores.containsKey("Bob")) {
            System.out.println("Bob is in the map");
        }
        
        // TreeMap - sorted by keys
        Map<String, Integer> sortedScores = new TreeMap<>();
        sortedScores.put("Charlie", 92);
        sortedScores.put("Alice", 95);
        sortedScores.put("Bob", 87);
        System.out.println("TreeMap (sorted): " + sortedScores);
        // Output: {Alice=95, Bob=87, Charlie=92}
        
        // LinkedHashMap - insertion order
        Map<String, Integer> orderedScores = new LinkedHashMap<>();
        orderedScores.put("Charlie", 92);
        orderedScores.put("Alice", 95);
        orderedScores.put("Bob", 87);
        System.out.println("LinkedHashMap (insertion order): " + orderedScores);
        // Output: {Charlie=92, Alice=95, Bob=87}
        
        // Iterating through a map
        System.out.println("\nIterating through entries:");
        for (Map.Entry<String, Integer> entry : scores.entrySet()) {
            System.out.println(entry.getKey() + ": " + entry.getValue());
        }
        
        // Working with keys and values
        System.out.println("All keys: " + scores.keySet());
        System.out.println("All values: " + scores.values());
        
        // Useful methods
        scores.putIfAbsent("Diana", 88);  // Only adds if key doesn't exist
        scores.remove("Bob");
        System.out.println("Final scores: " + scores);
    }
}
```

## 15.7 Collections Utility Class

The Collections class (note the 's' at the end) is a utility class packed with static methods for working with collections. It's like having a toolbox full of useful algorithms—sorting, searching, shuffling, reversing, and more. You never instantiate the Collections class; you just call its static methods.

Some of the most useful methods include `sort()` for sorting lists, `binarySearch()` for efficiently finding elements in sorted lists, `reverse()` for reversing the order of elements, `shuffle()` for randomizing the order, and `max()`/`min()` for finding the largest or smallest element.

There are also methods for creating thread-safe wrappers (`synchronizedList()`, `synchronizedSet()`, etc.) and unmodifiable collections (`unmodifiableList()`, `unmodifiableSet()`, etc.). The unmodifiable collections are great for ensuring that a collection can't be changed after creation—trying to add or remove elements will throw an exception.

One thing to remember: many of these methods modify the collection in place (like sort and reverse), while others return new collections or values (like min and max). Always check the documentation if you're unsure!

### Example: Collections Utility Methods

```java
import java.util.*;

public class CollectionsUtility {
    public static void main(String[] args) {
        // Sorting
        List<Integer> numbers = new ArrayList<>(Arrays.asList(5, 2, 8, 1, 9, 3));
        System.out.println("Original: " + numbers);
        
        Collections.sort(numbers);
        System.out.println("Sorted: " + numbers);
        
        Collections.reverse(numbers);
        System.out.println("Reversed: " + numbers);
        
        // Shuffling
        Collections.shuffle(numbers);
        System.out.println("Shuffled: " + numbers);
        
        // Finding min and max
        System.out.println("Min: " + Collections.min(numbers));
        System.out.println("Max: " + Collections.max(numbers));
        
        // Binary search (only works on sorted lists!)
        Collections.sort(numbers);
        int index = Collections.binarySearch(numbers, 5);
        System.out.println("Index of 5: " + index);
        
        // Fill - replace all elements with a value
        List<String> words = new ArrayList<>(Arrays.asList("a", "b", "c", "d"));
        Collections.fill(words, "X");
        System.out.println("After fill: " + words);
        
        // Frequency - count occurrences
        List<String> items = Arrays.asList("apple", "banana", "apple", "cherry", "apple");
        int appleCount = Collections.frequency(items, "apple");
        System.out.println("Apple count: " + appleCount);
        
        // Creating unmodifiable collections
        List<String> original = new ArrayList<>(Arrays.asList("A", "B", "C"));
        List<String> unmodifiable = Collections.unmodifiableList(original);
        System.out.println("Unmodifiable list: " + unmodifiable);
        
        // This would throw UnsupportedOperationException:
        // unmodifiable.add("D");
        
        // Rotating elements
        List<Integer> rotate = new ArrayList<>(Arrays.asList(1, 2, 3, 4, 5));
        Collections.rotate(rotate, 2);  // Rotate right by 2 positions
        System.out.println("Rotated: " + rotate);  // [4, 5, 1, 2, 3]
        
        // Sorting with custom comparator
        List<String> names = new ArrayList<>(Arrays.asList("Alice", "Bob", "Charlie", "Diana"));
        Collections.sort(names, (a, b) -> b.compareTo(a));  // Reverse alphabetical
        System.out.println("Reverse sorted: " + names);
    }
}
```

## 15.8 Iterators and Enhanced For-Loop

Iterating through collections is something you'll do constantly, and Java gives you several ways to do it. The enhanced for-loop (also called for-each loop) is the simplest and most readable for most cases. But sometimes you need more control, which is where Iterators come in.

The **enhanced for-loop** is clean and simple: `for (Type item : collection)`. It works with any collection and is perfect for when you just need to read each element. However, you can't modify the collection while using it (like removing elements), and you don't have access to the index.

An **Iterator** is an object that lets you traverse through a collection one element at a time. You get an iterator by calling `iterator()` on any collection. The Iterator has three main methods: `hasNext()` (returns true if there are more elements), `next()` (returns the next element), and `remove()` (safely removes the last element returned by next). The big advantage of iterators is that you can safely remove elements while iterating.

There's also a **ListIterator**, which is specific to Lists and gives you bidirectional traversal (you can go forward and backward) plus the ability to modify elements and add new ones while iterating.

For most everyday tasks, stick with the enhanced for-loop—it's cleaner and less error-prone. Use iterators when you need to remove elements during iteration or when you need more fine-grained control.

### Example: Different Ways to Iterate

```java
import java.util.*;

public class IterationExamples {
    public static void main(String[] args) {
        List<String> fruits = new ArrayList<>(Arrays.asList("Apple", "Banana", "Cherry", "Date"));
        
        // 1. Enhanced for-loop (for-each) - simplest and cleanest
        System.out.println("Enhanced for-loop:");
        for (String fruit : fruits) {
            System.out.println(fruit);
        }
        
        // 2. Iterator - allows safe removal
        System.out.println("\nUsing Iterator:");
        Iterator<String> iterator = fruits.iterator();
        while (iterator.hasNext()) {
            String fruit = iterator.next();
            System.out.println(fruit);
            
            // Safe to remove while iterating
            if (fruit.equals("Banana")) {
                iterator.remove();
            }
        }
        System.out.println("After removing Banana: " + fruits);
        
        // 3. Traditional for-loop with index (for Lists)
        System.out.println("\nTraditional for-loop:");
        for (int i = 0; i < fruits.size(); i++) {
            System.out.println(i + ": " + fruits.get(i));
        }
        
        // 4. ListIterator - bidirectional traversal
        System.out.println("\nListIterator (forward and backward):");
        ListIterator<String> listIterator = fruits.listIterator();
        
        // Forward
        System.out.print("Forward: ");
        while (listIterator.hasNext()) {
            System.out.print(listIterator.next() + " ");
        }
        System.out.println();
        
        // Backward
        System.out.print("Backward: ");
        while (listIterator.hasPrevious()) {
            System.out.print(listIterator.previous() + " ");
        }
        System.out.println();
        
        // 5. forEach method (Java 8+)
        System.out.println("\nforEach with lambda:");
        fruits.forEach(fruit -> System.out.println(fruit));
        
        // 6. Stream API (Java 8+)
        System.out.println("\nStream API:");
        fruits.stream()
              .filter(f -> f.startsWith("C"))
              .forEach(System.out::println);
        
        // DON'T DO THIS - ConcurrentModificationException!
        // for (String fruit : fruits) {
        //     if (fruit.equals("Cherry")) {
        //         fruits.remove(fruit);  // This will crash!
        //     }
        // }
    }
}
```

## 15.9 Comparable and Comparator Interfaces

Sorting is a common operation, but how does Java know how to sort your custom objects? That's where Comparable and Comparator come in. They're two different approaches to defining how objects should be ordered.

**Comparable** is an interface with a single method: `compareTo()`. When your class implements Comparable, you're defining its "natural ordering"—the default way it should be sorted. For example, String implements Comparable and sorts alphabetically. You'd implement Comparable when there's one obvious, standard way your objects should be ordered.

**Comparator** is a separate interface used to define alternative orderings. Instead of being part of the class itself, a Comparator is a separate object that knows how to compare two objects. This is great for when you want multiple different ways to sort the same class, or when you can't modify the class to implement Comparable (like classes from libraries).

Here's a rule of thumb: use Comparable for the default, natural ordering (like sorting people by their ID number). Use Comparator for alternative orderings (like sorting those same people by name, or by age, or by salary). You can have one Comparable implementation but as many Comparators as you need!

The compareTo/compare methods should return a negative number if the first object is "less than" the second, zero if they're equal, and a positive number if the first is "greater than" the second.

### Example: Comparable and Comparator

```java
import java.util.*;

// Student class implementing Comparable (natural ordering by ID)
class Student implements Comparable<Student> {
    private int id;
    private String name;
    private double gpa;
    
    public Student(int id, String name, double gpa) {
        this.id = id;
        this.name = name;
        this.gpa = gpa;
    }
    
    // Natural ordering: by ID
    @Override
    public int compareTo(Student other) {
        return this.id - other.id;
    }
    
    // Getters
    public int getId() { return id; }
    public String getName() { return name; }
    public double getGpa() { return gpa; }
    
    @Override
    public String toString() {
        return String.format("Student[id=%d, name=%s, gpa=%.2f]", id, name, gpa);
    }
}

public class ComparableComparatorExample {
    public static void main(String[] args) {
        List<Student> students = new ArrayList<>();
        students.add(new Student(103, "Alice", 3.8));
        students.add(new Student(101, "Bob", 3.5));
        students.add(new Student(102, "Charlie", 3.9));
        
        System.out.println("Original order:");
        students.forEach(System.out::println);
        
        // Sort using Comparable (natural ordering by ID)
        Collections.sort(students);
        System.out.println("\nSorted by ID (natural ordering):");
        students.forEach(System.out::println);
        
        // Sort using Comparator - by name
        Comparator<Student> nameComparator = new Comparator<Student>() {
            @Override
            public int compare(Student s1, Student s2) {
                return s1.getName().compareTo(s2.getName());
            }
        };
        
        Collections.sort(students, nameComparator);
        System.out.println("\nSorted by name:");
        students.forEach(System.out::println);
        
        // Sort using lambda expression - by GPA (descending)
        Collections.sort(students, (s1, s2) -> Double.compare(s2.getGpa(), s1.getGpa()));
        System.out.println("\nSorted by GPA (highest first):");
        students.forEach(System.out::println);
        
        // Using Comparator.comparing (Java 8+) - most elegant!
        students.sort(Comparator.comparing(Student::getName));
        System.out.println("\nSorted by name (using Comparator.comparing):");
        students.forEach(System.out::println);
        
        // Chaining comparators - sort by GPA, then by name
        students.sort(
            Comparator.comparing(Student::getGpa).reversed()
                      .thenComparing(Student::getName)
        );
        System.out.println("\nSorted by GPA (desc), then name:");
        students.forEach(System.out::println);
        
        // Sorting basic types
        List<String> names = Arrays.asList("Charlie", "Alice", "Bob");
        Collections.sort(names);  // Uses natural ordering
        System.out.println("\nSorted names: " + names);
        
        // Reverse ordering
        Collections.sort(names, Collections.reverseOrder());
        System.out.println("Reverse sorted names: " + names);
    }
}
```

## Common Pitfalls

**1. ConcurrentModificationException**: Don't modify a collection while iterating through it with a for-each loop. Use an Iterator and call its `remove()` method instead.

**2. Using HashMap with mutable keys**: If you modify an object after using it as a HashMap key, the hash code changes and you won't be able to find the value anymore. Use immutable objects as keys!

**3. Assuming HashSet/HashMap maintain order**: They don't! If you need order, use LinkedHashSet/LinkedHashMap or TreeSet/TreeMap.

**4. Forgetting that TreeSet/TreeMap require comparable elements**: If your objects don't implement Comparable and you don't provide a Comparator, you'll get a ClassCastException at runtime.

**5. Using == instead of equals() for containment checks**: Collections use `equals()` to check if an element exists, not `==`. Make sure your objects have a proper `equals()` implementation.

## Best Practices

**1. Program to interfaces**: Declare your variables as List, Set, or Map, not as ArrayList, HashSet, or HashMap. This gives you flexibility to change implementations later:
```java
List<String> list = new ArrayList<>();  // Good
ArrayList<String> list = new ArrayList<>();  // Not as flexible
```

**2. Choose the right collection**: Use ArrayList for most list needs, HashMap for most map needs, and HashSet for most set needs. Only switch to alternatives when you have a specific reason (need ordering, frequent insertions at the beginning, etc.).

**3. Initialize with capacity if you know the size**: Both ArrayList and HashMap can be initialized with an expected capacity to avoid resizing:
```java
List<String> bigList = new ArrayList<>(10000);
```

**4. Use immutable collections when appropriate**: If a collection shouldn't change, use `Collections.unmodifiableList()` or (in Java 9+) `List.of()`, `Set.of()`, `Map.of()`.

**5. Override equals() and hashCode() together**: If your class will be used in HashSet or as a HashMap key, always override both methods consistently.

## Practice Exercises

1. Create a program that reads a text file and counts the frequency of each word using a HashMap. Print the words sorted by frequency.

2. Implement a simple to-do list application using an ArrayList. Allow users to add tasks, remove tasks, mark tasks as complete, and view all tasks.

3. Use a TreeSet to store a collection of custom Person objects sorted by age. Demonstrate adding people and iterating through them in sorted order.

4. Create a PriorityQueue of tasks with different priority levels. Process the tasks in priority order.

5. Implement a simple cache using LinkedHashMap with access-order mode. When the cache reaches a maximum size, remove the least recently used item.

6. Write a method that takes two Lists and returns a new List containing only the elements that appear in both lists (intersection).

---

## Quiz Questions

This section contains quiz questions that test understanding of the Java Collections Framework. Each question is formatted for easy extraction and database import.

### Question 1: Collections Framework - Basic Concepts
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: java, collections, collections-framework, list, set

**Question Text**:
What is the main difference between a List and a Set in the Java Collections Framework?

**Answer Choices**:
- A) Lists are faster than Sets for all operations
- B) Lists allow duplicate elements while Sets do not
- C) Sets are ordered while Lists are not
- D) Lists can only store primitive types

**Correct Answer**: B

**Explanation**:
The fundamental difference is that Lists allow duplicate elements while Sets do not. When you try to add a duplicate element to a Set, it simply won't be added. Lists, on the other hand, happily accept duplicates. Both can store objects (not primitives directly, though autoboxing handles this). Lists maintain insertion order, while HashSet doesn't guarantee order (though TreeSet and LinkedHashSet have specific ordering). Performance depends on the specific implementation and operation.

---

### Question 2: ArrayList vs LinkedList - Performance
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, collections, list, arraylist, linkedlist, performance

**Question Text**:
Which operation is generally FASTER on a LinkedList compared to an ArrayList?

**Answer Choices**:
- A) Accessing an element by index (e.g., get(50))
- B) Adding an element at the end of the list
- C) Inserting an element at the beginning of the list
- D) Searching for an element by value

**Correct Answer**: C

**Explanation**:
Inserting at the beginning of a LinkedList is O(1) because it just requires updating a few references. For ArrayList, inserting at the beginning is O(n) because all existing elements must be shifted one position to the right. Accessing by index is faster on ArrayList (O(1) vs O(n) for LinkedList). Adding at the end is roughly comparable for both (amortized O(1)). Searching by value is O(n) for both, but ArrayList is typically faster in practice due to better cache locality.

---

### Question 3: HashSet - Duplicate Detection
**Type**: code-output
**Difficulty**: 3
**Tags**: java, collections, set, hashset, duplicates

**Question Text**:
What will be the output of this code?

```java
Set<String> set = new HashSet<>();
set.add("Java");
set.add("Python");
set.add("Java");
System.out.println(set.size());
```

**Answer Choices**:
- A) 1
- B) 2
- C) 3
- D) Compilation error

**Correct Answer**: B

**Explanation**:
A HashSet does not allow duplicates. When "Java" is added the first time, it's successfully added. When "Java" is added the second time, the HashSet detects it's a duplicate (using the equals() method) and doesn't add it again. The add() method returns false for the second attempt, but no exception is thrown. Therefore, the set contains only two elements: "Java" and "Python", so the size is 2.

---

### Question 4: TreeSet - Automatic Sorting
**Type**: code-output
**Difficulty**: 3
**Tags**: java, collections, set, treeset, sorting

**Question Text**:
What will be printed by this code?

```java
Set<Integer> numbers = new TreeSet<>();
numbers.add(5);
numbers.add(2);
numbers.add(8);
numbers.add(1);
System.out.println(numbers);
```

**Answer Choices**:
- A) [5, 2, 8, 1]
- B) [1, 2, 5, 8]
- C) [8, 5, 2, 1]
- D) The order is unpredictable

**Correct Answer**: B

**Explanation**:
TreeSet automatically maintains elements in sorted order using their natural ordering (or a provided Comparator). For integers, natural ordering means ascending numerical order. Regardless of the insertion order (5, 2, 8, 1), the TreeSet will store and display them as [1, 2, 5, 8]. This is one of the key features of TreeSet—you get automatic sorting at the cost of slightly slower performance (O(log n) operations instead of O(1) for HashSet).

---

### Question 5: HashMap - Key-Value Operations
**Type**: code-output
**Difficulty**: 4
**Tags**: java, collections, map, hashmap, key-value

**Question Text**:
What will be the output of this code?

```java
Map<String, Integer> map = new HashMap<>();
map.put("A", 10);
map.put("B", 20);
map.put("A", 30);
System.out.println(map.get("A"));
```

**Answer Choices**:
- A) 10
- B) 20
- C) 30
- D) null

**Correct Answer**: C

**Explanation**:
In a HashMap, keys must be unique. When you put a value with a key that already exists, it replaces the old value. The first put("A", 10) adds the key "A" with value 10. The third put("A", 30) updates the value for key "A" from 10 to 30. Therefore, map.get("A") returns 30. The map contains two entries: {"A"=30, "B"=20}. The old value (10) is lost.

---

### Question 6: PriorityQueue - Ordering
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, collections, queue, priorityqueue, ordering

**Question Text**:
In a PriorityQueue of integers (without a custom Comparator), which element will be at the head after these operations?

```java
PriorityQueue<Integer> pq = new PriorityQueue<>();
pq.offer(5);
pq.offer(2);
pq.offer(8);
pq.offer(1);
```

**Answer Choices**:
- A) 5 (the first element added)
- B) 1 (the smallest element)
- C) 8 (the largest element)
- D) 2 (the second element added)

**Correct Answer**: B

**Explanation**:
A PriorityQueue orders elements based on their natural ordering (or a provided Comparator). For integers, natural ordering is ascending (smallest to largest). The PriorityQueue is a min-heap by default, meaning the smallest element is always at the head. After adding 5, 2, 8, and 1, the element 1 (the smallest) will be at the head. When you call poll(), it returns 1, then 2, then 5, then 8. Note that this doesn't mean the PriorityQueue is fully sorted—only the head is guaranteed to be the minimum.

---

### Question 7: Collections.sort() - Behavior
**Type**: true-false
**Difficulty**: 3
**Tags**: java, collections, collections-utility, sorting

**Question Text**:
True or False: The Collections.sort() method returns a new sorted list and leaves the original list unchanged.

**Answer Choices**:
- A) True
- B) False

**Correct Answer**: B

**Explanation**:
False. Collections.sort() modifies the list in-place and returns void (nothing). It sorts the original list directly, so the list you pass in will be changed. If you want to keep the original list unchanged, you should create a copy first, then sort the copy. For example: List<Integer> sorted = new ArrayList<>(original); Collections.sort(sorted);. Alternatively, you can use streams: original.stream().sorted().collect(Collectors.toList()) which returns a new sorted list.

---

### Question 8: Iterator - Safe Removal
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, collections, iterator, iteration

**Question Text**:
What is the correct way to remove elements from a list while iterating through it?

**Answer Choices**:
- A) Use a for-each loop and call list.remove(element)
- B) Use a traditional for loop and call list.remove(i)
- C) Use an Iterator and call iterator.remove()
- D) Both A and B are correct

**Correct Answer**: C

**Explanation**:
The correct and safe way is to use an Iterator and call iterator.remove(). Using a for-each loop with list.remove() will throw a ConcurrentModificationException because you're modifying the collection while iterating. Using a traditional for loop with list.remove(i) can work but is error-prone—you need to decrement the index after removal, and it only works with Lists (not all Collections). Iterator.remove() is designed specifically for safe removal during iteration and works with any Collection.

---

### Question 9: Comparable vs Comparator - Usage
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, collections, comparable, comparator, sorting

**Question Text**:
When should you use a Comparator instead of implementing Comparable?

**Answer Choices**:
- A) When you need only one way to sort objects
- B) When you want multiple different sorting orders for the same class
- C) When sorting primitive types
- D) Comparator is always preferred over Comparable

**Correct Answer**: B

**Explanation**:
Use Comparator when you need multiple different ways to sort the same class. For example, you might want to sort Person objects by name, by age, or by salary—you'd use different Comparators for each. Comparable is for defining the "natural" or default ordering (implement it once in the class). Primitive types are already Comparable. While you can use Comparators for single sorting criteria, it's more elegant to implement Comparable for the natural ordering and use Comparators for alternative orderings. Both have their place depending on the use case.

---

### Question 10: LinkedHashMap - Ordering
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, collections, map, linkedhashmap, ordering

**Question Text**:
What is the primary advantage of LinkedHashMap over HashMap?

**Answer Choices**:
- A) LinkedHashMap is significantly faster
- B) LinkedHashMap can store more elements
- C) LinkedHashMap maintains insertion order (or access order)
- D) LinkedHashMap allows duplicate keys

**Correct Answer**: C

**Explanation**:
The primary advantage of LinkedHashMap is that it maintains the insertion order of entries (or optionally, access order if constructed with the appropriate flag). When you iterate through a LinkedHashMap, you get entries in the order they were inserted, which is predictable and useful. HashMap provides no ordering guarantees. The performance difference is minimal—LinkedHashMap is slightly slower and uses slightly more memory due to maintaining a linked list alongside the hash table. Both only allow unique keys (like all Map implementations).

---

### Question 11: Collections Utility - Binary Search
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, collections, collections-utility, binary-search, sorting

**Question Text**:
What is a requirement for using Collections.binarySearch() on a List?

**Answer Choices**:
- A) The list must be an ArrayList
- B) The list must be sorted in ascending order
- C) The list size must be a power of 2
- D) The list elements must be numbers

**Correct Answer**: B

**Explanation**:
Collections.binarySearch() requires that the list be sorted in ascending order (according to natural ordering or the provided Comparator). If the list is not sorted, the results are undefined—you might get an incorrect index or fail to find an element that's actually in the list. The list can be any implementation (ArrayList, LinkedList, etc.), though ArrayList is more efficient for binary search due to O(1) index access. The size can be any value, and elements can be any type as long as they're comparable and sorted. Always sort before using binary search!

---

### Question 12: Collection Modification - Unmodifiable Collections
**Type**: code-output
**Difficulty**: 3
**Tags**: java, collections, collections-utility, unmodifiable

**Question Text**:
What happens when you run this code?

```java
List<String> list = new ArrayList<>(Arrays.asList("A", "B", "C"));
List<String> unmod = Collections.unmodifiableList(list);
unmod.add("D");
```

**Answer Choices**:
- A) The list successfully adds "D"
- B) Nothing happens, the add is silently ignored
- C) UnsupportedOperationException is thrown
- D) Compilation error

**Correct Answer**: C

**Explanation**:
Collections.unmodifiableList() returns a wrapper around the original list that throws UnsupportedOperationException for any modification attempts (add, remove, set, clear, etc.). This is a runtime exception that will crash your program if not caught. The purpose of unmodifiable collections is to prevent accidental modifications and to safely share collections with other code without worrying about them being changed. Note that changes to the original list will be reflected in the unmodifiable view—it's a view, not a copy. For truly immutable collections in Java 9+, use List.of(), Set.of(), or Map.of().
