package rest

import (
	"html/template"
	"net/http"
	"testtask_frankrg/internal/files"
)

func PublishError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}

func PublishList(w http.ResponseWriter, list *[]files.File, template *template.Template) {
	template.Execute(w, *list)
}
