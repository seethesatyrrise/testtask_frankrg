package rest

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"path/filepath"
	"testtask_frankrg/internal/files"
	"testtask_frankrg/internal/utils"
)

type Rest struct {
	currentPath string
	rootPath    string
	Template    *template.Template
}

func NewRest(rootPath string) *Rest {
	return &Rest{currentPath: rootPath, rootPath: rootPath, Template: utils.GetTemplate()}
}

func (r *Rest) Register(router *mux.Router) {
	router.HandleFunc("/", r.list)
	router.HandleFunc("/move/{path}", r.move)
	router.HandleFunc("/moveup", r.moveUp)
	//router.HandleFunc("/upload", r.upload)
	//router.HandleFunc("/delete", r.delete)
	//router.HandleFunc("/rename", r.rename)
	//router.HandleFunc("/download", r.download)
	//router.HandleFunc("/mkdir", r.mkdir)
}

func (r *Rest) move(w http.ResponseWriter, req *http.Request) {
	filePath, _ := mux.Vars(req)["path"]
	//fmt.Println("path: ", filePath)
	r.currentPath = filepath.Join(r.currentPath, filePath)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Rest) moveUp(w http.ResponseWriter, req *http.Request) {
	if r.currentPath != r.rootPath {
		r.currentPath = filepath.Dir(r.currentPath)
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Rest) list(w http.ResponseWriter, req *http.Request) {
	//fmt.Println("curpath: ", r.currentPath)
	list, statusCode, err := files.GetList(r.currentPath)
	if err != nil {
		PublishError(w, err, statusCode)
		return
	}

	PublishList(w, &list, r.Template)
}
