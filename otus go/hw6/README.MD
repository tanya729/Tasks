# String slicer

## Usage
`$ go get github.com/lenniDespero/otus-golang/hw6`

    package main

	import (
    	"fmt"
    
    	"github.com/lenniDespero/otus-golang/hw5"
    )

	func main() {
			list := new(hw6.List)
        	list.PushBack(1)
        	second := list.PushBack(2)
        	list.PushFront(0)
        	list.PushBack(3)
        	list.Remove(second)
        	checker := list.PushBack("checker data")
        	list.Remove(checker)
            for i := list.First(); i != nil; i = i.Next() {
            	fmt.Println(i.Value())
            }
            fmt.Println("Количество элементов:", list.Len())
	}

Output will be like: 

    0
    1
    3
    Количество элементов: 3


For test run `$ go test [path_to_module]`  

Test coverage info: 

    $ go tool cover -func=coverage.out
    github.com/lenniDespero/otus-golang/hw6/DLinkedList.go:16:      Len             100.0%
    github.com/lenniDespero/otus-golang/hw6/DLinkedList.go:21:      First           100.0%
    github.com/lenniDespero/otus-golang/hw6/DLinkedList.go:26:      Last            100.0%
    github.com/lenniDespero/otus-golang/hw6/DLinkedList.go:31:      Value           100.0%
    github.com/lenniDespero/otus-golang/hw6/DLinkedList.go:36:      Next            100.0%
    github.com/lenniDespero/otus-golang/hw6/DLinkedList.go:41:      Prev            100.0%
    github.com/lenniDespero/otus-golang/hw6/DLinkedList.go:46:      PushFront       100.0%
    github.com/lenniDespero/otus-golang/hw6/DLinkedList.go:60:      PushBack        100.0%
    github.com/lenniDespero/otus-golang/hw6/DLinkedList.go:74:      Remove          100.0%
    total:                                                          (statements)    100.0%

