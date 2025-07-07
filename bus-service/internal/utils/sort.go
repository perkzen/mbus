package utils

import (
	"sort"
)

type HasDepartureAt interface {
	GetDepartureAt() string
}

func SortByDepartureAtAsc[T HasDepartureAt](items []T) {
	sort.SliceStable(items, func(i, j int) bool {
		t1, err1 := ParseClock(items[i].GetDepartureAt())
		t2, err2 := ParseClock(items[j].GetDepartureAt())

		if err1 != nil || err2 != nil {
			return false
		}
		return t1.Before(t2)
	})
}
