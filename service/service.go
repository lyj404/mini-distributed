package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/lyj404/mini-distributed/registry"
)

// 服务注册
func Start(ctx context.Context, host, port string, reg registry.Registration,
	 registerHandlersFunc func()) (context.Context, error){
	// 运行注册函数
	registerHandlersFunc()

	// 启动service
	ctx = startService(ctx, reg.ServiceName, host, port)

	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	
	return ctx, nil
}

// 服务启动
func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context{
	// 创建一个具有取消功能的上下文
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = ":" + port

	fmt.Printf("HTTP Server listening on port %s ...\n", port)
	
	// 异常终止
	go func ()  {
		// 启动HTTP服务，如果发生错误则退出
		log.Println(srv.ListenAndServe())
		err := registry.ShutdownService(fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			log.Panicln(err)
		}
		cancel()
	}()
	
	// 用户终止
	go func() {
		fmt.Printf("%v started. Press any key to stop.\n", serviceName)
		// 等待用户输入
		var s string
		fmt.Scanln(&s)

		err := registry.ShutdownService(fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			log.Panicln(err)
		}	

		// 手动停止发送信号
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}