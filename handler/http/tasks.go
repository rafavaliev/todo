package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"todo/task"
	"todo/user"
)

func taskMiddleware(taskService *task.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			t, err := taskService.FindByID(r.Context(), id)
			switch {
			case err == nil:
				break
			case errors.Is(err, task.ErrNotFound):
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				zap.S().With("error", err).Error("task middleware failed")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), task.TaskContextKey, t)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getTasks(service *task.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := r.Context().Value(PaginationCtxKey).(Pagination)
		includeAllStatuses := r.URL.Query().Get("include_statuses") == "all"

		opts := task.QueryOptions{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
		}
		if includeAllStatuses {
			opts.IncludeStatuses = []task.Status{task.CreatedStatus, task.ArchivedStatus, task.FinishedStatus}
		}
		tasks, err := service.FindAll(r.Context(), opts)
		if err != nil {
			zap.S().With("error", err).Error("fetch tasks failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		total, err := service.CountAll(r.Context(), opts)
		if err != nil {
			zap.S().With("error", err).Error("count tasks failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := ListResponse{
			Total:  total,
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
			Count:  len(tasks),
			Data:   tasks,
		}
		render.JSON(w, r, resp)
	}
}

func getTask(service *task.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := r.Context().Value(task.TaskContextKey).(*task.Task)
		render.JSON(w, r, t)
	}
}

func deleteTask(service *task.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := service.Delete(r.Context(), id)
		switch {
		case err == nil:
			break
		case errors.Is(err, task.ErrNotFound):
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			zap.S().With("error", err).Error("delete task failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func createTask(service *task.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		usr := r.Context().Value(user.UserContextKey).(user.User)

		newTask := task.Task{
			ID:          uuid.New().String(),
			Title:       req.Title,
			Description: req.Description,
			Status:      task.CreatedStatus,
			UserID:      usr.ID,
		}

		t, err := service.Create(r.Context(), &newTask)
		switch {
		case err == nil:
			break
		case errors.Is(err, task.ErrEmptyTitle) || errors.Is(err, task.ErrInvalidStatus):
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, APIErrorResponse{Error: err.Error()})
			return
		default:
			zap.S().With("error", err).Error("create task failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, t)
	}
}

func updateTask(service *task.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var req updateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, APIErrorResponse{Error: "invalid request json body"})
			return
		}
		if req.Title == nil && req.Description == nil && req.Status == nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, APIErrorResponse{Error: "at least one field for update must be provided"})
			return
		}

		updatedTask := task.UpdateTask{
			ID:          id,
			Title:       req.Title,
			Description: req.Description,
		}
		if req.Status != nil {
			switch *req.Status {
			case string(task.CreatedStatus):
				updatedTask.Status = &task.CreatedStatus
			case string(task.FinishedStatus):
				updatedTask.Status = &task.FinishedStatus
			case string(task.ArchivedStatus):
				updatedTask.Status = &task.ArchivedStatus
			}
		}

		t, err := service.Update(r.Context(), &updatedTask)
		switch {
		case err == nil:
			break
		case errors.Is(err, task.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, task.ErrInvalidStatus) || errors.Is(err, task.ErrEmptyTitle):
			render.JSON(w, r, APIErrorResponse{Error: err.Error()})
			w.WriteHeader(http.StatusBadRequest)
		default:
			zap.S().With("error", err).Error("update task failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, t)
	}
}

func completeTask(service *task.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		updatedTask := &task.UpdateTask{
			ID:     id,
			Status: &task.FinishedStatus,
		}
		t, err := service.Update(r.Context(), updatedTask)
		if err != nil {
			zap.S().With("error", err).Error("complete task failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, t)
	}
}
