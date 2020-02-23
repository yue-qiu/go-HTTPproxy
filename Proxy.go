package main

import (
	"fmt"
	"github.com/yue-qiu/go-HttpProxy/Adapter"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const (
	DefaultPort  = ":7080"
	TCP = "tcp"
)

func main()  {
	In, err := net.Listen(TCP, DefaultPort)
	if err != nil {
		log.Printf("listen failed: %v\n", err)
		os.Exit(1)
	}

	for  {
		conn, err := In.Accept()
		if err != nil {
			log.Printf("connect failed: %v\n", err)
			continue
		}
		go connHandler(conn)
	}
}

func connHandler(conn net.Conn)  {
	if conn == nil {
		return
	}
	defer conn.Close()

	var req, err = Adapter.NewRequest(conn)
	if err != nil {
		return
	}

	var server net.Conn
	var deadline = time.Now().Add(30 * time.Second)

	for {
		for tries := 0; time.Now().Before(deadline); tries++ {
			server, err = net.Dial(TCP, req.Addr)
			if err == nil {
				break
			} else {
				log.Printf("dial failed: %v, retrying...", err)
				time.Sleep(time.Second << uint(tries))
			}
		}
		if time.Now().Before(deadline) {
			_ = server.SetDeadline(time.Now().Add(3 * time.Second))
			if req.Method == "CONNECT" {
				_, _ = fmt.Fprint(conn, "HTTP/1.1 200 Connection established\r\n\r\n")
				go io.Copy(server, conn)
			} else {
				go func() {
					_, err := server.Write(req.Info)
					if err != nil {
						log.Printf("write failed: %v\n", err)
					}
				}()
			}
		} else {
			log.Printf("server %v not responding", req.Addr)
			return
		}
		break
	}
	io.Copy(conn, server)
}
