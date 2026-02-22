package router

import (
	"net/http"

	"github.com/ZipCodeCore/zip-athena/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New creates a new router with all application routes
func New(h *api.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Chi's built-in middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health and metrics endpoints (no auth required)
	r.Get("/health", h.HealthHandler)
	r.Get("/ready", h.ReadyHandler)
	r.Get("/metrics", h.MetricsHandler)

	// Serve static files
	r.Handle("/*", http.FileServer(http.Dir("./static")))

	// Authentication routes (no auth required)
	r.Post("/auth/register", h.RegisterHandler)
	r.Post("/auth/login", h.LoginHandler)
	r.Post("/auth/logout", h.LogoutHandler)
	r.Post("/auth/refresh", h.RefreshTokenHandler)

	// Public/student routes (may require auth via middleware in handlers)
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.CreateUserHandler)
		r.Get("/", h.GetUserHandler)
	})

	r.Route("/topics", func(r chi.Router) {
		r.Post("/", h.CreateTopicHandler)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			// Check if ID parameter is provided
			if r.URL.Query().Get("id") != "" {
				h.GetTopicHandler(w, r)
			} else {
				h.ListTopicsHandler(w, r)
			}
		})
	})

	r.Route("/questions", func(r chi.Router) {
		r.Post("/", h.CreateQuestionHandler)
		r.Get("/", h.GetQuestionsByTopicHandler)
	})

	// Quiz routes
	r.Post("/quiz/generate", h.GenerateQuizHandler)
	r.Post("/quiz/submit-answer", h.SubmitAnswerHandler)
	r.Post("/quiz/complete", h.CompleteQuizHandler)
	r.Post("/quiz/retake", h.RetakeQuizHandler)
	r.Get("/progress", h.GetProgressHandler)

	// Assigned Quiz routes (student)
	r.Get("/assigned-quizzes", h.GetAssignedQuizzesHandler)
	r.Post("/assigned-quizzes/start", h.TakeAssignedQuizHandler)

	r.Route("/quizzes", func(r chi.Router) {
		r.Get("/", h.GetQuizHandler)
	})

	r.Route("/quiz-attempts", func(r chi.Router) {
		r.Post("/", h.CreateQuizAttemptHandler)
		r.Get("/", h.GetQuizAttemptsByUserHandler)
	})

	// FSRS routes
	r.Post("/fsrs/review", h.ReviewQuestionHandler)
	r.Get("/fsrs/due", h.GetDueQuestionsHandler)
	r.Post("/fsrs/quiz/generate", h.GenerateAdaptiveQuizHandler)

	// Quiz History routes
	r.Get("/students/quiz-history", h.GetQuizHistoryHandler)
	r.Get("/students/quiz-history/{id}", h.GetQuizAttemptDetailsHandler)

	// Timed Test Taking routes
	r.Get("/tests/available", h.GetAvailableTestsHandler)
	r.Get("/tests/{id}", h.GetTestDetailsHandler)
	r.Post("/tests/{id}/start", h.StartTestHandler)
	r.Put("/test-attempts/{id}/answer", h.UpdateTestAnswerHandler)
	r.Post("/test-attempts/{id}/submit", h.SubmitTestHandler)
	r.Get("/test-attempts/{id}/results", h.GetTestResultsHandler)

	// Tag-Based Review routes
	r.Get("/tags", h.GetTagsHandler)
	r.Get("/tags/stats", h.GetTagStatsHandler)
	r.Post("/quiz/generate/by-tags", h.GenerateQuizByTagsHandler)

	// Learning Resources routes
	r.Get("/learning/outline", h.GetCurriculumOutlineHandler)
	r.Get("/learning/phase/{id}", h.GetPhaseResourcesHandler)
	r.Get("/learning/resource/{id}", h.GetLearningResourceHandler)
	r.Get("/learning/resource/{id}/content", h.GetLearningResourceContentHandler)
	r.Get("/learning/search", h.SearchLearningResourcesHandler)
	r.Post("/learning/track-view", h.TrackLearningViewHandler)
	r.Get("/quiz/question/{id}/learning", h.GetQuestionLearningLinksHandler)

	// Admin routes (protected by AdminMiddleware)
	r.Route("/admin", func(r chi.Router) {
		// Apply admin middleware to all routes in this group
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.AdminMiddleware(next.ServeHTTP)(w, r)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.ListUsersHandler)
			r.Post("/", h.CreateUserHandler)
			r.Put("/", h.UpdateUserHandler)
			r.Delete("/", h.DeleteUserHandler)
		})

		r.Route("/questions", func(r chi.Router) {
			r.Get("/", h.ListQuestionsHandler)
			r.Post("/", h.CreateQuestionHandler)
			r.Put("/", h.UpdateQuestionHandler)
			r.Delete("/", h.DeleteQuestionHandler)

			// JSON Import route
			r.Post("/import-json", h.ImportQuestionsJSONHandler)

			// CSV Import routes
			r.Post("/import", h.ImportQuestionsCSVHandler)
			r.Post("/import/validate", h.ValidateQuestionsCSVHandler)
			r.Post("/import/confirm", h.ConfirmImportQuestionsCSVHandler)

			// Quality Review routes
			r.Get("/flagged", h.GetFlaggedQuestionsHandler)
			r.Post("/auto-flag", h.AutoFlagQuestionsHandler)
			r.Post("/flag", h.FlagQuestionHandler)
			r.Delete("/flag", h.UnflagQuestionHandler)
			r.Put("/resolve", h.ResolveQuestionHandler)
		})

		r.Route("/quizzes", func(r chi.Router) {
			r.Get("/", h.ListQuizzesHandler)
			r.Post("/", h.CreateQuizHandler)
			r.Put("/", h.UpdateQuizHandler)
			r.Delete("/", h.DeleteQuizHandler)
		})

		r.Route("/topics", func(r chi.Router) {
			r.Put("/", h.UpdateTopicHandler)
			r.Delete("/", h.DeleteTopicHandler)
		})

		r.Route("/answers", func(r chi.Router) {
			r.Get("/", h.ListAnswersHandler)
			r.Post("/", h.CreateAnswerHandler)
			r.Put("/", h.UpdateAnswerHandler)
			r.Delete("/", h.DeleteAnswerHandler)
		})

		r.Route("/tests", func(r chi.Router) {
			r.Get("/", h.ListTestsHandler)
			r.Post("/", h.CreateTestHandler)
			r.Put("/", h.UpdateTestHandler)
			r.Delete("/", h.DeleteTestHandler)
		})

		r.Route("/quiz-questions", func(r chi.Router) {
			r.Get("/", h.GetQuizQuestionsHandler)
			r.Post("/", h.AddQuestionToQuizHandler)
			r.Delete("/", h.RemoveQuestionFromQuizHandler)
		})

		r.Route("/test-questions", func(r chi.Router) {
			r.Get("/", h.GetTestQuestionsHandler)
			r.Post("/", h.AddQuestionToTestHandler)
			r.Delete("/", h.RemoveQuestionFromTestHandler)
		})

		// Tag Management routes
		r.Route("/tags", func(r chi.Router) {
			r.Get("/", h.ListTagsHandler)
			r.Delete("/", h.DeleteTagHandler)
			r.Post("/merge", h.MergeTagsHandler)
			r.Get("/questions", h.GetQuestionsByTagHandler)
		})
	})

	// Instructor routes (also protected by AdminMiddleware)
	r.Route("/instructor", func(r chi.Router) {
		// Apply admin middleware to all routes in this group
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.AdminMiddleware(next.ServeHTTP)(w, r)
			})
		})

		// Student Analytics routes
		r.Get("/students", h.ListStudentsHandler)
		r.Get("/students/analytics", h.GetStudentAnalyticsHandler)
		r.Get("/students/quiz-history", h.GetStudentQuizHistoryHandler)
		r.Get("/students/weak-areas", h.GetStudentWeakAreasHandler)
		r.Get("/student-progress", h.GetStudentProgressHandler)
		r.Get("/cohort-analytics", h.GetCohortAnalyticsHandler)

		// Import/Export routes
		r.Get("/quizzes/export", h.ExportQuizHandler)
		r.Post("/quizzes/import", h.ImportQuizHandler)
		r.Get("/questions/export", h.ExportQuestionsHandler)
		r.Post("/questions/import", h.ImportQuestionsHandler)

		// Assigned Quiz routes
		r.Post("/assigned-quizzes/import", h.ImportAndCreateQuizFromJSONHandler)
		r.Post("/assigned-quizzes/assign", h.AssignQuizToStudentsHandler)
		r.Delete("/assigned-quizzes/unassign", h.UnassignQuizFromStudentHandler)
		r.Get("/assigned-quizzes/status", h.ListAssignedQuizzesStatusHandler)
		r.Get("/assigned-quizzes/status/quiz", h.GetAssignedQuizStatusHandler)

		// Topic Management routes
		r.Route("/topics", func(r chi.Router) {
			r.Post("/", h.CreateInstructorTopicHandler)
			r.Put("/", h.UpdateInstructorTopicHandler)
			r.Delete("/", h.DeleteInstructorTopicHandler)
			r.Put("/reorder", h.ReorderTopicsHandler)
		})

		// Test Management routes
		r.Post("/tests/{id}/publish", h.PublishTestHandler)
		r.Get("/tests/{id}/results", h.GetTestResultsInstructorHandler)
		r.Get("/tests/{id}/export", h.ExportTestResultsHandler)

		// Question Management routes
		r.Route("/questions", func(r chi.Router) {
			r.Get("/", h.ListInstructorQuestionsHandler)
			r.Post("/", h.CreateInstructorQuestionHandler)
			r.Put("/", h.UpdateInstructorQuestionHandler)
			r.Delete("/", h.DeactivateInstructorQuestionHandler)

			r.Get("/analytics", h.GetQuestionAnalyticsHandler)
			r.Get("/{id}/analytics", h.GetQuestionDetailedAnalyticsHandler)
			r.Get("/{id}/answer-distribution", h.GetAnswerDistributionHandler)
			r.Put("/{id}/difficulty", h.UpdateQuestionDifficultyHandler)
			r.Get("/{id}/preview", h.PreviewInstructorQuestionHandler)
		})
	})

	return r
}
