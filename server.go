package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"slices"
)

type Server struct {
	serviceMap       ServiceMap
	allowedRespCodes []int
}

func MakeServer(addr string, serviceMap ServiceMap, allowedResponseCodes []int) *http.Server {
	if addr == "" {
		panic("Error: no server address supplied")
	}

	server := &http.Server{Addr: addr, Handler: &Server{
		serviceMap:       serviceMap,
		allowedRespCodes: allowedResponseCodes,
	}}

	server.MaxHeaderBytes = MaxHeaderSize
	server.ReadTimeout = ReadTimeout
	server.ReadHeaderTimeout = ReadTimeout
	server.WriteTimeout = WriteTimeout
	server.IdleTimeout = IdleTimeout
	error_buff := bytes.NewBuffer([]byte{})
	server.ErrorLog = log.New(error_buff, "", log.Lshortfile)

	return server
}

func (server *Server) sendHTTPImpl(resp http.ResponseWriter, responseMap map[string]any) {
	/***1***/
	panicRes := recover()
	if panicRes != nil {
		responseMap["ExecutionOutcomeCode"] = 500
		responseMap["ExecutionOutcomeNote"] = fmt.Sprintf(
			`Panic sighted [%v : %s]`, panicRes, string(debug.Stack()),
		)
	}
	/***2***/
	handlerCode := responseMap["ExecutionOutcomeCode"].(int)
	if !slices.Contains(server.allowedRespCodes, handlerCode) && handlerCode != 500 {
		responseMap["ExecutionOutcomeCode"] = 500
		responseMap["ExecutionOutcomeNote"] = fmt.Sprintf(
			`Unexpected response code %d`, handlerCode,
		)
	}
	/***3***/
	if responseMap["ExecutionOutcomeCode"].(int) == 500 {
		xd05, xd10 := responseMap["ExecutionOutcomeNote"].(string)
		if !xd10 {
			xd05 = "Execution Outcome Note not a string"
		}
		Output_Logg("ERR", "DHI2", xd05)
		delete(responseMap, "ExecutionOutcomeNote")
	}
	/***4***/
	resp.Header().Set("Content-Type", "application/json")

	/***5***/
	jsonResp, _ := json.MarshalIndent(responseMap, "", "    ")
	jsonResp = append(jsonResp, '\n')
	resp.Write(jsonResp)
}

func (server *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	/***2***/
	responseMap := map[string]any{}
	responseMap["ExecutionOutcomeCode"] = 500

	defer server.sendHTTPImpl(resp, responseMap)

	//read the http body
	bodyBuff, ok := io.ReadAll(req.Body)
	if ok != nil {
		responseMap["ExecutionOutcomeCode"] = 500
		responseMap["ExecutionOutcomeNote"] = fmt.Sprintf(
			`Request read failed [%s]`, ok.Error(),
		)
		return
	}

	//validate the json format
	if !json.Valid(bodyBuff) {
		responseMap["ExecutionOutcomeCode"] = 400
		responseMap["ExecutionOutcomeNote"] = "Request JSON formatting invalid"
		return
	}

	//parse the json
	parsedJson := &RequestJson{}
	ok = json.Unmarshal(bodyBuff, parsedJson)
	if ok != nil {
		responseMap["ExecutionOutcomeCode"] = 400
		responseMap["ExecutionOutcomeNote"] = fmt.Sprintf(
			`Request unmarshal failed [%s]`, ok.Error(),
		)
		return
	}

	/***3***/
	handlerRes := server.route(parsedJson)
	responseMap["ExecutionOutcomeCode"] = handlerRes.Code
	responseMap["ExecutionOutcomeNote"] = handlerRes.Error
	if handlerRes.Yield != nil {
		responseMap["Yield"] = handlerRes.Yield
	}
}

func (s *Server) route(requestJson *RequestJson) ServiceResult {
	handler, ok := s.serviceMap[requestJson.HandlerID]
	if !ok {
		return ServiceResult{
			Code:  500,
			Error: fmt.Sprintf("Service (%v) is not supported", requestJson.HandlerID),
		}
	}
	return handler.Run(*requestJson)
}
