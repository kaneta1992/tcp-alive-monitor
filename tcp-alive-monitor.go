package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func DialTCP(addr string) (*net.TCPConn, error) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	return conn, err
}

func RetryDial(addr string) (*net.TCPConn, error) {
	var err error
	var conn *net.TCPConn
	for i := 0; i < retryTime; i++ {
		time.Sleep(time.Second * time.Duration(retrySecond))
		conn, err = DialTCP(addr)
		if err != nil {
			log.Printf("retry...")
			continue
		}
		return conn, nil
	}
	return nil, err
}

func Monitoring(conn *net.TCPConn) {
	r := bufio.NewReader(conn)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Printf("end %v", err)
			break
		}
		log.Printf("response: %s", line)
	}
}

func Retry(addr string) {
	for {
		conn, err := RetryDial(addr)
		if err != nil {
			log.Printf("Fatal: %v", err)
			break
		}
		SetKeepAlive(conn, time.Second*time.Duration(keepAlive))

		log.Printf("connected")

		Monitoring(conn)
		conn.Close()
	}
}

func SetKeepAlive(conn *net.TCPConn, t time.Duration) {
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(t)
}

var (
	idleSecond  int
	retrySecond int
	retryTime   int
	keepAlive   int
)

func main() {
	flag.IntVar(&idleSecond, "i", 5, "idle second")
	flag.IntVar(&idleSecond, "idle", 5, "idle second")
	flag.IntVar(&retrySecond, "r", 1, "retry second")
	flag.IntVar(&retrySecond, "retry", 1, "retry second")
	flag.IntVar(&retryTime, "t", 10, "retry time")
	flag.IntVar(&retryTime, "time", 10, "retry time")
	flag.IntVar(&keepAlive, "k", 5, "KeepAlive interval Second")
	flag.IntVar(&keepAlive, "keep", 5, "KeepAlive interval Second")
	flag.Parse()

	var addr string = flag.Arg(0)
	for {
		conn, err := DialTCP(addr)
		if err != nil {
			log.Printf("Fatal: %v", err)
			time.Sleep(time.Second * time.Duration(idleSecond))
			continue
		}
		SetKeepAlive(conn, time.Second*time.Duration(keepAlive))

		log.Printf("connected")

		fmt.Printf("open\n")

		Monitoring(conn)
		conn.Close()
		Retry(addr)

		fmt.Printf("close\n")
		time.Sleep(time.Second * time.Duration(idleSecond))
	}
}
