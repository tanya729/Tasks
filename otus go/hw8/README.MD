# String slicer

## Usage
`$ go get github.com/lenniDespero/otus-golang/hw8`

    package main

	import (
    	"fmt"
    
    	"github.com/lenniDespero/otus-golang/hw8"
    )

	func main() {
			f1 := func() error {
            		time.Sleep(1 * time.Second)
            		return errors.New("Error text f1")
            	}
            
            	f2 := func() error {
            		time.Sleep(4 * time.Second)
            		return errors.New("Error text f2")
            	}
            
            	f3 := func() error {
            		time.Sleep(1 * time.Second)
            		return nil
            	}
            
            	f4 := func() error {
            		time.Sleep(3 * time.Second)
            		return errors.New("Error text f4")
            	}
            
            	f5 := func() error {
            		time.Sleep(10 * time.Second)
            		return errors.New("Error text f5")
            	}
            
            	_, counter := hw8.Run([]hw8.WorkerFunction{
            		f1,
            		f2,
            		f3,
            		f4,
            		f5,
            	}, 2, 1)
            	fmt.Printf("Run %d tasks\n", counter)
	}

Output will be like: 

    Run 1 tasks


For test run `$ go test [path_to_module]`  

