package gobip39

// This file handles seed generation as detailed by
// BIP-0039 spec.

import (
	"golang.org/x/text/unicode/norm"
	"golang.org/x/crypto/pbkdf2"
	SHA512 "crypto/sha512"
)

const (
	Pbkdf2Iterations = 2048
	KeyLengthBits = 512
	KeyLengthBytes = KeyLengthBits / 8
)

// Generate the binary seed, with an optional passphrase, from a mnemonic sentence.
func GenerateBinarySeed(mnemonic string, passphrase ...string) []byte{
	_passphrase := ""

	if (passphrase != nil) {
		_passphrase = passphrase[0]
	}

	normalizedMnemonic := norm.NFKD.Bytes([]byte(mnemonic))
	normalizedPassphrase := norm.NFKD.Bytes([]byte("mnemonic" + _passphrase))

	// The reason we do not have to pass HMAC-SHA512.New is because of the fact that
	// pbkdf2.Key will call HMAC on passwords for you.
	return pbkdf2.Key(normalizedMnemonic, normalizedPassphrase, Pbkdf2Iterations, KeyLengthBytes, SHA512.New)
}
