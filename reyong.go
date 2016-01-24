// Copyright 2015 David Tweet.  All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license which can be found in the LICENSE file.

package reyong

import (
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"strings"
)

type Role struct {
	name    string
	note_a  rune
	note_b  rune
	rest    rune
	pattern []rune
}

type UnworkableSubpatterns [][]rune

type NoValidOptionsError struct {
	s string
}

func (e NoValidOptionsError) Error() string {
	return e.s
}

func (usps *UnworkableSubpatterns) BadEndingsFor(subpattern []rune) []rune {
	subp_length := len(subpattern)
	bad := make([]rune, 0)
	for _, usp := range *usps {
		if subp_length == len(usp)-1 && string(subpattern) == string(usp[:subp_length]) {
			bad = append(bad, usp[subp_length])
		}
	}
	return bad
}

func (usps *UnworkableSubpatterns) Add(pattern []rune) {
	log.WithFields(log.Fields{
		"subpattern": string(pattern),
		"length":     len(pattern),
	}).Info("Adding unworkable subpattern")
	new_pattern := make([]rune, len(pattern))
	copy(new_pattern, pattern)
	*usps = append(*usps, new_pattern)
}

func (r Role) NoStartingWithARest(i int) []rune {
	bad := make([]rune, 0)
	if i == 0 {
		bad = append(bad, r.rest)
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
	if i == len(r.pattern)-1 {
		bad = append(bad, r.pattern[0])
	}
	return bad
}

func (r Role) NoMoreThanThreeNotesWithoutARest(i int) []rune {
	pattern_len := len(r.pattern)
	bad := make([]rune, 0)
	if i >= 3 &&
		r.pattern[i-3] != r.rest &&
		r.pattern[i-2] != r.rest &&
		r.pattern[i-1] != r.rest {
		bad = append(bad, r.note_a, r.note_b)
	}
	if i == pattern_len-1 {
		for n := i - 2; n <= i; n++ {
			if (r.pattern[n] != r.rest || n == i) &&
				(r.pattern[(n+1)%pattern_len] != r.rest || (n+1)%pattern_len == i) &&
				(r.pattern[(n+2)%pattern_len] != r.rest || (n+2)%pattern_len == i) &&
				(r.pattern[(n+3)%pattern_len] != r.rest || (n+3)%pattern_len == i) {
				bad = append(bad, r.note_a, r.note_b)
			}
		}
	}
	return bad
}

func (r Role) NoRepeatingSingleNoteAndRestPairs(i int) []rune {
	pattern_len := len(r.pattern)
	bad := make([]rune, 0)
	if (i == 3 && r.pattern[i-2] == r.rest) ||
		(i > 3 && r.pattern[i-2] == r.rest && r.pattern[i-4] == r.rest) {
		bad = append(bad, r.rest)
	}
	if (i == pattern_len-1 || i == pattern_len-2) &&
		r.pattern[(i+2)%pattern_len] == r.rest &&
		(r.pattern[i-2] == r.rest || r.pattern[(i+4)%pattern_len] == r.rest) {
		bad = append(bad, r.rest)
	}
	return bad
}

func (r Role) NoSameNoteSeparatedByRestFollowedByRest(i int) []rune {
	// pattern_len := len(r.pattern)
	bad := make([]rune, 0)
	if i >= 3 && r.pattern[i-2] == r.rest && r.pattern[i-1] == r.pattern[i-3] {
		bad = append(bad, r.rest)
	}
	return bad
	//TODO: Add check for end-of-pattern.
}

func (r Role) NoSharedRests(i int, other_role Role) []rune {
	bad := make([]rune, 0)
	if other_role.pattern[i] == r.rest {
		bad = append(bad, r.rest)
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

func SetupPolos(p []rune, pattern_len int) Role {
	polos := Role{"polos", '1', '2', '.', make([]rune, pattern_len)}
	copy(polos.pattern, p)
	return polos
}

func SetupSangsih(p []rune, pattern_len int) Role {
	sangsih := Role{"sangsih", '3', '4', '.', make([]rune, pattern_len)}
	copy(sangsih.pattern, p)
	return sangsih
}

func GeneratePolosAndSangsih(initial_polos []rune, initial_sangsih []rune, pattern_len int) ([]rune, []rune) {
	polos := SetupPolos(initial_polos, pattern_len)
	sangsih := SetupSangsih(initial_sangsih, pattern_len)
	polos.FillPattern(sangsih)
	sangsih.FillPattern(polos)
	return polos.pattern, sangsih.pattern
}

func formatPatternInProgress(pattern []rune, i int) string {
	out := make([]rune, 0)
	for n, r := range pattern {
		if n == i {
			out = append(out, '(', r, ')')
		} else if r == 0 {
			out = append(out, 'X')
		} else {
			out = append(out, r)
		}
	}
	return string(out)
}

func formatBadOptionsList(bad_list []rune) string {
	out := make([]rune, 0)
	out = append(out, '[')
	for _, r := range bad_list {
		if len(out) > 1 {
			out = append(out, ',')
		}
		out = append(out, '\'', r, '\'')
	}
	out = append(out, ']')
	return string(out)
}

func (r Role) FillPattern(other_role Role) {
	unworkable_subpatterns := &UnworkableSubpatterns{}
	revising := false
	pattern_len := len(r.pattern)
	log.Info("Filling ", r.name, " pattern")
	for i := 0; i < pattern_len; i++ {
		bad := unworkable_subpatterns.BadEndingsFor(r.pattern[:i])
		if revising == true {
			log.WithFields(log.Fields{
				"index":         i,
				"bad_options":   formatBadOptionsList(bad),
				r.name:          formatPatternInProgress(r.pattern, i),
				other_role.name: formatPatternInProgress(other_role.pattern, i),
			}).Info("Revising pattern")
		}
		bad = append(bad, r.NoRepeats(i)...)
		bad = append(bad, r.NoMoreThanThreeNotesWithoutARest(i)...)
		bad = append(bad, r.NoRepeatingSingleNoteAndRestPairs(i)...)
		bad = append(bad, r.NoSameNoteSeparatedByRestFollowedByRest(i)...)
		bad = append(bad, r.HarmonizePolosAndSangsih(i, other_role)...)
		bad = append(bad, r.NoSharedRests(i, other_role)...)
		if r.name == "polos" {
			bad = append(bad, r.NoStartingWithARest(i)...)
		} else if r.name == "sangsih" {
			bad = append(bad, r.StartWithARest(i)...)
		}
		all_options := []rune{r.note_a, r.note_b, r.rest}
		valid_options := make([]rune, 0)
		for _, option := range all_options {
			if !strings.ContainsRune(string(bad), option) {
				valid_options = append(valid_options, option)
			}
		}
		if len(valid_options) == 0 {
			if i == 0 {
				panic("No valid options at zeroth index of pattern.")
			}
			unworkable_subpatterns.Add(r.pattern[:i])
			revising = true
			i = i - 2
		} else {
			if !strings.ContainsRune(string(valid_options), r.pattern[i]) {
				r.pattern[i] = valid_options[rand.Intn(len(valid_options))]
			}
			revising = false
		}
	}
}
