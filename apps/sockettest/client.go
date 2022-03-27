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
		"hello world!",
		"i love you!",
		"fuck you!",
	}

	for i := 0; i < 3; i++ { //3遍
		for _, v := range msg {
			bData, err := proto.Encode([]byte(v))
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

}
