package main

import "time"

const (
	Address_1                string        = ":80"
	Address_2                string        = ":443"
	MaxHeaderSize            int           = 1 * 1024 * 1024
	ReadTimeout              time.Duration = time.Minute * 5
	WriteTimeout             time.Duration = time.Minute * 5
	IdleTimeout              time.Duration = time.Minute * 5
	TimeZoneSecondOffset     int           = 0
)