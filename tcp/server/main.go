package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func main(){

	list,err:=net.Listen("tcp",":9000")
	defer list.Close()
	if err != nil {
		fmt.Println("listen error :",err)
		return
	}

	for {
		con,err:=list.Accept()
		defer con.Close()
		if err != nil {
			fmt.Println("accept error :",err)
			return
		}
		reader:=bufio.NewReader(con) //创建一个具有默认大小缓冲、从conn中读取的*reader。
		for{
			peek,err:=reader.Peek(4)//读取长度,且指针不移动,主要是为了如果本次没有接收完成,下一次用到重新读取长度

			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}else {
					break
				}
			}
			buffer := bytes.NewBuffer(peek)
			var length int32
			err = binary.Read(buffer,binary.BigEndian,&length)
			if err != nil {
				fmt.Println(err)
			}
			//表示还没有接收完成,继续接收
			if int32(reader.Buffered()) < length + 4 {
				fmt.Println("服务端获取数据进行了等待")
				continue //这里continue是因为reader会不断的从socket缓冲区中接收数据,这一次没有读取完整,下一次说不定数据就在里面了
			}

			data := make([]byte, length + 4)
			_, err = reader.Read(data)
			if err != nil {
				continue
			}
			fmt.Println("received msg:", string(data[4:]))
		}
	}

}


