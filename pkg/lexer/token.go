package lexer

import (
	"sort"
)

func init() {
	runeSequenceTree = make(runeTree)

	keys := make([]Token, 0, len(runeSequences))
	for k := range runeSequences {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return len(runeSequences[keys[i]]) < len(runeSequences[keys[j]])
	})

	for _, t := range keys {
		runeSet := runeSequences[t]
		runeSequenceTree.insert(
			t,
			rune(runeSet[0]),
			[]rune(runeSet[1:]),
		)
	}
}

type Token int

const (
	EOF = Token(iota)
	ILLEGAL
	IDENT
	NEWLINE

	// Built-in types
	T_SYMBOL
	T_STRING
	T_INT
	T_FLOAT
	T_BOOL
	T_TUPLE
	T_STRUCT
	T_MAP
	T_ENUM

	// Literals
	SYMBOL
	STRING
	TEMPLATE
	INT
	FLOAT
	BOOL
	NIL

	// Decimal arithmetic
	ADD
	SUB
	MULT
	DIV
	MOD

	// Bitwise arithmetic
	BIT_AND
	BIT_OR
	BIT_NOT
	BIT_LEFT
	BIT_RIGHT
	BIT_CLEAR

	ASSIGN
	SHORT_VAR

	// Decimal arithmetic assignment
	ASSIGN_ADD
	ASSIGN_SUB
	ASSIGN_MULT
	ASSIGN_DIV
	ASSIGN_MOD
	ASSIGN_INC
	ASSIGN_DEC

	// Bitwise arithmetic assignment
	ASSIGN_BIT_AND
	ASSIGN_BIT_OR
	ASSIGN_BIT_NOT
	ASSIGN_BIT_LEFT
	ASSIGN_BIT_RIGHT
	ASSIGN_BIT_CLEAR

	// Comparison
	AND
	OR
	EQ
	LT
	GT
	NOT
	NEQ
	LTEQ
	GTEQ

	// Punctuation
	OPEN_PAREN
	CLOSE_PAREN
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_BRACE
	CLOSE_BRACE
	COMMA
	ACCESS
	SEMI
	ARROW
	OMIT

	// Keywords
	PACKAGE
	FUNC
	VAR
	CONST
	TYPE
	ANY
	INTERFACE
	RETURN
	IF
	ELSE
	SWITCH
	DEFAULT
	CONTINUE
	BREAK
	FOR
	MATCH
)

func (t Token) String() string {
	return tokenNames[t]
}

var tokenNames = [...]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	NEWLINE: "NEWLINE",
	IDENT:   "IDENT",

	// Built-in types
	T_SYMBOL: "T_SYMBOL",
	T_STRING: "T_STRING",
	T_INT:    "T_INT",
	T_FLOAT:  "T_FLOAT",
	T_BOOL:   "T_BOOL",
	T_TUPLE:  "T_TUPLE",
	T_STRUCT: "T_STRUCT",
	T_MAP:    "T_MAP",
	T_ENUM:   "T_ENUM",

	// Literals
	SYMBOL:   "SYMBOL",
	STRING:   "STRING",
	TEMPLATE: "TEMPLATE",
	INT:      "INT",
	FLOAT:    "FLOAT",
	BOOL:     "BOOL",
	NIL:      "NIL",

	// Decimal arithmetic
	ADD:  "ADD",
	SUB:  "SUB",
	MULT: "MULT",
	DIV:  "DIV",
	MOD:  "MOD",

	// Bitwise arithmetic
	BIT_AND:   "BIT_AND",
	BIT_OR:    "BIT_OR",
	BIT_NOT:   "BIT_NOT",
	BIT_LEFT:  "BIT_LEFT",
	BIT_RIGHT: "BIT_RIGHT",
	BIT_CLEAR: "BIT_CLEAR",

	ASSIGN:    "ASSIGN",
	SHORT_VAR: "SHORT_VAR",

	// Decimal arithmetic assignment
	ASSIGN_ADD:  "ASSIGN_ADD",
	ASSIGN_SUB:  "ASSIGN_SUB",
	ASSIGN_MULT: "ASSIGN_MULT",
	ASSIGN_DIV:  "ASSIGN_DIV",
	ASSIGN_MOD:  "ASSIGN_MOD",
	ASSIGN_INC:  "ASSIGN_INC",
	ASSIGN_DEC:  "ASSIGN_DEC",

	// Bitwise arithmetic assignment
	ASSIGN_BIT_AND:   "ASSIGN_BIT_AND",
	ASSIGN_BIT_OR:    "ASSIGN_BIT_OR",
	ASSIGN_BIT_NOT:   "ASSIGN_BIT_NOT",
	ASSIGN_BIT_LEFT:  "ASSIGN_BIT_LEFT",
	ASSIGN_BIT_RIGHT: "ASSIGN_BIT_RIGHT",
	ASSIGN_BIT_CLEAR: "ASSIGN_BIT_CLEAR",

	// Comparison
	AND:  "AND",
	OR:   "OR",
	EQ:   "EQ",
	LT:   "LT",
	GT:   "GT",
	NOT:  "NOT",
	NEQ:  "NEQ",
	LTEQ: "LTEQ",
	GTEQ: "GTEQ",

	// Punctuation
	OPEN_PAREN:    "OPEN_PAREN",
	CLOSE_PAREN:   "CLOSE_PAREN",
	OPEN_BRACKET:  "OPEN_BRACKET",
	CLOSE_BRACKET: "CLOSE_BRACKET",
	OPEN_BRACE:    "OPEN_BRACE",
	CLOSE_BRACE:   "CLOSE_BRACE",
	COMMA:         "COMMA",
	ACCESS:        "ACCESS",
	SEMI:          "SEMI",
	ARROW:         "ARROW",
	OMIT:          "OMIT",

	// Keywords
	PACKAGE:   "PACKAGE",
	FUNC:      "FUNC",
	VAR:       "VAR",
	CONST:     "CONST",
	TYPE:      "TYPE",
	ANY:       "ANY",
	INTERFACE: "INTERFACE",
	RETURN:    "RETURN",
	IF:        "IF",
	ELSE:      "ELSE",
	SWITCH:    "SWITCH",
	DEFAULT:   "DEFAULT",
	CONTINUE:  "CONTINUE",
	BREAK:     "BREAK",
	FOR:       "FOR",
	MATCH:     "MATCH",
}

var runeSequenceTree runeTree

var runeSequences = map[Token]string{
	ADD:  "+",
	SUB:  "-",
	MULT: "*",
	DIV:  "/",
	MOD:  "%",

	BIT_AND:   "&",
	BIT_OR:    "|",
	BIT_NOT:   "^",
	BIT_LEFT:  "<<",
	BIT_RIGHT: ">>",
	BIT_CLEAR: "&^",

	ASSIGN:    "=",
	SHORT_VAR: ":=",

	ASSIGN_ADD:  "+=",
	ASSIGN_SUB:  "-=",
	ASSIGN_MULT: "*=",
	ASSIGN_DIV:  "/=",
	ASSIGN_MOD:  "%=",
	ASSIGN_INC:  "++",
	ASSIGN_DEC:  "--",

	ASSIGN_BIT_AND:   "&=",
	ASSIGN_BIT_OR:    "|=",
	ASSIGN_BIT_NOT:   "^=",
	ASSIGN_BIT_LEFT:  "<<=",
	ASSIGN_BIT_RIGHT: ">>=",
	ASSIGN_BIT_CLEAR: "&^=",

	AND:  "&&",
	OR:   "||",
	EQ:   "==",
	LT:   "<",
	GT:   ">",
	NOT:  "!",
	NEQ:  "!=",
	LTEQ: "<=",
	GTEQ: ">=",

	OPEN_PAREN:    "(",
	CLOSE_PAREN:   ")",
	OPEN_BRACKET:  "[",
	CLOSE_BRACKET: "]",
	OPEN_BRACE:    "{",
	CLOSE_BRACE:   "}",
	COMMA:         ",",
	ACCESS:        ".",
	SEMI:          ";",
	ARROW:         "=>",
	OMIT:          "_",
}

type runeTree map[rune]runeTreeNode

type runeTreeNode struct {
	t        *Token
	children runeTree
}

func (tree runeTree) insert(t Token, r rune, runes []rune) {
	node, ok := tree[r]
	if !ok {
		node = runeTreeNode{}
	}

	if len(runes) == 0 {
		node.t = &t
	} else {
		if node.children == nil {
			node.children = make(runeTree)
		}
		node.children.insert(t, runes[0], runes[1:])
	}

	tree[r] = node
}

var boolLiteral = [2]string{
	"true", "false",
}

var reserved = map[string]Token{
	"package":   PACKAGE,
	"func":      FUNC,
	"var":       VAR,
	"const":     CONST,
	"type":      TYPE,
	"any":       ANY,
	"nil":       NIL,
	"interface": INTERFACE,
	"return":    RETURN,
	"if":        IF,
	"else":      ELSE,
	"switch":    SWITCH,
	"default":   DEFAULT,
	"continue":  CONTINUE,
	"break":     BREAK,
	"for":       FOR,
	"match":     MATCH,

	"symbol": T_SYMBOL,
	"string": T_STRING,
	"int":    T_INT,
	"float":  T_FLOAT,
	"bool":   T_BOOL,
	"tuple":  T_TUPLE,
	"struct": T_STRUCT,
	"map":    T_MAP,
	"enum":   T_ENUM,
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
