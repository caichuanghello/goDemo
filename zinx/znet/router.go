package znet

import "zinx/ziface"

type BaseRouter struct {
}
//baseRouter的方法不需要实现,我们的有些router方法不需要三个中的其中几个方法,所以只继承baseRouter,然后实现需要的方法就行了
func (b *BaseRouter)PerHandle(request ziface.IRequest){

}

func (b *BaseRouter)Handle(request ziface.IRequest){

}

func (b *BaseRouter)PostHandle(request ziface.IRequest){

}