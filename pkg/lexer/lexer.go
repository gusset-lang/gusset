package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

var (
	ErrLexer = errors.New("lexer error")
)

const EOF_RUNE = rune(-1)

type Result struct {
	Item  *Item
	Error error
}

type Item struct {
	Pos    Position
	Token  Token
	String string
}

func itemResult(item Item) Result {
	return Result{
		Item:  &item,
		Error: nil,
	}
}

type runeMatcher func(rune) bool

type Position struct {
	Line       int
	Col        int
	MaxPrevCol int
}

func (p Position) IsAfter(t Position) bool {
	switch {
	case t.Line < p.Line:
		return true
	case t.Line == p.Line:
		return t.Col < p.Col
	default:
		return false
	}
}

func (p Position) Add(line, col int) Position {
	return Position{
		Line: p.Line + line,
		Col:  p.Col + col,
	}
}

func New(reader io.Reader, items chan Result) *Lexer {
	return &Lexer{
		Position{Line: 1, Col: 0},
		bufio.NewReader(reader),
		items,
	}
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
	result chan Result
}

func (l *Lexer) sendEOF() {
	l.result <- itemResult(Item{l.pos, EOF, ""})
}

func (l *Lexer) sendNewLine(pos Position) {
	l.result <- itemResult(Item{pos, NEWLINE, "\n"})
}

func (l *Lexer) sendItem(item *Item) {
	l.result <- itemResult(*item)
}

func (l *Lexer) sendError(err error) {
	l.result <- Result{
		Error: err,
	}
}

func (l *Lexer) next() (rune, error) {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return EOF_RUNE, nil
		}
		l.sendError(err)
		return 0, err
	}
	if r == '\n' {
		l.pos.MaxPrevCol = l.pos.Col
		l.pos.Col = 0
		l.pos.Line++
	} else {
		l.pos.Col++
	}
	return r, nil
}

func (l *Lexer) skip(n int) error {
	if _, err := l.reader.Discard(n); err != nil {
		l.sendError(err)
		return err
	}
	l.pos.Col += n
	return nil
}

func (l *Lexer) backup(r rune) error {
	if err := l.reader.UnreadRune(); err != nil {
		l.sendError(err)
		return err
	}
	if r == '\n' {
		l.pos.Col = l.pos.MaxPrevCol
		l.pos.MaxPrevCol = 0
		l.pos.Line--
	} else {
		l.pos.Col--
	}
	return nil
}

func (l *Lexer) peek() (rune, error) {
	b, err := l.reader.Peek(1)
	if err != nil {
		l.sendError(err)
		return 0, err
	}
	return rune(b[0]), nil
}

func (l *Lexer) peek2() ([2]rune, error) {
	b, err := l.reader.Peek(2)
	if err != nil {
		return [2]rune{}, err
	}
	return [2]rune{rune(b[0]), rune(b[1])}, nil
}

func (l *Lexer) matchRuneSequence(start Position, r rune) (bool, error) {
	node, ok := runeSequenceTree[r]
	if !ok {
		return false, nil
	}
	if node.t == nil {
		return false, fmt.Errorf("%w: matched leaf node of symbolic tree has no token", ErrLexer)
	}
	item := &Item{start, *node.t, runeSequences[*node.t]}

	if len(node.children) == 0 {
		l.sendItem(item)
		return true, nil
	}

	nextRunes, err := l.peek2()
	if err != nil {
		return false, err
	}

	secondNode, ok := node.children[nextRunes[0]]
	if !ok {
		l.sendItem(item)
		return true, nil
	}
	if secondNode.t == nil {
		return false, fmt.Errorf("%w: matched leaf node of symbolic tree has no token", ErrLexer)
	}
	item = &Item{start, *secondNode.t, runeSequences[*secondNode.t]}
	if err := l.skip(1); err != nil {
		return false, err
	}
	if len(secondNode.children) == 0 {
		l.sendItem(item)
		return true, nil
	}

	thirdNode, ok := node.children[nextRunes[1]]
	if !ok {
		l.sendItem(item)
		return true, nil
	}
	if thirdNode.t == nil {
		return false, fmt.Errorf("%w: matched leaf node of symbolic tree has no token", ErrLexer)
	}
	l.sendItem(&Item{start, *thirdNode.t, runeSequences[*thirdNode.t]})
	if err := l.skip(1); err != nil {
		return false, err
	}
	return true, nil
}

func (l *Lexer) collectSymbol(start Position) error {
	var seq strings.Builder

	for {
		r, err := l.next()
		if err != nil {
			return err
		}
		if r == EOF_RUNE {
			break
		}
		if r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r) {
			seq.WriteRune(r)
			continue
		}
		if err := l.backup(r); err != nil {
			return err
		}
		break
	}

	l.sendItem(&Item{start, SYMBOL, ":" + seq.String()})
	return nil
}

func (l *Lexer) itemFromAlphanum(startPos Position, initial rune) (*Item, error) {
	var seq strings.Builder
	seq.WriteRune(initial)

	collectSequence := func() *Item {
		seqString := seq.String()

		for _, bl := range boolLiteral {
			if seqString == bl {
				return &Item{startPos, BOOL, seqString}
			}
		}

		token := IDENT
		for res, resToken := range reserved {
			if seqString == res {
				token = resToken
			}
		}
		return &Item{startPos, token, seqString}
	}

	for {
		r, err := l.next()
		if err != nil {
			return nil, err
		}
		if r == EOF_RUNE {
			return collectSequence(), nil
		}

		if r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r) {
			seq.WriteRune(r)
			continue
		}

		if err := l.backup(r); err != nil {
			return nil, err
		}
		break
	}

	return collectSequence(), nil
}

func (l *Lexer) collectTemplateLiteral(start Position) error {
	var seq strings.Builder
	for {
		r, err := l.next()
		if err != nil {
			return err
		}
		if r == EOF_RUNE {
			// TODO: illegal
		}
		seq.WriteRune(r)
		if r == '`' && seq.Len() > 1 {
			break
		}
	}
	l.sendItem(&Item{start, TEMPLATE, "`" + seq.String()})
	return nil
}

func (l *Lexer) collectStructuredLiteral(start Position) error {
	var seq strings.Builder
	var delim rune
	for {
		r, err := l.next()
		if err != nil {
			return err
		}
		if r == EOF_RUNE {
			// todo: illegal
		}
		seq.WriteRune(r)
		if unicode.IsLetter(r) {
			continue
		}
		if !(r == '{' || r == '(' || r == '|') {
			// todo: illegal
		}
		delim = r
		break
	}

	if delim == 0 {
		return fmt.Errorf("%w: expected matched delimeter", ErrLexer)
	}

	matchedDelim, ok := structuredLiteralDelimiters[delim]
	if !ok {
		return fmt.Errorf("%w: expected to have a matched delimeter for %s", ErrLexer, string(delim))
	}

	delimOffset := 1

	for {
		r, err := l.next()
		if err != nil {
			return err
		}
		if r == EOF_RUNE {
			// TODO: illegal
		}
		seq.WriteRune(r)
		if r == delim {
			delimOffset++
		}
		if r == matchedDelim {
			delimOffset--
			if delimOffset == 0 {
				break
			}
		}
	}
	l.sendItem(&Item{start, STRUCTURED, "#" + seq.String()})
	return nil
}

func (l *Lexer) collectStringLiteral(start Position) error {
	var seq strings.Builder
	var prevRune rune

	for {
		r, err := l.next()
		if err != nil {
			return err
		}
		if r == EOF_RUNE {
			// TODO: illegal
		}
		seq.WriteRune(r)
		if r == '"' && prevRune != '\\' && seq.Len() > 1 {
			break
		}
		prevRune = r
	}
	l.sendItem(&Item{start, STRING, "\"" + seq.String()})
	return nil
}

func (l *Lexer) itemFromNumeric(start Position, initial rune) (*Item, error) {
	var item *Item
	var seq strings.Builder
	seq.WriteRune(initial)

	nextPeeked, err := l.peek()
	if err != nil {
		return nil, err
	}

	if initial == '0' && nextPeeked != '.' {
		next, err := l.next()
		if err != nil {
			return nil, err
		}

		writeWhileMatch := func(m runeMatcher) error {
			for {
				r, err := l.next()
				if err != nil {
					return err
				}
				if r == EOF_RUNE {
					return nil
				}
				if m(r) {
					seq.WriteRune(r)
				} else {
					return l.backup(r)
				}
			}
		}

		switch {
		case next == 'b' || next == 'B':
			seq.WriteRune(next)
			if err := writeWhileMatch(func(r rune) bool {
				return r == '0' || r == '1'
			}); err != nil {
				return nil, err
			}
		case next == 'x' || next == 'X':
			seq.WriteRune(next)
			if err := writeWhileMatch(func(r rune) bool {
				return unicode.Is(unicode.Hex_Digit, r)
			}); err != nil {
				return nil, err
			}
		case unicode.IsDigit(next):
			seq.WriteRune(next)
			if err := writeWhileMatch(unicode.IsDigit); err != nil {
				return nil, err
			}
		default:
			if err := l.backup(next); err != nil {
				return nil, err
			}
		}
		item = &Item{start, INT, seq.String()}
	} else {
		fractional := false
		exponent := false
		for {
			r, err := l.next()
			if err != nil {
				return nil, err
			}
			if r == EOF_RUNE {
				break
			}
			if r == '.' || r == 'e' || unicode.IsDigit(r) {
				if (fractional && r == '.') || (exponent && r == 'e') {
					return &Item{l.pos, ILLEGAL, string(r)}, nil
				}
				seq.WriteRune(r)

				if r == '.' {
					fractional = true
				}
				if r == 'e' {
					exponent = true
				}
			} else {
				if err := l.backup(r); err != nil {
					return nil, err
				}
				break
			}
		}
		item = &Item{start, INT, seq.String()}
		if fractional || exponent {
			item.Token = FLOAT
		}
	}

	return item, nil
}
