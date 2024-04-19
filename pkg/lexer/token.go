package lexer

type Token int

const (
	EOF = Token(iota)
	ILLEGAL
	IDENT
	NEWLINE

	OPEN_PAREN
	CLOSE_PAREN
	OPEN_BRACE
	CLOSE_BRACE
	OPEN_BRACKET
	CLOSE_BRACKET
	ASSIGN
	ACCESS
	COMMA
	COLON
	SEMI
	SHORT_VAR

	// Keywords
	PACKAGE
	FUNC
	VAR
	TYPE
	RETURN

	// Built-in types
	T_STRING
	T_INT
	T_FLOAT
	T_BOOL
	T_TUPLE
	T_STRUCT

	// Literals
	STRING
	INT
	FLOAT
	BOOL
)

func (t Token) String() string {
	return tokenNames[t]
}

var tokenNames = [...]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	NEWLINE: "NEWLINE",
	IDENT:   "IDENT",

	OPEN_PAREN:    "OPEN_PAREN",
	CLOSE_PAREN:   "CLOSE_PAREN",
	OPEN_BRACE:    "OPEN_BRACE",
	CLOSE_BRACE:   "CLOSE_BRACE",
	OPEN_BRACKET:  "OPEN_BRACKET",
	CLOSE_BRACKET: "CLOSE_BRACKET",
	ASSIGN:        "ASSIGN",
	ACCESS:        "ACCESS",
	COMMA:         "COMMA",
	COLON:         "COLON",
	SEMI:          "SEMI",
	SHORT_VAR:     "SHORT_VAR",

	PACKAGE: "PACKAGE",
	FUNC:    "FUNC",
	VAR:     "VAR",
	TYPE:    "TYPE",
	RETURN:  "RETURN",

	T_STRING: "T_STRING",
	T_INT:    "T_INT",
	T_FLOAT:  "T_FLOAT",
	T_BOOL:   "T_BOOL",
	T_TUPLE:  "T_TUPLE",
	T_STRUCT: "T_STRUCT",

	STRING: "STRING",
	INT:    "INT",
	FLOAT:  "FLOAT",
	BOOL:   "BOOL",
}

var boolLiteral = []string{
	"true", "false",
}

var reserved = map[string]Token{
	"package": PACKAGE,
	"func":    FUNC,
	"var":     VAR,
	"type":    TYPE,
	"return":  RETURN,

	"string": T_STRING,
	"int":    T_INT,
	"float":  T_FLOAT,
	"bool":   T_BOOL,
	"tuple":  T_TUPLE,
	"struct": T_STRUCT,
}

// ReservedWord gets the alphanum string reserved for the given token, if it exists.
// Otherwise an empty string is returned (not all tokens are reserved words).
func ReservedWord(t Token) string {
	for w, token := range reserved {
		if token == t {
			return w
		}
	}
	return ""
}
