package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/lyj404/mini-distributed/grades"
	"github.com/lyj404/mini-distributed/log"
	"github.com/lyj404/mini-distributed/registry"
	"github.com/lyj404/mini-distributed/service"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration {
		ServiceName: registry.GradingService,
		ServiceURL: serviceAddress,
		RequiredServices: []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartBeatURL: serviceAddress + "/heartbeat",
	}

	ctx, err := service.Start(context.Background(), host, port, r, grades.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}

	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %v\n",  logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}

	<- ctx.Done()
	fmt.Println("Shutting down grading servie")
}