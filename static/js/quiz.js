// Quiz JavaScript

// Check if user is logged in
const user = JSON.parse(localStorage.getItem('user'));
if (!user) {
    window.location.href = '/index.html';
}

// Display user info
document.getElementById('userDisplay').textContent = `Welcome, ${user.username}!`;

// Logout function
function logout() {
    localStorage.removeItem('user');
    window.location.href = '/index.html';
}

// Get URL parameters
const urlParams = new URLSearchParams(window.location.search);
const topicId = urlParams.get('topic_id');
const topicName = urlParams.get('topic_name');

if (!topicId) {
    window.location.href = '/dashboard.html';
}

// Quiz state
let quizData = null;
let currentQuestionIndex = 0;
let quizAttemptId = null;
let score = 0;

// Set topic title
document.getElementById('topicTitle').textContent = `${decodeURIComponent(topicName)} Quiz`;

// Generate quiz
async function generateQuiz() {
    try {
        const response = await fetch('/quiz/generate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                user_id: user.id,
                topic_id: parseInt(topicId),
                num_questions: 10
            }),
        });

        const data = await response.json();

        if (response.ok) {
            quizData = data;
            quizAttemptId = data.quiz_attempt_id;
            displayQuestion();
        } else {
            alert('Failed to generate quiz: ' + (data.error || 'Unknown error'));
            window.location.href = '/dashboard.html';
        }
    } catch (error) {
        console.error('Error generating quiz:', error);
        alert('An error occurred while generating the quiz.');
        window.location.href = '/dashboard.html';
    }
}

// Display current question
function displayQuestion() {
    if (currentQuestionIndex >= quizData.questions.length) {
        completeQuiz();
        return;
    }

    const question = quizData.questions[currentQuestionIndex];
    
    // Update progress
    const progress = ((currentQuestionIndex + 1) / quizData.questions.length) * 100;
    document.getElementById('progressBar').style.width = `${progress}%`;
    document.getElementById('progressBadge').textContent = `Question ${currentQuestionIndex + 1} of ${quizData.questions.length}`;

    // Update question
    document.getElementById('questionText').textContent = question.question.question_text;
    
    // Update difficulty badge
    const difficultyText = `Difficulty ${question.question.difficulty}/7`;
    const difficultyBadge = document.getElementById('difficultyBadge');
    difficultyBadge.textContent = difficultyText;
    
    // Color code by difficulty
    difficultyBadge.className = 'badge ';
    if (question.question.difficulty <= 2) {
        difficultyBadge.className += 'bg-success';
    } else if (question.question.difficulty <= 4) {
        difficultyBadge.className += 'bg-info';
    } else if (question.question.difficulty <= 6) {
        difficultyBadge.className += 'bg-warning';
    } else {
        difficultyBadge.className += 'bg-danger';
    }

    // Display answers
    const answersContainer = document.getElementById('answersContainer');
    answersContainer.innerHTML = '';

    question.answers.forEach((answer, index) => {
        const answerDiv = document.createElement('div');
        answerDiv.className = 'answer-option';
        answerDiv.innerHTML = `
            <label style="cursor: pointer; display: block; margin: 0;">
                <input type="radio" name="answer" value="${answer.id}" data-question-id="${question.question.id}">
                ${answer.answer_text}
            </label>
        `;
        
        answerDiv.addEventListener('click', function() {
            // Remove selected class from all answers
            document.querySelectorAll('.answer-option').forEach(opt => {
                opt.classList.remove('selected');
            });
            // Add selected class to this answer
            this.classList.add('selected');
            // Check the radio button
            this.querySelector('input[type="radio"]').checked = true;
        });

        answersContainer.appendChild(answerDiv);
    });

    // Add submit button
    const submitButton = document.createElement('button');
    submitButton.className = 'btn btn-primary mt-3';
    submitButton.textContent = 'Submit Answer';
    submitButton.addEventListener('click', submitAnswer);
    answersContainer.appendChild(submitButton);

    // Hide feedback
    document.getElementById('feedbackContainer').style.display = 'none';
}

// Submit answer
async function submitAnswer() {
    const selectedAnswer = document.querySelector('input[name="answer"]:checked');
    
    if (!selectedAnswer) {
        alert('Please select an answer');
        return;
    }

    const answerId = parseInt(selectedAnswer.value);
    const questionId = parseInt(selectedAnswer.dataset.questionId);

    try {
        const response = await fetch('/quiz/submit-answer', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                quiz_attempt_id: quizAttemptId,
                question_id: questionId,
                answer_id: answerId,
                user_id: user.id
            }),
        });

        const data = await response.json();

        if (response.ok) {
            // Show feedback
            const feedbackContainer = document.getElementById('feedbackContainer');
            const feedbackAlert = document.getElementById('feedbackAlert');
            
            if (data.correct) {
                score++;
                feedbackAlert.className = 'alert alert-success';
                feedbackAlert.innerHTML = '<strong>Correct!</strong> Well done!';
                
                // Highlight correct answer
                document.querySelector(`input[value="${answerId}"]`).closest('.answer-option').classList.add('correct');
            } else {
                feedbackAlert.className = 'alert alert-danger';
                feedbackAlert.innerHTML = `<strong>Incorrect.</strong> The correct answer was: <em>${data.correct_answer}</em>`;
                
                // Highlight incorrect answer
                document.querySelector(`input[value="${answerId}"]`).closest('.answer-option').classList.add('incorrect');
                
                // Find and highlight correct answer
                const question = quizData.questions[currentQuestionIndex];
                const correctAnswer = question.answers.find(a => a.is_correct);
                if (correctAnswer) {
                    document.querySelector(`input[value="${correctAnswer.id}"]`).closest('.answer-option').classList.add('correct');
                }
            }

            feedbackContainer.style.display = 'block';
            
            // Disable all answer options
            document.querySelectorAll('.answer-option').forEach(opt => {
                opt.style.pointerEvents = 'none';
            });
            document.querySelector('button.btn-primary').style.display = 'none';
            
            // Setup next button
            document.getElementById('nextButton').onclick = () => {
                currentQuestionIndex++;
                displayQuestion();
            };
        } else {
            alert('Failed to submit answer: ' + (data.error || 'Unknown error'));
        }
    } catch (error) {
        console.error('Error submitting answer:', error);
        alert('An error occurred while submitting your answer.');
    }
}

// Complete quiz
async function completeQuiz() {
    try {
        const response = await fetch('/quiz/complete', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                quiz_attempt_id: quizAttemptId
            }),
        });

        if (response.ok) {
            // Hide quiz content
            document.getElementById('quizContent').style.display = 'none';
            
            // Show completion screen
            const percentage = Math.round((score / quizData.questions.length) * 100);
            document.getElementById('finalScore').textContent = `${percentage}%`;
            document.getElementById('finalResults').textContent = `${score} / ${quizData.questions.length}`;
            document.getElementById('quizComplete').style.display = 'block';
        }
    } catch (error) {
        console.error('Error completing quiz:', error);
    }
}

// Start the quiz
generateQuiz();
