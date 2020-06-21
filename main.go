package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	rounds = flag.Int("rounds", 3, "The number of focus rounds per cycle/long break.")
	short  = flag.Int("short", 5, "The duration of the short breaks in minutes.")
	long   = flag.Int("long", 15, "The duration of the long breaks in minutes.")
	study  = flag.Int("study", 25, "The number of minutes per focus session.")
)

type session struct {
	startTime  time.Time
	endTime    time.Time
	pausedTime time.Duration
	focusTime  time.Duration
	rounds     int
}

func (s session) elapsed() time.Duration {
	return time.Since(s.startTime)
}

func main() {
	flag.Parse()

	// Wait for user input to start
	waitForUserInputToProceed()

	s := session{}
	s.startTime = time.Now()

	// start a timer
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	go tickTock(ticker, done)

	// TODO: present options
	// p - pause
	// s - stop
	// g - go/start
	waitForUserInputToProceed()
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
	fmt.Println(s.elapsed())
}

func waitForUserInputToProceed() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to begin...")
	_, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func tickTock(ticker *time.Ticker, done chan bool) {
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			fmt.Println("Tick at", t)
		}
	}
}
