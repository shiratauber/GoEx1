package Analyzer

import "os"

type Analyzer struct {
	outputFile  string
	file        *os.File
	current     string
	CallCounter int
	funcName    string
}

func New(path string) Analyzer {
	return Analyzer{}
}

func call() {
	var stamm string = " "
	println(stamm)
}
