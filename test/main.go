package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
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

var wg sync.WaitGroup


//
func main(){

	start := time.Now()
	var inputDir string
	var outputName string
	flag.StringVar(&inputDir, "i", "./dvr-file/flv", "input file dir")
	flag.StringVar(&outputName, "o", "out.mp4", "output file name")
	//解析命令行参数
	flag.Parse()

	exist, err := pathExists(inputDir)
	if err != nil {
		fmt.Printf("get dir error!: %v", err)
		return
	}
	if !exist {
		inputDir = os.Args[0]
	}
	inputDir, _  = filepath.Abs(inputDir)
	fmt.Println("argv: ", inputDir, outputName)
	if err = flvs2mp4(inputDir, outputName); err != nil {
		fmt.Printf("flv to mp4 error!: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Println("Running time:", elapsed)
}

// 判断文件或目录是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 命令行调用
func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	//fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

// 视频格式转换
func videoConvert(in string, out string) {
	defer wg.Done()
	//fmt.Println(in, out)
	cmdStr := fmt.Sprintf("ffmpeg -i %s -loglevel quiet -c copy -bsf:v h264_mp4toannexb -f mpegts %s", in, out)
	args := strings.Split(cmdStr, " ")
	msg, err := Cmd(args[0], args[1:])
	if err != nil {
		fmt.Printf("videoConvert failed, %v, output: %v\n", err, msg)
		return
	}
}

// 视频合成
func videoMerge(in []string, out string) {
	//fmt.Println(in, out)
	cmdStr := fmt.Sprintf("ffmpeg -i concat:%s -loglevel quiet -c copy -absf aac_adtstoasc -movflags faststart %s",
		strings.Join(in, "|"), out)
	args := strings.Split(cmdStr, " ")
	msg, err := Cmd(args[0], args[1:])
	if err != nil {
		fmt.Printf("videoMerge failed, %v, output: %v\n", err, msg)
		return
	}
}

func flvs2mp4(inDir string, outFile string)(err error) {
	tsFileDir := filepath.Join(inDir, "tsfile")
	if err = os.RemoveAll(tsFileDir); err != nil {
		return
	}
	if err = os.RemoveAll(outFile); err != nil {
		return
	}
	if err = os.Mkdir(tsFileDir,0666); err!=nil {
		return
	}

	infiles, _ := ioutil.ReadDir(inDir)
	for _, f := range infiles {
		if filepath.Ext(f.Name()) == ".flv" {
			tsfileName := filepath.Join(tsFileDir, strings.TrimSuffix(f.Name(), ".flv") + ".ts")
			wg.Add(1)
			go videoConvert(filepath.Join(inDir, f.Name()), tsfileName)
		}
	}
	wg.Wait()

	tsfiles, _ := ioutil.ReadDir(tsFileDir)
	tsfileNames := make([]string, 0, len(tsfiles))
	for _, f := range tsfiles {
		if filepath.Ext(f.Name()) == ".ts" {
			tsfileNames = append(tsfileNames, filepath.Join(tsFileDir, f.Name()))
		}
	}
	videoMerge(tsfileNames, outFile)

	return
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



