package hw4

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"
)

// Slicer return slice of to most used words in file or string
func Slicer(s string, count int) ([]string, error) {
	var words = make(map[string]int)
	_, err := os.Stat(s)
	if err != nil {
		words = processString(s, words)
	} else {
		words, err = processFile(s, words)
		if err != nil {
			return []string{}, err
		}
	}
	return sortWords(count, words), nil
}

// processString clear string and add words to map
func processString(s string, words map[string]int) map[string]int {
	clearString := clearText(s)
	return calculateWordsFromString(clearString, words)
}

// clearText remove symbols and whitespaces
func clearText(s string) string {
	for _, splitChar := range []string{",", " -", "—", ":", "'", ".", "!", "?", "/", "\\", "|"} {
		s = strings.Replace(s, splitChar, " ", -1)
	}
	return strings.ToLower(s)
}

// calculateWordsFromString add every word from string to map
func calculateWordsFromString(s string, words map[string]int) map[string]int {
	temp := strings.Fields(s)
	for _, val := range temp {
		if len(val) > 0 {
			words[val]++
		}
	}
	return words
}

// sortWords sort words from map by their counts and return slice
func sortWords(count int, words map[string]int) []string {
	type keyValue struct {
		Key   string
		Value int
	}
	wordsLen := len(words)
	var sortedStruct = make([]keyValue, 0, wordsLen)
	for key, value := range words {
		sortedStruct = append(sortedStruct, keyValue{key, value})
	}

	sort.Slice(sortedStruct, func(i, j int) bool {
		return sortedStruct[i].Value > sortedStruct[j].Value
	})
	var result = make([]string, 0, count)
	for i := 0; (i < count) && (i < wordsLen); i++ {
		result = append(result, sortedStruct[i].Key)
	}

	return result
}

// processFile will process every string from file
func processFile(s string, words map[string]int) (map[string]int, error) {
	file, err := os.Open(s)
	if err != nil {
		return map[string]int{}, err
	}
	reader := bufio.NewReader(file)
	defer file.Close()

	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return map[string]int{}, err
			}
		}
		words = processString(line, words)
	}
	return words, nil
}
