package main

import (
	"fmt"
	"log"
	"io"
	"net/http"
	"strconv"
	"strings"
)

//curl http://127.0.0.1:8888/tasks
//curl http://127.0.0.1:8888/tasks/1
//curl 'http://127.0.0.1:8888/tasks?id=1'
//curl -X POST http://localhost:8888/tasks -H "Content-Type: text/plain" -d "Постирать вещи"
//curl -X DELETE http://localhost:8888/tasks/2

// Task — базовая структура задачи
type Task struct {
	ID    int
	Title string
}

var tasks = []Task{
		{ID: 1, Title: "Купить молоко"},
		{ID: 2, Title: "Написать код на Go"},
		{ID: 3, Title: "Помыть машину"},
}

func taskHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Вспомогательная логика извлечения ID (работает и для GET, и для DELETE)
	var idParam string
	if request.URL.Query().Has("id") {
		idParam = request.URL.Query().Get("id")
	} else {
		path := request.URL.Path
		path = strings.TrimPrefix(path, "/tasks")
		path = strings.TrimPrefix(path, "/") // Убираем оставшийся слэш, если он есть
		idParam = strings.TrimSpace(path)
	}

	switch request.Method {
	case http.MethodGet:
		if idParam == "" || idParam == "/" {
			getAllTasks(writer)
			return
		}
		getOneTask(writer, idParam)
	
	case http.MethodPost:
		// POST принимаем только на чистый адрес /tasks
		cleanPath := strings.Trim(request.URL.Path, "/")
		if cleanPath != "tasks" {
			http.Error(writer, "Метод POST не поддерживается для этого пути", http.StatusMethodNotAllowed) 
			return
		}

		defer request.Body.Close()
		bodyBytes, err := io.ReadAll(request.Body)
		if err != nil {
			// Сюда мы попадем только если реально оборвался кабель или упала сеть
			http.Error(writer, "Критическая ошибка чтения сети", http.StatusInternalServerError) //500 внутр.ошибка сервера
			return
		}

		title := strings.TrimSpace(string(bodyBytes))

		if title == "" {
			http.Error(writer, "Поле title обязательно (передайте текст задачи в теле POST-запроса)", http.StatusBadRequest) //400
			return
		}

		newID := len(tasks) + 1
		newTask := Task{
			ID: 	newID,
			Title:	title,
		}
		tasks = append(tasks, newTask)

		writer.WriteHeader(http.StatusCreated) //201
		fmt.Fprintf(writer, "Создана задача: ID %d: - %s\n", newTask.ID, newTask.Title)

	case http.MethodDelete:
		if idParam == "" || idParam == "/" {
			http.Error(writer, "Не указан ID задачи для удаления", http.StatusBadRequest) //400
			return
		}

		deleteTask(writer, idParam)

	default:
		wHeader := writer.Header()
		wHeader.Set("Allow", "GET, POST, DELETE")
		http.Error(writer, "Метод не поддерживается (только GET и POST)", http.StatusMethodNotAllowed) //405
	}
}

func getAllTasks(writer http.ResponseWriter) {
	if len(tasks) == 0 {
		fmt.Fprintf(writer, "Список задач пуст.\n")
		return
	}
	for _, task := range tasks {
		// fmt.Fprintf пишет отформатированный текст прямо в writer (http.ResponseWriter)
		fmt.Fprintf(writer, "ID: %d — %s\n", task.ID, task.Title)
	}
}

func getOneTask(writer http.ResponseWriter, idStr string) {
	// Конвертируем строку с ID в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(writer, "Параметр 'id' должен быть числом", http.StatusBadRequest) //400
		return
	}

	//поиск задачи с id
	for _, task := range tasks {
		if task.ID == id {
			fmt.Fprintf(writer, "Найдена задача: ID: %d — %s\n", task.ID, task.Title)
			return
		}
	}

	//прошли цикл и ничего не нашли
	http.Error(writer, "Задача с таким ID не найдена", http.StatusNotFound) //404
}

func deleteTask(writer http.ResponseWriter, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(writer, "Параметр 'id' должен быть числом", http.StatusBadRequest) //400
		return
	}

	// Ищем задачу по ID и запоминаем её индекс
	idx := -1
	for i, task := range tasks {
		if task.ID == id {
			idx = i
			break
		}
	}

	//если задача не найдена
	if idx == -1 {
		http.Error(writer, "Задача для удаления не найдена", http.StatusNotFound) //404
		return
	}

	// Удаляем элемент из слайса
	tasks = append(tasks[:idx], tasks[idx+1:]...)

	fmt.Fprintf(writer, "Задача с ID %d успешно удалена\n", id)
}

func main() {
	fmt.Println("Список демонстрационных задач:")
	for _, task := range tasks {
		fmt.Printf("%d: %s\n", task.ID, task.Title)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", taskHandler)
	mux.HandleFunc("/tasks/", taskHandler)

	fmt.Println("Сервер запущен на http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", mux))	
}


/*
для рефакторинга: 

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Task — базовая структура задачи
type Task struct {
	ID    int
	Title string
}

// Глобальный слайс задач (в реальных приложениях защищается sync.Mutex)
var tasks = []Task{
	{ID: 1, Title: "Купить молоко"},
	{ID: 2, Title: "Написать код на Go"},
	{ID: 3, Title: "Помыть машину"},
}

// --- MIDDLEWARE ---

// loggingMiddleware логирует метод, путь и время выполнения каждого запроса
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Общая установка заголовка для всех эндпоинтов в цепочке
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		
		next.ServeHTTP(w, r)
		
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// --- HANDLERS ---

// getAllTasksHandler обрабатывает GET /tasks
func getAllTasksHandler(writer http.ResponseWriter, request *http.Request) {
	// Поддержка старого синтаксиса получения через query-параметр: /tasks?id=1
	if request.URL.Query().Has("id") {
		getOneTask(writer, request.URL.Query().Get("id"))
		return
	}

	if len(tasks) == 0 {
		fmt.Fprintf(writer, "Список задач пуст.\n")
		return
	}
	for _, task := range tasks {
		fmt.Fprintf(writer, "ID: %d — %s\n", task.ID, task.Title)
	}
}

// getOneTaskHandler обрабатывает GET /tasks/{id}
func getOneTaskHandler(writer http.ResponseWriter, request *http.Request) {
	// Извлекаем id, встроенным в Go 1.22+ методом PathValue
	idStr := request.PathValue("id")
	getOneTask(writer, idStr)
}

// createTaskHandler обрабатывает POST /tasks
func createTaskHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Критическая ошибка чтения сети", http.StatusInternalServerError)
		return
	}

	title := strings.TrimSpace(string(bodyBytes))
	if title == "" {
		http.Error(writer, "Поле title обязательно (передайте текст задачи в теле POST-запроса)", http.StatusBadRequest)
		return
	}

	newID := len(tasks) + 1
	// Находим максимальный ID для генерации нового, чтобы избежать дублей после удалений
	for _, t := range tasks {
		if t.ID >= newID {
			newID = t.ID + 1
		}
	}

	newTask := Task{
		ID:    newID,
		Title: title,
	}
	tasks = append(tasks, newTask)

	writer.WriteHeader(http.StatusCreated)
	fmt.Fprintf(writer, "Создана задача: ID %d: - %s\n", newTask.ID, newTask.Title)
}

// deleteTaskHandler обрабатывает DELETE /tasks/{id}
func deleteTaskHandler(writer http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(writer, "Параметр 'id' должен быть числом", http.StatusBadRequest)
		return
	}

	idx := -1
	for i, task := range tasks {
		if task.ID == id {
			idx = i
			break
		}
	}

	if idx == -1 {
		http.Error(writer, "Задача для удаления не найдена", http.StatusNotFound)
		return
	}

	tasks = append(tasks[:idx], tasks[idx+1:]...)
	fmt.Fprintf(writer, "Задача с ID %d успешно удалена\n", id)
}

// Вспомогательный метод поиска одной задачи
func getOneTask(writer http.ResponseWriter, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(writer, "Параметр 'id' должен быть числом", http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			fmt.Fprintf(writer, "Найдена задача: ID: %d — %s\n", task.ID, task.Title)
			return
		}
	}

	http.Error(writer, "Задача с таким ID не найдена", http.StatusNotFound)
}

// --- MAIN ---

func main() {
	mux := http.NewServeMux()

	// Регистрация маршрутов по правилам Go 1.22+ (Метод + Шаблон пути)
	mux.HandleFunc("GET /tasks", getAllTasksHandler)
	mux.HandleFunc("GET /tasks/{id}", getOneTaskHandler)
	mux.HandleFunc("POST /tasks", createTaskHandler)
	mux.HandleFunc("DELETE /tasks/{id}", deleteTaskHandler)

	// Оборачиваем весь роутер в middleware для логирования всех запросов
	wrappedMux := loggingMiddleware(mux)

	fmt.Println("Сервер запущен на http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", wrappedMux))
}


*/