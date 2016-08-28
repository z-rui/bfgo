//go:generate go tool yacc parse.y

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <bf source>\n\tRun the bf program.\n", os.Args[0])
		return
	}
	input, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	ls := newLexState(input)
	yyParse(ls)
	prog := Program{}
	prog.Compile(ls.ast)
	/*for _, inst := range prog {
		fmt.Println(inst.Unpack())
	}*/
	println(prog.Len())
	VMRun(prog)
}
