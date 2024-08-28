package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Структура для хранения строк и значения ключевого поля.
type line struct {
	content string
	key     string
}

type config struct {
	Key        int    //указание колонки для сортировки
	Numeric    bool   //сортировать по числовому значению
	Reverse    bool   //сортировать в обратном порядке
	Unique     bool   //не выводить повторяющиеся строки
	OutputFile string //наименование файла вывода
	InputFile  string //наименование файла ввода
}

func main() {
	configs, err := parseFlagsToConfigs()
	if err != nil {
		fmt.Printf("Ошибка в процессе чтения флагов: %v\n", err)
		os.Exit(1)
	}

	lines, err := readLinesFromFile(configs.InputFile, configs.Key)
	if err != nil {
		fmt.Printf("Ошибка в процессе чтения строк из исходного файла: %v\n", err)
		os.Exit(1)
	}

	if configs.Unique {
		lines = removeDuplicates(lines, configs.Numeric)
	}

	sortLines(lines, configs)

	if err := writeOutput(lines, configs.OutputFile); err != nil {
		fmt.Printf("Ошибка в процессе записи строк в файл вывода: %v", err)
		os.Exit(1)
	}
}

func parseFlagsToConfigs() (config, error) {
	key := flag.Int("k", 1, "указание колонки для сортировки")
	numeric := flag.Bool("n", false, "сортировать по числовому значению")
	reverse := flag.Bool("r", false, "сортировать в обратном порядке")
	unique := flag.Bool("u", false, "не выводить повторяющиеся строки")
	outputFile := flag.String("o", "out.txt", "имя выходного файла")

	flag.Parse()

	if *key < 1 {
		return config{}, fmt.Errorf("ключ -k должен быть не меньше 1")
	}

	if len(flag.Args()) < 1 {
		return config{}, fmt.Errorf("название исходного файла не указано")
	}

	if *outputFile == "" {
		return config{}, fmt.Errorf("имя выходного файла не может быть пустым")
	}

	*key = *key - 1

	return config{*key, *numeric, *reverse, *unique, *outputFile, flag.Args()[0]}, nil
}

// readLinesFromFile считывает строки из указанного файла и возвращает их в виде слайса структур line.
func readLinesFromFile(filename string, key int) ([]line, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка в ходе открытия файла: %v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []line

	for scanner.Scan() {
		content := scanner.Text()
		fields := strings.Fields(content)

		/*
			C ключом -k (который указывает, по какому полю нужно сортировать строки),
			строки, в которых отсутствует указанная колонка, будут рассматриваться как если бы они имели пустое значение в этой колонке.
		*/
		if key >= len(fields) {
			lines = append(lines, line{
				content: content,
				key:     "",
			})
		} else {
			lines = append(lines, line{
				content: content,
				key:     fields[key],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка сканера в ходе чтения файла: %v\n", err)
	}

	return lines, nil
}

// removeDuplicates удаляет дублирующиеся строки, оставляя только уникальные.
func removeDuplicates(lines []line, numeric bool) []line {
	checkedKeys := make(map[string]struct{})
	var uniqueLines []line

	for _, l := range lines {
		var keyToCheck string

		if numeric {
			// Попытка интерпретировать ключ как число
			if _, err := strconv.ParseFloat(l.key, 64); err == nil {
				keyToCheck = l.key
			} else {
				// Если ключ не числовой, представляем его как пустую строку
				keyToCheck = ""
			}
		} else {
			keyToCheck = l.key
		}

		if _, exists := checkedKeys[keyToCheck]; !exists {
			checkedKeys[keyToCheck] = struct{}{}
			uniqueLines = append(uniqueLines, l)
		}
	}

	return uniqueLines
}

// sortLines сортирует слайс строк на основе указанных флагов.
func sortLines(lines []line, configs config) {

	slices.SortStableFunc(lines, func(a, b line) int {
		var result int

		// Если указана сортировка по числовому значению
		if configs.Numeric {
			num1, err1 := strconv.ParseFloat(a.key, 64)
			num2, err2 := strconv.ParseFloat(b.key, 64)
			if err1 == nil && err2 == nil {
				if num1 < num2 {
					result = -1
				} else if num1 == num2 {
					result = 0
				} else {
					result = 1
				}

			} else {
				result = strings.Compare(a.key, b.key)
			}
		} else {
			result = strings.Compare(a.key, b.key)
		}

		if configs.Reverse {
			return -result
		}
		return result
	})
}

// writeOutput записывает отсортированные строки в указанный выходной файл или выводит на консоль.
func writeOutput(lines []line, outputFile string) error {
	var outputWriter *bufio.Writer

	// Открытие выходного файла для записи
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("ошибка создания выходного файла \"%s\": %v", outputFile, err)
	}
	defer file.Close()

	outputWriter = bufio.NewWriter(file)
	defer outputWriter.Flush()

	// Запись строк
	for _, line := range lines {
		if _, err := fmt.Fprintln(outputWriter, line.content); err != nil {
			return fmt.Errorf("ошибка записи строки: ошибка = %v, строка =\"%v\"", err, line.content)
		}
	}

	return nil
}
