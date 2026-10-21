package product

import (
	"context"
	"encoding/json"
	"errors"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage — он вернёт один продукт.
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return []models.Product{
				{ID: 1, Name: "Dog Food"},
			}, nil
		},
	}

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), mock)

	// Вызываем метод GetAllProducts, который является http.HandlerFunc
	handler.GetAllProducts(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
func TestGetAllProducts_Error(t *testing.T) {
	// Мокаем storage — он будет возвращать ошибку
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return nil, errors.New("DB error")
		},
	}

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
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
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 1, nil
		},
	}

	requestBody := `{"name": "Корм для собак", "price": 499.99, "stock": 10}`

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
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
	mock := &ProductsMock{}

	requestBody := `"na: "Корля соак", "price 499.99, "sto"0}`

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
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
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 0, errors.New("DB Error")
		},
	}

	requestBody := `{"name": "Корм для собак", "price": 499.99, "stock": 10}`

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
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
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return nil
		},
	}

	requestBody := `{"id": 10, "name": "Корм для собак", "price": 499.99, "stock": 10}`

	req := httptest.NewRequest(http.MethodPost, "/products/1", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), mock)

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
	mock := &ProductsMock{}

	requestBody := `d": 10, "name": "Корм для собак", "ice": 499.99, "stoc: 10}`

	req := httptest.NewRequest(http.MethodPost, "/products/1", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), mock)

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
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return errors.New("DB Error")
		},
	}

	requestBody := `{"id": 10, "name": "Корм для собак", "price": 499.99, "stock": 10}`

	req := httptest.NewRequest(http.MethodPost, "/products/1", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), mock)

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
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), mock)

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
	mock := &ProductsMock{}

	req := httptest.NewRequest(http.MethodDelete, "/products/", nil)
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), mock)

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
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return errors.New("DB Error")
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	handler := New(slog.Default(), mock)

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
