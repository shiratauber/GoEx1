package CompilationEngine

import (
	"GoEx1/jackTokenizer"
	"GoEx1/VMWriter"
	"GoEx1/SymbolTable"
		"os"
	"strconv"
	"strings"
)
type CompilationEngine struct {
	tokenizer jackTokenizer.JackTokenizer
	vmWriter  VMWriter.Writer
	symbolTable  SymbolTable.SymbolTable
	currentClass string
	currentSubroutine string
	labelIndex int
	labelCounterIf int
	labelCounterWhile int
}



func New(inputFile *os.File, outputFile *os.File) CompilationEngine {


	mofa := CompilationEngine{jackTokenizer.New(inputFile), VMWriter.New(outputFile),SymbolTable.New(),",","",0,0,0}
	compileClass(&mofa)
	return mofa

}

func compileClass(c *CompilationEngine)  {
	jackTokenizer.Advance(&c.tokenizer)
	if jackTokenizer.TokenType(c.tokenizer) != "KEYWORD" || jackTokenizer.KeyWord(c.tokenizer) != "CLASS" {
		print("ERROR IN ADVANCE OF COMPILECLASS")
	}

	jackTokenizer.Advance(&c.tokenizer)
	if jackTokenizer.TokenType(c.tokenizer) != "IDENTIFIER"{
		print("EXPECTED CLASSNAME")
	}

    c.currentClass=jackTokenizer.Identifier(c.tokenizer)
    RequireSymbol("{", c)

	// classVarDec* subroutineDec*
	CompileClassVarDec(c)
	CompileSubroutine(c)

	RequireSymbol("}", c)              //  }

	if jackTokenizer.HasMoreTokens(&c.tokenizer) {
		print("Unexpected tokens!")
		return
	}

	VMWriter.Close(c.vmWriter)





}
func RequireSymbol(s string, c *CompilationEngine)  {
	jackTokenizer.Advance(&c.tokenizer)
	if jackTokenizer.TokenType(c.tokenizer) != "SYMBOL" || jackTokenizer.Symbol(c.tokenizer)!= s {
		print("ERROR EXPECTED SYMBOL")
	}

}


// Compiles a static declaration or a field declaration.
// classVarDec ('static'|'field') type varName (','varNAme)* ';'
func CompileClassVarDec(c *CompilationEngine) {
	// print("compileClassVarDec")
	// first determine whether there is a classVarDec, nextToken is } or start subroutineDec
	jackTokenizer.Advance(&c.tokenizer)

	// next is }
	if jackTokenizer.TokenType(c.tokenizer) == "SYMBOL" && jackTokenizer.Symbol(c.tokenizer) == "}" {
		jackTokenizer.PointerBack(&c.tokenizer)
		return
	}

	// next is start subroutineDec or classVarDec, both start with keyword
	if jackTokenizer.TokenType(c.tokenizer) != "KEYWORD" {
		print("Expected keyword")
	}

	// next is subroutineDec
	if jackTokenizer.KeyWord(c.tokenizer) == "CONSTRUCTOR"|| jackTokenizer.KeyWord(c.tokenizer) == "FUNCTION"|| jackTokenizer.KeyWord(c.tokenizer) == "METHOD" {
		jackTokenizer.PointerBack(&c.tokenizer)
		return
	}

	// classVarDec exists
	if !(jackTokenizer.KeyWord(c.tokenizer) =="STATIC" ||jackTokenizer.KeyWord(c.tokenizer) == "FIELD") {
		print("Expected static or field")
	}

	kind := jackTokenizer.KeyWord(c.tokenizer)
	typeTok := CompileType()

	for true{
			// varName
			jackTokenizer.Advance(&c.tokenizer)
			if jackTokenizer.TokenType(c.tokenizer) != "IDENTIFIER" {
			print("Expected identifier")
		}

		name := jackTokenizer.Identifier(c.tokenizer)
		SymbolTable.Define(name, typeTok, kind, &c.symbolTable)

		// , or ;
		jackTokenizer.Advance(&c.tokenizer)

		if jackTokenizer.TokenType(c.tokenizer) != "SYMBOL" || !(jackTokenizer.Symbol(c.tokenizer) ==","|| jackTokenizer.Symbol(c.tokenizer)==";" ) {
			print("Expected , or ;")
		}

		if jackTokenizer.Symbol(c.tokenizer) == ";" {
			break
		}

	}

	CompileClassVarDec(c)
}

// Compiles a complete method, function or constructor.
func CompileSubroutine(c *CompilationEngine) {
	// determine whether there is a subroutine, next can be a '}'
	jackTokenizer.Advance(&c.tokenizer)

	// next is a '}'
	if jackTokenizer.TokenType(c.tokenizer) == "SYMBOL" && jackTokenizer.Symbol(c.tokenizer) == "}" {
		jackTokenizer.PointerBack(&c.tokenizer)
		return
	}

	// start of a subroutine
	if jackTokenizer.TokenType(c.tokenizer) != "KEYWORD" || !(jackTokenizer.KeyWord(c.tokenizer)=="CONSTRUCTOR" ||
		jackTokenizer.KeyWord(c.tokenizer) =="FUNCTION"||jackTokenizer.KeyWord(c.tokenizer)== "METHOD") {
		print("Expected constructor or function or method")
	}

	keyword := jackTokenizer.KeyWord(c.tokenizer)
	SymbolTable.StartSubroutine(&c.symbolTable)

	// for method this is the first argument
	if keyword == "METHOD" {
		SymbolTable.Define("this", c.currentClass, "ARG", &c.symbolTable)
	}

	typeTok := ""

	// 'void' or typeTok
	jackTokenizer.Advance(&c.tokenizer)
	if jackTokenizer.TokenType(c.tokenizer) == "KEYWORD" && jackTokenizer.KeyWord(c.tokenizer) == "VOID" {
	typeTok  = "void"
	} else {
		jackTokenizer.PointerBack(&c.tokenizer)
	typeTok  = CompileType()
	}

	// subroutineName which is a identifier
	jackTokenizer.Advance(&c.tokenizer)
	if jackTokenizer.TokenType(c.tokenizer) != "IDENTIFIER" {
		print("Expected subroutineName")
	}

	c.currentSubroutine = jackTokenizer.Identifier(c.tokenizer)

	// '('
	RequireSymbol("(", c)

	// parameterList
	CompileParameterList()

	// ')'
	RequireSymbol(")", c)

	// subroutineBody
	CompileSubroutineBody(keyword, c)

	CompileSubroutine(c)

}

// Compiles the body of a subroutine.
// '{'  varDec* statements '}'
func CompileSubroutineBody(keyword string, c *CompilationEngine) {
	// '{'
	RequireSymbol("{", c)
	// varDec*
	CompileVarDec()
	// write VM function declaration
	WriteFunctionDec(keyword, c)
	// statements
	CompileStatement()
	// '}'
	RequireSymbol("}", c)
}

// Writes function declaration, load pointer when keyword is METHOD or CONSTRUCTOR.
func WriteFunctionDec(keyword string, c *CompilationEngine) {
	VMWriter.WriteFunction(CurrentFunction(), SymbolTable.VarCount("VAR", &c.symbolTable), c.vmWriter)

	// METHOD and CONSTRUCTOR need to load this pointer
	if keyword == "METHOD" {
		// A Jack method with k arguments is compiled into a VM function that operates on k + 1 arguments.
		// The first argument (argument number 0) always refers to the this object.
		VMWriter.WritePush("argument", 0, c.vmWriter)
		VMWriter.WritePop("pointer", 0, c.vmWriter)
	} else if keyword == "CONSTRUCTOR" {
		// A Jack function or constructor with k arguments is compiled into a VM function that operates on k arguments.
		VMWriter.WritePush("constant", SymbolTable.VarCount("FIELD", &c.symbolTable), c.vmWriter)
		VMWriter.WriteCall("Memory.alloc", 1, c.vmWriter)
		VMWriter.WritePop("pointer", 0, c.vmWriter)
	}

}
func CompileStatement(c *CompilationEngine) {
	// determine whether there is a statement next can be a '}'
	jackTokenizer.Advance(&c.tokenizer)
	// next is a '}'
	if jackTokenizer.TokenType(c.tokenizer) == "SYMBOL" && jackTokenizer.Symbol(c.tokenizer) == "}"{
		jackTokenizer.PointerBack(&c.tokenizer)
		return
	}
	// next is 'let'|'if'|'while'|'do'|'return'
	if jackTokenizer.TokenType(c.tokenizer) != "KEYWORD" {
		print("Expected keyword")
	} else {
		switch jackTokenizer.KeyWord(c.tokenizer) {
		case "LET":
			{
				CompileLet()
				break
			}
		case "IF":
			{
				CompileIf()
				break
			}
		case "WHILE":
			{
				CompileWhile()
				break
			}
		case "DO" :
			{
				CompileDo()
				break
			}
		case "RETURN":
			{
				CompileReturn()
				break
			}
		default:
			{ //   default
				print("Expected let or if or while or do or return")
			}
		}
	}

	CompileStatement(c)
}
// Compiles a (possibly empty) parameter list,
// not including the enclosing "()".
// ((type varName)(',' type varName)*)?
func CompileParameterList(c *CompilationEngine) {
	// Check if there is parameterList, if next token is ')' than go back
	jackTokenizer.Advance(&c.tokenizer)
	if jackTokenizer.TokenType(c.tokenizer) == "SYMBOL" && jackTokenizer.Symbol(c.tokenizer) == ")"{
		jackTokenizer.PointerBack(&c.tokenizer)
		return
	}
	// there is parameter, at least one varName
	jackTokenizer.PointerBack(&c.tokenizer)
	for true{
		// typeTok
		typeTok  := CompileType()
		// varName
		jackTokenizer.Advance(&c.tokenizer)
		if jackTokenizer.TokenType(c.tokenizer) != "IDENTIFIER" {
			print("Expected identifier")
		}
		SymbolTable.Define(jackTokenizer.Identifier(c.tokenizer), typeTok , "ARG", &c.symbolTable)
		// ',' or ')'
		jackTokenizer.Advance(&c.tokenizer)
		if jackTokenizer.TokenType(c.tokenizer) != "SYMBOL" || !(jackTokenizer.Symbol(c.tokenizer) == "," ||
			jackTokenizer.Symbol(c.tokenizer) == ")") {
			print("Expected , or )")
		}
		if jackTokenizer.Symbol(c.tokenizer) == ")" {
			jackTokenizer.PointerBack(&c.tokenizer)
			break
		}
	}
}


func CompileVarDec(c *CompilationEngine) {
	// determine if there is a varDec
	jackTokenizer.Advance(&c.tokenizer)
	// no 'var' go back
	if jackTokenizer.TokenType(c.tokenizer) != "KEYWORD" || jackTokenizer.KeyWord(c.tokenizer) != "VAR"{
		jackTokenizer.PointerBack(&c.tokenizer)
		return
	}
	// typeTok
	typeTok  <- self$compileType()
	repeat{
		// varName
		self$tokenizer$advance()
		if (self$tokenizer$tokenType() != "IDENTIFIER") {
		self$throwException("Expected identifier")
		}
		self$symbolTable$define(self$tokenizer$identifier(), typeTok , "VAR")
		// ',' or ';'
		self$tokenizer$advance()
		if (self$tokenizer$tokenType() != "SYMBOL" | !(self$tokenizer$symbol() %in% c(",", ";"))) {
		self$throwException("Expected , or ;")
		}
		if (self$tokenizer$symbol() == ';') {
		break
	}
	}
	self$compileVarDec()
},



//compileStatements = function() {},

// Compiles a do statement.
// 'do' subroutineCall ';'
compileDo = function() {
// subroutineCall
self$compileSubroutineCall()
// ';'
self$requireSymbol(';')
// pop return value
self$vmWriter$writePop("temp", 0)
},

// Compiles a subroutine call.
// subroutineName '(' expressionList ')' | (className|varName) '.' subroutineName '(' expressionList ')'
compileSubroutineCall = function() {
self$tokenizer$advance()
if (self$tokenizer$tokenType() != "IDENTIFIER"){
self$throwException("Expected identifier")
}

name <- self$tokenizer$identifier()
nArgs <- 0

self$tokenizer$advance()
if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$symbol() == '(') {
// push this pointer
self$vmWriter$writePush("pointer", 0)
// '(' expressionList ')'
// expressionList
nArgs <- self$compileExpressionList() + 1
// ')'
self$requireSymbol(')')
// call subroutine
self$vmWriter$writeCall(paste(self$currentClass, '.', name, sep=""), nArgs)
} else if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$symbol() == '.') {
// (className|varName) '.' subroutineName '(' expressionList ')'

objName <- name
// subroutineName
self$tokenizer$advance()

if (self$tokenizer$tokenType() != "IDENTIFIER"){
self$throwException("Expected identifier")
}

name <- self$tokenizer$identifier()

// check for if it is built-in typeTok
typeTok  <- self$symbolTable$typeOf(objName)

if (typeTok  %in% c("int", "boolean", "char", "void")) {
self$throwException("No built-in type")
} else if (typeTok  == "") {
name <- paste(objName, ".", name, sep="")
} else {
nArgs <- 1
// push variable directly onto stack
self$vmWriter$writePush(self$getSeg(self$symbolTable$kindOf(objName)), self$symbolTable$indexOf(objName))
name <- paste(self$symbolTable$typeOf(objName), ".", name, sep="")
}

// '('
self$requireSymbol('(')
// expressionList
nArgs <- nArgs + self$compileExpressionList()
 ')'
self$requireSymbol(')')
 call subroutine
self$vmWriter$writeCall(name, nArgs)

} else {
self$throwException("Expected ( or .")
}
},

 Compiles a let statement
 'let' varName ('[' ']')? '=' expression ';'
compileLet = function() {   // let diff = y - x;
// varName
self$tokenizer$advance()
if (self$tokenizer$tokenType() != "IDENTIFIER") {
self$throwException("Expected varName")
}

varName <- self$tokenizer$identifier()
// print(paste("VARNAME :", varName))

// '[' or '='
self$tokenizer$advance()
if (self$tokenizer$tokenType() != "SYMBOL" | !(self$tokenizer$symbol() %in% c("[", "="))) {
self$throwException("Expected [ or =")
}

expExist <- FALSE

// '[' expression ']' , need to deal with array [base + offset]
if (self$tokenizer$symbol() == '[') {
expExist <- TRUE

// calc offset
self$compileExpression()

// ']'
self$requireSymbol(']')

// push array variable,base address into stack
self$vmWriter$writePush(self$getSeg(self$symbolTable$kindOf(varName)), self$symbolTable$indexOf(varName))

// base + offset
self$vmWriter$writeArithmetic("add")
}

if (expExist == TRUE) {
self$tokenizer$advance()
}

// expression
self$compileExpression()

// ';'
self$requireSymbol(';')

if (expExist == TRUE) {
// *(base + offset) = expression
// pop expression value to temp
self$vmWriter$writePop("temp", 0)
// pop base + index into 'that'
// self$vmWriter$writePop("pointer", 0)
self$vmWriter$writePop("pointer", 1)
// pop expression value into *(base + index)
self$vmWriter$writePush("temp", 0)
self$vmWriter$writePop("that", 0)
//print(paste("VARNEME :", varName))
} else {
// pop expression value directly
self$vmWriter$writePop(self$getSeg(self$symbolTable$kindOf(varName)), self$symbolTable$indexOf(varName))
// print(paste("VARNEME :", varName))
}

},

// Returns corresponding segment for input kind.
func getSeg(kind string) string{
	switch kind {
	case "FIELD":{
		return "this"
	}
	case "STATIC":{
		return "static"
	}
	case "VAR":{
		return "local"
	}
	case "ARG":{
		return "argument"
	}
	default:
		return "NONE"
	}
}

// Compiles a while statement.
// 'while' '(' expression ')' '{' statements '}'
func compileWhile(c *CompilationEngine) {
	whileExpLabel := "WHILE_EXP"+ string(c.labelCounterWhile)
	whileEndLabel := "WHILE_END"+ string(c.labelCounterWhile)
	c.labelCounterWhile = c.labelCounterWhile + 1

	// top label for while loop
	VMWriter.WriteLabel(whileExpLabel, c.vmWriter)

	// '('
	RequireSymbol("(", c)
	// expression while condition: true or false
	CompileExpression()
	// ')'
	RequireSymbol(")", c)

	// if ~(condition) go to continue label
	VMWriter.WriteArithmetic("not", c.vmWriter)
	VMWriter.WriteIf(whileEndLabel, c.vmWriter)

	// '{'
	RequireSymbol("{", c)
	// statements
	CompileStatement()
	// '}'
	RequireSymbol("}", c)

	// if (condition) go to top label
	VMWriter.WriteGoto(whileExpLabel, c.vmWriter)
	// or continue
	VMWriter.WriteLabel(whileEndLabel, c.vmWriter)

	// self$labelIndex <- self$labelIndex + 1
}

// Compiles a return statement.
// ‘return’ expression? ';'
func CompileReturn(c *CompilationEngine) {
	// check if there is any expression
	jackTokenizer.Advance(&c.tokenizer)

	if jackTokenizer.TokenType(c.tokenizer) == "SYMBOL" && jackTokenizer.Symbol(c.tokenizer) == ";" {
		// no expression push 0 to stack
		VMWriter.WritePush("constant", 0, c.vmWriter)
	}else {
		// expression exist
		jackTokenizer.PointerBack(&c.tokenizer)
		// expression
		CompileExpression()
		// ';'
		RequireSymbol(";", c)
	}

	VMWriter.WriteReturn(c.vmWriter)
}

// Compiles an if statement,
// possibly with a trailing else clause.
// 'if' '(' expression ')' '{' statements '}' ('else' '{' statements '}')?
func compileIf(c *CompilationEngine) {
	ifTrueLabel := "IF_TRUE"+ string(c.labelCounterIf)
	ifFalseLabel := "IF_FALSE"+ string(c.labelCounterIf)
	ifEndLabel := "IF_END"+ string(c.labelCounterIf)
	c.labelCounterIf = c.labelCounterIf + 1

	// '('
	RequireSymbol("(",c)
	// expression
	CompileExpression()
	// ')'
	RequireSymbol(")",c)

	VMWriter.WriteIf(ifTrueLabel,c.vmWriter)
	VMWriter.WriteGoto(ifFalseLabel, c.vmWriter)
	VMWriter.WriteLabel(ifTrueLabel, c.vmWriter)

	// '{'
	RequireSymbol("{", c)
	// statements
	CompileStatement()
	// '}'
	RequireSymbol("}",c)

	// check if there is 'else'
	jackTokenizer.Advance(&c.tokenizer)
	if jackTokenizer.TokenType(c.tokenizer) == "KEYWORD" && jackTokenizer.KeyWord(c.tokenizer) == "ELSE" {
	// ifEndLabel <- paste("IF_END", self$labelIndex, sep="")
	self$vmWriter$writeGoto(ifEndLabel)
	self$vmWriter$writeLabel(ifFalseLabel)

	// '{'
	self$requireSymbol('{')
	// statements
	self$compileStatement()
	// '}'
	self$requireSymbol('}')

	self$vmWriter$writeLabel(ifEndLabel)
	}else {   ##   only if
	self$tokenizer$pointerBack()
	self$vmWriter$writeLabel(ifFalseLabel)
	}

	# self$labelIndex <- self$labelIndex + 1
},
# compileIf = function() {
#     elseLabel <- self$newLabel()
#     endLabel <- self$newLabel()

#     ## '('
#     self$requireSymbol('(')
#     ## expression
#     self$compileExpression()
#     ## ')'
#     self$requireSymbol(')')
#     ## if ~(condition) go to else label
#     self$vmWriter$writeArithmetic("not")
#     self$vmWriter$writeIf(elseLabel)
#     ## '{'
#     self$requireSymbol('{')
#     ## statements
#     self$compileStatement()
#     ## '}'
#     self$requireSymbol('}')
#     ## if condition after statement finishing, go to end label
#     self$vmWriter$writeGoto(endLabel)

#     self$vmWriter$writeLabel(elseLabel)
#     ## check if there is 'else'
#     self$tokenizer$advance()
#     if (self$tokenizer$tokenType() == "KEYWORD" & self$tokenizer$keyWord() == "ELSE") {
#         ## '{'
#         self$requireSymbol('{')
#         ## statements
#         self$compileStatement()
#         ## '}'
#         self$requireSymbol('}')
#     }else {
#         self$tokenizer$pointerBack()
#     }

#     self$vmWriter$writeLabel(endLabel)
# },

## Compiles an expression
## term (op term)*
compileExpression = function() {
## term
self$compileTerm()
## (op term)*
repeat{
self$tokenizer$advance()
## op
if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$isOp()) {

opCommand <- ""

switch(self$tokenizer$symbol(),
'+'={
opCommand <- "add"
},
'-'={
opCommand <- "sub"
},
'*'={
opCommand <- "call Math.multiply 2"
},
'/'={
opCommand <- "call Math.divide 2"
},
'<'={
opCommand <- "lt"
},
'>'={
opCommand <- "gt"
},
'='={
opCommand <- "eq"
},
'&'={
opCommand <- "and"
},
'|'={
opCommand <- "or"
},
{
self$throwException("Unknown op")
})

## term
self$compileTerm()

self$vmWriter$writeCommand(opCommand)

}else {
self$tokenizer$pointerBack()
break
}
}
},

/*## Compiles a term. This routine is faced with a
## slight difficulty when trying to decide
## between some of the alternative parsing rules.
## Specifically, if the current token is an
## identifier, the routine must distinguish
## between a variable, an array entry and a
## subroutine call. A single look-ahead token,
## which may be one of "[" "(" "."
## suffices to distinguish between the three
## possibilities. Any other token is not part of
## this term and should not be advanced over.
## integerConstant|stringConstant|keywordConstant|varName|varName '[' expression ']'|subroutineCall| '(' expression ')'|unaryOp term*/
func CompileTerm(c *CompilationEngine) {
	jackTokenizer.Advance(&c.tokenizer)
	// check if it is an identifier
	if jackTokenizer.TokenType(c.tokenizer) == "IDENTIFIER" {
		// varName|varName '[' expression ']'|subroutineCall
		var tempId = jackTokenizer.Identifier(c.tokenizer)

		jackTokenizer.Advance(&c.tokenizer)
		arr2 := []string{"(", "."}

		if jackTokenizer.TokenType(c.tokenizer) == "SYMBOL" && jackTokenizer.Symbol(c.tokenizer) == "[" {
			//this is an array entry
			//expression
			CompileExpression()
			// ']'
			RequireSymbol("]", c)
			//push array variable,base address into stack
			VMWriter.WritePush(getSeg(SymbolTable.KindOf(tempId, &c.symbolTable)), SymbolTable.IndexOf(tempId, &c.symbolTable),c.vmWriter)
			// base + offset
			VMWriter.WriteArithmetic("add",c.vmWriter)

			// pop into 'that' pointer
			VMWriter.WritePop("pointer",1,c.vmWriter)

			// push *(base+index) onto stack
			VMWriter.WritePush("that",0,c.vmWriter)

		}else if jackTokenizer.TokenType(c.tokenizer) == "SYMBOL" &&  StringInSlice(jackTokenizer.Symbol(c.tokenizer),arr2) {
			//this is a subroutineCall
			jackTokenizer.PointerBack(&c.tokenizer)
			jackTokenizer.PointerBack(&c.tokenizer)
			CompileSubroutineCall()
		} else {
			//this is varName
			jackTokenizer.PointerBack(&c.tokenizer)
			// push variable directly onto stack
			VMWriter.WritePush(getSeg(SymbolTable.KindOf(tempId,&c.symbolTable)),SymbolTable.IndexOf(tempId,&c.symbolTable),c.vmWriter )
		}

		} else {
			// integerConstant|stringConstant|keywordConstant|'(' expression ')'|unaryOp term
		arr3 := []string{"FALSE", "NULL"}

		arr4 := []string{"-", "~"}
		if jackTokenizer.TokenType(c.tokenizer) == "INT_CONST" {
				//integerConstant just push its value onto stack
				VMWriter.WritePush("constant", jackTokenizer.IntVal(c.tokenizer),c.vmWriter)
			}else if jackTokenizer.TokenType(c.tokenizer) == "STRING_CONST" {
				// stringConstant new a string and append every char to the new stack
				var str = jackTokenizer.StringVal(c.tokenizer)
				var strLetters = strings.Split(str, "")
				VMWriter.WritePush("constant", len(strLetters), c.vmWriter)
				VMWriter.WriteCall("String.new", 1, c.vmWriter)
				for i := 0; i < len(strLetters); i++ {
					intVar, err := strconv.Atoi(strLetters[i])
					_=err
					VMWriter.WritePush("constant", intVar, c.vmWriter) // (int)str.charAt(i))
					VMWriter.WriteCall("String.appendChar", 2, c.vmWriter)
				}

			}else if jackTokenizer.TokenType(c.tokenizer) == "KEYWORD" && jackTokenizer.Symbol(c.tokenizer) == "TRUE"{

				// ~0 is true
				VMWriter.WritePush("constant", 0,c.vmWriter)
				VMWriter.WriteArithmetic("not",c.vmWriter)
			}else if jackTokenizer.TokenType(c.tokenizer)== "KEYWORD" && jackTokenizer.Symbol(c.tokenizer) == "THIS" {
// push this pointer onto stack
				VMWriter.WritePush("pointer",0,c.vmWriter)

			}else if jackTokenizer.TokenType(c.tokenizer) == "KEYWORD" && StringInSlice(jackTokenizer.KeyWord(c.tokenizer),arr3) {
                // 0 for false and null
				VMWriter.WritePush("constant", 0,c.vmWriter)
			}else if jackTokenizer.TokenType(c.tokenizer) == "SYMBOL" && jackTokenizer.Symbol(c.tokenizer) == "(" {
// expression
				CompileExpression()
// ')'
            RequireSymbol(")",c)
			}else if jackTokenizer.TokenType(c.tokenizer) == "SYMBOL" && StringInSlice(jackTokenizer.Symbol(c.tokenizer),arr4) {

				s := jackTokenizer.Symbol(c.tokenizer)

				// term
				CompileTerm(c)

				if s == "-" {
					VMWriter.WriteArithmetic("neg", c.vmWriter)
				}else { // ~
					VMWriter.WriteArithmetic("not", c.vmWriter)
				}

			}else {
				print("Expected integerConstant or stringConstant or keywordConstant or '(' expression ')' or unaryOp term")
			}
	}
}


// Compiles a (possibly empty) comma-separated list of expressions.
// (expression(','expression)*)?
func CompileExpressionList(c *CompilationEngine) int{
	nArgs := 0

	jackTokenizer.Advance(&c.tokenizer)
	// determine if there is any expression, if next is ')' then no
	if jackTokenizer.TokenType(c.tokenizer) == "SYMBOL" && jackTokenizer.Symbol(c.tokenizer) == ")" {
		jackTokenizer.PointerBack(&c.tokenizer)
	} else {
		nArgs = 1
		jackTokenizer.PointerBack(&c.tokenizer)
		// expression
		CompileExpression(c)
		// (','expression)*
		for true{
			jackTokenizer.Advance(&c.tokenizer)
			if jackTokenizer.TokenType(c.tokenizer) == "SYMBOL" && jackTokenizer.Symbol(c.tokenizer) == "," {
				// expression
				CompileExpression()
				nArgs = nArgs + 1
			}else {
				jackTokenizer.PointerBack(&c.tokenizer)
				break
			}
		}
	}

	return nArgs
}

//Returns current function name, className.subroutineName.
func CurrentFunction(c *CompilationEngine) string  {

	if len(c.currentClass)!=0 && len(c.currentSubroutine)!=0{
		return c.currentClass+"."+c.currentSubroutine
	}else {
		return ""
	}

}

func CompileType(c *CompilationEngine) string {
	jackTokenizer.Advance(&c.tokenizer)
	arr := []string{"INT", "CHAR", "BOOLEAN"}
	//var arr [3]=["CONSTRUCTOR", "FUNCTION", "METHOD"]
	if jackTokenizer.TokenType(c.tokenizer)=="KEYWORD" && StringInSlice(jackTokenizer.KeyWord(c.tokenizer),arr) {
		return c.tokenizer.CurrentToken
	}
	if jackTokenizer.TokenType(c.tokenizer)=="IDENTIFIER"{
		return jackTokenizer.Identifier(c.tokenizer)
	}
	print("Expected int or char or boolean or className")
	return ""
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
