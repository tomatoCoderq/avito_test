package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/tomatoCoderq/avito_task/src/internal/app"
)

func main() {
	_ = godotenv.Load(".env")

	PORT, _ := strconv.Atoi(os.Getenv("PORT"))
	ADDRESS := os.Getenv("ADDRESS")

	address := ADDRESS + ":" + strconv.Itoa(PORT)

	application := app.New(PORT, address)
	go func() {
		application.HttpServer.MustRun()
	}()

	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	
	<-stop

	application.HttpServer.Stop()
}
