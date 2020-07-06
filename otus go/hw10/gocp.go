package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	from := flag.String("from", "", "limit bytes to write")
	to := flag.String("to", "", "offset bytes to write")
	limit := flag.Int64("limit", -1, "limit bytes to write")
	offset := flag.Int64("offset", 0, "offset bytes to write")
	flag.Parse()
	if *from == "" || *to == "" {
		fmt.Println("Nothing to copy. Use -h to help.")
		os.Exit(1)
	}
	err := Copy(*from, *to, *limit, *offset)
	if err != nil {
		fmt.Printf("An error occurred while file copying: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("File copied successfully")
}

//Copy will copy From (string) file destination To (string) file destination Limit (int) bytes with Offset (int) bytes
func Copy(from string, to string, limit int64, offset int64) error {
	file, err := os.Open(from)
	if err != nil {

		return fmt.Errorf("no file or access denied: %s", err.Error())
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("can't read file information: %s", err.Error())
	}
	total := info.Size()

	input := io.Reader(file)
	if limit != -1 {
		input = io.LimitReader(file, limit)
		total = limit
	}

	if offset > 0 {
		pos, err := file.Seek(offset, 0)
		if err != nil || pos != offset {
			return fmt.Errorf("can't set offset for reading")
		}
		if offset+total >= info.Size() {
			if total == info.Size() {
				total = info.Size() - offset
			} else {
				total = offset + total - info.Size()
			}
		}
	}

	output, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("can't create file: %s", err.Error())
	}
	defer output.Close()
	defer println()
	for totalWritten := int64(0); totalWritten < total; {
		written, err := io.CopyN(output, input, 1024)
		totalWritten += written
		if err != nil {
			if err == io.EOF {
				printBar(totalWritten, totalWritten)
				break
			}
			return fmt.Errorf("error: %s", err.Error())
		}
		printBar(totalWritten, total)
	}

	return nil
}

func printBar(written int64, total int64) {
	fmt.Printf(
		"%%\rProcessing... [%d / %d] %.2f",
		written,
		total,
		float64(written)/float64(total)*100,
	)
}
