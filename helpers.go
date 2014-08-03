package mixpanel

import "time"

const (
	fmt_mptime = "2006-01-02T15:04:05"
)

func TimeToMPFmt(t time.Time) string {
	return t.Format(fmt_mptime)
}

func Now() string {
	return TimeToMPFmt(time.Now().UTC())
}
