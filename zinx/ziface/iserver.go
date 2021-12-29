package ziface

type IServer interface {

	Start()

	Stop()

	Serve()

	//给当前的服务注册路由方法,供客服端连接处理使用
	AddRouter(router IRouter)
}
