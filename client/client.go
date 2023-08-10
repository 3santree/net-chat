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
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("Connection Succeed!")
	// Close connection when server sends [exit]
	var wg sync.WaitGroup
	wg.Add(1)
	go reader(conn, &wg)
	go sender(conn)
	wg.Wait()
	fmt.Println("Connection Closed!")
}

// Send input + delimiter
func sender(c net.Conn) {
	for {
		fmt.Print("Client >> ")
		sc := bufio.NewScanner(os.Stdin)
		if sc.Scan() {
			c.Write([]byte(sc.Text() + "Message End-dnE egasseM"))
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
			fmt.Printf("\nServer #> %s\nClient >> ", string(r))
		}
	}
}
