package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func RetryDial(retrySecond, retryTime int, addr string) (net.Conn, error) {
	var err error
	var conn net.Conn
	for i := 0; i < retryTime; i++ {
		time.Sleep(time.Second * time.Duration(retrySecond))
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			log.Printf("retry...")
			continue
		}
		return conn, nil
	}
	return nil, err
}

func Monitoring(conn net.Conn) {
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

func Retry(addr string, retrySecond, retryTime int) {
	for {
		conn, err := RetryDial(retrySecond, retryTime, addr)
		if err != nil {
			log.Printf("Fatal: %v", err)
			break
		}

		log.Printf("connected")

		Monitoring(conn)
		conn.Close()
	}
}

func main() {
	var (
		idleSecond  int
		retrySecond int
		retryTime   int
	)
	flag.IntVar(&idleSecond, "i", 5, "idle second")
	flag.IntVar(&idleSecond, "idle", 5, "idle second")
	flag.IntVar(&retrySecond, "r", 1, "retry second")
	flag.IntVar(&retrySecond, "retry", 1, "retry second")
	flag.IntVar(&retryTime, "t", 10, "retry time")
	flag.IntVar(&retryTime, "time", 10, "retry time")
	flag.Parse()

	var addr string = flag.Arg(0)
	for {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Printf("Fatal: %v", err)
			time.Sleep(time.Second * time.Duration(idleSecond))
			continue
		}

		log.Printf("connected")

		fmt.Printf("open\r\n")

		Monitoring(conn)
		conn.Close()
		Retry(addr, retrySecond, retryTime)

		fmt.Printf("close\r\n")
		time.Sleep(time.Second * time.Duration(idleSecond))
	}
}
