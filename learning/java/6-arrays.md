# Java Arrays

Welcome to arrays! An array is a container that holds a fixed number of values of the same type. Think of it like a row of numbered boxes, where each box holds one value. Arrays are fundamental to programming—you'll use them constantly to store and organize collections of data.

Whether you're storing test scores, usernames, or sensor readings, arrays give you a powerful way to work with multiple related values at once.

## One-Dimensional Arrays

A one-dimensional array is a simple list of elements, all of the same type. Let's explore how to create and use them.

### Array Declaration

To declare an array, add square brackets `[]` to the type:

```java
int[] numbers;        // Preferred style
String[] names;
double[] prices;

// Alternative style (less common)
int numbers[];        // Also valid, but not preferred
```

**Pro tip**: Use `int[] numbers` rather than `int numbers[]` because it makes the type clearer—you're declaring "an array of ints called numbers."

### Array Initialization

After declaring an array, you need to create it with a specific size:

```java
// Create an array of 5 integers
int[] numbers = new int[5];

// Create an array of 3 strings
String[] names = new String[3];
```

When you create an array, Java initializes it with default values:
- Numbers (int, double, etc.): `0` or `0.0`
- Booleans: `false`
- Objects (String, etc.): `null`

### Array Initialization with Values

You can create and populate an array in one step:

```java
// Method 1: Using array initializer
int[] scores = {95, 87, 92, 78, 88};

// Method 2: Creating then assigning
int[] ages = new int[] {25, 30, 22, 35};

// Method 3: Specify size, then assign values
String[] fruits = new String[3];
fruits[0] = "Apple";
fruits[1] = "Banana";
fruits[2] = "Cherry";
```

The first method is the most concise and commonly used.

### Accessing Array Elements

Array indices start at 0. The first element is at index 0, the second at index 1, and so on:

```java
int[] scores = {95, 87, 92, 78, 88};

System.out.println(scores[0]);  // 95 (first element)
System.out.println(scores[1]);  // 87 (second element)
System.out.println(scores[4]);  // 88 (last element)

// Modify an element
scores[2] = 100;
System.out.println(scores[2]);  // 100
```

**Important**: Trying to access an index that doesn't exist causes an `ArrayIndexOutOfBoundsException`:

```java
int[] nums = {1, 2, 3};
System.out.println(nums[5]);  // ERROR! Array only has indices 0, 1, 2
```

### Array Length

Every array has a `length` property that tells you how many elements it contains:

```java
int[] numbers = {10, 20, 30, 40, 50};
System.out.println(numbers.length);  // 5

// Use length to access the last element
int lastElement = numbers[numbers.length - 1];  // 50
```

**Note**: It's `length` (not `length()`). It's a property, not a method!

### Iterating Through Arrays

Use loops to process all elements in an array:

```java
int[] scores = {95, 87, 92, 78, 88};

// Method 1: Traditional for loop
for (int i = 0; i < scores.length; i++) {
    System.out.println("Score " + i + ": " + scores[i]);
}

// Method 2: Enhanced for loop (for-each)
for (int score : scores) {
    System.out.println("Score: " + score);
}
```

The enhanced for loop is cleaner when you don't need the index. Use the traditional for loop when you need to know the position or modify elements.

### Practical Example: Calculate Average

```java
int[] testScores = {85, 92, 78, 90, 88};

int sum = 0;
for (int score : testScores) {
    sum += score;
}

double average = (double) sum / testScores.length;
System.out.println("Average score: " + average);  // 86.6
```

## Multi-Dimensional Arrays

Multi-dimensional arrays are arrays of arrays. The most common is a two-dimensional array, which you can visualize as a table with rows and columns.

### Two-Dimensional Arrays

Think of a 2D array as a grid:

```java
// Declare and create a 3x4 array (3 rows, 4 columns)
int[][] matrix = new int[3][4];

// Or initialize with values
int[][] grid = {
    {1, 2, 3, 4},
    {5, 6, 7, 8},
    {9, 10, 11, 12}
};
```

Visualized:
```
    Col0  Col1  Col2  Col3
Row0:  1     2     3     4
Row1:  5     6     7     8
Row2:  9    10    11    12
```

### Accessing 2D Array Elements

Use two indices: `[row][column]`:

```java
int[][] grid = {
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9}
};

System.out.println(grid[0][0]);  // 1 (first row, first column)
System.out.println(grid[1][2]);  // 6 (second row, third column)
System.out.println(grid[2][1]);  // 8 (third row, second column)

// Modify an element
grid[1][1] = 50;
System.out.println(grid[1][1]);  // 50
```

### Iterating Through 2D Arrays

Use nested loops:

```java
int[][] matrix = {
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9}
};

// Method 1: Traditional nested for loops
for (int row = 0; row < matrix.length; row++) {
    for (int col = 0; col < matrix[row].length; col++) {
        System.out.print(matrix[row][col] + " ");
    }
    System.out.println();  // New line after each row
}

// Method 2: Enhanced for loops
for (int[] row : matrix) {
    for (int value : row) {
        System.out.print(value + " ");
    }
    System.out.println();
}
```

### 2D Array Dimensions

```java
int[][] table = new int[3][4];  // 3 rows, 4 columns

System.out.println(table.length);      // 3 (number of rows)
System.out.println(table[0].length);   // 4 (number of columns in first row)
```

### Jagged Arrays

Jagged arrays are 2D arrays where each row can have a different length:

```java
int[][] jagged = {
    {1, 2},
    {3, 4, 5, 6},
    {7}
};

System.out.println(jagged[0].length);  // 2
System.out.println(jagged[1].length);  // 4
System.out.println(jagged[2].length);  // 1
```

Visualized:
```
Row 0: [1, 2]
Row 1: [3, 4, 5, 6]
Row 2: [7]
```

This is useful when different rows naturally have different numbers of elements, like storing family members for different families.

### Practical Example: Grade Table

```java
// Student grades: each row is a student, each column is a test
int[][] grades = {
    {85, 90, 78},  // Student 0
    {92, 88, 95},  // Student 1
    {76, 81, 79}   // Student 2
};

// Calculate each student's average
for (int i = 0; i < grades.length; i++) {
    int sum = 0;
    for (int j = 0; j < grades[i].length; j++) {
        sum += grades[i][j];
    }
    double average = (double) sum / grades[i].length;
    System.out.println("Student " + i + " average: " + average);
}
```

## Arrays Class and Utility Methods

The `java.util.Arrays` class provides useful methods for working with arrays. You must import it first:

```java
import java.util.Arrays;
```

### Sorting Arrays

Sort arrays in ascending order:

```java
int[] numbers = {5, 2, 8, 1, 9};
Arrays.sort(numbers);
System.out.println(Arrays.toString(numbers));  // [1, 2, 5, 8, 9]

String[] names = {"Charlie", "Alice", "Bob"};
Arrays.sort(names);
System.out.println(Arrays.toString(names));  // [Alice, Bob, Charlie]
```

You can also sort a portion of an array:

```java
int[] nums = {5, 2, 8, 1, 9};
Arrays.sort(nums, 1, 4);  // Sort indices 1-3
System.out.println(Arrays.toString(nums));  // [5, 1, 2, 8, 9]
```

### Searching Arrays

Use `binarySearch` to find elements (array must be sorted first!):

```java
int[] numbers = {1, 3, 5, 7, 9, 11};
int index = Arrays.binarySearch(numbers, 7);
System.out.println("Found at index: " + index);  // 3

// Search for value that doesn't exist
int notFound = Arrays.binarySearch(numbers, 6);
System.out.println(notFound);  // Negative number (insertion point)
```

**Important**: `binarySearch` only works correctly on sorted arrays! If the array isn't sorted, results are unpredictable.

### Copying Arrays

Create copies of arrays:

```java
int[] original = {1, 2, 3, 4, 5};

// Method 1: Arrays.copyOf - copy entire array
int[] copy1 = Arrays.copyOf(original, original.length);

// Method 2: Arrays.copyOf - copy with different length
int[] copy2 = Arrays.copyOf(original, 3);  // {1, 2, 3}
int[] copy3 = Arrays.copyOf(original, 7);  // {1, 2, 3, 4, 5, 0, 0}

// Method 3: Arrays.copyOfRange - copy a range
int[] copy4 = Arrays.copyOfRange(original, 1, 4);  // {2, 3, 4}
```

**Why copy?** Assignment doesn't copy arrays—it just makes another reference:

```java
int[] arr1 = {1, 2, 3};
int[] arr2 = arr1;  // arr2 points to the SAME array
arr2[0] = 99;
System.out.println(arr1[0]);  // 99! Both point to same array

// To actually copy:
int[] arr3 = Arrays.copyOf(arr1, arr1.length);
arr3[0] = 5;
System.out.println(arr1[0]);  // 99 (unchanged)
System.out.println(arr3[0]);  // 5
```

### Comparing Arrays

Check if two arrays are equal:

```java
int[] arr1 = {1, 2, 3};
int[] arr2 = {1, 2, 3};
int[] arr3 = {1, 2, 4};

System.out.println(Arrays.equals(arr1, arr2));  // true
System.out.println(Arrays.equals(arr1, arr3));  // false

// DON'T use == for arrays!
System.out.println(arr1 == arr2);  // false (compares references, not content)
```

### Filling Arrays

Fill an array with a specific value:

```java
int[] numbers = new int[5];
Arrays.fill(numbers, 7);
System.out.println(Arrays.toString(numbers));  // [7, 7, 7, 7, 7]

// Fill a range
int[] nums = {1, 2, 3, 4, 5};
Arrays.fill(nums, 1, 4, 0);  // Fill indices 1-3 with 0
System.out.println(Arrays.toString(nums));  // [1, 0, 0, 0, 5]
```

### Converting to String

Print arrays nicely:

```java
int[] numbers = {1, 2, 3, 4, 5};

// Without Arrays.toString()
System.out.println(numbers);  // Prints gibberish like [I@15db9742

// With Arrays.toString()
System.out.println(Arrays.toString(numbers));  // [1, 2, 3, 4, 5]

// For 2D arrays, use deepToString()
int[][] matrix = {{1, 2}, {3, 4}};
System.out.println(Arrays.deepToString(matrix));  // [[1, 2], [3, 4]]
```

## Common Array Operations and Algorithms

Let's explore common operations you'll perform on arrays.

### Finding Maximum and Minimum

```java
int[] numbers = {45, 23, 67, 12, 89, 34};

// Find maximum
int max = numbers[0];
for (int num : numbers) {
    if (num > max) {
        max = num;
    }
}
System.out.println("Max: " + max);  // 89

// Find minimum
int min = numbers[0];
for (int num : numbers) {
    if (num < min) {
        min = num;
    }
}
System.out.println("Min: " + min);  // 12
```

### Reversing an Array

```java
int[] arr = {1, 2, 3, 4, 5};

// Reverse in place
for (int i = 0; i < arr.length / 2; i++) {
    int temp = arr[i];
    arr[i] = arr[arr.length - 1 - i];
    arr[arr.length - 1 - i] = temp;
}

System.out.println(Arrays.toString(arr));  // [5, 4, 3, 2, 1]
```

### Searching for an Element

```java
int[] numbers = {10, 20, 30, 40, 50};
int target = 30;

// Linear search
int index = -1;
for (int i = 0; i < numbers.length; i++) {
    if (numbers[i] == target) {
        index = i;
        break;
    }
}

if (index != -1) {
    System.out.println("Found at index: " + index);  // 2
} else {
    System.out.println("Not found");
}
```

### Shifting Elements

```java
int[] arr = {1, 2, 3, 4, 5};

// Shift left (remove first element)
for (int i = 0; i < arr.length - 1; i++) {
    arr[i] = arr[i + 1];
}
arr[arr.length - 1] = 0;  // Clear last element
System.out.println(Arrays.toString(arr));  // [2, 3, 4, 5, 0]

// Shift right (make room at beginning)
int[] arr2 = {1, 2, 3, 4, 5};
for (int i = arr2.length - 1; i > 0; i--) {
    arr2[i] = arr2[i - 1];
}
arr2[0] = 0;  // Clear first element
System.out.println(Arrays.toString(arr2));  // [0, 1, 2, 3, 4]
```

### Removing Duplicates

```java
int[] arr = {1, 2, 2, 3, 4, 4, 5};
Arrays.sort(arr);  // Ensure sorted

// Count unique elements
int uniqueCount = 1;
for (int i = 1; i < arr.length; i++) {
    if (arr[i] != arr[i - 1]) {
        uniqueCount++;
    }
}

// Create new array with unique elements
int[] unique = new int[uniqueCount];
unique[0] = arr[0];
int j = 1;
for (int i = 1; i < arr.length; i++) {
    if (arr[i] != arr[i - 1]) {
        unique[j++] = arr[i];
    }
}

System.out.println(Arrays.toString(unique));  // [1, 2, 3, 4, 5]
```

### Merging Two Arrays

```java
int[] arr1 = {1, 3, 5};
int[] arr2 = {2, 4, 6};

int[] merged = new int[arr1.length + arr2.length];

// Copy first array
for (int i = 0; i < arr1.length; i++) {
    merged[i] = arr1[i];
}

// Copy second array
for (int i = 0; i < arr2.length; i++) {
    merged[arr1.length + i] = arr2[i];
}

System.out.println(Arrays.toString(merged));  // [1, 3, 5, 2, 4, 6]

// Sort the merged array
Arrays.sort(merged);
System.out.println(Arrays.toString(merged));  // [1, 2, 3, 4, 5, 6]
```

### Counting Occurrences

```java
int[] numbers = {1, 2, 3, 2, 4, 2, 5};
int target = 2;

int count = 0;
for (int num : numbers) {
    if (num == target) {
        count++;
    }
}

System.out.println("Number " + target + " appears " + count + " times");  // 3 times
```

## Common Pitfalls and Best Practices

### Pitfall 1: Index Out of Bounds

```java
int[] arr = {1, 2, 3};
System.out.println(arr[3]);  // ERROR! Valid indices are 0, 1, 2
```

**Fix**: Always check `i < arr.length`, not `i <= arr.length`

### Pitfall 2: Null Arrays

```java
int[] numbers = null;
System.out.println(numbers.length);  // NullPointerException!
```

**Fix**: Always initialize arrays before using them

### Pitfall 3: Array Assignment Creates Reference

```java
int[] arr1 = {1, 2, 3};
int[] arr2 = arr1;  // Both point to same array!
arr2[0] = 99;
System.out.println(arr1[0]);  // 99 (modified!)
```

**Fix**: Use `Arrays.copyOf()` to create a real copy

### Pitfall 4: Comparing with ==

```java
int[] arr1 = {1, 2, 3};
int[] arr2 = {1, 2, 3};
System.out.println(arr1 == arr2);  // false! (compares references)
```

**Fix**: Use `Arrays.equals()` to compare content

### Best Practices

1. **Use `Arrays.toString()`** for debugging and printing
2. **Check bounds** before accessing array elements
3. **Use enhanced for loops** when you don't need indices
4. **Initialize arrays** with meaningful size or values
5. **Sort before binary search** - `binarySearch` requires sorted arrays
6. **Use `Arrays.copyOf()`** for copying, not assignment
7. **Use `Arrays.equals()`** for comparing, not `==`

## Real-World Example: Student Grade Manager

```java
import java.util.Arrays;

public class GradeManager {
    public static void main(String[] args) {
        // Student names and their test scores
        String[] students = {"Alice", "Bob", "Charlie", "Diana"};
        int[][] scores = {
            {85, 90, 88},  // Alice
            {92, 88, 95},  // Bob
            {78, 82, 80},  // Charlie
            {95, 98, 92}   // Diana
        };
        
        // Calculate and display averages
        double[] averages = new double[students.length];
        for (int i = 0; i < students.length; i++) {
            int sum = 0;
            for (int score : scores[i]) {
                sum += score;
            }
            averages[i] = (double) sum / scores[i].length;
            System.out.printf("%s: %.2f%n", students[i], averages[i]);
        }
        
        // Find top student
        int topIndex = 0;
        for (int i = 1; i < averages.length; i++) {
            if (averages[i] > averages[topIndex]) {
                topIndex = i;
            }
        }
        System.out.println("\nTop student: " + students[topIndex]);
        
        // Sort students by average (selection sort for simplicity)
        for (int i = 0; i < averages.length - 1; i++) {
            int maxIdx = i;
            for (int j = i + 1; j < averages.length; j++) {
                if (averages[j] > averages[maxIdx]) {
                    maxIdx = j;
                }
            }
            // Swap students
            String tempName = students[i];
            students[i] = students[maxIdx];
            students[maxIdx] = tempName;
            // Swap averages
            double tempAvg = averages[i];
            averages[i] = averages[maxIdx];
            averages[maxIdx] = tempAvg;
        }
        
        System.out.println("\nRanked by average:");
        for (int i = 0; i < students.length; i++) {
            System.out.printf("%d. %s: %.2f%n", i + 1, students[i], averages[i]);
        }
    }
}
```

---

## Quiz Questions

### Question 1: Array Declaration and Initialization
**Type**: multiple-choice
**Difficulty**: 2
**Tags**: java, arrays, declaration, initialization

**Question Text**:
Which of the following correctly creates an array of 5 integers in Java?

**Answer Choices**:
- A) `int[] numbers = new int[5];`
- B) `int numbers = new int[5];`
- C) `int[] numbers = new [5];`
- D) `array<int> numbers = new array<int>(5);`

**Correct Answer**: A

**Explanation**:
The correct syntax is `int[] numbers = new int[5];` This declares an array variable of type int[], then creates a new array of size 5. Option B is missing the brackets, option C has incorrect syntax after 'new', and option D uses C++-style syntax that doesn't exist in Java. The created array will be initialized with default values (0 for integers).

---

### Question 2: Array Indexing
**Type**: code-output
**Difficulty**: 2
**Tags**: java, arrays, indexing, accessing

**Question Text**:
What does this code print?

```java
String[] fruits = {"Apple", "Banana", "Cherry"};
System.out.println(fruits[1]);
```

**Answer Choices**:
- A) Apple
- B) Banana
- C) Cherry
- D) Error

**Correct Answer**: B

**Explanation**:
Arrays in Java use zero-based indexing, meaning the first element is at index 0, the second at index 1, and so on. So fruits[0] is "Apple", fruits[1] is "Banana", and fruits[2] is "Cherry". Therefore, this code prints "Banana".

---

### Question 3: Array Length Property
**Type**: code-output
**Difficulty**: 3
**Tags**: java, arrays, length, property

**Question Text**:
What is the output of this code?

```java
int[] numbers = {10, 20, 30, 40, 50};
System.out.println(numbers[numbers.length - 1]);
```

**Answer Choices**:
- A) 10
- B) 40
- C) 50
- D) 5

**Correct Answer**: C

**Explanation**:
The `length` property returns the number of elements in the array, which is 5. So `numbers.length - 1` equals 4, and `numbers[4]` is the last element (50). This is a common pattern for accessing the last element of an array, since indices go from 0 to length-1.

---

### Question 4: Two-Dimensional Arrays
**Type**: code-output
**Difficulty**: 4
**Tags**: java, arrays, multidimensional, 2d-arrays

**Question Text**:
What does this code print?

```java
int[][] grid = {
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9}
};
System.out.println(grid[1][2]);
```

**Answer Choices**:
- A) 2
- B) 5
- C) 6
- D) 8

**Correct Answer**: C

**Explanation**:
In a 2D array, the first index represents the row and the second represents the column (both zero-based). `grid[1][2]` accesses row 1 (the second row: {4, 5, 6}) and column 2 (the third element: 6). So the answer is 6.

---

### Question 5: Arrays.sort() Method
**Type**: code-output
**Difficulty**: 3
**Tags**: java, arrays, sorting, arrays-class

**Question Text**:
What is the output after sorting?

```java
int[] numbers = {5, 2, 8, 1, 9};
Arrays.sort(numbers);
System.out.println(numbers[0] + " " + numbers[4]);
```

**Answer Choices**:
- A) 5 9
- B) 1 9
- C) 9 1
- D) 2 8

**Correct Answer**: B

**Explanation**:
`Arrays.sort()` sorts the array in ascending order, so the array becomes {1, 2, 5, 8, 9}. `numbers[0]` is the first element (1) and `numbers[4]` is the last element (9), so it prints "1 9". The sort modifies the original array rather than creating a new one.

---

### Question 6: Array Copying
**Type**: multiple-choice
**Difficulty**: 5
**Tags**: java, arrays, copying, references

**Question Text**:
What happens when you assign one array to another using `arr2 = arr1`?

**Answer Choices**:
- A) arr2 becomes an independent copy of arr1
- B) arr2 references the same array as arr1; changes to one affect the other
- C) The elements of arr1 are copied to arr2
- D) A compilation error occurs

**Correct Answer**: B

**Explanation**:
Array assignment in Java creates a reference, not a copy. Both variables point to the same array in memory. If you modify the array through either variable, the change is visible through both. To create an actual copy, you need to use `Arrays.copyOf()` or similar methods. This is a common source of bugs for beginners.

---

### Question 7: Jagged Arrays
**Type**: multiple-choice
**Difficulty**: 4
**Tags**: java, arrays, multidimensional, jagged-arrays

**Question Text**:
What is a jagged array in Java?

**Answer Choices**:
- A) An array with missing elements
- B) A 2D array where each row can have a different number of columns
- C) An array that contains only odd numbers
- D) An array that has been sorted

**Correct Answer**: B

**Explanation**:
A jagged array is a multidimensional array where each row can have a different length. For example: `int[][] jagged = {{1, 2}, {3, 4, 5, 6}, {7}};` has rows of lengths 2, 4, and 1 respectively. This is useful when different rows naturally have different numbers of elements, unlike rectangular 2D arrays where all rows have the same length.

---

### Question 8: Array Iteration and Sum
**Type**: code-output
**Difficulty**: 5
**Tags**: java, arrays, iteration, loops

**Question Text**:
What is the output of this code?

```java
int[] nums = {2, 4, 6, 8, 10};
int sum = 0;
for (int i = 0; i < nums.length; i += 2) {
    sum += nums[i];
}
System.out.println(sum);
```

**Answer Choices**:
- A) 30
- B) 18
- C) 20
- D) 12

**Correct Answer**: B

**Explanation**:
The loop increments by 2 (`i += 2`), so it only processes elements at even indices: i=0 (value 2), i=2 (value 6), and i=4 (value 10). The sum is 2 + 6 + 10 = 18. Elements at odd indices (4 and 8) are skipped. This demonstrates how loop control can be used to process only certain elements of an array.
