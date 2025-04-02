package database

import (
	"context"
	"database/sql"
	"errors"
	"go-users/internal/api"
	"go-users/internal/ownErrors"
	"testing"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPool struct {
	mock.Mock
}

func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgx.Row)
}

func (m *MockPool) Close() {}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...any) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	fixedTime := time.Date(2023, 11, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		prepare     func(*MockPool)
		user        *api.UserRequest
		expected    *api.User
		expectedErr string
	}{
		{
			name: "Database error",
			prepare: func(mp *MockPool) {
				mr := new(MockRow)
				mr.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("db error"))

				mp.On("QueryRow", context.Background(),
					"INSERT INTO users (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING id, first_name, last_name, email, created_at, updated_at",
					[]any{"John", "Doe", openapi_types.Email("john@example.com")},
				).Return(mr)
			},
			user: &api.UserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     openapi_types.Email("john@example.com"),
			},
			expectedErr: "failed to create user: db error",
		},
		{
			name: "Success",
			prepare: func(mp *MockPool) {
				mr := new(MockRow)
				mr.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						*args.Get(0).(*uint) = 1
						*args.Get(1).(*string) = "John"
						*args.Get(2).(*string) = "Doe"
						*args.Get(3).(*openapi_types.Email) = openapi_types.Email("john@example.com")
						*args.Get(4).(*time.Time) = fixedTime
						*args.Get(5).(*time.Time) = fixedTime
					}).Return(nil)

				mp.On("QueryRow", context.Background(),
					"INSERT INTO users (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING id, first_name, last_name, email, created_at, updated_at",
					[]any{"John", "Doe", openapi_types.Email("john@example.com")},
				).Return(mr)
			},
			user: &api.UserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     openapi_types.Email("john@example.com"),
			},
			expected: &api.User{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     openapi_types.Email("john@example.com"),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := new(MockPool)
			db := &db{pool: mp}

			tt.prepare(mp)

			result, err := db.CreateUser(context.Background(), tt.user)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			mp.AssertExpectations(t)
		})
	}
}

func TestGetUser(t *testing.T) {
	fixedTime := time.Date(2023, 11, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		id          uint
		prepare     func(*MockPool)
		expected    *api.User
		expectedErr error
	}{
		{
			name: "User found",
			id:   1,
			prepare: func(mp *MockPool) {
				mr := new(MockRow)
				mr.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						*args.Get(0).(*uint) = 1
						*args.Get(1).(*string) = "John"
						*args.Get(2).(*string) = "Doe"
						*args.Get(3).(*openapi_types.Email) = openapi_types.Email("john@example.com")
						*args.Get(4).(*time.Time) = fixedTime
						*args.Get(5).(*time.Time) = fixedTime
					}).Return(nil)

				mp.On("QueryRow", context.Background(),
					"SELECT id, first_name, last_name, email, created_at, updated_at FROM users WHERE id = $1",
					[]any{uint(1)},
				).Return(mr)
			},
			expected: &api.User{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     openapi_types.Email("john@example.com"),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
		},
		{
			name: "User not found",
			id:   2,
			prepare: func(mp *MockPool) {
				mr := new(MockRow)
				mr.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(sql.ErrNoRows)

				mp.On("QueryRow", context.Background(),
					"SELECT id, first_name, last_name, email, created_at, updated_at FROM users WHERE id = $1",
					[]any{uint(2)},
				).Return(mr)
			},
			expectedErr: ownErrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := new(MockPool)
			db := &db{pool: mp}

			tt.prepare(mp)

			result, err := db.GetUser(context.Background(), tt.id)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			mp.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	fixedTime := time.Date(2023, 11, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		id          uint
		user        *api.UserRequest
		prepare     func(*MockPool)
		expected    *api.User
		expectedErr error
	}{
		{
			name: "Success",
			id:   1,
			user: &api.UserRequest{
				FirstName: "John",
				LastName:  "Doe Updated",
				Email:     openapi_types.Email("john.updated@example.com"),
			},
			prepare: func(mp *MockPool) {
				mr := new(MockRow)
				mr.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						*args.Get(0).(*uint) = 1
						*args.Get(1).(*string) = "John"
						*args.Get(2).(*string) = "Doe Updated"
						*args.Get(3).(*openapi_types.Email) = openapi_types.Email("john.updated@example.com")
						*args.Get(4).(*time.Time) = fixedTime
						*args.Get(5).(*time.Time) = fixedTime.Add(1 * time.Hour)
					}).Return(nil)

				mp.On("QueryRow", context.Background(),
					"UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4 RETURNING id, first_name, last_name, email, created_at, updated_at",
					[]any{"John", "Doe Updated", openapi_types.Email("john.updated@example.com"), uint(1)},
				).Return(mr)
			},
			expected: &api.User{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe Updated",
				Email:     openapi_types.Email("john.updated@example.com"),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime.Add(1 * time.Hour),
			},
		},
		{
			name: "Not found",
			id:   2,
			user: &api.UserRequest{
				FirstName: "NonExistent",
				LastName:  "User",
				Email:     openapi_types.Email("nope@example.com"),
			},
			prepare: func(mp *MockPool) {
				mr := new(MockRow)
				mr.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(sql.ErrNoRows)

				mp.On("QueryRow", context.Background(),
					"UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4 RETURNING id, first_name, last_name, email, created_at, updated_at",
					[]any{"NonExistent", "User", openapi_types.Email("nope@example.com"), uint(2)},
				).Return(mr)
			},
			expectedErr: ownErrors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := new(MockPool)
			db := &db{pool: mp}

			tt.prepare(mp)

			result, err := db.UpdateUser(context.Background(), tt.user, tt.id)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			mp.AssertExpectations(t)
		})
	}
}
