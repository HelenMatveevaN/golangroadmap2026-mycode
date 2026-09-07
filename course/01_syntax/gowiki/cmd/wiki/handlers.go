package main

import (
	"html/template" //пакет для работы с шаблонами
	"net/http"

	"gowiki/internal/wiki"	
)

// Храним ссылку на сервис внутри пакета main для удобства хендлеров
var wikiService *wiki.Service

// renderTemplate объединяет HTML-файл шаблона и данные структуры Page
func renderTemplate(w http.ResponseWriter, tmpl string, p *wiki.Page) {
	// // Ищем HTML-шаблоны прямо в папке cmd/wiki/
	t, err := template.ParseFiles("cmd/wiki/" + tmpl + ".html")
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
	title := r.URL.Path[len("/view/"):]
	p, err := wikiService.GetPage(title)
	if err != nil {
		// Авторедирект на форму создания/редактирования
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	// Вместо fmt.Fprintf вызываем шаблонизатор
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div><p>[<a href=\"/edit/%s\">редактировать</a>]</p>", p.Title, p.Body, p.Title)
	renderTemplate(w, "view", p)
}

// editHandler — форма редактирования страницы
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := wikiService.GetPage(title)
	if err != nil {
		// Если файла нет, создаем пустую структуру страницы для формы
		p = &wiki.Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// saveHandler — сохранение отправленных из формы данных
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	// r.FormValue("body") достает текст, который пользователь написал в текстовом поле <textarea name="body">
	body := r.FormValue("body")

	err := wikiService.SavePage(title, []byte(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}