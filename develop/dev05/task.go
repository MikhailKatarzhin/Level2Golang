package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type config struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
	pattern    string
	files      []string
}

func main() {
	configs, err := parseFlagsToConfigs()
	if err != nil {
		fmt.Printf("Ошибка в процессе чтения флагов: %v", err)
		os.Exit(1)
	}

	if err := processFiles(configs); err != nil {
		fmt.Printf("Ошибка в процессе обработки файлов: %v", err)
		os.Exit(1)
	}
}

func parseFlagsToConfigs() (config, error) {
	after := flag.Int("A", 0, "печать +N строк после совпадения")
	before := flag.Int("B", 0, "печать +N строк до совпадения")
	context := flag.Int("C", 0, "печать ±N строк вокруг совпадения")
	count := flag.Bool("c", false, "количество строк")
	ignoreCase := flag.Bool("i", false, "игнорировать регистр")
	invert := flag.Bool("v", false, "исключать строки, соответствующие шаблону")
	fixed := flag.Bool("F", false, "точное совпадение со строкой")
	lineNum := flag.Bool("n", false, "печатать номер строки")

	flag.Parse()

	if len(flag.Args()) < 2 {
		return config{}, fmt.Errorf("несоответсвие шаблону команды: grep [флаги] шаблон_строки [файл...]")
	}

	return config{
		after:      *after,
		before:     *before,
		context:    *context,
		count:      *count,
		ignoreCase: *ignoreCase,
		invert:     *invert,
		fixed:      *fixed,
		lineNum:    *lineNum,
		pattern:    flag.Args()[0],
		files:      flag.Args()[1:],
	}, nil
}

func convertPatternToRegexp(configs config) (*regexp.Regexp, error) {
	if configs.ignoreCase {
		configs.pattern = "(?i)" + configs.pattern
	}

	if regex, err := regexp.Compile(configs.pattern); err != nil {
		return nil, fmt.Errorf("недействительный паттерн строки: %v", err)
	} else {
		return regex, nil
	}
}

func processFiles(configs config) error {
	regex, err := convertPatternToRegexp(configs)
	if err != nil {
		return fmt.Errorf("ошибка в процессе обработки паттерна строки: %v", err)
	}

	var processFileFail int

	for i, file := range configs.files {

		if err := processFile(i, regex, configs); err != nil {
			processFileFail++
			fmt.Printf("Ошибка в процессе обработки файла \"%s\": %v\n", file, err)
			continue
		}
	}

	return nil
}

func processFile(fileNumber int, regex *regexp.Regexp, configs config) error {
	scanResult, err := scanFile(fileNumber, regex, configs)
	if err != nil {
		return fmt.Errorf("В процессе сканирования файла произошла ошибка: %v", err)
	}

	if configs.count {
		fmt.Printf("В файле \"%s\" обнаружено %d соответствующих строк\n", configs.files[fileNumber], scanResult.nMatchedLines)
		return nil
	}

	linePrint := convertLineMatchToLinePrintArray(configs, scanResult.lineMatch)

	if configs.invert {
		for i, v := range linePrint {
			linePrint[i] = !v
		}
	}

	fmt.Printf("Результат фильтрации файла \"%s\"\n", configs.files[fileNumber])
	printLines(scanResult.lines, linePrint)

	return nil
}

type scanFileResult struct {
	lines         map[int]string
	lineMatch     []bool
	nMatchedLines int
}

func scanFile(fileNumber int, regex *regexp.Regexp, configs config) (scanFileResult, error) {
	f, err := os.Open(configs.files[fileNumber])
	if err != nil {
		return scanFileResult{}, fmt.Errorf("Ошибка при открытии файла: %v", err)
	}
	defer f.Close()

	lines := make(map[int]string)
	lineMatch := make([]bool, 0, 5)

	var lineNumber, nMatchedLines int

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		var match bool

		if configs.fixed {
			if configs.ignoreCase {
				match = strings.EqualFold(line, configs.pattern)
			} else {
				match = line == configs.pattern
			}
		} else {
			match = regex.MatchString(line)
		}

		if configs.lineNum {
			lines[lineNumber-1] = fmt.Sprintf("%d: %s", lineNumber, line)
		} else {
			lines[lineNumber-1] = line
		}

		if !match {
			lineMatch = append(lineMatch, false)
			continue
		}

		lineMatch = append(lineMatch, true)

		if configs.count {
			nMatchedLines++
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return scanFileResult{}, fmt.Errorf("в процессе работы сканера произошла ошибка: %v", err)
	}

	return scanFileResult{
			lines:         lines,
			lineMatch:     lineMatch,
			nMatchedLines: nMatchedLines,
		},
		nil
}

// convertLineMatchToLinePrintArray согласно флагам контекста, на основании булева массива соответсвующих строк возврщает булев массив с номерами строк, которые надо распечатать
func convertLineMatchToLinePrintArray(configs config, lineMatch []bool) []bool {
	linePrint := make([]bool, len(lineMatch))

	leftContext := maxOf2Int(configs.context, configs.before)
	rightContext := maxOf2Int(configs.context, configs.after)

	if leftContext == 0 && rightContext == 0 {
		return lineMatch
	}

	for i, v := range lineMatch {
		if v {
			j := maxOf2Int(i-leftContext, 0)
			k := minOf2Int(i+rightContext, len(linePrint))

			for ; j < k; j++ {
				linePrint[j] = true
			}
		}
	}
	return linePrint
}

func maxOf2Int(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minOf2Int(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printLines(lines map[int]string, linePrint []bool) {
	var linesForPrint strings.Builder

	for i, v := range linePrint {
		if v {
			linesForPrint.WriteString(lines[i])
			linesForPrint.WriteString("\n")
		}
	}

	fmt.Println(linesForPrint.String())

}
