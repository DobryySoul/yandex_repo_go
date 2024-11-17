package main

import "github.com/DobryySoul/yandex_repo/internal/application"

func main() {
	app := application.New()
	// app.Run()
	app.RunServer()
}
