package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
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
			peek,err:=reader.Peek(4)//读取长度

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
				continue
			}

			data := make([]byte, length + 4)
			_, err = reader.Read(data)
			if err != nil {
				continue
			}
			fmt.Println("received msg:", string(data[4:]))
			time.Sleep(time.Second)
		}
	}

}


