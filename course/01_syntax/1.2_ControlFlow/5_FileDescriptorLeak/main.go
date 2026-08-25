package main

import "os"

func processFile(filename string) error {
	file, err := os.Open("text.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func main() {
	//1-wrong, утечка файловых дескрипторов
	/*for i:=0; i < 1000; i++ {
		file, _ := os.Open("text.txt")
		defer file.Close()
	}*/

	//2-анонимн.ф-я
	func() {
		file, err := os.Open("text.txt")
		if err != nil {
			return
		}
		defer file.Close()
	}()

	//3-ф-я-воркер
	for i:=0; i < 1000; i++ {
		if err := processFile("text.txt"); err != nil {
			//обработка ошибки
		}
	}

}