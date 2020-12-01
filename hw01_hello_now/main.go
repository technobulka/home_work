package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const timeHost = "0.beevik-ntp.pool.ntp.org"

func main() {
	t, err := ntp.Time(timeHost)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"current time: %s\nexact time: %s\n",
		time.Now().Round(0),
		t.Round(0),
	)
}
