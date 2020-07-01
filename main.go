package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kevinfalting/go-pomo/session"
)

var (
	rounds = flag.Int("rounds", 3, "The number of focus rounds per cycle/long break.")
	short  = flag.Int("short", 5, "The duration of the short breaks in minutes.")
	long   = flag.Int("long", 15, "The duration of the long breaks in minutes.")
	focus  = flag.Int("focus", 25, "The number of minutes per focus session.")
	auto   = flag.Bool("auto", false, "Auto progress to next round without input.")
	// sounds  = flag.Bool("sounds", true, "Play sounds to indicate round changes.")
	seconds = flag.Bool("seconds", false, "Will instead count the break time in seconds instead of minutes")

	done    = make(chan bool)
	proceed = make(chan bool)

	allowProceed = false

	// S is the global state
	S = session.Session{}
)

func main() {
	flag.Parse()

	// Wait for user input to start
	waitForInput("Press enter to begin...")

	// Start the session
	config := session.Config{}
	config.Rounds = *rounds
	config.Short = session.ConfigTime(*short)
	config.Long = session.ConfigTime(*long)
	config.Focus = session.ConfigTime(*focus)
	config.Auto = *auto
	// config.Sounds = *sounds
	config.Seconds = *seconds
	S.Init(config)
	go tickTock()

	// loop as long as the user wants
	for {
		waitForInput("p: pause, r: resume, s: show stats, q: quit")
	}
}

func waitForInput(message string) {
	fmt.Println(message)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		exit(1, err)
	}

	input = strings.Trim(input, "\n")
	input = strings.ToLower(input)

	switch input {
	case "":
		// Proceeding
		if !S.GetStateStartTime().IsZero() && !S.IsPaused() && allowProceed {
			proceed <- true
		}
	case "p":
		// pause the timer
		S.Pause()
	case "q":
		exit(0, nil)
	case "r":
		// resume session
		S.Unpause()
		if !S.GetStateStartTime().IsZero() && !S.IsPaused() {
			proceed <- true
		}
	case "s":
		// show stats, make no state changes
		fmt.Println(S)
	default:
		waitForInput(fmt.Sprintf("Sorry, I don't know what to do with %q\n", input))
	}
}

// tickTock handles decisions based on the clock
func tickTock() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case t := <-ticker.C:

			// if the state is paused, we don't want to allow any other state changes to occur
			if S.IsPaused() {
				continue
			}

			if S.IsRoundOver() {
				if !S.ShouldAutoProgress() {
					fmt.Println("Press enter to proceed.")
					allowProceed = true
					<-proceed
					allowProceed = false
					fmt.Println("Proceeding....")
				}
				S.GoToNextState()
			}

			// TODO: Do something here to show progress... like a progress bar or how much time is left? But replace the line, don't keep writing newlines.

			_ = t // remove after testing

		}
	}
}

func exit(statusCode int, err error) {
	// If you try to quit at the first prompt, don't pass a done signal. Will panic.
	if !S.GetStateStartTime().IsZero() {
		done <- true
	}
	S.EndSession()
	fmt.Println(S)

	if err != nil {
		fmt.Println(err)
	}
	os.Exit(statusCode)
}
