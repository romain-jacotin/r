package r

import "unicode"
import "io"
import "bufio"
import "errors"
import "math"
import "strconv"
import "strings"

type character struct {
	// Rune
	r 		rune
	// File offest
	offset	int
	// Text editor position
	ncol	int
	nline	int
	// Data quantity
	nbyte	int
}

type Scanner struct {
	// File reader
	reader			*bufio.Reader
	// Last file offset
	currentOffset	int
	// Last text editor position
	nline			int
	ncol			int
	// Data quantity
	nrune			int
	// Pushback buffer to handle rune look ahead
	npush			int
	pushback		[16]character
}

// Token is the set of lexical tokens of the R programming language.
type Token struct {
	// Token type
	Type 		TokenType
	// Token value (if needed)
	intvalue 	int64
	realvalue	float64
	stringvalue	string
	// File offest
	offset		int
	// Text editor position
	ncol		int
	nline 		int
	// Data quantity
	nbyte		int
	nrune		int
}

type TokenType int32

// The list of tokens.
const (

	// Special tokens
	ERROR TokenType = iota
	END_OF_INPUT
	END_OF_LINE
	COMMENT        // #xxxxx[end of line]
	LINE_DIRECTIVE // #line xxx "sourcefile"

	// Constants
	CONST_CHARACTER // "xxx" | 'xxx' character constants
	CONST_INTEGER   // Integer constant: 03L 1L 42L 0x42afL 0X042AFL
	CONST_REAL      // Real constant: 1 10 0.1 .2 1e-7 1.2e+7 0x123p456
	CONST_COMPLEX   // Complex constant: 2i 4.1i 1e-2i
	CONST_NAN       // 'NaN'      = Not a number double constant
	CONST_INF       // 'Inf'      = Infinity double constant
	CONST_TRUE      // 'TRUE' logical constant
	CONST_FALSE     // 'FALSE' logical constant
	CONST_NULL      // 'NULL' constant

	// NA constants
	NA_CHARACTER  // 'NA_character_' contant
	NA_INTEGER    // 'NA_integer_' constant
	NA_REAL    	  // 'NA_real_' = NA double constant
	NA_COMPLEX    // 'NA_complex_' double constant
	NA_LOGICAL    // 'NA' logical constant

	// Reserved keywords
	KEYWORD_IF       // if
	KEYWORD_ELSE     // else
	KEYWORD_FOR      // for
	KEYWORD_IN       // in
	KEYWORD_REPEAT   // repeat
	KEYWORD_WHILE    // while
	KEYWORD_NEXT     // next
	KEYWORD_BREAK    // break
	KEYWORD_FUNCTION // function

	// Symbols & Identifiers
	SYMBOL // xxx | `xxx` symbol name

	// Infix operator
	INFIX 	// %xxxxxxx%

	// Assignment operators
	OP_RIGHT_ASSIGN  // ->
	OP_RIGHT_ASSIGN2 // ->>
	OP_LEFT_ASSIGN   // <-
	OP_LEFT_ASSIGN2  // <<-
	OP_EQUAL_ASSIGN  // =
	OP_COLON_ASSIGN  // :=

	OP_LEFT_SQUARE        // [
	OP_LEFT_SQUARE2       // [[
	OP_RIGHT_SQUARE       // ]
	OP_LEFT_ROUND         // (
	OP_RIGHT_ROUND        // )
	OP_LEFT_CURLY         // {
	OP_RIGHT_CURLY        // }
	OP_COLON              // :
	OP_NAMESPACE          // ::
	OP_NAMESPACE_INTERNAL // :::
	OP_DOLLAR             // $
	OP_AT                 // @

	// Arithmetic operators
	OP_ADD  // +
	OP_SUB  // -
	OP_MUL  // *
	OP_MUL2 // ** same as ^
	OP_DIV  // /
	OP_POW  // ^

	// Logical operators
	OP_GT   // >
	OP_GE   // >=
	OP_LT   // <
	OP_LE   // <=
	OP_EQ   // ==
	OP_NOT  // !
	OP_NE   // !=
	OP_AND  // &
	OP_AND2 // &&
	OP_OR   // |
	OP_OR2  // ||

	OP_TILDE    // ~
	OP_QUESTION // ?

	// SYMBOL_FORMALS just a modif token in actions for = left operand in parameters definition in function() creation
	// EQ_FORMALS just a modif token in actions for = operator in parameters definition in function() creation
	// EQ_SUB just a modif tolen in actions for = operator
	// SYMBOL_SUB just a modif tolen in actions for = left operand
	// SYMBOL_FUNCTION_CALL just a modif token in actions
	// SYMBOL_PACKAGE just a modif token in actions for :: and ::: left operand
	// SLOT just a modif token in actions for @ right operand
)

func (this TokenType) String() (s string) {
	switch this {

	case ERROR : s = "ERROR"
	case END_OF_INPUT : s = "END_OF_INPUT"
	case END_OF_LINE : s = "END_OF_LINE"
	case COMMENT : s = "COMMENT"
	case LINE_DIRECTIVE : s = "LINE_DIRECTIVE"
	
	// Constants

	case CONST_CHARACTER : s = "CONST_CHARACTER"
	case CONST_INTEGER : s = "CONST_INTEGER"
	case CONST_REAL : s = "CONST_REAL"
	case CONST_COMPLEX : s = "CONST_COMPLEX"
	case CONST_NAN : s = "CONST_NAN"
	case CONST_INF : s = "CONST_INF"
	case CONST_TRUE : s = "CONST_TRUE"
	case CONST_FALSE : s = "CONST_FALSE"
	case CONST_NULL : s = "CONST_NULL"

	// NA constants

	case NA_CHARACTER : s = "NA_CHARACTER"
	case NA_INTEGER : s = "NA_INTEGER"
	case NA_REAL : s = "NA_REAL"
	case NA_COMPLEX : s = "NA_COMPLEX"
	case NA_LOGICAL : s = "NA_LOGICAL"

	// Reserved keywords

	case KEYWORD_IF : s = "IF"
	case KEYWORD_ELSE : s = "ELSE"
	case KEYWORD_FOR : s = "FOR"
	case KEYWORD_IN : s = "IN"
	case KEYWORD_REPEAT : s = "REPEAT"
	case KEYWORD_WHILE : s = "WHILE"
	case KEYWORD_NEXT : s = "NEXT"
	case KEYWORD_BREAK : s = "BREAK"
	case KEYWORD_FUNCTION : s = "FUNCTION"

	// Symbols & Identifiers

	case SYMBOL : s = "SYMBOL"

	// Infix operator

	case INFIX : s = "INFIX"

	// Assignment operators

	case OP_RIGHT_ASSIGN : s = "RIGHT_ASSIGN"
	case OP_RIGHT_ASSIGN2 : s = "RIGHT_ASSIGN2"
	case OP_LEFT_ASSIGN : s = "LEFT_ASSIGN"
	case OP_LEFT_ASSIGN2 : s = "LEFT_ASSIGN2"
	case OP_EQUAL_ASSIGN : s = "EQUAL_ASSIGN"
	case OP_COLON_ASSIGN : s = "COLON_ASSIGN"

	case OP_LEFT_SQUARE : s = "LEFT_SQUARE"
	case OP_LEFT_SQUARE2 : s = "LEFT_SQUARE2"
	case OP_RIGHT_SQUARE : s = "RIGHT_SQUARE"
	case OP_LEFT_ROUND : s = "LEFT_ROUND"
	case OP_RIGHT_ROUND : s = "RIGHT_ROUND"
	case OP_LEFT_CURLY : s = "LEFT_CURLY"
	case OP_RIGHT_CURLY : s = "RIGHT_CURLY"
	case OP_COLON : s = "COLON"
	case OP_NAMESPACE : s = "NAMESPACE"
	case OP_NAMESPACE_INTERNAL : s = "NAMESPACE_INTERNAL"
	case OP_DOLLAR : s = "DOLLAR"
	case OP_AT : s = "AT"

	// Arithmetic operators

	case OP_ADD : s = "ADD"
	case OP_SUB : s = "SUB"
	case OP_MUL : s = "MUL"
	case OP_MUL2 : s = "MUL2"
	case OP_DIV : s = "DIV"
	case OP_POW : s = "POW"

	// Logical operators

	case OP_GT : s = "GT"
	case OP_GE : s = "GE"
	case OP_LT : s = "LT"
	case OP_LE : s = "LE"
	case OP_EQ : s = "EQ"
	case OP_NOT : s = "NOT"
	case OP_NE : s = "NE"
	case OP_AND : s = "AND"
	case OP_AND2 : s = "AND2"
	case OP_OR : s = "OR"
	case OP_OR2 : s = "OR2"

	case OP_TILDE : s = "TILDE"
	case OP_QUESTION : s = "QUESTION"
	}
	return
}

func (this *Scanner) getCharacter() (c *character, err error) {
	var ru rune
	var nb int

	if this.npush > 0 {
		c = &this.pushback[this.npush-1]
		this.npush--
	} else if ru, nb, err = this.reader.ReadRune(); err == nil {
		c = new(character)

		c.r = ru
		c.nbyte = nb
		c.offset = this.currentOffset
		c.ncol = this.ncol
		c.nline = this.nline

		this.currentOffset += c.nbyte
		this.nrune++
		if ru == '\n' {
			this.ncol = 1
			this.nline++
		} else {
			this.ncol++
		}
	}
	return
}

func (this *Scanner) ungetCharacter(c *character) (err error) {
	if this.npush < len(this.pushback) {
		this.pushback[this.npush] = *c
		this.npush++
	} else {
		err = errors.New("PUSHBACK Buffer full")
	}
	return
}

func (this *Scanner) isNextCharacter(r rune) (b bool, err error) {
	var c *character

	if c, err = this.getCharacter(); err != nil {
		return false, err
	}
	if c.r == r {
		b = true
	} else {
		b = false
		err = this.ungetCharacter(c)
	}
	return
}

func (this *Scanner) isNextCharacterLetter(r rune) (b bool, err error) {
	var c *character

	if c, err = this.getCharacter(); err != nil {
		return false, err
	}
	if c.r == unicode.ToLower(r) || c.r == unicode.ToUpper(r) {
		b = true
	} else {
		b = false
	}
	err = this.ungetCharacter(c)
	return
}

func (this *Scanner) isNextCharacterDigit() (b bool, err error) {
	var c *character

	if c, err = this.getCharacter(); err != nil {
		return false, err
	}
	if unicode.IsDigit(c.r) {
		b = true
	} else {
		b = false
	}
	err = this.ungetCharacter(c)
	return
}

func (this *Scanner) getCharacterAfterSpaces() (c *character, err error) {
	for {
		// Trying to read the next rune
		if c, err = this.getCharacter(); err != nil {
			return
		}
		if c.r != ' ' && c.r != '\t' && c.r != '\f' {
				break
		}
	}
	return
}

func (this *Scanner) processLineDirective(t *Token) {
	t.Type = LINE_DIRECTIVE
	return
}

func (this *Scanner) processComment(t *Token) {
	var c *character
	var err error

	t.stringvalue = "#"
	for {
		// Trying to read the next rune
		if c, err = this.getCharacter(); err != nil {
			if err == io.EOF {
				t.Type = COMMENT
			} else {
				t.Type = ERROR
			}
			return
		}
		switch c.r {
		case '\n' :
			t.Type = COMMENT
			if err = this.ungetCharacter(c); err != nil {
				t.Type = ERROR
			}
			if strings.HasPrefix(t.stringvalue, "#line") {
				this.processLineDirective(t)
			}
			return
		default:
			t.stringvalue += string(c.r)
		}
	}
	return
}

func (this *Scanner) processString(r rune, t *Token) {
	var c, cc, ccc								*character
	var hasunicode, hasoctal, hashexa, hascurly	bool
	var err 									error
	var v 										rune

	for {
		// Trying to read the next rune
		if c, err = this.getCharacter(); err != nil {
			t.Type = ERROR
			return
		}
		switch c.r {
		case '\n' :
			t.stringvalue += "\n"

		case '"', '\'' :
			if c.r == r {
				if hasunicode && (hasoctal || hashexa) {
					t.Type = ERROR
					return
				} 
				t.Type = CONST_CHARACTER
				return
			} else {
				t.stringvalue += string(c.r)
			}

		case '`' :
			if c.r == r {
				t.Type = SYMBOL
				return
			} else {
				t.stringvalue += string(c.r)
			}

		case '\\' :
			// Trying to read the next rune
			if cc, err = this.getCharacter(); err != nil {
				t.Type = ERROR
				return
			}
			switch cc.r {
			case ' ' : // \space space
				t.stringvalue = t.stringvalue + " "
			case 'n', '\n' : // \n	newline
				t.stringvalue = t.stringvalue + "\n"
			case 'r' : // \r	carriage return
				t.stringvalue = t.stringvalue + "\r"
			case 't' : // \t	tab
				t.stringvalue = t.stringvalue + "\t"
			case 'b' : // \b	backspace
				t.stringvalue = t.stringvalue + "\b"
			case 'a' : // \a	alert (bell)
				t.stringvalue = t.stringvalue + "\a"
			case 'f' : // \f	form feed
				t.stringvalue = t.stringvalue + "\f"
			case 'v' : // \v	vertical tab
				t.stringvalue = t.stringvalue + "\v"
			case '\\' : // \\	backslash \
				t.stringvalue = t.stringvalue + "\\"
			case '\'' : // \'	ASCII apostrophe '
				t.stringvalue = t.stringvalue + "'"
			case '"' : // \"	ASCII quotation mark "
				t.stringvalue = t.stringvalue + "\""
			case '`' : // \`	ASCII grave accent (backtick) `
				t.stringvalue = t.stringvalue + "`"

			case '0', '1', '2', '3', '4', '5', '6', '7' : // \nnn	character with given octal code (1, 2 or 3 digits)
				hasoctal = true
				v = cc.r - '0'
				for i := 0; i < 2; i++ {
					// Trying to read the next rune
					if ccc, err = this.getCharacter(); err != nil {
						t.Type = ERROR
						return
					}
					if ccc.r >= '0' && ccc.r <= '7' {
						v = 8*v + ccc.r - '0'
					} else {
						// pushback the last character if it is not '0-7'
						if err = this.ungetCharacter(ccc); err != nil {
							t.Type = ERROR
						}
						break
					}
				}
				if v == 0 { // nul character not allowed
						t.Type = ERROR
						return
				}
				t.stringvalue += string(v)

			case 'x' : // \xnn	character with given hex code (1 or 2 hex digits)
				hashexa = true
				v = 0
				for i := 0; i < 2; i++ {
					// Trying to read the next rune
					if ccc, err = this.getCharacter(); err != nil {
						t.Type = ERROR
						return
					}
					if this.isCharacterHexa(ccc.r) {
						v = v*16 + rune(this.hexaValue(ccc.r))
					} else {
						// pushback the last character if it is not '0-9a-fA-F'
						if err = this.ungetCharacter(ccc); err != nil {
							t.Type = ERROR
						}
						break
					}
				}
				if v == 0 { // nul character not allowed
					t.Type = ERROR
					return
				}
				t.stringvalue += string(v)

			case 'u' : // \unnnn	Unicode character with given code (1--4 hex digits)
				if r == '`' { // \\uxxxx sequences not supported inside backticks
					t.Type = ERROR
					return
				}
				hasunicode = true
				// Trying to read the next rune
				if ccc, err = this.getCharacter(); err != nil {
					t.Type = ERROR
					return
				}
				if ccc.r == '{' {
					hascurly = true
				} else {
					hascurly = false
					// pushback the last character if it is not '{'
					if err = this.ungetCharacter(ccc); err != nil {
						t.Type = ERROR
					}					
				}
				v = 0
				for i := 0; i < 4; i++ {
					// Trying to read the next rune
					if ccc, err = this.getCharacter(); err != nil {
						t.Type = ERROR
						return
					}
					if this.isCharacterHexa(ccc.r) {
						v = v*16 + rune(this.hexaValue(ccc.r))
					} else {
						// pushback the last character if it is not '0-9a-fA-F'
						if err = this.ungetCharacter(ccc); err != nil {
							t.Type = ERROR
						}
						break
					}
				}
				if v == 0 { // nul character not allowed
					t.Type = ERROR
					return
				}
				if hascurly {
					// Trying to read the next rune
					if ccc, err = this.getCharacter(); err != nil {
						t.Type = ERROR
						return
					}
					if ccc.r != '}' {
						t.Type = ERROR
						return
					}
				}
				t.stringvalue += string(v)

			case 'U' : // \Unnnnnnnn	Unicode character with given code (1--8 hex digits)
				if r == '`' { // \\Uxxxxxxxx sequences not supported inside backticks
					t.Type = ERROR
					return
				}
				hasunicode = true
				// Trying to read the next rune
				if ccc, err = this.getCharacter(); err != nil {
					t.Type = ERROR
					return
				}
				if ccc.r == '{' {
					hascurly = true
				} else {
					hascurly = false
					// pushback the last character if it is not '{'
					if err = this.ungetCharacter(ccc); err != nil {
						t.Type = ERROR
					}					
				}				
				v = 0
				for i := 0; i < 8; i++ {
					// Trying to read the next rune
					if ccc, err = this.getCharacter(); err != nil {
						t.Type = ERROR
						return
					}
					if this.isCharacterHexa(ccc.r) {
						v = v*16 + rune(this.hexaValue(ccc.r))
					} else {
						// pushback the last character if it is not '0-9a-fA-F'
						if err = this.ungetCharacter(ccc); err != nil {
							t.Type = ERROR
						}
						break
					}
				}
				if v == 0 { // nul character not allowed
					t.Type = ERROR
					return
				}
				if hascurly {
					// Trying to read the next rune
					if ccc, err = this.getCharacter(); err != nil {
						t.Type = ERROR
						return
					}
					if ccc.r != '}' {
						t.Type = ERROR
						return
					}
				}				
				t.stringvalue += string(v)

			default:
				t.Type = ERROR
				return
			}
		default:
			t.stringvalue += string(c.r)
		}
	}
	return
}

func (this *Scanner) processDecimal(c *character, t *Token) {
	var err 						error

	t.Type = CONST_REAL

	// process integer part of the decimal (at the left of '.')
	if c.r != '.' {
		t.stringvalue = string(c.r)
		for {
			// Trying to read the next rune
			if c, err = this.getCharacter(); err != nil {
				t.Type = ERROR
				return
			}
			if c.r >= '0' && c.r <= '9' {
				t.stringvalue = t.stringvalue + string(c.r)
			} else {
				break
			}
		}
	}

	// process fractional part of the decimal (at the right of '.')
	if c.r == '.' {
		t.stringvalue = t.stringvalue + "."
		for {
			// Trying to read the next rune
			if c, err = this.getCharacter(); err != nil {
				t.Type = ERROR
				return
			}
			if c.r >= '0' && c.r <= '9' {
				t.stringvalue = t.stringvalue + string(c.r)
			} else {
				break
			}
		}
	}

	// process exponential part of the decimal (at the right of 'e' or 'E')
	if c.r == 'e' || c.r == 'E' {
		t.stringvalue = t.stringvalue + string(c.r)
		if c, err = this.getCharacter(); err != nil {
			t.Type = ERROR
			return
		}
		switch c.r {
		case '+' :
			t.stringvalue = t.stringvalue + "+"
		case '-' :
			t.stringvalue = t.stringvalue + "-"
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			t.stringvalue = t.stringvalue + string(c.r)
		default:
			break
		}
		for {
			// Trying to read the next rune
			if c, err = this.getCharacter(); err != nil {
				t.Type = ERROR
				return
			}
			if c.r >= '0' && c.r <= '9' {
				t.stringvalue = t.stringvalue + string(c.r)
			} else {
				break
			}
		}
	}

	t.realvalue, _ = strconv.ParseFloat(t.stringvalue,64)
	t.intvalue = int64(t.realvalue)

	// Is it an integer ? (next character is 'L') or a complex ? (next character is 'i')
	if c.r == 'L' {
		t.Type = CONST_INTEGER
		t.stringvalue = t.stringvalue + string(c.r)
	} else if c.r == 'i' {
		t.Type = CONST_COMPLEX
		t.stringvalue = t.stringvalue + string(c.r)
	} else {
		// pushback the last character if it is not 'L' or 'i'
		if err = this.ungetCharacter(c); err != nil {
			t.Type = ERROR
		}
	}

	return
}

func (this *Scanner) isCharacterHexa(r rune) bool {
	if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') {
		return true
	}
	return false
}

func (this *Scanner) hexaValue(r rune) int {
	if r >= '0' && r <= '9' {
		return int(r-'0')
	} else if r >= 'a' && r <= 'f' {
		return int(r-'a'+10)
	} else if r >= 'A' && r <= 'F' {
		return int(r-'A'+10)
	}
	return -1
}

func (this *Scanner) processHexadecimal(c *character, t *Token) {
	var err 					error
	var expn, exph, sign, n		int64
	var v						float64

	exph = -1
	t.Type = CONST_REAL

	// must be '0'
	if c.r != '0' {
		t.Type = ERROR
		return
	}

	// remove the 'x' or 'X'
	if c, err = this.getCharacter(); err != nil {
		t.Type = ERROR
		return
	}
	t.stringvalue = "0"+string(c.r)
	if c.r != 'x' && c.r != 'X' {
		t.Type = ERROR
		return		
	}

	// Get the firs character after "0x"
	if c, err = this.getCharacter(); err != nil {
		t.Type = ERROR
		return
	}

	// process integer part of the decimal (at the left of '.')
	if c.r != '.' {
		v = float64(this.hexaValue(c.r))
		t.stringvalue = t.stringvalue + string(c.r)
		for {
			// Trying to read the next rune
			if c, err = this.getCharacter(); err != nil {
				t.Type = ERROR
				return
			}
			if this.isCharacterHexa(c.r) {
				v = v*16 + float64(this.hexaValue(c.r))
				t.stringvalue = t.stringvalue + string(c.r)
			} else {
				break
			}
		}
	}

	// process fractional part of the decimal (at the right of '.')
	if c.r == '.' {
		exph = 0
		t.stringvalue = t.stringvalue + "."
		for {
			// Trying to read the next rune
			if c, err = this.getCharacter(); err != nil {
				t.Type = ERROR
				return
			}
			if this.isCharacterHexa(c.r) {
				exph += 4
				v = v*16 + float64(this.hexaValue(c.r))
				t.stringvalue = t.stringvalue + string(c.r)
			} else {
				break
			}
		}
	}

	// process exponential part of the decimal (at the right of 'e' or 'E')
	if c.r == 'p' || c.r == 'P' {
		t.stringvalue = t.stringvalue + string(c.r)
		sign = 1
		if c, err = this.getCharacter(); err != nil {
			t.Type = ERROR
			return
		}
		switch c.r {
		case '+' :
			t.stringvalue = t.stringvalue + "+"
			sign = 1
		case '-' :
			t.stringvalue = t.stringvalue + "-"
			sign = -1
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			t.stringvalue = t.stringvalue + string(c.r)
			n = int64(c.r - '0')
		default:
			break
		}
		for {
			// Trying to read the next rune
			if c, err = this.getCharacter(); err != nil {
				t.Type = ERROR
				return
			}
			if c.r >= '0' && c.r <= '9' {
				if n < 9999 {
					n = n*10 + int64(c.r-'0')
				}
				t.stringvalue = t.stringvalue + string(c.r)
			} else {
				break
			}
		}
	}

	// Is it an integer ? (next character is 'L') or a complex ? (next character is 'i')
	if c.r == 'L' {
		t.Type = CONST_INTEGER
		t.stringvalue = t.stringvalue + string(c.r)
	} else if c.r == 'i' {
		t.Type = CONST_COMPLEX
		t.stringvalue = t.stringvalue + string(c.r)
	} else {
		// pushback the last character if it is not 'L' or 'i'
		if err = this.ungetCharacter(c); err != nil {
			t.Type = ERROR
		}
	}

	// Compute the final value

	if v != 0.0 {
		expn += sign * n
		if exph > 0 {
			expn -= exph
		}
		if expn < 0 {
    		v /= math.Pow(2,float64(-expn))
		} else {
    		v *= math.Pow(2,float64(expn))
		}
	}

	if v > math.MaxFloat64 {
		t.realvalue = math.Inf(1)
	} else {
		t.realvalue = v
		t.intvalue = int64(v)
	}
	return
}

func (this *Scanner) processSymbol(c *character, t *Token) {
	var err 	error
	t.stringvalue = string(c.r)

	for {
		// Trying to read the next rune
		if c, err = this.getCharacter(); err != nil {
			t.Type = ERROR
			return
		}
		if c.r == '.' || c.r == '_' || unicode.IsLetter(c.r) || unicode.IsDigit(c.r) {
			// Continue to read the symbol
			t.stringvalue += string(c.r)
		} else {
			// End of symbol
			t.Type = SYMBOL
			if err = this.ungetCharacter(c); err != nil {
				t.Type = ERROR
			} else {

				// Is it a language reserved keyword ?
				if t.stringvalue == "if" 		{ t.Type = KEYWORD_IF; break }
				if t.stringvalue == "else" 		{ t.Type = KEYWORD_ELSE; break }
				if t.stringvalue == "for" 		{ t.Type = KEYWORD_FOR; break }
				if t.stringvalue == "in" 		{ t.Type = KEYWORD_IN; break }
				if t.stringvalue == "repeat" 	{ t.Type = KEYWORD_REPEAT; break }
				if t.stringvalue == "while" 	{ t.Type = KEYWORD_WHILE; break }
				if t.stringvalue == "next" 		{ t.Type = KEYWORD_NEXT; break }
				if t.stringvalue == "break" 	{ t.Type = KEYWORD_BREAK; break }
				if t.stringvalue == "function"	{ t.Type = KEYWORD_FUNCTION; break }

				// Is it a constant reserved keyword ?
				if t.stringvalue == "NaN" { t.Type = CONST_NAN; break }
				if t.stringvalue == "Inf" { t.Type = CONST_INF; break }
				if t.stringvalue == "TRUE" { t.Type = CONST_TRUE; break }
				if t.stringvalue == "FALSE" { t.Type = CONST_FALSE; break }
				if t.stringvalue == "NULL" { t.Type = CONST_NULL; break }
				if t.stringvalue == "NA" { t.Type = NA_LOGICAL; break }
				if t.stringvalue == "NA_character_" { t.Type = NA_CHARACTER; break }
				if t.stringvalue == "NA_integer_" { t.Type = NA_INTEGER; break }
				if t.stringvalue == "NA_real_" { t.Type = NA_REAL; break }
				if t.stringvalue == "NA_complex_" { t.Type = NA_COMPLEX; break }
			}
			break
		}
	}
	return
}

func (this *Scanner) processInfix(t *Token) {
	var c *character
	var err error

	t.stringvalue = "%"
	for {
		// Trying to read the next rune
		if c, err = this.getCharacter(); err != nil {
			t.Type = ERROR
			return
		}
		switch c.r {
		case '\n' :
			t.Type = ERROR
			return
		case '%' :
			t.Type = INFIX
			t.stringvalue += string(c.r)
			return
		default:
			t.stringvalue += string(c.r)
		}
	}
	return
}

func (this *Scanner) NextToken() (t *Token) {
	var b	bool
	var c 	*character
	var err error

	// Allocate and init a new Token
	t = new(Token)
	t.Type = ERROR
	t.offset = this.currentOffset
	t.ncol = this.ncol
	t.nline = this.nline
	t.nbyte = 0
	t.nrune = 0

	// Skip all the space ' ', tabulation '\t' and form feed '\f' and return the next character
	if c, err = this.getCharacterAfterSpaces(); err != nil {
		if err == io.EOF {
			t.Type = END_OF_INPUT
		}
		return
	}

	// Now we have the next rune after skipping all the consecutive spaces, tabulation and form feed runes,

	// Is it a symbol ?
	if unicode.IsLetter(c.r) {
		this.processSymbol(c,t)
	} else {
		switch	c.r {

		// Comment or Line directive tokens
		case '#':	this.processComment(t)

		// String value
		case '"',
			'\'',
			'`' :	this.processString(c.r,t)

		// Infix operator
		case '%':	this.processInfix(t)

		// Symbol value or numeric can start with a "."
		case '.':
			// is it a numeric ?
			if b, err = this.isNextCharacterDigit(); err != nil { return }
			if b {
				// It is a numeric (decimal)
				this.processDecimal(c,t)
			} else {
				// It must be a symbol
				this.processSymbol(c,t)
			}

		// Numeric value
		case '1', '2', '3', '4', '5', '6', '7', '8', '9' : this.processDecimal(c, t)
		case '0' :
			// is it hexadecimal 0x or 0X ?
			if b, err = this.isNextCharacterLetter('x'); err != nil { return }
			if b {
				// It is a numeric (hexadecimal)
				this.processHexadecimal(c,t)
			} else {
				// It is a numeric (decimal)
				this.processDecimal(c, t)
			}

		// Single rune tokens
		case '\n' :	t.Type = END_OF_LINE; this.nline++
		case '+' :	t.Type = OP_ADD; t.stringvalue = "+"
    	case '/' :	t.Type = OP_DIV; t.stringvalue = "/"
    	case '^' :	t.Type = OP_POW; t.stringvalue = "^"
    	case '~' :	t.Type = OP_TILDE; t.stringvalue = "~"
		case '?' :	t.Type = OP_QUESTION; t.stringvalue = "?"
    	case '$' :	t.Type = OP_DOLLAR; t.stringvalue = "$"
    	case '@' :	t.Type = OP_AT; t.stringvalue = "@"
		case '(' :	t.Type = OP_LEFT_ROUND; t.stringvalue = "("
		case ')' :	t.Type = OP_RIGHT_ROUND; t.stringvalue = ")"
		case '{' :	t.Type = OP_LEFT_CURLY; t.stringvalue = "{"
		case '}' :	t.Type = OP_RIGHT_CURLY; t.stringvalue = "}"
		case ']' : 	t.Type = OP_RIGHT_SQUARE; t.stringvalue = "]"

		// Compound tokens

		// < <- <<-
		case '<':
			if b, err = this.isNextCharacter('<'); err != nil {	return }
			if b {
				// is it "<<-"
				if b, err = this.isNextCharacter('-'); err != nil { return }
				if b {
					// yes it is "<<-"
					t.Type = OP_LEFT_ASSIGN2
					t.stringvalue = "<<-"
				} else {
					// incorrect "<<" symbol
					return
				}
			} else {
				// is it "<-" ?
				if b, err = this.isNextCharacter('-'); err != nil { return }
				if b {
					// yes it is "<-"
					t.Type = OP_LEFT_ASSIGN
					t.stringvalue = "<-"
				} else {
					// is it "<=" ?
					if b, err = this.isNextCharacter('='); err != nil { return }
					if b {
						// yes it is "<="
						t.Type = OP_LE
						t.stringvalue = "<="
					} else {
						// yes it is "<"
						t.Type = OP_LT
						t.stringvalue = "<"
					}
				}
			}

		// - -> ->>
		case '-':
			if b, err = this.isNextCharacter('>'); err != nil {	return }
			if b {
				// is it "->>" ?
				if b, err = this.isNextCharacter('>'); err != nil { return }
				if b {
					// yes it is "->>"
					t.Type = OP_RIGHT_ASSIGN2
					t.stringvalue = "->>"
				} else {
					// yes it is "->"
					t.Type = OP_RIGHT_ASSIGN
					t.stringvalue = "->"
				}
			} else {
				// yes it is "-"
				t.Type = OP_SUB
				t.stringvalue = "-"
			}

		// > >=
		case '>':
			if b, err = this.isNextCharacter('='); err != nil { return }
			if b {
				// yes it is ">="
				t.Type = OP_GE
				t.stringvalue = ">="
			} else {
				// yes it is ">"
				t.Type = OP_GT
				t.stringvalue = ">"
			}

		// ! !=
		case '!':
			if b, err = this.isNextCharacter('='); err != nil { return }
			if b {
				// yes it is "!="
				t.Type = OP_NE
				t.stringvalue = "!="
			} else {
				// yes it is "!"
				t.Type = OP_NOT
				t.stringvalue = "!"
			}

		// * **
		case '*':
			if b, err = this.isNextCharacter('*'); err != nil { return }
			if b {
				// yes it is "**"
				t.Type = OP_MUL2
				t.stringvalue = "**"
			} else {
				// yes it is "*"
				t.Type = OP_MUL
				t.stringvalue = "*"
			}

		// = ==
		case '=':
			if b, err = this.isNextCharacter('='); err != nil { return }
			if b {
				// yes it is "=="
				t.Type = OP_EQ
				t.stringvalue = "=="
			} else {
				// yes it is "="
				t.Type = OP_EQUAL_ASSIGN
				t.stringvalue = "="
			}

		// : :: ::: :=
		case ':':
			if b, err = this.isNextCharacter(':'); err != nil { return }
			if b {
				// is it "->>" ?
				if b, err = this.isNextCharacter(':'); err != nil { return }
				if b {
					// yes it is ":::"
					t.Type = OP_NAMESPACE_INTERNAL
					t.stringvalue = ":::"
				} else {
					// yes it is "::"
					t.Type = OP_NAMESPACE
					t.stringvalue = "::"
				}
			} else {
				if b, err = this.isNextCharacter('='); err != nil { return }
				if b {
					// yes it is ":="
					t.Type = OP_COLON_ASSIGN
					t.stringvalue = ":="
				} else {
					// yes it is ":"
					t.Type = OP_COLON
					t.stringvalue = ":"
				}
			}

		// & &&
		case '&':
			if b, err = this.isNextCharacter('&'); err != nil { return }
			if b {
				// yes it is "&&"
				t.Type = OP_AND2
				t.stringvalue = "&&"
			} else {
				// yes it is "&"
				t.Type = OP_AND
				t.stringvalue = "&"
			}

		// | ||
		case '|':
			if b, err = this.isNextCharacter('|'); err != nil { return }
			if b {
				// yes it is "||"
				t.Type = OP_OR2
				t.stringvalue = "||"
			} else {
				// yes it is "|"
				t.Type = OP_OR
				t.stringvalue = "|"
			}

		// [ [[
		case '[':
			if b, err = this.isNextCharacter('['); err != nil { return }
			if b {
				// yes it is "[["
				t.Type = OP_LEFT_SQUARE2
				t.stringvalue = "[["
			} else {
				// yes it is "["
				t.Type = OP_LEFT_SQUARE
				t.stringvalue = "["
			}
		}
	}
	return
}

func NewScanner(r io.Reader) (s *Scanner) {
	s = new(Scanner)

	// properties init
	s.reader = bufio.NewReader(r)
	s.currentOffset = 0
	s.ncol = 1
	s.nline = 1
	s.nrune = 0

	return
}


