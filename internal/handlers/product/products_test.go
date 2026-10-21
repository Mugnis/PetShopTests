package product

import (
	"encoding/json"
	"errors"
	"go-pet-shop/internal/handlers/product/mocks"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage — он вернёт один продукт.
	ProductGetterMock := mocks.NewProducts(t)

	ProductGetterMock.On("GetAllProducts", mock.Anything).
		Return([]models.Product{{ID: 1, Name: "Dog Food"}}, nil).
		Once()

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), ProductGetterMock)

	// Вызываем метод GetAllProducts, который является http.HandlerFunc
	handler.GetAllProducts(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
func TestGetAllProducts_Error(t *testing.T) {
	// Мокаем storage — он будет возвращать ошибку
	ProductGetterMock := mocks.NewProducts(t)

	ProductGetterMock.On("GetAllProducts", mock.Anything).
		Return(nil, errors.New("DB Error")).
		Once()

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ProductGetterMock)
	handler.GetAllProducts(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Create Product
// =======================

func TestCreateProduct_Success(t *testing.T) {
	// TODO: Написать Unit-тест для создания продукта (200 OK)
	ProductCreatingMock := mocks.NewProducts(t)

	product := &models.Product{}
	product.Name = "Корм для собак"
	product.Price = 499.99
	product.Stock = 10

	ProductCreatingMock.On("CreateProduct", mock.Anything, *product).
		Return(1, nil).
		Once()

	requestBody := `{"name": "Корм для собак", "price": 499.99, "stock": 10}`

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ProductCreatingMock)
	handler.CreateProduct(w, req)

	// Ожидаем HTTP 200
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["status"] != "Product created successfully" {
		t.Errorf("unexpected status string: %v", response["status"])
	}

	if response["id"] != float64(1) {
		t.Errorf("unexpected id: %v", response["id"])
	}
}

func TestCreateProduct_BadRequest(t *testing.T) {
	// TODO: Написать Unit-тест для создания продукта с невалидным JSON (400 Bad Request)
	ProductCreatingMock := mocks.NewProducts(t)

	requestBody := `"na: "Корля соак", "price 499.99, "sto"0}`

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ProductCreatingMock)
	handler.CreateProduct(w, req)

	// Ожидаем HTTP 400
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["error"] != "Bad request" {
		t.Errorf("unexpected status string: %v", response["status"])
	}
}

func TestCreateProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для создания продукта при ошибке сервиса (500 Internal Server Error)
	ProductCreatingMock := mocks.NewProducts(t)

	product := &models.Product{}
	product.Name = "Корм для собак"
	product.Price = 499.99
	product.Stock = 10

	ProductCreatingMock.On("CreateProduct", mock.Anything, *product).
		Return(0, errors.New("DB Error")).
		Once()

	requestBody := `{"name": "Корм для собак", "price": 499.99, "stock": 10}`

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), ProductCreatingMock)
	handler.CreateProduct(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["error"] != "Internal server error" {
		t.Errorf("unexpected error: %v", response["status"])
	}
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) {
	// TODO: Написать Unit-тест для обновления продукта (200 OK)
	ProductCreatingMock := mocks.NewProducts(t)

	product := &models.Product{}
	product.ID = 1
	product.Name = "Корм для собак"
	product.Price = 499.99
	product.Stock = 10

	ProductCreatingMock.On("UpdateProduct", mock.Anything, *product).
		Return(nil).
		Once()

	requestBody := `{"id": 10, "name": "Корм для собак", "price": 499.99, "stock": 10}`

	req := httptest.NewRequest(http.MethodPost, "/products/1", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), ProductCreatingMock)

	router.Post("/products/{id}", handler.UpdateProduct)

	router.ServeHTTP(w, req)

	// Ожидаем HTTP 200
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["status"] != "Product updated successfully" {
		t.Errorf("unexpected status string: %v", response["status"])
	}

	if response["id"] != float64(1) {
		t.Errorf("unexpected id: %v", response["id"])
	}
}

func TestUpdateProduct_BadRequest(t *testing.T) {
	// TODO: Написать Unit-тест для обновления продукта с невалидным JSON (400 Bad Request)
	ProductCreatingMock := mocks.NewProducts(t)

	requestBody := `d": 10, "name": "Корм для собак", "ice": 499.99, "stoc: 10}`

	req := httptest.NewRequest(http.MethodPost, "/products/1", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), ProductCreatingMock)

	router.Post("/products/{id}", handler.UpdateProduct)

	router.ServeHTTP(w, req)

	// Ожидаем HTTP 400
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["error"] != "Bad request" {
		t.Errorf("unexpected error: %v", response["status"])
	}
}

func TestUpdateProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для обновления продукта при ошибке сервиса (500 Internal Server Error)
	ProductCreatingMock := mocks.NewProducts(t)

	product := &models.Product{}
	product.ID = 1
	product.Name = "Корм для собак"
	product.Price = 499.99
	product.Stock = 10

	ProductCreatingMock.On("UpdateProduct", mock.Anything, *product).
		Return(errors.New("DB Error")).
		Once()

	requestBody := `{"id": 10, "name": "Корм для собак", "price": 499.99, "stock": 10}`

	req := httptest.NewRequest(http.MethodPost, "/products/1", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), ProductCreatingMock)

	router.Post("/products/{id}", handler.UpdateProduct)

	router.ServeHTTP(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["error"] != "Internal server error" {
		t.Errorf("unexpected error: %v", response["status"])
	}
}

// =======================
// Delete Product
// =======================

func TestDeleteProduct_Success(t *testing.T) {
	// TODO: Написать Unit-тест для удаления продукта (200 OK)
	ProductDeleterMock := mocks.NewProducts(t)

	ProductDeleterMock.On("DeleteProduct", mock.Anything, 1).
		Return(nil).
		Once()

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), ProductDeleterMock)

	router.Delete("/products/{id}", handler.DeleteProduct)

	router.ServeHTTP(w, req)

	// Ожидаем HTTP 200
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["status"] != "Product deleted successfully" {
		t.Errorf("unexpected status string: %v", response["status"])
	}

	if response["id"] != float64(1) {
		t.Errorf("unexpected id: %v", response["id"])
	}
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	// TODO: Написать Unit-тест для удаления продукта с пустым id (400 Bad Request)
	ProductDeleterMock := mocks.NewProducts(t)

	req := httptest.NewRequest(http.MethodDelete, "/products/", nil)
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), ProductDeleterMock)

	router.Delete("/products/{id}", handler.DeleteProduct)
	router.Delete("/products/", handler.DeleteProduct)

	router.ServeHTTP(w, req)

	// Ожидаем HTTP 400
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["error"] != "Bad request" {
		t.Errorf("unexpected error: %v", response["status"])
	}
}

func TestDeleteProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для удаления продукта при ошибке сервиса (500 Internal Server Error)
	ProductDeleterMock := mocks.NewProducts(t)

	ProductDeleterMock.On("DeleteProduct", mock.Anything, 1).
		Return(errors.New("DB Error")).
		Once()

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), ProductDeleterMock)

	router.Delete("/products/{id}", handler.DeleteProduct)

	router.ServeHTTP(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["error"] != "Internal server error" {
		t.Errorf("unexpected error: %v", response["status"])
	}
}
