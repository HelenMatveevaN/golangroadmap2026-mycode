//Этот пакет будет отвечать только за саму сущность страницы 
//и конвертацию Markdown через blackfriday

package wiki

import (
	"bytes"
	"fmt"
	"html/template" //пакет для работы с шаблонами
	"strings"

	"github.com/russross/blackfriday/v2"

	yaml "gopkg.in/yaml.v3"
)

// MetaData описывает параметры статьи из Front Matter
type MetaData struct {
	Title  string `yaml:"title"`
	Author string `yaml:"author"`
}

// Page описывает полную доменную модель страницы с метаданными
type Page struct {
	Title 		 string
	Meta 		 MetaData 		// Новое поле для метаданных
	Body  		 []byte 		// Чистый Markdown (без Front Matter)
	DisplayBody  template.HTML 	// Готовый HTML для шаблона
}

// Storage — расширяем контракт. Теперь хранилище обязано уметь давать список файлов.
type Storage interface {
	Save(title string, body []byte) error
	Load(title string) ([]byte, error)
	List() ([]string, error)
}

// Service объединяет бизнес-логику страниц и абстрактное хранилище
type Service struct {
	store Storage // Мы подсовываем сюда ИНТЕРФЕЙС, а не конкретный диск
}

// NewService — конструктор для создания "мозга" приложения.
// "принимай интерфейсы, возвращай конкретные типы структуры"
// интерфейс на входе дает гибкость
func NewService(s Storage) *Service {
	return &Service{store: s}
}

// GetPage теперь умеет отделять Front Matter от основного текста Markdown
func (s *Service) GetPage(title string) (*Page, error) {
	// Просим хранилище загрузить сырые байты (вызывается метод Load)
	rawBody, err := s.store.Load(title)
	if err != nil {
		return nil, err
	}

	var meta MetaData
	markdownBody := rawBody

	// Проверяем, начинается ли файл с Front Matter разделителя "---"
	if bytes.HasPrefix(rawBody, []byte("---\n")) || bytes.HasPrefix(rawBody, []byte("---\r\n")) {
		// Разделяем файл по "---"
		parts := strings.SplitN(string(rawBody), "---", 3) //bytes->string->[]string
		if len(parts) >= 3 {
			// Вторая часть (индекс 1) — это наш YAML блок с метаданными
			yamlBlock := parts[1]
			if err := yaml.Unmarshal([]byte(yamlBlock), &meta); err != nil {
				return nil, fmt.Errorf("ошибка парсинга YAML: %w", err)
			}

			//это чистый Markdown текст статьи
			markdownBody = []byte(strings.TrimSpace(parts[2]))
		}
	}

	// Переводим чистый Markdown без метаданных в HTML
	htmlBytes := blackfriday.Run(markdownBody)

	if meta.Title == "" {
		meta.Title = title
	}

	return &Page{
		Title: 			title,
		Meta:			meta,
		Body: 			markdownBody,
		DisplayBody:	template.HTML(htmlBytes),  // Приводим к типу template.HTML для шаблона
	}, nil
}

// SavePage передает текст на запись в хранилище
func (s *Service) SavePage(title string, body []byte) error {
	return s.store.Save(title, body)
}

// CreateNewPage — логика для команды `new`. Создает файл с дефолтным Front Matter.
func (s *Service) CreateNewPage(title, author string) error {
	defaultContent := fmt.Sprintf("---\ntitle: %q\author: %q\n---\n\n# %s\n\nНачните писать здесь...", title, author, title)
	return s.store.Save(title, []byte(defaultContent))
}



// GetALlTitles возвращает список имен всех статей
func (s *Service) GetAllTitles() ([]string, error) {
	return s.store.List()
}