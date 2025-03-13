package utils

import (
	"errors"
	"mime/multipart"
)

func CheckContentType(file *multipart.FileHeader, contentTypes ...string) error {
	if len(contentTypes) > 0 {
		for _, contentType := range contentTypes {
			contentTypeFile := file.Header.Get("Content-Type")
			if contentTypeFile == contentType {
				return nil
			}
		}

		return errors.New("not allowed type file")
	} else {
		return errors.New("not found content type to be checking")
	}
}
