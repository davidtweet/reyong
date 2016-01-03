package reyong_test

import (
	"github.com/davidtweet/reyong"
	"testing"
)

func newpattern(p []rune) []rune {
	pattern := make([]rune, reyong.PATTERN_LEN)
	for i, r := range p {
		pattern[i] = r
	}
	return pattern
}

func TestNoStartingWithARest(t *testing.T) {
	bad := reyong.NoStartingWithARest(newpattern([]rune{}), 0)
	if !(bad[0] == reyong.REST) {
		t.Fail()
	}
}

func TestNoRepeats(t *testing.T) {
	bad := reyong.NoRepeats(newpattern([]rune("12")), 2)
	if !(bad[0] == '2') {
		t.Fail()
	}
}

func TestNoMoreThanThreeNotesWithoutARest(t *testing.T) {
	bad1 := reyong.NoMoreThanThreeNotesWithoutARest(
		newpattern([]rune("1.212")), 5)
	if !(bad1[0] == '1' && bad1[1] == '2') {
		t.Fail()
	}
	bad2 := reyong.NoMoreThanThreeNotesWithoutARest(
		newpattern([]rune("21.21.21.21.212.12.12.12.121.1")), 31)
	if !(bad2[0] == '1' && bad2[1] == '2') {
		t.Fail()
	}
	bad3 := reyong.NoMoreThanThreeNotesWithoutARest(
		newpattern([]rune("2.121.21.21.212.12.12.12.12.21")), 31)
	if !(bad3[0] == '1' && bad3[1] == '2') {
		t.Fail()
	}
	bad4 := reyong.NoMoreThanThreeNotesWithoutARest(
		newpattern([]rune("212.1.21.21.212.12.12.12.1.12.")), 31)
	if !(bad4[0] == '1' && bad4[1] == '2') {
		t.Fail()
	}
}

func TestNoRepeatingSingleNoteAndRestPairs(t *testing.T) {
	bad1 := reyong.NoRepeatingSingleNoteAndRestPairs(
		newpattern([]rune("1.2")), 3)
	if !(bad1[0] == '.') {
		t.Fail()
	}
	bad2 := reyong.NoRepeatingSingleNoteAndRestPairs(
		newpattern([]rune("1.212.1.2")), 9)
	if !(bad2[0] == '.') {
		t.Fail()
	}
	bad3 := reyong.NoRepeatingSingleNoteAndRestPairs(
		newpattern([]rune("2.121.21.212.212.12.12.12.121.1")), 31)
	if !(bad3[0] == '.') {
		t.Fail()
	}
}

func TestUnworkableSubpatterns(t *testing.T) {
	usps := &reyong.UnworkableSubpatterns{}
	pattern := []rune("21.1.12.12.2.12.1.212.2.212.1.1")
	usps.Add(pattern)
	bad1 := usps.BadEndingsFor([]rune("21.1.12.12.2.12.1.212.2.212.1."))
	if !(bad1[0] == '1') {
		t.Fail()
	}
	pattern[30] = '2'
	usps.Add(pattern)
	bad2 := usps.BadEndingsFor([]rune("21.1.12.12.2.12.1.212.2.212.1."))
	if !(bad2[0] == '1' && bad2[1] == '2') {
		t.Fail()
	}
}
