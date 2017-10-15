package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/urfave/cli"
)

func run(c *cli.Context) error {
	ipaddr := c.String("ipaddr")
	port := c.Int("port")
	servernumber := c.Int("number")
	connectnum := c.Int("connetcnt")

	for i := 0; i < servernumber; i++ {
		server := fmt.Sprintf("%v:%v", ipaddr, port+i)
		fmt.Printf("Server Addr: [%v]\n", server)
		go connectToServer(server, connectnum)
		time.Sleep(time.Second)
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
	app.Name = "c1000k-client"
	app.Usage = "c1000k-client"
	app.Copyright = "panyingyun@gmail.com"
	app.Version = "0.1.0"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "ipaddr,ip",
			Usage:  "Set Server start ip address here",
			Value:  "127.0.0.1",
			EnvVar: "IPADDR",
		},
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
		cli.IntFlag{
			Name:   "connetcnt,cn",
			Usage:  "Set Connet Count of Server here",
			Value:  1,
			EnvVar: "CONNETCNT",
		},
	}
	app.Run(os.Args)
}

func connectToServer(server string, connectnum int) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", server)
	for i := 0; i < connectnum; i++ {
		conn, _ := net.DialTCP("tcp", nil, tcpAddr)
		go onMessageRecived(conn)
		time.Sleep(100 * time.Millisecond)
	}
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
			break
		}
		fmt.Printf("Connect [%v] [%v]\n", conn.LocalAddr().String(), runtime.NumGoroutine())
		time.Sleep(6000 * time.Second)
		b := []byte(msg)
		conn.Write(b)
	}
}
