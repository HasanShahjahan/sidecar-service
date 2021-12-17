package proxy

import (
	"net/http/httputil"
	"os"
	"regexp"
	"time"
)

const (
	PathNotFound = "Not found!"
	IdFormat     = "[a-z]+[0-9]+[a-z0-9]+"
)

//Downstream - defines downstream address and allowed list
type Downstream struct {
	Address     string
	AllowedList []*regexp.Regexp
}

//Proxy - defines downstream details with reverse proxy
type Proxy struct {
	ServiceName string
	Downstream  *Downstream
	Handler     *httputil.ReverseProxy
	Timeout      *time.Duration
	LogPrefix    string
	ErrorLogger  *os.File
	AccessLogger *os.File
	FailedRequestsCount  int
	TotalRequestsCount   int
	InvalidRequestsCount int
}
