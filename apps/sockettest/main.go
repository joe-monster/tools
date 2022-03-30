package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"tools/apps/sockettest/proto"
)

//server
func main() {
	tcpaddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8282")
	listener, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	reader := bufio.NewReaderSize(conn, 4096) //弄一个4M的字节缓冲池
	for {

		peek, err := reader.Peek(4) //前4个字节是单个协议包大小
		if err != nil {
			if err == io.EOF { //客户端关闭连接
				break
			}
			log.Println(err)
			continue
		}

		buffer := bytes.NewBuffer(peek)
		var length int32
		err = binary.Read(buffer, binary.BigEndian, &length)
		if err != nil { //字节转数字出错
			log.Println(err)
			break
		}

		if length > 4096 { //一个tcp包如果超过4096字节，则直接丢弃
			reader.Reset(conn)
			continue
		}

		for { //数据不够一个自定义协议包大小，一直等待直到收到完整的一个包
			if int32(reader.Buffered()) >= length {
				break
			}
			log.Println(reader.Buffered())
		}

		data := make([]byte, length)
		_, err = reader.Read(data)
		if err != nil { //读取一个包长度的字节流出错
			log.Println(err)
			continue
		}

		//打印一下完整包的二进制数据
		//log.Println("received msg", data)

		//解码
		pro, err := proto.Decode(data)
		if err != nil { //解码出错
			log.Printf("%+v \n", err)
			continue
		}

		log.Printf("\n\tlength: %d\n\theader length: %d\n\tversion: %d\n\tact: %d\n\tseq: %d\n\tbody: %s\n\n", pro.Length, pro.HeaderLength, pro.Version, pro.Act, pro.Seq, string(pro.Body))
	}
}
