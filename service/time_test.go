package service

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func newDate(d, m, y int) time.Time {
	return time.Date(y, time.Month(m), d, rand.Intn(24), rand.Intn(60), rand.Intn(60), 0, &time.Location{})
}

func TestFirstAndLastWeekDay(t *testing.T) {
	// fday and lday are first week day and last week day
	var fday, lday = 13, 19
	tasks := []time.Time{
		newDate(13, 04, 2020),
		newDate(14, 04, 2020),
		newDate(15, 04, 2020),
		newDate(16, 04, 2020),
		newDate(17, 04, 2020),
		newDate(18, 04, 2020),
		newDate(19, 04, 2020),
	}

	for _, task := range tasks {
		assert.Equal(t, fday, firstWeekDay(task).Day())
		assert.Equal(t, lday, lastWeekDay(task).Day())
	}
}

func TestTimeListName(t *testing.T) {
	tasks := []struct {
		t    time.Time
		name string
	}{
		{
			t:    newDate(14, 04, 2020),
			name: "13.04.2020-19.04.2020",
		},
		{
			t:    newDate(22, 05, 2020),
			name: "18.05.2020-24.05.2020",
		},
	}

	for _, task := range tasks {
		assert.Equal(t, task.name, TimeListName(task.t))
	}
}
