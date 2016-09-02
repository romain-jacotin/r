package r

import "testing"
import "strings"
import "math"

// .romain if else break function ... jacotin
// NA NA_character_ NA_integer_ NA_complex_ NA_real_
// -> ->> <- <<- = :=
// [ [[ ] ( ) { } : :: ::: $ @
// + - * ** / ^
// > >= < <= == ! != & && | ||
// ~ ?
// %romain%  # los pollos hermanos  


// type Token struct {
//	Type 		TokenType
//	intvalue 	int64
//	realvalue	float64
//	stringvalue	string
//	offset		int
//	nrow 		int
//	nline 		int
//	nbyte		int
//	nrune		int }

func TestProcessHexadecimal(e *testing.T) {
	var t *Token
	var str = " 0x0001a2B3c 0X1a2B3c 0x1a2B3c.4d5E6fp0 0x1a2B3c.4d5E6fP0 0x1a2B3c.4d5E6fp+1 0x1a2B3c.4d5E6fP-1 0x.4d5E6fP0 0x.4d5E6fp0 0x.4d5E6fp+1 0x.4d5E6fP-1 0x0001a2B3cL 0X1a2B3cL 0x1a2B3c.4d5E6fp0L 0x1a2B3c.4d5E6fP0L 0x1a2B3c.4d5E6fp+1L 0x1a2B3c.4d5E6fP-1L 0x.4d5E6fP0L 0x.4d5E6fp0L 0x.4d5E6fp+1L 0x.4d5E6fP-1L 0x0001a2B3ci 0X1a2B3ci 0x1a2B3c.4d5E6fp0i 0x1a2B3c.4d5E6fP0i 0x1a2B3c.4d5E6fp+1i 0x1a2B3c.4d5E6fP-1i 0x.4d5E6fP0i 0x.4d5E6fp0i 0x.4d5E6fp+1i 0x.4d5E6fP-1i "
	var tests []Token = []Token{
		{ CONST_REAL, 1715004, 1715004, "0x0001a2B3c", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 1715004, 1715004, "0X1a2B3c", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 1715004, 1715004.30222219228744506836, "0x1a2B3c.4d5E6fp0", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 1715004, 1715004.30222219228744506836, "0x1a2B3c.4d5E6fP0", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 3430008, 3430008.60444438457489013672, "0x1a2B3c.4d5E6fp+1", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 857502, 857502.15111109614372253418, "0x1a2B3c.4d5E6fP-1", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 0, 0.30222219228744506836, "0x.4d5E6fP0", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 0, 0.30222219228744506836, "0x.4d5E6fp0", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 0, 0.60444438457489013672, "0x.4d5E6fp+1", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 0, 0.15111109614372253418, "0x.4d5E6fP-1", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 1715004, 1715004, "0x0001a2B3cL", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 1715004, 1715004, "0X1a2B3cL", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 1715004, 1715004.30222219228744506836, "0x1a2B3c.4d5E6fp0L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 1715004, 1715004.30222219228744506836, "0x1a2B3c.4d5E6fP0L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 3430008, 3430008.60444438457489013672, "0x1a2B3c.4d5E6fp+1L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 857502, 857502.15111109614372253418, "0x1a2B3c.4d5E6fP-1L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 0, 0.30222219228744506836, "0x.4d5E6fP0L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 0, 0.30222219228744506836, "0x.4d5E6fp0L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 0, 0.60444438457489013672, "0x.4d5E6fp+1L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 0, 0.15111109614372253418, "0x.4d5E6fP-1L", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 1715004, 1715004, "0x0001a2B3ci", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 1715004, 1715004, "0X1a2B3ci", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 1715004, 1715004.30222219228744506836, "0x1a2B3c.4d5E6fp0i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 1715004, 1715004.30222219228744506836, "0x1a2B3c.4d5E6fP0i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 3430008, 3430008.60444438457489013672, "0x1a2B3c.4d5E6fp+1i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 857502, 857502.15111109614372253418, "0x1a2B3c.4d5E6fP-1i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 0, 0.30222219228744506836, "0x.4d5E6fP0i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 0, 0.30222219228744506836, "0x.4d5E6fp0i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 0, 0.60444438457489013672, "0x.4d5E6fp+1i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 0, 0.15111109614372253418, "0x.4d5E6fP-1i", 0, 0, 0, 0, 0 },
	}

	r := strings.NewReader(str)
	s := NewScanner(r)
	for i := 0; i <len(tests);i++ {
		t = s.NextToken()
		if t.Type != tests[i].Type { e.Error("Test Hexa", t.stringvalue, "Type Failed", t.Type) }
		if t.intvalue != tests[i].intvalue { e.Error("Test Hexa", t.stringvalue, "intvalue Failed", t.intvalue) }
		if t.realvalue != tests[i].realvalue { e.Error("Test Hexa", t.stringvalue, "realvalue Failed", t.stringvalue) }
		if t.stringvalue != tests[i].stringvalue { e.Error("Test Hexa", t.stringvalue, "stringvalue Failed") }
	}
}

func TestProcessDecimal(e *testing.T) {
	var t *Token
	var str = " 1.79769313486231e+308 1e-323 1e-324 1e309 000123 123 123. 123.456 123.456e0 123.456E0 123.456e+1 123.456E-1 .456 .456e0 .456E0 .456e+1 .456E-1 000123L 123L 123.L 123.456L 123.456e0L 123.456E0L 123.456e+1L 123.456E-1L .456L .456e0L .456E0L .456e+1L .456E-1L 000123i 123i 123.i 123.456i 123.456e0i 123.456E0i 123.456e+1i 123.456E-1i .456i .456e0i .456E0i .456e+1i .456E-1i "
	var tests []Token = []Token{
		{ CONST_REAL, -9223372036854775808, 1.79769313486231e+308, "1.79769313486231e+308", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 0, 1e-323, "1e-323", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 0, 0, "1e-324", 0, 0, 0, 0, 0 },
		{ CONST_REAL, -9223372036854775808, math.Inf(0), "1e309", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 123, 123.0, "000123", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 123, 123.0, "123", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 123, 123.0, "123.", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 123, 123.456, "123.456", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 123, 123.456, "123.456e0", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 123, 123.456, "123.456E0", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 1234, 1234.56, "123.456e+1", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 12, 12.3456, "123.456E-1", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 0, 0.456, ".456", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 0, 0.456, ".456e0", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 0, 0.456, ".456E0", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 4, 4.56, ".456e+1", 0, 0, 0, 0, 0 },
		{ CONST_REAL, 0, 0.0456, ".456E-1", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 123, 123.0, "000123L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 123, 123.0, "123L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 123, 123.0, "123.L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 123, 123.456, "123.456L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 123, 123.456, "123.456e0L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 123, 123.456, "123.456E0L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 1234, 1234.56, "123.456e+1L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 12, 12.3456, "123.456E-1L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 0, 0.456, ".456L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 0, 0.456, ".456e0L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 0, 0.456, ".456E0L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 4, 4.56, ".456e+1L", 0, 0, 0, 0, 0 },
		{ CONST_INTEGER, 0, 0.0456, ".456E-1L", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 123, 123.0, "000123i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 123, 123.0, "123i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 123, 123.0, "123.i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 123, 123.456, "123.456i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 123, 123.456, "123.456e0i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 123, 123.456, "123.456E0i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 1234, 1234.56, "123.456e+1i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 12, 12.3456, "123.456E-1i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 0, 0.456, ".456i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 0, 0.456, ".456e0i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 0, 0.456, ".456E0i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 4, 4.56, ".456e+1i", 0, 0, 0, 0, 0 },
		{ CONST_COMPLEX, 0, 0.0456, ".456E-1i", 0, 0, 0, 0, 0 },
	}

	r := strings.NewReader(str)
	s := NewScanner(r)
	for i := 0; i <len(tests);i++ {
		t = s.NextToken()
		if t.Type != tests[i].Type { e.Error("Test Decimale", t.stringvalue, "Type Failed", t.Type) }
		if t.intvalue != tests[i].intvalue { e.Error("Test Decimale", t.stringvalue, "intvalue Failed", t.intvalue) }
		if t.realvalue != tests[i].realvalue { e.Error("Test Decimale", t.stringvalue, "realvalue Failed", t.realvalue) }
		if t.stringvalue != tests[i].stringvalue { e.Error("Test Decimale", t.stringvalue, "stringvalue Failed") }
	}
}

func TestProcessSymbol(e *testing.T) {
	var t *Token
	var str = ` .romain if else for in repeat while next break function ... ja_co.tin	
	 NA NA_character_ NA_integer_ NA_complex_ NA_real_ 
	 -> ->> <- <<- = := [ [[ ] ( ) { } : :: ::: $ @
 + - * ** / ^ > >= < <= == ! != & && | ||  ~ ? 
 %romain%  # los pollos hermanos 	 `

	var tests []Token = []Token{
		// Tests 0 to 12
		// .romain if else for in repeat while next break function ... ja_co.tin
		{ SYMBOL, 0, 0, ".romain", 0, 0, 0, 0, 0 },
		{ KEYWORD_IF, 0, 0, "if", 0, 0, 0, 0, 0 },
		{ KEYWORD_ELSE, 0, 0, "else", 0, 0, 0, 0, 0 },
		{ KEYWORD_FOR, 0, 0, "for", 0, 0, 0, 0, 0 },
		{ KEYWORD_IN, 0, 0, "in", 0, 0, 0, 0, 0 },
		{ KEYWORD_REPEAT, 0, 0, "repeat", 0, 0, 0, 0, 0 },
		{ KEYWORD_WHILE, 0, 0, "while", 0, 0, 0, 0, 0 },
		{ KEYWORD_NEXT, 0, 0, "next", 0, 0, 0, 0, 0 },
		{ KEYWORD_BREAK, 0, 0, "break", 0, 0, 0, 0, 0 },
		{ KEYWORD_FUNCTION, 0, 0, "function", 0, 0, 0, 0, 0 },
		{ SYMBOL, 0, 0, "...", 0, 0, 0, 0, 0 },
		{ SYMBOL, 0, 0, "ja_co.tin", 0, 0, 0, 0, 0 },
		{ END_OF_LINE, 0, 0, "", 0, 0, 0, 0, 0 },
		// Tests 13 to 18
		// NA NA_character_ NA_integer_ NA_complex_ NA_real_
		{ NA_LOGICAL, 0, 0, "NA", 0, 0, 0, 0, 0 },
		{ NA_CHARACTER, 0, 0, "NA_character_", 0, 0, 0, 0, 0 },
		{ NA_INTEGER, 0, 0, "NA_integer_", 0, 0, 0, 0, 0 },
		{ NA_COMPLEX, 0, 0, "NA_complex_", 0, 0, 0, 0, 0 },
		{ NA_REAL, 0, 0, "NA_real_", 0, 0, 0, 0, 0 },
		{ END_OF_LINE, 0, 0, "", 0, 0, 0, 0, 0 },
		// Tests 19 to 37
		// -> ->> <- <<- = := [ [[ ] ( ) { } : :: ::: $ @
		{ OP_RIGHT_ASSIGN, 0, 0, "->", 0, 0, 0, 0, 0 },
		{ OP_RIGHT_ASSIGN2, 0, 0, "->>", 0, 0, 0, 0, 0 },
		{ OP_LEFT_ASSIGN, 0, 0, "<-", 0, 0, 0, 0, 0 },
		{ OP_LEFT_ASSIGN2, 0, 0, "<<-", 0, 0, 0, 0, 0 },
		{ OP_EQUAL_ASSIGN, 0, 0, "=", 0, 0, 0, 0, 0 },
		{ OP_COLON_ASSIGN, 0, 0, ":=", 0, 0, 0, 0, 0 },
		{ OP_LEFT_SQUARE, 0, 0, "[", 0, 0, 0, 0, 0 },
		{ OP_LEFT_SQUARE2, 0, 0, "[[", 0, 0, 0, 0, 0 },
		{ OP_RIGHT_SQUARE, 0, 0, "]", 0, 0, 0, 0, 0 },
		{ OP_LEFT_ROUND, 0, 0, "(", 0, 0, 0, 0, 0 },
		{ OP_RIGHT_ROUND, 0, 0, ")", 0, 0, 0, 0, 0 },
		{ OP_LEFT_CURLY, 0, 0, "{", 0, 0, 0, 0, 0 },
		{ OP_RIGHT_CURLY, 0, 0, "}", 0, 0, 0, 0, 0 },
		{ OP_COLON, 0, 0, ":", 0, 0, 0, 0, 0 },
		{ OP_NAMESPACE, 0, 0, "::", 0, 0, 0, 0, 0 },
		{ OP_NAMESPACE_INTERNAL, 0, 0, ":::", 0, 0, 0, 0, 0 },
		{ OP_DOLLAR, 0, 0, "$", 0, 0, 0, 0, 0 },
		{ OP_AT, 0, 0, "@", 0, 0, 0, 0, 0 },
		{ END_OF_LINE, 0, 0, "", 0, 0, 0, 0, 0 },
		// Tests 38 to 57
		// + - * ** / ^ > >= < <= == ! != & && | || ~ ?
		{ OP_ADD, 0, 0, "+", 0, 0, 0, 0, 0 },
		{ OP_SUB, 0, 0, "-", 0, 0, 0, 0, 0 },
		{ OP_MUL, 0, 0, "*", 0, 0, 0, 0, 0 },
		{ OP_MUL2, 0, 0, "**", 0, 0, 0, 0, 0 },
		{ OP_DIV, 0, 0, "/", 0, 0, 0, 0, 0 },
		{ OP_POW, 0, 0, "^", 0, 0, 0, 0, 0 },
		{ OP_GT, 0, 0, ">", 0, 0, 0, 0, 0 },
		{ OP_GE, 0, 0, ">=", 0, 0, 0, 0, 0 },
		{ OP_LT, 0, 0, "<", 0, 0, 0, 0, 0 },
		{ OP_LE, 0, 0, "<=", 0, 0, 0, 0, 0 },
		{ OP_EQ, 0, 0, "==", 0, 0, 0, 0, 0 },
		{ OP_NOT, 0, 0, "!", 0, 0, 0, 0, 0 },
		{ OP_NE, 0, 0, "!=", 0, 0, 0, 0, 0 },
		{ OP_AND, 0, 0, "&", 0, 0, 0, 0, 0 },
		{ OP_AND2, 0, 0, "&&", 0, 0, 0, 0, 0 },
		{ OP_OR, 0, 0, "|", 0, 0, 0, 0, 0 },
		{ OP_OR2, 0, 0, "||", 0, 0, 0, 0, 0 },
		{ OP_TILDE, 0, 0, "~", 0, 0, 0, 0, 0 },
		{ OP_QUESTION, 0, 0, "?", 0, 0, 0, 0, 0 },
		{ END_OF_LINE, 0, 0, "", 0, 0, 0, 0, 0 },
		// Test 58 to 60
		//  %romain%  # los pollos hermanos
		{ INFIX, 0, 0, "%romain%", 0, 0, 0, 0, 0 },
		{ COMMENT, 0, 0, "# los pollos hermanos 	 ", 0, 0, 0, 0, 0 },

	}

	r := strings.NewReader(str)
	s := NewScanner(r)
	for i := 0; i <len(tests);i++ {
		t = s.NextToken()
		if t.Type != tests[i].Type { e.Error("Test Symbol[", i, "]", t.stringvalue, "Type Failed", t.Type) }
		if t.stringvalue != tests[i].stringvalue { e.Error("Test Symbol[", i, "]", t.stringvalue, "stringvalue Failed") }
	}
}

func TestProcessComment(e *testing.T) {
	var t *Token
	var str = `#line "romain"  
	dfg # dfdf dfdf d'() 
# dssdf sdfsdf(§) 
# dfdf dfdf d'()`
	var tests []Token = []Token{
		{ LINE_DIRECTIVE, 0, 0, "#line \"romain\"  ", 0, 0, 0, 0, 0 },
		{ END_OF_LINE, 0, 0, "", 0, 0, 0, 0, 0 },
		{ SYMBOL, 0, 0, "dfg", 0, 0, 0, 0, 0 },
		{ COMMENT, 0, 0, "# dfdf dfdf d'() ", 0, 0, 0, 0, 0 },
		{ END_OF_LINE, 0, 0, "", 0, 0, 0, 0, 0 },
		{ COMMENT, 0, 0, "# dssdf sdfsdf(§) ", 0, 0, 0, 0, 0 },
		{ END_OF_LINE, 0, 0, "", 0, 0, 0, 0, 0 },
		{ COMMENT, 0, 0, "# dfdf dfdf d'()", 0, 0, 0, 0, 0 },
	}

	r := strings.NewReader(str)
	s := NewScanner(r)
	for i := 0; i <len(tests);i++ {
		t = s.NextToken()
		if t.Type != tests[i].Type { e.Error("Test Symbol[", i, "]", t.stringvalue, "Type Failed", t.Type) }
		if t.stringvalue != tests[i].stringvalue { e.Error("Test Symbol[", i, "]", t.stringvalue, "stringvalue Failed") }
	}
}

func TestProcessString(e *testing.T) {
	var t *Token
	var str = ` "  hj
 k k\n df
 jhjh jhjh \
	dfg # dfdf dfdf d'()"  "+\7-" "*\07/" "\007" '\0074'   "\x7aD"   'm\u7m' 'm\u07m'  'm\u007m' 'm\u0007m' 'm\u{07}2m' 
 'n\U7n' 'n\U07n'  'n\U007n' 'n\U0007n' 'n\U{00000007}2n'   `
	var tests []Token = []Token{
		{ CONST_CHARACTER, 0, 0, "  hj\n k k\n df\n jhjh jhjh \n	dfg # dfdf dfdf d'()", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "+\a-", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "*\a/", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "\a", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "\a4", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "zD", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "m\am", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "m\am", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "m\am", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "m\am", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "m\a2m", 0, 0, 0, 0, 0 },
		{ END_OF_LINE, 0, 0, "", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "n\an", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "n\an", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "n\an", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "n\an", 0, 0, 0, 0, 0 },
		{ CONST_CHARACTER, 0, 0, "n\a2n", 0, 0, 0, 0, 0 },
	}

	r := strings.NewReader(str)
	s := NewScanner(r)
	for i := 0; i <len(tests);i++ {
		t = s.NextToken()
		if t.Type != tests[i].Type { e.Error("Test Symbol[", i, "]", t.stringvalue, "Type Failed", t.Type) }
		if t.stringvalue != tests[i].stringvalue { e.Error("Test Symbol[", i, "]", t.stringvalue, "stringvalue Failed") }
	}
}

