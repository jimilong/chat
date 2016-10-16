package main

import (
	"./protocal"
	"bufio"
	"fmt"
	"net"
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
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println("disconnected :" + ipStr)
		conn.Close()
	}()
	reader := bufio.NewReader(conn)

	for {
		message, err := protocal.Unpack(reader)
		if err != nil {
			return
		}
		fmt.Println(conn.RemoteAddr().String() + ":" + string(message))

		b, err := protocal.Pack(conn.RemoteAddr().String() + ":" + string(message))
		if err != nil {
			continue
		}
		//boradcastMessage
		for _, conn := range ConnMap {
			conn.Write(b)
		}

	}
}
