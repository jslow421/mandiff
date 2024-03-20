package main

import (
	"testing"
)

// Test the checkLanguage function
func TestIdentifyEnglishLineOfText(t *testing.T) {
	// Call the checkLanguage function
	isEnglish, err := checkIfTextIsLikelyEnglish("Mount the Ranger pressure infusor Model 145 on only a 3M model 90068/90124 pressure infusor I.V. pole/base")

	if err != nil {
		t.Errorf("Error calling checkLanguage: %v", err)
	}
	expected := true

	if isEnglish != expected {
		t.Errorf("Expected %t, but got %t", expected, isEnglish)
	}
}

func TestIdentifyingNonEnglishText(t *testing.T) {
	// Call the checkLanguage function
	isEnglish, err := checkIfTextIsLikelyEnglish("Installez le dispositif de perfusion sous pression Ranger, modèle 145, sur une potence/base I.V. 3M, modèle 90068/90124")

	if err != nil {
		t.Errorf("Error calling checkLanguage: %v", err)
	}
	expected := false

	if isEnglish != expected {
		t.Errorf("Expected %t, but got %t", expected, isEnglish)
	}
}
