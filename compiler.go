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
	case "add":
		n := stoks.getint()
		c += strings.Repeat("+", n)
	case "sub":
		n := stoks.getint()
		c += strings.Repeat("-", n)
	case "padd":
		n := stoks.getint()
		c += strings.Repeat(">", n)
	case "psub":
		n := stoks.getint()
		c += strings.Repeat("<", n)
	case "loop":
		c += "["
	case "end":
		c += "]"
	case "_trace":
		c += "!"
	case "_add":
		n := stoks.getint()
		m := stoks.getint()
		c += pset(n)
		c += strings.Repeat("+", m)
	case "_sub":
		n := stoks.getint()
		m := stoks.getint()
		c += pset(n)
		c += strings.Repeat("-", m)
	case "_prt":
		n := stoks.getint()
		c += pset(n) + "."
	case "_scn":
		n := stoks.getint()
		c += pset(n) + ","
	case "_set":
		n := stoks.getint()
		m := stoks.getint()
		c += pset(n)
		c += "[-]"
		c += strings.Repeat("+", m)
	case "_pset":
		n := stoks.getint()
		c += pset(n)
	case "_snd":
		n := stoks.getint()
		m := stoks.getint()
		c += pset(n) + "[" + pset(m) + "+" + "<[<<]>" + strings.Repeat(">>", 6) + "+" + pset(n) + "-]"       // move n to r7/m
		c += "<[<<]>" + strings.Repeat(">>", 6) + "[" + pset(n) + "+<[<<]>" + strings.Repeat(">>", 6) + "-]" // move r7 to n back
		c += pset(m)
	case "_chr":
		n := stoks.getint()
		m := stoks.get()[0]
		c += pset(n)
		c += "[-]" + strings.Repeat("+", int(m))
	case "_if":
		n := stoks.getint()
		c += pset(n)
		c += "["
	case "_endif":
		c += prst
		c += "]>"
	case "_loop":
		n := stoks.getint()
		c += pset(n)
		c += "["
	case "_end":
		n := stoks.getint()
		c += pset(n)
		c += "]"
	case "_exit":
		c += "<[<<]>[-]>>[-]"
	case "_goto":
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
