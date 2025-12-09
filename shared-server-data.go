package main

import "sync"

type ServerCompletion struct {
	serverCount int
	mutex       sync.Mutex
	errors      []error
	serversDone int
}

func (s *ServerCompletion) register() {
	s.serverCount++
}

func (s *ServerCompletion) finishServer(e error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if e != nil {
		s.errors = append(s.errors, e)
	}
	s.serversDone++
}

func (s *ServerCompletion) isDone() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.serversDone == s.serverCount
}

func (s *ServerCompletion) getErrors() []error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.errors
}
