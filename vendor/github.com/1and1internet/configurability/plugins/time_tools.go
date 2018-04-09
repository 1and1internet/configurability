package plugins

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const msPerDay = 86400000
const msPerHour = 3600000
const msPerMin = 60000
const msPerSec = 1000

type TimeValue struct {
	TimeStr        string
	NumberAsString string
	Units          string
	Milliseconds   int
	Error          error
}

func GetTimeValue(timestr string) *TimeValue {
	timeval := TimeValue{
		TimeStr: timestr,
	}
	timeval.TimeStrToTimeValue(timestr)
	return &timeval
}

func (result *TimeValue) TimeStrToTimeValue(timestr string) {
	// Convert values like 50ms, 3s, 18min, 22h, 1d
	// to a millisecond value that can be compared to other
	// times that were expressed as a string
	result.TimeStr = strings.ToLower(timestr)
	result.Milliseconds = 0

	re1 := regexp.MustCompile("[-0-9]+")
	re2 := regexp.MustCompile("(ms|s|min|h|d)")
	result.NumberAsString = re1.FindString(result.TimeStr)
	result.Units = re2.FindString(result.TimeStr)

	if result.NumberAsString != "" && fmt.Sprintf("%s%s", result.NumberAsString, result.Units) == result.TimeStr {
		numberAsInt, err := strconv.Atoi(result.NumberAsString)
		if err != nil {
			result.Error = errors.New(fmt.Sprintf("Failed to convert time string (1): %s", result.TimeStr))
		} else {
			switch result.Units {
			case "":
				result.Milliseconds = numberAsInt
			case "ms":
				result.Milliseconds = numberAsInt
			case "s":
				result.Milliseconds = numberAsInt * msPerSec
			case "min":
				result.Milliseconds = numberAsInt * msPerMin
			case "h":
				result.Milliseconds = numberAsInt * msPerHour
			case "d":
				result.Milliseconds = numberAsInt * msPerDay
			}
		}
	} else {
		result.Error = errors.New(fmt.Sprintf("Failed to convert time string (2): %s", result.TimeStr))
	}
	if result.Error != nil {
		log.Printf("%s time string conversion failed (1): %s", result.TimeStr, result.Error)
	}
}

func (timeValue *TimeValue) LessThan(otherValue *TimeValue) bool {
	if timeValue.Milliseconds < otherValue.Milliseconds {
		return true
	}
	return false
}
