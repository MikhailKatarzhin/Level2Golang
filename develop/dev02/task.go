package main

import (
	"fmt"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
  - "a4bc2d5e" => "aaaabccddddde"
  - "abcd" => "abcd"
  - "45" => "" (некорректная строка)
  - "" => ""

Дополнительное задание: поддержка escape - последовательностей

  - qwe\4\5 => qwe45 (*)

  - qwe\45 => qwe44444 (*)

  - qwe\\5 => qwe\\\\\ (*)

    Комментарий от стажёра - так как escape - последовательности указаны только для цифр и '\' , то в случаях любых иных символов функция вернёт ошибку

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
func StringPrimitiveDecoder(str string) (string, error) {
	var unpackedStr strings.Builder
	var previousR rune

	escapeMode := false
	repeatingEscape := false

	for i, r := range str {
		if escapeMode {
			if r == '\\' || unicode.IsDigit(r) {
				unpackedStr.WriteRune(r)
			} else {
				return "", fmt.Errorf("некорректная строка: неизвестная escape-последовательность при обрбаботке %dго символа", i)
			}
			escapeMode = false
		} else if r == '\\' {
			escapeMode = true
			repeatingEscape = true
		} else if unicode.IsDigit(r) {

			if i == 0 || unicode.IsDigit(previousR) {
				if repeatingEscape {
					repeatingEscape = false
				} else {
					return "", fmt.Errorf("некорректная строка: ошибка при обрбаботке %dго символа", i)
				}
			}

			for j := r - '0'; j > 1; j-- {
				unpackedStr.WriteRune(previousR)
			}

		} else {
			unpackedStr.WriteRune(r)

			if repeatingEscape {
				repeatingEscape = false
			}
		}
		previousR = r
	}

	if escapeMode {
		return "", fmt.Errorf("некорректная строка: незаконченная escape - последовательность")
	}

	return unpackedStr.String(), nil
}

func main() {
	testCases := []string{
		"a4bc2d5e",  // aaaabccddddde
		"abcd",      // abcd
		"45",        // error
		"",          // ""
		"qwe\\4\\5", // qwe45
		"qwe\\45",   // qwe44444
		"qwe\\\\5",  // qwe\\\\\
	}

	for _, testCase := range testCases {
		result, err := StringPrimitiveDecoder(testCase)
		fmt.Printf("\"%s\" => \"%s\"\n", testCase, result)

		if err != nil {
			fmt.Println("Ошибка в предыдущей строке:", err)
		}
	}
}
