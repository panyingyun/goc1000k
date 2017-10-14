package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

var quitSemaphore chan bool

func main() {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	for i := 0; i < 1000; i++ {
		conn, _ := net.DialTCP("tcp", nil, tcpAddr)
		go onMessageRecived(conn)
		time.Sleep(1000 * time.Microsecond)
	}
	<-quitSemaphore
}

func onMessageRecived(conn *net.TCPConn) {
	if conn == nil {
		fmt.Println("connected error!!!")
		return
	}
	b := []byte("time\n")
	conn.Write(b)
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			quitSemaphore <- true
			break
		}
		fmt.Printf("Connect [%v]\n", conn.LocalAddr().String())
		time.Sleep(60 * time.Second)
		b := []byte(msg)
		conn.Write(b)
	}
}
