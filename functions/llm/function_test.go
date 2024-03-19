package main

import (
	"testing"
)

func TestBuildPrompt(test *testing.T) {
	prompt := buildPrompt()
	if len(prompt) < 20 {
		test.Errorf("Prompt is likely too short: %s", prompt)
	}
}

func TestRunRequest(test *testing.T) {
	runRequest()
}
