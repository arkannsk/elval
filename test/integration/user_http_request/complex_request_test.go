package user_http_request

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComplexRequest_ParseRequest_FullSuccessWithRouter(t *testing.T) {
	mux := http.NewServeMux()
	var capturedReq *ComplexRequest

	mux.HandleFunc("/users/{id}/{version}", func(w http.ResponseWriter, r *http.Request) {
		capturedReq = &ComplexRequest{}
		err := capturedReq.ParseRequest(r)
		if err != nil {
			t.Logf("Parse error: %v", err) // Логируем ошибку для отладки
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	// ВАЖНО: Для слайсов используем многократные параметры ids=1&ids=2&ids=3
	testURL := server.URL + "/users/123e4567-e89b-12d3-a456-426614174000/1?page=1&limit=10&ids=1&ids=2&ids=3&tags=a&tags=b&active=true&created_after=2023-01-01T00:00:00Z&score=9.5"

	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	req.Header.Set("X-Request-ID", "req-123")
	req.Header.Set("X-Tenant-ID", "100")
	req.Header.Set("X-Rate-Limit", "50")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NotNil(t, capturedReq)

	// Path params
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", capturedReq.UserID)
	assert.Equal(t, 1, capturedReq.Version)

	// Query params
	assert.Equal(t, 1, capturedReq.Page)
	assert.Equal(t, 10, capturedReq.Limit)
	assert.Equal(t, []int{1, 2, 3}, capturedReq.IDs)
	assert.Equal(t, []string{"a", "b"}, capturedReq.Tags)
	assert.True(t, capturedReq.Active)

	expectedTime, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	assert.Equal(t, expectedTime, capturedReq.CreatedAfter)
	assert.Equal(t, 9.5, capturedReq.Score)

	// Header params
	assert.Equal(t, "req-123", capturedReq.RequestID)
	assert.Equal(t, int64(100), capturedReq.TenantID)
	assert.Equal(t, uint32(50), capturedReq.RateLimit)
}

func TestComplexRequest_ParseRequest_MissingOptionalParams(t *testing.T) {
	mux := http.NewServeMux()
	var capturedReq *ComplexRequest

	mux.HandleFunc("/users/{id}/{version}", func(w http.ResponseWriter, r *http.Request) {
		capturedReq = &ComplexRequest{}
		err := capturedReq.ParseRequest(r)
		if err != nil {
			t.Logf("Parse error: %v", err)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	// Запрос только с обязательными path параметрами
	testURL := server.URL + "/users/123e4567-e89b-12d3-a456-426614174000/1"

	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NotNil(t, capturedReq)

	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", capturedReq.UserID)
	assert.Equal(t, 1, capturedReq.Version)

	assert.Equal(t, 0, capturedReq.Page)
	assert.Equal(t, 0, capturedReq.Limit)
	assert.Nil(t, capturedReq.IDs)
	assert.Nil(t, capturedReq.Tags)
	assert.False(t, capturedReq.Active)
	assert.Equal(t, time.Time{}, capturedReq.CreatedAfter)
	assert.Equal(t, 0.0, capturedReq.Score)
	assert.Empty(t, capturedReq.RequestID)
	assert.Equal(t, int64(0), capturedReq.TenantID)
	assert.Equal(t, uint32(0), capturedReq.RateLimit)
}

func TestComplexRequest_ParseRequest_InvalidPathInt(t *testing.T) {
	mux := http.NewServeMux()
	var capturedErr error

	mux.HandleFunc("/users/{id}/{version}", func(w http.ResponseWriter, r *http.Request) {
		v := &ComplexRequest{}
		capturedErr = v.ParseRequest(r)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	testURL := server.URL + "/users/123e4567-e89b-12d3-a456-426614174000/not_a_number"

	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	client := &http.Client{}
	_, _ = client.Do(req)

	require.Error(t, capturedErr)
	pErr, ok := errors.AsType[*errs.ParseRequestError](capturedErr)
	require.True(t, ok, "Ожидается ошибка типа ParseRequestError")
	assert.Equal(t, "Version", pErr.Field)
	assert.Contains(t, pErr.Message, "invalid integer")
}

func TestComplexRequest_ParseRequest_InvalidQueryInt(t *testing.T) {
	mux := http.NewServeMux()
	var capturedErr error

	mux.HandleFunc("/users/{id}/{version}", func(w http.ResponseWriter, r *http.Request) {
		v := &ComplexRequest{}
		capturedErr = v.ParseRequest(r)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	testURL := server.URL + "/users/123e4567-e89b-12d3-a456-426614174000/1?page=abc"

	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	client := &http.Client{}
	_, _ = client.Do(req)

	require.Error(t, capturedErr)
	pErr, ok := errors.AsType[*errs.ParseRequestError](capturedErr)
	require.True(t, ok)
	assert.Equal(t, "Page", pErr.Field)
	assert.Contains(t, pErr.Message, "invalid integer")
}

// TestComplexRequest_ParseRequest_InvalidQueryFloat
// Проверяет ошибку при невалидном Query параметре (float)
func TestComplexRequest_ParseRequest_InvalidQueryFloat(t *testing.T) {
	mux := http.NewServeMux()
	var capturedErr error

	mux.HandleFunc("/users/{id}/{version}", func(w http.ResponseWriter, r *http.Request) {
		v := &ComplexRequest{}
		capturedErr = v.ParseRequest(r)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	testURL := server.URL + "/users/123e4567-e89b-12d3-a456-426614174000/1?score=not_float"

	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	client := &http.Client{}
	_, _ = client.Do(req)

	require.Error(t, capturedErr)
	pErr, ok := errors.AsType[*errs.ParseRequestError](capturedErr)
	require.True(t, ok)
	assert.Equal(t, "Score", pErr.Field)
	assert.Contains(t, pErr.Value, "not_float")
}

func TestComplexRequest_ParseRequest_InvalidQueryBool(t *testing.T) {
	mux := http.NewServeMux()
	var capturedErr error

	mux.HandleFunc("/users/{id}/{version}", func(w http.ResponseWriter, r *http.Request) {
		v := &ComplexRequest{}
		capturedErr = v.ParseRequest(r)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	testURL := server.URL + "/users/123e4567-e89b-12d3-a456-426614174000/1?active=maybe"

	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	client := &http.Client{}
	_, _ = client.Do(req)

	require.Error(t, capturedErr)
	pErr, ok := errors.AsType[*errs.ParseRequestError](capturedErr)
	require.True(t, ok)
	assert.Equal(t, "Active", pErr.Field)
}

func TestComplexRequest_ParseRequest_InvalidQueryTime(t *testing.T) {
	mux := http.NewServeMux()
	var capturedErr error

	mux.HandleFunc("/users/{id}/{version}", func(w http.ResponseWriter, r *http.Request) {
		v := &ComplexRequest{}
		capturedErr = v.ParseRequest(r)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	testURL := server.URL + "/users/123e4567-e89b-12d3-a456-426614174000/1?created_after=not_time"

	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	client := &http.Client{}
	_, _ = client.Do(req)

	require.Error(t, capturedErr)
	pErr, ok := errors.AsType[*errs.ParseRequestError](capturedErr)
	require.True(t, ok)
	assert.Equal(t, "CreatedAfter", pErr.Field)
	assert.Contains(t, pErr.Message, "invalid time format")
}

func TestComplexRequest_ParseRequest_InvalidSliceInt(t *testing.T) {
	mux := http.NewServeMux()
	var capturedErr error

	mux.HandleFunc("/users/{id}/{version}", func(w http.ResponseWriter, r *http.Request) {
		v := &ComplexRequest{}
		capturedErr = v.ParseRequest(r)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	testURL := server.URL + "/users/123e4567-e89b-12d3-a456-426614174000/1?ids=1&ids=abc&ids=3"

	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	client := &http.Client{}
	_, _ = client.Do(req)

	require.Error(t, capturedErr)
	pErr, ok := errors.AsType[*errs.ParseRequestError](capturedErr)
	require.True(t, ok)
	assert.Equal(t, "IDs", pErr.Field)
	assert.Equal(t, "abc", pErr.Value) // Значение, которое вызвало ошибку
}

func TestComplexRequest_ParseRequest_InvalidHeaderInt(t *testing.T) {
	mux := http.NewServeMux()
	var capturedErr error

	mux.HandleFunc("/users/{id}/{version}", func(w http.ResponseWriter, r *http.Request) {
		v := &ComplexRequest{}
		capturedErr = v.ParseRequest(r)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	testURL := server.URL + "/users/123e4567-e89b-12d3-a456-426614174000/1"

	req, _ := http.NewRequest(http.MethodGet, testURL, nil)
	req.Header.Set("X-Tenant-ID", "not_int")

	client := &http.Client{}
	_, _ = client.Do(req)

	require.Error(t, capturedErr)
	pErr, ok := errors.AsType[*errs.ParseRequestError](capturedErr)
	require.True(t, ok)
	assert.Equal(t, "TenantID", pErr.Field)
}
