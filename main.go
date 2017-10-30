package main

import (
	"github.com/1and1internet/configurability/file_helpers"
)

func main() {
	cd := file_helpers.CustomisationData{}
	cd.LoadCustomisationData()
}