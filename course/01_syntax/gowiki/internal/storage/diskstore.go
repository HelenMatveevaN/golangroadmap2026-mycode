package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DiskStore реализует физическое сохранение файлов на диск
type DiskStore struct{}

// NewDiskStore — конструктор хранилища
func NewDiskStore() *DiskStore {
	return &DiskStore{}
}

// Save пишет .md файл на диск
func(d *DiskStore) Save(title string, body []byte) error {
	filename := title + ".md"
	//0600 — права доступа (только чтение и запись для владельца)
	if err := os.WriteFile(filename, body, 0600); err != nil {
		return fmt.Errorf("failed to write md file: %w", err)
	}
	return nil
}

// Load читает .md файл с диска
func (d *DiskStore) Load(title string) ([]byte, error) {
	filename := title + ".md"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read md file: %w", err)
	}
	return body, nil
}

// List сканирует текущую папку и возвращает имена всех файлов с расширением .md
func (c *DiskStore) List() ([]string, error) {
	files, err := os.ReadDir(".") //current dir
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var titles []string
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".md" {
			// Отрезаем расширение .md, оставляя только заголовок (например, TestPage)
			title := strings.TrimSuffix(f.Name(), ".md")
			titles = append(titles, title)
		}
	}
	return titles, nil
}