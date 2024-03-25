package main

import (
	"testing"
)

func TestBuildPrompt(test *testing.T) {
	prompt, err :=
		buildDocumentComparisonPrompt(
			"This is a document that will be used for comparison",
			"This is a document that will also be used for comparison")
	if len(prompt) < 20 || err != nil {
		test.Errorf("Error building prompt: %s", err)
	}
}

func TestPromptTypeIsEducation(t *testing.T) {
	prompt := promptFactory("", "", "education.txt")
	expected := EducationExtraction

	if prompt.PromptType != expected {
		t.Errorf("Prompt types should be the same, got %v and %v", expected, prompt.PromptType)
	}
}

func TestPromptTypeIsComparison(t *testing.T) {
	prompt := promptFactory("document1.txt", "document2.txt", "")
	expected := DocumentComparison

	if prompt.PromptType != expected {
		t.Errorf("Prompt types should be the same, got %v and %v", expected, prompt.PromptType)
	}
}
