package types

import (
	"github.com/EscanBE/go-ienumerable/goe"
	"sort"
	"strings"
)

type SecurityReport struct {
	SecurityRecords []SecurityRecord
}

func (r *SecurityReport) Sort() {
	if len(r.SecurityRecords) < 1 {
		return
	}

	sort.Slice(r.SecurityRecords, func(i, j int) bool {
		r1 := r.SecurityRecords[i]
		r2 := r.SecurityRecords[j]

		if r1.Fatal != r2.Fatal {
			return r1.Fatal
		}

		cmp := strings.Compare(r1.Module, r2.Module)
		if cmp < 0 {
			return true
		}

		if cmp > 0 {
			return false
		}

		cmp = strings.Compare(r1.Content, r2.Content)
		return cmp < 0
	})
}

func (r *SecurityReport) Add(record SecurityRecord) {
	r.SecurityRecords = append(r.SecurityRecords, record)
}

func (r *SecurityReport) CountFatal() int {
	return goe.NewIEnumerable(r.SecurityRecords...).Where(func(v SecurityRecord) bool {
		return v.Fatal
	}).Count(nil)
}

func (r *SecurityReport) CountWarning() int {
	return goe.NewIEnumerable(r.SecurityRecords...).Where(func(v SecurityRecord) bool {
		return !v.Fatal
	}).Count(nil)
}
