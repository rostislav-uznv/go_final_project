package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("repeat rule is empty string")
	}

	date, err := time.Parse(layout, dstart)
	if err != nil {
		return "", fmt.Errorf("parse time error: %v", err)
	}

	repeatParts := strings.Split(repeat, " ")

	interval := 365
	switch repeatParts[0] {
	case "d":
		if len(repeatParts) == 1 {
			return "", errors.New("days number didn't set")
		}
		interval, err = strconv.Atoi(repeatParts[1])
		if err != nil {
			return "", fmt.Errorf("convert day number into a string error: %v", err)
		}
		if (interval < 1) || (interval > 400) {
			return "", errors.New("days wrong number")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	default:
		return "", errors.New("wrong repeat format")
	}

	return date.Format(layout), nil
}

func afterNow(date, now time.Time) bool {

	if date.Format(layout) > now.Format(layout) {
		return true
	}
	return false
}
