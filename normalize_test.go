package gogle

import "testing"

func TestNormalize(t *testing.T) {
	var input = []rune("hellÒ    \t \n world")
	var expectedOutput = "hello world"
	Normalize(&input)

	if string(input) != expectedOutput {
		t.Fatalf("Expected: %v\nGOT: %v\n", expectedOutput, string(input))
	}
}
