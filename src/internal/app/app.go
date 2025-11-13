package app

import (
	httpApp "github.com/tomatoCoderq/avito_task/src/internal/app/http"
)

type App struct {
	HttpServer *httpApp.App
}

func New(
	port int,
	address string,
	) *App {
		httpApp := httpApp.New(port, address)

	return &App{
		HttpServer: httpApp,
	}
}
