package main

import (
	"log"
	"net"
	"tools/apps/sockettest/proto"
)

//client
func main() {
	tcpaddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8282")

	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	msg := []string{
		"abc",                                  //用于构造小于一个tcp包
		"1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ", //用于构造大于1500字节的包
	}

	for _, v := range msg {
		vv := []byte{}
		for j := 0; j < 200; j++ {
			vv = append(vv, []byte(v)...)
		}

		bData, err := proto.Encode(vv)
		if err != nil {
			log.Println(err)
			continue
		}

		if _, err := conn.Write(bData); err != nil {
			log.Println(err)
			continue
		}
	}

}
