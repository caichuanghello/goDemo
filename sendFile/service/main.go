package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main(){

	lis,_:=net.Listen("tcp",":2222")
	defer lis.Close()
	for {
		conn,_:=lis.Accept()
		go handFunc(conn)
	}
}

func handFunc(c net.Conn) {
	//第一次接收的是文件名:大小
	fileinfo := make([]byte,100)
	n,_:=c.Read(fileinfo)
	fileName := string(fileinfo[:n])

	//告诉客服端,ok
	_,err:=c.Write([]byte("ok"))
	if err !=nil {
		fmt.Println("ok发送失败",err)
		return
	}
	//创建文件
	file,err:=os.OpenFile(fileName,os.O_RDWR|os.O_CREATE,0666)
	defer  file.Close()
	if err != nil {
		fmt.Println("文件创建失败!")
		return
	}
	//开始循环接收文件本身内容
	for {
		fileData :=make([]byte,1024*4)
		n,err:=c.Read(fileData)
		if err != nil {
			if err == io.EOF {
				fmt.Println("文件接收完成")
			} else {
				fmt.Println("文件上传失败",err)
			}
			return
		}
		//写入文件中
		file.Write(fileData[:n])
	}




}
