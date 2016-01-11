package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/davidtweet/reyong"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	verbose_logging_default = false
	pattern_length_default  = 32
	polos_pattern_default   = ""
	sangsih_pattern_default = ""
)

var verbose_logging bool
var pattern_length int
var polos_pattern string
var sangsih_pattern string

func main() {
	flag.BoolVar(&verbose_logging, "v", verbose_logging_default, "Log verbosely")
	flag.IntVar(&pattern_length, "pattern_length", pattern_length_default, "Length of pattern (multiple of 16)")
	flag.StringVar(&polos_pattern, "polos_pattern", polos_pattern_default, "Initial polos pattern (compose sangsih only)")
	flag.StringVar(&sangsih_pattern, "sangsih_pattern", sangsih_pattern_default, "Initial sangsih pattern (compose polos only)")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(os.Stderr)
	if verbose_logging {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if pattern_length <= 0 || pattern_length%16 != 0 {
		fmt.Println("Pattern length must be a multiple of 16.")
		flag.Usage()
		os.Exit(1)
	}
	if len(polos_pattern) > pattern_length || len(sangsih_pattern) > pattern_length {
		fmt.Println("Initial pattern is longer than pattern length. Change pattern or pattern_length.")
		flag.Usage()
		os.Exit(1)
	}
	for _, r := range polos_pattern {
		if !strings.Contains("12.", string(r)) {
			fmt.Println("Polos pattern contains invalid option %v", r)
		}
	}
	for _, r := range sangsih_pattern {
		if !strings.Contains("34.", string(r)) {
			fmt.Println("Sangsih pattern contains invalid option %v", r)
		}
	}

	rand.Seed(time.Now().UnixNano())
	polos_pattern, sangsih_pattern := reyong.GeneratePolosAndSangsih([]rune(polos_pattern), []rune(sangsih_pattern), pattern_length)
	fmt.Printf("polos:    [%v]\n", string(polos_pattern))
	fmt.Printf("sangsih:  [%v]\n", string(sangsih_pattern))
}
