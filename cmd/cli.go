package main

import (
	"fmt"
	"github.com/intc/go-pwgen/pkg/pwgen"
	"os"
	s "strconv"
)

func main() {
	var err error
	n := 1
	l := 13
	// args [pw_length] [pw_number]
	args := len(os.Args)
	if args > 1 {
		l, err = s.Atoi(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if args > 2 {
			n, err = s.Atoi(os.Args[2])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
	// use NoIlSet
	pwgen.ActivateNoIlSet()
	// print elements
	// pwgen.PrintElements()
	// run
	for n > 0 {
		str, _ := pwgen.PhonemeGen(l)
		fmt.Printf("%s\n", *str)
		n = n - 1
	}
}
