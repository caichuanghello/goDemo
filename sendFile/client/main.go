package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main(){

	fmt.Println("请输入需要上传的文件:")
	var filePath string
	fmt.Scan(&filePath)
	file,err:= os.OpenFile(filePath,os.O_RDONLY,0666)
	defer file.Close()
	if  err != nil {
		fmt.Println(filePath,"文件打开失败")
		return
	}
	fileInfo,_ := file.Stat()

	conn,err:=net.Dial("tcp","127.0.0.1:2222")
	defer conn.Close()
	if err != nil{
		fmt.Println("上传服务暂未启动")
		return
	}
	_,err=conn.Write([]byte(fileInfo.Name()))

	if err != nil {
		fmt.Println("文件发送失败")
		return
	}
	//接收服务器返回的信息
	res :=make([]byte,8)
	n,_:=conn.Read(res)

	if "ok" == string(res[:n]) {
		fmt.Println("开始发送文件...")
		//发送文件本身
		for {
			sendData :=make([]byte,1024*4) //每次发送4k大小的
			n,err:=file.Read(sendData) //读取文件内容
			if err!= nil {
				if err == io.EOF {
					fmt.Println("文件发送完成")
				} else{
					fmt.Println("文件读取失败",err)
				}
				return
			}
			conn.Write(sendData[:n]) //发送给服务器
		}
	} else {
		fmt.Println("文件发送失败")
		return
	}










}
