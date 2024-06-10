package main

import (
	"log/slog"
	_ "messenger/jobs"
	"messenger/web"
)

func main() {
	slog.Info("Initialization complete. Starting the web server.")
	web.Run()
	select {}
}
