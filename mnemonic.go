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

// Error type specifically for Mnemonic errors
type mnemonicError struct {
	Message string
}

func (err mnemonicError) Error() string {
	return err.Message
}

// Type to wrap Mnemonic-related methods
type Mnemonic struct {
	Entropy Entropy
	Checksum byte
	Sentence []uint32
}

// Generate Mnemonic from Entropy.
// An error is returned if Entropy's size is outside
// the domain of valid entropy size, or if the
// size is not a multiple of 32, in which case the
// Mnemonic returned is in an invalid state.
func GetMnemonicFromEntropy(ent Entropy) (Mnemonic, error) {
	if (ent.Size > MaximumEntropySize || ent.Size < MinimumEntropySize) {
		return Mnemonic{}, mnemonicError{Message: "Size of entropy was outside of domain [128, 256]."}
	}

	if (ent.Size % 32 != 0) {
		return Mnemonic{}, mnemonicError{Message: "Size of entropy was not a multiple of 32."}
	}

	// Get entropy + checksum
	checksum, checksumErr := ent.GenerateChecksum()

	if (checksumErr != nil) { return Mnemonic{}, mnemonicError{Message: checksumErr.Error()} }

	// Consider the following example:
	// Entropy size: 128 bits
	// Ergo, checksum size: 4 bits
	// Thus, our checksum may look like 00000011 (3)
	// When appended to our entropy of all 0's, fullData looks like:
	// [ 00000000, ... 00000000, 00000011 ]
	// However, since we get fullData's length (136) divided by WordBitLength (11),
	// the for loop runs from 0 to 12 (truncated integer division), meaning
	// the last bit read is the 12 * 11 bit. Or, in this case:
	// [ ..., 00 (0) 000011 ]
	// Hence, after some math, it's clear the checksum must be shifted left
	// to combat that this offset problem. We fix it by shifting it
	// 8 - length of entropy / 32; this shifts the checksum's bits
	// to the beginning of the byte (though, not in a different Endian
	// style). Thus, the checksum's real value is read.
	fullData := append(ent.Data, checksum << (8 - ent.Size / 32))

	// Get byte reader from concatenation of entropy and checksum
	byteReader := bytes.NewReader(fullData)

	// Get bit reader
	bitReader := bitreader.NewBitReader(byteReader)

	// Make empty sentence
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

// Helper method that calls GetEntropyFromBytes with the given
// hex data, then calls GetMnemonicFromEntropy if returned a valid
// entropy. Returns the resultant mnemonic if there is no error.
// An error is returned if the size (in bits) of the data
// is outside the domain of valid entropy size, or if the
// length of the data is not a multiple of 32, in which
// case the Mnemonic returned is in an invalid state.
func GetMnemonicFromBytes(data []byte) (Mnemonic, error) {
	entropy, err := GetEntropyFromBytes(data)

	if (err != nil) { return Mnemonic{}, mnemonicError{Message: err.Error()} }

	return GetMnemonicFromEntropy(entropy)
}

// Convenience method to generate Mnemonic with "size" bits of entropy.
// An error is returned when entropy size is invalid (outside domain
// [128, 256], entropy size is not a multiple of 32, or reading the
// 11 bit portions of the entropy + checksum data fails, in which case the
// Mnemomic returned will be in an undefined state.
func GenerateMnemonic(size uint16) (Mnemonic, error) {
	// Get new entropy
	ent, err := GenerateEntropy(size)

	if (err != nil) { return Mnemonic{}, mnemonicError{Message: err.Error()} }

	return GetMnemonicFromEntropy(ent)
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