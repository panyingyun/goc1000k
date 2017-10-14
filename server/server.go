package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

var ConnectCounter int = 0
var ReceiveCounter int = 0
var SendCounter int = 0
var sn sync.RWMutex

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:9999")
	if err != nil {
		handleServerError(err)
		return
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		handleServerError(err)
		return
	}

	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}
		sn.Lock()
		ConnectCounter++
		fmt.Printf("Connect [%v] [%v] [%v] [%v]\n", ConnectCounter, ReceiveCounter, SendCounter, tcpConn.RemoteAddr().String())
		sn.Unlock()

		go handleMessage(tcpConn)
	}

}

func handleMessage(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		sn.Lock()
		ConnectCounter--
		fmt.Printf("Connect [%v] [%v] [%v] [%v]\n", ConnectCounter, ReceiveCounter, SendCounter, ipStr)
		conn.Close()
		sn.Unlock()
	}()
	reader := bufio.NewReader(conn)

	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		sn.Lock()
		ReceiveCounter++
		sn.Unlock()
		//fmt.Println(message)
		msg := time.Now().String() + "\n"
		b := []byte(msg)
		conn.Write(b)
		sn.Lock()
		SendCounter++
		sn.Unlock()
		fmt.Printf("Connect [%v] [%v] [%v] [%v]\n", ConnectCounter, ReceiveCounter, SendCounter, ipStr)
	}
}

func handleServerError(err error) {
	fmt.Println("Server Error:", err)
	fmt.Printf("Connect [%v] [%v] [%v]\n", ConnectCounter, ReceiveCounter, SendCounter)
}
