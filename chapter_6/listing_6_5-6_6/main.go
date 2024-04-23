package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func fileSearch(dir string, filename string, wg *sync.WaitGroup) {
	// Читает все файлы из директории, переданной аргументом в функцию.
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		// Присоединить каждый файл к директории.
		fpath := filepath.Join(dir, file.Name())
		if strings.Contains(file.Name(), filename) {
			// Если совпадения есть, вывести их в консоль.
			fmt.Println(fpath)
		}
		// Если это директория...
		if file.IsDir() {
			// Добавить 1 в группу ожидания  до запуска новой горутины.
			wg.Add(1)
			// Создать горутину рекурсивно, осуществляя поиск
			// по новой директории.
			go fileSearch(fpath, filename, wg)
		}
	}
	// Вызвать Done() на группе ожидания после обработки всех файлов.
	wg.Done()
}

func main() {
	// Создать новую пустую группу ожидания.
	wg := sync.WaitGroup{}
	// Добавить дельту в размере 1 в группу ожидания.
	wg.Add(1)
	// Создать новую горутину, выполняющую файловый поиск
	// и передавая ей ссылку на группу ожидания.
	go fileSearch(os.Args[1], os.Args[2], &wg)
	// Подождать, пока поиск закончится.
	wg.Wait()
}
