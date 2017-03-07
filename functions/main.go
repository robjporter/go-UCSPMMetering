package functions

import (
	"strconv"
	"strings"
	"time"

	"github.com/robjporter/go-functions/as"
)

func CurrentMonthName() string {
	return time.Now().Month().String()
}

func IsYear(input string) string {
	if isNumber(input) {
		if isValidYear(input) {
			return input
		}
	}
	return CurrentYear()
}

func isValidYear(input string) bool {
	if result, err := strconv.ParseInt(input, 10, 64); err == nil {
		if result > 1999 && result < 2031 {
			return true
		}
	}
	return false
}

func isNumber(input string) bool {
	if _, err := strconv.ParseInt(input, 10, 64); err == nil {
		return true
	}
	return false
}

func IsMonth(input string) string {
	if monthContains(input, "jan", "january") {
		return "January"
	}
	if monthContains(input, "feb", "february") {
		return "February"
	}
	if monthContains(input, "mar", "march") {
		return "March"
	}
	if monthContains(input, "apr", "april") {
		return "April"
	}
	if monthContains(input, "may", "may") {
		return "May"
	}
	if monthContains(input, "jun", "june") {
		return "June"
	}
	if monthContains(input, "jul", "july") {
		return "July"
	}
	if monthContains(input, "aug", "august") {
		return "August"
	}
	if monthContains(input, "sep", "september") {
		return "September"
	}
	if monthContains(input, "oct", "october") {
		return "October"
	}
	if monthContains(input, "nov", "november") {
		return "November"
	}
	if monthContains(input, "dec", "december") {
		return "December"
	}
	return ""
}

func monthContains(input string, start string, end string) bool {
	input = strings.ToLower(input)
	if input == start || input == end {
		return true
	}
	pos := len(input)
	if pos <= len(end) {
		part := end[:pos]
		if input == part {
			return true
		}
	} else {
		return false
	}

	return false
}

func CurrentYear() string {
	return as.ToString(time.Now().Year())
}
