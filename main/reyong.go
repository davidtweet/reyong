package main

import (
	"fmt"
	"github.com/davidtweet/reyong"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	polos_pattern := reyong.GeneratePolos()
	fmt.Printf("polos:    [%v]\n", string(polos_pattern))
}
