package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type config struct {
	fields    string
	delimiter string
	separated bool
}

func main() {
	configs, err := parseFlagsToConfigs()
	if err != nil {
		fmt.Printf("Ошибка в процессе чтения флагов: %v", err)
		os.Exit(1)
	}

	// Парсинг полей, если они указаны
	var fields []int
	if configs.fields != "" {
		for _, f := range strings.Split(configs.fields, ",") {
			field, err := strconv.Atoi(f)
			if err != nil {
				fmt.Printf("Ошибка в ходе считывания номеров полей: %v", err)
				os.Exit(1)
			}

			fields = append(fields, field-1)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// Проверка, содержит ли строка разделитель
		if configs.separated && !strings.Contains(line, configs.delimiter) {
			continue
		}

		// Разделение строки по указанному разделителю
		splitLine := strings.Split(line, configs.delimiter)

		// Выбор и вывод нужных полей
		if len(fields) > 0 {
			var output strings.Builder
			for _, field := range fields {
				if field >= 0 && field < len(splitLine) {
					output.WriteString(splitLine[field])
				}
			}
			fmt.Println(output.String())
		} else {
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения ввода:", err)
	}
}

func parseFlagsToConfigs() (config, error) {
	fields := flag.String("f", "", "выбрать поля (колонки)")
	delimiter := flag.String("d", "\t", "использовать другой разделитель")
	separated := flag.Bool("s", false, "только строки с разделителем")

	flag.Parse()

	if len(flag.Args()) > 0 {
		return config{}, fmt.Errorf("обнаружены неизвестные аргументы")
	}

	return config{
		fields:    *fields,
		delimiter: *delimiter,
		separated: *separated,
	}, nil
}
