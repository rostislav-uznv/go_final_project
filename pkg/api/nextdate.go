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
	case "w":
		if len(repeatParts) != 2 {
			return "", errors.New("days number didn't set")
		}
		weekDays := strings.Split(repeatParts[1], ",")
		var weekDaysMarked [7]bool

		for _, v := range weekDays {
			weekDay, err := strconv.Atoi(v)
			if err != nil {
				return "", fmt.Errorf("convert weekday to digit, error: %v", err)
			}
			if weekDay < 1 || weekDay > 7 {
				return "", errors.New("days wrong number")
			}
			if weekDay == 7 {
				weekDaysMarked[0] = true
			} else {
				weekDaysMarked[weekDay] = true
			}
		}
		for {
			date = date.AddDate(0, 0, 1)
			if afterNow(date, now) && weekDaysMarked[int(date.Weekday())] {
				break
			}
		}

	case "m":
		if len(repeatParts) == 1 {
			return "", errors.New("days number didn't set")
		}
		var days [32]bool
		var months [13]bool
		var lastDay bool
		var secondLastDay bool

		for _, v := range strings.Split(repeatParts[1], ",") {
			day, err := strconv.Atoi(v)
			if err != nil {
				return "", fmt.Errorf("convert day to digit, error: %v", err)
			}
			if day < -2 || day > 31 || day == 0 {
				return "", errors.New("days wrong number")
			}
			if day < 0 {
				if day == -1 {
					lastDay = true
				}
				if day == -2 {
					secondLastDay = true
				}
			} else {
				days[day] = true
			}
		}
		if len(repeatParts) == 2 {
			for i := range months {
				months[i] = true
			}
			months[0] = false
		} else if len(repeatParts) == 3 {
			for _, v := range strings.Split(repeatParts[2], ",") {
				month, err := strconv.Atoi(v)
				if err != nil {
					return "", fmt.Errorf("convert month to digit, error: %v", err)
				}
				if month < 1 || month > 12 {
					return "", errors.New("months wrong number")
				}
				months[month] = true
			}
		} else {
			return "", errors.New("days and months wrong input format")
		}

		//check if rule is correct
		maxDaysInMonth := 0
		monthLengths := [13]int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
		for i := range months {
			if months[i] {
				if maxDaysInMonth < monthLengths[i] {
					maxDaysInMonth = monthLengths[i]
				}
			}
		}
		correctRuleExists := false
		if lastDay || secondLastDay {
			correctRuleExists = true
		}
		for i := range days {
			if i <= maxDaysInMonth && days[i] {
				correctRuleExists = true
			}
		}

		if !correctRuleExists {
			return "", errors.New("days and months wrong number, there is no such day in the months")
		}

		for {
			date = date.AddDate(0, 0, 1)
			afterNowValue := afterNow(date, now)
			if afterNowValue && days[date.Day()] && months[int(date.Month())] {
				break
			}
			lastDayOfMonth := int(time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day())
			secondLastDayOfMonth := lastDayOfMonth - 1

			if afterNowValue && lastDayOfMonth == int(date.Day()) && lastDay && months[int(date.Month())] {
				break
			}
			if afterNowValue && secondLastDayOfMonth == int(date.Day()) && secondLastDay && months[int(date.Month())] {
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
