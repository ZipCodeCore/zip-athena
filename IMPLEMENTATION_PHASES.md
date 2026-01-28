# Athena Implementation Phases

## Executive Summary

This document organizes all Feature Requirements from SPECIFICATION.md into logical implementation phases. Each phase builds upon the previous ones and represents a deployable increment of functionality.

### Status Overview

- ✅ **Phase 0**: Backend Foundation - COMPLETE
- ✅ **Phase 1**: Core Student Experience - COMPLETE  
- ✅ **Phase 2**: Spaced Repetition (FSRS) - COMPLETE
- 🔲 **Phase 3**: Instructor Analytics & Management
- 🔲 **Phase 4**: Advanced Student Features & Testing
- 🔲 **Phase 5**: Content Management & Quality
- 🔲 **Phase 6**: Performance & Polish

---

## Phase 3: Instructor Analytics & Management Features

### Objectives
Enable instructors to monitor student progress, identify struggling students, and gain insights into class performance. This phase transforms Athena from a student-only tool into a teaching platform.

### Dependencies
- Phase 1 (User authentication and quiz system)
- Phase 2 (FSRS algorithm and tracking)

### Feature Requirements Included

#### FR-INSTRUCTOR-001: Student List & Search
**Priority**: P0

**Description**: Instructors can view and search all students with filtering and sorting capabilities.

**GitHub Issue Template**:
```markdown
## Feature: Student List & Search

### Description
Implement an instructor dashboard page that displays all registered students with their key performance metrics. Enable search and filtering to help instructors quickly find and identify students.

### Acceptance Criteria
- [ ] Display list of all students with key metrics
- [ ] Search by name or email
- [ ] Sort by various metrics (accuracy, completion rate, last activity)
- [ ] Filter by cohort/enrollment period
- [ ] Flag struggling students (accuracy < 60%) with visual indicator

### Student List Columns
- Name
- Email
- Total Quiz Attempts
- Overall Accuracy %
- Last Activity Date
- Status Indicator (green/yellow/red based on performance)

### API Requirements
- GET /instructor/students - List all students with metrics
- GET /instructor/students/search?q={query} - Search students

### Unit Tests Required
1. Test student list retrieval with metrics calculation
2. Test search functionality (name and email)
3. Test sorting by each column
4. Test filtering by cohort
5. Test struggling student flagging logic (accuracy < 60%)
6. Test pagination for large student lists
7. Test empty state when no students exist
8. Test access control (only instructors can access)

### UI Components
- Student table with sortable columns
- Search input field
- Filter dropdowns
- Status indicators (colored badges)
- Responsive design for mobile

### Database Queries Needed
- Aggregate user statistics (quiz attempts, accuracy)
- Join user and quiz attempt data
- Filter and search operations
```

#### FR-INSTRUCTOR-002: Individual Student Analytics
**Priority**: P0

**Description**: Instructors can drill into individual student performance with detailed analytics.

**GitHub Issue Template**:
```markdown
## Feature: Individual Student Analytics

### Description
Create a detailed analytics view for individual students, allowing instructors to understand student progress, identify weak areas, and compare performance to cohort averages.

### Acceptance Criteria
- [ ] View student's progress dashboard
- [ ] See all quiz attempt history with details
- [ ] Identify weak topic areas with visual indicators
- [ ] Identify weak conceptual areas (tags)
- [ ] View test results (when Phase 4 complete)
- [ ] Compare to cohort averages with charts
- [ ] Export student report as PDF

### Analytics Views
1. **Summary**: Overall stats vs. cohort
2. **Topics**: Mastery by topic with trend lines
3. **Tags**: Conceptual strengths/weaknesses heat map
4. **History**: Timeline of all attempts
5. **Tests**: Formal assessment results (Phase 4)

### API Requirements
- GET /instructor/students/{id}/analytics - Get comprehensive analytics
- GET /instructor/students/{id}/quiz-history - Get quiz history
- GET /instructor/students/{id}/weak-areas - Get weak topics and tags
- GET /instructor/students/{id}/report - Export PDF report

### Unit Tests Required
1. Test analytics data aggregation for a student
2. Test weak topic identification algorithm
3. Test weak tag identification algorithm
4. Test cohort comparison calculations
5. Test quiz history retrieval with pagination
6. Test PDF report generation
7. Test data for students with no attempts
8. Test data for students with many attempts (performance)

### UI Components
- Student header with photo and basic info
- Tabbed interface for different views
- Charts (line, bar, radar for topic mastery)
- Tables for quiz history
- Export button
- Comparison overlays (student vs. cohort)

### Dependencies
- Charting library (e.g., Chart.js)
- PDF generation library (e.g., pdfmake or server-side)
```

#### FR-INSTRUCTOR-003: Cohort Analytics Dashboard
**Priority**: P0

**Description**: Instructors view aggregate class performance and identify trends.

**GitHub Issue Template**:
```markdown
## Feature: Cohort Analytics Dashboard

### Description
Provide instructors with a comprehensive view of class/cohort performance, enabling them to identify common problem areas and track overall progress.

### Acceptance Criteria
- [ ] Select cohort/class period from dropdown
- [ ] View aggregate statistics (avg accuracy, completion rate)
- [ ] Topic-by-topic class performance with charts
- [ ] Identify common struggling concepts (tags)
- [ ] See distribution charts (histograms of scores)
- [ ] Export cohort report as CSV/PDF

### Dashboard Sections
1. **Overview**: Class averages, completion rates
2. **Topics**: Avg accuracy per topic, struggling topics
3. **Students**: Leaderboard and struggling student list
4. **Trends**: Week-over-week progress line chart
5. **Questions**: Hardest questions for class (lowest accuracy)

### API Requirements
- GET /instructor/cohorts - List available cohorts
- GET /instructor/cohorts/{id}/analytics - Get cohort analytics
- GET /instructor/cohorts/{id}/topic-performance - Topic breakdown
- GET /instructor/cohorts/{id}/struggling-concepts - Common weak tags
- GET /instructor/cohorts/{id}/export - Export report

### Unit Tests Required
1. Test cohort analytics aggregation
2. Test topic performance calculations
3. Test struggling concept identification
4. Test score distribution calculation
5. Test trend analysis (week-over-week)
6. Test hardest questions identification
7. Test export functionality
8. Test with empty cohort (no students)
9. Test with large cohort (performance)

### UI Components
- Cohort selector dropdown
- Key metrics cards (avg accuracy, completion %)
- Topic performance bar chart
- Score distribution histogram
- Struggling students table
- Trends line chart
- Export buttons

### Calculations Needed
- Average accuracy across all students
- Completion rate (students who completed vs. started)
- Standard deviation of scores
- Percentile rankings
- Week-over-week change percentages
```

#### FR-INSTRUCTOR-007: Topic Management
**Priority**: P2

**Description**: Instructors can manage topic categories.

**GitHub Issue Template**:
```markdown
## Feature: Topic Management

### Description
Allow instructors to create, edit, and manage topic categories for organizing questions.

### Acceptance Criteria
- [ ] Create new topics with name and description
- [ ] Edit topic names and descriptions
- [ ] Deactivate topics (soft delete)
- [ ] Reorder topic display order
- [ ] Set topic prerequisites (future enhancement)

### API Requirements
- POST /instructor/topics - Create topic
- PUT /instructor/topics/{id} - Update topic
- DELETE /instructor/topics/{id} - Deactivate topic
- PUT /instructor/topics/reorder - Reorder topics

### Unit Tests Required
1. Test topic creation with validation
2. Test topic name uniqueness enforcement
3. Test topic update
4. Test topic deactivation (not deletion)
5. Test topic reordering
6. Test validation errors (empty name, duplicate)
7. Test access control (only instructors)

### UI Components
- Topic list with drag-and-drop reordering
- Create/edit topic modal
- Deactivate confirmation dialog
- Form validation

### Business Rules
- Topic names must be unique
- Cannot delete topic with existing questions (deactivate only)
- Display order is persisted
```

### Phase 3 Success Criteria
- All instructor endpoints functional
- Instructor can view and search students
- Instructor can drill into student analytics
- Instructor can view cohort performance
- All Phase 3 unit tests passing (28+ tests)
- Instructor UI is responsive and intuitive
- Documentation complete

### Estimated Testing
**Unit Tests**: ~28 tests minimum
- Student List & Search: 8 tests
- Individual Student Analytics: 8 tests
- Cohort Analytics: 9 tests
- Topic Management: 7 tests

---

## Phase 4: Advanced Student Features & Testing

### Objectives
Enhance the student experience with advanced features including timed tests, quiz history review, and tag-based practice. Enable formal assessments for instructors.

### Dependencies
- Phase 1 (Core student experience)
- Phase 2 (FSRS algorithm)
- Phase 3 (Instructor management)

### Feature Requirements Included

#### FR-STUDENT-004: Quiz Attempt History
**Priority**: P1

**Description**: Students can review past quiz attempts with detailed results.

**GitHub Issue Template**:
```markdown
## Feature: Quiz Attempt History

### Description
Allow students to view and review their past quiz attempts to learn from mistakes and track progress over time.

### Acceptance Criteria
- [ ] List all past quiz attempts with date, topic, score
- [ ] Sort by date (newest first)
- [ ] Filter by topic
- [ ] Drill into specific attempt to see individual questions
- [ ] Show correct answers for incorrect questions
- [ ] Show time spent per question (if tracked)

### API Requirements
- GET /students/quiz-history - Get user's quiz history
- GET /students/quiz-history/{id} - Get specific attempt details
- GET /students/quiz-history?topic_id={id} - Filter by topic

### Unit Tests Required
1. Test quiz history retrieval for user
2. Test sorting by date
3. Test filtering by topic
4. Test drill-down into specific attempt
5. Test display of correct answers
6. Test pagination for users with many attempts
7. Test empty state (no quiz attempts yet)

### UI Components
- Quiz history table with columns: Date, Topic, Score, Questions
- Filter dropdown for topics
- Detail modal/page for viewing questions
- Correct/incorrect indicators
- Review mode (show what was wrong)

### Database Queries
- Join quiz_attempts with questions and answers
- Filter by user_id and optionally topic_id
- Sort by attempt date descending
```

#### FR-STUDENT-006: Timed Test Taking
**Priority**: P1

**Description**: Students complete formal timed assessments with countdown timer.

**GitHub Issue Template**:
```markdown
## Feature: Timed Test Taking

### Description
Implement formal timed tests that simulate exam conditions, with countdown timer, no-pause policy, and auto-submission.

### Acceptance Criteria
- [ ] View available tests with descriptions
- [ ] See test duration and question count before starting
- [ ] Accept test rules before starting (no pause, single attempt)
- [ ] Display countdown timer during test
- [ ] Auto-submit when timer expires
- [ ] Cannot review questions mid-test
- [ ] Receive final score immediately after submission
- [ ] Warning at 5 minutes remaining

### Test Flow
1. Student views available tests
2. Clicks "Start Test" and reads rules
3. Accepts rules (checkbox confirmation)
4. Timer starts (e.g., 60 minutes)
5. Questions displayed sequentially (can navigate back/forward)
6. Timer visible on all pages
7. Warning at 5 minutes remaining
8. Auto-submit at 0:00 or manual submit
9. Display final score and percentile

### API Requirements
- GET /tests/available - List available tests
- GET /tests/{id} - Get test details
- POST /tests/{id}/start - Start test attempt
- PUT /test-attempts/{id}/answer - Submit answer during test
- POST /test-attempts/{id}/submit - Submit test (manual or auto)
- GET /test-attempts/{id}/results - Get test results

### Unit Tests Required
1. Test available tests retrieval
2. Test test start creates attempt record
3. Test timer countdown logic
4. Test auto-submit when timer expires
5. Test manual submit before timer expires
6. Test cannot start test twice
7. Test answer submission during test
8. Test score calculation
9. Test percentile calculation
10. Test warning trigger at 5 minutes

### UI Components
- Available tests list
- Test rules modal with acceptance checkbox
- Countdown timer (always visible)
- Question navigation (back/forward buttons)
- Submit test button
- Warning modal at 5 minutes
- Auto-submit notification
- Results display

### Technical Requirements
- Client-side timer with server validation
- Test attempt state tracking
- Prevent page refresh losing progress
- Local storage backup of answers
```

#### FR-STUDENT-007: Tag-Based Review
**Priority**: P2

**Description**: Students can practice specific concepts via tags.

**GitHub Issue Template**:
```markdown
## Feature: Tag-Based Review

### Description
Allow students to focus on specific concepts by generating quizzes filtered to particular tags, overriding the normal FSRS algorithm.

### Acceptance Criteria
- [ ] View list of all tags with question counts
- [ ] Filter quizzes to specific tag(s)
- [ ] See accuracy rate per tag
- [ ] Generate tag-specific quiz (overrides FSRS)
- [ ] Support multi-tag selection (AND/OR logic)

### API Requirements
- GET /tags - List all tags with question counts
- GET /tags/stats?user_id={id} - Get accuracy per tag
- POST /quiz/generate/by-tags - Generate tag-filtered quiz

### Unit Tests Required
1. Test tag list retrieval with counts
2. Test tag stats calculation per user
3. Test quiz generation with single tag
4. Test quiz generation with multiple tags (AND)
5. Test quiz generation with multiple tags (OR)
6. Test tag accuracy calculation
7. Test empty state (no questions for tag)

### UI Components
- Tag cloud or list with counts
- Tag selection interface (checkboxes or pills)
- AND/OR toggle for multi-tag
- Accuracy indicators per tag
- Generate quiz button
- Tag filter on quiz history

### Business Logic
- Tag filtering overrides FSRS prioritization
- Still respects difficulty progression
- Updates FSRS state after quiz completion
```

#### FR-INSTRUCTOR-005: Test Management
**Priority**: P1

**Description**: Instructors can create and manage formal tests.

**GitHub Issue Template**:
```markdown
## Feature: Test Management

### Description
Enable instructors to create formal tests, assign them to cohorts, and view results.

### Acceptance Criteria
- [ ] Create new test with title, description
- [ ] Add questions to test (search and select)
- [ ] Set time limit (minutes)
- [ ] Set availability window (start/end dates)
- [ ] Assign to specific cohort(s)
- [ ] View test results by student
- [ ] Export test results as CSV

### Test Creation Flow
1. Set test metadata (title, description, time limit)
2. Search/filter question bank
3. Add questions to test (drag-and-drop or click)
4. Preview test as student would see
5. Set availability window
6. Assign to cohort(s)
7. Publish

### API Requirements
- POST /instructor/tests - Create test
- PUT /instructor/tests/{id} - Update test
- DELETE /instructor/tests/{id} - Delete test
- POST /instructor/tests/{id}/questions - Add questions
- DELETE /instructor/tests/{id}/questions/{qid} - Remove question
- PUT /instructor/tests/{id}/publish - Publish test
- GET /instructor/tests/{id}/results - View results
- GET /instructor/tests/{id}/export - Export results CSV

### Unit Tests Required
1. Test test creation
2. Test adding questions to test
3. Test removing questions from test
4. Test time limit validation
5. Test availability window validation
6. Test cohort assignment
7. Test publish action
8. Test results retrieval
9. Test export functionality
10. Test access control

### UI Components
- Test creation form
- Question search and selection interface
- Drag-and-drop question ordering
- Preview mode
- Availability date pickers
- Cohort selection multi-select
- Results table
- Export button

### Data Model
- Test metadata (title, description, time_limit)
- Test questions (ordered)
- Test assignments (test_id, cohort_id)
- Test attempts (student attempts)
```

#### FR-INSTRUCTOR-006: Question Performance Analytics
**Priority**: P1

**Description**: Instructors can analyze question effectiveness and identify problematic questions.

**GitHub Issue Template**:
```markdown
## Feature: Question Performance Analytics

### Description
Provide instructors with analytics on how questions perform, helping identify questions that are too easy, too hard, or potentially ambiguous.

### Acceptance Criteria
- [ ] View list of all questions with avg accuracy
- [ ] Identify questions that are too easy (>90% correct)
- [ ] Identify questions that are too hard (<30% correct)
- [ ] See answer distribution (which distractors chosen)
- [ ] Flag potentially ambiguous questions
- [ ] Adjust difficulty ratings based on data

### Question Analytics View
- Question ID and text preview
- Topic and Tags
- Current Difficulty Rating (1-7)
- Number of Attempts
- Average Accuracy %
- Recommended Difficulty (calculated from data)
- Answer Distribution (histogram)
- Flag for Review (if outlier)

### API Requirements
- GET /instructor/questions/analytics - Get analytics for all questions
- GET /instructor/questions/{id}/analytics - Get detailed analytics
- GET /instructor/questions/{id}/answer-distribution - Answer stats
- PUT /instructor/questions/{id}/difficulty - Update difficulty

### Unit Tests Required
1. Test question analytics aggregation
2. Test average accuracy calculation
3. Test too-easy identification (>90%)
4. Test too-hard identification (<30%)
5. Test answer distribution calculation
6. Test recommended difficulty algorithm
7. Test flagging logic for review
8. Test difficulty update
9. Test with questions never attempted
10. Test with questions attempted many times

### UI Components
- Sortable questions table
- Filter by too easy/too hard
- Detail modal with answer distribution chart
- Difficulty adjustment slider
- Flag indicator
- Recommended difficulty badge

### Analytics Calculations
- Average accuracy = correct / total attempts
- Recommended difficulty based on accuracy:
  - >90% → decrease difficulty
  - 70-90% → appropriate
  - 30-70% → appropriate
  - <30% → increase difficulty
- Answer distribution = count per option
```

### Phase 4 Success Criteria
- Students can review quiz history
- Timed tests functional with countdown
- Tag-based practice available
- Instructors can create and manage tests
- Question analytics help improve content
- All Phase 4 unit tests passing (35+ tests)
- UI supports all new features
- Documentation complete

### Estimated Testing
**Unit Tests**: ~35 tests minimum
- Quiz Attempt History: 7 tests
- Timed Test Taking: 10 tests
- Tag-Based Review: 7 tests
- Test Management: 10 tests
- Question Analytics: 10 tests

---

## Phase 5: Content Management & Quality

### Objectives
Enable efficient content creation and management, support bulk imports, and establish quality review processes. This phase focuses on scaling the question bank and maintaining content quality.

### Dependencies
- Phase 1 (Basic CRUD)
- Phase 3 (Instructor management)

### Feature Requirements Included

#### FR-INSTRUCTOR-004: Question Management
**Priority**: P0

**Description**: Instructors can create, edit, and manage questions through a user interface.

**GitHub Issue Template**:
```markdown
## Feature: Question Management UI

### Description
Implement a comprehensive question management interface allowing instructors to create, edit, and manage questions with full metadata.

### Acceptance Criteria
- [ ] Create new question with text, answers, correct answer
- [ ] Assign topic, difficulty (1-7), tags
- [ ] Edit existing questions
- [ ] Deactivate questions (don't delete)
- [ ] Preview questions as students see them
- [ ] Rich text editor for code formatting

### Question Form Fields
- Question Text (rich text editor with code syntax highlighting)
- Topic Dropdown
- Difficulty (1-7 slider with descriptions)
- Tags (multi-select with autocomplete)
- Answer Option 1 (text + correct checkbox)
- Answer Option 2 (text + correct checkbox)
- Answer Option 3 (text + correct checkbox)
- Answer Option 4 (text + correct checkbox)
- Status (active/inactive toggle)

### API Requirements
- GET /instructor/questions - List questions with pagination
- POST /instructor/questions - Create question
- PUT /instructor/questions/{id} - Update question
- DELETE /instructor/questions/{id} - Deactivate question
- GET /instructor/questions/{id}/preview - Preview as student

### Unit Tests Required
1. Test question creation with validation
2. Test question text required
3. Test minimum 2 answers required
4. Test exactly 1 correct answer required
5. Test question update
6. Test question deactivation
7. Test cannot delete question with attempts
8. Test tag assignment
9. Test difficulty validation (1-7)
10. Test preview generation

### UI Components
- Question list table with search and filter
- Create/edit question form
- Rich text editor (e.g., TinyMCE, Quill)
- Code syntax highlighter
- Tag autocomplete input
- Difficulty slider with labels
- Answer options with radio buttons for correct
- Preview modal
- Deactivate confirmation

### Validation Rules
- Question text required, max 1000 chars
- 2-4 answer options required
- Exactly 1 correct answer
- Topic required
- Difficulty 1-7 required
- Tags optional
```

#### FR-CONTENT-001: Bulk Question Import
**Priority**: P1

**Description**: Content admins can import questions in bulk via CSV upload.

**GitHub Issue Template**:
```markdown
## Feature: Bulk Question Import

### Description
Enable bulk import of questions via CSV file upload, with validation, preview, and error reporting.

### Acceptance Criteria
- [ ] Upload CSV file with questions
- [ ] Map CSV columns to question fields
- [ ] Validate data before import
- [ ] Preview questions before finalizing
- [ ] Report errors for invalid rows
- [ ] Support 500+ questions per import

### CSV Format
```csv
Topic,Question,Answer1,Answer2,Answer3,Answer4,CorrectAnswer,Difficulty,Tags
Java,"What is the output of: int x = 5; System.out.println(++x);","5","6","7","Error","6",2,"variables,operators"
Python,"What does len([1,2,3]) return?","1","2","3","4","3",1,"lists,functions"
```

### Import Flow
1. Upload CSV file
2. Parse and validate CSV
3. Show preview table with validation status
4. User reviews and confirms
5. Import valid rows
6. Display summary (success count, error count)
7. Download error report for invalid rows

### API Requirements
- POST /admin/questions/import - Upload CSV
- POST /admin/questions/import/validate - Validate without importing
- POST /admin/questions/import/confirm - Confirm and import

### Unit Tests Required
1. Test CSV parsing
2. Test column mapping
3. Test validation (required fields)
4. Test validation (correct answer exists in options)
5. Test validation (difficulty 1-7)
6. Test duplicate question detection
7. Test batch insert performance (500 questions)
8. Test error reporting for invalid rows
9. Test partial import (some valid, some invalid)
10. Test rollback on error

### UI Components
- File upload dropzone
- Column mapping interface
- Preview table with validation indicators
- Error list with line numbers
- Progress bar during import
- Success/error summary
- Download error report button

### Validation Rules
- Topic must exist or be created
- Question text required
- 2-4 answers required
- Correct answer must match one of the answer options
- Difficulty must be 1-7
- Tags comma-separated, trimmed
```

#### FR-CONTENT-002: Question Tagging System
**Priority**: P0

**Description**: Content admins manage tag taxonomy for organizing questions by concept.

**GitHub Issue Template**:
```markdown
## Feature: Question Tagging System

### Description
Implement a comprehensive tagging system for organizing questions by concepts, enabling tag CRUD operations and tag suggestions.

### Acceptance Criteria
- [ ] Create new tags with names and descriptions
- [ ] Edit tag names and descriptions
- [ ] Merge duplicate tags
- [ ] View questions per tag
- [ ] Tag autocomplete in question forms
- [ ] Suggest tags for questions (future: ML-based)

### Common Tag Examples
- Variables, Data Types, Operators
- Control Structures, Loops, Conditionals
- Functions, Methods, Parameters
- OOP, Classes, Inheritance
- Error Handling, Exceptions
- SQL Queries, Joins, Aggregation

### API Requirements
- GET /admin/tags - List all tags
- POST /admin/tags - Create tag
- PUT /admin/tags/{id} - Update tag
- DELETE /admin/tags/{id} - Delete tag
- POST /admin/tags/merge - Merge tags
- GET /admin/tags/{id}/questions - Questions with tag

### Unit Tests Required
1. Test tag creation
2. Test tag name uniqueness
3. Test tag update
4. Test tag deletion (only if no questions)
5. Test tag merge operation
6. Test questions by tag retrieval
7. Test tag autocomplete
8. Test tag description optional

### UI Components
- Tag list table
- Create/edit tag modal
- Merge tags interface
- Questions per tag view
- Tag autocomplete in question forms
- Tag cloud visualization

### Business Rules
- Tag names must be unique (case-insensitive)
- Cannot delete tag with questions (must merge first)
- Tag merge transfers all questions to target tag
- Tags case-insensitive but stored with proper casing
```

#### FR-CONTENT-003: Question Quality Review
**Priority**: P2

**Description**: Content admins review flagged questions for quality issues.

**GitHub Issue Template**:
```markdown
## Feature: Question Quality Review

### Description
Implement a quality review queue for flagged questions, enabling content admins to review and resolve quality issues.

### Acceptance Criteria
- [ ] View queue of flagged questions
- [ ] See student feedback/comments (future)
- [ ] See performance statistics
- [ ] Approve or edit questions
- [ ] Remove from queue when resolved
- [ ] Auto-flag based on criteria

### Auto-Flag Criteria
- Average accuracy < 30% (too hard or ambiguous)
- Average accuracy > 90% (too easy)
- Answer distribution skewed (one distractor never chosen)
- High skip rate (future)
- Student reports (future)

### API Requirements
- GET /admin/questions/flagged - Get flagged questions queue
- POST /admin/questions/{id}/flag - Flag question
- DELETE /admin/questions/{id}/flag - Unflag question
- PUT /admin/questions/{id}/resolve - Mark as resolved

### Unit Tests Required
1. Test flagged questions queue retrieval
2. Test auto-flag for too easy (>90%)
3. Test auto-flag for too hard (<30%)
4. Test manual flag
5. Test unflag operation
6. Test resolve operation
7. Test flag reason tracking

### UI Components
- Flagged questions queue table
- Flag reason column
- Quick edit modal
- Approve/reject buttons
- Resolve confirmation
- Statistics display

### Flag Reasons
- TOO_HARD: <30% accuracy
- TOO_EASY: >90% accuracy
- AMBIGUOUS: Mixed feedback
- POOR_DISTRACTORS: Skewed distribution
- MANUAL: Instructor flagged
```

### Phase 5 Success Criteria
- Question management UI complete
- Bulk import supports 500+ questions
- Tagging system enables concept organization
- Quality review process established
- All Phase 5 unit tests passing (30+ tests)
- Content team can efficiently manage questions
- Documentation complete

### Estimated Testing
**Unit Tests**: ~30 tests minimum
- Question Management: 10 tests
- Bulk Import: 10 tests
- Tagging System: 8 tests
- Quality Review: 7 tests

---

## Phase 6: Performance & Polish

### Objectives
Optimize performance, enhance user experience, implement advanced analytics, and prepare for production scale. This phase focuses on polish and production readiness.

### Dependencies
- All previous phases

### Feature Requirements Included

#### FR-STUDENT-005: Enhanced Progress Dashboard
**Priority**: P0 (Enhancement of existing feature)

**Description**: Enhance existing progress dashboard with advanced visualizations and insights.

**GitHub Issue Template**:
```markdown
## Feature: Enhanced Progress Dashboard

### Description
Enhance the student progress dashboard with comprehensive metrics, visualizations, and insights to motivate students and guide their learning.

### Acceptance Criteria
- [ ] Total questions attempted across all topics
- [ ] Total correct/incorrect counts with trends
- [ ] Overall accuracy percentage with history
- [ ] Per-topic accuracy percentages with progress bars
- [ ] Recent activity chart (last 30 days)
- [ ] Tag-based strength/weakness analysis
- [ ] Progress toward course completion
- [ ] Current streak tracking
- [ ] Achievements and milestones

### Dashboard Sections
1. **Overview**: Total stats, current streak, recent achievement
2. **Topics**: Progress bar per topic with mastery %
3. **Recent Activity**: Last 10 quiz attempts with sparkline
4. **Insights**: Weakest tags, suggested focus areas
5. **Achievements**: Milestones and badges

### API Requirements
- GET /students/dashboard - Get comprehensive dashboard data
- GET /students/progress/trends - Get historical trends
- GET /students/achievements - Get earned achievements

### Unit Tests Required
1. Test dashboard data aggregation
2. Test trend calculation (30 days)
3. Test tag strength/weakness analysis
4. Test progress toward completion calculation
5. Test streak tracking
6. Test achievement unlocking
7. Test with new user (no data)
8. Test with power user (lots of data)

### UI Components
- Key metrics cards
- Topic progress bars with percentages
- Activity line chart (30 days)
- Tag radar chart
- Achievement badges
- Recommendations panel
- Motivational messages

### Achievements Examples
- First Quiz Complete
- 10 Quizzes Complete
- 100% Accuracy on Quiz
- 7-Day Streak
- Topic Master (90% on all questions)
```

#### Performance Optimization
**Priority**: P0

**GitHub Issue Template**:
```markdown
## Feature: Performance Optimization

### Description
Optimize database queries, implement caching, and improve response times for production scale.

### Acceptance Criteria
- [ ] All API endpoints respond in <200ms (95th percentile)
- [ ] Database query optimization with indexes
- [ ] Implement caching for frequently accessed data
- [ ] Pagination for all list endpoints
- [ ] Connection pooling configured
- [ ] Query performance monitoring

### Optimization Tasks
1. Add database indexes on foreign keys
2. Add indexes on frequently queried columns
3. Implement Redis caching for:
   - Topic lists
   - Tag lists
   - User statistics
4. Implement pagination (default 50 items)
5. Optimize N+1 queries
6. Use prepared statements everywhere
7. Connection pool tuning

### Performance Tests Required
1. Load test with 100 concurrent users
2. Test quiz generation time (<100ms)
3. Test dashboard load time (<200ms)
4. Test analytics page load (<500ms)
5. Test bulk import (500 questions <10s)
6. Test query response times
7. Test cache hit rates

### Monitoring
- Add performance logging
- Track slow queries
- Monitor cache hit rates
- Alert on >1s response times
```

#### Mobile Responsive Design
**Priority**: P1

**GitHub Issue Template**:
```markdown
## Feature: Mobile Responsive Design

### Description
Ensure all pages are fully responsive and usable on mobile devices.

### Acceptance Criteria
- [ ] All pages responsive (320px to 1920px)
- [ ] Touch-friendly buttons and controls
- [ ] Readable fonts on small screens
- [ ] No horizontal scrolling
- [ ] Mobile-optimized navigation
- [ ] Progressive Web App (PWA) support

### Pages to Optimize
- Login/Register
- Student Dashboard
- Quiz Taking
- Quiz History
- Instructor Dashboard
- Student Analytics
- Question Management

### Testing Requirements
1. Test on iPhone (Safari)
2. Test on Android (Chrome)
3. Test on iPad
4. Test various screen sizes
5. Test touch interactions
6. Test PWA installation
7. Test offline capabilities (future)

### UI Updates
- Hamburger menu for mobile
- Collapsible sections
- Touch-friendly quiz interface
- Swipe gestures for navigation
- Optimized charts for small screens
```

#### Advanced Analytics
**Priority**: P2

**GitHub Issue Template**:
```markdown
## Feature: Advanced Analytics

### Description
Implement advanced analytics features including predictive modeling and detailed insights.

### Acceptance Criteria
- [ ] Predicted exam scores based on quiz performance
- [ ] Learning velocity tracking
- [ ] Optimal study time recommendations
- [ ] Concept dependency mapping
- [ ] Retention rate predictions

### Analytics Features
1. **Predictive Scores**: Predict test performance
2. **Learning Velocity**: Rate of improvement
3. **Study Recommendations**: Best times/topics to study
4. **Retention Curves**: Memory retention over time
5. **Mastery Predictions**: Time to topic mastery

### API Requirements
- GET /analytics/predictions?user_id={id}
- GET /analytics/learning-velocity?user_id={id}
- GET /analytics/recommendations?user_id={id}

### Algorithms Needed
- Linear regression for score prediction
- Exponential curve fitting for retention
- Time-series analysis for velocity
- Collaborative filtering for recommendations

### Unit Tests Required
1. Test score prediction accuracy
2. Test learning velocity calculation
3. Test recommendation generation
4. Test retention curve fitting
5. Test with insufficient data
```

#### Production Readiness
**Priority**: P0

**GitHub Issue Template**:
```markdown
## Feature: Production Readiness

### Description
Implement security, monitoring, and deployment features for production launch.

### Acceptance Criteria
- [ ] HTTPS/SSL configuration
- [ ] Rate limiting on API endpoints
- [ ] Input validation and sanitization
- [ ] SQL injection prevention
- [ ] XSS protection
- [ ] CSRF tokens
- [ ] Logging and monitoring
- [ ] Error tracking (e.g., Sentry)
- [ ] Health check endpoint
- [ ] Graceful shutdown
- [ ] Database migration system
- [ ] Environment-based configuration

### Security Tasks
1. Implement rate limiting (100 req/min per user)
2. Add CSRF tokens to forms
3. Sanitize all inputs
4. Use parameterized queries (already done)
5. Add security headers
6. Implement proper CORS
7. Password complexity requirements
8. Session timeout (7 days)

### Monitoring Tasks
1. Application logging (info, warn, error)
2. Request/response logging
3. Performance metrics
4. Error tracking integration
5. Health check endpoint (/health)
6. Database connection monitoring

### Deployment Tasks
1. Docker configuration
2. Database migration scripts
3. Environment variables
4. CI/CD pipeline
5. Backup strategy
6. Rollback procedures

### Testing Required
1. Security audit
2. Penetration testing
3. Load testing (500 concurrent users)
4. Disaster recovery test
5. Backup/restore test
```

### Phase 6 Success Criteria
- All pages load in <200ms
- Mobile responsive on all devices
- Advanced analytics provide value
- Production security implemented
- Monitoring and logging active
- All Phase 6 tests passing
- Performance benchmarks met
- Security audit passed
- Documentation complete
- Ready for production deployment

### Estimated Testing
**Unit Tests**: ~20 tests minimum
**Load Tests**: 5 scenarios
**Security Tests**: 10 checks
**E2E Tests**: 15 user flows

---

## Testing Strategy Across All Phases

### Testing Pyramid
- **Unit Tests**: 150+ tests (70%)
- **Integration Tests**: 50+ tests (20%)
- **E2E Tests**: 20+ tests (10%)

### Test Categories
1. **Model Tests**: Validate data structures
2. **Database Tests**: CRUD operations
3. **API Tests**: Endpoint functionality
4. **Algorithm Tests**: FSRS, analytics calculations
5. **Integration Tests**: Multi-component workflows
6. **Performance Tests**: Load and stress testing
7. **Security Tests**: Vulnerability scanning
8. **E2E Tests**: Full user workflows

### Continuous Testing
- All tests run on every commit
- Code coverage minimum: 80%
- Performance regression detection
- Security scanning in CI/CD

---

## Implementation Roadmap

### Timeline Estimates
- **Phase 3**: 3-4 weeks (Instructor features)
- **Phase 4**: 4-5 weeks (Testing & student features)
- **Phase 5**: 3-4 weeks (Content management)
- **Phase 6**: 3-4 weeks (Polish & production)

**Total**: 13-17 weeks for full implementation

### Resource Requirements
- 2 Backend Developers
- 1 Frontend Developer
- 1 QA Engineer
- 1 Product Manager (part-time)

### Parallel Work Opportunities
- Phase 3 & 4 can partially overlap (different developers)
- Phase 5 can start once Phase 3 instructor foundation is ready
- Phase 6 optimization can start once features stabilize

---

## Risk Mitigation

### Technical Risks
1. **FSRS Complexity**: Already mitigated in Phase 2
2. **Performance at Scale**: Phase 6 addresses
3. **Data Privacy**: Security audit in Phase 6

### Process Risks
1. **Scope Creep**: Strict phase boundaries
2. **Testing Debt**: Tests required for each PR
3. **Documentation Lag**: Docs updated with each phase

### Mitigation Strategies
- Feature flags for gradual rollout
- Comprehensive testing at each phase
- Regular code reviews
- Sprint demos for stakeholder feedback

---

## Success Metrics

### Technical Metrics
- Test coverage >80%
- API response times <200ms (95th percentile)
- Zero critical security vulnerabilities
- Database query efficiency (all indexed)

### Product Metrics
- Student engagement (daily active users)
- Quiz completion rates
- Instructor adoption
- Question bank growth (1000+ questions)

### Quality Metrics
- Bug rate <1% of features
- User-reported issues <5/week
- Uptime >99.9%
- Customer satisfaction score >4.5/5

---

## Conclusion

This phased implementation plan breaks down the comprehensive SPECIFICATION.md feature requirements into manageable, testable increments. Each phase builds upon the previous ones, with clear success criteria and testing requirements.

The plan balances delivering value early (Phase 1-2 already complete) with building toward a comprehensive, production-ready platform (Phase 6). Each feature includes specific unit test descriptions to ensure quality and prove functionality.

**Next Steps**:
1. Review and approve this implementation plan
2. Create GitHub issues from templates in this document
3. Prioritize and schedule work for Phase 3
4. Begin implementation following the defined roadmap
