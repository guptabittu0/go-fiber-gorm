package handler_test

import (
	"go-fiber-gorm/internal/handler"
	"go-fiber-gorm/internal/model"
	"go-fiber-gorm/test"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of user service for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(req *model.CreateUserRequest) (*model.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *MockUserService) GetByID(id uint) (*model.UserResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *MockUserService) Update(id uint, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *MockUserService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) GetAll(page, limit int) ([]model.UserResponse, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.UserResponse), args.Get(1).(int64), args.Error(2)
}

func setupTestApp(userService *MockUserService) *fiber.App {
	app := fiber.New()

	// Create handler with mock service
	userHandler := handler.NewUserHandler(userService)

	// Setup routes
	api := app.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users")

	// Define routes
	users.Post("/", userHandler.Create)
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)
	users.Get("/", userHandler.GetAll)

	return app
}

func TestCreateUser(t *testing.T) {
	// Create mock service
	mockService := new(MockUserService)

	// Setup test response
	mockResponse := &model.UserResponse{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "user",
	}

	// Setup mock behavior
	mockService.On("Create", mock.Anything).Return(mockResponse, nil)

	// Setup test app
	app := setupTestApp(mockService)

	// Create test request
	req := test.TestRequest{
		Method: "POST",
		URL:    "/api/v1/users",
		Body: map[string]interface{}{
			"name":     "Test User",
			"email":    "test@example.com",
			"password": "password123",
		},
	}

	// Make request
	resp := test.MakeTestRequest(t, app, req)

	// Verify status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Parse and verify response
	var response map[string]interface{}
	test.ParseResponse(t, resp, &response)

	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])
	assert.Equal(t, "Test User", data["name"])
	assert.Equal(t, "test@example.com", data["email"])

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}
