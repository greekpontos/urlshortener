package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/url"
	"sync"
)

type URLShortener struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		urls: make(map[string]string),
	}
}

// Shorten создает короткий идентификатор для URL
func (us *URLShortener) Shorten(originalURL string) (string, error) {
	if !isValidURL(originalURL) {
		return "", errors.New("invalid or empty URL")
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	var shortID string
	for {
		shortID = generateShortID()
		if _, exists := us.urls[shortID]; !exists {
			break
		}
	}

	us.urls[shortID] = originalURL
	return shortID, nil
}

// GetOriginal возвращает оригинальный URL по короткому ID
func (us *URLShortener) GetOriginal(shortID string) (string, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	original, exists := us.urls[shortID]
	if !exists {
		return "", errors.New("short URL not found")
	}

	return original, nil
}

// generateShortID генерирует случайную строку от 6 до 8 символов
func generateShortID() string {
	b := make([]byte, 5)
	_, _ = rand.Read(b)
	str := base64.RawURLEncoding.EncodeToString(b)
	if len(str) > 8 {
		return str[:8]
	}
	return str
}

// isValidURL проверяет корректность URL (HTTP или HTTPS)
func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
