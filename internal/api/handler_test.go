package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"go-users/internal/config"
	"go-users/internal/ownErrors"
)

func TestNewHandler(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var db DB

	tests := []struct {
		name           string
		openAPICfg     config.OpenAPI
		setupFunc      func() string
		sendRequest    func(ts *httptest.Server) *http.Response
		expectedError  string
		expectedStatus int
	}{
		{
			name: "Successful handler creation with spec path",
			openAPICfg: config.OpenAPI{
				SpecPath:  "",
				APIPrefix: "/api",
			},
			setupFunc: func() string {
				tmpDir := t.TempDir()
				specFile := filepath.Join(tmpDir, "openapi.yaml")
				err := os.WriteFile(specFile, []byte("openapi: '3.0.0'"), 0644)
				require.NoError(t, err)
				return specFile
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "OpenAPI spec file not found",
			openAPICfg: config.OpenAPI{
				SpecPath:  "/nonexistent/path/openapi.yaml",
				APIPrefix: "/api",
			},
			setupFunc:     func() string { return "" },
			expectedError: "OpenAPI specification file not found",
		},
		{
			name: "Empty spec path behaves as file not found",
			openAPICfg: config.OpenAPI{
				SpecPath: "",
			},
			setupFunc: func() string {
				return ""
			},
			expectedError: "OpenAPI specification file not found",
		},
		{
			name: "Middleware and request validation error",
			openAPICfg: config.OpenAPI{
				SpecPath:  "",
				APIPrefix: "/api",
			},
			setupFunc: func() string {
				tmpDir := t.TempDir()
				specFile := filepath.Join(tmpDir, "openapi.yaml")
				err := os.WriteFile(specFile, []byte("openapi: '3.0.0'"), 0644)
				require.NoError(t, err)
				return specFile
			},
			sendRequest: func(ts *httptest.Server) *http.Response {
				client := &http.Client{}
				req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/users", nil)
				require.NoError(t, err)
				resp, err := client.Do(req)
				require.NoError(t, err)
				return resp
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "SpecPath is empty, defaults to executable path",
			openAPICfg: config.OpenAPI{
				SpecPath:  "",
				APIPrefix: "/api",
			},
			setupFunc: func() string {
				tmpDir := t.TempDir()
				execPath := filepath.Join(tmpDir, "app")
				err := os.WriteFile(execPath, []byte("#!/bin/bash\necho 'mock'"), 0755)
				require.NoError(t, err)

				specDir := filepath.Join(tmpDir, "api")
				err = os.Mkdir(specDir, 0755)
				require.NoError(t, err)

				specFile := filepath.Join(specDir, "openapi.yaml")
				err = os.WriteFile(specFile, []byte("openapi: '3.0.0'"), 0644)
				require.NoError(t, err)

				originalExecutable := os.Args[0]
				os.Args[0] = execPath
				t.Cleanup(func() {
					os.Args[0] = originalExecutable
				})

				return specFile
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Middleware and response error handler",
			openAPICfg: config.OpenAPI{
				SpecPath:  "",
				APIPrefix: "/api",
			},
			setupFunc: func() string {
				tmpDir := t.TempDir()
				specFile := filepath.Join(tmpDir, "openapi.yaml")
				err := os.WriteFile(specFile, []byte("openapi: '3.0.0'"), 0644)
				require.NoError(t, err)
				return specFile
			},
			sendRequest: func(ts *httptest.Server) *http.Response {
				resp, err := http.Get(ts.URL + "/nonexistent-endpoint")
				require.NoError(t, err)
				return resp
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.openAPICfg.SpecPath == "" && tc.setupFunc != nil {
				tc.openAPICfg.SpecPath = tc.setupFunc()
			}

			handler, err := NewHandler(tc.openAPICfg, logger, db)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, handler)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, handler)

				ts := httptest.NewServer(handler)
				defer ts.Close()

				if tc.sendRequest != nil {
					resp := tc.sendRequest(ts)
					assert.Equal(t, tc.expectedStatus, resp.StatusCode)
				}
			}
		})
	}
}

func TestUserHandler_Health(t *testing.T) {
	handler := &UserHandler{}
	resp, err := handler.Health(context.Background(), HealthRequestObject{})

	assert.NoError(t, err, "Health should not return an error")
	assert.NotNil(t, resp, "Health response should not be nil")

	healthResp, ok := resp.(Health200JSONResponse)
	assert.True(t, ok, "Response should be of type Health200JSONResponse")
	assert.Equal(t, "OK", *healthResp.Status, "Health status should be 'OK'")
}

// MockUserRepository is a mock implementation of a database.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Close() {
	panic("implement me")
}

func (m *MockUserRepository) GetUser(ctx context.Context, id uint) (*User, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *UserRequest) (*User, error) {
	args := m.Called(ctx, user)
	if result := args.Get(0); result != nil {
		return result.(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, u *UserRequest, id uint) (*User, error) {
	args := m.Called(ctx, u, id)
	if result := args.Get(0); result != nil {
		return result.(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestUserHandler_PostUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	fixedTime := time.Date(2020, time.November, 10, 12, 0, 0, 0, time.UTC)
	handler := &UserHandler{repo: mockRepo}

	testCases := []struct {
		name           string
		input          PostUserRequestObject
		mockResponse   *User
		mockError      error
		expectedOutput PostUserResponseObject
		expectedError  error
	}{
		{
			name: "Successful creation",
			input: PostUserRequestObject{
				Body: &PostUserJSONRequestBody{
					FirstName: "John",
					LastName:  "Doe",
					Email:     types.Email("john.doe@example.com"),
				},
			},
			mockResponse: &User{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     types.Email("john.doe@example.com"),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockError: nil,
			expectedOutput: PostUser201JSONResponse{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     types.Email("john.doe@example.com"),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expectedError: nil,
		},
		{
			name: "Missing request body",
			input: PostUserRequestObject{
				Body: nil,
			},
			expectedOutput: PostUser400JSONResponse{
				Error: stringPtr("Missing request body"),
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockResponse != nil {
				mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(tc.mockResponse, tc.mockError)
			}

			resp, err := handler.PostUser(context.Background(), tc.input)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.expectedOutput, resp)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	fixedTime := time.Date(2023, time.November, 10, 12, 0, 0, 0, time.UTC)
	handler := &UserHandler{repo: mockRepo}

	testCases := []struct {
		name           string
		inputID        uint
		mockResponse   *User
		mockError      error
		expectedOutput GetUserResponseObject
		expectedError  error
	}{
		{
			name:    "Successful retrieval",
			inputID: 1,
			mockResponse: &User{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockError: nil,
			expectedOutput: GetUser200JSONResponse{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     types.Email("john.doe@example.com"),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expectedError: nil,
		},
		{
			name:           "User not found",
			inputID:        2,
			mockResponse:   nil,
			mockError:      ownErrors.ErrNotFound,
			expectedOutput: GetUser404JSONResponse{Error: stringPtr("User not found")},
			expectedError:  nil,
		},
		{
			name:           "Internal server error",
			inputID:        3,
			mockResponse:   nil,
			mockError:      errors.New("database connection error"),
			expectedOutput: GetUser500JSONResponse{Error: stringPtr("Internal server error")},
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.On("GetUser", mock.Anything, tc.inputID).Return(tc.mockResponse, tc.mockError)
			resp, err := handler.GetUser(context.Background(), GetUserRequestObject{Id: tc.inputID})

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, resp)
			mockRepo.AssertCalled(t, "GetUser", mock.Anything, tc.inputID)
		})
	}
}

func TestUserHandler_PutUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	fixedTime := time.Date(2023, time.November, 10, 12, 0, 0, 0, time.UTC)
	handler := &UserHandler{repo: mockRepo}

	testCases := []struct {
		name           string
		inputID        uint
		inputBody      *PutUserJSONRequestBody
		mockResponse   *User
		mockError      error
		expectedOutput PutUserResponseObject
		expectedError  error
	}{
		{
			name:    "Successful update",
			inputID: 1,
			inputBody: &PutUserJSONRequestBody{
				FirstName: "John",
				LastName:  "Doe",
				Email:     types.Email("john.doe@example.com"),
			},
			mockResponse: &User{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     types.Email("john.doe@example.com"),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockError: nil,
			expectedOutput: PutUser200JSONResponse{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     types.Email("john.doe@example.com"),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expectedError: nil,
		},
		{
			name:           "Missing request body",
			inputID:        1,
			inputBody:      nil,
			mockResponse:   nil,
			mockError:      nil,
			expectedOutput: PutUser400JSONResponse{Error: stringPtr("Missing request body")},
			expectedError:  nil,
		},
		{
			name:    "User not found",
			inputID: 2,
			inputBody: &PutUserJSONRequestBody{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     types.Email("jane.doe@example.com"),
			},
			mockResponse:   nil,
			mockError:      ownErrors.ErrNotFound,
			expectedOutput: PutUser404JSONResponse{Error: stringPtr("User not found")},
			expectedError:  nil,
		},
		{
			name:    "Repository error",
			inputID: 3,
			inputBody: &PutUserJSONRequestBody{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     types.Email("jane.doe@example.com"),
			},
			mockResponse:   nil,
			mockError:      errors.New("database update failed"),
			expectedOutput: PutUser500JSONResponse{Error: stringPtr("Internal server error")},
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.inputBody != nil {
				mockRepo.On("UpdateUser", mock.Anything, mock.Anything, tc.inputID).Return(tc.mockResponse, tc.mockError)
			}

			resp, err := handler.PutUser(context.Background(), PutUserRequestObject{
				Id:   tc.inputID,
				Body: tc.inputBody,
			})

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.expectedOutput, resp)

			if tc.inputBody != nil {
				mockRepo.AssertCalled(t, "UpdateUser", mock.Anything, mock.Anything, tc.inputID)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
