package main

import (
	"testing"
)

func TestURLShortener_ShortenAndGet(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"валидный длинный URL из задания", "https://example.com", false},
		{"валидный HTTP URL", "http://example.com", false},
		{"невалидный URL без протокола", "not-a-url", true},
		{"пустая строка", "", true},
	}

	us := NewURLShortener()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortID, err := us.Shorten(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Shorten() ошибка = %v, ожидали ошибку = %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(shortID) < 6 || len(shortID) > 8 {
					t.Errorf("короткий ID неверной длины: %s", shortID)
				}

				orig, err := us.GetOriginal(shortID)
				if err != nil {
					t.Errorf("GetOriginal() вернул ошибку для созданного ID: %v", err)
				}
				if orig != tt.url {
					t.Errorf("GetOriginal() = %s, хотели %s", orig, tt.url)
				}
			}
		})
	}
}

func TestURLShortener_GetOriginal_NotFound(t *testing.T) {
	us := NewURLShortener()
	_, err := us.GetOriginal("nonexistent")
	if err == nil {
		t.Error("ожидалась ошибка для несуществующего shortID, получили nil")
	}
}
