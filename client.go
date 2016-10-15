package main

import (
	"io"
	"log"
	"net"
	"os"
)

//!+
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//done := make(chan bool)
	go func() {
		//input and send
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			log.Fatal(err)
		}

		//done <- true // signal the main goroutine
	}()

	//receive and show
	io.Copy(os.Stdout, conn) // NOTE: ignoring errors
	log.Println("done")
	//<-done // wait for background goroutine to finish
}

//!-
