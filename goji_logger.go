package log4go

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/peterbourgon/g2s"
)

var LastPathParamIntegerRegex = regexp.MustCompile(`\_[0-9]{1,}$`)
var PathParamIntegerRegex = regexp.MustCompile(`\_[0-9]{1,}\_`)

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

	writeHttpLoging(fmt.Sprintf("Started %s '%s' from %s", req.Method, req.URL.Path, req.RemoteAddr))
	lresp := wrapWriter(resp)
	startAt := time.Now()
	g.h.ServeHTTP(lresp, req)
	lresp.maybeWriteHeader()

	latency := float64(time.Since(startAt)) / float64(time.Millisecond)
	writeHttpLoging(fmt.Sprintf("Returning %d in %s", lresp.status(), fmt.Sprintf("%6.4f ms", latency)))
}

func (g gojiStatds) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	startAt := time.Now()
	if g.WriteLog {
		writeHttpLoging(fmt.Sprintf("Started %s '%s' from %s", req.Method, req.URL.Path, req.RemoteAddr))
	}
	lresp := wrapWriter(resp)

	g.h.ServeHTTP(lresp, req)
	lresp.maybeWriteHeader()
	elapsedTime := time.Since(startAt)
	latency := float64(elapsedTime) / float64(time.Millisecond)

	if g.WriteLog {
		writeHttpLoging(fmt.Sprintf("Returning %d in %s", lresp.status(), fmt.Sprintf("%6.4f ms", latency)))
	}
	statsd, err := g2s.Dial("udp", g.Config.IpPort)
	if err != nil {
		return
	}

	apiName := ReplaceIntegerPathParameters(req.URL.Path, PathParamIntegerRegex, LastPathParamIntegerRegex)
	keyBase := ""
	if g.Config.Prefix != "" {
		keyBase += g.Config.Prefix + "."
	}
	keyBase += apiName + "."
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

func writeHttpLoging(msg string) {

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

func ReplaceIntegerPathParameters(input string, pathRegex *regexp.Regexp, lastPathRegex *regexp.Regexp) (output string) {
	output = strings.Replace(input, "/", "_", -1)
	matchs := pathRegex.FindAllStringSubmatch(output, -1)

	for _, match := range matchs {
		if match[0] != "" {
			output = strings.Replace(output, match[0], "_id_", -1)
		}
	}
	output = replaceLastIntegerPathParameters(output, lastPathRegex)
	return output
}

func replaceLastIntegerPathParameters(input string, r *regexp.Regexp) (output string) {
	match := r.FindStringSubmatch(input)
	if len(match) > 0 && match[0] != "" {
		output = strings.Replace(input, match[0], "_id", -1)
	}
	return output
}
