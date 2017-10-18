package gobip39

// This file wraps entropy and checksum specifications as detailed by
// BIP-0039 spec.

import (
	"crypto/sha256"
	"crypto/rand"
	"github.com/32bitkid/bitreader"
	"bytes"
)

const (
	MinimumEntropySize = 128
	MaximumEntropySize = 256
	MinimumChecksumSize = MinimumEntropySize / 32
	MaximumCheckSumSize = MaximumEntropySize / 32
)

// Error type specifically for entropy errors
type entropyError struct {
	error
	Message string
}

func (err entropyError) Error() string {
	return err.Message
}

// Type to wrap entropy-related methods
type entropy struct {
	Size uint16
	Data []byte
}

// Generate entropy with specified amount of bits.
// Returns the generated entropy and potential error.
// GenerateEntropy will return an error if the passed size
// does not conform to specification limits.
func GenerateEntropy(size uint16) (entropy, error) {
	// If size is outside allowed domain
	if (size < MinimumEntropySize || size > MaximumEntropySize) {
		return entropy{}, entropyError{Message: "Size of entropy is out of range [128, 256]."}
	}

	// If size is not divisible by 32
	if (size % 32 != 0) {
		return entropy{}, entropyError{Message: "Size of entropy is not divisible by 32."}
	}

	// Read random bytes
	randomBytes := make([]byte, size / 8)
	rand.Read(randomBytes)

	// Make SHA256 hash and read random data into it
	shaHash := sha256.New()
	shaHash.Write(randomBytes)

	return entropy{size, shaHash.Sum(nil)}, nil
}

// Generate checksum will return a checksum of the given entropy
// via the specification, returning said checksum as a byte. An empty
// byte and error will be returned if an error is encountered while
// trying to read the checksum.
func (ent entropy) GenerateChecksum() (byte, error) {
	// Convert entropy data to a Reader
	byteReader := bytes.NewReader(ent.Data)

	bitReader := bitreader.NewBitReader(byteReader)

	// Read the first (entropy's bits / 32) bits of the entropy data
	checksum, err := bitReader.Read32(uint(ent.Size / 32))

	if (err != nil) { return 0, entropyError{ Message: err.Error() }}

	return byte(checksum), nil
}