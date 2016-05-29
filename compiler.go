package betterbf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Compile compiles the given code.
func Compile(code string) (string, error) {
	code = comment.ReplaceAllString(code, "")
	// one-sled boilerplate
	// 1 exit-call + 1 routineid
	ret := strings.Repeat(">>+", 16-1) + ">" + "\n" +
		strings.Repeat("+", 64) + "\n" +
		"[[->>+<<]>++>-]<[<<]>+>>+<<[\nloop started\n"
	strtokens := regexp.MustCompile(`\s+`).Split(code, -1)
	stoks := &stokens{tslice: strtokens, off: 0}
	for stoks.off < uint(len(stoks.tslice)) {
		if s := stoks.get(); s == "" {
			continue
		} else if s != "routine" {
			return "", fmt.Errorf("routine decl expected, got %s", s)
		} else {
			n := stoks.get()
			snum, err := strconv.Atoi(n)
			if err != nil {
				return "", err
			}
			rtoks := &stokens{tslice: stoks.slice_until("endroutine"), off: 0}
			stoks.off++
			ret += fmt.Sprintf("routine %d selector\n", snum)
			ret += "<[<<]>>>>>+<<" + strings.Repeat("-", snum) + "[>>-]<[<<]>>>" +
				strings.Repeat("+", snum) + ">>[[-]" +
				fmt.Sprintf("\nroutine %d code start\n", snum)
			for rtoks.off < uint(len(rtoks.tslice)) {
				offbefore := rtoks.off
				c := compileOp(rtoks.get(), rtoks)
				offafter := rtoks.off
				debug := ""
				for _, st := range rtoks.tslice[offbefore:offafter] {
					debug += st + " "
				}
				ret += debug + c + "\n"
			}
			ret += fmt.Sprintf("routine %d code end\n", snum)
			ret += fmt.Sprintf("<[<<]>>>>>]\nroutine %d selector end\n", snum)
		}
	}
	ret += "end phase(exit_call check)\n<[<<]>]\n"
	return ret, nil
}

func compileOp(op string, stoks *stokens) (c string) {
	switch op {
	case "_add":
		n := stoks.getint()
		c += strings.Repeat("+", n)
	case "_sub":
		n := stoks.getint()
		c += strings.Repeat("-", n)
	case "_padd":
		n := stoks.getint()
		c += strings.Repeat(">", n)
	case "_psub":
		n := stoks.getint()
		c += strings.Repeat("<", n)
	case "_loop":
		c += "["
	case "_end":
		c += "]"
	case "trace":
		c += "!"
	case "add":
		n := stoks.getint()
		m := stoks.getint()
		c += pset(n)
		c += strings.Repeat("+", m)
	case "sub":
		n := stoks.getint()
		m := stoks.getint()
		c += pset(n)
		c += strings.Repeat("-", m)
	case "prt":
		n := stoks.getint()
		c += pset(n) + "."
	case "scn":
		n := stoks.getint()
		c += pset(n) + ","
	case "set":
		n := stoks.getint()
		m := stoks.getint()
		c += pset(n)
		c += "[-]"
		c += strings.Repeat("+", m)
	case "pset":
		n := stoks.getint()
		c += pset(n)
	case "snd":
		n := stoks.getint()
		m := stoks.getint()
		c += psetr(10) + "[-]" // clear r10
		c += pset(n) + "[-" + pset(m) + "+" +
			psetr(10) + "+" + pset(n) + "]" // #M += #N, r10 += #N, #N=0
		c += psetr(10) + "[-" + pset(n) + "+" + psetr(10) + "]"
		c += pset(m)
	case "chr":
		n := stoks.getint()
		m := stoks.get()[0]
		c += pset(n)
		c += "[-]" + strings.Repeat("+", int(m))
	case "if":
		n := stoks.getint()
		c += pset(n)
		c += "["
	case "endif":
		c += psetr(10) + "[-]]"
	case "loop":
		n := stoks.getint()
		c += pset(n)
		c += "["
	case "end":
		n := stoks.getint()
		c += pset(n)
		c += "]"
	case "exit":
		c += "<[<<]>[-]>>[-]"
	case "goto":
		n := stoks.getint()
		c += "<[<<]>>>[-]"
		c += strings.Repeat("+", n)
	default:
		panic("undefined operation: " + op)
	}
	return
}

func compileByArgs(op string, args ...string) string {
	return compileOp(op, &stokens{tslice: args})
}
