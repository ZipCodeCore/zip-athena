# Question Bulk Loader

This directory contains question files that can be bulk-loaded into the Athena database.

## Usage

### Web UI Import (Recommended)

1. Log in to the Admin Dashboard at `/admin.html`
2. Navigate to the **Questions** tab
3. Click the **Import Questions** button
4. Select a topic from the dropdown menu
5. Choose a JSON file to upload
6. Click **Import Questions**

The import panel will display:
- Number of questions successfully imported
- Number of duplicate questions skipped
- Any errors encountered during import

### Command-Line Import

```bash
# Load all JSON files from the questions directory
go run ./cmd/loadquestions/main.go

# Dry run to validate without importing
go run ./cmd/loadquestions/main.go --dry-run

# Load from a specific file
go run ./cmd/loadquestions/main.go --file questions/java-fundamentals.json

# Use a different database
go run ./cmd/loadquestions/main.go --db ./path/to/athena.db

# Use a different questions directory
go run ./cmd/loadquestions/main.go --dir ./path/to/questions
```

Like, etc, etc:

```bash
go run ./cmd/loadquestions/main.go --file questions/java-basics2.json
```

## JSON Format

Each question file should be a JSON file with the following structure:

```json
{
  "topic": "Topic Name",
  "questions": [
    {
      "text": "The question text?",
      "difficulty": 3,
      "tags": ["tag1", "tag2"],
      "answers": [
        {"text": "Correct answer", "is_correct": true},
        {"text": "Wrong answer 1", "is_correct": false},
        {"text": "Wrong answer 2", "is_correct": false},
        {"text": "Wrong answer 3", "is_correct": false}
      ]
    }
  ]
}
```

### Field Requirements

| Field | Required | Description |
|-------|----------|-------------|
| `topic` | Yes | Must match an existing topic name in the database (case-insensitive) |
| `text` | Yes | The question text |
| `difficulty` | Yes | Integer from 1 (easiest) to 7 (hardest) |
| `tags` | No | Array of tag strings for categorization |
| `answers` | Yes | 2-6 answer options |
| `is_correct` | Yes | Exactly one answer must be marked as correct |

### Available Topics

Run this to see available topics:
```bash
sqlite3 athena.db "SELECT name FROM topics;"
```

Default topics:
- Java Fundamentals
- Python Basics
- SQL Database

## Features

- **Duplicate Detection**: Questions with identical text are skipped to prevent duplicates
- **Validation**: Questions are validated before import (difficulty range, answer requirements)
- **Dry Run**: Use `--dry-run` to validate files without importing
- **Batch Processing**: All JSON files in the directory are processed

## Example Files

- `java-fundamentals.json` - 30 Java questions covering core concepts
- `python-basics.json` - 30 Python questions covering fundamentals
- `git-github-basics.json` - 25 beginner Git/GitHub questions
- `shell-commands-basics.json` - 25 beginner shell command questions
