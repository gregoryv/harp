package warp

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// debug logger is global to one package
var debug = log.New(ioutil.Discard, "D ", log.LstdFlags|log.Lshortfile)

func init() {
	// enable debug log if environment D=t is set
	if yes, _ := strconv.ParseBool(os.Getenv("D")); yes {
		debug.SetOutput(os.Stderr)
	}
}

func SetDebugOutput(w io.Writer) {
	debug.SetOutput(w)
}
