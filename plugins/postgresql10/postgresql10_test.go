package main

import (
	"fmt"
	"testing"
)

func TestMaxMemory(t *testing.T) {
	allInfo := CustomisationInfo{}
	allInfo.SetMaxMemory()
	if allInfo.MaxMemory.CorrectOptimisedStrValue != "16GB" {
		t.Fatal(fmt.Sprintf("%v != 16GB (%v)", allInfo.MaxMemory.CorrectOptimisedStrValue, allInfo.MaxMemory))
	}
}
