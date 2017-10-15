package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/urfave/cli"
)

var ConnectCounter int64 = 0
var ReceiveCounter int64 = 0
var SendCounter int64 = 0

func run(c *cli.Context) error {
	port := c.Int("port")
	n := c.Int("number")
	for i := 0; i < n; i++ {
		go startServer(port + i)
		time.Sleep(10 * time.Millisecond)
	}
	//quit when receive end signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	fmt.Printf("signal received signal %v\n", <-sigChan)
	fmt.Println("shutting down server")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "c1000k-server"
	app.Usage = "c1000k-server"
	app.Copyright = "panyingyun@gmail.com"
	app.Version = "0.1.1"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "port,p",
			Usage:  "Set Server start port here",
			Value:  10000,
			EnvVar: "PORT",
		},
		cli.IntFlag{
			Name:   "number,n",
			Usage:  "Set Number of Server here",
			Value:  1,
			EnvVar: "NUMBER",
		},
	}
	app.Run(os.Args)
}

//启动TCP服务
func startServer(port int) {
	addr := fmt.Sprintf("0.0.0.0:%v", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
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

	fmt.Printf("Start Server Listen On Address:[%v]\n", addr)
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}
		atomic.AddInt64(&ConnectCounter, 1)
		printCounter()
		go handleMessage(tcpConn)
	}
}

func handleMessage(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		atomic.AddInt64(&ConnectCounter, -1)
		printCounterAndIP(ipStr)
		conn.Close()
	}()
	reader := bufio.NewReader(conn)

	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		atomic.AddInt64(&ReceiveCounter, 1)
		//fmt.Println(message)
		msg := time.Now().String() + "\n"
		b := []byte(msg)
		conn.Write(b)
		atomic.AddInt64(&SendCounter, 1)
		printCounter()
	}
}

func handleServerError(err error) {
	fmt.Println("Server Error:", err)
	printCounter()
}

func printCounter() {
	RConnectCounter := atomic.LoadInt64(&ConnectCounter)
	RReceiveCounter := atomic.LoadInt64(&ReceiveCounter)
	RSendCounter := atomic.LoadInt64(&SendCounter)
	RGoroutine := runtime.NumGoroutine()
	fmt.Printf("Connect [%v] [%v] [%v] [%v]\n", RConnectCounter, RReceiveCounter, RSendCounter, RGoroutine)
}

func printCounterAndIP(ip string) {
	RConnectCounter := atomic.LoadInt64(&ConnectCounter)
	RReceiveCounter := atomic.LoadInt64(&ReceiveCounter)
	RSendCounter := atomic.LoadInt64(&SendCounter)
	RGoroutine := runtime.NumGoroutine()
	fmt.Printf("Connect [%v] [%v] [%v] [%v] [%v]\n", RConnectCounter, RReceiveCounter, RSendCounter, RGoroutine, ip)
}
