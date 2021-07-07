package main

import (
	"fmt"
	"go_project/internal"
	"go_project/internal/middleware"
	"os"
	"os/signal"
	"syscall"
)

var httpHandler = middleware.App{}

func main() {

	fmt.Println(internal.MsgResponseStartProcess)

	addr := ":" + os.Getenv(internal.APP_PORT)

	//httpHandler := middleware.App{}
	_ = httpHandler.Initialize(internal.ValueEmpty, internal.ValueEmpty)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	_ = httpHandler.Run(addr)
}
