package main

import (
	Calc "calc"
	Uniq "uniq"
)

func main() {
	uniq := true
	if uniq {
		Uniq.ParseFile()
	} else {
		Calc.Calc()
	}
}
