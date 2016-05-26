package betterbf

import "strings"

func pset(n int) string {
	return "<[<<]" + strings.Repeat(">", 2*(n+32)+1)
}
