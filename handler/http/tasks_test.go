package http

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"todo/search"
	"todo/task"
	"todo/user"
)

func testHTTPCall(method, url string, body io.Reader, username, password string) (string, int, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", 0, err
	}
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	return strings.TrimSpace(string(bodyBytes)), resp.StatusCode, nil

}

func Test_Generic(t *testing.T) {

	// start test server with mock db
	logger := zap.S()
	searchService := &search.Service{
		Repo: search.MockUserIndexRepository{
			FindFn: func(ctx context.Context, userID uint) (*search.UserIndex, error) {
				if userID != 42 {
					return nil, fmt.Errorf("user not found")
				}
				return &search.UserIndex{
					UserID: userID,
					Index:  search.Index{"task": []string{"1", "2"}},
				}, nil
			},
			UpdateFn: func(ctx context.Context, userIndex *search.UserIndex) error {
				if userIndex.UserID != 42 {
					return fmt.Errorf("user not found")
				}
				if reflect.DeepEqual(userIndex.Index, search.Index{"3": []string{"3"}, "task": []string{"3"}}) {
					return fmt.Errorf("unexpected search index")
				}
				return nil
			},
			CreateFn: nil,
		}}
	taskService := &task.Service{
		Repo: task.MockRepository{
			FindAllFn: func(ctx context.Context, options task.QueryOptions) ([]*task.Task, error) {
				return []*task.Task{
					{
						ID:     "1",
						UserID: 42,
						Title:  "task 1",
						Status: task.FinishedStatus,
					},
					{
						ID:     "2",
						UserID: 42,
						Title:  "task 2",
						Status: task.CreatedStatus,
					},
				}, nil
			},
			CountAllFn: func(ctx context.Context, options task.QueryOptions) (int64, error) {
				return 2, nil
			},
			FindByIDsFn: func(ctx context.Context, options task.QueryOptions) ([]*task.Task, error) {
				tasks := make([]*task.Task, 0, len(options.IDs))
				for _, id := range options.IDs {
					tasks = append(tasks, &task.Task{
						ID:     id,
						UserID: 42,
						Title:  "task " + id,
						Status: task.FinishedStatus,
					})
				}
				return tasks, nil
			},
			FindByIDFn: nil,
			CreateFn: func(ctx context.Context, userId uint, task *task.Task) (*task.Task, error) {
				if userId != 42 {
					return nil, fmt.Errorf("unexpected user id: %d", userId)
				}
				task.ID = "3"
				return task, nil
			},
			UpdateFn: nil,
			DeleteFn: nil,
		},
		SearchService: searchService,
	}

	userRepo := &user.MockRepository{
		CreateFn: func(ctx context.Context, user *user.User) (*user.User, error) {
			if user.Username != "rafa" {
				return nil, fmt.Errorf("unexpected username: %s", user.Username)
			}
			user.ID = uint(42)
			return user, nil
		},
		FindByUsernameFn: func(ctx context.Context, username string) (*user.User, error) {
			if username != "rafa" {
				return nil, fmt.Errorf("unexpected username: %s", username)
			}
			hashedPassword := sha256.Sum256([]byte("salttest"))
			return &user.User{ID: uint(42), Username: username, HashedPassword: hashedPassword[:]}, nil
		},
	}

	handler := NewHandler(logger, taskService, searchService, userRepo)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	t.Run("create user", func(t *testing.T) {
		buf := bytes.NewBufferString(`{"username":"rafa","password": "test"}`)
		resp, code, err := testHTTPCall("POST", srv.URL+"/v1/signup", buf, "", "")
		if err != nil {
			t.Fatal(err)
		}
		if code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", code)
		}
		wantResp := `{"id":42,"username":"rafa"}`
		if resp != wantResp {
			t.Fatalf("unexpected response: \n`%s`\n, want: `\n%s`", resp, wantResp)
		}
	})

	t.Run("fetch something without login", func(t *testing.T) {
		_, code, err := testHTTPCall("GET", srv.URL+"/v1/tasks", nil, "", "")
		if err != nil {
			t.Fatal(err)
		}
		if code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", code)
		}
	})

	t.Run("fetch tasks", func(t *testing.T) {
		resp, code, err := testHTTPCall("GET", srv.URL+"/v1/tasks", nil, "rafa", "test")
		if err != nil {
			t.Fatal(err)
		}
		if code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", code)
		}
		wantResp := `{"total":2,"count":2,"offset":0,"limit":10,"data":[{"id":"1","title":"task 1","description":"","status":"finished","user_id":42,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"id":"2","title":"task 2","description":"","status":"created","user_id":42,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]}`
		if resp != wantResp {
			t.Fatalf("unexpected response: `%s`", resp)
		}
	})

	t.Run("create task", func(t *testing.T) {
		buf := bytes.NewBufferString(`{"title":"task 3"}`)
		resp, code, err := testHTTPCall("POST", srv.URL+"/v1/tasks", buf, "rafa", "test")
		if err != nil {
			t.Fatal(err)
		}
		if code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", code)
		}
		wantResp := `{"id":"3","title":"task 3","description":"","status":"created","user_id":42,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`
		if resp != wantResp {
			t.Fatalf("unexpected response: `%s`", resp)
		}
	})

	t.Run("fetch tasks with search", func(t *testing.T) {
		resp, code, err := testHTTPCall("GET", srv.URL+"/v1/search?query=task", nil, "rafa", "test")
		if err != nil {
			t.Fatal(err)
		}
		if code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", code)
		}
		wantResp := `{"total":2,"count":2,"offset":0,"limit":10,"data":[{"id":"1","title":"task 1","description":"","status":"finished","user_id":42,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"id":"2","title":"task 2","description":"","status":"finished","user_id":42,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]}`
		if resp != wantResp {
			t.Fatalf("unexpected response: \n`%s`\nwant:\n`%s`", resp, wantResp)
		}
	})
}
