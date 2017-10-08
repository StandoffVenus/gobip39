package gobip39

// This file wraps entropy and checksum specifications as detailed by
// BIP-0039 spec.

import (
	"crypto/sha256"
	"crypto/rand"
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
// via the specification, returning said checksum as a byte array.
func (ent entropy) GenerateChecksum() []byte {
	// Get hash data from 0 to amount of entropy's bits / 32
	return ent.Data[:ent.Size / 32]
}