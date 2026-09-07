//Этот пакет будет отвечать только за саму сущность страницы 
//и конвертацию Markdown через blackfriday

package wiki

import (
	"html/template" //пакет для работы с шаблонами

	"github.com/russross/blackfriday/v2"	
)

// Page описывает доменную модель вики-страницы
type Page struct {
	Title 		 string
	Body  		 []byte // Сырой Markdown текст из файла
	DisplayBody  template.HTML // Готовый HTML для вывода в браузер
}

// Storage описывает контракт (требования) к хранилищу данных
// "реализуйте эти два метода": Применяем правило Дейва Чейни.
type Storage interface {
	Save(title string, body []byte) error
	Load(title string) ([]byte, error)
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

// GetPage загружает сырые байты из хранилища и трансформирует Markdown в HTML.
func (s *Service) GetPage(title string) (*Page, error) {
	// 1. Просим хранилище загрузить сырые байты (вызывается метод Load)
	body, err := s.store.Load(title)
	if err != nil {
		return nil, err
	}

	// 2. Сам сервис переводит эти байты в HTML через blackfriday
	htmlBytes := blackfriday.Run(body)

	return &Page{
		Title: 			title, 
		Body: 			body,
		DisplayBody:	template.HTML(htmlBytes),  // Приводим к типу template.HTML для шаблона
	}, nil
}

// SavePage передает текст на запись в хранилище
func (s *Service) SavePage(title string, body []byte) error {
	return s.store.Save(title, body)
}