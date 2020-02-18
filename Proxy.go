package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

const (
	DefaultPort  = ":7080"
	TCP = "tcp"
	UDP = "udp"
)

func main()  {
	In, err := net.Listen(TCP, DefaultPort)
	if err != nil {
		log.Printf("listenning network: %v\n", err)
		return
	}

	for  {
		conn, err := In.Accept()
		if err != nil {
			log.Printf("distributing connection: %v\n", err)
			return
		}
		go connHandler(conn)
	}
}

func connHandler(conn net.Conn)  {
	if conn == nil {
		return
	}
	defer conn.Close()

	var info [1024]byte
	_, err := conn.Read(info[:])
	if err != nil {
		log.Printf("Reading conn: %v\n", err)
		return
	}

	//fmt.Println(string(info[:n]))
	var method, host, address string
	fmt.Println("【ATTENTION】" + string(info[:bytes.IndexByte(info[:], '\r')]))
	_, err = fmt.Sscanf(string(info[:bytes.IndexByte(info[:], '\r')]), "%s%s", &method, &host)
	hostPortUrl, err := url.Parse(host)
	if err != nil {
		log.Printf("Parsing URI: %v\n", err)
		return
	}

	fmt.Println(host + " URI:" + hostPortUrl.String() + " Host:" + hostPortUrl.Host + " Scheme:" + hostPortUrl.Scheme + " Opaque:" + hostPortUrl.Opaque)
	if hostPortUrl.Opaque == "443" {
		address = hostPortUrl.String()
	} else {
		if strings.Index(hostPortUrl.Host, ":") == -1 {
			address = hostPortUrl.Host + ":80"  // default port is 80
		} else {
			address = hostPortUrl.Host
		}
	}

	var server net.Conn
	server, err = net.Dial("tcp", address)
	if err != nil {
		log.Printf("Dialing: %v\n", err)
		return
	}

	if method == "CONNECT" {
		_, _ = fmt.Fprint(conn, "HTTP /1.1 200 Connection established\r\n\r\n")
	}


	go io.Copy(server, conn)
	io.Copy(conn, server)
}

