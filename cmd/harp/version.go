package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gregoryv/harp"
)

func version() {
	if !showVersion {
		return
	}
	fmt.Println(harp.Version())
	os.Exit(0)
}

var showVersion bool

func init() {
	flag.BoolVar(&showVersion, "v", showVersion, "show version and exit")
}
