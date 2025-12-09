package main

/*
	Name: DaemonCore-go
	Version: 0.0.3
*/

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"slices"
	"syscall"
	"time"
)

func shutdownDaemon(daemon *Daemon) {
	requestMap := map[string]string{}
	requestMap["Command"] = "shutdown"
	daemon.clap <- requestMap

	daemonShutdownMessage := fmt.Sprintf(
		`PROJECT: Daemon %s: Shutdown signalledd`, daemon.Name,
	)
	Output_Logg("OUT", "Main", daemonShutdownMessage)
	xd20 := make(chan bool, 1)
	if daemon.ShutdownGrace != 0 { //no idea what this is doing
		go func() {
			time.Sleep(daemon.ShutdownGrace)
			xd20 <- true
		}()
	}
	select {
	case xf05 := <-daemon.flap:
		{
			if xf05["ExctnOtcmCode"] != "200" {
				xg05 := fmt.Sprintf(
					`PROJECT: Daemon %s: Encountered error [%s]`, daemon.Name, xf05["ExctnOtcmNote"],
				)
				Output_Logg("ERR", "Main", xg05)
				return
			}
		}
	case _ = <-xd20:
		{
		}
	}
	xd25 := fmt.Sprintf(
		`PROJECT: Daemon %s: Shutdown successful`, daemon.Name,
	)
	Output_Logg("OUT", "Main", xd25)
}

func concatErrorMessages(errors []error) string {
	ret := ""
	for _, error := range errors {
		ret += error.Error() + " "
	}
	return ret
}

func main() {
	/***1***/
	Output_Logg("OUT", "Main", "PROJECT: Starting up")

	DHI0_SPRegister = append(DHI0_SPRegister, &ResponseHandler{
		Code: "sr05",
		Program: walterHandler,
	})

	var DaemonRegister []*Daemon = []*Daemon{
		{Name: "DHI0", Program: runPanicServers, State: STARTING},
	}
	if len(DaemonRegister) == 0 {
		Output_Logg("OUT", "Main", "PROJECT: No Daemon(s) to run. Shutting down now")
		return
	}
	var SupportedShutdownSignal []syscall.Signal = []syscall.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	}

	/***2***/
	defer func() {
		reversedRegister := DaemonRegister
		slices.Reverse(reversedRegister)
		for _, daemon := range reversedRegister {
			if daemon.State != RUNNING {
				continue
			}
			shutdownDaemon(daemon)
		}
	}()

	/***3***/
	xb05 := make(chan bool, 1)
	for _, daemon := range DaemonRegister {
		/***1***/
		if daemon.Program == nil {
			xd05 := fmt.Sprintf(
				`PROJECT: Daemon %s: Skipping (Daemon has no program to run)`,
				daemon.Name,
			)
			Output_Logg("OUT", "Main", xd05)
			continue
		}
		/***2***/
		daemon.clap = make(chan map[string]string, 1)
		daemon.flap = make(chan map[string]string, 1)
		xc25 := fmt.Sprintf(
			`PROJECT: Daemon %s: Starting up... Please wait`, daemon.Name,
		)
		Output_Logg("OUT", "Main", xc25)
		/***3***/
		go func() {
			defer func() {
				xe05 := recover()
				if xe05 == nil {
					return
				}
				daemon.State = 2
				xe10 := fmt.Sprintf(
					`Paniced [%v : %s]`, xe05, debug.Stack(),
				)
				xe15 := map[string]string{}
				xe15["ExctnOtcmCode"] = "500"
				xe15["ExctnOtcmNote"] = xe10
				daemon.flap <- xe15
				xb05 <- true
			}()
			daemon.State = 1
			xd05 := daemon.Program(daemon.clap, daemon.flap)
			daemon.State = 2
			xd10 := map[string]string{}
			xd10["ExctnOtcmCode"] = "200"
			if xd05 != nil {
				xd10["ExctnOtcmCode"] = "500"
				xd10["ExctnOtcmNote"] = concatErrorMessages(xd05)
			}
			daemon.flap <- xd10
			xb05 <- true
		}()
		/***4***/
		xc30 := make(chan bool, 1)
		if daemon.StartupGrace != 0 {
			go func() {
				time.Sleep(daemon.StartupGrace)
				xc30 <- true
			}()
		}
		select {
		case xe05 := <-daemon.flap:
			{
				if xe05["StartupCode"] != "200" {
					xf05 := fmt.Sprintf(
						`PROJECT: Daemon %s: Startup failed [%s]`,
						daemon.Name, xe05["StartupNote"],
					)
					Output_Logg("ERR", "Main", xf05)
					return
				}
				break
			}
		case _ = <-xc30:
			{
				xe10 := fmt.Sprintf(
					`PROJECT: Daemon %s: Startup failed [%s]`,
					daemon.Name, "Startup grace period expired",
				)
				Output_Logg("ERR", "Main", xe10)
				return
			}
		}
		/***5***/
		xc35 := fmt.Sprintf(`PROJECT: Daemon %s: Up and running`, daemon.Name)
		Output_Logg("OUT", "Main", xc35)
	}
	/***4***/
	xb10 := make(chan os.Signal, 1)
	for _, xc10 := range SupportedShutdownSignal {
		signal.Notify(xb10, xc10)
	}
	for {
		select {
		case _ = <-xb05:
			{
				for _, xf10 := range DaemonRegister {
					select {
					case xh10 := <-xf10.flap:
						{
							if xh10["ExctnOtcmCode"] != "200" {
								xi05 := fmt.Sprintf(
									`PROJECT: Daemon %s: Encountered error [%s]`, xf10.Name, xh10["ExctnOtcmNote"],
								)
								Output_Logg(
									"ERR", "Main", xi05,
								)
								return
							}
							xh15 := fmt.Sprintf(
								`PROJECT: Daemon %s: Finished`, xf10.Name,
							)
							Output_Logg("OUT", "Main", xh15)
						}
					default:
						{
						}
					}
				}
			}
		case _ = <-xb10:
			{
				fmt.Println("")
				return
			}
		}
	}
}
