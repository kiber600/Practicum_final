package api

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	now := r.URL.Query().Get("now")
	dstart := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if now == "" {
		log.Println("Empty date now")
		http.Error(w, "Empty date now", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		log.Println("Empty date repeat")
		http.Error(w, "Empty date repeat", http.StatusBadRequest)
		return
	}

	if now == "today" {
		now = time.Now().Format(dateFormat)
	}

	nowTime, err := time.Parse(dateFormat, now)
	if err != nil {
		fmt.Printf("Error parsing time: %v\n", err)
		http.Error(w, "Invalid date now", http.StatusInternalServerError)
		return
	}

	nextDt, err := nextDate(nowTime, dstart, repeat)
	if err != nil {
		fmt.Printf("Error getting naxt date: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDt))

}

func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	if dstart == "" {
		dstart = now.Format(dateFormat)
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {

		return "", fmt.Errorf("Error parsing time in nextDate: %w\n", err)
	}

	splitRepeat := strings.Split(repeat, " ")

	var num int

	switch splitRepeat[0] {
	case "d":

		if len(splitRepeat) < 2 {
			return "", fmt.Errorf("Repeat value is not correct(days)")
		} else {
			num, err = strconv.Atoi(splitRepeat[1])
			if err != nil {

				return "", fmt.Errorf("Repeat value is not correct(days): %w\n", err)
			}
		}
		nextDate, err := addDay(now, date, num)
		if err != nil {
			return "", fmt.Errorf("Error in nextDate: %w\n", err)
		}
		return nextDate, nil

	case "y":
		num = 1
		nextDate, err := addYear(now, date, num)
		if err != nil {
			return "", fmt.Errorf("Error in nextDate: %w\n", err)
		}
		return nextDate, nil

	case "m":
		var days []string
		var resultDays []string
		var months []string
		if len(splitRepeat) < 2 {
			return "", fmt.Errorf("Repeat value is not correct(months)")
		}
		if len(splitRepeat) == 2 {
			days = strings.Split(splitRepeat[1], ",")

			for _, day := range days {
				num, err := strconv.Atoi(day)
				if err != nil {
					return "", fmt.Errorf("Repeat value is not correct(months): %w\n", err)
				}

				nextDate, err := nextMonth(now, date, num)
				if err != nil {
					return "", fmt.Errorf("Error in nextDate: %w\n", err)
				}
				resultDays = append(resultDays, nextDate)
			}

			sort.Strings(resultDays)
			nextDate := resultDays[0]
			return nextDate, nil
		}

		if len(splitRepeat) == 3 {
			var resultM []string
			days = strings.Split(splitRepeat[1], ",")
			months = strings.Split(splitRepeat[2], ",")
			for _, day := range days {
				num, err = strconv.Atoi(day)
				if err != nil {
					return "", fmt.Errorf("Repeat day value is not correct(months): %w\n", err)
				}

				for _, month := range months {
					mon, err := strconv.Atoi(month)
					if err != nil {
						return "", fmt.Errorf("Repeat months value is not correct(months): %w\n", err)
					}
					nextDate, err := certainMonths(now, date, mon, num)
					if err != nil {
						log.Println("Error type m 1 1")
						return "", nil
					}
					resultM = append(resultM, nextDate)
				}
			}
			sort.Strings(resultM)
			log.Println(resultM)
			nextDate := resultM[0]
			return nextDate, nil
		}
		return "", nil
	case "w":
		var weekDays []string

		if len(splitRepeat) < 2 {
			return "", fmt.Errorf("Repeat value is not correct(days)")
		}

		if len(splitRepeat[1]) == 0 {
			return "", fmt.Errorf("Repeat value is not correct (week)")
		}

		splitDate := strings.Split(splitRepeat[1], ",")

		for _, days := range splitDate {
			num, err = strconv.Atoi(days)
			if err != nil {
				return "", fmt.Errorf("Repeat  value is not correct(week): %w\n", err)
			}

			nextDate, err := weekDay(now, date, num)
			if err != nil {
				return "", fmt.Errorf("Error getting next date (week): %w\n", err)
			}
			weekDays = append(weekDays, nextDate)

		}
		if len(weekDays) == 0 {
			return "", fmt.Errorf("Repeat value is not correct (week)")
		}
		nextDate := weekDays[0]
		return nextDate, nil

	default:
		log.Println("invalid format of literal")
		return "", fmt.Errorf("invalid format of literal")
	}

}

func addDay(now time.Time, date time.Time, days int) (string, error) {

	if days > 400 {
		return "", fmt.Errorf("The maximum interval has been exceeded")
	}
	nextDate := date
	for {
		nextDate = nextDate.AddDate(0, 0, days)
		if afterNow(nextDate, now) {
			break
		}
	}
	return nextDate.Format(dateFormat), nil
}

func addYear(now time.Time, date time.Time, years int) (string, error) {

	nextDate := date
	for {
		nextDate = nextDate.AddDate(years, 0, 0)
		if afterNow(nextDate, now) {
			break
		}

	}
	return nextDate.Format(dateFormat), nil
}

func afterNow(nextDate time.Time, now time.Time) bool {
	return nextDate.After(now)
}

func weekDay(now time.Time, date time.Time, day int) (string, error) {
	const daysInWeek = 7
	if date.Before(now) {
		date = now
	}

	if day < 1 || day > 7 {
		return "", fmt.Errorf("invalid day of week")
	}

	days := (day - int(date.Weekday()) + daysInWeek) % daysInWeek
	if days == 0 {
		days = 7
	}

	nextD := date
	for {
		nextD = nextD.AddDate(0, 0, days)
		if afterNow(nextD, now) {
			break
		}
	}
	return nextD.Format(dateFormat), nil
}

func nextMonth(now time.Time, date time.Time, day int) (string, error) {
	months := map[string]int{
		"January":   31,
		"February":  28,
		"March":     31,
		"April":     30,
		"May":       31,
		"June":      30,
		"July":      31,
		"August":    31,
		"September": 30,
		"October":   31,
		"November":  30,
		"December":  31,
	}

	var nextDate time.Time

	if int(date.Year())%4 == 0 {
		months["February"] = 29
	} else {
		months["February"] = 28
	}

	if day < -2 {

		log.Println("Inncorected day of month repeat")
		return "", fmt.Errorf("Inncorected day of repeat")

	}
	if day > 31 {
		log.Println("Days more 32")
		return "", fmt.Errorf("Days more 32")
	}

	if date.Before(now) {
		date = now
	}

	if day < 0 {
		day = -day
		if date.Day() < months[date.Month().String()]-(day-1) {
			nextDate = time.Date(int(date.Year()), date.Month(), months[date.Month().String()]-(day-1), 0, 0, 0, 0, time.UTC)
			return nextDate.Format(dateFormat), nil
		}
		if date.Day() > months[date.Month().String()]-(day-1) {
			nextDate = time.Date(int(date.Year()), date.AddDate(0, 1, 0).Month(), months[date.Month().String()]-(day-1), 0, 0, 0, 0, time.UTC)
			return nextDate.Format(dateFormat), nil
		}
	}

	if date.Day() < day && day < months[date.Month().String()] {
		nextDate = time.Date(int(date.Year()), date.Month(), day, 0, 0, 0, 0, time.UTC)
		return nextDate.Format(dateFormat), nil
	}

	if date.Day() <= day && day > months[date.Month().String()] {
		nextDate = time.Date(int(date.Year()), date.AddDate(0, 1, 0).Month(), day, 0, 0, 0, 0, time.UTC)
		return nextDate.Format(dateFormat), nil
	}

	if date.Day() >= day && day < months[date.AddDate(0, 1, 0).Month().String()] {
		nextDate = time.Date(int(date.Year()), date.AddDate(0, 1, 0).Month(), day, 0, 0, 0, 0, time.UTC)
		return nextDate.Format(dateFormat), nil
	}

	return "", nil
}

func certainMonths(now time.Time, date time.Time, mon int, day int) (string, error) {

	var monthNames = map[int]string{
		1:  "January",
		2:  "February",
		3:  "March",
		4:  "April",
		5:  "May",
		6:  "June",
		7:  "July",
		8:  "August",
		9:  "September",
		10: "October",
		11: "November",
		12: "December",
	}

	months := map[string]int{
		"January":   31,
		"February":  28,
		"March":     31,
		"April":     30,
		"May":       31,
		"June":      30,
		"July":      31,
		"August":    31,
		"September": 30,
		"October":   31,
		"November":  30,
		"December":  31,
	}

	if int(date.Year())%4 == 0 {
		months["February"] = 29
	} else {
		months["February"] = 28
	}

	if date.Before(now) {
		date = now
	}

	if day < -2 {

		log.Println("Inncorected day of month repeat")
		return "", fmt.Errorf("Inncorected day of repeat")

	}

	if mon < 1 || mon > 12 {
		return "", fmt.Errorf("Incorrect of number months")
	}
	if day > months[monthNames[mon]] {
		log.Println(day, "-days in month more- ", monthNames[mon])
		return "", fmt.Errorf("days in month more")
	}
	if day < 0 {
		nextDate := time.Date(date.Year(), time.Month(mon), months[time.Month(mon).String()]-(day-1), 0, 0, 0, 0, time.UTC)
		return nextDate.Format(dateFormat), nil
	}
	if mon == int(date.Month()) {
		nextDate := time.Date(date.Year(), time.Month(mon+1), day, 0, 0, 0, 0, time.UTC)
		return nextDate.Format(dateFormat), nil
	}
	if mon < int(date.Month()) {
		nextDate := time.Date(date.AddDate(1, 0, 0).Year(), time.Month(mon), day, 0, 0, 0, 0, time.UTC)
		return nextDate.Format(dateFormat), nil
	}

	nextDate := time.Date(date.Year(), time.Month(mon), day, 0, 0, 0, 0, time.UTC)
	return nextDate.Format(dateFormat), nil
}
