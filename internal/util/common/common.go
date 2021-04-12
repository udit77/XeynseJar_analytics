package common

import (
	"time"
)

func GetMonthBeforeTimeStamp(month int) int64 {
	t := time.Now()
	t2 := t.AddDate(0, (-1)*month, 0)
	return t2.Unix()
}

func GetDayBeforeTimeStamp(day int) int64 {
	t := time.Now()
	t2 := t.AddDate(0, 0, (-1)*day)
	return t2.Unix()
}

func GetTodayTimeStampByLocation() int64 {
	location, _ := time.LoadLocation("Asia/Kolkata")
	t1 := time.Now()
	return time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, location).Unix()
}

func GetTodayTimeStringByLocation() string {
	location, _ := time.LoadLocation("Asia/Kolkata")
	t1 := time.Now()
	return time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, location).Format("2006-01-02 15:04:05.999999+07:00")
}

func GetWeekTimeStringByLocation() string {
	location, _ := time.LoadLocation("Asia/Kolkata")
	t1 := time.Now()
	weekDay := int(t1.Weekday())
	if weekDay == 0 {
		return time.Date(t1.Year(), t1.Month(), t1.Day()-6, 0, 0, 0, 0, location).Format("2006-01-02 15:04:05.999999+07:00")
	}
	return time.Date(t1.Year(), t1.Month(), t1.Day()-(weekDay-1), 0, 0, 0, 0, location).Format("2006-01-02 15:04:05.999999+07:00")
}

func GetMonthTimeStringByLocation() string {
	location, _ := time.LoadLocation("Asia/Kolkata")
	t1 := time.Now()
	t2 := time.Date(t1.Year(), t1.Month(), 1, 0, 0, 0, 0, location).Format("2006-01-02 15:04:05.999999+07:00")
	return t2
}

func GetTimeByLocation(t1 time.Time) time.Time {
	location, _ := time.LoadLocation("Asia/Kolkata")
	return t1.In(location)
}

func GetStartTimeForDay(t1 time.Time) int64 {
	location, _ := time.LoadLocation("Asia/Kolkata")
	return time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, location).Unix()
}

func GetEndTimeForDay(t1 time.Time) int64 {
	location, _ := time.LoadLocation("Asia/Kolkata")
	return time.Date(t1.Year(), t1.Month(), t1.Day(), 23, 59, 59, 0, location).Unix()
}

func GetLastTimeWithDayOffset(t1 time.Time, offset int) int64 {
	location, _ := time.LoadLocation("Asia/Kolkata")
	return time.Date(t1.Year(), t1.Month(), t1.Day()+offset, 23, 0, 0, 0, location).Unix()
}

func GetTimeSeriesMapping(t1 time.Time) string {
	location, _ := time.LoadLocation("Asia/Kolkata")
	return time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour()+1, 0, 0, 0, location).Format(time.Kitchen)
}
