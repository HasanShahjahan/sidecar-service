package proxy

import (
	"fmt"
	"github.com/HasanShahjahan/sidecar-service/internal/config"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

//NewProxy - initialize proxy with configurable values
func NewProxy(allowedList []*regexp.Regexp) *Proxy {
	proxy := &Proxy{
		ServiceName: config.Config.ServiceName,
		Downstream: &Downstream{
			Address:     config.Config.DownStreamURL,
			AllowedList: allowedList,
		},
		Handler: nil,
		Timeout: &config.Config.Timeout,
		LogPrefix: config.Config.LogPrefix,
		ErrorLogger:          nil,
		AccessLogger:         nil,
		FailedRequestsCount:  0,
		TotalRequestsCount:   0,
		InvalidRequestsCount: 0,
	}

	proxy.initializeHandler()
	return proxy
}
//ServeHTTP - create a Server with given address and handler
func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxy.TotalRequestsCount++
	if proxy.ValidatePath(r.RequestURI) {
		proxy.Handler.ServeHTTP(w, r)
		return
	}

	proxy.InvalidRequestsCount++
	fmt.Fprintf(
		proxy.AccessLogger,
		"%s %s %s %s %s %d Path is not allowed by proxy!\n",
		proxy.LogPrefix,
		time.Now().Format(time.RFC3339Nano),
		r.RemoteAddr,
		r.Method,
		r.RequestURI,
		http.StatusNotFound,
	)

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, PathNotFound)
}

//ValidatePath - validates path based on down stream allowed list.
func (proxy *Proxy) ValidatePath(path string) bool {
	for _, allowedPath := range proxy.Downstream.AllowedList {
		if allowedPath.MatchString(strings.Trim(strings.ToLower(strings.Split(path, "?")[0]), "/")) {
			return true
		}
	}
	return false
}

func (proxy *Proxy) initializeHandler() {
	url, _ := url.Parse(proxy.Downstream.Address)
	proxy.Handler = httputil.NewSingleHostReverseProxy(url)
	proxy.Handler.ErrorLog = log.New(proxy.ErrorLogger, proxy.LogPrefix, log.LstdFlags)

	if false {
		proxy.Handler.Director = func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", proxy.Downstream.Address)
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			req.URL.Path = url.Path + req.URL.Path
		}
	}

	if proxy.Timeout != nil {
		proxy.Handler.Transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: *proxy.Timeout,
			}).DialContext,
		}
	}

	proxy.Handler.ModifyResponse = func(r *http.Response) error {
		if proxy.AccessLogger != nil {
			fmt.Fprintf(
				proxy.AccessLogger,
				"%s %s %s %s %s %s\n",
				proxy.LogPrefix,
				time.Now().Format(time.RFC3339Nano),
				r.Request.RemoteAddr,
				r.Request.Method,
				r.Request.RequestURI,
				r.Status,
			)
		}

		if r.StatusCode >= 500 {
			proxy.FailedRequestsCount++
		}

		return nil
	}
	proxy.Handler.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("error: %+v", err)
		if proxy.ErrorLogger != nil {
			fmt.Fprintln(proxy.ErrorLogger, err)
		}

		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
}

func (proxy *Proxy) initializeLog(AccessLogFile string, ErrorLogFile string) {
	if ErrorLogFile != "" {
		f, err := os.OpenFile(ErrorLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		proxy.ErrorLogger = f
	}

	if AccessLogFile != "" {
		f, err := os.OpenFile(AccessLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		proxy.AccessLogger = f
	}
}