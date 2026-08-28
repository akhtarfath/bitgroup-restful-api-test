package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/akhtarfath/config"
	"github.com/gofiber/fiber/v3"
)

func newTestApp(t *testing.T) *fiber.App {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}

	if err := config.Load(); err != nil {
		t.Fatalf("config.Load: %v", err)
	}

	app := fiber.New()
	New(app)
	return app
}

func doJSON(t *testing.T, app *fiber.App, method, url, token string, body any) (int, map[string]any) {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}

	req := httptest.NewRequest(method, url, &buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	defer resp.Body.Close()

	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("%s %s: decode response: %v", method, url, err)
	}

	return resp.StatusCode, out
}

func TestUnauthorizedRequestsAreRejected(t *testing.T) {
	app := newTestApp(t)

	for _, url := range []string{"/products", "/products/1", "/categories", "/categories/1"} {
		code, body := doJSON(t, app, "GET", url, "", nil)
		if code != http.StatusUnauthorized {
			t.Errorf("%s without token: want %d, got %d", url, http.StatusUnauthorized, code)
		}
		if body["status"] != "error" {
			t.Errorf("%s without token: want error status, got %v", url, body["status"])
		}
	}
}

func TestLoginWithInvalidCredentials(t *testing.T) {
	app := newTestApp(t)

	code, body := doJSON(t, app, "POST", "/login", "", map[string]string{
		"username": "admin",
		"password": "wrong",
	})
	if code != http.StatusUnauthorized {
		t.Errorf("login with wrong password: want %d, got %d", http.StatusUnauthorized, code)
	}
	if body["message"] != "Invalid username or password" {
		t.Errorf("unexpected message: %v", body["message"])
	}
}

func TestAuthorizedCRUDFlow(t *testing.T) {
	app := newTestApp(t)

	// Login and capture the bearer token
	code, body := doJSON(t, app, "POST", "/login", "", map[string]string{
		"username": "admin",
		"password": "admin",
	})
	if code != http.StatusOK {
		t.Fatalf("login: want %d, got %d (%v)", http.StatusOK, code, body["message"])
	}
	data, _ := body["data"].(map[string]any)
	token, _ := data["token"].(string)
	if token == "" {
		t.Fatal("login did not return a token")
	}

	// Seed list present and accessible
	code, body = doJSON(t, app, "GET", "/categories", token, nil)
	if code != http.StatusOK {
		t.Fatalf("list categories: want %d, got %d", http.StatusOK, code)
	}

	// Create a category
	code, body = doJSON(t, app, "POST", "/categories", token, map[string]string{"name": "Books"})
	if code != http.StatusCreated {
		t.Fatalf("create category: want %d, got %d (%v)", http.StatusCreated, code, body["message"])
	}
	catData, _ := body["data"].(map[string]any)
	catID, _ := catData["id"].(string)
	if catID == "" {
		t.Fatal("create category did not return an id")
	}

	// Create a category referencing it
	code, body = doJSON(t, app, "POST", "/products", token, map[string]any{
		"name":        "Clean Code",
		"price":       35,
		"category_id": catID,
	})
	if code != http.StatusCreated {
		t.Fatalf("create product: want %d, got %d (%v)", http.StatusCreated, code, body["message"])
	}
	prodData, _ := body["data"].(map[string]any)
	prodID, _ := prodData["id"].(string)
	if prodID == "" {
		t.Fatal("create product did not return an id")
	}

	// Create a product with a missing category is rejected
	code, _ = doJSON(t, app, "POST", "/products", token, map[string]any{
		"name":        "Orphan",
		"price":       1,
		"category_id": "does-not-exist",
	})
	if code != http.StatusBadRequest {
		t.Errorf("product with bad category: want %d, got %d", http.StatusBadRequest, code)
	}

	// Get product by id
	code, _ = doJSON(t, app, "GET", "/products/"+prodID, token, nil)
	if code != http.StatusOK {
		t.Errorf("get product: want %d, got %d", http.StatusOK, code)
	}

	// Update product
	code, _ = doJSON(t, app, "PUT", "/products/"+prodID, token, map[string]any{
		"name":        "Clean Code 2nd Ed",
		"price":       40,
		"category_id": catID,
	})
	if code != http.StatusOK {
		t.Errorf("update product: want %d, got %d", http.StatusOK, code)
	}

	// Delete product, then it must 404
	code, _ = doJSON(t, app, "DELETE", "/products/"+prodID, token, nil)
	if code != http.StatusOK {
		t.Errorf("delete product: want %d, got %d", http.StatusOK, code)
	}
	code, _ = doJSON(t, app, "GET", "/products/"+prodID, token, nil)
	if code != http.StatusNotFound {
		t.Errorf("get deleted product: want %d, got %d", http.StatusNotFound, code)
	}
}
