package selftest

import (
	"testing"
	"time"
)

var nowFunc = time.Now

func SetNow(t time.Time) {
	nowFunc = func() time.Time { return t }
}

func NewNow() time.Time {
	return nowFunc()
}

func TestNewNow(t *testing.T) {
	cases := []struct{ t time.Time }{
		{time.Date(2018, time.January, 10, 18, 34, 32, 30, time.UTC)},
		{time.Date(2019, time.January, 10, 18, 34, 32, 30, time.UTC)},
	}

	for _, c := range cases {
		SetNow(c.t)
		now := NewNow()

		if c.t != now {
			t.Errorf("got %s, expected %s", now, c.t)
		}
	}

}
