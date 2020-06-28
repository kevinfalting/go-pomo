package session

import (
	"fmt"
	"time"
)

// State type
type State int

// Represents the various states the application can be in.
const (
	StateFocus State = iota
	StateShort State = iota
	StateLong  State = iota
)

// ConfigTime is just a time.Duration
type ConfigTime time.Duration

// Config holds all the configuration information from flags
type Config struct {
	Rounds  int
	Short   ConfigTime
	Long    ConfigTime
	Focus   ConfigTime
	Auto    bool
	Sounds  bool
	Seconds bool
	Unit    time.Duration
}

// stats contain statistics for the session so that we can generate a report.
type stats struct {
	startTime     time.Time
	pausedTime    time.Duration
	focusTime     time.Duration
	breakTime     time.Duration
	elapsedRounds int
}

// Session is the source of truth for the state of the program
type Session struct {
	config Config
	stats  stats

	state          State
	stateStartTime time.Time

	paused         bool
	pauseStartTime time.Time
}

// Init initializes the Session
func (s *Session) Init(c Config) {
	s.config.Rounds = c.Rounds
	s.config.Short = c.Short
	s.config.Long = c.Long
	s.config.Focus = c.Focus
	s.config.Auto = c.Auto
	s.config.Sounds = c.Sounds

	s.stats.startTime = time.Now()
	s.stats.elapsedRounds = 0

	s.state = StateFocus
	s.stateStartTime = time.Now()
	s.config.Unit = 1 * time.Minute
	if c.Seconds {
		s.config.Unit = 1 * time.Second
	}
}

func (s Session) String() string {
	// TODO: return a formatted report of all the statistics from the session.
	r := "\n----Session Stats----\n"
	r += fmt.Sprintf("Time in Focus: %v\n", s.stats.focusTime)
	r += fmt.Sprintf("Time spent paused: %v\n", s.stats.pausedTime)
	r += fmt.Sprintf("Time spent in breaks: %v\n", s.stats.breakTime)
	r += fmt.Sprintf("Total elapsed time: %v\n", time.Now().Sub(s.stats.startTime))
	r += fmt.Sprintf("Total elapsed rounds: %d\n", s.stats.elapsedRounds)

	r += fmt.Sprintf("Current State: %d\n", s.state)

	return r
}

// GetState returns the current state of the session. Check IsPaused to see if it's paused or not.
func (s *Session) GetState() State {
	return s.state
}

// GetStateStartTime returns the start time for the current state
func (s *Session) GetStateStartTime() time.Time {
	return s.stateStartTime
}

// IsPaused returns if the session is paused or not.
func (s *Session) IsPaused() bool {
	return s.paused
}

// Pause handles pausing the session
func (s *Session) Pause() {
	fmt.Println("Paused...")
	s.pauseStartTime = time.Now()
	s.paused = true
}

// Unpause handles unpausing the session
func (s *Session) Unpause() {
	fmt.Println("Resuming...")
	p := time.Now().Sub(s.pauseStartTime)
	s.stats.pausedTime += p
	s.stateStartTime = s.stateStartTime.Add(p) // Remove the time paused from the current state
	s.paused = false
}

// EndSession prepares for closing the application
func (s *Session) EndSession() {
	if s.IsPaused() {
		s.Unpause()
	}
	t := time.Now().Sub(s.stateStartTime)
	switch s.state {
	case StateFocus:
		s.stats.focusTime += t
	case StateShort, StateLong:
		s.stats.breakTime += t
	}
}

// IsRoundOver returns true if the current round is over
func (s *Session) IsRoundOver() bool {
	var l ConfigTime
	switch s.state {
	case StateFocus:
		l = s.config.Focus
	case StateShort:
		l = s.config.Short
	case StateLong:
		l = s.config.Long
	}

	return time.Now().Sub(s.stateStartTime) > time.Duration(l)*s.config.Unit
}

// GoToNextState progresses the session to the next state
func (s *Session) GoToNextState() {
	r := (s.stats.elapsedRounds + 1) % s.config.Rounds
	t := time.Now().Sub(s.stateStartTime)
	switch s.state {
	case StateFocus:
		s.stats.focusTime += t
		if r == 0 {
			s.state = StateLong
		} else {
			s.state = StateShort
		}
	case StateShort, StateLong:
		s.stats.breakTime += t
		s.state = StateFocus
	}

	s.stats.elapsedRounds++
	s.stateStartTime = time.Now()
}

// ShouldAutoProgress returns if the application should automatically progress to the next state without user input
func (s *Session) ShouldAutoProgress() bool {
	return s.config.Auto
}
