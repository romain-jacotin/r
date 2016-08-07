package token

// Token is the set of lexical tokens of the R programming language.
type Token int

// The list of tokens.
const (

	// Special tokens

	ERROR Token = iota
	END_OF_LINE
	END_OF_INPUT
	COMMENT        // #xxxxx[end of line]
	LINE_DIRECTIVE // #line xxx "sourcefile"

	// Constants

	CHARACTER_CONST // "xxx" | 'xxx' character constants
	NA_CHARACTER    // NA character contant

	INTEGER_CONST // Integer constant: 03L 1L 42L 0x42afL 0X042AFL
	NA_INTEGER    // NA integer constant

	DOUBLE_CONST // Double constant: 1 10 0.1 .2 1e-7 1.2e+7 0x123p456
	NA_REAL      // 'NA_real_' = NA double constant
	NAN          // 'NaN'      = Not a number double constant
	INF          // 'Inf'      = Infinity double constant

	COMPLEX_CONST // Complex constant: 2i 4.1i 1e-2i
	NA_COMPLEX    // 'NA_complex_' double constant

	TRUE_CONST  // 'TRUE' logical constant
	FALSE_CONST // 'FALSE' logical constant
	NA          // 'NA' logical constant

	NULL_CONST // 'NULL' constant

	// Reserved keywords

	IF       // if
	ELSE     // else
	FOR      // for
	IN       // in
	REPEAT   // repeat
	WHILE    // while
	NEXT     // next
	BREAK    // break
	FUNCTION // function

	// Symbols & Identifiers

	SYMBOL // xxx | `xxx` symbol name

	// Assignment operators

	RIGHT_ASSIGN  // ->
	RIGHT_ASSIGN2 // ->>
	LEFT_ASSIGN   // <-
	LEFT_ASSIGN2  // <<-
	EQUAL_ASSIGN  // =
	COLON_ASSIGN  // :=

	LSQUAR             // [
	LSQUAR2            // [[
	RSQUAR             // ]
	LPAREN             // (
	RPAREN             // )
	LCURLY             // {
	RCURLY             // }
	COLON              // :
	NAMESPACE          // ::
	NAMESPACE_INTERNAL // :::
	DOLLAR             // $
	AT                 // @

	// Arithmetic operators

	ADD  // +
	SUB  // -
	MUL  // *
	MUL2 // ** same as ^
	QUO  // /
	POW  // ^

	// Logical operators

	GT   // >
	GE   // >=
	LT   // <
	LE   // <=
	EQ   // ==
	NOT  // !
	NE   // !=
	AND  // &
	AND2 // &&
	OR   // |
	OR2  // ||

	TILDE    // ~
	QUESTION // ?

	// SYMBOL_FORMALS just a modif token in actions for = left operand in parameters definition in function() creation
	// EQ_FORMALS just a modif token in actions for = operator in parameters definition in function() creation
	// EQ_SUB just a modif tolen in actions for = operator
	// SYMBOL_SUB just a modif tolen in actions for = left operand
	// SYMBOL_FUNCTION_CALL just a modif token in actions
	// SYMBOL_PACKAGE just a modif token in actions for :: and ::: left operand
	// SLOT just a modif token in actions for @ right operand
)

var tokens = [...]string{
	ERROR:          "Invalid input",
	END_OF_LINE:    "End of line",
	END_OF_INPUT:   "End of input",
	COMMENT:        "# comment",
	LINE_DIRECTIVE: "#line directive",
	// Constants
	CHARACTER_CONST: "Character constant",
	NA_CHARACTER:    "NA_character_ constant",
	INTEGER_CONST:   "integer constant", // Integer constant: 03L 1L 42L 0x42afL 0X042AFL
	NA_INTEGER:      "NA_integer_ constant",
	DOUBLE_CONST:    "double constant", // Double constant: 1 10 0.1 .2 1e-7 1.2e+7 0x123p456
	NA_REAL:         "NA_real_ constant",
	NAN:             "NaN double constant",
	INF:             "Inf double constant",
	COMPLEX_CONST:   "complex constant", // Complex constant: 2i 4.1i 1e-2i
	NA_COMPLEX:      "NA_complex_ constant",
	TRUE_CONST:      "TRUE",
	FALSE_CONST:     "FALSE",
	NA:              "NA logical constant",
	NULL_CONST:      "NULL constant",
	// Reserved keywords
	IF:       "if",
	ELSE:     "else",
	FOR:      "for",
	IN:       "in",
	REPEAT:   "repeat",
	WHILE:    "while",
	NEXT:     "next",
	BREAK:    "break",
	FUNCTION: "function",
	// Symbols & Identifiers
	SYMBOL: "symbol name",
	// Assignment operators
	RIGHT_ASSIGN:  "->",
	RIGHT_ASSIGN2: "->>",
	LEFT_ASSIGN:   "<-",
	LEFT_ASSIGN2:  "<<-",
	EQUAL_ASSIGN:  "=",
	COLON_ASSIGN:  ":=",
	// Others compound tokens
	LSQUAR:             "[",
	LSQUAR2:            "[[",
	RSQUAR:             "]",
	LPAREN:             "(",
	RPAREN:             ")",
	LCURLY:             "{",
	RCURLY:             "}",
	COLON:              ":",
	NAMESPACE:          "::",
	NAMESPACE_INTERNAL: ":::",
	DOLLAR:             "$",
	AT:                 "@",
	// Arithmetic operators
	ADD:  "+",
	SUB:  "-",
	MUL:  "*",
	MUL2: "**",
	QUO:  "/",
	POW:  "^",
	// Logical operators
	GT:   ">",
	GE:   ">=",
	LT:   "<",
	LE:   "<=",
	EQ:   "==",
	NOT:  "!",
	NE:   "!=",
	AND:  "&",
	AND2: "&&",
	OR:   "|",
	OR2:  "||",
	// others
	TILDE:    "~",
	QUESTION: "?",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// Precedence returns the operator precedence of the binary operator op.
func (op Token) Precedence() int {
	switch op {
	case LOR:
		return 1
	case LAND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case ADD, SUB, OR, XOR:
		return 4
	case MUL, QUO, REM, SHL, SHR, AND, AND_NOT:
		return 5
	}
	return LowestPrec
}
