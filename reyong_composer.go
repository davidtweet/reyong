// Copyright 2015 David Tweet.  All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license which can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	PATTERN_LENGTH = 32
	REST           = '.'
	POLOS_NOTE_A   = '1'
	POLOS_NOTE_B   = '2'
	SANGSIH_NOTE_A = '3'
	SANGSIH_NOTE_B = '4'
)

type patternOptions map[rune]bool

type noValidOptionsError struct {
	s string
}

func (e noValidOptionsError) Error() string {
	return e.s
}

func (po patternOptions) Random() (r rune, e error) {
	candidate_list := make([]rune, 0)
	for candidate, included := range po {
		if included == true {
			candidate_list = append(candidate_list, candidate)
		}
	}
	if len(candidate_list) > 0 {
		return candidate_list[rand.Intn(len(candidate_list))], nil
	} else {
		return 0, noValidOptionsError{"No valid options."}
	}
}

func GeneratePolos() []rune {
	pattern := make([]rune, PATTERN_LENGTH)
	// No starting a polos pattern with a rest.
	disqualified := []rune{REST}
	log_revisions := false
	for i := 0; i < PATTERN_LENGTH; i++ {
		if len(disqualified) > 0 && i != 0 {
			log_revisions = true
		}
		if log_revisions == true {
			fmt.Printf("Revising polos index %v. ", i)
			fmt.Printf("Disqualified is [%v]. ", string(disqualified))
			fmt.Printf("Pattern is [%v].\n", string(pattern))
		}
		polos_options := patternOptions{
			POLOS_NOTE_A: true,
			POLOS_NOTE_B: true,
			REST:         true,
		}
		// No more than three notes without a rest.
		if i >= 3 &&
			pattern[i-3] != REST &&
			pattern[i-2] != REST &&
			pattern[i-1] != REST {
			disqualified = append(disqualified, POLOS_NOTE_A, POLOS_NOTE_B)
		}
		// No repeating notes or rests.
		if i > 0 {
			disqualified = append(disqualified, pattern[i-1])
		}
		if i == PATTERN_LENGTH-1 {
			// No more than three notes without a rest, considering that
			// the pattern wraps around from the end to the beginning.
			if pattern[0] != REST &&
				pattern[1] != REST &&
				(pattern[2] != REST || pattern[i-1] != REST) {
				disqualified = append(disqualified, POLOS_NOTE_A, POLOS_NOTE_B)
			}
			// Don't let the first and last note of the pattern be the same.
			disqualified = append(disqualified, pattern[0])
		}
		for _, option := range disqualified {
			polos_options[option] = false
		}
		option, err := polos_options.Random()
		// If we paint ourselves into a corner and have no valid options,
		// go back and replace the last note with something different.
		if err != nil {
			if i == 0 {
				panic("No valid options at zeroth index of pattern.")
			}
			disqualified = []rune{pattern[i-1]}
			i = i - 2
		} else {
			pattern[i] = option
			disqualified = make([]rune, 0)
		}
	}
	return pattern
}

func main() {
	rand.Seed(time.Now().UnixNano())
	polos_pattern := GeneratePolos()
	fmt.Printf("polos:    [%v]\n", string(polos_pattern))
}
