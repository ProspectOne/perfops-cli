// Copyright 2017 The PerfOps-CLI Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"fmt"
	"sync"
	"time"
	"unicode/utf8"
)

// Spinner represents the indicator.
type Spinner struct {
	mu        sync.Mutex
	frames    []rune
	length    int
	pos       int
	active    bool
	lastFrame int
	stopChan  chan struct{}
}

// NewSpinner returns a spinner.
func NewSpinner() *Spinner {
	const frames = `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`
	s := &Spinner{
		frames:    []rune(frames),
		length:    len(frames),
		lastFrame: -1,
		stopChan:  make(chan struct{}, 1),
	}
	s.length = len(s.frames)
	return s
}

// Start will start the indicator.
func (s *Spinner) Start() {
	if s.active {
		return
	}
	s.active = true
	s.lastFrame = -1
	s.pos = 0
	go func() {
		for {
			for i := 0; i < s.length; i++ {
				select {
				case <-s.stopChan:
					return
				default:
					s.mu.Lock()
					s.erase()
					fmt.Printf("\r%s ", s.next())
					s.mu.Unlock()

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()
}

// Stop will stop the indicator.
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.active {
		s.active = false
		s.erase()
		s.stopChan <- struct{}{}
	}
}

func (s *Spinner) current() string {
	if s.lastFrame < 0 {
		return ""
	}
	r := s.frames[s.lastFrame%s.length]
	return string(r)
}

func (s *Spinner) next() string {
	r := s.frames[s.pos%s.length]
	s.lastFrame = s.pos
	s.pos++
	return string(r)
}

func (s *Spinner) erase() {
	n := utf8.RuneCountInString(s.current()) + 1
	if n == 1 {
		return
	}
	for _, c := range []string{"\b", " ", "\b"} {
		for i := 0; i < n; i++ {
			fmt.Printf(c)
		}
	}
}
