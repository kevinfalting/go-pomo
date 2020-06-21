package main

import (
	"flag"
)

var (
	rounds = flag.Int("rounds", 3, "The number of study rounds per cycle/long break.")
	short  = flag.Int("short", 5, "The duration of the short breaks in minutes.")
	long   = flag.Int("long", 15, "The duration of the long breaks in minutes.")
	study  = flag.Int("study", 25, "The number of minutes per study session.")
)

func main() {
	flag.Parse()
}
