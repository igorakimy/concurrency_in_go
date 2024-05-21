package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

var r, _ = regexp.Compile("GET (.+) HTTP/1.1\r\n")

func handleHttpRequest(conn net.Conn) {
	// Создать буфер, чтобы сохранить HTTP-запрос.
	buff := make([]byte, 1024)
	// Прочитать содержимое из входящего соединения в буфер.
	size, _ := conn.Read(buff)
	// Если запрос валидный, прочитать файл из директории resources.
	if r.Match(buff[:size]) {
		file, err := os.ReadFile(
			fmt.Sprintf("../resources/%s", r.FindSubmatch(buff[:size])[1]),
		)
		// Если файл существует, ответить клиенту
		// соответствующий заголовок и содержимое файла.
		if err != nil {
			conn.Write([]byte(fmt.Sprintf(
				"HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n",
				len(file),
			)))
			conn.Write(file)
			// Если же файл не существует, вернуть ответ с ошибкой.
		} else {
			conn.Write([]byte(
				"HTTP/1.1 404 Not Found\r\n\r\n<html>Not Found</html>",
			))
		}
		// Если HTTP-запрос не валидный, вернуть соответствующую ошибку.
	} else {
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
	}
	// Закрыть соединение после обработки запроса.
	conn.Close()
}

func StartHttpWorkers(n int, incomingConnections <-chan net.Conn) {
	// Запустить n-ое количество горутин.
	for i := 0; i < n; i++ {
		go func() {
			// Потребляет соединения из канала рабочей очереди
			// пока канал не будет закрыт.
			for c := range incomingConnections {
				// Обработать HTTP-запрос из полученного соединения.
				handleHttpRequest(c)
			}
		}()
	}
}

func main() {
	// Создать канал рабочей очереди.
	incomingConnections := make(chan net.Conn)
	// Запустить пул воркеров с тремя горутинами.
	StartHttpWorkers(3, incomingConnections)

	// Привязать прослушивание TCP соединение к порту 8080.
	server, _ := net.Listen("tcp", "localhost:8080")
	defer server.Close()
	for {
		// Блокируется до тех пор, пока
		// не появится новое подключение от клиента.
		conn, _ := server.Accept()
		// Передавать подключение в канал работающей очереди.
		incomingConnections <- conn
	}
}
