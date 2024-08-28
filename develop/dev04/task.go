package main

import (
	"fmt"
	"slices"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func findAnagramsFromSlice(words []string) map[string][]string {
	anagrams := make(map[string][]string)

	// checkedWords для хранения уже считанных слов с целью поддержания их уникальности
	checkedWords := make(map[string]struct{})

	for _, word := range words {
		// Приводим слово к нижнему регистру
		lowerWord := strings.ToLower(word)

		if _, ok := checkedWords[lowerWord]; !ok {
			checkedWords[lowerWord] = struct{}{}

			arr := []byte(lowerWord)
			slices.SortFunc(arr, func(a, b byte) int {
				if a < b {
					return -1
				} else if a > b {
					return 1
				}
				return 0
			})

			anagrams[string(arr)] = append(anagrams[string(arr)], lowerWord)
		}
	}

	// Генерирация итоговой мапы, соответствующуй заданию
	result := make(map[string][]string)

	for _, group := range anagrams {
		if len(group) > 1 {
			// Сортировка группы
			slices.Sort(group)
			// Закрепление первого слова из отсортированной группы ключом группы в словваре
			result[group[0]] = group
		}
	}

	return result
}

func main() {
	words := []string{"пятак", "столик", "пяткА", "тяпка", "тяпка", "тЯпка", "листок", "слиток", "слово", "волос", "ослов", "зола"}

	anagramGroups := findAnagramsFromSlice(words)

	for key, group := range anagramGroups {
		fmt.Println(key, ":", group)
	}
}
