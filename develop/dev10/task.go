package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

// В windows вместо Ctrl + D нужно сочетание Ctrl + Z

// Для теста рекомендую подключаться к towel.blinkenlights.nl 23

func main() {
	// Обработка флагов командной строки
	timeoutFlag := flag.Duration("timeout", 10*time.Second, "Время ожидания")
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Println("Шаблон команды запуска: go run task.go [--timeout={10}s] {host} {port}")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	// Установка соединения с сервером
	conn, err := net.DialTimeout("tcp", address, *timeoutFlag)
	if err != nil {
		fmt.Println("Ошибка подключения: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Соединение с %s устьановлено \n", address)

	done := make(chan struct{})

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Printf("Ошибка при копировании данных из соединения в консоль: %v\n", err)
		} else {
			fmt.Println("\nСоединение закрыто хостом")
		}

		done <- struct{}{}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err := fmt.Fprintln(conn, scanner.Text())
			if err != nil {
				fmt.Println("Ошибка записи в соединение: ", err)
				break
			}
		}

		if scanner.Err() != nil {
			fmt.Println("Ошибка чтения STDIN:", scanner.Err())
		}

		conn.Close()
		done <- struct{}{}
	}()

	// Ожидание завершения
	<-done
	fmt.Println("Завершение...")
}
