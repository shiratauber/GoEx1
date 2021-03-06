package SymbolTable

import (
	"GoEx1/Symbol"
)

type SymbolTable struct {
	classSymbols      map[string]Symbol.Symbol //STATIC, FIELD
	subroutineSymbols map[string]Symbol.Symbol //ARG, VAR
	indices           map[string]int
}

func New() SymbolTable {
	mofa := SymbolTable{make(map[string]Symbol.Symbol), make(map[string]Symbol.Symbol), make(map[string]int)}
	mofa.indices["ARG"] = 0
	mofa.indices["FIELD"] = 0
	mofa.indices["STATIC"] = 0
	mofa.indices["VAR"] = 0
	return mofa

}
func StartSubroutine(s *SymbolTable) {
	s.subroutineSymbols = make(map[string]Symbol.Symbol)
	s.indices["ARG"] = 0
	s.indices["VAR"] = 0
}

func Define(name string, t string, kind string, s *SymbolTable) {
	if kind == "ARG" || kind == "VAR" {
		var index = s.indices[kind]
		var sym = Symbol.New(t, kind, index)
		s.indices[kind] = index + 1
		s.subroutineSymbols[name] = sym

	} else if kind == "STATIC" || kind == "FIELD" {
		var index = s.indices[kind]
		var sym = Symbol.New(t, kind, index)
		s.indices[kind] = index + 1
		s.classSymbols[name] = sym
	}
}

//Returns the number of variables of the given kind already defined in the current scope.
func VarCount(kind string, s *SymbolTable) int {
	return s.indices[kind]
}

//Returns the kind of the named identifier in the current scope.
//If the identifier is unknown in the current scope returns NONE.
func KindOf(name string, s *SymbolTable) string {
	var sym = lookUp(name, s)
	if Symbol.GetIndex(sym) == -1 {
		return "NONE"
	} else {
		return Symbol.GetKind(sym)
	}
}

// Returns the type of the named identifier in the current scope.

func TypeOf(name string, s *SymbolTable) string {
	var sym = lookUp(name, s)
	if Symbol.GetIndex(sym) == -1 {
		return ""
	} else {
		return Symbol.GetType(sym)
	}
}

func IndexOf(name string, s *SymbolTable) int {
	var sym = lookUp(name, s)
	if Symbol.GetIndex(sym) == -1 {
		return -1
	} else {
		return Symbol.GetIndex(sym)
	}
}

func lookUp(name string, s *SymbolTable) Symbol.Symbol {
	var _, ok = s.subroutineSymbols[name]
	var _, ok2 = s.classSymbols[name]

	if ok == true {
		return s.subroutineSymbols[name]
	} else if ok2 == true {
		return s.classSymbols[name]
	} else {
		//print("SYMBOL NOT FOUND ON TABLES :", name+"\n")
		return Symbol.New("", "", -1)
	}

}
