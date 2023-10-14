package rest

import (
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"net/http"
	"os"
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
	files.CreateRootIfNotExist(rootPath)
	return &Rest{currentPath: rootPath, rootPath: rootPath, Template: utils.GetTemplate()}
}

func (r *Rest) Register(router *mux.Router) {
	router.HandleFunc("/", r.list)
	router.HandleFunc("/move/{path}", r.move)
	router.HandleFunc("/moveup", r.moveUp)
	router.HandleFunc("/mkdir/{name}", r.mkDir)
	router.HandleFunc("/rename/{old_name}/{new_name}", r.rename)
	router.HandleFunc("/remove/{name}", r.remove)
	router.HandleFunc("/upload", r.upload).Methods("POST")
	router.HandleFunc("/download/{name}", r.download)
}

func (r *Rest) move(w http.ResponseWriter, req *http.Request) {
	filePath, _ := mux.Vars(req)["path"]
	r.currentPath = filepath.Join(r.currentPath, filePath)
	utils.Logger.Info("go into " + filePath)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Rest) moveUp(w http.ResponseWriter, req *http.Request) {
	if r.currentPath != r.rootPath {
		r.currentPath = filepath.Dir(r.currentPath)
		utils.Logger.Info("go up")
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Rest) mkDir(w http.ResponseWriter, req *http.Request) {
	name, _ := mux.Vars(req)["name"]
	if name != "" {
		err := files.MkDir(filepath.Join(r.currentPath, name))
		if err != nil {
			utils.Logger.Error(err.Error())
		} else {
			utils.Logger.Info("make dir " + name)
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
			utils.Logger.Error(err.Error())
		} else {
			utils.Logger.Info("rename " + oldName + " into " + newName)
		}
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Rest) remove(w http.ResponseWriter, req *http.Request) {
	name, _ := mux.Vars(req)["name"]
	if name != "" {
		err := files.Remove(filepath.Join(r.currentPath, name))
		if err != nil {
			utils.Logger.Error(err.Error())
		} else {
			utils.Logger.Info("remove " + name)
		}
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Rest) upload(w http.ResponseWriter, req *http.Request) {
	file, handler, err := req.FormFile("newFile")
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	defer file.Close()

	f, err := os.OpenFile(filepath.Join(r.currentPath, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	defer f.Close()
	_, _ = io.Copy(f, file)
	utils.Logger.Info("upload file " + handler.Filename)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Rest) download(w http.ResponseWriter, req *http.Request) {
	name, _ := mux.Vars(req)["name"]
	if name != "" {
		path := filepath.Join(r.currentPath, name)
		w.Header().Set("Content-Disposition", "attachment; filename="+path)
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, req, path)
		utils.Logger.Info("download file " + name)
	}
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Rest) list(w http.ResponseWriter, req *http.Request) {
	list, statusCode, err := files.GetList(r.currentPath)
	if err != nil {
		PublishError(w, err, statusCode)
		return
	}
	PublishList(w, &list, r.Template)
}
