package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	// Парсим аргументы командной строки
	urlStr, dir, err := parseFlagsToString()
	if err != nil {
		fmt.Printf("Ошибка обработки флагов: %v", err)
	}

	if err := download(urlStr, dir); err != nil {
		fmt.Printf("Ошибка загрузки: %v", err)
	}
}

func parseFlagsToString() (string, string, error) {
	urlStr := flag.String("url", "", "URL сайта загрузки")
	dir := flag.String("dir", ".", "Directory to save the downloaded files")
	flag.Parse()

	if *urlStr == "" {
		return "", *dir, fmt.Errorf("флаг -url должен быть заполнен")
	}
	return *urlStr, *dir, nil
}

// Функция для загрузки ресурса по URL
func download(urlStr string, baseDir string) error {
	resp, err := http.Get(urlStr)
	if err != nil {
		return fmt.Errorf("ошибка доступа по адресу %s: %v", urlStr, err)
	}
	defer resp.Body.Close()

	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("ошибка обработки строки адреса %s: %v", urlStr, err)
	}

	// Определяем путь для сохранения файла
	filePath := path.Join(baseDir, u.Host, u.Path)
	if strings.HasSuffix(urlStr, "/") {
		filePath = path.Join(filePath, "index.html")
	}

	err = os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("ошибка создания сопутствующих папок для %s: %v", filePath, err)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла %s: %v", filePath, err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл %s: %v", filePath, err)
	}

	fmt.Printf("Загружено: %s -> %s\n", urlStr, filePath)

	// Если контент HTML, парсим его и загружаем связанные ресурсы
	if resp.Header.Get("Content-Type") == "text/html" {
		downloadAssets(resp.Body, u, baseDir)
	}

	return nil
}

// Функция для загрузки ресурсов, связанных с HTML
func downloadAssets(body io.Reader, base *url.URL, baseDir string) {
	tokenizer := html.NewTokenizer(body)

	for {
		tt := tokenizer.Next()

		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			t := tokenizer.Token()

			var src string
			for _, attr := range t.Attr {
				if attr.Key == "src" || attr.Key == "href" {
					src = attr.Val
					break
				}
			}

			if src == "" {
				continue
			}

			assetURL, err := base.Parse(src)
			if err != nil {
				fmt.Printf("error parsing asset URL %s: %v\n", src, err)
				continue
			}

			// Рекурсивная загрузка ресурсов
			err = download(assetURL.String(), baseDir)
			if err != nil {
				fmt.Printf("Ошибка при загрузке %s: %v\n", assetURL.String(), err)
			}
		}
	}
}
