package lexer

import (
	"io"
	"unicode"
)

func Exec(l *Lexer) {
	for {
		start := l.pos
		r, err := l.next()
		if err != nil {
			if err == io.EOF {
				l.sendEOF()
			}
			break
		}
		if r == EOF_RUNE {
			l.sendEOF()
			break
		}

		if r == '\n' {
			l.sendNewLine(start)
			continue
		}

		if unicode.IsSpace(r) {
			continue
		}

		matchedRuneSeq, err := l.matchRuneSequence(start, r)
		if err != nil {
			break
		}
		if matchedRuneSeq {
			continue
		}

		if r == '"' {
			if err := l.collectStringLiteral(start); err != nil {
				break
			}
			continue
		}

		if r == ':' {
			if err := l.collectSymbol(start); err != nil {
				break
			}
			continue
		}

		if r == '`' {
			if err := l.collectTemplateLiteral(start); err != nil {
				break
			}
			continue
		}

		if r == '#' {
			if err := l.collectStructuredLiteral(start); err != nil {
				break
			}
			continue
		}

		if unicode.IsDigit(r) {
			item, err := l.itemFromNumeric(start, r)
			if err != nil {
				break
			}
			l.sendItem(item)
			continue
		}

		item, err := l.itemFromAlphanum(start, r)
		if err != nil {
			break
		}
		if item != nil {
			l.sendItem(item)
			continue
		}
	}

	close(l.result)
}
