// Copyright 2015 David Tweet.  All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license which can be found in the LICENSE file.

package reyong

import (
	"fmt"
	"math/rand"
)

const (
	PATTERN_LEN = 32
	REST        = '.'
)

type Role struct {
	name    string
	note_a  rune
	note_b  rune
	pattern []rune
}

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

func (r Role) NoStartingWithARest(i int) []rune {
	bad := make([]rune, 0)
	if i == 0 {
		bad = append(bad, REST)
	}
	return bad
}

func (r Role) StartWithARest(i int) []rune {
	bad := make([]rune, 0)
	if i == 0 {
		bad = append(bad, r.note_a, r.note_b)
	}
	return bad
}

func (r Role) NoRepeats(i int) []rune {
	bad := make([]rune, 0)
	if i > 0 {
		bad = append(bad, r.pattern[i-1])
	}
	if i == PATTERN_LEN-1 {
		bad = append(bad, r.pattern[0])
	}
	return bad
}

func (r Role) NoMoreThanThreeNotesWithoutARest(i int) []rune {
	bad := make([]rune, 0)
	if i >= 3 &&
		r.pattern[i-3] != REST &&
		r.pattern[i-2] != REST &&
		r.pattern[i-1] != REST {
		bad = append(bad, r.note_a, r.note_b)
	}
	if i == PATTERN_LEN-1 {
		for n := i - 2; n <= i; n++ {
			if (r.pattern[n] != REST || n == i) &&
				(r.pattern[(n+1)%PATTERN_LEN] != REST || (n+1)%PATTERN_LEN == i) &&
				(r.pattern[(n+2)%PATTERN_LEN] != REST || (n+2)%PATTERN_LEN == i) &&
				(r.pattern[(n+3)%PATTERN_LEN] != REST || (n+3)%PATTERN_LEN == i) {
				bad = append(bad, r.note_a, r.note_b)
			}
		}
	}
	return bad
}

func (r Role) NoRepeatingSingleNoteAndRestPairs(i int) []rune {
	bad := make([]rune, 0)
	if (i == 3 && r.pattern[i-2] == REST) ||
		(i > 3 && r.pattern[i-2] == REST && r.pattern[i-4] == REST) {
		bad = append(bad, REST)
	}
	if (i == PATTERN_LEN-1 || i == PATTERN_LEN-2) &&
		r.pattern[(i+2)%PATTERN_LEN] == REST &&
		(r.pattern[i-2] == REST || r.pattern[(i+4)%PATTERN_LEN] == REST) {
		bad = append(bad, REST)
	}
	return bad
}

func (r Role) NoSharedRests(i int, other_role Role) []rune {
	bad := make([]rune, 0)
	if other_role.pattern[i] == REST {
		bad = append(bad, REST)
	}
	return bad
}

func (r Role) HarmonizePolosAndSangsih(i int, other_role Role) []rune {
	bad := make([]rune, 0)
	if other_role.pattern[i] == other_role.note_b {
		bad = append(bad, r.note_b)
	} else if other_role.pattern[i] == other_role.note_a {
		bad = append(bad, r.note_a)
	}
	return bad
}

func SetupPolos(p []rune) Role {
	polos := Role{"polos", '1', '2', make([]rune, PATTERN_LEN)}
	copy(polos.pattern, p)
	return polos
}

func SetupSangsih(p []rune) Role {
	sangsih := Role{"sangsih", '3', '4', make([]rune, PATTERN_LEN)}
	copy(sangsih.pattern, p)
	return sangsih
}

func GeneratePolosAndSangsih() ([]rune, []rune) {
	polos := SetupPolos([]rune{})
	sangsih := SetupSangsih([]rune{})
	polos.FillPattern(sangsih)
	sangsih.FillPattern(polos)
	return polos.pattern, sangsih.pattern
}

func (r Role) FillPattern(other_role Role) {
	unworkable_subpatterns := &UnworkableSubpatterns{}
	backtracking := false
	for i := 0; i < PATTERN_LEN; i++ {
		bad := unworkable_subpatterns.BadEndingsFor(r.pattern[:i])
		if backtracking == true {
			fmt.Printf("Revising %v index %v. ", r.name, i)
			fmt.Printf("Bad is [%v].\n", string(bad))
			fmt.Printf("Pattern is [%v].\n", string(r.pattern))
			fmt.Printf("Other patt:[%v].\n", string(other_role.pattern))
			fmt.Printf("Unworkable subpatterns are:\n")
			for _, patt := range *unworkable_subpatterns {
				fmt.Println(string(patt))
			}
		}
		options := PatternOptions{
			r.note_a: true,
			r.note_b: true,
			REST:     true,
		}
		bad = append(bad, r.NoRepeats(i)...)
		bad = append(bad, r.NoMoreThanThreeNotesWithoutARest(i)...)
		bad = append(bad, r.NoRepeatingSingleNoteAndRestPairs(i)...)
		if r.name == "polos" {
			bad = append(bad, r.NoStartingWithARest(i)...)
		} else if r.name == "sangsih" {
			bad = append(bad, r.StartWithARest(i)...)
			bad = append(bad, r.HarmonizePolosAndSangsih(i, other_role)...)
			bad = append(bad, r.NoSharedRests(i, other_role)...)
		}
		for _, option := range bad {
			options[option] = false
		}
		option, err := options.Random()
		if err != nil {
			if i == 0 {
				panic("No valid options at zeroth index of pattern.")
			}
			fmt.Printf("Had no options at index %v\n", i)
			unworkable_subpatterns.Add(r.pattern[:i])
			backtracking = true
			i = i - 2
		} else {
			r.pattern[i] = option
			backtracking = false
		}
	}
}
