package app

import (
	"context"
	"github.com/gorilla/mux"
	"path/filepath"
	"testtask_frankrg/internal/rest"
	"testtask_frankrg/internal/server"
)

type App struct {
	router *mux.Router
	server *server.Server
	cfg    *AppConfig
}

func New(ctx context.Context) (app *App, err error) {
	app = &App{}
	app.cfg, err = newConfig()
	if err != nil {
		return nil, err
	}

	rootPath, _ := filepath.Abs("./root")
	rest := rest.NewRest(rootPath)

	router := mux.NewRouter()
	rest.Register(router)

	app.router = router

	app.server, err = server.NewServer(&app.cfg.ServerConfig, app.router)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run(ctx context.Context) (err error) {
	if err := app.server.Start(); err != nil {
		return err
	}
	return nil
}

func (app *App) Shutdown(ctx context.Context) error {
	if app.server != nil {
		if err := app.server.Shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}
