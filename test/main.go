package main

import "zinx/znet"

func main(){

	s:=znet.NewServer("ewServer(\"[zinx v0.1]\")")

	s.Server()

}