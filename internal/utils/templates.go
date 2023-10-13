package utils

import (
	"html/template"
)

func GetTemplate() *template.Template {
	t, err := template.ParseFiles("internal/utils/index.html")
	if err != nil {
		Logger.Error(err.Error())
	}
	return t
}
