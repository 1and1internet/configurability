package main

import (
	"fmt"
	"testing"

	"github.com/1and1internet/configurability/plugins"
)

func TestMaxMemory(t *testing.T) {
	allInfo := CustomisationInfo{}
	allInfo.SetMaxMemory()
	if allInfo.MaxMemory.CorrectOptimisedStrValue != "16GB" {
		t.Fatal(fmt.Sprintf("%v != 16GB", allInfo.MaxMemory.CorrectOptimisedStrValue))
	}
}

func TestSetMemVal_reqval(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "4GB",
	}
	systemMax := plugins.GetMemoryValue("9GB")
	confline.SetMemVal("8GB", "4GB", "1kB", "16GB", *systemMax)
	if confline.Value != "8GB" {
		t.Fatal(fmt.Sprintf("%v != 8GB", confline.Value))
	}
	if confline.UseOrig != false {
		t.Fatal(fmt.Sprintf("%v != false", confline.UseOrig))
	}
}

func TestSetMemVal_sysmaxlim(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "4GB",
	}
	systemMax := plugins.GetMemoryValue("1234")
	confline.SetMemVal("8GB", "4GB", "1kB", "16GB", *systemMax)
	if confline.Value != "4GB" {
		t.Fatal(fmt.Sprintf("%v != 4GB", confline.Value))
	}
	if confline.UseOrig != false {
		t.Fatal(fmt.Sprintf("%v != false", confline.UseOrig))
	}
}

func TestSetMemVal_usedefault(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "4GB",
	}
	systemMax := plugins.GetMemoryValue("9GB")
	confline.SetMemVal("4GB", "4GB", "1kB", "16GB", *systemMax)
	if confline.UseOrig != true {
		t.Fatal(fmt.Sprintf("%v != true", confline.UseOrig))
	}
}

func TestSetMemVal_minlimit(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "4GB",
	}
	systemMax := plugins.GetMemoryValue("90GB")
	confline.SetMemVal("1", "4GB", "1kB", "16GB", *systemMax)
	if confline.Value != "4GB" {
		t.Fatal(fmt.Sprintf("%v != 4GB", confline.Value))
	}
	if confline.UseOrig != false {
		t.Fatal(fmt.Sprintf("%v != false", confline.UseOrig))
	}
}

func TestSetMemVal_maxlimit(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "4GB",
	}
	systemMax := plugins.GetMemoryValue("90GB")
	confline.SetMemVal("17GB", "4GB", "1kB", "16GB", *systemMax)
	if confline.Value != "4GB" {
		t.Fatal(fmt.Sprintf("%v != 4GB", confline.Value))
	}
	if confline.UseOrig != false {
		t.Fatal(fmt.Sprintf("%v != false", confline.UseOrig))
	}
}

func TestSetTimeVal(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "1min",
	}
	confline.SetTimeVal("10min", "1min", "0", "1d")
	if confline.Value != "10min" {
		t.Fatal(fmt.Sprintf("%v != 10min", confline.Value))
	}
	if confline.UseOrig != false {
		t.Fatal(fmt.Sprintf("%v != false", confline.UseOrig))
	}
}

func TestSetTimeVal_usedefault(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "1min",
	}
	confline.SetTimeVal("1min", "1min", "0", "1d")
	if confline.UseOrig != true {
		t.Fatal(fmt.Sprintf("%v != true", confline.UseOrig))
	}
}

func TestSetTimeVal_minlimit(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "1min",
	}
	confline.SetTimeVal("10s", "2min", "1min", "1d")
	if confline.Value != "2min" {
		t.Fatal(fmt.Sprintf("%v != 2min", confline.Value))
	}
	if confline.UseOrig != false {
		t.Fatal(fmt.Sprintf("%v != false", confline.UseOrig))
	}
}

func TestSetTimeVal_maxlimit(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "1min",
	}
	confline.SetTimeVal("2d", "2min", "1min", "1d")
	if confline.Value != "2min" {
		t.Fatal(fmt.Sprintf("%v != 2min", confline.Value))
	}
	if confline.UseOrig != false {
		t.Fatal(fmt.Sprintf("%v != false", confline.UseOrig))
	}
}

func TestIntVal(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "10",
	}
	confline.SetIntVal(20, 10, 1, 100)
	if confline.Value != "20" {
		t.Fatal(fmt.Sprintf("%v != 20", confline.Value))
	}
	if confline.UseOrig != false {
		t.Fatal(fmt.Sprintf("%v != false", confline.UseOrig))
	}
}

func TestIntVal_usedefault(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "10",
	}
	confline.SetIntVal(10, 10, 1, 100)
	if confline.Value != "10" {
		t.Fatal(fmt.Sprintf("%v != 10", confline.Value))
	}
	if confline.UseOrig != true {
		t.Fatal(fmt.Sprintf("%v != true", confline.UseOrig))
	}
}

func TestIntVal_maxlimit(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "10",
	}
	confline.SetIntVal(1000, 10, 1, 100)
	if confline.Value != "10" {
		t.Fatal(fmt.Sprintf("%v != 10", confline.Value))
	}
	if confline.UseOrig != true {
		t.Fatal(fmt.Sprintf("%v != true", confline.UseOrig))
	}
}

func TestIntVal_minlimit(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "10",
	}
	confline.SetIntVal(1, 10, 5, 100)
	if confline.Value != "10" {
		t.Fatal(fmt.Sprintf("%v != 10", confline.Value))
	}
	if confline.UseOrig != true {
		t.Fatal(fmt.Sprintf("%v != true", confline.UseOrig))
	}
}

func TestFloatVal(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "12.34",
	}
	confline.SetFloatVal("34.56", "12.34", "1.0", "100.0")
	if confline.Value != "34.56" {
		t.Fatal(fmt.Sprintf("%v != 34.56", confline.Value))
	}
	if confline.UseOrig != false {
		t.Fatal(fmt.Sprintf("%v != false", confline.UseOrig))
	}
}

func TestFloatVal_usedefault(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "12.34",
	}
	confline.SetFloatVal("NOTAFLOAT", "12.34", "1.0", "100.0")
	if confline.Value != "12.34" {
		t.Fatal(fmt.Sprintf("%v != 12.34", confline.Value))
	}
	if confline.UseOrig != true {
		t.Fatal(fmt.Sprintf("%v != true", confline.UseOrig))
	}
}

func TestFloatVal_maxlimit(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "12.34",
	}
	confline.SetFloatVal("134.56", "12.34", "1.0", "100.0")
	if confline.Value != "12.34" {
		t.Fatal(fmt.Sprintf("%v != 12.34", confline.Value))
	}
	if confline.UseOrig != true {
		t.Fatal(fmt.Sprintf("%v != true", confline.UseOrig))
	}
}

func TestFloatVal_minlimit(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "12.34",
	}
	confline.SetFloatVal("1.0", "12.34", "10.0", "100.0")
	if confline.Value != "12.34" {
		t.Fatal(fmt.Sprintf("%v != 12.34", confline.Value))
	}
	if confline.UseOrig != true {
		t.Fatal(fmt.Sprintf("%v != true", confline.UseOrig))
	}
}

func TestStrVal(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "standard_value",
	}
	confline.SetStrVal("nonstandard_value", "standard_value")
	if confline.Value != "nonstandard_value" {
		t.Fatal(fmt.Sprintf("%v != nonstandard_value", confline.Value))
	}
	if confline.UseOrig != false {
		t.Fatal(fmt.Sprintf("%v != false", confline.UseOrig))
	}
}

func TestStrVal_emptystring(t *testing.T) {
	confline := &ConfigLine{
		UseOrig: true,
		Value:   "standard_value",
	}
	confline.SetStrVal("", "standard_value")
	if confline.Value != "standard_value" {
		t.Fatal(fmt.Sprintf("%v != standard_value", confline.Value))
	}
	if confline.UseOrig != true {
		t.Fatal(fmt.Sprintf("%v != true", confline.UseOrig))
	}
}
