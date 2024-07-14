package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/lyj404/mini-distributed/log"
	"github.com/lyj404/mini-distributed/portal"
	"github.com/lyj404/mini-distributed/registry"
	"github.com/lyj404/mini-distributed/service"
)

func main() {
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}

	host, port := "localhost", "5000"
	serviceAdress := fmt.Sprintf("http://%s:%s", host, port)
	
	r := registry.Registration{
		ServiceName: registry.PortalService,
		ServiceURL: serviceAdress,
		RequiredServices: []registry.ServiceName{
			registry.LogService,
			registry.GradingService,
		},
		ServiceUpdateURL: serviceAdress + "/services",
		HeartBeatURL: serviceAdress + "/heartbeat",
	}

	ctx, err := service.Start(context.Background(), host, port, r, portal.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err != nil {
		log.SetClientLogger(logProvider, r.ServiceName)
	}
	<- ctx.Done()
	fmt.Println("Shutting down portal service")
}