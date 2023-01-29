package http

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	"net/http"
	"todo/user"
)

const salt = "salt"

func authMiddleware(userRepo user.Repository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}

			// Calculate SHA-256 hashes for the provided and expected  usernames and passwords.
			passwordHash := sha256.Sum256([]byte(fmt.Sprintf("%s%s", "salt", password)))

			u, err := userRepo.FindByUsername(r.Context(), username)
			if err == user.ErrNotFound {
				zap.S().With("username", username).Error("user not found")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if err != nil {
				zap.S().Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(u.Username)) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], u.HashedPassword[:]) == 1

			// If the username and password are correct, then call
			// the next handler in the chain. Make sure to return
			// afterwards, so that none of the code below is run.
			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), user.UserContextKey, *u)))
				return
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

		})
	}
}

func singup(userRepo user.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, APIErrorResponse{Error: "could not parse request"})
			return
		}
		if req.Username == "" || req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, APIErrorResponse{Error: "username and password are required"})
			return
		}

		hashedPassword := sha256.Sum256([]byte(fmt.Sprintf("%s%s", salt, req.Password)))
		usr := &user.User{
			Username:       req.Username,
			HashedPassword: hashedPassword[:],
		}

		usr, err := userRepo.Create(r.Context(), usr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			zap.S().Errorf("could not create user: %v", err)
			return
		}

		response := signupResponse{
			Username: usr.Username,
			ID:       usr.ID,
		}
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, response)
	}
}
