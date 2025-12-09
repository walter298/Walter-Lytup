package main

import  (
	"time"
)

type DaemonState int

const (
	STARTING DaemonState = 0
	RUNNING DaemonState  = 1
	DONE                 = 2
)

type Daemon struct  {
	Name   string
	Program  func  (<-  chan map[string]string, chan <- map[string]string) ([]error)
	State  DaemonState  
	StartupGrace     time.Duration
	ShutdownGrace    time.Duration

	// internal use: don't set properties below
	clap   chan map[string]string
	flap   chan map[string]string
}