package main

import (
	"testing"
)

func TestIdentifyEnglishLineOfText(t *testing.T) {
	isEnglish := checkIfTextIsLikelyEnglish("Mount the Ranger pressure infusor Model 145 on only a 3M model 90068/90124 pressure infusor I.V. pole/base")

	expected := true

	if isEnglish != expected {
		t.Errorf("Expected %t, but got %t", expected, isEnglish)
	}
}

func TestIdentifyingNonEnglishText(t *testing.T) {
	isEnglish := checkIfTextIsLikelyEnglish("Installez le dispositif de perfusion sous pression Ranger, modèle 145, sur une potence/base I.V. 3M, modèle 90068/90124")

	expected := false

	if isEnglish != expected {
		t.Errorf("Expected %t, but got %t", expected, isEnglish)
	}
}
