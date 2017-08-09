package log4go

import "time"

import "github.com/getsentry/raven-go"

//SentryLogWriter This is the standard writer that prints to standard output.
type SentryLogWriter struct {
	format string
	w      chan *LogRecord
	o      chan string
	client *raven.Client
}

//NewSentryLogWriter This creates a new SentryLogWriter
func NewSentryLogWriter(dsn string) *SentryLogWriter {
	rClient, err := raven.New(dsn)
	if err != nil {
		panic(err)
	}

	sentryWriter := &SentryLogWriter{
		format: FORMAT_DEFAULT,
		w:      make(chan *LogRecord, LogBufferLength),
		client: rClient,
	}
	go sentryWriter.run(nil)
	return sentryWriter
}
func (c *SentryLogWriter) SetFormat(format string) {
	c.format = format
}
func (c *SentryLogWriter) run(o chan string) {
	for rec := range c.w {
		eventID := c.client.CaptureMessage(rec.Message, map[string]string{
			"level":  rec.Level.String(),
			"source": rec.Source,
		})
		if o != nil {
			o <- eventID
		}
	}
}

//LogWrite This is the SentryLogWriter's output method.  This will block if the output
// buffer is full.
func (c *SentryLogWriter) LogWrite(rec *LogRecord) {
	c.w <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (c *SentryLogWriter) Close() {
	close(c.w)
	time.Sleep(50 * time.Millisecond) // Try to give console I/O time to complete
}
