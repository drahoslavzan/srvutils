package dev

import "strings"

func Assert(condition bool, errs ...string) {
	if condition {
		return
	}

	err := "assertion failed"
	if len(errs) > 0 {
		err = strings.Join(errs, ": ")
	}

	panic(err)
}
