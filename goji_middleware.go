package log4go

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

type gojiLogger struct {
	h http.Handler
	c *web.C
}

func (gmLog gojiLogger) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	reqID := middleware.GetReqID(*gmLog.c)
	Debug("[%s] Started %s '%s' from %s", reqID, req.Method, req.RequestURI, req.RemoteAddr)
	lresp := wrapWriter(resp)

	gmLog.h.ServeHTTP(lresp, req)
	lresp.maybeWriteHeader()

	latency := float64(time.Since(start)) / float64(time.Millisecond)
	Debug("[%s] Returning %d in %s", reqID, lresp.status(), fmt.Sprintf("%6.4f ms", latency))
}

func NewGojiLog() func(*web.C, http.Handler) http.Handler {
	fn := func(c *web.C, h http.Handler) http.Handler {
		return gojiLogger{h: h, c: c}
	}
	return fn
}
