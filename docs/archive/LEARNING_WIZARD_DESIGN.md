# Learning Wizard Design Document

## Overview

The Learning Wizard is an integrated feature within Athena that allows students to access and study learning materials directly from their dashboard. It provides a seamless learning experience by:

1. Displaying learning materials in a structured, easy-to-navigate interface
2. Linking quiz questions to relevant learning materials
3. Allowing students to review materials when they get quiz questions wrong

## Architecture

### 1. Database Schema Extensions

```sql
-- Learning Resources table
CREATE TABLE learning_resources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    language TEXT NOT NULL,              -- 'java' or 'python'
    phase INTEGER NOT NULL,               -- 1-5
    section_number TEXT NOT NULL,         -- e.g., '3', '5', '10'
    section_title TEXT NOT NULL,          -- e.g., 'Operators and Expressions'
    file_path TEXT NOT NULL,              -- Relative path from learning/ directory
    display_order INTEGER NOT NULL,       -- For ordering within a phase
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(language, section_number)
);

-- Link questions to learning resources
CREATE TABLE question_learning_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    question_id INTEGER NOT NULL,
    learning_resource_id INTEGER NOT NULL,
    relevance_score REAL DEFAULT 1.0,     -- How relevant (0-1) this resource is to the question
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (question_id) REFERENCES questions(id),
    FOREIGN KEY (learning_resource_id) REFERENCES learning_resources(id),
    UNIQUE(question_id, learning_resource_id)
);

-- Track which resources students have viewed
CREATE TABLE user_learning_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    learning_resource_id INTEGER NOT NULL,
    first_viewed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_viewed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    view_count INTEGER DEFAULT 1,
    completion_percentage REAL DEFAULT 0,  -- Future feature: track reading progress
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (learning_resource_id) REFERENCES learning_resources(id),
    UNIQUE(user_id, learning_resource_id)
);
```

### 2. Go Models

```go
// LearningResource represents a learning guide document
type LearningResource struct {
    ID            int       `json:"id"`
    Language      string    `json:"language"`       // java, python
    Phase         int       `json:"phase"`          // 1-5
    SectionNumber string    `json:"section_number"` // e.g., "3", "5", "10"
    SectionTitle  string    `json:"section_title"`
    FilePath      string    `json:"file_path"`
    DisplayOrder  int       `json:"display_order"`
    IsActive      bool      `json:"is_active"`
    CreatedAt     time.Time `json:"created_at"`
}

// QuestionLearningLink links questions to learning resources
type QuestionLearningLink struct {
    ID                  int     `json:"id"`
    QuestionID          int     `json:"question_id"`
    LearningResourceID  int     `json:"learning_resource_id"`
    RelevanceScore      float64 `json:"relevance_score"`
}

// UserLearningProgress tracks student's learning material viewing
type UserLearningProgress struct {
    ID                   int       `json:"id"`
    UserID               int       `json:"user_id"`
    LearningResourceID   int       `json:"learning_resource_id"`
    FirstViewedAt        time.Time `json:"first_viewed_at"`
    LastViewedAt         time.Time `json:"last_viewed_at"`
    ViewCount            int       `json:"view_count"`
    CompletionPercentage float64   `json:"completion_percentage"`
}

// PhaseOutline represents a phase in the curriculum
type PhaseOutline struct {
    Phase       int                `json:"phase"`
    Title       string             `json:"title"`
    Description string             `json:"description"`
    Resources   []LearningResource `json:"resources"`
}
```

### 3. API Endpoints

#### Learning Resources
- `GET /api/learning/outline` - Get complete curriculum outline organized by phases
- `GET /api/learning/phases` - Get list of all phases
- `GET /api/learning/phase/{phase}` - Get resources for a specific phase
- `GET /api/learning/resource/{id}` - Get a specific learning resource's metadata
- `GET /api/learning/content/{id}` - Get the markdown content of a learning resource
- `GET /api/learning/search?q={query}` - Search learning resources
- `POST /api/learning/track-view` - Record that user viewed a resource

#### Learning-Quiz Integration
- `GET /api/quiz/question/{id}/learning` - Get learning resources linked to a question
- `GET /api/learning/resource/{id}/questions` - Get quiz questions related to a resource
- `POST /api/learning/link-question` - (Admin) Link a question to a learning resource

## UI/UX Design

### Learning Wizard Layout

The Learning Wizard will be accessible from the dashboard with three main views:

#### 1. **Curriculum Overview Mode**
```
+----------------------------------------------------------+
|  [Dashboard] [Learn] [Practice] [Profile]               |
+----------------------------------------------------------+
|                                                           |
|  Learning Materials                           [Language] |
|  +-----------------+  +-----------------------------+    |
|  | Phase Outline   |  | Welcome to Learning Center  |    |
|  |                 |  |                             |    |
|  | ▼ Phase 1:      |  | Select a topic from the     |    |
|  |   Fundamentals  |  | outline to begin studying.  |    |
|  |   □ Variables   |  |                             |    |
|  |   □ Operators   |  | Your progress:              |    |
|  |   □ Methods     |  | ▓▓▓▓░░░░░░ 40% complete     |    |
|  |   □ Arrays      |  |                             |    |
|  |   □ Strings     |  |                             |    |
|  |                 |  |                             |    |
|  | ▶ Phase 2: OOP  |  |                             |    |
|  |                 |  |                             |    |
|  | ▶ Phase 3:      |  |                             |    |
|  |   Advanced Data |  |                             |    |
|  +-----------------+  +-----------------------------+    |
+----------------------------------------------------------+
```

#### 2. **Study Mode** (Reading a Topic)
```
+----------------------------------------------------------+
|  [Dashboard] [Learn] [Practice] [Profile]               |
+----------------------------------------------------------+
| ← Back to Outline                  Java: Operators      |
+----------------------------------------------------------+
| Phase Outline   | Topic Content                         |
|                 | +-----------------------------------+ |
| ▼ Phase 1       | | # Operators and Expressions       | |
|   ✓ Variables   | |                                   | |
|   ► Operators   | | Let's explore how to perform      | |
|   □ Methods     | | operations on data in Java...     | |
|   □ Arrays      | |                                   | |
|   □ Strings     | | ## Arithmetic Operators           | |
|                 | | Java provides standard math...    | |
| ▶ Phase 2: OOP  | |                                   | |
|                 | | ```java                           | |
| ▶ Phase 3       | | int sum = 5 + 3;  // 8            | |
|                 | | int diff = 10 - 4; // 6           | |
+-----------------+ | ```                               | |
| Actions:        | |                                   | |
| [Take Quiz]     | | [Show more...]                    | |
| [Practice]      | +-----------------------------------+ |
+-----------------+ +-----------------------------------+ |
```

#### 3. **Quiz Review Mode** (After Wrong Answer)
```
+----------------------------------------------------------+
| Quiz Results - Question 3                                |
+----------------------------------------------------------+
| ❌ Incorrect                                             |
|                                                          |
| Question: What is the result of 5 % 2 in Java?          |
| Your answer: 2                                           |
| Correct answer: 1                                        |
|                                                          |
| Explanation:                                             |
| The modulo operator (%) returns the remainder...         |
|                                                          |
| 📚 Related Learning Material:                            |
| +------------------------------------------------------+ |
| | Java: Operators and Expressions                      | |
| | Section 3.1: Arithmetic Operators                    | |
| |                                                      | |
| | [Review This Topic]                                  | |
| +------------------------------------------------------+ |
|                                                          |
| [Next Question]                                          |
+----------------------------------------------------------+
```

### 4. Learning Wizard HTML Page

The main learning wizard interface (`learning.html`)

### 5. Dashboard Integration

Add a "Learn" button/tab to the main dashboard that opens the learning wizard.

## Implementation Phases

### Phase 1: Backend Infrastructure
1. Create database migrations for new tables
2. Add Go models
3. Implement API endpoints for serving learning content
4. Create seed data linking existing resources

### Phase 2: Basic UI
1. Create `learning.html` with two-pane layout
2. Implement phase outline navigation (left pane)
3. Implement content display (right pane)
4. Add markdown rendering for learning materials

### Phase 3: Quiz Integration
1. Update Question model to support learning resource links
2. Add API endpoints for question-learning linkage
3. Modify quiz results UI to show related learning materials
4. Add "Review Material" buttons after wrong answers

### Phase 4: Tracking & Analytics
1. Implement view tracking
2. Add progress indicators
3. Create learning analytics for students
4. Show completion percentages

### Phase 5: Enhanced Features
1. Add search functionality
2. Implement bookmarks/favorites
3. Add note-taking capability
4. Create study plans/paths

## File Structure

```
learning/
├── java/
│   ├── 2-variables-and-data-types.md
│   ├── 3-operators-and-expressions.md
│   ├── 4-control-structures.md
│   ├── 5-methods.md
│   ├── 6-arrays.md
│   └── 7-strings.md
└── python/
    ├── 2-variables-and-data-types.md
    ├── 3-operators-and-expressions.md
    ├── 4-strings.md
    ├── 5-control-structures.md
    ├── 6-data-structures.md
    └── 10-functions.md

static/
├── learning.html          (New - Learning Wizard UI)
├── dashboard.html         (Updated - Add Learn button)
└── js/
    ├── learning.js        (New - Learning wizard logic)
    └── dashboard.js       (Updated)

api/
├── learning.go           (New - Learning resource handlers)
└── handlers.go           (Updated - Link to learning materials)

database/
├── migrations/           (New - Database schema updates)
└── seed_learning.sql    (New - Seed learning resources)
```

## Markdown Rendering

Use a JavaScript markdown library (like `marked.js`) to render markdown content in the browser:

```html
<script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
<script>
    // Render markdown
    document.getElementById('content').innerHTML = marked.parse(markdownText);
</script>
```

Add syntax highlighting for code blocks:
```html
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/highlight.js/styles/github.css">
<script src="https://cdn.jsdelivr.net/npm/highlight.js"></script>
<script>
    // After markdown rendering
    document.querySelectorAll('pre code').forEach((block) => {
        hljs.highlightBlock(block);
    });
</script>
```

## Success Criteria

1. ✅ Students can browse curriculum by phase
2. ✅ Students can read learning materials within Athena
3. ✅ Quiz questions link to relevant learning materials
4. ✅ Students can review materials after wrong answers
5. ✅ Progress is tracked and displayed
6. ✅ Materials are organized and easy to navigate

## Future Enhancements

- **AI-Powered Recommendations**: Suggest which topics to study based on quiz performance
- **Spaced Repetition Integration**: Remind students to review specific topics
- **Interactive Examples**: Add runnable code examples
- **Video Integration**: Embed video tutorials
- **Collaborative Notes**: Allow students to share notes
- **Practice Problems**: Add practice exercises separate from quizzes
