# Local NTP time

## Usage
`$ go get github.com/lenniDespero/otus-golang/hw1`

    package main

	import "github.com/lenniDespero/otus-golang/hw1"

	func main() {
		hw1.CurrentTime()
	}
Output will be like "Current time is: {Your local time}" or "Error: {your error message}"