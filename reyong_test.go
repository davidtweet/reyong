package reyong_test

import (
	"github.com/davidtweet/reyong"
	"testing"
)

func TestNoStartingWithARest(t *testing.T) {
	polos := reyong.SetupPolos([]rune{})
	bad := polos.NoStartingWithARest(0)
	if !(bad[0] == reyong.REST) {
		t.Fail()
	}
}

func TestNoRepeats(t *testing.T) {
	polos := reyong.SetupPolos([]rune("12"))
	bad := polos.NoRepeats(2)
	if !(bad[0] == '2') {
		t.Fail()
	}
}

func TestNoMoreThanThreeNotesWithoutARest(t *testing.T) {
	polos1 := reyong.SetupPolos([]rune("1.212"))
	bad1 := polos1.NoMoreThanThreeNotesWithoutARest(5)
	if !(bad1[0] == '1' && bad1[1] == '2') {
		t.Fail()
	}
	polos2 := reyong.SetupPolos([]rune("21.21.21.21.212.12.12.12.121.1"))
	bad2 := polos2.NoMoreThanThreeNotesWithoutARest(31)
	if !(bad2[0] == '1' && bad2[1] == '2') {
		t.Fail()
	}
	polos3 := reyong.SetupPolos([]rune("2.121.21.21.212.12.12.12.12.21"))
	bad3 := polos3.NoMoreThanThreeNotesWithoutARest(31)
	if !(bad3[0] == '1' && bad3[1] == '2') {
		t.Fail()
	}
	polos4 := reyong.SetupPolos([]rune("212.1.21.21.212.12.12.12.1.12."))
	bad4 := polos4.NoMoreThanThreeNotesWithoutARest(31)
	if !(bad4[0] == '1' && bad4[1] == '2') {
		t.Fail()
	}
}

func TestNoRepeatingSingleNoteAndRestPairs(t *testing.T) {
	polos1 := reyong.SetupPolos([]rune("1.2"))
	bad1 := polos1.NoRepeatingSingleNoteAndRestPairs(3)
	if !(bad1[0] == '.') {
		t.Fail()
	}
	polos2 := reyong.SetupPolos([]rune("1.212.1.2"))
	bad2 := polos2.NoRepeatingSingleNoteAndRestPairs(9)
	if !(bad2[0] == '.') {
		t.Fail()
	}
	polos3 := reyong.SetupPolos([]rune("2.121.21.212.212.12.12.12.121.1"))
	bad3 := polos3.NoRepeatingSingleNoteAndRestPairs(31)
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
