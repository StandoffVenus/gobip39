package test

import (
  "testing"
  "gobip39/wordlist"
)

var stringArray []string = []string{"first", "second", "third", "fourth"}

func TestWordlist_FindWordIn_ShouldReturnNegative1WhenNoMatch(t *testing.T) {
  if (wordlist.FindWordIn(stringArray, "fifth") != -1) {
    t.Error("FindWordIn didn't return -1 when string \"fifth\" was not within array.")
  }
}

func TestWordlist_FindWordIn_ShouldReturnCorrectIndexWhenMatch(t *testing.T) {
  firstIndex := wordlist.FindWordIn(stringArray, "first")
  thirdIndex := wordlist.FindWordIn(stringArray, "third")

  if (firstIndex != 0) {
    t.Error("FindWordIn didn't return 0 for string \"first\"'s index.")
  }

  if (thirdIndex != 2) {
    t.Error("FindWordIn didn't return 2 for string \"first\"'s index.")
  }
}