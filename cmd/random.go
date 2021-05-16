package main

import (
	"math/rand"
	"time"
)

func setupRandSeed() {
	rand.Seed(time.Now().UTC().UnixNano())
}
