package main

import (
	"fmt"
	//"net"
	//"time"
	"hash/crc32"
)

func main() {
	// conn, _ := net.Dial("tcp", "127.0.0.1:8000")
	// buf := make([]byte, 32*1024)

	// for {
	// 	fmt.Println("loop start")
	// 	nr, er := conn.Read(buf)
	// 	if nr > 0 {
	// 		fmt.Println(buf[:nr])
	// 	}
	// 	if er != nil {
	// 		fmt.Println(er)
	// 		break
	// 	}
	// }

	fmt.Println(crc32.ChecksumIEEE([]byte("aaa")))

	fmt.Println("over")

	/*msgs := make([]bool, 0, 10)
	msgs = append(msgs, true, false)
	fmt.Println(msgs)
	msgs = msgs[:0]
	fmt.Println(msgs)

	in := make(chan bool)
	flag := false

	go func() {
		time.Sleep(5 * time.Second)
		close(in)
		fmt.Println("close in channel")
	}()

	for {
		fmt.Println("wait in channel")
		select {
		case <-in:
			fmt.Println("come from in channel")
			flag = true
		}
		if flag {
			break
		}
	}

	fmt.Println("over 88")*/
}
