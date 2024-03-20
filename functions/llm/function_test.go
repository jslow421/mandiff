package main

import (
	"testing"
)

func TestBuildPrompt(test *testing.T) {
	prompt, err := buildPrompt("document1.txt", "document2.txt")
	if len(prompt) < 20 || err != nil {
		test.Errorf("Prompt is likely too short: %s", prompt)
	}
}
