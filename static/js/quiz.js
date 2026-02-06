// Quiz JavaScript

// Check if user is logged in
const user = JSON.parse(localStorage.getItem('user'));
if (!user) {
    window.location.href = '/index.html';
}

// Display user info
document.getElementById('userDisplay').textContent = `Welcome, ${user.username}!`;

// Show admin button if user is admin
if (user.is_admin) {
    document.getElementById('instructorButton').style.display = 'inline-block';
}

// Logout function
function logout() {
    localStorage.removeItem('user');
    window.location.href = '/index.html';
}

// Get URL parameters
const urlParams = new URLSearchParams(window.location.search);
const topicId = urlParams.get('topic_id');
const topicName = urlParams.get('topic_name');

// Check if there's a retake quiz in localStorage
const currentQuiz = localStorage.getItem('currentQuiz');
// Check if there's a resume quiz in localStorage
const resumeQuiz = localStorage.getItem('resumeQuiz');

if (!topicId && !currentQuiz && !resumeQuiz) {
    window.location.href = '/dashboard.html';
}

// Quiz state
let quizData = null;
let currentQuestionIndex = 0;
let quizAttemptId = null;
let score = 0;
let answeredQuestions = {}; // Track which questions have been answered
let selectedDifficulty = 3; // Default to level 3 (all difficulties)

// Set topic title and initialize
if (resumeQuiz) {
    // This is a resumed quiz
    const resumeData = JSON.parse(resumeQuiz);
    document.getElementById('topicTitle').textContent = `Resume Quiz`;
    quizAttemptId = resumeData.attemptId;
    answeredQuestions = resumeData.questionResults || {};
    
    // Initialize score with previously correct answers
    score = Object.values(answeredQuestions).filter(result => result === true).length;
    
    // Clear the resume data from localStorage
    localStorage.removeItem('resumeQuiz');
    
    // Load the quiz attempt and resume
    loadAndResumeQuiz(resumeData);
} else if (currentQuiz) {
    // This is a retake quiz
    const quizObj = JSON.parse(currentQuiz);
    document.getElementById('topicTitle').textContent = `Quiz Retake`;
    // Initialize with retake data
    quizData = quizObj;
    quizAttemptId = quizObj.quiz_attempt_id;
    // Clear the retake quiz from localStorage
    localStorage.removeItem('currentQuiz');
    displayQuestion();
} else {
    document.getElementById('topicTitle').textContent = `${decodeURIComponent(topicName)} Quiz`;
    // Show difficulty selection modal for new quiz
    showDifficultyModal();
}

// Show difficulty selection modal
function showDifficultyModal() {
    const modal = new bootstrap.Modal(document.getElementById('difficultyModal'));
    modal.show();
    
    // Handle start quiz button click
    document.getElementById('startQuizBtn').addEventListener('click', function() {
        const selected = document.querySelector('input[name="difficulty"]:checked');
        if (selected) {
            selectedDifficulty = parseInt(selected.value);
            modal.hide();
            generateQuiz();
        }
    }, { once: true }); // Use once to avoid multiple event listeners
}

// Load and resume a quiz
async function loadAndResumeQuiz(resumeData) {
    try {
        // The resume data contains questions in the format [{question: {...}, answers: [...]}]
        quizData = {
            questions: resumeData.questions,
            quiz_attempt_id: resumeData.attemptId
        };
        
        // Find the first unanswered question
        currentQuestionIndex = 0;
        for (let i = 0; i < quizData.questions.length; i++) {
            const quizQuestion = quizData.questions[i];
            const question = quizQuestion.question;
            const questionKey = question.topic_question_id || question.id.toString();
            if (!answeredQuestions[questionKey]) {
                currentQuestionIndex = i;
                break;
            }
        }
        
        // If all questions answered, start at beginning
        if (currentQuestionIndex === 0 && Object.keys(answeredQuestions).length > 0) {
            const firstQuestion = quizData.questions[0].question;
            const firstQuestionKey = firstQuestion.topic_question_id || firstQuestion.id.toString();
            if (answeredQuestions[firstQuestionKey]) {
                // All questions answered but quiz not completed - start from end
                currentQuestionIndex = quizData.questions.length - 1;
            }
        }
        
        displayQuestion();
    } catch (error) {
        console.error('Error loading quiz for resume:', error);
        alert('Unable to resume quiz. Returning to dashboard.');
        window.location.href = '/dashboard.html';
    }
}

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
                num_questions: 10,
                difficulty_level: selectedDifficulty
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
                
                // Load and display related learning materials
                loadLearningResources(questionId, feedbackAlert);
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
            const completedQuiz = await response.json();
            
            // Hide quiz content
            document.getElementById('quizContent').style.display = 'none';
            
            // Show completion screen using the backend's correct count
            const totalQuestions = completedQuiz.total_questions;
            const correctAnswers = completedQuiz.correct_answers;
            const percentage = Math.round((correctAnswers / totalQuestions) * 100);
            document.getElementById('finalScore').textContent = `${percentage}%`;
            document.getElementById('finalResults').textContent = `${correctAnswers} / ${totalQuestions}`;
            document.getElementById('quizComplete').style.display = 'block';
        }
    } catch (error) {
        console.error('Error completing quiz:', error);
    }
}

// Load and display related learning resources for a question
async function loadLearningResources(questionId, feedbackElement) {
    try {
        const response = await fetch(`/quiz/question/${questionId}/learning`);
        if (!response.ok) {
            return; // Silently fail if no resources found
        }

        const data = await response.json();
        if (!data.resources || data.resources.length === 0) {
            return; // No resources to display
        }

        // Build learning resources HTML
        let resourcesHtml = '<div class="mt-3 p-3 bg-light rounded"><strong>📚 Review These Topics:</strong><ul class="mt-2 mb-0">';
        
        data.resources.forEach(resource => {
            const languageBadge = resource.language === 'java' 
                ? '<span class="badge bg-warning text-dark">Java</span>' 
                : '<span class="badge bg-primary">Python</span>';
            
            resourcesHtml += `
                <li class="mb-2">
                    ${languageBadge}
                    <a href="/learning.html?resource=${resource.id}" target="_blank" class="ms-2">
                        ${resource.section_number}. ${resource.section_title}
                    </a>
                    <small class="text-muted">(Phase ${resource.phase})</small>
                </li>
            `;
        });
        
        resourcesHtml += '</ul></div>';
        
        // Append to feedback
        feedbackElement.innerHTML += resourcesHtml;
        
    } catch (error) {
        console.error('Error loading learning resources:', error);
        // Silently fail - don't disrupt quiz flow
    }
}

// Quiz initialization is handled in the initialization section above
// (either through generateQuiz() or loading from localStorage)
