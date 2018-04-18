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
	Original                 string
	OptNumberPartAsInt       int
	OrigNumberPartAsString   string
	Units                    string
	ActualIntMemsize         int
	Converted                bool
	CorrectOptimisedStrValue string
	Error                    error
}

func GetMemoryValue(memstr string) *MemValue {
	memvalue := MemValue{
		Original: memstr,
	}
	memvalue.MemStrToMemValue()
	return &memvalue
}

func (memval *MemValue) memStrToInt() {
	memval.ActualIntMemsize = memval.OptNumberPartAsInt
	if memval.Units != "" {
		memval.ActualIntMemsize = memval.ActualIntMemsize * 1024
		if memval.Units != "kB" {
			memval.ActualIntMemsize = memval.ActualIntMemsize * 1024
			if memval.Units != "MB" {
				memval.ActualIntMemsize = memval.ActualIntMemsize * 1024
				if memval.Units != "GB" {
					memval.ActualIntMemsize = memval.ActualIntMemsize * 1024
				}
			}
		}
	}
}

func (memval *MemValue) optimise() {
	for memval.OptNumberPartAsInt >= 1024 && memval.Units != "TB" {
		memval.OptNumberPartAsInt = memval.OptNumberPartAsInt / 1024
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
	result.OrigNumberPartAsString = re1.FindString(result.Original)
	result.Units = re2.FindString(result.Original)
	if result.OrigNumberPartAsString != "" && fmt.Sprintf("%s%s", result.OrigNumberPartAsString, result.Units) == result.Original {
		result.Units = strings.ToUpper(result.Units)
		if result.Units == "KB" {
			result.Units = "kB"
		}
		optNumberPartAsInt, err := strconv.Atoi(result.OrigNumberPartAsString)
		if err != nil {
			result.Error = errors.New(fmt.Sprintf("Failed to convert memory string (1) %s\n\t%s", result.Original, err))
		} else {
			result.OptNumberPartAsInt = optNumberPartAsInt
			result.optimise()
			result.memStrToInt()
			result.Converted = true
			result.CorrectOptimisedStrValue = fmt.Sprintf("%d%s", result.OptNumberPartAsInt, result.Units)
		}
	} else {
		result.Error = errors.New(fmt.Sprintf("Failed to convert memory string (2) [%s]", result.Original))
	}
	if result.Error != nil {
		log.Printf("%s memory string conversion failed: %s\n", result.Error, result.Error)
	}
}

func GetMaxMemoryOfContainerAsString() string {
	cgroup_mem_limit_fname := "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	_, err := os.Stat(cgroup_mem_limit_fname)
	if err == nil {
		cgroup_mem_limit, errRead := ioutil.ReadFile(cgroup_mem_limit_fname)
		if errRead == nil {
			cgroup_mem_limit_str := strings.Trim(string(cgroup_mem_limit), "\n")
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
		memValue.ActualIntMemsize < otherValue.ActualIntMemsize {
		return true
	}
	return false
}
