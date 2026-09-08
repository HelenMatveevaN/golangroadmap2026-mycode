package wiki

import (
	"errors"
	"testing"
)

/*
➜  gowiki git:(main) ✗ go test ./internal/wiki
ok  	gowiki/internal/wiki	0.472s

➜  gowiki git:(main) ✗ go test -cover ./internal/wiki
ok  	gowiki/internal/wiki	0.479s	coverage: 85.7% of statements

*/

// mockStorage — это наше виртуальное хранилище-заглушка (Mock).
// Вместо физического диска оно сохраняет файлы прямо в оперативную память.
type mockStorage struct {
	files map[string][]byte
}

// Реализуем метод Save для нашего виртуального диска
func (m *mockStorage) Save(title string, body []byte) error {
	m.files[title] = body
	return nil
}

// Реализуем метод Load для нашего виртуального диска
func (m *mockStorage) Load(title string) ([]byte, error) {
	body, ok := m.files[title]
	if !ok {
		// Если файла нет в мапе, возвращаем встроенную ошибку Go (как os.ErrNotExist)
		return nil, errors.New("file not found")
	}
	return body, nil
}

// Реализуем метод List для нашего виртуального диска
func (m *mockStorage) List() ([]string, error) {
	var titles []string
	for title := range m.files {
		titles = append(titles, title)
	}
	return titles, nil
}

// TestGetPage_WithFrontMatter проверяет, как ядро системы отделяет метаданные от Markdown
func TestGetPage_WithFrontMatter(t *testing.T) {
	// 1. Создаем виртуальное хранилище в памяти
	mockStore := &mockStorage{
		files: make(map[string][]byte),
	}

	// Подготавливаем тестовый текст статьи, оформленный по всем правилам Front Matter
	testContent := `---
title: "Тестовая статья"
author: "Мартин Фаулер"
---
# Мой заголовок

Привет, это **проверка** генератора сайтов!`

	// Записываем статью на наш виртуальный диск
	_ = mockStore.Save("test-slug", []byte(testContent))

	// 2. Инициализируем наш "мозг" (Service), подсовывая ему виртуальный диск вместо настоящего
	service := NewService(mockStore)

	// 3. Вызываем тестируемый метод GetPage
	page, err := service.GetPage("test-slug")
	if err != nil {
		t.Fatalf("GetPage вернул ошибку, хотя файл в памяти существует: %v", err)
	}

	// 4. Проверяем, что парсер Front Matter вытащил метаданные без ошибок
	if page.Meta.Title != "Тестовая статья" {
		t.Errorf("Неверный заголовок из метаданных. Ожидали: 'Тестовая статья', получили: %q", page.Meta.Title)
	}
	if page.Meta.Author != "Мартин Фаулер" {
		t.Errorf("Неверный автор из метаданных. Ожидали: 'Мартин Фаулер', получили: %q", page.Meta.Author)
	}

	// 5. Проверяем, что библиотека blackfriday успешно превратила Markdown в живой HTML-код
	expectedHTML := "<h1>Мой заголовок</h1>\n\n<p>Привет, это <strong>проверка</strong> генератора сайтов!</p>\n"
	
	if string(page.DisplayBody) != expectedHTML {
		t.Errorf("Markdown отрендерился неверно.\nОжидали HTML:\n%q\nПолучили HTML:\n%q", expectedHTML, string(page.DisplayBody))
	}
}

// TestService_SaveAndCreate проверяет оставшиеся методы сохранения и создания для 100% покрытия
func TestService_SaveAndCreate(t *testing.T) {
	mockStore := &mockStorage{
		files: make(map[string][]byte),
	}
	service := NewService(mockStore)

	// 1. Тестируем SavePage
	err := service.SavePage("save-test", []byte("hello"))
	if err != nil {
		t.Errorf("SavePage вернул ошибку: %v", err)
	}

	// 2. Тестируем CreateNewPage
	err = service.CreateNewPage("new-test", "Author")
	if err != nil {
		t.Errorf("CreateNewPage вернул ошибку: %v", err)
	}

	// 3. Тестируем GetAllTitles (метод List)
	titles, err := service.GetAllTitles()
	if err != nil {
		t.Errorf("GetAllTitles вернул ошибку: %v", err)
	}
	if len(titles) != 2 {
		t.Errorf("Ожидали 2 файла в памяти, получили: %d", len(titles))
	}
}