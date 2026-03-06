package transport

import (
	"Goworkspace/Project/domain"
	"Goworkspace/Project/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func SetupTestRout() http.Handler {
	st := storage.NewMemoryStorage()
	svc := domain.NewService(st)
	return NewRouter(svc)
}

func TestIntegration_CreateSuccess(t *testing.T) {
	router := SetupTestRout()

	body := []byte(`{"name": "Alex"}`)
	req := httptest.NewRequest(http.MethodPost, "/item", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got: %d", rec.Code)
	}

	var resp ResponseResult

	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unexpected error json: %v", err)
	}

	if resp.Item == nil || resp.Item.Name != "Alex" {
		t.Fatalf("unexpected responce, got: %+v", resp)
	}

}

func TestIntegration_CreateBadRequest(t *testing.T) {
	router := SetupTestRout()

	body := []byte(`{"name": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/item", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected BadRequest, got: %d", rec.Code)
	}
}

func TestIntegration_GetNotFound(t *testing.T) {
	router := SetupTestRout()

	req := httptest.NewRequest(http.MethodGet, "/item/1", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expacted ErrorNotFound, got: %v", rec.Code)
	}
}

func TestIntegration_DeleteNotFound(t *testing.T) {
	router := SetupTestRout()

	req := httptest.NewRequest(http.MethodDelete, "/item/1", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expacted ErrorNotFound, got: %v", rec.Code)
	}
}

func TestIntegration_CRUD_Flow(t *testing.T) {
	router := SetupTestRout()

	// CREATE
	createReq := httptest.NewRequest(
		http.MethodPost,
		"/item",
		bytes.NewBuffer([]byte(`{"name":"Alex"}`)),
	)
	createReq.Header.Set("Content-Type", "application/json")

	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("create failed: %d", createRec.Code)
	}

	// GET
	getReq := httptest.NewRequest(http.MethodGet, "/item/1", nil)
	getRec := httptest.NewRecorder()

	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("get failed: %d", getRec.Code)
	}

	// DELETE
	delReq := httptest.NewRequest(http.MethodDelete, "/item/1", nil)
	delRec := httptest.NewRecorder()

	router.ServeHTTP(delRec, delReq)

	if delRec.Code != http.StatusOK {
		t.Fatalf("delete failed: %d", delRec.Code)
	}

	// GET after delete
	getAfterDel := httptest.NewRequest(http.MethodGet, "/item/1", nil)
	getAfterDelRec := httptest.NewRecorder()

	router.ServeHTTP(getAfterDelRec, getAfterDel)

	if getAfterDelRec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getAfterDelRec.Code)
	}
}

func doRequest(t *testing.T, router http.Handler, method, path string, body []byte, expectedCode int) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != expectedCode {
		t.Errorf("Expected code %d, got: %d", expectedCode, recorder.Code)
		return recorder
	}

	return recorder
}

func TestIntegration_CRUD_Flow_Concurency(t *testing.T) {
	router := SetupTestRout()
	const n = 10

	ids := make(map[int]int)
	var (
		wg sync.WaitGroup
		mu sync.RWMutex
	)

	// Create 10 items
	for i := 1; i <= n; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			body, err := json.Marshal(map[string]string{"name": fmt.Sprintf("input-%d", i)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			recorder := doRequest(t, router, http.MethodPost, "/item", body, http.StatusCreated)

			var response ResponseResult

			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Errorf("Unexpected error json: %v", err)
			}
			mu.Lock()
			ids[i] = response.Item.ID
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	// Get 10 items
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			mu.RLock()
			id := ids[i]
			mu.RUnlock()

			doRequest(t, router, http.MethodGet, fmt.Sprintf("/item/%d", id), nil, http.StatusOK)
		}(i)
	}
	wg.Wait()

	// Delete 10 items
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			mu.RLock()
			id := ids[i]
			mu.RUnlock()

			doRequest(t, router, http.MethodDelete, fmt.Sprintf("/item/%d", id), nil, http.StatusOK)

			mu.Lock()
			delete(ids, i)
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	if len(ids) != 0 {
		t.Fatalf("Expected value ids=0; got: %d", len(ids))
	}
}
