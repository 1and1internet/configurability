package main

import (
	"configurability/file_helpers"
)

func main() {
	cd := file_helpers.CustomisationData{}
	cd.LoadCustomisationData()
}