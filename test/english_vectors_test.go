package test

import (
	"testing"
	"io/ioutil"
	"encoding/json"
	"encoding/hex"
	"gobip39"
	"gobip39/wordlist"
	"strings"
	"bytes"
)

/*
	Vector[0] holds the entropy data (in hex)
	Vector[1] holds the sentence
	Vector[2] holds the seed
	Vector[3] holds mnemonic data
 */
type Vector []string

// The passphrase is "TREZOR" for all tests
const PASSPHRASE = "TREZOR"

func TestEnglishVectors(t *testing.T) {
	file, err := ioutil.ReadFile("./vectors.json")

	if (err != nil) {
		t.Error("Failed to read from required vector file 'english_vectors.json':", err.Error())
	}

	// Anonymous struct that holds vectors
	var vectors struct {
		Vectors []Vector `json:"english"`
	}

	if marshalErr := json.Unmarshal(file, &vectors); marshalErr != nil {
		t.Error("Failed to unmarshal file data:", marshalErr.Error())
	}

	for _, v := range vectors.Vectors {
		// Entropy data (decoded from hex)
		entropyHex, decodeErr := hex.DecodeString(v[0])

		if (decodeErr != nil) {
			t.Error("Failed to decode hex data.\nData:", v[0], "\nError:", decodeErr.Error())
		}

		entropy, entropyErr := gobip39.GetEntropyFromBytes(entropyHex)

		if (entropyErr != nil) {
			t.Error("Failed to generate Entropy from hex data:", entropyErr.Error())
		}

		// Turn this Entropy into Mnemonic
		mnemonic, mnemonicErr := gobip39.GetMnemonicFromEntropy(entropy)

		if (mnemonicErr != nil) {
			t.Error("Failed to generate Mnemonic from Entropy:", mnemonicErr.Error())
		}

		// Get sentence from mnemonic
		sentence, sentenceErr := mnemonic.GetSentenceFrom(wordlist.English)

		if (sentenceErr != nil) {
			t.Error("Failed to generate sentence from Mnemonic:", sentenceErr.Error())
		}

		// Sentence does not match testing vector
		if joinedSentence := strings.Join(sentence, " "); joinedSentence != v[1] {
			t.Error("Expected mnemonic sentence", joinedSentence, "to equal", v[1])
		}

		// Assert that this Mnemonic's binary seed is equivalent to the expected seed
		actualSeed := gobip39.GenerateBinarySeed(strings.Join(sentence, " "), PASSPHRASE)

		expectedSeed, seedDecodeErr := hex.DecodeString(v[2])

		if (seedDecodeErr != nil) {
			t.Error("Failed to to decode hex data.\nData:", v[2], "\nError:", seedDecodeErr.Error())
		}

		if !bytes.Equal(actualSeed, expectedSeed) {
			t.Error("Expected binary seed data", actualSeed, "to equal", expectedSeed)
		}
	}
}
