package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type config struct {
	netcatMode bool
	address    string
	port       string
	udp        bool
}

func main() {
	configs, err := parseFlagsToConfigs()
	if err != nil {
		fmt.Printf("В процессе обработки флагов произошла ошибка: %v", err)
	}

	if configs.netcatMode {
		if err := runNetcat(configs); err != nil {
			fmt.Printf("В процессе выполнения netcat произошла ошибка: %v", err)
		}
	} else {
		runShell()
	}
}

func parseFlagsToConfigs() (config, error) {
	netcatMode := flag.Bool("netcat", false, "Запуск netcat клиента")
	address := flag.String("address", "", "Адрес для подключения")
	port := flag.String("port", "", "Порт для подключения")
	udp := flag.Bool("udp", false, "Заменить тип соединения на UDP? (По стандарту: TCP)")

	flag.Parse()

	if *netcatMode && (*address == "" || *port == "") {
		return config{}, fmt.Errorf("для работы netcat необходимо указать и адрес, и порт")
	}

	return config{
			netcatMode: *netcatMode,
			address:    *address,
			port:       *port,
			udp:        *udp,
		},
		nil
}

// Shell-related functions
func runShell() {
	for {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if len(input) == 0 {
			continue
		}

		if strings.Contains(input, "|") {
			commands := strings.Split(input, "|")

			var cmds [][]string
			for _, command := range commands {
				cmds = append(cmds, strings.Fields(command))
			}

			runCommandWithPipes(cmds)

		} else {
			args := strings.Split(input, " ")
			cmd := args[0]

			switch cmd {
			case "exit":
				return
			case "cd":
				if len(args) < 2 {
					fmt.Println("cd: недостаточно аргументов")
					continue
				}

				if err := os.Chdir(args[1]); err != nil {
					fmt.Println("cd:", err)
				}
			case "pwd":
				if dir, err := os.Getwd(); err != nil {
					fmt.Println("pwd:", err)
				} else {
					fmt.Println(dir)
				}
			case "echo":
				var result strings.Builder

				result.WriteString(args[1])
				for _, v := range args[2:] {
					result.WriteRune(' ')
					result.WriteString(v)
				}

				fmt.Println(result.String())
			case "kill":
				if len(args) < 2 {
					fmt.Println("kill: отсутствует PID")
					continue
				}

				pid, err := strconv.Atoi(args[1])
				if err != nil {
					fmt.Println("kill: некорректный PID")
					continue
				}

				proc, err := os.FindProcess(pid)
				if err != nil {
					fmt.Println("kill:", err)
					continue
				}

				if err := proc.Kill(); err != nil {
					fmt.Println("kill:", err)
				}
			case "ps":
				if runtime.GOOS == "windows" {
					cmd := exec.Command("tasklist")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						fmt.Println("ps:", err)
					}
				} else {
					if err := exec.Command("ps").Run(); err != nil {
						fmt.Println("ps:", err)
					}
				}
			default:
				runCommand(args)
			}
		}
	}
}

func runCommand(args []string) {
	if len(args) == 0 {
		return
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

func runCommandWithPipes(commands [][]string) {
	if len(commands) == 0 {
		return
	}

	// Создаем канал для передачи данных между командами
	var chans []chan string
	for range commands {
		chans = append(chans, make(chan string))
	}

	for i, args := range commands {
		// Запуск команды в горутине
		go func(i int, args []string) {
			defer close(chans[i])

			cmd := exec.Command(args[0], args[1:]...)

			if i > 0 {
				go func() {
					for line := range chans[i-1] {
						cmd.Stdin = strings.NewReader(line)
					}
				}()
			}

			// Запись вывода команды в текущий канал или в stdout, если это последняя команда
			output, err := cmd.Output()
			if err != nil {
				fmt.Printf("Ошибка выполнения команды: %v\n", err)
				return
			}

			if i < len(commands)-1 {
				chans[i] <- string(output)
			} else {
				fmt.Print(string(output))
			}
		}(i, args)
	}

	// Дожидаемся завершения всех горутин
	for i := range chans {
		for range chans[i] {
		}
	}
}

// Netcat-related functions

func runNetcat(configs config) error {
	addr := fmt.Sprintf("%s:%s", configs.address, configs.port)

	network := "tcp"

	if configs.udp {
		network = "udp"
	}

	conn, err := net.Dial(network, addr)
	if err != nil {
		return fmt.Errorf("ошибка подключения: %v", err)
	}
	defer conn.Close()

	go func() {
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("Соединение закрыто")
				} else {
					fmt.Printf("ошибка чтения из подключения: %v\n", err)
				}
				return
			}
			fmt.Print(line)
		}
	}()

	writer := bufio.NewWriter(conn)
	stdin := bufio.NewReader(os.Stdin)
	for {
		input, err := stdin.ReadString('\n')
		if err != nil {
			return fmt.Errorf("Ошибка чтения из stdin: %v", err)
		}

		input = strings.TrimSpace(input)
		if input == "exit" {
			return nil
		}

		_, err = writer.WriteString(input + "\n")
		if err != nil {
			return fmt.Errorf("Ошибка при записи через соединение: %v", err)
		}
		writer.Flush()
	}
}
