package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
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


	file,err:=os.OpenFile("D:\\vue\\vue学习第二天.flv",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	defer file.Close()
	if err !=nil{
		fmt.Println("文件打开失败")
	}
	for i:=1;i<10;i++{
		filename := "D:\\vue\\"+strconv.Itoa(i)+".flv"
		file1,err:= os.Open(filename)
		if err !=nil{
			fmt.Println(filename,"not found")
			break
		}
		res,err:=ioutil.ReadAll(file1)
		n,err:=file.Write(res)

		if err != nil {
			fmt.Println(filename,"合并失败")
			break
		}
		fmt.Println(filename,"success:",n)
		file1.Close()
	}

	fmt.Println("合并完成")



	//file,_:=os.Open("log.txt")
	//
	//reader:=bufio.NewReader(file)
	//
	//rem,_:=reader.ReadBytes(',')
	//
	//fmt.Println(string(rem))


	/*
	//获取当前时间
	tim := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("当前时间:",tim)


	//当前时间戳(秒)
	timestamp:=time.Now().Unix()
	fmt.Println("时间戳:",timestamp)

	//把一个时间戳类型转化成时间格式
	times := int64(1555555555)
	timesString:=time.Unix(times,0).Format("2006-01-02 15:04:05")
	fmt.Println("时间戳为",times,"的日期为:",timesString)

	//把时间格式的字符串转化成时间戳
	tims :="2019-04-18 10:45:55"
	loc,_:=time.LoadLocation("Asia/Shanghai")
	ti,_:=time.ParseInLocation("2006-01-02 15:04:05",tims,loc)
	tis:=ti.Unix()
	fmt.Println(tims,"对应的是时间戳是:",tis)

	//比较两个时间
	t1:=time.Now()
	t2,_:= time.Parse("2006-01-02 15:04:05","2022-02-03 15:12:11")
	fmt.Println("t1是否在t2之后:",t1.After(t2))
	fmt.Println("t1是否在t2之前:",t1.Before(t2))
	fmt.Println("t1与t2的间隔:",int((t2.Sub(t1)).Seconds()))

	//日期相加
	fmt.Println("后天是:",time.Now().AddDate(0,0,2).Format("2006-01-02 15:04:05"))

	fmt.Println("过24小时一分钟之后秒是:",time.Now().Add(time.Second*60+time.Hour*24).Format("2006-01-02 15:04:05"))

	 */

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



