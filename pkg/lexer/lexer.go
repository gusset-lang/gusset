package lexer

import (
	"bufio"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
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

type Position struct {
	Line int
	Col  int
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

type runeMatcher func(rune) bool

type Lexer struct {
	pos    Position
	reader *bufio.Reader
	result chan Result
}

// nextRune reads a rune and handles the EOF or unexpected errors.
// Bool result indicates whether execution should continue.
func (l *Lexer) nextRune() (rune, error) {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return EOF_RUNE, nil
		}
		l.sendError(err)
		return 0, err
	}
	return r, nil
}

func (l *Lexer) sendEOF() {
	l.result <- itemResult(Item{l.pos, EOF, ""})
}

func (l *Lexer) sendNewLine() {
	l.result <- itemResult(Item{l.pos, NEWLINE, "\n"})
}

func (l *Lexer) sendItem(item *Item) {
	l.result <- itemResult(*item)
}

func (l *Lexer) sendShortVar() {
	l.result <- itemResult(Item{l.pos, SHORT_VAR, ":="})
}

func (l *Lexer) sendColon() {
	l.result <- itemResult(Item{l.pos, COLON, ":"})
}

func (l *Lexer) sendError(err error) {
	l.result <- Result{
		Error: err,
	}
}

func (l *Lexer) sendIllegal(char rune) {
	l.result <- itemResult(Item{l.pos, ILLEGAL, string(char)})
}

// backup unreads a rune
func (l *Lexer) backup() error {
	if err := l.reader.UnreadRune(); err != nil {
		l.sendError(err)
		return err
	}
	return nil
}

func (l *Lexer) peek() (rune, error) {
	// Peek the first byte to determine the size of the rune
	b, err := l.reader.Peek(1)
	if err != nil {
		return 0, err
	}

	runeSize := 1
	for i := range b {
		if b[i] < 0x80 || b[i] >= 0xC0 {
			break
		}
		runeSize++
	}

	// Peek the required number of bytes to read the full rune
	b, err = l.reader.Peek(runeSize)
	if err != nil {
		return 0, err
	}

	// Decode the rune from the peeked bytes
	r, _ := utf8.DecodeRune(b)
	return r, nil
}

func (l *Lexer) advance(n int) {
	l.pos.Col += n
}

// nextLine moves the position to the first column of the next line.
func (l *Lexer) nextLine() {
	l.pos.Line++
	l.pos.Col = 0
}

// nextCol moves the position to the next column.
func (l *Lexer) nextCol() {
	l.pos.Col++
}

// itemFromRune finds if the rune is any single-rune token that
// does not conflict with a multi-rune token.
func (l *Lexer) itemFromRune(r rune) *Item {
	var item *Item
	switch r {
	case '(':
		item = &Item{l.pos, OPEN_PAREN, "("}
	case ')':
		item = &Item{l.pos, CLOSE_PAREN, ")"}
	case '{':
		item = &Item{l.pos, OPEN_BRACE, "{"}
	case '}':
		item = &Item{l.pos, CLOSE_BRACE, "}"}
	case '[':
		item = &Item{l.pos, OPEN_BRACKET, "["}
	case ']':
		item = &Item{l.pos, CLOSE_BRACKET, "]"}
	case '=':
		item = &Item{l.pos, ASSIGN, "="} // TODO: this will conflict with equality operator
	case ',':
		item = &Item{l.pos, COMMA, ","}
	case ';':
		item = &Item{l.pos, SEMI, ";"}
	}
	return item
}

// itemFromAlphanum finds an item from an alphanumeric sequence.
// The return int is the number of columns to advance.
func (l *Lexer) itemFromAlphanum(startPos Position) (*Item, int, error) {
	var seq strings.Builder
	var runeCount int

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
		r, err := l.nextRune()
		if err != nil {
			return nil, 0, err
		}
		if r == EOF_RUNE {
			return collectSequence(), runeCount, nil
		}

		if r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r) {
			seq.WriteRune(r)
			runeCount++
			continue
		}

		if err := l.backup(); err != nil {
			return nil, 0, err
		}

		if seq.Len() == 0 {
			return nil, 0, nil
		}

		return collectSequence(), runeCount, nil
	}
}

func (l *Lexer) itemFromStringLiteral(startPos Position) (*Item, int, error) {
	var seq strings.Builder
	var runeCount int

	for {
		r, err := l.nextRune()
		if err != nil {
			return nil, 0, err
		}
		seq.WriteRune(r)
		runeCount++
		if r == '"' && runeCount > 1 {
			break
		}
	}

	return &Item{startPos, STRING, seq.String()}, runeCount, nil
}

func (l *Lexer) itemFromNumeric(startPos Position) (*Item, int, error) {
	var item *Item
	var seq strings.Builder
	var runeCount int
	digit, err := l.nextRune()
	if err != nil {
		return nil, 0, err
	}
	seq.WriteRune(digit)
	runeCount++

	nextPeeked, err := l.peek()
	if err != nil {
		return nil, 0, err
	}

	if digit == '0' && nextPeeked != '.' {
		next, err := l.nextRune()
		if err != nil {
			return nil, 0, err
		}

		writeWhileMatch := func(m runeMatcher) error {
			for {
				r, err := l.nextRune()
				if err != nil {
					return err
				}
				if r == EOF_RUNE {
					return nil
				}
				if m(r) {
					seq.WriteRune(r)
					runeCount++
				} else {
					return l.backup()
				}
			}
		}

		switch {
		case next == 'b' || next == 'B':
			seq.WriteRune(next)
			runeCount++
			if err := writeWhileMatch(func(r rune) bool {
				return r == '0' || r == '1'
			}); err != nil {
				return nil, 0, err
			}
		case next == 'x' || next == 'X':
			seq.WriteRune(next)
			runeCount++
			if err := writeWhileMatch(func(r rune) bool {
				return unicode.Is(unicode.Hex_Digit, r)
			}); err != nil {
				return nil, 0, err
			}
		case unicode.IsDigit(next):
			seq.WriteRune(next)
			runeCount++
			if err := writeWhileMatch(unicode.IsDigit); err != nil {
				return nil, 0, err
			}
		default:
			if err := l.backup(); err != nil {
				return nil, 0, err
			}
		}
		item = &Item{startPos, INT, seq.String()}
	} else {
		fractional := false
		exponent := false
		for {
			r, err := l.nextRune()
			if err != nil {
				return nil, 0, err
			}
			if r == EOF_RUNE {
				break
			}
			if r == '.' || r == 'e' || unicode.IsDigit(r) {
				if (fractional && r == '.') || (exponent && r == 'e') {
					return &Item{startPos.Add(0, runeCount), ILLEGAL, string(r)}, 1, nil
				}
				seq.WriteRune(r)
				runeCount++

				if r == '.' {
					fractional = true
				}
				if r == 'e' {
					exponent = true
				}
			} else {
				if err := l.backup(); err != nil {
					return nil, 0, err
				}
				break
			}
		}
		item = &Item{startPos, INT, seq.String()}
		if fractional || exponent {
			item.Token = FLOAT
		}
	}

	return item, runeCount, nil
}

func New(reader io.Reader, items chan Result) *Lexer {
	return &Lexer{
		Position{Line: 1, Col: 0},
		bufio.NewReader(reader),
		items,
	}
}
