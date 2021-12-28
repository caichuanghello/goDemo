package main

import (
	"fmt"
<<<<<<< HEAD
	"io"
=======
>>>>>>> 85b911c1ceaf61c4bdef68cdf20357c5fc288e3b
	"net"
)

func main(){
<<<<<<< HEAD
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
=======

	var s = 10
	for i:=0;i<20;i++ {
		if i==s {
			for j:=0;j<5;j++ {
				if j==2 {
					fmt.Println("j==2")
					break
				}
				fmt.Println(j)
			}
			fmt.Println("end..")
			break
		}
	}
	fmt.Println("运行这了..")

>>>>>>> 85b911c1ceaf61c4bdef68cdf20357c5fc288e3b
}


func handFun(conn net.Conn){
<<<<<<< HEAD
	defer conn.Close()
	buf :=make([]byte,1024)
	for {
		n,err:=conn.Read(buf)
		if err != nil || err != io.EOF{
			fmt.Println("conn read error :",err)
			return
		}
		if n==0{
			return
		}
		fmt.Println(string(buf[:n]))
	}
=======

>>>>>>> 85b911c1ceaf61c4bdef68cdf20357c5fc288e3b
}
