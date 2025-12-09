package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func  Output_Logg(Type, Source, Output string) {
	Type  = strings.ToLower (Type)
	xb05 := fmt.Sprintf (
		`[%s//%s] %s`, time.Now ().In (time.FixedZone("TTT", TimeZoneSecondOffset)).Format ("2006-01-02 15:04:05.000 -07:00"), Source, Output,
	)
	if Type == "out" {
		os.Stdout.Write ([ ]byte (xb05 + "\n") )
	}  else {
		os.Stderr.Write ([ ]byte (xb05 + "\n") )
	}
}