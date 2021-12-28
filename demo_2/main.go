package main

import (
	"fmt"

	"io"

	"net"
)

func main(){

	lis,err:=net.Listen("tcp",":9000")
	if err !=nil {
		panic(err)
	}
	defer lis.Close()
	for {
		conn,err:=lis.Accept()
		if err !=nil{
			fmt.Println("accept error:",err)
		}
		go handFun(conn)
	}

}


func handFun(conn net.Conn){

	defer conn.Close()
	buf :=make([]byte,1024)
	for {
		n,err:=conn.Read(buf)
		if err != nil && err != io.EOF{
			fmt.Println("conn read error :",err)
			return
		}
		if n==0{
			return
		}
		fmt.Println(string(buf[:n]))
	}

}
