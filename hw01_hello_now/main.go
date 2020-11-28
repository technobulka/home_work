package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"current time: %s\nexact time: %s\n",
		time.Now().Round(0),
		t.Round(0),
	)
}
