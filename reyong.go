// Copyright 2015 David Tweet.  All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license which can be found in the LICENSE file.

package reyong

import (
	"fmt"
	"math/rand"
)

const (
	PATTERN_LEN    = 32
	REST           = '.'
	POLOS_NOTE_A   = '1'
	POLOS_NOTE_B   = '2'
	SANGSIH_NOTE_A = '3'
	SANGSIH_NOTE_B = '4'
)

//type Instrument struct {
//	name   string
//	note_a rune
//	note_b rune
//}

type UnworkableSubpatterns [][]rune

type PatternOptions map[rune]bool

type NoValidOptionsError struct {
	s string
}

func (e NoValidOptionsError) Error() string {
	return e.s
}

func SameElements(s1 []rune, s2 []rune) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, _ := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func (usps *UnworkableSubpatterns) BadEndingsFor(subpattern []rune) []rune {
	subp_length := len(subpattern)
	bad := make([]rune, 0)
	for _, usp := range *usps {
		if subp_length == len(usp)-1 && SameElements(subpattern, usp[:subp_length]) {
			bad = append(bad, usp[subp_length])
		}
	}
	return bad
}

func (usps *UnworkableSubpatterns) Add(pattern []rune) {
	new_pattern := make([]rune, len(pattern))
	copy(new_pattern, pattern)
	*usps = append(*usps, new_pattern)
}

func (po PatternOptions) Random() (r rune, e error) {
	candidate_list := make([]rune, 0)
	for candidate, included := range po {
		if included == true {
			candidate_list = append(candidate_list, candidate)
		}
	}
	if len(candidate_list) > 0 {
		return candidate_list[rand.Intn(len(candidate_list))], nil
	} else {
		return 0, NoValidOptionsError{"No valid options."}
	}
}

func NoStartingWithARest(pattern []rune, i int) []rune {
	bad := make([]rune, 0)
	if i == 0 {
		bad = append(bad, REST)
	}
	return bad
}

func NoRepeats(pattern []rune, i int) []rune {
	bad := make([]rune, 0)
	if i > 0 {
		bad = append(bad, pattern[i-1])
	}
	if i == PATTERN_LEN-1 {
		bad = append(bad, pattern[0])
	}
	return bad
}

func NoMoreThanThreeNotesWithoutARest(pattern []rune, i int) []rune {
	bad := make([]rune, 0)
	if i >= 3 &&
		pattern[i-3] != REST &&
		pattern[i-2] != REST &&
		pattern[i-1] != REST {
		bad = append(bad, POLOS_NOTE_A, POLOS_NOTE_B)
	}
	if i == PATTERN_LEN-1 {
		for n := i - 2; n <= i; n++ {
			if (pattern[n] != REST || n == i) &&
				(pattern[(n+1)%PATTERN_LEN] != REST || (n+1)%PATTERN_LEN == i) &&
				(pattern[(n+2)%PATTERN_LEN] != REST || (n+2)%PATTERN_LEN == i) &&
				(pattern[(n+3)%PATTERN_LEN] != REST || (n+3)%PATTERN_LEN == i) {
				bad = append(bad, POLOS_NOTE_A, POLOS_NOTE_B)
			}
		}
	}
	return bad
}

func NoRepeatingSingleNoteAndRestPairs(pattern []rune, i int) []rune {
	bad := make([]rune, 0)
	if (i == 3 && pattern[i-2] == REST) ||
		(i > 3 && pattern[i-2] == REST && pattern[i-4] == REST) {
		bad = append(bad, REST)
	}
	if i == PATTERN_LEN-1 &&
		pattern[i-2] == REST &&
		pattern[(i+2)%PATTERN_LEN] == REST {
		bad = append(bad, REST)
	}
	return bad
}

func GeneratePolos() []rune {
	pattern := make([]rune, PATTERN_LEN)
	unworkable_subpatterns := &UnworkableSubpatterns{}
	backtracking := false
	for i := 0; i < PATTERN_LEN; i++ {
		bad := unworkable_subpatterns.BadEndingsFor(pattern[:i])
		if backtracking == true {
			fmt.Printf("Revising polos index %v. ", i)
			fmt.Printf("Bad is [%v]. ", string(bad))
			fmt.Printf("Pattern is [%v].\n", string(pattern))
			fmt.Printf("Unworkable subpatterns are:\n")
			for _, patt := range *unworkable_subpatterns {
				fmt.Println(string(patt))
			}
		}
		polos_options := PatternOptions{
			POLOS_NOTE_A: true,
			POLOS_NOTE_B: true,
			REST:         true,
		}
		bad = append(bad, NoStartingWithARest(pattern, i)...)
		bad = append(bad, NoRepeats(pattern, i)...)
		bad = append(bad, NoMoreThanThreeNotesWithoutARest(pattern, i)...)
		bad = append(bad, NoRepeatingSingleNoteAndRestPairs(pattern, i)...)
		for _, option := range bad {
			polos_options[option] = false
		}
		option, err := polos_options.Random()
		if err != nil {
			if i == 0 {
				panic("No valid options at zeroth index of pattern.")
			}
			fmt.Printf("Had no options at index %v\n", i)
			unworkable_subpatterns.Add(pattern[:i])
			backtracking = true
			i = i - 2
		} else {
			pattern[i] = option
			backtracking = false
		}
	}
	return pattern
}
