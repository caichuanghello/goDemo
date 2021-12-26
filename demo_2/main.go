package main

import (
	"fmt"
	"net"
)

func main(){

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

}


func handFun(conn net.Conn){

}
