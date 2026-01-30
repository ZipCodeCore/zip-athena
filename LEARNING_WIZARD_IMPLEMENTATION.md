# Learning Wizard Implementation Guide

## Overview

The Learning Wizard has been successfully implemented in Athena, providing students with an integrated learning experience that connects study materials with quiz questions.

## Features Implemented

### 1. Learning Center Interface ✅
- **Two-pane layout**: Left pane shows curriculum outline, right pane shows content
- **Phase-based organization**: Learning materials organized into 5 phases from fundamentals to advanced topics
- **Language filtering**: Toggle between Java, Python, or both languages
- **Search functionality**: Search for topics across all learning materials
- **Progress tracking**: Automatically tracks which resources students have viewed
- **Markdown rendering**: Beautiful rendering of learning materials with syntax highlighting

### 2. Database Schema ✅
Three new tables added:
- `learning_resources`: Stores metadata about learning materials
- `question_learning_links`: Links quiz questions to relevant learning materials
- `user_learning_progress`: Tracks student viewing history

### 3. API Endpoints ✅
- `GET /learning/outline` - Get complete curriculum outline
- `GET /learning/phase/{phase}` - Get resources for a specific phase
- `GET /learning/resource/{id}` - Get resource metadata
- `GET /learning/resource/{id}/content` - Get markdown content
- `GET /learning/search?q={query}` - Search resources
- `POST /learning/track-view` - Track user views
- `GET /quiz/question/{id}/learning` - Get resources linked to a question

### 4. Quiz Integration ✅
- When students get answers wrong, related learning materials are automatically displayed
- Students can click links to review topics in the Learning Center
- Seamless navigation between quiz and learning materials

### 5. Dashboard Integration ✅
- New "Learning Center" button on student dashboard
- Quick access to study materials before taking quizzes

## Setup Instructions

### 1. Apply Database Migration

Run the setup script:
```bash
chmod +x setup_learning.sh
./setup_learning.sh
```

Or manually apply the migration:
```bash
sqlite3 athena.db < database/migrations/004_learning_resources.sql
sqlite3 athena.db < database/seed_learning_resources.sql
```

### 2. Restart the Server

```bash
go run main.go
```

### 3. Access the Learning Center

Navigate to: `http://localhost:9206/learning.html`

## For Instructors: Linking Questions to Learning Materials

### Manual Method (SQL)

To link a quiz question to a learning resource:

```sql
-- Find the resource ID
SELECT id, section_title FROM learning_resources 
WHERE language = 'java' AND section_number = '3';

-- Find the question ID
SELECT id, question_text FROM questions WHERE topic_question_id = '3.17';

-- Create the link
INSERT INTO question_learning_links (question_id, learning_resource_id, relevance_score)
VALUES (123, 5, 1.0);
```

### Recommended Linking Strategy

1. **Tag-Based Linking**: Questions with tags like `operators`, `arithmetic` should link to "Operators and Expressions"
2. **Topic-Based Linking**: Use the `topic_question_id` prefix (e.g., `3.17` → section 3)
3. **Relevance Score**: Use 1.0 for highly relevant, 0.5 for somewhat relevant

### Example Linking Script

```sql
-- Link all Java operator questions (topic 3) to Operators resource
INSERT INTO question_learning_links (question_id, learning_resource_id, relevance_score)
SELECT q.id, lr.id, 1.0
FROM questions q
JOIN learning_resources lr ON lr.section_number = '3' AND lr.language = 'java'
WHERE q.topic_question_id LIKE '3.%' 
  AND NOT EXISTS (
    SELECT 1 FROM question_learning_links 
    WHERE question_id = q.id AND learning_resource_id = lr.id
  );

-- Link all Python control structure questions to Control Structures resource
INSERT INTO question_learning_links (question_id, learning_resource_id, relevance_score)
SELECT q.id, lr.id, 1.0
FROM questions q
JOIN topics t ON q.topic_id = t.id
JOIN learning_resources lr ON lr.section_number = '5' AND lr.language = 'python'
WHERE t.name LIKE '%Control%' 
  AND NOT EXISTS (
    SELECT 1 FROM question_learning_links 
    WHERE question_id = q.id AND learning_resource_id = lr.id
  );
```

## Student Workflow

1. **Study First**:
   - Go to Learning Center
   - Browse by phase or search for topics
   - Read through learning materials
   - Materials are automatically tracked

2. **Take Quiz**:
   - Click "Take Quiz" from Learning Center or Dashboard
   - Answer questions
   - If wrong, see related learning materials

3. **Review**:
   - Click on suggested learning materials
   - Review the concepts
   - Retake quiz to improve

## File Structure

```
├── static/
│   ├── learning.html              (New - Learning wizard UI)
│   ├── dashboard.html             (Updated - Added Learning Center link)
│   └── js/
│       ├── learning.js            (New - Learning wizard logic)
│       └── quiz.js                (Updated - Shows learning links)
├── api/
│   └── learning.go                (New - Learning API handlers)
├── models/
│   └── models.go                  (Updated - Learning resource models)
├── database/
│   ├── migrations/
│   │   └── 004_learning_resources.sql  (New - Schema)
│   └── seed_learning_resources.sql     (New - Initial data)
├── learning/
│   ├── java/                      (Existing learning materials)
│   │   ├── 2-variables-and-data-types.md
│   │   ├── 3-operators-and-expressions.md
│   │   ├── 4-control-structures.md
│   │   ├── 5-methods.md
│   │   ├── 6-arrays.md
│   │   └── 7-strings.md
│   └── python/                    (Existing learning materials)
│       ├── 2-variables-and-data-types.md
│       ├── 3-operators-and-expressions.md
│       ├── 4-strings.md
│       ├── 5-control-structures.md
│       ├── 6-data-structures.md
│       └── 10-functions.md
├── main.go                        (Updated - Added learning routes)
├── setup_learning.sh              (New - Setup script)
└── LEARNING_WIZARD_DESIGN.md      (New - Design documentation)
```

## Current Learning Resources

### Java (Phase 1 - Active)
1. Variables and Data Types
2. Operators and Expressions  
3. Control Structures
4. Methods (Functions)
5. Arrays
6. Strings

### Python (Phase 1 - Active)
1. Variables and Data Types
2. Operators and Expressions
3. Strings
4. Control Structures
5. Data Structures
6. Functions

**Note**: Phases 2-5 resources are seeded but marked as `is_active = 0` until materials are created.

## Future Enhancements

### Phase 2 (Planned)
- [ ] Interactive code examples with in-browser execution
- [ ] Embedded video tutorials
- [ ] Note-taking functionality
- [ ] Bookmarks/favorites

### Phase 3 (Planned)
- [ ] AI-powered recommendations based on quiz performance
- [ ] Personalized study plans
- [ ] Progress analytics and insights
- [ ] Collaborative study features

### Phase 4 (Planned)
- [ ] Auto-linking questions using AI/ML
- [ ] Difficulty adjustment based on learning progress
- [ ] Spaced repetition integration
- [ ] Practice problems separate from quizzes

## Troubleshooting

### Learning materials not showing?
- Check that migration was applied: `sqlite3 athena.db ".tables"` should show `learning_resources`
- Check that resources were seeded: `sqlite3 athena.db "SELECT COUNT(*) FROM learning_resources;"`

### Links not appearing in quiz?
- Check that question-learning links exist: `sqlite3 athena.db "SELECT COUNT(*) FROM question_learning_links;"`
- Verify the question ID is correct

### Resource content not loading?
- Check file paths in database match actual files in `learning/` directory
- Verify file permissions allow reading

## API Examples

### Get curriculum outline
```bash
curl http://localhost:9206/learning/outline
```

### Get Java resources only
```bash
curl http://localhost:9206/learning/outline?language=java
```

### Get Phase 1 resources
```bash
curl http://localhost:9206/learning/phase/1
```

### Search for "operators"
```bash
curl http://localhost:9206/learning/search?q=operators
```

### Get resource content
```bash
curl http://localhost:9206/learning/resource/1/content
```

## Success Metrics

Monitor these metrics to gauge effectiveness:

1. **Engagement**:
   - Number of students using Learning Center
   - Average time spent on learning materials
   - Number of resources viewed per student

2. **Learning Outcomes**:
   - Quiz performance before/after using Learning Center
   - Improvement in scores after reviewing materials
   - Reduction in wrong answers on retakes

3. **Resource Effectiveness**:
   - Which resources are viewed most
   - Correlation between resource views and quiz success
   - Student feedback on resource quality

## Support

For questions or issues:
1. Check the design document: `LEARNING_WIZARD_DESIGN.md`
2. Review API code: `api/learning.go`
3. Check database schema: `database/migrations/004_learning_resources.sql`

---

**Status**: ✅ Fully Implemented
**Version**: 1.0
**Date**: January 30, 2026
