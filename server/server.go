package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

var ConnectCount int = 0
var sn sync.RWMutex

func main() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "0.0.0.0:9999")
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}
		sn.Lock()
		ConnectCount++
		fmt.Println("ConnectCount = ", ConnectCount)
		fmt.Println("connected : " + tcpConn.RemoteAddr().String())
		sn.Unlock()

		go handleMessage(tcpConn)
	}

}

func handleMessage(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		sn.Lock()
		ConnectCount--
		fmt.Println("ConnectCount = ", ConnectCount)
		fmt.Println("disconnected :" + ipStr)
		conn.Close()
		sn.Unlock()
	}()
	reader := bufio.NewReader(conn)

	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		//fmt.Println(message)
		msg := time.Now().String() + "\n"
		b := []byte(msg)
		conn.Write(b)
	}
}
