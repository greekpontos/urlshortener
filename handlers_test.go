package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		expectedURL    string
	}{
		{
			name:           "Успешный POST с URL из задания",
			method:         "POST",
			body:           `{"url": "https://example.com"}`,
			expectedStatus: http.StatusOK,
			expectedURL:    "https://example.com",
		},
		{
			name:           "Неверный JSON",
			method:         "POST",
			body:           `{"url": "https://...`,
			expectedStatus: http.StatusBadRequest,
			expectedURL:    "",
		},
		{
			name:           "Невалидный URL",
			method:         "POST",
			body:           `{"url": "bad-url"}`,
			expectedStatus: http.StatusBadRequest,
			expectedURL:    "",
		},
		{
			name:           "Неверный метод HTTP",
			method:         "GET",
			body:           ``,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedURL:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/shorten", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()

			ShortenHandler(rec, req)

			res := rec.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("ожидали статус %d, получили %d", tt.expectedStatus, res.StatusCode)
			}

			if res.StatusCode == http.StatusOK {
				var resp ShortenResponse
				if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
					t.Fatalf("не удалось распарсить JSON ответа: %v", err)
				}
				if resp.ShortURL == "" || resp.OriginalURL != tt.expectedURL {
					t.Errorf("неверные данные в JSON ответе: %+v", resp)
				}
				if len(resp.ShortURL) < 6 || len(resp.ShortURL) > 8 {
					t.Errorf("длина короткого URL вне диапазона 6-8: %s", resp.ShortURL)
				}
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	targetURL := "https://example.com"
	id, err := shortener.Shorten(targetURL)
	if err != nil {
		t.Fatalf("ошибка подготовки теста: %v", err)
	}

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedLoc    string
	}{
		{"Успешный редирект", "GET", "/" + id, http.StatusFound, targetURL},
		{"ID не найден", "GET", "/notfound123", http.StatusNotFound, ""},
		{"Пустой ID", "GET", "/", http.StatusNotFound, ""},
		{"Неверный метод", "POST", "/" + id, http.StatusMethodNotAllowed, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			RedirectHandler(rec, req)

			res := rec.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("ожидали статус %d, получили %d", tt.expectedStatus, res.StatusCode)
			}

			if res.StatusCode == http.StatusFound {
				loc := res.Header.Get("Location")
				if loc != tt.expectedLoc {
					t.Errorf("ожидали редирект на %s, улетели на %s", tt.expectedLoc, loc)
				}
			}
		})
	}
}
