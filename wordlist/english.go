package wordlist

import (
	"bufio"
	"os"
)

// This file contains the English Wordlist

// struct english to give our global variable
// English a more descriptive type.
type english struct {}

func (wl english) Language() string {
	return "English"
}

func (wl english) Words() ([WordlistSize]string, error) {
	wordlistFile, err := os.Open("english.txt")

	if (err != nil) {
		return [2048]string{}, wordlistError{Message: err.Error()}
	}

	reader := bufio.NewReader(wordlistFile)

	// Forward declare this fixed-size string array
	var words [2048]string

	for i := 0; i < 2048; i++ {
		word, readErr := reader.ReadString('\n')

		if (readErr != nil) {
			return [2048]string{}, wordlistError{Message: readErr.Error()}
		}

		// Read up to second to last character in word because it will be the newline
		words[i] = word[:len(word) - 2]
	}

	return words, nil
}

func (wl english) GetWordAt(index uint32) (string, error) {
	if (index < 0 || index >= WordlistSize) {
		return "", wordlistError{Message: "Index out of range."}
	}

	wordlistFile, err := os.Open("english.txt")

	if (err != nil) {
		return "", wordlistError{Message: err.Error()}
	}

	reader := bufio.NewReader(wordlistFile)

	// Skip over index number of lines
	for i := uint32(0); i < index; i++; {
		_, skippingReadErr := reader.ReadString('\n')

		if (skippingReadErr != nil) {
			return "", wordlistError{Message: skippingReadErr.Error()}
		}
	}

	// By now, we should be on the right line to read from
	word, finalReadErr := reader.ReadString('\n')

	if (finalReadErr != nil) {
		return "", wordlistError{Message: finalReadErr.Error()}
	}

	return word[:len(word) - 2], nil
}

func (wl english) FindWord(word string) int {
	words, err := wl.Words()[:]

	if (err != nil) {
		return -1
	}

	index := FindWordIn(words[:], word)
	return index
}

// Export single variable to allow users access to english struct.
var English english = english{}