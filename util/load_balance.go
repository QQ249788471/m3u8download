package util

import (
	"math/rand"
	"time"
)

func LB() {
	n := rand.Intn(50) + 50

	// fmt.Println("LB", n)

	time.Sleep(time.Millisecond * time.Duration(n))
}
