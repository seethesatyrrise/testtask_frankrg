package files

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func GetList(path string) ([]File, int, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, http.StatusNotFound, errors.New("path doesn't exist")
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var list []File
	for _, file := range files {
		if info, err := file.Info(); err == nil {
			list = append(list, File{
				Name:  info.Name(),
				IsDir: info.IsDir(),
			})
		}
	}
	fmt.Println(list)

	return list, http.StatusOK, nil
}

func MkDir(path string) error {
	err := os.Mkdir(path, 0777)
	if err != nil {
		return err
	}
	return nil
}

func Rename(oldPath, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}
	return nil
}
