package main

import (
	"fmt"
	"net/http"
	"sync"
)

type ServerManager struct {
	servers []*http.Server
	waitGroup sync.WaitGroup
	mutex sync.Mutex
	errors []error
}

func (s *ServerManager) launchServer(server *http.Server) {
	var e error = nil
	defer func() { 
		s.waitGroup.Done()
		if e != nil {
			s.mutex.Lock()
			defer s.mutex.Unlock()
			s.errors = append(s.errors, e)
		}
	}()

	ok := server.ListenAndServe()
	if ok != nil && ok != http.ErrServerClosed {
		e = fmt.Errorf(`HTTP(-) interface listener unexpectedly shutdown [%s]`, ok.Error())
	}
}

func (s *ServerManager) Run() {
	for _, server := range s.servers {
		s.waitGroup.Add(1)
		go s.launchServer(server)
	}
	s.waitGroup.Wait()

	for _, e := range s.errors {
		fmt.Println(e)
	}
}

func (s *ServerManager) AddAsyncServer(server *http.Server) {
	s.servers = append(s.servers, server)
}