package gobip39

// This file wraps mnemonic [sentence] specifications as detailed by
// BIP-0039 spec.

import (
	"bytes"
	"github.com/32bitkid/bitreader"
	"gobip39/wordlist"
)

const (
	WordBitLength = 11
	MinimumSentenceSize = (MinimumEntropySize + MinimumChecksumSize) / WordBitLength
	MaximumSentenceSize = (MaximumEntropySize + MaximumCheckSumSize) / WordBitLength
)

// Error type specifically for mnemonic errors
type mnemonicError struct {
	error
	Message string
}

func (err mnemonicError) Error() string {
	return err.Message
}

// Type to wrap mnemonic-related methods
type Mnemonic struct {
	Entropy entropy
	Checksum byte
	Sentence []uint32
}

// Generate mnemonic sentence based on the size of entropy bits.
// An error is returned when there is an error, in which case the
// Mnemomic returned will be in an undefined state.
func GenerateMnemonic(size uint16) (Mnemonic, error) {
	// Get new entropy
	ent, err := GenerateEntropy(size)

	if (err != nil) { return Mnemonic{}, err }

	// Get entropy + checksum
	checksum, checksumErr := ent.GenerateChecksum()

	if (checksumErr != nil) { return Mnemonic{}, checksumErr }

	fullData := append(ent.Data, checksum)

	// Get byte reader from concatenation of entropy and checksum
	byteReader := bytes.NewReader(fullData)

	// Get bit reader
	bitReader := bitreader.NewBitReader(byteReader)

	// Get sentence
	words := make([]uint32, len(fullData) * 8 / WordBitLength) // (length of bits of entropy + checksum) / length of word

	for i := 0; i < len(words); i++ {
		// Try to read 11 bits
		var bitError error
		words[i], bitError = bitReader.Read32(WordBitLength)

		// If there's an error, return it as a mnemonicError
		if (bitError != nil) {
			return Mnemonic{}, mnemonicError{Message: bitError.Error()}
		}
	}

	return Mnemonic{ent, checksum, words}, nil
}

// Get the sentence that the mnemonic's indices correspond to in a Wordlist.
// An error is returned in the case that reading from the wordlist fails.
func (mnemonic Mnemonic) GetSentenceFrom(wordlist wordlist.Wordlist) ([]string, error) {
	words := make([]string, len(mnemonic.Sentence))

	for i := 0; i < len(words); i++ {
		var getErr error
		words[i], getErr = wordlist.GetWordAt(mnemonic.Sentence[i])

		if (getErr != nil) {
			return []string{}, mnemonicError{Message: getErr.Error()}
		}
	}

	return words[:], nil
}