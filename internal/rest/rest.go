package rest

import (
	"fmt"
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
	router.HandleFunc("/mkdir/{name}", r.mkDir)
	router.HandleFunc("/rename/{old_name}/{new_name}", r.rename)
	//router.HandleFunc("/upload", r.upload)
	//router.HandleFunc("/delete", r.delete)
	//router.HandleFunc("/download", r.download)

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

func (r *Rest) mkDir(w http.ResponseWriter, req *http.Request) {
	name, _ := mux.Vars(req)["name"]
	if name != "" {
		err := files.MkDir(filepath.Join(r.currentPath, name))
		if err != nil {
			fmt.Println(err)
		}
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Rest) rename(w http.ResponseWriter, req *http.Request) {
	oldName, _ := mux.Vars(req)["old_name"]
	newName, _ := mux.Vars(req)["new_name"]
	if oldName != "" && newName != "" {
		err := files.Rename(filepath.Join(r.currentPath, oldName), filepath.Join(r.currentPath, newName))
		if err != nil {
			fmt.Println(err)
		}
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
