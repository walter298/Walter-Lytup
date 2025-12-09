package main

import (
	"fmt"
	"net/http"
	"strings"
)

type RequestMap = map[string]string
type ServerSlice = []*http.Server
type ListenHandler = func(*http.Server) error

type TLSServer struct {
	server *http.Server
	crt string
	key string
}
type TLSServerSlice = []*TLSServer

type ServerManager struct {
	completion ServerCompletion
	servers ServerSlice
	tlsServers TLSServerSlice
}

func (s *ServerManager) isDone() bool {
	return s.completion.isDone()
}

func (s *ServerManager) launchServerImpl(server *http.Server, listen_func ListenHandler) {
	var e error = nil
	defer func() { s.completion.finishServer(e) }()

	Output_Logg("OUT", "DHI1", fmt.Sprintf(`HTTP(-): %s`, strings.Split((*server).Addr, ":")[1]))

	ok := listen_func(server)
	if ok != nil && ok != http.ErrServerClosed {
		e = fmt.Errorf(`HTTP(-) interface listener unexpectedly shutdown [%s]`, ok.Error())
	}
}

func (s *ServerManager) launchServerTLS(tlsServer *TLSServer) {
	listen := func(s *http.Server) error {
		return s.ListenAndServeTLS(tlsServer.crt, tlsServer.key)
	}
	go s.launchServerImpl(tlsServer.server, listen)
}

func (s *ServerManager) launchServer(server *http.Server) {
	listen := func(s *http.Server) error {
		return s.ListenAndServe()
	}
	go s.launchServerImpl(server, listen)
}

func (s *ServerManager) Run() {
	for _, server := range s.servers {
		s.launchServer(server)
	}
	for _, tlsServer := range s.tlsServers {
		s.launchServerTLS(tlsServer)
	}
}

func (s *ServerManager) AddAsyncServer(server *http.Server) {
	s.completion.register()
	s.servers = append(s.servers, server)
}

func (s *ServerManager) AddAsyncServerTLS(crt string, key string, server *http.Server) {
	s.completion.register()
	tlsServer := &TLSServer{ 
		crt:crt,
		key:key,
		server:server,
	}
	s.tlsServers = append(s.tlsServers, tlsServer)
}
