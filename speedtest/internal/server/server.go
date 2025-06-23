package server

import (
	"net/http"
	"speedtest/internal/config"
	"speedtest/internal/speedtester"
)

func Routes(tester *speedtester.Tester, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /speedtest", handleSpeedTest(tester, cfg))
	return mux
}
