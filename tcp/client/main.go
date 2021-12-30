package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)


func converToBianry(n int) string {
	result := ""
	for ; n > 0; n /= 2 {
		lsb := n % 2
		result = strconv.Itoa(lsb) + result
	}
	return result
}

func main(){

	con,err:=net.Dial("tcp","127.0.0.1:9000")
	defer con.Close()
	if err != nil {
		fmt.Println("connect error:",err)
		return
	}

	for i:=0;i<10;i++{
		msg:="hello world!hello world!hello world!hello world!hello world!hello world!"
		res,_:=MyEncode(msg)
		con.Write(res)
	}

	fmt.Println("发送完成")

}

//方式二
func MyEncode(message string) ([]byte, error) {
	magicNum := make([]byte, 4)
	binary.BigEndian.PutUint32(magicNum, uint32(len(message))) //把数字转化成网络字节序,然后大端方式

	b:=bytes.NewBuffer(magicNum) //=以一个已存在的字节创建缓冲区
	err:=binary.Write(b,binary.BigEndian,[]byte(message)) //写入主题内容
	if err != nil{
		return nil,err
	}
	return b.Bytes(),nil
}

func Encode(message string) ([]byte, error) {
	var length = int32(len(message)) //计算出长度
	var pkg = new(bytes.Buffer) //创建一个字节缓冲区
	err := binary.Write(pkg, binary.BigEndian, length)//以大端序往直接缓冲区中写入数据
	if err != nil {
		return nil, err
	}
	err = binary.Write(pkg, binary.BigEndian, []byte(message)) //同上
	if err != nil {
		return nil, err
	}

	return pkg.Bytes(), nil //返回缓冲区中的所有数据
}

