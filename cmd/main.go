package main

import (
	"fmt"
	"go_project/internal/middleware"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	fmt.Println("Start process -->   ")

	//errors := make(chan error)

	addr := ":9090"

	httpHandler := middleware.App{}
	_ = httpHandler.Initialize("", "")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	_ = httpHandler.Run(addr)
	//errors <- fmt.Errorf("%s", <-ch)

	//go func() {
	//	errors <- httpHandler.Run(addr)
	//}()
}
