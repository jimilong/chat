package main

import (
	"bufio"
	"chat/protobuf/example"
	"chat/protocal"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

var ConnMap map[string]*net.TCPConn

func main() {
	var tcpAddr *net.TCPAddr
	ConnMap = make(map[string]*net.TCPConn)
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)

	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		fmt.Println("A client connected : " + tcpConn.RemoteAddr().String())
		ConnMap[tcpConn.RemoteAddr().String()] = tcpConn
		go handleConn(tcpConn)
	}

}

func handleConn(conn *net.TCPConn) {
	timeout := make(chan bool)
	in := make(chan bool)

	defer func() {
		ipStr := conn.RemoteAddr().String()
		fmt.Println("disconnected :" + ipStr)
		conn.Close()
	}()

	go handleTimeout(timeout, in)
	go handleMsg(conn, in)

	<-timeout
}

func handleTimeout(timeout chan bool, in chan bool) {
	tick := time.Tick(30 * time.Second)
	msgs := make([]bool, 0, 10)
	flag := false
	for {
		select {
		case msg := <-in:
			if len(msgs) < 1 {
				msgs = append(msgs, msg)
			}

			fmt.Println(msgs)
		case <-tick:
			if len(msgs) == 0 {
				fmt.Println("time-out")
				timeout <- true
				flag = true
			} else {
				msgs = msgs[:0]
			}
		}
		if flag {
			break
		}
	}
	fmt.Println("handleTimeout over 88")
}

func handleMsg(conn net.Conn, in chan bool) {
	//ipStr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)

	for {
		message, err := protocal.Unpack(reader)
		if err != nil {
			return
		}
		//通知超时处理方法
		in <- true

		//解码
		j2 := &example.MyData{}
		err = proto.Unmarshal([]byte(message), j2)
		if err != nil {
			continue
		}

		fmt.Println(j2.Client + ":" + j2.Input)

		//b, err := protocal.Pack(conn.RemoteAddr().String()+":"+string(message), "aa", 111)
		b, err := protocal.Pack(string(message), "aa", 111)
		if err != nil {
			continue
		}
		//boradcastMessage
		for _, conn := range ConnMap {
			conn.Write(b)
		}

	}
}
