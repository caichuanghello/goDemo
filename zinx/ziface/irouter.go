package ziface

type IRouter interface {
	//在处理conn业务之前的钩子方法
	PreHandle(request IRequest)
	//在处理conn业务的主方法hook
	Handle(request IRequest)
	//在处理conn业务之后的方法
	PostHandle(request IRequest)
}