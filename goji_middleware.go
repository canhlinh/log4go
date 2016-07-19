package log4go

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/peterbourgon/g2s"
)

type StatsdConfig struct {
	// IP and port of the statsd server. Optional. Default to "127.0.0.1:8125".
	IpPort string

	// Prefix added to the metric keys. Optional.
	Prefix string
}

type gojiLogger struct {
	h http.Handler
}

type gojiStatds struct {
	h        http.Handler
	WriteLog bool
	Config   StatsdConfig
}

func (g gojiLogger) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	writeGoji(fmt.Sprintf("Started %s '%s' from %s", req.Method, req.RequestURI, req.RemoteAddr))
	lresp := wrapWriter(resp)
	start := time.Now()
	g.h.ServeHTTP(lresp, req)
	lresp.maybeWriteHeader()

	latency := float64(time.Since(start)) / float64(time.Millisecond)
	writeGoji(fmt.Sprintf("Returning %d in %s", lresp.status(), fmt.Sprintf("%6.4f ms", latency)))
}

func (g gojiStatds) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	if g.WriteLog {
		writeGoji(fmt.Sprintf("Started %s '%s' from %s", req.Method, req.RequestURI, req.RemoteAddr))
	}
	lresp := wrapWriter(resp)

	g.h.ServeHTTP(lresp, req)
	lresp.maybeWriteHeader()
	elapsedTime := time.Since(start)
	latency := float64(elapsedTime) / float64(time.Millisecond)

	if g.WriteLog {
		writeGoji(fmt.Sprintf("Returning %d in %s", lresp.status(), fmt.Sprintf("%6.4f ms", latency)))
	}
	statsd, err := g2s.Dial("udp", g.Config.IpPort)
	if err != nil {
		return
	}

	apiURI := strings.Replace(req.RequestURI, "/", "_", -1)
	keyBase := ""
	if g.Config.Prefix != "" {
		keyBase += g.Config.Prefix + "."
	}
	keyBase += apiURI + "."
	keyBase += req.Method + "."
	statsd.Counter(1.0, keyBase+"status_code."+strconv.Itoa(lresp.status()), 1)
	statsd.Timing(1.0, keyBase+"elapsed_time.", elapsedTime)
}

func NewGojiLog() func(http.Handler) http.Handler {

	fn := func(h http.Handler) http.Handler {
		return gojiLogger{h: h}
	}
	return fn
}

func NewGojiStatsd(config StatsdConfig, writeLog bool) func(http.Handler) http.Handler {

	if len(config.IpPort) == 0 {
		config.IpPort = "127.0.0.1:8125"
	}
	fn := func(h http.Handler) http.Handler {
		return gojiStatds{h: h, WriteLog: writeLog, Config: config}
	}
	return fn
}

func writeGoji(msg string) {

	msg = Global.formatColor(DEBUG, msg)
	rec := &LogRecord{
		Level:   DEBUG,
		Created: time.Now(),
		Source:  "Web Controller",
		Message: msg,
	}
	for _, filt := range Global {
		filt.LogWrite(rec)
	}
}
