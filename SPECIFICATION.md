# Athena - Technical Knowledge Assessment Platform
## Project Specification Document

**Version:** 1.0  
**Date:** January 27, 2026  
**Organization:** ZipCodeWilmington  
**Document Type:** Product Specification & Functional Requirements  

---

## Executive Summary

Athena is a SaaS web application designed to help software development students at ZipCodeWilmington achieve technical mastery through adaptive, intelligent quiz-based learning. The platform leverages spaced repetition algorithms (FSRS - Free Spaced Repetition Scheduler) to optimize knowledge retention and track student progress across multiple programming languages and technical domains including Java, Python, SQL, and Agile methodologies.

---

# Part 1: Project Charter

## 1.1 Project Vision

To create an intelligent, adaptive learning platform that accelerates software development students' mastery of technical concepts through scientifically-proven spaced repetition techniques, detailed performance analytics, and personalized learning pathways.

## 1.2 Primary Goals

1. **Accelerate Learning**: Reduce time-to-mastery for core programming concepts by 30% through adaptive spaced repetition
2. **Improve Retention**: Increase long-term knowledge retention rates by 50% compared to traditional study methods
3. **Identify Gaps**: Provide real-time insights into student knowledge gaps to enable targeted intervention
4. **Scale Education**: Support unlimited concurrent students with minimal instructor overhead
5. **Data-Driven Instruction**: Enable instructors to make evidence-based curriculum adjustments based on aggregate student performance

## 1.3 Key Success Metrics (KPIs)

### Student Success Metrics
- **Mastery Rate**: Percentage of questions answered correctly (Target: 80%+)
- **Retention Rate**: Percentage of previously mastered concepts retained after 30 days (Target: 75%+)
- **Quiz Completion Rate**: Percentage of assigned quizzes completed (Target: 90%+)
- **Average Time to Mastery**: Days required to achieve 80% accuracy on a topic (Target: <14 days)

### Engagement Metrics
- **Daily Active Users (DAU)**: Student logins per day (Target: 70% of enrolled students)
- **Average Quiz Attempts per Week**: Quiz sessions per student (Target: 5+)
- **Session Duration**: Average time spent per quiz session (Target: 15-20 minutes)

### Platform Performance Metrics
- **System Uptime**: Platform availability (Target: 99.5%+)
- **Response Time**: API response time for quiz generation (Target: <500ms)
- **Data Accuracy**: Correctness of spaced repetition scheduling (Target: 95%+)

### Instructor Efficiency Metrics
- **Time to Identify Struggling Students**: Hours to detect at-risk students (Target: <24 hours)
- **Curriculum Adjustment Speed**: Days to implement topic-based interventions (Target: <7 days)

## 1.4 Core User Roles

### 1. **Student User** (Primary Role)
- **Description**: Software development students enrolled in ZipCodeWilmington courses
- **Primary Goals**: Master programming concepts, track personal progress, identify weak areas
- **Access Level**: Personal dashboard, quiz attempts, test attempts, progress reports
- **Success Criteria**: Achieve 80%+ mastery rate across all assigned topics

### 2. **Instructor/Admin** (Secondary Role)
- **Description**: Teaching staff responsible for course delivery and student success
- **Primary Goals**: Monitor student progress, identify struggling students, adjust curriculum
- **Access Level**: All student data, analytics dashboards, question bank management, test creation
- **Success Criteria**: 85%+ of students achieve course completion standards

### 3. **Content Administrator** (Supporting Role)
- **Description**: Subject matter experts who manage the question database
- **Primary Goals**: Create, update, and maintain high-quality quiz questions and answers
- **Access Level**: Question bank CRUD operations, topic management, difficulty calibration
- **Success Criteria**: Maintain 500+ quality questions per major topic area

### 4. **System Administrator** (Supporting Role)
- **Description**: Technical staff responsible for platform operations
- **Primary Goals**: Ensure system uptime, manage user accounts, maintain data integrity
- **Access Level**: System configuration, user management, database maintenance, logs
- **Success Criteria**: Maintain 99.5%+ uptime and <500ms response times

### 5. **Course Director** (Strategic Role)
- **Description**: Leadership responsible for overall curriculum effectiveness
- **Primary Goals**: Evaluate program effectiveness, make strategic curriculum decisions
- **Access Level**: Aggregate analytics, cohort comparisons, longitudinal trend data
- **Success Criteria**: Demonstrate measurable improvement in student outcomes year-over-year

---

# Part 2: Project Specification

## 2.1 Project Summary

Athena is a three-tier web application consisting of:
- **Frontend**: JavaScript/JSON-based single-page application (SPA)
- **Backend**: RESTful API server (Go-based) handling business logic
- **Database**: Relational database (PostgreSQL/MySQL) storing questions, user data, and attempt history

The platform implements an adaptive learning system using the FSRS (Free Spaced Repetition Scheduler) algorithm to intelligently schedule question reviews based on individual student performance patterns.

### Core Value Proposition
Unlike generic quiz platforms (Kahoot!, Quizlet), Athena is purpose-built for technical software development education with:
- Deep tag-based categorization for identifying conceptual gaps
- Difficulty-progressive question sequencing
- Integration with ZipCodeWilmington's curriculum structure
- Comprehensive instructor analytics for intervention

## 2.2 User Personas

### Persona 1: Alex - The Determined Student
**Demographics**
- Age: 24
- Background: Career changer, previously in retail management
- Technical Level: Beginner to intermediate

**Goals**
- Master Java fundamentals within 12 weeks
- Build confidence in technical interviews
- Identify and fix knowledge gaps quickly

**Pain Points**
- Struggles to know what to study when
- Forgets concepts learned weeks ago
- Doesn't know if ready for assessments

**How Athena Helps**
- Automatic scheduling of reviews prevents forgetting
- Clear progress metrics show readiness
- Tag-based insights reveal conceptual gaps (e.g., "control structures")

**Usage Pattern**
- Completes 1-2 quiz sessions daily (15 minutes each)
- Reviews progress dashboard before tests
- Uses difficulty ratings to gauge mastery

---

### Persona 2: Jordan - The Struggling Learner
**Demographics**
- Age: 19
- Background: Recent high school graduate, first exposure to coding
- Technical Level: Novice

**Goals**
- Pass course requirements to complete bootcamp
- Build foundational understanding before advanced topics
- Get personalized help when stuck

**Pain Points**
- Overwhelmed by volume of material
- Falls behind quickly when concepts aren't clear
- Doesn't know how to ask for help effectively

**How Athena Helps**
- Spaced repetition focuses on weak areas automatically
- Difficulty progression prevents overwhelming complexity
- Instructor dashboard flags Jordan for intervention

**Usage Pattern**
- Completes required weekly quizzes
- Struggles with difficulty 4+ questions
- Benefits from instructor outreach based on analytics

---

### Persona 3: Taylor - The Advanced Student
**Demographics**
- Age: 28
- Background: BS in Computer Science, seeking specialized training
- Technical Level: Intermediate to advanced

**Goals**
- Reinforce fundamentals while learning advanced concepts
- Maintain knowledge across multiple languages (Java, Python, SQL)
- Prepare efficiently for technical interviews

**Pain Points**
- Wastes time on concepts already mastered
- Needs efficient review to prevent degradation
- Wants challenging material

**How Athena Helps**
- Algorithm quickly identifies mastered concepts, reduces review frequency
- Difficulty levels allow focus on complex questions
- Multi-topic support maintains breadth of knowledge

**Usage Pattern**
- Completes quizzes quickly with high accuracy
- Focuses on difficulty 5-7 questions
- Uses platform for pre-interview refreshers

---

### Persona 4: Dr. Morgan - The Instructor
**Demographics**
- Age: 35
- Background: 10 years industry experience + 3 years teaching
- Role: Lead Instructor

**Goals**
- Ensure 85%+ of cohort achieves mastery
- Identify struggling students within first 2 weeks
- Adjust curriculum based on common gaps

**Pain Points**
- Can't monitor 30+ students individually
- Discovers struggles too late to intervene
- Unclear which topics need more classroom time

**How Athena Helps**
- Real-time dashboard shows struggling students
- Tag analytics reveal common misconceptions
- Historical data informs curriculum improvements

**Usage Pattern**
- Reviews cohort dashboard 2-3x weekly
- Drills into individual student records when alerted
- Analyzes test results to adjust upcoming lessons

---

### Persona 5: Sam - The Content Administrator
**Demographics**
- Age: 42
- Background: Senior software engineer, curriculum SME
- Role: Question Bank Manager

**Goals**
- Maintain 1000+ high-quality questions across topics
- Ensure questions align with curriculum objectives
- Calibrate difficulty ratings accurately

**Pain Points**
- Hard to know if questions are too easy/hard
- Difficult to maintain consistency across topics
- Time-consuming to create quality distractors (wrong answers)

**How Athena Helps**
- Question statistics show average accuracy rates
- Tag system enables consistency checking
- Template system speeds question creation

**Usage Pattern**
- Reviews question performance metrics monthly
- Creates 10-20 new questions per topic quarterly
- Adjusts difficulty ratings based on student performance data

---

## 2.3 Key Features

### Student-Facing Features

#### F1: Adaptive Quiz System
- **Description**: Intelligent quiz generation using FSRS algorithm
- **Details**: 
  - Generates personalized quizzes based on past performance
  - Balances new questions with review of previously seen material
  - Prioritizes incorrectly answered questions for earlier review
  - Schedules correctly answered questions with increasing intervals
- **User Value**: Optimizes study time by focusing on weak areas

#### F2: Progress Dashboard
- **Description**: Visual representation of learning progress
- **Details**:
  - Overall mastery percentage across all topics
  - Per-topic progress bars
  - Recent quiz attempt history
  - Strength/weakness indicators by tag
  - Streak tracking for motivation
- **User Value**: Clear understanding of current standing and readiness

#### F3: Difficulty-Progressive Learning
- **Description**: Questions ranked 1-7 by difficulty
- **Details**:
  - Start with easier questions (difficulty 1-3)
  - Unlock harder questions as mastery increases
  - Filter quiz generation by difficulty range
  - Visual difficulty indicators
- **User Value**: Builds confidence while preventing overwhelming complexity

#### F4: Tag-Based Insights
- **Description**: Conceptual gap identification through question tags
- **Details**:
  - Questions tagged with concepts (e.g., "control structures", "loops", "recursion")
  - Dashboard shows accuracy by tag
  - Suggested focus areas based on low-performing tags
  - Cross-topic tag aggregation
- **User Value**: Identifies specific concepts needing additional study

#### F5: Quiz Attempt Review
- **Description**: Detailed review of completed quizzes
- **Details**:
  - Review all questions from past attempts
  - See correct/incorrect answers marked clearly
  - Explanations for correct answers (future enhancement)
  - Time spent per question tracking
- **User Value**: Learn from mistakes and reinforce correct knowledge

#### F6: Timed Test Mode
- **Description**: Formal assessments with fixed question sets
- **Details**:
  - Countdown timer during test
  - Fixed question set (same questions for all students in cohort)
  - Single attempt per test (no retakes without instructor approval)
  - Higher weight in overall scoring
  - Simulates real assessment conditions
- **User Value**: Prepares students for formal evaluations and interviews

### Instructor/Admin Features

#### F7: Student Progress Analytics
- **Description**: Comprehensive view of individual student performance
- **Details**:
  - Mastery percentage by topic
  - Quiz completion rates
  - Time-to-mastery trends
  - Comparison to cohort averages
  - Historical performance graphs
- **User Value**: Data-driven student intervention

#### F8: Cohort Analytics Dashboard
- **Description**: Aggregate view of class performance
- **Details**:
  - Class-wide mastery rates by topic
  - Struggling student identification (red flags)
  - Common mistake patterns (tag-based)
  - Quiz completion rates
  - Test score distributions
- **User Value**: Inform curriculum adjustments and classroom focus

#### F9: Question Bank Management
- **Description**: CRUD operations for quiz questions
- **Details**:
  - Create questions with multiple-choice answers
  - Assign difficulty ratings (1-7)
  - Tag questions with multiple concepts
  - Associate questions with topics
  - Activate/deactivate questions without deletion
- **User Value**: Maintain high-quality, relevant question content

#### F10: Test Creation & Management
- **Description**: Create and administer formal tests
- **Details**:
  - Select specific questions for test
  - Set time limits (e.g., 60 minutes)
  - Schedule test availability windows
  - Assign tests to specific cohorts
  - Manual grading override capabilities
- **User Value**: Standardized assessment across cohorts

#### F11: Question Performance Analytics
- **Description**: Statistical analysis of question effectiveness
- **Details**:
  - Average accuracy rate per question
  - Identify questions that are too easy/hard
  - Flag potentially ambiguous questions
  - Distribution of student answers (distractor analysis)
- **User Value**: Continuous improvement of question quality

### System Features

#### F12: Spaced Repetition Engine (FSRS)
- **Description**: Implementation of Free Spaced Repetition Scheduler algorithm
- **Details**:
  - Tracks individual question attempt history per student
  - Calculates optimal review intervals based on:
    - Number of correct/incorrect attempts
    - Time since last attempt
    - Question difficulty
    - Student's historical performance on similar questions
  - Dynamically adjusts scheduling based on performance
  - Prevents "cramming" by enforcing minimum intervals
- **User Value**: Scientifically-optimized learning schedule

#### F13: Multi-Topic Support
- **Description**: Support for multiple programming languages and concepts
- **Details**:
  - Core topics: Java, Python, SQL, Agile, Data Structures, Algorithms
  - Separate question banks per topic
  - Cross-topic quiz generation option
  - Topic-specific progress tracking
- **User Value**: Comprehensive coverage of bootcamp curriculum

#### F14: Authentication & Authorization
- **Description**: Secure user management and role-based access
- **Details**:
  - Student accounts with email/password
  - Instructor accounts with elevated privileges
  - Session management with JWT tokens
  - Role-based API endpoint protection
- **User Value**: Secure, personalized experience

---

# Part 3: Functional Requirements Document (FRD)

## 3.1 Spaced Repetition System - FSRS Algorithm

### 3.1.1 Algorithm Overview

The Free Spaced Repetition Scheduler (FSRS) is a modern, optimized algorithm based on the following principles:

1. **Stability**: How well information is retained in memory
2. **Difficulty**: Inherent difficulty of the material (different from question difficulty rating)
3. **Retrievability**: Current probability of successful recall

### 3.1.2 Implementation Requirements

#### FR-SRS-001: Question Attempt Tracking
**Priority**: P0 (Critical)  
**Description**: System must track every question attempt with timestamp and result

**Acceptance Criteria**:
- Record question ID, user ID, timestamp, correct/incorrect for every attempt
- Store time taken to answer each question
- Support batch recording for quiz completion
- Maintain complete historical record (no deletions)

**Data Model**:
```
QuestionAttempt {
  ID: string
  UserID: string
  QuestionID: string (format: topicID.questionID, e.g., "3.17")
  IsCorrect: boolean
  AttemptDateTime: timestamp
  TimeSpentSeconds: int
  QuizAttemptID: string (optional, links to quiz session)
}
```

#### FR-SRS-002: FSRS Parameter Calculation
**Priority**: P0 (Critical)  
**Description**: Calculate and maintain FSRS parameters for each user-question pair

**Acceptance Criteria**:
- Calculate stability (S) after each attempt
- Calculate difficulty (D) based on attempt history
- Update retrievability (R) based on time elapsed
- Store parameters for efficient quiz generation

**Algorithm Parameters**:
- **Initial Stability**: S = 1 day (configurable)
- **Stability Increase on Success**: S_new = S_old * success_multiplier
- **Stability Decrease on Failure**: S_new = S_old * failure_multiplier
- **Retrievability Decay**: R = e^(-t/S) where t = time since last attempt

**Data Model**:
```
UserQuestionState {
  UserID: string
  QuestionID: string
  Stability: float (days)
  Difficulty: float (0-10 scale)
  LastRetrievability: float (0-1 scale)
  LastAttemptDateTime: timestamp
  TotalAttempts: int
  CorrectAttempts: int
  NextScheduledReview: timestamp
}
```

#### FR-SRS-003: Quiz Generation Algorithm
**Priority**: P0 (Critical)  
**Description**: Generate optimal quiz question sets using FSRS data

**Acceptance Criteria**:
- Select questions where R < threshold (default 0.7)
- Balance between new questions and reviews
- Respect difficulty progression
- Limit quiz size (default: 10-20 questions)
- Avoid repeating same question in single session

**Quiz Generation Logic**:
1. Query all UserQuestionState records for student
2. Filter questions due for review (R < 0.7 or NextScheduledReview < now)
3. Sort by priority: (1 - R) * difficulty_weight
4. Select top N questions within student's difficulty range
5. Add new questions (never attempted) to reach target quiz size
6. Randomize order within difficulty tiers

#### FR-SRS-004: Review Scheduling
**Priority**: P0 (Critical)  
**Description**: Calculate optimal next review date after each attempt

**Acceptance Criteria**:
- Update NextScheduledReview immediately after attempt
- Successful attempts → increase interval (e.g., 1d → 3d → 7d → 14d → 30d)
- Failed attempts → short interval (e.g., reset to 1 day)
- Never schedule same question multiple times in one day

**Scheduling Formula**:
```
if correct:
  new_interval = current_stability * 2.5 * ease_factor
else:
  new_interval = max(1 day, current_stability * 0.3)

NextScheduledReview = now + new_interval
```

### 3.1.3 FSRS Configuration Settings

| Parameter | Default Value | Description |
|-----------|---------------|-------------|
| Initial Stability | 1 day | Interval for first review |
| Success Multiplier | 2.5 | Stability growth on correct answer |
| Failure Multiplier | 0.3 | Stability reduction on incorrect answer |
| Retrievability Threshold | 0.7 | Trigger point for including in quiz |
| Min Interval | 1 day | Minimum time between reviews |
| Max Interval | 180 days | Maximum time between reviews |
| New Question Ratio | 0.3 | Proportion of new vs. review questions |

---

## 3.2 User Role-Based Features

### 3.2.1 Student Role Features

#### FR-STUDENT-001: User Registration & Login
**Priority**: P0  
**Description**: Students can create accounts and authenticate

**Acceptance Criteria**:
- Register with email, password, full name
- Password requirements: 8+ characters, mix of letters and numbers
- Email verification required (future enhancement)
- Session persists for 7 days
- Password reset functionality

#### FR-STUDENT-002: Topic Selection
**Priority**: P0  
**Description**: Students can view available topics and select for practice

**Acceptance Criteria**:
- Display all active topics (Java, Python, SQL, etc.)
- Show number of questions per topic
- Show student's mastery percentage per topic
- Allow selection of single or multiple topics for quiz

#### FR-STUDENT-003: Adaptive Quiz Taking
**Priority**: P0  
**Description**: Students take quizzes generated by FSRS algorithm

**Acceptance Criteria**:
- Request quiz generation with topic selection
- Display questions one at a time
- Show multiple choice answers (4 options typical)
- Immediate feedback on correct/incorrect
- Display running score during quiz
- Save progress if quiz interrupted
- Calculate total points earned at completion

**Quiz Flow**:
1. Student clicks "Start Quiz" with topic(s) selected
2. System generates quiz using FSRS algorithm
3. Display first question with 4 multiple-choice answers
4. Student selects answer
5. System immediately shows correct/incorrect
6. Student clicks "Next" to proceed
7. Repeat until all questions answered
8. Display final score and statistics

#### FR-STUDENT-004: Quiz Attempt History
**Priority**: P1  
**Description**: Students can review past quiz attempts

**Acceptance Criteria**:
- List all past quiz attempts with date, topic, score
- Sort by date (newest first)
- Filter by topic
- Drill into specific attempt to see individual questions
- Show correct answers for incorrect questions

#### FR-STUDENT-005: Progress Dashboard
**Priority**: P0  
**Description**: Students view comprehensive progress metrics

**Acceptance Criteria**:
- Total questions attempted across all topics
- Total correct/incorrect counts
- Overall accuracy percentage
- Per-topic accuracy percentages
- Recent activity chart (last 30 days)
- Tag-based strength/weakness analysis
- Progress toward course completion

**Dashboard Sections**:
1. **Overview**: Total stats, current streak
2. **Topics**: Progress bar per topic
3. **Recent Activity**: Last 10 quiz attempts
4. **Insights**: Weakest tags, suggested focus areas
5. **Achievements**: Milestones and badges (future)

#### FR-STUDENT-006: Timed Test Taking
**Priority**: P1  
**Description**: Students complete formal timed assessments

**Acceptance Criteria**:
- View available tests with descriptions
- See test duration and question count
- Accept test rules before starting (no pause, single attempt)
- Display countdown timer during test
- Auto-submit when timer expires
- Cannot review questions mid-test
- Receive final score immediately after submission

**Test Flow**:
1. Student clicks "Start Test" (must read rules)
2. System starts timer (e.g., 60 minutes)
3. Display questions sequentially (can navigate back/forward)
4. Timer visible on all pages
5. Warning at 5 minutes remaining
6. Auto-submit at 0:00 or manual submit
7. Display final score and percentile

#### FR-STUDENT-007: Tag-Based Review
**Priority**: P2  
**Description**: Students can practice specific concepts via tags

**Acceptance Criteria**:
- View list of all tags with question counts
- Filter quizzes to specific tag(s)
- See accuracy rate per tag
- Generate tag-specific quiz (overrides normal FSRS)

### 3.2.2 Instructor/Admin Role Features

#### FR-INSTRUCTOR-001: Student List & Search
**Priority**: P0  
**Description**: Instructors can view and search all students

**Acceptance Criteria**:
- Display list of all students with key metrics
- Search by name or email
- Sort by various metrics (accuracy, completion rate, etc.)
- Filter by cohort/enrollment period
- Flag struggling students (accuracy < 60%)

**Student List Columns**:
- Name
- Email
- Total Quiz Attempts
- Overall Accuracy %
- Last Activity Date
- Status Indicator (green/yellow/red)

#### FR-INSTRUCTOR-002: Individual Student Analytics
**Priority**: P0  
**Description**: Instructors can drill into individual student performance

**Acceptance Criteria**:
- View student's progress dashboard
- See all quiz attempt history
- Identify weak topic areas
- Identify weak conceptual areas (tags)
- View test results
- Compare to cohort averages
- Export student report as PDF

**Analytics Views**:
1. **Summary**: Overall stats vs. cohort
2. **Topics**: Mastery by topic with trends
3. **Tags**: Conceptual strengths/weaknesses
4. **History**: Timeline of all attempts
5. **Tests**: Formal assessment results

#### FR-INSTRUCTOR-003: Cohort Analytics Dashboard
**Priority**: P0  
**Description**: Instructors view aggregate class performance

**Acceptance Criteria**:
- Select cohort/class period
- View aggregate statistics (avg accuracy, completion rate)
- Topic-by-topic class performance
- Identify common struggling concepts (tags)
- See distribution charts (histograms of scores)
- Export cohort report

**Dashboard Sections**:
1. **Overview**: Class averages, completion rates
2. **Topics**: Avg accuracy per topic, struggling topics
3. **Students**: Leaderboard and struggling student list
4. **Trends**: Week-over-week progress
5. **Questions**: Hardest questions for class

#### FR-INSTRUCTOR-004: Question Management
**Priority**: P0  
**Description**: Instructors can create, edit, and manage questions

**Acceptance Criteria**:
- Create new question with text, answers, correct answer
- Assign topic, difficulty (1-7), tags
- Edit existing questions
- Deactivate questions (don't delete)
- Bulk import questions (CSV)
- Preview questions as students see them

**Question Form Fields**:
- Question Text (rich text editor for code formatting)
- Topic Dropdown
- Difficulty (1-7 slider)
- Tags (multi-select)
- Answer Option 1 (text + correct checkbox)
- Answer Option 2 (text + correct checkbox)
- Answer Option 3 (text + correct checkbox)
- Answer Option 4 (text + correct checkbox)
- Status (active/inactive)

#### FR-INSTRUCTOR-005: Test Management
**Priority**: P1  
**Description**: Instructors can create and manage formal tests

**Acceptance Criteria**:
- Create new test with title, description
- Add questions to test (search and select)
- Set time limit (minutes)
- Set availability window (start/end dates)
- Assign to specific cohort(s)
- View test results by student
- Export test results

**Test Creation Flow**:
1. Set test metadata (title, description, time limit)
2. Search/filter question bank
3. Add questions to test (drag-and-drop or click)
4. Preview test as student would see
5. Set availability window
6. Assign to cohort(s)
7. Publish

#### FR-INSTRUCTOR-006: Question Performance Analytics
**Priority**: P1  
**Description**: Instructors can analyze question effectiveness

**Acceptance Criteria**:
- View list of all questions with avg accuracy
- Identify questions that are too easy (>90% correct)
- Identify questions that are too hard (<30% correct)
- See answer distribution (which distractors chosen)
- Flag potentially ambiguous questions
- Adjust difficulty ratings based on data

**Question Analytics View**:
- Question ID and text
- Topic and Tags
- Current Difficulty Rating
- Number of Attempts
- Average Accuracy %
- Recommended Difficulty (calculated)
- Answer Distribution (histogram)
- Flag for Review (if outlier)

#### FR-INSTRUCTOR-007: Topic Management
**Priority**: P2  
**Description**: Instructors can manage topic categories

**Acceptance Criteria**:
- Create new topics
- Edit topic names and descriptions
- Deactivate topics
- Reorder topic display
- Set topic prerequisites (future)

### 3.2.3 Content Administrator Features

#### FR-CONTENT-001: Bulk Question Import
**Priority**: P1  
**Description**: Content admins can import questions in bulk

**Acceptance Criteria**:
- Upload CSV file with questions
- Map CSV columns to question fields
- Validate data before import
- Preview questions before finalizing
- Report errors for invalid rows
- Support 500+ questions per import

**CSV Format**:
```
Topic,Question,Answer1,Answer2,Answer3,Answer4,CorrectAnswer,Difficulty,Tags
Java,"What is the output of: int x = 5; System.out.println(++x);","5","6","7","Error","6",2,"variables,operators"
```

#### FR-CONTENT-002: Question Tagging System
**Priority**: P0  
**Description**: Content admins manage tag taxonomy

**Acceptance Criteria**:
- Create new tags with descriptions
- Edit tag names
- Merge duplicate tags
- View questions per tag
- Suggest tags for questions (future: ML-based)

**Common Tag Examples**:
- Variables, Data Types, Operators
- Control Structures, Loops, Conditionals
- Functions, Methods, Parameters
- OOP, Classes, Inheritance
- Error Handling, Exceptions
- SQL Queries, Joins, Aggregation

#### FR-CONTENT-003: Question Quality Review
**Priority**: P2  
**Description**: Content admins review flagged questions

**Acceptance Criteria**:
- View queue of flagged questions
- See student feedback/comments (future)
- See performance statistics
- Approve or edit questions
- Remove from queue when resolved

---

## 3.3 User Story Map - Student Journey

### 3.3.1 Journey Overview

The student journey through Athena follows a progression from **Onboarding** → **Daily Practice** → **Mastery** → **Assessment** → **Continued Growth**.

### 3.3.2 Detailed User Story Map

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        STUDENT LEARNING JOURNEY                          │
└─────────────────────────────────────────────────────────────────────────┘

PHASE 1: ONBOARDING (Day 1)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Epic: Account Setup & First Exposure

User Story 1.1: Account Creation
  As a new student
  I want to create an account with my email
  So that I can access the quiz platform
  
  Tasks:
  - Navigate to registration page
  - Enter email, password, full name
  - Verify email (future)
  - Complete profile setup

  Acceptance: Student has active account and can log in

User Story 1.2: Platform Orientation
  As a new student
  I want to see a tutorial of how the platform works
  So that I understand how to use it effectively
  
  Tasks:
  - View welcome modal with key features
  - Take guided tour of dashboard
  - Understand quiz vs. test distinction
  - Learn about spaced repetition benefits

  Acceptance: Student completes orientation, can navigate platform

User Story 1.3: Diagnostic Assessment
  As a new student
  I want to take an initial diagnostic quiz
  So that the system knows my starting knowledge level
  
  Tasks:
  - Select primary topic (e.g., Java)
  - Take 20-question diagnostic quiz
  - Answer questions across difficulty levels
  - Receive initial baseline score

  Acceptance: System establishes baseline for FSRS algorithm


PHASE 2: DAILY PRACTICE (Weeks 1-4)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Epic: Regular Quiz Practice & Knowledge Building

User Story 2.1: Daily Quiz Access
  As a student
  I want to take a daily adaptive quiz
  So that I can practice and improve my knowledge
  
  Tasks:
  - Log in to platform
  - View recommended topics for today
  - Click "Start Daily Quiz"
  - Answer 10-15 questions
  - Receive immediate feedback on each question
  - View session summary

  Acceptance: Student completes quiz, score recorded

User Story 2.2: Review Mistakes
  As a student
  I want to review questions I answered incorrectly
  So that I can learn from my mistakes
  
  Tasks:
  - Click on completed quiz attempt
  - See list of all questions
  - Identify incorrect answers (marked in red)
  - View correct answer for each mistake
  - Understand why answer was correct

  Acceptance: Student can review all past attempts

User Story 2.3: Track Progress
  As a student
  I want to see my progress over time
  So that I know I'm improving
  
  Tasks:
  - Navigate to Progress Dashboard
  - View overall accuracy percentage
  - See progress bars per topic
  - Check recent activity chart
  - Note areas of strength/weakness

  Acceptance: Dashboard displays accurate, up-to-date metrics

User Story 2.4: Identify Weak Areas
  As a student
  I want to see which concepts I'm struggling with
  So that I can focus my study time
  
  Tasks:
  - View "Insights" section on dashboard
  - See tag-based weakness analysis
  - Note concepts with <60% accuracy
  - Generate focused quiz on weak tags
  - Practice weak areas

  Acceptance: Student identifies specific concepts to study

User Story 2.5: Build Streak Motivation
  As a student
  I want to maintain a daily quiz streak
  So that I stay motivated and consistent
  
  Tasks:
  - View current streak count
  - Receive reminder notification (future)
  - Complete quiz before day ends
  - See streak milestone achievements

  Acceptance: Streak count increments daily, resets if day missed


PHASE 3: MASTERY BUILDING (Weeks 5-8)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Epic: Deepening Knowledge & Handling Complexity

User Story 3.1: Progress to Harder Questions
  As an improving student
  I want to tackle more difficult questions
  So that I can deepen my understanding
  
  Tasks:
  - Achieve 80%+ accuracy on difficulty 1-3 questions
  - System unlocks difficulty 4-5 questions
  - Attempt harder questions in daily quizzes
  - Apply knowledge to complex scenarios
  - Maintain high accuracy

  Acceptance: Student successfully answers difficulty 4+ questions

User Story 3.2: Multi-Topic Practice
  As a student
  I want to practice multiple topics in one session
  So that I maintain breadth of knowledge
  
  Tasks:
  - Select multiple topics for quiz (Java + SQL)
  - Take mixed-topic quiz
  - Context-switch between topics
  - Maintain accuracy across topics

  Acceptance: Student completes multi-topic quiz successfully

User Story 3.3: Long-Term Retention
  As a student who learned a concept weeks ago
  I want to be reminded of that concept periodically
  So that I don't forget it
  
  Tasks:
  - (Automatic) FSRS schedules old question for review
  - Student sees question from Week 2 in Week 7 quiz
  - Successfully recalls correct answer
  - System extends review interval further

  Acceptance: Student maintains >75% accuracy on old concepts

User Story 3.4: Concept Mastery Certification
  As a student
  I want to see that I've mastered a concept
  So that I feel confident and accomplished
  
  Tasks:
  - Achieve 90%+ accuracy on all questions with a specific tag
  - Receive "Mastery Badge" for that concept
  - View collection of mastery badges
  - Share achievement (future: social features)

  Acceptance: Student earns badges, visible on profile


PHASE 4: FORMAL ASSESSMENT (Weeks 4, 8, 12)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Epic: High-Stakes Testing & Evaluation

User Story 4.1: Test Preparation
  As a student approaching a test
  I want to know if I'm ready
  So that I can decide if I need more study
  
  Tasks:
  - View upcoming test schedule
  - Check progress dashboard for test topics
  - Note mastery level for each topic
  - Take additional practice quizzes if needed
  - Review weak areas

  Acceptance: Student feels prepared, confidence level measured

User Story 4.2: Taking Timed Test
  As a student
  I want to take a formal timed test
  So that I can demonstrate my knowledge
  
  Tasks:
  - Navigate to "Tests" section
  - Click "Start Test" for Week 4 Assessment
  - Accept test rules (60 min, no pause, single attempt)
  - Answer 50 questions within time limit
  - Monitor countdown timer
  - Submit before timer expires
  - View final score

  Acceptance: Student completes test, score recorded

User Story 4.3: Test Results Review
  As a student who completed a test
  I want to see my results and understand mistakes
  So that I can learn from the test
  
  Tasks:
  - View test score and percentile rank
  - See breakdown by topic
  - Review incorrect answers (after availability window)
  - Compare score to previous tests
  - Identify areas for future study

  Acceptance: Student understands performance, learns from test

User Story 4.4: Instructor Feedback
  As a student who struggled on a test
  I want my instructor to reach out
  So that I can get help to improve
  
  Tasks:
  - (Automatic) Instructor sees low test score in dashboard
  - Instructor schedules 1:1 meeting
  - Student receives notification
  - Student meets with instructor
  - Receives personalized study plan

  Acceptance: Struggling student receives timely intervention


PHASE 5: CONTINUED GROWTH (Weeks 9-12+)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Epic: Maintaining Knowledge & Expanding to New Topics

User Story 5.1: New Topic Introduction
  As a student who mastered Java
  I want to start learning Python
  So that I can expand my skill set
  
  Tasks:
  - Add Python to active topics
  - Take Python diagnostic quiz
  - Begin daily practice with Python questions
  - Maintain Java knowledge through spaced reviews

  Acceptance: Student practices multiple topics concurrently

User Story 5.2: Efficient Maintenance Reviews
  As a student who mastered a topic
  I want to review it efficiently
  So that I don't waste time on things I know
  
  Tasks:
  - (Automatic) FSRS schedules mastered questions infrequently
  - Student sees only 2-3 review questions per session
  - Quickly confirms knowledge
  - Spends majority of time on new/weak areas

  Acceptance: Review burden decreases as mastery increases

User Story 5.3: Graduation Readiness
  As a student nearing course completion
  I want to verify I'm ready for job interviews
  So that I can confidently apply for positions
  
  Tasks:
  - View comprehensive progress report
  - Verify 80%+ mastery across all topics
  - Take final comprehensive test
  - Receive course completion certificate
  - Export portfolio report

  Acceptance: Student meets graduation criteria, receives certificate

User Story 5.4: Post-Graduation Continued Access
  As a graduated student
  I want to continue using Athena
  So that I can prepare for interviews and maintain skills
  
  Tasks:
  - Retain account access after graduation
  - Take periodic refresher quizzes
  - Use before technical interviews
  - Maintain long-term retention

  Acceptance: Alumni can access platform indefinitely
```

### 3.3.3 Instructor Journey (Parallel to Student Journey)

```
INSTRUCTOR MONITORING JOURNEY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Week 1: Cohort Onboarding
  User Story: As an instructor, I want to see all new students in my cohort
  Tasks:
  - View cohort enrollment list
  - Verify all students have accounts
  - Monitor completion of diagnostic assessments
  - Identify students who haven't started

Week 2-3: Early Intervention
  User Story: As an instructor, I want to identify struggling students early
  Tasks:
  - Review cohort dashboard
  - Filter students with <60% accuracy
  - Drill into individual student analytics
  - Schedule intervention meetings
  - Assign targeted practice

Week 4: First Assessment
  User Story: As an instructor, I want to evaluate cohort progress
  Tasks:
  - Create Week 4 Test with selected questions
  - Set time limit and availability window
  - Monitor test completion rates
  - Review test results
  - Identify common wrong answers
  - Adjust upcoming curriculum

Week 5-7: Curriculum Adjustment
  User Story: As an instructor, I want to address common gaps
  Tasks:
  - Analyze tag-based cohort weaknesses
  - Note that 60% of class struggles with "recursion"
  - Plan additional classroom time on recursion
  - Create focused practice quiz on recursion
  - Assign to all students
  - Monitor improvement

Week 8: Mid-Course Assessment
  User Story: As an instructor, I want to measure progress
  Tasks:
  - Create Week 8 comprehensive test
  - Compare results to Week 4 test
  - Measure improvement rates
  - Celebrate student growth
  - Identify persistent strugglers

Week 9-12: Final Push
  User Story: As an instructor, I want all students to succeed
  Tasks:
  - Intensive support for struggling students
  - Monitor daily quiz completion rates
  - Provide encouragement and motivation
  - Prepare final comprehensive test
  - Evaluate graduation readiness
```

---

## 3.4 Data Model Specification

### 3.4.1 Core Entities

#### User
```
User {
  ID: string (UUID)
  Email: string (unique)
  PasswordHash: string
  FullName: string
  Role: enum (student, instructor, content_admin, system_admin)
  TotalPoints: int (sum of all correct answers)
  CreatedAt: timestamp
  LastLoginAt: timestamp
  CohortID: string (optional, for students)
  IsActive: boolean
}
```

#### Topic
```
Topic {
  ID: string (UUID)
  Name: string (e.g., "Java", "Python", "SQL")
  Description: string
  DisplayOrder: int
  IsActive: boolean
  CreatedAt: timestamp
  QuestionCount: int (computed)
}
```

#### Question
```
Question {
  ID: string (UUID)
  QuestionID: string (format: "topicID.questionNumber", e.g., "3.17")
  TopicID: string (foreign key to Topic)
  QuestionText: string (supports markdown/code)
  DifficultyRating: int (1-7)
  Tags: array of strings (e.g., ["loops", "control-structures"])
  IsActive: boolean
  CreatedAt: timestamp
  CreatedBy: string (UserID)
  LastModifiedAt: timestamp
}
```

#### Answer
```
Answer {
  ID: string (UUID)
  QuestionID: string (foreign key)
  AnswerText: string
  IsCorrect: boolean
  DisplayOrder: int (1-4)
}
```

#### QuizAttempt
```
QuizAttempt {
  ID: string (UUID)
  UserID: string (foreign key)
  StartDateTime: timestamp
  CompletedDateTime: timestamp (null if in progress)
  TotalQuestions: int
  CorrectAnswers: int
  IncorrectAnswers: int
  TopicIDs: array of strings (topics included in quiz)
  PointsEarned: int
  QuestionResults: array of {QuestionID, IsCorrect, TimeSpentSeconds}
}
```

#### QuestionAttempt
```
QuestionAttempt {
  ID: string (UUID)
  UserID: string (foreign key)
  QuestionID: string (foreign key)
  IsCorrect: boolean
  AttemptDateTime: timestamp
  TimeSpentSeconds: int
  QuizAttemptID: string (optional, foreign key)
  SelectedAnswerID: string (foreign key to Answer)
}
```

#### UserQuestionState (for FSRS)
```
UserQuestionState {
  UserID: string (composite key)
  QuestionID: string (composite key)
  Stability: float (days)
  Difficulty: float (0-10, FSRS difficulty, not question difficulty)
  LastRetrievability: float (0-1)
  LastAttemptDateTime: timestamp
  NextScheduledReview: timestamp
  TotalAttempts: int
  CorrectAttempts: int
  IncorrectAttempts: int
  AverageTimeSpentSeconds: int
  LastUpdated: timestamp
}
```

#### Test
```
Test {
  ID: string (UUID)
  Title: string
  Description: string
  TimeLimitMinutes: int
  QuestionIDs: array of strings (ordered)
  AvailableFrom: timestamp
  AvailableUntil: timestamp
  CohortIDs: array of strings
  CreatedBy: string (UserID)
  CreatedAt: timestamp
  IsActive: boolean
}
```

#### TestAttempt
```
TestAttempt {
  ID: string (UUID)
  TestID: string (foreign key)
  UserID: string (foreign key)
  StartDateTime: timestamp
  CompletedDateTime: timestamp
  TotalQuestions: int
  CorrectAnswers: int
  Score: float (percentage)
  TimeSpentMinutes: int
  QuestionResults: array of {QuestionID, IsCorrect, SelectedAnswerID}
  PercentileRank: float (compared to cohort)
}
```

### 3.4.2 Entity Relationships

```
User (1) ──── (many) QuizAttempt
User (1) ──── (many) QuestionAttempt
User (1) ──── (many) UserQuestionState
User (1) ──── (many) TestAttempt

Topic (1) ──── (many) Question

Question (1) ──── (many) Answer
Question (1) ──── (many) QuestionAttempt
Question (1) ──── (many) UserQuestionState

QuizAttempt (1) ──── (many) QuestionAttempt

Test (1) ──── (many) TestAttempt
Test (many) ──── (many) Question (via QuestionIDs array)
```

---

## 3.5 API Endpoints Specification

### 3.5.1 Authentication Endpoints

```
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/logout
POST   /api/auth/refresh-token
POST   /api/auth/forgot-password
POST   /api/auth/reset-password
```

### 3.5.2 User Endpoints

```
GET    /api/users/me
PUT    /api/users/me
GET    /api/users/:id (instructor only)
GET    /api/users (instructor only)
PUT    /api/users/:id/points
```

### 3.5.3 Topic Endpoints

```
GET    /api/topics
GET    /api/topics/:id
POST   /api/topics (instructor only)
PUT    /api/topics/:id (instructor only)
DELETE /api/topics/:id (instructor only)
```

### 3.5.4 Question Endpoints

```
GET    /api/questions
GET    /api/questions/:id
POST   /api/questions (instructor only)
PUT    /api/questions/:id (instructor only)
DELETE /api/questions/:id (instructor only)
GET    /api/questions/by-topic/:topicId
POST   /api/questions/bulk-import (content admin only)
```

### 3.5.5 Quiz Endpoints

```
POST   /api/quiz/generate
POST   /api/quiz/submit
GET    /api/quiz/attempts
GET    /api/quiz/attempts/:id
GET    /api/quiz/attempts/user/:userId (instructor only)
```

### 3.5.6 Test Endpoints

```
GET    /api/tests
GET    /api/tests/:id
POST   /api/tests (instructor only)
PUT    /api/tests/:id (instructor only)
POST   /api/tests/:id/start
POST   /api/tests/:id/submit
GET    /api/test-attempts
GET    /api/test-attempts/:id
GET    /api/test-attempts/user/:userId (instructor only)
```

### 3.5.7 Analytics Endpoints

```
GET    /api/analytics/student/:userId
GET    /api/analytics/cohort/:cohortId
GET    /api/analytics/questions
GET    /api/analytics/questions/:questionId
GET    /api/analytics/tags
```

---

## 3.6 Technical Architecture

### 3.6.1 Technology Stack

**Frontend**
- Framework: React.js with TypeScript
- State Management: Redux or Context API
- UI Library: Material-UI or Tailwind CSS
- Charts: Chart.js or Recharts
- HTTP Client: Axios

**Backend**
- Language: Go (current) or Node.js/Python (alternative)
- Framework: Gin (Go) or Express.js or FastAPI
- Authentication: JWT tokens
- API Style: RESTful JSON

**Database**
- Primary: PostgreSQL (production)
- Development: In-memory (current) or SQLite
- ORM: GORM (Go) or Sequelize or SQLAlchemy
- Caching: Redis (future optimization)

**Infrastructure**
- Hosting: AWS, GCP, or Azure
- Container: Docker
- Orchestration: Kubernetes (for scale) or Docker Compose (for MVP)
- CI/CD: GitHub Actions
- Monitoring: Prometheus + Grafana

### 3.6.2 Security Requirements

- **SEC-001**: All passwords hashed with bcrypt (cost factor 12+)
- **SEC-002**: HTTPS only in production
- **SEC-003**: JWT tokens expire after 7 days
- **SEC-004**: Rate limiting on API endpoints (100 req/min per user)
- **SEC-005**: SQL injection prevention via parameterized queries
- **SEC-006**: XSS prevention via input sanitization
- **SEC-007**: CORS configured for frontend domain only
- **SEC-008**: Role-based access control on all endpoints
- **SEC-009**: Audit logging for admin actions
- **SEC-010**: Regular security dependency updates

### 3.6.3 Performance Requirements

- **PERF-001**: API response time <500ms for 95th percentile
- **PERF-002**: Quiz generation <300ms
- **PERF-003**: Support 500 concurrent users
- **PERF-004**: Database queries <100ms
- **PERF-005**: Frontend initial load <3 seconds
- **PERF-006**: Handle 1M+ question attempts in database

---

## 3.7 Success Criteria & Launch Plan

### 3.7.1 MVP (Minimum Viable Product) Features

**Phase 0 (Current): Backend Foundation** ✓
- [x] Data models
- [x] In-memory database
- [x] Basic API endpoints
- [x] 20 comprehensive tests

**Phase 1 (MVP): Core Student Experience**
- [ ] User authentication
- [ ] Basic quiz generation (without FSRS)
- [ ] Multiple choice question answering
- [ ] Progress tracking
- [ ] Simple student dashboard
- [ ] PostgreSQL database integration

**Phase 2: Spaced Repetition**
- [ ] FSRS algorithm implementation
- [ ] UserQuestionState tracking
- [ ] Adaptive quiz generation
- [ ] Review scheduling

**Phase 3: Instructor Features**
- [ ] Instructor dashboard
- [ ] Student progress analytics
- [ ] Cohort analytics
- [ ] Question management UI

**Phase 4: Testing & Quality**
- [ ] Timed test functionality
- [ ] Test creation UI
- [ ] Question performance analytics
- [ ] Tag-based insights

**Phase 5: Polish & Scale**
- [ ] Mobile responsive design
- [ ] Performance optimization
- [ ] Advanced analytics
- [ ] Gamification (badges, leaderboards)

### 3.7.2 Launch Criteria

**Before Beta Launch**:
1. All Phase 1 & 2 features complete
2. 500+ questions across 3 core topics (Java, Python, SQL)
3. User testing with 10 students
4. Security audit completed
5. 95%+ API test coverage

**Before Production Launch**:
1. All Phase 1-4 features complete
2. 1000+ questions across 5+ topics
3. Pilot cohort of 30 students (4 weeks)
4. Instructor feedback incorporated
5. Performance testing (500 concurrent users)
6. Documentation complete

### 3.7.3 Post-Launch Roadmap

**Q2 2026**:
- Mobile app (iOS/Android)
- Social features (study groups, peer challenges)
- Explanations for correct answers
- Code execution environment for programming questions

**Q3 2026**:
- Integration with LMS (Canvas, Moodle)
- Advanced analytics (predictive modeling)
- AI-generated questions (experimental)
- API for third-party integrations

**Q4 2026**:
- Multi-tenancy (support multiple bootcamps)
- White-label option
- Advanced gamification
- Machine learning for difficulty calibration

---

## 3.8 Appendix: Comparative Analysis

### 3.8.1 Athena vs. Competitors

| Feature | Athena | Anki | Quizlet | Kahoot! | Knowt |
|---------|--------|------|---------|---------|-------|
| Spaced Repetition | ✓ (FSRS) | ✓ (SM-2) | ✗ | ✗ | ✓ |
| Multiple Choice | ✓ | ✗ | ✓ | ✓ | ✓ |
| Code Formatting | ✓ | Limited | Limited | ✗ | Limited |
| Instructor Analytics | ✓ | ✗ | Limited | ✓ | Limited |
| Tag-Based Insights | ✓ | Limited | ✗ | ✗ | ✗ |
| Timed Tests | ✓ | ✗ | Limited | ✓ | Limited |
| Difficulty Progression | ✓ | ✗ | ✗ | ✗ | ✗ |
| Custom Question Bank | ✓ | ✓ | ✓ | ✓ | ✓ |
| Focus | Tech Ed | General | General | Engagement | Students |

### 3.8.2 Key Differentiators

1. **Purpose-Built for Coding Education**: Unlike general-purpose quiz tools, Athena is designed specifically for software development education with code-friendly question formatting and technical domain tags.

2. **Modern Spaced Repetition**: FSRS algorithm is more accurate than traditional SM-2 used in Anki, leading to better retention with fewer reviews.

3. **Comprehensive Instructor Tools**: Goes beyond student-facing features to provide deep analytics for teaching staff, enabling data-driven intervention.

4. **Tag-Based Gap Analysis**: Unique approach to identifying conceptual weaknesses through cross-question tag analysis.

5. **Dual Mode (Quiz + Test)**: Combines low-stakes practice with high-stakes assessment in one platform.

---

## Document Approval

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Product Manager | ___________ | ___________ | _______ |
| Lead Instructor | ___________ | ___________ | _______ |
| Technical Lead | ___________ | ___________ | _______ |
| ZipCode Leadership | ___________ | ___________ | _______ |

---

**Document Version History**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-27 | Product Team | Initial specification |

---

**End of Specification Document**
