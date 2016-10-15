package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type client chan string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	message  = make(chan string)
)

//msg transmit center
func broadcaster() {
	clients := make(map[client]bool)

	for {
		select {
		case msg := <-message:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func clientWriter(conn net.Conn, ch chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func clientInput(conn net.Conn, out chan bool, in chan bool) {
	input := bufio.NewScanner(conn)
	who := conn.RemoteAddr().String()
	for input.Scan() {
		message <- who + ": " + input.Text()
		in <- true
	}

	out <- true
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

func handleConn(conn net.Conn) {
	ch := make(chan string)
	in := make(chan bool, 1)
	out := make(chan bool)
	timeout := make(chan bool)

	//send msg to client
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()

	ch <- "You are " + who
	message <- who + "has arrived"
	entering <- ch

	//超时处理
	go handleTimeout(timeout, in)

	//receive msg from client
	go clientInput(conn, out, in)

	for {
		select {
		case <-out:
			leaving <- ch
			message <- who + ": have left"
			conn.Close()
		case <-timeout:
			leaving <- ch
			message <- who + ": lose contact"
			conn.Close()
		}
		break
	}

	fmt.Println(who + " client over")
}

func main() {
	listner, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go handleConn(conn)
	}
}
