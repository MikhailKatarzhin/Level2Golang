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
	Key         int    // указание колонки для сортировки
	Numeric     bool   // сортировать по числовому значению (нечисловые значения трактуются как)
	Reverse     bool   // сортировать в обратном порядке
	Unique      bool   // не выводить повторяющиеся строки (с сохранением первого образца)
	OutputFile  string // наименование файла вывода
	InputFile   string // наименование файла ввода
	Month       bool   // сортировать по названию месяца, т.е. выполнять сравнение по трёх-символьным сокращениям англоязычных названий месяцев, т.е. JAN < ... < DEC , или их полному названию
	TrimSpaces  bool   // согласно man sort удаляет лишние пробелы перед и после строки
	CheckSorted bool   // проверять отсортированы ли данные
	Suffixes    bool   // сортировать по числовому значению с учетом суффиксов (согласно приставкам СИ)
}

// Порядковое представление месяца согласно его названию полному или короткому
var months = map[string]int{
	"january":   1,
	"jan":       1,
	"1":         1,
	"february":  2,
	"feb":       2,
	"2":         2,
	"march":     3,
	"mar":       3,
	"3":         3,
	"april":     4,
	"apr":       4,
	"4":         4,
	"may":       5,
	"5":         5,
	"june":      6,
	"jun":       6,
	"6":         6,
	"july":      7,
	"jul":       7,
	"7":         7,
	"august":    8,
	"aug":       8,
	"8":         8,
	"september": 9,
	"sep":       9,
	"9":         9,
	"october":   10,
	"oct":       10,
	"10":        10,
	"november":  11,
	"nov":       11,
	"11":        11,
	"december":  12,
	"dec":       12,
	"12":        12,
}

func main() {
	configs, err := parseFlagsToConfigs()
	if err != nil {
		fmt.Printf("Ошибка в процессе чтения флагов: %v\n", err)
		os.Exit(1)
	}

	lines, err := readLinesFromFile(configs)
	if err != nil {
		fmt.Printf("Ошибка в процессе чтения строк из исходного файла: %v\n", err)
		os.Exit(1)
	}

	if configs.Month {
		convertMonthKeys(lines)
	}

	if configs.Suffixes {
		parseNumberWithSuffix(lines)
	}

	if configs.CheckSorted {
		checkSorted(lines, configs)
		os.Exit(0)
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

	month := flag.Bool("M", false, "сортировать по названию месяца")
	trimSpaces := flag.Bool("b", false, "игнорировать хвостовые пробелы")
	checkSorted := flag.Bool("c", false, "проверять отсортированы ли данные")
	suffixes := flag.Bool("h", false, "сортировать по числовому значению с учётом суффиксов")

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

	if *suffixes {
		*numeric = true
	}

	*key = *key - 1

	return config{
		Key:        *key,
		Numeric:    *numeric,
		Reverse:    *reverse,
		Unique:     *unique,
		OutputFile: *outputFile,
		InputFile:  flag.Args()[0],

		Month:       *month,
		TrimSpaces:  *trimSpaces,
		CheckSorted: *checkSorted,
		Suffixes:    *suffixes,
	}, nil
}

func convertMonthKeys(lines []line) {
	for i, l := range lines {
		if value, ok := months[strings.ToLower(l.key)]; ok {
			lines[i].key = fmt.Sprintf("%d", value)
		} else {
			lines[i].key = "0"
		}
	}

}

// readLinesFromFile считывает строки из указанного файла и возвращает их в виде слайса структур line.
func readLinesFromFile(configs config) ([]line, error) {

	file, err := os.Open(configs.InputFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка в ходе открытия файла: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []line

	for scanner.Scan() {
		content := scanner.Text()

		if configs.TrimSpaces {
			content = strings.TrimLeft(content, " ")
			content = strings.TrimRight(content, " ")
		}

		fields := strings.Fields(content)

		/*
			C ключом -k (который указывает, по какому полю нужно сортировать строки),
			строки, в которых отсутствует указанная колонка, будут рассматриваться как если бы они имели пустое значение в этой колонке.
		*/
		if configs.Key >= len(fields) {
			lines = append(lines, line{
				content: content,
				key:     "",
			})
		} else {
			lines = append(lines, line{
				content: content,
				key:     fields[configs.Key],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка сканера в ходе чтения файла: %v", err)
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

func parseNumberWithSuffix(lines []line) {
	suffixes := map[string]float64{
		"da": 1e1,
		"h":  1e2,
		"k":  1e3,
		"M":  1e6,
		"G":  1e9,
		"T":  1e12,
		"P":  1e15,
		"E":  1e18,
		"Z":  1e21,
		"Y":  1e24,
		"R":  1e27,
		"Q":  1e30,
		"d":  1e-1,
		"c":  1e-2,
		"m":  1e-3,
		"µ":  1e-6,
		"n":  1e-9,
		"p":  1e-12,
		"f":  1e-15,
		"a":  1e-18,
		"z":  1e-21,
		"y":  1e-24,
		"r":  1e-27,
		"q":  1e-30,
	}

	for i, l := range lines {
		for suffix, multiplier := range suffixes {
			if strings.HasSuffix(l.key, suffix) {
				numStr := strings.TrimSuffix(l.key, suffix)
				if num, err := strconv.ParseFloat(numStr, 64); err == nil {
					lines[i].key = fmt.Sprintf("%f", num*multiplier)
				}

				break
			}
		}
	}
}

// sortLines сортирует слайс строк на основе указанных флагов.
func sortLines(lines []line, configs config) {

	slices.SortStableFunc(lines, func(a, b line) int {
		var result int

		if configs.Numeric || configs.Month {
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

// checkSorted проверяет, отсортированы ли строки в соответствии с заданными параметрами.
func checkSorted(lines []line, configs config) {
	compare := func(a, b line) int {
		var result int

		if configs.Numeric || configs.Month {
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
	}

	for i := 0; i < len(lines)-1; i++ {
		result := compare(lines[i], lines[i+1])
		if result != 1 {
			fmt.Printf("Файл не отсортирован, ошибка на строке %d: %q\n", i+1, lines[i].content)
			return
		}
	}
}
