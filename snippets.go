package betterbf

import (
	"fmt"
	"strings"
)

const prst = "<[<<]"

func pset(n int) string {
	return prst + ">" + strings.Repeat(">>", 16+n) + fmt.Sprintf("(pset %d)", n)
}

func psetr(n int) string { // n: 1~15
	return prst + ">" + strings.Repeat(">>", n)
}
