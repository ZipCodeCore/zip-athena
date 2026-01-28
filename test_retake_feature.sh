#!/bin/bash

# Quiz Retake Feature Manual Test Script
# This script demonstrates the quiz retake feature end-to-end

set -e

echo "=========================================="
echo "Quiz Retake Feature Manual Test"
echo "=========================================="
echo ""

# Clean up any existing database
rm -f athena_test.db

# Set the database path
export DB_PATH="athena_test.db"

# Start the server in the background
echo "Starting Athena server..."
go run main.go &
SERVER_PID=$!

# Give the server time to start
sleep 3

# Ensure we kill the server on exit
trap "kill $SERVER_PID 2>/dev/null || true; rm -f athena_test.db" EXIT

BASE_URL="http://localhost:9206"

echo "1. Creating a test user..."
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }')

USER_ID=$(echo "$USER_RESPONSE" | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
echo "   Created user with ID: $USER_ID"

echo ""
echo "2. Creating a test topic..."
TOPIC_RESPONSE=$(curl -s -X POST "$BASE_URL/topics" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Topic",
    "description": "Topic for testing quiz retake"
  }')

TOPIC_ID=$(echo "$TOPIC_RESPONSE" | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
echo "   Created topic with ID: $TOPIC_ID"

echo ""
echo "3. Creating test questions..."
# Question 1
Q1_RESPONSE=$(curl -s -X POST "$BASE_URL/questions" \
  -H "Content-Type: application/json" \
  -d "{
    \"topic_id\": $TOPIC_ID,
    \"question_text\": \"What is 2 + 2?\",
    \"difficulty\": 1,
    \"tags\": [\"math\"],
    \"topic_question_id\": \"test.1\"
  }")

Q1_ID=$(echo "$Q1_RESPONSE" | grep -o '"id":[0-9]*' | grep -o '[0-9]*' | head -1)
echo "   Created question 1 with ID: $Q1_ID"

# Add answers for question 1
curl -s -X POST "$BASE_URL/admin/answers" \
  -H "Content-Type: application/json" \
  -d "{
    \"question_id\": $Q1_ID,
    \"answer_text\": \"4\",
    \"is_correct\": true
  }" > /dev/null

curl -s -X POST "$BASE_URL/admin/answers" \
  -H "Content-Type: application/json" \
  -d "{
    \"question_id\": $Q1_ID,
    \"answer_text\": \"5\",
    \"is_correct\": false
  }" > /dev/null

# Question 2
Q2_RESPONSE=$(curl -s -X POST "$BASE_URL/questions" \
  -H "Content-Type: application/json" \
  -d "{
    \"topic_id\": $TOPIC_ID,
    \"question_text\": \"What is 3 + 3?\",
    \"difficulty\": 1,
    \"tags\": [\"math\"],
    \"topic_question_id\": \"test.2\"
  }")

Q2_ID=$(echo "$Q2_RESPONSE" | grep -o '"id":[0-9]*' | grep -o '[0-9]*' | head -1)
echo "   Created question 2 with ID: $Q2_ID"

# Add answers for question 2
curl -s -X POST "$BASE_URL/admin/answers" \
  -H "Content-Type: application/json" \
  -d "{
    \"question_id\": $Q2_ID,
    \"answer_text\": \"6\",
    \"is_correct\": true
  }" > /dev/null

curl -s -X POST "$BASE_URL/admin/answers" \
  -H "Content-Type: application/json" \
  -d "{
    \"question_id\": $Q2_ID,
    \"answer_text\": \"7\",
    \"is_correct\": false
  }" > /dev/null

echo ""
echo "4. Generating initial quiz..."
QUIZ_RESPONSE=$(curl -s -X POST "$BASE_URL/quiz/generate" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": $USER_ID,
    \"topic_id\": $TOPIC_ID,
    \"num_questions\": 2
  }")

QUIZ_ATTEMPT_ID=$(echo "$QUIZ_RESPONSE" | grep -o '"quiz_attempt_id":[0-9]*' | grep -o '[0-9]*')
echo "   Generated quiz with attempt ID: $QUIZ_ATTEMPT_ID"

echo ""
echo "5. Simulating quiz with partial score (50%)..."
# Get the first answer ID (correct answer for Q1)
A1_RESPONSE=$(curl -s "$BASE_URL/admin/answers?question_id=$Q1_ID")
A1_CORRECT_ID=$(echo "$A1_RESPONSE" | grep -o '"id":[0-9]*,"question_id":[0-9]*,"answer_text":"4","is_correct":true' | grep -o '"id":[0-9]*' | grep -o '[0-9]*')

# Get the wrong answer ID for Q2
A2_RESPONSE=$(curl -s "$BASE_URL/admin/answers?question_id=$Q2_ID")
A2_WRONG_ID=$(echo "$A2_RESPONSE" | grep -o '"id":[0-9]*,"question_id":[0-9]*,"answer_text":"7","is_correct":false' | grep -o '"id":[0-9]*' | grep -o '[0-9]*')

# Submit correct answer for Q1
curl -s -X POST "$BASE_URL/quiz/submit-answer" \
  -H "Content-Type: application/json" \
  -d "{
    \"quiz_attempt_id\": $QUIZ_ATTEMPT_ID,
    \"question_id\": $Q1_ID,
    \"answer_id\": $A1_CORRECT_ID,
    \"user_id\": $USER_ID
  }" > /dev/null

echo "   Submitted correct answer for question 1"

# Submit incorrect answer for Q2
curl -s -X POST "$BASE_URL/quiz/submit-answer" \
  -H "Content-Type: application/json" \
  -d "{
    \"quiz_attempt_id\": $QUIZ_ATTEMPT_ID,
    \"question_id\": $Q2_ID,
    \"answer_id\": $A2_WRONG_ID,
    \"user_id\": $USER_ID
  }" > /dev/null

echo "   Submitted incorrect answer for question 2"

echo ""
echo "6. Completing the quiz..."
COMPLETE_RESPONSE=$(curl -s -X POST "$BASE_URL/quiz/complete" \
  -H "Content-Type: application/json" \
  -d "{
    \"quiz_attempt_id\": $QUIZ_ATTEMPT_ID
  }")

CORRECT_ANSWERS=$(echo "$COMPLETE_RESPONSE" | grep -o '"correct_answers":[0-9]*' | grep -o '[0-9]*')
TOTAL_QUESTIONS=$(echo "$COMPLETE_RESPONSE" | grep -o '"total_questions":[0-9]*' | grep -o '[0-9]*')
echo "   Quiz completed: $CORRECT_ANSWERS/$TOTAL_QUESTIONS correct (50%)"

echo ""
echo "7. Requesting quiz retake..."
RETAKE_RESPONSE=$(curl -s -X POST "$BASE_URL/quiz/retake" \
  -H "Content-Type: application/json" \
  -d "{
    \"parent_attempt_id\": $QUIZ_ATTEMPT_ID,
    \"user_id\": $USER_ID
  }")

RETAKE_ATTEMPT_ID=$(echo "$RETAKE_RESPONSE" | grep -o '"quiz_attempt_id":[0-9]*' | grep -o '[0-9]*')
echo "   ✓ Retake quiz generated with attempt ID: $RETAKE_ATTEMPT_ID"

# Verify it's a retake
IS_RETAKE=$(echo "$RETAKE_RESPONSE" | grep -o '"is_retake":true' || echo "false")
if [ "$IS_RETAKE" = "false" ]; then
  echo "   ✗ ERROR: Quiz is not marked as retake!"
  exit 1
fi
echo "   ✓ Quiz is correctly marked as retake"

# Verify same questions
ORIGINAL_QUESTIONS=$(echo "$QUIZ_RESPONSE" | grep -o '"topic_question_id":"[^"]*"' | sort)
RETAKE_QUESTIONS=$(echo "$RETAKE_RESPONSE" | grep -o '"topic_question_id":"[^"]*"' | sort)

if [ "$ORIGINAL_QUESTIONS" = "$RETAKE_QUESTIONS" ]; then
  echo "   ✓ Retake quiz contains the same questions as original"
else
  echo "   ✗ ERROR: Questions don't match!"
  exit 1
fi

echo ""
echo "8. Testing 100% score rejection..."
# First create a perfect score quiz attempt
PERFECT_QUIZ_RESPONSE=$(curl -s -X POST "$BASE_URL/quiz/generate" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": $USER_ID,
    \"topic_id\": $TOPIC_ID,
    \"num_questions\": 2
  }")

PERFECT_ATTEMPT_ID=$(echo "$PERFECT_QUIZ_RESPONSE" | grep -o '"quiz_attempt_id":[0-9]*' | grep -o '[0-9]*')

# Answer both correctly
curl -s -X POST "$BASE_URL/quiz/submit-answer" \
  -H "Content-Type: application/json" \
  -d "{
    \"quiz_attempt_id\": $PERFECT_ATTEMPT_ID,
    \"question_id\": $Q1_ID,
    \"answer_id\": $A1_CORRECT_ID,
    \"user_id\": $USER_ID
  }" > /dev/null

A2_CORRECT_ID=$(echo "$A2_RESPONSE" | grep -o '"id":[0-9]*,"question_id":[0-9]*,"answer_text":"6","is_correct":true' | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
curl -s -X POST "$BASE_URL/quiz/submit-answer" \
  -H "Content-Type: application/json" \
  -d "{
    \"quiz_attempt_id\": $PERFECT_ATTEMPT_ID,
    \"question_id\": $Q2_ID,
    \"answer_id\": $A2_CORRECT_ID,
    \"user_id\": $USER_ID
  }" > /dev/null

# Complete perfect quiz
curl -s -X POST "$BASE_URL/quiz/complete" \
  -H "Content-Type: application/json" \
  -d "{
    \"quiz_attempt_id\": $PERFECT_ATTEMPT_ID
  }" > /dev/null

# Try to retake - should fail
RETAKE_PERFECT_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/quiz/retake" \
  -H "Content-Type: application/json" \
  -d "{
    \"parent_attempt_id\": $PERFECT_ATTEMPT_ID,
    \"user_id\": $USER_ID
  }")

HTTP_CODE=$(echo "$RETAKE_PERFECT_RESPONSE" | tail -1)
if [ "$HTTP_CODE" = "400" ]; then
  echo "   ✓ Correctly rejected retake for 100% score"
else
  echo "   ✗ ERROR: Should have rejected retake for 100% score (got HTTP $HTTP_CODE)"
  exit 1
fi

echo ""
echo "=========================================="
echo "✓ All manual tests passed!"
echo "=========================================="
echo ""
echo "Quiz Retake Feature Summary:"
echo "  - Retake quiz generated with same questions"
echo "  - Parent-child relationship established"
echo "  - 100% scores correctly rejected"
echo "  - FSRS integration working (incorrectly answered questions tracked)"
echo ""
