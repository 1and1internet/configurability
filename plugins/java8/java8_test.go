package main

import (
	"fmt"
	"testing"
)

func TestRounding(t *testing.T) {
	size := GetRoundedTo1024(2147483661)
	if size != 2147483648 {
		t.Fatal(fmt.Sprintf("%v != %v", size, 2147483648))
	}
	size = GetRoundedTo1024(2147484161)
	if size != 2147484672 {
		t.Fatal(fmt.Sprintf("%v != %v", size, 2147484672))
	}
}

func TestNiceMemoryString(t *testing.T) {
	size := GetMemoryInMultiplesOf1024AsTidySuffixedString(2147483648)
	if size != "2048M" {
		t.Fatal(fmt.Sprintf("%v != 2048M", size))
	}
}

func TestMaxMemory(t *testing.T) {
	allInfo := CustomisationInfo{}
	allInfo.GetMaxMemory()
	if allInfo.MaxMemoryBytes != 8589934592 {
		t.Fatal(fmt.Sprintf("%v != 8589934592", allInfo.MaxMemoryBytes))
	}
}
