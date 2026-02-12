# GitHub Issues Creation Guide

## Overview

This document provides templates and guidelines for creating GitHub issues from the Feature Requirements defined in SPECIFICATION.md and organized in IMPLEMENTATION_PHASES.md.

## Issue Naming Convention

Use the following format for issue titles:
```
[Phase X] Feature Name - Brief Description
```

Examples:
- `[Phase 3] FR-INSTRUCTOR-001 - Student List & Search`
- `[Phase 4] FR-STUDENT-006 - Timed Test Taking`
- `[Phase 5] FR-CONTENT-001 - Bulk Question Import`

## Standard Issue Template

```markdown
## Feature: [Feature Name]

**Feature Requirement**: [FR-XXX-###]  
**Phase**: [Phase Number]  
**Priority**: [P0/P1/P2]  
**Dependencies**: [List any dependencies on other features]

### Description
[Brief description of the feature and its purpose]

### Acceptance Criteria
- [ ] [Acceptance criterion 1]
- [ ] [Acceptance criterion 2]
- [ ] [Acceptance criterion 3]
- [ ] ...

### API Requirements
- [HTTP Method] [Endpoint] - [Description]
- [HTTP Method] [Endpoint] - [Description]

### Database Changes
- [ ] New tables: [list]
- [ ] Schema modifications: [list]
- [ ] New indexes: [list]

### Unit Tests Required
1. [Test description 1]
2. [Test description 2]
3. [Test description 3]
...

**Minimum Tests**: [Number]

### UI Components (if applicable)
- [Component 1]
- [Component 2]
- [Component 3]

### Technical Notes
[Any technical considerations, algorithms, or implementation notes]

### Definition of Done
- [ ] All acceptance criteria met
- [ ] All unit tests written and passing
- [ ] Code reviewed and approved
- [ ] Documentation updated
- [ ] Manual testing completed
- [ ] No new security vulnerabilities
- [ ] Performance benchmarks met (if applicable)
```

---

## Example Issues by Phase

### Phase 3 Example: Student List & Search

```markdown
## Feature: Student List & Search

**Feature Requirement**: FR-INSTRUCTOR-001  
**Phase**: Phase 3  
**Priority**: P0  
**Dependencies**: 
- Phase 1: User authentication
- Phase 2: Quiz attempt tracking

### Description
Implement an instructor dashboard page that displays all registered students with their key performance metrics. Enable search and filtering to help instructors quickly find and identify struggling students.

### Acceptance Criteria
- [ ] Display list of all students with name, email, quiz attempts, accuracy, last activity
- [ ] Search by name or email (case-insensitive, partial match)
- [ ] Sort by columns: name, email, accuracy, attempts, last activity
- [ ] Filter by cohort/enrollment period
- [ ] Flag struggling students (accuracy < 60%) with red indicator
- [ ] Show status indicators: green (>75%), yellow (60-75%), red (<60%)
- [ ] Pagination for lists >50 students
- [ ] Responsive design for mobile/tablet

### API Requirements
- GET /instructor/students - List all students with metrics
  - Query params: ?page={n}&per_page={n}&sort={field}&order={asc|desc}
- GET /instructor/students/search?q={query} - Search students
- GET /instructor/students?cohort={id} - Filter by cohort

### Database Changes
- [ ] Add index on users.email for search performance
- [ ] Add index on users.created_at for cohort filtering

### Unit Tests Required
1. Test student list retrieval with basic metrics
2. Test metrics calculation (quiz attempts count, accuracy percentage)
3. Test search by name (exact and partial match)
4. Test search by email (exact and partial match)
5. Test sorting by each column (asc and desc)
6. Test filtering by cohort ID
7. Test struggling student flagging logic (accuracy < 60%)
8. Test status indicator logic (green/yellow/red thresholds)
9. Test pagination with various page sizes
10. Test empty state when no students exist
11. Test access control (only instructors can access)
12. Test performance with 500+ students

**Minimum Tests**: 12

### UI Components
- Instructor dashboard page (/instructor/students.html)
- Student table with sortable column headers
- Search input field with debounced search
- Filter dropdown for cohorts
- Status badge component (green/yellow/red)
- Pagination controls
- Empty state message
- Loading spinner

### Technical Notes
- Use debounced search (300ms delay) to reduce API calls
- Cache student metrics for 5 minutes to improve performance (consider materialized view for production)
- Consider using virtual scrolling for very large lists (>1000 students)
- Ensure mobile layout stacks table into cards

### Definition of Done
- [ ] All acceptance criteria met
- [ ] All 12 unit tests written and passing
- [ ] Code reviewed and approved
- [ ] API documentation updated
- [ ] UI screenshots added to PR
- [ ] Manual testing on Chrome, Firefox, Safari
- [ ] Mobile responsive verified
- [ ] No new security vulnerabilities (CodeQL passed)
- [ ] API response time <200ms (tested with 500 students)
```

---

### Phase 4 Example: Timed Test Taking

```markdown
## Feature: Timed Test Taking

**Feature Requirement**: FR-STUDENT-006  
**Phase**: Phase 4  
**Priority**: P1  
**Dependencies**: 
- Phase 1: User authentication, quiz system
- Phase 4: FR-INSTRUCTOR-005 (Test Management)

### Description
Implement formal timed tests that simulate exam conditions with countdown timer, no-pause policy, auto-submission, and immediate results. Students should have a clear, stress-free interface that focuses them on the test content.

### Acceptance Criteria
- [ ] Display list of available tests with title, description, duration, question count
- [ ] Show test instructions and rules before starting
- [ ] Require acceptance of rules (checkbox) before starting test
- [ ] Create test attempt record when test starts
- [ ] Display countdown timer (always visible, never refreshes)
- [ ] Show questions sequentially with navigation controls
- [ ] Save answers automatically as they are selected
- [ ] Warn at 5 minutes remaining (modal notification)
- [ ] Auto-submit when timer reaches 0:00
- [ ] Allow manual submit before timer expires
- [ ] Prevent starting same test twice
- [ ] Display final score and percentile immediately after submission
- [ ] Persist test state in case of page refresh (localStorage backup)

### API Requirements
- GET /tests/available?user_id={id} - List available tests for user
- GET /tests/{id} - Get test details (before starting)
- POST /tests/{id}/start - Start test attempt (creates TestAttempt record)
  - Response: { test_attempt_id, start_time, duration_seconds, questions }
- PUT /test-attempts/{id}/answer - Submit/update answer during test
  - Body: { question_id, answer_id }
- POST /test-attempts/{id}/submit - Submit test (manual or auto)
  - Response: { score, total, percentage, percentile }
- GET /test-attempts/{id}/status - Check test status (for validation)

### Database Changes
- [ ] Add test_attempts table:
  - id, user_id, test_id, start_time, end_time, duration_seconds
  - submitted, score, total, percentage, percentile
- [ ] Add test_answers table:
  - test_attempt_id, question_id, answer_id, answered_at
- [ ] Add index on test_attempts(user_id, test_id) for duplicate prevention
- [ ] Add index on test_answers(test_attempt_id) for retrieval

### Unit Tests Required
1. Test available tests retrieval for user
2. Test test details retrieval
3. Test test start creates TestAttempt record
4. Test test start records start_time
5. Test cannot start same test twice (unique constraint)
6. Test answer submission during test
7. Test answer update during test (can change answer)
8. Test manual submit before timer expires
9. Test auto-submit calculation (server-side time check)
10. Test score calculation (correct/total)
11. Test percentile calculation based on historical attempts
12. Test timer expiry validation on answer submission
13. Test test state retrieval after page refresh
14. Test preventing answer after submission
15. Test preventing answer after timer expiry

**Minimum Tests**: 15

### UI Components
- Available tests list page (/tests.html)
- Test instructions modal with rules
- Acceptance checkbox ("I understand the rules")
- Test taking page (/test-attempt.html)
- Countdown timer component (MM:SS format, always visible)
- Question display with multiple choice options
- Navigation controls (Previous/Next buttons)
- Question navigator (sidebar showing answered/unanswered)
- Submit test button
- 5-minute warning modal
- Auto-submit notification modal
- Results display page
- Progress indicator (Question X of Y)

### Technical Notes

**Timer Implementation**:
- Client-side countdown using setInterval (1 second)
- Store start_time and duration in localStorage
- Recalculate remaining time on page load (handle refresh)
- Server validates all submissions against start_time + duration
- Warning at remaining_time <= 300 seconds

**Auto-submit Logic**:
```javascript
// Client-side
if (remainingSeconds <= 0) {
  autoSubmit();
}

// Server-side validation
elapsed = now - start_time;
if (elapsed > duration_seconds) {
  reject("Test time expired");
}
```

**LocalStorage Backup**:
```javascript
localStorage.setItem('test_attempt', JSON.stringify({
  id: attemptId,
  startTime: startTime,
  duration: duration,
  answers: answersArray
}));
```

**Percentile Calculation**:
```
1. Get all completed attempts for this test
2. Sort by score
3. Count attempts below user's score
4. Percentile = (below + 0.5) / total * 100
```

### Definition of Done
- [ ] All acceptance criteria met
- [ ] All 15 unit tests written and passing
- [ ] Timer accuracy tested (doesn't drift over 60 minutes)
- [ ] Page refresh handling tested
- [ ] Auto-submit tested (wait for timer to expire)
- [ ] Code reviewed and approved
- [ ] Documentation updated
- [ ] Manual testing on Chrome, Firefox, Safari
- [ ] Mobile responsive verified
- [ ] No new security vulnerabilities
- [ ] Server-side time validation prevents cheating
```

---

### Phase 5 Example: Bulk Question Import

```markdown
## Feature: Bulk Question Import

**Feature Requirement**: FR-CONTENT-001  
**Phase**: Phase 5  
**Priority**: P1  
**Dependencies**: 
- Phase 3: Question management foundation
- Phase 5: FR-CONTENT-002 (Tagging system)

### Description
Enable content administrators to import questions in bulk via CSV file upload, with validation, preview, and error reporting. Support imports of 500+ questions efficiently.

### Acceptance Criteria
- [ ] Upload CSV file (drag-and-drop or file picker)
- [ ] Parse CSV and detect column headers
- [ ] Allow column mapping to question fields
- [ ] Validate all rows before import
- [ ] Display preview table with validation status
- [ ] Show validation errors with line numbers
- [ ] Allow user to confirm import after preview
- [ ] Import only valid rows
- [ ] Display summary (success count, error count)
- [ ] Download error report for invalid rows
- [ ] Support 500+ questions per import
- [ ] Handle duplicate question detection

### CSV Format
```csv
Topic,Question,Answer1,Answer2,Answer3,Answer4,CorrectAnswer,Difficulty,Tags
Java,"What is the output of: int x = 5; System.out.println(++x);","5","6","7","Error","6",2,"variables,operators"
Python,"What does len([1,2,3]) return?","1","2","3","4","3",1,"lists,functions"
SQL,"Which clause filters rows?","SELECT","WHERE","FROM","JOIN","WHERE",1,"queries,basics"
```

### Import Flow
1. User uploads CSV file
2. System parses CSV and shows preview (first 10 rows)
3. User maps columns (if headers don't match)
4. System validates all rows
5. User reviews preview table with validation indicators
6. User clicks "Import" to confirm
7. System imports valid rows
8. System displays summary and error report

### API Requirements
- POST /admin/questions/import/upload - Upload CSV file
  - Body: multipart/form-data with file
  - Response: { preview: [...], total: N, valid: N, invalid: N }
- POST /admin/questions/import/validate - Validate mapped data
  - Body: { mapping: {...}, data: [...] }
  - Response: { validations: [...], errors: [...] }
- POST /admin/questions/import/confirm - Execute import
  - Body: { data: [...] }
  - Response: { imported: N, skipped: N, errors: [...] }
- GET /admin/questions/import/errors/{id} - Download error report CSV

### Database Changes
- [ ] Add import_logs table for tracking imports:
  - id, user_id, filename, total, imported, skipped, created_at
- [ ] Add transaction support for bulk insert
- [ ] Add index on questions(topic_id, question_text) for duplicate detection

### Unit Tests Required
1. Test CSV upload and parsing
2. Test column header detection
3. Test column mapping (custom headers)
4. Test validation: required fields (topic, question, answers)
5. Test validation: correct answer exists in answer options
6. Test validation: difficulty range (1-7)
7. Test validation: minimum 2 answers required
8. Test validation: exactly 1 correct answer
9. Test duplicate question detection (same topic + question text)
10. Test topic creation if doesn't exist
11. Test tag parsing (comma-separated)
12. Test bulk insert performance (500 questions <10 seconds)
13. Test partial import (some valid, some invalid)
14. Test transaction rollback on error
15. Test error report generation
16. Test import log creation

**Minimum Tests**: 16

### UI Components
- Import page (/admin/import.html)
- File upload dropzone (drag-and-drop)
- Column mapping interface (dropdowns)
- Preview table (first 10 rows)
- Validation indicator column (✓ or ✗)
- Error list panel (expandable)
- Import button (disabled until validated)
- Progress bar during import
- Success/error summary modal
- Download error report button

### Technical Notes

**Validation Rules**:
```javascript
{
  topic: { required: true },
  question: { required: true, maxLength: 1000 },
  answer1: { required: true },
  answer2: { required: true },
  answer3: { optional: true },
  answer4: { optional: true },
  correctAnswer: { required: true, mustMatchOneAnswer: true },
  difficulty: { required: true, min: 1, max: 7, type: 'integer' },
  tags: { optional: true, format: 'comma-separated' }
}
```

**Duplicate Detection**:
- Check if question with same topic_id and question_text exists
- Option to skip or update existing questions
- Flag duplicates in preview

**Performance Optimization**:
- Batch insert in chunks of 100 questions
- Use prepared statements
- Single transaction for entire import
- Show progress percentage

**Error Report Format**:
```csv
Line,Error,Topic,Question
2,Correct answer 'X' not in answer options,Java,"What is ..."
5,Difficulty must be 1-7,Python,"..."
```

### Definition of Done
- [ ] All acceptance criteria met
- [ ] All 16 unit tests written and passing
- [ ] Import of 500 questions <10 seconds (tested)
- [ ] CSV parsing handles edge cases (quotes, commas in text)
- [ ] Code reviewed and approved
- [ ] Documentation updated (CSV format guide)
- [ ] Manual testing with various CSV files
- [ ] Error handling tested (malformed CSV, empty file)
- [ ] No new security vulnerabilities
- [ ] File upload size limit enforced (10 MB)
```

---

## Issue Labels

Use the following labels to categorize issues:

### Phase Labels
- `phase-3-instructor` - Phase 3 features
- `phase-4-student-testing` - Phase 4 features
- `phase-5-content` - Phase 5 features
- `phase-6-polish` - Phase 6 features

### Priority Labels
- `priority-p0` - Critical (must have)
- `priority-p1` - Important (should have)
- `priority-p2` - Nice to have (could have)

### Type Labels
- `feature` - New feature
- `enhancement` - Enhancement to existing feature
- `bug` - Bug fix
- `documentation` - Documentation only
- `testing` - Testing infrastructure

### Status Labels
- `status-backlog` - Not started
- `status-in-progress` - Being worked on
- `status-review` - In code review
- `status-testing` - In testing
- `status-done` - Completed

### Component Labels
- `backend` - Backend/API work
- `frontend` - Frontend/UI work
- `database` - Database changes
- `security` - Security-related
- `performance` - Performance-related

---

## Creating Issues in Bulk

### Using GitHub CLI

```bash
# Install GitHub CLI if not already installed
# https://cli.github.com/

# Login
gh auth login

# Create issue from template
gh issue create \
  --title "[Phase 3] FR-INSTRUCTOR-001 - Student List & Search" \
  --body-file issues/instructor-001.md \
  --label "phase-3-instructor,priority-p0,feature,backend,frontend" \
  --assignee "@me"
```

### Using GitHub API

```bash
# Set your GitHub token
TOKEN="your_github_token"
REPO="ZipCodeCore/zip-athena"

# Create issue
curl -X POST \
  -H "Authorization: token $TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$REPO/issues \
  -d '{
    "title": "[Phase 3] FR-INSTRUCTOR-001 - Student List & Search",
    "body": "...",
    "labels": ["phase-3-instructor", "priority-p0", "feature"]
  }'
```

---

## Issue Workflow

### 1. Creation
- Create issue from template
- Add appropriate labels
- Link to Feature Requirement in SPECIFICATION.md
- Reference in IMPLEMENTATION_PHASES.md

### 2. Refinement
- Team reviews issue
- Clarifies acceptance criteria
- Estimates effort (S/M/L/XL)
- Assigns to sprint/milestone

### 3. Development
- Assign to developer
- Move to "In Progress"
- Create feature branch: `feature/FR-XXX-###-short-name`
- Implement feature
- Write unit tests

### 4. Review
- Create pull request
- Reference issue: "Closes #123"
- Code review
- Address feedback

### 5. Testing
- All unit tests pass
- Manual testing completed
- QA verification
- Performance check

### 6. Completion
- Merge to main
- Issue automatically closed
- Update documentation
- Deploy to staging

---

## Sprint Planning

### Recommended Sprint Structure

**Sprint Duration**: 2 weeks

**Capacity Planning**:
- Small (S): 1-2 days
- Medium (M): 3-4 days
- Large (L): 5-7 days
- Extra Large (XL): 8-10 days

**Sprint Velocity**:
- Junior dev: ~3-4 medium issues per sprint
- Mid-level dev: ~4-5 medium issues per sprint
- Senior dev: ~5-6 medium issues per sprint

### Phase 3 Sprint Example

**Sprint 1 (Instructor Foundation)**:
1. FR-INSTRUCTOR-001: Student List & Search (M)
2. FR-INSTRUCTOR-007: Topic Management (S)
3. Database optimizations for analytics (S)

**Sprint 2 (Individual Analytics)**:
1. FR-INSTRUCTOR-002: Individual Student Analytics (L)
2. PDF export functionality (M)

**Sprint 3 (Cohort Analytics)**:
1. FR-INSTRUCTOR-003: Cohort Analytics Dashboard (L)
2. Chart library integration (S)

---

## Milestone Planning

### Phase 3 Milestone
- **Name**: Instructor Analytics & Management
- **Target Date**: [Set based on team capacity]
- **Issues**: All FR-INSTRUCTOR-* features
- **Success Criteria**:
  - All Phase 3 tests passing (28+)
  - Instructor can view and analyze students
  - Documentation complete

### Phase 4 Milestone
- **Name**: Advanced Student Features & Testing
- **Target Date**: [+4-5 weeks from Phase 3]
- **Issues**: FR-STUDENT-004, FR-STUDENT-006, FR-STUDENT-007, FR-INSTRUCTOR-005, FR-INSTRUCTOR-006
- **Success Criteria**:
  - Timed tests functional
  - Quiz history accessible
  - Tag-based practice available
  - All Phase 4 tests passing (35+)

### Phase 5 Milestone
- **Name**: Content Management & Quality
- **Target Date**: [+3-4 weeks from Phase 4]
- **Issues**: FR-INSTRUCTOR-004, FR-CONTENT-001, FR-CONTENT-002, FR-CONTENT-003
- **Success Criteria**:
  - Bulk import supports 500+ questions
  - Quality review process active
  - All Phase 5 tests passing (30+)

### Phase 6 Milestone
- **Name**: Performance & Polish
- **Target Date**: [+3-4 weeks from Phase 5]
- **Issues**: Performance, mobile, analytics, production readiness
- **Success Criteria**:
  - Performance benchmarks met
  - Mobile responsive
  - Production ready
  - Security audit passed

---

## Templates Directory Structure

Store issue templates in `.github/ISSUE_TEMPLATE/`:

```
.github/
└── ISSUE_TEMPLATE/
    ├── feature_request.md
    ├── bug_report.md
    ├── phase_3_instructor.md
    ├── phase_4_student.md
    ├── phase_5_content.md
    └── phase_6_polish.md
```

---

## Summary

This guide provides everything needed to create well-structured GitHub issues from the Feature Requirements. Each issue includes:

1. ✅ Clear description and context
2. ✅ Specific acceptance criteria
3. ✅ API requirements
4. ✅ Database changes
5. ✅ Unit test descriptions
6. ✅ UI components
7. ✅ Technical notes
8. ✅ Definition of Done

By following this template structure, all issues will be:
- **Actionable**: Developers know exactly what to build
- **Testable**: Clear unit test requirements
- **Measurable**: Explicit success criteria
- **Traceable**: Links to original requirements

**Next Steps**:
1. Review this template guide
2. Create issues for Phase 3 (next phase to implement)
3. Organize issues into sprints
4. Begin development following the defined workflow
