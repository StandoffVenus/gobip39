package test

import (
	"testing"
	"gobip39"
	"reflect"
	"bytes"
	"crypto/sha256"
)

// Name of entropyError
const ENTROPY_ERROR = "entropyError"

func TestEntropy_GenerateEntropy_FailsOnLowerBoundViolation(t *testing.T) {
	_, err := gobip39.GenerateEntropy(gobip39.MinimumEntropySize - 1)

	if (!isEntropyError(err)) {
		t.Error("Expected GenerateEntropy to return an entropyError when size <", gobip39.MinimumEntropySize)
	}
}

func TestEntropy_GenerateEntropy_FailsOnUpperBoundViolation(t *testing.T) {
	_, err := gobip39.GenerateEntropy(gobip39.MaximumEntropySize + 1)

	if (!isEntropyError(err)){
		t.Error("Expected GenerateEntropy to return an entropyError when size >", gobip39.MaximumEntropySize)
	}
}

func TestEntropy_GenerateEntropy_FailsOnNon32MultipleSize(t *testing.T) {
	_, err := gobip39.GenerateEntropy(gobip39.MinimumEntropySize + 1)

	if (!isEntropyError(err)) {
		t.Error("Expected GenerateEntropy to return an entropyError when size % 32 != 0.")
	}
}

func TestEntropy_GenerateEntropy_ReturnsAnEntropyInstance(t *testing.T) {
	var size uint16 = gobip39.MinimumEntropySize

	entropy, err := gobip39.GenerateEntropy(size)

	if (err != nil) {
		t.Error("Expected GenerateEntropy to return nil error.")
	}

	if (entropy.Data == nil) {
		t.Error("Expected GenerateEntropy to return an Entropy with non-nil Data.")
	}

	if (entropy.Size != size) {
		t.Error("Entropy's size is not equal to size passed; expected", size, "\b.")
	}

	if (uint16(len(entropy.Data)) * 8 != entropy.Size) {
		t.Error("Entropy's Size and length of Data do not match; expected", len(entropy.Data) * 8, "==", entropy.Size, "\b.")
	}
}

func TestEntropy_GetEntropyFromBytes_FailsWhenLengthExceedsUInt16MaxBits(t *testing.T) {
	// Add 1 to max size of uint16
	maxPlusOne := int32(^uint16(0))

	arr := make([]byte, maxPlusOne)

	_, err := gobip39.GetEntropyFromBytes(arr)

	if (!isEntropyError(err)) {
		t.Error("Expected GetEntropyFromBytes to return an entropyError when the byte array length exceeds uint16's maximum.")
	}
}

func TestEntropy_GetEntropyFromBytes_FailsWhenBitLengthViolatesUpperBound(t *testing.T) {
	// Add 1 to max bytes
	maxPlusOne := (gobip39.MaximumEntropySize + 8) / 8

	arr := make([]byte, maxPlusOne)

	_, err := gobip39.GetEntropyFromBytes(arr)

	if (!isEntropyError(err)) {
		t.Error("Expected GetEntropyFromBytes to return an entropyError when the byte array length (in bits) exceeds", gobip39.MaximumEntropySize, "\b.")
	}
}

func TestEntropy_GetEntropyFromBytes_FailsWhenBitLengthViolatesLowerBound(t *testing.T) {
	// Subtract 1 from max bytes
	minMinusOne := (gobip39.MinimumEntropySize - 8) / 8

	arr := make([]byte, minMinusOne)

	_, err := gobip39.GetEntropyFromBytes(arr)

	if (!isEntropyError(err)) {
		t.Error("Expected GetEntropyFromBytes to return an entropyError when the byte array length (in bits) falls below", gobip39.MinimumEntropySize, "\b.")
	}
}

func TestEntropy_GetEntropyFromBytes_FailsWhenBitLengthIsNot32Multiple(t *testing.T) {
	// Add 1 from minimum bytes to mess up length % 32 = 0
	non32Multiple := (gobip39.MinimumEntropySize + 8) / 8

	arr := make([]byte, non32Multiple)

	_, err := gobip39.GetEntropyFromBytes(arr)

	if (!isEntropyError(err)) {
		t.Error("Expected GetEntropyFromBytes to return an entropyError when the byte array length (in bits) is not a multiple of 32.")
	}
}

func TestEntropy_GetEntropyFromBytes_ReturnsEntropyInstance(t *testing.T) {
	arr := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}

	ent, err := gobip39.GetEntropyFromBytes(arr)

	if (isEntropyError(err)) {
		t.Error("Expected GetEntropyFromBytes to return nil error on valid byte array.")
	}

	if (ent.Size != uint16(len(arr)) * 8) {
		t.Error("Entropy does not have correct Size; expected", len(arr) * 8, "==", ent.Size, "\b.")
	}

	if (!bytes.Equal(arr, ent.Data)) {
		t.Error("Entropy Data does not equal passed byte array; expected", arr, "==", ent.Data, "\b.")
	}
}

func TestEntropy_GenerateChecksum_ReturnsByteHoldingChecksum(t *testing.T) {
	arr := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}

	ent, _ := gobip39.GetEntropyFromBytes(arr)

	checksum, err := ent.GenerateChecksum()

	if (err != nil) {
		t.Error("GenerateChecksum should not generate errors on valid Entropy.")
	}

	// Need to SHA256 the Entropy.Data
	shaHash := sha256.New()
	shaHash.Write(ent.Data)
	firstByte := shaHash.Sum(nil)[0]

	// Need to disregard last 4 bits since checksum only will read the first 4 bits.
	// However, this is future proofed, hence the subtraction from 8.
	// This is how checksum's length is calculated.
	expectedChecksum := firstByte >> (8 - ent.Size / 32)

	if checksum != expectedChecksum {
		t.Error("Checksum holds the wrong value; expected", checksum, "==", expectedChecksum, "\b.")
	}
}

func isEntropyError(e error) bool {
	if (e == nil) {
		return false
	}

	return reflect.TypeOf(e).Name() == ENTROPY_ERROR
}