// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cr0sh/betterbf"
)

const usage = `BetterBF compiler v1.1
usage: %s <filename>
`

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(usage, os.Args[0])
		return
	}
	ext := filepath.Ext(os.Args[1])
	if ext != ".bbf" {
		fmt.Println("Invalid extension", ext)
		return
	}
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	compiled, err := betterbf.Compile(string(b))
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if err := ioutil.WriteFile(os.Args[1][:len(os.Args[1])-3]+"bf", []byte(compiled), 0644); err != nil {
		fmt.Println("write error:", err)
		return
	}
}
