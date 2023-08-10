package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		fmt.Println("Waiting for Connection...")
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Connection from ", conn.RemoteAddr().String())
		handle(conn)

	}
}

// For single communication
func handle(c net.Conn) {
	defer c.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go reader(c, &wg)
	go sender(c, &wg)

	wg.Wait()
	fmt.Println("Connection Closed!")
}

// Send input + delimiter
func sender(c net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		fmt.Print("Server #> ")
		sc := bufio.NewScanner(os.Stdin)
		if sc.Scan(); sc.Text() != "exit" {
			c.Write([]byte(sc.Text() + "Message End-dnE egasseM"))
		} else {
			//fmt.Println("Sender Exit!")
			c.Write([]byte("exit" + "Message End-dnE egasseM"))
			break
		}
	}
}

// Read response until delimiter
func reader(c net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		res := make([]byte, 1024)
		buf := bytes.Buffer{}
		for !bytes.Contains(buf.Bytes(), []byte("Message End-dnE egasseM")) {
			_, err := c.Read(res)
			if err != nil {
				log.Fatal(err)
			}
			buf.Write(res)
		}
		i := bytes.Index(buf.Bytes(), []byte("Message End-dnE egasseM"))
		r := buf.Bytes()[:i]
		if string(r) == "exit" {
			fmt.Println("Clossing Connection!")
			c.Write([]byte("exit" + "Message End-dnE egasseM"))
			break
		} else {
			fmt.Printf("\nClient >> %s\nServer #> ", string(r))
		}
	}
}
