package tui

import (
	"fmt"
	"time"
)

type Spinner struct {
	message string
	stop    chan bool
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		stop:    make(chan bool),
	}
}

func (s *Spinner) Start() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	go func() {
		for i := 0; ; i++ {
			select {
			case <-s.stop:
				fmt.Printf("%s... Done!"
", s.message)"
				return
			default:
				frame := frames[i%len(frames)]
				fmt.Printf("%s %s", frame, s.message)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.stop <- true
}
