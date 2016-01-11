package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/davidtweet/reyong"
	"math/rand"
	"os"
	"time"
)

const (
	verbose_logging_default = false
	pattern_length_default  = 32
)

var verbose_logging bool
var pattern_length int

func main() {
	flag.BoolVar(&verbose_logging, "v", verbose_logging_default, "Log verbosely")
	flag.IntVar(&pattern_length, "pattern_length", pattern_length_default, "Length of pattern (multiple of 16)")
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

	rand.Seed(time.Now().UnixNano())
	polos_pattern, sangsih_pattern := reyong.GeneratePolosAndSangsih(pattern_length)
	fmt.Printf("polos:    [%v]\n", string(polos_pattern))
	fmt.Printf("sangsih:  [%v]\n", string(sangsih_pattern))
}
