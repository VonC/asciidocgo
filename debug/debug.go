// to check: https://gist.github.com/alphazero/2718939
// For now:
// - http://play.golang.org/p/8o0WywyaDT
// - http://play.golang.org/p/QFheQeChIn
// - http://play.golang.org/p/p5Z8X3nxXL <===
// All from https://groups.google.com/forum/#!topic/golang-nuts/ct99dtK2Jo4/discussion
package debug

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var debug = *flag.Bool("d", false, "turn on debug info")
var debugLog = log.New(os.Stderr, "", log.LstdFlags)

func Switch() {
	debug = !debug
	if debug {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func Debug(msg string) {
	log.Println("DEBUG: " + msg)
}
