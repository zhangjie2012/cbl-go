package cbl

import (
	"fmt"
	"math"
	"time"
)

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second)-1, t.Location())
}

// StartOfMonth month start datetime
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth month end datetime
func EndOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 999999999, t.Location())
}

// RangeMonth get month every day date
func RangeMonth(t time.Time) []time.Time {
	s := StartOfMonth(t)

	dates := []time.Time{}
	for d := s; d.Month() == s.Month(); d = d.AddDate(0, 0, 1) {
		t := time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())
		dates = append(dates, t)
	}
	return dates
}

func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year()+1, time.January, 0, 23, 59, 59, 999999999, t.Location())
}

func RangeYear(t time.Time) []time.Time {
	s := StartOfYear(t)

	dates := []time.Time{}
	for d := s; d.Year() == s.Year(); d = d.AddDate(0, 0, 1) {
		t := time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())
		dates = append(dates, t)
	}
	return dates
}

func FormatDayTime(t time.Time, showSec bool) string {
	h := t.Hour()

	desc := ""
	for {
		if h >= 5 && h < 8 {
			desc = "早晨"
			break
		}
		if h >= 8 && h < 12 {
			desc = "上午"
			break
		}
		if h >= 12 && h < 13 {
			desc = "中午"
			break
		}
		if h >= 13 && h < 18 {
			desc = "下午"
			break
		}
		if h >= 18 && h < 19 {
			desc = "傍晚"
			break
		}
		if h >= 19 && h < 23 {
			desc = "晚上"
			break
		}
		if (h >= 23 && h < 24) || (h >= 0 && h < 3) {
			desc = "深夜"
			break
		}
		if h >= 3 && h < 5 {
			desc = "凌晨"
			break
		}
		break
	}

	s := ""
	if showSec {
		s = fmt.Sprintf("%s %02d:%02d:%02d", desc, t.Hour(), t.Minute(), t.Second())
	} else {
		s = fmt.Sprintf("%s %02d:%02d", desc, t.Hour(), t.Minute())
	}
	return s
}

func DayOfWeekCN(n int) string {
	switch n {
	case 0:
		return "星期日"
	case 1:
		return "星期一"
	case 2:
		return "星期二"
	case 3:
		return "星期三"
	case 4:
		return "星期四"
	case 5:
		return "星期五"
	case 6:
		return "星期六"
	default:
		return "未知"
	}
}

func DurationString(d time.Duration) string {
	sec := d.Seconds()

	var interval float64
	interval = math.Floor(sec / 31536000)
	if interval > 1 {
		return fmt.Sprintf("%d 年", int64(interval))
	}

	interval = math.Floor(sec / 2592000)
	if interval > 1 {
		return fmt.Sprintf("%d 月", int64(interval))
	}

	interval = math.Floor(sec / 86400)
	if interval > 1 {
		return fmt.Sprintf("%d 天", int64(interval))
	}

	interval = math.Floor(sec / 3600)
	if interval > 1 {
		return fmt.Sprintf("%d 小时", int64(interval))
	}

	interval = math.Floor(sec / 60)
	if interval > 1 {
		return fmt.Sprintf("%d 分钟", int64(interval))
	}

	interval = math.Floor(sec)
	return fmt.Sprintf("%d 秒", int64(interval))
}
