package main

import (
	"testing"
	"time"
)

// 1. Тестируем дефолтный сервер
func TestNewServer_Defaults(t *testing.T) {
	// Создаем сервер без опций
	srv := NewServer("localhost", 8080)

	// Проверяем, что хост и порт записались правильно
	if srv.host != "localhost" || srv.port != 8080 {
		t.Errorf("Ожидался адрес localhost:8080, получили %s:%d", srv.host, srv.port)
	}

	// Проверяем, что применился дефолтный таймаут 30 секунд
	if srv.timeout != 30*time.Second {
		t.Errorf("Ожидался дефолтный таймаут 30s, получили %v", srv.timeout)
	}

	// Проверяем дефолтные соединения (должно быть 10)
	if srv.maxConn != 10 {
		t.Errorf("Ожидалось дефолтное maxConn 10, получили %d", srv.maxConn)
	}
}

// 2. Тестируем сервер с кастомными опциями
func TestNewServer_WithOptions(t *testing.T) {
	// Создаем сервер и передаем функции-опции
	srv := NewServer("127.0.0.1", 9000, WithTimeout(5*time.Second), WithMaxConn(100))

	// Проверяем, перетерся ли дефолтный таймаут на 5 секунд
	if srv.timeout != 5*time.Second {
		t.Errorf("Ожидался кастомный таймаут 5s, получили %v", srv.timeout)
	}

	// Проверяем, перетерлось ли maxConn на 100
	if srv.maxConn != 100 {
		t.Errorf("Ожидалось кастомное maxConn 100, получили %d", srv.maxConn)
	}
}
