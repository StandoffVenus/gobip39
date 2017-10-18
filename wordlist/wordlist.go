package wordlist

const (
	WordlistSize = 2048
)

// Error type specifically for wordlist errors
type wordlistError struct {
	Message string
}

func (err wordlistError) Error() string {
	return err.Message
}

// A Wordlist must implement
// 	Language - returns the language that that Wordlist's words are in as a string.
// 	Words - returns all the words from that Wordlist as a string array of size 2048.
// 	GetWordAt - returns word at specified uint32 index in Wordlist's words and potential errors.
// 	  * GetWordAt will return an error if the index is greater than 2047 or less than 0.
// 	  * When GetWordAt returns an error, it return an empty string ("")
// 	  * GetWordAt uses uint32 as the index type because of Mnemonic's use of BitReader.Read32.
// 	    Read32 returns type uint32. This type works for Wordlist's required methods, and although,
// 	    makes things a little more inconvenient, will remain as such.
// 	FindWord - returns index of specified word in Wordlist's words as int. If not found, it returns -1.
//
// Wordlists must follow BIP-0039 specification: https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki#wordlist
//
// Look at English.go for an example.
type Wordlist interface {
	Language() string
	Words() [WordlistSize]string
	GetWordAt(uint32) (string, error)
	FindWord(string) int
}

// Wordlist binary search.
// Takes array of strings to search for passed string in.
// Returns int of b's position in a, -1 if not found.
func FindWordIn(a []string, b string) int {
	low := 0
	high := len(a) - 1
	var index int

	for {
		if (high < low) { break }

		index = (low + high) / 2

		if (a[index] < b) {
			low = index + 1
		} else if (a[index] > b) {
			high = index - 1
		} else {
			return index
		}
	}

	return -1
}

