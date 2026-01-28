# Quiz Retake Feature Documentation

## Overview

The quiz retake feature allows students to retake quizzes they didn't score 100% on, using the exact same questions. This helps students practice and improve their understanding of topics they struggled with.

## Key Features

### 1. **Exact Question Reproduction**
- Retake quizzes contain the exact same questions as the original quiz
- Questions appear in the same order
- Ensures fair assessment and focused practice

### 2. **Parent-Child Relationship Tracking**
- Each retake is linked to its parent attempt
- Attempt numbers increment (1st attempt, 2nd attempt, etc.)
- Complete audit trail of student progress

### 3. **FSRS Integration**
- Questions answered incorrectly are tracked in the spaced repetition system
- Incorrect questions get shorter review intervals
- Correct questions get progressively longer intervals
- Adaptive quiz generation prioritizes due questions

### 4. **Smart Retake Requirements**
- Only completed quizzes can be retaken
- Students with 100% scores cannot retake (no need)
- Students can only retake their own quizzes (security)

## Database Schema Changes

### QuizAttempt Model

Three new fields were added to support retakes:

```go
type QuizAttempt struct {
    // ... existing fields ...
    ParentAttemptID *int  `json:"parent_attempt_id,omitempty"` // Links to original attempt
    AttemptNumber   int   `json:"attempt_number"`              // 1, 2, 3, etc.
    IsRetake        bool  `json:"is_retake"`                   // Flags retake attempts
}
```

### SQLite Schema

```sql
CREATE TABLE quiz_attempts (
    -- ... existing columns ...
    parent_attempt_id INTEGER,
    attempt_number INTEGER DEFAULT 1,
    is_retake BOOLEAN DEFAULT 0,
    FOREIGN KEY (parent_attempt_id) REFERENCES quiz_attempts(id)
);
```

## API Endpoints

### POST /quiz/retake

Generate a retake quiz with the same questions as a previous attempt.

**Request Body:**
```json
{
    "parent_attempt_id": 123,
    "user_id": 456
}
```

**Response:**
```json
{
    "quiz_attempt_id": 124,
    "questions": [
        {
            "question": {
                "id": 1,
                "topic_question_id": "go.1",
                "question_text": "What is a goroutine?",
                // ... other fields
            },
            "answers": [
                {
                    "id": 1,
                    "answer_text": "A lightweight thread",
                    "is_correct": true
                },
                // ... other answers
            ]
        }
        // ... more questions
    ]
}
```

**Error Responses:**
- `400 Bad Request` - Cannot retake incomplete quiz
- `400 Bad Request` - Cannot retake quiz with 100% score
- `401 Unauthorized` - User doesn't own the quiz attempt
- `404 Not Found` - Parent quiz attempt not found

## Usage Examples

### Student Flow

1. **Take Initial Quiz**
   ```bash
   POST /quiz/generate
   {
     "user_id": 1,
     "topic_id": 2,
     "num_questions": 10
   }
   ```

2. **Complete Quiz (Score: 70%)**
   ```bash
   POST /quiz/complete
   {
     "quiz_attempt_id": 123
   }
   ```

3. **Request Retake**
   ```bash
   POST /quiz/retake
   {
     "parent_attempt_id": 123,
     "user_id": 1
   }
   ```

4. **Take Retake Quiz**
   - Same questions as original
   - Attempt number: 2
   - Fresh start with 0 correct answers

5. **Complete Retake (Score: 90%)**
   - FSRS states updated for all questions
   - Points awarded for correct answers
   - Can retake again if score < 100%

### Multiple Retakes

Students can retake a quiz multiple times:

```
Attempt 1 (Original) - Score: 60%
   ↓
Attempt 2 (Retake 1) - Score: 80%
   ↓
Attempt 3 (Retake 2) - Score: 100%
```

Each retake links to its immediate parent, creating a chain.

## FSRS Integration Details

### On Quiz Completion

For each question in the quiz:

1. **Get or Create FSRS State**
   - If first time seeing question: create new state
   - Otherwise: load existing state

2. **Apply Rating Based on Correctness**
   - Correct answer → Rating 3 (Good)
   - Incorrect answer → Rating 1 (Again)

3. **Update FSRS Parameters**
   - Stability, Difficulty, Due date
   - Reps (repetitions), Lapses (failures)
   - State (New → Learning → Review → Relearning)

4. **Schedule Next Review**
   - Correct answers: longer intervals
   - Incorrect answers: shorter intervals

### Adaptive Quiz Generation

The existing `/fsrs/quiz/generate` endpoint automatically prioritizes:

1. **Due questions** (past their review date)
2. **New questions** (never seen before)
3. **Future questions** (scheduled but not yet due)

Questions answered incorrectly will have shorter intervals and appear sooner in adaptive quizzes.

## Testing

### Automated Tests

Run the comprehensive test suite:

```bash
go test -v ./api -run TestRetake
go test -v ./api -run TestQuizCompletion
go test -v ./api -run TestFSRS
```

### Manual Integration Test

Run the manual test script:

```bash
./test_retake_feature.sh
```

This script:
1. Creates test user and topic
2. Creates test questions with answers
3. Generates and completes quiz with 50% score
4. Requests retake and verifies same questions
5. Tests 100% score rejection
6. Verifies FSRS integration

## Implementation Notes

### Design Decisions

1. **Same Questions Guaranteed**
   - QuizAttempt stores `QuestionIDs` (topic_question_id strings)
   - Retake uses parent's QuestionIDs to reconstruct exact question set
   - Ensures fairness and focused practice

2. **Parent-Child Linking**
   - Each retake points to immediate parent, not original
   - Allows tracking of progress across multiple retakes
   - Maintains clear audit trail

3. **FSRS Rating Simplification**
   - Correct → "Good" (3) instead of Easy/Good/Hard
   - Incorrect → "Again" (1)
   - Simplifies UX while maintaining FSRS effectiveness

4. **100% Score Rejection**
   - Students who score 100% don't need retakes
   - Prevents unnecessary quiz attempts
   - Encourages moving to new material

### Future Enhancements

Potential improvements for future versions:

1. **Time Limits on Retakes**
   - Require waiting period between retakes
   - Prevent rapid-fire retakes without studying

2. **Limited Retake Count**
   - Cap retakes at 3-5 attempts per quiz
   - Encourage seeking help rather than endless retries

3. **Retake Analytics**
   - Show improvement trends across retakes
   - Identify questions consistently failed
   - Suggest topics for review

4. **Partial Retakes**
   - Only retake incorrectly answered questions
   - More efficient for near-perfect scores

## Troubleshooting

### Common Issues

**Q: Retake has different questions than original**
- Check that parent_attempt_id is correct
- Verify QuestionIDs are stored in parent attempt

**Q: Cannot retake quiz**
- Ensure quiz is completed (CompletedTime is set)
- Check score is not 100%
- Verify user owns the quiz attempt

**Q: FSRS states not updating**
- Check that CompleteQuizHandler is called
- Verify questions exist in database
- Check for errors in server logs

## Related Files

- `models/models.go` - QuizAttempt model definition
- `database/sqlite.go` - Database schema and operations
- `api/quiz.go` - Quiz and retake handlers
- `api/quiz_retake_test.go` - Retake functionality tests
- `api/quiz_fsrs_test.go` - FSRS integration tests
- `main.go` - Route registration
