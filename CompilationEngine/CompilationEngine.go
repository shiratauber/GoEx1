package CompilationEnginesource

import (
	"GoEx1/Symbol"
	"GoEx1/jackTokenizer"
	"GoEx1/VMWriter"
	"GoEx1/SymbolTable"
		"os"
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
	CompileSubroutineBody(keyword)

	CompileSubroutine(c)

}

// Compiles the body of a subroutine.
// '{'  varDec* statements '}'
compileSubroutineBody = function(keyword) {
// '{'
self$requireSymbol('{')
// varDec*
self$compileVarDec()
// write VM function declaration
self$writeFunctionDec(keyword)
// statements
self$compileStatement()
// '}'
self$requireSymbol('}')
},

// Writes function declaration, load pointer when keyword is METHOD or CONSTRUCTOR.
writeFunctionDec = function(keyword) {
self$vmWriter$writeFunction(self$currentFunction(), self$symbolTable$varCount("VAR"))

// METHOD and CONSTRUCTOR need to load this pointer
if (keyword == "METHOD") {
// A Jack method with k arguments is compiled into a VM function that operates on k + 1 arguments.
// The first argument (argument number 0) always refers to the this object.
self$vmWriter$writePush("argument", 0)
self$vmWriter$writePop("pointer", 0)
} else if (keyword == "CONSTRUCTOR") {
// A Jack function or constructor with k arguments is compiled into a VM function that operates on k arguments.
self$vmWriter$writePush("constant", self$symbolTable$varCount("FIELD"))
self$vmWriter$writeCall("Memory.alloc", 1)
self$vmWriter$writePop("pointer", 0)
}

},

// Compiles a single statement.
compileStatement = function() {
// determine whether there is a statement next can be a '}'
self$tokenizer$advance()

// next is a '}'
if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$symbol() == '}'){
self$tokenizer$pointerBack()
return()
}

// next is 'let'|'if'|'while'|'do'|'return'
if (self$tokenizer$tokenType() != "KEYWORD") {
self$throwException("Expected keyword")
} else {
switch(self$tokenizer$keyWord(),
"LET"={
self$compileLet()
},
"IF"={
self$compileIf()
},
"WHILE"={
self$compileWhile()
},
"DO"={
self$compileDo()
},
"RETURN"={
self$compileReturn()
},
{   //   default
self$throwException("Expected let or if or while or do or return")
}
)
}

self$compileStatement()
},

// Compiles a (possibly empty) parameter list,
// not including the enclosing "()".
// ((type varName)(',' type varName)*)?
compileParameterList = function() {
// Check if there is parameterList, if next token is ')' than go back
self$tokenizer$advance()
if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$symbol() == ')'){
self$tokenizer$pointerBack()
return()
}

// there is parameter, at least one varName
self$tokenizer$pointerBack()
repeat{
// typeTok
typeTok  <- self$compileType()

// varName
self$tokenizer$advance()
if (self$tokenizer$tokenType() != "IDENTIFIER") {
self$throwException("Expected identifier")
}

self$symbolTable$define(self$tokenizer$identifier(), typeTok , "ARG")

// ',' or ')'
self$tokenizer$advance()
if (self$tokenizer$tokenType() != "SYMBOL" | !(self$tokenizer$symbol() %in% c(",", ")"))) {
self$throwException("Expected , or )")
}

if (self$tokenizer$symbol() == ')') {
self$tokenizer$pointerBack()
break
}
}
},

// Compiles a var declaration.
// 'var' type varName (',' varName)*;
compileVarDec = function() {
// determine if there is a varDec
self$tokenizer$advance()
// no 'var' go back
if (self$tokenizer$tokenType() != "KEYWORD" | self$tokenizer$keyWord() != "VAR"){
self$tokenizer$pointerBack()
return()
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
getSeg = function(kind) {
switch(kind,
"FIELD"={
return("this")
},
"STATIC"={
return("static")
},
"VAR"={
return("local")
},
"ARG"={
return("argument")
},
{   //   default
return("NONE")
}
)
},

// Compiles a while statement.
// 'while' '(' expression ')' '{' statements '}'
compileWhile = function() {
whileExpLabel <- paste("WHILE_EXP", self$labelCounterWhile, sep="")
whileEndLabel <- paste("WHILE_END", self$labelCounterWhile, sep="")
self$labelCounterWhile <- self$labelCounterWhile + 1

// top label for while loop
self$vmWriter$writeLabel(whileExpLabel)

// '('
self$requireSymbol('(')
// expression while condition: true or false
self$compileExpression()
// ')'
self$requireSymbol(')')

// if ~(condition) go to continue label
self$vmWriter$writeArithmetic("not")
self$vmWriter$writeIf(whileEndLabel)

// '{'
self$requireSymbol('{')
// statements
self$compileStatement()
// '}'
self$requireSymbol('}')

// if (condition) go to top label
self$vmWriter$writeGoto(whileExpLabel)
// or continue
self$vmWriter$writeLabel(whileEndLabel)

// self$labelIndex <- self$labelIndex + 1
},


## Compiles a return statement.
## ‘return’ expression? ';'
compileReturn = function() {
## check if there is any expression
self$tokenizer$advance()

if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$symbol() == ';') {
## no expression push 0 to stack
self$vmWriter$writePush("constant", 0)
}else {
## expression exist
self$tokenizer$pointerBack()
## expression
self$compileExpression()
## ';'
self$requireSymbol(';')
}

self$vmWriter$writeReturn()
},

## Compiles an if statement,
## possibly with a trailing else clause.
# 'if' '(' expression ')' '{' statements '}' ('else' '{' statements '}')?
compileIf = function() {
ifTrueLabel <- paste("IF_TRUE", self$labelCounterIf, sep="")
ifFalseLabel <- paste("IF_FALSE", self$labelCounterIf, sep="")
ifEndLabel <- paste("IF_END", self$labelCounterIf, sep="")
self$labelCounterIf <- self$labelCounterIf + 1

## '('
self$requireSymbol('(')
## expression
self$compileExpression()
## ')'
self$requireSymbol(')')

self$vmWriter$writeIf(ifTrueLabel)
self$vmWriter$writeGoto(ifFalseLabel)
self$vmWriter$writeLabel(ifTrueLabel)

## '{'
self$requireSymbol('{')
## statements
self$compileStatement()
## '}'
self$requireSymbol('}')

## check if there is 'else'
self$tokenizer$advance()
if (self$tokenizer$tokenType() == "KEYWORD" & self$tokenizer$keyWord() == "ELSE") {
# ifEndLabel <- paste("IF_END", self$labelIndex, sep="")
self$vmWriter$writeGoto(ifEndLabel)
self$vmWriter$writeLabel(ifFalseLabel)

## '{'
self$requireSymbol('{')
## statements
self$compileStatement()
## '}'
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

## Compiles a term. This routine is faced with a
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
##
## integerConstant|stringConstant|keywordConstant|varName|varName '[' expression ']'|subroutineCall| '(' expression ')'|unaryOp term
compileTerm = function() {
self$tokenizer$advance()
## check if it is an identifier
if (self$tokenizer$tokenType() == "IDENTIFIER") {
## varName|varName '[' expression ']'|subroutineCall
tempId <- self$tokenizer$identifier()

self$tokenizer$advance()
if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$symbol() == '[') {
## this is an array entry
## expression
self$compileExpression()
## ']'
self$requireSymbol(']')

## push array variable,base address into stack
self$vmWriter$writePush(self$getSeg(self$symbolTable$kindOf(tempId)), self$symbolTable$indexOf(tempId))


## base + offset
self$vmWriter$writeArithmetic("add")

## pop into 'that' pointer
self$vmWriter$writePop("pointer", 1)
## push *(base+index) onto stack
self$vmWriter$writePush("that", 0)
} else if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$symbol() %in% c("(", ".")) {
## this is a subroutineCall
self$tokenizer$pointerBack()
self$tokenizer$pointerBack()
self$compileSubroutineCall()
} else {
## this is varName
self$tokenizer$pointerBack()
## push variable directly onto stack
self$vmWriter$writePush(self$getSeg(self$symbolTable$kindOf(tempId)), self$symbolTable$indexOf(tempId))
}

} else {
## integerConstant|stringConstant|keywordConstant|'(' expression ')'|unaryOp term
if (self$tokenizer$tokenType() == "INT_CONST") {
## integerConstant just push its value onto stack
self$vmWriter$writePush("constant", self$tokenizer$intVal())
}else if (self$tokenizer$tokenType() == "STRING_CONST") {
## stringConstant new a string and append every char to the new stack
str <- self$tokenizer$stringVal()

strLetters <- strsplit(str, "")[[1]]

self$vmWriter$writePush("constant", length(strLetters))
self$vmWriter$writeCall("String.new", 1)

for (i in 1:length(strLetters)) {
self$vmWriter$writePush("constant", as.numeric(charToRaw(strLetters[i]))) ## (int)str.charAt(i))
self$vmWriter$writeCall("String.appendChar", 2)
}

}else if(self$tokenizer$tokenType() == "KEYWORD" & self$tokenizer$keyWord() == "TRUE") {
## ~0 is true
self$vmWriter$writePush("constant", 0)
self$vmWriter$writeArithmetic("not")

}else if(self$tokenizer$tokenType() == "KEYWORD" & self$tokenizer$keyWord() == "THIS") {
## push this pointer onto stack
self$vmWriter$writePush("pointer", 0)

}else if(self$tokenizer$tokenType() == "KEYWORD" & self$tokenizer$keyWord() %in% c("FALSE", "NULL")) {
## 0 for false and null
self$vmWriter$writePush("constant", 0)
}else if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$symbol() == '(') {
## expression
self$compileExpression()
## ')'
self$requireSymbol(')')
}else if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$symbol() %in% c("-", "~")) {

s <- self$tokenizer$symbol()

## term
self$compileTerm()

if (s == '-') {
self$vmWriter$writeArithmetic("neg")
}else { ## ~
self$vmWriter$writeArithmetic("not")
}

}else {
self$throwException("Expected integerConstant or stringConstant or keywordConstant or '(' expression ')' or unaryOp term")
}
}
},

## Compiles a (possibly empty) comma-separated list of expressions.
## (expression(','expression)*)?
compileExpressionList = function() {
nArgs <- 0

self$tokenizer$advance()
## determine if there is any expression, if next is ')' then no
if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$symbol() == ')') {
self$tokenizer$pointerBack()
} else {
nArgs <- 1
self$tokenizer$pointerBack()
## expression
self$compileExpression()
## (','expression)*
repeat{
self$tokenizer$advance()
if (self$tokenizer$tokenType() == "SYMBOL" & self$tokenizer$symbol() == ',') {
## expression
self$compileExpression()
nArgs <- nArgs + 1
}else {
self$tokenizer$pointerBack()
break
}
}
}

return(nArgs)
},

//Returns current function name, className.subroutineName.
func currentFunction(c *CompilationEngine) string  {

	if len(c.currentClass)!=0 && len(c.currentSubroutine)!=0{
		return c.currentClass+"."+c.currentSubroutine
	}else {
		return ""
	}

}


## Compiles a type.
compileType = function() {
self$tokenizer$advance()

if (self$tokenizer$tokenType() == "KEYWORD" & self$tokenizer$keyWord() %in% c("INT", "CHAR", "BOOLEAN")) {
return(self$tokenizer$currentToken)
}

if (self$tokenizer$tokenType() == "IDENTIFIER") {
return(self$tokenizer$identifier())
}

self$throwException("Expected int or char or boolean or className")
},

## Throws an exception, showing an error message and quiting.
throwException = function(errorMessage) {
print(errorMessage)
quit()
}

)
)
