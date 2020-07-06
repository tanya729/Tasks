package hw1

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

const server = "0.pool.ntp.org"

// CurrentTime returns current time or error from ntp server
func CurrentTime() {
	time, err := ntp.Time(server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Current time is: %s\n", time.Local())
}
