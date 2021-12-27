package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/ziface"
)

//IServer接口实现
type Server struct {
	Name string
	IPVersion string
	IP string
	Port int
}
//这边目前写死,以后优化可以由
func CallBackToClient(conn net.Conn,data []byte,cnt int) error{
	if _,err:=conn.Write(data[:cnt]); err !=nil{
		return errors.New("CallBackToClient error")
	}

	return nil
}

func(s *Server)Start(){
	fmt.Printf("[start] Server Lister at IP: %s,port %d",s.IP, s.Port)

	go func() {
		lis,err:=net.Listen(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))

		if err != nil {
			fmt.Println("listen err:",err)
			return
		}

		for {
			con,err:=lis.Accept()
			if err != nil {
				fmt.Println("accept error :",err)
				continue
			}
			go handleFunc(con)
		}
	}()


}


func handleFunc(con net.Conn){

	var connID uint32
	connID =0

	delconn:=NewConnection(con,connID,CallBackToClient)
	connID++

	//启动当前的连接处理业务
	go delconn.Start()

	//for {
	//	res :=make([]byte,512)
	//	n,err:=con.Read(res)
	//	if err !=nil {
	//		fmt.Println("recv buf err",err)
	//		continue
	//	}
	//	con.Write(res[:n])
	//}


}

func(s *Server)Stop(){

}

func(s *Server)Serve(){
	s.Start()

	select {

	}
}


func NewServer(name string) ziface.IServer{
	s :=&Server{
		Name: name,
		IP: "0.0.0.0",
		Port: 8999,
		IPVersion: "tcp4",
	}

	return s
}
