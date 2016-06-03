package main
import log "github.com/canhlinh/log4go"

func main() {
	// Load the configuration (isn't this easy?)
	log.LoadConfiguration("example.xml")

	// And now we're ready!
	log.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	log.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
	log.Info("About that time, eh chaps?")
}
