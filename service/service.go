package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// 服务注册
func Start(ctx context.Context, serviceName, host, port string,
	 registerHandlersFunc func()) (context.Context, error){
	// 运行注册函数
	registerHandlersFunc()

	// 启动service
	ctx = startService(ctx, serviceName, host, port)

	return ctx, nil
}

// 服务启动
func startService(ctx context.Context, serviceName, host, port string) context.Context{
	// 创建一个具有取消功能的上下文
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = ":" + port

	fmt.Printf("HTTP Server listening on port %s ...\n", port)
	
	// 异常终止
	go func ()  {
		// 启动HTTP服务，如果发生错误则退出
		log.Println(srv.ListenAndServe())
		cancel()
	}()
	
	// 用户终止
	go func() {
		fmt.Printf("%v started. Press any key to stop.\n", serviceName)
		// 等待用户输入
		var s string
		fmt.Scanln(&s)

		// 手动停止发送信号
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}