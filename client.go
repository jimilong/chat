package main

import (
	"./protocal"
	"bufio"
	"fmt"
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

		//fmt.Println(string(msg))
		if len(msg) > 0 {
			b, err := protocal.Pack(string(msg), "test.service", 111)
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
		fmt.Println(msg)
		if err != nil {
			quitSemaphore <- true
			break
		}
	}
}
