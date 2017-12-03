package wordlist

import (
	"bufio"
	"os"
	"path"
	"strings"
)

// This file contains the English Wordlist

// struct english to give our global variable
// English a more descriptive type.
type english struct {}

func (wl english) Language() string {
	return "English"
}

func (wl english) Words() ([WordlistSize]string, error) {
	directory, directoryErr := getCurrentDirectory()

	if (directoryErr != nil) {
		return [WordlistSize]string{}, wordlistError{Message: directoryErr.Error()}
	}

	wordlistFile, err := os.Open(path.Join(directory, "english.txt"))

	if (err != nil) {
		return [WordlistSize]string{}, wordlistError{Message: err.Error()}
	}

	reader := bufio.NewReader(wordlistFile)

	// Forward declare this fixed-size string array
	var words [WordlistSize]string

	for i := 0; i < WordlistSize; i++ {
		word, readErr := reader.ReadString('\n')

		if (readErr != nil) {
			return [WordlistSize]string{}, wordlistError{Message: readErr.Error()}
		}

		// Read up to second to last character in word because it will be the newline
		// and then we must remove a possible \r because Windows
		words[i] = strings.Replace(word[:len(word) - 1], "\r", "", -1)
	}

	return words, nil
}

func (wl english) GetWordAt(index uint32) (string, error) {
	if (index < 0 || index >= WordlistSize) {
		return "", wordlistError{Message: "Index out of range."}
	}

	directory, directoryErr := getCurrentDirectory()

	if (directoryErr != nil) {
		return "", wordlistError{Message: directoryErr.Error()}
	}

	wordlistFile, err := os.Open(path.Join(directory, "english.txt"))

	if (err != nil) {
		return "", wordlistError{Message: err.Error()}
	}

	reader := bufio.NewReader(wordlistFile)

	// Skip over index number of lines
	for i := uint32(0); i < index; i++ {
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

	// Trim possible \r from word
	return strings.Replace(word[:len(word) - 1], "\r", "", -1), nil
}

func (wl english) FindWord(word string) int {
	words, err := wl.Words()

	if (err != nil) {
		return -1
	}

	index := FindWordIn(words[:], word)
	return index
}

// Export single variable to allow users access to english struct.
var English english = english{}