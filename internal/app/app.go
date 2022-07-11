package app

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dwilkolek/moxy/config"
	"github.com/elliotchance/sshtunnel"
)

func Run(config *config.MoxyConfig) {
	var wg sync.WaitGroup

	tunnelPort := setupTunnel(config.Tunnel)

	wg.Add(1)
	for service, conf := range config.Services {
		setupHttpServerForService(service, conf, fmt.Sprintf("http://localhost:%d", tunnelPort))
	}
	wg.Wait()
}

func setupTunnel(tunnelConfig config.TunnelConfig) int {
	tunnel := sshtunnel.NewSSHTunnel(
		tunnelConfig.UserAndHost,
		sshtunnel.PrivateKeyFile(tunnelConfig.PathToPrivateKey),
		tunnelConfig.Destination,
		"0",
	)

	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

	go tunnel.Start()
	time.Sleep(100 * time.Millisecond)
	tunnel.Log.Printf("Tunnel started and exposed on port: %d\n", tunnel.Local.Port)

	return tunnel.Local.Port
}

func setupHttpServerForService(service string, conf config.ProxyConfig, to string) {
	go func() {
		logger := log.New(os.Stdout, service+"\t", log.Ldate|log.Lmicroseconds)
		origin, _ := url.Parse(to)

		logger.Printf("Setting up at port: %d \n", conf.Port)

		director := func(req *http.Request) {
			for header, value := range conf.Headers {
				req.Header.Add(header, value)
				if strings.ToLower(header) == "host" {
					req.Host = value
				}
			}

			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}

		proxy := &httputil.ReverseProxy{Director: director}
		server := http.NewServeMux()

		server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s %s %s\n", r.Method, r.URL.Path, r.URL.RawQuery)
			defer r.Body.Close()
			if conf.AllowCors {
				if r.Method == http.MethodOptions {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
					w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
					w.WriteHeader(http.StatusNoContent)
				} else {
					// Set CORS headers for the main request.
					w.Header().Set("Access-Control-Allow-Origin", "*")
					proxy.ServeHTTP(w, r)
				}
			} else {
				proxy.ServeHTTP(w, r)
			}

		})

		logger.Printf("Starting server at port: %d\n", conf.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), server); err != nil {
			log.Fatal(err)
		}
	}()
}
