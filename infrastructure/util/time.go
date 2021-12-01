package util

import (
	"github.com/hako/durafmt"
	"strings"
	"time"
	"unicode"
)

func GetCurrentShanghaiTime() time.Time {
	return time.Now().In(time.FixedZone("UTC-8", 8*60*60))
}

func ToShanghaiTime(t time.Time) time.Time {
	return t.In(time.FixedZone("UTC-8", 8*60*60))
}

func GetTimeInHour(t time.Time, hour int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, t.Location())
}

func GetDurationFriendly(tl time.Time, tr time.Time) string {
	runes := []rune(durafmt.Parse(tl.Sub(tr)).LimitFirstN(2).String())

	sb := strings.Builder{}
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case unicode.IsDigit(r):
			sb.WriteRune(runes[i])
		case r == 'y':
			i += 3
			sb.WriteString("年")
		case r == 'w':
			i += 3
			sb.WriteString("週")
		case r == 'd':
			i += 2
			sb.WriteString("日")
		case r == 'h':
			i += 3
			sb.WriteString("時")
		case r == 'm' && i+2 < len(runes) && runes[i+2] == 'l':
			i += 10
			sb.WriteString("毫秒")
		case r == 'm':
			i += 5
			sb.WriteString("分")
		case r == 's' && i+1 < len(runes) && unicode.IsLetter(runes[i+1]):
			sb.WriteString("秒")
			i += 5
		}
	}
	return sb.String()
}
