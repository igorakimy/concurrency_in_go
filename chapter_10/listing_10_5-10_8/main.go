package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CodeDepth struct {
	file  string
	level int
}

func deepestNestedBlock(filename string) CodeDepth {
	// Прочитать содержимое файла в буфер памяти.
	code, _ := os.ReadFile(filename)
	maximum := 0
	level := 0
	// Итерировать по каждому символу в файле.
	for _, c := range code {
		if c == '{' {
			// Когда символ является открывающей фигурной скобкой,
			// увеличить переменную уровня вложенности level на 1.
			level += 1
			// Записать максимальное значение уровня вложенности в переменную.
			maximum = int(math.Max(float64(maximum), float64(level)))
		} else if c == '}' {
			// Когда встретиться закрывающая фигурная скобка,
			// уменьшить уровень вложенности на 1.
			level -= 1
		}
	}
	// Вернуть результат с именем файла.
	return CodeDepth{filename, maximum}
}

func forkIfNeeded(path string, info os.FileInfo,
	wg *sync.WaitGroup, results chan CodeDepth) {
	// Проверить, что указанный путь является файлом и имеет расширение .go.
	if !info.IsDir() && strings.HasSuffix(path, ".go") {
		// Добавить 1 к группе ожидания.
		wg.Add(1)
		// Создать новую горутину.
		go func() {
			// Вызвать функцию и записать возвращаемый результать
			// в результирующий канал.
			results <- deepestNestedBlock(path)
			// Отметить, что работа выполнена на группе ожидания.
			wg.Done()
		}()
	}
}

func joinResults(partialResults chan CodeDepth) chan CodeDepth {
	// Создать канал, который будет содержать финальный результат.
	finalResult := make(chan CodeDepth)
	maximum := CodeDepth{"", 0}
	go func() {
		// Получать результаты из канала, пока он не будет закрыт.
		for pr := range partialResults {
			if pr.level > maximum.level {
				// Записать значение наиболее глубоко вложенного блока.
				maximum = pr
			}
		}
		// После того как канал будет закрыт,
		// записать результат в выходной канал.
		finalResult <- maximum
	}()
	return finalResult
}

func main() {
	// Прочитать корневую директорию из аргументов.
	dir := os.Args[1]
	// Создать общий канал, используемый всеми разветвляющими горутинами.
	partialResults := make(chan CodeDepth)
	wg := sync.WaitGroup{}

	// Проходиться по корневой директории и для каждого файла
	// вызывает функцию, которая создает горутину.
	_ = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		forkIfNeeded(path, info, &wg, partialResults)
		return nil
	})
	// Вызывать соединяющую функцию и получить канал,
	// который будет содержать финальный результат.
	finalResult := joinResults(partialResults)

	// Подождать, пока все разветвляющие горутины завершат свою работу.
	wg.Wait()

	// Закрыть общий канал, дающий сигнал соединяющей горутине,
	// о том, что работа завершена.
	close(partialResults)

	// Получить финальный результат и вывести его в консоль.
	result := <-finalResult
	fmt.Printf("%s has the deepest nested code block of %d\n",
		result.file, result.level)
}
