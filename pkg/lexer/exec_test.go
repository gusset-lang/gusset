package lexer_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/gusset-lang/gusset/pkg/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tokens []lexer.Token

type lexerTestCase struct {
	name   string
	tokens tokens
	idents []string
	input  string
}

func varTestCase(ident string, t lexer.Token) lexerTestCase {
	return lexerTestCase{
		name: fmt.Sprintf("var %s %s", ident, t),
		tokens: tokens{
			lexer.VAR,
			lexer.IDENT,
			t,
			lexer.EOF,
		},
		idents: []string{ident},
		input:  fmt.Sprintf("var %s %s", ident, lexer.ReservedWord(t)),
	}
}

func shortVarWithStringLiteralTestCase(ident string, literal string) lexerTestCase {
	return lexerTestCase{
		name:   fmt.Sprintf("short var with string literal %s", ident),
		tokens: tokens{lexer.IDENT, lexer.SHORT_VAR, lexer.STRING, lexer.EOF},
		idents: []string{ident},
		input:  fmt.Sprintf("%s := %s", ident, literal),
	}
}

func shortVarWithNumericLiteralTestCase(ident string, t lexer.Token, literal string) lexerTestCase {
	return lexerTestCase{
		name:   fmt.Sprintf("short var with numeric literal %s %s", ident, t),
		tokens: tokens{lexer.IDENT, lexer.SHORT_VAR, t, lexer.EOF},
		idents: []string{ident},
		input:  fmt.Sprintf("%s := %s", ident, literal),
	}
}

func assignTestCase(ident string, t lexer.Token, rhs string) lexerTestCase {
	idents := []string{ident}
	if t == lexer.IDENT {
		idents = append(idents, rhs)
	}
	return lexerTestCase{
		name:   fmt.Sprintf("assign with ident %s rhs %s", ident, t),
		tokens: tokens{lexer.IDENT, lexer.ASSIGN, t, lexer.EOF},
		idents: idents,
		input:  fmt.Sprintf("%s = %s", ident, rhs),
	}
}

func userDefinedAlphanumTypeTestCase(ident string, t lexer.Token) lexerTestCase {
	return lexerTestCase{
		name:   fmt.Sprintf("user-defined type %s %s", ident, t),
		tokens: tokens{lexer.TYPE, lexer.IDENT, t, lexer.EOF},
		idents: []string{ident},
		input:  fmt.Sprintf("type %s %s", ident, lexer.ReservedWord(t)),
	}
}

func userDefinedTupleTypeTestCase(ident string, types ...lexer.Token) lexerTestCase {
	var ts strings.Builder
	tokens := tokens{
		lexer.TYPE,
		lexer.IDENT,
		lexer.T_TUPLE,
		lexer.OPEN_PAREN,
	}
	for i, t := range types {
		tokens = append(tokens, t)
		ts.WriteString(lexer.ReservedWord(t))
		if i+1 != len(types) {
			tokens = append(tokens, lexer.COMMA)
			ts.WriteRune(',')
		}
	}
	tokens = append(tokens, lexer.CLOSE_PAREN, lexer.EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("user-defined tuple type %s", ident),
		tokens: tokens,
		idents: []string{ident},
		input:  fmt.Sprintf("type %s tuple(%s)", ident, ts.String()),
	}
}

func tupleLiteralTestCase(ident string, identItems []string, values tokens, input string) lexerTestCase {
	var tokens tokens
	var idents []string
	if ident != "" {
		tokens = append(tokens, lexer.IDENT)
		idents = append(idents, ident)
	}
	idents = append(idents, identItems...)
	tokens = append(tokens, lexer.OPEN_PAREN)
	for i, t := range values {
		tokens = append(tokens, t)
		if i < len(values)-1 {
			tokens = append(tokens, lexer.COMMA)
		}
	}
	tokens = append(tokens, lexer.CLOSE_PAREN, lexer.EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("tuple literal %s", ident),
		tokens: tokens,
		idents: idents,
		input:  input,
	}
}

func structTypeTestCase(ident string, fieldNames []string, types tokens, input string, multiline bool) lexerTestCase {
	idents := []string{}
	var tokens tokens
	if ident != "" {
		tokens = append(tokens, lexer.TYPE, lexer.IDENT)
		idents = append(idents, ident)
	}
	tokens = append(tokens, lexer.T_STRUCT, lexer.OPEN_BRACE)
	if multiline {
		tokens = append(tokens, lexer.NEWLINE)
	}
	for i, n := range fieldNames {
		tokens = append(tokens, lexer.IDENT, types[i])
		if multiline {
			tokens = append(tokens, lexer.NEWLINE)
		} else {
			if i < len(fieldNames)-1 {
				tokens = append(tokens, lexer.SEMI)
			}
		}
		idents = append(idents, n)
	}
	tokens = append(tokens, lexer.CLOSE_BRACE, lexer.EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("user-defined struct type %s", ident),
		tokens: tokens,
		idents: idents,
		input:  input,
	}
}

func structLiteralTestCase(ident string, fieldNames []string, types tokens, input string, multiline bool) lexerTestCase {
	idents := []string{ident}
	tokens := tokens{lexer.IDENT, lexer.OPEN_BRACE}
	if multiline {
		tokens = append(tokens, lexer.NEWLINE)
	}
	for i, n := range fieldNames {
		tokens = append(tokens, lexer.IDENT, lexer.COLON, types[i])
		idents = append(idents, n)
		if multiline {
			tokens = append(tokens, lexer.COMMA, lexer.NEWLINE)
		} else {
			if i < len(fieldNames)-1 {
				tokens = append(tokens, lexer.COMMA)
			}
		}
	}
	tokens = append(tokens, lexer.CLOSE_BRACE, lexer.EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("struct literal %s", ident),
		tokens: tokens,
		idents: idents,
		input:  input,
	}
}

func funcTestCase(ident string, argNames []string, types tokens, input string) lexerTestCase {
	tokens := tokens{
		lexer.FUNC,
		lexer.IDENT,
		lexer.OPEN_PAREN,
	}
	idents := []string{ident}
	for i, argName := range argNames {
		tokens = append(tokens, lexer.IDENT, types[i])
		idents = append(idents, argName)
		if i < len(argNames)-1 {
			tokens = append(tokens, lexer.COMMA)
		}
	}
	if len(types) > len(argNames) {
		// there are types that are return values
		tokens = append(tokens, lexer.CLOSE_PAREN, lexer.OPEN_PAREN)
		for i := len(argNames); i < len(types); i++ {
			tokens = append(tokens, types[i])
			if i < len(types)-1 {
				tokens = append(tokens, lexer.COMMA)
			}
		}
	}

	tokens = append(tokens, lexer.CLOSE_PAREN, lexer.OPEN_BRACE, lexer.CLOSE_BRACE, lexer.EOF)

	return lexerTestCase{
		name:   fmt.Sprintf("func %s", ident),
		tokens: tokens,
		idents: idents,
		input:  input,
	}
}

// func userDefinedFuncTypeTestCase(ident string, args map[string]lexer.Token, input string) lexerTestCase {
// 	tokens := tokens{
// 		lexer.TYPE,
// 		lexer.IDENT,
// 		lexer.FUNC,

// 	}
// 	return lexerTestCase{
// 		name: fmt.Sprintf("user-defined func type %s", ident),
// 		tokens: nil,
// 	}
// }

func packageTestCase(ident string) lexerTestCase {
	return lexerTestCase{
		name: fmt.Sprintf("pacakge %s", ident),
		tokens: tokens{
			lexer.PACKAGE,
			lexer.IDENT,
			lexer.NEWLINE,
			lexer.NEWLINE,
			lexer.EOF,
		},
		idents: []string{ident},
		input:  fmt.Sprintf("package %s\n\n", ident),
	}
}

func returnTestCase(name string, values tokens, input string) lexerTestCase {
	tokens := tokens{lexer.RETURN}
	for i, t := range values {
		tokens = append(tokens, t)
		if i < len(values)-1 {
			tokens = append(tokens, lexer.COMMA)
		}
	}
	tokens = append(tokens, lexer.EOF)
	return lexerTestCase{
		name:   fmt.Sprintf("return %s", name),
		tokens: tokens,
		input:  input,
	}
}

var testCases = []lexerTestCase{
	packageTestCase("main"),
	packageTestCase("something_else"),
	varTestCase("Test", lexer.T_STRING),
	varTestCase("Test", lexer.T_INT),
	varTestCase("Test", lexer.T_FLOAT),
	varTestCase("Test", lexer.T_BOOL),
	assignTestCase("test1", lexer.IDENT, "other"),
	assignTestCase("test1", lexer.STRING, `"myval"`),
	assignTestCase("test1", lexer.INT, "100"),
	assignTestCase("test1", lexer.FLOAT, "0.123"),
	assignTestCase("test1", lexer.BOOL, "true"),
	assignTestCase("test1", lexer.BOOL, "false"),
	shortVarWithStringLiteralTestCase("test1", `"gus"`),
	shortVarWithStringLiteralTestCase("test2", `""`),
	shortVarWithStringLiteralTestCase("test3", `"123 gus"`),
	shortVarWithStringLiteralTestCase("test4", `"%Ƀ:=-ɸ_"`),
	shortVarWithNumericLiteralTestCase("test1", lexer.INT, "0"),
	shortVarWithNumericLiteralTestCase("test2", lexer.INT, "1"),
	shortVarWithNumericLiteralTestCase("test3", lexer.INT, "456"),
	shortVarWithNumericLiteralTestCase("test4", lexer.INT, "0b0011"),
	shortVarWithNumericLiteralTestCase("test5", lexer.INT, "0B0011"),
	shortVarWithNumericLiteralTestCase("test6", lexer.INT, "0x00ff00"),
	shortVarWithNumericLiteralTestCase("test7", lexer.INT, "0X00ff00"),
	shortVarWithNumericLiteralTestCase("test1", lexer.FLOAT, "0.0"),
	shortVarWithNumericLiteralTestCase("test2", lexer.FLOAT, "0.123"),
	shortVarWithNumericLiteralTestCase("test3", lexer.FLOAT, "0.1e16"),
	userDefinedAlphanumTypeTestCase("Test", lexer.T_STRING),
	userDefinedAlphanumTypeTestCase("Test", lexer.T_INT),
	userDefinedAlphanumTypeTestCase("Test", lexer.T_FLOAT),
	userDefinedAlphanumTypeTestCase("Test", lexer.T_BOOL),
	userDefinedTupleTypeTestCase("Test", lexer.T_STRING, lexer.T_INT),
	structTypeTestCase(
		"Row",
		[]string{"Col1", "Col2"},
		tokens{lexer.T_INT, lexer.T_STRING},
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
		tokens{lexer.T_INT, lexer.T_STRING},
		"struct {F1 int; F2 string}",
		false,
	),
	tupleLiteralTestCase(
		"",
		nil,
		tokens{lexer.INT, lexer.FLOAT, lexer.BOOL},
		"(100, 0.1, false)",
	),
	tupleLiteralTestCase(
		"Sizes",
		nil,
		tokens{lexer.INT, lexer.INT},
		`Sizes(20, 104)`,
	),
	tupleLiteralTestCase(
		"multi",
		[]string{"ref"},
		tokens{lexer.INT, lexer.FLOAT, lexer.STRING, lexer.BOOL, lexer.IDENT},
		`multi(20, 0.1, "test", false, ref)`,
	),
	structLiteralTestCase(
		"Row",
		[]string{"Col1", "Col2", "Col3", "Col4"},
		tokens{lexer.INT, lexer.FLOAT, lexer.STRING, lexer.BOOL},
		multilineInput(`
			Row{
				Col1: 123,
				Col2: 0.5,
				Col3: "test",
				Col4: true,
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
		tokens{lexer.STRING, lexer.INT},
		`RowInline{Col1: "val", Col2: 123}`,
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
		tokens{lexer.T_INT, lexer.T_FLOAT, lexer.T_BOOL, lexer.T_STRING},
		"func testMultiArg(first int, second float, third bool, fourth string) {}",
	),
	funcTestCase(
		"testSingleArg",
		[]string{"first"},
		tokens{lexer.T_INT},
		"func testSingleArg(first int) {}",
	),
	funcTestCase(
		"testNoArgReturn",
		nil,
		tokens{lexer.T_STRING},
		"func testNoArgReturn() (string) {}",
	),
	funcTestCase(
		"testSingleArgSingleReturn",
		[]string{"first"},
		tokens{lexer.T_INT, lexer.T_STRING},
		"func testSingleArgSingleReturn(first int) (string) {}",
	),
	funcTestCase(
		"testSingleArgMultiReturn",
		[]string{"first"},
		tokens{lexer.T_INT, lexer.T_STRING, lexer.T_BOOL},
		"func testSingleArgMultiReturn(first int) (string, bool) {}",
	),
	returnTestCase("void", nil, "return"),
	returnTestCase("int", tokens{lexer.INT}, "return 123"),
	returnTestCase("float", tokens{lexer.FLOAT}, "return 0.123"),
	returnTestCase("string", tokens{lexer.STRING}, `return "test"`),
	returnTestCase("bool", tokens{lexer.BOOL}, "return false"),
	returnTestCase("multi", tokens{lexer.INT, lexer.BOOL}, "return 0, false"),
}

func TestLexer(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := make(chan lexer.Result)
			inputReader := bytes.NewReader([]byte(testCase.input))
			lex := lexer.New(inputReader, result)
			go lexer.Exec(lex)
			currentToken := 0
			currentIdent := 0
			lastPos := lexer.Position{0, 0}
			for r := range result {
				require.NoError(t, r.Error)
				require.NotNil(t, r.Item)
				assert.True(t, r.Item.Pos.IsAfter(lastPos))
				require.Equalf(t, testCase.tokens[currentToken], r.Item.Token, "expected token [%s]; received [%s]", testCase.tokens[currentToken], r.Item.Token)
				if r.Item.Token == lexer.IDENT {
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
