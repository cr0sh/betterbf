// +build ignore

package main

import (
	"fmt"
	"io/ioutil"

	"github.com/cr0sh/betterbf"
)

func main() {
	b, _ := ioutil.ReadFile("sample.bbf")
	fmt.Println(betterbf.Compile(string(b)))
}
