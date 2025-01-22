// Formatted with gofmt -s
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
				fmt.Printf("\r%s... Done!\n", s.message)
				return
			default:
				frame := frames[i%len(frames)]
				fmt.Printf("\r%s %s", frame, s.message)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.stop <- true
}
