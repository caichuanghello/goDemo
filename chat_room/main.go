package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	C net.Conn
	Name string
	Addr string
}
//存储所有在线的用户信息
var onlineMap = make(map[string]Client)

func main(){

	li,err:=net.Listen("tcp",":9000")
	if err != nil {
		fmt.Println("listen error :",err)
		return
	}
	defer li.Close()
	for {
		conn,err:=li.Accept()
		if err !=nil{
			fmt.Println("accept error :",err)
			return
		}
		go handFun(conn)
	}
}

func handFun(conn net.Conn){

	fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" ["+conn.RemoteAddr().String() + "] 以上线")
	defer conn.Close()

	//每个连接的唯一标识
	name :=conn.RemoteAddr().String() //默认用户名ip+port
	clnt :=Client{
		C : conn,
		Name: name,
		Addr: name,
	}

	//添加到全局在线的map,广播通知
	onlineMap[name] = clnt
	msg := time.Now().Format("2006-01-02 15:04:05")+" : ["+conn.RemoteAddr().String() + "] 进入聊天室\n"
	//通知所有的人有新用户上线
	go sendMsgToAllClient(name,msg)

	//这是获取用户的输入,然后广播给所有人
	var buf = make([]byte,1024)
	//超时标记
	var isOnline = make(chan bool)

	go func(){
		for{
			//循环读取用户输入,告知所有用户
			n,err:=conn.Read(buf)
			//用户不是命令退出,是直接关掉客户端退出的
			if n == 0 {
				doSomeThingWhenClose(name,0)
				return
			}

			if err != nil {
				fmt.Println("conn read error :",err)
				break
			}
			inputByte := string(buf[:n])

			//主动断开连接命令
			if inputByte == "exit\n"{ //如果输入了exit,则直接断开连接,因为接收到的最后会带一个\n,所以去掉自后一个\n
				doSomeThingWhenClose(name,0)
				return
			}

			//修改用户昵称命令 rename xxx
			if strings.HasPrefix(inputByte,"rename") {
				nickename := strings.TrimSpace(strings.TrimLeft(inputByte,"rename")) //截取昵称,去掉空格
				if len(nickename) == 0 || len(nickename) >20 {
					onlineMap[name].C.Write([]byte("昵称不能为空或者昵称太长..\n"))
				} else{
					//map修改需要等个拿出来在修改,不能单独修改
					client := onlineMap[name]
					client.Name = nickename
					onlineMap[name] = client
				}
			} else{
				go sendMsgToAllClient(name,"["+onlineMap[name].Name+"]"+inputByte)
			}
			//维持心跳
			isOnline<-true
		}
	}()

	for {
		select {
		case <-isOnline:
		case <-time.After(10*time.Second):
			doSomeThingWhenClose(name,1)
			return
		}
	}
}

//给所有
func sendMsgToAllClient(name string,msg string){

	for k,cli:=range onlineMap{

		if k == name { //不用给自己发送上线通知
			continue
		}
		_,err:=cli.C.Write([]byte(msg))
		if err != nil {
			fmt.Println("write error :",err)
			continue
		}
	}
}

//当用户下线的时候,需要执行的操作,删除map数据,服务端打印
func doSomeThingWhenClose(name string,isTellSelf int){
	_, ok := onlineMap[name]
	if ok {
		//控制台打印输出
		fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" : ["+onlineMap[name].Name+"] 断开连接")
		//给自己发送断开信息
		if isTellSelf == 1 {
			onlineMap[name].C.Write([]byte("您已超时,已断开连接..."))
		} else {
			onlineMap[name].C.Close()//直接关闭,不是通过defer来关闭
		}
		delete(onlineMap,name)
	}
}

