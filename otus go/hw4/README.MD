# String slicer

## Usage
`$ go get github.com/lenniDespero/otus-golang/hw4`

    package main

	import (
    	"fmt"
    
    	"github.com/lenniDespero/otus-golang/hw4"
    )

	func main() {
		words, err := hw4.Slicer(`один два, три четырЕ - пять. шЕсть, семь вОсемь девять дЕсять.
                      Один два, три чЕтыре - пять. шесть, сЕмь восЕмь девять,
                      один два, три четырЕ - пять. шесть, семь вОсЕмь -
                      Один два, три чЕтыре - пять. шесть, сЕмь
                      один два, три четыре - пять. шесть,
                      Один два, три четырЕ - пять.
                      один два, три чЕтыре
                      Один два, три  -
                      один два,
                      Один`, 10)
            //Or you can read text from file        
           	//words := hw4.Slicer("./test.txt", 10)
        	if err != nil {
        		print(err)
        	}
        	for _, word := range words {
        		print(word)
        		print("\n")
        	}
	}

Output will be like: 

    [один два три четыре пять шесть семь восемь девять десять]

For test run `$ go test [path_to_module]`