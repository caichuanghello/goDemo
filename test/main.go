package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
)

func unimplemented(conn net.Conn){
	var buf string

	buf = "HTTP/1.0 501 Method Not Implemented\r\n"
	_, _ = conn.Write([]byte(buf))
	buf = "Server: httpd/0.1.0\r\n"
	_, _ = conn.Write([]byte(buf))
	buf = "Content-Type: text/html\r\n"
	_, _ = conn.Write([]byte(buf))
	buf = "\r\n"
	_, _ = conn.Write([]byte(buf))
	buf = "<HTML><HEAD><TITLE>Method Not Implemented\r\n"
	_, _ = conn.Write([]byte(buf))
	buf = "</TITLE></HEAD>\r\n"
	_, _ = conn.Write([]byte(buf))
	buf = "<BODY><P>HTTP request method not supported.\r\n"
	_, _ = conn.Write([]byte(buf))
	buf = "</BODY></HTML>\r\n"
	_, _ = conn.Write([]byte(buf))
}

func accept_request_thread(conn net.Conn)  {
	defer conn.Close()
	var i  int

	buf := make([]byte, 1024)
	n, err := conn.Read(buf) // 从conn中读取客户端发送的数据内容
	if err != nil {
		fmt.Printf("客户端退出 err=%v\n", err)
		return
	}
	fmt.Println("buf:",string(buf))

	// 获取方法
	i = 0
	var method_bt strings.Builder
	for(i < n && buf[i] != ' '){
		method_bt.WriteByte(buf[i])
		i++
	}
	fmt.Println("buf[i]:",string(buf[:i]))
	method := method_bt.String()


	if(method != "GET"){
		unimplemented(conn)
		return
	}


	for(i < n && buf[i] == ' '){
		i++
	}

	//api/camera/get_ptz?camera_id=1324566666789876543
	var url_bt strings.Builder
	for(i < n && buf[i] != ' '){
		url_bt.WriteByte(buf[i])
		i++
	}
	url := url_bt.String()

	if(method == "GET"){
		//url ---> /api/camera/get_ptz?camera_id=1324566666789876543
		// 跳到第一个？
		var path, query_string string
		j := strings.IndexAny(url, "?")
		if(j != -1){
			path = url[:j]
			if(j + 1 < len(url)){
				query_string = url[j+1:]
			}
		}else{
			path = url
		}

		fmt.Print(path + "请求已经创建\t")
		resp := execute(path, query_string)// =1324566666789876543
		fmt.Println("返回", string(resp))
		header(conn, "application/json", len(resp))
		_ , err := conn.Write(resp)
		if(err != nil){
			fmt.Println(err)
		}
	}
}

//回应客户端必须先设置好head头，浏览器才能解析
func header(conn net.Conn, content_type string , length int )  {
	var buf string

	buf = "HTTP/1.0 200 OK\r\n"
	_, _ = conn.Write([]byte(buf))
	buf = "Server: httpd/0.1.0\r\n"
	_, _ = conn.Write([]byte(buf))
	buf = "Content-Type: " + content_type + "\r\n"
	_, _ = conn.Write([]byte(buf))
	_, _ = fmt.Sscanf(buf, "Content-Length: %d\r\n", length)
	_, _ = conn.Write([]byte(buf))
	buf = "\r\n"
	_, _ = conn.Write([]byte(buf))
}



func  execute(path string, query_string string) ([]byte)  {
	query_params := make(map[string]string)
	parse_query_params(query_string, query_params)

	if("/api/camera/get_ptz" == path){
		/*
		 * do something
		 */
		camera_id := query_params["camera_id"]

		resp := make(map[string]interface{})
		resp["camera_id"] = camera_id
		resp["code"] = 200
		resp["msg"] = "ok"

		rs, err := json.Marshal(resp)
		if err != nil{
			log.Fatalln(err)
		}
		return rs
	}else if("get_abc" == path){
		/*
		 * do something
		 */

		return []byte("abcdcvfdswa")
	}


	return []byte("do't match")
}

/*map作为函数入参是作为指针进行传递的
函数里面对map进行修改时，会同时修改源map的值，但是将map修改为nil时，则达不到预期效果。*/
// camera_id=1324566666789876543&tt=%E5%88%9B%E5%BB%BA%E6%88%90%E5%8A%9F
func parse_query_params(query_string string, query_params map[string]string)  {
	kvs := strings.Split(query_string, "&")
	if(len(kvs) == 0){
		return
	}

	for _, kv := range kvs {
		kv := strings.Split(kv, "=")
		if(len(kv) != 2){
			continue
		}
		query_params[kv[0]] = kv[1]
	}
}
func main01() {

	listen, err := net.Listen("tcp", ":8888")  // 创建用于监听的 socket
	if err != nil {
		fmt.Println("listen err=", err)
		return
	}
	defer listen.Close()  // 服务器结束前关闭 listener

	// 循环等待客户端链接
	for   {
		fmt.Println("阻塞等待客户端链接...")
		conn, err := listen.Accept() // 创建用户数据通信的socket
		if err != nil {
			panic("Accept() err=  " +  err.Error())
		}
		// 这里准备起一个协程，为客户端服务
		go accept_request_thread(conn)
	}
}
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func main(){




	//var (
	//	name    string
	//	age     int
	//)
	//fmt.Println("请输入名称:")
	//fmt.Scan(&name)
	//fmt.Println("请输入年龄:")
	//fmt.Scan(&age)
	//
	//fmt.Printf("扫描结果 name:%s age:%d \n", name, age)





	//return
	//
	//lis,_:=net.Listen("tcp",":8888")
	//defer lis.Close()
	//
	//for {
	//	con,err:=lis.Accept()
	//	if err != nil {
	//		panic("Accept() err=  " +  err.Error())
	//	}
	//	go handleFunc(con)
	//}

}




func handleFunc(con net.Conn){

	defer con.Close()

	buf := make([]byte, 1024)
	con.Read(buf) // 从conn中读取客户端发送的数据内容


	resp :=`<h1>Hello World!</h1>`

	header(con,"text/html",len(resp))

	con.Write([]byte(resp))


	return

}

