package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"runtime/debug"
	"slices"
	"time"
)

func runPanicServers(Clap <-chan RequestMap, Flap chan<- RequestMap) (Errors []error) {
	fmt.Println("Running panic servers...")

	/***1***/
	serverManager := ServerManager{}
	requestMap := map[string]string{}
	servers := []*http.Server{}

	if DHI0_Addr1 != "" {
		panicServer := MakeServer(DHI0_Addr1)
		servers = append(servers, panicServer)
		serverManager.AddAsyncServer(panicServer)
	}
	if DHI0_Addr2 != "" && DHI0_Addr2_Crt != "" && DHI0_Addr2_Key != "" {
		panicServerTLS := MakeServer(DHI0_Addr2)
		servers = append(servers, panicServerTLS)
		serverManager.AddAsyncServerTLS(DHI0_Addr2_Crt, DHI0_Addr2_Key, panicServerTLS)
	}

	if len(servers) < 1 {
		requestMap["StartupCode"] = "500"
		requestMap["StartupNote"] = "HTTP and HTTPS addresses not configured"
		Flap <- requestMap
		fmt.Println("Setting StartupCode to 500!")
		return
	}

	serverManager.Run()

	if DHI0_RedirectHTTP && !regexp.MustCompile(`^https\:\/\/.+$`).MatchString(DHI0_RedirectDestination) {
		requestMap["StartupCode"] = "500"
		requestMap["StartupNote"] = "Conf parameter DHI0_RedirectDestination not valid"
		fmt.Println("Setting StartupCode to 500!")
		Flap <- requestMap
		return
	}

	fmt.Println("Setting StartupCode to 200!")

	requestMap["StartupCode"] = "200"
	Flap <- requestMap

	/***2***/
	for !serverManager.isDone() {
		select {
		case <-Clap:
			{
				for _, server := range servers {
					server.Close()
				}
			}
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}

	return serverManager.completion.getErrors()
}

type PanicManager struct{}

func (d *PanicManager) ServeHTTP(R http.ResponseWriter, r *http.Request) {
	/***1***/
	if r.TLS == nil && DHI0_RedirectHTTP {
		http.Redirect(R, r, DHI0_RedirectDestination, http.StatusTemporaryRedirect)
	}
	/***2***/
	xb05 := map[string]any{}
	xb05["ExecutionOutcomeCode"] = 500
	defer func() {
		/***1***/
		xc01 := recover()
		if xc01 != nil {
			xb05["ExecutionOutcomeCode"] = 500
			xb05["ExecutionOutcomeNote"] = fmt.Sprintf(
				`Panic sighted [%v : %s]`, xc01, string(debug.Stack()),
			)
		}
		/***2***/
		xc05 := xb05["ExecutionOutcomeCode"].(int)
		if !slices.Contains(DNI0_AllowedResponseCode, xc05) && xc05 != 500 {
			xb05["ExecutionOutcomeCode"] = 500
			xb05["ExecutionOutcomeNote"] = fmt.Sprintf(
				`Unexpected response code %d`, xc05,
			)
		}
		/***3***/
		if xb05["ExecutionOutcomeCode"].(int) == 500 {
			xd05, xd10 := xb05["ExecutionOutcomeNote"].(string)
			if !xd10 {
				xd05 = "Execution Outcome Note not a string"
			}
			Output_Logg("ERR", "DHI2", xd05)
			delete(xb05, "ExecutionOutcomeNote")
		}
		/***4***/
		for _, xd05 := range DNI0_ResponseHeaders {
			R.Header().Set(xd05[0], xd05[1])
		}
		/***5***/
		xc10, _ := json.MarshalIndent(xb05, "", "    ")
		xc10 = append(xc10, '\n')
		R.Write(xc10)
	}()
	/***3***/
	xb15, xb20 := io.ReadAll(r.Body)
	if xb20 != nil {
		xb05["ExecutionOutcomeCode"] = 500
		xb05["ExecutionOutcomeNote"] = fmt.Sprintf(
			`Request read failed [%s]`, xb20.Error(),
		)
		return
	}
	if !json.Valid(xb15) {
		xb05["ExecutionOutcomeCode"] = 400
		xb05["ExecutionOutcomeNote"] = fmt.Sprintf(`Request JSON formatting invalid`)
		return
	}
	xb25 := &Request{}
	xb30 := json.Unmarshal(xb15, xb25)
	if xb30 != nil {
		xb05["ExecutionOutcomeCode"] = 400
		xb05["ExecutionOutcomeNote"] = fmt.Sprintf(
			`Request unmarshal failed [%s]`, xb30.Error(),
		)
		return
	}
	if xb25.SrID == "" {
		xb05["ExecutionOutcomeCode"] = 400
		xb05["ExecutionOutcomeNote"] = fmt.Sprintf(`No service specified`)
		return
	}

	/***4***/
	xb35, xb40, xb45 := route(r, xb25, R)
	xb05["ExecutionOutcomeCode"] = xb35
	xb05["ExecutionOutcomeNote"] = xb40
	if xb45 != nil {
		xb05["Yield"] = xb45
	}
}

func route(r *http.Request, s *Request, R http.ResponseWriter) (
	C int, N string, Y any,
) {
	fmt.Println("Handling Request...")
	
	/***1***/
	C = 500
	var ServiceProvider *ResponseHandler
	for _, xc10 := range DHI0_SPRegister {
		if s.SrID == xc10.Code {
			ServiceProvider = xc10
		}
	}
	if ServiceProvider == nil {
		C = 400
		N = fmt.Sprintf(`Service specified not supported`)
		return
	}
	/***2***/
	C, N, Y = ServiceProvider.Program(r, s.SrID, s.Seed)
	return
}
