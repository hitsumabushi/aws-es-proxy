package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"go.uber.org/zap"
)

const (
	service = "es"
	region  = "us-west-2"
)

var logger *zap.Logger
var cred *credentials.Credentials
var signer *v4.Signer
var userAgent string

func init() {
	s := session.Must(session.NewSession())
	cred = s.Config.Credentials
	signer = v4.NewSigner(cred)
}

func getUserAgent() string {
	if userAgent == "" {
		userAgent = "aws-es-proxy-go/" + version + "_" + commit
	}
	return userAgent
}

func newReverseProxy(c Config) *ReverseProxy {
	// disable http/2
	http.DefaultTransport.(*http.Transport).TLSNextProto = map[string]func(authority string, c *tls.Conn) http.RoundTripper{}

	director := func(req *http.Request) {
		logger.Debug("access", zap.String("path", req.URL.Path), zap.String("query", req.URL.RawQuery))
		// request routing
		req.URL.Scheme = "https"
		// check path match
		for k, v := range c.ServerMap {
			if strings.HasPrefix(req.URL.Path, k) {
				// change Host
				req.Host = v.Host
				req.URL.Host = v.Host
				// update Path
				req.URL.Path = strings.Replace(req.URL.Path, k, "", 1)
				req.RequestURI = strings.Replace(req.RequestURI, k, "", 1)
				logger.Debug("match key", zap.String("key", k), zap.String("entry", v.String()), zap.String("New request path", req.URL.Path))
				req.Header.Set("User-Agent", getUserAgent())
				break
			}
		}
	}

	return &ReverseProxy{
		Director: director,
		Signer:   signer,
		Service:  "es",
		Config:   &c,
	}
}

func main() {
	configPath := flag.String("config", "", "config file path")
	port := flag.Uint("port", 8080, "listening port")
	debug := flag.Bool("debug", false, "debug log setting")
	flag.Parse()

	// check flags
	if *configPath == "" {
		fmt.Println("Please path -config parameter")
		os.Exit(1)
	}

	// load config
	config := loadConfig(*configPath)

	// configure logger
	if *debug {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	logger.Info("Start server",
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("buildGoVersion", buildGoVersion),
		zap.String("buildDate", buildDate),
		zap.String("config", *configPath),
		zap.Uint("port", *port),
	)

	// proxy
	pxy := newReverseProxy(config)

	// http.HandleFunc("/_list", func(w http.ResponseWriter, r *http.Request) {
	// 	for k, v := range config.ServerMap {
	// 		fmt.Fprintf(w, "%s -> %s\n", k, v.String())
	// 	}
	// })
	http.Handle("/", pxy)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
