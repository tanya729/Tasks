package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Client struct {
	destination string
	timeout     time.Duration
	connection  net.Conn
	in          io.Reader
	out         io.Writer
}

func main() {
	var timeout string
	flag.StringVar(&timeout, "timeout", "90s", "timeout connection")
	flag.Parse()
	if flag.Arg(1) == "" {
		log.Fatalln("Need enter host and port")
	}
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatalln("Parse duration error:", err)
	}
	client := NewClient(net.JoinHostPort(flag.Arg(0), flag.Arg(1)), duration, os.Stdin, os.Stdout)
	if err := client.connect(); err != nil {
		fmt.Println("\nConnection closed:", err)
	}
}

// NewClient will create new Client
func NewClient(dest string, timeout time.Duration, in io.Reader, out io.Writer) *Client {
	c := &Client{
		destination: dest,
		timeout:     timeout,
		in:          in,
		out:         out,
	}
	return c
}

//connect
func (c *Client) connect() error {
	var err error
	c.connection, err = net.DialTimeout("tcp", c.destination, c.timeout)
	if err != nil {
		return err
	}
	defer c.connection.Close()
	fmt.Printf("Connected to %s with timeout %.0f seconds\n", c.destination, c.timeout.Seconds())
	fmt.Println("Press Ctrl-D / Ctrl-C to exit")
	errorChan := make(chan error)
	inputChan := make(chan string)
	outputChan := make(chan string)

	// Run the support goroutine to handle system interrupt with Ctrl-C
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
		<-signals
		errorChan <- errors.New("aborted by system interrupt")
	}()

	go forwarder(c.connection, inputChan, errorChan)
	go forwarder(c.in, outputChan, errorChan)
	for {
		select {
		case in := <-inputChan:
			fmt.Print(in)
		case out := <-outputChan:
			c.connection.Write([]byte(out))
		case err := <-errorChan:
			return err
		}
	}
}

//forwarder will forward input data to channels
func forwarder(in io.Reader, channel chan string, errChan chan error) {
	r := bufio.NewReader(in)
	for {
		str, err := r.ReadString('\n')
		if err != nil {
			errChan <- err
		}
		channel <- str
	}
}
