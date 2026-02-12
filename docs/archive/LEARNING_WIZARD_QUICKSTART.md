# Learning Wizard - Quick Start Guide

## What's New?

The Learning Wizard is a new feature in Athena that integrates learning materials directly into the platform, creating a seamless study-to-quiz workflow.

## For Students

### Access the Learning Center
1. Log in to Athena
2. From the dashboard, click **"📚 Learning Center"**
3. Browse topics by phase or search for specific subjects
4. Click on any topic to read the material
5. Take quizzes to test your knowledge

### Study-Quiz-Review Workflow
```
1. Study → Read learning materials in Learning Center
2. Quiz → Take quiz on the topic
3. Review → If you get answers wrong, click suggested materials to review
4. Repeat → Retake quiz to improve
```

### Features
- ✅ Organized curriculum (5 phases from beginner to advanced)
- ✅ Both Java and Python materials
- ✅ Search functionality
- ✅ Progress tracking
- ✅ Syntax-highlighted code examples
- ✅ Direct links from quiz questions to review materials

## For Instructors

### Quick Setup (5 minutes)

1. **Apply database migration**:
   ```bash
   chmod +x setup_learning.sh
   ./setup_learning.sh
   ```

2. **Restart server**:
   ```bash
   go run main.go
   ```

3. **Done!** Students can now access `/learning.html`

### Link Questions to Materials

Create links so students see the right materials when they get questions wrong:

```sql
-- Link Java operator questions to Operators learning material
INSERT INTO question_learning_links (question_id, learning_resource_id, relevance_score)
SELECT q.id, lr.id, 1.0
FROM questions q
JOIN learning_resources lr ON lr.section_number = '3' AND lr.language = 'java'
WHERE q.topic_question_id LIKE '3.%';
```

### Add New Learning Materials

1. Create markdown file in `learning/java/` or `learning/python/`
2. Add entry to `learning_resources` table:
   ```sql
   INSERT INTO learning_resources 
   (language, phase, section_number, section_title, file_path, display_order, is_active)
   VALUES ('java', 2, '8', 'Classes and Objects', 'java/8-classes-objects.md', 1, 1);
   ```

## Architecture Overview

```
┌─────────────────┐
│   Dashboard     │
│  [Learn] [Quiz] │
└────────┬────────┘
         │
    ┌────┴──────────────────────┐
    │                           │
┌───▼────────────┐    ┌────────▼──────┐
│ Learning Center│    │  Quiz System  │
│                │    │               │
│ • Browse       │◄───┤ • Questions   │
│ • Search       │    │ • Answers     │
│ • Read         │    │ • Review links│
│ • Track        │    │               │
└────────────────┘    └───────────────┘
         │
         │ (Stores in DB)
         ▼
┌─────────────────┐
│ learning_resources
│ question_learning_links
│ user_learning_progress
└─────────────────┘
```

## Key Files

| File | Purpose |
|------|---------|
| `static/learning.html` | Learning wizard UI |
| `static/js/learning.js` | Learning wizard logic |
| `api/learning.go` | Backend API handlers |
| `database/migrations/004_learning_resources.sql` | Database schema |
| `database/seed_learning_resources.sql` | Initial data |
| `LEARNING_WIZARD_DESIGN.md` | Full design documentation |
| `LEARNING_WIZARD_IMPLEMENTATION.md` | Implementation guide |

## Current Status

### ✅ Implemented
- Learning Center UI with two-pane layout
- 12 learning resources (6 Java, 6 Python) for Phase 1
- API endpoints for serving content
- Database schema and models
- Quiz integration (shows learning links on wrong answers)
- Progress tracking
- Search functionality
- Language filtering

### 📋 Ready to Add
- Link existing quiz questions to learning materials
- Create learning materials for Phases 2-5
- Add more Phase 1 materials as needed

### 🔮 Future Enhancements
- Interactive code examples
- Video integration
- AI-powered recommendations
- Collaborative features
- Note-taking

## Testing

### Test the Learning Center
1. Navigate to `http://localhost:9206/learning.html`
2. Should see Phase 1 with 6 Java and 6 Python topics
3. Click on "3. Operators and Expressions" (Java)
4. Should see formatted markdown content
5. Try searching for "control"
6. Should find Control Structures resources

### Test Quiz Integration
1. Take a quiz on a topic you've studied
2. Get an answer wrong intentionally
3. Should see "📚 Review These Topics:" section
4. Click the link to open learning material

### Test Progress Tracking
```sql
-- Check if views are being tracked
SELECT * FROM user_learning_progress WHERE user_id = 1;
```

## Support & Documentation

- **Design**: See [LEARNING_WIZARD_DESIGN.md](LEARNING_WIZARD_DESIGN.md)
- **Implementation**: See [LEARNING_WIZARD_IMPLEMENTATION.md](LEARNING_WIZARD_IMPLEMENTATION.md)
- **Phased Plan**: See [learning/PHASED_PLAN.md](learning/PHASED_PLAN.md)

## Questions?

Common issues and solutions in `LEARNING_WIZARD_IMPLEMENTATION.md` → Troubleshooting section.

---

**Ready to go!** Run `./setup_learning.sh` and start exploring the Learning Center.
