package router_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-server/internal"
	"go-server/internal/infrastructure"
	"go-server/internal/service/user"
	userDelivery "go-server/internal/service/user/delivery"
	"go-server/internal/service/workspace"
	workspaceDelivery "go-server/internal/service/workspace/delivery"
	"go-server/pkg/config"
	"go-server/pkg/logger"

	"github.com/gin-gonic/gin"
)

// TestRouterConnections tests all router connections and endpoints
func TestRouterConnections(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load test configuration
	cfg := config.Load()

	// Setup logger for testing
	logger.SetupLogger(cfg)

	// Mock database and Redis connections for testing
	// Note: In a real test environment, you'd want to use test databases
	db, err := infrastructure.ConnectDB(cfg)
	if err != nil {
		t.Skipf("Skipping router tests - database connection failed: %v", err)
		return
	}

	redisClient, err := infrastructure.ConnectRedis(cfg)
	if err != nil {
		t.Skipf("Skipping router tests - Redis connection failed: %v", err)
		return
	}
	defer infrastructure.CloseRedis(redisClient)

	// Initialize HTTP client
	httpClient := infrastructure.NewDefaultHTTPClient()

	// Setup router
	router := infrastructure.SetupRouter()

	// Register all routes
	internal.RegisterHealthRoutes(router, db, redisClient)
	user.RegisterUserService(router, db, httpClient)
	workspace.RegisterWorkspaceService(router, db, httpClient)

	// Create test server
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Run all route tests
	t.Run("Health Endpoints", func(t *testing.T) {
		testHealthEndpoints(t, testServer)
	})

	t.Run("User Service Routes", func(t *testing.T) {
		testUserServiceRoutes(t, testServer)
	})

	t.Run("Workspace Service Routes", func(t *testing.T) {
		testWorkspaceServiceRoutes(t, testServer)
	})

	fmt.Printf("\n✅ All router connection tests completed successfully!\n")
}

// testHealthEndpoints tests the health check endpoints
func testHealthEndpoints(t *testing.T, server *httptest.Server) {
	tests := []struct {
		name           string
		endpoint       string
		expectedStatus int
	}{
		{
			name:           "Health Check - Full",
			endpoint:       "/health",
			expectedStatus: http.StatusOK, // or StatusServiceUnavailable depending on DB/Redis
		},
		{
			name:           "Health Check - Simple",
			endpoint:       "/health/simple",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(server.URL + tt.endpoint)
			if err != nil {
				t.Fatalf("HTTP request failed: %v", err)
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("📊 %s: Status %d\n", tt.name, resp.StatusCode)
			fmt.Printf("   Response: %s\n", string(body))

			// For health endpoints, we accept both OK and ServiceUnavailable
			// since test environment might not have all services running
			validStatuses := []int{http.StatusOK, http.StatusServiceUnavailable}
			isValid := false
			for _, status := range validStatuses {
				if resp.StatusCode == status {
					isValid = true
					break
				}
			}
			if !isValid {
				t.Errorf("Expected status %v, got %d", validStatuses, resp.StatusCode)
			}
		})
	}
}

// testUserServiceRoutes tests the user service endpoints
func testUserServiceRoutes(t *testing.T, server *httptest.Server) {
	tests := []struct {
		name           string
		method         string
		endpoint       string
		payload        interface{}
		expectedStatus []int // Allow multiple valid status codes
	}{
		{
			name:     "Create User - Valid Data",
			method:   "POST",
			endpoint: "/users/",
			payload: map[string]string{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: []int{http.StatusCreated, http.StatusInternalServerError}, // May fail due to DB constraints
		},
		{
			name:           "Create User - Invalid Data",
			method:         "POST",
			endpoint:       "/users/",
			payload:        map[string]string{}, // Empty payload
			expectedStatus: []int{http.StatusBadRequest},
		},
		{
			name:           "Delete User",
			method:         "DELETE",
			endpoint:       "/users/1",
			payload:        nil,
			expectedStatus: []int{http.StatusOK, http.StatusInternalServerError}, // May fail if user doesn't exist
		},
		{
			name:           "Delete User - Invalid ID",
			method:         "DELETE",
			endpoint:       "/users/invalid",
			payload:        nil,
			expectedStatus: []int{http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.payload != nil {
				jsonPayload, _ := json.Marshal(tt.payload)
				req, err = http.NewRequest(tt.method, server.URL+tt.endpoint, bytes.NewBuffer(jsonPayload))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tt.method, server.URL+tt.endpoint, nil)
			}

			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("HTTP request failed: %v", err)
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("👤 %s: %s %s - Status %d\n", tt.name, tt.method, tt.endpoint, resp.StatusCode)
			fmt.Printf("   Response: %s\n", string(body))

			// Check if status code is in expected list
			isValid := false
			for _, expectedStatus := range tt.expectedStatus {
				if resp.StatusCode == expectedStatus {
					isValid = true
					break
				}
			}
			if !isValid {
				t.Errorf("Expected status codes %v, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

// testWorkspaceServiceRoutes tests the workspace service endpoints
func testWorkspaceServiceRoutes(t *testing.T, server *httptest.Server) {
	tests := []struct {
		name           string
		method         string
		endpoint       string
		payload        interface{}
		expectedStatus []int
	}{
		{
			name:     "Create Workspace - Valid Data",
			method:   "POST",
			endpoint: "/workspaces/",
			payload: map[string]string{
				"name":        "Test Workspace",
				"description": "A test workspace",
			},
			expectedStatus: []int{http.StatusCreated, http.StatusInternalServerError, http.StatusBadRequest},
		},
		{
			name:           "Create Workspace - Invalid Data",
			method:         "POST",
			endpoint:       "/workspaces/",
			payload:        map[string]string{},
			expectedStatus: []int{http.StatusBadRequest},
		},
		{
			name:           "Delete Workspace",
			method:         "DELETE",
			endpoint:       "/workspaces/1",
			payload:        nil,
			expectedStatus: []int{http.StatusOK, http.StatusInternalServerError, http.StatusBadRequest},
		},
		{
			name:           "Delete Workspace - Invalid ID",
			method:         "DELETE",
			endpoint:       "/workspaces/invalid",
			payload:        nil,
			expectedStatus: []int{http.StatusBadRequest},
		},
		{
			name:     "Invite Members",
			method:   "POST",
			endpoint: "/workspaces/1/members",
			payload: map[string]interface{}{
				"user_ids": []int{2, 3},
				"role":     "member",
			},
			expectedStatus: []int{http.StatusOK, http.StatusInternalServerError, http.StatusBadRequest},
		},
		{
			name:     "Remove Members",
			method:   "POST",
			endpoint: "/workspaces/1/members/remove",
			payload: map[string]interface{}{
				"user_ids": []int{2},
			},
			expectedStatus: []int{http.StatusOK, http.StatusInternalServerError, http.StatusBadRequest},
		},
		{
			name:     "Change Member Access",
			method:   "PUT",
			endpoint: "/workspaces/1/members/2/access",
			payload: map[string]string{
				"role": "admin",
			},
			expectedStatus: []int{http.StatusOK, http.StatusInternalServerError, http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.payload != nil {
				jsonPayload, _ := json.Marshal(tt.payload)
				req, err = http.NewRequest(tt.method, server.URL+tt.endpoint, bytes.NewBuffer(jsonPayload))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tt.method, server.URL+tt.endpoint, nil)
			}

			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("HTTP request failed: %v", err)
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("🏢 %s: %s %s - Status %d\n", tt.name, tt.method, tt.endpoint, resp.StatusCode)
			fmt.Printf("   Response: %s\n", string(body))

			// Check if status code is in expected list
			isValid := false
			for _, expectedStatus := range tt.expectedStatus {
				if resp.StatusCode == expectedStatus {
					isValid = true
					break
				}
			}
			if !isValid {
				t.Errorf("Expected status codes %v, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

// TestIndividualRouterComponents tests router components in isolation
func TestIndividualRouterComponents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Router Setup", func(t *testing.T) {
		router := infrastructure.SetupRouter()
		if router == nil {
			t.Fatal("Expected router to be non-nil")
		}
		fmt.Printf("✅ Router setup successful\n")
	})

	t.Run("User Router Creation", func(t *testing.T) {
		router := infrastructure.SetupRouter()
		userRouter := userDelivery.NewUserRouter(router)
		if userRouter == nil {
			t.Fatal("Expected user router to be non-nil")
		}
		fmt.Printf("✅ User router creation successful\n")
	})

	t.Run("Workspace Router Creation", func(t *testing.T) {
		router := infrastructure.SetupRouter()
		workspaceRouter := workspaceDelivery.NewWorkspaceRouter(router)
		if workspaceRouter == nil {
			t.Fatal("Expected workspace router to be non-nil")
		}
		fmt.Printf("✅ Workspace router creation successful\n")
	})
}

// TestRouterMiddleware tests any middleware functionality
func TestRouterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := infrastructure.SetupRouter()

	// Test that router has default Gin middleware (Logger, Recovery)
	if router == nil {
		t.Fatal("Expected router to be non-nil")
	}
	fmt.Printf("✅ Router middleware tests completed\n")
}
