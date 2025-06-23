package server

import (
	"net/http"
	"speedtest/internal/speedtester"
)

func Routes(tester *speedtester.Tester) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /speedtest", handleSpeedTest(tester))
	return mux
} 