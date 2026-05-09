package main

import (
	"gin-admin-base/internal/interfaces/app"
)

func main() {
	application := app.New()
	defer application.Shutdown()

	application.Run()
}
