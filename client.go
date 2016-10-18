package main

import (
	"./protocal"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	//"time"
)

var quitSemaphore chan bool

type MyData struct {
	Client string `json:"client"`
	Input  string `json:"input"`
}

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

		var aa MyData
		aa.Input = string(msg)
		aa.Client = "longmin"
		result, err := json.Marshal(aa)
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
		fmt.Println(msg)

		j2 := make(map[string]interface{})
		err = json.Unmarshal([]byte(msg), &j2)

		if err != nil {
			fmt.Println(err)
			quitSemaphore <- true
			break
		}

		fmt.Println(j2["input"])

	}
}
