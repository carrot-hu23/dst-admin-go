package utils

import (
	"time"
)

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func Truncate(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}

func Get_stamp_day(start_time, end_time time.Time) (args []int64) {

	t1 := Bod(start_time)
	t2 := Bod(end_time)

	sInt := t1.UnixMilli() - 8*60*60*1000
	eInt := t2.UnixMilli() - 8*60*60*1000

	args = append(args, sInt)
	for {
		sInt += 86400 * 1000
		if sInt > eInt {
			return
		}
		args = append(args, sInt)
	}
}

func Get_stamp_month(start_time, end_time time.Time) (args []int64) {

	t1 := Bod(start_time)
	t2 := Bod(end_time)

	sInt := t1.Unix()
	eInt := t2.Unix()
	args = append(args, sInt)
	for {
		sInt += 2592000000
		if sInt > eInt {
			return
		}
		args = append(args, sInt)
	}
}
