package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/lyj404/mini-distributed/log"
	"github.com/lyj404/mini-distributed/service"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	ctx, err := service.Start(context.Background(), "Log service", host, port, log.RegisterHandlers)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()

	fmt.Println("Shutting down log service. ")
}