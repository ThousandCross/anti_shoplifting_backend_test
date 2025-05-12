package utils

import "time"

var layout = "2022/01/01 00:00:00"

func StringToTime(str string) time.Time {
	t, _ := time.Parse(layout, str)
	return t
}

func TimeToString(t time.Time) string {
	str := t.Format(layout)
	return str
}
