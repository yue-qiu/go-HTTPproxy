package Adapter

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
)

type Request struct {
	Info []byte
	Url *url.URL
	Method string
	Addr string
}


func NewRequest(conn net.Conn) (*Request, error) {
	var req = Request{}
	info := [4096]byte{}
	rawUrl := ""
	n, err := conn.Read(info[:])
	if err != nil {
		log.Printf("read failed: : %v\n", err)
		return nil, errors.New("read failed")
	}

	req.Info = info[:n]
	_, err = fmt.Sscanf(string(req.Info[:bytes.IndexByte(req.Info, '\r')]), "%s%s", &req.Method, &rawUrl)
	req.Url, err = url.Parse(rawUrl)
	if err != nil {
		log.Printf("parse failed: %v\n", err)
		return nil, errors.New("parse failed")
	}

	fmt.Println("rawUrl:" + rawUrl + " URI:" + req.Url.String() + " Host:" + req.Url.Host + " Scheme:" + req.Url.Scheme + " Opaque:" + req.Url.Opaque)
	if req.Url.Opaque == "443" {
		req.Addr = req.Url.String()
	} else {
		if req.Url.Port() == "" {
			// default port is 80
			req.Addr = req.Url.Host + ":80"
		} else {
			req.Addr = req.Url.Host
		}
	}

	return &req, nil
}
