package utils

import (
	"io"
	"net/http"
	"os"
)

func DownloadImage(URL string, filename string) error {
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Сохраняем изображение в файл
	_, err = io.Copy(out, resp.Body)
	return err
}
