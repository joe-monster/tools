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

	//将数据流做成结构体数据，解决粘包问题
	reader := bufio.NewReader(conn)
	for {

		peek, err := reader.Peek(4)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				continue
			} else {
				break
			}
		}

		buffer := bytes.NewBuffer(peek)
		var length int32
		err = binary.Read(buffer, binary.BigEndian, &length)
		if err != nil { //前4个字节是整个包大小，转数字出错
			log.Println(err)
			break
		}

		if int32(reader.Buffered()) < length { //数据不够一个包大小
			continue
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
