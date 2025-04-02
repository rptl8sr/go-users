package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"go-users/internal/config"
	"go-users/internal/ownErrors"
	"go-users/internal/router"
)

type DB interface {
	CreateUser(ctx context.Context, user *UserRequest) (*User, error)
	GetUser(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, u *UserRequest, id uint) (*User, error)
	Close()
}

// _ ensures that UserHandler implements the StrictServerInterface at compile time.
var _ StrictServerInterface = (*UserHandler)(nil)

// UserHandler implements the StrictServerInterface
type UserHandler struct {
	repo DB
}

// NewHandler creates a new HTTP handler
func NewHandler(openAPICfg config.OpenAPI, logger *slog.Logger, repo DB) (http.Handler, error) {
	specPath := openAPICfg.SpecPath
	if specPath == "" {
		execPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %w", err)
		}
		specPath = filepath.Join(filepath.Dir(execPath), OpenAPISpecPath)
	}

	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("OpenAPI specification file not found at %s", specPath)
	}

	r := router.New(router.Options{
		Logger: logger,
	})

	handler := &UserHandler{repo: repo}

	RegisterSwaggerRoutes(r)

	r.Route(openAPICfg.APIPrefix, func(r chi.Router) {
		ownStrictHandler := NewStrictHandlerWithOptions(handler, nil, StrictHTTPServerOptions{
			RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				logger.Error("request validation error", "error", err)
				http.Error(w, "Invalid request", http.StatusBadRequest)
			},
			ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				logger.Error("response processing error", "error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			},
		})
		HandlerFromMux(ownStrictHandler, r)
	})

	return r, nil
}

// Health implements the service's health endpoint
func (h *UserHandler) Health(_ context.Context, _ HealthRequestObject) (HealthResponseObject, error) {
	status := "OK"
	return Health200JSONResponse{
		Status: &status,
	}, nil
}

// PostUser creates a new user
func (h *UserHandler) PostUser(ctx context.Context, request PostUserRequestObject) (PostUserResponseObject, error) {
	if request.Body == nil {
		errorMsg := "Missing request body"
		return PostUser400JSONResponse{
			Error: &errorMsg,
		}, nil
	}

	user, err := h.repo.CreateUser(ctx, request.Body)
	if err != nil {
		if errors.Is(err, ownErrors.ErrUserAlreadyExists) {
			errorMsg := "User already exists"
			return PostUser409JSONResponse{
				Error: &errorMsg,
			}, nil
		}
		log.Printf("error: %v", err)
		errorMsg := err.Error()
		return PostUser500JSONResponse{
			Error: &errorMsg,
		}, nil
	}

	return PostUser201JSONResponse(*user), nil
}

// GetUser fetches a user by their ID
func (h *UserHandler) GetUser(ctx context.Context, request GetUserRequestObject) (GetUserResponseObject, error) {
	user, err := h.repo.GetUser(ctx, request.Id)
	if err != nil {
		if errors.Is(err, ownErrors.ErrNotFound) {
			errorMsg := "User not found"
			return GetUser404JSONResponse{
				Error: &errorMsg,
			}, nil
		}

		errorMsg := "Internal server error"
		return GetUser500JSONResponse{
			Error: &errorMsg,
		}, nil
	}

	return GetUser200JSONResponse(*user), nil
}

// PutUser updates user entity with replacing user data
func (h *UserHandler) PutUser(ctx context.Context, request PutUserRequestObject) (PutUserResponseObject, error) {
	if request.Body == nil {
		errorMsg := "Missing request body"
		return PutUser400JSONResponse{
			Error: &errorMsg,
		}, nil
	}

	user, err := h.repo.UpdateUser(ctx, request.Body, request.Id)
	if err != nil {
		if errors.Is(err, ownErrors.ErrNotFound) {
			errorMsg := "User not found"
			return PutUser404JSONResponse{
				Error: &errorMsg,
			}, nil
		}
		errorMsg := "Internal server error"
		return PutUser500JSONResponse{
			Error: &errorMsg,
		}, nil
	}

	return PutUser200JSONResponse(*user), nil
}
