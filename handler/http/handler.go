package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
	"todo/search"
	"todo/task"
	"todo/user"
)

// NewHandler return a new router with some handy middleware and api routes
func NewHandler(log *zap.SugaredLogger, taskService *task.Service, searchService *search.Service, userRepo user.Repository) chi.Router {
	r := chi.NewRouter()

	r.Use(
		traceMiddleware(),
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.Timeout(60*time.Second),
	)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/v1/signup", singup(userRepo))

	r.With(authMiddleware(userRepo), profilingMiddleware(log)).Route("/v1", func(r chi.Router) {
		r.Route("/tasks", func(r chi.Router) {
			r.With(paginationMiddleware()).Get("/", getTasks(taskService))
			r.Post("/", createTask(taskService))
			r.With(taskMiddleware(taskService)).Get("/{id}", getTask(taskService))
			r.With(taskMiddleware(taskService)).Patch("/{id}", updateTask(taskService))
			r.With(taskMiddleware(taskService)).Post("/{id}/complete", completeTask(taskService))
			r.With(taskMiddleware(taskService)).Delete("/{id}", deleteTask(taskService))
		})
		r.Route("/search", func(r chi.Router) {
			r.With(paginationMiddleware()).Get("/", searchTasks(searchService, taskService))
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	return r
}
