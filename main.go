package main

/*
Shira Tauber 213936271
Shvut Lazare 213195977
exercise 1
group 43
*/

import (
	"GoEx1/Analyzer"
	"GoEx1/CodeWriter"
	"GoEx1/CompilationEngine"
	"GoEx1/Parser"
	"GoEx1/Tokenizer"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"strings"
)

func main() {
	//tokenizer()
	//analayzer()
	jackCompiler()
	//VmToHack()

}

func jackCompiler() {
	var path string
	fmt.Scanln(&path)
	files, err := ioutil.ReadDir(path)
	Check(err)
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".jack" {
			//token := Tokenizer.New(path + "\\" + f.Name())
			inputFile, err := os.Open(path + "\\" + f.Name())
			Check(err)

			var split []string = strings.Split(path+"\\"+f.Name(), ".")
			var first string = split[0]
			//create the output file
			outputFile, err := os.Create(first + ".vm")
			Check(err)
			//open the output file
			outputFile, err = os.OpenFile(outputFile.Name(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
			Check(err)

			coEn := CompilationEngine.New(inputFile, outputFile)
			_ = coEn
		}

	}

}

func analayzer(path string) {
	//var path string
	//fmt.Scanln(&path)
	files, err := ioutil.ReadDir(path)
	Check(err)
	for _, f := range files {
		//var split string = strings.Split(f.Name(), ".")[0]
		if filepath.Ext(f.Name()) == ".xml" && filepath.Base(f.Name())[len(filepath.Base(f.Name()))-6:] == "TT.xml" {
			ana := Analyzer.New(path + "\\" + f.Name())
			scanner := bufio.NewScanner(ana.InputFile)
			Analyzer.WriteClass(scanner, ana)
			Analyzer.Close(ana)
		}

	}

}

func tokenizer() {
	var path string
	fmt.Scanln(&path)
	files, err := ioutil.ReadDir(path)
	Check(err)
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".jack" {
			token := Tokenizer.New(path + "\\" + f.Name())
			var s string = "<tokens>" + "\n"
			if _, err := token.OutputFile.WriteString(s); err != nil {
				panic(err)
			}
			scanner := bufio.NewScanner(token.InputFile)
			Tokenizer.Advance(scanner, token)
			s = "</tokens>" + "\n"
			if _, err := token.OutputFile.WriteString(s); err != nil {
				panic(err)
			}
			Tokenizer.Close(token)
		}

	}
	analayzer(path)
}

func VmToHack() {
	var path string
	fmt.Scanln(&path)
	code := CodeWriter.New(path)
	var numVm int = 0
	files, err := ioutil.ReadDir(path)
	Check(err)

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".vm" {
			numVm += 1
		}
	}

	if numVm == 0 {
		return
	}

	if numVm > 1 {
		CodeWriter.WriteInit(true, &code)
		code.CallCounter++
	} else {
		//	CodeWriter.WriteInit(false, &code)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".vm" {
			pars := parser.New(path + "\\" + f.Name())
			scanner := bufio.NewScanner(pars.File)
			for true {
				if parser.HasMoreCommands(&pars, scanner) {
					parser.Advance(scanner, &pars)
					WriteByCommand(pars, &code)
				} else {
					parser.Close(pars)
					break
				}
			}
		}

	}
	CodeWriter.Close(code)
}

func WriteByCommand(pars parser.Parser, code *CodeWriter.CodeWriter) {
	if parser.CommandType(pars) == "C_ARITHMETIC" {
		CodeWriter.WriteArithmetic(pars.Current, *code, strconv.Itoa(pars.LineNumber))
	} else if parser.CommandType(pars) == "C_PUSH" || parser.CommandType(pars) == "C_POP" {
		CodeWriter.WritePushPop(parser.CommandType(pars), parser.Arg1(pars), parser.Arg2(pars), *code, pars.FileName)
	} else if parser.CommandType(pars) == "C_LABEL" {
		CodeWriter.WriteLabel(parser.Arg1(pars), *code)
	} else if parser.CommandType(pars) == "C_GOTO" {
		CodeWriter.WriteGoTo(parser.Arg1(pars), *code)
	} else if parser.CommandType(pars) == "C_IF" {
		CodeWriter.WriteIf(parser.Arg1(pars), *code)
	} else if parser.CommandType(pars) == "C_RETURN" {
		CodeWriter.WriteReturn(*code)
	} else if parser.CommandType(pars) == "C_FUNCTION" {
		CodeWriter.WriteFunction(parser.Arg1(pars), strconv.Itoa(parser.Arg2(pars)), *code)
	} else if parser.CommandType(pars) == "C_CALL" {
		CodeWriter.WriteCall(parser.Arg1(pars), strconv.Itoa(parser.Arg2(pars)), *code)
		code.CallCounter++

	}
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
