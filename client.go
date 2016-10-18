package main

import (
	"bufio"
	"chat/protobuf/example"
	"chat/protocal"
	//"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	//"time"
)

var quitSemaphore chan bool

func main() {
	var tcpAddr *net.TCPAddr
	quitSemaphore = make(chan bool)
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	fmt.Println("connected!")

	go onMsgRecv(conn)
	go sendMsg(conn)

	<-quitSemaphore
}

func sendMsg(conn *net.TCPConn) {
	msg := make([]byte, 0, 8*1024)
	reader := bufio.NewReader(os.Stdin)

	for {
		msg, _, _ = reader.ReadLine()

		aa := &example.MyData{
			Client: "longmin",
			Input:  string(msg),
		}
		//编码
		result, err := proto.Marshal(aa)
		if err != nil {
			fmt.Println(err)
		}

		//fmt.Println(string(msg))
		if len(msg) > 0 {
			b, err := protocal.Pack(string(result), "test.service", 111)
			if err != nil {
				quitSemaphore <- true
				break
			}
			_, err = conn.Write(b)
			if err != nil {
				quitSemaphore <- true
				break
			}
			msg = msg[:0]
		}
	}
}

func onMsgRecv(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := protocal.Unpack(reader)
		if err != nil {
			quitSemaphore <- true
			break
		}
		//fmt.Println(msg)
		//解码
		j2 := &example.MyData{}
		err = proto.Unmarshal([]byte(msg), j2)

		if err != nil {
			fmt.Println(err)
			quitSemaphore <- true
			break
		}

		fmt.Println(j2.Client+":", j2.Input)

	}
}
