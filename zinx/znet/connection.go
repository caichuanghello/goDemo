package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn net.Conn
	ConnID uint32
	isClosed bool
	handlAPI ziface.HandleFunc

	//告知当前连接已近退出/停止 channel
	ExitChan chan bool

	Router ziface.IRouter
}

func NewConnection(conn net.Conn,connID uint32, callback ziface.HandleFunc) *Connection {
	c :=&Connection{
		Conn: conn,
		ConnID: connID,
		handlAPI: callback,
		isClosed: false,
		ExitChan: make(chan bool,1),
	}

	return c
}



func (c *Connection)StartReader(){
	fmt.Println("READER GOROUTER IS RUNNING,...")
	defer fmt.Println("connID =",c.GetConnID())
	defer c.Stop()
	defer c.Conn.Close()

	for {
		res :=make([]byte,512)
		n,err:=c.Conn.Read(res)
		if err !=nil {
			fmt.Println("recv buf err",err)
			continue
		}

		//从路由中找到注册绑定的con对应的router调用
		req :=Request{
			conn: c,
			data: res,
		}
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}


func (c *Connection)Start(){

	fmt.Println("Conn start()... connID = ",c.GetConnID())

	//启动从当前连接中读取数据的业务
	go c.StartReader()

	//启动从当前连接中写数据的业务

	//go c.StartWriter()




}

func (c *Connection)Stop(){
	fmt.Println("Conn stop()... ConnID =",c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	close(c.ExitChan)
}

func (c *Connection)GetTCPConnection() *net.Conn {
	return &c.Conn
}

func (c *Connection)GetConnID() uint32{
	return c.ConnID
}
func (c *Connection)RemoteAddr() net.Addr{
	return c.Conn.RemoteAddr()
}
func (c *Connection)Send() error{
	return nil
}

