package plugins

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type MemValue struct {
	Original        string
	NumberAsInt     int
	NumberAsString  string
	Units           string
	IntMemsize      int
	Converted       bool
	CorrectStrValue string
	Error           error
}

func GetMemoryValue(memstr string) *MemValue {
	memvalue := MemValue{
		Original: memstr,
	}
	memvalue.MemStrToMemValue()
	return &memvalue
}

func (memval *MemValue) memStrToInt() {
	memval.IntMemsize = memval.NumberAsInt
	if memval.Units != "" {
		memval.IntMemsize = memval.IntMemsize * 1024
		if memval.Units != "kB" {
			memval.IntMemsize = memval.IntMemsize * 1024
			if memval.Units != "MB" {
				memval.IntMemsize = memval.IntMemsize * 1024
				if memval.Units != "GB" {
					memval.IntMemsize = memval.IntMemsize * 1024
				}
			}
		}
	}
}

func (memval *MemValue) optimise() {
	for memval.NumberAsInt >= 1024 && memval.Units != "TB" {
		memval.NumberAsInt = memval.NumberAsInt / 1024
		switch memval.Units {
		case "":
			memval.Units = "kB"
		case "kB":
			memval.Units = "MB"
		case "MB":
			memval.Units = "GB"
		case "GB":
			memval.Units = "TB"
		}
	}
}

func (result *MemValue) MemStrToMemValue() {
	// Convert things like
	// 1024 to 1kB
	// 1mb to 1MB
	// 8388608 to 8MB
	// 13813160064 to 12GB

	result.Converted = false

	re1 := regexp.MustCompile("[-0-9]+")
	re2 := regexp.MustCompile("(kB|KB|Kb|kb|MB|mb|Mb|mB|GB|gb|Gb|gB|TB|tb|Tb|tB)")
	result.NumberAsString = re1.FindString(result.Original)
	result.Units = re2.FindString(result.Original)
	if result.NumberAsString != "" && fmt.Sprintf("%s%s", result.NumberAsString, result.Units) == result.Original {
		result.Units = strings.ToUpper(result.Units)
		if result.Units == "KB" {
			result.Units = "kB"
		}
		numberAsInt, err := strconv.Atoi(result.NumberAsString)
		if err != nil {
			result.Error = errors.New(fmt.Sprintf("Failed to convert memory string (1) %s\n\t%s", result.Original, err))
		} else {
			result.NumberAsInt = numberAsInt
			result.optimise()
			result.memStrToInt()
			result.Converted = true
			result.CorrectStrValue = fmt.Sprintf("%s%s", result.NumberAsString, result.Units)
		}
	} else {
		result.Error = errors.New(fmt.Sprintf("Failed to convert memory string (2) [%s]", result.Original))
	}
	if result.Error != nil {
		log.Printf("%s memory string conversion failed: %s\n", result.Error, result.Error)
	}
}

func GetMaxMemoryOfContainerAsString(imposedLimit string) string {
	cgroup_mem_limit_fname := "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	_, err := os.Stat(cgroup_mem_limit_fname)
	if err == nil {
		cgroup_mem_limit, errRead := ioutil.ReadFile(cgroup_mem_limit_fname)
		if errRead == nil {
			cgroup_mem_limit_str := strings.Trim(string(cgroup_mem_limit), "\n")
			if imposedLimit != "" && imposedLimit < cgroup_mem_limit_str {
				log.Printf("WARNING: Imposing memory limit of %s", imposedLimit)
				return imposedLimit
			}
			return cgroup_mem_limit_str
		} else {
			return "0"
		}
	} else {
		return "0"
	}
}

func (memValue *MemValue) LessThan(otherValue *MemValue) bool {
	if memValue.Converted &&
		otherValue.Converted &&
		memValue.IntMemsize < otherValue.IntMemsize {
		return true
	}
	return false
}
