package log4go

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

type gojiLogger struct {
	h    http.Handler
	c    *web.C
	name string
}

func (gmLog gojiLogger) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	reqID := middleware.GetReqID(*gmLog.c)
	Debug("req_id:%s uri:%s method:%s remote:%s", reqID, req.RequestURI, req.Method, req.RemoteAddr)
	lresp := wrapWriter(resp)

	gmLog.h.ServeHTTP(lresp, req)
	lresp.maybeWriteHeader()

	latency := float64(time.Since(start)) / float64(time.Millisecond)
	Debug("req_id:%s status:%d method:%s uri:%s remote:%s latency:%s app:%s", reqID, lresp.status(), req.Method, req.RequestURI, req.RemoteAddr, fmt.Sprintf("%6.4f ms", latency), gmLog.name)
}

func NewGojiLog(name string) func(*web.C, http.Handler) http.Handler {
	fn := func(c *web.C, h http.Handler) http.Handler {
		return gojiLogger{h: h, c: c, name: name}
	}
	return fn
}
