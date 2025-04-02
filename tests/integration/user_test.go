package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	defaultBaseURL = "http://localhost:8080/api/v1"
)

type TestClient struct {
	client  *http.Client
	baseURL string
}

func NewTestClient() *TestClient {
	return &TestClient{
		client:  &http.Client{},
		baseURL: defaultBaseURL,
	}
}

func TestCreateUser(t *testing.T) {
	client := NewTestClient()

	tests := []struct {
		name           string
		user           map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid user",
			user: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      fmt.Sprintf("john.doe%d@example.com", rand.Intn(1000)),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Empty email",
			user: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.user)
			assert.NoError(t, err)

			resp, err := client.client.Post(client.baseURL+"/users", "application/json", bytes.NewBuffer(body))
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == http.StatusCreated {
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)

				assert.Equal(t, tt.user["first_name"], result["first_name"])
				assert.Equal(t, tt.user["last_name"], result["last_name"])
				assert.Equal(t, tt.user["email"], result["email"])
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	client := NewTestClient()

	user := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
		"email":      fmt.Sprintf("john.doe%d@example.com", rand.Intn(1000)),
	}

	body, err := json.Marshal(user)
	assert.NoError(t, err)

	resp, err := client.client.Post(client.baseURL+"/users", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&createResult)
	assert.NoError(t, err)
	userID := createResult["id"]

	resp, err = client.client.Get(fmt.Sprintf("%s/users/%v", client.baseURL, userID))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, user["first_name"], result["first_name"])
	assert.Equal(t, user["last_name"], result["last_name"])
	assert.Equal(t, user["email"], result["email"])
}

func TestUpdateUser(t *testing.T) {
	client := NewTestClient()

	user := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
		"email":      fmt.Sprintf("john.doe%d@example.com", rand.Intn(1000)),
	}

	body, err := json.Marshal(user)
	assert.NoError(t, err)

	resp, err := client.client.Post(client.baseURL+"/users", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&createResult)
	assert.NoError(t, err)
	userID := createResult["id"]

	updatedUser := map[string]interface{}{
		"first_name": "John Updated",
		"last_name":  "Doe Updated",
		"email":      fmt.Sprintf("john.updated%d@example.com", rand.Intn(1000)),
	}

	body, err = json.Marshal(updatedUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/users/%v", client.baseURL, userID), bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, updatedUser["first_name"], result["first_name"])
	assert.Equal(t, updatedUser["last_name"], result["last_name"])
	assert.Equal(t, updatedUser["email"], result["email"])
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
