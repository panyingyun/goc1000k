package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

var quitSemaphore chan bool

func main() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "192.168.0.80:9999")
	for i := 0; i < 15000; i++ {
		conn, _ := net.DialTCP("tcp", nil, tcpAddr)
		go onMessageRecived(conn)
	}
	<-quitSemaphore
}

func onMessageRecived(conn *net.TCPConn) {
	if conn == nil {
		fmt.Println("connected error!!!")
		return
	}
	fmt.Println("connected!")
	b := []byte("time\n")
	conn.Write(b)
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		//fmt.Println(msg)
		if err != nil {
			conn.Close()
			quitSemaphore <- true
			break
		}
		time.Sleep(60 * time.Second)
		b := []byte(msg)
		conn.Write(b)
	}
}
