# String unpacker

## Usage
`$ go get github.com/lenniDespero/otus-golang/hw3`

    package main

	import "github.com/lenniDespero/otus-golang/hw3"

	func main() {
		baseStrings := []string{`a4bc2d5e`, `abcd`, `45`, `qwe\4\5`, `qwe\45`, `qwe\\5`, `qwe\42`}
        	for _, string := range baseStrings {
        		result, err := hw3.Unpack(string)
        		if err != nil {
        			fmt.Printf("%s => %s (❌  not unpacked) : %s\n", string, result, err)
        		} else {
        			fmt.Printf("%s => %s (✅  unpacked)\n", string, result)
        		}
        
        	}
	}

Output will be like: 

    a4bc2d5e => aaaabccddddde (✅  unpacked)
    abcd => abcd (✅  unpacked)
    45 =>  (❌  not unpacked) : Something strange with unpacking string '45' on character - 4, position - 0.
    qwe\4\5 => qwe45 (✅  unpacked)
    qwe\45 => qwe44444 (✅  unpacked)
    qwe\\5 => qwe\\\\\ (✅  unpacked)
    qwe\42 => qwe44 (✅  unpacked)
