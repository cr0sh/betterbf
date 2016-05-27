package betterbf

import (
	"fmt"
	"strconv"
)

type tokenID uint

/* Token lists
add
sub
padd
psub
prst
loop
end
prt
scn
_set
_pset
_if
_endif
_goto
_exit
_rst
_snd
_chr
*/

type codegen interface {
	gen(*stokens) (string, error)
}

type stokens struct {
	tslice []string
	off    uint
}

func (t *stokens) slice_until(search string) []string {
	tt := make([]string, 0)
	for t.tslice[t.off] != search {
		tt = append(tt, t.tslice[t.off])
		t.off++
	}
	return tt
}

func (t *stokens) get() string {
	t.off += 1
	return t.tslice[t.off-1]
}

func (t *stokens) getint() int {
	s := t.get()
	n, err := strconv.Atoi(s)
	if err != nil {
		panic("strconv.Atoi conversion failed")
	}
	return int(n)
}

func (t *stokens) get_n(n uint) []string {
	tt := t.tslice[t.off : t.off+n]
	t.off += n
	return tt
}

type routine struct {
	ID   uint
	code []codegen
}

type tokenError struct {
	line    uint
	message string
}

func (e tokenError) Error() string {
	return fmt.Sprintf("%s (line %d)", e.message, e.line)
}
