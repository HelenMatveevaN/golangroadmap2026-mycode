package main

import (
	"html/template" //пакет для работы с шаблонами
	"log"
	"net/http"
	"os"

	"github.com/russross/blackfriday/v2"
)

// Page описывает структуру вики-страницы
type Page struct {
	Title 		 string
	Body  		 []byte // Сырой Markdown текст из файла
	DisplayBody  template.HTML // Готовый HTML для вывода в браузер
}

// save записывает страницу на диск в виде текстового файла
func (p *Page) save() error {
	filename := p.Title + ".md"
	// Записываем файл на диск. 0600 — права доступа (только чтение и запись для владельца)
	return os.WriteFile(filename, p.Body, 0600)
}

// loadPage читает .md файл и превращает его текст в HTML
func loadPage(title string) (*Page, error) {
	filename := title + ".md"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Магия: blackfriday превращает Markdown-байты в HTML-байты
	htmlBytes := blackfriday.Run(body)

	return &Page{
		Title: 			title, 
		Body: 			body,
		DisplayBody:	template.HTML(htmlBytes),  // Приводим к типу template.HTML для шаблона
	}, nil
}

// renderTemplate объединяет HTML-файл шаблона и данные структуры Page
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	// Парсим файл (например, "view.html")
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, "Ошибка сервера: " + err.Error(), http.StatusInternalServerError)
		return
	}
	// Рендерим шаблон в ответ пользователю
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// viewHandler — просмотр страницы
func viewHandler(w http.ResponseWriter, r *http.Request) {
	// r.URL.Path[6:] отрезает первые 6 символов "/view/" из пути, оставляя только заголовок "TestPage"
	title := r.URL.Path[6:]
	p, err := loadPage(title)
	if err != nil {
		// Если страницы нет, автоматически перенаправляем пользователя на страницу её создания/редактирования!
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	// Вместо fmt.Fprintf вызываем шаблонизатор
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div><p>[<a href=\"/edit/%s\">редактировать</a>]</p>", p.Title, p.Body, p.Title)
	renderTemplate(w, "view", p)
}

// editHandler — форма редактирования страницы
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[6:]

	p, err := loadPage(title)
	if err != nil {
		// Если файла нет, создаем пустую структуру страницы для формы
		p = &Page{Title: title}
	}

	// Вместо fmt.Fprintf вызываем шаблонизатор
	/*fmt.Fprintf(w, `
		<h1>Редактирование страницы: %s</h1>
		<form action="/save/%s" method="POST">
			<textarea name="body" rows="20" cols="80">%s</textarea><br>
			<input type="submit" value="Сохранить">
		</form>
		`, p.Title, p.Title, p.Body)*/
	renderTemplate(w, "edit", p)
}

// saveHandler — сохранение отправленных из формы данных
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[6:]

	// r.FormValue("body") достает текст, который пользователь написал в текстовом поле <textarea name="body">
	body := r.FormValue("body")

	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// После успешного сохранения возвращаем пользователя на страницу просмотра
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	// Регистрируем все три ручки
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	log.Println("🌐 Сервер gowiki с поддержкой Markdown запущен на http://localhost:8085/view/TestPage")
	log.Fatal(http.ListenAndServe(":8085", nil))
}