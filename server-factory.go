package main

import (
	"bytes"
	"log"
	"net/http"
)

func MakeServer(addr string) *http.Server {
	server := &http.Server{Addr: addr, Handler: &PanicManager{}}

	server.MaxHeaderBytes = DHI0_MaxHeaderSize
	server.ReadTimeout = DHI0_ReadTimeout
	server.ReadHeaderTimeout = DHI0_ReadTimeout
	server.WriteTimeout = DHI0_WrttTimeout
	server.IdleTimeout = DHI0_IdleTimeout
	error_buff := bytes.NewBuffer([]byte{})
	server.ErrorLog = log.New(error_buff, "", log.Lshortfile)

	return server
}