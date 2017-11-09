package gobip39

// This file wraps entropy and checksum specifications as detailed by
// BIP-0039 spec.

import (
	"crypto/sha256"
	"crypto/rand"
	"github.com/32bitkid/bitreader"
	"bytes"
	"fmt"
)

const (
	MinimumEntropySize = 128
	MaximumEntropySize = 256
	MinimumChecksumSize = MinimumEntropySize / 32
	MaximumCheckSumSize = MaximumEntropySize / 32
)

// Error type specifically for entropy errors
type entropyError struct {
	Message string
}

func (err entropyError) Error() string {
	return err.Message
}

// Type to wrap entropy-related methods
type Entropy struct {
	Size uint16
	Data []byte
}

// Generate entropy with specified amount of bits.
// GenerateEntropy will return an error if the passed size
// does not conform to specification limits:
// 1) Size of Entropy (in bits) is within domain [128, 256],
// and 2) Size of Entropy is not a multiple of 32.
// In the case of error, the returned Entropy will be
// in an invalid state.
func GenerateEntropy(size uint16) (Entropy, error) {
	// If size is outside allowed domain
	if (size < MinimumEntropySize || size > MaximumEntropySize) {
		return Entropy{}, entropyError{Message: "Size of entropy is out of domain [128, 256]."}
	}

	// If size is not divisible by 32
	if (size % 32 != 0) {
		return Entropy{}, entropyError{Message: "Size of entropy is not a multiple of 32."}
	}

	// Read random bytes
	randomBytes := make([]byte, size / 8)
	rand.Read(randomBytes)

	// Return Entropy from these random bytes
	return GetEntropyFromBytes(randomBytes)
}

// Generate Entropy from hex data.
// An error is returned if the size (in bits) of the data
// is outside the domain of valid entropy size, or if the
// length of the data is not a multiple of 32, in which
// the Entropy returned is in an invalid state.
func GetEntropyFromBytes(data []byte) (Entropy, error) {
	// Early check to prevent errors from wrap around on uint16 conversion
	if (int(^uint16(0)) < len(data)) {
		return Entropy{}, entropyError{Message: "Length of data (in bits) was outside of domain [128, 256]."}
	}

	dataLength := uint16(len(data)) * 8

	if (dataLength > MaximumEntropySize || dataLength < MinimumEntropySize) {
		return Entropy{}, entropyError{Message: "Length of data (in bits) was outside of domain [128, 256]."}
	}

	if (dataLength % 32 != 0) {
		return Entropy{}, entropyError{Message: "Length of data (in bits) was not a multiple of 32."}
	}

	return Entropy{Size: dataLength, Data: data}, nil
}

// Generate checksum will return a checksum of the given entropy
// via the specification, returning said checksum as a byte. An empty
// byte and error will be returned if an error is encountered while
// trying to read the checksum.
func (ent Entropy) GenerateChecksum() (byte, error) {
	// We must SHA the data for the checksum
	shaHash := sha256.New()
	shaHash.Write(ent.Data)

	// Convert data to a Reader
	byteReader := bytes.NewReader(shaHash.Sum(nil))

	bitReader := bitreader.NewBitReader(byteReader)

	// Read the first (entropy's bits / 32) bits of the SHA256 digest
	fmt.Println(ent.Size)

	checksum, err := bitReader.Read32(uint(ent.Size / 32))

	if (err != nil) { return 0, entropyError{ Message: err.Error() }}

	return byte(checksum), nil
}