package lexer

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tokens []Token

func optionalToken(t Token) *Token {
	return &t
}

type lexerTestCase struct {
	name   string
	tokens tokens
	idents []string
	input  string
}

func varTestCase(ident string, t Token) lexerTestCase {
	return lexerTestCase{
		name: fmt.Sprintf("var %s %s", ident, t),
		tokens: tokens{
			VAR,
			IDENT,
			t,
			EOF,
		},
		idents: []string{ident},
		input:  fmt.Sprintf("var %s %s", ident, ReservedWord(t)),
	}
}

func constTestCase(ident string, t Token) lexerTestCase {
	return lexerTestCase{
		name: fmt.Sprintf("const %s %s", ident, t),
		tokens: tokens{
			CONST,
			IDENT,
			t,
			EOF,
		},
		idents: []string{ident},
		input:  fmt.Sprintf("const %s %s", ident, ReservedWord(t)),
	}
}

func shortVarWithStringLiteralTestCase(ident string, literal string) lexerTestCase {
	return lexerTestCase{
		name:   fmt.Sprintf("short var with string literal %s", ident),
		tokens: tokens{IDENT, SHORT_VAR, STRING, EOF},
		idents: []string{ident},
		input:  fmt.Sprintf("%s := %s", ident, literal),
	}
}

func shortVarWithNumericLiteralTestCase(ident string, t Token, literal string) lexerTestCase {
	return lexerTestCase{
		name:   fmt.Sprintf("short var with numeric literal %s %s", ident, t),
		tokens: tokens{IDENT, SHORT_VAR, t, EOF},
		idents: []string{ident},
		input:  fmt.Sprintf("%s := %s", ident, literal),
	}
}

func assignTestCase(ident string, t Token, rhs string) lexerTestCase {
	idents := []string{ident}
	if t == IDENT {
		idents = append(idents, rhs)
	}
	return lexerTestCase{
		name:   fmt.Sprintf("assign with ident %s rhs %s", ident, t),
		tokens: tokens{IDENT, ASSIGN, t, EOF},
		idents: idents,
		input:  fmt.Sprintf("%s = %s", ident, rhs),
	}
}

func userDefinedBasicTypeTestCase(ident string, t Token) lexerTestCase {
	return lexerTestCase{
		name:   fmt.Sprintf("user-defined value type %s %s", ident, t),
		tokens: tokens{TYPE, IDENT, t, EOF},
		idents: []string{ident},
		input:  fmt.Sprintf("type %s %s", ident, ReservedWord(t)),
	}
}

func userDefinedTupleTypeTestCase(ident string, types ...Token) lexerTestCase {
	var ts strings.Builder
	tokens := tokens{
		TYPE,
		IDENT,
		T_TUPLE,
		OPEN_PAREN,
	}
	for i, t := range types {
		tokens = append(tokens, t)
		ts.WriteString(ReservedWord(t))
		if i+1 != len(types) {
			tokens = append(tokens, COMMA)
			ts.WriteRune(',')
		}
	}
	tokens = append(tokens, CLOSE_PAREN, EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("user-defined tuple type %s", ident),
		tokens: tokens,
		idents: []string{ident},
		input:  fmt.Sprintf("type %s tuple(%s)", ident, ts.String()),
	}
}

func userDefinedMapTypeTestCase(ident string, keyTypeTokens tokens, valueTypeTokens tokens, idents []string, input string) lexerTestCase {
	idents = append([]string{ident}, idents...)
	tokens := tokens{
		TYPE,
		IDENT,
		T_MAP,
		OPEN_BRACKET,
	}
	tokens = append(tokens, keyTypeTokens...)
	tokens = append(tokens, CLOSE_BRACKET)
	tokens = append(tokens, valueTypeTokens...)
	tokens = append(tokens, EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("user-defined map type %s", ident),
		tokens: tokens,
		idents: idents,
		input:  input,
	}
}

func tupleLiteralTestCase(ident string, identItems []string, values tokens, input string) lexerTestCase {
	var tokens tokens
	var idents []string
	if ident != "" {
		tokens = append(tokens, IDENT)
		idents = append(idents, ident)
	}
	idents = append(idents, identItems...)
	tokens = append(tokens, OPEN_PAREN)
	for i, t := range values {
		tokens = append(tokens, t)
		if i < len(values)-1 {
			tokens = append(tokens, COMMA)
		}
	}
	tokens = append(tokens, CLOSE_PAREN, EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("tuple literal %s", ident),
		tokens: tokens,
		idents: idents,
		input:  input,
	}
}

func arrayLiteralTokens(arrayType Token, inner ...Token) tokens {
	list := tokens{OPEN_BRACKET, INT, CLOSE_BRACKET, arrayType, OPEN_BRACE}
	list = append(list, inner...)
	list = append(list, CLOSE_BRACE, EOF)
	return list
}

func sliceLiteralTokens(sliceType Token, inner ...Token) tokens {
	list := tokens{OPEN_BRACKET, CLOSE_BRACKET, sliceType, OPEN_BRACE}
	list = append(list, inner...)
	list = append(list, CLOSE_BRACE, EOF)
	return list
}

func symbolLiteralTestCase(value string) lexerTestCase {
	return lexerTestCase{
		name:   fmt.Sprintf("symbol literal %s", value),
		tokens: tokens{SYMBOL, EOF},
		input:  fmt.Sprintf(":%s", value),
	}
}

func structTypeTestCase(ident string, fieldNames []string, types tokens, input string, multiline bool) lexerTestCase {
	idents := []string{}
	var tokens tokens
	if ident != "" {
		tokens = append(tokens, TYPE, IDENT)
		idents = append(idents, ident)
	}
	tokens = append(tokens, T_STRUCT, OPEN_BRACE)
	if multiline {
		tokens = append(tokens, NEWLINE)
	}
	for i, n := range fieldNames {
		tokens = append(tokens, IDENT, types[i])
		if multiline {
			tokens = append(tokens, NEWLINE)
		} else {
			if i < len(fieldNames)-1 {
				tokens = append(tokens, SEMI)
			}
		}
		idents = append(idents, n)
	}
	tokens = append(tokens, CLOSE_BRACE, EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("user-defined struct type %s", ident),
		tokens: tokens,
		idents: idents,
		input:  input,
	}
}

func structLiteralTestCase(ident string, fieldNames []string, types tokens, input string, multiline bool) lexerTestCase {
	idents := []string{ident}
	tokens := tokens{IDENT, OPEN_BRACE}
	if multiline {
		tokens = append(tokens, NEWLINE)
	}
	for i, n := range fieldNames {
		tokens = append(tokens, IDENT, ARROW, types[i])
		idents = append(idents, n)
		if multiline {
			tokens = append(tokens, COMMA, NEWLINE)
		} else {
			if i < len(fieldNames)-1 {
				tokens = append(tokens, COMMA)
			}
		}
	}
	tokens = append(tokens, CLOSE_BRACE, EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("struct literal %s", ident),
		tokens: tokens,
		idents: idents,
		input:  input,
	}
}

func mapLiteralTestCase(
	ident string,
	keyTypeTokens tokens,
	valueTypeTokens tokens,
	keys *Token,
	values *Token,
	count int,
	input string,
	multiline bool,
) lexerTestCase {
	tokens := tokens{T_MAP, OPEN_BRACKET}
	tokens = append(tokens, keyTypeTokens...)
	tokens = append(tokens, CLOSE_BRACKET)
	tokens = append(tokens, valueTypeTokens...)
	tokens = append(tokens, OPEN_BRACE)
	if multiline {
		tokens = append(tokens, NEWLINE)
	}
	for i := 0; i < count; i++ {
		tokens = append(tokens, *keys, ARROW, *values)
		if multiline {
			tokens = append(tokens, COMMA, NEWLINE)
		} else {
			if i < count-1 {
				tokens = append(tokens, COMMA)
			}
		}
	}
	tokens = append(tokens, CLOSE_BRACE, EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("map literal %s", ident),
		tokens: tokens,
		idents: []string{ident},
		input:  input,
	}
}

func funcTestCase(ident string, argNames []string, types tokens, input string) lexerTestCase {
	tokens := tokens{
		FUNC,
		IDENT,
		OPEN_PAREN,
	}
	idents := []string{ident}
	for i, argName := range argNames {
		tokens = append(tokens, IDENT, types[i])
		idents = append(idents, argName)
		if i < len(argNames)-1 {
			tokens = append(tokens, COMMA)
		}
	}
	if len(types) > len(argNames) {
		// there are types that are return values
		tokens = append(tokens, CLOSE_PAREN, OPEN_PAREN)
		for i := len(argNames); i < len(types); i++ {
			tokens = append(tokens, types[i])
			if i < len(types)-1 {
				tokens = append(tokens, COMMA)
			}
		}
	}

	tokens = append(tokens, CLOSE_PAREN, OPEN_BRACE, CLOSE_BRACE, EOF)

	return lexerTestCase{
		name:   fmt.Sprintf("func %s", ident),
		tokens: tokens,
		idents: idents,
		input:  input,
	}
}

func packageTestCase(ident string) lexerTestCase {
	return lexerTestCase{
		name: fmt.Sprintf("pacakge %s", ident),
		tokens: tokens{
			PACKAGE,
			IDENT,
			NEWLINE,
			NEWLINE,
			EOF,
		},
		idents: []string{ident},
		input:  fmt.Sprintf("package %s\n\n", ident),
	}
}

func returnTestCase(name string, values tokens, input string) lexerTestCase {
	tokens := tokens{RETURN}
	for i, t := range values {
		tokens = append(tokens, t)
		if i < len(values)-1 {
			tokens = append(tokens, COMMA)
		}
	}
	tokens = append(tokens, EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("return %s", name),
		tokens: tokens,
		input:  input,
	}
}

var testCases = []lexerTestCase{
	packageTestCase("main"),
	packageTestCase("something_else"),
	varTestCase("Test", T_STRING),
	varTestCase("Test", T_INT),
	varTestCase("Test", T_FLOAT),
	varTestCase("Test", T_BOOL),
	varTestCase("Test", ANY),
	{
		name:   "var any nil",
		tokens: tokens{VAR, IDENT, ANY, ASSIGN, NIL, EOF},
		idents: []string{"t"},
		input:  "var t any = nil",
	},
	constTestCase("Test", T_STRING),
	constTestCase("Test", T_INT),
	constTestCase("Test", T_FLOAT),
	constTestCase("Test", T_BOOL),
	assignTestCase("test1", IDENT, "other"),
	assignTestCase("test1", STRING, `"myval"`),
	assignTestCase("test1", INT, "100"),
	assignTestCase("test1", FLOAT, "0.123"),
	assignTestCase("test1", BOOL, "true"),
	assignTestCase("test1", BOOL, "false"),
	assignTestCase("testNil", NIL, "nil"),
	symbolLiteralTestCase("lower"),
	symbolLiteralTestCase("Upper"),
	symbolLiteralTestCase("lower_lower"),
	symbolLiteralTestCase("Something_More"),
	symbolLiteralTestCase("CamelCase"),
	symbolLiteralTestCase("Base64"),
	symbolLiteralTestCase("base_64"),
	symbolLiteralTestCase("_follow"),
	shortVarWithStringLiteralTestCase("test1", `"gus"`),
	shortVarWithStringLiteralTestCase("test2", `""`),
	shortVarWithStringLiteralTestCase("test3", `"123 gus"`),
	shortVarWithStringLiteralTestCase("test4", `"%Ƀ:=-ɸ_"`),
	shortVarWithStringLiteralTestCase("test4", `"\"test\""`),
	{
		name:   "template",
		tokens: tokens{VAR, IDENT, ASSIGN, TEMPLATE, NEWLINE, EOF},
		idents: []string{"message"},
		input:  "var message = `this is\n a test\n`\n",
	},
	shortVarWithNumericLiteralTestCase("test1", INT, "0"),
	shortVarWithNumericLiteralTestCase("test2", INT, "1"),
	shortVarWithNumericLiteralTestCase("test3", INT, "456"),
	shortVarWithNumericLiteralTestCase("test4", INT, "0b0011"),
	shortVarWithNumericLiteralTestCase("test5", INT, "0B0011"),
	shortVarWithNumericLiteralTestCase("test6", INT, "0x00ff00"),
	shortVarWithNumericLiteralTestCase("test7", INT, "0X00ff00"),
	shortVarWithNumericLiteralTestCase("test1", FLOAT, "0.0"),
	shortVarWithNumericLiteralTestCase("test2", FLOAT, "0.123"),
	shortVarWithNumericLiteralTestCase("test3", FLOAT, "0.1e16"),
	userDefinedBasicTypeTestCase("Test", T_STRING),
	userDefinedBasicTypeTestCase("Test", T_INT),
	userDefinedBasicTypeTestCase("Test", T_FLOAT),
	userDefinedBasicTypeTestCase("Test", T_BOOL),
	userDefinedTupleTypeTestCase("Test", T_STRING, T_INT),
	structTypeTestCase(
		"Row",
		[]string{"Col1", "Col2"},
		tokens{T_INT, T_STRING},
		multilineInput(`
			type Row struct {
				Col1 int
				Col2 string
			}
		`),
		true,
	),
	structTypeTestCase(
		"",
		[]string{"F1", "F2"},
		tokens{T_INT, T_STRING},
		"struct {F1 int; F2 string}",
		false,
	),
	tupleLiteralTestCase(
		"",
		nil,
		tokens{INT, FLOAT, BOOL},
		"(100, 0.1, false)",
	),
	tupleLiteralTestCase(
		"Sizes",
		nil,
		tokens{INT, INT},
		`Sizes(20, 104)`,
	),
	tupleLiteralTestCase(
		"multi",
		[]string{"ref"},
		tokens{INT, FLOAT, STRING, BOOL, IDENT},
		`multi(20, 0.1, "test", false, ref)`,
	),
	structLiteralTestCase(
		"Row",
		[]string{"Col1", "Col2", "Col3", "Col4"},
		tokens{INT, FLOAT, STRING, BOOL},
		multilineInput(`
			Row{
				Col1 => 123,
				Col2 => 0.5,
				Col3 => "test",
				Col4 => true,
			}
		`),
		true,
	),
	structLiteralTestCase(
		"RowNoFields",
		nil,
		nil,
		"RowNoFields{}",
		false,
	),
	structLiteralTestCase(
		"RowInline",
		[]string{"Col1", "Col2"},
		tokens{STRING, INT},
		`RowInline{Col1 => "val", Col2 => 123}`,
		false,
	),
	funcTestCase(
		"testEmpty",
		nil,
		nil,
		"func testEmpty() {}",
	),
	funcTestCase(
		"testMultiArg",
		[]string{"first", "second", "third", "fourth"},
		tokens{T_INT, T_FLOAT, T_BOOL, T_STRING},
		"func testMultiArg(first int, second float, third bool, fourth string) {}",
	),
	funcTestCase(
		"testSingleArg",
		[]string{"first"},
		tokens{T_INT},
		"func testSingleArg(first int) {}",
	),
	funcTestCase(
		"testNoArgReturn",
		nil,
		tokens{T_STRING},
		"func testNoArgReturn() (string) {}",
	),
	funcTestCase(
		"testSingleArgSingleReturn",
		[]string{"first"},
		tokens{T_INT, T_STRING},
		"func testSingleArgSingleReturn(first int) (string) {}",
	),
	funcTestCase(
		"testSingleArgMultiReturn",
		[]string{"first"},
		tokens{T_INT, T_STRING, T_BOOL},
		"func testSingleArgMultiReturn(first int) (string, bool) {}",
	),
	{
		name: "method receiver",
		tokens: tokens{
			FUNC, OPEN_PAREN, IDENT, IDENT, CLOSE_PAREN,
			IDENT, OPEN_PAREN, IDENT, T_STRING, CLOSE_PAREN,
			OPEN_PAREN, T_INT, COMMA, T_BOOL, CLOSE_PAREN,
			OPEN_BRACE, CLOSE_BRACE, EOF,
		},
		idents: []string{"r", "Receiver", "Exec", "s"},
		input:  "func (r Receiver) Exec(s string) (int, bool) {}",
	},
	returnTestCase("void", nil, "return"),
	returnTestCase("int", tokens{INT}, "return 123"),
	returnTestCase("float", tokens{FLOAT}, "return 0.123"),
	returnTestCase("string", tokens{STRING}, `return "test"`),
	returnTestCase("bool", tokens{BOOL}, "return false"),
	returnTestCase("multi", tokens{INT, BOOL}, "return 0, false"),
	returnTestCase("nil", tokens{NIL}, "return nil"),
	{
		name:   "arrayLiteral empty",
		tokens: arrayLiteralTokens(T_STRING),
		input:  "[0]string{}",
	},
	{
		name:   "arrayLiteral strings",
		tokens: arrayLiteralTokens(T_STRING, STRING, COMMA, STRING, COMMA, STRING),
		input:  `[3]string{"first", "second", "third"}`,
	},
	{
		name:   "arrayLiteral struct long",
		tokens: arrayLiteralTokens(IDENT, IDENT, OPEN_BRACE, IDENT, ARROW, INT, CLOSE_BRACE),
		idents: []string{"Row", "Row", "Col1"},
		input:  "[1]Row{Row{Col1 => 1}}",
	},
	{
		name:   "arrayLiteral struct short",
		tokens: arrayLiteralTokens(IDENT, OPEN_BRACE, IDENT, ARROW, INT, CLOSE_BRACE),
		idents: []string{"Row", "Col1"},
		input:  "[1]Row{{Col1 => 1}}",
	},
	{
		name:   "sliceLiteral empty",
		tokens: sliceLiteralTokens(T_INT),
		input:  "[]int{}",
	},
	{
		name:   "sliceLiteral ints",
		tokens: sliceLiteralTokens(T_INT, INT, COMMA, INT),
		input:  "[]int{1, 2}",
	},
	{
		name: "sliceLiteral multiple structs long",
		tokens: sliceLiteralTokens(
			IDENT,
			IDENT, OPEN_BRACE, CLOSE_BRACE, COMMA,
			IDENT, OPEN_BRACE, CLOSE_BRACE,
		),
		idents: []string{"Row", "Row", "Row"},
		input:  "[]Row{Row{}, Row{}}",
	},
	{
		name: "sliceLiteral multiple structs short",
		tokens: sliceLiteralTokens(
			IDENT,
			OPEN_BRACE, CLOSE_BRACE, COMMA,
			OPEN_BRACE, CLOSE_BRACE,
		),
		idents: []string{"Row"},
		input:  "[]Row{{}, {}}",
	},
	{
		name: "sliceLiteral multiline",
		tokens: sliceLiteralTokens(
			T_STRING,
			NEWLINE, STRING, COMMA, NEWLINE, STRING, COMMA, NEWLINE,
		),
		input: multilineInput(
			`[]string{
				"one",
				"two",
			}`,
		),
	},
	userDefinedMapTypeTestCase("StringString", tokens{T_STRING}, tokens{T_STRING}, nil, "type StringString map[string]string"),
	userDefinedMapTypeTestCase("StringBool", tokens{T_STRING}, tokens{T_BOOL}, nil, "type StringBool map[string]bool"),
	userDefinedMapTypeTestCase("StringInt", tokens{T_STRING}, tokens{T_INT}, nil, "type StringInt map[string]int"),
	userDefinedMapTypeTestCase("StringAny", tokens{T_STRING}, tokens{ANY}, nil, "type StringAny map[string]any"),
	userDefinedMapTypeTestCase("IntInt", tokens{T_INT}, tokens{T_INT}, nil, "type IntInt map[int]int"),
	userDefinedMapTypeTestCase("IntSliceInt", tokens{T_INT}, tokens{OPEN_BRACKET, CLOSE_BRACKET, T_INT}, nil, "type IntSliceInt map[int][]int"),
	userDefinedMapTypeTestCase("StructInt", tokens{IDENT}, tokens{T_INT}, []string{"Item"}, "type StructInt map[Item]int"),
	userDefinedMapTypeTestCase(
		"StructSliceFunc",
		tokens{IDENT},
		tokens{OPEN_BRACKET, CLOSE_BRACKET, FUNC, OPEN_PAREN, T_INT, CLOSE_PAREN, T_STRING},
		[]string{"Item"},
		"type StructSliceFunc map[Item][]func(int) string",
	),
	userDefinedMapTypeTestCase("StringInterface", tokens{T_STRING}, tokens{INTERFACE, OPEN_BRACE, CLOSE_BRACE}, nil, "type StringInterface map[string]interface{}"),
	mapLiteralTestCase(
		"StringString-empty",
		tokens{T_STRING},
		tokens{T_STRING},
		nil,
		nil,
		0,
		"map[string]string{}",
		false,
	),
	mapLiteralTestCase(
		"StringString-single",
		tokens{T_STRING},
		tokens{T_STRING},
		optionalToken(STRING),
		optionalToken(STRING),
		1,
		`map[string]string{"first" => "test1"}`,
		false,
	),
	mapLiteralTestCase(
		"StringString-multi",
		tokens{T_STRING},
		tokens{T_STRING},
		optionalToken(STRING),
		optionalToken(STRING),
		3,
		`map[string]string{"first" => "test1", "second" => "test2", "third" => "test3"}`,
		false,
	),
	mapLiteralTestCase(
		"StringString-single-multiline",
		tokens{T_STRING},
		tokens{T_STRING},
		optionalToken(STRING),
		optionalToken(STRING),
		1,
		multilineInput(`map[string]string{
			"first" => "test1",
		}`),
		true,
	),
	mapLiteralTestCase(
		"StringString-multi-multiline",
		tokens{T_STRING},
		tokens{T_STRING},
		optionalToken(STRING),
		optionalToken(STRING),
		3,
		multilineInput(`map[string]string{
			"first" => "test1",
			"second" => "test2",
			"third" => "test3",
		}`),
		true,
	),
	mapLiteralTestCase(
		"SymbolInt",
		tokens{T_SYMBOL},
		tokens{T_INT},
		optionalToken(SYMBOL),
		optionalToken(INT),
		2,
		`map[symbol]int{:first => 1, :second => 2}`,
		false,
	),
	{
		name:   "arrow-func no-args",
		tokens: tokens{IDENT, OPEN_PAREN, STRING, COMMA, OPEN_PAREN, CLOSE_PAREN, ARROW, OPEN_BRACE, CLOSE_BRACE, CLOSE_PAREN, EOF},
		idents: []string{"exists"},
		input:  `exists("test", () => {})`,
	},
	{
		name:   "arrow-func args",
		tokens: tokens{IDENT, OPEN_PAREN, STRING, COMMA, OPEN_PAREN, IDENT, CLOSE_PAREN, ARROW, OPEN_BRACE, CLOSE_BRACE, CLOSE_PAREN, EOF},
		idents: []string{"exists", "ok"},
		input:  `exists("test", (ok) => {})`,
	},
	{
		name: "arrow-func implicit return",
		tokens: tokens{
			IDENT, ACCESS, IDENT, OPEN_PAREN,
			OPEN_PAREN, IDENT, CLOSE_PAREN,
			ARROW, IDENT, ACCESS, IDENT, CLOSE_PAREN, EOF,
		},
		idents: []string{"users", "Filter", "user", "user", "inactive"},
		input:  `users.Filter((user) => user.inactive)`,
	},
	{
		name:   "named interface empty",
		tokens: tokens{TYPE, IDENT, INTERFACE, OPEN_BRACE, CLOSE_BRACE, EOF},
		idents: []string{"Empty"},
		input:  "type Empty interface {}",
	},
	{
		name: "interface singleline",
		tokens: tokens{
			INTERFACE, OPEN_BRACE,
			IDENT, OPEN_PAREN, CLOSE_PAREN, T_STRING, SEMI,
			IDENT, OPEN_PAREN, T_BOOL, CLOSE_PAREN, T_INT,
			CLOSE_BRACE, EOF,
		},
		idents: []string{"First", "Second"},
		input:  `interface{First() string; Second(bool) int}`,
	},
	{
		name: "interface multiline",
		tokens: tokens{
			INTERFACE, OPEN_BRACE, NEWLINE,
			IDENT, OPEN_PAREN, CLOSE_PAREN, T_STRING, NEWLINE,
			IDENT, OPEN_PAREN, T_BOOL, CLOSE_PAREN, T_INT, NEWLINE,
			CLOSE_BRACE, EOF,
		},
		idents: []string{"First", "Second"},
		input: multilineInput(
			`interface{
				First() string
				Second(bool) int
			}`),
	},
	{
		name: "enum Color",
		tokens: tokens{
			TYPE, IDENT, T_ENUM, OPEN_BRACE, NEWLINE,
			IDENT, NEWLINE, IDENT, NEWLINE, IDENT, NEWLINE,
			CLOSE_BRACE, EOF,
		},
		idents: []string{"Color", "Red", "Green", "Blue"},
		input: multilineInput(
			`type Color enum {
				Red
				Green
				Blue
			}`,
		),
	},
	{
		name: "enum ScreenColor",
		tokens: tokens{
			TYPE, IDENT, T_ENUM, OPEN_BRACE, NEWLINE,
			IDENT, OPEN_PAREN, T_INT, COMMA, T_INT, COMMA, T_INT, CLOSE_PAREN, NEWLINE,
			IDENT, OPEN_PAREN, T_INT, COMMA, T_INT, COMMA, T_INT, COMMA, T_INT, CLOSE_PAREN, NEWLINE,
			CLOSE_BRACE, EOF,
		},
		idents: []string{"ScreenColor", "RGB", "CMYK"},
		input: multilineInput(
			`type ScreenColor enum {
				RGB(int, int, int)
				CMYK(int, int, int, int)
			}`,
		),
	},
	{
		name: "enum Value",
		tokens: tokens{
			TYPE, IDENT, T_ENUM, OPEN_PAREN, T_INT, CLOSE_PAREN, OPEN_BRACE, NEWLINE,
			IDENT, NEWLINE, IDENT, NEWLINE, IDENT, NEWLINE,
			CLOSE_BRACE, EOF,
		},
		idents: []string{"Value", "First", "Second", "Third"},
		input: multilineInput(
			`type Value enum(int) {
				First
				Second
				Third
			}`,
		),
	},
	{
		name: "enum Name",
		tokens: tokens{
			TYPE, IDENT, T_ENUM, OPEN_PAREN, T_STRING, CLOSE_PAREN, OPEN_BRACE, NEWLINE,
			IDENT, ASSIGN, STRING, NEWLINE,
			IDENT, ASSIGN, STRING, NEWLINE,
			IDENT, ASSIGN, STRING, NEWLINE,
			CLOSE_BRACE, EOF,
		},
		idents: []string{"Name", "First", "Second", "Third"},
		input: multilineInput(
			`type Name enum(string) {
				First = "first"
				Second = "second"
				Third = "third"
			}`,
		),
	},
	{
		name: "enum Loc",
		tokens: tokens{
			TYPE, IDENT, T_ENUM, OPEN_PAREN, T_INT, COMMA, T_INT, CLOSE_PAREN, OPEN_BRACE, NEWLINE,
			IDENT, ASSIGN, OPEN_PAREN, INT, COMMA, INT, CLOSE_PAREN, NEWLINE,
			IDENT, ASSIGN, OPEN_PAREN, INT, COMMA, INT, CLOSE_PAREN, NEWLINE,
			CLOSE_BRACE, EOF,
		},
		idents: []string{"Loc", "Zero", "PlayerStart"},
		input: multilineInput(
			`type Loc enum(int, int) {
				Zero = (0, 0)
				PlayerStart = (10, 20)
			}`,
		),
	},
}

func TestLexer(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := make(chan Result)
			inputReader := bytes.NewReader([]byte(testCase.input))
			lex := New(inputReader, result)
			go Exec(lex)
			currentToken := 0
			currentIdent := 0
			lastPos := Position{0, 0, 0}
			for r := range result {
				require.NoError(t, r.Error)
				require.NotNil(t, r.Item)
				assert.True(t, r.Item.Pos.IsAfter(lastPos))
				require.Equalf(t, testCase.tokens[currentToken], r.Item.Token, "expected token [%s]; received [%s]", testCase.tokens[currentToken], r.Item.Token)
				if r.Item.Token == IDENT {
					require.Equalf(t, testCase.idents[currentIdent], r.Item.String, "expected ident [%s]; received [%s]", testCase.idents[currentIdent], r.Item.String)
					currentIdent++
				}
				currentToken += 1
				lastPos = r.Item.Pos
			}
		})
	}
}

func multilineInput(input string) string {
	return strings.TrimSpace(input)
}
