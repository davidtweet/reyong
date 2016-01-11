package main

import (
	"fmt"
	"github.com/davidtweet/reyong"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	polos_pattern, sangsih_pattern := reyong.GeneratePolosAndSangsih()
	fmt.Printf("polos:    [%v]\n", string(polos_pattern))
	fmt.Printf("sangsih:  [%v]\n", string(sangsih_pattern))
}
