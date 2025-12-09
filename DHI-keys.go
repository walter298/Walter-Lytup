package main

import  "time"

const (
	DHI0_Addr1 string = ":80"
	DHI0_Addr2 string = ":443"
	DHI0_Addr2_Key string = "tls.key"
	DHI0_Addr2_Crt string = "tls.crt"
	DHI0_RedirectHTTP bool = true
	DHI0_RedirectDestination string = "https://localhost"
	DHI0_MaxHeaderSize int = 1*1024 * 1024
	DHI0_ReadTimeout time.Duration  = time.Minute*5
	DHI0_WrttTimeout time.Duration  = time.Minute*5
	DHI0_IdleTimeout time.Duration  = time.Minute*5
	TimeZoneSecondOffset int = 0
)