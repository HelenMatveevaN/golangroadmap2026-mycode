package main

import (
	"fmt"
	"log"
	"net/http"
	
	"gowiki/internal/storage"
	"gowiki/internal/wiki"
)

func main() {
	//init
	store := storage.NewDiskStore()

	//передача "рук" в "мозг"
	//Конструктор NewService на входе ожидает абстрактный интерфейс 
	//wiki.Storage
	//(это просто список требований: «мне нужен любой объект с методами 
	//Load и Save»).
	wikiService = wiki.NewService(store)


	// Регистрируем все три ручки
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	fmt.Println("🌐 Чистый сервер gowiki запущен на http://localhost:8085/view/TestPage")
	log.Fatal(http.ListenAndServe(":8085", nil))
}