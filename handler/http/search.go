package http

import (
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"todo/search"
	"todo/task"
)

func searchTasks(searchService *search.Service, taskService *task.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			render.JSON(w, r, ListResponse{})
		}

		taskIds, err := searchService.Search(r.Context(), query)
		if err != nil {
			zap.S().With("error", err).Error("search failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(taskIds) == 0 {
			render.JSON(w, r, ListResponse{})
			return
		}

		pagination := r.Context().Value(PaginationCtxKey).(Pagination)
		includeAllStatuses := r.URL.Query().Get("include_statuses") == "all"
		opts := task.QueryOptions{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
		}
		if includeAllStatuses {
			opts.IncludeStatuses = []task.Status{task.CreatedStatus, task.ArchivedStatus, task.FinishedStatus}
		}
		tasks, err := taskService.Search(r.Context(), query, opts)

		render.JSON(w, r, ListResponse{
			Total:  int64(len(taskIds)),
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
			Count:  len(tasks),
			Data:   tasks,
		})
	}
}
