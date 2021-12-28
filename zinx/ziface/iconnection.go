package ziface

import "net"

//定义连接模块的 抽象层

type IConnection interface {
	//启动连接,让当前的连接准备开始工作
	Start()
	//停止连接,结束当前连接的工作
	Stop()
	//获取客户端的tcp状态IP port
	GetTCPConnection() *net.Conn

	GetConnID() uint32

	RemoteAddr() net.Addr

	//发送数据
	Send(data []byte) error

}

//定义一个处理连接业务的方法
type HandleFunc func(net.Conn,[]byte,int) error