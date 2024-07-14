package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/lyj404/mini-distributed/grades"
	"github.com/lyj404/mini-distributed/registry"
	"github.com/lyj404/mini-distributed/service"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration {
		ServiceName: registry.GradingService,
		ServiceURL: serviceAddress,
	}

	ctx, err := service.Start(context.Background(), host, port, r, grades.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	<- ctx.Done()
	fmt.Println("Shutting down grading servie")
}