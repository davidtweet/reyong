package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	PATTERN_LENGTH = 32
	REST           = '.'
	POLOS_NOTE_A   = '2'
	POLOS_NOTE_B   = '3'
	SANGSIH_NOTE_A = '5'
	SANGSIH_NOTE_B = '6'
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
	disqualified := []rune{REST}
	counting := false
	for i := 0; i < PATTERN_LENGTH; i++ {
		if len(disqualified) > 0 && i != 0 {
			counting = true
		}
		if counting == true {
			fmt.Printf("Revising polos index %v. Disqualified is [%v]. Pattern is [%v].\n", i, string(disqualified), string(pattern))
		}
		polos_options := patternOptions{
			POLOS_NOTE_A: true,
			POLOS_NOTE_B: true,
			REST:         true,
		}
		if i >= 3 &&
			pattern[i-3] != REST &&
			pattern[i-2] != REST &&
			pattern[i-1] != REST {
			disqualified = append(disqualified, POLOS_NOTE_A, POLOS_NOTE_B)
		}
		if i > 0 {
			disqualified = append(disqualified, pattern[i-1])
		}
		if i == PATTERN_LENGTH-1 {
			if pattern[0] != REST &&
				pattern[1] != REST &&
				(pattern[2] != REST || pattern[i-1] != REST) {
				disqualified = append(disqualified, POLOS_NOTE_A, POLOS_NOTE_B)
			}
			disqualified = append(disqualified, pattern[0])
		}
		for _, option := range disqualified {
			polos_options[option] = false
		}
		option, err := polos_options.Random()
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
	fmt.Println(string(polos_pattern))
}
