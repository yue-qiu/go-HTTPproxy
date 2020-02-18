package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
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

	var info [4096]byte
	n, err := conn.Read(info[:])
	if err != nil {
		log.Printf("read failed: : %v\n", err)
		return
	}

	var method, rawURL, address string
	_, err = fmt.Sscanf(string(info[:bytes.IndexByte(info[:], '\r')]), "%s%s", &method, &rawURL)
	URI, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("parse failed: %v\n", err)
		return
	}

	fmt.Println("rawURL:" + rawURL + " URI:" + URI.String() + " Host:" + URI.Host + " Scheme:" + URI.Scheme + " Opaque:" + URI.Opaque)
	if URI.Opaque == "443" {
		address = URI.String()
	} else {
		if strings.Index(URI.Host, ":") == -1 {
			address = URI.Host + ":80"  // default port is 80
		} else {
			address = URI.Host
		}
	}

	var server net.Conn
	server, err = net.Dial(TCP, address)
	if err != nil {
		log.Printf("dial failed: : %v\n", err)
		_, _ = fmt.Fprint(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
		return
	}

	if method == "CONNECT" {
		_, _ = fmt.Fprint(conn, "HTTP/1.1 200 Connection established\r\n\r\n")
		go io.Copy(server, conn)
	} else {
		go func() {
			_, err := server.Write(info[:n])
			if err != nil {
				log.Printf("write failed: %v\n", err)
			}
		}()
	}

	io.Copy(conn, server)
}
